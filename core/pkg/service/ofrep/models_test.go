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

func TestSuccessResultCodeDefault(t *testing.T) {
	value := evaluator.AnyValue{
		Value:    false, // zero value for boolean
		Variant:  "",    // empty variant
		Reason:   model.FallbackReason,
		FlagKey:  "noDefaultFlag",
		Metadata: map[string]interface{}{},
	}

	// when
	evaluationSuccess := SuccessResponseFrom(value)

	// then verify the reason is converted to DEFAULT
	if evaluationSuccess.Reason != model.DefaultReason {
		t.Errorf("expected reason %v, got %v", model.DefaultReason, evaluationSuccess.Reason)
	}

	if evaluationSuccess.Value != nil {
		t.Errorf("expected nil value for code default, got %v", evaluationSuccess.Value)
	}

	if evaluationSuccess.Variant != "" {
		t.Errorf("expected empty variant for code default, got %v", evaluationSuccess.Variant)
	}

	if evaluationSuccess.Key != value.FlagKey {
		t.Errorf("expected key %v, got %v", value.FlagKey, evaluationSuccess.Key)
	}

	// Verify JSON marshaling omits the value field
	marshaled, err := json.Marshal(evaluationSuccess)
	if err != nil {
		t.Errorf("error marshalling: %v", err)
	}

	// Check that "value" field is not in the JSON
	if string(marshaled) == "" {
		t.Errorf("marshalled output is empty")
	}

	// Parse back to verify structure
	var result map[string]interface{}
	err = json.Unmarshal(marshaled, &result)
	if err != nil {
		t.Errorf("error unmarshalling: %v", err)
	}

	if _, hasValue := result["value"]; hasValue {
		t.Errorf("value field should be omitted for code defaults, but found in: %s", string(marshaled))
	}

	if _, hasVariant := result["variant"]; hasVariant {
		t.Errorf("variant field should be omitted for code defaults, but found in: %s", string(marshaled))
	}

	if reason, ok := result["reason"].(string); !ok || reason != model.DefaultReason {
		t.Errorf("reason should be DEFAULT, got: %v", result["reason"])
	}
}

func TestBulkEvaluationResponse(t *testing.T) {
	tests := []struct {
		name   string
		input  []evaluator.AnyValue
		verify func(*testing.T, []byte)
	}{
		{
			name:  "empty input",
			input: nil,
			verify: func(t *testing.T, data []byte) {
				var result BulkEvaluationResponse
				err := json.Unmarshal(data, &result)
				if err != nil {
					t.Errorf("error unmarshalling: %v", err)
					return
				}
				if len(result.Flags) != 0 {
					t.Errorf("expected 0 flags, got %d", len(result.Flags))
				}
			},
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
			verify: func(t *testing.T, data []byte) {
				var result BulkEvaluationResponse
				err := json.Unmarshal(data, &result)
				if err != nil {
					t.Errorf("error unmarshalling: %v", err)
					return
				}
				if len(result.Flags) != 2 {
					t.Errorf("expected 2 flags, got %d", len(result.Flags))
					return
				}

				// Verify first flag (success case)
				firstFlag, ok := result.Flags[0].(map[string]interface{})
				if !ok {
					t.Errorf("first flag: expected map[string]interface{}, got %T", result.Flags[0])
					return
				}
				if firstFlag["key"] != "key" {
					t.Errorf("first flag: expected key 'key', got '%v'", firstFlag["key"])
				}
				if firstFlag["value"] != false {
					t.Errorf("first flag: expected value false, got %v", firstFlag["value"])
				}
				if firstFlag["variant"] != "false" {
					t.Errorf("first flag: expected variant 'false', got '%v'", firstFlag["variant"])
				}
				if firstFlag["reason"] != model.StaticReason {
					t.Errorf("first flag: expected reason %v, got %v", model.StaticReason, firstFlag["reason"])
				}
				metadata, ok := firstFlag["metadata"].(map[string]interface{})
				if !ok || metadata["key"] != "value" {
					t.Errorf("first flag: expected metadata['key']='value', got %v", firstFlag["metadata"])
				}

				// Verify second flag (error case)
				secondFlag, ok := result.Flags[1].(map[string]interface{})
				if !ok {
					t.Errorf("second flag: expected map[string]interface{}, got %T", result.Flags[1])
					return
				}
				if secondFlag["key"] != "errorFlag" {
					t.Errorf("second flag: expected key 'errorFlag', got '%v'", secondFlag["key"])
				}
				if secondFlag["errorCode"] != model.FlagNotFoundErrorCode {
					t.Errorf("second flag: expected errorCode %v, got %v", model.FlagNotFoundErrorCode, secondFlag["errorCode"])
				}
			},
		},
		{
			name: "mixed with code defaults",
			input: []evaluator.AnyValue{
				{
					Value:   "on",
					Variant: "on",
					Reason:  model.StaticReason,
					FlagKey: "featureA",
					Metadata: map[string]interface{}{
						"description": "feature A",
					},
				},
				{
					Value:    false, // code default (no defaultVariant set)
					Variant:  "",
					Reason:   model.FallbackReason,
					FlagKey:  "featureNoDefault",
					Metadata: map[string]interface{}{},
				},
				{
					Value:   42,
					Variant: "high",
					Reason:  model.TargetingMatchReason,
					FlagKey: "priority",
					Metadata: map[string]interface{}{
						"tier": "premium",
					},
				},
			},
			verify: func(t *testing.T, data []byte) {
				var result BulkEvaluationResponse
				err := json.Unmarshal(data, &result)
				if err != nil {
					t.Errorf("error unmarshalling: %v", err)
					return
				}
				if len(result.Flags) != 3 {
					t.Errorf("expected 3 flags, got %d", len(result.Flags))
					return
				}

				// Verify first flag (normal static evaluation)
				firstFlag, ok := result.Flags[0].(map[string]interface{})
				if !ok {
					t.Errorf("first flag: expected map[string]interface{}, got %T", result.Flags[0])
					return
				}
				if firstFlag["key"] != "featureA" {
					t.Errorf("first flag: expected key 'featureA', got '%v'", firstFlag["key"])
				}
				if firstFlag["value"] != "on" {
					t.Errorf("first flag: expected value 'on', got %v", firstFlag["value"])
				}
				if firstFlag["variant"] != "on" {
					t.Errorf("first flag: expected variant 'on', got '%v'", firstFlag["variant"])
				}
				if firstFlag["reason"] != model.StaticReason {
					t.Errorf("first flag: expected reason %v, got %v", model.StaticReason, firstFlag["reason"])
				}

				// Verify second flag (code default, should omit value and variant fields)
				secondFlag, ok := result.Flags[1].(map[string]interface{})
				if !ok {
					t.Errorf("second flag: expected map[string]interface{}, got %T", result.Flags[1])
					return
				}
				if secondFlag["key"] != "featureNoDefault" {
					t.Errorf("second flag: expected key 'featureNoDefault', got '%v'", secondFlag["key"])
				}
				if secondFlag["reason"] != model.DefaultReason {
					t.Errorf("second flag: expected reason %v, got %v", model.DefaultReason, secondFlag["reason"])
				}
				if _, hasValue := secondFlag["value"]; hasValue {
					t.Errorf("second flag: value field should be omitted for code defaults, but found: %v", secondFlag["value"])
				}
				if _, hasVariant := secondFlag["variant"]; hasVariant {
					t.Errorf("second flag: variant field should be omitted for code defaults, but found: %v", secondFlag["variant"])
				}

				// Verify third flag (targeting match evaluation)
				thirdFlag, ok := result.Flags[2].(map[string]interface{})
				if !ok {
					t.Errorf("third flag: expected map[string]interface{}, got %T", result.Flags[2])
					return
				}
				if thirdFlag["key"] != "priority" {
					t.Errorf("third flag: expected key 'priority', got '%v'", thirdFlag["key"])
				}
				if thirdFlag["value"] != float64(42) {
					t.Errorf("third flag: expected value 42, got %v", thirdFlag["value"])
				}
				if thirdFlag["variant"] != "high" {
					t.Errorf("third flag: expected variant 'high', got '%v'", thirdFlag["variant"])
				}
				if thirdFlag["reason"] != model.TargetingMatchReason {
					t.Errorf("third flag: expected reason %v, got %v", model.TargetingMatchReason, thirdFlag["reason"])
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := BulkEvaluationResponseFrom(test.input, model.Metadata{})
			marshal, err := json.Marshal(response)
			if err != nil {
				t.Errorf("error marshalling the response: %v", err)
				return
			}

			if test.verify != nil {
				test.verify(t, marshal)
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

func TestZeroValuesAreMarshaled(t *testing.T) {
	// This test verifies that zero-values (false, 0, empty string) are properly
	// communicated in OFREP responses, even with omitempty tags on fields
	tests := []struct {
		name           string
		value          interface{}
		variant        string
		expectedInJSON bool
	}{
		{
			name:           "boolean false is included",
			value:          false,
			variant:        "false",
			expectedInJSON: true,
		},
		{
			name:           "numeric zero is included",
			value:          float64(0),
			variant:        "zero",
			expectedInJSON: true,
		},
		{
			name:           "empty string is included",
			value:          "",
			variant:        "empty",
			expectedInJSON: true,
		},
		{
			name:           "nil value is omitted",
			value:          nil,
			variant:        "",
			expectedInJSON: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			anyValue := evaluator.AnyValue{
				Value:    test.value,
				Variant:  test.variant,
				Reason:   model.StaticReason,
				FlagKey:  "testFlag",
				Metadata: map[string]interface{}{},
			}

			if test.value == nil {
				anyValue.Reason = model.FallbackReason
			}

			success := SuccessResponseFrom(anyValue)
			marshaled, err := json.Marshal(success)
			if err != nil {
				t.Errorf("error marshalling: %v", err)
				return
			}

			var result map[string]interface{}
			err = json.Unmarshal(marshaled, &result)
			if err != nil {
				t.Errorf("error unmarshalling: %v", err)
				return
			}

			_, hasValue := result["value"]
			if test.expectedInJSON && !hasValue {
				t.Errorf("expected value field to be in JSON for %s, but it was omitted: %s", test.name, string(marshaled))
			}
			if !test.expectedInJSON && hasValue {
				t.Errorf("expected value field to be omitted from JSON for %s, but it was present: %v", test.name, result["value"])
			}

			// For non-nil values, verify the actual value matches
			if test.expectedInJSON && test.value != nil {
				if result["value"] != test.value {
					t.Errorf("expected value %v, but got %v", test.value, result["value"])
				}
			}
		})
	}
}
