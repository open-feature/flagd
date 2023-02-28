package file

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	msync "sync"
	"testing"
	"time"

	"github.com/open-feature/flagd/internal/pkg/sync"

	"github.com/open-feature/flagd/internal/pkg/logger"
)

const (
	fetchDirName      = "test"
	fetchFileName     = "to_fetch.json"
	fetchFileContents = "fetch me"
)

func TestSimpleReSync(t *testing.T) {
	tests := map[string]struct {
		fileContents     string
		expectedDataSync sync.DataSync
	}{
		"simple-read": {
			fileContents: "hello",
			expectedDataSync: sync.DataSync{
				FlagData: "hello",
				Source:   fmt.Sprintf("%s/%s", fetchDirName, fetchFileName),
				Type:     sync.ALL,
			},
		},
	}

	handler := Sync{
		URI:    fmt.Sprintf("%s/%s", fetchDirName, fetchFileName),
		Logger: logger.NewLogger(nil, false),
		Source: fmt.Sprintf("%s/%s", fetchDirName, fetchFileName),
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			defer t.Cleanup(cleanupFilePath)
			setupDir(t)
			createFile(t)
			writeToFile(t, tt.fileContents)

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
				if !reflect.DeepEqual(tt.expectedDataSync, s) {
					t.Errorf("resync failed, incorrect datasync value, got %v want %v", s, tt.expectedDataSync)
				}
			case <-time.After(5 * time.Second):
				t.Error("timed out waiting for datasync")
			}
		})
	}
}

func TestSimpleSync(t *testing.T) {
	tests := map[string]struct {
		manipulationFuncs []func(t *testing.T)
		expectedDataSync  []sync.DataSync
	}{
		"simple-read": {
			manipulationFuncs: []func(t *testing.T){
				func(t *testing.T) {
					writeToFile(t, fetchFileContents)
				},
			},
			expectedDataSync: []sync.DataSync{
				{
					FlagData: fetchFileContents,
					Source:   fmt.Sprintf("%s/%s", fetchDirName, fetchFileName),
					Type:     sync.ALL,
				},
			},
		},
		"update-event": {
			manipulationFuncs: []func(t *testing.T){
				func(t *testing.T) {
					writeToFile(t, fetchFileContents)
				},
				func(t *testing.T) {
					writeToFile(t, "new content")
				},
			},
			expectedDataSync: []sync.DataSync{
				{
					FlagData: fetchFileContents,
					Source:   fmt.Sprintf("%s/%s", fetchDirName, fetchFileName),
					Type:     sync.ALL,
				},
				{
					FlagData: "new content",
					Source:   fmt.Sprintf("%s/%s", fetchDirName, fetchFileName),
					Type:     sync.ALL,
				},
			},
		},
		"delete-event": {
			manipulationFuncs: []func(t *testing.T){
				func(t *testing.T) {
					writeToFile(t, fetchFileContents)
				},
				func(t *testing.T) {
					deleteFile(t, fetchDirName, fetchFileName)
				},
			},
			expectedDataSync: []sync.DataSync{
				{
					FlagData: fetchFileContents,
					Source:   fmt.Sprintf("%s/%s", fetchDirName, fetchFileName),
					Type:     sync.ALL,
				},
				{
					FlagData: defaultState,
					Source:   fmt.Sprintf("%s/%s", fetchDirName, fetchFileName),
					Type:     sync.DELETE,
				},
			},
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			defer t.Cleanup(cleanupFilePath)
			setupDir(t)
			createFile(t)

			ctx := context.Background()

			dataSyncChan := make(chan sync.DataSync, len(tt.expectedDataSync))

			go func() {
				handler := Sync{
					URI:    fmt.Sprintf("%s/%s", fetchDirName, fetchFileName),
					Logger: logger.NewLogger(nil, false),
					Mux:    &msync.RWMutex{},
					Source: fmt.Sprintf("%s/%s", fetchDirName, fetchFileName),
				}
				err := handler.Init(ctx)
				if err != nil {
					log.Fatalf("Error init sync: %s", err.Error())
					return
				}
				err = handler.Sync(ctx, dataSyncChan)
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
					if data.Type != syncEvent.Type {
						t.Errorf("expected type: %b, but received type: %b", syncEvent.Type, data.Type)
					}
				case <-time.After(10 * time.Second):
					t.Errorf("event not found, timeout out after 10 seconds")
				}
			}
		})
	}
}

func TestFilePathSync_Fetch(t *testing.T) {
	tests := map[string]struct {
		fpSync         Sync
		handleResponse func(t *testing.T, fetched string, err error)
	}{
		"success": {
			fpSync: Sync{
				URI:    fmt.Sprintf("%s/%s", fetchDirName, fetchFileName),
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
			fpSync: Sync{
				URI:    fmt.Sprintf("%s/%s", fetchDirName, "not_found"),
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
			setupDir(t)
			createFile(t)
			writeToFile(t, fetchFileContents)
			defer t.Cleanup(cleanupFilePath)

			data, err := tt.fpSync.fetch(context.Background())

			tt.handleResponse(t, data, err)
		})
	}
}

func TestIsReadySyncFlag(t *testing.T) {
	fpSync := Sync{
		URI:    fmt.Sprintf("%s/%s", fetchDirName, fetchFileName),
		Logger: logger.NewLogger(nil, false),
		Mux:    &msync.RWMutex{},
	}

	setupDir(t)
	createFile(t)
	writeToFile(t, fetchFileContents)
	defer t.Cleanup(cleanupFilePath)
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

func cleanupFilePath() {
	if err := os.RemoveAll(fetchDirName); err != nil {
		log.Fatalf("rmdir: %v", err)
	}
}

func deleteFile(t *testing.T, dirName string, fileName string) {
	if err := os.Remove(fmt.Sprintf("%s/%s", dirName, fileName)); err != nil {
		t.Fatal(err)
	}
}

func setupDir(t *testing.T) {
	if err := os.Mkdir(fetchDirName, os.ModePerm); err != nil {
		t.Fatal(err)
	}
}

func createFile(t *testing.T) {
	if _, err := os.Create(fmt.Sprintf("%s/%s", fetchDirName, fetchFileName)); err != nil {
		t.Fatal(err)
	}
}

func writeToFile(t *testing.T, fileContents string) {
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
