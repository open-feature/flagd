package file

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/open-feature/flagd/core/pkg/logger"
)

// Implements file.Watcher using a timer and os.FileInfo
type fileInfoWatcher struct {
	// Event Chan
	evChan chan fsnotify.Event
	// Errors Chan
	erChan chan error
	// logger
	logger *logger.Logger
	// Func to wrap os.Stat (injection point for test helpers)
	statFunc func(string) (fs.FileInfo, error)
	// thread-safe interface to underlying files we are watching
	mu      sync.RWMutex
	watches map[string]fs.FileInfo // filename -> info
	// done signals the timer goroutine to stop and aborts any in-flight send,
	// so Close never closes evChan/erChan while the goroutine is still sending
	done chan struct{}
	// stopOnce guards done so Close is safe to call more than once
	stopOnce sync.Once
	// wg tracks the timer goroutine so Close can wait for it to exit
	wg sync.WaitGroup
}

// NewFsNotifyWatcher returns a new fsNotifyWatcher
func NewFileInfoWatcher(ctx context.Context, logger *logger.Logger) Watcher {
	fiw := &fileInfoWatcher{
		evChan:   make(chan fsnotify.Event, 32),
		erChan:   make(chan error, 32),
		statFunc: getFileInfo,
		logger:   logger,
		watches:  make(map[string]fs.FileInfo),
		done:     make(chan struct{}),
	}
	fiw.run(ctx, (1 * time.Second))
	return fiw
}

// fileInfoWatcher explicitly implements file.Watcher
var _ Watcher = &fileInfoWatcher{}

// Close stops the timer goroutine and closes the event and error channels
func (f *fileInfoWatcher) Close() error {
	// signal the timer goroutine to stop and abort any in-flight send, then wait
	// for it to exit before closing the channels it sends on. Closing before the
	// goroutine has stopped would race with its sends -> send on closed channel.
	f.stopOnce.Do(func() {
		close(f.done)
	})
	f.wg.Wait()
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
	defer f.mu.RUnlock()
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

// run is a blocking function that starts the filewatcher's timer thread
func (f *fileInfoWatcher) run(ctx context.Context, s time.Duration) {
	// timer thread
	f.wg.Add(1)
	go func() {
		defer f.wg.Done()
		// execute update on the configured interval of time
		ticker := time.NewTicker(s)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-f.done:
				return
			case <-ticker.C:
				if err := f.update(); err != nil {
					select {
					case f.erChan <- err:
					case <-f.done:
					}
					return
				}
			}
		}
	}()
}

func (f *fileInfoWatcher) update() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for path, info := range f.watches {
		newInfo, err := f.statFunc(path)
		if err != nil {
			// if the file isn't there, it must have been removed
			// fire off a remove event and remove it from the watches
			if errors.Is(err, os.ErrNotExist) {
				select {
				case f.evChan <- fsnotify.Event{
					Name: path,
					Op:   fsnotify.Remove,
				}:
				case <-f.done:
					return nil
				}
				delete(f.watches, path)
				continue
			}
			return err
		}

		// if the new stat doesn't match the old stat, figure out what changed
		if info != newInfo {
			event := f.generateEvent(path, newInfo)
			if event != nil {
				select {
				case f.evChan <- *event:
				case <-f.done:
					return nil
				}
			}
			f.watches[path] = newInfo
		}
	}
	return nil
}

// generateEvent figures out what changed and generates an fsnotify.Event for it. (if we care)
// file removal are handled above in the update() method
func (f *fileInfoWatcher) generateEvent(path string, newInfo fs.FileInfo) *fsnotify.Event {
	info := f.watches[path]
	switch {
	// new mod time is more recent than old mod time, generate a write event
	case newInfo.ModTime().After(info.ModTime()):
		return &fsnotify.Event{
			Name: path,
			Op:   fsnotify.Write,
		}
		// the file modes changed, generate a chmod event
	case info.Mode() != newInfo.Mode():
		return &fsnotify.Event{
			Name: path,
			Op:   fsnotify.Chmod,
		}
	// nothing changed that we care about
	default:
		return nil
	}
}

// getFileInfo returns the fs.FileInfo for the given path
func getFileInfo(path string) (fs.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error from os.Open(%s): %w", path, err)
	}

	info, err := f.Stat()
	if err != nil {
		return info, fmt.Errorf("error from fs.Stat(%s): %w", path, err)
	}

	if err := f.Close(); err != nil {
		return info, fmt.Errorf("err from fs.Close(%s): %w", path, err)
	}

	return info, nil
}
