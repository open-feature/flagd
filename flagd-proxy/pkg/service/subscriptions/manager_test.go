package subscriptions

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	isync "github.com/open-feature/flagd/core/pkg/sync"
)

type syncMock struct {
	isync.ISync
	dataSyncChanIn chan isync.DataSync
	errChanIn      chan error
	resyncData     *isync.DataSync
	resyncError    error

	initError     error
	ctxCloseError error

	mu sync.Mutex
}

func newMockSync() *syncMock {
	return &syncMock{
		dataSyncChanIn: make(chan isync.DataSync, 1),
		errChanIn:      make(chan error, 1),
	}
}

func (s *syncMock) Init(_ context.Context) error {
	return s.initError
}

func (s *syncMock) Sync(ctx context.Context, dataSync chan<- isync.DataSync) error {
	for {
		select {
		case <-ctx.Done():
			s.mu.Lock()
			defer s.mu.Unlock()
			return s.ctxCloseError
		case d := <-s.dataSyncChanIn:
			dataSync <- d
		case e := <-s.errChanIn:
			return e
		}
	}
}

func (s *syncMock) ReSync(_ context.Context, dataSync chan<- isync.DataSync) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.resyncData != nil {
		dataSync <- *s.resyncData
	}
	return s.resyncError
}

type syncBuilderMock struct {
	mock      isync.ISync
	initError error
}

func (s *syncBuilderMock) SyncsFromConfig(_ []isync.SourceConfig, _ *logger.Logger) ([]isync.ISync, error) {
	return nil, nil
}

func (s *syncBuilderMock) SyncFromURI(_ string, _ *logger.Logger) (isync.ISync, error) {
	return s.mock, s.initError
}

func newSyncHandler() (*multiplexer, string) {
	coreDataSyncChan := make(chan isync.DataSync, 1)
	dataSyncChan := make(chan isync.DataSync, 1)
	errChan := make(chan error, 1)
	key := "key"

	return &multiplexer{
		dataSync: coreDataSyncChan,
		subs: map[interface{}]storedChannels{
			key: {
				errChan:  errChan,
				dataSync: dataSyncChan,
			},
		},
		mu: &sync.RWMutex{},
	}, key
}

func Test_watchResource(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	syncStore := NewManager(ctx, logger.NewLogger(nil, false))
	syncMock := newMockSync()
	syncStore.syncBuilder = &syncBuilderMock{
		mock: syncMock,
	}

	target := "test-target"
	syncHandler, key := newSyncHandler()

	syncStore.multiplexers[target] = syncHandler

	go syncStore.watchResource(target)

	// sync update should be broadcasted to all registered sync subs:
	in := isync.DataSync{
		FlagData: "im a flag",
		Source:   "im a flag source",
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
		t.Error("incorrect number of subs in multiplexer", syncHandler.subs)
	}
	syncStore.mu.Unlock()

	// cancellation of context will result in the multiplexer being deleted
	cancel()
	// allow for the goroutine to catch the lock first
	time.Sleep(1 * time.Second)
	syncStore.mu.Lock()
	if syncStore.multiplexers[target] != nil {
		t.Error("multiplexer has not been closed down after cancellation", syncHandler.subs)
	}
	syncStore.mu.Unlock()
}

func Test_watchResource_initFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	syncStore := NewManager(ctx, logger.NewLogger(nil, false))
	syncMock := newMockSync()

	// return an error on startup
	syncMock.initError = errors.New("a terrible error")
	syncStore.syncBuilder = &syncBuilderMock{
		mock: syncMock,
	}

	target := "test-target"
	syncHandler, key := newSyncHandler()

	syncStore.multiplexers[target] = syncHandler

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
	if syncStore.multiplexers[target] != nil {
		t.Error("multiplexer has not been closed down after cancellation", syncHandler.subs)
	}
	syncStore.mu.Unlock()
}

func Test_watchResource_SyncFromURIFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	syncStore := NewManager(ctx, logger.NewLogger(nil, false))
	syncMock := newMockSync()

	// return an error on startup
	syncBuilder := &syncBuilderMock{
		mock:      syncMock,
		initError: errors.New("a terrible error"),
	}
	syncStore.syncBuilder = syncBuilder

	target := "test-target"
	syncHandler, key := newSyncHandler()

	syncStore.multiplexers[target] = syncHandler

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
	if syncStore.multiplexers[target] != nil {
		t.Error("multiplexer has not been closed down after cancellation", syncHandler.subs)
	}
	syncStore.mu.Unlock()
}

func Test_watchResource_SyncErrorOnClose(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	syncStore := NewManager(ctx, logger.NewLogger(nil, false))
	syncMock := newMockSync()

	// return an error on startup
	syncMock.ctxCloseError = errors.New("a terrible error")
	syncStore.syncBuilder = &syncBuilderMock{
		mock: syncMock,
	}

	target := "test-target"
	syncHandler, key := newSyncHandler()

	syncStore.multiplexers[target] = syncHandler

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
	if syncStore.multiplexers[target] != nil {
		t.Error("multiplexer has not been closed down after cancellation", syncHandler.subs)
	}
	syncStore.mu.Unlock()
}

func Test_watchResource_SyncHandlerDoesNotExist(_ *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	syncStore := NewManager(ctx, logger.NewLogger(nil, false))
	syncMock := newMockSync()

	// return an error on startup
	syncMock.ctxCloseError = errors.New("a terrible error")
	syncStore.syncBuilder = &syncBuilderMock{
		mock: syncMock,
	}

	target := "test-target"

	// sync store will early return and not block
	syncStore.watchResource(target)
}

func Test_watchResource_Cleanup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	syncStore := NewManager(ctx, logger.NewLogger(nil, false))
	syncMock := newMockSync()

	// return an error on startup
	syncMock.ctxCloseError = errors.New("a terrible error")
	syncStore.syncBuilder = &syncBuilderMock{
		mock: syncMock,
	}

	target := "test-target"

	syncHandler, _ := newSyncHandler()
	syncHandler.subs = map[interface{}]storedChannels{}
	doneChan := make(chan struct{}, 1)
	syncHandler.cancelFunc = func() {
		doneChan <- struct{}{}
	}
	syncStore.mu.Lock()
	syncStore.multiplexers[target] = syncHandler
	syncStore.mu.Unlock()
	go func() {
		syncStore.cleanup()
	}()

	select {
	case <-doneChan:
		return
	case <-time.After(10 * time.Second):
		t.Error("multiplexers not being cleaned up, timed out after 10 seconds")
	}
}

func Test_FetchAllFlags(t *testing.T) {
	tests := map[string]struct {
		expectErr  bool
		mockData   *isync.DataSync
		mockError  error
		setMock    bool
		setHandler bool
	}{
		"resync route": {
			expectErr: false,
			setMock:   true,
			mockData: &isync.DataSync{
				FlagData: "im a flag",
				Source:   "im a flag source",
			},
			setHandler: true,
		},
		"resync route sync does not exist": {
			expectErr:  true,
			setMock:    false,
			mockData:   nil,
			setHandler: true,
		},
		"resync route returns error": {
			expectErr:  true,
			setMock:    true,
			mockData:   nil,
			mockError:  errors.New("disaster"),
			setHandler: true,
		},
		"register subscription route timeout": {
			expectErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			syncStore := NewManager(ctx, logger.NewLogger(nil, false))
			syncMock := newMockSync()
			syncMock.resyncData = tt.mockData
			syncMock.resyncError = tt.mockError
			syncStore.syncBuilder = &syncBuilderMock{
				mock: syncMock,
			}

			target := "test-target"
			syncHandler, key := newSyncHandler()
			if tt.setMock {
				syncHandler.syncRef = syncMock
			}
			if tt.setHandler {
				syncStore.multiplexers[target] = syncHandler
			}

			data, err := syncStore.FetchAllFlags(ctx, key, target)
			if err != nil && !tt.expectErr {
				t.Error(err)
			}
			if err == nil && tt.expectErr {
				t.Error("did not receive expected error")
			}
			if tt.mockData != nil && !reflect.DeepEqual(data.FlagData, tt.mockData.FlagData) {
				t.Error("data does not match expected value", tt.mockData.FlagData, data.FlagData)
			}
		})
	}
}

func Test_registerSubscriptionResyncPath(t *testing.T) {
	tests := map[string]struct {
		data      *isync.DataSync
		err       error
		expectErr bool
	}{
		"happy path": {
			data: &isync.DataSync{
				FlagData: "im a flag",
				Source:   "im a flag source",
			},
			expectErr: false,
		},
		"resync fails": {
			err:       errors.New("disaster"),
			expectErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			syncStore := NewManager(ctx, logger.NewLogger(nil, false))

			syncMock := newMockSync()
			syncMock.resyncData = tt.data
			syncMock.resyncError = tt.err

			syncStore.syncBuilder = &syncBuilderMock{
				mock: syncMock,
			}

			target := "test-target"
			syncHandler, _ := newSyncHandler()
			syncHandler.syncRef = syncMock
			key := struct{}{}
			syncStore.multiplexers[target] = syncHandler
			dataChan := make(chan isync.DataSync, 1)
			errChan := make(chan error, 1)

			go syncStore.RegisterSubscription(ctx, target, key, dataChan, errChan)

			select {
			case d := <-dataChan:
				if !reflect.DeepEqual(d, *tt.data) {
					t.Error("received unexpected data", d, *tt.data)
				}
			case err := <-errChan:
				if !tt.expectErr {
					t.Error(err)
				}
			case <-time.After(3 * time.Second):
				t.Error("timed out waiting for data chan")
			}
		})
	}
}

func Test_syncMetrics(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	syncStore := NewManager(ctx, logger.NewLogger(nil, false))
	syncMock := newMockSync()
	syncStore.syncBuilder = &syncBuilderMock{
		mock: syncMock,
	}

	subs := syncStore.GetActiveSubscriptionsInt64()
	if subs != 0 {
		t.Error("there are no subscriptions registered, active subs should be 0")
	}

	target := "test-target"
	syncHandler, _ := newSyncHandler()

	syncStore.multiplexers[target] = syncHandler

	subs = syncStore.GetActiveSubscriptionsInt64()
	if subs != 1 {
		t.Error("active subs metric should equal 1")
	}
}
