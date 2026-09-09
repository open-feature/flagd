package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	syncv1grpc "buf.build/gen/go/open-feature/flagd/grpc/go/flagd/sync/v1/syncv1grpc"
	syncv1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/sync/v1"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/service"
	"github.com/open-feature/flagd/flagd-proxy/pkg/service/subscriptions"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Test_SyncFlags_churnOnMissingResource drives the real gRPC handler, coordinator and
// file sync with no mocks. Subscribing to a resource that does not exist makes every
// watcher fail and return, and clients retry into that window. On an unguarded subs map
// this kills the process with "fatal error: concurrent map iteration and map write" (#2030).
func Test_SyncFlags_churnOnMissingResource(t *testing.T) {
	ctx := t.Context()

	flagFile := filepath.Join(t.TempDir(), "flags.json")
	target := "file:" + flagFile

	port := freePort(t)
	log := logger.NewLogger(nil, false)
	s := NewServer(ctx, log, subscriptions.NewManager(ctx, log))
	s.config = service.Configuration{Port: port, ManagementPort: freePort(t), ReadinessProbe: func() bool { return true }}
	serving := make(chan error, 1)
	go func() { serving <- s.startServer() }()
	t.Cleanup(func() {
		select {
		case err := <-serving:
			// it exited on its own, so there is nothing to stop and Shutdown would
			// dereference a grpcServer that startServer never got as far as assigning
			if err != nil {
				t.Errorf("proxy server exited with %v", err)
			}
			return
		default:
		}
		s.Shutdown()
		if err := <-serving; err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			t.Errorf("proxy server exited with %v", err)
		}
	})

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	waitForListener(t, addr)
	// freePort only reserves a port briefly, so the listener we just reached may belong to
	// whoever took it from us. Surface that now rather than 30s later as a missing config
	select {
	case err := <-serving:
		serving <- err // hand it back: the cleanup reads it to decide whether to Shutdown
		t.Fatalf("proxy server exited before serving: %v", err)
	default:
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()
	client := syncv1grpc.NewFlagSyncServiceClient(conn)

	// a provider-like retry loop: subscribe, fail, retry
	attempt := func() (string, error) {
		sctx, scancel := context.WithTimeout(ctx, 2*time.Second)
		defer scancel()
		stream, err := client.SyncFlags(sctx, &syncv1.SyncFlagsRequest{Selector: target})
		if err != nil {
			return "", err
		}
		resp, err := stream.Recv()
		if err != nil {
			return "", err
		}
		return resp.GetFlagConfiguration(), nil
	}

	// the file does not exist, so every watcher fails and returns, opening the window
	// repeatedly. A subscription landing in it attaches to a dead multiplexer and then
	// neither errors nor delivers: it just hangs to its own deadline.
	results := make(chan bool, 400)
	for range 20 {
		var wg sync.WaitGroup
		for range 20 {
			wg.Go(func() {
				start := time.Now()
				cfg, err := attempt()
				results <- err != nil && cfg == "" && time.Since(start) > 1500*time.Millisecond
			})
		}
		wg.Wait()
	}
	close(results)
	hung, total := 0, 0
	for w := range results {
		total++
		if w {
			hung++
		}
	}
	// hangs are possible for unrelated reasons (broadcastError does a non-blocking send on
	// the handler's unbuffered channel), so the count is reported rather than asserted
	t.Logf("%d/%d subscriptions hung to the deadline while the resource was missing", hung, total)

	// the resource now exists. A subscription must start receiving it without the proxy
	// being restarted, which is the contract #2030 broke.
	const flags = `{"flags":{"probe":{"state":"ENABLED","variants":{"on":true},"defaultVariant":"on"}}}`
	if err := os.WriteFile(flagFile, []byte(flags), 0o600); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if cfg, err := attempt(); err == nil && cfg != "" {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatal("no subscription received the configuration after the resource was created")
}
