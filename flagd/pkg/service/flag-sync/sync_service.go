package sync

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"slices"
	"time"

	"buf.build/gen/go/open-feature/flagd/grpc/go/flagd/sync/v1/syncv1grpc"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type ISyncService interface {
	// Start the sync service
	Start(context.Context) error

	// Emit updates for sync listeners
	Emit(source string)
}

type SvcConfigurations struct {
	Logger                       *logger.Logger
	Port                         uint16
	Sources                      []string
	Store                        store.IStore
	ContextValues                map[string]any
	CertPath                     string
	KeyPath                      string
	SocketPath                   string
	StreamDeadline               time.Duration
	DisableSyncMetadata          bool
	MetricsRecorder              telemetry.IMetricsRecorder
	KeepAliveMinTime             time.Duration
	KeepAlivePermitWithoutStream bool

	// HTTPPort serves the flag configuration over HTTP. Zero disables it.
	HTTPPort              uint16
	ServiceName           string
	CORS                  []string
	MaxRequestHeaderBytes int64
}

type Service struct {
	listener net.Listener
	logger   *logger.Logger
	server   *grpc.Server

	certPath     string
	keyPath      string
	httpServer   *http.Server
	httpListener net.Listener
	modTime      *modTime

	startupTracker syncTracker
}

func loadTLSCredentials(certPath string, keyPath string) (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load key pair from certificate paths '%s' and '%s': %w", certPath, keyPath, err)
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
		MinVersion:   tls.VersionTLS12,
	}

	return credentials.NewTLS(config), nil
}

// keepAliveEnforcementPolicy builds the gRPC keepalive enforcement policy so
// provider keepalive pings on the long-lived SyncFlags stream aren't rejected
// with GOAWAY ENHANCE_YOUR_CALM.
func keepAliveEnforcementPolicy(cfg SvcConfigurations) keepalive.EnforcementPolicy {
	return keepalive.EnforcementPolicy{
		MinTime:             cfg.KeepAliveMinTime,
		PermitWithoutStream: cfg.KeepAlivePermitWithoutStream,
	}
}

func NewSyncService(cfg SvcConfigurations) (*Service, error) {
	var err error
	l := cfg.Logger

	serverOpts := []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.KeepaliveEnforcementPolicy(keepAliveEnforcementPolicy(cfg)),
	}

	if cfg.CertPath != "" && cfg.KeyPath != "" {
		tlsCredentials, err := loadTLSCredentials(cfg.CertPath, cfg.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS cert and key: %w", err)
		}
		serverOpts = append(serverOpts, grpc.Creds(tlsCredentials))
	}

	server := grpc.NewServer(serverOpts...)

	// Normalized on cfg itself so the HTTP server built below sees the same recorder.
	if cfg.MetricsRecorder == nil {
		cfg.MetricsRecorder = &telemetry.NoopMetricsRecorder{}
	}
	metricsRecorder := cfg.MetricsRecorder

	syncv1grpc.RegisterFlagSyncServiceServer(server, &syncHandler{
		store:               cfg.Store,
		log:                 l,
		contextValues:       cfg.ContextValues,
		deadline:            cfg.StreamDeadline,
		disableSyncMetadata: cfg.DisableSyncMetadata,
		metricsRecorder:     metricsRecorder,
	})

	var lis net.Listener
	if cfg.SocketPath != "" {
		l.Info(fmt.Sprintf("starting flag sync service at %s", cfg.SocketPath))
		lis, err = net.Listen("unix", cfg.SocketPath)
	} else {
		l.Info(fmt.Sprintf("starting flag sync service on port %d", cfg.Port))
		lis, err = net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	}
	if err != nil {
		return nil, fmt.Errorf("error creating listener: %w", err)
	}

	svc := &Service{
		listener: lis,
		logger:   l,
		server:   server,
		certPath: cfg.CertPath,
		keyPath:  cfg.KeyPath,
		modTime:  &modTime{},
		startupTracker: syncTracker{
			sources:  slices.Clone(cfg.Sources),
			doneChan: make(chan interface{}),
		},
	}

	if cfg.HTTPPort != 0 {
		// Bound here, like the gRPC listener above, so a port conflict fails while the runtime is
		// still being built.
		l.Info(fmt.Sprintf("starting flag sync http service on port %d", cfg.HTTPPort))
		svc.httpListener, err = net.Listen("tcp", fmt.Sprintf(":%d", cfg.HTTPPort))
		if err != nil {
			lis.Close()
			return nil, fmt.Errorf("error creating http listener: %w", err)
		}

		svc.httpServer = newHTTPServer(cfg, httpHandler{
			store:   cfg.Store,
			log:     l,
			modTime: svc.modTime,
		})
	}

	return svc, nil
}

func (s *Service) Start(ctx context.Context) error {
	// derive errgroup so we track ctx for exit as well as startup errors
	g, lCtx := errgroup.WithContext(ctx)

	// One gate shared by every server, so the initial-sync timeout is armed once and both are
	// released together.
	ready := make(chan struct{})
	g.Go(func() error {
		s.waitForInitialSync()
		close(ready)
		return nil
	})

	g.Go(func() error {
		<-ready

		err := s.server.Serve(s.listener)
		if err != nil {
			s.logger.Warn(fmt.Sprintf("error from sync server start: %v", err))
		}
		return nil
	})

	if s.httpServer != nil {
		g.Go(func() error {
			<-ready

			err := s.serveHTTP()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				s.logger.Warn(fmt.Sprintf("error from sync http server start: %v", err))
			}
			return nil
		})
	}

	g.Go(func() error {
		<-lCtx.Done()
		s.shutdown()

		return nil
	})

	err := g.Wait()
	if err != nil {
		return fmt.Errorf("error from sync service: %w", err)
	}

	return nil
}

func (s *Service) Emit(source string) {
	s.startupTracker.trackAndRemove(source)
	// Stamped here rather than from a store watch, which would re-marshal the whole configuration
	// just to date it. Emit can run ahead of a real content change; the ETag stays exact.
	s.modTime.set(time.Now())
}

// waitForInitialSync delays serving until every known sync source has reported its first sync, so
// clients cannot be handed a partial configuration at startup.
func (s *Service) waitForInitialSync() {
	select {
	case <-time.After(5 * time.Second):
		s.logger.Warn("timeout while waiting for all sync sources to complete their initial sync. " +
			"continuing sync service")
	case <-s.startupTracker.getDone():
	}
}

func (s *Service) shutdown() {
	s.logger.Info("shutting down gRPC sync service")
	s.server.Stop()

	if s.httpServer != nil {
		s.logger.Info("shutting down http sync service")
		if err := s.httpServer.Close(); err != nil {
			s.logger.Warn(fmt.Sprintf("error from sync http server shutdown: %v", err))
		}
	}
}

// syncTracker is a helper to track sync payloads at the startup
// It simply starts with known set of sync sources and remove
type syncTracker struct {
	sources  []string
	doneChan chan interface{}
}

func (t *syncTracker) getDone() <-chan interface{} {
	return t.doneChan
}

// trackAndRemove tracks sources and remove channel if all sources that are tracking are complete.
func (t *syncTracker) trackAndRemove(source string) {
	index := slices.Index(t.sources, source)
	if index < 0 {
		return
	}

	t.sources = slices.Delete(t.sources, index, index+1)

	if len(t.sources) == 0 {
		close(t.doneChan)
	}
}
