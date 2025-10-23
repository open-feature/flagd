package sync

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"time"

	"github.com/open-feature/flagd/core/pkg/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"buf.build/gen/go/open-feature/flagd/grpc/go/flagd/sync/v1/syncv1grpc"
	syncv1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/sync/v1"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/store"
	flagdService "github.com/open-feature/flagd/flagd/pkg/service"
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
}

func (s syncHandler) SyncFlags(req *syncv1.SyncFlagsRequest, server syncv1grpc.FlagSyncService_SyncFlagsServer) error {
	watcher := make(chan store.FlagQueryResult, 1)
	selectorExpression := s.getSelectorExpression(server.Context(), req)
	selector := store.NewSelector(selectorExpression)
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

			flagMap := s.convertMap(payload.Flags)

			flags, err := json.Marshal(flagMap)
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

// getSelectorExpression extracts the selector expression from the request.
// It first checks the Flagd-Selector header (metadata), then falls back to the request body selector.
// A deprecation warning is logged when the request body selector is used.
func (s syncHandler) getSelectorExpression(ctx context.Context, req interface{}) string {
	// Try to get selector from metadata (header)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get(flagdService.FLAGD_SELECTOR_HEADER); len(values) > 0 {
			return values[0]
		}
	}

	// Fall back to request body selector for backward compatibility
	var bodySelector string
	switch r := req.(type) {
	case *syncv1.SyncFlagsRequest:
		bodySelector = r.GetSelector()
	case *syncv1.FetchAllFlagsRequest:
		bodySelector = r.GetSelector()
	}

	// Log deprecation warning if using request body selector
	if bodySelector != "" {
		s.log.Warn("Using selector from request body is deprecated. Please use the 'Flagd-Selector' header instead. " +
			"Request body selector support will be removed in a future major version.")
	}

	return bodySelector
}

func (s syncHandler) convertMap(flags []model.Flag) map[string]model.Flag {
	flagMap := make(map[string]model.Flag, len(flags))
	for _, flag := range flags {
		flagMap[flag.Key] = flag
	}
	return flagMap
}

func (s syncHandler) FetchAllFlags(ctx context.Context, req *syncv1.FetchAllFlagsRequest) (
	*syncv1.FetchAllFlagsResponse, error,
) {
	selectorExpression := s.getSelectorExpression(ctx, req)
	selector := store.NewSelector(selectorExpression)
	flags, _, err := s.store.GetAll(ctx, &selector)
	if err != nil {
		s.log.Error(fmt.Sprintf("error retrieving flags from store: %v", err))
		return nil, status.Error(codes.Internal, "error retrieving flags from store")
	}

	flagsString, err := json.Marshal(s.convertMap(flags))

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
