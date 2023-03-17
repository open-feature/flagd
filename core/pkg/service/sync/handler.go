package sync

import (
	"context"

	syncv1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/sync/v1"
	"github.com/bufbuild/connect-go"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	syncStore "github.com/open-feature/flagd/core/pkg/sync-store"
)

type handler struct {
	syncStore *syncStore.SyncStore
	logger    *logger.Logger
}

func (l *handler) FetchAllFlags(ctx context.Context, req *connect.Request[syncv1.FetchAllFlagsRequest]) (
	*connect.Response[syncv1.FetchAllFlagsResponse],
	error,
) {
	data, err := l.syncStore.FetchAllFlags(ctx, req, req.Msg.GetSelector())
	if err != nil {
		return connect.NewResponse(&syncv1.FetchAllFlagsResponse{}), err
	}

	return connect.NewResponse(&syncv1.FetchAllFlagsResponse{
		FlagConfiguration: data.FlagData,
	}), nil
}

func (l *handler) SyncFlags(
	ctx context.Context,
	req *connect.Request[syncv1.SyncFlagsRequest],
	stream *connect.ServerStream[syncv1.SyncFlagsResponse],
) error {
	errChan := make(chan error)
	dataSync := make(chan sync.DataSync)
	l.syncStore.RegisterSubscription(ctx, req.Msg.GetSelector(), req, dataSync, errChan)
	for {
		select {
		case e := <-errChan:
			return e
		case d := <-dataSync:
			if err := stream.Send(&syncv1.SyncFlagsResponse{
				FlagConfiguration: d.FlagData,
				State:             dataSyncToGrpcState(d),
			}); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func dataSyncToGrpcState(s sync.DataSync) syncv1.SyncState {
	return syncv1.SyncState(s.Type + 1)
}
