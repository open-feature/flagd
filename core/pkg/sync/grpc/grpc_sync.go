package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	credentials2 "github.com/open-feature/flagd/core/pkg/sync/grpc/credentials"
	"math"
	"strings"
	msync "sync"
	"time"

	"buf.build/gen/go/open-feature/flagd/grpc/go/sync/v1/syncv1grpc"
	v1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/sync/v1"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	"google.golang.org/grpc"
)

const (
	// Prefix for GRPC URL inputs. GRPC does not define a standard prefix. This prefix helps to differentiate remote
	// URLs for REST APIs (i.e - HTTP) from GRPC endpoints.
	Prefix       = "grpc://"
	PrefixSecure = "grpcs://"

	// Connection retry constants
	// Back off period is calculated with backOffBase ^ #retry-iteration. However, when #retry-iteration count reach
	// backOffLimit, retry delay fallback to constantBackOffDelay
	backOffLimit         = 3
	backOffBase          = 4
	constantBackOffDelay = 60

	tlsVersion = tls.VersionTLS12
)

// type aliases for interfaces required by this component - needed for mock generation with gomock

type FlagSyncServiceClient interface {
	syncv1grpc.FlagSyncServiceClient
}
type FlagSyncServiceClientResponse interface {
	syncv1grpc.FlagSyncService_SyncFlagsClient
}

var once msync.Once

type Sync struct {
	URI               string
	ProviderID        string
	CertPath          string
	Logger            *logger.Logger
	CredentialBuilder credentials2.Builder

	client FlagSyncServiceClient
	ready  bool
}

func (g *Sync) Init(ctx context.Context) error {
	tCredentials, err := g.CredentialBuilder.Build(g.URI, g.CertPath)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("error building transport credentials: %s", err.Error()))
		return err
	}

	target, ok := sourceToGRPCTarget(g.URI)
	if !ok {
		return fmt.Errorf("invalid grpc source: %s", g.URI)
	}

	// Derive reusable client connection
	rpcCon, err := grpc.DialContext(ctx, target, grpc.WithTransportCredentials(tCredentials))
	if err != nil {
		g.Logger.Error(fmt.Sprintf("error initiating grpc client connection: %s", err.Error()))
		return err
	}

	// Setup service client
	g.client = syncv1grpc.NewFlagSyncServiceClient(rpcCon)

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
		Source:   g.URI,
		Type:     sync.ALL,
	}
	return nil
}

func (g *Sync) IsReady() bool {
	return g.ready
}

func (g *Sync) Sync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	// Initialize SyncFlags client. This fails if server connection establishment fails (ex:- grpc server offline)
	syncClient, err := g.client.SyncFlags(ctx, &v1.SyncFlagsRequest{ProviderId: g.ProviderID})
	if err != nil {
		return err
	}

	// Initial stream listening. Error will be logged and continue and retry connection establishment
	err = g.handleFlagSync(syncClient, dataSync)
	if err == nil {
		// This should not happen as handleFlagSync expects to return with an error
		return err
	}

	g.Logger.Warn(fmt.Sprintf("error with stream listener: %s", err.Error()))

	// retry connection establishment
	for {
		syncClient, ok := g.connectWithRetry(ctx)
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

		g.Logger.Warn(fmt.Sprintf("connection re-establishment attempt in-progress for grpc target: %s", g.URI))

		syncClient, err := g.client.SyncFlags(ctx, &v1.SyncFlagsRequest{ProviderId: g.ProviderID})
		if err != nil {
			g.Logger.Debug(fmt.Sprintf("error opening service client: %s", err.Error()))
			continue
		}

		g.Logger.Info(fmt.Sprintf("connection re-established with grpc target: %s", g.URI))
		return syncClient, true
	}
}

// handleFlagSync wraps the stream listening and push updates through dataSync channel
func (g *Sync) handleFlagSync(stream syncv1grpc.FlagSyncService_SyncFlagsClient, dataSync chan<- sync.DataSync) error {
	// Set ready state once only
	once.Do(func() {
		g.ready = true
	})

	for {
		data, err := stream.Recv()
		if err != nil {
			return err
		}

		switch data.State {
		case v1.SyncState_SYNC_STATE_ALL:
			dataSync <- sync.DataSync{
				FlagData: data.FlagConfiguration,
				Source:   g.URI,
				Type:     sync.ALL,
			}

			g.Logger.Debug("received full configuration payload")
		case v1.SyncState_SYNC_STATE_ADD:
			dataSync <- sync.DataSync{
				FlagData: data.FlagConfiguration,
				Source:   g.URI,
				Type:     sync.ADD,
			}

			g.Logger.Debug("received an add payload")
		case v1.SyncState_SYNC_STATE_UPDATE:
			dataSync <- sync.DataSync{
				FlagData: data.FlagConfiguration,
				Source:   g.URI,
				Type:     sync.UPDATE,
			}

			g.Logger.Debug("received an update payload")
		case v1.SyncState_SYNC_STATE_DELETE:
			dataSync <- sync.DataSync{
				FlagData: data.FlagConfiguration,
				Source:   g.URI,
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

// sourceToGRPCTarget is a helper to derive GRPC target from a provided URL
// For example, function returns the target localhost:9090 for the input grpc://localhost:9090
func sourceToGRPCTarget(url string) (string, bool) {
	var separator string

	switch {
	case strings.Contains(url, Prefix):
		separator = Prefix
	case strings.Contains(url, PrefixSecure):
		separator = PrefixSecure
	default:
		return "", false
	}

	index := strings.Split(url, separator)

	if len(index) == 2 && len(index[1]) != 0 {
		return index[1], true
	}

	return "", false
}
