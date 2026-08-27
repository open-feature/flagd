package sse

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/launchdarkly/eventsource"

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

func newTestStore(t *testing.T) *store.Store {
	t.Helper()
	s, err := store.NewStore(logger.NewLogger(nil, false), []string{"src1", "src2"})
	require.NoError(t, err)
	return s
}

// newTracker builds a Tracker over a real eventsource server, so publishes are exercised.
func newTracker(t *testing.T, s store.IStore) *Tracker {
	t.Helper()
	es := eventsource.NewServer()
	t.Cleanup(es.Close)

	tr := NewTracker(logger.NewLogger(nil, false), s, es)
	t.Cleanup(tr.Close)
	return tr
}

// subscribe registers a channel and ties its release to the test cleanup.
func subscribe(t *testing.T, tr *Tracker, expr string) {
	t.Helper()
	release, err := tr.Subscribe(expr, mustSelector(t, expr))
	require.NoError(t, err)
	t.Cleanup(release)
}

// versionOf polls until the channel has a seeded ETag, since seeding is asynchronous.
func versionOf(t *testing.T, tr *Tracker, channel string) (string, int64) {
	t.Helper()
	var etag string
	var lastModified int64
	require.Eventually(t, func() bool {
		var ok bool
		etag, lastModified, ok = tr.Version(channel)
		return ok
	}, 3*time.Second, 5*time.Millisecond, "no version seeded for channel %q", channel)
	return etag, lastModified
}

func TestFingerprint_StableAndSensitive(t *testing.T) {
	sel := mustSelector(t, "flagSetId=fs1")
	a := testFlag("fs1", "a", "on")
	b := testFlag("fs1", "b", "off")

	// deterministic for identical input
	assert.Equal(t, fingerprint(&sel, []model.Flag{a, b}), fingerprint(&sel, []model.Flag{a, b}))
	// order independent
	assert.Equal(t, fingerprint(&sel, []model.Flag{a, b}), fingerprint(&sel, []model.Flag{b, a}))

	// sensitive to a definition change
	aModified := testFlag("fs1", "a", "off")
	assert.NotEqual(t, fingerprint(&sel, []model.Flag{a}), fingerprint(&sel, []model.Flag{aModified}))

	// duplicates are not collapsed: a flag appearing twice is a different result set
	assert.NotEqual(t, fingerprint(&sel, []model.Flag{a}), fingerprint(&sel, []model.Flag{a, a}))

	// the flagSetId carried on a flag is not part of the identity; the selector the watch was
	// opened with is
	relabelled := testFlag("some-other-flag-set", "a", "on")
	assert.Equal(t, fingerprint(&sel, []model.Flag{a}), fingerprint(&sel, []model.Flag{relabelled}))

	other := mustSelector(t, "source=src1")
	assert.NotEqual(t, fingerprint(&sel, []model.Flag{a}), fingerprint(&other, []model.Flag{a}))
}

// TestTracker_PassesSelectorToStore pins that filtering is delegated to the store by handing it
// the channel's selector, rather than re-implemented by grouping flags in memory.
func TestTracker_PassesSelectorToStore(t *testing.T) {
	rec := &recordingStore{}
	tr := newTracker(t, rec)

	subscribe(t, tr, "flagSetId=fs1")

	selectors := rec.selectors()
	require.Len(t, selectors, 1)
	expected := mustSelector(t, "flagSetId=fs1")
	assert.Equal(t, expected.ToLogString(), selectors[0].ToLogString())
}

func TestTracker_SharesOneWatchPerChannel(t *testing.T) {
	rec := &recordingStore{}
	tr := newTracker(t, rec)

	subscribe(t, tr, "flagSetId=fs1")
	subscribe(t, tr, "flagSetId=fs1")
	assert.Len(t, rec.selectors(), 1, "subscribers on one channel must share a store watch")

	subscribe(t, tr, "flagSetId=fs2")
	assert.Len(t, rec.selectors(), 2, "a distinct channel needs its own watch")
}

func TestTracker_Version_NoSubscription(t *testing.T) {
	tr := newTracker(t, newTestStore(t))

	_, _, ok := tr.Version("flagSetId=fs1")
	assert.False(t, ok, "no live stream for the channel means no ETag and no 304")
}

func TestTracker_ReleaseTearsDownLastSubscriberOnly(t *testing.T) {
	tr := newTracker(t, newTestStore(t))

	releaseA, err := tr.Subscribe("flagSetId=fs1", mustSelector(t, "flagSetId=fs1"))
	require.NoError(t, err)
	releaseB, err := tr.Subscribe("flagSetId=fs1", mustSelector(t, "flagSetId=fs1"))
	require.NoError(t, err)

	releaseA()
	assert.Len(t, tr.Channels(), 1, "one subscriber leaving must not tear down a shared channel")

	releaseB()
	assert.Empty(t, tr.Channels())
	_, _, ok := tr.Version("flagSetId=fs1")
	assert.False(t, ok)
}

func TestTracker_IdenticalUpdateIsSuppressed(t *testing.T) {
	s := newTestStore(t)
	flags := []model.Flag{testFlag("fs1", "a", "on")}
	s.Update("src1", flags, model.Metadata{"flagSetId": "fs1"}, false)

	tr := newTracker(t, s)
	subscribe(t, tr, "flagSetId=fs1")

	_, firstModified := versionOf(t, tr, "flagSetId=fs1")

	// Update re-inserts unconditionally, so an identical re-sync still wakes the watch
	s.Update("src1", flags, model.Metadata{"flagSetId": "fs1"}, false)
	time.Sleep(200 * time.Millisecond)

	_, stillModified := versionOf(t, tr, "flagSetId=fs1")
	assert.Equal(t, firstModified, stillModified, "lastModified must not move when nothing changed")
}

func TestTracker_VersionTracksChanges(t *testing.T) {
	s := newTestStore(t)
	s.Update("src1", []model.Flag{testFlag("fs1", "a", "on")}, model.Metadata{"flagSetId": "fs1"}, false)

	tr := newTracker(t, s)
	subscribe(t, tr, "flagSetId=fs1")

	before, _ := versionOf(t, tr, "flagSetId=fs1")

	s.Update("src1", []model.Flag{testFlag("fs1", "a", "off")}, model.Metadata{"flagSetId": "fs1"}, false)
	require.Eventually(t, func() bool {
		etag, _, _ := tr.Version("flagSetId=fs1")
		return etag != before
	}, 3*time.Second, 5*time.Millisecond, "the ETag must move when the config changes")
}

// closingStore closes the watcher without emitting, simulating store.Watch's error-close path.
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

func TestTracker_StoreErrorInvalidatesVersion(t *testing.T) {
	tr := newTracker(t, closingStore{})

	subscribe(t, tr, "flagSetId=fs1")

	assert.Eventually(t, func() bool {
		_, _, ok := tr.Version("flagSetId=fs1")
		return !ok
	}, 3*time.Second, 5*time.Millisecond,
		"a dead watch must stop serving its ETag, or 304s would be stale forever")
}

func TestTracker_CloseRejectsNewSubscribers(t *testing.T) {
	es := eventsource.NewServer()
	defer es.Close()
	tr := NewTracker(logger.NewLogger(nil, false), newTestStore(t), es)

	release, err := tr.Subscribe("flagSetId=fs1", mustSelector(t, "flagSetId=fs1"))
	require.NoError(t, err)

	require.NotPanics(t, tr.Close)
	assert.Empty(t, tr.Channels())

	_, err = tr.Subscribe("flagSetId=fs2", mustSelector(t, "flagSetId=fs2"))
	assert.Error(t, err)

	assert.NotPanics(t, release) // releasing after Close is a no-op
}

// recordingStore records the selector handed to each Watch call, mimicking the real store's
// contract: an immediate initial snapshot on a buffered channel, closed on ctx cancellation.
type recordingStore struct {
	mu   sync.Mutex
	sels []store.Selector
}

func (r *recordingStore) selectors() []store.Selector {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]store.Selector(nil), r.sels...)
}

func (r *recordingStore) Get(context.Context, string, *store.Selector) (model.Flag, model.Metadata, error) {
	return model.Flag{}, nil, nil
}

func (r *recordingStore) GetAll(context.Context, *store.Selector) ([]model.Flag, model.Metadata, error) {
	return nil, nil, nil
}

func (r *recordingStore) Watch(ctx context.Context, selector *store.Selector, watcher chan<- store.FlagQueryResult) {
	r.mu.Lock()
	r.sels = append(r.sels, *selector)
	r.mu.Unlock()

	go func() {
		watcher <- store.FlagQueryResult{}
		<-ctx.Done()
		close(watcher)
	}()
}

func (r *recordingStore) Update(string, []model.Flag, model.Metadata, bool) {}
