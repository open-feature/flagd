package ofrep

import (
	"encoding/json"
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

func TestBulkEvaluationResponse(t *testing.T) {
	tests := []struct {
		name             string
		input            []evaluator.AnyValue
		marshalledOutput string
	}{
		{
			name:             "empty input",
			input:            nil,
			marshalledOutput: "{\"flags\":[],\"metadata\":{}}",
		},
		{
			name: "valid values",
			input: []evaluator.AnyValue{
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
					Value:    false,
					Variant:  "false",
					Reason:   model.ErrorReason,
					FlagKey:  "errorFlag",
					Error:    errors.New(model.FlagNotFoundErrorCode),
					Metadata: map[string]interface{}{},
				},
			},
			marshalledOutput: "{\"flags\":[{\"value\":false,\"key\":\"key\",\"reason\":\"STATIC\",\"variant\":\"false\",\"metadata\":{\"key\":\"value\"}},{\"key\":\"errorFlag\",\"errorCode\":\"FLAG_NOT_FOUND\",\"errorDetails\":\"flag `errorFlag` does not exist\",\"metadata\":{}}],\"metadata\":{}}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := BulkEvaluationResponseFrom(test.input, model.Metadata{})

			marshal, err := json.Marshal(response)
			if err != nil {
				t.Errorf("error marshalling the response: %v", err)
			}

			if test.marshalledOutput != string(marshal) {
				t.Errorf("expected %s, got %s", test.marshalledOutput, string(marshal))
			}
		})
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
			expectedStatus: 404,
			expectedCode:   model.FlagNotFoundErrorCode,
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
