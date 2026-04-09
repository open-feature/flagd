package service

import (
	"context"
	"testing"

	evalV2 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/evaluation/v2"
	"connectrpc.com/connect"
	mock "github.com/open-feature/flagd/core/pkg/evaluator/mock"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/structpb"
)

// fallbackResult captures the fields we care about from any resolve response.
type fallbackResult struct {
	value    any
	variant  *string
	reason   string
	metadata *structpb.Struct
}

// TestFlagEvaluationServiceV2_Fallback tests that the V2 service correctly
// handles the fallback case (null targeting + no defaultVariant) for all flag
// types. Each response type's SetReasonOnly method must set reason to DEFAULT
// with nil value and variant, signaling the provider to use the code default.
// note: V2 signals fallback via SetReasonOnly (V1 uses an error instead);
// other V2 behaviors are covered in flag_evaluator_v1_test.go.
func TestFlagEvaluationServiceV2_Fallback(t *testing.T) {
	tests := []struct {
		name string
		// run sets up the mock, calls the service, and returns the common response fields
		run func(eval *mock.MockIEvaluator, s *FlagEvaluationServiceV2) (fallbackResult, error)
	}{
		{
			name: "boolean",
			run: func(eval *mock.MockIEvaluator, s *FlagEvaluationServiceV2) (fallbackResult, error) {
				eval.EXPECT().ResolveBooleanValue(gomock.Any(), gomock.Any(), "flag", gomock.Any()).Return(
					false, "", model.FallbackReason, metadata, nil,
				)
				got, err := s.ResolveBoolean(context.Background(), connect.NewRequest(&evalV2.ResolveBooleanRequest{
					FlagKey: "flag", Context: &structpb.Struct{},
				}))
				return fallbackResult{got.Msg.Value, got.Msg.Variant, got.Msg.Reason, got.Msg.Metadata}, err
			},
		},
		{
			name: "string",
			run: func(eval *mock.MockIEvaluator, s *FlagEvaluationServiceV2) (fallbackResult, error) {
				eval.EXPECT().ResolveStringValue(gomock.Any(), gomock.Any(), "flag", gomock.Any()).Return(
					"", "", model.FallbackReason, metadata, nil,
				)
				got, err := s.ResolveString(context.Background(), connect.NewRequest(&evalV2.ResolveStringRequest{
					FlagKey: "flag", Context: &structpb.Struct{},
				}))
				return fallbackResult{got.Msg.Value, got.Msg.Variant, got.Msg.Reason, got.Msg.Metadata}, err
			},
		},
		{
			name: "int",
			run: func(eval *mock.MockIEvaluator, s *FlagEvaluationServiceV2) (fallbackResult, error) {
				eval.EXPECT().ResolveIntValue(gomock.Any(), gomock.Any(), "flag", gomock.Any()).Return(
					int64(0), "", model.FallbackReason, metadata, nil,
				)
				got, err := s.ResolveInt(context.Background(), connect.NewRequest(&evalV2.ResolveIntRequest{
					FlagKey: "flag", Context: &structpb.Struct{},
				}))
				return fallbackResult{got.Msg.Value, got.Msg.Variant, got.Msg.Reason, got.Msg.Metadata}, err
			},
		},
		{
			name: "float",
			run: func(eval *mock.MockIEvaluator, s *FlagEvaluationServiceV2) (fallbackResult, error) {
				eval.EXPECT().ResolveFloatValue(gomock.Any(), gomock.Any(), "flag", gomock.Any()).Return(
					float64(0), "", model.FallbackReason, metadata, nil,
				)
				got, err := s.ResolveFloat(context.Background(), connect.NewRequest(&evalV2.ResolveFloatRequest{
					FlagKey: "flag", Context: &structpb.Struct{},
				}))
				return fallbackResult{got.Msg.Value, got.Msg.Variant, got.Msg.Reason, got.Msg.Metadata}, err
			},
		},
		{
			name: "object",
			run: func(eval *mock.MockIEvaluator, s *FlagEvaluationServiceV2) (fallbackResult, error) {
				eval.EXPECT().ResolveObjectValue(gomock.Any(), gomock.Any(), "flag", gomock.Any()).Return(
					nil, "", model.FallbackReason, metadata, nil,
				)
				got, err := s.ResolveObject(context.Background(), connect.NewRequest(&evalV2.ResolveObjectRequest{
					FlagKey: "flag", Context: &structpb.Struct{},
				}))
				return fallbackResult{got.Msg.Value, got.Msg.Variant, got.Msg.Reason, got.Msg.Metadata}, err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			eval := mock.NewMockIEvaluator(ctrl)
			metrics, exp := getMetricReader()
			s := NewFlagEvaluationServiceV2(
				logger.NewLogger(nil, false), eval, &eventingConfiguration{},
				metrics, nil, nil, 0,
			)

			result, err := tt.run(eval, s)

			require.NoError(t, err)
			require.Nil(t, result.value, "value should be nil so the provider falls back to the code default")
			require.Nil(t, result.variant, "variant should be nil")
			require.Equal(t, model.DefaultReason, result.reason)
			require.Equal(t, responseStruct, result.metadata)

			var data metricdata.ResourceMetrics
			require.NoError(t, exp.Collect(context.TODO(), &data))
			require.Equal(t, 1, len(data.ScopeMetrics))
		})
	}
}
