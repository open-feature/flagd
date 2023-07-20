package runtime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	msync "sync"
	"time"

	"github.com/open-feature/flagd/core/pkg/eval"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/service"
	flageval "github.com/open-feature/flagd/core/pkg/service/flag-evaluation"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/sync"
	"github.com/open-feature/flagd/core/pkg/sync/file"
	"github.com/open-feature/flagd/core/pkg/sync/grpc"
	"github.com/open-feature/flagd/core/pkg/sync/grpc/credentials"
	httpSync "github.com/open-feature/flagd/core/pkg/sync/http"
	"github.com/open-feature/flagd/core/pkg/sync/kubernetes"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	"github.com/robfig/cron"
	"go.uber.org/zap"
)

// from_config is a collection of structures and parsers responsible for deriving flagd runtime

const (
	syncProviderFile       = "file"
	syncProviderGrpc       = "grpc"
	syncProviderKubernetes = "kubernetes"
	syncProviderHTTP       = "http"
	svcName                = "flagd"
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
	TLS         bool   `json:"tls,omitempty"`
	ProviderID  string `json:"providerID,omitempty"`
	Selector    string `json:"selector,omitempty"`
}

// Config is the configuration structure derived from startup arguments.
type Config struct {
	MetricExporter    string
	MetricsPort       uint16
	OtelCollectorURI  string
	ServiceCertPath   string
	ServiceKeyPath    string
	ServicePort       uint16
	ServiceSocketPath string

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
func FromConfig(logger *logger.Logger, version string, config Config) (*Runtime, error) {
	telCfg := telemetry.Config{
		MetricsExporter: config.MetricExporter,
		CollectorTarget: config.OtelCollectorURI,
	}

	// register error handling for OpenTelemetry
	telemetry.RegisterErrorHandling(logger)

	// register trace provider for the runtime
	err := telemetry.BuildTraceProvider(context.Background(), logger, svcName, version, telCfg)
	if err != nil {
		return nil, fmt.Errorf("error building trace provider: %w", err)
	}

	// build metrics recorder with startup configurations
	recorder, err := telemetry.BuildMetricsRecorder(context.Background(), svcName, version, telCfg)
	if err != nil {
		return nil, fmt.Errorf("error building metrics recorder: %w", err)
	}

	// build flag store
	s := store.NewFlags()
	sources := []store.SourceDetails{}
	for _, sync := range config.SyncProviders {
		sources = append(sources, store.SourceDetails{
			Source:   sync.URI,
			Selector: sync.Selector,
		})
	}
	s.FlagSources = sources

	// derive evaluator
	evaluator := setupJSONEvaluator(logger, s)

	// derive service
	connectService := flageval.NewConnectService(
		logger.WithFields(zap.String("component", "service")),
		evaluator,
		recorder)

	// build sync providers
	syncLogger := logger.WithFields(zap.String("component", "sync"))
	iSyncs, err := syncProvidersFromConfig(syncLogger, config.SyncProviders)
	if err != nil {
		return nil, err
	}

	return &Runtime{
		Logger:    logger.WithFields(zap.String("component", "runtime")),
		Evaluator: evaluator,
		Service:   connectService,
		ServiceConfig: service.Configuration{
			Port:        config.ServicePort,
			MetricsPort: config.MetricsPort,
			ServiceName: svcName,
			KeyPath:     config.ServiceKeyPath,
			CertPath:    config.ServiceCertPath,
			SocketPath:  config.ServiceSocketPath,
			CORS:        config.CORS,
			Options:     telemetry.BuildConnectOptions(telCfg),
		},
		SyncImpl: iSyncs,
	}, nil
}

func setupJSONEvaluator(logger *logger.Logger, s *store.Flags) *eval.JSONEvaluator {
	evaluator := eval.NewJSONEvaluator(
		logger,
		s,
		eval.WithEvaluator(
			"fractionalEvaluation",
			eval.NewFractionalEvaluator(logger).FractionalEvaluation,
		),
		eval.WithEvaluator(
			"starts_with",
			eval.NewStringComparisonEvaluator(logger).StartsWithEvaluation,
		),
		eval.WithEvaluator(
			"ends_with",
			eval.NewStringComparisonEvaluator(logger).EndsWithEvaluation,
		),
		eval.WithEvaluator(
			"sem_ver",
			eval.NewSemVerComparisonEvaluator(logger).SemVerEvaluation,
		),
	)
	return evaluator
}

// syncProvidersFromConfig is a helper to build ISync implementations from SourceConfig
func syncProvidersFromConfig(logger *logger.Logger, sources []SourceConfig) ([]sync.ISync, error) {
	syncImpls := []sync.ISync{}

	for _, syncProvider := range sources {
		switch syncProvider.Provider {
		case syncProviderFile:
			syncImpls = append(syncImpls, NewFile(syncProvider, logger))
			logger.Debug(fmt.Sprintf("using filepath sync-provider for: %q", syncProvider.URI))
		case syncProviderKubernetes:
			k, err := NewK8s(syncProvider.URI, logger)
			if err != nil {
				return nil, err
			}
			syncImpls = append(syncImpls, k)
			logger.Debug(fmt.Sprintf("using kubernetes sync-provider for: %s", syncProvider.URI))
		case syncProviderHTTP:
			syncImpls = append(syncImpls, NewHTTP(syncProvider, logger))
			logger.Debug(fmt.Sprintf("using remote sync-provider for: %s", syncProvider.URI))
		case syncProviderGrpc:
			syncImpls = append(syncImpls, NewGRPC(syncProvider, logger))
			logger.Debug(fmt.Sprintf("using grpc sync-provider for: %s", syncProvider.URI))

		default:
			return nil, fmt.Errorf("invalid sync provider: %s, must be one of with '%s', '%s', '%s' or '%s'",
				syncProvider.Provider, syncProviderFile, syncProviderKubernetes, syncProviderHTTP, syncProviderKubernetes)
		}
	}
	return syncImpls, nil
}

func NewGRPC(config SourceConfig, logger *logger.Logger) *grpc.Sync {
	return &grpc.Sync{
		URI: config.URI,
		Logger: logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "grpc"),
		),
		CredentialBuilder: &credentials.CredentialBuilder{},
		CertPath:          config.CertPath,
		ProviderID:        config.ProviderID,
		Secure:            config.TLS,
		Selector:          config.Selector,
	}
}

func NewHTTP(config SourceConfig, logger *logger.Logger) *httpSync.Sync {
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

func NewK8s(uri string, logger *logger.Logger) (*kubernetes.Sync, error) {
	reader, dynamic, err := kubernetes.GetClients()
	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes clients: %w", err)
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

func NewFile(config SourceConfig, logger *logger.Logger) *file.Sync {
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
		return syncProvidersParsed, fmt.Errorf("error parsing sync providers: %w", err)
	}
	for _, sp := range syncProvidersParsed {
		if sp.URI == "" {
			return syncProvidersParsed, errors.New("sync provider argument parse: uri is a required field")
		}
		if sp.Provider == "" {
			return syncProvidersParsed, errors.New("sync provider argument parse: provider is a required field")
		}
	}
	return syncProvidersParsed, nil
}

// ParseSyncProviderURIs uri flag based sync sources to SourceConfig array. Replaces uri prefixes where necessary to
// derive SourceConfig
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
		case regGRPC.Match(uriB):
			syncProvidersParsed = append(syncProvidersParsed, SourceConfig{
				URI:      regGRPC.ReplaceAllString(uri, ""),
				Provider: syncProviderGrpc,
			})
		case regGRPCSecure.Match(uriB):
			syncProvidersParsed = append(syncProvidersParsed, SourceConfig{
				URI:      regGRPCSecure.ReplaceAllString(uri, ""),
				Provider: syncProviderGrpc,
				TLS:      true,
			})
		default:
			return syncProvidersParsed, fmt.Errorf("invalid sync uri argument: %s, must start with 'file:', "+
				"'http(s)://', 'grpc(s)://', or 'core.openfeature.dev'", uri)
		}
	}
	return syncProvidersParsed, nil
}
