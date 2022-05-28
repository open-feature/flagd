package runtime

import (
	"context"

	"github.com/open-feature/flagd/pkg/service"
	"github.com/open-feature/flagd/pkg/sync"
	log "github.com/sirupsen/logrus"
)

func Start(syncr sync.ISync, server service.IService, ctx context.Context) {

	// This is a very simple example of how the interface can be used for service and sync
	// The service interface will serve requests whilst the sync interface is responsible
	// for refreshing the configuration data
	messageBuffer := make(chan string)
	requestBuffer := make(chan struct{})
  
	go server.Serve(
		func(ir service.IServiceRequest) service.IServiceResponse {
			if ir.GetRequestType() == service.SERVICE_REQUEST_ALL_FLAGS {
				requestBuffer <- struct{}{}
				return ir.GenerateServiceResponse(<-messageBuffer)
			}
			return nil
		})

	for {
		select {
		case <-ctx.Done():
			log.Info("Runtime context has been cancelled")
			return
		case <-requestBuffer:
			data, err := syncr.Fetch()
			if err != nil {
				log.Warn(err.Error())
			}
			messageBuffer <- data
		}
	}
}
