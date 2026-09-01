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

	syncv1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/sync/v1"
	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// flagsClient returns a client trusting the test CA when the endpoint is served over TLS.
func flagsClient(t *testing.T, caCertPath string) *http.Client {
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
		// a hand-built Transport is HTTP/1.1 only unless asked; the endpoint must work either way
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
		{title: "plaintext", port: 18070, scheme: "http", host: "localhost"},
		{
			title:    "with TLS",
			port:     18071,
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
				Port:        uint16(tt.port),
				HTTPEnabled: true,
				CertPath:    tt.certPath,
				KeyPath:     tt.keyPath,
			})
			emit()

			client := flagsClient(t, tt.caPath)
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

// The body is the flag configuration document, not the FetchAllFlags envelope around it.
func TestHTTPService_ServesTheFlagsDocumentNotTheEnvelope(t *testing.T) {
	const port = 18072

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	_, _, emit, _ := startSyncService(t, ctx, SvcConfigurations{Port: port, HTTPEnabled: true})
	emit()

	baseURL := fmt.Sprintf("http://localhost:%d", port)

	resp := getWithRetry(t, flagsClient(t, ""), baseURL+flagsPath)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	// the document itself, with its own top-level "flags" key
	assert.JSONEq(t, `{"flags":{
		"flagA":{"state":"ENABLED","defaultVariant":"false","variants":{"true":true,"false":false},
			"source":"testSource1","metadata":{"keyA":"valueA","keyDuped":"value"}},
		"flagB":{"state":"ENABLED","defaultVariant":"true","variants":{"true":true,"false":false},
			"source":"testSource2","metadata":{"keyB":"valueB","keyDuped":"value"}}
	}}`, string(body))
	assert.NotContains(t, string(body), "flagConfiguration")

	// and byte-identical to what the connect surface carries inside flag_configuration
	fetched, err := connectClient(t, baseURL, "").FetchAllFlags(
		context.Background(), connect.NewRequest(&syncv1.FetchAllFlagsRequest{}))
	require.NoError(t, err)
	assert.Equal(t, fetched.Msg.GetFlagConfiguration(), string(body))
}

// A poller should revalidate rather than re-download.
func TestHTTPService_ConditionalRequest(t *testing.T) {
	const port = 18073

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	_, _, emit, _ := startSyncService(t, ctx, SvcConfigurations{Port: port, HTTPEnabled: true})
	emit()

	url := fmt.Sprintf("http://localhost:%d%s", port, flagsPath)
	client := flagsClient(t, "")

	first := getWithRetry(t, client, url)
	first.Body.Close()
	etag := first.Header.Get("ETag")
	require.NotEmpty(t, etag)
	require.NotEmpty(t, first.Header.Get("Last-Modified"))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	req.Header.Set("If-None-Match", etag)

	revalidated, err := client.Do(req)
	require.NoError(t, err)
	defer revalidated.Body.Close()

	assert.Equal(t, http.StatusNotModified, revalidated.StatusCode)
	body, err := io.ReadAll(revalidated.Body)
	require.NoError(t, err)
	assert.Empty(t, body)
}

// HTTPEnabled=false drops the route without disturbing the connect surface beside it.
func TestHTTPService_Disabled(t *testing.T) {
	const port = 18074

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	_, _, emit, _ := startSyncService(t, ctx, SvcConfigurations{Port: port, HTTPEnabled: false})
	emit()

	baseURL := fmt.Sprintf("http://localhost:%d", port)

	// connect still answers, which is also how we know the server is up
	_ = fetchWithRetry(t, connectClient(t, baseURL, ""), connect.NewRequest(&syncv1.FetchAllFlagsRequest{}))

	resp, err := flagsClient(t, "").Get(baseURL + flagsPath)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// Last-Modified is only advertised once a sync source has actually reported in.
func TestHTTPService_ModTimeStampedOnEmit(t *testing.T) {
	const port = 18075

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	svc, _, emit, _ := startSyncService(t, ctx, SvcConfigurations{Port: port, HTTPEnabled: true})
	require.True(t, svc.modTime.get().IsZero(), "stamped before any source reported")

	emit()

	require.False(t, svc.modTime.get().IsZero(), "not stamped after sources reported")

	resp := getWithRetry(t, flagsClient(t, ""), fmt.Sprintf("http://localhost:%d%s", port, flagsPath))
	defer resp.Body.Close()
	assert.NotEmpty(t, resp.Header.Get("Last-Modified"))
}
