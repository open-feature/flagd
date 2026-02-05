package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	syncv1 "buf.build/gen/go/open-feature-forking/flagd/grpc/go/flagd/sync/v1/syncv1grpc"
	rpc "buf.build/gen/go/open-feature-forking/flagd/grpc/go/sync/v1/syncv1grpc"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/service"
	"github.com/open-feature/flagd/flagd-proxy/pkg/service/subscriptions"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	server        *http.Server
	metricsServer *http.Server
	Logger        *logger.Logger
	// oldHandler will not be required anymore when https://github.com/open-feature/flagd/issues/1088 is being worked on
	oldHandler        *oldHandler
	handler           *handler
	config            service.Configuration
	grpcServer        *grpc.Server
	metricServerReady bool
}

func NewServer(ctx context.Context, logger *logger.Logger, store subscriptions.Manager) *Server {
	theOldHandler := &oldHandler{
		logger:    logger,
		syncStore: store,
		ctx:       ctx,
	}
	theNewHandler := &handler{
		logger:    logger,
		syncStore: store,
		ctx:       ctx,
	}
	return &Server{
		oldHandler: theOldHandler,
		handler:    theNewHandler,
		Logger:     logger,
	}
}

func (s *Server) Serve(ctx context.Context, svcConf service.Configuration) error {
	s.config = svcConf
	s.metricServerReady = true

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(s.startServer)
	g.Go(s.startMetricsServer)
	g.Go(func() error {
		<-gCtx.Done()
		if s.server != nil {
			if err := s.server.Shutdown(gCtx); err != nil {
				return fmt.Errorf("error shutting down flag evaluation server: %w", err)
			}
		}
		return nil
	})
	g.Go(func() error {
		<-gCtx.Done()
		if s.metricsServer != nil {
			if err := s.metricsServer.Shutdown(gCtx); err != nil {
				return fmt.Errorf("error shutting down metrics server: %w", err)
			}
		}
		return nil
	})
	g.Go(s.captureMetrics)

	err := g.Wait()
	if err != nil {
		return fmt.Errorf("errgroup closed with error: %w", err)
	}
	return nil
}

func (s *Server) Shutdown() {
	s.metricServerReady = false

	// Stop the GRPc server gracefully
	s.grpcServer.GracefulStop()
}

func (s *Server) startServer() error {
	var lis net.Listener
	var err error
	address := fmt.Sprintf(":%d", s.config.Port)
	lis, err = net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("error setting up listener for address %s: %w", address, err)
	}

	s.grpcServer = grpc.NewServer()
	rpc.RegisterFlagSyncServiceServer(s.grpcServer, s.oldHandler)
	syncv1.RegisterFlagSyncServiceServer(s.grpcServer, s.handler)

	if err := s.grpcServer.Serve(
		lis,
	); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error returned from grpc server: %w", err)
	}

	return nil
}

func (s *Server) startMetricsServer() error {
	s.Logger.Info(fmt.Sprintf("binding metrics to %d", s.config.ManagementPort))

	grpcServer := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	mux := http.NewServeMux()
	mux.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	mux.Handle("/readyz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.metricServerReady && s.config.ReadinessProbe() {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusPreconditionFailed)
		}
	}))
	mux.Handle("/metrics", promhttp.Handler())

	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// if this is 'application/grpcServer' and HTTP2, handle with gRPC, otherwise HTTP.
		if request.ProtoMajor == 2 && strings.HasPrefix(request.Header.Get("Content-Type"), "application/grpcServer") {
			grpcServer.ServeHTTP(writer, request)
		} else {
			mux.ServeHTTP(writer, request)
			return
		}
	})

	s.metricsServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.config.ManagementPort),
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           h2c.NewHandler(handler, &http2.Server{}), // we need to use h2c to support plaintext HTTP2
	}
	if err := s.metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error returned from metrics server: %w", err)
	}
	return nil
}
