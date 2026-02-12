package sync

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"time"

	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"

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
	metricsRecorder     telemetry.IMetricsRecorder
}

func (s syncHandler) SyncFlags(req *syncv1.SyncFlagsRequest, server syncv1grpc.FlagSyncService_SyncFlagsServer) error {
	startTime := time.Now()
	selectorExpression := s.getSelectorExpression(server.Context(), req)

	// Build metric attributes
	attrs := []attribute.KeyValue{}
	if selectorExpression != "" {
		attrs = append(attrs, attribute.String("selector", selectorExpression))
	}
	if req.GetProviderId() != "" {
		attrs = append(attrs, attribute.String("provider_id", req.GetProviderId()))
	}

	// Record stream start
	s.metricsRecorder.SyncStreamStart(server.Context(), attrs)

	// Track exit reason for duration metric
	var exitReason string
	defer func() {
		duration := time.Since(startTime)
		reasonAttrs := append([]attribute.KeyValue{}, attrs...)
		reasonAttrs = append(reasonAttrs, attribute.String("reason", exitReason))
		s.metricsRecorder.SyncStreamEnd(server.Context(), attrs)
		s.metricsRecorder.SyncStreamDuration(server.Context(), duration, reasonAttrs)
	}()

	watcher := make(chan store.FlagQueryResult, 1)
	selector := store.NewSelector(selectorExpression)
	ctx := server.Context()

	syncContextMap := make(map[string]any)
	maps.Copy(syncContextMap, s.contextValues)
	syncContext, err := structpb.NewStruct(syncContextMap)
	if err != nil {
		exitReason = "error"
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
				exitReason = "error"
				return fmt.Errorf("error constructing metadata response")
			}

			flags, err := s.generateResponse(payload.Flags)
			if err != nil {
				s.log.Error(fmt.Sprintf("error retrieving flags from store: %v", err))
				exitReason = "error"
				return status.Error(codes.DataLoss, "error marshalling flags")
			}

			err = server.Send(&syncv1.SyncFlagsResponse{FlagConfiguration: string(flags), SyncContext: syncContext})
			if err != nil {
				s.log.Debug(fmt.Sprintf("error sending stream response: %v", err))
				exitReason = "client_disconnect"
				return fmt.Errorf("error sending stream response: %w", err)
			}
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				s.log.Debug(fmt.Sprintf("server-side deadline of %s exceeded, exiting stream request with grpc error code 4", s.deadline.String()))
				exitReason = "deadline_exceeded"
				return status.Error(codes.DeadlineExceeded, "stream closed due to server-side timeout")
			}
			s.log.Debug("context complete and exiting stream request")
			exitReason = "normal_close"
			return nil
		}
	}
}

func (s syncHandler) generateResponse(payload []model.Flag) ([]byte, error) {
	flagConfig := map[string]interface{}{
		"flags": s.convertMap(payload),
	}

	flags, err := json.Marshal(flagConfig)
	return flags, err
}

// getSelectorExpression extracts the selector expression from the request.
// It first checks the Flagd-Selector header (metadata), then falls back to the request body selector.
//
// The req parameter accepts *syncv1.SyncFlagsRequest or *syncv1.FetchAllFlagsRequest.
// Using interface{} here is intentional as both protobuf-generated types implement GetSelector()
// but do not share a common interface.
func (s syncHandler) getSelectorExpression(ctx context.Context, req interface{}) string {
	// Try to get selector from metadata (header)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get(flagdService.FLAGD_SELECTOR_HEADER); len(values) > 0 {
			headerSelector := values[0]
			s.log.Debug(fmt.Sprintf("using selector from request header: %s", headerSelector))
			return headerSelector
		}
	}

	// Fall back to request body selector for backward compatibility
	// Eventually we will want to log a deprecation warning here and then remote it entirely
	var bodySelector string
	type selectorGetter interface {
		GetSelector() string
	}

	if r, ok := req.(selectorGetter); ok {
		bodySelector = r.GetSelector()
	}

	if bodySelector != "" {
		s.log.Debug(fmt.Sprintf("using selector from request body: %s", bodySelector))
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

	flagsString, err := s.generateResponse(flags)

	if err != nil {
		return nil, err
	}

	return &syncv1.FetchAllFlagsResponse{
		FlagConfiguration: string(flagsString),
	}, nil
}

// Deprecated - GetMetadata is deprecated and will be removed in a future release.
// Use the sync_context field in syncv1.SyncFlagsResponse, providing same info.
//nolint:staticcheck // SA1019 temporarily suppress deprecation warning
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
