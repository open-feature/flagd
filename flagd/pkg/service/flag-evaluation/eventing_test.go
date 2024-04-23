package service

import (
	"sync"
	"testing"

	iservice "github.com/open-feature/flagd/core/pkg/service"
	"github.com/stretchr/testify/require"
)

func TestSubscribe(t *testing.T) {
	// given
	eventing := &eventingConfiguration{
		subs: make(map[interface{}]chan iservice.Notification),
		mu:   &sync.RWMutex{},
	}

	idA := "a"
	chanA := make(chan iservice.Notification, 1)

	idB := "b"
	chanB := make(chan iservice.Notification, 1)

	// when
	eventing.Subscribe(idA, chanA)
	eventing.Subscribe(idB, chanB)

	// then
	require.Equal(t, chanA, eventing.subs[idA], "incorrect subscription association")
	require.Equal(t, chanB, eventing.subs[idB], "incorrect subscription association")
}

func TestUnsubscribe(t *testing.T) {
	// given
	eventing := &eventingConfiguration{
		subs: make(map[interface{}]chan iservice.Notification),
		mu:   &sync.RWMutex{},
	}

	idA := "a"
	chanA := make(chan iservice.Notification, 1)
	idB := "b"
	chanB := make(chan iservice.Notification, 1)

	// when
	eventing.Subscribe(idA, chanA)
	eventing.Subscribe(idB, chanB)

	eventing.Unsubscribe(idA)

	// then
	require.Empty(t, eventing.subs[idA],
		"expected subscription cleared, but value present: %v", eventing.subs[idA])
	require.Equal(t, chanB, eventing.subs[idB], "incorrect subscription association")
}
