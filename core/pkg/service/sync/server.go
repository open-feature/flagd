package sync

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	rpc "buf.build/gen/go/open-feature/flagd/bufbuild/connect-go/sync/v1/syncv1connect"
	"github.com/open-feature/flagd/core/pkg/logger"
	iservice "github.com/open-feature/flagd/core/pkg/service"
	syncStore "github.com/open-feature/flagd/core/pkg/sync-store"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	server  http.Server
	Logger  *logger.Logger
	handler handler
	config  iservice.Configuration
}

func NewServer(ctx context.Context, logger *logger.Logger) *Server {
	syncStore := syncStore.NewSyncStore(ctx, logger)
	return &Server{
		handler: handler{
			logger:    logger,
			syncStore: syncStore,
		},
		Logger: logger,
	}
}

func (s *Server) Serve(ctx context.Context, svcConf iservice.Configuration) error {
	s.config = svcConf
	lis, err := s.setupServer()
	if err != nil {
		return err
	}

	go s.bindMetrics()

	errChan := make(chan error, 1)
	go func() {
		if err := s.server.Serve(
			lis,
		); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return s.server.Shutdown(ctx)
	}
}

func (s *Server) setupServer() (net.Listener, error) {

	var lis net.Listener
	var err error
	mux := http.NewServeMux()
	address := fmt.Sprintf(":%d", s.config.Port)
	lis, err = net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	path, handler := rpc.NewFlagSyncServiceHandler(&s.handler)
	mux.Handle(path, handler)

	s.server = http.Server{
		ReadHeaderTimeout: time.Second,
		Handler:           h2c.NewHandler(mux, &http2.Server{}),
	}
	return lis, nil
}

func (s *Server) bindMetrics() {
	s.Logger.Info(fmt.Sprintf("binding metrics to %d", s.config.MetricsPort))
	server := &http.Server{
		ReadHeaderTimeout: 3 * time.Second,
		Addr:              fmt.Sprintf(":%d", s.config.MetricsPort),
	}
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/healthz":
			w.WriteHeader(http.StatusOK)
		case "/readyz":
			if s.config.ReadinessProbe() {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusPreconditionFailed)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
