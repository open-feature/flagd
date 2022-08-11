package cmd

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron"
	"golang.org/x/sync/errgroup"

	log "github.com/sirupsen/logrus"

	"github.com/open-feature/flagd/pkg/eval"
	"github.com/open-feature/flagd/pkg/runtime"
	"github.com/open-feature/flagd/pkg/service"
	"github.com/open-feature/flagd/pkg/sync"
	"github.com/spf13/cobra"
)

var (
	serviceProvider   string
	syncProvider      string
	evaluator         string
	uri               []string
	servicePort       int32
	socketServicePath string
	bearerToken       string
	serverCertPath    string
	serverKeyPath     string
)

func findService(name string) (service.IService, error) {
	registeredServices := map[string]service.IService{
		"http": &service.HTTPService{
			HTTPServiceConfiguration: &service.HTTPServiceConfiguration{
				Port:           servicePort,
				ServerKeyPath:  serverKeyPath,
				ServerCertPath: serverCertPath,
			},
			GRPCService: &service.GRPCService{},
			Logger: log.WithFields(log.Fields{
				"service":   "http",
				"component": "service",
			}),
		},
		"grpc": &service.GRPCService{
			GRPCServiceConfiguration: &service.GRPCServiceConfiguration{
				Port:           servicePort,
				ServerKeyPath:  serverKeyPath,
				ServerCertPath: serverCertPath,
			},
			Logger: log.WithFields(log.Fields{
				"service":   "grpc",
				"component": "service",
			}),
		},
	}

	v, ok := registeredServices[name]
	if !ok {
		return nil, errors.New("no service-provider set")
	}
	log.Debugf("Using %s service-provider\n", name)
	return v, nil
}

func findSync(name string) ([]sync.ISync, error) {
	results := make([]sync.ISync, 0, len(uri))
	for _, u := range uri {
		registeredSync := map[string]sync.ISync{
			"filepath": &sync.FilePathSync{
				URI: u,
				Logger: log.WithFields(log.Fields{
					"sync":      "filepath",
					"component": "sync",
				}),
			},
			"remote": &sync.HTTPSync{
				URI:         u,
				BearerToken: bearerToken,
				Client: &http.Client{
					Timeout: time.Second * 10,
				},
				Logger: log.WithFields(log.Fields{
					"sync":      "remote",
					"component": "sync",
				}),
				Cron: cron.New(),
			},
		}
		v, ok := registeredSync[name]
		if !ok {
			return nil, errors.New("no sync-provider set")
		}
		results = append(results, v)
		log.Debugf("Using %s sync-provider on %q\n", name, u)
	}

	return results, nil
}

func findEvaluator(name string) (eval.IEvaluator, error) {
	registeredEvaluators := map[string]eval.IEvaluator{
		"json": &eval.JSONEvaluator{
			Logger: log.WithFields(log.Fields{
				"evaluator": "json",
				"component": "evaluator",
			}),
		},
	}

	v, ok := registeredEvaluators[name]
	if !ok {
		return nil, errors.New("no evaluator set")
	}

	log.Debugf("Using %s evaluator\n", name)
	return v, nil
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start flagd",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Configure loggers -------------------------------------------------------
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		log.SetLevel(log.DebugLevel)
		// Configure service-provider impl------------------------------------------
		var serviceImpl service.IService
		foundService, err := findService(serviceProvider)
		if err != nil {
			log.Errorf("Unable to find service '%s'", serviceProvider)
			return
		}
		serviceImpl = foundService

		// Configure sync-provider impl--------------------------------------------
		var syncImpl []sync.ISync
		foundSync, err := findSync(syncProvider)
		if err != nil {
			log.Errorf("Unable to find sync '%s'", syncProvider)
			return
		}
		syncImpl = foundSync

		// Configure evaluator-provider impl------------------------------------------
		var evalImpl eval.IEvaluator
		foundEval, err := findEvaluator(evaluator)
		if err != nil {
			log.Errorf("Unable to find evaluator '%s'", evaluator)
			return
		}
		evalImpl = foundEval

		// Serve ------------------------------------------------------------------
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer cancel()
		g, gCtx := errgroup.WithContext(ctx)
		runtime.Start(gCtx, syncImpl, serviceImpl, evalImpl, log.WithFields(log.Fields{
			"component": "runtime",
		}))
		if err := g.Wait(); err != nil {
			log.Printf("exit reason: %s \n", err)
		} else {
			log.Println("server shutdown")
		}
	},
}

func init() {
	startCmd.Flags().Int32VarP(
		&servicePort, "port", "p", 8013, "Port to listen on")
	startCmd.Flags().StringVarP(
		&socketServicePath, "socketpath", "d", "/tmp/flagd.sock", "flagd socket path")
	startCmd.Flags().StringVarP(
		&serviceProvider, "service-provider", "s", "http", "Set a serve provider e.g. http or grpc")
	startCmd.Flags().StringVarP(
		&syncProvider, "sync-provider", "y", "filepath", "Set a sync provider e.g. filepath or remote")
	startCmd.Flags().StringVarP(
		&evaluator, "evaluator", "e", "json", "Set an evaluator e.g. json")
	startCmd.Flags().StringVarP(
		&serverCertPath, "server-cert-path", "c", "", "Server side tls certificate path")
	startCmd.Flags().StringVarP(
		&serverKeyPath, "server-key-path", "k", "", "Server side tls key path")
	startCmd.Flags().StringSliceVarP(
		&uri, "uri", "f", []string{}, "Set a sync provider uri to read data from this can be a filepath or url. "+
			"Using multiple providers is supported where collisions between "+
			"flags with the same key, the later will be used.")
	startCmd.Flags().StringVarP(
		&bearerToken, "bearer-token", "b", "", "Set a bearer token to use for remote sync")

	_ = startCmd.MarkFlagRequired("uri")
	rootCmd.AddCommand(startCmd)
}
