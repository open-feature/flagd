package blob

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/gcsblob" // needed to initialize GCS driver
	//nolint:gosec
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

func (hs *Sync) Init(ctx context.Context) error {
	return nil
}

func (hs *Sync) IsReady() bool {
	return hs.ready
}

func (hs *Sync) Sync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	hs.Logger.Info(fmt.Sprintf("starting sync from %s/%s with interval %d", hs.Bucket, hs.Object, hs.Interval))
	// Initial fetch
	hs.Logger.Debug(fmt.Sprintf("initial sync of the %s/%s", hs.Bucket, hs.Object))
	err := hs.ReSync(ctx, dataSync)
	if err != nil {
		return err
	}
	hs.ready = true

	hs.Logger.Debug(fmt.Sprintf("polling %s/%s every %d seconds", hs.Bucket, hs.Object, hs.Interval))
	_ = hs.Cron.AddFunc(fmt.Sprintf("*/%d * * * *", hs.Interval), func() {
		hs.Logger.Debug(fmt.Sprintf("fetching configuration from %s/%s", hs.Bucket, hs.Object))
		bucket, err := hs.getBucket(ctx)
		if err != nil {
			hs.Logger.Warn(fmt.Sprintf("couldn't get bucket: %v", err))
			return
		}
		defer bucket.Close()
		updated, err := hs.fetchObjectModificationTime(ctx, bucket)
		if err != nil {
			hs.Logger.Warn(fmt.Sprintf("couldn't get object attributes: %v", err))
			return
		}
		if hs.lastUpdated == updated {
			hs.Logger.Debug("configuration hasn't changed, skipping fetching full object")
			return
		}
		msg, err := hs.fetchObject(ctx, bucket)
		if err != nil {
			hs.Logger.Warn(fmt.Sprintf("couldn't get object: %v", err))
			return
		}
		hs.Logger.Info(fmt.Sprintf("configuration updated: %s", msg))
		dataSync <- sync.DataSync{FlagData: msg, Source: hs.Bucket + hs.Object, Type: sync.ALL}
		hs.lastUpdated = updated
	})

	hs.Cron.Start()

	<-ctx.Done()
	hs.Cron.Stop()

	return nil
}

func (hs *Sync) ReSync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	bucket, err := hs.getBucket(ctx)
	if err != nil {
		return err
	}
	defer bucket.Close()
	updated, err := hs.fetchObjectModificationTime(ctx, bucket)
	if err != nil {
		return err
	}
	msg, err := hs.fetchObject(ctx, bucket)
	if err != nil {
		return err
	}
	hs.Logger.Info(fmt.Sprintf("configuration updated: %s", msg))
	dataSync <- sync.DataSync{FlagData: msg, Source: hs.Bucket + hs.Object, Type: sync.ALL}
	hs.lastUpdated = updated
	return nil
}

func (hs *Sync) getBucket(ctx context.Context) (*blob.Bucket, error) {
	if hs.Bucket == "" {
		return nil, errors.New("no bucket string set")
	}
	return hs.BlobURLMux.OpenBucket(ctx, hs.Bucket)
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
	if hs.Object == "" {
		return "", errors.New("no object string set")
	}
	r, err := bucket.NewReader(ctx, hs.Object, nil)
	if err != nil {
		return "", fmt.Errorf("error creating reader for object %s/%s: %w", hs.Bucket, hs.Object, err)
	}
	defer r.Close()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, r)
	if err != nil {
		return "", fmt.Errorf("error reading object %s/%s: %w", hs.Bucket, hs.Object, err)
	}

	return string(buf.Bytes()), nil
}
