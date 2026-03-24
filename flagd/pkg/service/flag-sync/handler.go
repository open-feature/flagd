package sync

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"sort"
	"strings"
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
	selectors := parseSelectorList(selectorExpression)
	ctx := server.Context()
	syncContext, err := s.buildSyncContext()
	if err != nil {
		return status.Error(codes.DataLoss, "error constructing sync context")
	}

	var cancel context.CancelFunc
	ctx, cancel = s.withDeadline(ctx)
	if cancel != nil {
		defer cancel()
	}

	s.watchSelectors(ctx, selectors, watcher)
	return s.streamFlagUpdates(ctx, selectors, watcher, syncContext, server)
}

func (s syncHandler) buildSyncContext() (*structpb.Struct, error) {
	syncContextMap := make(map[string]any)
	maps.Copy(syncContextMap, s.contextValues)
	return structpb.NewStruct(syncContextMap)
}

func (s syncHandler) withDeadline(ctx context.Context) (context.Context, context.CancelFunc) {
	if s.deadline == 0 {
		return ctx, nil
	}

	streamDeadline := time.Now().Add(s.deadline)
	return context.WithDeadline(ctx, streamDeadline)
}

func (s syncHandler) watchSelectors(ctx context.Context, selectors []store.Selector, watcher chan store.FlagQueryResult) {
	switch len(selectors) {
	case 0:
		s.store.Watch(ctx, nil, watcher)
	case 1:
		s.store.Watch(ctx, &selectors[0], watcher)
	default:
		// For multi-selector requests, watch all updates and recompute merged view in order.
		s.store.Watch(ctx, nil, watcher)
	}
}

func (s syncHandler) streamFlagUpdates(
	ctx context.Context,
	selectors []store.Selector,
	watcher chan store.FlagQueryResult,
	syncContext *structpb.Struct,
	server syncv1grpc.FlagSyncService_SyncFlagsServer,
) error {
	for {
		select {
		case payload := <-watcher:
			flagsToSend, err := s.resolveFlagsForSelectors(ctx, selectors, payload.Flags)
			if err != nil {
				s.log.Error(fmt.Sprintf("error retrieving merged flags from store: %v", err))
				return status.Error(codes.Internal, "error retrieving flags from store")
			}

			flags, err := s.generateResponse(flagsToSend)
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

func (s syncHandler) resolveFlagsForSelectors(ctx context.Context, selectors []store.Selector, payloadFlags []model.Flag) ([]model.Flag, error) {
	if len(selectors) == 1 {
		return payloadFlags, nil
	}
	return s.fetchMergedFlags(ctx, selectors)
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
		if headerSelector := flagdService.SelectorExpressionFromGRPCMetadata(md); headerSelector != "" {
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
	flags, err := s.fetchMergedFlags(ctx, parseSelectorList(selectorExpression))
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

func parseSelectorList(selectorExpression string) []store.Selector {
	if strings.TrimSpace(selectorExpression) == "" {
		return nil
	}

	parts := strings.Split(selectorExpression, ",")
	selectors := make([]store.Selector, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		selector := store.NewSelector(trimmed)
		selectors = append(selectors, selector)
	}
	return selectors
}

func (s syncHandler) fetchMergedFlags(ctx context.Context, selectors []store.Selector) ([]model.Flag, error) {
	switch len(selectors) {
	case 0:
		flags, _, err := s.store.GetAll(ctx, nil)
		return flags, err
	case 1:
		flags, _, err := s.store.GetAll(ctx, &selectors[0])
		return flags, err
	default:
		type flagIdentifier struct {
			flagSetID string
			key       string
		}

		merged := map[flagIdentifier]model.Flag{}
		for _, selector := range selectors {
			flags, _, err := s.store.GetAll(ctx, &selector)
			if err != nil {
				return nil, err
			}
			for _, flag := range flags {
				merged[flagIdentifier{flagSetID: flag.FlagSetId, key: flag.Key}] = flag
			}
		}

		out := make([]model.Flag, 0, len(merged))
		for _, flag := range merged {
			out = append(out, flag)
		}
		sort.Slice(out, func(i, j int) bool {
			if out[i].FlagSetId != out[j].FlagSetId {
				return out[i].FlagSetId < out[j].FlagSetId
			}
			return out[i].Key < out[j].Key
		})
		return out, nil
	}
}

// Deprecated - GetMetadata is deprecated and will be removed in a future release.
// Use the sync_context field in syncv1.SyncFlagsResponse, providing same info.
//
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
