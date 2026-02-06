package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	evalV2 "buf.build/gen/go/open-feature-forking/flagd/protocolbuffers/go/flagd/evaluation/v2"
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
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/structpb"
)

type FlagEvaluationServiceV2 struct {
	logger                     *logger.Logger
	eval                       evaluator.IEvaluator
	metrics                    telemetry.IMetricsRecorder
	eventingConfiguration      IEvents
	flagEvalTracer             trace.Tracer
	contextValues              map[string]any
	headerToContextKeyMappings map[string]string
	deadline                   time.Duration
}

// NewFlagEvaluationServiceV2 creates a FlagEvaluationServiceV2 with provided parameters
func NewFlagEvaluationServiceV2(log *logger.Logger,
	eval evaluator.IEvaluator,
	eventingCfg IEvents,
	metricsRecorder telemetry.IMetricsRecorder,
	contextValues map[string]any,
	headerToContextKeyMappings map[string]string,
	streamDeadline time.Duration,
) *FlagEvaluationServiceV2 {
	svc := &FlagEvaluationServiceV2{
		logger:                     log,
		eval:                       eval,
		metrics:                    &telemetry.NoopMetricsRecorder{},
		eventingConfiguration:      eventingCfg,
		flagEvalTracer:             otel.Tracer("flagd.evaluation.v2"),
		contextValues:              contextValues,
		headerToContextKeyMappings: headerToContextKeyMappings,
		deadline:                   streamDeadline,
	}

	if metricsRecorder != nil {
		svc.metrics = metricsRecorder
	}

	return svc
}

// nolint: dupl
func (s *FlagEvaluationServiceV2) EventStream(
	ctx context.Context,
	req *connect.Request[evalV2.EventStreamRequest],
	stream *connect.ServerStream[evalV2.EventStreamResponse],
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
			err := stream.Send(&evalV2.EventStreamResponse{
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
			err = stream.Send(&evalV2.EventStreamResponse{
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

func (s *FlagEvaluationServiceV2) ResolveBoolean(
	ctx context.Context,
	req *connect.Request[evalV2.ResolveBooleanRequest],
) (*connect.Response[evalV2.ResolveBooleanResponse], error) {
	ctx, span := s.flagEvalTracer.Start(ctx, "resolveBoolean", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	selectorExpression := req.Header().Get(flagdService.FLAGD_SELECTOR_HEADER)
	selector := store.NewSelector(selectorExpression)
	ctx = context.WithValue(ctx, store.SelectorContextKey{}, selector)

	res := connect.NewResponse(&evalV2.ResolveBooleanResponse{})
	err := resolveV2(
		ctx,
		s.logger,
		s.eval.ResolveBooleanValue,
		req.Header(),
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&booleanResponseV2{evalV2Resp: res},
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

func (s *FlagEvaluationServiceV2) ResolveString(
	ctx context.Context,
	req *connect.Request[evalV2.ResolveStringRequest],
) (*connect.Response[evalV2.ResolveStringResponse], error) {
	ctx, span := s.flagEvalTracer.Start(ctx, "resolveString", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	selectorExpression := req.Header().Get(flagdService.FLAGD_SELECTOR_HEADER)
	selector := store.NewSelector(selectorExpression)
	ctx = context.WithValue(ctx, store.SelectorContextKey{}, selector)

	res := connect.NewResponse(&evalV2.ResolveStringResponse{})
	err := resolveV2(
		ctx,
		s.logger,
		s.eval.ResolveStringValue,
		req.Header(),
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&stringResponseV2{evalV2Resp: res},
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

func (s *FlagEvaluationServiceV2) ResolveInt(
	ctx context.Context,
	req *connect.Request[evalV2.ResolveIntRequest],
) (*connect.Response[evalV2.ResolveIntResponse], error) {
	ctx, span := s.flagEvalTracer.Start(ctx, "resolveInt", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	selectorExpression := req.Header().Get(flagdService.FLAGD_SELECTOR_HEADER)
	selector := store.NewSelector(selectorExpression)
	ctx = context.WithValue(ctx, store.SelectorContextKey{}, selector)

	res := connect.NewResponse(&evalV2.ResolveIntResponse{})
	err := resolveV2(
		ctx,
		s.logger,
		s.eval.ResolveIntValue,
		req.Header(),
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&intResponseV2{evalV2Resp: res},
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

func (s *FlagEvaluationServiceV2) ResolveFloat(
	ctx context.Context,
	req *connect.Request[evalV2.ResolveFloatRequest],
) (*connect.Response[evalV2.ResolveFloatResponse], error) {
	ctx, span := s.flagEvalTracer.Start(ctx, "resolveFloat", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	selectorExpression := req.Header().Get(flagdService.FLAGD_SELECTOR_HEADER)
	selector := store.NewSelector(selectorExpression)
	ctx = context.WithValue(ctx, store.SelectorContextKey{}, selector)

	res := connect.NewResponse(&evalV2.ResolveFloatResponse{})
	err := resolveV2(
		ctx,
		s.logger,
		s.eval.ResolveFloatValue,
		req.Header(),
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&floatResponseV2{evalV2Resp: res},
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

func (s *FlagEvaluationServiceV2) ResolveObject(
	ctx context.Context,
	req *connect.Request[evalV2.ResolveObjectRequest],
) (*connect.Response[evalV2.ResolveObjectResponse], error) {
	ctx, span := s.flagEvalTracer.Start(ctx, "resolveObject", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	selectorExpression := req.Header().Get(flagdService.FLAGD_SELECTOR_HEADER)
	selector := store.NewSelector(selectorExpression)
	ctx = context.WithValue(ctx, store.SelectorContextKey{}, selector)

	res := connect.NewResponse(&evalV2.ResolveObjectResponse{})
	err := resolveV2(
		ctx,
		s.logger,
		s.eval.ResolveObjectValue,
		req.Header(),
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&objectResponseV2{evalV2Resp: res},
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

func resolveV2[T constraints](ctx context.Context, logger *logger.Logger, resolver resolverSignature[T], header http.Header, flagKey string,
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
		if evalErr.Error() == model.FlagNotFoundErrorCode || evalErr.Error() == model.ParseErrorCode {
			logger.WarnWithID(reqID, fmt.Sprintf("flag not found or parse error, returning DEFAULT reason: %v", evalErr))
			reason = model.DefaultReason
		} else {
			logger.WarnWithID(reqID, fmt.Sprintf("returning error response, reason: %v", evalErr))
			reason = model.ErrorReason
			evalErrFormatted = errFormatV2(evalErr)
		}
	}

	if metrics != nil {
		metrics.RecordEvaluation(ctx, evalErr, reason, variant, flagKey)
	}

	spanFromContext := trace.SpanFromContext(ctx)
	spanFromContext.SetAttributes(telemetry.SemConvFeatureFlagAttributes(flagKey, variant)...)

	// For V2, when we have FLAG_NOT_FOUND or PARSE_ERROR, only set reason (no value/variant)
	if evalErr != nil && (evalErr.Error() == model.FlagNotFoundErrorCode || evalErr.Error() == model.ParseErrorCode) {
		if respV2, ok := resp.(responseV2[T]); ok {
			if err := respV2.SetReasonOnly(reason, metadata); err != nil {
				logger.ErrorWithID(reqID, err.Error())
				return fmt.Errorf("error setting response result: %w", err)
			}
		} else if err := resp.SetResult(result, variant, reason, metadata); err != nil {
			logger.ErrorWithID(reqID, err.Error())
			return fmt.Errorf("error setting response result: %w", err)
		}
	} else {
		if err := resp.SetResult(result, variant, reason, metadata); err != nil && evalErr == nil {
			logger.ErrorWithID(reqID, err.Error())
			return fmt.Errorf("error setting response result: %w", err)
		}
	}

	return evalErrFormatted
}

// errFormatV2 formats errors for V2 API, excluding FLAG_NOT_FOUND and PARSE_ERROR which are not errors in V2
func errFormatV2(err error) error {
	ReadableErrorMsg := model.GetErrorMessage(err.Error())
	switch err.Error() {
	case model.FlagDisabledErrorCode:
		return connect.NewError(connect.CodeNotFound, fmt.Errorf("%s", ReadableErrorMsg))
	case model.TypeMismatchErrorCode:
		return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%s", ReadableErrorMsg))
	case model.GeneralErrorCode:
		return connect.NewError(connect.CodeUnknown, fmt.Errorf("%s", ReadableErrorMsg))
	}

	return err
}
