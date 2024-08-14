package blob

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	synctesting "github.com/open-feature/flagd/core/pkg/sync/testing"
	"go.uber.org/mock/gomock"
)

const (
	scheme = "xyz"
	bucket = "b"
	object = "o"
)

func TestSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCron := synctesting.NewMockCron(ctrl)
	mockCron.EXPECT().AddFunc(gomock.Any(), gomock.Any()).DoAndReturn(func(spec string, cmd func()) error {
		return nil
	})
	mockCron.EXPECT().Start().Times(1)

	blobSync := &Sync{
		Bucket: scheme + "://" + bucket,
		Object: object,
		Cron:   mockCron,
		Logger: logger.NewLogger(nil, false),
	}
	blobMock := NewMockBlob(scheme, func() *Sync {
		return blobSync
	})
	blobSync.BlobURLMux = blobMock.URLMux()

	ctx := context.Background()
	dataSyncChan := make(chan sync.DataSync, 1)

	config := "my-config"
	blobMock.AddObject(object, config)

	go func() {
		err := blobSync.Sync(ctx, dataSyncChan)
		if err != nil {
			log.Fatalf("Error start sync: %s", err.Error())
			return
		}
	}()

	data := <-dataSyncChan // initial sync
	if data.FlagData != config {
		t.Errorf("expected content: %s, but received content: %s", config, data.FlagData)
	}
	tickWithConfigChange(t, mockCron, dataSyncChan, blobMock, "new config")
	tickWithoutConfigChange(t, mockCron, dataSyncChan)
	tickWithConfigChange(t, mockCron, dataSyncChan, blobMock, "new config 2")
	tickWithoutConfigChange(t, mockCron, dataSyncChan)
	tickWithoutConfigChange(t, mockCron, dataSyncChan)
}

func tickWithConfigChange(t *testing.T, mockCron *synctesting.MockCron, dataSyncChan chan sync.DataSync, blobMock *MockBlob, newConfig string) {
	time.Sleep(1 * time.Millisecond) // sleep so the new file has different modification date
	blobMock.AddObject(object, newConfig)
	mockCron.Tick()
	select {
	case data, ok := <-dataSyncChan:
		if ok {
			if data.FlagData != newConfig {
				t.Errorf("expected content: %s, but received content: %s", newConfig, data.FlagData)
			}
		} else {
			t.Errorf("data channel unexpecdly closed")
		}
	default:
		t.Errorf("data channel has no expected update")
	}
}

func tickWithoutConfigChange(t *testing.T, mockCron *synctesting.MockCron, dataSyncChan chan sync.DataSync) {
	mockCron.Tick()
	select {
	case data, ok := <-dataSyncChan:
		if ok {
			t.Errorf("unexpected update: %s", data.FlagData)
		} else {
			t.Errorf("data channel unexpecdly closed")
		}
	default:
	}
}

func TestReSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCron := synctesting.NewMockCron(ctrl)

	blobSync := &Sync{
		Bucket: scheme + "://" + bucket,
		Object: object,
		Cron:   mockCron,
		Logger: logger.NewLogger(nil, false),
	}
	blobMock := NewMockBlob(scheme, func() *Sync {
		return blobSync
	})
	blobSync.BlobURLMux = blobMock.URLMux()

	ctx := context.Background()
	dataSyncChan := make(chan sync.DataSync, 1)

	config := "my-config"
	blobMock.AddObject(object, config)

	err := blobSync.ReSync(ctx, dataSyncChan)
	if err != nil {
		log.Fatalf("Error start sync: %s", err.Error())
		return
	}

	data := <-dataSyncChan
	if data.FlagData != config {
		t.Errorf("expected content: %s, but received content: %s", config, data.FlagData)
	}
}
