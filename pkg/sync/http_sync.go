package sync

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

type HttpSync struct {
	URI         string
	Client      *http.Client
	BearerToken string
	LastBodySHA string
}

func (fs *HttpSync) fetchBodyFromURL(url string) ([]byte, error) {

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(nil))
	if err != nil {
		return []byte(""), err
	}

	req.Header.Add("Accept", "application/json")

	if fs.BearerToken != "" {
		bearer := "Bearer " + fs.BearerToken
		req.Header.Set("Authorization", bearer)
	}

	resp, err := fs.Client.Do(req)
	if err != nil {
		return []byte(""), err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}

	return body, nil
}
func (fs *HttpSync) generateSha(body []byte) string {
	hasher := sha1.New()
	hasher.Write(body)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}
func (fs *HttpSync) Fetch() (string, error) {
	if fs.URI == "" {
		return "", errors.New("no filepath string set")
	}

	body, err := fs.fetchBodyFromURL(fs.URI)
	if len(body) != 0 {
		fs.LastBodySHA = fs.generateSha(body)
	}
	return string(body), err
}

func (fs *HttpSync) Watch(w chan IWatcher) {

	c := cron.New()
	/*
		This initial implementation uses a cron to poll the remote endpoint.
		Whilst not a true watch stream, it will give similar functionality to the fsnotify watcher.
	*/
	c.AddFunc("*/5 * * * *", func() {
		body, err := fs.fetchBodyFromURL(fs.URI)
		if err != nil {
			log.Error(err)
			return
		}
		if len(body) == 0 {
			w <- &Watcher{
				Event: Event{
					EventType: E_EVENT_TYPE_DELETE,
				},
			}
		} else {
			if fs.LastBodySHA == "" {
				w <- &Watcher{
					Event: Event{
						EventType: E_EVENT_TYPE_CREATE,
					},
				}
			} else {
				currentSHA := fs.generateSha(body)
				if fs.LastBodySHA != currentSHA {
					log.Infof("Old hash: %s New Hash: %s", fs.LastBodySHA, currentSHA)
					log.Infof("http watcher event: %s has been modified", fs.URI)
					w <- &Watcher{
						Event: Event{
							EventType: E_EVENT_TYPE_MODIFY,
						},
					}
				}
				fs.LastBodySHA = currentSHA
			}
		}
	})
	c.Start()

}
