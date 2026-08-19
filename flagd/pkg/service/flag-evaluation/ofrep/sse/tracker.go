package sse

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/launchdarkly/eventsource"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
)

type version struct {
	etag         string
	lastModified int64
}

// subscription is a single store.Watch, shared by every client on the same channel. It is
// reachable through Tracker.subs only while refs > 0.
type subscription struct {
	selector store.Selector // immutable after construction: the store retains &selector
	cancel   context.CancelFunc

	// ver is nil until the first snapshot lands, and again if the watch dies.
	ver atomic.Pointer[version]

	refs int // guarded by Tracker.mu
}

type Tracker struct {
	logger *logger.Logger
	store  store.IStore
	es     *eventsource.Server

	mu     sync.Mutex
	subs   map[string]*subscription
	closed bool

	wg      sync.WaitGroup
	eventID atomic.Int64
}

func NewTracker(log *logger.Logger, s store.IStore, es *eventsource.Server) *Tracker {
	return &Tracker{
		logger: log,
		store:  s,
		es:     es,
		subs:   map[string]*subscription{},
	}
}

// Subscribe registers interest in a channel, starting a store watch for its selector if this is
// the first subscriber. The returned release must be called exactly once on disconnect.
func (t *Tracker) Subscribe(channel string, selector store.Selector) (func(), error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return nil, fmt.Errorf("ofrep sse tracker is shutting down")
	}

	sub, existing := t.subs[channel]
	if !existing {
		// Subscriptions are shared, so they outlive any single request and are torn down by
		// release or Close rather than by a request context.
		ctx, cancel := context.WithCancel(context.Background())
		sub = &subscription{selector: selector, cancel: cancel}
		t.subs[channel] = sub
		t.wg.Add(1)

		watcher := make(chan store.FlagQueryResult, 1)
		t.store.Watch(ctx, &sub.selector, watcher)
		go t.watch(ctx, channel, sub, watcher)
	}
	sub.refs++

	var once sync.Once
	return func() { once.Do(func() { t.release(channel, sub) }) }, nil
}

// release drops one reference, tearing the subscription down when the last one goes. The map
// delete happens under the same lock as the refcount check and before cancel, so a concurrent
// Subscribe cannot attach to a dying subscription.
func (t *Tracker) release(channel string, sub *subscription) {
	t.mu.Lock()
	sub.refs--
	last := sub.refs <= 0
	if last && t.subs[channel] == sub {
		delete(t.subs, channel)
	}
	t.mu.Unlock()

	if last {
		sub.cancel()
	}
}

func (t *Tracker) watch(ctx context.Context, channel string, sub *subscription, watcher <-chan store.FlagQueryResult) {
	defer t.wg.Done()

	first := true
	// This loop must run until the channel closes: store.Watch's send does not select on
	// ctx.Done(), so abandoning the channel would leak the store's goroutine.
	for res := range watcher {
		fp := fingerprint(res.Flags)
		if prev := sub.ver.Load(); prev != nil && prev.etag == fp {
			// A no-op wakeup: an identical re-sync, or a coarser radix watch channel firing
			// for a change outside this selector.
			continue
		}
		sub.ver.Store(&version{etag: fp, lastModified: time.Now().Unix()})

		if first {
			// seed only; the connecting client re-fetches unconditionally anyway (ADR-0008)
			first = false
			continue
		}
		t.publish(channel, sub)
	}

	if ctx.Err() != nil {
		return // our own teardown
	}

	// The store hit a selector/iterator error, so this channel will never fire again. Stop
	// serving its ETag, which would otherwise answer 304 from a frozen fingerprint forever.
	sub.ver.Store(nil)
	if t.logger != nil {
		t.logger.Error(fmt.Sprintf(
			"ofrep sse watch for channel %q ended unexpectedly; its cached ETag is invalidated", channel))
	}
}

// publish must be called only after the new version is stored, so a client that refetches the
// instant it receives the event cannot read the previous ETag and be served a stale 304.
func (t *Tracker) publish(channel string, sub *subscription) {
	v := sub.ver.Load()
	if v == nil {
		return
	}

	id := strconv.FormatInt(t.eventID.Add(1), 10)
	t.es.Publish([]string{channel}, newRefetchEvent(id, v.etag, v.lastModified))
	if t.logger != nil {
		t.logger.Debug(fmt.Sprintf("published refetch event to channel %q (etag=%s)", channel, v.etag))
	}
}

// Version returns the current config ETag and last-modified time (unix seconds) for the channel
// matching the given selector expression. ok is false when no stream for it is live, so the
// caller skips conditional handling.
func (t *Tracker) Version(channel string) (etag string, lastModified int64, ok bool) {
	t.mu.Lock()
	sub, exists := t.subs[channel]
	t.mu.Unlock()

	if !exists {
		return "", 0, false
	}
	v := sub.ver.Load()
	if v == nil {
		return "", 0, false
	}
	return v.etag, v.lastModified, true
}

// Channels returns the channels with at least one live subscriber.
func (t *Tracker) Channels() []string {
	t.mu.Lock()
	defer t.mu.Unlock()

	channels := make([]string, 0, len(t.subs))
	for channel := range t.subs {
		channels = append(channels, channel)
	}
	return channels
}

// Close stops accepting subscriptions and waits for every watch goroutine to exit. It must
// complete before the eventsource server is closed, since a publish after that blocks forever.
func (t *Tracker) Close() {
	t.mu.Lock()
	t.closed = true
	subs := make([]*subscription, 0, len(t.subs))
	for _, sub := range t.subs {
		subs = append(subs, sub)
	}
	t.subs = map[string]*subscription{}
	t.mu.Unlock()

	for _, sub := range subs {
		sub.cancel()
	}
	t.wg.Wait()
}

// fingerprint produces a deterministic, restart-stable hash of a group of flag definitions.
// nilFlagSetId is normalised to "" so identical config yields the same fingerprint across
// restarts.
func fingerprint(flags []model.Flag) string {
	sorted := make([]model.Flag, len(flags))
	copy(sorted, flags)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].FlagSetId != sorted[j].FlagSetId {
			return sorted[i].FlagSetId < sorted[j].FlagSetId
		}
		if sorted[i].Source != sorted[j].Source {
			return sorted[i].Source < sorted[j].Source
		}
		return sorted[i].Key < sorted[j].Key
	})

	nilID := store.NilFlagSetId()
	h := sha256.New()
	for _, f := range sorted {
		fsid := f.FlagSetId
		if fsid == nilID {
			fsid = ""
		}
		h.Write([]byte(fsid))
		h.Write([]byte{0})
		h.Write([]byte(f.Key))
		h.Write([]byte{0})
		// model.Flag.MarshalJSON is stable (map keys sorted) and covers the definition fields.
		b, _ := json.Marshal(f)
		h.Write(b)
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))
}
