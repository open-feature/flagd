package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	parseUrl "net/url"
	"path/filepath"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	"github.com/open-feature/flagd/core/pkg/utils"
	"golang.org/x/crypto/sha3" //nolint:gosec
)

type Sync struct {
	URI         string
	Client      Client
	Cron        Cron
	LastBodySHA string
	Logger      *logger.Logger
	BearerToken string
	AuthHeader  string
	Interval    uint32
	ready       bool
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

func (hs *Sync) ReSync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	msg, err := hs.Fetch(ctx)
	if err != nil {
		return err
	}
	dataSync <- sync.DataSync{FlagData: msg, Source: hs.URI}
	return nil
}

func (hs *Sync) Init(_ context.Context) error {
	if hs.BearerToken != "" {
		hs.Logger.Warn("Deprecation Alert: bearerToken option is deprecated, please use authHeader instead")
	}
	return nil
}

func (hs *Sync) IsReady() bool {
	return hs.ready
}

func (hs *Sync) Sync(ctx context.Context, dataSync chan<- sync.DataSync) error {
	// Initial fetch
	fetch, err := hs.Fetch(ctx)
	if err != nil {
		return err
	}

	// Set ready state
	hs.ready = true

	hs.Logger.Debug(fmt.Sprintf("polling %s every %d seconds", hs.URI, hs.Interval))
	_ = hs.Cron.AddFunc(fmt.Sprintf("*/%d * * * *", hs.Interval), func() {
		hs.Logger.Debug(fmt.Sprintf("fetching configuration from %s", hs.URI))
		body, err := hs.fetchBodyFromURL(ctx, hs.URI)
		if err != nil {
			hs.Logger.Error(err.Error())
			return
		}

		if body == "" {
			hs.Logger.Debug("configuration deleted")
			return
		}

		currentSHA := hs.generateSha([]byte(body))

		if hs.LastBodySHA == "" {
			hs.Logger.Debug("new configuration created")
			dataSync <- sync.DataSync{FlagData: body, Source: hs.URI}
		} else if hs.LastBodySHA != currentSHA {
			hs.Logger.Debug("configuration modified")
			dataSync <- sync.DataSync{FlagData: body, Source: hs.URI}
		}

		hs.LastBodySHA = currentSHA
	})

	hs.Cron.Start()

	dataSync <- sync.DataSync{FlagData: fetch, Source: hs.URI}

	<-ctx.Done()
	hs.Cron.Stop()

	return nil
}

func (hs *Sync) fetchBodyFromURL(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, bytes.NewBuffer(nil))
	if err != nil {
		return "", fmt.Errorf("error creating request for url %s: %w", url, err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept", "application/yaml")

	if hs.AuthHeader != "" {
		req.Header.Set("Authorization", hs.AuthHeader)
	} else if hs.BearerToken != "" {
		bearer := fmt.Sprintf("Bearer %s", hs.BearerToken)
		req.Header.Set("Authorization", bearer)
	}

	resp, err := hs.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error calling endpoint %s: %w", url, err)
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			hs.Logger.Debug(fmt.Sprintf("error closing the response body: %s", err.Error()))
		}
	}()

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		return "", fmt.Errorf("error fetching from url %s: %s", url, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read body to bytes: %w", err)
	}

	json, err := utils.ConvertToJSON(body, getFileExtensions(url), resp.Header.Get("Content-Type"))
	if err != nil {
		return "", fmt.Errorf("error converting response body to json: %w", err)
	}
	return json, nil
}

// getFileExtensions returns the file extension from the URL path
func getFileExtensions(url string) string {
	u, err := parseUrl.Parse(url)
	if err != nil {
		return ""
	}

	return filepath.Ext(u.Path)
}

func (hs *Sync) generateSha(body []byte) string {
	hasher := sha3.New256()
	hasher.Write(body)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func (hs *Sync) Fetch(ctx context.Context) (string, error) {
	if hs.URI == "" {
		return "", errors.New("no HTTP URL string set")
	}

	body, err := hs.fetchBodyFromURL(ctx, hs.URI)
	if err != nil {
		return "", err
	}
	if body != "" {
		hs.LastBodySHA = hs.generateSha([]byte(body))
	}

	return body, nil
}
