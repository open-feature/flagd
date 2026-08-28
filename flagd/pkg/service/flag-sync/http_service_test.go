package sync

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

// startHTTPSyncService leaves sources un-emitted so callers can exercise the startup gate.
func startHTTPSyncService(
	t *testing.T, ctx context.Context, port int, certPath, keyPath string,
) (svc *Service, flagStore store.IStore, emit func(), done chan interface{}) {
	t.Helper()

	flagStore, sources := getSimpleFlagStore(t)

	svc, err := NewSyncService(SvcConfigurations{
		Logger:      logger.NewLogger(nil, false),
		Port:        uint16(port),
		Sources:     sources,
		Store:       flagStore,
		HTTPEnabled: true,
		CertPath:    certPath,
		KeyPath:     keyPath,
	})
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

// httpSyncClient returns a client trusting the test CA when the endpoint is served over TLS.
func httpSyncClient(t *testing.T, caCertPath string) *http.Client {
	t.Helper()

	client := &http.Client{Timeout: 2 * time.Second}
	if caCertPath == "" {
		return client
	}

	pemCA, err := os.ReadFile(caCertPath)
	require.NoError(t, err)

	pool := x509.NewCertPool()
	require.True(t, pool.AppendCertsFromPEM(pemCA))

	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: pool, MinVersion: tls.VersionTLS12},
		// a hand-built Transport is HTTP/1.1 only unless asked, which would hide h2 routing bugs
		ForceAttemptHTTP2: true,
	}
	return client
}

// getWithRetry polls until the endpoint answers, since Serve is delayed by the startup tracker.
func getWithRetry(t *testing.T, client *http.Client, url string) *http.Response {
	t.Helper()

	var lastErr error
	for i := 0; i < 50; i++ {
		resp, err := client.Get(url)
		if err == nil {
			return resp
		}
		lastErr = err
		time.Sleep(20 * time.Millisecond)
	}
	require.NoError(t, lastErr)
	return nil
}

func TestHTTPService_Serves(t *testing.T) {
	tests := []struct {
		title    string
		port     int
		certPath string
		keyPath  string
		caPath   string
		scheme   string
		host     string
	}{
		{title: "plaintext", port: 18020, scheme: "http", host: "localhost"},
		{
			title:    "with TLS",
			port:     18021,
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
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			_, _, emit, _ := startHTTPSyncService(t, ctx, tt.port, tt.certPath, tt.keyPath)
			emit()

			client := httpSyncClient(t, tt.caPath)
			resp := getWithRetry(t, client, fmt.Sprintf("%s://%s:%d%s", tt.scheme, tt.host, tt.port, flagsPath))
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
			assert.NotEmpty(t, resp.Header.Get("ETag"))
			assert.Contains(t, string(body), "flagA")
			assert.Contains(t, string(body), "flagB")
		})
	}
}

// The listener is bound eagerly so the TCP connect succeeds; it is the HTTP exchange that must not
// complete until every source has reported in.
func TestHTTPService_WaitsForInitialSync(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, _, emit, _ := startHTTPSyncService(t, ctx, 18022, "", "")

	url := fmt.Sprintf("http://localhost:%d%s", 18022, flagsPath)
	impatient := &http.Client{Timeout: 300 * time.Millisecond}

	_, err := impatient.Get(url)
	require.Error(t, err, "endpoint answered before the initial sync completed")

	emit()

	resp := getWithRetry(t, &http.Client{Timeout: 2 * time.Second}, url)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHTTPService_ShutsDownWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	_, _, emit, done := startHTTPSyncService(t, ctx, 18023, "", "")
	emit()

	url := fmt.Sprintf("http://localhost:%d%s", 18023, flagsPath)
	resp := getWithRetry(t, &http.Client{Timeout: 2 * time.Second}, url)
	resp.Body.Close()

	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for the sync service to shut down")
	}

	_, err := (&http.Client{Timeout: time.Second}).Get(url)
	assert.Error(t, err, "endpoint still answering after shutdown")
}

// HTTPEnabled=false must opt out of the endpoint without disturbing the gRPC service alongside it.
func TestHTTPService_Disabled(t *testing.T) {
	flagStore, sources := getSimpleFlagStore(t)

	svc, err := NewSyncService(SvcConfigurations{
		Logger:      logger.NewLogger(nil, false),
		Port:        19024,
		Sources:     sources,
		Store:       flagStore,
		HTTPEnabled: false,
	})
	require.NoError(t, err)

	assert.Nil(t, svc.httpServer)
	assert.Nil(t, svc.httpListener)
	assert.NotNil(t, svc.grpcListener)
}

// Last-Modified is only advertised once a sync source has actually reported in.
func TestHTTPService_ModTimeStampedOnEmit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	svc, _, emit, _ := startHTTPSyncService(t, ctx, 18025, "", "")
	require.True(t, svc.modTime.get().IsZero(), "stamped before any source reported")

	emit()

	require.False(t, svc.modTime.get().IsZero(), "not stamped after sources reported")

	url := fmt.Sprintf("http://localhost:%d%s", 18025, flagsPath)
	resp := getWithRetry(t, &http.Client{Timeout: 2 * time.Second}, url)
	defer resp.Body.Close()
	assert.NotEmpty(t, resp.Header.Get("Last-Modified"))
}

// Both servers wait on one timer, so a stalled source must not warn once per server.
func TestInitialSyncGate_WarnsOnce(t *testing.T) {
	core, logs := observer.New(zap.WarnLevel)

	flagStore, sources := getSimpleFlagStore(t)
	svc, err := NewSyncService(SvcConfigurations{
		Logger:      logger.NewLogger(zap.New(core), false),
		Port:        18030,
		Sources:     sources,
		Store:       flagStore,
		HTTPEnabled: true,
	})
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() { _ = svc.Start(ctx) }()

	// never emit, so the 5s timeout is what releases the servers
	require.Eventually(t, func() bool {
		return logs.FilterMessageSnippet("timeout while waiting for all sync sources").Len() > 0
	}, 8*time.Second, 100*time.Millisecond)

	time.Sleep(500 * time.Millisecond)
	require.Equal(t, 1, logs.FilterMessageSnippet("timeout while waiting for all sync sources").Len(),
		"initial-sync timeout warned more than once")
}

// The mux sends every h2 connection to gRPC, so ALPN ordering is what keeps the config endpoint
// reachable over TLS: an h2-capable client must be steered to http/1.1 rather than into gRPC.
func TestHTTPService_TLSNegotiatesHTTP11(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, _, emit, _ := startHTTPSyncService(t, ctx, 18026,
		"./test-cert/server-cert.pem", "./test-cert/server-key.pem")
	emit()

	// httpSyncClient sets ForceAttemptHTTP2, so this client offers h2 and would be routed to gRPC
	// if the server preferred it.
	client := httpSyncClient(t, "./test-cert/ca-cert.pem")
	resp := getWithRetry(t, client, fmt.Sprintf("https://0.0.0.0:%d%s", 18026, flagsPath))
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 1, resp.ProtoMajor, "h2 was negotiated, which routes to gRPC")
	assert.Contains(t, string(body), "flagA")
}
