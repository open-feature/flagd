package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	evalV1 "buf.build/gen/go/open-feature-forking/flagd/protocolbuffers/go/flagd/evaluation/v1"
	"connectrpc.com/connect"
	"github.com/open-feature/flagd/core/pkg/evaluator"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/service"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	flagdService "github.com/open-feature/flagd/flagd/pkg/service"
	"github.com/rs/xid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/structpb"
)

type FlagEvaluationService struct {
	logger                     *logger.Logger
	eval                       evaluator.IEvaluator
	metrics                    telemetry.IMetricsRecorder
	eventingConfiguration      IEvents
	flagEvalTracer             trace.Tracer
	contextValues              map[string]any
	headerToContextKeyMappings map[string]string
	deadline                   time.Duration
}

// NewFlagEvaluationService creates a FlagEvaluationService with provided parameters
func NewFlagEvaluationService(log *logger.Logger,
	eval evaluator.IEvaluator,
	eventingCfg IEvents,
	metricsRecorder telemetry.IMetricsRecorder,
	contextValues map[string]any,
	headerToContextKeyMappings map[string]string,
	streamDeadline time.Duration,
) *FlagEvaluationService {
	svc := &FlagEvaluationService{
		logger:                     log,
		eval:                       eval,
		metrics:                    &telemetry.NoopMetricsRecorder{},
		eventingConfiguration:      eventingCfg,
		flagEvalTracer:             otel.Tracer("flagd.evaluation.v1"),
		contextValues:              contextValues,
		headerToContextKeyMappings: headerToContextKeyMappings,
		deadline:                   streamDeadline,
	}

	if metricsRecorder != nil {
		svc.metrics = metricsRecorder
	}

	return svc
}

// nolint:dupl,funlen
func (s *FlagEvaluationService) ResolveAll(
	ctx context.Context,
	req *connect.Request[evalV1.ResolveAllRequest],
) (*connect.Response[evalV1.ResolveAllResponse], error) {
	reqID := xid.New().String()
	defer s.logger.ClearFields(reqID)

	ctx, span := s.flagEvalTracer.Start(ctx, "resolveAll", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	res := &evalV1.ResolveAllResponse{
		Flags: make(map[string]*evalV1.AnyFlag),
	}

	selectorExpression := req.Header().Get(flagdService.FLAGD_SELECTOR_HEADER)
	selector := store.NewSelector(selectorExpression)
	evaluationContext := mergeContexts(req.Msg.GetContext().AsMap(), s.contextValues, req.Header(), s.headerToContextKeyMappings)
	ctx = context.WithValue(ctx, store.SelectorContextKey{}, selector)

	resolutions, flagSetMetadata, err := s.eval.ResolveAllValues(ctx, reqID, evaluationContext)
	if err != nil {
		s.logger.WarnWithID(reqID, fmt.Sprintf("error resolving all flags: %v", err))
		return nil, fmt.Errorf("error resolving flags. Tracking ID: %s", reqID)
	}

	span.SetAttributes(attribute.Int("feature_flag.count", len(resolutions)))
	for _, resolved := range resolutions {
		// register the impression and reason for each flag evaluated
		s.metrics.RecordEvaluation(ctx, resolved.Error, resolved.Reason, resolved.Variant, resolved.FlagKey)
		switch v := resolved.Value.(type) {
		case bool:
			res.Flags[resolved.FlagKey] = &evalV1.AnyFlag{
				Reason:  resolved.Reason,
				Variant: resolved.Variant,
				Value: &evalV1.AnyFlag_BoolValue{
					BoolValue: v,
				},
			}
		case string:
			res.Flags[resolved.FlagKey] = &evalV1.AnyFlag{
				Reason:  resolved.Reason,
				Variant: resolved.Variant,
				Value: &evalV1.AnyFlag_StringValue{
					StringValue: v,
				},
			}
		case float64:
			res.Flags[resolved.FlagKey] = &evalV1.AnyFlag{
				Reason:  resolved.Reason,
				Variant: resolved.Variant,
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
			res.Flags[resolved.FlagKey] = &evalV1.AnyFlag{
				Reason:  resolved.Reason,
				Variant: resolved.Variant,
				Value: &evalV1.AnyFlag_ObjectValue{
					ObjectValue: val,
				},
			}
		}
		metadata, err := structpb.NewStruct(resolved.Metadata)
		if err != nil {
			s.logger.WarnWithID(reqID, fmt.Sprintf("error resolving all flags: %v", err))
			return nil, fmt.Errorf("error resolving flags. Tracking ID: %s", reqID)
		}

		res.Flags[resolved.FlagKey].Metadata = metadata
	}
	res.Metadata, err = structpb.NewStruct(flagSetMetadata)
	if err != nil {
		s.logger.WarnWithID(reqID, fmt.Sprintf("error resolving all flags: %v", err))
		return nil, fmt.Errorf("error resolving flags. Tracking ID: %s", reqID)
	}

	return connect.NewResponse(res), nil
}

// nolint: dupl
func (s *FlagEvaluationService) EventStream(
	ctx context.Context,
	req *connect.Request[evalV1.EventStreamRequest],
	stream *connect.ServerStream[evalV1.EventStreamResponse],
) error {
	// attach server-side stream deadline to context
	s.logger.Debug("starting event stream for request")

	if s.deadline != 0 {
		streamDeadline := time.Now().Add(s.deadline)
		deadlineCtx, cancel := context.WithDeadline(ctx, streamDeadline)
		ctx = deadlineCtx
		defer cancel()
	}

	s.logger.Debug("starting event stream for request")
	requestNotificationChan := make(chan service.Notification, 1)
	selectorExpression := req.Header().Get(flagdService.FLAGD_SELECTOR_HEADER)
	selector := store.NewSelector(selectorExpression)
	s.eventingConfiguration.Subscribe(ctx, req, &selector, requestNotificationChan)
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
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				s.logger.Debug(fmt.Sprintf("server-side deadline of %s exceeded, exiting stream request with grpc error code 4", s.deadline.String()))
				return connect.NewError(connect.CodeDeadlineExceeded, fmt.Errorf("%s", "stream closed due to server-side timeout"))
			}
			return nil
		}
	}
}

func (s *FlagEvaluationService) ResolveBoolean(
	ctx context.Context,
	req *connect.Request[evalV1.ResolveBooleanRequest],
) (*connect.Response[evalV1.ResolveBooleanResponse], error) {
	ctx, span := s.flagEvalTracer.Start(ctx, "resolveBoolean", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	selectorExpression := req.Header().Get(flagdService.FLAGD_SELECTOR_HEADER)
	selector := store.NewSelector(selectorExpression)
	ctx = context.WithValue(ctx, store.SelectorContextKey{}, selector)

	res := connect.NewResponse(&evalV1.ResolveBooleanResponse{})
	err := resolve(
		ctx,
		s.logger,
		s.eval.ResolveBooleanValue,
		req.Header(),
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&booleanResponse{evalV1Resp: res},
		s.metrics,
		s.contextValues,
		s.headerToContextKeyMappings,
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
	ctx, span := s.flagEvalTracer.Start(ctx, "resolveString", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	selectorExpression := req.Header().Get(flagdService.FLAGD_SELECTOR_HEADER)
	selector := store.NewSelector(selectorExpression)
	ctx = context.WithValue(ctx, store.SelectorContextKey{}, selector)

	res := connect.NewResponse(&evalV1.ResolveStringResponse{})
	err := resolve(
		ctx,
		s.logger,
		s.eval.ResolveStringValue,
		req.Header(),
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&stringResponse{evalV1Resp: res},
		s.metrics,
		s.contextValues,
		s.headerToContextKeyMappings,
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
	ctx, span := s.flagEvalTracer.Start(ctx, "resolveInt", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	selectorExpression := req.Header().Get(flagdService.FLAGD_SELECTOR_HEADER)
	selector := store.NewSelector(selectorExpression)
	ctx = context.WithValue(ctx, store.SelectorContextKey{}, selector)

	res := connect.NewResponse(&evalV1.ResolveIntResponse{})
	err := resolve(
		ctx,
		s.logger,
		s.eval.ResolveIntValue,
		req.Header(),
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&intResponse{evalV1Resp: res},
		s.metrics,
		s.contextValues,
		s.headerToContextKeyMappings,
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
	ctx, span := s.flagEvalTracer.Start(ctx, "resolveFloat", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	selectorExpression := req.Header().Get(flagdService.FLAGD_SELECTOR_HEADER)
	selector := store.NewSelector(selectorExpression)
	ctx = context.WithValue(ctx, store.SelectorContextKey{}, selector)

	res := connect.NewResponse(&evalV1.ResolveFloatResponse{})
	err := resolve(
		ctx,
		s.logger,
		s.eval.ResolveFloatValue,
		req.Header(),
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&floatResponse{evalV1Resp: res},
		s.metrics,
		s.contextValues,
		s.headerToContextKeyMappings,
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
	ctx, span := s.flagEvalTracer.Start(ctx, "resolveObject", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	selectorExpression := req.Header().Get(flagdService.FLAGD_SELECTOR_HEADER)
	selector := store.NewSelector(selectorExpression)
	ctx = context.WithValue(ctx, store.SelectorContextKey{}, selector)

	res := connect.NewResponse(&evalV1.ResolveObjectResponse{})
	err := resolve(
		ctx,
		s.logger,
		s.eval.ResolveObjectValue,
		req.Header(),
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&objectResponse{evalV1Resp: res},
		s.metrics,
		s.contextValues,
		s.headerToContextKeyMappings,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("error evaluating flag with key %s", req.Msg.GetFlagKey()))
	}

	return res, err
}
