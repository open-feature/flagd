package blob

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	"github.com/open-feature/flagd/core/pkg/utils"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/azureblob" // needed to initialize Azure Blob Storage driver
	_ "gocloud.dev/blob/gcsblob"   // needed to initialize GCS driver
	_ "gocloud.dev/blob/s3blob"    // needed to initialize s3 driver
)

type Sync struct {
	Bucket      string
	Object      string
	BlobURLMux  *blob.URLMux
	Cron        Cron
	Logger      *logger.Logger
	Interval    uint32
	ready       bool
	lastUpdated time.Time
}

// Cron defines the behaviour required of a cron
type Cron interface {
	AddFunc(spec string, cmd func()) error
	Start()
	Stop()
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
	hs.Logger.Info(fmt.Sprintf("starting sync from %s/%s with interval %ds", hs.Bucket, hs.Object, hs.Interval))
	_ = hs.Cron.AddFunc(fmt.Sprintf("*/%d * * * *", hs.Interval), func() {
		err := hs.sync(ctx, dataSync, false)
		if err != nil {
			hs.Logger.Warn(fmt.Sprintf("sync failed: %v", err))
		}
	})
	// Initial fetch
	hs.Logger.Debug(fmt.Sprintf("initial sync of the %s/%s", hs.Bucket, hs.Object))
	err := hs.sync(ctx, dataSync, false)
	if err != nil {
		return err
	}

	hs.ready = true
	hs.Cron.Start()
	<-ctx.Done()
	hs.Cron.Stop()

	return nil
}

func (hs *Sync) ReSync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	return hs.sync(ctx, dataSync, true)
}

func (hs *Sync) sync(ctx context.Context, dataSync chan<- sync.DataSync, skipCheckingModTime bool) error {
	bucket, err := hs.getBucket(ctx)
	if err != nil {
		return fmt.Errorf("couldn't get bucket: %v", err)
	}
	defer bucket.Close()
	var updated time.Time
	if !skipCheckingModTime {
		updated, err = hs.fetchObjectModificationTime(ctx, bucket)
		if err != nil {
			return fmt.Errorf("couldn't get object attributes: %v", err)
		}
		if hs.lastUpdated.Equal(updated) {
			hs.Logger.Debug("configuration hasn't changed, skipping fetching full object")
			return nil
		}
		if hs.lastUpdated.After(updated) {
			hs.Logger.Warn("configuration changed but the modification time decreased instead of increasing")
		}
	}
	msg, err := hs.fetchObject(ctx, bucket)
	if err != nil {
		return fmt.Errorf("couldn't get object: %v", err)
	}
	hs.Logger.Debug(fmt.Sprintf("configuration updated: %s", msg))
	if !skipCheckingModTime {
		hs.lastUpdated = updated
	}
	dataSync <- sync.DataSync{FlagData: msg, Source: hs.Bucket + hs.Object}
	return nil
}

func (hs *Sync) getBucket(ctx context.Context) (*blob.Bucket, error) {
	b, err := hs.BlobURLMux.OpenBucket(ctx, hs.Bucket)
	if err != nil {
		return nil, fmt.Errorf("error opening bucket %s: %v", hs.Bucket, err)
	}
	return b, nil
}

func (hs *Sync) fetchObjectModificationTime(ctx context.Context, bucket *blob.Bucket) (time.Time, error) {
	if hs.Object == "" {
		return time.Time{}, errors.New("no object string set")
	}
	attrs, err := bucket.Attributes(ctx, hs.Object)
	if err != nil {
		return time.Time{}, fmt.Errorf("error fetching attributes for object %s/%s: %w", hs.Bucket, hs.Object, err)
	}
	return attrs.ModTime, nil
}

func (hs *Sync) fetchObject(ctx context.Context, bucket *blob.Bucket) (string, error) {
	r, err := bucket.NewReader(ctx, hs.Object, nil)
	if err != nil {
		return "", fmt.Errorf("error opening reader for object %s/%s: %w", hs.Bucket, hs.Object, err)
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("error downloading object %s/%s: %w", hs.Bucket, hs.Object, err)
	}

	json, err := utils.ConvertToJSON(data, filepath.Ext(hs.Object), r.ContentType())
	if err != nil {
		return "", fmt.Errorf("error converting blob data to json: %w", err)
	}
	return json, nil
}
