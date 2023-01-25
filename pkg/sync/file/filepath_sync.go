package file

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/open-feature/flagd/pkg/sync"

	"gopkg.in/yaml.v3"

	"github.com/fsnotify/fsnotify"
	"github.com/open-feature/flagd/pkg/logger"
)

type Sync struct {
	URI          string
	Logger       *logger.Logger
	ProviderArgs sync.ProviderArgs
	// FileType indicates the file type e.g., json, yaml/yml etc.,
	fileType string
}

//nolint:funlen
func (fs *Sync) Sync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	fs.Logger.Info("Starting filepath sync notifier")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(fs.URI)
	if err != nil {
		return err
	}

	// file watcher is ready(and stable), fetch and emit the initial results
	fetch, err := fs.fetch(ctx)
	if err != nil {
		return err
	}

	dataSync <- sync.DataSync{FlagData: fetch, Source: fs.URI}

	fs.Logger.Info(fmt.Sprintf("watching filepath: %s", fs.URI))
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				fs.Logger.Info("filepath notifier closed")
				return errors.New("filepath notifier closed")
			}

			fs.Logger.Info(fmt.Sprintf("filepath event: %s %s", event.Name, event.Op.String()))

			switch event.Op {
			case fsnotify.Create:
				fs.Logger.Debug("new configuration created")
				msg, err := fs.fetch(ctx)
				if err != nil {
					fs.Logger.Error(fmt.Sprintf("error fetching after Create notification: %s", err.Error()))
					continue
				}

				dataSync <- sync.DataSync{FlagData: msg, Source: fs.URI}
			case fsnotify.Write:
				fs.Logger.Debug("configuration modified")
				msg, err := fs.fetch(ctx)
				if err != nil {
					fs.Logger.Error(fmt.Sprintf("error fetching after write notification: %s", err.Error()))
					continue
				}

				dataSync <- sync.DataSync{FlagData: msg, Source: fs.URI}
			case fsnotify.Remove:
				// K8s exposes config maps as symlinks.
				// Updates cause a remove event, we need to re-add the watcher in this case.
				err = watcher.Add(fs.URI)
				if err != nil {
					fs.Logger.Error(fmt.Sprintf("error restoring watcher, file may have been deleted: %s", err.Error()))
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return errors.New("watcher error")
			}

			fs.Logger.Error(err.Error())
		case <-ctx.Done():
			fs.Logger.Debug("exiting file watcher")
			return nil
		}
	}
}

func (fs *Sync) fetch(_ context.Context) (string, error) {
	if fs.URI == "" {
		return "", errors.New("no filepath string set")
	}
	if fs.fileType == "" {
		uriSplit := strings.Split(fs.URI, ".")
		fs.fileType = uriSplit[len(uriSplit)-1]
	}
	rawFile, err := os.ReadFile(fs.URI)
	if err != nil {
		return "", err
	}

	switch fs.fileType {
	case "yaml", "yml":
		return yamlToJSON(rawFile)
	case "json":
		return string(rawFile), nil
	default:
		return "", fmt.Errorf("filepath extension for URI: '%s' is not supported", fs.URI)
	}
}

// yamlToJSON is a generic helper function to convert
// yaml to json
func yamlToJSON(rawFile []byte) (string, error) {
	var ms map[string]interface{}
	// yaml.Unmarshal unmarshals to map[interface]interface{}
	if err := yaml.Unmarshal(rawFile, &ms); err != nil {
		return "", fmt.Errorf("unmarshal yaml: %w", err)
	}

	r, err := json.Marshal(ms)
	if err != nil {
		return "", fmt.Errorf("convert yaml to json: %w", err)
	}

	return string(r), err
}
