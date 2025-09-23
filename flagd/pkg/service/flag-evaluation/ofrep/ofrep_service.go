package ofrep

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/open-feature/flagd/core/pkg/evaluator"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	"github.com/rs/cors"
	"golang.org/x/sync/errgroup"
)

type IOfrepService interface {
	// Start the OFREP service with context for shutdown
	Start(context.Context) error
}

type SvcConfiguration struct {
	Logger          *logger.Logger
	Port            uint16
	ServiceName     string
	MetricsRecorder telemetry.IMetricsRecorder
}

type Service struct {
	logger *logger.Logger
	port   uint16
	server *http.Server
}

func NewOfrepService(
	evaluator evaluator.IEvaluator, origins []string, cfg SvcConfiguration, contextValues map[string]any, headerToContextKeyMappings map[string]string,
) (*Service, error) {
	corsMW := cors.New(cors.Options{
		AllowedOrigins: origins,
		AllowedMethods: []string{http.MethodPost},
	})

	h := corsMW.Handler(NewOfrepHandler(
		cfg.Logger,
		evaluator,
		contextValues,
		headerToContextKeyMappings,
		cfg.MetricsRecorder,
		cfg.ServiceName,
	))

	server := http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           h,
		ReadHeaderTimeout: 3 * time.Second,
	}

	return &Service{
		logger: cfg.Logger,
		port:   cfg.Port,
		server: &server,
	}, nil
}

func (s Service) Start(ctx context.Context) error {
	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		s.logger.Info(fmt.Sprintf("ofrep service listening at %d", s.port))
		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("error from ofrep service: %w", err)
		}

		return nil
	})

	group.Go(func() error {
		<-gCtx.Done()
		s.logger.Info("shutting down ofrep service")
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
