package sync

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	rpc "buf.build/gen/go/open-feature/flagd/bufbuild/connect-go/sync/v1/syncv1connect"
	syncv1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/sync/v1"
	"github.com/bufbuild/connect-go"
	"github.com/open-feature/flagd/core/pkg/logger"
	iservice "github.com/open-feature/flagd/core/pkg/service"
	"github.com/open-feature/flagd/core/pkg/sync"
	syncStore "github.com/open-feature/flagd/core/pkg/sync-store"
)

type SyncServer struct {
	SyncStore     syncStore.SyncStore
	server        http.Server
	Configuration SyncServerConfiguration
	Logger        *logger.Logger
}

type SyncServerConfiguration struct {
	Port        uint16
	MetricsPort uint16
}

func (s *SyncServer) Serve(ctx context.Context, svcConf iservice.Configuration) error {
	lis, err := s.setupServer()
	if err != nil {
		return err
	}

	go s.bindMetrics(svcConf)

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

func (s *SyncServer) setupServer() (net.Listener, error) {
	var lis net.Listener
	var err error
	mux := http.NewServeMux()
	address := fmt.Sprintf(":%d", s.Configuration.Port)
	lis, err = net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	path, handler := rpc.NewFlagSyncServiceHandler(s)
	mux.Handle(path, handler)

	s.server = http.Server{
		ReadHeaderTimeout: time.Second,
		Handler:           handler,
	}
	return lis, nil
}

func (s *SyncServer) bindMetrics(svcConf iservice.Configuration) {
	s.Logger.Info(fmt.Sprintf("binding metrics to %d", s.Configuration.MetricsPort))
	server := &http.Server{
		ReadHeaderTimeout: 3 * time.Second,
		Addr:              fmt.Sprintf(":%d", s.Configuration.MetricsPort),
	}
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/healthz":
			w.WriteHeader(http.StatusOK)
		case "/readyz":
			if svcConf.ReadinessProbe() {
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

func (s *SyncServer) FetchAllFlags(ctx context.Context, req *connect.Request[syncv1.FetchAllFlagsRequest]) (*connect.Response[syncv1.FetchAllFlagsResponse], error) {
	data, err := s.SyncStore.FetchAllFlags(ctx, nil, req.Msg.GetProviderId())
	if err != nil {
		return connect.NewResponse(&syncv1.FetchAllFlagsResponse{}), err
	}

	return connect.NewResponse(&syncv1.FetchAllFlagsResponse{
		FlagConfiguration: data.FlagData,
	}), nil
}

func (s *SyncServer) SyncFlags(ctx context.Context, req *connect.Request[syncv1.SyncFlagsRequest], stream *connect.ServerStream[syncv1.SyncFlagsResponse]) error {
	errChan := make(chan error)
	dataSync := make(chan sync.DataSync)
	s.SyncStore.RegisterSubscription(ctx, req.Msg.GetProviderId(), req, dataSync, errChan)
	for {
		select {
		case e := <-errChan:
			return e
		case d := <-dataSync:
			if err := stream.Send(&syncv1.SyncFlagsResponse{
				FlagConfiguration: d.FlagData,
				State:             syncv1.SyncState(d.Type + 1),
			}); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}
