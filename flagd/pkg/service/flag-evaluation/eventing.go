package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/notifications"
	iservice "github.com/open-feature/flagd/core/pkg/service"
	"github.com/open-feature/flagd/core/pkg/store"
)

// IEvents is an interface for event subscriptions
type IEvents interface {
	Subscribe(ctx context.Context, id any, selector *store.Selector, notifyChan chan iservice.Notification)
	Unsubscribe(id any)
	EmitToAll(n iservice.Notification)
}

var _ IEvents = &eventingConfiguration{}

// eventingConfiguration is a wrapper for notification subscriptions
type eventingConfiguration struct {
	mu     *sync.RWMutex
	subs   map[any]chan iservice.Notification
	store  store.IStore
	logger *logger.Logger
}

func (eventing *eventingConfiguration) Subscribe(ctx context.Context, id any, selector *store.Selector, notifier chan iservice.Notification) {
	eventing.mu.Lock()
	defer eventing.mu.Unlock()

	// proxy events from our store watcher to the notify channel, so that RPC mode event streams
	watcher := make(chan store.FlagQueryResult, 1)
	go func() {
		// store the previous flags to compare against new notifications, to compute proper diffs for RPC mode
		var oldFlags map[string]model.Flag
		for result := range watcher {
			newFlags := make(map[string]model.Flag)
			for _, flag := range result.Flags {
                // we should be either selecting on a flag set here, or using the source-priority - duplicates are already handled, so we don't have to worry about overwrites
				newFlags[flag.Key] = flag
			}

			// ignore the first notification (nil old flags), the watcher emits on initialization, but for RPC we don't care until there's a change
			if oldFlags != nil {
				notifications := notifications.NewFromFlags(oldFlags, newFlags)
				notifier <- iservice.Notification{
					Type: iservice.ConfigurationChange,
					Data: map[string]interface{}{
						"flags": map[string]interface{}(notifications),
					},
				}
			}
			oldFlags = newFlags
		}

		eventing.logger.Debug(fmt.Sprintf("closing notify channel for id %v", id))
		close(notifier)
	}()

	eventing.store.Watch(ctx, selector, watcher)
	eventing.subs[id] = notifier
}

func (eventing *eventingConfiguration) EmitToAll(n iservice.Notification) {
	eventing.mu.RLock()
	defer eventing.mu.RUnlock()

	for _, send := range eventing.subs {
		send <- n
	}
}

func (eventing *eventingConfiguration) Unsubscribe(id any) {
	eventing.mu.Lock()
	defer eventing.mu.Unlock()

	delete(eventing.subs, id)
}
