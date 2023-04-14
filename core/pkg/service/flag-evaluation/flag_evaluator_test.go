package service

import (
	"context"
	"errors"
	"testing"

	schemaV1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/schema/v1"
	"github.com/bufbuild/connect-go"
	"github.com/golang/mock/gomock"
	"github.com/open-feature/flagd/core/pkg/eval"
	mock "github.com/open-feature/flagd/core/pkg/eval/mock"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestConnectService_ResolveAll(t *testing.T) {
	tests := map[string]struct {
		req     *schemaV1.ResolveAllRequest
		evalRes []eval.AnyValue
		wantErr error
		wantRes *schemaV1.ResolveAllResponse
	}{
		"happy-path": {
			req: &schemaV1.ResolveAllRequest{},
			evalRes: []eval.AnyValue{
				{
					Value:   true,
					Variant: "bool-true",
					Reason:  "true",
					FlagKey: "bool",
				},
				{
					Value:   float64(12.12),
					Variant: "float",
					Reason:  "float",
					FlagKey: "float",
				},
				{
					Value:   "hello",
					Variant: "string",
					Reason:  "string",
					FlagKey: "string",
				},
				{
					Value:   "hello",
					Variant: "object",
					Reason:  "string",
					FlagKey: "object",
				},
			},
			wantErr: nil,
			wantRes: &schemaV1.ResolveAllResponse{
				Flags: map[string]*schemaV1.AnyFlag{
					"bool": {
						Value: &schemaV1.AnyFlag_BoolValue{
							BoolValue: true,
						},
						Reason: "STATIC",
					},
					"float": {
						Value: &schemaV1.AnyFlag_DoubleValue{
							DoubleValue: float64(12.12),
						},
						Reason: "STATIC",
					},
					"string": {
						Value: &schemaV1.AnyFlag_StringValue{
							StringValue: "hello",
						},
						Reason: "STATIC",
					},
				},
			},
		},
	}
	ctrl := gomock.NewController(t)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			eval := mock.NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveAllValues(gomock.Any(), gomock.Any(), gomock.Any()).Return(
				tt.evalRes,
			).AnyTimes()
			metrics, exp := getMetricReader()
			s := NewFlagEvaluationService(
				logger.NewLogger(nil, false),
				eval,
				&eventingConfiguration{},
				metrics,
			)
			got, err := s.ResolveAll(context.Background(), connect.NewRequest(tt.req))
			if err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("ConnectService.ResolveAll() error = %v, wantErr %v", err.Error(), tt.wantErr.Error())
				return
			}
			var data metricdata.ResourceMetrics
			err = exp.Collect(context.TODO(), &data)
			require.Nil(t, err)
			// the impression metric is registered
			require.Equal(t, len(data.ScopeMetrics), 1)
			for _, flag := range tt.evalRes {
				switch v := flag.Value.(type) {
				case bool:
					val := got.Msg.Flags[flag.FlagKey].Value.(*schemaV1.AnyFlag_BoolValue)
					require.Equal(t, v, val.BoolValue)
				case string:
					val := got.Msg.Flags[flag.FlagKey].Value.(*schemaV1.AnyFlag_StringValue)
					require.Equal(t, v, val.StringValue)
				case float64:
					val := got.Msg.Flags[flag.FlagKey].Value.(*schemaV1.AnyFlag_DoubleValue)
					require.Equal(t, v, val.DoubleValue)
				}
			}
		})
	}
}

type resolveBooleanArgs struct {
	evalFields   resolveBooleanEvalFields
	functionArgs resolveBooleanFunctionArgs
	want         *schemaV1.ResolveBooleanResponse
	wantErr      error
	mCount       int
}
type resolveBooleanFunctionArgs struct {
	ctx context.Context
	req *schemaV1.ResolveBooleanRequest
}
type resolveBooleanEvalFields struct {
	result  bool
	variant string
	reason  string
}

func TestFlag_Evaluation_ResolveBoolean(t *testing.T) {
	ctrl := gomock.NewController(t)

	tests := map[string]resolveBooleanArgs{
		"happy path": {
			mCount: 1,
			evalFields: resolveBooleanEvalFields{
				result:  true,
				variant: "on",
				reason:  model.DefaultReason,
			},
			functionArgs: resolveBooleanFunctionArgs{
				context.Background(),
				&schemaV1.ResolveBooleanRequest{
					FlagKey: "bool",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveBooleanResponse{
				Value:   true,
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
		"eval returns error": {
			mCount: 1,
			evalFields: resolveBooleanEvalFields{
				result:  true,
				variant: ":(",
				reason:  model.ErrorReason,
			},
			functionArgs: resolveBooleanFunctionArgs{
				context.Background(),
				&schemaV1.ResolveBooleanRequest{
					FlagKey: "bool",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveBooleanResponse{
				Value:   true,
				Variant: ":(",
				Reason:  model.ErrorReason,
			},
			wantErr: errors.New("eval interface error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			eval := mock.NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveBooleanValue(gomock.Any(), gomock.Any(), tt.functionArgs.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.wantErr,
			).AnyTimes()
			metrics, exp := getMetricReader()
			s := NewFlagEvaluationService(
				logger.NewLogger(nil, false),
				eval,
				&eventingConfiguration{},
				metrics,
			)
			got, err := s.ResolveBoolean(tt.functionArgs.ctx, connect.NewRequest(tt.functionArgs.req))
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("Flag_Evaluation.ResolveBoolean() error = %v, wantErr %v", err.Error(), tt.wantErr.Error())
				return
			}
			var data metricdata.ResourceMetrics
			err = exp.Collect(context.TODO(), &data)
			require.Nil(t, err)
			// the impression metric is registered
			require.Equal(t, len(data.ScopeMetrics), tt.mCount)
			require.Equal(t, tt.want, got.Msg)
		})
	}
}

func BenchmarkFlag_Evaluation_ResolveBoolean(b *testing.B) {
	ctrl := gomock.NewController(b)
	tests := map[string]resolveBooleanArgs{
		"happy path": {
			evalFields: resolveBooleanEvalFields{
				result:  true,
				variant: "on",
				reason:  model.DefaultReason,
			},
			functionArgs: resolveBooleanFunctionArgs{
				context.Background(),
				&schemaV1.ResolveBooleanRequest{
					FlagKey: "bool",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveBooleanResponse{
				Value:   true,
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
	}
	for name, tt := range tests {
		eval := mock.NewMockIEvaluator(ctrl)
		eval.EXPECT().ResolveBooleanValue(gomock.Any(), gomock.Any(), tt.functionArgs.req.FlagKey, gomock.Any()).Return(
			tt.evalFields.result,
			tt.evalFields.variant,
			tt.evalFields.reason,
			tt.wantErr,
		).AnyTimes()
		metrics, exp := getMetricReader()
		s := NewFlagEvaluationService(
			logger.NewLogger(nil, false),
			eval,
			&eventingConfiguration{},
			metrics,
		)
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				got, err := s.ResolveBoolean(tt.functionArgs.ctx, connect.NewRequest(tt.functionArgs.req))
				if (err != nil) && !errors.Is(err, tt.wantErr) {
					b.Errorf("Flag_Evaluation.ResolveBoolean() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				require.Equal(b, tt.want, got.Msg)
				var data metricdata.ResourceMetrics
				err = exp.Collect(context.TODO(), &data)
				require.Nil(b, err)
				// the impression metric is registered
				require.Equal(b, len(data.ScopeMetrics), 1)
			}
		})
	}
}

type resolveStringArgs struct {
	evalFields   resolveStringEvalFields
	functionArgs resolveStringFunctionArgs
	want         *schemaV1.ResolveStringResponse
	wantErr      error
	mCount       int
}
type resolveStringFunctionArgs struct {
	ctx context.Context
	req *schemaV1.ResolveStringRequest
}
type resolveStringEvalFields struct {
	result  string
	variant string
	reason  string
}

func TestFlag_Evaluation_ResolveString(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := map[string]resolveStringArgs{
		"happy path": {
			mCount: 1,
			evalFields: resolveStringEvalFields{
				result:  "true",
				variant: "on",
				reason:  model.DefaultReason,
			},
			functionArgs: resolveStringFunctionArgs{
				context.Background(),
				&schemaV1.ResolveStringRequest{
					FlagKey: "string",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveStringResponse{
				Value:   "true",
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
		"eval returns error": {
			mCount: 1,
			evalFields: resolveStringEvalFields{
				result:  "true",
				variant: ":(",
				reason:  model.ErrorReason,
			},
			functionArgs: resolveStringFunctionArgs{
				context.Background(),
				&schemaV1.ResolveStringRequest{
					FlagKey: "string",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveStringResponse{
				Value:   "true",
				Variant: ":(",
				Reason:  model.ErrorReason,
			},
			wantErr: errors.New("eval interface error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			eval := mock.NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveStringValue(
				gomock.Any(), gomock.Any(), tt.functionArgs.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.wantErr,
			)
			metrics, exp := getMetricReader()
			s := NewFlagEvaluationService(
				logger.NewLogger(nil, false),
				eval,
				&eventingConfiguration{},
				metrics,
			)
			got, err := s.ResolveString(tt.functionArgs.ctx, connect.NewRequest(tt.functionArgs.req))
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("Flag_Evaluation.ResolveString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var data metricdata.ResourceMetrics
			err = exp.Collect(context.TODO(), &data)
			require.Nil(t, err)
			// the impression metric is registered
			require.Equal(t, len(data.ScopeMetrics), tt.mCount)
			require.Equal(t, tt.want, got.Msg)
		})
	}
}

func BenchmarkFlag_Evaluation_ResolveString(b *testing.B) {
	ctrl := gomock.NewController(b)
	tests := map[string]resolveStringArgs{
		"happy path": {
			evalFields: resolveStringEvalFields{
				result:  "true",
				variant: "on",
				reason:  model.DefaultReason,
			},
			functionArgs: resolveStringFunctionArgs{
				context.Background(),
				&schemaV1.ResolveStringRequest{
					FlagKey: "string",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveStringResponse{
				Value:   "true",
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
	}
	for name, tt := range tests {
		eval := mock.NewMockIEvaluator(ctrl)
		eval.EXPECT().ResolveStringValue(gomock.Any(), gomock.Any(), tt.functionArgs.req.FlagKey, gomock.Any()).Return(
			tt.evalFields.result,
			tt.evalFields.variant,
			tt.evalFields.reason,
			tt.wantErr,
		).AnyTimes()
		metrics, exp := getMetricReader()
		s := NewFlagEvaluationService(
			logger.NewLogger(nil, false),
			eval,
			&eventingConfiguration{},
			metrics,
		)
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				got, err := s.ResolveString(tt.functionArgs.ctx, connect.NewRequest(tt.functionArgs.req))
				if (err != nil) && !errors.Is(err, tt.wantErr) {
					b.Errorf("Flag_Evaluation.ResolveString() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				require.Equal(b, tt.want, got.Msg)
				var data metricdata.ResourceMetrics
				err = exp.Collect(context.TODO(), &data)
				require.Nil(b, err)
				// the impression metric is registered
				require.Equal(b, len(data.ScopeMetrics), 1)
			}
		})
	}
}

type resolveFloatArgs struct {
	evalFields   resolveFloatEvalFields
	functionArgs resolveFloatFunctionArgs
	want         *schemaV1.ResolveFloatResponse
	wantErr      error
	mCount       int
}
type resolveFloatFunctionArgs struct {
	ctx context.Context
	req *schemaV1.ResolveFloatRequest
}
type resolveFloatEvalFields struct {
	result  float64
	variant string
	reason  string
}

func TestFlag_Evaluation_ResolveFloat(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := map[string]resolveFloatArgs{
		"happy path": {
			mCount: 1,
			evalFields: resolveFloatEvalFields{
				result:  12,
				variant: "on",
				reason:  model.DefaultReason,
			},
			functionArgs: resolveFloatFunctionArgs{
				context.Background(),
				&schemaV1.ResolveFloatRequest{
					FlagKey: "float",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveFloatResponse{
				Value:   12,
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
		"eval returns error": {
			mCount: 1,
			evalFields: resolveFloatEvalFields{
				result:  12,
				variant: ":(",
				reason:  model.ErrorReason,
			},
			functionArgs: resolveFloatFunctionArgs{
				context.Background(),
				&schemaV1.ResolveFloatRequest{
					FlagKey: "float",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveFloatResponse{
				Value:   12,
				Variant: ":(",
				Reason:  model.ErrorReason,
			},
			wantErr: errors.New("eval interface error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			eval := mock.NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveFloatValue(gomock.Any(), gomock.Any(), tt.functionArgs.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.wantErr,
			).AnyTimes()
			metrics, exp := getMetricReader()
			s := NewFlagEvaluationService(
				logger.NewLogger(nil, false),
				eval,
				&eventingConfiguration{},
				metrics,
			)
			got, err := s.ResolveFloat(tt.functionArgs.ctx, connect.NewRequest(tt.functionArgs.req))
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("Flag_Evaluation.ResolveNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got.Msg)
			var data metricdata.ResourceMetrics
			err = exp.Collect(context.TODO(), &data)
			require.Nil(t, err)
			// the impression metric is registered
			require.Equal(t, len(data.ScopeMetrics), tt.mCount)
		})
	}
}

func BenchmarkFlag_Evaluation_ResolveFloat(b *testing.B) {
	ctrl := gomock.NewController(b)
	tests := map[string]resolveFloatArgs{
		"happy path": {
			evalFields: resolveFloatEvalFields{
				result:  12,
				variant: "on",
				reason:  model.DefaultReason,
			},
			functionArgs: resolveFloatFunctionArgs{
				context.Background(),
				&schemaV1.ResolveFloatRequest{
					FlagKey: "float",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveFloatResponse{
				Value:   12,
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
	}
	for name, tt := range tests {
		eval := mock.NewMockIEvaluator(ctrl)
		eval.EXPECT().ResolveFloatValue(gomock.Any(), gomock.Any(), tt.functionArgs.req.FlagKey, gomock.Any()).Return(
			tt.evalFields.result,
			tt.evalFields.variant,
			tt.evalFields.reason,
			tt.wantErr,
		).AnyTimes()
		metrics, exp := getMetricReader()
		s := NewFlagEvaluationService(
			logger.NewLogger(nil, false),
			eval,
			&eventingConfiguration{},
			metrics,
		)
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				got, err := s.ResolveFloat(tt.functionArgs.ctx, connect.NewRequest(tt.functionArgs.req))
				if (err != nil) && !errors.Is(err, tt.wantErr) {
					b.Errorf("Flag_Evaluation.ResolveNumber() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				require.Equal(b, tt.want, got.Msg)
				var data metricdata.ResourceMetrics
				err = exp.Collect(context.TODO(), &data)
				require.Nil(b, err)
				// the impression metric is registered
				require.Equal(b, len(data.ScopeMetrics), 1)
			}
		})
	}
}

type resolveIntArgs struct {
	evalFields   resolveIntEvalFields
	functionArgs resolveIntFunctionArgs
	want         *schemaV1.ResolveIntResponse
	wantErr      error
	mCount       int
}
type resolveIntFunctionArgs struct {
	ctx context.Context
	req *schemaV1.ResolveIntRequest
}
type resolveIntEvalFields struct {
	result  int64
	variant string
	reason  string
}

func TestFlag_Evaluation_ResolveInt(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := map[string]resolveIntArgs{
		"happy path": {
			mCount: 1,
			evalFields: resolveIntEvalFields{
				result:  12,
				variant: "on",
				reason:  model.DefaultReason,
			},
			functionArgs: resolveIntFunctionArgs{
				context.Background(),
				&schemaV1.ResolveIntRequest{
					FlagKey: "int",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveIntResponse{
				Value:   12,
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
		"eval returns error": {
			mCount: 1,
			evalFields: resolveIntEvalFields{
				result:  12,
				variant: ":(",
				reason:  model.ErrorReason,
			},
			functionArgs: resolveIntFunctionArgs{
				context.Background(),
				&schemaV1.ResolveIntRequest{
					FlagKey: "int",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveIntResponse{
				Value:   12,
				Variant: ":(",
				Reason:  model.ErrorReason,
			},
			wantErr: errors.New("eval interface error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			eval := mock.NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveIntValue(gomock.Any(), gomock.Any(), tt.functionArgs.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.wantErr,
			).AnyTimes()
			metrics, exp := getMetricReader()
			s := NewFlagEvaluationService(
				logger.NewLogger(nil, false),
				eval,
				&eventingConfiguration{},
				metrics,
			)
			got, err := s.ResolveInt(tt.functionArgs.ctx, connect.NewRequest(tt.functionArgs.req))
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("Flag_Evaluation.ResolveNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got.Msg)
			var data metricdata.ResourceMetrics
			err = exp.Collect(context.TODO(), &data)
			require.Nil(t, err)
			// the impression metric is registered
			require.Equal(t, len(data.ScopeMetrics), tt.mCount)
		})
	}
}

func BenchmarkFlag_Evaluation_ResolveInt(b *testing.B) {
	ctrl := gomock.NewController(b)
	tests := map[string]resolveIntArgs{
		"happy path": {
			evalFields: resolveIntEvalFields{
				result:  12,
				variant: "on",
				reason:  model.DefaultReason,
			},
			functionArgs: resolveIntFunctionArgs{
				context.Background(),
				&schemaV1.ResolveIntRequest{
					FlagKey: "int",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveIntResponse{
				Value:   12,
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
	}
	for name, tt := range tests {
		eval := mock.NewMockIEvaluator(ctrl)
		eval.EXPECT().ResolveIntValue(gomock.Any(), gomock.Any(), tt.functionArgs.req.FlagKey, gomock.Any()).Return(
			tt.evalFields.result,
			tt.evalFields.variant,
			tt.evalFields.reason,
			tt.wantErr,
		).AnyTimes()
		metrics, exp := getMetricReader()
		s := NewFlagEvaluationService(
			logger.NewLogger(nil, false),
			eval,
			&eventingConfiguration{},
			metrics,
		)
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				got, err := s.ResolveInt(tt.functionArgs.ctx, connect.NewRequest(tt.functionArgs.req))
				if (err != nil) && !errors.Is(err, tt.wantErr) {
					b.Errorf("Flag_Evaluation.ResolveNumber() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				require.Equal(b, tt.want, got.Msg)
				var data metricdata.ResourceMetrics
				err = exp.Collect(context.TODO(), &data)
				require.Nil(b, err)
				// the impression metric is registered
				require.Equal(b, len(data.ScopeMetrics), 1)
			}
		})
	}
}

type resolveObjectArgs struct {
	evalFields   resolveObjectEvalFields
	functionArgs resolveObjectFunctionArgs
	want         *schemaV1.ResolveObjectResponse
	wantErr      error
	mCount       int
}
type resolveObjectFunctionArgs struct {
	ctx context.Context
	req *schemaV1.ResolveObjectRequest
}
type resolveObjectEvalFields struct {
	result  map[string]interface{}
	variant string
	reason  string
}

func TestFlag_Evaluation_ResolveObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := map[string]resolveObjectArgs{
		"happy path": {
			mCount: 1,
			evalFields: resolveObjectEvalFields{
				result: map[string]interface{}{
					"food": "bars",
				},
				variant: "on",
				reason:  model.DefaultReason,
			},
			functionArgs: resolveObjectFunctionArgs{
				context.Background(),
				&schemaV1.ResolveObjectRequest{
					FlagKey: "object",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveObjectResponse{
				Value:   nil,
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
		"eval returns error": {
			mCount: 1,
			evalFields: resolveObjectEvalFields{
				result: map[string]interface{}{
					"food": "bars",
				},
				variant: ":(",
				reason:  model.ErrorReason,
			},
			functionArgs: resolveObjectFunctionArgs{
				context.Background(),
				&schemaV1.ResolveObjectRequest{
					FlagKey: "object",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveObjectResponse{
				Variant: ":(",
				Reason:  model.ErrorReason,
			},
			wantErr: errors.New("eval interface error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			eval := mock.NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveObjectValue(gomock.Any(), gomock.Any(), tt.functionArgs.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.wantErr,
			).AnyTimes()
			metrics, exp := getMetricReader()
			s := NewFlagEvaluationService(
				logger.NewLogger(nil, false),
				eval,
				&eventingConfiguration{},
				metrics,
			)

			outParsed, err := structpb.NewStruct(tt.evalFields.result)
			if err != nil {
				t.Error(err)
			}
			tt.want.Value = outParsed
			got, err := s.ResolveObject(tt.functionArgs.ctx, connect.NewRequest(tt.functionArgs.req))
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("Flag_Evaluation.ResolveObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got.Msg)
			var data metricdata.ResourceMetrics
			err = exp.Collect(context.TODO(), &data)
			require.Nil(t, err)
			// the impression metric is registered
			require.Equal(t, len(data.ScopeMetrics), tt.mCount)
		})
	}
}

func BenchmarkFlag_Evaluation_ResolveObject(b *testing.B) {
	ctrl := gomock.NewController(b)
	tests := map[string]resolveObjectArgs{
		"happy path": {
			evalFields: resolveObjectEvalFields{
				result: map[string]interface{}{
					"food": "bars",
				},
				variant: "on",
				reason:  model.DefaultReason,
			},
			functionArgs: resolveObjectFunctionArgs{
				context.Background(),
				&schemaV1.ResolveObjectRequest{
					FlagKey: "object",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveObjectResponse{
				Value:   nil,
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
	}
	for name, tt := range tests {
		eval := mock.NewMockIEvaluator(ctrl)
		eval.EXPECT().ResolveObjectValue(gomock.Any(), gomock.Any(), tt.functionArgs.req.FlagKey, gomock.Any()).Return(
			tt.evalFields.result,
			tt.evalFields.variant,
			tt.evalFields.reason,
			tt.wantErr,
		).AnyTimes()
		metrics, exp := getMetricReader()
		s := NewFlagEvaluationService(
			logger.NewLogger(nil, false),
			eval,
			&eventingConfiguration{},
			metrics,
		)
		if name != "eval returns error" {
			outParsed, err := structpb.NewStruct(tt.evalFields.result)
			if err != nil {
				b.Error(err)
			}
			tt.want.Value = outParsed
		}
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				got, err := s.ResolveObject(tt.functionArgs.ctx, connect.NewRequest(tt.functionArgs.req))
				if (err != nil) && !errors.Is(err, tt.wantErr) {
					b.Errorf("Flag_Evaluation.ResolveObject() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				require.Equal(b, tt.want, got.Msg)
				var data metricdata.ResourceMetrics
				err = exp.Collect(context.TODO(), &data)
				require.Nil(b, err)
				// the impression metric is registered
				require.Equal(b, len(data.ScopeMetrics), 1)
			}
		})
	}
}

func getMetricReader() (*telemetry.MetricsRecorder, metric.Reader) {
	exp := metric.NewManualReader()
	rs := resource.NewWithAttributes("testSchema")
	return telemetry.NewOTelRecorder(exp, rs, "testSvc"), exp
}
