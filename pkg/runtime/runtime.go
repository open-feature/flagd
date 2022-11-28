package runtime

import (
	"context"
	"fmt"
	msync "sync"

	"github.com/open-feature/flagd/pkg/eval"
	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/service"
	"github.com/open-feature/flagd/pkg/sync"
)

type Runtime struct {
	config       Config
	Service      service.IService
	SyncImpl     []sync.ISync
	syncNotifier chan sync.INotify

	mu        msync.Mutex
	Evaluator eval.IEvaluator
	Logger    *logger.Logger
}

type Config struct {
	ServicePort       int32
	MetricsPort       int32
	ServiceSocketPath string
	ServiceCertPath   string
	ServiceKeyPath    string

	SyncProvider    string
	ProviderArgs    sync.ProviderArgs
	SyncURI         []string
	SyncBearerToken string

	Evaluator string
	CORS      []string
}

func (r *Runtime) startSyncer(ctx context.Context, syncr sync.ISync) error {
	if err := r.updateState(ctx, syncr); err != nil {
		return err
	}

	go syncr.Notify(ctx, r.syncNotifier)

	for {
		select {
		case <-ctx.Done():
			return nil
		case w := <-r.syncNotifier:
			switch w.GetEvent().EventType {
			case sync.DefaultEventTypeCreate:
				r.Logger.Debug("New configuration created")
				if err := r.updateState(ctx, syncr); err != nil {
					r.Logger.Error(err.Error())
				}
			case sync.DefaultEventTypeModify:
				r.Logger.Debug("Configuration modified")
				if err := r.updateState(ctx, syncr); err != nil {
					r.Logger.Error(err.Error())
				}
			case sync.DefaultEventTypeDelete:
				r.Logger.Debug("Configuration deleted")
			case sync.DefaultEventTypeReady:
				r.Logger.Debug("Notifier ready")
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
	notifications, err := r.Evaluator.SetState(syncr.Source(), msg)
	if err != nil {
		return fmt.Errorf("set state: %w", err)
	}
	for _, n := range notifications {
		r.Logger.Info(fmt.Sprintf("configuration change (%s) for flagKey %s (%s)", n.Type, n.FlagKey, n.Source))
		r.Service.Notify(service.Notification{
			Type: service.ConfigurationChange,
			Data: n.ToMap(),
		})
	}
	return nil
}
