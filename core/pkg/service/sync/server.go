package sync

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	rpc "buf.build/gen/go/open-feature/flagd/grpc/go/sync/v1/syncv1grpc"
	"github.com/open-feature/flagd/core/pkg/logger"
	iservice "github.com/open-feature/flagd/core/pkg/service"
	syncStore "github.com/open-feature/flagd/core/pkg/sync-store"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Server struct {
	server        *http.Server
	metricsServer *http.Server
	Logger        *logger.Logger
	handler       *handler
	config        iservice.Configuration
}

func NewServer(ctx context.Context, logger *logger.Logger) *Server {
	syncStore := syncStore.NewSyncStore(ctx, logger)
	return &Server{
		handler: &handler{
			logger:    logger,
			syncStore: syncStore,
		},
		Logger: logger,
	}
}

func (s *Server) Serve(ctx context.Context, svcConf iservice.Configuration) error {
	s.config = svcConf

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

func (s *Server) startServer() error {
	var lis net.Listener
	var err error
	address := fmt.Sprintf(":%d", s.config.Port)
	lis, err = net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("error setting up listener for address %s: %w", address, err)
	}
	grpcServer := grpc.NewServer()
	rpc.RegisterFlagSyncServiceServer(grpcServer, s.handler)

	if err := grpcServer.Serve(
		lis,
	); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error returned from grpc server: %w", err)
	}

	return nil
}

func (s *Server) startMetricsServer() error {
	s.Logger.Info(fmt.Sprintf("binding metrics to %d", s.config.MetricsPort))
	s.metricsServer = &http.Server{
		ReadHeaderTimeout: 3 * time.Second,
		Addr:              fmt.Sprintf(":%d", s.config.MetricsPort),
	}
	s.metricsServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/healthz":
			w.WriteHeader(http.StatusOK)
		case "/readyz":
			if s.config.ReadinessProbe() {
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
		return fmt.Errorf("error returned from metrics server: %w", err)
	}
	return nil
}
