package file

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/sync/errgroup"
)

// Implements file.Watcher using a timer and os.FileInfo
type fileInfoWatcher struct {
	// Event Chan
	evChan chan fsnotify.Event
	// Errors Chan
	erChan chan error
	// timer thread errgroup
	eg errgroup.Group
	// Func to wrap os.Stat (injection point for test helpers)
	statFunc func(string) (fs.FileInfo, error)
	// thread-safe interface to underlying files we are watching
	mu      sync.RWMutex
	watches map[string]fs.FileInfo // filename -> info
}

// NewFsNotifyWatcher returns a new fsNotifyWatcher
func NewFileInfoWatcher() *fileInfoWatcher {
	return &fileInfoWatcher{
		evChan: make(chan fsnotify.Event),
		erChan: make(chan error),
	}
}

// fileInfoWatcher explicitly implements file.Watcher
var _ Watcher = &fileInfoWatcher{}

// Close calls close on the underlying fsnotify.Watcher
func (f *fileInfoWatcher) Close() error {
	// close all channels and exit
	close(f.evChan)
	close(f.erChan)
	return nil
}

// Add calls Add on the underlying fsnotify.Watcher
func (f *fileInfoWatcher) Add(name string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// exit early if name already exists
	if _, ok := f.watches[name]; ok {
		return nil
	}

	info, err := f.statFunc(name)
	if err != nil {
		return err
	}

	f.watches[name] = info

	return nil
}

// Remove calls Remove on the underlying fsnotify.Watcher
func (f *fileInfoWatcher) Remove(name string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// no need to exit early, deleting non-existent key is a no-op
	delete(f.watches, name)

	return nil
}

// Watchlist calls watchlist on the underlying fsnotify.Watcher
func (f *fileInfoWatcher) WatchList() []string {
	f.mu.RLock()
	defer f.mu.Unlock()
	out := []string{}
	for name := range f.watches {
		n := name
		out = append(out, n)
	}
	return out
}

// Events returns the underlying watcher's Events chan
func (f *fileInfoWatcher) Events() chan fsnotify.Event {
	return f.evChan
}

// Errors returns the underlying watcher's Errors chan
func (f *fileInfoWatcher) Errors() chan error {
	return f.erChan
}

// Run is a blocking function that starts the filewatcher's timer thread
func (f *fileInfoWatcher) Run(ctx context.Context, s time.Duration) error {
	// timer thread
	f.eg.Go(func() error {
		// execute update on the configured interval of time
		ticker := time.NewTicker(s)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				if err := f.update(); err != nil {
					return err
				}
			}
		}
	})

	return f.eg.Wait()
}

func (f *fileInfoWatcher) update() error {
	event := &fsnotify.Event{}
	f.mu.Lock()
	defer f.mu.Unlock()

	for path, info := range f.watches {
		newInfo, err := getFileInfo(path)
		if err != nil {
			// if the file isn't there, it must have been removed
			// fire off a remove event and remove it from the watches
			if errors.Is(err, os.ErrNotExist) {
				f.evChan <- fsnotify.Event{
					Name: path,
					Op:   fsnotify.Remove,
				}
				delete(f.watches, path)
			}
			return err
		}

		// if the new stat doesn't match the old stat, figure out what changed
		if info != newInfo {
			event, err = f.generateEvent(path, newInfo)
			if err != nil {
				f.erChan <- err
			} else {
				if event != nil {
					f.evChan <- *event
				}
			}
			f.watches[path] = newInfo
		}
	}
	return nil
}

// generateEvent figures out what changed and generates an fsnotify.Event for it. (if we care)
func (f *fileInfoWatcher) generateEvent(path string, newInfo fs.FileInfo) (*fsnotify.Event, error) {
	info := f.watches[path]
	switch {
	// new mod time is more recent than old mod time, generate a write event
	case newInfo.ModTime().After(info.ModTime()):
		return &fsnotify.Event{
			Name: path,
			Op:   fsnotify.Write,
		}, nil
		// the file modes changed, generate a chmod event
	case info.Mode() != newInfo.Mode():
		return &fsnotify.Event{
			Name: path,
			Op:   fsnotify.Chmod,
		}, nil
	}
	return nil, nil
}

// getFileInfo returns the fs.FileInfo for the given path
// TODO: verify this works correctly on windows
func getFileInfo(path string) (fs.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	info, err := f.Stat()
	if err != nil {
		return info, err
	}

	if err := f.Close(); err != nil {
		return info, err
	}

	return info, nil
}
