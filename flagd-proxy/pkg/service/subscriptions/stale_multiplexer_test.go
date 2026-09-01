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
	return waitForNthSync(t, &b.mu, &b.syncs, n)
}

// waitForNthSync blocks until a builder has handed out n syncs and returns the nth (1-based).
func waitForNthSync[T any](t *testing.T, mu *sync.Mutex, syncs *[]T, n int) T {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		if len(*syncs) >= n {
			s := (*syncs)[n-1]
			mu.Unlock()
			return s
		}
		mu.Unlock()
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("watcher %d was never started", n)
	var zero T
	return zero
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

// Test_RegisterSubscription_afterWatcherStopped covers #2030 end to end: a client that
// reconnects after its sync failed must receive data again. It usually gets there via the
// removed map entry rather than isDead, so it does not by itself prove the dead-multiplexer
// path; Test_RegisterSubscription_whileStoppingWatcherBroadcasts pins that one down.
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
	w.store.multiplexers[w.target].kill()
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

// gatedSync blocks inside Sync after its context is cancelled until the test releases it.
// That holds open the window between the cleanup loop cancelling a multiplexer and
// watchResource's defer removing it from the map, which is where a reconnecting client
// would otherwise attach to a multiplexer that can never deliver (#2030).
type gatedSync struct {
	isync.ISync
	dataIn    chan isync.DataSync
	cancelled chan struct{}
	// release is nil on every sync but the first, which is the only one the test gates
	release chan struct{}
}

func (g *gatedSync) Init(_ context.Context) error { return nil }

func (g *gatedSync) ReSync(_ context.Context, _ chan<- isync.DataSync) error { return nil }

func (g *gatedSync) Sync(ctx context.Context, dataSync chan<- isync.DataSync) error {
	for {
		select {
		case <-ctx.Done():
			close(g.cancelled)
			if g.release != nil {
				<-g.release
			}
			return nil
		case d := <-g.dataIn:
			dataSync <- d
		}
	}
}

// gatedSyncBuilder gates only the first sync it hands out; the replacement runs normally.
type gatedSyncBuilder struct {
	mu      sync.Mutex
	syncs   []*gatedSync
	release chan struct{}
}

func (b *gatedSyncBuilder) SyncsFromConfig(_ []isync.SourceConfig, _ *logger.Logger) ([]isync.ISync, error) {
	return nil, nil
}

func (b *gatedSyncBuilder) SyncFromURI(_ string, _ *logger.Logger) (isync.ISync, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	g := &gatedSync{
		dataIn:    make(chan isync.DataSync, 1),
		cancelled: make(chan struct{}),
	}
	if len(b.syncs) == 0 {
		g.release = b.release
	}
	b.syncs = append(b.syncs, g)
	return g, nil
}

func (b *gatedSyncBuilder) waitForSync(t *testing.T, n int) *gatedSync {
	t.Helper()
	return waitForNthSync(t, &b.mu, &b.syncs, n)
}

// Test_RegisterSubscription_afterCleanupLoopCancelled drives the real cleanup loop rather
// than simulating it, so it covers the second way a multiplexer dies: every cancellation
// site has to mark the multiplexer dead, not just watchResource's own defer (#2030).
func Test_RegisterSubscription_afterCleanupLoopCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	builder := &gatedSyncBuilder{release: make(chan struct{})}
	defer close(builder.release)

	store := NewManager(ctx, logger.NewLogger(nil, false))
	store.syncBuilder = builder
	target := "ns/flags"

	subCtx, dropSubscriber := context.WithCancel(ctx)
	dataChan := make(chan isync.DataSync, 1)
	store.RegisterSubscription(subCtx, target, "first", dataChan, make(chan error, 1))

	first := builder.waitForSync(t, 1)
	first.dataIn <- isync.DataSync{FlagData: "initial"}
	waitForData(t, dataChan, "watcher never started")

	store.mu.RLock()
	sh := store.multiplexers[target]
	store.mu.RUnlock()

	// leave the multiplexer with no subscribers, which is what the cleanup loop acts on
	dropSubscriber()
	select {
	case <-first.cancelled:
	case <-time.After(20 * time.Second):
		t.Fatal("cleanup loop never cancelled the idle multiplexer")
	}

	// the gate keeps watchResource inside Sync, so the cancelled multiplexer is still mapped
	requireStillMapped(t, store, target, sh)
	secondData := make(chan isync.DataSync, 1)
	store.RegisterSubscription(ctx, target, "second", secondData, make(chan error, 1))

	second := builder.waitForSync(t, 2)
	second.dataIn <- isync.DataSync{FlagData: "after the cleanup loop ran"}
	waitForData(t, secondData, "subscription attached to a multiplexer the cleanup loop had cancelled")
}

// Test_RegisterSubscription_whileStoppingWatcherBroadcasts pins the second half of the
// invariant kill() carries: a stopping watcher parks in broadcastError before it can remove
// its own map entry, so the multiplexer has to already read as dead in that window (#2030).
func Test_RegisterSubscription_whileStoppingWatcherBroadcasts(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w := newWatchedTarget(t, ctx)

	w.store.mu.RLock()
	sh := w.store.multiplexers[w.target]
	w.store.mu.RUnlock()

	// park the stopping watcher inside broadcastError, which it reaches after kill() but
	// before it takes Coordinator.mu to delete itself. The lock is held on its own
	// goroutine: the test goroutine then only ever takes Coordinator.mu, so it cannot
	// invert the Coordinator.mu -> multiplexer.mu order and hang instead of failing
	release := parkMultiplexer(sh, 5*time.Second)
	defer close(release)
	w.sync.errChanIn <- errors.New("resource not found")

	deadline := time.Now().Add(3 * time.Second)
	for {
		w.store.mu.RLock()
		dead := sh.isDead()
		w.store.mu.RUnlock()
		if dead {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("stopping watcher never marked its multiplexer dead")
		}
		time.Sleep(5 * time.Millisecond)
	}

	// this only proves anything while the dead entry is still mapped: otherwise the rebuild
	// comes from a map miss and the isDead branch is never reached
	requireStillMapped(t, w.store, w.target, sh)

	// off the test goroutine, so a regression that attaches to the dead multiplexer blocks
	// there instead of taking the whole package down with the -timeout
	registered := make(chan chan isync.DataSync, 1)
	go func() { registered <- w.subscribeAgain(ctx) }()
	var secondData chan isync.DataSync
	select {
	case secondData = <-registered:
	case <-time.After(3 * time.Second):
		t.Fatal("RegisterSubscription attached to the dead multiplexer instead of rebuilding")
	}

	second := w.builder.waitForSync(t, 2)
	second.dataSyncChanIn <- isync.DataSync{FlagData: "after the watcher stopped"}
	waitForData(t, secondData, "subscription attached to a multiplexer whose watcher had stopped")
}

// parkMultiplexer holds sh.mu on a goroutine of its own until the returned channel is closed
// or the hold expires. The timeout is load-bearing: a caller stuck on Coordinator.mu cannot
// close the channel, so without it a lock-order regression deadlocks instead of failing.
func parkMultiplexer(sh *multiplexer, hold time.Duration) chan struct{} {
	release := make(chan struct{})
	parked := make(chan struct{})
	go func() {
		sh.mu.Lock()
		close(parked)
		select {
		case <-release:
		case <-time.After(hold):
		}
		sh.mu.Unlock()
	}()
	<-parked
	return release
}

// requireStillMapped fails unless target still resolves to sh, so the tests that depend on a
// dead-but-present multiplexer cannot quietly degrade into map-miss tests.
func requireStillMapped(t *testing.T, store *Coordinator, target string, sh *multiplexer) {
	t.Helper()
	store.mu.RLock()
	defer store.mu.RUnlock()
	if store.multiplexers[target] != sh {
		t.Fatal("the cancelled multiplexer was already removed, so this test proves nothing")
	}
}

// Test_kill_multiplexerWithoutWatcher covers the state RegisterSubscription leaves behind when
// a subscriber goes away before watchResource has taken Coordinator.mu to set cancelFunc: the
// cleanup loop reaps it on the next tick and must not panic on the nil cancelFunc or done.
func Test_kill_multiplexerWithoutWatcher(t *testing.T) {
	sh := &multiplexer{
		dataSync: make(chan isync.DataSync),
		subs:     map[interface{}]storedChannels{},
		mu:       &sync.RWMutex{},
	}

	sh.kill()

	if sh.isDead() {
		t.Fatal("a multiplexer whose watcher never started must not read as dead")
	}
}

// Test_kill_marksBeforeCancelling pins the ordering half of kill's contract: a watcher must
// never observe its own cancellation while the multiplexer still reads as alive (#2030).
func Test_kill_marksBeforeCancelling(t *testing.T) {
	sh := &multiplexer{mu: &sync.RWMutex{}, done: make(chan struct{})}
	sh.cancelFunc = func() {
		if !sh.isDead() {
			t.Error("watcher cancelled while the multiplexer still read as alive")
		}
	}

	sh.kill()
}
