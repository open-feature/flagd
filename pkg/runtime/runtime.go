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

var (
	mu msync.Mutex
	ev eval.IEvaluator
)

func updateState(ctx context.Context, syncr sync.ISync) error {
	msg, err := syncr.Fetch(ctx)
	if err != nil {
		return fmt.Errorf("fetch: %w", err)
	}
	mu.Lock()
	defer mu.Unlock()
	err = ev.SetState(msg)
	if err != nil {
		return fmt.Errorf("set state: %w", err)
	}
	return nil
}

func startSyncer(ctx context.Context, notifier chan sync.INotify, syncr sync.ISync, logger *log.Entry) {
	if err := updateState(ctx, syncr); err != nil {
		logger.Error(err)
	}

	go syncr.Notify(ctx, notifier)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-notifier:
				switch w.GetEvent().EventType {
				case sync.DefaultEventTypeCreate:
					logger.Info("New configuration created")
					if err := updateState(ctx, syncr); err != nil {
						log.Error(err)
					}
				case sync.DefaultEventTypeModify:
					logger.Info("Configuration modified")
					if err := updateState(ctx, syncr); err != nil {
						log.Error(err)
					}
				case sync.DefaultEventTypeDelete:
					logger.Info("Configuration deleted")
				}
			}
		}
	}()
}

func Start(ctx context.Context, syncr []sync.ISync, server service.IService,
	evaluator eval.IEvaluator, logger *log.Entry,
) {
	ev = evaluator

	syncNotifier := make(chan sync.INotify)

	for _, s := range syncr {
		startSyncer(ctx, syncNotifier, s, logger)
	}

	go func() { _ = server.Serve(ctx, ev) }()
}
