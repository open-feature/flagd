package store

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	isync "github.com/open-feature/flagd/core/pkg/sync"
)

type syncMock struct {
	isync.ISync
	dataSyncChanIn chan isync.DataSync
	errChanIn      chan error
	resyncData     isync.DataSync
	resyncError    error

	initError     error
	ctxCloseError error
}

func newMockSync() *syncMock {
	return &syncMock{
		dataSyncChanIn: make(chan isync.DataSync, 1),
		errChanIn:      make(chan error, 1),
	}
}

func (s *syncMock) Init(ctx context.Context) error {
	return s.initError
}

func (s *syncMock) Sync(ctx context.Context, dataSync chan<- isync.DataSync) error {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("here")
			return s.ctxCloseError
		case d := <-s.dataSyncChanIn:
			dataSync <- d
		case e := <-s.errChanIn:
			return e
		}
	}
}

func (s *syncMock) ReSync(ctx context.Context, dataSync chan<- isync.DataSync) error {
	dataSync <- s.resyncData
	return s.resyncError
}

type syncBuilderMock struct {
	mock      isync.ISync
	initError error
}

func (s *syncBuilderMock) SyncFromURI(uri string, logger logger.Logger) (isync.ISync, error) {
	return s.mock, s.initError
}

func newSyncHandler() (*syncHandler, interface{}) {
	coreDataSyncChan := make(chan isync.DataSync, 1)
	dataSyncChan := make(chan isync.DataSync, 1)
	errChan := make(chan error, 1)
	key := struct{}{}

	return &syncHandler{
		dataSync: coreDataSyncChan,
		subs: map[interface{}]storedChannels{
			key: {
				errChan:  errChan,
				dataSync: dataSyncChan,
			},
		},
	}, key
}

func Test_watchResource(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	syncStore := NewSyncStore(ctx, logger.NewLogger(nil, false))
	syncMock := newMockSync()
	syncStore.syncBuilder = &syncBuilderMock{
		mock: syncMock,
	}

	target := "test-target"
	syncHandler, key := newSyncHandler()

	syncStore.syncHandlers[target] = syncHandler

	go syncStore.watchResource(target)

	// sync update should be broadcasted to all registered sync subs:
	in := isync.DataSync{
		FlagData: "im a flag",
		Source:   "im a flag source",
		Type:     isync.ALL,
	}
	syncMock.dataSyncChanIn <- in

	select {
	case d := <-syncHandler.subs[key].dataSync:
		if !reflect.DeepEqual(d, in) {
			t.Error("unexpected sync data", in, d)
		}
	case <-time.After(3 * time.Second):
		t.Errorf("timed out waiting for broadcast of %v", in)
	}

	// errors should be broadcasted to all registered sync subs
	err := errors.New("very bad error")
	syncMock.errChanIn <- err
	select {
	case e := <-syncHandler.subs[key].errChan:
		if !errors.Is(e, err) {
			t.Error("unexpected sync error", e, err)
		}
	case <-time.After(3 * time.Second):
		t.Errorf("timed out waiting for broadcast of error")
	}

	// no context cancellation should have occurred, and there should still be registered sync sub
	syncStore.mu.Lock()
	if len(syncHandler.subs) != 1 {
		t.Error("incorrect number of subs in syncHandler", syncHandler.subs)
	}
	syncStore.mu.Unlock()

	// cancellation of context will result in the syncHandler being deleted
	cancel()
	// allow for the goroutine to catch the lock first
	time.Sleep(1 * time.Second)
	syncStore.mu.Lock()
	if syncStore.syncHandlers[target] != nil {
		t.Error("syncHandler has not been closed down after cancellation", syncHandler.subs)
	}
	syncStore.mu.Unlock()
}

func Test_watchResource_initFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	syncStore := NewSyncStore(ctx, logger.NewLogger(nil, false))
	syncMock := newMockSync()

	// return an error on startup
	syncMock.initError = errors.New("a terrible error")
	syncStore.syncBuilder = &syncBuilderMock{
		mock: syncMock,
	}

	target := "test-target"
	syncHandler, key := newSyncHandler()

	syncStore.syncHandlers[target] = syncHandler

	go syncStore.watchResource(target)

	// the error channel should immediately receive an error response and close
	select {
	case e := <-syncHandler.subs[key].errChan:
		if !errors.Is(e, syncMock.initError) {
			t.Error("unexpected sync error", e, syncMock.initError)
		}
	case <-time.After(3 * time.Second):
		t.Errorf("timed out waiting for broadcast of error")
	}

	// this should then close the internal context and the watcher should be removed
	time.Sleep(1 * time.Second)
	syncStore.mu.Lock()
	if syncStore.syncHandlers[target] != nil {
		t.Error("syncHandler has not been closed down after cancellation", syncHandler.subs)
	}
	syncStore.mu.Unlock()
}

func Test_watchResource_SyncFromURIFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	syncStore := NewSyncStore(ctx, logger.NewLogger(nil, false))
	syncMock := newMockSync()

	// return an error on startup
	syncBuilder := &syncBuilderMock{
		mock:      syncMock,
		initError: errors.New("a terrible error"),
	}
	syncStore.syncBuilder = syncBuilder

	target := "test-target"
	syncHandler, key := newSyncHandler()

	syncStore.syncHandlers[target] = syncHandler

	go syncStore.watchResource(target)

	// the error channel should immediately receive an error response and close
	select {
	case e := <-syncHandler.subs[key].errChan:
		if !errors.Is(e, syncBuilder.initError) {
			t.Error("unexpected sync error", e, syncBuilder.initError)
		}
	case <-time.After(3 * time.Second):
		t.Errorf("timed out waiting for broadcast of error")
	}

	// this should then close the internal context and the watcher should be removed
	time.Sleep(1 * time.Second)
	syncStore.mu.Lock()
	if syncStore.syncHandlers[target] != nil {
		t.Error("syncHandler has not been closed down after cancellation", syncHandler.subs)
	}
	syncStore.mu.Unlock()
}

func Test_watchResource_SyncErrorOnClose(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	syncStore := NewSyncStore(ctx, logger.NewLogger(nil, false))
	syncMock := newMockSync()

	// return an error on startup
	syncMock.ctxCloseError = errors.New("a terrible error")
	syncStore.syncBuilder = &syncBuilderMock{
		mock: syncMock,
	}

	target := "test-target"
	syncHandler, key := newSyncHandler()

	syncStore.syncHandlers[target] = syncHandler

	go syncStore.watchResource(target)
	cancel()
	// the error channel should immediately receive an error response and close
	select {
	case e := <-syncHandler.subs[key].errChan:
		if !errors.Is(e, syncMock.ctxCloseError) {
			t.Error("unexpected sync error", e, syncMock.initError)
		}
	case <-time.After(3 * time.Second):
		t.Errorf("timed out waiting for broadcast of error")
	}

	// this should then close the internal context and the watcher should be removed
	time.Sleep(1 * time.Second)
	syncStore.mu.Lock()
	if syncStore.syncHandlers[target] != nil {
		t.Error("syncHandler has not been closed down after cancellation", syncHandler.subs)
	}
	syncStore.mu.Unlock()
}
