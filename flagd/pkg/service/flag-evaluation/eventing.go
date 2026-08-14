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
	subs   map[any]*subscription
	store  store.IStore
	logger *logger.Logger
}

type subscription struct {
	mu       sync.Mutex
	notifier chan iservice.Notification
	done     chan struct{}
	stopOnce sync.Once
	closed   bool
}

func (s *subscription) send(n iservice.Notification) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}

	// if the subscription is already stopped, deliver nothing;
	select {
	case <-s.done:
		return
	default:
	}

	select {
	case s.notifier <- n:
	case <-s.done:
	}
}

func (s *subscription) close() {
	s.stop()

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.closed {
		close(s.notifier)
		s.closed = true
	}
}

func (s *subscription) stop() {
	s.stopOnce.Do(func() { close(s.done) })
}

func (eventing *eventingConfiguration) Subscribe(ctx context.Context, id any, selector *store.Selector, notifier chan iservice.Notification) {
	eventing.mu.Lock()
	defer eventing.mu.Unlock()
	sub := &subscription{notifier: notifier, done: make(chan struct{})}

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
				// if there are no changes, don't emit a notification
				if len(notifications) == 0 {
					oldFlags = newFlags
					continue
				}
				sub.send(iservice.Notification{
					Type: iservice.ConfigurationChange,
					Data: map[string]interface{}{
						// don't use our custom type or it cannot be serialized, convert to map
						"flags": map[string]interface{}(notifications),
					},
				})
			}
			oldFlags = newFlags
		}

		eventing.logger.Debug(fmt.Sprintf("closing notify channel for id %v", id))
		eventing.mu.Lock()
		if eventing.subs[id] == sub {
			delete(eventing.subs, id)
		}
		eventing.mu.Unlock()
		sub.close()
	}()

	eventing.store.Watch(ctx, selector, watcher)
	eventing.subs[id] = sub
}

func (eventing *eventingConfiguration) EmitToAll(n iservice.Notification) {
	eventing.mu.RLock()
	subs := make([]*subscription, 0, len(eventing.subs))
	for _, sub := range eventing.subs {
		subs = append(subs, sub)
	}
	eventing.mu.RUnlock()

	for _, sub := range subs {
		sub.send(n)
	}
}

func (eventing *eventingConfiguration) Unsubscribe(id any) {
	eventing.mu.Lock()
	sub := eventing.subs[id]
	delete(eventing.subs, id)
	eventing.mu.Unlock()

	if sub != nil {
		sub.stop()
	}
}
