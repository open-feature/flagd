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
	// watcherCtx is the watchResource context, nil until it starts. Guarded by Coordinator.mu.
	watcherCtx context.Context
}

// isDead reports whether this multiplexer's watcher was cancelled and can no longer deliver;
// one that has not started yet is not dead. Callers must hold Coordinator.mu (#2030).
func (h *multiplexer) isDead() bool {
	return h.watcherCtx != nil && h.watcherCtx.Err() != nil
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
