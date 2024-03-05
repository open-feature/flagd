package sync

import (
	"fmt"
	"net"

	"buf.build/gen/go/open-feature/flagd/grpc/go/flagd/sync/v1/syncv1grpc"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/store"
	"google.golang.org/grpc"
)

type ISyncService interface {
	Serve() error
	Emit()
	Shutdown()
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
	mux      *syncMultiplexer
	server   *grpc.Server
}

func NewSyncService(cfg SvcConfigurations) (*Service, error) {
	l := cfg.Logger
	mux := newMux(cfg.Store, cfg.Sources)

	server := grpc.NewServer()
	syncv1grpc.RegisterFlagSyncServiceServer(server, &syncHandler{
		mux: mux,
		log: l,
	})

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

func (s *Service) Serve() error {
	err := s.server.Serve(s.listener)
	if err != nil {
		return fmt.Errorf("error from server: %w", err)
	}

	return nil
}

func (s *Service) Emit() {
	err := s.mux.pushUpdates()
	if err != nil {
		s.logger.Warn(fmt.Sprintf("error: %v", err))
		return
	}
}

func (s *Service) Shutdown() {
	s.server.Stop()
}

// NoopSyncService as a filler implementation of the sync service
type NoopSyncService struct{}

func (n *NoopSyncService) Serve() error {
	// NOOP
	return nil
}

func (n *NoopSyncService) Emit() {
	// NOOP
}

func (n *NoopSyncService) Shutdown() {
	// NOOP
}
