package grpc

import (
	"context"
	"fmt"
	"io"
	"strings"

	"buf.build/gen/go/kavindudodan/flagd/grpc/go/sync/v1/servicev1grpc"
	v1 "buf.build/gen/go/kavindudodan/flagd/protocolbuffers/go/sync/v1"

	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/sync"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Prefix for GRPC URL inputs. GRPC does not define a prefix through standard. This prefix helps to differentiate
// remote URLs for REST APIs (i.e - HTTP) from GRPC endpoints.
const Prefix = "grpc://"

type Sync struct {
	Target string
	Key    string
	Logger *logger.Logger
}

func (g *Sync) Sync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	dial, err := grpc.Dial(g.Target, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
				Source:   g.Target,
				Type:     sync.ALL,
			}

			g.Logger.Debug("received full configuration payload")
			continue
		case v1.SyncState_SYNC_STATE_ADD:
			dataSync <- sync.DataSync{
				FlagData: data.Flags,
				Source:   g.Target,
				Type:     sync.ADD,
			}

			g.Logger.Debug("received an add payload")
			continue
		case v1.SyncState_SYNC_STATE_UPDATE:
			dataSync <- sync.DataSync{
				FlagData: data.Flags,
				Source:   g.Target,
				Type:     sync.UPDATE,
			}

			g.Logger.Debug("received an update payload")
			continue
		case v1.SyncState_SYNC_STATE_DELETE:
			dataSync <- sync.DataSync{
				FlagData: data.Flags,
				Source:   g.Target,
				Type:     sync.DELETE,
			}

			g.Logger.Debug("received a delete payload")
			continue
		case v1.SyncState_SYNC_STATE_PING:
			g.Logger.Debug("received server ping")
		default:
			g.Logger.Warn(fmt.Sprintf("receivied unknown state: %s", data.State.String()))
		}
	}
}

// URLToGRPCTarget is a helper to derive GRPC target from a provided URL
// For example, function returns the target localhost:9090 for the input grpc://localhost:9090
func URLToGRPCTarget(url string) string {
	index := strings.Split(url, Prefix)

	if len(index) == 2 {
		return index[1]
	}

	return index[0]
}
