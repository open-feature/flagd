package http

import (
	"bytes"
	"context"
	"crypto/sha1" //nolint:gosec
	"encoding/base64"
	"errors"
	"io"
	"net/http"

	"github.com/open-feature/flagd/pkg/sync"

	"github.com/open-feature/flagd/pkg/logger"
)

type Sync struct {
	URI          string
	Client       Client
	Cron         Cron
	BearerToken  string
	LastBodySHA  string
	Logger       *logger.Logger
	ProviderArgs sync.ProviderArgs
}

// Client defines the behaviour required of a http client
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

// Cron defines the behaviour required of a cron
type Cron interface {
	AddFunc(spec string, cmd func()) error
	Start()
	Stop()
}

func (hs *Sync) Sync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	// Initial fetch
	fetch, err := hs.Fetch(ctx)
	if err != nil {
		hs.Logger.Error(err.Error())
		return err
	}

	dataSync <- sync.DataSync{FlagData: fetch, Source: hs.URI}

	_ = hs.Cron.AddFunc("*/5 * * * *", func() {
		body, err := hs.fetchBodyFromURL(ctx, hs.URI)
		if err != nil {
			hs.Logger.Error(err.Error())
			return
		}
		if len(body) == 0 {
			hs.Logger.Debug("Configuration deleted")
		} else {
			if hs.LastBodySHA == "" {
				hs.Logger.Debug("New configuration created")
				msg, err := hs.Fetch(ctx)
				if err != nil {
					hs.Logger.Error(err.Error())
				}
				dataSync <- sync.DataSync{
					FlagData: msg,
					Source:   hs.URI,
				}
			} else {
				currentSHA := hs.generateSha(body)
				if hs.LastBodySHA != currentSHA {
					hs.Logger.Debug("Configuration modified")
					msg, err := hs.Fetch(ctx)
					if err != nil {
						hs.Logger.Error(err.Error())
					}
					dataSync <- sync.DataSync{
						FlagData: msg,
						Source:   hs.URI,
					}
				}
				hs.LastBodySHA = currentSHA
			}
		}
	})

	hs.Cron.Start()
	<-ctx.Done()
	hs.Cron.Stop()

	return nil
}

func (hs *Sync) fetchBodyFromURL(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, bytes.NewBuffer(nil))
	if err != nil {
		return []byte(""), err
	}

	req.Header.Add("Accept", "application/json")

	if hs.BearerToken != "" {
		bearer := "Bearer " + hs.BearerToken
		req.Header.Set("Authorization", bearer)
	}

	resp, err := hs.Client.Do(req)
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

func (hs *Sync) generateSha(body []byte) string {
	hasher := sha1.New() //nolint:gosec
	hasher.Write(body)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

func (hs *Sync) Fetch(ctx context.Context) (string, error) {
	if hs.URI == "" {
		return "", errors.New("no HTTP URL string set")
	}

	body, err := hs.fetchBodyFromURL(ctx, hs.URI)
	if err != nil {
		return "", err
	}
	if len(body) != 0 {
		hs.LastBodySHA = hs.generateSha(body)
	}
	return string(body), nil
}
