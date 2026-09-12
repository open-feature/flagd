package subscriptions

import (
	"context"
	"fmt"
	"sync"

	"github.com/open-feature/flagd/core/pkg/logger"
	sourceSync "github.com/open-feature/flagd/core/pkg/sync"
)

// multiplexer distributes updates for a target to all of its subscribers
type multiplexer struct {
	// subs, done and cancelFunc are guarded by mu; use the methods below rather than taking it
	// directly, so a Coordinator.mu call site is never mistaken for a multiplexer one.
	// subs is additionally only ever written with Coordinator.mu held, which is what lets the
	// teardown broadcast to exactly the subscribers it saw under that lock
	subs       map[interface{}]storedChannels
	done       bool
	cancelFunc context.CancelFunc
	dataSync   chan sourceSync.DataSync
	// syncRef is written by watchResource and read by the resync paths, all under Coordinator.mu
	syncRef sourceSync.ISync
	mu      sync.RWMutex
}

// watchedBy records the watcher's cancel, which is also what makes the multiplexer killable.
func (h *multiplexer) watchedBy(cancel context.CancelFunc) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.cancelFunc = cancel
}

// markDeadLocked flags the multiplexer and hands back the cancel for the caller to run with mu
// released, or nil when its watcher has not started: marking that one would let
// RegisterSubscription replace the entry, and the pending watchResource would adopt the
// replacement and orphan its watcher (#2030). Callers hold mu.
func (h *multiplexer) markDeadLocked() context.CancelFunc {
	if h.cancelFunc == nil {
		return nil
	}
	h.done = true
	return h.cancelFunc
}

// kill is the only way to cancel a watcher: the mark precedes the cancel, which runs with mu
// released so no cancel callback executes under its own lock (#2030).
func (h *multiplexer) kill() {
	h.mu.Lock()
	cancel := h.markDeadLocked()
	h.mu.Unlock()
	if cancel != nil {
		cancel()
	}
}

// killIfIdle kills the multiplexer when no subscribers remain, deciding and marking under one
// hold of mu rather than counting and killing separately, and reports the count it saw (#2030).
func (h *multiplexer) killIfIdle() (int, bool) {
	h.mu.Lock()
	subs := len(h.subs)
	var cancel context.CancelFunc
	if subs == 0 {
		cancel = h.markDeadLocked()
	}
	h.mu.Unlock()
	if cancel == nil {
		return subs, false
	}
	cancel()
	return subs, true
}

// isDead reports whether this multiplexer's watcher was cancelled and can no longer deliver (#2030).
func (h *multiplexer) isDead() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.done
}

// addSub registers a subscriber's channels.
func (h *multiplexer) addSub(key interface{}, chans storedChannels) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.subs[key] = chans
}

// removeSub drops a subscriber that has gone away.
func (h *multiplexer) removeSub(key interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.subs, key)
}

// subCount reports how many subscribers are still attached.
func (h *multiplexer) subCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.subs)
}

func (h *multiplexer) broadcastError(logger *logger.Logger, err error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for k, ec := range h.subs {
		select {
		case ec.errChan <- err:
			continue
		default:
			logger.Error(fmt.Sprintf("unable to write error to channel for key %p", k))
		}
	}
}

func (h *multiplexer) broadcastData(logger *logger.Logger, data sourceSync.DataSync) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for k, ds := range h.subs {
		select {
		case ds.dataSync <- data:
			continue
		default:
			logger.Error(fmt.Sprintf("unable to write data to channel for key %p", k))
		}
	}
}
