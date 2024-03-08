package sync

import (
	"context"
	"fmt"
	"net"

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
	Emit()
}

type SvcConfigurations struct {
	Logger  *logger.Logger
	Port    uint16
	Sources []string
	Store   *store.Flags
}

type Service struct {
	listener net.Listener
	logger   *logger.Logger
	mux      *Multiplexer
	server   *grpc.Server
}

func NewSyncService(cfg SvcConfigurations) (*Service, error) {
	l := cfg.Logger
	mux, err := NewMux(cfg.Store, cfg.Sources)
	if err != nil {
		return nil, fmt.Errorf("error initializing multiplexer: %w", err)
	}

	server := grpc.NewServer()
	syncv1grpc.RegisterFlagSyncServiceServer(server, &syncHandler{
		mux: mux,
		log: l,
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
	}, nil
}

func (s *Service) Start(ctx context.Context) error {
	// derive errgroup so we track ctx for exit as well as startup errors
	g, lCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		err := s.server.Serve(s.listener)
		if err != nil {
			s.logger.Info(fmt.Sprintf("error from sync server start: %v", err))
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

func (s *Service) Emit() {
	err := s.mux.Publish()
	if err != nil {
		s.logger.Warn(fmt.Sprintf("error while publishing sync streams: %v", err))
		return
	}
}

func (s *Service) shutdown() {
	err := s.listener.Close()
	if err != nil {
		s.logger.Warn(fmt.Sprintf("error closing the listener: %v", err))
	}
	s.server.Stop()
}

// NoopSyncService as a filler implementation of the sync service.
// This can be used as a default implementation and avoid unnecessary null checks or service enabled checks in runtime.
type NoopSyncService struct{}

func (n *NoopSyncService) Start(context.Context) error {
	// NOOP
	return nil
}

func (n *NoopSyncService) Emit() {
	// NOOP
}
