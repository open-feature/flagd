package grpc

import (
	"context"
	"fmt"
	"math"
	"strings"
	msync "sync"
	"time"

	"buf.build/gen/go/open-feature/flagd/grpc/go/sync/v1/syncv1grpc"
	v1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/sync/v1"

	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/sync"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	Mux        *msync.RWMutex

	syncClient syncv1grpc.FlagSyncService_SyncFlagsClient
	client     syncv1grpc.FlagSyncServiceClient
	options    []grpc.DialOption
	ready      bool
}

func (g *Sync) connectClient(ctx context.Context) error {
	// initial dial and connection. Failure here must result in a startup failure
	dial, err := grpc.DialContext(ctx, g.Target, g.options...)
	if err != nil {
		return err
	}

	g.client = syncv1grpc.NewFlagSyncServiceClient(dial)

	syncClient, err := g.client.SyncFlags(ctx, &v1.SyncFlagsRequest{ProviderId: g.ProviderID})
	if err != nil {
		g.Logger.Error(fmt.Sprintf("error calling streaming operation: %s", err.Error()))
		return err
	}
	g.syncClient = syncClient
	return nil
}

func (g *Sync) ReSync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	res, err := g.client.FetchAllFlags(ctx, &v1.FetchAllFlagsRequest{})
	if err != nil {
		g.Logger.Error(fmt.Sprintf("fetching all flags: %s", err.Error()))
		return err
	}
	dataSync <- sync.DataSync{
		FlagData: res.GetFlagConfiguration(),
		Source:   g.Target,
		Type:     sync.ALL,
	}
	return nil
}

func (g *Sync) Init(ctx context.Context) error {
	g.options = []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// initial dial and connection. Failure here must result in a startup failure
	return g.connectClient(ctx)
}

func (g *Sync) IsReady() bool {
	g.Mux.RLock()
	defer g.Mux.RUnlock()
	return g.ready
}

func (g *Sync) setReady(val bool) {
	g.Mux.Lock()
	defer g.Mux.Unlock()
	g.ready = val
}

func (g *Sync) Sync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	// initial stream listening
	g.setReady(true)
	err := g.handleFlagSync(g.syncClient, dataSync)
	if err == nil {
		return nil
	}
	g.Logger.Warn(fmt.Sprintf("error with stream listener: %s", err.Error()))
	// retry connection establishment
	for {
		g.setReady(false)
		syncClient, ok := g.connectWithRetry(ctx)
		if !ok {
			// We shall exit
			return nil
		}
		g.setReady(true)
		err = g.handleFlagSync(syncClient, dataSync)
		if err != nil {
			g.setReady(false)
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
	ctx context.Context,
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

		if err := g.connectClient(ctx); err != nil {
			g.Logger.Debug(fmt.Sprintf("error dialing target: %s", err.Error()))
			continue
		}

		syncClient, err := g.client.SyncFlags(ctx, &v1.SyncFlagsRequest{ProviderId: g.ProviderID})
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
