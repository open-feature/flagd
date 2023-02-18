package grpc

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	"buf.build/gen/go/open-feature/flagd/grpc/go/sync/v1/syncv1grpc"
	v1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/sync/v1"

	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/sync"
	"google.golang.org/grpc"
)

const (
	// Prefix for GRPC URL inputs. GRPC does not define a prefix through standard. This prefix helps to differentiate
	// remote URLs for REST APIs (i.e - HTTP) from GRPC endpoints.
	Prefix = "grpc://"

	// Connection retry constants
	// Back off period is calculated with backOffBase ^ #retry-iteration. However, when #retry-iteration count reach
	// backOffLimit, retry delay fallback to constantBackOffDelay
	backOffLimit         = 3
	backOffBase          = 4
	constantBackOffDelay = 60
)

type Sync struct {
	Target     string
	ProviderID string
	Logger     *logger.Logger
}

func (g *Sync) Init(ctx context.Context) error {
	return nil
}

func (g *Sync) Sync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// initial dial and connection. Failure here must result in a startup failure
	dial, err := grpc.DialContext(ctx, g.Target, options...)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("error establishing grpc connection: %s", err.Error()))
		return err
	}

	serviceClient := syncv1grpc.NewFlagSyncServiceClient(dial)
	syncClient, err := serviceClient.SyncFlags(ctx, &v1.SyncFlagsRequest{ProviderId: g.ProviderID})
	if err != nil {
		g.Logger.Error(fmt.Sprintf("error calling streaming operation: %s", err.Error()))
		return err
	}

	// initial stream listening
	err = g.handleFlagSync(syncClient, dataSync)
	g.Logger.Warn(fmt.Sprintf("error with stream listener: %s", err.Error()))

	// retry connection establishment
	for {
		syncClient, ok := g.connectWithRetry(ctx, options...)
		if !ok {
			// We shall exit
			return nil
		}

		err = g.handleFlagSync(syncClient, dataSync)
		if err != nil {
			g.Logger.Warn(fmt.Sprintf("error with stream listener: %s", err.Error()))
			continue
		}
	}
}

// connectWithRetry is a helper that performs exponential back off after retrying connection attempts periodically until
// a successful connection is established. Caller must not expect an error. Hence, errors are handled, logged
// internally. However, if the provided context is done, method exit with a non-ok state which must be verified by the
// caller
func (g *Sync) connectWithRetry(
	ctx context.Context, options ...grpc.DialOption,
) (syncv1grpc.FlagSyncService_SyncFlagsClient, bool) {
	var iteration int

	for {
		var sleep time.Duration
		if iteration >= backOffLimit {
			sleep = constantBackOffDelay
		} else {
			iteration++
			sleep = time.Duration(math.Pow(backOffBase, float64(iteration)))
		}

		// Block the next connection attempt and check the context
		select {
		case <-time.After(sleep * time.Second):
			break
		case <-ctx.Done():
			// context done means we shall exit
			return nil, false
		}

		g.Logger.Warn(fmt.Sprintf("connection re-establishment attempt in-progress for grpc target: %s", g.Target))

		dial, err := grpc.DialContext(ctx, g.Target, options...)
		if err != nil {
			g.Logger.Debug(fmt.Sprintf("error dialing target: %s", err.Error()))
			continue
		}

		serviceClient := syncv1grpc.NewFlagSyncServiceClient(dial)
		syncClient, err := serviceClient.SyncFlags(ctx, &v1.SyncFlagsRequest{ProviderId: g.ProviderID})
		if err != nil {
			g.Logger.Debug(fmt.Sprintf("error opening service client: %s", err.Error()))
			continue
		}

		g.Logger.Info(fmt.Sprintf("connection re-established with grpc target: %s", g.Target))
		return syncClient, true
	}
}

// handleFlagSync wraps the stream listening and push updates through dataSync channel
func (g *Sync) handleFlagSync(stream syncv1grpc.FlagSyncService_SyncFlagsClient, dataSync chan<- sync.DataSync) error {
	for {
		data, err := stream.Recv()
		if err != nil {
			return err
		}

		switch data.State {
		case v1.SyncState_SYNC_STATE_ALL:
			dataSync <- sync.DataSync{
				FlagData: data.FlagConfiguration,
				Source:   g.Target,
				Type:     sync.ALL,
			}

			g.Logger.Debug("received full configuration payload")
		case v1.SyncState_SYNC_STATE_ADD:
			dataSync <- sync.DataSync{
				FlagData: data.FlagConfiguration,
				Source:   g.Target,
				Type:     sync.ADD,
			}

			g.Logger.Debug("received an add payload")
		case v1.SyncState_SYNC_STATE_UPDATE:
			dataSync <- sync.DataSync{
				FlagData: data.FlagConfiguration,
				Source:   g.Target,
				Type:     sync.UPDATE,
			}

			g.Logger.Debug("received an update payload")
		case v1.SyncState_SYNC_STATE_DELETE:
			dataSync <- sync.DataSync{
				FlagData: data.FlagConfiguration,
				Source:   g.Target,
				Type:     sync.DELETE,
			}

			g.Logger.Debug("received a delete payload")
		case v1.SyncState_SYNC_STATE_PING:
			g.Logger.Debug("received server ping")
		default:
			g.Logger.Debug(fmt.Sprintf("received unknown state: %s", data.State.String()))
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
