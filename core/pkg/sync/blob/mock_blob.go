package blob

import (
	"context"
	"log"
	"net/url"

	"gocloud.dev/blob"
	"gocloud.dev/blob/memblob"
)

type MockBlob struct {
	mux    *blob.URLMux
	bucket *blob.Bucket
	scheme string
	opener *fakeOpener
}

type fakeOpener struct {
	object      string
	content     string
	keepModTime bool
	getSync     func() *Sync
}

func (f *fakeOpener) OpenBucketURL(ctx context.Context, u *url.URL) (*blob.Bucket, error) {
	bucketUrl, err := url.Parse("mem://")
	if err != nil {
		log.Fatalf("couldn't parse url: %s: %v", "mem://", err)
	}
	opener := &memblob.URLOpener{}
	bucket, err := opener.OpenBucketURL(context.Background(), bucketUrl)
	if err != nil {
		log.Fatalf("couldn't open in memory bucket: %v", err)
	}
	if f.object != "" {
		err = bucket.WriteAll(ctx, f.object, []byte(f.content), nil)
		if err != nil {
			log.Fatalf("couldn't write in memory file: %v", err)
		}
	}
	if f.keepModTime && f.object != "" {
		attrs, err := bucket.Attributes(ctx, f.object)
		if err != nil {
			log.Fatalf("couldn't get memory file attributes: %v", err)
		}
		f.getSync().lastUpdated = attrs.ModTime
	} else {
		f.keepModTime = true
	}
	return bucket, nil
}

func NewMockBlob(scheme string, getSync func() *Sync) *MockBlob {
	mux := new(blob.URLMux)
	opener := &fakeOpener{getSync: getSync}
	mux.RegisterBucket(scheme, opener)
	return &MockBlob{
		mux:    mux,
		scheme: scheme,
		opener: opener,
	}
}

func (mb *MockBlob) URLMux() *blob.URLMux {
	return mb.mux
}

func (mb *MockBlob) AddObject(object, content string) {
	mb.opener.object = object
	mb.opener.content = content
	mb.opener.keepModTime = false
}
