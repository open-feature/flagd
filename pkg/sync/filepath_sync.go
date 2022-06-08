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

func (fs *FilePathSync) Watch(w chan IWatcher) {
	log.Info("Starting filepath sync watcher")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		defer close(done)
		log.Info("Watching filepath: ", fs.URI)
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Info("Filepath watcher closed")
					return
				}
				var evtType E_EVENT_TYPE
				switch event.Op {
				case fsnotify.Create:
					evtType = E_EVENT_TYPE_CREATE
				case fsnotify.Write:
					evtType = E_EVENT_TYPE_MODIFY
				case fsnotify.Remove:
					evtType = E_EVENT_TYPE_DELETE
				}
				log.Infof("Filepath watcher event: %s %s", event.Name, event.Op.String())
				w <- &Watcher{
					Event: Event{
						EventType: evtType,
					},
				}
				log.Info("Filepath watcher event sent")
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
