package grpc

import (
	"context"
	"fmt"
	"math"
	msync "sync"
	"time"

	"buf.build/gen/go/open-feature/flagd/grpc/go/flagd/sync/v1/syncv1grpc"
	v1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/sync/v1"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	grpccredential "github.com/open-feature/flagd/core/pkg/sync/grpc/credentials"
	_ "github.com/open-feature/flagd/core/pkg/sync/grpc/nameresolvers" // initialize custom resolvers e.g. envoy.Init()
	"google.golang.org/grpc"
)

const (
	// Prefix for GRPC URL inputs. GRPC does not define a standard prefix. This prefix helps to differentiate remote
	// URLs for REST APIs (i.e - HTTP) from GRPC endpoints.
	Prefix          = "grpc://"
	PrefixSecure    = "grpcs://"
	SupportedScheme = "(envoy|dns|uds|xds)"

	// Connection retry constants
	// Back off period is calculated with backOffBase ^ #retry-iteration. However, when #retry-iteration count reach
	// backOffLimit, retry delay fallback to constantBackOffDelay
	backOffLimit         = 3
	backOffBase          = 4
	constantBackOffDelay = 60
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
	CertPath          string
	CredentialBuilder grpccredential.Builder
	Logger            *logger.Logger
	ProviderID        string
	Secure            bool
	Selector          string
	URI               string
	MaxMsgSize        int

	client FlagSyncServiceClient
	ready  bool
}

func (g *Sync) Init(_ context.Context) error {
	tCredentials, err := g.CredentialBuilder.Build(g.Secure, g.CertPath)
	if err != nil {
		err := fmt.Errorf("error building transport credentials: %w", err)
		g.Logger.Error(err.Error())
		return err
	}

	// Derive reusable client connection
	// Set MaxMsgSize if passed
	var rpcCon *grpc.ClientConn

	if g.MaxMsgSize > 0 {
		g.Logger.Info(fmt.Sprintf("setting max receive message size %d bytes default 4MB", g.MaxMsgSize))
		dialOptions := grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(g.MaxMsgSize))
		rpcCon, err = grpc.NewClient(g.URI, grpc.WithTransportCredentials(tCredentials), dialOptions)
	} else {
		rpcCon, err = grpc.NewClient(g.URI, grpc.WithTransportCredentials(tCredentials))
	}

	if err != nil {
		err := fmt.Errorf("error initiating grpc client connection: %w", err)
		g.Logger.Error(err.Error())
		return err
	}

	// Setup service client
	g.client = syncv1grpc.NewFlagSyncServiceClient(rpcCon)

	return nil
}

func (g *Sync) ReSync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	res, err := g.client.FetchAllFlags(ctx, &v1.FetchAllFlagsRequest{ProviderId: g.ProviderID, Selector: g.Selector})
	if err != nil {
		err = fmt.Errorf("error fetching all flags: %w", err)
		g.Logger.Error(err.Error())
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
	syncClient, err := g.client.SyncFlags(ctx, &v1.SyncFlagsRequest{ProviderId: g.ProviderID, Selector: g.Selector})
	if err != nil {
		return fmt.Errorf("unable to sync flags: %w", err)
	}

	// Initial stream listening. Error will be logged and continue and retry connection establishment
	err = g.handleFlagSync(syncClient, dataSync)
	if err == nil {
		// This should not happen as handleFlagSync expects to return with an error
		return nil
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
			sleep = constantBackOffDelay * time.Second
		} else {
			iteration++
			sleep = time.Duration(math.Pow(backOffBase, float64(iteration))) * time.Second
		}

		// Block the next connection attempt and check the context
		select {
		case <-time.After(sleep):
			break
		case <-ctx.Done():
			// context done means we shall exit
			return nil, false
		}

		g.Logger.Warn(fmt.Sprintf("connection re-establishment attempt in-progress for grpc target: %s", g.URI))

		syncClient, err := g.client.SyncFlags(ctx, &v1.SyncFlagsRequest{ProviderId: g.ProviderID, Selector: g.Selector})
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
	once.Do(func() {
		g.ready = true
	})

	for {
		data, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("error receiving payload from stream: %w", err)
		}

		dataSync <- sync.DataSync{
			FlagData: data.FlagConfiguration,
			Source:   g.URI,
			Selector: g.Selector,
			Type:     sync.ALL,
		}

		g.Logger.Debug("received full configuration payload")
	}
}
