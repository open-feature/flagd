package runtime

import (
	"context"
	"errors"
	"github.com/open-feature/flagd/core/pkg/otel"
	"os"
	"os/signal"
	msync "sync"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/open-feature/flagd/core/pkg/eval"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/service"
	"github.com/open-feature/flagd/core/pkg/sync"
)

type Runtime struct {
	Evaluator   eval.IEvaluator
	Logger      *logger.Logger
	Service     service.IFlagEvaluationService
	SyncImpl    []sync.ISync
	config      Config
	metrics     *otel.MetricsRecorder
	mu          msync.Mutex
	serviceName string
}

type Config struct {
	ServicePort       uint16
	MetricsPort       uint16
	ServiceSocketPath string
	ServiceCertPath   string
	ServiceKeyPath    string

	SyncProviders []sync.SourceConfig
	CORS          []string
}

// nolint: funlen
func (r *Runtime) Start() error {
	if r.Service == nil {
		return errors.New("no service set")
	}
	if len(r.SyncImpl) == 0 {
		return errors.New("no sync implementation set")
	}
	if r.Evaluator == nil {
		return errors.New("no evaluator set")
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	g, gCtx := errgroup.WithContext(ctx)
	dataSync := make(chan sync.DataSync, len(r.SyncImpl))
	// Initialize DataSync channel watcher
	g.Go(func() error {
		for {
			select {
			case data := <-dataSync:
				// resync events are triggered when a delete occurs during flag merges in the store
				// resync events may trigger further resync events, however for a flag to be deleted from the store
				// its source must match, preventing the opportunity for resync events to snowball
				if resyncRequired := r.updateWithNotify(data); resyncRequired {
					for _, s := range r.SyncImpl {
						p := s
						go func() {
							g.Go(func() error {
								return p.ReSync(gCtx, dataSync)
							})
						}()
					}
				}
			case <-gCtx.Done():
				return nil
			}
		}
	})
	// Init sync providers
	for _, s := range r.SyncImpl {
		if err := s.Init(gCtx); err != nil {
			return err
		}
	}
	// Start sync provider
	for _, s := range r.SyncImpl {
		p := s
		g.Go(func() error {
			return p.Sync(gCtx, dataSync)
		})
	}
	g.Go(func() error {
		return r.Service.Serve(gCtx, r.Evaluator, service.Configuration{
			ReadinessProbe: r.isReady,
			Port:           r.config.ServicePort,
			MetricsPort:    r.config.MetricsPort,
			ServiceName:    r.serviceName,
		})
	})
	<-gCtx.Done()
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func (r *Runtime) isReady() bool {
	// if all providers can watch for flag changes, we are ready.
	for _, p := range r.SyncImpl {
		if !p.IsReady() {
			return false
		}
	}
	return true
}

// updateWithNotify helps to update state and notify listeners
func (r *Runtime) updateWithNotify(payload sync.DataSync) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	notifications, resyncRequired, err := r.Evaluator.SetState(payload)
	if err != nil {
		r.Logger.Error(err.Error())
		return false
	}

	r.Service.Notify(service.Notification{
		Type: service.ConfigurationChange,
		Data: map[string]interface{}{
			"flags": notifications,
		},
	})

	return resyncRequired
}
