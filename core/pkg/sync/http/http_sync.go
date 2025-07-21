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
	eTag        string
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
	msg, _, err := hs.fetchBody(ctx, true)
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
	fetch, _, err := hs.fetchBody(ctx, true)
	if err != nil {
		return err
	}

	// Set ready state
	hs.ready = true

	hs.Logger.Debug(fmt.Sprintf("polling %s every %d seconds", hs.URI, hs.Interval))
	_ = hs.Cron.AddFunc(fmt.Sprintf("*/%d * * * *", hs.Interval), func() {
		hs.Logger.Debug(fmt.Sprintf("fetching configuration from %s", hs.URI))
		previousBodySHA := hs.LastBodySHA
		body, noChange, err := hs.fetchBody(ctx, false)
		if err != nil {
			hs.Logger.Error(fmt.Sprintf("error fetching: %s", err.Error()))
			return
		}

		if body == "" && !noChange {
			hs.Logger.Debug("configuration deleted")
			return
		}

		if previousBodySHA == "" {
			hs.Logger.Debug("configuration created")
			dataSync <- sync.DataSync{FlagData: body, Source: hs.URI}
		} else if previousBodySHA != hs.LastBodySHA {
			hs.Logger.Debug("configuration updated")
			dataSync <- sync.DataSync{FlagData: body, Source: hs.URI}
		}
	})

	hs.Cron.Start()

	dataSync <- sync.DataSync{FlagData: fetch, Source: hs.URI}

	<-ctx.Done()
	hs.Cron.Stop()

	return nil
}

func (hs *Sync) fetchBody(ctx context.Context, fetchAll bool) (string, bool, error) {
	if hs.URI == "" {
		return "", false, errors.New("no HTTP URL string set")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", hs.URI, bytes.NewBuffer(nil))
	if err != nil {
		return "", false, fmt.Errorf("error creating request for url %s: %w", hs.URI, err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept", "application/yaml")

	if hs.AuthHeader != "" {
		req.Header.Set("Authorization", hs.AuthHeader)
	} else if hs.BearerToken != "" {
		bearer := fmt.Sprintf("Bearer %s", hs.BearerToken)
		req.Header.Set("Authorization", bearer)
	}

	if hs.eTag != "" && !fetchAll {
		req.Header.Set("If-None-Match", hs.eTag)
	}

	resp, err := hs.Client.Do(req)
	if err != nil {
		return "", false, fmt.Errorf("error calling endpoint %s: %w", hs.URI, err)
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			hs.Logger.Error(fmt.Sprintf("error closing the response body: %s", err.Error()))
		}
	}()

	if resp.StatusCode == 304 {
		hs.Logger.Debug("no changes detected")
		return "", true, nil
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		return "", false, fmt.Errorf("error fetching from url %s: %s", hs.URI, resp.Status)
	}

	if resp.Header.Get("ETag") != "" {
		hs.eTag = resp.Header.Get("ETag")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false, fmt.Errorf("unable to read body to bytes: %w", err)
	}

	json, err := utils.ConvertToJSON(body, getFileExtensions(hs.URI), resp.Header.Get("Content-Type"))
	if err != nil {
		return "", false, fmt.Errorf("error converting response body to json: %w", err)
	}

	if json != "" {
		hs.LastBodySHA = hs.generateSha([]byte(body))
	}

	return json, false, nil
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
	body, _, err := hs.fetchBody(ctx, false)
	return body, err
}
