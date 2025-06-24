package sync

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

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
	deadline      time.Duration
}

func (s syncHandler) SyncFlags(req *syncv1.SyncFlagsRequest, server syncv1grpc.FlagSyncService_SyncFlagsServer) error {
	muxPayload := make(chan payload, 1)
	selector := req.GetSelector()
	ctx := server.Context()

	// attach server-side stream deadline to context
	if s.deadline != 0 {
		streamDeadline := time.Now().Add(s.deadline)
		deadlineCtx, cancel := context.WithDeadline(server.Context(), streamDeadline)
		ctx = deadlineCtx
		defer cancel()
	}

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

			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				s.log.Debug(fmt.Sprintf("server-side deadline of %s exceeded, exiting stream request with grpc error code 4", s.deadline.String()))
				return status.Error(codes.DeadlineExceeded, "stream closed due to server-side timeout")
			}
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
