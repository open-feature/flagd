//lint:file-ignore SA4003 old proto is deprecated but we want to serve it for a while

package service

import (
	"context"
	"errors"
	"testing"

	schemaV1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/schema/v1"
	"connectrpc.com/connect"
	"github.com/open-feature/flagd/core/pkg/evaluator"
	mock "github.com/open-feature/flagd/core/pkg/evaluator/mock"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/structpb"
)

type evalCommons struct {
	variant  string
	reason   string
	metadata map[string]interface{}
}

var metadata = map[string]interface{}{
	"scope": "some-scope",
}

var responseStruct *structpb.Struct

func init() {
	pbStruct, err := structpb.NewStruct(metadata)
	if err != nil {
		panic("failure to generate protobuf structure from metadata")
	}

	responseStruct = pbStruct
}

var happyCommon = evalCommons{
	variant:  "on",
	reason:   model.DefaultReason,
	metadata: metadata,
}

var sadCommon = evalCommons{
	variant:  ":(",
	reason:   model.ErrorReason,
	metadata: metadata,
}

func TestConnectService_ResolveAll(t *testing.T) {
	tests := map[string]struct {
		req      *schemaV1.ResolveAllRequest
		evalRes  []evaluator.AnyValue
		metadata model.Metadata
		wantErr  error
		wantRes  *schemaV1.ResolveAllResponse
	}{
		"happy-path": {
			req: &schemaV1.ResolveAllRequest{},
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
				tt.evalRes, tt.metadata, nil,
			).AnyTimes()
			metrics, exp := getMetricReader()
			s := NewOldFlagEvaluationService(
				logger.NewLogger(nil, false),
				eval,
				&eventingConfiguration{},
				metrics,
				nil,
				"",
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
	result bool
	evalCommons
}

func TestFlag_Evaluation_ResolveBoolean(t *testing.T) {
	ctrl := gomock.NewController(t)

	tests := map[string]resolveBooleanArgs{
		"happy path": {
			mCount: 1,
			evalFields: resolveBooleanEvalFields{
				result:      true,
				evalCommons: happyCommon,
			},
			functionArgs: resolveBooleanFunctionArgs{
				context.Background(),
				&schemaV1.ResolveBooleanRequest{
					FlagKey: "bool",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveBooleanResponse{
				Value:    true,
				Reason:   model.DefaultReason,
				Variant:  "on",
				Metadata: responseStruct,
			},
			wantErr: nil,
		},
		"eval returns error": {
			mCount: 1,
			evalFields: resolveBooleanEvalFields{
				result:      true,
				evalCommons: sadCommon,
			},
			functionArgs: resolveBooleanFunctionArgs{
				context.Background(),
				&schemaV1.ResolveBooleanRequest{
					FlagKey: "bool",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveBooleanResponse{
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
			s := NewOldFlagEvaluationService(
				logger.NewLogger(nil, false),
				eval,
				&eventingConfiguration{},
				metrics,
				nil,
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

func BenchmarkFlag_Evaluation_ResolveBoolean(b *testing.B) {
	ctrl := gomock.NewController(b)
	tests := map[string]resolveBooleanArgs{
		"happy path": {
			evalFields: resolveBooleanEvalFields{
				result:      true,
				evalCommons: happyCommon,
			},
			functionArgs: resolveBooleanFunctionArgs{
				context.Background(),
				&schemaV1.ResolveBooleanRequest{
					FlagKey: "bool",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveBooleanResponse{
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
		s := NewOldFlagEvaluationService(
			logger.NewLogger(nil, false),
			eval,
			&eventingConfiguration{},
			metrics,
			nil,
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
	result string
	evalCommons
}

func TestFlag_Evaluation_ResolveString(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := map[string]resolveStringArgs{
		"happy path": {
			mCount: 1,
			evalFields: resolveStringEvalFields{
				result:      "true",
				evalCommons: happyCommon,
			},
			functionArgs: resolveStringFunctionArgs{
				context.Background(),
				&schemaV1.ResolveStringRequest{
					FlagKey: "string",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveStringResponse{
				Value:    "true",
				Reason:   model.DefaultReason,
				Variant:  "on",
				Metadata: responseStruct,
			},
			wantErr: nil,
		},
		"eval returns error": {
			mCount: 1,
			evalFields: resolveStringEvalFields{
				result:      "true",
				evalCommons: sadCommon,
			},
			functionArgs: resolveStringFunctionArgs{
				context.Background(),
				&schemaV1.ResolveStringRequest{
					FlagKey: "string",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveStringResponse{
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
			s := NewOldFlagEvaluationService(
				logger.NewLogger(nil, false),
				eval,
				&eventingConfiguration{},
				metrics,
				nil,
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

func BenchmarkFlag_Evaluation_ResolveString(b *testing.B) {
	ctrl := gomock.NewController(b)
	tests := map[string]resolveStringArgs{
		"happy path": {
			evalFields: resolveStringEvalFields{
				result:      "true",
				evalCommons: happyCommon,
			},
			functionArgs: resolveStringFunctionArgs{
				context.Background(),
				&schemaV1.ResolveStringRequest{
					FlagKey: "string",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveStringResponse{
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
		s := NewOldFlagEvaluationService(
			logger.NewLogger(nil, false),
			eval,
			&eventingConfiguration{},
			metrics,
			nil,
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
	result float64
	evalCommons
}

func TestFlag_Evaluation_ResolveFloat(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := map[string]resolveFloatArgs{
		"happy path": {
			mCount: 1,
			evalFields: resolveFloatEvalFields{
				result:      12,
				evalCommons: happyCommon,
			},
			functionArgs: resolveFloatFunctionArgs{
				context.Background(),
				&schemaV1.ResolveFloatRequest{
					FlagKey: "float",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveFloatResponse{
				Value:    12,
				Reason:   model.DefaultReason,
				Variant:  "on",
				Metadata: responseStruct,
			},
			wantErr: nil,
		},
		"eval returns error": {
			mCount: 1,
			evalFields: resolveFloatEvalFields{
				result:      12,
				evalCommons: sadCommon,
			},
			functionArgs: resolveFloatFunctionArgs{
				context.Background(),
				&schemaV1.ResolveFloatRequest{
					FlagKey: "float",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveFloatResponse{
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
			s := NewOldFlagEvaluationService(
				logger.NewLogger(nil, false),
				eval,
				&eventingConfiguration{},
				metrics,
				nil,
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

func BenchmarkFlag_Evaluation_ResolveFloat(b *testing.B) {
	ctrl := gomock.NewController(b)
	tests := map[string]resolveFloatArgs{
		"happy path": {
			evalFields: resolveFloatEvalFields{
				result:      12,
				evalCommons: happyCommon,
			},
			functionArgs: resolveFloatFunctionArgs{
				context.Background(),
				&schemaV1.ResolveFloatRequest{
					FlagKey: "float",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveFloatResponse{
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
		s := NewOldFlagEvaluationService(
			logger.NewLogger(nil, false),
			eval,
			&eventingConfiguration{},
			metrics,
			nil,
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
	result int64
	evalCommons
}

func TestFlag_Evaluation_ResolveInt(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := map[string]resolveIntArgs{
		"happy path": {
			mCount: 1,
			evalFields: resolveIntEvalFields{
				result:      12,
				evalCommons: happyCommon,
			},
			functionArgs: resolveIntFunctionArgs{
				context.Background(),
				&schemaV1.ResolveIntRequest{
					FlagKey: "int",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveIntResponse{
				Value:    12,
				Reason:   model.DefaultReason,
				Variant:  "on",
				Metadata: responseStruct,
			},
			wantErr: nil,
		},
		"eval returns error": {
			mCount: 1,
			evalFields: resolveIntEvalFields{
				result:      12,
				evalCommons: sadCommon,
			},
			functionArgs: resolveIntFunctionArgs{
				context.Background(),
				&schemaV1.ResolveIntRequest{
					FlagKey: "int",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveIntResponse{
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
			s := NewOldFlagEvaluationService(
				logger.NewLogger(nil, false),
				eval,
				&eventingConfiguration{},
				metrics,
				nil,
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

func BenchmarkFlag_Evaluation_ResolveInt(b *testing.B) {
	ctrl := gomock.NewController(b)
	tests := map[string]resolveIntArgs{
		"happy path": {
			evalFields: resolveIntEvalFields{
				result:      12,
				evalCommons: happyCommon,
			},
			functionArgs: resolveIntFunctionArgs{
				context.Background(),
				&schemaV1.ResolveIntRequest{
					FlagKey: "int",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveIntResponse{
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
		s := NewOldFlagEvaluationService(
			logger.NewLogger(nil, false),
			eval,
			&eventingConfiguration{},
			metrics,
			nil,
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
	result map[string]interface{}
	evalCommons
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
				evalCommons: happyCommon,
			},
			functionArgs: resolveObjectFunctionArgs{
				context.Background(),
				&schemaV1.ResolveObjectRequest{
					FlagKey: "object",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveObjectResponse{
				Value:    nil,
				Reason:   model.DefaultReason,
				Variant:  "on",
				Metadata: responseStruct,
			},
			wantErr: nil,
		},
		"eval returns error": {
			mCount: 1,
			evalFields: resolveObjectEvalFields{
				result: map[string]interface{}{
					"food": "bars",
				},
				evalCommons: sadCommon,
			},
			functionArgs: resolveObjectFunctionArgs{
				context.Background(),
				&schemaV1.ResolveObjectRequest{
					FlagKey: "object",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveObjectResponse{
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
			s := NewOldFlagEvaluationService(
				logger.NewLogger(nil, false),
				eval,
				&eventingConfiguration{},
				metrics,
				nil,
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

func BenchmarkFlag_Evaluation_ResolveObject(b *testing.B) {
	ctrl := gomock.NewController(b)
	tests := map[string]resolveObjectArgs{
		"happy path": {
			evalFields: resolveObjectEvalFields{
				result: map[string]interface{}{
					"food": "bars",
				},
				evalCommons: happyCommon,
			},
			functionArgs: resolveObjectFunctionArgs{
				context.Background(),
				&schemaV1.ResolveObjectRequest{
					FlagKey: "object",
					Context: &structpb.Struct{},
				},
			},
			want: &schemaV1.ResolveObjectResponse{
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
		s := NewOldFlagEvaluationService(
			logger.NewLogger(nil, false),
			eval,
			&eventingConfiguration{},
			metrics,
			nil,
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

func getMetricReader() (*telemetry.MetricsRecorder, metric.Reader) {
	exp := metric.NewManualReader()
	rs := resource.NewWithAttributes("testSchema")
	return telemetry.NewOTelRecorder(exp, rs, "testSvc"), exp
}

// TestFlag_Evaluation_ErrorCodes test validate error mapping from known errors to connect.Code and avoid accidental
// changes. This is essential as SDK implementations rely on connect. Code to differentiate GRPC errors vs Flag errors.
// For any change in error codes, we must change respective SDK.
func TestFlag_Evaluation_ErrorCodes(t *testing.T) {
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

func Test_Readable_ErrorMessage(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{
			name: "Testing flag not found error",
			code: model.FlagNotFoundErrorCode,
			want: model.ReadableErrorMessage[model.FlagNotFoundErrorCode],
		},
		{
			name: "Testing parse error",
			code: model.ParseErrorCode,
			want: model.ReadableErrorMessage[model.ParseErrorCode],
		},
		{
			name: "Testing type mismatch error",
			code: model.TypeMismatchErrorCode,
			want: model.ReadableErrorMessage[model.TypeMismatchErrorCode],
		},
		{
			name: "Testing general error",
			code: model.GeneralErrorCode,
			want: model.ReadableErrorMessage[model.GeneralErrorCode],
		},
		{
			name: "Testing flag disabled error",
			code: model.FlagDisabledErrorCode,
			want: model.ReadableErrorMessage[model.FlagDisabledErrorCode],
		},
		{
			name: "Testing invalid context error",
			code: model.InvalidContextCode,
			want: model.ReadableErrorMessage[model.InvalidContextCode],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := model.GetErrorMessage(tt.code); got != tt.want {
				t.Errorf("GetErrorMessage() Wanted: %v , but got: %v as a ReadableErrorMessage", tt.want, got)
			}
		})
	}
}
