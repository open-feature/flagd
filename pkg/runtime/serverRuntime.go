package runtime

import (
	"context"
	"errors"
	"fmt"
	"github.com/open-feature/flagd/pkg/server"
	"github.com/open-feature/flagd/pkg/sync/kubernetes"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/sync"
)

type ServerConfig struct {
	Address     string
	Secure      bool
	CertPath    string
	KeyPath     string
	SyncSources []string
}

type ServerRuntime struct {
	syncImpl []sync.ISync
	logger   *logger.Logger
	config   ServerConfig
}

func NewServerRuntime(config ServerConfig, rootLogger *logger.Logger) (*ServerRuntime, error) {
	impls, err := buildSyncImpls(config.SyncSources, rootLogger)
	if err != nil {
		return nil, err
	}

	return &ServerRuntime{
		syncImpl: impls,
		logger:   rootLogger.WithFields(zap.String("component", "Server Runtime")),
		config:   config,
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
	dataSync := make(chan sync.DataSync, len(sr.syncImpl))

	// Start server
	g.Go(func() error {
		return s.Listen(gCtx, dataSync)
	})

	// Start sync providers
	for _, s := range sr.syncImpl {
		p := s
		g.Go(func() error {
			return p.Sync(gCtx, dataSync)
		})
	}

	<-gCtx.Done()
	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func buildSyncImpls(sources []string, rootLogger *logger.Logger) ([]sync.ISync, error) {
	if len(sources) == 0 {
		return nil, errors.New("no sync provider sources provided")
	}

	var syncs []sync.ISync
	for _, source := range sources {
		switch sourceBytes := []byte(source); {
		case regCrd.Match(sourceBytes):
			syncs = append(syncs, &kubernetes.Sync{
				Logger: rootLogger.WithFields(
					zap.String("component", "sync"),
					zap.String("sync", "kubernetes"),
				),
				URI: regCrd.ReplaceAllString(source, ""),
			})
			rootLogger.Debug(fmt.Sprintf("using kubernetes sync-provider for: %s", source))
		default:
			return nil, fmt.Errorf("server supports only crd sync providers. recieved : %s", source)
		}
	}

	return syncs, nil
}
