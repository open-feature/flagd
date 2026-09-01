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
	// subs is written holding both Coordinator.mu and mu, and read holding either
	subs       map[interface{}]storedChannels
	dataSync   chan sourceSync.DataSync
	cancelFunc context.CancelFunc
	syncRef    sourceSync.ISync
	mu         *sync.RWMutex
	// done is closed by kill once the watcher is cancelled. Nil until the watcher starts.
	// Written once by watchResource under Coordinator.mu; readers hold it too, except
	// watchResource's own defer, which is the writer.
	done    chan struct{}
	dieOnce sync.Once
}

// kill is the only way to cancel a watcher: mark first, so it never reads alive once cancelled.
// A watcher that has not started yet is skipped here and reaped on a later tick instead (#2030).
func (h *multiplexer) kill() {
	if h.done != nil {
		h.dieOnce.Do(func() { close(h.done) })
	}
	if h.cancelFunc != nil {
		h.cancelFunc()
	}
}

// isDead reports whether this multiplexer's watcher was cancelled and can no longer deliver;
// one that has not started yet is not dead. Callers must hold Coordinator.mu (#2030).
func (h *multiplexer) isDead() bool {
	select {
	case <-h.done: // a nil channel blocks, so a watcher that never started is not dead
		return true
	default:
		return false
	}
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
