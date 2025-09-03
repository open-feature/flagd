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
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	"github.com/open-feature/flagd/core/pkg/utils"
	"github.com/robfig/cron"
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3" //nolint:gosec
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type Sync struct {
	uri         string
	client      Client
	cron        Cron
	lastBodySHA string
	logger      *logger.Logger
	bearerToken string
	authHeader  string
	interval    uint32
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
	dataSync <- sync.DataSync{FlagData: msg, Source: hs.uri}
	return nil
}

func (hs *Sync) Init(_ context.Context) error {
	if hs.bearerToken != "" {
		hs.logger.Warn("Deprecation Alert: bearerToken option is deprecated, please use authHeader instead")
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

	hs.logger.Debug(fmt.Sprintf("polling %s every %d seconds", hs.uri, hs.interval))
	_ = hs.cron.AddFunc(fmt.Sprintf("*/%d * * * *", hs.interval), func() {
		hs.logger.Debug(fmt.Sprintf("fetching configuration from %s", hs.uri))
		previousBodySHA := hs.lastBodySHA
		body, noChange, err := hs.fetchBody(ctx, false)
		if err != nil {
			hs.logger.Error(fmt.Sprintf("error fetching: %s", err.Error()))
			return
		}

		if body == "" && !noChange {
			hs.logger.Debug("configuration deleted")
			return
		}

		if previousBodySHA == "" {
			hs.logger.Debug("configuration created")
			dataSync <- sync.DataSync{FlagData: body, Source: hs.uri}
		} else if previousBodySHA != hs.lastBodySHA {
			hs.logger.Debug("configuration updated")
			dataSync <- sync.DataSync{FlagData: body, Source: hs.uri}
		}
	})

	hs.cron.Start()

	dataSync <- sync.DataSync{FlagData: fetch, Source: hs.uri}

	<-ctx.Done()
	hs.cron.Stop()

	return nil
}

func (hs *Sync) fetchBody(ctx context.Context, fetchAll bool) (string, bool, error) {
	if hs.uri == "" {
		return "", false, errors.New("no HTTP URL string set")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", hs.uri, bytes.NewBuffer(nil))
	if err != nil {
		return "", false, fmt.Errorf("error creating request for url %s: %w", hs.uri, err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept", "application/yaml")

	if hs.authHeader != "" {
		req.Header.Set("Authorization", hs.authHeader)
	} else if hs.bearerToken != "" {
		bearer := fmt.Sprintf("Bearer %s", hs.bearerToken)
		req.Header.Set("Authorization", bearer)
	}

	if hs.eTag != "" && !fetchAll {
		req.Header.Set("If-None-Match", hs.eTag)
	}

	resp, err := hs.client.Do(req)
	if err != nil {
		return "", false, fmt.Errorf("error calling endpoint %s: %w", hs.uri, err)
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			hs.logger.Error(fmt.Sprintf("error closing the response body: %s", err.Error()))
		}
	}()

	if resp.StatusCode == 304 {
		hs.logger.Debug("no changes detected")
		return "", true, nil
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		return "", false, fmt.Errorf("error fetching from url %s: %s", hs.uri, resp.Status)
	}

	if resp.Header.Get("ETag") != "" {
		hs.eTag = resp.Header.Get("ETag")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false, fmt.Errorf("unable to read body to bytes: %w", err)
	}

	json, err := utils.ConvertToJSON(body, getFileExtensions(hs.uri), resp.Header.Get("Content-Type"))
	if err != nil {
		return "", false, fmt.Errorf("error converting response body to json: %w", err)
	}

	if json != "" {
		hs.lastBodySHA = hs.generateSha([]byte(body))
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

func NewHTTP(config sync.SourceConfig, logger *logger.Logger) *Sync {
	// Default to 5 seconds
	var interval uint32 = 5
	if config.Interval != 0 {
		interval = config.Interval
	}

	var client *http.Client
	if config.OAuthConfig != nil {
		oauth := clientcredentials.Config{
			ClientID:     config.OAuthConfig.ClientId,
			ClientSecret: config.OAuthConfig.ClientSecret,
			TokenURL:     config.OAuthConfig.TokenUrl,
			AuthStyle:    oauth2.AuthStyleAutoDetect,
		}
		client = oauth.Client(context.Background())
	} else {
		client = &http.Client{
			Timeout: time.Second * 10,
		}
	}

	return &Sync{
		uri:    config.URI,
		client: client,
		logger: logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "http"),
		),
		bearerToken: config.BearerToken,
		authHeader:  config.AuthHeader,
		interval:    interval,
		cron:        cron.New(),
	}
}
