package blob

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	synctesting "github.com/open-feature/flagd/core/pkg/sync/testing"
)

func TestBlobSync(t *testing.T) {
	tests := map[string]struct {
		scheme           string
		bucket           string
		object           string
		content          string
		convertedContent string
	}{
		"json file type": {
			scheme:           "xyz",
			bucket:           "b",
			object:           "flags.json",
			content:          "{\"flags\":{}}",
			convertedContent: "{\"flags\":{}}",
		},
		"yaml file type": {
			scheme:           "xyz",
			bucket:           "b",
			object:           "flags.yaml",
			content:          "flags: []",
			convertedContent: "{\"flags\":[]}",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockPoller := synctesting.NewMockPoller()

			blobSync := &Sync{
				Bucket: tt.scheme + "://" + tt.bucket,
				Object: tt.object,
				Poller: mockPoller,
				Logger: logger.NewLogger(nil, false),
			}
			blobMock := NewMockBlob(tt.scheme, func() *Sync {
				return blobSync
			})
			blobSync.BlobURLMux = blobMock.URLMux()

			ctx := context.Background()
			dataSyncChan := make(chan sync.DataSync, 1)

			blobMock.AddObject(tt.object, tt.content)

			go func() {
				err := blobSync.Sync(ctx, dataSyncChan)
				if err != nil {
					log.Fatalf("Error start sync: %s", err.Error())
					return
				}
			}()

			data := <-dataSyncChan // initial sync
			if data.FlagData != tt.convertedContent {
				t.Errorf("expected content: %s, but received content: %s", tt.convertedContent, data.FlagData)
			}
			tickWithConfigChange(t, mockPoller, dataSyncChan, blobMock, tt.object, tt.convertedContent)
			tickWithoutConfigChange(t, mockPoller, dataSyncChan)
			tickWithConfigChange(t, mockPoller, dataSyncChan, blobMock, tt.object, tt.convertedContent)
			tickWithoutConfigChange(t, mockPoller, dataSyncChan)
			tickWithoutConfigChange(t, mockPoller, dataSyncChan)
		})
	}
}

func tickWithConfigChange(t *testing.T, mockPoller *synctesting.MockPoller, dataSyncChan chan sync.DataSync, blobMock *MockBlob, object string, newConfig string) {
	time.Sleep(1 * time.Millisecond) // sleep so the new file has different modification date
	blobMock.AddObject(object, newConfig)
	mockPoller.Tick()
	select {
	case data, ok := <-dataSyncChan:
		if ok {
			if data.FlagData != newConfig {
				t.Errorf("expected content: %s, but received content: %s", newConfig, data.FlagData)
			}
		} else {
			t.Errorf("data channel unexpectedly closed")
		}
	default:
		t.Errorf("data channel has no expected update")
	}
}

func tickWithoutConfigChange(t *testing.T, mockPoller *synctesting.MockPoller, dataSyncChan chan sync.DataSync) {
	mockPoller.Tick()
	select {
	case data, ok := <-dataSyncChan:
		if ok {
			t.Errorf("unexpected update: %s", data.FlagData)
		} else {
			t.Errorf("data channel unexpectedly closed")
		}
	default:
	}
}

func TestReSync(t *testing.T) {
	const (
		scheme = "xyz"
		bucket = "b"
		object = "flags.json"
	)
	mockPoller := synctesting.NewMockPoller()

	blobSync := &Sync{
		Bucket: scheme + "://" + bucket,
		Object: object,
		Poller: mockPoller,
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
