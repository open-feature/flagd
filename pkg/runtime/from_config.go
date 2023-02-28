package runtime

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	msync "sync"
	"time"

	"github.com/open-feature/flagd/pkg/eval"
	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/service"
	"github.com/open-feature/flagd/pkg/sync"
	"github.com/open-feature/flagd/pkg/sync/file"
	"github.com/open-feature/flagd/pkg/sync/grpc"
	httpSync "github.com/open-feature/flagd/pkg/sync/http"
	"github.com/open-feature/flagd/pkg/sync/kubernetes"
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
	regCrd  *regexp.Regexp
	regURL  *regexp.Regexp
	regGRPC *regexp.Regexp
	regFile *regexp.Regexp
)

func init() {
	regCrd = regexp.MustCompile("^core.openfeature.dev/")
	regURL = regexp.MustCompile("^https?://")
	regGRPC = regexp.MustCompile("^" + grpc.Prefix)
	regFile = regexp.MustCompile("^file:")
}

func FromConfig(logger *logger.Logger, config Config) (*Runtime, error) {
	rt := Runtime{
		config:    config,
		Logger:    logger.WithFields(zap.String("component", "runtime")),
		Evaluator: eval.NewJSONEvaluator(logger),
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
			r.SyncImpl = append(
				r.SyncImpl,
				r.newK8s(syncProvider, logger),
			)
			rtLogger.Debug(fmt.Sprintf("using kubernetes sync-provider for: %s", syncProvider.URI))
		case syncProviderHTTP:
			r.SyncImpl = append(
				r.SyncImpl,
				r.newHTTP(syncProvider, logger),
			)
			rtLogger.Debug(fmt.Sprintf("using remote sync-provider for: %q", syncProvider.URI))
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

func (r *Runtime) newGRPC(config sync.ProviderConfig, logger *logger.Logger) *grpc.Sync {
	return &grpc.Sync{
		Target: grpc.URLToGRPCTarget(config.URI),
		Logger: logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "grpc"),
		),
		Mux: &msync.RWMutex{},
	}
}

func (r *Runtime) newHTTP(config sync.ProviderConfig, logger *logger.Logger) *httpSync.Sync {
	return &httpSync.Sync{
		URI: config.URI,
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
		Logger: logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "remote"),
		),
		Config: config,
		Cron:   cron.New(),
	}
}

func (r *Runtime) newK8s(config sync.ProviderConfig, logger *logger.Logger) *kubernetes.Sync {
	return &kubernetes.Sync{
		Logger: logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "kubernetes"),
		),
		URI:    config.URI,
		Config: config,
	}
}

func (r *Runtime) newFile(config sync.ProviderConfig, logger *logger.Logger) *file.Sync {
	return &file.Sync{
		URI: config.URI,
		Logger: logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "filepath"),
		),
		Config: config,
		Mux:    &msync.RWMutex{},
	}
}

func SyncProviderArgPass(syncProviders string) ([]sync.ProviderConfig, error) {
	syncProvidersParsed := []sync.ProviderConfig{}
	if err := json.Unmarshal([]byte(syncProviders), &syncProvidersParsed); err != nil {
		return syncProvidersParsed, fmt.Errorf("unable to parse sync providers: %w", err)
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

func SyncProvidersFromURIs(uris []string) ([]sync.ProviderConfig, error) {
	syncProvidersParsed := []sync.ProviderConfig{}
	for _, uri := range uris {
		switch uriB := []byte(uri); {
		case regFile.Match(uriB):
			syncProvidersParsed = append(syncProvidersParsed, sync.ProviderConfig{
				URI:      regFile.ReplaceAllString(uri, ""),
				Provider: syncProviderFile,
			})
		case regCrd.Match(uriB):
			syncProvidersParsed = append(syncProvidersParsed, sync.ProviderConfig{
				URI:      regCrd.ReplaceAllString(uri, ""),
				Provider: syncProviderKubernetes,
			})
		case regURL.Match(uriB):
			syncProvidersParsed = append(syncProvidersParsed, sync.ProviderConfig{
				URI:      uri,
				Provider: syncProviderHTTP,
			})
		case regGRPC.Match(uriB):
			syncProvidersParsed = append(syncProvidersParsed, sync.ProviderConfig{
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
