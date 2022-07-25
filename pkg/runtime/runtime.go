package runtime

import (
	"context"
	msync "sync"

	"github.com/open-feature/flagd/pkg/eval"
	"github.com/open-feature/flagd/pkg/service"
	"github.com/open-feature/flagd/pkg/sync"
	log "github.com/sirupsen/logrus"
)

var (
	mu msync.Mutex
	ev eval.IEvaluator
)

func updateState(ctx context.Context, syncr sync.ISync) error {
	msg, err := syncr.Fetch(ctx)
	if err != nil {
		return err
	}
	mu.Lock()
	_ = ev.SetState(msg)
	mu.Unlock()
	return nil
}

func startSyncer(ctx context.Context, notifier chan sync.INotify, syncr sync.ISync) {
	if err := updateState(ctx, syncr); err != nil {
		log.Error(err)
	}

	go syncr.Notify(ctx, notifier)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-notifier:
				switch w.GetEvent().EventType {
				case sync.EEventTypeCreate:
					log.Info("New configuration created")
					if err := updateState(ctx, syncr); err != nil {
						log.Error(err)
					}
				case sync.EEventTypeModify:
					log.Info("Configuration modified")
					if err := updateState(ctx, syncr); err != nil {
						log.Error(err)
					}
				case sync.EEventTypeDelete:
					log.Info("Configuration deleted")
				}
			}
		}
	}()
}

func Start(ctx context.Context, syncr []sync.ISync, server service.IService, evaluator eval.IEvaluator) {
	ev = evaluator

	syncNotifier := make(chan sync.INotify)

	for _, s := range syncr {
		startSyncer(ctx, syncNotifier, s)
	}

	go func() { _ = server.Serve(ctx, ev) }()
}
