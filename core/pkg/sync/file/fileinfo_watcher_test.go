package file

import (
	"io/fs"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/google/go-cmp/cmp"
)

func Test_fileInfoWatcher_Close(t *testing.T) {
	type fields struct{}
	tests := []struct {
		name    string
		watcher *fileInfoWatcher
		wantErr bool
	}{
		{
			name:    "all chans close",
			watcher: makeTestWatcher(t, map[string]fs.FileInfo{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.watcher.Close(); (err != nil) != tt.wantErr {
				t.Errorf("fileInfoWatcher.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, ok := (<-tt.watcher.Errors()); ok != false {
				t.Error("fileInfoWatcher.Close() failed to close error chan")
			}
			if _, ok := (<-tt.watcher.Events()); ok != false {
				t.Error("fileInfoWatcher.Close() failed to close events chan")
			}
		})
	}
}

func Test_fileInfoWatcher_Add(t *testing.T) {
	tests := []struct {
		name    string
		watcher *fileInfoWatcher
		add     []string
		want    map[string]fs.FileInfo
		wantErr bool
	}{
		{
			name:    "add one watch",
			watcher: makeTestWatcher(t, map[string]fs.FileInfo{}),
			add:     []string{"/foo"},
			want: map[string]fs.FileInfo{
				"/foo": &mockFileInfo{},
			},
		},
	}
	for _, tt := range tests {
		tt.watcher.statFunc = makeStatFunc(t, &mockFileInfo{})
		t.Run(tt.name, func(t *testing.T) {
			for _, path := range tt.add {
				if err := tt.watcher.Add(path); (err != nil) != tt.wantErr {
					t.Errorf("fileInfoWatcher.Add() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if !cmp.Equal(tt.watcher.watches, tt.want, cmp.AllowUnexported(mockFileInfo{})) {
				t.Errorf("fileInfoWatcher.Add(): want-, got+: %v ", cmp.Diff(tt.want, tt.watcher.watches))
			}
		})
	}
}

// makeTestWatcher returns a pointer to a fileInfoWatcher suitable for testing
func makeTestWatcher(t *testing.T, watches map[string]fs.FileInfo) *fileInfoWatcher {
	t.Helper()

	return &fileInfoWatcher{
		evChan:  make(chan fsnotify.Event),
		erChan:  make(chan error),
		watches: watches,
	}
}

// makeStateFunc returns an os.Stat wrapper that parrots back whatever its
// constructor is given
func makeStatFunc(t *testing.T, fi fs.FileInfo) func(string) (fs.FileInfo, error) {
	t.Helper()
	return func(s string) (fs.FileInfo, error) {
		return fi, nil
	}
}

// mockFileInfo implements fs.FileInfo for mocks
type mockFileInfo struct {
	name    string      // base name of the file
	size    int64       // length in bytes for regular files; system-dependent for others
	mode    fs.FileMode // file mode bits
	modTime time.Time   // modification time
}

// explicitly impements fs.FileInfo
var _ fs.FileInfo = &mockFileInfo{}

func (mfi *mockFileInfo) Name() string {
	return mfi.name
}

func (mfi *mockFileInfo) Size() int64 {
	return mfi.size
}

func (mfi *mockFileInfo) Mode() fs.FileMode {
	return mfi.mode
}

func (mfi *mockFileInfo) ModTime() time.Time {
	return mfi.modTime
}

func (mfi *mockFileInfo) IsDir() bool {
	return false
}

func (mfi *mockFileInfo) Sys() any {
	return "foo"
}
