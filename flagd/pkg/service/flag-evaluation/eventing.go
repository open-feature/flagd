package service

import (
	"sync"

	iservice "github.com/open-feature/flagd/core/pkg/service"
)

// IEvents is an interface for event subscriptions
type IEvents interface {
	Subscribe(id any, notifyChan chan iservice.Notification)
	Unsubscribe(id any)
	EmitToAll(n iservice.Notification)
}

// eventingConfiguration is a wrapper for notification subscriptions
type eventingConfiguration struct {
	mu   *sync.RWMutex
	subs map[any]chan iservice.Notification
}

func (eventing *eventingConfiguration) Subscribe(id any, notifyChan chan iservice.Notification) {
	eventing.mu.Lock()
	defer eventing.mu.Unlock()

	eventing.subs[id] = notifyChan
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
