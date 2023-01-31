package file

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/open-feature/flagd/pkg/sync"

	"github.com/open-feature/flagd/pkg/logger"
)

const (
	fetchDirName      = "test"
	fetchFileName     = "to_fetch.json"
	fetchFileContents = "fetch me"
)

func TestSimpleSync(t *testing.T) {

	tests := map[string]struct {
		manipulationFuncs []func(t *testing.T)
		expectedDataSync  []sync.DataSync
	}{
		"simple-read": {
			manipulationFuncs: []func(t *testing.T){
				func(t *testing.T) {
					writeToFile(t, fetchDirName, fetchFileName, fetchFileContents)
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
					writeToFile(t, fetchDirName, fetchFileName, fetchFileContents)
				},
				func(t *testing.T) {
					writeToFile(t, fetchDirName, fetchFileName, "new content")
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
					writeToFile(t, fetchDirName, fetchFileName, fetchFileContents)
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
					FlagData: "",
					Source:   fmt.Sprintf("%s/%s", fetchDirName, fetchFileName),
					Type:     sync.ALL,
				},
			},
		},
	}

	handler := Sync{
		URI:    fmt.Sprintf("%s/%s", fetchDirName, fetchFileName),
		Logger: logger.NewLogger(nil, false),
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			defer t.Cleanup(cleanupFilePath)
			setupDir(t, fetchDirName)
			createFile(t, fetchDirName, fetchFileName)

			ctx := context.Background()
			dataSyncChan := make(chan sync.DataSync, len(tt.expectedDataSync))

			go func() {
				err := handler.Sync(ctx, dataSyncChan)
				if err != nil {
					log.Fatalf("Error start sync: %s", err.Error())
					return
				}
			}()

			for i, manipulation := range tt.manipulationFuncs {
				syncEvent := tt.expectedDataSync[i]
				manipulation(t)
				select {
				case data := <-dataSyncChan:
					fmt.Println(data)
					if data.FlagData != syncEvent.FlagData {
						t.Errorf("expected content: %s, but received content: %s", syncEvent.FlagData, data.FlagData)
					}
					if data.Source != syncEvent.Source {
						t.Errorf("expected source: %s, but received source: %s", syncEvent.Source, data.Source)
					}
					if data.Type != syncEvent.Type {
						t.Errorf("expected type: %b, but received type: %b", syncEvent.Type, data.Type)
					}
				case <-time.After(3 * time.Second):
					t.Errorf("event not found, timeout out after 3 seconds")
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
			setupDir(t, fetchDirName)
			createFile(t, fetchDirName, fetchFileName)
			writeToFile(t, fetchDirName, fetchFileName, fetchFileContents)
			defer t.Cleanup(cleanupFilePath)

			data, err := tt.fpSync.fetch(context.Background())

			tt.handleResponse(t, data, err)
		})
	}
}

func cleanupFilePath() {
	if err := os.RemoveAll(fetchDirName); err != nil {
		log.Fatalf("rmdir: %v", err)
	}
}

func deleteFile(t *testing.T, dirName string, fileName string) {
	if err := os.Remove(fmt.Sprintf(fmt.Sprintf("%s/%s", dirName, fileName))); err != nil {
		t.Fatal(err)
	}
}

func setupDir(t *testing.T, dirName string) {
	if err := os.Mkdir(dirName, os.ModePerm); err != nil {
		t.Fatal(err)
	}
}

func createFile(t *testing.T, dirName string, fileName string) {
	if _, err := os.Create(fmt.Sprintf("%s/%s", dirName, fileName)); err != nil {
		t.Fatal(err)
	}
}

func writeToFile(t *testing.T, dirName string, fileName string, fileContents string) {
	file, err := os.OpenFile(fmt.Sprintf("%s/%s", dirName, fileName), os.O_RDWR, 0o644)
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
