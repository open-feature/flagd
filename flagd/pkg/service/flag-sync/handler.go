package sync

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"time"

	syncv1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/sync/v1"
	"connectrpc.com/connect"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/protobuf/types/known/structpb"
)

// syncHandler implements the sync contract
type syncHandler struct {
	store               store.IStore
	log                 *logger.Logger
	contextValues       map[string]any
	deadline            time.Duration
	disableSyncMetadata bool
	metricsRecorder     telemetry.IMetricsRecorder
}

// syncFlagsSender is the subset of connect.ServerStream used here; it has no exported constructor,
// so tests cannot build one.
type syncFlagsSender interface {
	Send(*syncv1.SyncFlagsResponse) error
}

func (s syncHandler) SyncFlags(
	ctx context.Context,
	req *connect.Request[syncv1.SyncFlagsRequest],
	stream *connect.ServerStream[syncv1.SyncFlagsResponse],
) error {
	return s.syncFlags(ctx, req.Header(), req.Msg.GetProviderId(), req.Msg.GetSelector(), stream)
}

func (s syncHandler) syncFlags(
	ctx context.Context, header http.Header, providerID, bodySelector string, stream syncFlagsSender,
) error {
	startTime := time.Now()
	selectorExpression := s.getSelectorExpression(header, bodySelector)

	// Build metric attributes
	attrs := []attribute.KeyValue{}
	if selectorExpression != "" {
		attrs = append(attrs, attribute.String("selector", selectorExpression))
	}
	if providerID != "" {
		attrs = append(attrs, attribute.String("provider_id", providerID))
	}

	// Record stream start
	s.metricsRecorder.SyncStreamStart(ctx, attrs)

	// Track exit reason for duration metric
	var exitReason string
	defer func() {
		duration := time.Since(startTime)
		reasonAttrs := append([]attribute.KeyValue{}, attrs...)
		reasonAttrs = append(reasonAttrs, attribute.String("reason", exitReason))
		s.metricsRecorder.SyncStreamEnd(ctx, attrs)
		s.metricsRecorder.SyncStreamDuration(ctx, duration, reasonAttrs)
	}()

	watcher := make(chan store.FlagQueryResult, 1)
	selector, err := newSelector(selectorExpression)
	if err != nil {
		exitReason = "error"
		return s.connectError(err)
	}

	syncContextMap := make(map[string]any)
	maps.Copy(syncContextMap, s.contextValues)
	syncContext, err := structpb.NewStruct(syncContextMap)
	if err != nil {
		exitReason = "error"
		return connect.NewError(connect.CodeDataLoss, errors.New("error constructing sync context"))
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
			flags, err := generateResponse(payload.Flags)
			if err != nil {
				s.log.Error(fmt.Sprintf("error retrieving flags from store: %v", err))
				exitReason = "error"
				return connect.NewError(connect.CodeDataLoss, errors.New("error marshalling flags"))
			}

			err = stream.Send(&syncv1.SyncFlagsResponse{FlagConfiguration: string(flags), SyncContext: syncContext})
			if err != nil {
				s.log.Debug(fmt.Sprintf("error sending stream response: %v", err))
				exitReason = "client_disconnect"
				return fmt.Errorf("error sending stream response: %w", err)
			}
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				s.log.Debug(fmt.Sprintf("server-side deadline of %s exceeded, exiting stream request with grpc error code 4", s.deadline.String()))
				exitReason = "deadline_exceeded"
				return connect.NewError(connect.CodeDeadlineExceeded, errors.New("stream closed due to server-side timeout"))
			}
			s.log.Debug("context complete and exiting stream request")
			exitReason = "normal_close"
			return nil
		}
	}
}

func (s syncHandler) FetchAllFlags(ctx context.Context, req *connect.Request[syncv1.FetchAllFlagsRequest]) (
	*connect.Response[syncv1.FetchAllFlagsResponse], error,
) {
	expression := s.getSelectorExpression(req.Header(), req.Msg.GetSelector())

	flagsString, err := fetchAllFlags(ctx, s.store, expression)
	if err != nil {
		return nil, s.connectError(err)
	}

	return connect.NewResponse(&syncv1.FetchAllFlagsResponse{
		FlagConfiguration: string(flagsString),
	}), nil
}

// Deprecated - GetMetadata is deprecated and will be removed in a future release.
// Use the sync_context field in syncv1.SyncFlagsResponse, providing same info.
//
//nolint:staticcheck // SA1019 temporarily suppress deprecation warning
func (s syncHandler) GetMetadata(_ context.Context, _ *connect.Request[syncv1.GetMetadataRequest]) (
	*connect.Response[syncv1.GetMetadataResponse], error,
) {
	if s.disableSyncMetadata {
		return nil, connect.NewError(connect.CodeUnimplemented, errors.New("metadata endpoint disabled"))
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

	return connect.NewResponse(&syncv1.GetMetadataResponse{
			Metadata: metadata,
		}),
		nil
}

// connectError maps a fetch failure onto a status; both selector failures are InvalidArgument.
func (s syncHandler) connectError(err error) error {
	var fetchErr fetchError
	if !errors.As(err, &fetchErr) {
		return err
	}

	switch fetchErr.kind {
	case fetchSelectorMalformed, fetchSelectorInvalid:
		return connect.NewError(connect.CodeInvalidArgument, fetchErr.cause)
	case fetchMarshal:
		s.log.Error(fmt.Sprintf("error marshalling flags: %v", fetchErr.cause))
		return connect.NewError(connect.CodeDataLoss, errors.New("error marshalling flags"))
	default:
		s.log.Error(fmt.Sprintf("error retrieving flags from store: %v", fetchErr.cause))
		return connect.NewError(connect.CodeInternal, errors.New("error retrieving flags from store"))
	}
}

// getSelectorExpression resolves the selector and logs which slot supplied it. The body selector is
// kept for backward compatibility and should eventually be deprecated.
func (s syncHandler) getSelectorExpression(header http.Header, bodySelector string) string {
	expression := resolveSelector(header, bodySelector)
	if expression == "" {
		return ""
	}

	if resolveSelector(header, "") != "" {
		s.log.Debug(fmt.Sprintf("using selector from request header: %s", expression))
	} else {
		s.log.Debug(fmt.Sprintf("using selector from request body: %s", expression))
	}
	return expression
}

type flagConfiguration struct {
	Flags map[string]model.Flag `json:"flags"`
}

func generateResponse(payload []model.Flag) ([]byte, error) {
	flags, err := json.Marshal(flagConfiguration{Flags: convertMap(payload)})
	return flags, err
}

func convertMap(flags []model.Flag) map[string]model.Flag {
	flagMap := make(map[string]model.Flag, len(flags))
	for _, flag := range flags {
		flagMap[flag.Key] = flag
	}
	return flagMap
}
