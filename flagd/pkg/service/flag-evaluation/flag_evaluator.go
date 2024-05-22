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
}

// NewOldFlagEvaluationService creates a OldFlagEvaluationService with provided parameters
func NewOldFlagEvaluationService(log *logger.Logger,
	eval evaluator.IEvaluator, eventingCfg IEvents, metricsRecorder *telemetry.MetricsRecorder,
) *OldFlagEvaluationService {
	svc := &OldFlagEvaluationService{
		logger:                log,
		eval:                  eval,
		metrics:               &telemetry.NoopMetricsRecorder{},
		eventingConfiguration: eventingCfg,
		flagEvalTracer:        otel.Tracer("flagEvaluationService"),
	}

	if metricsRecorder != nil {
		svc.metrics = metricsRecorder
	}

	return svc
}

// nolint:dupl
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
	values := s.eval.ResolveAllValues(sCtx, reqID, evalCtx)
	span.SetAttributes(attribute.Int("feature_flag.count", len(values)))
	for _, value := range values {
		// register the impression and reason for each flag evaluated
		if s.metrics != nil {
			s.metrics.RecordEvaluation(sCtx, value.Error, value.Reason, value.Variant, value.FlagKey)
		}
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

func (s *OldFlagEvaluationService) ResolveBoolean(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveBooleanRequest],
) (*connect.Response[schemaV1.ResolveBooleanResponse], error) {
	sCtx, span := s.flagEvalTracer.Start(ctx, "resolveBoolean", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()
	res := connect.NewResponse(&schemaV1.ResolveBooleanResponse{})
	err := resolve[bool](
		sCtx,
		s.logger,
		s.eval.ResolveBooleanValue,
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&booleanResponse{schemaV1Resp: res},
		s.metrics,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("error evaluating flag with key %s", req.Msg.GetFlagKey()))
	}

	return res, err
}

func (s *OldFlagEvaluationService) ResolveString(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveStringRequest],
) (*connect.Response[schemaV1.ResolveStringResponse], error) {
	sCtx, span := s.flagEvalTracer.Start(ctx, "resolveString", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	res := connect.NewResponse(&schemaV1.ResolveStringResponse{})
	err := resolve[string](
		sCtx,
		s.logger,
		s.eval.ResolveStringValue,
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&stringResponse{schemaV1Resp: res},
		s.metrics,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("error evaluating flag with key %s", req.Msg.GetFlagKey()))
	}

	return res, err
}

func (s *OldFlagEvaluationService) ResolveInt(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveIntRequest],
) (*connect.Response[schemaV1.ResolveIntResponse], error) {
	sCtx, span := s.flagEvalTracer.Start(ctx, "resolveInt", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	res := connect.NewResponse(&schemaV1.ResolveIntResponse{})
	err := resolve[int64](
		sCtx,
		s.logger,
		s.eval.ResolveIntValue,
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&intResponse{schemaV1Resp: res},
		s.metrics,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("error evaluating flag with key %s", req.Msg.GetFlagKey()))
	}

	return res, err
}

func (s *OldFlagEvaluationService) ResolveFloat(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveFloatRequest],
) (*connect.Response[schemaV1.ResolveFloatResponse], error) {
	sCtx, span := s.flagEvalTracer.Start(ctx, "resolveFloat", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	res := connect.NewResponse(&schemaV1.ResolveFloatResponse{})
	err := resolve[float64](
		sCtx,
		s.logger,
		s.eval.ResolveFloatValue,
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&floatResponse{schemaV1Resp: res},
		s.metrics,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("error evaluating flag with key %s", req.Msg.GetFlagKey()))
	}

	return res, err
}

func (s *OldFlagEvaluationService) ResolveObject(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveObjectRequest],
) (*connect.Response[schemaV1.ResolveObjectResponse], error) {
	sCtx, span := s.flagEvalTracer.Start(ctx, "resolveObject", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	res := connect.NewResponse(&schemaV1.ResolveObjectResponse{})
	err := resolve[map[string]any](
		sCtx,
		s.logger,
		s.eval.ResolveObjectValue,
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
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
	evaluationContext *structpb.Struct, resp response[T], metrics telemetry.IMetricsRecorder,
) error {
	reqID := xid.New().String()
	defer logger.ClearFields(reqID)

	logger.WriteFields(
		reqID,
		zap.String("flag-key", flagKey),
		zap.Strings("context-keys", formatContextKeys(evaluationContext)),
	)

	var evalErrFormatted error
	result, variant, reason, metadata, evalErr := resolver(ctx, reqID, flagKey, evaluationContext.AsMap())
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

func formatContextKeys(context *structpb.Struct) []string {
	res := []string{}
	for k := range context.AsMap() {
		res = append(res, k)
	}
	return res
}

func errFormat(err error) error {
	switch err.Error() {
	case model.FlagNotFoundErrorCode, model.FlagDisabledErrorCode:
		return connect.NewError(connect.CodeNotFound, fmt.Errorf("%s, %s", ErrorPrefix, err.Error()))
	case model.TypeMismatchErrorCode:
		return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%s, %s", ErrorPrefix, err.Error()))
	case model.ParseErrorCode:
		return connect.NewError(connect.CodeDataLoss, fmt.Errorf("%s, %s", ErrorPrefix, err.Error()))
	case model.GeneralErrorCode:
		return connect.NewError(connect.CodeUnknown, fmt.Errorf("%s, %s", ErrorPrefix, err.Error()))
	}

	return err
}
