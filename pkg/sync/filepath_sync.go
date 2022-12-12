package sync

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/fsnotify/fsnotify"
	"github.com/open-feature/flagd/pkg/logger"
)

const (
	FileTypeYAML = "yaml"
	FileTypeJSON = "json"
	FileTypeYML  = "yml"
)

type FilePathSync struct {
	URI          string
	Logger       *logger.Logger
	ProviderArgs ProviderArgs
	// FileType indicates the file type e.g., json, yaml/yml etc.,
	FileType string
}

func (fs *FilePathSync) Source() string {
	return fs.URI
}

func (fs *FilePathSync) Fetch(_ context.Context) (string, error) {
	if fs.URI == "" {
		return "", errors.New("no filepath string set")
	}
	rawFile, err := os.ReadFile(fs.URI)
	if err != nil {
		return "", err
	}

	switch fs.FileType {
	case FileTypeYAML:
		fallthrough
	case FileTypeYML:
		return yamlToJSON(rawFile)
	case FileTypeJSON:
		return string(rawFile), nil
	default:
		return "", fmt.Errorf("filepath extension '%v' is not supported", fs.FileType)
	}
}

func (fs *FilePathSync) Notify(ctx context.Context, w chan<- INotify) {
	fs.Logger.Info("Starting filepath sync notifier")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fs.Logger.Fatal(err.Error())
	}
	defer watcher.Close()

	err = watcher.Add(fs.URI)
	if err != nil {
		fs.Logger.Error(err.Error())
		return
	}

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		fs.Logger.Info(fmt.Sprintf("Notifying filepath: %s", fs.URI))
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
						fs.Logger.Error(fmt.Sprintf("Error restoring watcher, file may have been deleted: %s", err.Error()))
					}
					evtType = DefaultEventTypeDelete
				}
				fs.Logger.Info(fmt.Sprintf("Filepath notifier event: %s %s", event.Name, event.Op.String()))
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
				fs.Logger.Error(err.Error())
			}
		}
	}()

	w <- &Notifier{Event: Event[DefaultEventType]{DefaultEventTypeReady}} // signal readiness to the caller
	<-ctx.Done()
}

// yamlToJSON is a generic helper function to convert
// yaml to json
func yamlToJSON(rawFile []byte) (string, error) {
	var ms map[string]interface{}
	// yaml.Unmarshal unmarshals to map[interface]interface{}
	if err := yaml.Unmarshal(rawFile, &ms); err != nil {
		return "", fmt.Errorf("unmarshal yaml: %w", err)
	}

	// Adding spaces here because our evaluator transposer function
	// doesn't understand json without indentations quite well
	r, err := json.MarshalIndent(ms, "", "  ")
	if err != nil {
		return "", fmt.Errorf("convert yaml to json: %w", err)
	}

	return string(r), err
}
