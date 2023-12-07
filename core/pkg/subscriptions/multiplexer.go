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
	subs       map[interface{}]storedChannels
	dataSync   chan sourceSync.DataSync
	cancelFunc context.CancelFunc
	syncRef    sourceSync.ISync
	mu         *sync.RWMutex
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
