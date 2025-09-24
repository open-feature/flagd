//nolint:dupl
package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	evaluationV1 "buf.build/gen/go/open-feature/flagd/connectrpc/go/flagd/evaluation/v1/evaluationv1connect"
	schemaConnectV1 "buf.build/gen/go/open-feature/flagd/connectrpc/go/schema/v1/schemav1connect"
	"github.com/open-feature/flagd/core/pkg/evaluator"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/service"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	"github.com/open-feature/flagd/flagd/pkg/service/middleware"
	corsmw "github.com/open-feature/flagd/flagd/pkg/service/middleware/cors"
	h2cmw "github.com/open-feature/flagd/flagd/pkg/service/middleware/h2c"
	metricsmw "github.com/open-feature/flagd/flagd/pkg/service/middleware/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	ErrorPrefix = "FlagdError:"

	flagdSchemaPrefix = "/flagd"
)

// bufSwitchHandler combines the handlers of the old and new evaluation schema and combines them into one
// this way we support both the new and the (deprecated) old schemas until only the new schema is supported
// NOTE: this will not be required anymore when it is time to work on https://github.com/open-feature/flagd/issues/1088
type bufSwitchHandler struct {
	old http.Handler
	new http.Handler
}

func (b bufSwitchHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if strings.HasPrefix(request.URL.Path, flagdSchemaPrefix) {
		b.new.ServeHTTP(writer, request)
	} else {
		b.old.ServeHTTP(writer, request)
	}
}

type ConnectService struct {
	logger                *logger.Logger
	eval                  evaluator.IEvaluator
	metrics               telemetry.IMetricsRecorder
	eventingConfiguration IEvents

	server        *http.Server
	metricsServer *http.Server

	serverMtx        sync.RWMutex
	metricsServerMtx sync.RWMutex

	readinessEnabled bool

	selectorFallbackKey string
}

// NewConnectService creates a ConnectService with provided parameters
func NewConnectService(logger *logger.Logger, evaluator evaluator.IEvaluator, store store.IStore, mRecorder telemetry.IMetricsRecorder, selectorFallbackKey string) *ConnectService {
	cs := &ConnectService{
		logger:  logger,
		eval:    evaluator,
		metrics: &telemetry.NoopMetricsRecorder{},
		eventingConfiguration: &eventingConfiguration{
			subs:   make(map[interface{}]chan service.Notification),
			mu:     &sync.RWMutex{},
			store:  store,
			logger: logger,
		},
		selectorFallbackKey: selectorFallbackKey,
	}
	if mRecorder != nil {
		cs.metrics = mRecorder
	}
	return cs
}

// Serve serves services with provided configuration options
func (s *ConnectService) Serve(ctx context.Context, svcConf service.Configuration) error {
	g, gCtx := errgroup.WithContext(ctx)
	s.readinessEnabled = true

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
				return fmt.Errorf("error returned from flag evaluation server shutdown: %w", err)
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
				return fmt.Errorf("error returned from metrics server shutdown: %w", err)
			}
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return fmt.Errorf("errgroup closed with error: %w", err)
	}
	return nil
}

// Notify emits change event notifications for subscriptions
func (s *ConnectService) Notify(n service.Notification) {
	s.eventingConfiguration.EmitToAll(n)
}

// nolint: funlen
func (s *ConnectService) setupServer(svcConf service.Configuration) (net.Listener, error) {
	var lis net.Listener
	var err error

	if svcConf.SocketPath != "" {
		lis, err = net.Listen("unix", svcConf.SocketPath)
	} else {
		address := fmt.Sprintf(":%d", svcConf.Port)
		lis, err = net.Listen("tcp", address)
	}
	if err != nil {
		return nil, fmt.Errorf("error creating listener for flag evaluation service: %w", err)
	}

	// register handler for old flag evaluation schema
	// can be removed as a part of https://github.com/open-feature/flagd/issues/1088
	fes := NewOldFlagEvaluationService(
		s.logger.WithFields(zap.String("component", "flagservice")),
		s.eval,
		s.eventingConfiguration,
		s.metrics,
		svcConf.ContextValues,
		s.selectorFallbackKey,
	)

	marshalOpts := WithJSON(
		// json parsing configuration - we emit "unpopulated" fields (falsy fields are not dropped)
		protojson.MarshalOptions{EmitUnpopulated: true},
		protojson.UnmarshalOptions{DiscardUnknown: true},
	)

	_, oldHandler := schemaConnectV1.NewServiceHandler(fes, append(svcConf.Options, marshalOpts)...)

	// register handler for new flag evaluation schema

	newFes := NewFlagEvaluationService(s.logger.WithFields(zap.String("component", "flagd.evaluation.v1")),
		s.eval,
		s.eventingConfiguration,
		s.metrics,
		svcConf.ContextValues,
		svcConf.HeaderToContextKeyMappings,
		svcConf.StreamDeadline,
		s.selectorFallbackKey,
	)

	_, newHandler := evaluationV1.NewServiceHandler(newFes, append(svcConf.Options, marshalOpts)...)

	bs := bufSwitchHandler{
		old: oldHandler,
		new: newHandler,
	}

	s.serverMtx.Lock()
	s.server = &http.Server{
		ReadHeaderTimeout: time.Second,
		Handler:           bs,
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

func (s *ConnectService) Shutdown() {
	s.readinessEnabled = false
	s.eventingConfiguration.EmitToAll(service.Notification{
		Type: service.Shutdown,
		Data: map[string]interface{}{},
	})
}

func (s *ConnectService) startServer(svcConf service.Configuration) error {
	lis, err := s.setupServer(svcConf)
	if err != nil {
		return err
	}
	s.logger.Info(fmt.Sprintf("Flag IResolver listening at %s", lis.Addr()))
	if svcConf.CertPath != "" && svcConf.KeyPath != "" {
		if err := s.server.ServeTLS(
			lis,
			svcConf.CertPath,
			svcConf.KeyPath,
		); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("error returned from flag evaluation server: %w", err)
		}
	} else {
		if err := s.server.Serve(
			lis,
		); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("error returned from flag evaluation server: %w", err)
		}
	}
	return nil
}

func (s *ConnectService) startMetricsServer(svcConf service.Configuration) error {
	s.logger.Info(fmt.Sprintf("metrics and probes listening at %d", svcConf.ManagementPort))

	srv := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(srv, health.NewServer())

	mux := http.NewServeMux()
	mux.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	mux.Handle("/readyz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.readinessEnabled && svcConf.ReadinessProbe() {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusPreconditionFailed)
		}
	}))
	mux.Handle("/metrics", promhttp.Handler())

	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// if this is 'application/grpc' and HTTP2, handle with gRPC, otherwise HTTP.
		if request.ProtoMajor == 2 && strings.HasPrefix(request.Header.Get("Content-Type"), "application/grpc") {
			srv.ServeHTTP(writer, request)
		} else {
			mux.ServeHTTP(writer, request)
			return
		}
	})

	s.metricsServerMtx.Lock()
	s.metricsServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", svcConf.ManagementPort),
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           h2c.NewHandler(handler, &http2.Server{}), // we need to use h2c to support plaintext HTTP2
	}
	s.metricsServerMtx.Unlock()

	if err := s.metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error returned from metrics server: %w", err)
	}
	return nil
}
