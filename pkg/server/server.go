package server

import (
	"buf.build/gen/go/open-feature/flagd/grpc/go/sync/v1/syncv1grpc"
	v1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/sync/v1"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/sync"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
)

type Server struct {
	Logger *logger.Logger

	Secure   bool
	CertPath string
	KeyPath  string
	Address  string
}

func (s *Server) Listen(ctx context.Context, sync <-chan sync.DataSync) error {
	listen, err := net.Listen("tcp", s.Address)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Error when listening to address : %s\n", err.Error()))
		return err
	}

	options, err := s.buildOptions()
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Error building dial options : %s\n", err.Error()))
		return err
	}

	server := grpc.NewServer(options...)
	defer server.Stop()

	listener := NewListener()

	group, lcCtxt := errgroup.WithContext(ctx)

	group.Go(func() error {
		for {
			select {
			case data := <-sync:
				fmt.Printf("New data :%s", data)
				listener.persist(data.FlagData)
			case <-ctx.Done():
				return nil
			}
		}
	})

	syncv1grpc.RegisterFlagSyncServiceServer(server, &internal{
		Logger: s.Logger,
		Ls:     &listener,
	})

	group.Go(func() error {
		err = server.Serve(listen)
		if err != nil {
			s.Logger.Error(fmt.Sprintf("Error when starting the server : %s\n", err.Error()))
			return err
		}

		return nil
	})

	<-lcCtxt.Done()
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

type internal struct {
	Logger *logger.Logger
	Ls     *Listener
}

func (i *internal) SyncFlags(req *v1.SyncFlagsRequest, stream syncv1grpc.FlagSyncService_SyncFlagsServer) error {
	i.Logger.Info(fmt.Sprintf("Request with ID: %s", req.ProviderId))

	// Initially send the current state
	err := stream.Send(&v1.SyncFlagsResponse{
		FlagConfiguration: i.Ls.currentState(),
		State:             v1.SyncState_SYNC_STATE_ALL,
	})
	if err != nil {
		return err
	}

	emit := i.Ls.getEmit()

	// Then wait for updates
	for {
		select {
		case _ = <-emit:
			stream.Send(&v1.SyncFlagsResponse{
				FlagConfiguration: i.Ls.currentState(),
				State:             v1.SyncState_SYNC_STATE_ALL,
			})
		}
	}

}

// todo we need a sync mechanism better than listener

type Listener struct {
	emit chan string
	data string
}

func NewListener() Listener {
	return Listener{
		emit: make(chan string),
		data: "",
	}
}

func (s *Listener) persist(input string) {
	s.data = input
	s.emit <- s.data
}

func (s *Listener) getEmit() <-chan string {
	return s.emit
}

func (s *Listener) currentState() string {
	return s.data
}
