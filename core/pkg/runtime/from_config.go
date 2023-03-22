package runtime

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/open-feature/flagd/core/pkg/service"
	"net/http"
	"regexp"
	msync "sync"
	"time"

	"go.opentelemetry.io/otel/exporters/prometheus"

	"github.com/open-feature/flagd/core/pkg/eval"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/otel"
	flageval "github.com/open-feature/flagd/core/pkg/service/flag-evaluation"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/sync"
	"github.com/open-feature/flagd/core/pkg/sync/file"
	"github.com/open-feature/flagd/core/pkg/sync/grpc"
	httpSync "github.com/open-feature/flagd/core/pkg/sync/http"
	"github.com/open-feature/flagd/core/pkg/sync/kubernetes"
	"github.com/robfig/cron"
	"go.uber.org/zap"
)

// from_config is a collection of structures and parsers responsible for deriving flagd runtime

const (
	syncProviderFile       = "file"
	syncProviderGrpc       = "grpc"
	syncProviderKubernetes = "kubernetes"
	syncProviderHTTP       = "http"
	svcName                = "openfeature/flagd"
)

var (
	regCrd        *regexp.Regexp
	regURL        *regexp.Regexp
	regGRPC       *regexp.Regexp
	regGRPCSecure *regexp.Regexp
	regFile       *regexp.Regexp
)

// SourceConfig is configuration option for flagd. This maps to startup parameter sources
type SourceConfig struct {
	URI      string `json:"uri"`
	Provider string `json:"provider"`

	BearerToken string `json:"bearerToken,omitempty"`
	CertPath    string `json:"certPath,omitempty"`
	ProviderID  string `json:"providerID,omitempty"`
	Selector    string `json:"selector,omitempty"`
}

// Config is the configuration structure derived from startup arguments.
type Config struct {
	ServicePort       uint16
	MetricsPort       uint16
	ServiceSocketPath string
	ServiceCertPath   string
	ServiceKeyPath    string

	SyncProviders []SourceConfig
	CORS          []string
}

func init() {
	regCrd = regexp.MustCompile("^core.openfeature.dev/")
	regURL = regexp.MustCompile("^https?://")
	regGRPC = regexp.MustCompile("^" + grpc.Prefix)
	regGRPCSecure = regexp.MustCompile("^" + grpc.PrefixSecure)
	regFile = regexp.MustCompile("^file:")
}

// FromConfig builds a runtime from startup configurations
func FromConfig(logger *logger.Logger, config Config) (*Runtime, error) {
	// build connect service
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	connectService := &flageval.ConnectService{
		ConnectServiceConfiguration: &flageval.ConnectServiceConfiguration{
			ServerKeyPath:    config.ServiceKeyPath,
			ServerCertPath:   config.ServiceCertPath,
			ServerSocketPath: config.ServiceSocketPath,
			CORS:             config.CORS,
		},
		Logger: logger.WithFields(
			zap.String("component", "service"),
		),
		Metrics: otel.NewOTelRecorder(exporter, svcName),
	}

	// build flag store
	s := store.NewFlags()
	sources := []string{}
	for _, sync := range config.SyncProviders {
		sources = append(sources, sync.URI)
	}
	s.FlagSources = sources

	// build sync providers
	syncLogger := logger.WithFields(zap.String("component", "sync"))
	iSyncs, err := syncProvidersFromConfig(syncLogger, config.SyncProviders)
	if err != nil {
		return nil, err
	}

	return &Runtime{
		Logger:    logger.WithFields(zap.String("component", "runtime")),
		Evaluator: eval.NewJSONEvaluator(logger, s),
		Service:   connectService,
		ServiceConfig: service.Configuration{
			Port:        config.ServicePort,
			MetricsPort: config.MetricsPort,
			ServiceName: svcName,
		},
		SyncImpl: iSyncs,
	}, nil
}

// syncProvidersFromConfig is a helper to build ISync implementations from SourceConfig
func syncProvidersFromConfig(logger *logger.Logger, sources []SourceConfig) ([]sync.ISync, error) {
	syncImpls := []sync.ISync{}

	for _, syncProvider := range sources {
		switch syncProvider.Provider {
		case syncProviderFile:
			syncImpls = append(syncImpls, newFile(syncProvider, logger))
			logger.Debug(fmt.Sprintf("using filepath sync-provider for: %q", syncProvider.URI))
		case syncProviderKubernetes:
			k, err := newK8s(syncProvider.URI, logger)
			if err != nil {
				return nil, err
			}
			syncImpls = append(syncImpls, k)
			logger.Debug(fmt.Sprintf("using kubernetes sync-provider for: %s", syncProvider.URI))
		case syncProviderHTTP:
			syncImpls = append(syncImpls, newHTTP(syncProvider, logger))
			logger.Debug(fmt.Sprintf("using remote sync-provider for: %s", syncProvider.URI))
		case syncProviderGrpc:
			syncImpls = append(syncImpls, newGRPC(syncProvider, logger))
			logger.Debug(fmt.Sprintf("using grpc sync-provider for: %s", syncProvider.URI))

		default:
			return nil, fmt.Errorf("invalid sync provider: %s, must be one of with '%s', '%s', '%s' or '%s'",
				syncProvider.Provider, syncProviderFile, syncProviderKubernetes, syncProviderHTTP, syncProviderKubernetes)
		}
	}
	return syncImpls, nil
}

func newGRPC(config SourceConfig, logger *logger.Logger) *grpc.Sync {
	return &grpc.Sync{
		URI: config.URI,
		Logger: logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "grpc"),
		),
		CertPath:   config.CertPath,
		ProviderID: config.ProviderID,
		Selector:   config.Selector,
	}
}

func newHTTP(config SourceConfig, logger *logger.Logger) *httpSync.Sync {
	return &httpSync.Sync{
		URI: config.URI,
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
		Logger: logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "remote"),
		),
		BearerToken: config.BearerToken,
		Cron:        cron.New(),
	}
}

func newK8s(uri string, logger *logger.Logger) (*kubernetes.Sync, error) {
	reader, dynamic, err := kubernetes.GetClients()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewK8sSync(
		logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "kubernetes"),
		),
		regCrd.ReplaceAllString(uri, ""),
		reader,
		dynamic,
	), nil
}

func newFile(config SourceConfig, logger *logger.Logger) *file.Sync {
	return &file.Sync{
		URI: config.URI,
		Logger: logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "filepath"),
		),
		Mux: &msync.RWMutex{},
	}
}

// ParseSources parse a json formatted SourceConfig array string and performs validations on the content
func ParseSources(sourcesFlag string) ([]SourceConfig, error) {
	syncProvidersParsed := []SourceConfig{}

	if err := json.Unmarshal([]byte(sourcesFlag), &syncProvidersParsed); err != nil {
		return syncProvidersParsed, fmt.Errorf("unable to parse sync providers: %w", err)
	}
	for i, sp := range syncProvidersParsed {
		if sp.URI == "" {
			return syncProvidersParsed, errors.New("sync provider argument parse: uri is a required field")
		}
		if sp.Provider == "" {
			return syncProvidersParsed, errors.New("sync provider argument parse: provider is a required field")
		}
		switch uriB := []byte(sp.URI); {
		case regFile.Match(uriB):
			syncProvidersParsed[i].URI = regFile.ReplaceAllString(syncProvidersParsed[i].URI, "")
		case regCrd.Match(uriB):
			syncProvidersParsed[i].URI = regCrd.ReplaceAllString(syncProvidersParsed[i].URI, "")
		}
	}
	return syncProvidersParsed, nil
}

// ParseSyncProviderURIs uri flag based sync sources to SourceConfig array. Replaces uri prefixes where necessary
func ParseSyncProviderURIs(uris []string) ([]SourceConfig, error) {
	syncProvidersParsed := []SourceConfig{}

	for _, uri := range uris {
		switch uriB := []byte(uri); {
		case regFile.Match(uriB):
			syncProvidersParsed = append(syncProvidersParsed, SourceConfig{
				URI:      regFile.ReplaceAllString(uri, ""),
				Provider: syncProviderFile,
			})
		case regCrd.Match(uriB):
			syncProvidersParsed = append(syncProvidersParsed, SourceConfig{
				URI:      regCrd.ReplaceAllString(uri, ""),
				Provider: syncProviderKubernetes,
			})
		case regURL.Match(uriB):
			syncProvidersParsed = append(syncProvidersParsed, SourceConfig{
				URI:      uri,
				Provider: syncProviderHTTP,
			})
		case regGRPC.Match(uriB), regGRPCSecure.Match(uriB):
			syncProvidersParsed = append(syncProvidersParsed, SourceConfig{
				URI:      uri,
				Provider: syncProviderGrpc,
			})
		default:
			return syncProvidersParsed, fmt.Errorf("invalid sync uri argument: %s, must start with 'file:', "+
				"'http(s)://', 'grpc(s)://', or 'core.openfeature.dev'", uri)
		}
	}
	return syncProvidersParsed, nil
}
