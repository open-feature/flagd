package runtime

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func (r *Runtime) Start() error {
	if r.Service == nil {
		return errors.New("runtime start: no Service set")
	}
	if len(r.SyncImpl) == 0 {
		return errors.New("runtime start: no SyncImplementation set")
	}
	if r.Evaluator == nil {
		return errors.New("runtime start: no Evaluator set")
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
		return err
	}
	return nil
}
