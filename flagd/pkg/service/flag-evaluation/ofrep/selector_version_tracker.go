package ofrep

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
)

// FlagStore is an interface for querying flags and watching for changes.
type FlagStore interface {
	GetAll(ctx context.Context, selector *store.Selector) ([]model.Flag, model.Metadata, error)
	WatchSelector(selector *store.Selector) <-chan struct{}
}

// trackedSelector holds the state for a single tracked selector
type trackedSelector struct {
	etag   string
	cancel context.CancelFunc
}

// SelectorVersionTracker tracks content hashes for selectors to enable ETag-based caching.
type SelectorVersionTracker struct {
	logger      *logger.Logger
	flagStore   FlagStore
	mu          sync.RWMutex
	selectors   map[string]*trackedSelector // single map for all per-selector state
	insertOrder []string                    // FIFO order for eviction
	maxCapacity int
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewSelectorVersionTracker creates a new version tracker with watch-based invalidation.
// maxCapacity limits the number of tracked selectors (0 = unlimited).
func NewSelectorVersionTracker(logger *logger.Logger, flagStore FlagStore, maxCapacity int) *SelectorVersionTracker {
	ctx, cancel := context.WithCancel(context.Background())
	return &SelectorVersionTracker{
		logger:      logger,
		flagStore:   flagStore,
		selectors:   make(map[string]*trackedSelector),
		insertOrder: make([]string, 0),
		maxCapacity: maxCapacity,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Close shuts down the tracker and stops all watch goroutines
func (t *SelectorVersionTracker) Close() {
	t.cancel()
}

// ETag returns the current ETag for a selector.
// returns empty string if the selector has never been tracked.
func (t *SelectorVersionTracker) ETag(selectorExpression string) string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if s, ok := t.selectors[selectorExpression]; ok {
		return s.etag
	}
	return ""
}

// Track starts tracking a selector and returns its current content-based ETag.
// if already tracking, returns the cached ETag without recomputing.
func (t *SelectorVersionTracker) Track(selectorExpression string) string {
	t.mu.Lock()
	defer t.mu.Unlock()

	// if already tracking, return cached ETag
	if s, exists := t.selectors[selectorExpression]; exists {
		return s.etag
	}

	// evict oldest if at capacity
	if t.maxCapacity > 0 && len(t.selectors) >= t.maxCapacity {
		t.evictOldest()
	}

	// compute content-based ETag
	etag := t.computeETag(selectorExpression)

	// start watching for changes
	var watchCancel context.CancelFunc
	if t.flagStore != nil {
		selector := store.NewSelector(selectorExpression)
		watchCh := t.flagStore.WatchSelector(&selector)
		var watchCtx context.Context
		watchCtx, watchCancel = context.WithCancel(t.ctx)
		go t.watchAndRecompute(watchCtx, selectorExpression, watchCh)
	}

	t.selectors[selectorExpression] = &trackedSelector{etag: etag, cancel: watchCancel}
	t.insertOrder = append(t.insertOrder, selectorExpression)

	t.logger.Debug(fmt.Sprintf("tracking selector '%s' with ETag %s", selectorExpression, etag))
	return etag
}

// computeETag generates a content-based ETag by hashing the flags for a selector.
// this ensures ETags are consistent across replicas with the same flag content.
func (t *SelectorVersionTracker) computeETag(selectorExpression string) string {
	if t.flagStore == nil {
		return ""
	}

	selector := store.NewSelector(selectorExpression)
	flags, metadata, err := t.flagStore.GetAll(t.ctx, &selector)
	if err != nil {
		t.logger.Warn(fmt.Sprintf("error getting flags for selector '%s': %v", selectorExpression, err))
		return ""
	}

	// create a hashable representation that includes the key
	// (model.Flag.MarshalJSON omits the key, so we need to include it explicitly)
	type flagForHash struct {
		Key            string
		State          string
		DefaultVariant string
		Variants       map[string]any
		Targeting      json.RawMessage
		Metadata       model.Metadata
	}
	hashableFlags := make([]flagForHash, len(flags))
	for i, f := range flags {
		hashableFlags[i] = flagForHash{
			Key:            f.Key,
			State:          f.State,
			DefaultVariant: f.DefaultVariant,
			Variants:       f.Variants,
			Targeting:      f.Targeting,
			Metadata:       f.Metadata,
		}
	}

	// serialize flags and metadata to create deterministic hash
	data, err := json.Marshal(struct {
		Flags    []flagForHash
		Metadata model.Metadata
	}{Flags: hashableFlags, Metadata: metadata})
	if err != nil {
		t.logger.Warn(fmt.Sprintf("error marshaling flags for selector '%s': %v", selectorExpression, err))
		return ""
	}

	hash := md5.Sum(data)
	return fmt.Sprintf("\"%x\"", hash)
}

// evictOldest removes the oldest tracked selector (FIFO).
// caller must hold the mutex.
func (t *SelectorVersionTracker) evictOldest() {
	if len(t.insertOrder) == 0 {
		return
	}
	oldest := t.insertOrder[0]
	t.insertOrder = t.insertOrder[1:]
	if s, ok := t.selectors[oldest]; ok {
		if s.cancel != nil {
			s.cancel()
		}
		delete(t.selectors, oldest)
	}
	t.logger.Warn(fmt.Sprintf("evicted selector '%s' from version tracker, consider increasing ofrep-cache-capacity", oldest))
}

// watchAndRecompute monitors a watch channel and recomputes the ETag when flags change
func (t *SelectorVersionTracker) watchAndRecompute(ctx context.Context, selectorExpression string, watchCh <-chan struct{}) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-watchCh:
			t.recomputeETag(selectorExpression)

			// re-establish watch for future changes
			if t.flagStore != nil {
				selector := store.NewSelector(selectorExpression)
				watchCh = t.flagStore.WatchSelector(&selector)
			} else {
				return
			}
		}
	}
}

// recomputeETag recomputes and updates the content-based ETag for a selector
func (t *SelectorVersionTracker) recomputeETag(selectorExpression string) {
	etag := t.computeETag(selectorExpression)

	t.mu.Lock()
	defer t.mu.Unlock()

	s, ok := t.selectors[selectorExpression]
	if !ok {
		return
	}
	s.etag = etag

	t.logger.Debug(fmt.Sprintf("recomputed ETag for selector '%s': %s", selectorExpression, etag))
}
