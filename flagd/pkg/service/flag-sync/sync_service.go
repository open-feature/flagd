package sync

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
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
)

type ISyncService interface {
	// Start the sync service
	Start(context.Context) error

	// Emit updates for sync listeners
	Emit(source string)
}

type SvcConfigurations struct {
	Logger              *logger.Logger
	Port                uint16
	Sources             []string
	Store               store.IStore
	ContextValues       map[string]any
	CertPath            string
	KeyPath             string
	SocketPath          string
	StreamDeadline      time.Duration
	DisableSyncMetadata bool
	MetricsRecorder     telemetry.IMetricsRecorder
}

type Service struct {
	listener net.Listener
	logger   *logger.Logger
	server   *grpc.Server

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

func NewSyncService(cfg SvcConfigurations) (*Service, error) {
	var err error
	l := cfg.Logger

	var server *grpc.Server
	if cfg.CertPath != "" && cfg.KeyPath != "" {
		tlsCredentials, err := loadTLSCredentials(cfg.CertPath, cfg.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS cert and key: %w", err)
		}
		server = grpc.NewServer(
			grpc.Creds(tlsCredentials),
			grpc.StatsHandler(otelgrpc.NewServerHandler()),
		)
	} else {
		server = grpc.NewServer(
			grpc.StatsHandler(otelgrpc.NewServerHandler()),
		)
	}

	metricsRecorder := cfg.MetricsRecorder
	if metricsRecorder == nil {
		metricsRecorder = &telemetry.NoopMetricsRecorder{}
	}

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

	return &Service{
		listener: lis,
		logger:   l,
		server:   server,
		startupTracker: syncTracker{
			sources:  slices.Clone(cfg.Sources),
			doneChan: make(chan interface{}),
		},
	}, nil
}

func (s *Service) Start(ctx context.Context) error {
	// derive errgroup so we track ctx for exit as well as startup errors
	g, lCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		// delay server start until we see all syncs from known sync sources OR timeout
		select {
		case <-time.After(5 * time.Second):
			s.logger.Warn("timeout while waiting for all sync sources to complete their initial sync. " +
				"continuing sync service")
			break
		case <-s.startupTracker.getDone():
			break
		}

		err := s.server.Serve(s.listener)
		if err != nil {
			s.logger.Warn(fmt.Sprintf("error from sync server start: %v", err))
		}
		return nil
	})

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
}

func (s *Service) shutdown() {
	s.logger.Info("shutting down gRPC sync service")
	s.server.Stop()
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
