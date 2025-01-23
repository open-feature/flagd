package file

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	msync "sync"

	"github.com/fsnotify/fsnotify"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	"github.com/open-feature/flagd/core/pkg/utils"
)

const (
	FSNOTIFY = "fsnotify"
	FILEINFO = "fileinfo"
)

type Watcher interface {
	Close() error
	Add(name string) error
	Remove(name string) error
	WatchList() []string
	Events() chan fsnotify.Event
	Errors() chan error
}

type Sync struct {
	URI    string
	Logger *logger.Logger
	// watchType indicates how to watch the file FSNOTIFY|FILEINFO
	watchType string
	watcher   Watcher
	ready     bool
	Mux       *msync.RWMutex
}

func NewFileSync(uri string, watchType string, logger *logger.Logger) *Sync {
	return &Sync{
		URI:       uri,
		watchType: watchType,
		Logger:    logger,
		Mux:       &msync.RWMutex{},
	}
}

// default state is used to prevent EOF errors when handling filepath delete events + empty files
const defaultState = "{}"

func (fs *Sync) ReSync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	fs.sendDataSync(ctx, sync.ALL, dataSync)
	return nil
}

func (fs *Sync) Init(ctx context.Context) error {
	fs.Logger.Info("Starting filepath sync notifier")

	switch fs.watchType {
	case FSNOTIFY, "":
		w, err := NewFSNotifyWatcher()
		if err != nil {
			return fmt.Errorf("error creating fsnotify watcher: %w", err)
		}
		fs.watcher = w
	case FILEINFO:
		w := NewFileInfoWatcher(ctx, fs.Logger)
		fs.watcher = w
	default:
		return fmt.Errorf("unknown watcher type: '%s'", fs.watchType)
	}

	if err := fs.watcher.Add(fs.URI); err != nil {
		return fmt.Errorf("error adding watcher %s: %w", fs.URI, err)
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
		case event, ok := <-fs.watcher.Events():
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

		case err, ok := <-fs.watcher.Errors():
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

	if syncType == sync.DELETE {
		// Skip fetching and emit default state to avoid EOF errors
		dataSync <- sync.DataSync{FlagData: defaultState, Source: fs.URI, Type: syncType}
		return
	}

	msg := defaultState
	m, err := fs.fetch(ctx)
	if err != nil {
		fs.Logger.Error(fmt.Sprintf("Error fetching %s: %s", fs.URI, err.Error()))
	}
	if m == "" {
		fs.Logger.Warn(fmt.Sprintf("file %s is empty", fs.URI))
	} else {
		msg = m
	}

	dataSync <- sync.DataSync{FlagData: msg, Source: fs.URI, Type: syncType}
}

func (fs *Sync) fetch(_ context.Context) (string, error) {
	if fs.URI == "" {
		return "", errors.New("no filepath string set")
	}

	file, err := os.Open(fs.URI)
	if err != nil {
		return "", fmt.Errorf("error opening file %s: %w", fs.URI, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %w", fs.URI, err)
	}

	// File extension is used to determine the content type, so media type is unnecessary
	json, err := utils.ConvertToJSON(data, filepath.Ext(fs.URI), "")
	if err != nil {
		return "", fmt.Errorf("error converting file content to json: %w", err)
	}
	return json, nil
}
