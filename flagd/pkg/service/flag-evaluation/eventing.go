package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	iservice "github.com/open-feature/flagd/core/pkg/service"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/utils"
)

// IEvents is an interface for event subscriptions
type IEvents interface {
	Subscribe(ctx context.Context, id any, selector *store.Selector, notifyChan chan iservice.Notification)
	Unsubscribe(id any)
	EmitToAll(n iservice.Notification)
}

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
	watcher := make(chan store.Payload, 1)
	go func() {
		// store the previous flags to compare against new notifications, to compute proper diffs for RPC mode
		oldFlags := map[string]model.Flag{}
		for payload := range watcher {
			newFlags := payload.Flags
			notifications := utils.BuildNotifications(oldFlags, newFlags)
			notifier <- iservice.Notification{
				Type: iservice.ConfigurationChange,
				Data: map[string]interface{}{
					"flags": notifications,
				},
			}
			oldFlags = newFlags
		}

		eventing.logger.Debug(fmt.Sprintf("closing notify channel for id %v", id))
		close(notifier)
	}()

	_, _, err := eventing.store.GetAll(ctx, selector, watcher)
	if err != nil {
		eventing.logger.Error(fmt.Sprintf("unable to subscribe to events for id %v: %v", id, err))
		close(notifier)
		return
	}
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
