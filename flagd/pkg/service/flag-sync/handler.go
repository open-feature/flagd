package sync

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"buf.build/gen/go/open-feature/flagd/grpc/go/flagd/sync/v1/syncv1grpc"
	syncv1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/sync/v1"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/store"
	"google.golang.org/protobuf/types/known/structpb"
)

// syncHandler implements the sync contract
type syncHandler struct {
	//mux                 *Multiplexer
	store               store.IStore
	log                 *logger.Logger
	contextValues       map[string]any
	deadline            time.Duration
	disableSyncMetadata bool

	selectorFallbackKey string
}

func (s syncHandler) SyncFlags(req *syncv1.SyncFlagsRequest, server syncv1grpc.FlagSyncService_SyncFlagsServer) error {
	watcher := make(chan store.FlagQueryResult, 1)
	selectorExpression := req.GetSelector()
	selector := store.NewSelectorWithFallback(selectorExpression, s.selectorFallbackKey)
	ctx := server.Context()

	syncContextMap := make(map[string]any)
	maps.Copy(syncContextMap, s.contextValues)
	syncContext, err := structpb.NewStruct(syncContextMap)
	if err != nil {
		return status.Error(codes.DataLoss, "error constructing sync context")
	}

	// attach server-side stream deadline to context
	if s.deadline != 0 {
		streamDeadline := time.Now().Add(s.deadline)
		deadlineCtx, cancel := context.WithDeadline(ctx, streamDeadline)
		ctx = deadlineCtx
		defer cancel()
	}

	s.store.Watch(ctx, &selector, watcher)

	for {
		select {
		case payload := <-watcher:
			if err != nil {
				s.log.Error(fmt.Sprintf("error from struct creation: %v", err))
				return fmt.Errorf("error constructing metadata response")
			}
			flags, err := json.Marshal(payload.Flags)
			if err != nil {
				s.log.Error(fmt.Sprintf("error retrieving flags from store: %v", err))
				return status.Error(codes.DataLoss, "error marshalling flags")
			}

			err = server.Send(&syncv1.SyncFlagsResponse{FlagConfiguration: string(flags), SyncContext: syncContext})
			if err != nil {
				s.log.Debug(fmt.Sprintf("error sending stream response: %v", err))
				return fmt.Errorf("error sending stream response: %w", err)
			}
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				s.log.Debug(fmt.Sprintf("server-side deadline of %s exceeded, exiting stream request with grpc error code 4", s.deadline.String()))
				return status.Error(codes.DeadlineExceeded, "stream closed due to server-side timeout")
			}
			s.log.Debug("context complete and exiting stream request")
			return nil
		}
	}
}

func (s syncHandler) FetchAllFlags(ctx context.Context, req *syncv1.FetchAllFlagsRequest) (
	*syncv1.FetchAllFlagsResponse, error,
) {
	selectorExpression := req.GetSelector()
	selector := store.NewSelector(selectorExpression)
	flags, _, err := s.store.GetAll(ctx, &selector)
	if err != nil {
		s.log.Error(fmt.Sprintf("error retrieving flags from store: %v", err))
		return nil, status.Error(codes.Internal, "error retrieving flags from store")
	}

	flagsString, err := json.Marshal(flags)

	if err != nil {
		return nil, err
	}

	return &syncv1.FetchAllFlagsResponse{
		FlagConfiguration: string(flagsString),
	}, nil
}

// Deprecated - GetMetadata is deprecated and will be removed in a future release.
// Use the sync_context field in syncv1.SyncFlagsResponse, providing same info.
func (s syncHandler) GetMetadata(_ context.Context, _ *syncv1.GetMetadataRequest) (
	*syncv1.GetMetadataResponse, error,
) {
	if s.disableSyncMetadata {
		return nil, status.Error(codes.Unimplemented, "metadata endpoint disabled")
	}
	metadataSrc := make(map[string]any)
	for k, v := range s.contextValues {
		metadataSrc[k] = v
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
