package file

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	msync "sync"

	"gopkg.in/yaml.v3"

	"github.com/fsnotify/fsnotify"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
)

type Sync struct {
	URI    string
	Logger *logger.Logger
	// FileType indicates the file type e.g., json, yaml/yml etc.,
	fileType string
	watcher  *fsnotify.Watcher
	ready    bool
	Mux      *msync.RWMutex
}

// default state is used to prevent EOF errors when handling filepath delete events + empty files
const defaultState = "{}"

func (fs *Sync) ReSync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	fs.sendDataSync(ctx, sync.ALL, dataSync)
	return nil
}

func (fs *Sync) Init(_ context.Context) error {
	fs.Logger.Info("Starting filepath sync notifier")
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	fs.watcher = w
	err = fs.watcher.Add(fs.URI)
	if err != nil {
		return err
	}
	return nil
}

func (fs *Sync) IsReady() bool {
	fs.Mux.RLock()
	defer fs.Mux.RUnlock()
	return fs.ready
}

func (fs *Sync) setReady(val bool) {
	fs.Mux.Lock()
	defer fs.Mux.Unlock()
	fs.ready = val
}

//nolint:funlen
func (fs *Sync) Sync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	defer fs.watcher.Close()
	fs.sendDataSync(ctx, sync.ALL, dataSync)
	fs.setReady(true)
	fs.Logger.Info(fmt.Sprintf("watching filepath: %s", fs.URI))
	for {
		select {
		case event, ok := <-fs.watcher.Events:
			if !ok {
				fs.Logger.Info("filepath notifier closed")
				return errors.New("filepath notifier closed")
			}

			fs.Logger.Info(fmt.Sprintf("filepath event: %s %s", event.Name, event.Op.String()))
			switch {
			case event.Has(fsnotify.Create) || event.Has(fsnotify.Write):
				fs.sendDataSync(ctx, sync.ALL, dataSync)
			case event.Has(fsnotify.Remove):
				// K8s exposes config maps as symlinks.
				// Updates cause a remove event, we need to re-add the watcher in this case.
				err := fs.watcher.Add(fs.URI)
				if err != nil {
					// the watcher could not be re-added, so the file must have been deleted
					fs.Logger.Error(fmt.Sprintf("error restoring watcher, file may have been deleted: %s", err.Error()))
					fs.sendDataSync(ctx, sync.DELETE, dataSync)
					continue
				}

				// Counterintuitively, remove events are the only meaningful ones seen in K8s.
				// K8s handles mounted ConfigMap updates by modifying symbolic links, which is an atomic operation.
				// At the point the remove event is fired, we have our new data, so we can send it down the channel.
				fs.sendDataSync(ctx, sync.ALL, dataSync)
			case event.Has(fsnotify.Chmod):
				// on linux the REMOVE event will not fire until all file descriptors are closed, this cannot happen
				// while the file is being watched, os.Stat is used here to infer deletion
				if _, err := os.Stat(fs.URI); errors.Is(err, os.ErrNotExist) {
					fs.Logger.Error(fmt.Sprintf("file has been deleted: %s", err.Error()))
					fs.sendDataSync(ctx, sync.DELETE, dataSync)
				}
			}

		case err, ok := <-fs.watcher.Errors:
			if !ok {
				fs.setReady(false)
				return errors.New("watcher error")
			}

			fs.Logger.Error(err.Error())
		case <-ctx.Done():
			fs.Logger.Debug("exiting file watcher")
			return nil
		}
	}
}

func (fs *Sync) sendDataSync(ctx context.Context, syncType sync.Type, dataSync chan<- sync.DataSync) {
	fs.Logger.Debug(fmt.Sprintf("Configuration %s:  %s", fs.URI, syncType.String()))

	msg := defaultState
	if syncType != sync.DELETE {
		m, err := fs.fetch(ctx)
		if err != nil {
			fs.Logger.Error(fmt.Sprintf("Error fetching %s: %s", fs.URI, err.Error()))
		}
		if m == "" {
			fs.Logger.Warn(fmt.Sprintf("file %s is empty", fs.URI))
		} else {
			msg = m
		}
	}

	dataSync <- sync.DataSync{FlagData: msg, Source: fs.URI, Type: syncType}
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
