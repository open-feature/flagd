package service

import (
	"sync"

	iservice "github.com/open-feature/flagd/core/pkg/service"
)

// eventingConfiguration is a wrapper for notification subscriptions
type eventingConfiguration struct {
	mu   *sync.RWMutex
	subs map[interface{}]chan iservice.Notification
}

func (eventing *eventingConfiguration) subscribe(id interface{}, notifyChan chan iservice.Notification) {
	eventing.mu.Lock()
	defer eventing.mu.Unlock()

	eventing.subs[id] = notifyChan
}

func (eventing *eventingConfiguration) emitToAll(n iservice.Notification) {
	eventing.mu.RLock()
	defer eventing.mu.RUnlock()

	for _, send := range eventing.subs {
		send <- n
	}
}

func (eventing *eventingConfiguration) unSubscribe(id interface{}) {
	eventing.mu.Lock()
	defer eventing.mu.Unlock()

	delete(eventing.subs, id)
}
