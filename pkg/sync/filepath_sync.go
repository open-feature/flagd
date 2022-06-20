package sync

import (
	"errors"
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"github.com/fsnotify/fsnotify"
)

type FilePathSync struct {
	URI string
}

func (fs *FilePathSync) Fetch() (string, error) {
	if fs.URI == "" {
		return "", errors.New("no filepath string set")
	}
	rawFile, err := ioutil.ReadFile(fs.URI)
	if err != nil {
		return "", err
	}
	return string(rawFile), nil
}

func (fs *FilePathSync) Notify(w chan<- INotify) {
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
				var evtType E_EVENT_TYPE
				switch event.Op {
				case fsnotify.Create:
					evtType = E_EVENT_TYPE_CREATE
				case fsnotify.Write:
					evtType = E_EVENT_TYPE_MODIFY
				case fsnotify.Remove:
					// K8s exposes config maps as symlinks.
					// Updates cause a remove event, we need to re-add the watcher in this case.
					err = watcher.Add(fs.URI)
					if (err != nil) {
						log.Println("Error restoring watcher:", err)
						log.Fatal(err)
					}
					evtType = E_EVENT_TYPE_DELETE
				}
				log.Infof("Filepath notifier event: %s %s", event.Name, event.Op.String())
				w <- &Notifier{
					Event: Event{
						EventType: evtType,
					},
				}
				log.Info("Filepath notifier event sent")
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	err = watcher.Add(fs.URI)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
