package service

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/service"
	"github.com/open-feature/flagd/flagd-proxy/pkg/service/subscriptions"
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
