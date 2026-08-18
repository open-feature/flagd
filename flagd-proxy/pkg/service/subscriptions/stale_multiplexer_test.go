package subscriptions

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	isync "github.com/open-feature/flagd/core/pkg/sync"
)

// freshSyncBuilder hands out a new sync per call, as the real builder does. These tests
// start a second watcher, and a single shared mock would let both watchers compete for
// the same input channel.
type freshSyncBuilder struct {
	mu    sync.Mutex
	syncs []*syncMock
}

func (b *freshSyncBuilder) SyncsFromConfig(_ []isync.SourceConfig, _ *logger.Logger) ([]isync.ISync, error) {
	return nil, nil
}

func (b *freshSyncBuilder) SyncFromURI(_ string, _ *logger.Logger) (isync.ISync, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	s := newMockSync()
	b.syncs = append(b.syncs, s)
	return s, nil
}

// waitForSync blocks until the nth sync (1-based) has been handed out and returns it.
func (b *freshSyncBuilder) waitForSync(t *testing.T, n int) *syncMock {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		b.mu.Lock()
		if len(b.syncs) >= n {
			s := b.syncs[n-1]
			b.mu.Unlock()
			return s
		}
		b.mu.Unlock()
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("watcher %d was never started", n)
	return nil
}

// waitForData fails unless the subscriber receives a data sync before the deadline.
func waitForData(t *testing.T, ds <-chan isync.DataSync, msg string) {
	t.Helper()
	select {
	case <-ds:
	case <-time.After(3 * time.Second):
		t.Fatal(msg)
	}
}

// watchedTarget is a coordinator with one target already being watched, which is the
// starting point for the tests that then stop that watcher one way or another.
type watchedTarget struct {
	store   *Coordinator
	builder *freshSyncBuilder
	target  string
	sync    *syncMock
	errChan chan error
}

func newWatchedTarget(t *testing.T, ctx context.Context) *watchedTarget {
	t.Helper()
	store := NewManager(ctx, logger.NewLogger(nil, false))
	builder := &freshSyncBuilder{}
	store.syncBuilder = builder
	target := "ns/flags"

	dataChan := make(chan isync.DataSync, 1)
	errChan := make(chan error, 1)
	store.RegisterSubscription(ctx, target, "first", dataChan, errChan)

	// the watcher is running once it forwards data
	first := builder.waitForSync(t, 1)
	first.dataSyncChanIn <- isync.DataSync{FlagData: "initial"}
	waitForData(t, dataChan, "watcher never started")

	return &watchedTarget{store: store, builder: builder, target: target, sync: first, errChan: errChan}
}

// subscribeAgain registers a second subscription and returns its data channel, standing in
// for the client reconnecting after its stream broke.
func (w *watchedTarget) subscribeAgain(ctx context.Context) chan isync.DataSync {
	dataChan := make(chan isync.DataSync, 1)
	w.store.RegisterSubscription(ctx, w.target, "second", dataChan, make(chan error, 1))
	return dataChan
}

// Test_RegisterSubscription_afterWatcherStopped covers #2030: once watchResource has
// returned, its multiplexer is dead, so a subscription arriving afterwards must get a
// new watcher instead of silently attaching to the dead one and never receiving data.
func Test_RegisterSubscription_afterWatcherStopped(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w := newWatchedTarget(t, ctx)

	// the resource does not exist yet, so the sync fails and watchResource returns
	w.sync.errChanIn <- errors.New("resource not found")
	select {
	case <-w.errChan:
	case <-time.After(3 * time.Second):
		t.Fatal("sync error was never broadcast")
	}

	// the client retries while the dead multiplexer is still in the map
	secondData := w.subscribeAgain(ctx)

	// the resource now exists: a live watcher forwards this to its subscribers
	second := w.builder.waitForSync(t, 2)
	second.dataSyncChanIn <- isync.DataSync{FlagData: "after the resource was created"}
	waitForData(t, secondData, "subscription is wedged: no watcher was started for it")
}

// Test_RegisterSubscription_afterIdleShutdown covers the second way a multiplexer dies:
// the cleanup loop cancels multiplexers that have no subscriptions left. A subscription
// arriving after that cancellation must also get a new watcher.
func Test_RegisterSubscription_afterIdleShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w := newWatchedTarget(t, ctx)

	// what the cleanup loop does to an idle multiplexer
	w.store.mu.Lock()
	w.store.multiplexers[w.target].cancelFunc()
	w.store.mu.Unlock()

	secondData := w.subscribeAgain(ctx)

	second := w.builder.waitForSync(t, 2)
	second.dataSyncChanIn <- isync.DataSync{FlagData: "after shutdown"}
	waitForData(t, secondData, "subscription is wedged: no watcher was started after the idle shutdown")
}

// Test_multiplexerSubsGuardedConsistently covers a second race in the same area: subs is
// written under Coordinator.mu but read under multiplexer.mu, so a subscriber going away
// while the watcher broadcasts tears the map. Run under -race, which the Makefile uses.
func Test_multiplexerSubsGuardedConsistently(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	syncStore := NewManager(ctx, logger.NewLogger(nil, false))
	builder := &freshSyncBuilder{}
	syncStore.syncBuilder = builder
	target := "ns/flags"

	// a long-lived subscriber keeps the multiplexer and its watcher alive
	baseData := make(chan isync.DataSync, 1)
	baseErr := make(chan error, 1)
	syncStore.RegisterSubscription(ctx, target, "base", baseData, baseErr)
	syncSrc := builder.waitForSync(t, 1)

	// keep broadcasting while subscribers come and go
	done := make(chan struct{})
	var pushers sync.WaitGroup
	pushers.Add(1)
	go func() {
		defer pushers.Done()
		for {
			select {
			case <-done:
				return
			case syncSrc.dataSyncChanIn <- isync.DataSync{FlagData: "update"}:
			}
		}
	}()

	for i := 0; i < 50; i++ {
		subCtx, subCancel := context.WithCancel(ctx)
		syncStore.RegisterSubscription(subCtx, target, i, make(chan isync.DataSync, 1), make(chan error, 1))
		subCancel()
	}

	close(done)
	pushers.Wait()
}

// Test_watchResource_doesNotDeleteReplacement covers the delete being unconditional: a
// stopping watcher must not remove a multiplexer that a later subscription already
// rebuilt for the same target, which would strand that subscription in turn.
func Test_watchResource_doesNotDeleteReplacement(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	syncStore := NewManager(ctx, logger.NewLogger(nil, false))
	syncMock := newMockSync()
	syncStore.syncBuilder = &syncBuilderMock{mock: syncMock}
	target := "ns/flags"

	stopping, _ := newSyncHandler()
	replacement, _ := newSyncHandler()

	// a watcher is running for `stopping`, which is what the map holds
	syncStore.mu.Lock()
	syncStore.multiplexers[target] = stopping
	syncStore.mu.Unlock()
	go syncStore.watchResource(target)
	syncMock.dataSyncChanIn <- isync.DataSync{FlagData: "initial"}
	waitForData(t, stopping.subs["key"].dataSync, "watcher never started")

	// a replacement takes its place before the stopping watcher finishes
	syncStore.mu.Lock()
	syncStore.multiplexers[target] = replacement
	syncStore.mu.Unlock()

	syncMock.errChanIn <- errors.New("resource not found")

	// the stopping watcher removes its entry either in its defer or, before this fix, from a
	// goroutine woken later, so poll rather than sampling once
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		syncStore.mu.RLock()
		current, present := syncStore.multiplexers[target]
		syncStore.mu.RUnlock()
		if !present || current != replacement {
			t.Fatal("the stopping watcher removed the replacement multiplexer")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// stalledResyncSync models the real kubernetes sync: ReSync ends in an uncancellable send
// on an unbuffered channel, so a stalled subscriber parks the calling goroutine.
type stalledResyncSync struct {
	syncMock
	entered chan struct{}
	release chan struct{}
}

func (b *stalledResyncSync) ReSync(_ context.Context, _ chan<- isync.DataSync) error {
	close(b.entered)
	<-b.release
	return nil
}

// Test_watchResource_broadcastsErrorWhileResyncStalls keeps the sync error reachable when a
// ReSync is stalled: the broadcast must not sit behind a lock that ReSync can hold (#2030).
func Test_watchResource_broadcastsErrorWhileResyncStalls(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	syncStore := NewManager(ctx, logger.NewLogger(nil, false))
	stalled := &stalledResyncSync{
		syncMock: *newMockSync(),
		entered:  make(chan struct{}),
		release:  make(chan struct{}),
	}
	defer close(stalled.release)
	syncStore.syncBuilder = &syncBuilderMock{mock: stalled}
	target := "ns/flags"

	firstData := make(chan isync.DataSync, 1)
	firstErr := make(chan error, 1)
	syncStore.RegisterSubscription(ctx, target, "first", firstData, firstErr)

	stalled.dataSyncChanIn <- isync.DataSync{FlagData: "initial"}
	select {
	case <-firstData:
	case <-time.After(3 * time.Second):
		t.Fatal("watcher never started")
	}

	// a second subscriber triggers a ReSync that stalls, standing in for a stuck client
	syncStore.RegisterSubscription(ctx, target, "second", make(chan isync.DataSync, 1), make(chan error, 1))
	select {
	case <-stalled.entered:
	case <-time.After(3 * time.Second):
		t.Fatal("the resync never started")
	}

	// the sync now fails; the error must still be delivered
	stalled.errChanIn <- errors.New("resource not found")
	select {
	case <-firstErr:
	case <-time.After(3 * time.Second):
		t.Fatal("sync error never reached the subscriber; the coordinator lock is jammed")
	}
}
