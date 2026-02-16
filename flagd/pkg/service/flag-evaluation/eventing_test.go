package service

import (
	"context"
	"sync"
	"testing"

	"github.com/open-feature/flagd/core/pkg/logger"
	iservice "github.com/open-feature/flagd/core/pkg/service"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestSubscribe(t *testing.T) {
	// given
	sources := []string{"source1", "source2"}
	log := logger.NewLogger(nil, false)
	s, err := store.NewStore(log, sources)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}

	eventing := &eventingConfiguration{
		subs:  make(map[interface{}]chan iservice.Notification),
		mu:    &sync.RWMutex{},
		store: s,
	}

	idA := "a"
	chanA := make(chan iservice.Notification, 1)

	idB := "b"
	chanB := make(chan iservice.Notification, 1)

	// when
	eventing.Subscribe(context.Background(), idA, nil, chanA)
	eventing.Subscribe(context.Background(), idB, nil, chanB)

	// then
	require.Equal(t, chanA, eventing.subs[idA], "incorrect subscription association")
	require.Equal(t, chanB, eventing.subs[idB], "incorrect subscription association")
}

func TestUnsubscribe(t *testing.T) {
	// given
	sources := []string{"source1", "source2"}
	log := logger.NewLogger(nil, false)
	s, err := store.NewStore(log, sources)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}
	eventing := &eventingConfiguration{
		subs:  make(map[interface{}]chan iservice.Notification),
		mu:    &sync.RWMutex{},
		store: s,
	}

	idA := "a"
	chanA := make(chan iservice.Notification, 1)
	idB := "b"
	chanB := make(chan iservice.Notification, 1)

	// when
	eventing.Subscribe(context.Background(), idA, nil, chanA)
	eventing.Subscribe(context.Background(), idB, nil, chanB)

	eventing.Unsubscribe(idA)

	// then
	require.Empty(t, eventing.subs[idA],
		"expected subscription cleared, but value present: %v", eventing.subs[idA])
	require.Equal(t, chanB, eventing.subs[idB], "incorrect subscription association")
}

func TestNotificationDataStructpbConversion(t *testing.T) {
	notification := iservice.Notification{
		Type: iservice.ConfigurationChange,
		Data: map[string]interface{}{
			"flags": map[string]interface{}{
				"flag1": map[string]interface{}{"state": "ENABLED"},
			},
		},
	}

	_, err := structpb.NewStruct(notification.Data)
	require.Nil(t, err)
}
