//nolint:dupl
package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/open-feature/flagd/core/pkg/telemetry"

	"golang.org/x/sync/errgroup"

	"github.com/open-feature/flagd/core/pkg/service/middleware"

	schemaConnectV1 "buf.build/gen/go/open-feature/flagd/bufbuild/connect-go/schema/v1/schemav1connect"
	"github.com/open-feature/flagd/core/pkg/eval"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/service"
	corsmw "github.com/open-feature/flagd/core/pkg/service/middleware/cors"
	h2cmw "github.com/open-feature/flagd/core/pkg/service/middleware/h2c"
	metricsmw "github.com/open-feature/flagd/core/pkg/service/middleware/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

const ErrorPrefix = "FlagdError:"

type ConnectService struct {
	logger                *logger.Logger
	eval                  eval.IEvaluator
	metrics               *telemetry.MetricsRecorder
	eventingConfiguration *eventingConfiguration

	server        *http.Server
	metricsServer *http.Server

	serverMtx        sync.RWMutex
	metricsServerMtx sync.RWMutex
}

// NewConnectService creates a ConnectService with provided parameters
func NewConnectService(
	logger *logger.Logger, evaluator eval.IEvaluator, mRecorder *telemetry.MetricsRecorder,
) *ConnectService {
	return &ConnectService{
		logger:  logger,
		eval:    evaluator,
		metrics: mRecorder,
		eventingConfiguration: &eventingConfiguration{
			subs: make(map[interface{}]chan service.Notification),
			mu:   &sync.RWMutex{},
		},
	}
}

// Serve serves services with provided configuration options
func (s *ConnectService) Serve(ctx context.Context, svcConf service.Configuration) error {
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return s.startServer(svcConf)
	})
	g.Go(func() error {
		return s.startMetricsServer(svcConf)
	})
	g.Go(func() error {
		<-gCtx.Done()
		s.serverMtx.RLock()
		defer s.serverMtx.RUnlock()
		if s.server != nil {
			if err := s.server.Shutdown(gCtx); err != nil {
				return err
			}
		}
		return nil
	})
	g.Go(func() error {
		<-gCtx.Done()
		s.metricsServerMtx.RLock()
		defer s.metricsServerMtx.RUnlock()
		if s.metricsServer != nil {
			if err := s.metricsServer.Shutdown(gCtx); err != nil {
				return err
			}
		}
		return nil
	})
	return g.Wait()
}

// Notify emits change event notifications for subscriptions
func (s *ConnectService) Notify(n service.Notification) {
	s.eventingConfiguration.emitToAll(n)
}

func (s *ConnectService) setupServer(svcConf service.Configuration) (net.Listener, error) {
	var lis net.Listener
	var err error
	mux := http.NewServeMux()
	if svcConf.SocketPath != "" {
		lis, err = net.Listen("unix", svcConf.SocketPath)
	} else {
		address := fmt.Sprintf(":%d", svcConf.Port)
		lis, err = net.Listen("tcp", address)
	}
	if err != nil {
		return nil, err
	}
	fes := NewFlagEvaluationService(
		s.logger.WithFields(zap.String("component", "flagservice")),
		s.eval,
		s.eventingConfiguration,
		s.metrics,
	)
	path, handler := schemaConnectV1.NewServiceHandler(fes)
	mux.Handle(path, handler)

	s.serverMtx.Lock()
	s.server = &http.Server{
		ReadHeaderTimeout: time.Second,
		Handler:           handler,
	}
	s.serverMtx.Unlock()

	// Add middlewares

	metricsMiddleware := metricsmw.NewHTTPMetric(metricsmw.Config{
		Service:        svcConf.ServiceName,
		MetricRecorder: s.metrics,
		Logger:         s.logger,
		HandlerID:      "",
	})

	s.AddMiddleware(metricsMiddleware)

	corsMiddleware := corsmw.New(svcConf.CORS)
	s.AddMiddleware(corsMiddleware)

	if svcConf.CertPath == "" || svcConf.KeyPath == "" {
		h2cMiddleware := h2cmw.New()
		s.AddMiddleware(h2cMiddleware)
	}

	return lis, nil
}

func (s *ConnectService) AddMiddleware(mw middleware.IMiddleware) {
	s.server.Handler = mw.Handler(s.server.Handler)
}

func (s *ConnectService) startServer(svcConf service.Configuration) error {
	lis, err := s.setupServer(svcConf)
	if err != nil {
		return err
	}
	s.logger.Info(fmt.Sprintf("Flag Evaluation listening at %s", lis.Addr()))
	if svcConf.CertPath != "" && svcConf.KeyPath != "" {
		if err := s.server.ServeTLS(
			lis,
			svcConf.CertPath,
			svcConf.KeyPath,
		); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	} else {
		if err := s.server.Serve(
			lis,
		); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	return nil
}

func (s *ConnectService) startMetricsServer(svcConf service.Configuration) error {
	s.logger.Info(fmt.Sprintf("metrics and probes listening at %d", svcConf.MetricsPort))
	s.metricsServerMtx.Lock()
	s.metricsServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", svcConf.MetricsPort),
		ReadHeaderTimeout: 3 * time.Second,
	}
	s.metricsServerMtx.Unlock()
	s.metricsServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/healthz":
			w.WriteHeader(http.StatusOK)
		case "/readyz":
			if svcConf.ReadinessProbe() {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusPreconditionFailed)
			}
		case "/metrics":
			promhttp.Handler().ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	if err := s.metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
