package flag_sync

import (
	"buf.build/gen/go/open-feature/flagd/connectrpc/go/flagd/sync/v1/syncv1connect"
	"context"
	"fmt"
	"net/http"

	"github.com/open-feature/flagd/core/pkg/logger"
)

// Service wrapper

type SvcConfigurations struct {
	Logger *logger.Logger
	Port   uint16
}

type SyncService struct {
	logger  *logger.Logger
	server  *http.Server
	sources []string
}

func NewSyncService(sources []string, cfg SvcConfigurations) SyncService {
	l := cfg.Logger
	path, h := syncv1connect.NewFlagSyncServiceHandler(&syncHandler{})
	l.Info(fmt.Sprintf("serving flag syncs at %s", path))

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: h,
	}

	return SyncService{
		l,
		&server,
		sources,
	}
}

func (s *SyncService) Serve() error {
	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (s *SyncService) Shutdown() error {
	err := s.server.Shutdown(context.Background())
	if err != nil {
		return err
	}

	return nil
}
