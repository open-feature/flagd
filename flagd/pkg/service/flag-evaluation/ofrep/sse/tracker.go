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
	"go.uber.org/zap"
)

type version struct {
	configEtag   string
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

	wg sync.WaitGroup
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
		fp := fingerprint(&sub.selector, res.Flags)
		prev := sub.ver.Load()
		if prev != nil && prev.configEtag == fp {
			// A no-op wakeup: an identical re-sync, or a coarser radix watch channel firing
			// for a change outside this selector.
			continue
		}
		sub.ver.CompareAndSwap(prev, &version{configEtag: fp, lastModified: time.Now().Unix()})

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
		t.logger.Error("ofrep sse watch ended unexpectedly; its cached ETag is invalidated", zap.String("channel", channel))
	}
}

// publish must be called only after the new version is stored, so a client that refetches the
// instant it receives the event cannot read the previous ETag and be served a stale 304.
func (t *Tracker) publish(channel string, sub *subscription) {
	v := sub.ver.Load()
	if v == nil {
		return
	}

	id := strconv.FormatInt(time.Now().UnixMilli(), 10)
	t.es.Publish([]string{channel}, newRefetchEvent(id, v.configEtag, v.lastModified))
	if t.logger != nil {
		t.logger.Debug("published refetch event", zap.String("channel", channel), zap.String("etag", v.configEtag))
	}
}

// Version returns the current config ETag and last-modified time (unix seconds) for the channel
// matching the given selector expression. ok is false when no stream for it is live, so the
// caller skips conditional handling.
func (t *Tracker) Version(channel string) (configEtag string, lastModified int64, ok bool) {
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
	return v.configEtag, v.lastModified, true
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
	subs := t.subs
	t.subs = map[string]*subscription{}
	t.mu.Unlock()

	for _, sub := range subs {
		sub.cancel()
	}
	t.wg.Wait()
}

// fingerprint produces a deterministic hash of the result set a channel's watch resolves to.
func fingerprint(selector *store.Selector, flags []model.Flag) string {
	digests := make([]string, 0, len(flags))
	for _, f := range flags {
		h := sha256.New()
		h.Write([]byte(f.Key))
		h.Write([]byte{0})
		h.Write([]byte(f.Source))
		h.Write([]byte{0})
		b, _ := json.Marshal(f)
		h.Write(b)
		digests = append(digests, string(h.Sum(nil)))
	}
	sort.Strings(digests)

	h := sha256.New()
	// the selector path scopes the hash to this watch, so two channels resolving to the same
	// flags never share an ETag unless they come from the same selector source.
	h.Write([]byte(selector.ToLogString()))
	h.Write([]byte{0})
	for _, d := range digests {
		h.Write([]byte(d))
	}
	return hex.EncodeToString(h.Sum(nil))
}
