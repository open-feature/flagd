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
	"os"
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

	oauthCredential *oauthCredentialHandler
}

type oauthCredentialHandler struct {
	clientId           string
	clientSecret       string
	tokenUrl           string
	folderSource       string
	reloadDelaySeconds int

	lastUpdate time.Time
}

func (och *oauthCredentialHandler) loadOAuthConfiguration() (*clientcredentials.Config, error) {
	if och.folderSource == "" {
		oauth := &clientcredentials.Config{
			ClientID:     och.clientId,
			ClientSecret: och.clientSecret,
			TokenURL:     och.tokenUrl,
			AuthStyle:    oauth2.AuthStyleInParams,
		}
		return oauth, nil
	}
	// we load from files
	id, err := os.ReadFile(och.folderSource + "/" + och.clientId)
	if err != nil {
		return nil, err
	}

	secret, err := os.ReadFile(och.folderSource + "/" + och.clientSecret)
	if err != nil {
		return nil, err
	}
	return &clientcredentials.Config{
		ClientID:     string(id),
		ClientSecret: string(secret),
		TokenURL:     och.tokenUrl,
		AuthStyle:    oauth2.AuthStyleInParams,
	}, nil
}

func (och *oauthCredentialHandler) getOAuthClient() (*http.Client, error) {
	oauth, err := och.loadOAuthConfiguration()
	if err != nil {
		return nil, err
	}
	return oauth.Client(context.Background()), nil
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
	client := hs.getClient()
	resp, err := client.Do(req)
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

func (hs *Sync) getClient() Client {
	if hs.oauthCredential != nil {
		delta := time.Since(hs.oauthCredential.lastUpdate).Seconds()
		if delta <= float64(hs.oauthCredential.reloadDelaySeconds) && hs.client != nil {
			return hs.client
		}
	} else if hs.client != nil {
		return hs.client
	}
	// if we cannot reuse the cached one, let's create a new one
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	if hs.oauthCredential != nil {
		if c, err := hs.oauthCredential.getOAuthClient(); err == nil {
			client = c
			hs.oauthCredential.lastUpdate = time.Now()
		} else {
			hs.logger.Error(fmt.Sprintf("Cannot init OAuth. Default to normal HTTP client: %v", err))
		}
	}
	hs.client = client
	return client
}

func NewHTTP(config sync.SourceConfig, logger *logger.Logger) *Sync {
	// Default to 5 seconds
	var interval uint32 = 5
	if config.Interval != 0 {
		interval = config.Interval
	}

	var oauthCredential *oauthCredentialHandler
	if config.OAuthConfig != nil {
		oauthCredential = &oauthCredentialHandler{
			clientId:           config.OAuthConfig.ClientId,
			clientSecret:       config.OAuthConfig.ClientSecret,
			tokenUrl:           config.OAuthConfig.TokenUrl,
			folderSource:       config.OAuthConfig.Folder,
			reloadDelaySeconds: config.OAuthConfig.ReloadDelayS,
			lastUpdate:         time.Now(),
		}
	}

	return &Sync{
		uri: config.URI,
		logger: logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "http"),
		),
		bearerToken:     config.BearerToken,
		authHeader:      config.AuthHeader,
		interval:        interval,
		cron:            cron.New(),
		oauthCredential: oauthCredential,
	}
}
