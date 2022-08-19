package runtime

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func (r *Runtime) Start() {
	if r.Service == nil {
		r.Logger.Error("no Service set")
		return
	}
	if len(r.SyncImpl) == 0 {
		r.Logger.Error("no SyncImplementation set")
		return
	}
	if r.Evaluator == nil {
		r.Logger.Error("no Evaluator set")
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	g, gCtx := errgroup.WithContext(ctx)

	for _, s := range r.SyncImpl {
		p := s
		g.Go(func() error {
			return r.startSyncer(gCtx, p)
		})
	}

	g.Go(func() error {
		return r.Service.Serve(gCtx, r.Evaluator)
	})

	<-gCtx.Done()
	if err := g.Wait(); err != nil {
		r.Logger.Error(err)
	}
}
