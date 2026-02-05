package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"buf.build/gen/go/open-feature-forking/flagd/grpc/go/sync/v1/syncv1grpc"
	"github.com/open-feature/flagd/flagd-proxy/tests/loadtest/pkg/client"
	trigger "github.com/open-feature/flagd/flagd-proxy/tests/loadtest/pkg/trigger"
	"github.com/open-feature/flagd/flagd-proxy/tests/loadtest/pkg/watcher"
)

type Handler struct {
	config  Config
	trigger trigger.Trigger
}

type Config struct {
	FilePath string `json:"filePath"`
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	OutFile  string `json:"outFile"`
}

type TestConfig struct {
	Watchers int           `json:"watchers"`
	Repeats  int           `json:"repeats"`
	Delay    time.Duration `json:"delay"`
}

type TestResult struct {
	TotalTime      time.Duration `json:"totalTime"`
	TimePerWatcher time.Duration `json:"timePerWatcher"`
}

type ProfilingResults struct {
	Tests                 []TestResult  `json:"tests"`
	AverageTotalDuration  time.Duration `json:"averageTotalDuration"`
	AverageTimePerWatcher time.Duration `json:"averageTimePerWatcher"`
	Watchers              int           `json:"watchers"`
	Repeats               int           `json:"repeats"`
}

func NewHandler(config Config, trigger trigger.Trigger) *Handler {
	return &Handler{
		config:  config,
		trigger: trigger,
	}
}

func (h *Handler) Profile(ctx context.Context, configs []TestConfig) ([]ProfilingResults, error) {
	out := []ProfilingResults{}
	for _, config := range configs {
		if err := h.trigger.Setup(); err != nil {
			return []ProfilingResults{}, fmt.Errorf("unable to setup trigger: %w", err)
		}
		results := []TestResult{}
		for i := 1; i <= config.Repeats; i++ {
			fmt.Printf("starting profile %d\n", i)
			res := h.runTest(ctx, config.Watchers)
			results = append(results, res)
			fmt.Println("-----------------------")
			time.Sleep(config.Delay)
		}
		timePer := time.Duration(0)
		totalTime := time.Duration(0)
		for _, res := range results {
			timePer += res.TimePerWatcher
			totalTime += res.TotalTime
		}
		out = append(out, ProfilingResults{
			Watchers:              config.Watchers,
			Repeats:               config.Repeats,
			Tests:                 results,
			AverageTotalDuration:  totalTime / time.Duration(config.Repeats),
			AverageTimePerWatcher: timePer / time.Duration(config.Repeats),
		})
	}

	return out, h.writeFile(out)
}

//nolint:funlen
func (h *Handler) runTest(ctx context.Context, watchers int) TestResult {
	readyWg := sync.WaitGroup{}
	readyWg.Add(watchers)
	readySuccess := make(chan bool, 1)

	finishedWg := sync.WaitGroup{}
	finishedWg.Add(watchers)
	finishedSuccess := make(chan bool, 1)

	errChan := make(chan error, 1)

	fmt.Printf("starting %d watchers...\n", watchers)
	var c syncv1grpc.FlagSyncServiceClient
	var err error

	for i := 0; i < watchers; i++ {
		if i%250 == 0 {
			c, err = client.NewClient(client.Config{
				Host: h.config.Host,
				Port: h.config.Port,
			})
			if err != nil {
				log.Fatal(err)
			}
		}
		go func() {
			w := watcher.NewWatcher(c, h.config.FilePath)
			go func() {
				if err := w.StartWatcher(ctx); err != nil {
					fmt.Println(err)
				}
			}()
			<-w.Ready
			readyWg.Done()
			if err := w.Wait(); err != nil {
				log.Fatal(err)
			}
			finishedWg.Done()
		}()
	}

	go func() {
		readyWg.Wait()
		readySuccess <- true
	}()

	go func() {
		finishedWg.Wait()
		finishedSuccess <- true
	}()

	fmt.Println("waiting for watchers to be ready...")
	select {
	case err := <-errChan:
		log.Fatal(err)
	case <-readySuccess:
		fmt.Println("all watchers ready, starting timer and writing to file...")
	}

	start := time.Now()

	if err := h.trigger.Update(); err != nil {
		log.Fatal(err)
	}

	select {
	case err := <-errChan:
		log.Fatal(err)
	case <-finishedSuccess:
		fmt.Println("all watchers ready, starting timer and writing to file...")
	}

	end := time.Now()
	fmt.Println("done")

	timeTaken := end.Sub(start)
	fmt.Println("process took", timeTaken)
	fmt.Println("time per run", timeTaken/time.Duration(watchers))
	return TestResult{
		TotalTime:      timeTaken,
		TimePerWatcher: timeTaken / time.Duration(watchers),
	}
}

func (h *Handler) writeFile(results []ProfilingResults) error {
	resB, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		return fmt.Errorf("unable to marshal profiling results: %w", err)
	}
	if err = os.WriteFile(h.config.OutFile, resB, 0o600); err != nil {
		return fmt.Errorf("unable to write output file %s: %w", h.config.OutFile, err)
	}
	return nil
}
