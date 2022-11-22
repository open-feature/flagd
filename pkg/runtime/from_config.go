package runtime

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/open-feature/flagd/pkg/eval"
	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/service"
	"github.com/open-feature/flagd/pkg/sync"
	"github.com/open-feature/flagd/pkg/sync/kubernetes"
	"github.com/robfig/cron"
	"go.uber.org/zap"
)

func FromConfig(logger *logger.Logger, config Config) (*Runtime, error) {
	rt := Runtime{
		config:       config,
		Logger:       logger.WithFields(zap.String("component", "runtime")),
		syncNotifier: make(chan sync.INotify),
	}
	if err := rt.setEvaluatorFromConfig(logger); err != nil {
		return nil, err
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

func (r *Runtime) setEvaluatorFromConfig(logger *logger.Logger) error {
	switch r.config.Evaluator {
	case "json":
		r.Evaluator = eval.NewJSONEvaluator(logger)
	default:
		return errors.New("no evaluator set")
	}
	logger.Debug(fmt.Sprintf("Using %s evaluator", r.config.Evaluator))
	return nil
}

func (r *Runtime) setSyncImplFromConfig(logger *logger.Logger) error {
	rtLogger := logger.WithFields(zap.String("component", "runtime"))
	r.SyncImpl = make([]sync.ISync, 0, len(r.config.SyncURI))
	switch r.config.SyncProvider {
	case "filepath":
		for _, u := range r.config.SyncURI {
			r.SyncImpl = append(r.SyncImpl, &sync.FilePathSync{
				URI: u,
				Logger: logger.WithFields(
					zap.String("component", "sync"),
					zap.String("sync", "filepath"),
				),
				ProviderArgs: r.config.ProviderArgs,
			})
			rtLogger.Debug(fmt.Sprintf("Using %s sync-provider on %q", r.config.SyncProvider, u))
		}
	case "kubernetes":
		r.SyncImpl = append(r.SyncImpl, &kubernetes.Sync{
			Logger: logger.WithFields(
				zap.String("component", "sync"),
				zap.String("sync", "kubernetes"),
			),
			ProviderArgs: r.config.ProviderArgs,
		})
		rtLogger.Debug(fmt.Sprintf("Using %s sync-provider", r.config.SyncProvider))
	case "remote":
		for _, u := range r.config.SyncURI {
			r.SyncImpl = append(r.SyncImpl, &sync.HTTPSync{
				URI:         u,
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
			})
			rtLogger.Debug(fmt.Sprintf("Using %s sync-provider on %q", r.config.SyncProvider, u))
		}
	default:
		return fmt.Errorf("invalid sync provider argument: %s", r.config.SyncProvider)
	}
	return nil
}
