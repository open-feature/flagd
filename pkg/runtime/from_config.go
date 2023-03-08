package runtime

import (
	"fmt"
	"net/http"
	"regexp"
	msync "sync"
	"time"

	"github.com/open-feature/flagd/pkg/eval"
	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/service"
	"github.com/open-feature/flagd/pkg/store"
	"github.com/open-feature/flagd/pkg/sync"
	"github.com/open-feature/flagd/pkg/sync/file"
	"github.com/open-feature/flagd/pkg/sync/grpc"
	httpSync "github.com/open-feature/flagd/pkg/sync/http"
	"github.com/open-feature/flagd/pkg/sync/kubernetes"
	"github.com/robfig/cron"
	"go.uber.org/zap"
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
	s := store.NewFlags()
	s.FlagSources = config.SyncURI
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
	r.SyncImpl = make([]sync.ISync, 0, len(r.config.SyncURI))
	for _, uri := range r.config.SyncURI {
		switch uriB := []byte(uri); {
		case regFile.Match(uriB):
			r.SyncImpl = append(
				r.SyncImpl,
				r.newFile(uri, logger),
			)
			rtLogger.Debug(fmt.Sprintf("using filepath sync-provider for: %q", uri))
		case regCrd.Match(uriB):
			k, err := r.newK8s(uri, logger)
			if err != nil {
				return err
			}
			r.SyncImpl = append(
				r.SyncImpl,
				k,
			)
			rtLogger.Debug(fmt.Sprintf("using kubernetes sync-provider for: %s", uri))
		case regURL.Match(uriB):
			r.SyncImpl = append(
				r.SyncImpl,
				r.newHTTP(uri, logger),
			)
			rtLogger.Debug(fmt.Sprintf("using remote sync-provider for: %q", uri))
		case regGRPC.Match(uriB):
			r.SyncImpl = append(
				r.SyncImpl,
				r.newGRPC(uri, logger),
			)
		default:
			return fmt.Errorf("invalid sync uri argument: %s, must start with 'file:', 'http(s)://', 'grpc://',"+
				" or 'core.openfeature.dev'", uri)
		}
	}
	return nil
}

func (r *Runtime) newGRPC(uri string, logger *logger.Logger) *grpc.Sync {
	return &grpc.Sync{
		Target: grpc.URLToGRPCTarget(uri),
		Logger: logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "grpc"),
		),
	}
}

func (r *Runtime) newHTTP(uri string, logger *logger.Logger) *httpSync.Sync {
	return &httpSync.Sync{
		URI:         uri,
		BearerToken: r.config.SyncBearerToken,
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
		Logger: logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "remote"),
		),
		ProviderArgs: r.config.ProviderArgs,
		Cron:         cron.New(),
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
		r.config.ProviderArgs,
		reader,
		dynamic,
	), nil
}

func (r *Runtime) newFile(uri string, logger *logger.Logger) *file.Sync {
	return &file.Sync{
		URI: regFile.ReplaceAllString(uri, ""),
		Logger: logger.WithFields(
			zap.String("component", "sync"),
			zap.String("sync", "filepath"),
		),
		ProviderArgs: r.config.ProviderArgs,
		Mux:          &msync.RWMutex{},
	}
}
