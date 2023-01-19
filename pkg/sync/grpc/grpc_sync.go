package grpc

import (
	"buf.build/gen/go/kavindudodan/flagd/grpc/go/sync/v1/servicev1grpc"
	v1 "buf.build/gen/go/kavindudodan/flagd/protocolbuffers/go/sync/v1"
	"context"
	"fmt"
	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/sync"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
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

	stream, err := client.SyncFlags(context.Background(), &v1.SyncFlagsRequest{
		// todo - key from configurations
		Key: "test",
	})

	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error calling streaming operation: %s", err.Error()))
		return err
	}

	for {
		data, err := stream.Recv()

		if err == io.EOF {
			g.Logger.Error("Server streaming ended")
			// todo - attempt reconnection rather than returning error
			return err
		}

		if err != nil {
			g.Logger.Warn(fmt.Sprintf("Error with stream response: %s. Continuing receiver", err.Error()))
			continue
		}

		switch data.State {
		case v1.SyncState_SYNC_STATE_ALL:
			// todo - feed data sync to store
			continue
		case v1.SyncState_SYNC_STATE_ADD:
			// todo - feed data sync to store
			continue
		case v1.SyncState_SYNC_STATE_UPDATE:
			// todo - feed data sync to store
			continue
		case v1.SyncState_SYNC_STATE_DELETE:
			// todo - feed data sync to store
			continue
		case v1.SyncState_SYNC_STATE_CLEAN:
			// todo - feed data sync to store
			continue
		case v1.SyncState_SYNC_STATE_PING:
			g.Logger.Debug("Received server ping")
		case v1.SyncState_SYNC_STATE_UNSPECIFIED:
		default:
			g.Logger.Debug(fmt.Sprintf("Receivied unknown state: %d", data.State))
		}
	}
}
