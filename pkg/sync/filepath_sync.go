package sync

import (
	"context"
	"errors"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/fsnotify/fsnotify"
)

type FilePathSync struct {
	URI              string
	Logger           *log.Entry
	SyncProviderArgs SyncProviderArgs
}

func (fs *FilePathSync) Fetch(_ context.Context) (string, error) {
	if fs.URI == "" {
		return "", errors.New("no filepath string set")
	}
	rawFile, err := os.ReadFile(fs.URI)
	if err != nil {
		return "", err
	}
	return string(rawFile), nil
}

func (fs *FilePathSync) Notify(ctx context.Context, w chan<- INotify) {
	fs.Logger.Info("Starting filepath sync notifier")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(fs.URI)
	if err != nil {
		fs.Logger.Println(err)
		return
	}

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		fs.Logger.Info("Notifying filepath: ", fs.URI)
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					fs.Logger.Info("Filepath notifier closed")
					return
				}
				var evtType DefaultEventType
				switch event.Op {
				case fsnotify.Create:
					evtType = DefaultEventTypeCreate
				case fsnotify.Write:
					evtType = DefaultEventTypeModify
				case fsnotify.Remove:
					// K8s exposes config maps as symlinks.
					// Updates cause a remove event, we need to re-add the watcher in this case.
					err = watcher.Add(fs.URI)
					if err != nil {
						fs.Logger.Errorf("Error restoring watcher, file may have been deleted: %s", err.Error())
					}
					evtType = DefaultEventTypeDelete
				}
				fs.Logger.Infof("Filepath notifier event: %s %s", event.Name, event.Op.String())
				w <- &Notifier{
					Event: Event[DefaultEventType]{
						EventType: evtType,
					},
				}
				fs.Logger.Info("Filepath notifier event sent")
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fs.Logger.Println("error:", err)
			}
		}
	}()

	w <- &Notifier{Event: Event[DefaultEventType]{DefaultEventTypeReady}} // signal readiness to the caller
	<-ctx.Done()
}
