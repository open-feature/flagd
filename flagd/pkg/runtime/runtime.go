package runtime

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	msync "sync"
	"syscall"
	"time"

	"github.com/open-feature/flagd/core/pkg/evaluator"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/service"
	"github.com/open-feature/flagd/core/pkg/sync"
	flagsync "github.com/open-feature/flagd/flagd/pkg/service/flag-sync"
	"golang.org/x/sync/errgroup"
)

type Runtime struct {
	Evaluator     evaluator.IEvaluator
	Logger        *logger.Logger
	FlagSync      flagsync.ISyncService
	Service       service.IFlagEvaluationService
	ServiceConfig service.Configuration
	SyncImpl      []sync.ISync

	mu msync.Mutex
}

//nolint:funlen
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
				// resync events are triggered when a delete occurs during flag merges in the store
				// resync events may trigger further resync events, however for a flag to be deleted from the store
				// its source must match, preventing the opportunity for resync events to snowball
				if resyncRequired := r.updateAndEmit(data); resyncRequired {
					for _, s := range r.SyncImpl {
						p := s
						g.Go(func() error {
							err := p.ReSync(gCtx, dataSync)
							if err != nil {
								return fmt.Errorf("error resyncing sources: %w", err)
							}
							return nil
						})
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
			return fmt.Errorf("sync provider Init returned error: %w", err)
		}
	}
	// Start sync provider
	for _, s := range r.SyncImpl {
		p := s
		g.Go(func() error {
			if err := p.Sync(gCtx, dataSync); err != nil {
				return fmt.Errorf("sync provider returned error: %w", err)
			}
			return nil
		})
	}

	defer func() {
		r.Logger.Info("Shutting down server...")
		r.Service.Shutdown()
		r.FlagSync.Shutdown()
		r.Logger.Info("Server successfully shutdown.")
	}()

	g.Go(func() error {
		// Readiness probe rely on the runtime
		r.ServiceConfig.ReadinessProbe = r.isReady
		if err := r.Service.Serve(gCtx, r.ServiceConfig); err != nil {
			return fmt.Errorf("error returned from serving flag evaluation service: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		// startup delay - allow all sync sources to finish the initial sync
		<-time.After(5 * time.Second)

		err := r.FlagSync.Start()
		if err != nil {
			return fmt.Errorf("error from server: %w", err)
		}

		return nil
	})

	<-gCtx.Done()
	if err := g.Wait(); err != nil {
		return fmt.Errorf("errgroup closed with error: %w", err)
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

// updateAndEmit helps to update state, notify changes and trigger sync updates
func (r *Runtime) updateAndEmit(payload sync.DataSync) bool {
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

	// Emit flag syncs only when re-syncs are not needed
	if !resyncRequired {
		r.FlagSync.Emit()
	}

	return resyncRequired
}
