package service

import (
	"context"
	"fmt"

	"buf.build/gen/go/open-feature/flagd/grpc/go/flagd/sync/v1/syncv1grpc"
	rpc "buf.build/gen/go/open-feature/flagd/grpc/go/sync/v1/syncv1grpc"
	syncv12 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/sync/v1"
	syncv1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/sync/v1"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	"github.com/open-feature/flagd/flagd-proxy/pkg/service/subscriptions"
)

type handler struct {
	syncv1grpc.UnimplementedFlagSyncServiceServer
	syncStore subscriptions.Manager
	logger    *logger.Logger
	// ctx is used to handle SIG[INT|TERM]
	ctx context.Context
}

func (nh *handler) SyncFlags(
	request *syncv12.SyncFlagsRequest,
	server syncv1grpc.FlagSyncService_SyncFlagsServer,
) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errChan := make(chan error)
	dataSync := make(chan sync.DataSync)
	nh.syncStore.RegisterSubscription(ctx, request.GetSelector(), request, dataSync, errChan)
	for {
		select {
		case e := <-errChan:
			return e
		case d := <-dataSync:
			if err := server.Send(&syncv12.SyncFlagsResponse{
				FlagConfiguration: d.FlagData,
			}); err != nil {
				return fmt.Errorf("error sending configuration change event: %w", err)
			}
		case <-server.Context().Done():
			return nil
		case <-nh.ctx.Done():
			return nil
		}
	}
}

func (nh *handler) FetchAllFlags(
	ctx context.Context,
	request *syncv12.FetchAllFlagsRequest,
) (*syncv12.FetchAllFlagsResponse, error) {
	data, err := nh.syncStore.FetchAllFlags(ctx, request, request.GetSelector())
	if err != nil {
		return &syncv12.FetchAllFlagsResponse{}, fmt.Errorf("error fetching all flags from sync store: %w", err)
	}

	return &syncv12.FetchAllFlagsResponse{
		FlagConfiguration: data.FlagData,
	}, nil
}

// oldHandler is the implementation of the old sync schema.
// this will not be required anymore when it is time to work on https://github.com/open-feature/flagd/issues/1088
type oldHandler struct {
	rpc.UnimplementedFlagSyncServiceServer
	syncStore subscriptions.Manager
	logger    *logger.Logger
	// ctx is used to handle SIG[INT|TERM]
	ctx context.Context
}

//nolint:staticcheck
func (l *oldHandler) FetchAllFlags(ctx context.Context, req *syncv1.FetchAllFlagsRequest) (
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

//nolint:staticcheck
func (l *oldHandler) SyncFlags(
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
				return fmt.Errorf("error sending configuration change event: %w", err)
			}
		case <-stream.Context().Done():
			return nil
		case <-l.ctx.Done():
			return nil
		}
	}
}
