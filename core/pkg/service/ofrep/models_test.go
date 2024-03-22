package ofrep

import (
	"errors"
	"reflect"
	"testing"

	"github.com/open-feature/flagd/core/pkg/evaluator"
	"github.com/open-feature/flagd/core/pkg/model"
)

func TestSuccessResult(t *testing.T) {
	// given
	value := evaluator.AnyValue{
		Value:   false,
		Variant: "false",
		Reason:  model.StaticReason,
		FlagKey: "key",
		Metadata: map[string]interface{}{
			"key": "value",
		},
	}

	// when
	evaluationSuccess := SuccessResponseFrom(value)

	if evaluationSuccess.Key != value.FlagKey {
		t.Errorf("expected %v, got %v", value.FlagKey, evaluationSuccess.Key)
	}

	if evaluationSuccess.Value != value.Value {
		t.Errorf("expected %v, got %v", value.Value, evaluationSuccess.Value)
	}

	if evaluationSuccess.Variant != value.Variant {
		t.Errorf("expected %v, got %v", value.Variant, evaluationSuccess.Variant)
	}

	if evaluationSuccess.Reason != value.Reason {
		t.Errorf("expected %v, got %v", value.Reason, evaluationSuccess.Reason)
	}

	if !reflect.DeepEqual(evaluationSuccess.Metadata, value.Metadata) {
		t.Errorf("metadata mismatch")
	}
}

func TestBulkEvaluationResult(t *testing.T) {
	// given
	values := []evaluator.AnyValue{
		{
			Value:   false,
			Variant: "false",
			Reason:  model.StaticReason,
			FlagKey: "key",
			Metadata: map[string]interface{}{
				"key": "value",
			},
		},
		{
			Value:   false,
			Variant: "false",
			Reason:  model.ErrorReason,
			FlagKey: "errorFlag",
			Error:   errors.New(model.FlagNotFoundErrorCode),
		},
	}

	bulkResponse := BulkEvaluationResponseFrom(values)

	if len(bulkResponse.Flags) != 2 {
		t.Errorf("expected bulk response to contain 2 results, but got %d", len(bulkResponse.Flags))
	}
}

func TestErrorStatus(t *testing.T) {
	tests := []struct {
		name           string
		modelError     string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "parsing error",
			modelError:     model.ParseErrorCode,
			expectedStatus: 400,
			expectedCode:   model.ParseErrorCode,
		},
		{
			name:           "flag disabled",
			modelError:     model.FlagDisabledErrorCode,
			expectedStatus: 400,
			expectedCode:   model.GeneralErrorCode,
		},
		{
			name:           "general error",
			modelError:     model.GeneralErrorCode,
			expectedStatus: 400,
			expectedCode:   model.GeneralErrorCode,
		},
		{
			name:           "flag not found",
			modelError:     model.FlagNotFoundErrorCode,
			expectedStatus: 404,
			expectedCode:   model.FlagNotFoundErrorCode,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			status, evaluationError := EvaluationErrorResponseFrom(evaluator.AnyValue{
				Value:   "value",
				Variant: "variant",
				Reason:  model.ErrorReason,
				FlagKey: "key",
				Error:   errors.New(test.modelError),
			})

			if status != test.expectedStatus {
				t.Errorf("expected status %d, but got %d", test.expectedStatus, status)
			}

			if evaluationError.ErrorCode != test.expectedCode {
				t.Errorf("expected error code %s, but got %s", test.expectedCode, evaluationError.ErrorCode)
			}
		})
	}
}
