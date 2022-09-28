package runtime

import (
	"context"
	"fmt"
	msync "sync"

	"github.com/open-feature/flagd/pkg/eval"
	"github.com/open-feature/flagd/pkg/service"
	"github.com/open-feature/flagd/pkg/sync"
	log "github.com/sirupsen/logrus"
)

type Runtime struct {
	config       Config
	Service      service.IService
	SyncImpl     []sync.ISync
	syncNotifier chan sync.INotify
	mu           msync.Mutex
	Evaluator    eval.IEvaluator
	Logger       *log.Entry
}

type Config struct {
	ServicePort       int32
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
					log.Error(err)
				}
			case sync.DefaultEventTypeModify:
				r.Logger.Debug("Configuration modified")
				if err := r.updateState(ctx, syncr); err != nil {
					log.Error(err)
				}
				r.Service.Notify(service.ConfigurationChange)
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
	err = r.Evaluator.SetState(msg)
	if err != nil {
		return fmt.Errorf("set state: %w", err)
	}
	return nil
}
