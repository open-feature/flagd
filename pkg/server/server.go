package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"buf.build/gen/go/open-feature/flagd/grpc/go/sync/v1/syncv1grpc"
	v1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/sync/v1"

	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/sync"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	// type of the listener
	serverListenType = "tcp"

	// Time between server to client pings
	pingDelay time.Duration = 20 * time.Second
)

type Server struct {
	Logger *logger.Logger

	Secure   bool
	CertPath string
	KeyPath  string
	Address  string
}

func (s *Server) Listen(ctx context.Context, sync <-chan sync.DataSync) error {
	options, err := s.buildOptions()
	if err != nil {
		s.Logger.Error(fmt.Sprintf("error building dial options : %s\n", err.Error()))
		return err
	}

	server := grpc.NewServer(options...)

	store := NewDataStore()
	syncv1grpc.RegisterFlagSyncServiceServer(server, &StreamHandler{
		Logger: s.Logger,
		DS:     store,
	})

	group, lcCtxt := errgroup.WithContext(ctx)

	group.Go(func() error {
		for {
			select {
			case data := <-sync:
				store.cache(dataType(data.FlagData))
			case <-lcCtxt.Done():
				s.Logger.Debug("exiting server with context done")
				server.Stop()
				return nil
			}
		}
	})

	group.Go(func() error {
		listen, err := net.Listen(serverListenType, s.Address)
		if err != nil {
			s.Logger.Error(fmt.Sprintf("error when listening to address : %s\n", err.Error()))
			return err
		}

		err = server.Serve(listen)
		if err != nil {
			s.Logger.Error(fmt.Sprintf("error when starting the server : %s\n", err.Error()))
			return err
		}

		return nil
	})

	err = group.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) buildOptions() ([]grpc.ServerOption, error) {
	var options []grpc.ServerOption

	if !s.Secure {
		return options, nil
	}

	keyPair, err := tls.LoadX509KeyPair(s.CertPath, s.KeyPath)
	if err != nil {
		return nil, err
	}

	options = append(options, grpc.Creds(credentials.NewServerTLSFromCert(&keyPair)))

	return options, nil
}

type StreamHandler struct {
	Logger *logger.Logger
	DS     *DataStore
}

func (sh *StreamHandler) SyncFlags(req *v1.SyncFlagsRequest, stream syncv1grpc.FlagSyncService_SyncFlagsServer) error {
	sh.Logger.Debug(fmt.Sprintf("stream registering for provider identifier: %s", req.ProviderId))

	subID := StorageID()
	syncChan := make(chan dataType)

	sh.DS.subscribe(subID, syncChan)
	defer sh.DS.unsubscribe(subID)

	// Initially send the current state
	err := stream.Send(&v1.SyncFlagsResponse{
		FlagConfiguration: sh.DS.currentState().string(),
		State:             v1.SyncState_SYNC_STATE_ALL,
	})
	if err != nil {
		sh.Logger.Warn(fmt.Sprintf("error writing to stream: %s", err.Error()))
		return err
	}

	// Then wait for updates
	for {
		select {
		case data := <-syncChan:
			err := stream.Send(&v1.SyncFlagsResponse{
				FlagConfiguration: data.string(),
				State:             v1.SyncState_SYNC_STATE_ALL,
			})
			if err != nil {
				sh.Logger.Warn(fmt.Sprintf("exiting stream listener, stream send failed: %s", err.Error()))
				return err
			}
		case <-time.After(pingDelay):
			err := stream.Send(&v1.SyncFlagsResponse{
				State: v1.SyncState_SYNC_STATE_PING,
			})
			if err != nil {
				sh.Logger.Warn(fmt.Sprintf("exiting stream listener, server ping failed: %s", err.Error()))
				return err
			}
		}
	}
}
