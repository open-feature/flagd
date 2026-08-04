package ofrep

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/open-feature/flagd/core/pkg/evaluator"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	"github.com/open-feature/flagd/flagd/pkg/service/flag-evaluation/ofrep/sse"
	corsmw "github.com/open-feature/flagd/flagd/pkg/service/middleware/cors"
	"golang.org/x/sync/errgroup"
)

// ssePath is the endpoint OFREP clients subscribe to for change notifications (ADR-0008).
const ssePath = "/ofrep/v1/sse"

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

	// SSE (ADR-0008) settings
	SSEEnabled            bool
	SSEInactivityDelaySec int
	SSEPublicURL          string
}

type Service struct {
	logger *logger.Logger
	port   uint16
	server *http.Server
	sse    *sse.Service
}

func NewOfrepService(
	evaluator evaluator.IEvaluator, flagStore store.IStore, origins []string, cfg SvcConfiguration, contextValues map[string]any, headerToContextKeyMappings map[string]string,
) (*Service, error) {
	corsMiddleware := corsmw.New(origins)

	// SSE requires a flag store to watch; without one we cannot serve or advertise it, so treat
	// a missing store as SSE disabled rather than crashing later in the tracker goroutine.
	sseEnabled := cfg.SSEEnabled
	if sseEnabled && flagStore == nil {
		cfg.Logger.Warn("OFREP SSE requested but no flag store was provided; disabling SSE")
		sseEnabled = false
	}

	var sseService *sse.Service
	sseCfg := SSEConfig{
		Enabled:            sseEnabled,
		InactivityDelaySec: cfg.SSEInactivityDelaySec,
		PublicURL:          cfg.SSEPublicURL,
	}
	if sseEnabled {
		sseService = sse.New(flagStore, sse.Config{Logger: cfg.Logger})
		sseCfg.Versioner = sseService.Tracker()
	}

	ofrepHandler := NewOfrepHandler(
		cfg.Logger,
		evaluator,
		contextValues,
		headerToContextKeyMappings,
		cfg.MetricsRecorder,
		cfg.ServiceName,
		sseCfg,
	)

	// Route the long-lived SSE stream separately from the request/response evaluate routes so
	// the request-body limit is not applied to the stream. Everything else falls through to the
	// existing OFREP handler.
	mux := http.NewServeMux()
	var evaluateHandler http.Handler = ofrepHandler
	if cfg.MaxRequestBodyBytes > 0 {
		evaluateHandler = http.MaxBytesHandler(evaluateHandler, cfg.MaxRequestBodyBytes)
	}
	if sseService != nil {
		mux.Handle(ssePath, sseService.Handler())
	}
	mux.Handle("/", evaluateHandler)

	var h http.Handler = mux
	h = corsMiddleware.Handler(h)

	server := http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           h,
		ReadHeaderTimeout: 3 * time.Second,
		// slowloris/slow-client DoS protection; safe for server streams
		ReadTimeout:    5 * time.Second,
		MaxHeaderBytes: int(cfg.MaxRequestHeaderBytes),
	}

	return &Service{
		logger: cfg.Logger,
		port:   cfg.Port,
		server: &server,
		sse:    sseService,
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

	if s.sse != nil {
		group.Go(func() error {
			return s.sse.Start(gCtx)
		})
	}

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
