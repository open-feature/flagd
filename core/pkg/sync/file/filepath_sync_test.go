package file

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	msync "sync"
	"testing"
	"time"

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
