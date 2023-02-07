package grpc

import (
	"context"
	"fmt"
	"io"
	"math"
	"strings"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	"buf.build/gen/go/open-feature/flagd/grpc/go/sync/v1/syncv1grpc"
	v1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/sync/v1"

	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/sync"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

const (
	// Prefix for GRPC URL inputs. GRPC does not define a prefix through standard. This prefix helps to differentiate
	// remote URLs for REST APIs (i.e - HTTP) from GRPC endpoints.
	Prefix = "grpc://"

	// Connection retry constants
	// Backoff period is calculated with backOffBase ^ #retry-iteration. However, when backoffLimit is reached, fallback
	// to constantBackoffDelay
	backoffLimit         = 3
	backOffBase          = 4
	constantBackoffDelay = 60
)

type Sync struct {
	Target     string
	ProviderID string
	Logger     *logger.Logger
}

func (g *Sync) Sync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// initial dial and connection. Failure here must result in a startup failure
	dial, err := grpc.DialContext(ctx, g.Target, options...)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error establishing grpc connection: %s", err.Error()))
		return err
	}

	serviceClient := syncv1grpc.NewFlagSyncServiceClient(dial)
	syncClient, err := serviceClient.SyncFlags(context.Background(), &v1.SyncFlagsRequest{ProviderId: g.ProviderID})
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error calling streaming operation: %s", err.Error()))
		return err
	}

	// initial stream listening
	err = g.streamListener(ctx, syncClient, dataSync)
	if err != nil {
		g.Logger.Warn(fmt.Sprintf("Error with stream listener: %s", err.Error()))
	}

	// retry connection establishment
	for {
		g.Logger.Warn(fmt.Sprintf("Connection re-establishment attempt in-progress for grpc target: %s", g.Target))

		syncClient = g.connectWithRetry(ctx, options...)
		err = g.streamListener(ctx, syncClient, dataSync)
		if err != nil {
			g.Logger.Warn(fmt.Sprintf("Error with stream listener: %s", err.Error()))
			continue
		}
	}
}

// connectWithRetry is a helper to perform exponential backoff till provided configurations and then retry connection
// periodically till a successful connection is established
func (g *Sync) connectWithRetry(
	ctx context.Context, options ...grpc.DialOption,
) syncv1grpc.FlagSyncService_SyncFlagsClient {
	var iteration int

	for {
		var sleep time.Duration
		if iteration >= backoffLimit {
			sleep = constantBackoffDelay
		} else {
			iteration++
			sleep = time.Duration(math.Pow(backOffBase, float64(iteration)))
		}

		time.Sleep(sleep * time.Second)

		dial, err := grpc.DialContext(ctx, g.Target, options...)
		if err != nil {
			g.Logger.Debug(fmt.Sprintf("Error dialing target: %s", err.Error()))
			continue
		}

		serviceClient := syncv1grpc.NewFlagSyncServiceClient(dial)
		syncClient, err := serviceClient.SyncFlags(context.Background(), &v1.SyncFlagsRequest{ProviderId: g.ProviderID})
		if err != nil {
			g.Logger.Debug(fmt.Sprintf("Error openning service client: %s", err.Error()))
			continue
		}

		g.Logger.Info(fmt.Sprintf("Connection re-established with grpc target: %s", g.Target))
		return syncClient
	}
}

// streamListener wraps the grpc listening on provided stream and push updates through dataSync channel
func (g *Sync) streamListener(
	ctx context.Context, stream syncv1grpc.FlagSyncService_SyncFlagsClient, dataSync chan<- sync.DataSync,
) error {
	group, localContext := errgroup.WithContext(ctx)
	group.Go(func() error {
		return g.handleFlagSync(stream, dataSync)
	})

	<-localContext.Done()

	err := group.Wait()
	if err == io.EOF {
		g.Logger.Info("Stream closed by the server")
		return err
	}

	return err
}

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
			g.Logger.Debug(fmt.Sprintf("receivied unknown state: %s", data.State.String()))
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
