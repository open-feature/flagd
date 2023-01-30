package file

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/open-feature/flagd/pkg/sync"

	"github.com/open-feature/flagd/pkg/logger"
)

const (
	dirName           = "test"
	fetchFileName     = "to_fetch.json"
	fetchFileContents = "fetch me"
)

func TestSimpleSync(t *testing.T) {
	handler := Sync{
		URI:    fmt.Sprintf("%s/%s", dirName, fetchFileName),
		Logger: logger.NewLogger(nil, false),
	}

	defer t.Cleanup(cleanupFilePath)
	setupFilePathFetch(t)

	ctx := context.Background()
	dataSyncChan := make(chan sync.DataSync)

	go func() {
		err := handler.Sync(ctx, dataSyncChan)
		if err != nil {
			log.Fatalf("Error start sync: %s", err.Error())
			return
		}
	}()

	data := <-dataSyncChan

	if data.FlagData != fetchFileContents {
		t.Errorf("expected content: %s, but received content: %s", fetchFileContents, data.FlagData)
	}
}

func TestFilePathSync_Fetch(t *testing.T) {
	tests := map[string]struct {
		fpSync         Sync
		handleResponse func(t *testing.T, fetched string, err error)
	}{
		"success": {
			fpSync: Sync{
				URI:    fmt.Sprintf("%s/%s", dirName, fetchFileName),
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
				URI:    fmt.Sprintf("%s/%s", dirName, "not_found"),
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
			setupFilePathFetch(t)
			defer t.Cleanup(cleanupFilePath)

			data, err := tt.fpSync.fetch(context.Background())

			tt.handleResponse(t, data, err)
		})
	}
}

func cleanupFilePath() {
	if err := os.RemoveAll(dirName); err != nil {
		log.Fatalf("rmdir: %v", err)
	}
}

func setupFilePathFetch(t *testing.T) {
	if err := os.Mkdir(dirName, os.ModePerm); err != nil {
		t.Fatal(err)
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
