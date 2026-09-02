package sync

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"buf.build/gen/go/open-feature/flagd/grpc/go/flagd/sync/v1/syncv1grpc"
	v1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/sync/v1"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
)

// Regression test for issue #1998: frequent pings with no active stream must not be GOAWAYed. The
// run outlasts ReadTimeout to catch an idle teardown too. h2c and TLS use different h2 servers.
func TestSyncService_ToleratesAggressiveClientKeepalive(t *testing.T) {
	tests := []struct {
		title    string
		port     int
		certPath string
		keyPath  string
		host     string
	}{
		{title: "h2c", port: 18050, host: "localhost"},
		{
			title:    "TLS",
			port:     18052,
			certPath: "./test-cert/server-cert.pem",
			keyPath:  "./test-cert/server-key.pem",
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

			var tlsConfig *tls.Config
			if tt.certPath != "" {
				tlsConfig = probeTLSConfig(t, "./test-cert/ca-cert.pem")
			}

			probe := dialKeepaliveProbe(t, fmt.Sprintf("%s:%d", tt.host, tt.port), tlsConfig)
			defer probe.close()

			// no stream open, so this covers PermitWithoutStream as well as MinTime
			const pings = 60
			for i := 0; i < pings; i++ {
				probe.ping()
				if goAway := probe.goAway(); goAway != nil {
					t.Fatalf("server sent GOAWAY (%v) after %d keepalive pings", goAway.ErrCode, i+1)
				}
				time.Sleep(100 * time.Millisecond)
			}

			require.Nil(t, probe.goAway(), "server tore the connection down over client keepalive")
			assert.Greater(t, int(probe.acked.Load()), pings/2,
				"server stopped acking keepalive pings, so the connection was not healthy")
		})
	}
}

// The other side of #1998: an idle stream must not be reaped either.
func TestSyncService_StreamOutlivesReadTimeout(t *testing.T) {
	const port = 18051

	ctx := t.Context()

	_, flagStore, emit, _ := startSyncService(t, ctx, SvcConfigurations{Port: port})
	emit()

	client := getSyncClient(t, "", "", false, port, ctx)

	var stream syncv1grpc.FlagSyncService_SyncFlagsClient
	require.Eventually(t, func() bool {
		s, err := client.SyncFlags(ctx, &v1.SyncFlagsRequest{})
		if err != nil {
			return false
		}
		stream = s
		return true
	}, 2*time.Second, 20*time.Millisecond)

	_, err := stream.Recv()
	require.NoError(t, err, "initial configuration")

	// http.Server.ReadTimeout is 5s; an idle stream must survive it
	time.Sleep(6 * time.Second)

	updated := []model.Flag{{
		Key:            "flagA",
		State:          "ENABLED",
		DefaultVariant: "true",
		Variants:       testVariants,
	}}
	flagStore.Update(testSource1, updated, model.Metadata{}, false)

	resp, err := stream.Recv()
	require.NoError(t, err, "stream was reaped while idle")
	assert.Contains(t, resp.GetFlagConfiguration(), "flagA")
}

// Control: a bare grpc.NewServer is the server #1998 described, so the probe must see its GOAWAY.
// Without this, a probe that stopped detecting anything would still look green.
func TestKeepaliveProbe_DetectsGrpcGoEnforcement(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	server := grpc.NewServer()
	go func() { _ = server.Serve(lis) }()
	t.Cleanup(server.Stop)

	probe := dialKeepaliveProbe(t, lis.Addr().String(), nil)
	defer probe.close()

	var goAway *http2.GoAwayFrame
	for i := 0; i < 40 && goAway == nil; i++ {
		probe.ping()
		time.Sleep(50 * time.Millisecond)
		goAway = probe.goAway()
	}

	require.NotNil(t, goAway, "probe failed to observe grpc-go's keepalive enforcement")
	assert.Equal(t, http2.ErrCodeEnhanceYourCalm, goAway.ErrCode)
}

// keepaliveProbe drives a raw connection, since grpc-go clamps client keepalive to 10s minimum.
type keepaliveProbe struct {
	conn   net.Conn
	framer *http2.Framer
	writeM sync.Mutex

	acked   atomic.Int64
	goAwayM sync.Mutex
	goAwayF *http2.GoAwayFrame
}

// probeTLSConfig offers only h2, so a negotiation failure is loud rather than a silent downgrade.
func probeTLSConfig(t *testing.T, caCertPath string) *tls.Config {
	t.Helper()

	pemCA, err := os.ReadFile(caCertPath)
	require.NoError(t, err)

	pool := x509.NewCertPool()
	require.True(t, pool.AppendCertsFromPEM(pemCA))

	return &tls.Config{RootCAs: pool, MinVersion: tls.VersionTLS12, NextProtos: []string{"h2"}}
}

func dialKeepaliveProbe(t *testing.T, addr string, tlsConfig *tls.Config) *keepaliveProbe {
	t.Helper()

	// the listener is bound in NewSyncService but Serve is delayed by the startup tracker
	var conn net.Conn
	var err error
	for i := 0; i < 50; i++ {
		if tlsConfig != nil {
			conn, err = tls.Dial("tcp", addr, tlsConfig)
		} else {
			conn, err = net.Dial("tcp", addr)
		}
		if err == nil {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	require.NoError(t, err)

	if tlsConn, ok := conn.(*tls.Conn); ok {
		require.Equal(t, "h2", tlsConn.ConnectionState().NegotiatedProtocol,
			"ALPN did not settle on h2, so this is not exercising the HTTP/2 path")
	}

	_, err = io.WriteString(conn, http2.ClientPreface)
	require.NoError(t, err)

	p := &keepaliveProbe{conn: conn, framer: http2.NewFramer(conn, conn)}

	// the preface must be followed by a SETTINGS frame
	p.write(func() error { return p.framer.WriteSettings() })
	go p.readLoop()

	return p
}

func (p *keepaliveProbe) write(fn func() error) {
	p.writeM.Lock()
	defer p.writeM.Unlock()
	_ = fn()
}

func (p *keepaliveProbe) ping() {
	var data [8]byte
	p.write(func() error { return p.framer.WritePing(false, data) })
}

func (p *keepaliveProbe) goAway() *http2.GoAwayFrame {
	p.goAwayM.Lock()
	defer p.goAwayM.Unlock()
	return p.goAwayF
}

func (p *keepaliveProbe) close() {
	p.conn.Close()
}

// readLoop acks server SETTINGS, counts ping acks, and records the first GOAWAY.
func (p *keepaliveProbe) readLoop() {
	for {
		frame, err := p.framer.ReadFrame()
		if err != nil {
			return
		}
		switch f := frame.(type) {
		case *http2.SettingsFrame:
			if !f.IsAck() {
				p.write(func() error { return p.framer.WriteSettingsAck() })
			}
		case *http2.PingFrame:
			if f.IsAck() {
				p.acked.Add(1)
			}
		case *http2.GoAwayFrame:
			p.goAwayM.Lock()
			p.goAwayF = f
			p.goAwayM.Unlock()
			return
		}
	}
}
