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
	evaluationV2 "buf.build/gen/go/open-feature/flagd/connectrpc/go/flagd/evaluation/v2/evaluationv2connect"
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

	flagdSchemaPrefix   = "/flagd"
	flagdV2SchemaPrefix = "/flagd.evaluation.v2"
)

// bufSwitchHandler combines the handlers of the old and new evaluation schemas
// this way we support both the new (v2) and the old (v1 and deprecated) schemas
// NOTE: this will not be required anymore when it is time to work on https://github.com/open-feature/flagd/issues/1088
type bufSwitchHandler struct {
	old http.Handler
	v1  http.Handler
	v2  http.Handler
}

func (b bufSwitchHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if strings.HasPrefix(request.URL.Path, flagdV2SchemaPrefix) {
		b.v2.ServeHTTP(writer, request)
	} else if strings.HasPrefix(request.URL.Path, flagdSchemaPrefix) {
		b.v1.ServeHTTP(writer, request)
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

	readinessEnabled bool
}

// NewConnectService creates a ConnectService with provided parameters
func NewConnectService(
	logger *logger.Logger, evaluator evaluator.IEvaluator, store store.IStore, mRecorder telemetry.IMetricsRecorder,
) *ConnectService {
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
		return s.startServer(gCtx, svcConf)
	})
	g.Go(func() error {
		return s.startMetricsServer(gCtx, svcConf)
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
	)

	marshalOpts := WithJSON(
		// json parsing configuration - we emit "unpopulated" fields (falsy fields are not dropped)
		protojson.MarshalOptions{EmitUnpopulated: true},
		protojson.UnmarshalOptions{DiscardUnknown: true},
	)

	_, oldHandler := schemaConnectV1.NewServiceHandler(fes, append(svcConf.Options, marshalOpts)...)

	// register handler for new flag evaluation schema (v1)

	v1Fes := NewFlagEvaluationService(s.logger.WithFields(zap.String("component", "flagd.evaluation.v1")),
		s.eval,
		s.eventingConfiguration,
		s.metrics,
		svcConf.ContextValues,
		svcConf.HeaderToContextKeyMappings,
		svcConf.StreamDeadline,
	)

	_, v1Handler := evaluationV1.NewServiceHandler(v1Fes, append(svcConf.Options, marshalOpts)...)

	// register handler for evaluation v2 schema (with optional value and variant)

	v2Fes := NewFlagEvaluationServiceV2(s.logger.WithFields(zap.String("component", "flagd.evaluation.v2")),
		s.eval,
		s.eventingConfiguration,
		s.metrics,
		svcConf.ContextValues,
		svcConf.HeaderToContextKeyMappings,
		svcConf.StreamDeadline,
	)

	_, v2Handler := evaluationV2.NewServiceHandler(v2Fes, append(svcConf.Options, marshalOpts)...)

	bs := bufSwitchHandler{
		old: oldHandler,
		v1:  v1Handler,
		v2:  v2Handler,
	}

	var svcHandler http.Handler = bs
	if svcConf.MaxRequestBodyBytes > 0 {
		svcHandler = http.MaxBytesHandler(svcHandler, svcConf.MaxRequestBodyBytes)
	}

	s.server = &http.Server{
		ReadHeaderTimeout: time.Second,
		Handler:           svcHandler,
		MaxHeaderBytes:    int(svcConf.MaxRequestHeaderBytes),
	}

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

func serveWithShutdown(ctx context.Context, server *http.Server, serveFn func() error) error {
	errChan := make(chan error, 1)
	go func() { errChan <- serveFn() }()

	select {
	case err := <-errChan:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	case <-ctx.Done():
		// use a fresh context; ctx is already cancelled
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("error shutting down server: %w", err)
		}
		// wait for server to fully stop
		<-errChan
		return nil
	}
}

func (s *ConnectService) startServer(ctx context.Context, svcConf service.Configuration) error {
	lis, err := s.setupServer(svcConf)
	if err != nil {
		return err
	}
	s.logger.Info(fmt.Sprintf("Flag IResolver listening at %s", lis.Addr()))

	if svcConf.CertPath != "" && svcConf.KeyPath != "" {
		return serveWithShutdown(ctx, s.server, func() error {
			return s.server.ServeTLS(lis, svcConf.CertPath, svcConf.KeyPath)
		})
	}
	return serveWithShutdown(ctx, s.server, func() error {
		return s.server.Serve(lis)
	})
}

func (s *ConnectService) startMetricsServer(ctx context.Context, svcConf service.Configuration) error {
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

	s.metricsServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", svcConf.ManagementPort),
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           h2c.NewHandler(handler, &http2.Server{}), // we need to use h2c to support plaintext HTTP2
	}

	return serveWithShutdown(ctx, s.metricsServer, s.metricsServer.ListenAndServe)
}
