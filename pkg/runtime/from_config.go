package runtime

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/open-feature/flagd/pkg/eval"
	"github.com/open-feature/flagd/pkg/service"
	"github.com/open-feature/flagd/pkg/sync"
	"github.com/open-feature/flagd/pkg/sync/kubernetes"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

func FromConfig(config Config) (*Runtime, error) {
	rt := Runtime{
		config: config,
		Logger: log.WithFields(log.Fields{
			"component": "runtime",
		}),
		syncNotifier: make(chan sync.INotify),
	}
	if err := rt.setEvaluatorFromConfig(); err != nil {
		return nil, err
	}
	if err := rt.setServiceFromConfig(); err != nil {
		return nil, err
	}
	if err := rt.setSyncImplFromConfig(); err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *Runtime) setServiceFromConfig() error {
	switch r.config.ServiceProvider {
	case "http":
		r.Service = &service.HTTPService{
			HTTPServiceConfiguration: &service.HTTPServiceConfiguration{
				Port:             r.config.ServicePort,
				ServerKeyPath:    r.config.ServiceKeyPath,
				ServerCertPath:   r.config.ServiceCertPath,
				ServerSocketPath: r.config.ServiceSocketPath,
			},
			GRPCService: &service.GRPCService{},
			Logger: log.WithFields(log.Fields{
				"service":   "http",
				"component": "service",
			}),
		}
	case "grpc":
		r.Service = &service.GRPCService{
			GRPCServiceConfiguration: &service.GRPCServiceConfiguration{
				Port:             r.config.ServicePort,
				ServerKeyPath:    r.config.ServiceKeyPath,
				ServerCertPath:   r.config.ServiceCertPath,
				ServerSocketPath: r.config.ServiceSocketPath,
			},
			Logger: log.WithFields(log.Fields{
				"service":   "grpc",
				"component": "service",
			}),
		}
	default:
		return errors.New("no service-provider set")
	}
	log.Debugf("Using %s service-provider\n", r.config.ServiceProvider)
	return nil
}

func (r *Runtime) setEvaluatorFromConfig() error {
	switch r.config.Evaluator {
	case "json":
		r.Evaluator = &eval.JSONEvaluator{
			Logger: log.WithFields(log.Fields{
				"evaluator": "json",
				"component": "evaluator",
			}),
		}
	default:
		return errors.New("no evaluator set")
	}
	log.Debugf("Using %s evaluator\n", r.config.Evaluator)
	return nil
}

func (r *Runtime) setSyncImplFromConfig() error {
	r.SyncImpl = make([]sync.ISync, 0, len(r.config.SyncURI))
	switch r.config.SyncProvider {
	case "filepath":
		for _, u := range r.config.SyncURI {
			r.SyncImpl = append(r.SyncImpl, &sync.FilePathSync{
				URI: u,
				Logger: log.WithFields(log.Fields{
					"sync":      "filepath",
					"component": "sync",
				}),
				SyncProviderArgs: r.config.SyncProviderArgs,
			})
			log.Debugf("Using %s sync-provider on %q\n", r.config.SyncProvider, u)
		}
	case "kubernetes":
		r.SyncImpl = append(r.SyncImpl, &kubernetes.KubernetesSync{
			Logger: log.WithFields(log.Fields{
				"sync":      "kubernetes",
				"component": "sync",
			}),
			SyncProviderArgs: r.config.SyncProviderArgs,
		})
		log.Debugf("Using %s sync-provider\n", r.config.SyncProvider)
	case "remote":
		for _, u := range r.config.SyncURI {
			r.SyncImpl = append(r.SyncImpl, &sync.HTTPSync{
				URI:         u,
				BearerToken: r.config.SyncBearerToken,
				Client: &http.Client{
					Timeout: time.Second * 10,
				},
				Logger: log.WithFields(log.Fields{
					"sync":      "remote",
					"component": "sync",
				}),
				SyncProviderArgs: r.config.SyncProviderArgs,
				Cron:             cron.New(),
			})
			log.Debugf("Using %s sync-provider on %q\n", r.config.SyncProvider, u)
		}
	default:
		return fmt.Errorf("invalid sync provider argument: %s", r.config.SyncProvider)
	}
	return nil
}
