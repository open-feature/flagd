package sse

import (
	"context"
	"testing"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testFlag(flagSetID, key, defaultVariant string) model.Flag {
	return model.Flag{
		Key:            key,
		FlagSetId:      flagSetID,
		State:          "ENABLED",
		DefaultVariant: defaultVariant,
		Variants:       map[string]any{"on": true, "off": false},
		Source:         "src1",
	}
}

func mustSelector(t *testing.T, expr string) store.Selector {
	t.Helper()
	s, err := store.NewSelector(expr)
	require.NoError(t, err)
	return s
}

func TestFingerprint_StableAndSensitive(t *testing.T) {
	a := testFlag("fs1", "a", "on")
	b := testFlag("fs1", "b", "off")

	// deterministic for identical input
	assert.Equal(t, fingerprint([]model.Flag{a, b}), fingerprint([]model.Flag{a, b}))
	// order independent
	assert.Equal(t, fingerprint([]model.Flag{a, b}), fingerprint([]model.Flag{b, a}))

	// sensitive to a definition change
	aModified := testFlag("fs1", "a", "off")
	assert.NotEqual(t, fingerprint([]model.Flag{a}), fingerprint([]model.Flag{aModified}))

	// nilFlagSetId is normalised so it does not leak the random UUID into the hash
	nilGroup := fingerprint([]model.Flag{testFlag(store.NilFlagSetId(), "a", "on")})
	emptyGroup := fingerprint([]model.Flag{testFlag("", "a", "on")})
	assert.Equal(t, emptyGroup, nilGroup)
}

func TestTracker_Update_ChangedChannels(t *testing.T) {
	tr := &Tracker{versions: map[string]version{}}

	fs1a := testFlag("fs1", "a", "on")
	fs2b := testFlag("fs2", "b", "on")

	// first update seeds every channel (this is the init snapshot Run skips)
	first := tr.update([]model.Flag{fs1a, fs2b})
	assert.Contains(t, first, allChannel)
	assert.Contains(t, first, "fs1")
	assert.Contains(t, first, "fs2")

	// no change -> no channels
	assert.Empty(t, tr.update([]model.Flag{fs1a, fs2b}))

	// change only fs1 -> catch-all + fs1, not fs2
	fs1aModified := testFlag("fs1", "a", "off")
	changed := tr.update([]model.Flag{fs1aModified, fs2b})
	assert.Contains(t, changed, allChannel)
	assert.Contains(t, changed, "fs1")
	assert.NotContains(t, changed, "fs2")

	// deleting a whole flag set notifies that set's channel + catch-all
	removed := tr.update([]model.Flag{fs1aModified})
	assert.Contains(t, removed, allChannel)
	assert.Contains(t, removed, "fs2")
	assert.NotContains(t, removed, "fs1")
}

func TestTracker_Update_NilFlagSetIdNotSubscribable(t *testing.T) {
	tr := &Tracker{versions: map[string]version{}}

	changed := tr.update([]model.Flag{testFlag(store.NilFlagSetId(), "a", "on")})
	assert.Contains(t, changed, allChannel, "catch-all must fire for flags without a flagSetId")
	assert.NotContains(t, changed, store.NilFlagSetId(), "internal nilFlagSetId must not be a subscribable channel")
}

func TestTracker_Run_NilStoreDoesNotPanic(t *testing.T) {
	tr := NewTracker(logger.NewLogger(nil, false), nil, nil)
	// Run must return promptly instead of dereferencing the nil store (which would panic in
	// the background goroutine and terminate the process).
	require.NotPanics(t, func() { tr.Run(context.Background()) })
}

// closingStore closes the watcher immediately without emitting, simulating store.Watch's
// error-close path.
type closingStore struct{}

func (closingStore) Get(context.Context, string, *store.Selector) (model.Flag, model.Metadata, error) {
	return model.Flag{}, nil, nil
}
func (closingStore) GetAll(context.Context, *store.Selector) ([]model.Flag, model.Metadata, error) {
	return nil, nil, nil
}
func (closingStore) Watch(_ context.Context, _ *store.Selector, watcher chan<- store.FlagQueryResult) {
	close(watcher)
}
func (closingStore) Update(string, []model.Flag, model.Metadata, bool) {}

func TestTracker_Run_UnexpectedCloseInvalidatesVersions(t *testing.T) {
	tr := NewTracker(logger.NewLogger(nil, false), closingStore{}, nil)
	tr.versions = map[string]version{allKey: {etag: "frozen"}}

	// context is NOT cancelled -> the watcher closing is unexpected (store error path)
	tr.Run(context.Background())

	_, _, ok := tr.Version(store.Selector{})
	assert.False(t, ok, "frozen versions must be invalidated so the bulk handler stops serving stale 304s")
}

func TestTracker_Run_ContextCancelKeepsVersions(t *testing.T) {
	tr := NewTracker(logger.NewLogger(nil, false), closingStore{}, nil)
	tr.versions = map[string]version{allKey: {etag: "current"}}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // expected shutdown

	tr.Run(ctx)

	etag, _, ok := tr.Version(store.Selector{})
	assert.True(t, ok, "versions must be retained on a normal context-cancel shutdown")
	assert.Equal(t, "current", etag)
}

func TestTracker_Version(t *testing.T) {
	tr := &Tracker{versions: map[string]version{}}
	tr.update([]model.Flag{testFlag("fs1", "a", "on")})

	// catch-all (empty selector)
	allEtag, _, ok := tr.Version(store.Selector{})
	require.True(t, ok)
	assert.NotEmpty(t, allEtag)

	// flagSetId selector
	fsEtag, _, ok := tr.Version(mustSelector(t, "flagSetId=fs1"))
	require.True(t, ok)
	assert.NotEmpty(t, fsEtag)

	// source selector is tracked for ETag lookups
	_, _, ok = tr.Version(mustSelector(t, "source=src1"))
	assert.True(t, ok)

	// unknown source -> not tracked
	_, _, ok = tr.Version(mustSelector(t, "source=missing"))
	assert.False(t, ok)

	// unknown flagSetId -> not tracked
	_, _, ok = tr.Version(mustSelector(t, "flagSetId=nope"))
	assert.False(t, ok)
}
