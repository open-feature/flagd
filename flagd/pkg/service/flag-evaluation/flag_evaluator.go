package service

import (
	"context"
	"fmt"
	"time"

	schemaV1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/schema/v1"
	"connectrpc.com/connect"
	"github.com/open-feature/flagd/core/pkg/evaluator"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/service"
	"github.com/open-feature/flagd/core/pkg/telemetry"
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
}

// NewOldFlagEvaluationService creates a OldFlagEvaluationService with provided parameters
func NewOldFlagEvaluationService(
	log *logger.Logger,
	eval evaluator.IEvaluator,
	eventingCfg IEvents,
	metricsRecorder telemetry.IMetricsRecorder,
	contextValues map[string]any,
) *OldFlagEvaluationService {
	svc := &OldFlagEvaluationService{
		logger:                log,
		eval:                  eval,
		metrics:               &telemetry.NoopMetricsRecorder{},
		eventingConfiguration: eventingCfg,
		flagEvalTracer:        otel.Tracer("flagEvaluationService"),
		contextValues:         contextValues,
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
	sCtx, span := s.flagEvalTracer.Start(ctx, "resolveAll", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()
	res := &schemaV1.ResolveAllResponse{
		Flags: make(map[string]*schemaV1.AnyFlag),
	}
	evalCtx := map[string]any{}
	if e := req.Msg.GetContext(); e != nil {
		evalCtx = e.AsMap()
	}
	for k, v := range s.contextValues {
		evalCtx[k] = v
	}

	values, err := s.eval.ResolveAllValues(sCtx, reqID, evalCtx)
	if err != nil {
		s.logger.WarnWithID(reqID, fmt.Sprintf("error resolving all flags: %v", err))
		return nil, fmt.Errorf("error resolving flags. Tracking ID: %s", reqID)
	}

	span.SetAttributes(attribute.Int("feature_flag.count", len(values)))
	for _, value := range values {
		// register the impression and reason for each flag evaluated
		s.metrics.RecordEvaluation(sCtx, value.Error, value.Reason, value.Variant, value.FlagKey)

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
	requestNotificationChan := make(chan service.Notification, 1)
	s.eventingConfiguration.Subscribe(req, requestNotificationChan)
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
	sCtx, span := s.flagEvalTracer.Start(ctx, "resolveBoolean", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()
	res := connect.NewResponse(&schemaV1.ResolveBooleanResponse{})
	evalCtx := map[string]any{}
	if e := req.Msg.GetContext(); e != nil {
		evalCtx = e.AsMap()
	}
	for k, v := range s.contextValues {
		evalCtx[k] = v
	}

	err := resolve[bool](
		sCtx,
		s.logger,
		s.eval.ResolveBooleanValue,
		req.Msg.GetFlagKey(),
		evalCtx,
		&booleanResponse{schemaV1Resp: res},
		s.metrics,
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
	sCtx, span := s.flagEvalTracer.Start(ctx, "resolveString", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	evalCtx := map[string]any{}
	if e := req.Msg.GetContext(); e != nil {
		evalCtx = e.AsMap()
	}
	for k, v := range s.contextValues {
		evalCtx[k] = v
	}

	res := connect.NewResponse(&schemaV1.ResolveStringResponse{})
	err := resolve[string](
		sCtx,
		s.logger,
		s.eval.ResolveStringValue,
		req.Msg.GetFlagKey(),
		evalCtx,
		&stringResponse{schemaV1Resp: res},
		s.metrics,
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
	sCtx, span := s.flagEvalTracer.Start(ctx, "resolveInt", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	evalCtx := map[string]any{}
	if e := req.Msg.GetContext(); e != nil {
		evalCtx = e.AsMap()
	}
	for k, v := range s.contextValues {
		evalCtx[k] = v
	}

	res := connect.NewResponse(&schemaV1.ResolveIntResponse{})
	err := resolve[int64](
		sCtx,
		s.logger,
		s.eval.ResolveIntValue,
		req.Msg.GetFlagKey(),
		evalCtx,
		&intResponse{schemaV1Resp: res},
		s.metrics,
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
	sCtx, span := s.flagEvalTracer.Start(ctx, "resolveFloat", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	evalCtx := map[string]any{}
	if e := req.Msg.GetContext(); e != nil {
		evalCtx = e.AsMap()
	}
	for k, v := range s.contextValues {
		evalCtx[k] = v
	}

	res := connect.NewResponse(&schemaV1.ResolveFloatResponse{})
	err := resolve[float64](
		sCtx,
		s.logger,
		s.eval.ResolveFloatValue,
		req.Msg.GetFlagKey(),
		evalCtx,
		&floatResponse{schemaV1Resp: res},
		s.metrics,
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
	sCtx, span := s.flagEvalTracer.Start(ctx, "resolveObject", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	evalCtx := map[string]any{}
	if e := req.Msg.GetContext(); e != nil {
		evalCtx = e.AsMap()
	}
	for k, v := range s.contextValues {
		evalCtx[k] = v
	}

	res := connect.NewResponse(&schemaV1.ResolveObjectResponse{})
	err := resolve[map[string]any](
		sCtx,
		s.logger,
		s.eval.ResolveObjectValue,
		req.Msg.GetFlagKey(),
		evalCtx,
		&objectResponse{schemaV1Resp: res},
		s.metrics,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("error evaluating flag with key %s", req.Msg.GetFlagKey()))
	}

	return res, err
}

// resolve is a generic flag resolver
func resolve[T constraints](ctx context.Context, logger *logger.Logger, resolver resolverSignature[T], flagKey string,
	evaluationContext map[string]any, resp response[T], metrics telemetry.IMetricsRecorder,
) error {
	reqID := xid.New().String()
	defer logger.ClearFields(reqID)

	logger.WriteFields(
		reqID,
		zap.String("flag-key", flagKey),
		zap.Strings("context-keys", formatContextKeys(evaluationContext)),
	)

	var evalErrFormatted error
	result, variant, reason, metadata, evalErr := resolver(ctx, reqID, flagKey, evaluationContext)
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
