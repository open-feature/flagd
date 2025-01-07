package sync

import (
	"context"
	"crypto/tls"
	"fmt"
	"google.golang.org/grpc/credentials"
	"net"
	"slices"
	"time"

	"buf.build/gen/go/open-feature/flagd/grpc/go/flagd/sync/v1/syncv1grpc"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/store"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type ISyncService interface {
	// Start the sync service
	Start(context.Context) error

	// Emit updates for sync listeners
	Emit(isResync bool, source string)
}

type SvcConfigurations struct {
	Logger        *logger.Logger
	Port          uint16
	Sources       []string
	Store         *store.Flags
	ContextValues map[string]any
	CertPath      string
	KeyPath       string
}

type Service struct {
	listener net.Listener
	logger   *logger.Logger
	mux      *Multiplexer
	server   *grpc.Server

	startupTracker syncTracker
}

func loadTLSCredentials(certPath string, keyPath string) (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load key pair from given certificate paths: '%s' and '%s'", certPath, keyPath)
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

func NewSyncService(cfg SvcConfigurations) (*Service, error) {
	l := cfg.Logger
	mux, err := NewMux(cfg.Store, cfg.Sources)
	if err != nil {
		return nil, fmt.Errorf("error initializing multiplexer: %w", err)
	}

	var server *grpc.Server
	tlsCredentials, err := loadTLSCredentials(cfg.CertPath, cfg.KeyPath)
	if err != nil {
		l.Info("Couldn't load credentials or no credentials given. Starting server without TLS connection...")
		server = grpc.NewServer()
	} else {
		server = grpc.NewServer(grpc.Creds(tlsCredentials))
	}

	syncv1grpc.RegisterFlagSyncServiceServer(server, &syncHandler{
		mux:           mux,
		log:           l,
		contextValues: cfg.ContextValues,
	})

	l.Info(fmt.Sprintf("starting flag sync service on port %d", cfg.Port))
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("error creating listener: %w", err)
	}

	return &Service{
		listener: listener,
		logger:   l,
		mux:      mux,
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

func (s *Service) Emit(isResync bool, source string) {
	s.startupTracker.trackAndRemove(source)

	if !isResync {
		err := s.mux.Publish()
		if err != nil {
			s.logger.Warn(fmt.Sprintf("error while publishing sync streams: %v", err))
			return
		}
	}
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
