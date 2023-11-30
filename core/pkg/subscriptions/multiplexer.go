package subscriptions

import (
	"context"
	"fmt"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	sync2 "sync"
)

// multiplexer distributes updates for a target to all of its subscribers
type multiplexer struct {
	subs       map[interface{}]storedChannels
	dataSync   chan sync.DataSync
	cancelFunc context.CancelFunc
	syncRef    sync.ISync
	mu         *sync2.RWMutex
}

func (h *multiplexer) writeError(logger *logger.Logger, err error) {
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

func (h *multiplexer) writeData(logger *logger.Logger, data sync.DataSync) {
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
