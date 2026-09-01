package sync

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"buf.build/gen/go/open-feature/flagd/connectrpc/go/flagd/sync/v1/syncv1connect"
	syncv1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/sync/v1"
	"connectrpc.com/connect"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

// startSyncService leaves sources un-emitted so callers can exercise the startup gate.
func startSyncService(
	t *testing.T, ctx context.Context, cfg SvcConfigurations,
) (svc *Service, flagStore store.IStore, emit func(), done chan interface{}) {
	t.Helper()

	flagStore, sources := getSimpleFlagStore(t)
	if cfg.Logger == nil {
		cfg.Logger = logger.NewLogger(nil, false)
	}
	cfg.Sources = sources
	cfg.Store = flagStore

	svc, err := NewSyncService(cfg)
	require.NoError(t, err)

	done = make(chan interface{})
	go func() {
		_ = svc.Start(ctx)
		close(done)
	}()

	return svc, flagStore, func() {
		for _, source := range sources {
			svc.Emit(source)
		}
	}, done
}

// connectClient talks the Connect protocol over HTTP/1.1, trusting the test CA when TLS is on.
func connectClient(t *testing.T, baseURL, caCertPath string) syncv1connect.FlagSyncServiceClient {
	t.Helper()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	if caCertPath != "" {
		pemCA, err := os.ReadFile(caCertPath)
		require.NoError(t, err)

		pool := x509.NewCertPool()
		require.True(t, pool.AppendCertsFromPEM(pemCA))

		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: pool, MinVersion: tls.VersionTLS12},
		}
	}

	return syncv1connect.NewFlagSyncServiceClient(httpClient, baseURL)
}

// fetchWithRetry polls until the endpoint answers, since Serve is delayed by the startup tracker.
func fetchWithRetry(
	t *testing.T, client syncv1connect.FlagSyncServiceClient, req *connect.Request[syncv1.FetchAllFlagsRequest],
) *connect.Response[syncv1.FetchAllFlagsResponse] {
	t.Helper()

	var lastErr error
	for i := 0; i < 50; i++ {
		resp, err := client.FetchAllFlags(context.Background(), req)
		if err == nil {
			return resp
		}
		lastErr = err
		time.Sleep(20 * time.Millisecond)
	}
	require.NoError(t, lastErr)
	return nil
}

// An ordinary HTTP client reaches FetchAllFlags, without speaking gRPC.
func TestConnectService_FetchAllFlagsOverHTTP(t *testing.T) {
	tests := []struct {
		title    string
		port     int
		certPath string
		keyPath  string
		caPath   string
		scheme   string
		host     string
	}{
		{title: "plaintext", port: 18040, scheme: "http", host: "localhost"},
		{
			title:    "with TLS",
			port:     18041,
			certPath: "./test-cert/server-cert.pem",
			keyPath:  "./test-cert/server-key.pem",
			caPath:   "./test-cert/ca-cert.pem",
			scheme:   "https",
			// the test cert carries a single SAN of IP:0.0.0.0
			host: "0.0.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			_, _, emit, _ := startSyncService(t, ctx, SvcConfigurations{
				Port:     uint16(tt.port),
				CertPath: tt.certPath,
				KeyPath:  tt.keyPath,
			})
			emit()

			client := connectClient(t, fmt.Sprintf("%s://%s:%d", tt.scheme, tt.host, tt.port), tt.caPath)
			resp := fetchWithRetry(t, client, connect.NewRequest(&syncv1.FetchAllFlagsRequest{}))

			config := resp.Msg.GetFlagConfiguration()
			assert.Contains(t, config, "flagA")
			assert.Contains(t, config, "flagB")
		})
	}
}

// Header beats body.
func TestConnectService_SelectorPrecedence(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	_, _, emit, _ := startSyncService(t, ctx, SvcConfigurations{Port: 18042})
	emit()

	client := connectClient(t, "http://localhost:18042", "")

	t.Run("body selector", func(t *testing.T) {
		resp := fetchWithRetry(t, client, connect.NewRequest(&syncv1.FetchAllFlagsRequest{
			Selector: "source=" + testSource1,
		}))

		assert.Contains(t, resp.Msg.GetFlagConfiguration(), "flagA")
		assert.NotContains(t, resp.Msg.GetFlagConfiguration(), "flagB")
	})

	t.Run("header wins", func(t *testing.T) {
		req := connect.NewRequest(&syncv1.FetchAllFlagsRequest{Selector: "source=" + testSource1})
		req.Header().Set(selectorHeaderKey, "source="+testSource2)

		resp := fetchWithRetry(t, client, req)

		assert.Contains(t, resp.Msg.GetFlagConfiguration(), "flagB")
		assert.NotContains(t, resp.Msg.GetFlagConfiguration(), "flagA")
	})

	t.Run("invalid selector", func(t *testing.T) {
		req := connect.NewRequest(&syncv1.FetchAllFlagsRequest{Selector: "invalidKey=val"})

		// warm the endpoint, so a startup race cannot be read as the expected failure
		fetchWithRetry(t, client, connect.NewRequest(&syncv1.FetchAllFlagsRequest{}))

		_, err := client.FetchAllFlags(context.Background(), req)
		require.Error(t, err)
		assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	})
}

// SyncFlags streams over the Connect protocol, not just gRPC.
func TestConnectService_SyncFlagsStream(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	_, _, emit, _ := startSyncService(t, ctx, SvcConfigurations{
		Port:          18043,
		ContextValues: map[string]any{"env": "test"},
	})
	emit()

	client := connectClient(t, "http://localhost:18043", "")
	fetchWithRetry(t, client, connect.NewRequest(&syncv1.FetchAllFlagsRequest{}))

	stream, err := client.SyncFlags(ctx, connect.NewRequest(&syncv1.SyncFlagsRequest{}))
	require.NoError(t, err)
	defer stream.Close()

	require.True(t, stream.Receive(), "expected a configuration on the stream: %v", stream.Err())
	assert.Contains(t, stream.Msg().GetFlagConfiguration(), "flagA")
	assert.Equal(t, map[string]any{"env": "test"}, stream.Msg().GetSyncContext().AsMap())
}

// The listener binds eagerly, so it is the RPC that must block until every source has reported in.
func TestConnectService_WaitsForInitialSync(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	_, _, emit, _ := startSyncService(t, ctx, SvcConfigurations{Port: 18044})

	impatient := syncv1connect.NewFlagSyncServiceClient(
		&http.Client{Timeout: 300 * time.Millisecond}, "http://localhost:18044")

	_, err := impatient.FetchAllFlags(context.Background(), connect.NewRequest(&syncv1.FetchAllFlagsRequest{}))
	require.Error(t, err, "endpoint answered before the initial sync completed")

	emit()

	resp := fetchWithRetry(t, connectClient(t, "http://localhost:18044", ""),
		connect.NewRequest(&syncv1.FetchAllFlagsRequest{}))
	assert.Contains(t, resp.Msg.GetFlagConfiguration(), "flagA")
}

func TestConnectService_ShutsDownWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())

	_, _, emit, done := startSyncService(t, ctx, SvcConfigurations{Port: 18045})
	emit()

	client := connectClient(t, "http://localhost:18045", "")
	fetchWithRetry(t, client, connect.NewRequest(&syncv1.FetchAllFlagsRequest{}))

	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for the sync service to shut down")
	}

	_, err := client.FetchAllFlags(context.Background(), connect.NewRequest(&syncv1.FetchAllFlagsRequest{}))
	assert.Error(t, err, "endpoint still answering after shutdown")
}

// connect has no default receive limit, so the service has to set one.
func TestConnectService_RejectsOversizedRequest(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	_, _, emit, _ := startSyncService(t, ctx, SvcConfigurations{
		Port:                18046,
		MaxRequestBodyBytes: 16,
	})
	emit()

	client := connectClient(t, "http://localhost:18046", "")
	fetchWithRetry(t, client, connect.NewRequest(&syncv1.FetchAllFlagsRequest{}))

	_, err := client.FetchAllFlags(context.Background(), connect.NewRequest(&syncv1.FetchAllFlagsRequest{
		Selector: "source=" + strings.Repeat("x", 512),
	}))

	require.Error(t, err)
	assert.Equal(t, connect.CodeResourceExhausted, connect.CodeOf(err))
}

// A stalled source must not release the server early, and must warn exactly once.
func TestConnectService_InitialSyncGateWarnsOnce(t *testing.T) {
	core, logs := observer.New(zap.WarnLevel)

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	// never emit, so the 5s timeout is what releases the server
	startSyncService(t, ctx, SvcConfigurations{
		Port:   18047,
		Logger: logger.NewLogger(zap.New(core), false),
	})

	require.Eventually(t, func() bool {
		return logs.FilterMessageSnippet("timeout while waiting for all sync sources").Len() > 0
	}, 8*time.Second, 100*time.Millisecond)

	time.Sleep(500 * time.Millisecond)
	require.Equal(t, 1, logs.FilterMessageSnippet("timeout while waiting for all sync sources").Len(),
		"initial-sync timeout warned more than once")
}
