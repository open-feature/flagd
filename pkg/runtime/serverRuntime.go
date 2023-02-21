package runtime

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/open-feature/flagd/pkg/server"
	"github.com/open-feature/flagd/pkg/sync/kubernetes"
	"go.uber.org/zap"

	"golang.org/x/sync/errgroup"

	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/sync"
)

type ServerConfig struct {
	Address     string
	Secure      bool
	CertPath    string
	KeyPath     string
	SyncSources string
}

type ServerRuntime struct {
	syncProvider sync.ISync
	logger       *logger.Logger
	config       ServerConfig
}

func NewServerRuntime(config ServerConfig, rootLogger *logger.Logger) (*ServerRuntime, error) {
	syncImpl, err := buildSyncImpl(config.SyncSources, rootLogger)
	if err != nil {
		return nil, err
	}

	return &ServerRuntime{
		syncProvider: syncImpl,
		logger:       rootLogger.WithFields(zap.String("component", "Server Runtime")),
		config:       config,
	}, nil
}

func (sr *ServerRuntime) Start() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Build server
	s := server.Server{
		Logger:   sr.logger.WithFields(zap.String("component", "Server")),
		Secure:   sr.config.Secure,
		CertPath: sr.config.CertPath,
		KeyPath:  sr.config.KeyPath,
		Address:  sr.config.Address,
	}

	g, gCtx := errgroup.WithContext(ctx)
	dataSync := make(chan sync.DataSync)

	// Start server
	g.Go(func() error {
		return s.Listen(gCtx, dataSync)
	})

	// Start sync provider
	g.Go(func() error {
		return sr.syncProvider.Sync(gCtx, dataSync)
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func buildSyncImpl(source string, rootLogger *logger.Logger) (sync.ISync, error) {
	if len(source) == 0 {
		return nil, errors.New("no sync provider sources provided")
	}

	switch sourceBytes := []byte(source); {
	case regCrd.Match(sourceBytes):
		rootLogger.Debug(fmt.Sprintf("using kubernetes sync-provider for: %s", source))
		return &kubernetes.Sync{
			Logger: rootLogger.WithFields(
				zap.String("component", "sync"),
				zap.String("sync", "kubernetes"),
			),
			URI: regCrd.ReplaceAllString(source, ""),
		}, nil
	default:
		return nil, fmt.Errorf("server supports only crd sync provider, but received : %s", source)
	}
}
