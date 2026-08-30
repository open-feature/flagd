package sse

import (
	"context"
	"encoding/json"
	"net/http"
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

// testSSEPath mirrors the prefix the OFREP service mounts the stream on.
const testSSEPath = "/ofrep/v1/sse"

// newStreamServer mounts the service the way the OFREP service does, so requests exercise the
// real routing: the channel is the final path segment.
func newStreamServer(t *testing.T, svc *Service) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	svc.Register(mux, testSSEPath)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

// streamURL builds the stream URL for a channel on ts.
func streamURL(ts *httptest.Server, channel string) string {
	return ts.URL + ChannelPath(testSSEPath, channel)
}

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

	ts := newStreamServer(t, svc)

	// the channel token is a selector expression, same syntax as Flagd-Selector
	stream, err := eventsource.SubscribeWithURL(streamURL(ts, "flagSetId=fs1"))
	require.NoError(t, err)
	defer stream.Close()

	// eventsource writes the 200 before registering, so Subscribe returning is not enough
	time.Sleep(100 * time.Millisecond)

	s.Update("src1", []model.Flag{testFlag("fs1", "a", "on")}, model.Metadata{"flagSetId": "fs1"}, false)

	select {
	case ev := <-stream.Events:
		assert.Equal(t, eventName, ev.Event())

		var payload refetchPayload
		require.NoError(t, json.Unmarshal([]byte(ev.Data()), &payload))
		assert.Equal(t, refetchEventType, payload.Type)
		assert.NotEmpty(t, payload.Etag, "the event must carry the config version")
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

	ts := newStreamServer(t, svc)

	bySource, err := eventsource.SubscribeWithURL(streamURL(ts, "source=src1"))
	require.NoError(t, err)
	defer bySource.Close()

	byFlagSet, err := eventsource.SubscribeWithURL(streamURL(ts, "flagSetId=fs2"))
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
	ts := newStreamServer(t, svc)

	resp, err := http.Get(streamURL(ts, "bogusKey=1"))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Empty(t, svc.Tracker().Channels(), "a rejected request must not create a subscription")
}

// TestService_ChannelComesFromPath pins the routing: the channel is the final path segment, the
// bare path is the catch-all, and a selector carrying a "/" survives the round trip escaped.
func TestService_ChannelComesFromPath(t *testing.T) {
	log := logger.NewLogger(nil, false)
	s, err := store.NewStore(log, []string{"./mySource"})
	require.NoError(t, err)

	svc := New(s, Config{Logger: log, HeartbeatInterval: time.Hour})
	ts := newStreamServer(t, svc)

	for _, channel := range []string{"", "flagSetId=fs1", "source=./mySource"} {
		stream, err := eventsource.SubscribeWithURL(streamURL(ts, channel))
		require.NoError(t, err, "channel %q", channel)
		defer stream.Close()
	}

	assert.Eventually(t, func() bool {
		return len(svc.Tracker().Channels()) == 3
	}, 3*time.Second, 5*time.Millisecond)
	assert.ElementsMatch(t, []string{"", "flagSetId=fs1", "source=./mySource"}, svc.Tracker().Channels(),
		"the subscription must key off the decoded path segment, not its escaped form")
}
