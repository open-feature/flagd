package sync

import (
	"bytes"
	"context"
	"crypto/sha1" //nolint:gosec
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

type HTTPSync struct {
	URI         string
	Client      *http.Client
	BearerToken string
	LastBodySHA string
}

func (fs *HTTPSync) fetchBodyFromURL(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, bytes.NewBuffer(nil))
	if err != nil {
		return []byte(""), err
	}

	req.Header.Add("Accept", "application/json")

	if fs.BearerToken != "" {
		bearer := "Bearer " + fs.BearerToken
		req.Header.Set("Authorization", bearer)
	}

	resp, err := fs.Client.Do(req)
	defer func() { _ = resp.Body.Close() }()
	if err != nil {
		return []byte(""), err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}

	return body, nil
}

func (fs *HTTPSync) generateSha(body []byte) string {
	hasher := sha1.New() //nolint:gosec
	hasher.Write(body)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

func (fs *HTTPSync) Fetch(ctx context.Context) (string, error) {
	if fs.URI == "" {
		return "", errors.New("no HTTP URL string set")
	}

	body, err := fs.fetchBodyFromURL(ctx, fs.URI)
	if len(body) != 0 {
		fs.LastBodySHA = fs.generateSha(body)
	}
	return string(body), err
}

func (fs *HTTPSync) Notify(ctx context.Context, w chan<- INotify) {
	c := cron.New()

	_ = c.AddFunc("*/5 * * * *", func() {
		body, err := fs.fetchBodyFromURL(ctx, fs.URI)
		if err != nil {
			log.Error(err)
			return
		}
		if len(body) == 0 {
			w <- &Notifier{
				Event: Event{
					EventType: EEventTypeDelete,
				},
			}
		} else {
			if fs.LastBodySHA == "" {
				w <- &Notifier{
					Event: Event{
						EventType: EEventTypeCreate,
					},
				}
			} else {
				currentSHA := fs.generateSha(body)
				if fs.LastBodySHA != currentSHA {
					log.Infof("http notifier event: %s has been modified", fs.URI)
					w <- &Notifier{
						Event: Event{
							EventType: EEventTypeModify,
						},
					}
				}
				fs.LastBodySHA = currentSHA
			}
		}
	})
	c.Start()
}
