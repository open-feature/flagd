package sync

import (
	"context"
	"fmt"
	"maps"

	"buf.build/gen/go/open-feature/flagd/grpc/go/flagd/sync/v1/syncv1grpc"
	syncv1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/sync/v1"
	"github.com/open-feature/flagd/core/pkg/logger"
	"google.golang.org/protobuf/types/known/structpb"
)

// syncHandler implements the sync contract
type syncHandler struct {
	mux           *Multiplexer
	log           *logger.Logger
	contextValues map[string]any
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

			metadataSrc := make(map[string]any)
			maps.Copy(metadataSrc, s.contextValues)

			if sources := s.mux.SourcesAsMetadata(); sources != "" {
				metadataSrc["sources"] = sources
			}

			metadata, err := structpb.NewStruct(metadataSrc)
			if err != nil {
				s.log.Warn(fmt.Sprintf("error from struct creation: %v", err))
				return fmt.Errorf("error constructing metadata response")
			}

			err = server.Send(&syncv1.SyncFlagsResponse{
				FlagConfiguration: payload.flags,
				SyncContext:       metadata,
			})
			if err != nil {
				s.log.Debug(fmt.Sprintf("error sending stream response: %v", err))
				return fmt.Errorf("error sending stream response: %w", err)
			}
		case <-ctx.Done():
			s.mux.Unregister(ctx, selector)
			s.log.Debug("context complete and exiting stream request")
			return nil
		}
	}
}

func (s syncHandler) FetchAllFlags(_ context.Context, req *syncv1.FetchAllFlagsRequest) (
	*syncv1.FetchAllFlagsResponse, error,
) {
	flags, err := s.mux.GetAllFlags(req.GetSelector())
	if err != nil {
		return nil, err
	}

	return &syncv1.FetchAllFlagsResponse{
		FlagConfiguration: flags,
	}, nil
}

// Deprecated - GetMetadata is deprecated and will be removed in a future release.
// User the sync_context field in syncv1.SyncFlagsResponse instead.
func (s syncHandler) GetMetadata(_ context.Context, _ *syncv1.GetMetadataRequest) (
	*syncv1.GetMetadataResponse, error,
) {
	metadataSrc := make(map[string]any)
	for k, v := range s.contextValues {
		metadataSrc[k] = v
	}
	if sources := s.mux.SourcesAsMetadata(); sources != "" {
		metadataSrc["sources"] = sources
	}

	metadata, err := structpb.NewStruct(metadataSrc)
	if err != nil {
		s.log.Warn(fmt.Sprintf("error from struct creation: %v", err))
		return nil, fmt.Errorf("error constructing metadata response")
	}

	return &syncv1.GetMetadataResponse{
			Metadata: metadata,
		},
		nil
}
