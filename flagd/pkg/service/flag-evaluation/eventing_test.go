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
		subs:   make(map[interface{}]*subscription),
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
	require.Equal(t, chanA, eventing.subs[idA].notifier, "incorrect subscription association")
	require.Equal(t, chanB, eventing.subs[idB].notifier, "incorrect subscription association")
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
	require.Equal(t, chanB, eventing.subs[idB].notifier, "incorrect subscription association")
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

// TestEmitToAllVsSubscriberCancel checks EmitToAll racing subscriber
// cancellation (which closes notifiers) never sends on a closed channel.
func TestEmitToAllVsSubscriberCancel(t *testing.T) {
	// cancelling subscribers while hammering EmitToAll must not panic.
	eventing, _ := newTestEventingConfig(t, []string{"source1"})

	var drain, work sync.WaitGroup
	cancels := make([]context.CancelFunc, 0, 50)
	for i := 0; i < 50; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		cancels = append(cancels, cancel)
		ch := make(chan iservice.Notification, 1)
		eventing.Subscribe(ctx, i, nil, ch)
		drain.Add(1)
		go func() {
			defer drain.Done()
			for range ch { //nolint:revive // drain so EmitToAll never blocks on a full buffer
			}
		}()
	}
	time.Sleep(50 * time.Millisecond) // let subscription goroutines start watching

	work.Add(2)
	go func() {
		defer work.Done()
		for _, c := range cancels {
			c() // cancel -> store closes watcher -> goroutine closes notifier
		}
	}()
	go func() {
		defer work.Done()
		for i := 0; i < 200; i++ {
			// hammer EmitToAll into the teardown; panics before the fix
			eventing.EmitToAll(iservice.Notification{Type: iservice.Shutdown})
		}
	}()
	work.Wait()
	drain.Wait()
}

func TestBlockedSubscriberDoesNotBlockCleanup(t *testing.T) {
	eventing, _ := newTestEventingConfig(t, []string{"source1"})

	stalledCtx, cancelStalled := context.WithCancel(context.Background())
	stalledNotifier := make(chan iservice.Notification)
	eventing.Subscribe(stalledCtx, "stalled", nil, stalledNotifier)
	stalledSub := eventing.subs["stalled"]

	emitDone := make(chan struct{})
	go func() {
		defer close(emitDone)
		eventing.EmitToAll(iservice.Notification{Type: iservice.Shutdown})
	}()

	require.Eventually(t, func() bool {
		if stalledSub.mu.TryLock() {
			stalledSub.mu.Unlock()
			return false
		}
		return true
	}, time.Second, time.Millisecond, "EmitToAll did not block on the stalled subscriber")

	otherCtx, cancelOther := context.WithCancel(context.Background())
	otherNotifier := make(chan iservice.Notification, 1)
	eventing.Subscribe(otherCtx, "other", nil, otherNotifier)
	cancelOther()

	require.Eventually(t, func() bool {
		eventing.mu.RLock()
		defer eventing.mu.RUnlock()
		_, ok := eventing.subs["other"]
		return !ok
	}, time.Second, time.Millisecond, "other subscriber cleanup was blocked")

	select {
	case _, ok := <-otherNotifier:
		require.False(t, ok, "other subscriber notifier was not closed")
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for other subscriber notifier to close")
	}

	cancelStalled()
	select {
	case <-emitDone:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for EmitToAll to stop after cancellation")
	}
}

func TestNoDeliveryAfterUnsubscribe(t *testing.T) {
	// nothing should be delivered after Unsubscribe returns.
	// looped because the unfixed select only leaks ~50% of the time.
	for i := 0; i < 50; i++ {
		eventing, _ := newTestEventingConfig(t, []string{"source1"})
		ch := make(chan iservice.Notification, 1)
		eventing.Subscribe(context.Background(), "id", nil, ch)
		sub := eventing.subs["id"]

		sub.mu.Lock() // park EmitToAll's send after it snapshots this sub
		done := make(chan struct{})
		go func() {
			defer close(done)
			eventing.EmitToAll(iservice.Notification{Type: iservice.ConfigurationChange})
		}()
		time.Sleep(20 * time.Millisecond)
		eventing.Unsubscribe("id") // closes done while the send is parked
		sub.mu.Unlock()

		<-done
		require.Empty(t, ch, "notification delivered after Unsubscribe returned")
	}
}
