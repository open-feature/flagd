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
		scheme string
		bucket string
		object string
		// initial and changed content in the object's native format, along
		// with the JSON we expect flagd to publish for each.
		initialContent   string
		initialConverted string
		changedContent   string
		changedConverted string
	}{
		"json file type": {
			scheme:           "xyz",
			bucket:           "b",
			object:           "flags.json",
			initialContent:   "{\"flags\":{}}",
			initialConverted: "{\"flags\":{}}",
			changedContent:   "{\"flags\":{\"a\":{}}}",
			changedConverted: "{\"flags\":{\"a\":{}}}",
		},
		"yaml file type": {
			scheme:           "xyz",
			bucket:           "b",
			object:           "flags.yaml",
			initialContent:   "flags: []",
			initialConverted: "{\"flags\":[]}",
			changedContent:   "flags: {a: {}}",
			changedConverted: "{\"flags\":{\"a\":{}}}",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockPoller := synctesting.NewMockPoller()

			blobMock := NewMockBlob(tt.scheme)
			blobSync := &Sync{
				Bucket:     tt.scheme + "://" + tt.bucket,
				Object:     tt.object,
				BlobURLMux: blobMock.URLMux(),
				Poller:     mockPoller,
				Logger:     logger.NewLogger(nil, false),
			}

			ctx := context.Background()
			dataSyncChan := make(chan sync.DataSync, 1)

			blobMock.AddObject(tt.initialContent)

			go func() {
				err := blobSync.Sync(ctx, dataSyncChan)
				if err != nil {
					log.Fatalf("Error start sync: %s", err.Error())
					return
				}
			}()

			data := <-dataSyncChan // initial sync
			if data.FlagData != tt.initialConverted {
				t.Errorf("expected content: %s, but received content: %s", tt.initialConverted, data.FlagData)
			}
			// A genuine content change is published.
			tickWithConfigChange(t, mockPoller, dataSyncChan, blobMock, tt.changedContent, tt.changedConverted)
			// An unchanged object (matching ETag) is skipped.
			tickWithoutConfigChange(t, mockPoller, dataSyncChan)
			// A new ETag/ModTime but identical bytes is skipped by the body hash backstop.
			tickWithMetadataOnlyChange(t, mockPoller, dataSyncChan, blobMock, tt.changedContent)
			// Back to the original content is a genuine change again.
			tickWithConfigChange(t, mockPoller, dataSyncChan, blobMock, tt.initialContent, tt.initialConverted)
			tickWithoutConfigChange(t, mockPoller, dataSyncChan)
		})
	}
}

// tickWithConfigChange writes new content (bumping ETag/ModTime) and expects the change to be published.
func tickWithConfigChange(t *testing.T, mockPoller *synctesting.MockPoller, dataSyncChan chan sync.DataSync, blobMock *MockBlob, newContent, expectedConverted string) {
	blobMock.AddObject(newContent)
	mockPoller.Tick()
	select {
	case data, ok := <-dataSyncChan:
		if ok {
			if data.FlagData != expectedConverted {
				t.Errorf("expected content: %s, but received content: %s", expectedConverted, data.FlagData)
			}
		} else {
			t.Errorf("data channel unexpectedly closed")
		}
	default:
		t.Errorf("data channel has no expected update")
	}
}

// tickWithoutConfigChange ticks without touching the object; the matching ETag means nothing is fetched or published.
func tickWithoutConfigChange(t *testing.T, mockPoller *synctesting.MockPoller, dataSyncChan chan sync.DataSync) {
	mockPoller.Tick()
	expectNoUpdate(t, dataSyncChan)
}

// tickWithMetadataOnlyChange rewrites the object with identical bytes, bumping ETag/ModTime;
// the body hash backstop must recognize the content is unchanged and skip publishing.
func tickWithMetadataOnlyChange(t *testing.T, mockPoller *synctesting.MockPoller, dataSyncChan chan sync.DataSync, blobMock *MockBlob, sameContent string) {
	blobMock.AddObject(sameContent)
	mockPoller.Tick()
	expectNoUpdate(t, dataSyncChan)
}

func expectNoUpdate(t *testing.T, dataSyncChan chan sync.DataSync) {
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

	blobMock := NewMockBlob(scheme)
	blobSync := &Sync{
		Bucket:     scheme + "://" + bucket,
		Object:     object,
		BlobURLMux: blobMock.URLMux(),
		Poller:     mockPoller,
		Logger:     logger.NewLogger(nil, false),
	}

	ctx := context.Background()
	dataSyncChan := make(chan sync.DataSync, 1)

	config := "my-config"
	blobMock.AddObject(config)

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

// runSync performs a single change-detecting sync and reports whether the full
// object was fetched and whether anything was published.
func runSync(t *testing.T, s *Sync, blobMock *MockBlob, ch chan sync.DataSync) (fetched, published bool, data string) {
	t.Helper()
	before := blobMock.Reads()
	if err := s.sync(context.Background(), ch, false); err != nil {
		t.Fatalf("sync failed: %v", err)
	}
	fetched = blobMock.Reads() > before
	select {
	case d := <-ch:
		published = true
		data = d.FlagData
	default:
	}
	return fetched, published, data
}

func TestBlobSync_ChangeDetection(t *testing.T) {
	t0 := time.Now()
	t1 := t0.Add(time.Second)

	body1 := "{\"flags\":{\"a\":{}}}"
	body2 := "{\"flags\":{\"b\":{}}}"

	type state struct {
		etag    string
		modTime time.Time
		content string
	}

	tests := map[string]struct {
		base        state
		next        state
		wantFetch   bool
		wantPublish bool
		wantData    string
	}{
		"ETag unchanged -> no fetch, no publish": {
			base:        state{etag: "v1", modTime: t0, content: body1},
			next:        state{etag: "v1", modTime: t0, content: body1},
			wantFetch:   false,
			wantPublish: false,
		},
		"ETag changed, ModTime bumped, body identical -> fetch but no publish": {
			base:        state{etag: "v1", modTime: t0, content: body1},
			next:        state{etag: "v2", modTime: t1, content: body1},
			wantFetch:   true,
			wantPublish: false,
		},
		"ETag absent, ModTime unchanged -> no fetch, no publish": {
			base:        state{etag: "", modTime: t0, content: body1},
			next:        state{etag: "", modTime: t0, content: body1},
			wantFetch:   false,
			wantPublish: false,
		},
		"ETag absent, ModTime bumped, body identical -> fetch but no publish": {
			base:        state{etag: "", modTime: t0, content: body1},
			next:        state{etag: "", modTime: t1, content: body1},
			wantFetch:   true,
			wantPublish: false,
		},
		"ETag absent, ModTime bumped, body changed -> fetch and publish": {
			base:        state{etag: "", modTime: t0, content: body1},
			next:        state{etag: "", modTime: t1, content: body2},
			wantFetch:   true,
			wantPublish: true,
			wantData:    body2,
		},
		"ETag changed, body changed -> fetch and publish": {
			base:        state{etag: "v1", modTime: t0, content: body1},
			next:        state{etag: "v2", modTime: t1, content: body2},
			wantFetch:   true,
			wantPublish: true,
			wantData:    body2,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			blobMock := NewMockBlob("fake")
			blobMock.SetObject(tt.base.etag, tt.base.modTime, tt.base.content)
			s := &Sync{
				Bucket:     "fake://bucket",
				Object:     "flags.json",
				BlobURLMux: blobMock.URLMux(),
				Logger:     logger.NewLogger(nil, false),
			}
			ch := make(chan sync.DataSync, 1)

			// Baseline sync establishes the cached ETag/ModTime/body hash.
			fetched, published, data := runSync(t, s, blobMock, ch)
			if !fetched || !published || data != tt.base.content {
				t.Fatalf("baseline sync: expected fetch+publish of %q, got fetched=%v published=%v data=%q",
					tt.base.content, fetched, published, data)
			}

			// Apply the scenario's next observed state.
			blobMock.SetObject(tt.next.etag, tt.next.modTime, tt.next.content)

			fetched, published, data = runSync(t, s, blobMock, ch)
			if fetched != tt.wantFetch {
				t.Errorf("fetch: expected %v, got %v", tt.wantFetch, fetched)
			}
			if published != tt.wantPublish {
				t.Errorf("publish: expected %v, got %v (data=%q)", tt.wantPublish, published, data)
			}
			if tt.wantPublish && data != tt.wantData {
				t.Errorf("published data: expected %q, got %q", tt.wantData, data)
			}
		})
	}
}
