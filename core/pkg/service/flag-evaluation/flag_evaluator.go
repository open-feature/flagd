package service

import (
	"context"
	"fmt"
	"time"

	schemaV1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/schema/v1"
	"github.com/bufbuild/connect-go"
	"github.com/open-feature/flagd/core/pkg/eval"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/service"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	"github.com/rs/xid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/structpb"
)

type FlagEvaluationService struct {
	logger                *logger.Logger
	eval                  eval.IEvaluator
	metrics               *telemetry.MetricsRecorder
	eventingConfiguration *eventingConfiguration
}

// NewFlagEvaluationService creates a FlagEvaluationService with provided parameters
func NewFlagEvaluationService(log *logger.Logger,
	eval eval.IEvaluator, eventingCfg *eventingConfiguration, metricsRecorder *telemetry.MetricsRecorder,
) *FlagEvaluationService {
	return &FlagEvaluationService{
		logger:                log,
		eval:                  eval,
		metrics:               metricsRecorder,
		eventingConfiguration: eventingCfg,
	}
}

func (s *FlagEvaluationService) ResolveAll(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveAllRequest],
) (*connect.Response[schemaV1.ResolveAllResponse], error) {
	reqID := xid.New().String()
	defer s.logger.ClearFields(reqID)
	res := &schemaV1.ResolveAllResponse{
		Flags: make(map[string]*schemaV1.AnyFlag),
	}
	values := s.eval.ResolveAllValues(reqID, req.Msg.GetContext())
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

func (s *FlagEvaluationService) EventStream(
	ctx context.Context,
	req *connect.Request[schemaV1.EventStreamRequest],
	stream *connect.ServerStream[schemaV1.EventStreamResponse],
) error {
	requestNotificationChan := make(chan service.Notification, 1)
	s.eventingConfiguration.subscribe(req, requestNotificationChan)
	defer s.eventingConfiguration.unSubscribe(req)

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

func resolve[T constraints](
	goCtx context.Context,
	logger *logger.Logger,
	resolver func(reqID, flagKey string, ctx *structpb.Struct) (T, string, string, error),
	flagKey string,
	ctx *structpb.Struct,
	resp response[T],
	metrics *telemetry.MetricsRecorder,
) error {
	reqID := xid.New().String()
	defer logger.ClearFields(reqID)

	logger.WriteFields(
		reqID,
		zap.String("flag-key", flagKey),
		zap.Strings("context-keys", formatContextKeys(ctx)),
	)

	var evalErrFormatted error
	result, variant, reason, evalErr := resolver(reqID, flagKey, ctx)
	if evalErr != nil {
		logger.WarnWithID(reqID, fmt.Sprintf("returning error response, reason: %v", evalErr))
		reason = model.ErrorReason
		evalErrFormatted = errFormat(evalErr)
	}

	metrics.RecordEvaluation(goCtx, evalErr, reason, variant, flagKey)

	if err := resp.SetResult(result, variant, reason); err != nil && evalErr == nil {
		logger.ErrorWithID(reqID, err.Error())
		return err
	}

	return evalErrFormatted
}

func (s *FlagEvaluationService) ResolveBoolean(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveBooleanRequest],
) (*connect.Response[schemaV1.ResolveBooleanResponse], error) {
	res := connect.NewResponse(&schemaV1.ResolveBooleanResponse{})
	err := resolve[bool](
		ctx,
		s.logger,
		s.eval.ResolveBooleanValue,
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&booleanResponse{res},
		s.metrics,
	)

	return res, err
}

func (s *FlagEvaluationService) ResolveString(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveStringRequest],
) (*connect.Response[schemaV1.ResolveStringResponse], error) {
	res := connect.NewResponse(&schemaV1.ResolveStringResponse{})
	err := resolve[string](
		ctx,
		s.logger,
		s.eval.ResolveStringValue,
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&stringResponse{res},
		s.metrics,
	)

	return res, err
}

func (s *FlagEvaluationService) ResolveInt(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveIntRequest],
) (*connect.Response[schemaV1.ResolveIntResponse], error) {
	res := connect.NewResponse(&schemaV1.ResolveIntResponse{})
	err := resolve[int64](
		ctx,
		s.logger,
		s.eval.ResolveIntValue,
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&intResponse{res},
		s.metrics,
	)

	return res, err
}

func (s *FlagEvaluationService) ResolveFloat(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveFloatRequest],
) (*connect.Response[schemaV1.ResolveFloatResponse], error) {
	res := connect.NewResponse(&schemaV1.ResolveFloatResponse{})
	err := resolve[float64](
		ctx,
		s.logger,
		s.eval.ResolveFloatValue,
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&floatResponse{res},
		s.metrics,
	)

	return res, err
}

func (s *FlagEvaluationService) ResolveObject(
	ctx context.Context,
	req *connect.Request[schemaV1.ResolveObjectRequest],
) (*connect.Response[schemaV1.ResolveObjectResponse], error) {
	res := connect.NewResponse(&schemaV1.ResolveObjectResponse{})
	err := resolve[map[string]any](
		ctx,
		s.logger,
		s.eval.ResolveObjectValue,
		req.Msg.GetFlagKey(),
		req.Msg.GetContext(),
		&objectResponse{res},
		s.metrics,
	)

	return res, err
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
	case model.FlagNotFoundErrorCode:
		return connect.NewError(connect.CodeNotFound, fmt.Errorf("%s, %s", ErrorPrefix, err.Error()))
	case model.TypeMismatchErrorCode:
		return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%s, %s", ErrorPrefix, err.Error()))
	case model.FlagDisabledErrorCode:
		return connect.NewError(connect.CodeUnavailable, fmt.Errorf("%s, %s", ErrorPrefix, err.Error()))
	case model.ParseErrorCode:
		return connect.NewError(connect.CodeDataLoss, fmt.Errorf("%s, %s", ErrorPrefix, err.Error()))
	case model.GeneralErrorCode:
		return connect.NewError(connect.CodeUnknown, fmt.Errorf("%s, %s", ErrorPrefix, err.Error()))
	}

	return err
}
