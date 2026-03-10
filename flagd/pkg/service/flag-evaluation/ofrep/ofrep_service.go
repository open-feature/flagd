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
	corsmw "github.com/open-feature/flagd/flagd/pkg/service/middleware/cors"
	"golang.org/x/sync/errgroup"
)

type IOfrepService interface {
	// Start the OFREP service with context for shutdown
	Start(context.Context) error
}

type SvcConfiguration struct {
	Logger                *logger.Logger
	Port                  uint16
	ServiceName           string
	MetricsRecorder       telemetry.IMetricsRecorder
	MaxRequestBodyBytes   int64
	MaxRequestHeaderBytes int64
}

type Service struct {
	logger *logger.Logger
	port   uint16
	server *http.Server
}

func NewOfrepService(
	evaluator evaluator.IEvaluator, origins []string, cfg SvcConfiguration, contextValues map[string]any, headerToContextKeyMappings map[string]string,
) (*Service, error) {
	corsMiddleware := corsmw.New(origins)

	var h http.Handler = NewOfrepHandler(
		cfg.Logger,
		evaluator,
		contextValues,
		headerToContextKeyMappings,
		cfg.MetricsRecorder,
		cfg.ServiceName,
	)
	if cfg.MaxRequestBodyBytes > 0 {
		h = http.MaxBytesHandler(h, cfg.MaxRequestBodyBytes)
	}
	h = corsMiddleware.Handler(h)

	server := http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           h,
		ReadHeaderTimeout: 3 * time.Second,
		MaxHeaderBytes:    int(cfg.MaxRequestHeaderBytes),
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
