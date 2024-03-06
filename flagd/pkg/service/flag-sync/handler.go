package sync

import (
	"context"
	"fmt"

	"buf.build/gen/go/open-feature/flagd/grpc/go/flagd/sync/v1/syncv1grpc"
	"buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/sync/v1"
	"github.com/open-feature/flagd/core/pkg/logger"
	"google.golang.org/protobuf/types/known/structpb"
)

// syncHandler implements the sync contract
type syncHandler struct {
	mux *Multiplexer
	log *logger.Logger
}

func (s syncHandler) SyncFlags(req *syncv1.SyncFlagsRequest, server syncv1grpc.FlagSyncService_SyncFlagsServer) error {
	muxPayload := make(chan payload, 1)
	selector := req.GetSelector()

	ctx := server.Context()

	err := s.mux.Register(ctx, selector, muxPayload)
	if err != nil {
		return err
	}

	for {
		select {
		case payload := <-muxPayload:
			err := server.Send(&syncv1.SyncFlagsResponse{FlagConfiguration: payload.flags})
			if err != nil {
				s.log.Debug(fmt.Sprintf("error sending stream response: %v", err))
				return fmt.Errorf("error sending stream response: %w", err)
			}
		case <-ctx.Done():
			s.mux.Unregister(ctx, selector)
			s.log.Debug("context done, exiting stream request")
			return nil
		}
	}
}

func (s syncHandler) FetchAllFlags(_ context.Context, req *syncv1.FetchAllFlagsRequest) (
	*syncv1.FetchAllFlagsResponse, error,
) {
	flags, err := s.mux.GetALlFlags(req.GetSelector())
	if err != nil {
		return nil, err
	}

	return &syncv1.FetchAllFlagsResponse{
		FlagConfiguration: flags,
	}, nil
}

func (s syncHandler) GetMetadata(_ context.Context, _ *syncv1.GetMetadataRequest) (
	*syncv1.GetMetadataResponse, error,
) {
	metadata, err := structpb.NewStruct(map[string]interface{}{
		"sources": s.mux.SourcesAsMetadata(),
	})
	if err != nil {
		s.log.Warn(fmt.Sprintf("error from struct creation: %v", err))
		return nil, fmt.Errorf("error constructing response")
	}

	return &syncv1.GetMetadataResponse{
			Metadata: metadata,
		},
		nil
}
