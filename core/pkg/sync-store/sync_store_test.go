package sync_store

import (
	"context"
	"errors"
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
}

func (s *syncMock) Init(ctx context.Context) error {
	return nil
}

func (s *syncMock) Sync(ctx context.Context, dataSync chan<- isync.DataSync) error {
	for {
		select {
		case <-ctx.Done():
			return nil
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
	mock isync.ISync
}

func (s *syncBuilderMock) SyncFromURI(uri string, logger logger.Logger) (isync.ISync, error) {
	return s.mock, nil
}

func Test_watchResource(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	syncStore := NewSyncStore(ctx, logger.NewLogger(nil, false))

	syncMock := &syncMock{
		dataSyncChanIn: make(chan isync.DataSync, 1),
		errChanIn:      make(chan error, 1),
	}

	syncStore.SyncBuilder = &syncBuilderMock{
		mock: syncMock,
	}

	coreDataSyncChan := make(chan isync.DataSync, 1)
	dataSyncChan := make(chan isync.DataSync, 1)
	errChan := make(chan error, 1)
	key := struct{}{}
	target := "test-target"

	syncHandler := syncHandler{
		dataSync: coreDataSyncChan,
		subs: map[interface{}]storedChannels{
			key: {
				errChan:  errChan,
				dataSync: dataSyncChan,
			},
		},
	}
	syncStore.syncHandlers[target] = &syncHandler

	go syncStore.watchResource(ctx, target)

	// sync update should be broadcasted to all registered sync subs:
	in := isync.DataSync{
		FlagData: "im a flag",
		Source:   "im a flag source",
		Type:     isync.ALL,
	}
	syncMock.dataSyncChanIn <- in

	select {
	case d := <-dataSyncChan:
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
	case e := <-errChan:
		if !errors.Is(e, err) {
			t.Error("unexpected sync error", e, err)
		}
	case <-time.After(3 * time.Second):
		t.Errorf("timed out waiting for broadcast of error")
	}

	// no context cancellation should have ocurred, and there should still be registered sync sub
	syncStore.mu.Lock()
	if len(syncHandler.subs) != 1 {
		t.Error("incorrect number of subs in syncHandler", syncHandler.subs)
	}
	syncStore.mu.Unlock()

	// cancellation of context will result in the sub being deleted
	cancel()
	time.Sleep(3 * time.Second)
	syncStore.mu.Lock()
	if len(syncHandler.subs) != 0 {
		t.Error("incorrect number of subs in syncHandler after cancellation", syncHandler.subs)
	}
	syncStore.mu.Unlock()
}
