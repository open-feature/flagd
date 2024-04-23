package service

import (
	"context"
	"fmt"
	"time"

	evalV1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/evaluation/v1"
	"connectrpc.com/connect"
	"github.com/open-feature/flagd/core/pkg/evaluator"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/service"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	"github.com/rs/xid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/structpb"
)

type FlagEvaluationService struct {
	logger                *logger.Logger
	eval                  evaluator.IEvaluator
	metrics               *telemetry.MetricsRecorder
	eventingConfiguration IEvents
	flagEvalTracer        trace.Tracer
}

// NewFlagEvaluationService creates a FlagEvaluationService with provided parameters
func NewFlagEvaluationService(log *logger.Logger,
	eval evaluator.IEvaluator,
	eventingCfg IEvents,
	metricsRecorder *telemetry.MetricsRecorder,
) *FlagEvaluationService {
	return &FlagEvaluationService{
		logger:                log,
		eval:                  eval,
		metrics:               metricsRecorder,
		eventingConfiguration: eventingCfg,
		flagEvalTracer:        otel.Tracer("flagd.evaluation.v1"),
	}
}

// nolint:dupl,funlen
func (s *FlagEvaluationService) ResolveAll(
	ctx context.Context,
	req *connect.Request[evalV1.ResolveAllRequest],
) (*connect.Response[evalV1.ResolveAllResponse], error) {
	reqID := xid.New().String()
	defer s.logger.ClearFields(reqID)

	sCtx, span := s.flagEvalTracer.Start(ctx, "resolveAll", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	res := &evalV1.ResolveAllResponse{
		Flags: make(map[string]*evalV1.AnyFlag),
	}

	evalCtx := map[string]any{}
	if e := req.Msg.GetContext(); e != nil {
		evalCtx = e.AsMap()
	}

	values := s.eval.ResolveAllValues(sCtx, reqID, evalCtx)
	span.SetAttributes(attribute.Int("feature_flag.count", len(values)))
	for _, value := range values {
		// register the impression and reason for each flag evaluated
		s.metrics.RecordEvaluation(sCtx, value.Error, value.Reason, value.Variant, value.FlagKey)
		switch v := value.Value.(type) {
		case bool:
			res.Flags[value.FlagKey] = &evalV1.AnyFlag{
				Reason:  value.Reason,
				Variant: value.Variant,
				Value: &evalV1.AnyFlag_BoolValue{
					BoolValue: v,
				},
			}
		case string:
			res.Flags[value.FlagKey] = &evalV1.AnyFlag{
				Reason:  value.Reason,
				Variant: value.Variant,
				Value: &evalV1.AnyFlag_StringValue{
					StringValue: v,
				},
			}
		case float64:
			res.Flags[value.FlagKey] = &evalV1.AnyFlag{
				Reason:  value.Reason,
				Variant: value.Variant,
				Value: &evalV1.AnyFlag_DoubleValue{
					DoubleValue: v,
				},
			}
		case map[string]any:
			val, err := structpb.NewStruct(v)
			if err != nil {
				s.logger.ErrorWithID(reqID, fmt.Sprintf("struct response construction: %v", err))
				continue
			}
			res.Flags[value.FlagKey] = &evalV1.AnyFlag{
				Reason:  value.Reason,
				Variant: value.Variant,
				Value: &evalV1.AnyFlag_ObjectValue{
					ObjectValue: val,
				},
			}
		}
	}
	return connect.NewResponse(res), nil
}

func (s *FlagEvaluationService) EventStream(
	ctx context.Context,
	req *connect.Request[evalV1.EventStreamRequest],
	stream *connect.ServerStream[evalV1.EventStreamResponse],
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
			err := stream.Send(&evalV1.EventStreamResponse{
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
			err = stream.Send(&evalV1.EventStreamResponse{
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

func (s *FlagEvaluationService) ResolveBoolean(
	ctx context.Context,
	req *connect.Request[evalV1.ResolveBooleanRequest],
) (*connect.Response[evalV1.ResolveBooleanResponse], error) {
	sCtx, span := s.flagEvalTracer.Start(ctx, "resolveBoolean", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()
	res := connect.NewResponse(&evalV1.ResolveBooleanResponse{})
	err := resolve[bool](
		sCtx,
		s.logger,
		s.eval.ResolveBooleanValue,
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&booleanResponse{evalV1Resp: res},
		s.metrics,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("error evaluating flag with key %s", req.Msg.GetFlagKey()))
	}

	return res, err
}

func (s *FlagEvaluationService) ResolveString(
	ctx context.Context,
	req *connect.Request[evalV1.ResolveStringRequest],
) (*connect.Response[evalV1.ResolveStringResponse], error) {
	sCtx, span := s.flagEvalTracer.Start(ctx, "resolveString", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	res := connect.NewResponse(&evalV1.ResolveStringResponse{})
	err := resolve[string](
		sCtx,
		s.logger,
		s.eval.ResolveStringValue,
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&stringResponse{evalV1Resp: res},
		s.metrics,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("error evaluating flag with key %s", req.Msg.GetFlagKey()))
	}

	return res, err
}

func (s *FlagEvaluationService) ResolveInt(
	ctx context.Context,
	req *connect.Request[evalV1.ResolveIntRequest],
) (*connect.Response[evalV1.ResolveIntResponse], error) {
	sCtx, span := s.flagEvalTracer.Start(ctx, "resolveInt", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	res := connect.NewResponse(&evalV1.ResolveIntResponse{})
	err := resolve[int64](
		sCtx,
		s.logger,
		s.eval.ResolveIntValue,
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&intResponse{evalV1Resp: res},
		s.metrics,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("error evaluating flag with key %s", req.Msg.GetFlagKey()))
	}

	return res, err
}

func (s *FlagEvaluationService) ResolveFloat(
	ctx context.Context,
	req *connect.Request[evalV1.ResolveFloatRequest],
) (*connect.Response[evalV1.ResolveFloatResponse], error) {
	sCtx, span := s.flagEvalTracer.Start(ctx, "resolveFloat", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	res := connect.NewResponse(&evalV1.ResolveFloatResponse{})
	err := resolve[float64](
		sCtx,
		s.logger,
		s.eval.ResolveFloatValue,
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&floatResponse{evalV1Resp: res},
		s.metrics,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("error evaluating flag with key %s", req.Msg.GetFlagKey()))
	}

	return res, err
}

func (s *FlagEvaluationService) ResolveObject(
	ctx context.Context,
	req *connect.Request[evalV1.ResolveObjectRequest],
) (*connect.Response[evalV1.ResolveObjectResponse], error) {
	sCtx, span := s.flagEvalTracer.Start(ctx, "resolveObject", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	res := connect.NewResponse(&evalV1.ResolveObjectResponse{})
	err := resolve[map[string]any](
		sCtx,
		s.logger,
		s.eval.ResolveObjectValue,
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&objectResponse{evalV1Resp: res},
		s.metrics,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("error evaluating flag with key %s", req.Msg.GetFlagKey()))
	}

	return res, err
}
