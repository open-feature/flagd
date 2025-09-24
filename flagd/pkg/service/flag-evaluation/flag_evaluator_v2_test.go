package service

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"

	evalV1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/evaluation/v1"
	"connectrpc.com/connect"
	"github.com/open-feature/flagd/core/pkg/evaluator"
	mock "github.com/open-feature/flagd/core/pkg/evaluator/mock"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestConnectServiceV2_ResolveAll(t *testing.T) {
	tests := map[string]struct {
		req         *evalV1.ResolveAllRequest
		evalRes     []evaluator.AnyValue
		metadataRes model.Metadata
		evalErr     error
		wantErr     bool
		wantRes     *evalV1.ResolveAllResponse
	}{
		"happy-path": {
			req: &evalV1.ResolveAllRequest{},
			evalRes: []evaluator.AnyValue{
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
			metadataRes: model.Metadata{
				"key": "value",
			},
			wantRes: &evalV1.ResolveAllResponse{
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"key": structpb.NewStringValue("value"),
					},
				},
				Flags: map[string]*evalV1.AnyFlag{
					"bool": {
						Value: &evalV1.AnyFlag_BoolValue{
							BoolValue: true,
						},
						Reason: "STATIC",
					},
					"float": {
						Value: &evalV1.AnyFlag_DoubleValue{
							DoubleValue: float64(12.12),
						},
						Reason: "STATIC",
					},
					"string": {
						Value: &evalV1.AnyFlag_StringValue{
							StringValue: "hello",
						},
						Reason: "STATIC",
					},
				},
			},
		},
		"resolver error": {
			req:     &evalV1.ResolveAllRequest{},
			evalRes: []evaluator.AnyValue{},
			evalErr: errors.New("some error from internal evaluator"),
			wantErr: true,
		},
	}
	ctrl := gomock.NewController(t)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			eval := mock.NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveAllValues(gomock.Any(), gomock.Any(), gomock.Any()).Return(
				tt.evalRes, tt.metadataRes, tt.evalErr,
			).AnyTimes()

			metrics, exp := getMetricReader()
			s := NewFlagEvaluationService(logger.NewLogger(nil, false), eval, &eventingConfiguration{}, metrics, nil, nil, 0, "")

			// when
			got, err := s.ResolveAll(context.Background(), connect.NewRequest(tt.req))

			// then
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but git none")
				}

				return
			}

			var data metricdata.ResourceMetrics
			err = exp.Collect(context.TODO(), &data)
			require.Nil(t, err)
			// the impression metric is registered
			require.Equal(t, len(data.ScopeMetrics), 1)
			require.EqualValues(t, tt.wantRes.Metadata, got.Msg.Metadata)
			for _, flag := range tt.evalRes {
				switch v := flag.Value.(type) {
				case bool:
					val := got.Msg.Flags[flag.FlagKey].Value.(*evalV1.AnyFlag_BoolValue)
					require.Equal(t, v, val.BoolValue)
				case string:
					val := got.Msg.Flags[flag.FlagKey].Value.(*evalV1.AnyFlag_StringValue)
					require.Equal(t, v, val.StringValue)
				case float64:
					val := got.Msg.Flags[flag.FlagKey].Value.(*evalV1.AnyFlag_DoubleValue)
					require.Equal(t, v, val.DoubleValue)
				}
			}
		})
	}
}

type resolveBooleanArgsV2 struct {
	evalFields   resolveBooleanEvalFieldsV2
	functionArgs resolveBooleanFunctionArgsV2
	want         *evalV1.ResolveBooleanResponse
	wantErr      error
	mCount       int
}
type resolveBooleanFunctionArgsV2 struct {
	ctx context.Context
	req *evalV1.ResolveBooleanRequest
}
type resolveBooleanEvalFieldsV2 struct {
	result bool
	evalCommons
}

func TestFlag_EvaluationV2_ResolveBoolean(t *testing.T) {
	ctrl := gomock.NewController(t)

	tests := map[string]resolveBooleanArgsV2{
		"happy path": {
			mCount: 1,
			evalFields: resolveBooleanEvalFieldsV2{
				result:      true,
				evalCommons: happyCommon,
			},
			functionArgs: resolveBooleanFunctionArgsV2{
				context.Background(),
				&evalV1.ResolveBooleanRequest{
					FlagKey: "bool",
					Context: &structpb.Struct{},
				},
			},
			want: &evalV1.ResolveBooleanResponse{
				Value:    true,
				Reason:   model.DefaultReason,
				Variant:  "on",
				Metadata: responseStruct,
			},
			wantErr: nil,
		},
		"eval returns error": {
			mCount: 1,
			evalFields: resolveBooleanEvalFieldsV2{
				result:      true,
				evalCommons: sadCommon,
			},
			functionArgs: resolveBooleanFunctionArgsV2{
				context.Background(),
				&evalV1.ResolveBooleanRequest{
					FlagKey: "bool",
					Context: &structpb.Struct{},
				},
			},
			want: &evalV1.ResolveBooleanResponse{
				Value:    true,
				Variant:  ":(",
				Reason:   model.ErrorReason,
				Metadata: responseStruct,
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
				tt.evalFields.metadata,
				tt.wantErr,
			).AnyTimes()
			metrics, exp := getMetricReader()
			s := NewFlagEvaluationService(
				logger.NewLogger(nil, false),
				eval,
				&eventingConfiguration{},
				metrics,
				nil,
				nil,
				0,
				"",
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

func BenchmarkFlag_EvaluationV2_ResolveBoolean(b *testing.B) {
	ctrl := gomock.NewController(b)
	tests := map[string]resolveBooleanArgsV2{
		"happy path": {
			evalFields: resolveBooleanEvalFieldsV2{
				result:      true,
				evalCommons: happyCommon,
			},
			functionArgs: resolveBooleanFunctionArgsV2{
				context.Background(),
				&evalV1.ResolveBooleanRequest{
					FlagKey: "bool",
					Context: &structpb.Struct{},
				},
			},
			want: &evalV1.ResolveBooleanResponse{
				Value:    true,
				Reason:   model.DefaultReason,
				Variant:  "on",
				Metadata: responseStruct,
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
			tt.evalFields.metadata,
			tt.wantErr,
		).AnyTimes()
		metrics, exp := getMetricReader()
		s := NewFlagEvaluationService(
			logger.NewLogger(nil, false),
			eval,
			&eventingConfiguration{},
			metrics,
			nil,
			nil,
			0,
			"",
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

type resolveStringArgsV2 struct {
	evalFields   resolveStringEvalFieldsV2
	functionArgs resolveStringFunctionArgsV2
	want         *evalV1.ResolveStringResponse
	wantErr      error
	mCount       int
}
type resolveStringFunctionArgsV2 struct {
	ctx context.Context
	req *evalV1.ResolveStringRequest
}
type resolveStringEvalFieldsV2 struct {
	result string
	evalCommons
}

func TestFlag_EvaluationV2_ResolveString(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := map[string]resolveStringArgsV2{
		"happy path": {
			mCount: 1,
			evalFields: resolveStringEvalFieldsV2{
				result:      "true",
				evalCommons: happyCommon,
			},
			functionArgs: resolveStringFunctionArgsV2{
				context.Background(),
				&evalV1.ResolveStringRequest{
					FlagKey: "string",
					Context: &structpb.Struct{},
				},
			},
			want: &evalV1.ResolveStringResponse{
				Value:    "true",
				Reason:   model.DefaultReason,
				Variant:  "on",
				Metadata: responseStruct,
			},
			wantErr: nil,
		},
		"eval returns error": {
			mCount: 1,
			evalFields: resolveStringEvalFieldsV2{
				result:      "true",
				evalCommons: sadCommon,
			},
			functionArgs: resolveStringFunctionArgsV2{
				context.Background(),
				&evalV1.ResolveStringRequest{
					FlagKey: "string",
					Context: &structpb.Struct{},
				},
			},
			want: &evalV1.ResolveStringResponse{
				Value:    "true",
				Variant:  ":(",
				Reason:   model.ErrorReason,
				Metadata: responseStruct,
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
				tt.evalFields.metadata,
				tt.wantErr,
			)
			metrics, exp := getMetricReader()
			s := NewFlagEvaluationService(
				logger.NewLogger(nil, false),
				eval,
				&eventingConfiguration{},
				metrics,
				nil,
				nil,
				0,
				"",
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

func BenchmarkFlag_EvaluationV2_ResolveString(b *testing.B) {
	ctrl := gomock.NewController(b)
	tests := map[string]resolveStringArgsV2{
		"happy path": {
			evalFields: resolveStringEvalFieldsV2{
				result:      "true",
				evalCommons: happyCommon,
			},
			functionArgs: resolveStringFunctionArgsV2{
				context.Background(),
				&evalV1.ResolveStringRequest{
					FlagKey: "string",
					Context: &structpb.Struct{},
				},
			},
			want: &evalV1.ResolveStringResponse{
				Value:    "true",
				Reason:   model.DefaultReason,
				Variant:  "on",
				Metadata: responseStruct,
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
			tt.evalFields.metadata,
			tt.wantErr,
		).AnyTimes()
		metrics, exp := getMetricReader()
		s := NewFlagEvaluationService(
			logger.NewLogger(nil, false),
			eval,
			&eventingConfiguration{},
			metrics,
			nil,
			nil,
			0,
			"",
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

type resolveFloatArgsV2 struct {
	evalFields   resolveFloatEvalFieldsV2
	functionArgs resolveFloatFunctionArgsV2
	want         *evalV1.ResolveFloatResponse
	wantErr      error
	mCount       int
}
type resolveFloatFunctionArgsV2 struct {
	ctx context.Context
	req *evalV1.ResolveFloatRequest
}
type resolveFloatEvalFieldsV2 struct {
	result float64
	evalCommons
}

func TestFlag_EvaluationV2_ResolveFloat(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := map[string]resolveFloatArgsV2{
		"happy path": {
			mCount: 1,
			evalFields: resolveFloatEvalFieldsV2{
				result:      12,
				evalCommons: happyCommon,
			},
			functionArgs: resolveFloatFunctionArgsV2{
				context.Background(),
				&evalV1.ResolveFloatRequest{
					FlagKey: "float",
					Context: &structpb.Struct{},
				},
			},
			want: &evalV1.ResolveFloatResponse{
				Value:    12,
				Reason:   model.DefaultReason,
				Variant:  "on",
				Metadata: responseStruct,
			},
			wantErr: nil,
		},
		"eval returns error": {
			mCount: 1,
			evalFields: resolveFloatEvalFieldsV2{
				result:      12,
				evalCommons: sadCommon,
			},
			functionArgs: resolveFloatFunctionArgsV2{
				context.Background(),
				&evalV1.ResolveFloatRequest{
					FlagKey: "float",
					Context: &structpb.Struct{},
				},
			},
			want: &evalV1.ResolveFloatResponse{
				Value:    12,
				Variant:  ":(",
				Reason:   model.ErrorReason,
				Metadata: responseStruct,
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
				tt.evalFields.metadata,
				tt.wantErr,
			).AnyTimes()
			metrics, exp := getMetricReader()
			s := NewFlagEvaluationService(
				logger.NewLogger(nil, false),
				eval,
				&eventingConfiguration{},
				metrics,
				nil,
				nil,
				0,
				"",
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

func BenchmarkFlag_EvaluationV2_ResolveFloat(b *testing.B) {
	ctrl := gomock.NewController(b)
	tests := map[string]resolveFloatArgsV2{
		"happy path": {
			evalFields: resolveFloatEvalFieldsV2{
				result:      12,
				evalCommons: happyCommon,
			},
			functionArgs: resolveFloatFunctionArgsV2{
				context.Background(),
				&evalV1.ResolveFloatRequest{
					FlagKey: "float",
					Context: &structpb.Struct{},
				},
			},
			want: &evalV1.ResolveFloatResponse{
				Value:    12,
				Reason:   model.DefaultReason,
				Variant:  "on",
				Metadata: responseStruct,
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
			tt.evalFields.metadata,
			tt.wantErr,
		).AnyTimes()
		metrics, exp := getMetricReader()
		s := NewFlagEvaluationService(
			logger.NewLogger(nil, false),
			eval,
			&eventingConfiguration{},
			metrics,
			nil,
			nil,
			0,
			"",
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

type resolveIntArgsV2 struct {
	evalFields   resolveIntEvalFieldsV2
	functionArgs resolveIntFunctionArgsV2
	want         *evalV1.ResolveIntResponse
	wantErr      error
	mCount       int
}
type resolveIntFunctionArgsV2 struct {
	ctx context.Context
	req *evalV1.ResolveIntRequest
}
type resolveIntEvalFieldsV2 struct {
	result int64
	evalCommons
}

func TestFlag_EvaluationV2_ResolveInt(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := map[string]resolveIntArgsV2{
		"happy path": {
			mCount: 1,
			evalFields: resolveIntEvalFieldsV2{
				result:      12,
				evalCommons: happyCommon,
			},
			functionArgs: resolveIntFunctionArgsV2{
				context.Background(),
				&evalV1.ResolveIntRequest{
					FlagKey: "int",
					Context: &structpb.Struct{},
				},
			},
			want: &evalV1.ResolveIntResponse{
				Value:    12,
				Reason:   model.DefaultReason,
				Variant:  "on",
				Metadata: responseStruct,
			},
			wantErr: nil,
		},
		"eval returns error": {
			mCount: 1,
			evalFields: resolveIntEvalFieldsV2{
				result:      12,
				evalCommons: sadCommon,
			},
			functionArgs: resolveIntFunctionArgsV2{
				context.Background(),
				&evalV1.ResolveIntRequest{
					FlagKey: "int",
					Context: &structpb.Struct{},
				},
			},
			want: &evalV1.ResolveIntResponse{
				Value:    12,
				Variant:  ":(",
				Reason:   model.ErrorReason,
				Metadata: responseStruct,
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
				tt.evalFields.metadata,
				tt.wantErr,
			).AnyTimes()
			metrics, exp := getMetricReader()
			s := NewFlagEvaluationService(
				logger.NewLogger(nil, false),
				eval,
				&eventingConfiguration{},
				metrics,
				nil,
				nil,
				0,
				"",
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

func BenchmarkFlag_EvaluationV2_ResolveInt(b *testing.B) {
	ctrl := gomock.NewController(b)
	tests := map[string]resolveIntArgsV2{
		"happy path": {
			evalFields: resolveIntEvalFieldsV2{
				result:      12,
				evalCommons: happyCommon,
			},
			functionArgs: resolveIntFunctionArgsV2{
				context.Background(),
				&evalV1.ResolveIntRequest{
					FlagKey: "int",
					Context: &structpb.Struct{},
				},
			},
			want: &evalV1.ResolveIntResponse{
				Value:    12,
				Reason:   model.DefaultReason,
				Variant:  "on",
				Metadata: responseStruct,
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
			tt.evalFields.metadata,
			tt.wantErr,
		).AnyTimes()
		metrics, exp := getMetricReader()
		s := NewFlagEvaluationService(
			logger.NewLogger(nil, false),
			eval,
			&eventingConfiguration{},
			metrics,
			nil,
			nil,
			0,
			"",
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

type resolveObjectArgsV2 struct {
	evalFields   resolveObjectEvalFieldsV2
	functionArgs resolveObjectFunctionArgsV2
	want         *evalV1.ResolveObjectResponse
	wantErr      error
	mCount       int
}
type resolveObjectFunctionArgsV2 struct {
	ctx context.Context
	req *evalV1.ResolveObjectRequest
}
type resolveObjectEvalFieldsV2 struct {
	result map[string]interface{}
	evalCommons
}

func TestFlag_EvaluationV2_ResolveObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := map[string]resolveObjectArgsV2{
		"happy path": {
			mCount: 1,
			evalFields: resolveObjectEvalFieldsV2{
				result: map[string]interface{}{
					"food": "bars",
				},
				evalCommons: happyCommon,
			},
			functionArgs: resolveObjectFunctionArgsV2{
				context.Background(),
				&evalV1.ResolveObjectRequest{
					FlagKey: "object",
					Context: &structpb.Struct{},
				},
			},
			want: &evalV1.ResolveObjectResponse{
				Value:    nil,
				Reason:   model.DefaultReason,
				Variant:  "on",
				Metadata: responseStruct,
			},
			wantErr: nil,
		},
		"eval returns error": {
			mCount: 1,
			evalFields: resolveObjectEvalFieldsV2{
				result: map[string]interface{}{
					"food": "bars",
				},
				evalCommons: sadCommon,
			},
			functionArgs: resolveObjectFunctionArgsV2{
				context.Background(),
				&evalV1.ResolveObjectRequest{
					FlagKey: "object",
					Context: &structpb.Struct{},
				},
			},
			want: &evalV1.ResolveObjectResponse{
				Variant:  ":(",
				Reason:   model.ErrorReason,
				Metadata: responseStruct,
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
				tt.evalFields.metadata,
				tt.wantErr,
			).AnyTimes()
			metrics, exp := getMetricReader()
			s := NewFlagEvaluationService(
				logger.NewLogger(nil, false),
				eval,
				&eventingConfiguration{},
				metrics,
				nil,
				nil,
				0,
				"",
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

func BenchmarkFlag_EvaluationV2_ResolveObject(b *testing.B) {
	ctrl := gomock.NewController(b)
	tests := map[string]resolveObjectArgsV2{
		"happy path": {
			evalFields: resolveObjectEvalFieldsV2{
				result: map[string]interface{}{
					"food": "bars",
				},
				evalCommons: happyCommon,
			},
			functionArgs: resolveObjectFunctionArgsV2{
				context.Background(),
				&evalV1.ResolveObjectRequest{
					FlagKey: "object",
					Context: &structpb.Struct{},
				},
			},
			want: &evalV1.ResolveObjectResponse{
				Value:    nil,
				Reason:   model.DefaultReason,
				Variant:  "on",
				Metadata: responseStruct,
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
			tt.evalFields.metadata,
			tt.wantErr,
		).AnyTimes()
		metrics, exp := getMetricReader()
		s := NewFlagEvaluationService(
			logger.NewLogger(nil, false),
			eval,
			&eventingConfiguration{},
			metrics,
			nil,
			nil,
			0,
			"",
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

// TestFlag_EvaluationV2_ErrorCodes test validate error mapping from known errors to connect.Code and avoid accidental
// changes. This is essential as SDK implementations rely on connect. Code to differentiate GRPC errors vs Flag errors.
// For any change in error codes, we must change respective SDK.
func TestFlag_EvaluationV2_ErrorCodes(t *testing.T) {
	tests := []struct {
		err  error
		code connect.Code
	}{
		{
			err:  errors.New(model.FlagNotFoundErrorCode),
			code: connect.CodeNotFound,
		},
		{
			err:  errors.New(model.TypeMismatchErrorCode),
			code: connect.CodeInvalidArgument,
		},
		{
			err:  errors.New(model.ParseErrorCode),
			code: connect.CodeDataLoss,
		},
		{
			err:  errors.New(model.FlagDisabledErrorCode),
			code: connect.CodeNotFound,
		},
		{
			err:  errors.New(model.GeneralErrorCode),
			code: connect.CodeUnknown,
		},
	}

	for _, test := range tests {
		err := errFormat(test.err)

		var connectErr *connect.Error
		ok := errors.As(err, &connectErr)

		if !ok {
			t.Error("formatted error is not of type connect.Error")
		}

		if connectErr.Code() != test.code {
			t.Errorf("expected code %s, but got code %s for model error %s", test.code, connectErr.Code(),
				test.err.Error())
		}
	}
}

func Test_mergeContexts(t *testing.T) {
	type args struct {
		headers                    http.Header
		headerToContextKeyMappings map[string]string
		clientContext              map[string]any
		configContext              map[string]any
	}

	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{
			name: "merge contexts with no headers, with no header-context mappings",
			args: args{
				clientContext:              map[string]any{"k1": "v1", "k2": "v2"},
				configContext:              map[string]any{"k2": "v22", "k3": "v3"},
				headers:                    http.Header{},
				headerToContextKeyMappings: map[string]string{},
			},
			// static context should "win"
			want: map[string]any{"k1": "v1", "k2": "v22", "k3": "v3"},
		},
		{
			name: "merge contexts with headers, with no header-context mappings",
			args: args{
				clientContext:              map[string]any{"k1": "v1", "k2": "v2"},
				configContext:              map[string]any{"k2": "v22", "k3": "v3"},
				headers:                    http.Header{"X-key": []string{"value"}, "X-token": []string{"token"}},
				headerToContextKeyMappings: map[string]string{},
			},
			// static context should "win"
			want: map[string]any{"k1": "v1", "k2": "v22", "k3": "v3"},
		},
		{
			name: "merge contexts with no headers, with header-context mappings",
			args: args{
				clientContext:              map[string]any{"k1": "v1", "k2": "v2"},
				configContext:              map[string]any{"k2": "v22", "k3": "v3"},
				headers:                    http.Header{},
				headerToContextKeyMappings: map[string]string{"X-key": "k2"},
			},
			// static context should "win"
			want: map[string]any{"k1": "v1", "k2": "v22", "k3": "v3"},
		},
		{
			name: "merge contexts with headers, with header-context mappings",
			args: args{
				clientContext:              map[string]any{"k1": "v1", "k2": "v2"},
				configContext:              map[string]any{"k2": "v22", "k3": "v3"},
				headers:                    http.Header{"X-key": []string{"value"}, "X-token": []string{"token"}},
				headerToContextKeyMappings: map[string]string{"X-key": "k2"},
			},
			// header context should "win"
			want: map[string]any{"k1": "v1", "k2": "value", "k3": "v3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeContexts(tt.args.clientContext, tt.args.configContext, tt.args.headers, tt.args.headerToContextKeyMappings)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot:  %+v\nwant: %+v", got, tt.want)
			}
		})
	}
}
