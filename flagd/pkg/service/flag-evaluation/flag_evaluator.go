package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	schemaV1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/schema/v1"
	"connectrpc.com/connect"
	"github.com/open-feature/flagd/core/pkg/evaluator"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/service"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	flagdService "github.com/open-feature/flagd/flagd/pkg/service"
	"github.com/rs/xid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/structpb"
)

type resolverSignature[T constraints] func(context context.Context, reqID, flagKey string, ctx map[string]any) (
	T, string, string, map[string]interface{}, error)

// OldFlagEvaluationService implements the methods required for the soon-to-be deprecated flag evaluation schema
// this can be removed as a part of https://github.com/open-feature/flagd/issues/1088
type OldFlagEvaluationService struct {
	logger                *logger.Logger
	eval                  evaluator.IEvaluator
	metrics               telemetry.IMetricsRecorder
	eventingConfiguration IEvents
	flagEvalTracer        trace.Tracer
	contextValues         map[string]any
	selectorFallbackKey   string
}

// NewOldFlagEvaluationService creates a OldFlagEvaluationService with provided parameters
func NewOldFlagEvaluationService(
	log *logger.Logger,
	eval evaluator.IEvaluator,
	eventingCfg IEvents,
	metricsRecorder telemetry.IMetricsRecorder,
	contextValues map[string]any,
	selectorFallback string,
) *OldFlagEvaluationService {
	svc := &OldFlagEvaluationService{
		logger:                log,
		eval:                  eval,
		metrics:               &telemetry.NoopMetricsRecorder{},
		eventingConfiguration: eventingCfg,
		flagEvalTracer:        otel.Tracer("flagEvaluationService"),
		contextValues:         contextValues,
		selectorFallbackKey:   selectorFallback,
	}

	if metricsRecorder != nil {
		svc.metrics = metricsRecorder
	}

	return svc
}

// nolint:dupl,funlen,staticcheck
func (s *OldFlagEvaluationService) ResolveAll(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveAllRequest],
) (*connect.Response[schemaV1.ResolveAllResponse], error) {
	reqID := xid.New().String()
	defer s.logger.ClearFields(reqID)
	ctx, span := s.flagEvalTracer.Start(ctx, "resolveAll", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()
	res := &schemaV1.ResolveAllResponse{
		Flags: make(map[string]*schemaV1.AnyFlag),
	}

	selectorExpression := req.Header().Get(flagdService.FLAGD_SELECTOR_HEADER)
	selector := store.NewSelectorWithFallback(selectorExpression, s.selectorFallbackKey)
	ctx = context.WithValue(ctx, store.SelectorContextKey{}, selector)

	values, _, err := s.eval.ResolveAllValues(ctx, reqID, mergeContexts(req.Msg.GetContext().AsMap(), s.contextValues, req.Header(), make(map[string]string)))
	if err != nil {
		s.logger.WarnWithID(reqID, fmt.Sprintf("error resolving all flags: %v", err))
		return nil, fmt.Errorf("error resolving flags. Tracking ID: %s", reqID)
	}

	span.SetAttributes(attribute.Int("feature_flag.count", len(values)))
	for _, value := range values {
		// register the impression and reason for each flag evaluated
		s.metrics.RecordEvaluation(ctx, value.Error, value.Reason, value.Variant, value.FlagKey)

		switch v := value.Value.(type) {
		case bool:
			res.Flags[value.FlagKey] = &schemaV1.AnyFlag{
				Reason:  value.Reason,
				Variant: value.Variant,
				Value: &schemaV1.AnyFlag_BoolValue{
					BoolValue: v,
				},
			}
		case string:
			res.Flags[value.FlagKey] = &schemaV1.AnyFlag{
				Reason:  value.Reason,
				Variant: value.Variant,
				Value: &schemaV1.AnyFlag_StringValue{
					StringValue: v,
				},
			}
		case float64:
			res.Flags[value.FlagKey] = &schemaV1.AnyFlag{
				Reason:  value.Reason,
				Variant: value.Variant,
				Value: &schemaV1.AnyFlag_DoubleValue{
					DoubleValue: v,
				},
			}
		case map[string]any:
			val, err := structpb.NewStruct(v)
			if err != nil {
				s.logger.ErrorWithID(reqID, fmt.Sprintf("struct response construction: %v", err))
				continue
			}
			res.Flags[value.FlagKey] = &schemaV1.AnyFlag{
				Reason:  value.Reason,
				Variant: value.Variant,
				Value: &schemaV1.AnyFlag_ObjectValue{
					ObjectValue: val,
				},
			}
		}
	}
	return connect.NewResponse(res), nil
}

// nolint:dupl,staticcheck
func (s *OldFlagEvaluationService) EventStream(
	ctx context.Context,
	req *connect.Request[schemaV1.EventStreamRequest],
	stream *connect.ServerStream[schemaV1.EventStreamResponse],
) error {
	s.logger.Debug(fmt.Sprintf("starting event stream for request"))

	requestNotificationChan := make(chan service.Notification, 1)
	selectorExpression := req.Header().Get(flagdService.FLAGD_SELECTOR_HEADER)
	selector := store.NewSelectorWithFallback(selectorExpression, s.selectorFallbackKey)
	s.eventingConfiguration.Subscribe(ctx, req, &selector, requestNotificationChan)
	defer s.eventingConfiguration.Unsubscribe(req)

	requestNotificationChan <- service.Notification{
		Type: service.ProviderReady,
	}
	for {
		select {
		case <-time.After(20 * time.Second):
			err := stream.Send(&schemaV1.EventStreamResponse{
				Type: string(service.KeepAlive),
			})
			if err != nil {
				s.logger.Error(err.Error())
			}
		case notification := <-requestNotificationChan:
			d, err := structpb.NewStruct(notification.Data)
			if err != nil {
				s.logger.Error(err.Error())
			}
			err = stream.Send(&schemaV1.EventStreamResponse{
				Type: string(notification.Type),
				Data: d,
			})
			if err != nil {
				s.logger.Error(err.Error())
			}
		case <-ctx.Done():
			return nil
		}
	}
}

//nolint:staticcheck
func (s *OldFlagEvaluationService) ResolveBoolean(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveBooleanRequest],
) (*connect.Response[schemaV1.ResolveBooleanResponse], error) {
	ctx, span := s.flagEvalTracer.Start(ctx, "resolveBoolean", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()
	res := connect.NewResponse(&schemaV1.ResolveBooleanResponse{})
	selectorExpression := req.Header().Get(flagdService.FLAGD_SELECTOR_HEADER)
	selector := store.NewSelectorWithFallback(selectorExpression, s.selectorFallbackKey)
	ctx = context.WithValue(ctx, store.SelectorContextKey{}, selector)

	err := resolve[bool](
		ctx,
		s.logger,
		s.eval.ResolveBooleanValue,
		req.Header(),
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&booleanResponse{schemaV1Resp: res},
		s.metrics,
		s.contextValues,
		make(map[string]string),
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("error evaluating flag with key %s", req.Msg.GetFlagKey()))
	}

	return res, err
}

//nolint:staticcheck
func (s *OldFlagEvaluationService) ResolveString(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveStringRequest],
) (*connect.Response[schemaV1.ResolveStringResponse], error) {
	ctx, span := s.flagEvalTracer.Start(ctx, "resolveString", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	selectorExpression := req.Header().Get(flagdService.FLAGD_SELECTOR_HEADER)
	selector := store.NewSelectorWithFallback(selectorExpression, s.selectorFallbackKey)
	ctx = context.WithValue(ctx, store.SelectorContextKey{}, selector)

	res := connect.NewResponse(&schemaV1.ResolveStringResponse{})
	err := resolve[string](
		ctx,
		s.logger,
		s.eval.ResolveStringValue,
		req.Header(),
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&stringResponse{schemaV1Resp: res},
		s.metrics,
		s.contextValues,
		make(map[string]string),
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("error evaluating flag with key %s", req.Msg.GetFlagKey()))
	}

	return res, err
}

//nolint:staticcheck
func (s *OldFlagEvaluationService) ResolveInt(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveIntRequest],
) (*connect.Response[schemaV1.ResolveIntResponse], error) {
	ctx, span := s.flagEvalTracer.Start(ctx, "resolveInt", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	selectorExpression := req.Header().Get(flagdService.FLAGD_SELECTOR_HEADER)
	selector := store.NewSelectorWithFallback(selectorExpression, s.selectorFallbackKey)
	ctx = context.WithValue(ctx, store.SelectorContextKey{}, selector)

	res := connect.NewResponse(&schemaV1.ResolveIntResponse{})
	err := resolve[int64](
		ctx,
		s.logger,
		s.eval.ResolveIntValue,
		req.Header(),
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&intResponse{schemaV1Resp: res},
		s.metrics,
		s.contextValues,
		make(map[string]string),
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("error evaluating flag with key %s", req.Msg.GetFlagKey()))
	}

	return res, err
}

//nolint:staticcheck
func (s *OldFlagEvaluationService) ResolveFloat(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveFloatRequest],
) (*connect.Response[schemaV1.ResolveFloatResponse], error) {
	ctx, span := s.flagEvalTracer.Start(ctx, "resolveFloat", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	selectorExpression := req.Header().Get(flagdService.FLAGD_SELECTOR_HEADER)
	selector := store.NewSelectorWithFallback(selectorExpression, s.selectorFallbackKey)
	ctx = context.WithValue(ctx, store.SelectorContextKey{}, selector)

	res := connect.NewResponse(&schemaV1.ResolveFloatResponse{})
	err := resolve[float64](
		ctx,
		s.logger,
		s.eval.ResolveFloatValue,
		req.Header(),
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&floatResponse{schemaV1Resp: res},
		s.metrics,
		s.contextValues,
		make(map[string]string),
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("error evaluating flag with key %s", req.Msg.GetFlagKey()))
	}

	return res, err
}

//nolint:staticcheck
func (s *OldFlagEvaluationService) ResolveObject(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveObjectRequest],
) (*connect.Response[schemaV1.ResolveObjectResponse], error) {
	ctx, span := s.flagEvalTracer.Start(ctx, "resolveObject", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	selectorExpression := req.Header().Get(flagdService.FLAGD_SELECTOR_HEADER)
	selector := store.NewSelectorWithFallback(selectorExpression, s.selectorFallbackKey)
	ctx = context.WithValue(ctx, store.SelectorContextKey{}, selector)

	res := connect.NewResponse(&schemaV1.ResolveObjectResponse{})
	err := resolve[map[string]any](
		ctx,
		s.logger,
		s.eval.ResolveObjectValue,
		req.Header(),
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&objectResponse{schemaV1Resp: res},
		s.metrics,
		s.contextValues,
		make(map[string]string),
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("error evaluating flag with key %s", req.Msg.GetFlagKey()))
	}

	return res, err
}

// mergeContexts combines context values from headers, static context (from cli) and request context.
// highest priority > header-context-from-cli > static-context-from-cli > request-context > lowest priority
func mergeContexts(reqCtx, configFlagsCtx map[string]any, headers http.Header, headerToContextKeyMappings map[string]string) map[string]any {
	merged := make(map[string]any)
	for k, v := range reqCtx {
		merged[k] = v
	}
	for k, v := range configFlagsCtx {
		merged[k] = v
	}
	for header, contextKey := range headerToContextKeyMappings {
		if values, ok := headers[header]; ok {
			merged[contextKey] = values[0]
		}
	}
	return merged
}

// resolve is a generic flag resolver
func resolve[T constraints](ctx context.Context, logger *logger.Logger, resolver resolverSignature[T], header http.Header, flagKey string,
	evaluationContext *structpb.Struct, resp response[T], metrics telemetry.IMetricsRecorder,
	configContextValues map[string]any, configHeaderToContextKeyMappings map[string]string,
) error {
	reqID := xid.New().String()
	defer logger.ClearFields(reqID)

	mergedContext := mergeContexts(evaluationContext.AsMap(), configContextValues, header, configHeaderToContextKeyMappings)

	logger.WriteFields(
		reqID,
		zap.String("flag-key", flagKey),
		zap.Strings("context-keys", formatContextKeys(mergedContext)),
	)

	var evalErrFormatted error
	result, variant, reason, metadata, evalErr := resolver(ctx, reqID, flagKey, mergedContext)
	if evalErr != nil {
		logger.WarnWithID(reqID, fmt.Sprintf("returning error response, reason: %v", evalErr))
		reason = model.ErrorReason
		evalErrFormatted = errFormat(evalErr)
	}

	if metrics != nil {
		metrics.RecordEvaluation(ctx, evalErr, reason, variant, flagKey)
	}

	spanFromContext := trace.SpanFromContext(ctx)
	spanFromContext.SetAttributes(telemetry.SemConvFeatureFlagAttributes(flagKey, variant)...)

	if err := resp.SetResult(result, variant, reason, metadata); err != nil && evalErr == nil {
		logger.ErrorWithID(reqID, err.Error())
		return fmt.Errorf("error setting response result: %w", err)
	}

	return evalErrFormatted
}

func formatContextKeys(context map[string]any) []string {
	res := []string{}
	for k := range context {
		res = append(res, k)
	}
	return res
}

func errFormat(err error) error {
	ReadableErrorMsg := model.GetErrorMessage(err.Error())
	switch err.Error() {
	case model.FlagNotFoundErrorCode, model.FlagDisabledErrorCode:
		return connect.NewError(connect.CodeNotFound, fmt.Errorf("%s", ReadableErrorMsg))
	case model.TypeMismatchErrorCode:
		return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%s", ReadableErrorMsg))
	case model.ParseErrorCode:
		return connect.NewError(connect.CodeDataLoss, fmt.Errorf("%s", ReadableErrorMsg))
	case model.GeneralErrorCode:
		return connect.NewError(connect.CodeUnknown, fmt.Errorf("%s", ReadableErrorMsg))
	}

	return err
}
