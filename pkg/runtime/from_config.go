package runtime

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/open-feature/flagd/pkg/eval"
	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/service"
	"github.com/open-feature/flagd/pkg/sync"
	"github.com/open-feature/flagd/pkg/sync/kubernetes"
	"github.com/robfig/cron"
	"go.uber.org/zap"
)

const (
	crdRegex  = "^core.openfeature.dev/"
	fileRegex = "^file://"
	urlRegex  = "^https?://"
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
	regCrd, err := regexp.Compile(crdRegex)
	if err != nil {
		return err
	}
	regUrl, err := regexp.Compile(urlRegex)
	if err != nil {
		return err
	}
	regFile, err := regexp.Compile(fileRegex)
	if err != nil {
		return err
	}
	rtLogger := logger.WithFields(zap.String("component", "runtime"))
	r.SyncImpl = make([]sync.ISync, 0, len(r.config.SyncURI))
	for _, uri := range r.config.SyncURI {
		if regFile.Match([]byte(uri)) {
			r.SyncImpl = append(r.SyncImpl, &sync.FilePathSync{
				URI: regFile.ReplaceAllString(uri, ""),
				Logger: logger.WithFields(
					zap.String("component", "sync"),
					zap.String("sync", "filepath"),
				),
				ProviderArgs: r.config.ProviderArgs,
			})
			rtLogger.Debug(fmt.Sprintf("Using filepath sync-provider for %q", uri))
		} else if regCrd.Match([]byte(uri)) {
			r.SyncImpl = append(r.SyncImpl, &kubernetes.Sync{
				Logger: logger.WithFields(
					zap.String("component", "sync"),
					zap.String("sync", "kubernetes"),
				),
				URI:          regCrd.ReplaceAllString(uri, ""),
				ProviderArgs: r.config.ProviderArgs,
			})
			rtLogger.Debug(fmt.Sprintf("Using kubernetes sync-provider for %s", uri))
		} else if regUrl.Match([]byte(uri)) {
			r.SyncImpl = append(r.SyncImpl, &sync.HTTPSync{
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
			})
			rtLogger.Debug(fmt.Sprintf("Using remote sync-provider for %q", uri))
		} else {
			return fmt.Errorf("invalid sync uri argument: %s", uri)
		}
	}
	return nil
}
