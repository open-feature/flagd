package ofrep

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/open-feature/flagd/core/pkg/evaluator"
	"github.com/open-feature/flagd/core/pkg/logger"
	"golang.org/x/sync/errgroup"
)

type IOfrepService interface {
	// Start the OFREP service with context for shutdown
	Start(context.Context) error
}

type SvcConfiguration struct {
	Logger *logger.Logger
	Port   uint16
}

type Service struct {
	Logger *logger.Logger
	server *http.Server
}

func NewOfrepService(evaluator evaluator.IEvaluator, cfg SvcConfiguration) (*Service, error) {
	server := http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           NewHandler(cfg.Logger, evaluator),
		ReadHeaderTimeout: 3 * time.Second,
	}

	return &Service{
		Logger: cfg.Logger,
		server: &server,
	}, nil
}

func (s Service) Start(ctx context.Context) error {
	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("error from ofrep service: %w", err)
		}

		return nil
	})

	group.Go(func() error {
		<-gCtx.Done()
		fmt.Println("context done")
		err := s.server.Close()
		if err != nil {
			return fmt.Errorf("error from ofrep server shutdown: %w", err)
		}

		return nil
	})

	err := group.Wait()
	if err != nil {
		return fmt.Errorf("error from ofrep service: %w", err)
	}

	return nil
}
