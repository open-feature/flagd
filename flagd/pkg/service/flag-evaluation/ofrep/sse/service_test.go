package sse

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/launchdarkly/eventsource"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestService_PublishesRefetchOnChange verifies the end-to-end wiring: a client subscribed with
// a flagSetId selector receives an ADR-0008 refetchEvaluation event when that flag set changes.
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

	// the channel token is a selector expression, same syntax as Flagd-Selector
	stream, err := eventsource.Subscribe(ts.URL+"?channels="+url.QueryEscape("flagSetId=fs1"), "")
	require.NoError(t, err)
	defer stream.Close()

	// eventsource writes the 200 before registering, so Subscribe returning is not enough
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

// TestService_SelectorScopedNotifications asserts a client is woken only for changes to the
// flags its selector selects, including a source= selector, which previously fell back to the
// catch-all channel and was woken by every change.
func TestService_SelectorScopedNotifications(t *testing.T) {
	log := logger.NewLogger(nil, false)
	s, err := store.NewStore(log, []string{"src1", "src2"})
	require.NoError(t, err)

	s.Update("src1", []model.Flag{testFlag("fs1", "a", "on")}, model.Metadata{"flagSetId": "fs1"}, false)
	s.Update("src2", []model.Flag{testFlag("fs2", "b", "on")}, model.Metadata{"flagSetId": "fs2"}, false)

	svc := New(s, Config{Logger: log, HeartbeatInterval: time.Hour})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() { _ = svc.Start(ctx) }()

	ts := httptest.NewServer(svc.Handler())
	defer ts.Close()

	bySource, err := eventsource.Subscribe(ts.URL+"?channels="+url.QueryEscape("source=src1"), "")
	require.NoError(t, err)
	defer bySource.Close()

	byFlagSet, err := eventsource.Subscribe(ts.URL+"?channels="+url.QueryEscape("flagSetId=fs2"), "")
	require.NoError(t, err)
	defer byFlagSet.Close()

	time.Sleep(100 * time.Millisecond)

	s.Update("src2", []model.Flag{testFlag("fs2", "b", "off")}, model.Metadata{"flagSetId": "fs2"}, false)

	select {
	case ev := <-byFlagSet.Events:
		assert.Contains(t, ev.Data(), refetchEventType)
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for a refetch event on the flagSetId=fs2 stream")
	}

	select {
	case ev := <-bySource.Events:
		t.Fatalf("source=src1 must not be notified about a src2 change, got %q", ev.Data())
	case <-time.After(500 * time.Millisecond):
	}
}

func TestService_Handler_InvalidSelectorReturns400(t *testing.T) {
	log := logger.NewLogger(nil, false)
	s, err := store.NewStore(log, []string{"src1"})
	require.NoError(t, err)

	svc := New(s, Config{Logger: log, HeartbeatInterval: time.Hour})
	ts := httptest.NewServer(svc.Handler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "?channels=" + url.QueryEscape("bogusKey=1"))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Empty(t, svc.Tracker().Channels(), "a rejected request must not create a subscription")
}
