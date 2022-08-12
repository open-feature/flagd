package runtime

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	msync "sync"
	"time"

	"github.com/open-feature/flagd/pkg/eval"
	"github.com/open-feature/flagd/pkg/service"
	"github.com/open-feature/flagd/pkg/sync"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

type Runtime struct {
	config       RuntimeConfig
	Service      service.IService
	SyncImpl     []sync.ISync
	syncNotifier chan sync.INotify
	mu           msync.Mutex
	Evaluator    eval.IEvaluator
	Logger       *log.Entry
}

type RuntimeConfig struct {
	ServiceProvider   string
	ServicePort       int32
	ServiceSocketPath string
	ServiceCertPath   string
	ServiceKeyPath    string

	SyncProvider    string
	SyncUri         []string
	SyncBearerToken string

	Evaluator string
}

func (r *Runtime) SetService() error {
	switch r.config.ServiceProvider {
	case "http":
		r.Service = &service.HTTPService{
			HTTPServiceConfiguration: &service.HTTPServiceConfiguration{
				Port:           r.config.ServicePort,
				ServerKeyPath:  r.config.ServiceKeyPath,
				ServerCertPath: r.config.ServiceCertPath,
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
				Port:           r.config.ServicePort,
				ServerKeyPath:  r.config.ServiceKeyPath,
				ServerCertPath: r.config.ServiceCertPath,
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

func (r *Runtime) SetEvaluator() error {
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

func (r *Runtime) SetSyncImpl() error {
	r.SyncImpl = make([]sync.ISync, 0, len(r.config.SyncUri))
	syncLogger := log.WithFields(log.Fields{
		"sync":      "filepath",
		"component": "sync",
	})
	switch r.config.SyncProvider {
	case "filepath":
		for _, u := range r.config.SyncUri {
			r.SyncImpl = append(r.SyncImpl, &sync.FilePathSync{
				URI:    u,
				Logger: syncLogger,
			})
			log.Debugf("Using %s sync-provider on %q\n", r.config.SyncProvider, u)
		}
	case "remote":
		for _, u := range r.config.SyncUri {
			r.SyncImpl = append(r.SyncImpl, &sync.HTTPSync{
				URI:         u,
				BearerToken: r.config.SyncBearerToken,
				Client: &http.Client{
					Timeout: time.Second * 10,
				},
				Logger: syncLogger,
				Cron:   cron.New(),
			})
			log.Debugf("Using %s sync-provider on %q\n", r.config.SyncProvider, u)
		}
	default:
		return fmt.Errorf("invalid sync provider argument: %s", r.config.SyncProvider)
	}
	return nil
}

func (r *Runtime) startSyncer(ctx context.Context, syncr sync.ISync) error {
	if err := r.updateState(ctx, syncr); err != nil {
		r.Logger.Error(err)
	}

	go syncr.Notify(ctx, r.syncNotifier)

	for {
		select {
		case <-ctx.Done():
			return nil
		case w := <-r.syncNotifier:
			switch w.GetEvent().EventType {
			case sync.DefaultEventTypeCreate:
				r.Logger.Info("New configuration created")
				if err := r.updateState(ctx, syncr); err != nil {
					log.Error(err)
				}
			case sync.DefaultEventTypeModify:
				r.Logger.Info("Configuration modified")
				if err := r.updateState(ctx, syncr); err != nil {
					log.Error(err)
				}
			case sync.DefaultEventTypeDelete:
				r.Logger.Info("Configuration deleted")
			case sync.DefaultEventTypeReady:
				r.Logger.Info("Notifier ready")
			}
		}
	}
}

func (r *Runtime) updateState(ctx context.Context, syncr sync.ISync) error {
	msg, err := syncr.Fetch(ctx)
	if err != nil {
		return fmt.Errorf("fetch: %w", err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	err = r.Evaluator.SetState(msg)
	if err != nil {
		return fmt.Errorf("set state: %w", err)
	}
	return nil
}
