package sync

import (
	"context"
	"errors"
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"github.com/fsnotify/fsnotify"
)

type FilePathSync struct {
	URI string
}

func (fs *FilePathSync) Fetch(_ context.Context) (string, error) {
	if fs.URI == "" {
		return "", errors.New("no filepath string set")
	}
	rawFile, err := ioutil.ReadFile(fs.URI)
	if err != nil {
		return "", err
	}
	return string(rawFile), nil
}

func (fs *FilePathSync) Notify(ctx context.Context, w chan<- INotify) {
	log.Info("Starting filepath sync notifier")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		defer close(done)
		log.Info("Notifying filepath: ", fs.URI)
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Info("Filepath notifier closed")
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
						log.Fatalf("Error restoring watcher: %s, exiting...", err.Error())
					}
					evtType = DefaultEventTypeDelete
				}
				log.Infof("Filepath notifier event: %s %s", event.Name, event.Op.String())
				w <- &Notifier{
					Event: Event[DefaultEventType]{
						EventType: evtType,
					},
				}
				log.Info("Filepath notifier event sent")
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			case <-ctx.Done():
				return
			}
		}
	}()
	err = watcher.Add(fs.URI)
	if err != nil {
		log.Println(err)
		return
	}
	<-done
}
