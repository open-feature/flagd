package runtime

import (
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
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/sync"
	"github.com/open-feature/flagd/core/pkg/sync/file"
	"github.com/open-feature/flagd/core/pkg/sync/grpc"
	httpSync "github.com/open-feature/flagd/core/pkg/sync/http"
	"github.com/open-feature/flagd/core/pkg/sync/kubernetes"
	"github.com/robfig/cron"
	"go.uber.org/zap"
)

const (
	syncProviderFile       = "file"
	syncProviderGrpc       = "grpc"
	syncProviderKubernetes = "kubernetes"
	syncProviderHTTP       = "http"
)

var (
	regCrd        *regexp.Regexp
	regURL        *regexp.Regexp
	regGRPC       *regexp.Regexp
	regGRPCSecure *regexp.Regexp
	regFile       *regexp.Regexp
)

func init() {
	regCrd = regexp.MustCompile("^core.openfeature.dev/")
	regURL = regexp.MustCompile("^https?://")
	regGRPC = regexp.MustCompile("^" + grpc.Prefix)
	regGRPCSecure = regexp.MustCompile("^" + grpc.PrefixSecure)
	regFile = regexp.MustCompile("^file:")
}

func FromConfig(logger *logger.Logger, config Config) (*Runtime, error) {
	s := store.NewFlags()
	sources := []string{}
	for _, sync := range config.SyncProviders {
		sources = append(sources, sync.URI)
	}
	s.FlagSources = sources
	rt := Runtime{
		config:    config,
		Logger:    logger.WithFields(zap.String("component", "runtime")),
		Evaluator: eval.NewJSONEvaluator(logger, s),
	}
	if err := rt.setSyncImplFromConfig(logger); err != nil {
		return nil, err
	}
	rt.setService(logger)
	return &rt, nil
}

func (r *Runtime) setService(logger *logger.Logger) {
	r.Service = &service.ConnectService{
		ConnectServiceConfiguration: &service.ConnectServiceConfiguration{
			Port:             r.config.ServicePort,
			MetricsPort:      r.config.MetricsPort,
			ServerKeyPath:    r.config.ServiceKeyPath,
			ServerCertPath:   r.config.ServiceCertPath,
			ServerSocketPath: r.config.ServiceSocketPath,
			CORS:             r.config.CORS,
		},
		Logger: logger.WithFields(
			zap.String("component", "service"),
		),
	}
}

func (r *Runtime) setSyncImplFromConfig(logger *logger.Logger) error {
	rtLogger := logger.WithFields(zap.String("component", "runtime"))
	r.SyncImpl = make([]sync.ISync, 0, len(r.config.SyncProviders))
	for _, syncProvider := range r.config.SyncProviders {
		switch syncProvider.Provider {
		case syncProviderFile:
			r.SyncImpl = append(
				r.SyncImpl,
				r.newFile(syncProvider, logger),
			)
			rtLogger.Debug(fmt.Sprintf("using filepath sync-provider for: %q", syncProvider.URI))
		case syncProviderKubernetes:
			k, err := r.newK8s(syncProvider.URI, logger)
			if err != nil {
				return err
			}
			r.SyncImpl = append(
				r.SyncImpl,
				k,
			)
			rtLogger.Debug(fmt.Sprintf("using kubernetes sync-provider for: %s", syncProvider.URI))
		case syncProviderHTTP:
			r.SyncImpl = append(
				r.SyncImpl,
				r.newHTTP(syncProvider, logger),
			)
			rtLogger.Debug(fmt.Sprintf("using remote sync-provider for: %s", syncProvider.URI))
		case syncProviderGrpc:
			r.SyncImpl = append(
				r.SyncImpl,
				r.newGRPC(syncProvider, logger),
			)
		default:
			return fmt.Errorf("invalid sync uri argument: %s, must start with 'file:', 'http(s)://', 'grpc://',"+
				" or 'core.openfeature.dev'", syncProvider.URI)
		}
	}
	return nil
}

func (r *Runtime) newGRPC(config sync.SourceConfig, logger *logger.Logger) *grpc.Sync {
	return &grpc.Sync{
		URI: config.URI,
		Logger: logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "grpc"),
		),
		CertPath: config.CertPath,
	}
}

func (r *Runtime) newHTTP(config sync.SourceConfig, logger *logger.Logger) *httpSync.Sync {
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

func (r *Runtime) newK8s(uri string, logger *logger.Logger) (*kubernetes.Sync, error) {
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

func (r *Runtime) newFile(config sync.SourceConfig, logger *logger.Logger) *file.Sync {
	return &file.Sync{
		URI: config.URI,
		Logger: logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "filepath"),
		),
		Mux: &msync.RWMutex{},
	}
}

func SyncProviderArgParse(syncProviders string) ([]sync.SourceConfig, error) {
	syncProvidersParsed := []sync.SourceConfig{}
	if err := json.Unmarshal([]byte(syncProviders), &syncProvidersParsed); err != nil {
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

func SyncProvidersFromURIs(uris []string) ([]sync.SourceConfig, error) {
	syncProvidersParsed := []sync.SourceConfig{}
	for _, uri := range uris {
		switch uriB := []byte(uri); {
		case regFile.Match(uriB):
			syncProvidersParsed = append(syncProvidersParsed, sync.SourceConfig{
				URI:      regFile.ReplaceAllString(uri, ""),
				Provider: syncProviderFile,
			})
		case regCrd.Match(uriB):
			syncProvidersParsed = append(syncProvidersParsed, sync.SourceConfig{
				URI:      regCrd.ReplaceAllString(uri, ""),
				Provider: syncProviderKubernetes,
			})
		case regURL.Match(uriB):
			syncProvidersParsed = append(syncProvidersParsed, sync.SourceConfig{
				URI:      uri,
				Provider: syncProviderHTTP,
			})
		case regGRPC.Match(uriB), regGRPCSecure.Match(uriB):
			syncProvidersParsed = append(syncProvidersParsed, sync.SourceConfig{
				URI:      uri,
				Provider: syncProviderGrpc,
			})
		default:
			return syncProvidersParsed, fmt.Errorf("invalid sync uri argument: %s, must start with 'file:', "+
				"'http(s)://', 'grpc://', or 'core.openfeature.dev'", uri)
		}
	}
	return syncProvidersParsed, nil
}
