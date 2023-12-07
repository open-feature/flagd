package sync

import (
	"context"
	"fmt"

	rpc "buf.build/gen/go/open-feature/flagd/grpc/go/sync/v1/syncv1grpc"
	syncv1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/sync/v1"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/subscriptions"
	"github.com/open-feature/flagd/core/pkg/sync"
)

type handler struct {
	rpc.UnimplementedFlagSyncServiceServer
	syncStore subscriptions.Manager
	logger    *logger.Logger
}

func (l *handler) FetchAllFlags(ctx context.Context, req *syncv1.FetchAllFlagsRequest) (
	*syncv1.FetchAllFlagsResponse,
	error,
) {
	data, err := l.syncStore.FetchAllFlags(ctx, req, req.GetSelector())
	if err != nil {
		return &syncv1.FetchAllFlagsResponse{}, fmt.Errorf("error fetching all flags from sync store: %w", err)
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
				State:             dataSyncToGrpcState(d),
			}); err != nil {
				return fmt.Errorf("error sending configuration change event: %w", err)
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}

func dataSyncToGrpcState(s sync.DataSync) syncv1.SyncState {
	return syncv1.SyncState(s.Type + 1)
}
