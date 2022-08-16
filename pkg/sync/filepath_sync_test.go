package sync_test

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/open-feature/flagd/pkg/sync"
	log "github.com/sirupsen/logrus"
)

const (
	dirName           = "test"
	createFileName    = "to_create"
	modifyFileName    = "to_modify"
	deleteFileName    = "to_delete"
	fetchFileName     = "to_fetch"
	fetchFileContents = "fetch me"
)

func TestFilePathSync_Notify(t *testing.T) {
	tests := map[string]struct {
		triggerEvent      func(t *testing.T)
		expectedEventType sync.DefaultEventType
	}{
		"create event": {
			triggerEvent: func(t *testing.T) {
				if _, err := os.Create(fmt.Sprintf("%s/%s", dirName, createFileName)); err != nil {
					t.Fatal(err)
				}
			},
			expectedEventType: sync.DefaultEventTypeCreate,
		},
		"modify event": {
			triggerEvent: func(t *testing.T) {
				file, err := os.OpenFile(fmt.Sprintf("%s/%s", dirName, modifyFileName), os.O_RDWR, 0o644)
				if err != nil {
					t.Fatal(err)
				}
				defer func(file *os.File) {
					if err := file.Close(); err != nil {
						log.Fatalf("close file: %v", err)
					}
				}(file)

				_, err = file.WriteAt([]byte("foo"), 0)
				if err != nil {
					t.Fatal(err)
				}
			},
			expectedEventType: sync.DefaultEventTypeModify,
		},
		"delete event": {
			triggerEvent: func(t *testing.T) {
				if err := os.Remove(fmt.Sprintf("%s/%s", dirName, deleteFileName)); err != nil {
					t.Fatal(err)
				}
			},
			expectedEventType: sync.DefaultEventTypeDelete,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			setupFilePathNotify(t)
			defer t.Cleanup(cleanupFilePath)

			// prevent deadlock with a timeout if expected event doesn't arrive
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			fpSync := sync.FilePathSync{
				URI:    dirName,
				Logger: log.WithFields(log.Fields{}),
			}
			inotifyChan := make(chan sync.INotify)

			go func() {
				fpSync.Notify(ctx, inotifyChan)
			}()

			w := <-inotifyChan // first emitted event by Notify is to signal readiness
			if w.GetEvent().EventType != sync.DefaultEventTypeReady {
				t.Errorf("expected event type to be %d, got %d", sync.DefaultEventTypeReady, w.GetEvent().EventType)
			}

			tt.triggerEvent(t)

			for {
				select {
				case event, ok := <-inotifyChan:
					if !ok {
						t.Fatal("inotify chan closed")
					}
					if event.GetEvent().EventType != tt.expectedEventType {
						t.Errorf(
							"expected event of type %d, got %d", tt.expectedEventType, event.GetEvent().EventType,
						)
					}
					return
				case <-ctx.Done():
					t.Error("context timed out")
					return
				}
			}
		})
	}
}

func TestFilePathSync_Fetch(t *testing.T) {
	tests := map[string]struct {
		fpSync         sync.FilePathSync
		handleResponse func(t *testing.T, fetched string, err error)
	}{
		"success": {
			fpSync: sync.FilePathSync{
				URI:    fmt.Sprintf("%s/%s", dirName, fetchFileName),
				Logger: log.WithFields(log.Fields{}),
			},
			handleResponse: func(t *testing.T, fetched string, err error) {
				if err != nil {
					t.Error(err)
				}

				if fetched != fetchFileContents {
					t.Errorf("expected fetched to be '%s', got '%s'", fetchFileContents, fetched)
				}
			},
		},
		"not found": {
			fpSync: sync.FilePathSync{
				URI:    fmt.Sprintf("%s/%s", dirName, "not_found"),
				Logger: log.WithFields(log.Fields{}),
			},
			handleResponse: func(t *testing.T, fetched string, err error) {
				if err == nil {
					t.Error("expected an error, got nil")
				}
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			setupFilePathFetch(t)
			defer t.Cleanup(cleanupFilePath)

			fetched, err := tt.fpSync.Fetch(context.Background())

			tt.handleResponse(t, fetched, err)
		})
	}
}

func BenchmarkFilePathSync_Fetch(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	tests := map[string]struct {
		fpSync         sync.FilePathSync
		handleResponse func(b *testing.B, fetched string, err error)
	}{
		"success": {
			fpSync: sync.FilePathSync{
				URI:    fmt.Sprintf("%s/%s", dirName, fetchFileName),
				Logger: log.WithFields(log.Fields{}),
			},
			handleResponse: func(b *testing.B, fetched string, err error) {
				if err != nil {
					b.Error(err)
				}

				if fetched != fetchFileContents {
					b.Errorf("expected fetched to be '%s', got '%s'", fetchFileContents, fetched)
				}
			},
		},
		"not found": {
			fpSync: sync.FilePathSync{
				URI:    fmt.Sprintf("%s/%s", dirName, "not_found"),
				Logger: log.WithFields(log.Fields{}),
			},
			handleResponse: func(b *testing.B, fetched string, err error) {
				if err == nil {
					b.Error("expected an error, got nil")
				}
			},
		},
	}

	for name, tt := range tests {
		b.Run(name, func(b *testing.B) {
			setupFilePathFetchBenchmark(b)
			defer b.Cleanup(cleanupFilePath)
			for i := 0; i < b.N; i++ {
				fetched, err := tt.fpSync.Fetch(context.Background())
				tt.handleResponse(b, fetched, err)
			}
		})

	}
}

func setupFilePathNotify(t *testing.T) {
	if err := os.Mkdir(dirName, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Create(fmt.Sprintf("%s/%s", dirName, modifyFileName)); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Create(fmt.Sprintf("%s/%s", dirName, deleteFileName)); err != nil {
		t.Fatal(err)
	}
}

func cleanupFilePath() {
	if err := os.RemoveAll(dirName); err != nil {
		log.Fatalf("rmdir: %v", err)
	}
}

func setupFilePathFetch(t *testing.T) {
	if _, err := os.Stat(dirName); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dirName, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
	}
	if _, err := os.Create(fmt.Sprintf("%s/%s", dirName, fetchFileName)); err != nil {
		t.Fatal(err)
	}

	file, err := os.OpenFile(fmt.Sprintf("%s/%s", dirName, fetchFileName), os.O_RDWR, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Fatalf("close file: %v", err)
		}
	}(file)

	_, err = file.WriteAt([]byte(fetchFileContents), 0)
	if err != nil {
		t.Fatal(err)
	}
}

func setupFilePathFetchBenchmark(b *testing.B) {
	if err := os.Mkdir(dirName, os.ModePerm); err != nil {
		b.Fatal(err)
	}

	if _, err := os.Create(fmt.Sprintf("%s/%s", dirName, fetchFileName)); err != nil {
		b.Fatal(err)
	}

	file, err := os.OpenFile(fmt.Sprintf("%s/%s", dirName, fetchFileName), os.O_RDWR, 0o644)
	if err != nil {
		b.Fatal(err)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Fatalf("close file: %v", err)
		}
	}(file)

	_, err = file.WriteAt([]byte(fetchFileContents), 0)
	if err != nil {
		b.Fatal(err)
	}
}
