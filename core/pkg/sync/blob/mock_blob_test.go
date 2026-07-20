package blob

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync/atomic"
	"time"

	"gocloud.dev/blob"
	"gocloud.dev/blob/driver"
	"gocloud.dev/gcerrors"
)

// MockBlob is a controllable blob source for tests. It is backed by a fake
// driver.Bucket so tests can set the object's ETag, ModTime and body
// independently, and can inspect how many full-object fetches occurred.
type MockBlob struct {
	mux    *blob.URLMux
	scheme string
	drv    *fakeBlobDriver
}

// NewMockBlob registers a fake bucket under scheme and returns a handle for
// controlling the object it serves.
func NewMockBlob(scheme string) *MockBlob {
	drv := &fakeBlobDriver{}
	mux := new(blob.URLMux)
	mux.RegisterBucket(scheme, &fakeDriverOpener{drv: drv})
	return &MockBlob{
		mux:    mux,
		scheme: scheme,
		drv:    drv,
	}
}

func (mb *MockBlob) URLMux() *blob.URLMux {
	return mb.mux
}

// AddObject sets the object's content and assigns a fresh ETag and ModTime
func (mb *MockBlob) AddObject(content string) {
	mb.drv.revision++
	mb.drv.etag = fmt.Sprintf("etag-%d", mb.drv.revision)
	mb.drv.modTime = time.Now()
	mb.drv.content = []byte(content)
}

// SetObject sets the object's ETag, ModTime and content explicitly
func (mb *MockBlob) SetObject(etag string, modTime time.Time, content string) {
	mb.drv.etag = etag
	mb.drv.modTime = modTime
	mb.drv.content = []byte(content)
}

// Reads reports how many full-object fetches (NewRangeReader calls) have happened
func (mb *MockBlob) Reads() int {
	return int(mb.drv.reads.Load())
}

// fakeBlobDriver is a controllable driver.Bucket. It keeps its state
// across OpenBucket calls and never regenerates ETags or ModTimes on
// its own, so tests decide exactly what a sync observes. Object state is
// only mutated between syncs, but reads is bumped from within concurrent
// syncs, so it is atomic to stay clean under `go test -race`.
type fakeBlobDriver struct {
	etag     string
	modTime  time.Time
	content  []byte
	reads    atomic.Int64 // number of NewRangeReader calls, i.e. full-object fetches
	revision int          // monotonically bumped by AddObject to mint fresh ETags
}

func (d *fakeBlobDriver) Attributes(_ context.Context, _ string) (*driver.Attributes, error) {
	return &driver.Attributes{
		ContentType: "application/json",
		ModTime:     d.modTime,
		Size:        int64(len(d.content)),
		ETag:        d.etag,
	}, nil
}

func (d *fakeBlobDriver) NewRangeReader(
	_ context.Context, _ string, offset, length int64, _ *driver.ReaderOptions,
) (driver.Reader, error) {
	d.reads.Add(1)
	data := d.content
	if offset > int64(len(data)) {
		offset = int64(len(data))
	}
	data = data[offset:]
	if length >= 0 && length < int64(len(data)) {
		data = data[:length]
	}
	return &fakeBlobReader{
		Reader: bytes.NewReader(data),
		attrs: driver.ReaderAttributes{
			ContentType: "application/json",
			ModTime:     d.modTime,
			Size:        int64(len(d.content)),
		},
	}, nil
}

var errNotImplemented = errors.New("not implemented")

func (d *fakeBlobDriver) ErrorCode(error) gcerrors.ErrorCode { return gcerrors.OK }
func (d *fakeBlobDriver) As(any) bool                        { return false }
func (d *fakeBlobDriver) ErrorAs(error, any) bool            { return false }
func (d *fakeBlobDriver) Close() error                       { return nil }

func (d *fakeBlobDriver) ListPaged(context.Context, *driver.ListOptions) (*driver.ListPage, error) {
	return nil, errNotImplemented
}

func (d *fakeBlobDriver) NewTypedWriter(
	context.Context, string, string, *driver.WriterOptions,
) (driver.Writer, error) {
	return nil, errNotImplemented
}

func (d *fakeBlobDriver) Copy(context.Context, string, string, *driver.CopyOptions) error {
	return errNotImplemented
}
func (d *fakeBlobDriver) Delete(context.Context, string) error { return errNotImplemented }

func (d *fakeBlobDriver) SignedURL(context.Context, string, *driver.SignedURLOptions) (string, error) {
	return "", errNotImplemented
}

// fakeBlobReader adapts a bytes.Reader to the driver.Reader interface.
type fakeBlobReader struct {
	*bytes.Reader
	attrs driver.ReaderAttributes
}

func (r *fakeBlobReader) Close() error                         { return nil }
func (r *fakeBlobReader) Attributes() *driver.ReaderAttributes { return &r.attrs }
func (r *fakeBlobReader) As(any) bool                          { return false }

type fakeDriverOpener struct{ drv *fakeBlobDriver }

func (o *fakeDriverOpener) OpenBucketURL(_ context.Context, _ *url.URL) (*blob.Bucket, error) {
	return blob.NewBucket(o.drv), nil
}
