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

func updateState(syncr sync.ISync) error {
	msg, err := syncr.Fetch()
	if err != nil {
		return err
	}
	mu.Lock()
	_ = ev.SetState(msg)
	mu.Unlock()
	return nil
}

func Start(ctx context.Context, syncr sync.ISync, server service.IService, evaluator eval.IEvaluator) {
	ev = evaluator
	if err := updateState(syncr); err != nil {
		log.Error(err)
	}

	syncNotifier := make(chan sync.INotify)

	go syncr.Notify(syncNotifier)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-syncNotifier:
				switch w.GetEvent().EventType {
				case sync.EEventTypeCreate:
					log.Info("New configuration created")
					if err := updateState(syncr); err != nil {
						log.Error(err)
					}
				case sync.EEventTypeModify:
					log.Info("Configuration modified")
					if err := updateState(syncr); err != nil {
						log.Error(err)
					}
				case sync.EEventTypeDelete:
					log.Info("Configuration deleted")
				}
			}
		}
	}()

	go func() { _ = server.Serve(ctx, ev) }()
}
