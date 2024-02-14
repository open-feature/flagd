package sync

import (
	"context"

	rpc "buf.build/gen/go/open-feature/flagd/grpc/go/flagd/sync/v1/syncv1grpc"
	syncv1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/sync/v1"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	syncStore "github.com/open-feature/flagd/core/pkg/sync-store"
)

type handler struct {
	rpc.UnimplementedFlagSyncServiceServer
	syncStore *syncStore.SyncStore
	logger    *logger.Logger
}

func (l *handler) FetchAllFlags(ctx context.Context, req *syncv1.FetchAllFlagsRequest) (
	*syncv1.FetchAllFlagsResponse,
	error,
) {
	data, err := l.syncStore.FetchAllFlags(ctx, req, req.GetSelector())
	if err != nil {
		return &syncv1.FetchAllFlagsResponse{}, err
	}

	return &syncv1.FetchAllFlagsResponse{
		FlagConfiguration: data.FlagData,
	}, nil
}

func (l *handler) SyncFlags(
	req *syncv1.SyncFlagsRequest,
	stream rpc.FlagSyncService_SyncFlagsServer,
) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errChan := make(chan error)
	dataSync := make(chan sync.DataSync)
	l.syncStore.RegisterSubscription(ctx, req.GetSelector(), req, dataSync, errChan)
	for {
		select {
		case e := <-errChan:
			return e
		case d := <-dataSync:
			if err := stream.Send(&syncv1.SyncFlagsResponse{
				FlagConfiguration: d.FlagData,
			}); err != nil {
				return err
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}
