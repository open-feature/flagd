package subscriptions

import (
	"context"
	"sync"

	"github.com/open-feature/flagd/core/pkg/logger"
	sourceSync "github.com/open-feature/flagd/core/pkg/sync"
	"go.uber.org/zap"
)

// multiplexer distributes updates for a target to all of its subscribers
type multiplexer struct {
	subs       map[any]storedChannels
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
			logger.Error("unable to write error to channel", zap.Any("key", k))
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
			logger.Error("unable to write data to channel", zap.Any("key", k))
		}
	}
}
