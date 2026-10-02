package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	evalV1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/evaluation/v1"
	schemaV1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/schema/v1"
	"connectrpc.com/connect"
	"github.com/open-feature/flagd/core/pkg/evaluator"
	evaluatormock "github.com/open-feature/flagd/core/pkg/evaluator/mock"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/structpb"
)

type capturedEvaluation struct {
	err       error
	reason    string
	variant   string
	key       string
	flagSetID string
}

type capturingEvaluationMetrics struct {
	telemetry.NoopMetricsRecorder
	evaluations []capturedEvaluation
}

func (m *capturingEvaluationMetrics) RecordEvaluation(
	_ context.Context, err error, reason, variant, key, flagSetID string,
) {
	m.evaluations = append(m.evaluations, capturedEvaluation{
		err:       err,
		reason:    reason,
		variant:   variant,
		key:       key,
		flagSetID: flagSetID,
	})
}

type boolResponseCapture struct{}

func (boolResponseCapture) SetResult(bool, string, string, map[string]interface{}) error {
	return nil
}

func TestTypedSingleEvaluationRecordsFlagSetID(t *testing.T) {
	tests := []struct {
		name    string
		resolve func(telemetry.IMetricsRecorder) error
		wantID  string
		wantErr bool
	}{
		{
			name: "legacy success",
			resolve: func(metrics telemetry.IMetricsRecorder) error {
				return resolve(
					t.Context(), logger.NewLogger(nil, false),
					func(context.Context, string, string, map[string]any) (bool, string, string, map[string]interface{}, error) {
						return true, "on", model.StaticReason, model.Metadata{"flagSetId": "set-a"}, nil
					},
					http.Header{}, "flag", &structpb.Struct{}, boolResponseCapture{}, metrics, nil, nil,
				)
			},
			wantID: "set-a",
		},
		{
			name: "v2 selected lookup error",
			resolve: func(metrics telemetry.IMetricsRecorder) error {
				return resolveV2(
					t.Context(), logger.NewLogger(nil, false),
					func(context.Context, string, string, map[string]any) (bool, string, string, map[string]interface{}, error) {
						return false, "", model.ErrorReason, model.Metadata{"flagSetId": "selected-set"}, errors.New(model.FlagNotFoundErrorCode)
					},
					http.Header{}, "missing", &structpb.Struct{}, boolResponseCapture{}, metrics, nil, nil,
				)
			},
			wantID:  "selected-set",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			metrics := &capturingEvaluationMetrics{}
			err := test.resolve(metrics)
			require.Equal(t, test.wantErr, err != nil)
			require.Len(t, metrics.evaluations, 1)
			require.Equal(t, test.wantID, metrics.evaluations[0].flagSetID)
		})
	}
}

func TestTypedBulkEvaluationRecordsEachFlagSetID(t *testing.T) {
	tests := []struct {
		name string
		call func(*evaluatormock.MockIEvaluator, telemetry.IMetricsRecorder) error
	}{
		{
			name: "legacy schema",
			call: func(eval *evaluatormock.MockIEvaluator, metrics telemetry.IMetricsRecorder) error {
				service := NewOldFlagEvaluationService(
					logger.NewLogger(nil, false), eval, &eventingConfiguration{}, metrics, nil,
				)
				_, err := service.ResolveAll(t.Context(), connect.NewRequest(&schemaV1.ResolveAllRequest{}))
				return err
			},
		},
		{
			name: "evaluation v1",
			call: func(eval *evaluatormock.MockIEvaluator, metrics telemetry.IMetricsRecorder) error {
				service := NewFlagEvaluationService(
					logger.NewLogger(nil, false), eval, &eventingConfiguration{}, metrics, nil, nil, 0,
				)
				_, err := service.ResolveAll(t.Context(), connect.NewRequest(&evalV1.ResolveAllRequest{}))
				return err
			},
		},
	}

	values := []evaluator.AnyValue{
		{
			Value: true, Variant: "on", Reason: model.StaticReason, FlagKey: "flag-a",
			Metadata: model.Metadata{"flagSetId": "set-a"},
		},
		{
			Value: "value", Variant: "default", Reason: model.StaticReason, FlagKey: "flag-b",
			Metadata: model.Metadata{"flagSetId": "set-b"},
		},
		{
			Value: true, Variant: "on", Reason: model.StaticReason, FlagKey: "without-set",
			Metadata: model.Metadata{"flagSetId": 42},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			eval := evaluatormock.NewMockIEvaluator(gomock.NewController(t))
			eval.EXPECT().ResolveAllValues(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(values, model.Metadata{}, nil)
			metrics := &capturingEvaluationMetrics{}

			require.NoError(t, test.call(eval, metrics))
			require.Len(t, metrics.evaluations, 3)
			require.Equal(t, []string{"set-a", "set-b", ""}, []string{
				metrics.evaluations[0].flagSetID,
				metrics.evaluations[1].flagSetID,
				metrics.evaluations[2].flagSetID,
			})
		})
	}
}
