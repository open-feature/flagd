package file

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	msync "sync"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
)

const (
	fetchFileName     = "to_fetch.json"
	fetchFileContents = "fetch me"
)

func TestSimpleReSync(t *testing.T) {
	fetchDirName := t.TempDir()
	source := filepath.Join(fetchDirName, fetchFileName)
	expectedDataSync := sync.DataSync{
		FlagData: "hello",
		Source:   source,
	}
	handler := Sync{
		URI:    source,
		Logger: logger.NewLogger(nil, false),
	}

	createFile(t, fetchDirName)
	writeToFile(t, fetchDirName, "hello")
	ctx := context.Background()
	dataSyncChan := make(chan sync.DataSync, 1)

	go func() {
		err := handler.ReSync(ctx, dataSyncChan)
		if err != nil {
			log.Fatalf("Error start sync: %s", err.Error())
			return
		}
	}()

	select {
	case s := <-dataSyncChan:
		if !reflect.DeepEqual(expectedDataSync, s) {
			t.Errorf("resync failed, incorrect datasync value, got %v want %v", s, expectedDataSync)
		}
	case <-time.After(5 * time.Second):
		t.Error("timed out waiting for datasync")
	}
}

func TestSimpleSync(t *testing.T) {
	readDirName := t.TempDir()
	updateDirName := t.TempDir()
	deleteDirName := t.TempDir()
	tests := map[string]struct {
		manipulationFuncs []func(t *testing.T)
		expectedDataSync  []sync.DataSync
		fetchDirName      string
	}{
		"simple-read": {
			fetchDirName: readDirName,
			manipulationFuncs: []func(t *testing.T){
				func(t *testing.T) {
					writeToFile(t, readDirName, fetchFileContents)
				},
			},
			expectedDataSync: []sync.DataSync{
				{
					FlagData: fetchFileContents,
					Source:   fmt.Sprintf("%s/%s", readDirName, fetchFileName),
				},
			},
		},
		"update-event": {
			fetchDirName: updateDirName,
			manipulationFuncs: []func(t *testing.T){
				func(t *testing.T) {
					writeToFile(t, updateDirName, fetchFileContents)
				},
				func(t *testing.T) {
					writeToFile(t, updateDirName, "new content")
				},
			},
			expectedDataSync: []sync.DataSync{
				{
					FlagData: fetchFileContents,
					Source:   fmt.Sprintf("%s/%s", updateDirName, fetchFileName),
				},
				{
					FlagData: "new content",
					Source:   fmt.Sprintf("%s/%s", updateDirName, fetchFileName),
				},
			},
		},
		"delete-event": {
			fetchDirName: deleteDirName,
			manipulationFuncs: []func(t *testing.T){
				func(t *testing.T) {
					writeToFile(t, deleteDirName, fetchFileContents)
				},
				func(t *testing.T) {
					deleteFile(t, deleteDirName, fetchFileName)
				},
			},
			expectedDataSync: []sync.DataSync{
				{
					FlagData: fetchFileContents,
					Source:   fmt.Sprintf("%s/%s", deleteDirName, fetchFileName),
				},
				{
					FlagData: defaultState,
					Source:   fmt.Sprintf("%s/%s", deleteDirName, fetchFileName),
				},
			},
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			createFile(t, tt.fetchDirName)

			ctx := context.Background()

			dataSyncChan := make(chan sync.DataSync, len(tt.expectedDataSync))

			syncHandler := Sync{
				URI:    fmt.Sprintf("%s/%s", tt.fetchDirName, fetchFileName),
				Logger: logger.NewLogger(nil, false),
				Mux:    &msync.RWMutex{},
			}

			go func() {
				err := syncHandler.Init(ctx)
				if err != nil {
					log.Fatalf("Error init sync: %s", err.Error())
					return
				}
				err = syncHandler.Sync(ctx, dataSyncChan)
				if err != nil {
					log.Fatalf("Error start sync: %s", err.Error())
					return
				}
			}()

			// file sync perform an initial fetch and then watch for file events
			init := <-dataSyncChan
			if init.FlagData != defaultState {
				t.Errorf("initial fetch for empty file expected to return default state: %s", defaultState)
			}

			for i, manipulation := range tt.manipulationFuncs {
				syncEvent := tt.expectedDataSync[i]
				manipulation(t)
				select {
				case data := <-dataSyncChan:
					if data.FlagData != syncEvent.FlagData {
						t.Errorf("expected content: %s, but received content: %s", syncEvent.FlagData, data.FlagData)
					}
					if data.Source != syncEvent.Source {
						t.Errorf("expected source: %s, but received source: %s", syncEvent.Source, data.Source)
					}
				case <-time.After(10 * time.Second):
					t.Errorf("event not found, timeout out after 10 seconds")
				}

				// validate readiness - readiness must not change
				if syncHandler.ready != true {
					t.Errorf("readiness must be set to true, but found: %t", syncHandler.ready)
				}
			}
		})
	}
}

func TestFilePathSync_Fetch(t *testing.T) {
	successDirName := t.TempDir()
	failureDirName := t.TempDir()
	tests := map[string]struct {
		fpSync         Sync
		handleResponse func(t *testing.T, fetched string, err error)
		fetchDirName   string
	}{
		"success": {
			fetchDirName: successDirName,
			fpSync: Sync{
				URI:    fmt.Sprintf("%s/%s", successDirName, fetchFileName),
				Logger: logger.NewLogger(nil, false),
			},
			handleResponse: func(t *testing.T, fetched string, err error) {
				if err != nil {
					t.Error(err)
				}

				if fetched != fetchFileContents {
					t.Errorf("expected fetched to be: '%s', got: '%s'", fetchFileContents, fetched)
				}
			},
		},
		"not found": {
			fetchDirName: failureDirName,
			fpSync: Sync{
				URI:    fmt.Sprintf("%s/%s", failureDirName, "not_found"),
				Logger: logger.NewLogger(nil, false),
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
			createFile(t, tt.fetchDirName)
			writeToFile(t, tt.fetchDirName, fetchFileContents)

			data, err := tt.fpSync.fetch(context.Background())

			tt.handleResponse(t, data, err)
		})
	}
}

func TestIsReadySyncFlag(t *testing.T) {
	fetchDirName := t.TempDir()
	fpSync := Sync{
		URI:    fmt.Sprintf("%s/%s", fetchDirName, fetchFileName),
		Logger: logger.NewLogger(nil, false),
		Mux:    &msync.RWMutex{},
	}

	createFile(t, fetchDirName)
	writeToFile(t, fetchDirName, fetchFileContents)
	if fpSync.IsReady() != false {
		t.Errorf("expected not to be ready")
	}
	ctx := context.TODO()
	err := fpSync.Init(ctx)
	if err != nil {
		log.Printf("Error init sync: %s", err.Error())
		return
	}
	if fpSync.IsReady() != false {
		t.Errorf("expected not to be ready")
	}
	dataSyncChan := make(chan sync.DataSync, 1)

	go func() {
		err = fpSync.Sync(ctx, dataSyncChan)
		if err != nil {
			log.Fatalf("Error start sync: %s", err.Error())
			return
		}
	}()
	time.Sleep(1 * time.Second)
	if fpSync.IsReady() != true {
		t.Errorf("expected to be ready")
	}
}

func deleteFile(t *testing.T, dirName string, fileName string) {
	if err := os.Remove(fmt.Sprintf("%s/%s", dirName, fileName)); err != nil {
		t.Fatal(err)
	}
}

func createFile(t *testing.T, fetchDirName string) {
	f, err := os.Create(fmt.Sprintf("%s/%s", fetchDirName, fetchFileName))
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Fatalf("close file: %v", err)
		}
	}(f)
	if err != nil {
		t.Fatal(err)
	}
}

func writeToFile(t *testing.T, fetchDirName, fileContents string) {
	file, err := os.OpenFile(fmt.Sprintf("%s/%s", fetchDirName, fetchFileName), os.O_RDWR, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Fatalf("close file: %v", err)
		}
	}(file)

	_, err = file.WriteAt([]byte(fileContents), 0)
	if err != nil {
		t.Fatal(err)
	}
}

// reAddMockWatcher is a minimal file.Watcher used to exercise reAddWatcher's
// retry behavior without touching the real filesystem. Its Add fails the first
// addFailUntil calls, then succeeds and records the watched name.
type reAddMockWatcher struct {
	addFailUntil int
	addErr       error
	addCalls     int
	watchList    []string
}

func (m *reAddMockWatcher) Close() error { return nil }

func (m *reAddMockWatcher) Add(name string) error {
	m.addCalls++
	if m.addCalls <= m.addFailUntil {
		return m.addErr
	}
	m.watchList = append(m.watchList, name)
	return nil
}

func (m *reAddMockWatcher) Remove(string) error         { return nil }
func (m *reAddMockWatcher) WatchList() []string         { return m.watchList }
func (m *reAddMockWatcher) Events() chan fsnotify.Event { return make(chan fsnotify.Event) }
func (m *reAddMockWatcher) Errors() chan error          { return make(chan error) }

func newReAddSync(w Watcher) *Sync {
	return &Sync{
		URI:     "/tmp/does-not-matter/flags.json",
		Logger:  logger.NewLogger(nil, false),
		watcher: w,
		Mux:     &msync.RWMutex{},
	}
}

// Happy path: the file is momentarily absent (Add fails once) then restored, so
// the retry re-establishes the watch.
func TestReAddWatcher_RetriesUntilSuccess(t *testing.T) {
	mw := &reAddMockWatcher{addFailUntil: 1, addErr: errors.New("no such file or directory")}
	fs := newReAddSync(mw)

	if err := fs.reAddWatcher(context.Background()); err != nil {
		t.Fatalf("expected watch to be re-established after a transient failure, got error: %v", err)
	}
	if mw.addCalls != 2 {
		t.Errorf("expected 2 Add attempts (1 fail + 1 success), got %d", mw.addCalls)
	}
	if len(mw.watchList) != 1 || mw.watchList[0] != fs.URI {
		t.Errorf("expected watch list to contain %q, got %v", fs.URI, mw.watchList)
	}
}

// Edge case: the file is truly deleted, so every Add fails and the retry gives
// up after the bounded number of attempts, returning the last error.
func TestReAddWatcher_ExhaustsAttempts(t *testing.T) {
	wantErr := errors.New("no such file or directory")
	mw := &reAddMockWatcher{addFailUntil: reAddMaxAttempts + 1, addErr: wantErr}
	fs := newReAddSync(mw)

	err := fs.reAddWatcher(context.Background())
	if err == nil {
		t.Fatal("expected an error after exhausting all re-add attempts, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("expected the last Add error, got %v", err)
	}
	if mw.addCalls != reAddMaxAttempts {
		t.Errorf("expected exactly %d Add attempts, got %d", reAddMaxAttempts, mw.addCalls)
	}
}

// Cancellation: a context cancelled during the backoff aborts the retry
// promptly instead of running out the full attempt budget.
func TestReAddWatcher_ContextCancelled(t *testing.T) {
	mw := &reAddMockWatcher{addFailUntil: reAddMaxAttempts + 1, addErr: errors.New("no such file or directory")}
	fs := newReAddSync(mw)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancelled before the first backoff completes

	start := time.Now()
	err := fs.reAddWatcher(ctx)
	elapsed := time.Since(start)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if mw.addCalls >= reAddMaxAttempts {
		t.Errorf("expected retry to abort early, but it made %d attempts", mw.addCalls)
	}
	if elapsed >= reAddMaxAttempts*reAddBackoff {
		t.Errorf("expected prompt return on cancellation, took %s", elapsed)
	}
}
