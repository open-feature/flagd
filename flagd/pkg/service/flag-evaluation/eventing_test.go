package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	iservice "github.com/open-feature/flagd/core/pkg/service"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

// newTestEventingConfig creates an eventingConfiguration backed by a test store.
func newTestEventingConfig(t *testing.T, sources []string) (*eventingConfiguration, store.IStore) {
	t.Helper()
	log := logger.NewLogger(nil, false)
	s, err := store.NewStore(log, sources)
	require.NoError(t, err)
	return &eventingConfiguration{
		subs:   make(map[interface{}]chan iservice.Notification),
		mu:     &sync.RWMutex{},
		store:  s,
		logger: log,
	}, s
}

func TestSubscribe(t *testing.T) {
	// given
	eventing, _ := newTestEventingConfig(t, []string{"source1", "source2"})

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
	eventing, _ := newTestEventingConfig(t, []string{"source1", "source2"})

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

// TestNotificationCompatibleWithStructpb verifies that notification data from
// flag change events can be converted to protobuf structs, as required by the
// EventStream handlers. This is a regression test for
// https://github.com/open-feature/flagd/discussions/1869
func TestNotificationCompatibleWithStructpb(t *testing.T) {
	sources := []string{"source1"}
	eventing, s := newTestEventingConfig(t, sources)

	notifyChan := make(chan iservice.Notification, 1)
	eventing.Subscribe(context.Background(), "test", nil, notifyChan)
	// allow the subscription goroutine to process the initial watch result
	time.Sleep(100 * time.Millisecond)

	// first update sets up oldFlags.
	s.Update(sources[0], []model.Flag{
		{Key: "flag1", DefaultVariant: "off"},
	}, model.Metadata{}, false)

	// second update triggers a ConfigurationChange with a real diff.
	s.Update(sources[0], []model.Flag{
		{Key: "flag1", DefaultVariant: "on"},
	}, model.Metadata{}, false)

	select {
	case n := <-notifyChan:
		require.Equal(t, iservice.ConfigurationChange, n.Type)
		// contains a named map type instead of plain map[string]interface{}.
		_, err := structpb.NewStruct(n.Data)
		require.NoError(t, err, "notification data must be compatible with structpb.NewStruct")
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for notification")
	}
}

// TestNoNotificationWhenFlagsUnchanged verifies that no ConfigurationChange
// notification is sent when a store update contains the same flags as before.
func TestNoNotificationWhenFlagsUnchanged(t *testing.T) {
	sources := []string{"source1"}
	eventing, s := newTestEventingConfig(t, sources)

	notifyChan := make(chan iservice.Notification, 1)
	eventing.Subscribe(context.Background(), "test", nil, notifyChan)
	time.Sleep(100 * time.Millisecond)

	// first update creates flag1 — this produces a notification (create).
	s.Update(sources[0], []model.Flag{
		{Key: "flag1", DefaultVariant: "off"},
	}, model.Metadata{}, false)

	// drain the first notification (flag creation).
	select {
	case <-notifyChan:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for first notification")
	}

	// second update with the same flags — should not produce a notification.
	s.Update(sources[0], []model.Flag{
		{Key: "flag1", DefaultVariant: "off"},
	}, model.Metadata{}, false)

	select {
	case n := <-notifyChan:
		t.Fatalf("unexpected notification received: %v", n)
	case <-time.After(500 * time.Millisecond):
		// expected: no notification sent
	}
}
