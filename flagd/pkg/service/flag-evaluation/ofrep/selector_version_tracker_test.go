package ofrep

import (
	"context"
	"testing"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockFlagStore struct {
	flags   []model.Flag
	watchCh chan struct{}
}

func (m *mockFlagStore) GetAll(_ context.Context, _ *store.Selector) ([]model.Flag, model.Metadata, error) {
	return m.flags, model.Metadata{}, nil
}

func (m *mockFlagStore) WatchSelector(_ *store.Selector) <-chan struct{} {
	return m.watchCh
}

func TestNewSelectorVersionTracker(t *testing.T) {
	log := logger.NewLogger(nil, false)
	tracker := NewSelectorVersionTracker(log, nil, 0)
	require.NotNil(t, tracker)
	defer tracker.Close()
}

func TestSelectorVersionTracker_Track(t *testing.T) {
	log := logger.NewLogger(nil, false)
	mockStore := &mockFlagStore{
		flags:   []model.Flag{{Key: "test-flag", State: "ENABLED"}},
		watchCh: make(chan struct{}),
	}
	tracker := NewSelectorVersionTracker(log, mockStore, 0)
	defer tracker.Close()

	selector := "flagSetId=test-set"

	// track the selector
	etag := tracker.Track(selector)
	require.NotEmpty(t, etag)

	// tracking again should return the same etag
	etag2 := tracker.Track(selector)
	assert.Equal(t, etag, etag2)

	// ETag should return the same value
	assert.Equal(t, etag, tracker.ETag(selector))

	// untracked selector should return empty etag
	assert.Empty(t, tracker.ETag("flagSetId=different-set"))
}

func TestSelectorVersionTracker_EmptySelector(t *testing.T) {
	log := logger.NewLogger(nil, false)
	mockStore := &mockFlagStore{
		flags:   []model.Flag{{Key: "test-flag", State: "ENABLED"}},
		watchCh: make(chan struct{}),
	}
	tracker := NewSelectorVersionTracker(log, mockStore, 0)
	defer tracker.Close()

	// test with empty selector
	etag := tracker.Track("")
	require.NotEmpty(t, etag)

	assert.Equal(t, etag, tracker.ETag(""))
}

func TestSelectorVersionTracker_ContentBasedETags(t *testing.T) {
	log := logger.NewLogger(nil, false)

	// two stores with the same content should produce the same ETag
	flags := []model.Flag{{Key: "test-flag", State: "ENABLED", DefaultVariant: "on"}}

	store1 := &mockFlagStore{flags: flags, watchCh: make(chan struct{})}
	store2 := &mockFlagStore{flags: flags, watchCh: make(chan struct{})}

	tracker1 := NewSelectorVersionTracker(log, store1, 0)
	tracker2 := NewSelectorVersionTracker(log, store2, 0)
	defer tracker1.Close()
	defer tracker2.Close()

	etag1 := tracker1.Track("")
	etag2 := tracker2.Track("")

	assert.Equal(t, etag1, etag2, "same content should produce same ETag across replicas")
}

func TestSelectorVersionTracker_DifferentContentDifferentETags(t *testing.T) {
	log := logger.NewLogger(nil, false)

	store1 := &mockFlagStore{
		flags:   []model.Flag{{Key: "flag-a", State: "ENABLED"}},
		watchCh: make(chan struct{}),
	}
	store2 := &mockFlagStore{
		flags:   []model.Flag{{Key: "flag-b", State: "ENABLED"}},
		watchCh: make(chan struct{}),
	}

	tracker1 := NewSelectorVersionTracker(log, store1, 0)
	tracker2 := NewSelectorVersionTracker(log, store2, 0)
	defer tracker1.Close()
	defer tracker2.Close()

	etag1 := tracker1.Track("")
	etag2 := tracker2.Track("")

	// different content = different ETags
	assert.NotEqual(t, etag1, etag2, "different content should produce different ETags")
}

func TestSelectorVersionTracker_WatchBasedRecompute(t *testing.T) {
	log := logger.NewLogger(nil, false)

	watchCh := make(chan struct{})
	mockStore := &mockFlagStore{
		flags:   []model.Flag{{Key: "test-flag", State: "ENABLED"}},
		watchCh: watchCh,
	}

	tracker := NewSelectorVersionTracker(log, mockStore, 0)
	defer tracker.Close()

	// track a selector (this starts a watch goroutine)
	selector := "flagSetId=test"
	initialETag := tracker.Track(selector)
	require.NotEmpty(t, initialETag)

	// change the flags
	mockStore.flags = []model.Flag{{Key: "test-flag", State: "DISABLED"}}

	// simulate a store update
	close(watchCh)

	time.Sleep(50 * time.Millisecond)

	// ETag should change because content changed
	newETag := tracker.ETag(selector)
	assert.NotEqual(t, initialETag, newETag, "ETag should change after content changes")
}

func TestSelectorVersionTracker_Eviction(t *testing.T) {
	log := logger.NewLogger(nil, false)
	mockStore := &mockFlagStore{
		flags:   []model.Flag{{Key: "test-flag", State: "ENABLED"}},
		watchCh: make(chan struct{}),
	}

	// create tracker with capacity of 2
	tracker := NewSelectorVersionTracker(log, mockStore, 2)
	defer tracker.Close()

	// track 2 selectors
	tracker.Track("selector-1")
	tracker.Track("selector-2")

	// both should be tracked
	assert.NotEmpty(t, tracker.ETag("selector-1"))
	assert.NotEmpty(t, tracker.ETag("selector-2"))

	// track a 3rd - should evict selector-1 (oldest)
	tracker.Track("selector-3")

	assert.Empty(t, tracker.ETag("selector-1"), "selector-1 should be evicted")
	assert.NotEmpty(t, tracker.ETag("selector-2"))
	assert.NotEmpty(t, tracker.ETag("selector-3"))
}
