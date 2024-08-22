package file

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
)

// Implements file.Watcher by wrapping fsnotify.Watcher
// This is only necessary because fsnotify.Watcher directly exposes its Errors
// and Events channels rather than returning them by method invocation
type fsNotifyWatcher struct {
	watcher *fsnotify.Watcher
}

// NewFsNotifyWatcher returns a new fsNotifyWatcher
func NewFSNotifyWatcher() (Watcher, error) {
	fsn, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("fsnotify: %w", err)
	}
	return &fsNotifyWatcher{
		watcher: fsn,
	}, nil
}

// explicitly implements file.Watcher
var _ Watcher = &fsNotifyWatcher{}

// Close calls close on the underlying fsnotify.Watcher
func (f *fsNotifyWatcher) Close() error {
	if err := f.watcher.Close(); err != nil {
		return fmt.Errorf("fsnotify: %w", err)
	}
	return nil
}

// Add calls Add on the underlying fsnotify.Watcher
func (f *fsNotifyWatcher) Add(name string) error {
	if err := f.watcher.Add(name); err != nil {
		return fmt.Errorf("fsnotify: %w", err)
	}
	return nil
}

// Remove calls Remove on the underlying fsnotify.Watcher
func (f *fsNotifyWatcher) Remove(name string) error {
	if err := f.watcher.Remove(name); err != nil {
		return fmt.Errorf("fsnotify: %w", err)
	}
	return nil
}

// Watchlist calls watchlist on the underlying fsnotify.Watcher
func (f *fsNotifyWatcher) WatchList() []string {
	return f.watcher.WatchList()
}

// Events returns the underlying watcher's Events chan
func (f *fsNotifyWatcher) Events() chan fsnotify.Event {
	return f.watcher.Events
}

// Errors returns the underlying watcher's Errors chan
func (f *fsNotifyWatcher) Errors() chan error {
	return f.watcher.Errors
}
