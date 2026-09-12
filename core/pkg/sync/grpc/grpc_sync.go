package grpc

import (
	"context"
	"fmt"
	"math"
	"sync/atomic"
	"time"

	"buf.build/gen/go/open-feature/flagd/grpc/go/flagd/sync/v1/syncv1grpc"
	v1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/sync/v1"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	grpccredential "github.com/open-feature/flagd/core/pkg/sync/grpc/credentials"
	_ "github.com/open-feature/flagd/core/pkg/sync/grpc/nameresolvers" // initialize custom resolvers e.g. envoy.Init()
	"github.com/open-feature/flagd/core/pkg/sync/syncmetrics"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
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

type Sync struct {
	GrpcDialOptionsOverride []grpc.DialOption
	CertPath                string
	CredentialBuilder       grpccredential.Builder
	Logger                  *logger.Logger
	ProviderID              string
	Secure                  bool
	Selector                string
	URI                     string
	MaxMsgSize              int
	IncrementalUpdates      bool
	Headers                 map[string]string

	// SyncMetricsRecorder, when non-nil, is used to record the source-agnostic
	// client-side sync metrics (feature_flag.flagd.sync.client.flag_config.*). Nil is
	// safe — the recorder's methods tolerate nil receivers.
	SyncMetricsRecorder *syncmetrics.Recorder

	client        FlagSyncServiceClient
	ready         atomic.Bool
	streamMetrics *clientStreamMetrics
}

func (g *Sync) Init(_ context.Context) error {
	// Build the gRPC-client-specific stream-lifecycle recorder off the global MeterProvider
	// (telemetry.NewOTelRecorder registers it globally via otel.SetMeterProvider). Nil is
	// safe — the returned recorder degrades to no-ops.
	g.streamMetrics = newClientStreamMetrics(nil)

	var rpcCon *grpc.ClientConn // Reusable client connection
	var err error

	// Instrument the outbound sync channel (flagd -> flag-server) so that standard
	// rpc.client.* metrics land on the same reader as flagd's other metrics.
	// Applied on both the default and the override paths — grpc-go supports multiple
	// stats handlers, so an override that already installs its own still composes.
	statsHandler := grpc.WithStatsHandler(otelgrpc.NewClientHandler())

	if len(g.GrpcDialOptionsOverride) > 0 {
		g.Logger.Debug("GRPC DialOptions override provided")
		opts := make([]grpc.DialOption, 0, len(g.GrpcDialOptionsOverride)+1)
		opts = append(opts, g.GrpcDialOptionsOverride...)
		opts = append(opts, statsHandler)
		rpcCon, err = grpc.NewClient(g.URI, opts...)
	} else {
		var tCredentials credentials.TransportCredentials
		tCredentials, err = g.CredentialBuilder.Build(g.Secure, g.CertPath)
		if err != nil {
			err = fmt.Errorf("error building transport credentials: %w", err)
			g.Logger.Error(err.Error())
			return err
		}

		// Set MaxMsgSize if passed
		if g.MaxMsgSize > 0 {
			g.Logger.Info(fmt.Sprintf("setting max receive message size %d bytes default 4MB", g.MaxMsgSize))
			dialOptions := grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(g.MaxMsgSize))
			rpcCon, err = grpc.NewClient(g.URI, grpc.WithTransportCredentials(tCredentials), dialOptions, statsHandler)
		} else {
			rpcCon, err = grpc.NewClient(g.URI, grpc.WithTransportCredentials(tCredentials), statsHandler)
		}
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

func (g *Sync) contextWithHeaders(ctx context.Context) context.Context {
	if len(g.Headers) == 0 {
		return ctx
	}
	pairs := make([]string, 0, len(g.Headers)*2)
	for k, v := range g.Headers {
		pairs = append(pairs, k, v)
	}
	return metadata.AppendToOutgoingContext(ctx, pairs...)
}

func (g *Sync) ReSync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	res, err := g.client.FetchAllFlags(g.contextWithHeaders(ctx), &v1.FetchAllFlagsRequest{ProviderId: g.ProviderID, Selector: g.Selector})
	if err != nil {
		err = fmt.Errorf("error fetching all flags: %w", err)
		g.Logger.Error(err.Error())
		return err
	}
	update := sync.DataSync{
		FlagData:           res.GetFlagConfiguration(),
		Source:             g.URI,
		IncrementalUpdates: g.IncrementalUpdates,
	}
	select {
	case dataSync <- update:
		g.SyncMetricsRecorder.RecordFlagConfigReceived(ctx, syncmetrics.SourceGRPC, g.URI, g.Selector)
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func (g *Sync) IsReady() bool {
	return g.ready.Load()
}

func (g *Sync) Sync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	g.Logger.Info(fmt.Sprintf("starting sync from %s", g.URI))

	// Initialize SyncFlags client. This fails if server connection establishment fails (ex:- grpc server offline)
	g.Logger.Debug(fmt.Sprintf("initial stream connection to %s", g.URI))
	syncClient, err := g.client.SyncFlags(g.contextWithHeaders(ctx), &v1.SyncFlagsRequest{ProviderId: g.ProviderID, Selector: g.Selector})
	if err != nil {
		return fmt.Errorf("unable to sync flags: %w", err)
	}
	g.streamMetrics.recordStreamOpened(ctx, g.URI, g.Selector)

	g.Logger.Debug(fmt.Sprintf("watching %s for changes", g.URI))

	// Initial stream listening. Error will be logged and continue and retry connection establishment
	err = g.handleFlagSync(ctx, syncClient, dataSync)
	g.streamMetrics.recordStreamClosed(ctx, g.URI, g.Selector)
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
		g.streamMetrics.recordStreamOpened(ctx, g.URI, g.Selector)
		g.streamMetrics.recordReconnect(ctx, g.URI, g.Selector)

		err = g.handleFlagSync(ctx, syncClient, dataSync)
		g.streamMetrics.recordStreamClosed(ctx, g.URI, g.Selector)
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

		syncClient, err := g.client.SyncFlags(g.contextWithHeaders(ctx), &v1.SyncFlagsRequest{ProviderId: g.ProviderID, Selector: g.Selector})
		if err != nil {
			g.Logger.Debug(fmt.Sprintf("error opening service client: %s", err.Error()))
			continue
		}

		g.Logger.Info(fmt.Sprintf("connection re-established with grpc target: %s", g.URI))
		return syncClient, true
	}
}

// handleFlagSync wraps the stream listening and push updates through dataSync channel
func (g *Sync) handleFlagSync(ctx context.Context, stream syncv1grpc.FlagSyncService_SyncFlagsClient, dataSync chan<- sync.DataSync) error {
	g.ready.Store(true)

	for {
		data, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("error receiving payload from stream: %w", err)
		}

		update := sync.DataSync{
			FlagData:           data.FlagConfiguration,
			SyncContext:        data.SyncContext,
			Source:             g.URI,
			Selector:           g.Selector,
			IncrementalUpdates: g.IncrementalUpdates,
		}
		select {
		case dataSync <- update:
			g.SyncMetricsRecorder.RecordFlagConfigReceived(ctx, syncmetrics.SourceGRPC, g.URI, g.Selector)
		case <-ctx.Done():
			return ctx.Err()
		}

		g.Logger.Debug("received full configuration payload")
	}
}
