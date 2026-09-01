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

	"connectrpc.com/connect"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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

	HTTPEnabled           bool
	ServiceName           string
	CORS                  []string
	Options               []connect.HandlerOption
	MaxRequestHeaderBytes int64
	MaxRequestBodyBytes   int64
}

type Service struct {
	listener net.Listener
	logger   *logger.Logger
	server   *http.Server
	modTime  *modTime

	startupTracker syncTracker
}

func loadTLSConfig(certPath string, keyPath string) (*tls.Config, error) {
	serverCert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load key pair from certificate paths '%s' and '%s': %w", certPath, keyPath, err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
		MinVersion:   tls.VersionTLS12,
		// h2 first so gRPC clients get HTTP/2; http/1.1 keeps the other surfaces reachable
		NextProtos: []string{"h2", "http/1.1"},
	}, nil
}

func NewSyncService(cfg SvcConfigurations) (*Service, error) {
	var err error
	l := cfg.Logger

	// Normalized on cfg so the server built below sees the same recorder.
	if cfg.MetricsRecorder == nil {
		cfg.MetricsRecorder = &telemetry.NoopMetricsRecorder{}
	}

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

	stamp := &modTime{}

	// nil leaves the routes unregistered
	var flagsHandler http.Handler
	if cfg.HTTPEnabled {
		flagsHandler = httpHandler{store: cfg.Store, log: l, modTime: stamp}
	}

	server := newConnectServer(cfg, syncHandler{
		store:               cfg.Store,
		log:                 l,
		contextValues:       cfg.ContextValues,
		deadline:            cfg.StreamDeadline,
		disableSyncMetadata: cfg.DisableSyncMetadata,
		metricsRecorder:     cfg.MetricsRecorder,
	}, flagsHandler)

	// Loaded here rather than in ServeTLS so a bad path fails at construction.
	if cfg.CertPath != "" && cfg.KeyPath != "" {
		tlsConfig, tlsErr := loadTLSConfig(cfg.CertPath, cfg.KeyPath)
		if tlsErr != nil {
			lis.Close()
			return nil, fmt.Errorf("failed to load TLS cert and key: %w", tlsErr)
		}
		server.TLSConfig = tlsConfig
	}

	return &Service{
		listener: lis,
		logger:   l,
		server:   server,
		modTime:  stamp,
		startupTracker: syncTracker{
			sources:  slices.Clone(cfg.Sources),
			doneChan: make(chan interface{}),
		},
	}, nil
}

func (s *Service) Start(ctx context.Context) error {
	// derive errgroup so we track ctx for exit as well as startup errors
	g, lCtx := errgroup.WithContext(ctx)

	ready := make(chan struct{})
	g.Go(func() error {
		s.waitForInitialSync()
		close(ready)
		return nil
	})

	g.Go(func() error {
		<-ready

		err := s.serve()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
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

// serve leaves the ServeTLS paths empty because the certificates are already in TLSConfig.
func (s *Service) serve() error {
	if s.server.TLSConfig != nil {
		return s.server.ServeTLS(s.listener, "", "")
	}
	return s.server.Serve(s.listener)
}

func (s *Service) Emit(source string) {
	s.startupTracker.trackAndRemove(source)
	// Stamped here rather than from a store watch, which would re-marshal just to date it. This can
	// run ahead of a real content change; the ETag stays exact.
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
	s.logger.Info("shutting down flag sync service")

	// Close, not Shutdown: a long-lived SyncFlags stream would hold a graceful shutdown to its timeout.
	if err := s.server.Close(); err != nil {
		s.logger.Warn("error from sync server shutdown", zap.Error(err))
	}

	// Serve may never have started, in which case http.Server does not know the listener.
	s.listener.Close()
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
