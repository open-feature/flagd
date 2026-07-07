package blob

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	"github.com/open-feature/flagd/core/pkg/sync/internal/bloburi"
	"github.com/open-feature/flagd/core/pkg/sync/internal/polling"
	"github.com/open-feature/flagd/core/pkg/utils"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/azureblob" // needed to initialize Azure Blob Storage driver
	_ "gocloud.dev/blob/gcsblob"   // needed to initialize GCS driver
	_ "gocloud.dev/blob/s3blob"    // needed to initialize s3 driver
	"golang.org/x/crypto/sha3"     //nolint:gosec
)

type Sync struct {
	Bucket      string
	Object      string
	BlobURLMux  *blob.URLMux
	Poller      polling.Poller
	Logger      *logger.Logger
	Interval    uint32
	ready       bool
	lastUpdated time.Time
	lastETag    string
	lastBodySHA string
}

func (hs *Sync) Init(_ context.Context) error {
	if hs.Bucket == "" {
		return errors.New("no bucket string set")
	}
	if hs.Object == "" {
		return errors.New("no object string set")
	}
	return nil
}

func (hs *Sync) IsReady() bool {
	return hs.ready
}

func (hs *Sync) Sync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	hs.Logger.Info(fmt.Sprintf("starting sync from %s/%s (interval: %ds)", hs.Bucket, hs.Object, hs.Interval))

	// Initial fetch
	hs.Logger.Debug(fmt.Sprintf("initial fetch from %s/%s", hs.Bucket, hs.Object))
	err := hs.sync(ctx, dataSync, false)
	if err != nil {
		return err
	}

	hs.ready = true

	hs.Logger.Debug(fmt.Sprintf("polling %s/%s every %ds (offset: %ds)",
		hs.Bucket, hs.Object, hs.Interval, hs.Poller.Offset()))

	hs.Poller.Start(ctx, func() {
		err := hs.sync(ctx, dataSync, false)
		if err != nil {
			hs.Logger.Warn(fmt.Sprintf("sync failed: %v", err))
		}
	})

	return nil
}

func (hs *Sync) ReSync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	return hs.sync(ctx, dataSync, true)
}

func (hs *Sync) sync(ctx context.Context, dataSync chan<- sync.DataSync, skipChangeDetection bool) error {
	bucket, err := hs.getBucket(ctx)
	if err != nil {
		return fmt.Errorf("couldn't get bucket: %v", err)
	}
	defer bucket.Close()

	var attrs *blob.Attributes
	if !skipChangeDetection {
		attrs, err = hs.fetchObjectAttributes(ctx, bucket)
		if err != nil {
			return fmt.Errorf("couldn't get object attributes: %v", err)
		}
		if hs.attributesUnchanged(attrs) {
			hs.Logger.Debug("configuration hasn't changed, skipping fetching full object")
			return nil
		}
	}

	msg, bodySHA, err := hs.fetchObject(ctx, bucket)
	if err != nil {
		return fmt.Errorf("couldn't get object: %v", err)
	}

	// only publish if the content actually differs from what we last saw.
	if !skipChangeDetection && bodySHA == hs.lastBodySHA {
		hs.Logger.Debug("configuration hasn't changed, skipping publishing")
		hs.updateState(attrs, bodySHA)
		return nil
	}

	hs.Logger.Debug(fmt.Sprintf("configuration updated: %s", msg))
	if !skipChangeDetection {
		hs.updateState(attrs, bodySHA)
	}
	dataSync <- sync.DataSync{FlagData: msg, Source: bloburi.Join(hs.Bucket, hs.Object)}
	return nil
}

// attributesUnchanged reports whether the object can be considered unchanged
// based on its attributes alone, allowing us to skip fetching the full object.
// It prefers the ETag and falls back to the modification time for stores that don't expose one.
func (hs *Sync) attributesUnchanged(attrs *blob.Attributes) bool {
	if attrs.ETag != "" {
		return hs.lastETag == attrs.ETag
	}
	if hs.lastUpdated.Equal(attrs.ModTime) {
		return true
	}
	if hs.lastUpdated.After(attrs.ModTime) {
		hs.Logger.Warn("configuration changed but the modification time decreased instead of increasing")
	}
	return false
}

// updateState records the attributes and body hash of the object we just
// observed so subsequent syncs can detect whether anything actually changed.
func (hs *Sync) updateState(attrs *blob.Attributes, bodySHA string) {
	if attrs != nil {
		hs.lastETag = attrs.ETag
		hs.lastUpdated = attrs.ModTime
	}
	hs.lastBodySHA = bodySHA
}

func (hs *Sync) getBucket(ctx context.Context) (*blob.Bucket, error) {
	b, err := hs.BlobURLMux.OpenBucket(ctx, hs.Bucket)
	if err != nil {
		return nil, fmt.Errorf("error opening bucket %s: %v", hs.Bucket, err)
	}
	return b, nil
}

func (hs *Sync) fetchObjectAttributes(ctx context.Context, bucket *blob.Bucket) (*blob.Attributes, error) {
	if hs.Object == "" {
		return nil, errors.New("no object string set")
	}
	attrs, err := bucket.Attributes(ctx, hs.Object)
	if err != nil {
		return nil, fmt.Errorf("error fetching attributes for object %s/%s: %w", hs.Bucket, hs.Object, err)
	}
	return attrs, nil
}

// fetchObject downloads the object and returns its JSON representation along
// with a SHA-3 hash of the raw bytes, used to detect no-op changes.
func (hs *Sync) fetchObject(ctx context.Context, bucket *blob.Bucket) (string, string, error) {
	r, err := bucket.NewReader(ctx, hs.Object, nil)
	if err != nil {
		return "", "", fmt.Errorf("error opening reader for object %s/%s: %w", hs.Bucket, hs.Object, err)
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return "", "", fmt.Errorf("error downloading object %s/%s: %w", hs.Bucket, hs.Object, err)
	}

	json, err := utils.ConvertToJSON(data, filepath.Ext(hs.Object), r.ContentType())
	if err != nil {
		return "", "", fmt.Errorf("error converting blob data to json: %w", err)
	}
	return json, hs.generateSha(data), nil
}

func (hs *Sync) generateSha(body []byte) string {
	hasher := sha3.New256()
	hasher.Write(body)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
