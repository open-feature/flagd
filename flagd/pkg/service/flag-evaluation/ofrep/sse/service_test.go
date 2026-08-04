package sse

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/launchdarkly/eventsource"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestService_PublishesRefetchOnChange verifies the end-to-end wiring: a client subscribed to a
// flagSetId channel receives an ADR-0008 refetchEvaluation event when that flag set changes.
func TestService_PublishesRefetchOnChange(t *testing.T) {
	log := logger.NewLogger(nil, false)
	s, err := store.NewStore(log, []string{"src1"})
	require.NoError(t, err)

	// long heartbeat so keep-alive comments do not interfere with the assertion
	svc := New(s, Config{Logger: log, HeartbeatInterval: time.Hour})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() { _ = svc.Start(ctx) }()

	ts := httptest.NewServer(svc.Handler())
	defer ts.Close()

	// allow the tracker's initial (empty) snapshot to be consumed and skipped
	time.Sleep(100 * time.Millisecond)

	stream, err := eventsource.Subscribe(ts.URL+"?channels=fs1", "")
	require.NoError(t, err)
	defer stream.Close()

	// allow the subscription to register server-side before publishing
	time.Sleep(100 * time.Millisecond)

	s.Update("src1", []model.Flag{testFlag("fs1", "a", "on")}, model.Metadata{"flagSetId": "fs1"}, false)

	select {
	case ev := <-stream.Events:
		assert.Equal(t, eventName, ev.Event())
		assert.Contains(t, ev.Data(), refetchEventType)
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for refetch event on channel fs1")
	}
}
