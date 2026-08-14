package service

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/service"
	"github.com/open-feature/flagd/flagd-proxy/pkg/service/subscriptions"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// freePort reserves an ephemeral TCP port from the kernel and returns it after
// closing the listener, so the server under test can bind to it.
func freePort(t *testing.T) uint16 {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to reserve a free port: %v", err)
	}
	port := lis.Addr().(*net.TCPAddr).Port
	if err := lis.Close(); err != nil {
		t.Fatalf("failed to release reserved port: %v", err)
	}
	return uint16(port)
}

// waitForListener blocks until a TCP connection to addr succeeds or the
// deadline elapses.
func waitForListener(t *testing.T, addr string) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if err == nil {
			conn.Close()
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("server did not start listening on %s within the deadline", addr)
}

// TestMetricsServerRoutesGRPCHealth verifies that a real gRPC request (which
// carries Content-Type "application/grpc") sent to the management port is
// routed to the registered gRPC health server via h2c, instead of falling
// through to the plain HTTP mux. With the Content-Type prefix typo the h2c
// dispatcher never matched real gRPC traffic, so the health server was
// unreachable on the management port.
func TestMetricsServerRoutesGRPCHealth(t *testing.T) {
	ctx := context.Background()

	managementPort := freePort(t)
	s := NewServer(ctx, logger.NewLogger(nil, false), subscriptions.NewManager(ctx, logger.NewLogger(nil, false)))
	s.config = service.Configuration{
		ManagementPort: managementPort,
		ReadinessProbe: func() bool { return true },
	}
	s.metricServerReady = true

	serveErr := make(chan error, 1)
	go func() { serveErr <- s.startMetricsServer() }()

	managementAddr := fmt.Sprintf("127.0.0.1:%d", managementPort)
	waitForListener(t, managementAddr)

	// insecure transport credentials speak plaintext HTTP/2, which is what the
	// server's h2c handler expects.
	conn, err := grpc.NewClient(managementAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to create gRPC client: %v", err)
	}
	defer conn.Close()

	checkCtx, checkCancel := context.WithTimeout(ctx, 3*time.Second)
	defer checkCancel()

	resp, err := grpc_health_v1.NewHealthClient(conn).Check(checkCtx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("gRPC health check on the management port failed; the request was not routed to the gRPC health server: %v", err)
	}
	if resp.GetStatus() != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Errorf("expected health status %v, got %v", grpc_health_v1.HealthCheckResponse_SERVING, resp.GetStatus())
	}

	if err := s.metricsServer.Shutdown(ctx); err != nil {
		t.Errorf("failed to shut down metrics server: %v", err)
	}
	if err := <-serveErr; err != nil {
		t.Errorf("startMetricsServer returned an unexpected error: %v", err)
	}
}

// keepAliveConfig builds a Configuration carrying only the keepalive settings,
// so the tests below can name scenarios without restating the field mapping.
func keepAliveConfig(minTime time.Duration, permitWithoutStream bool) service.Configuration {
	return service.Configuration{
		KeepAliveMinTime:             minTime,
		KeepAlivePermitWithoutStream: permitWithoutStream,
	}
}

func TestKeepAliveEnforcementPolicy(t *testing.T) {
	for name, cfg := range map[string]service.Configuration{
		"flag defaults":                  keepAliveConfig(30*time.Second, true),
		"custom min time":                keepAliveConfig(10*time.Second, true),
		"permit without stream disabled": keepAliveConfig(30*time.Second, false),
		"zero min time deferred to grpc": keepAliveConfig(0, false),
	} {
		t.Run(name, func(t *testing.T) {
			policy := keepAliveEnforcementPolicy(cfg)

			if policy.MinTime != cfg.KeepAliveMinTime {
				t.Errorf("expected MinTime %v, got %v", cfg.KeepAliveMinTime, policy.MinTime)
			}
			if policy.PermitWithoutStream != cfg.KeepAlivePermitWithoutStream {
				t.Errorf("expected PermitWithoutStream %v, got %v", cfg.KeepAlivePermitWithoutStream, policy.PermitWithoutStream)
			}
		})
	}
}

// TestServerKeepAliveEnforcement proves the enforcement policy is applied to the
// sync gRPC server. The Go gRPC client clamps its keepalive interval to a 10s
// minimum, so we drive the HTTP/2 connection directly with a raw framer to flood
// keepalive pings and observe how the server's configured policy reacts: a
// permissive policy tolerates the pings, while a strict one tears the connection
// down with GOAWAY ENHANCE_YOUR_CALM. Without the policy the server falls back to
// gRPC's defaults (MinTime 5m, PermitWithoutStream false), which is what severed
// the providers' sync streams roughly every 90 seconds.
func TestServerKeepAliveEnforcement(t *testing.T) {
	tests := []struct {
		name                string
		cfg                 service.Configuration
		wantEnhanceYourCalm bool
	}{
		{"permissive policy tolerates frequent pings", keepAliveConfig(time.Millisecond, true), false},
		{"strict min time rejects frequent pings with GOAWAY", keepAliveConfig(time.Hour, true), true},
		{"pings without an active stream are rejected when not permitted", keepAliveConfig(time.Millisecond, false), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)

			port := freePort(t)
			s := NewServer(ctx, logger.NewLogger(nil, false), subscriptions.NewManager(ctx, logger.NewLogger(nil, false)))
			s.config = tt.cfg
			s.config.Port = port
			s.config.ReadinessProbe = func() bool { return true }

			serveErr := make(chan error, 1)
			go func() { serveErr <- s.startServer() }()

			addr := fmt.Sprintf("127.0.0.1:%d", port)
			waitForListener(t, addr)
			t.Cleanup(func() {
				// s.grpcServer is only assigned once startServer's listener is
				// up. freePort releases the port before the server rebinds it,
				// so a competing bind leaves the field nil while waitForListener
				// still succeeds against the competitor — stopping it
				// unconditionally would panic and bury the real listen error.
				if s.grpcServer != nil {
					s.grpcServer.Stop()
				}
				if err := <-serveErr; err != nil {
					t.Errorf("startServer returned an unexpected error: %v", err)
				}
			})

			if got := floodKeepalivePings(t, addr); got != tt.wantEnhanceYourCalm {
				t.Errorf("expected GOAWAY ENHANCE_YOUR_CALM %v, got %v", tt.wantEnhanceYourCalm, got)
			}
		})
	}
}

// floodKeepalivePings opens a raw HTTP/2 connection to the sync server and sends
// keepalive PING frames in rapid succession without opening a stream. It returns
// true if the server responds with GOAWAY ENHANCE_YOUR_CALM within a short
// window, and false if the connection is left healthy.
func floodKeepalivePings(t *testing.T, addr string) bool {
	t.Helper()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("failed to dial %s: %v", addr, err)
	}
	defer conn.Close()

	if _, err := io.WriteString(conn, http2.ClientPreface); err != nil {
		t.Fatalf("failed to write the HTTP/2 client preface: %v", err)
	}

	framer := http2.NewFramer(conn, conn)

	var writeMu sync.Mutex
	write := func(fn func() error) {
		writeMu.Lock()
		defer writeMu.Unlock()
		_ = fn()
	}

	// the client connection preface must be followed by a SETTINGS frame
	write(func() error { return framer.WriteSettings() })

	serverReady, enhanceYourCalm := watchForGoAway(framer, write)

	// Wait for the server's own SETTINGS frame before flooding. A successful dial
	// only means the kernel accepted the connection into the listen backlog, which
	// happens before startServer builds the gRPC server; pings sent in that window
	// would be buffered and then read back-to-back, registering as a MinTime
	// violation no matter how permissive the policy is. The server's SETTINGS
	// frame proves the HTTP/2 transport exists and its reader is live, so the
	// pings below are paced by this loop rather than by the backlog.
	select {
	case <-serverReady:
	case <-time.After(3 * time.Second):
		t.Fatal("server did not send its SETTINGS frame within the deadline")
	}

	// flood pings faster than any strict MinTime; grpc-go sends GOAWAY after more
	// than two "ping strikes", so a handful of rapid pings is enough
	var pingData [8]byte
	for i := 0; i < 8; i++ {
		write(func() error { return framer.WritePing(false, pingData) })
		select {
		case got := <-enhanceYourCalm:
			return got
		case <-time.After(10 * time.Millisecond):
		}
	}

	select {
	case got := <-enhanceYourCalm:
		return got
	case <-time.After(2 * time.Second):
		return false
	}
}

// watchForGoAway reads frames from the framer in the background, acking server
// SETTINGS. It returns two channels: serverReady is closed once the server's
// SETTINGS frame has been seen — deliberately never closed on any other path,
// so a server that fails to come up trips the caller's timeout rather than
// letting it flood a dead connection — and enhanceYourCalm reports whether the
// first GOAWAY carries the ENHANCE_YOUR_CALM code.
func watchForGoAway(framer *http2.Framer, write func(func() error)) (<-chan struct{}, <-chan bool) {
	serverReady := make(chan struct{})
	// a server may legitimately send more than one non-ACK SETTINGS frame
	var readyOnce sync.Once
	markReady := func() { readyOnce.Do(func() { close(serverReady) }) }

	enhanceYourCalm := make(chan bool, 1)
	go func() {
		for {
			frame, err := framer.ReadFrame()
			if err != nil {
				return
			}
			switch f := frame.(type) {
			case *http2.SettingsFrame:
				if !f.IsAck() {
					write(func() error { return framer.WriteSettingsAck() })
					markReady()
				}
			case *http2.GoAwayFrame:
				enhanceYourCalm <- f.ErrCode == http2.ErrCodeEnhanceYourCalm
				return
			}
		}
	}()
	return serverReady, enhanceYourCalm
}
