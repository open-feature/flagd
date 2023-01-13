package http

import (
	"bytes"
	"context"
	"crypto/sha1" //nolint:gosec
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/open-feature/flagd/pkg/sync"
	"io"
	"net/http"

	"github.com/open-feature/flagd/pkg/logger"
)

type Sync struct {
	URI          string
	Client       HTTPClient
	Cron         Cron
	BearerToken  string
	LastBodySHA  string
	Logger       *logger.Logger
	ProviderArgs sync.ProviderArgs
}

// HTTPClient defines the behaviour required of a http client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Cron defines the behaviour required of a cron
type Cron interface {
	AddFunc(spec string, cmd func()) error
	Start()
	Stop()
}

func (fs *Sync) Source() string {
	return fs.URI
}

func (fs *Sync) fetchBodyFromURL(ctx context.Context, url string) ([]byte, error) {
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
	if err != nil {
		return []byte(""), err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}

	return body, nil
}

func (fs *Sync) generateSha(body []byte) string {
	hasher := sha1.New() //nolint:gosec
	hasher.Write(body)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

func (fs *Sync) Fetch(ctx context.Context) (string, error) {
	if fs.URI == "" {
		return "", errors.New("no HTTP URL string set")
	}

	body, err := fs.fetchBodyFromURL(ctx, fs.URI)
	if err != nil {
		return "", err
	}
	if len(body) != 0 {
		fs.LastBodySHA = fs.generateSha(body)
	}
	return string(body), nil
}

func (fs *Sync) Notify(ctx context.Context, w chan<- sync.INotify) {
	_ = fs.Cron.AddFunc("*/5 * * * *", func() {
		body, err := fs.fetchBodyFromURL(ctx, fs.URI)
		if err != nil {
			fs.Logger.Error(err.Error())
			return
		}
		if len(body) == 0 {
			w <- &sync.Notifier{
				Event: sync.Event[sync.DefaultEventType]{
					sync.DefaultEventTypeDelete,
				},
			}
		} else {
			if fs.LastBodySHA == "" {
				w <- &sync.Notifier{
					Event: sync.Event[sync.DefaultEventType]{
						sync.DefaultEventTypeCreate,
					},
				}
			} else {
				currentSHA := fs.generateSha(body)
				if fs.LastBodySHA != currentSHA {
					fs.Logger.Info(fmt.Sprintf("http notifier event: %s has been modified", fs.URI))
					w <- &sync.Notifier{
						Event: sync.Event[sync.DefaultEventType]{
							sync.DefaultEventTypeModify,
						},
					}
				}
				fs.LastBodySHA = currentSHA
			}
		}
	})
	w <- &sync.Notifier{
		Event: sync.Event[sync.DefaultEventType]{
			sync.DefaultEventTypeReady,
		},
	}
	
	fs.Cron.Start()
	<-ctx.Done()
	fs.Cron.Stop()
}
