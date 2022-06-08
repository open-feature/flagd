package runtime

import (
	"context"
	msync "sync"

	"github.com/open-feature/flagd/pkg/service"
	"github.com/open-feature/flagd/pkg/sync"
	log "github.com/sirupsen/logrus"
)

var (
	mu          msync.Mutex
	syncPayload string
)

func getSyncPayload() string {
	mu.Lock()
	s := syncPayload
	mu.Unlock()
	return s
}

func setSyncPayload(syncr sync.ISync) error {
	msg, err := syncr.Fetch()
	if err != nil {
		return err
	}
	mu.Lock()
	syncPayload = msg
	mu.Unlock()
	return nil
}

func Start(syncr sync.ISync, server service.IService, ctx context.Context) {

	if err := setSyncPayload(syncr); err != nil {
		log.Error(err)
	}

	syncWatcher := make(chan sync.IWatcher)

	go syncr.Watch(syncWatcher)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-syncWatcher:
				switch w.GetEvent().EventType {
				case sync.E_EVENT_TYPE_CREATE:
					log.Info("New configuration created")
					if err := setSyncPayload(syncr); err != nil {
						log.Error(err)
					}
				case sync.E_EVENT_TYPE_MODIFY:
					log.Info("Configuration modified")
					if err := setSyncPayload(syncr); err != nil {
						log.Error(err)
					}
				case sync.E_EVENT_TYPE_DELETE:
					log.Info("Configuration deleted")
				}
			}
		}
	}()

	go server.Serve(
		func(ir service.IServiceRequest) service.IServiceResponse {
			if ir.GetRequestType() == service.SERVICE_REQUEST_ALL_FLAGS {
				return ir.GenerateServiceResponse(getSyncPayload())
			}
			return nil
		})

}
