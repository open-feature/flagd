package runtime

import (
	"context"
	"errors"
	"os"
	"os/signal"
	msync "sync"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/open-feature/flagd/pkg/eval"
	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/service"
	"github.com/open-feature/flagd/pkg/sync"
)

type Runtime struct {
	config   Config
	Service  service.IService
	SyncImpl []sync.ISync

	mu        msync.Mutex
	Evaluator eval.IEvaluator
	Logger    *logger.Logger
}

type Config struct {
	ServicePort       int32
	MetricsPort       int32
	ServiceSocketPath string
	ServiceCertPath   string
	ServiceKeyPath    string

	ProviderArgs    sync.ProviderArgs
	SyncURI         []string
	RemoteSyncType  string
	SyncBearerToken string

	CORS []string
}

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

	// Start sync providers
	for _, s := range r.SyncImpl {
		p := s
		g.Go(func() error {
			return p.Sync(gCtx, dataSync)
		})
	}

	g.Go(func() error {
		return r.Service.Serve(gCtx, r.Evaluator)
	})

	<-gCtx.Done()
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
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
