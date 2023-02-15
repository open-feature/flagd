package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc/credentials/insecure"

	"buf.build/gen/go/open-feature/flagd/grpc/go/sync/v1/syncv1grpc"
	v1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/sync/v1"

	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/sync"
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

type Sync struct {
	CertPath   string
	Logger     *logger.Logger
	ProviderID string
	Source     string

	// rpcCon is a reusable grpc client connection. Lazy initialization by waiting for runtime to call Sync
	rpcCon *grpc.ClientConn
}

// Sync initialize internals and start internal sync implementation
func (g *Sync) Sync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	tCredentials, err := buildTransportCredentials(g.Source, g.CertPath)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("error building transport credentials: %s", err.Error()))
		return err
	}

	target, ok := sourceToGRPCTarget(g.Source)
	if !ok {
		return fmt.Errorf("invalid grpc source: %s", g.Source)
	}

	// Derive reusable client connection
	g.rpcCon, err = grpc.DialContext(ctx, target, grpc.WithTransportCredentials(tCredentials))
	if err != nil {
		g.Logger.Error(fmt.Sprintf("error initiating grpc client connection: %s", err.Error()))
		return err
	}

	// Cleanup when exiting the sync
	defer func(rpcCon *grpc.ClientConn) {
		err := rpcCon.Close()
		if err != nil {
			g.Logger.Warn(fmt.Sprintf("error while closing the client connection: %s", err.Error()))
		}
	}(g.rpcCon)
	return g.syncInternal(ctx, dataSync)
}

// syncInternal connects to grpc stream and push updates through sync channel. It attempts to reconnect if connection
// fails. However, initial connection must be successful. This makes sure provided configurations are valid.
func (g *Sync) syncInternal(ctx context.Context, dataSync chan<- sync.DataSync) error {
	serviceClient := syncv1grpc.NewFlagSyncServiceClient(g.rpcCon)
	syncClient, err := serviceClient.SyncFlags(ctx, &v1.SyncFlagsRequest{ProviderId: g.ProviderID})
	if err != nil {
		g.Logger.Error(fmt.Sprintf("error initializing the client: %s", err.Error()))
		return err
	}

	err = g.handleFlagSync(syncClient, dataSync)
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
func (g *Sync) connectWithRetry(ctx context.Context) (syncv1grpc.FlagSyncService_SyncFlagsClient, bool) {
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

		g.Logger.Warn(fmt.Sprintf("connection re-establishment attempt in-progress for grpc source: %s", g.Source))

		serviceClient := syncv1grpc.NewFlagSyncServiceClient(g.rpcCon)
		syncClient, err := serviceClient.SyncFlags(ctx, &v1.SyncFlagsRequest{ProviderId: g.ProviderID})
		if err != nil {
			g.Logger.Debug(fmt.Sprintf("error opening service client: %s", err.Error()))
			continue
		}

		g.Logger.Info(fmt.Sprintf("connection re-established with grpc source: %s", g.Source))
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
				Source:   g.Source,
				Type:     sync.ALL,
			}

			g.Logger.Debug("received full configuration payload")
		case v1.SyncState_SYNC_STATE_ADD:
			dataSync <- sync.DataSync{
				FlagData: data.FlagConfiguration,
				Source:   g.Source,
				Type:     sync.ADD,
			}

			g.Logger.Debug("received an add payload")
		case v1.SyncState_SYNC_STATE_UPDATE:
			dataSync <- sync.DataSync{
				FlagData: data.FlagConfiguration,
				Source:   g.Source,
				Type:     sync.UPDATE,
			}

			g.Logger.Debug("received an update payload")
		case v1.SyncState_SYNC_STATE_DELETE:
			dataSync <- sync.DataSync{
				FlagData: data.FlagConfiguration,
				Source:   g.Source,
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

// buildTransportCredentials is a helper to build grpc credentials.TransportCredentials based on source and cert path
func buildTransportCredentials(source string, certPath string) (credentials.TransportCredentials, error) {
	if strings.Contains(source, Prefix) {
		return insecure.NewCredentials(), nil
	}

	if !strings.Contains(source, PrefixSecure) {
		return nil, fmt.Errorf("invalid source. grpc source must must contain prefix %s or %s", Prefix, PrefixSecure)
	}

	var tCred credentials.TransportCredentials

	if certPath == "" {
		// Rely on CA certs provided from system
		tCred = credentials.NewTLS(&tls.Config{MinVersion: tlsVersion})
	} else {
		// Rely on provided certificate
		certBytes, err := os.ReadFile(certPath)
		if err != nil {
			return nil, err
		}

		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(certBytes) {
			return nil, fmt.Errorf("invalid certificate provided at path: %s", certPath)
		}

		tCred = credentials.NewTLS(&tls.Config{
			MinVersion: tlsVersion,
			RootCAs:    cp,
		})
	}

	return tCred, nil
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
