package grpc

import (
	"context"
	"fmt"
	"io"

	"buf.build/gen/go/kavindudodan/flagd/grpc/go/sync/v1/servicev1grpc"
	v1 "buf.build/gen/go/kavindudodan/flagd/protocolbuffers/go/sync/v1"

	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/sync"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Sync struct {
	URI    string
	Key    string
	Logger *logger.Logger
}

func (g *Sync) Sync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	// todo - Add certificates and/or tokens
	dial, err := grpc.Dial("localhost:8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error establishing connection: %s", err.Error()))
		return err
	}

	client := servicev1grpc.NewFlagServiceClient(dial)

	stream, err := client.SyncFlags(context.Background(), &v1.SyncFlagsRequest{Key: g.Key})
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error calling streaming operation: %s", err.Error()))
		return err
	}

	group, localContext := errgroup.WithContext(ctx)

	group.Go(func() error {
		return g.streamHandler(stream, dataSync)
	})

	<-localContext.Done()

	err = group.Wait()
	if err == io.EOF {
		// todo - we can retry connection if this happens
		g.Logger.Info("Stream closed by the server. Exiting without retry attempts.")
		return err
	}

	return err
}

func (g *Sync) streamHandler(stream servicev1grpc.FlagService_SyncFlagsClient, dataSync chan<- sync.DataSync) error {
	for {
		data, err := stream.Recv()
		if err != nil {
			g.Logger.Warn(fmt.Sprintf("Error with stream response: %s", err.Error()))
			return err
		}

		switch data.State {
		case v1.SyncState_SYNC_STATE_ALL:
			dataSync <- sync.DataSync{
				FlagData: data.Flags,
				Source:   g.URI,
			}
			continue
		case v1.SyncState_SYNC_STATE_PING:
			g.Logger.Info("Received server ping")
		default:
			g.Logger.Info(fmt.Sprintf("Receivied unknown state: %s", data.State.String()))
		}
	}
}
