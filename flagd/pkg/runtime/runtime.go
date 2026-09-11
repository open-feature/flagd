package runtime

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	msync "sync"
	"syscall"

	"github.com/open-feature/flagd/core/pkg/evaluator"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/service"
	"github.com/open-feature/flagd/core/pkg/sync"
	"github.com/open-feature/flagd/flagd/pkg/service/flag-evaluation/ofrep"
	flagsync "github.com/open-feature/flagd/flagd/pkg/service/flag-sync"
	"golang.org/x/sync/errgroup"
)

type Runtime struct {
	Evaluator         evaluator.IEvaluator
	Logger            *logger.Logger
	SyncService       flagsync.ISyncService
	OfrepService      ofrep.IOfrepService
	EvaluationService service.IFlagEvaluationService
	ServiceConfig     service.Configuration
	Syncs             []sync.ISync

	// sourceReadiness records whether each configured sync provider has successfully updated the evaluator.
	sourceReadiness []bool

	mu msync.Mutex
}

type providerDataSync struct {
	payload       sync.DataSync
	providerIndex int
}

//nolint:funlen
func (r *Runtime) Start() error {
	if r.EvaluationService == nil {
		return errors.New("no service set")
	}
	if len(r.Syncs) == 0 {
		return errors.New("no sync implementation set")
	}
	if r.Evaluator == nil {
		return errors.New("no evaluator set")
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	g, gCtx := errgroup.WithContext(ctx)
	dataSync := make(chan providerDataSync, len(r.Syncs))
	// Initialize DataSync channel watcher
	g.Go(func() error {
		for {
			select {
			case data := <-dataSync:
				r.updateAndEmit(data.payload, data.providerIndex)
			case <-gCtx.Done():
				return nil
			}
		}
	})
	// Init sync providers
	for _, s := range r.Syncs {
		if err := s.Init(gCtx); err != nil {
			return fmt.Errorf("sync provider Init returned error: %w", err)
		}
	}
	// Start sync providers. Each provider has its own channel so the runtime can preserve
	// provider identity even when two configured providers share a source URI.
	for index, s := range r.Syncs {
		p, providerIndex := s, index
		providerData := make(chan sync.DataSync, 1)
		g.Go(func() error {
			for {
				select {
				case data, ok := <-providerData:
					if !ok {
						return nil
					}
					select {
					case dataSync <- providerDataSync{payload: data, providerIndex: providerIndex}:
					case <-gCtx.Done():
						return nil
					}
				case <-gCtx.Done():
					return nil
				}
			}
		})
		g.Go(func() error {
			if err := p.Sync(gCtx, providerData); err != nil {
				return fmt.Errorf("sync provider returned error: %w", err)
			}
			return nil
		})
	}

	defer func() {
		r.Logger.Info("Shutting down server...")
		r.EvaluationService.Shutdown()
		r.Logger.Info("Server successfully shutdown.")
	}()

	g.Go(func() error {
		// Readiness probe rely on the runtime
		r.ServiceConfig.ReadinessProbe = r.isReady
		if err := r.EvaluationService.Serve(gCtx, r.ServiceConfig); err != nil {
			return fmt.Errorf("error returned from serving flag evaluation service: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		err := r.OfrepService.Start(gCtx)
		if err != nil {
			return fmt.Errorf("error from ofrep server: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		err := r.SyncService.Start(gCtx)
		if err != nil {
			return fmt.Errorf("error from sync server: %w", err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("errgroup closed with error: %w", err)
	}
	return nil
}

func (r *Runtime) isReady() bool {
	// if all providers can watch for flag changes, we are ready.
	for _, p := range r.Syncs {
		if !p.IsReady() {
			return false
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	for _, ready := range r.sourceReadiness {
		if !ready {
			return false
		}
	}
	return true
}

// updateAndEmit helps to update state, notify changes and trigger sync updates
func (r *Runtime) updateAndEmit(payload sync.DataSync, providerIndex int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	err := r.Evaluator.SetState(payload)
	if err != nil {
		r.Logger.Error(fmt.Sprintf("error setting state: %v", err))
		return
	}

	if providerIndex < len(r.sourceReadiness) {
		r.sourceReadiness[providerIndex] = true
	}
	r.SyncService.Emit(payload.Source)
}
