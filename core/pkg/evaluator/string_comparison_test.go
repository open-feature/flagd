package evaluator

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/stretchr/testify/assert"
)

func TestJSONEvaluator_startsWithEvaluation(t *testing.T) {
	const source = "testSource"
	var sources = []string{source}
	ctx := context.Background()

	tests := map[string]struct {
		flags           []model.Flag
		flagKey         string
		context         map[string]any
		expectedValue   string
		expectedVariant string
		expectedReason  string
		expectedError   error
	}{
		"two strings provided - match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants: colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"starts_with": ["user@faas.com", "user@faas"]
											  },
											  "red", null
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"email": "user@faas.com",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"resolve target property using nested operation - match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants: colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"starts_with": [{"var": "email"}, "user@faas"]
											  },
											  "red", null
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"email": "user@faas.com",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"two strings provided - no match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants: colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"starts_with": ["user@faas.com", "nope"]
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"email": "user@faas.com",
			},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"resolve target property using nested operation - no match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants: colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"starts_with": [{"var": "email"}, "nope"]
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"email": "user@faas.com",
			},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"error during parsing - return default": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants: colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"starts_with": "no-array"
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"email": "user@faas.com",
			},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
	}

	const reqID = "default"
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			log := logger.NewLogger(nil, false)
			s, err := store.NewStore(log, sources)
			if err != nil {
				t.Fatalf("NewStore failed: %v", err)
			}
			je := NewJSON(log, s)
			je.store.Update(source, tt.flags, model.Metadata{})

			value, variant, reason, _, err := resolve[string](ctx, reqID, tt.flagKey, tt.context, je.evaluateVariant)

			if value != tt.expectedValue {
				t.Errorf("expected value '%s', got '%s'", tt.expectedValue, value)
			}

			if variant != tt.expectedVariant {
				t.Errorf("expected variant '%s', got '%s'", tt.expectedVariant, variant)
			}

			if reason != tt.expectedReason {
				t.Errorf("expected reason '%s', got '%s'", tt.expectedReason, reason)
			}

			if !errors.Is(err, tt.expectedError) {
				t.Errorf("expected err '%v', got '%v'", tt.expectedError, err)
			}
		})
	}
}

func TestJSONEvaluator_endsWithEvaluation(t *testing.T) {
	const source = "testSource"
	var sources = []string{source}
	ctx := context.Background()

	tests := map[string]struct {
		flags           []model.Flag
		flagKey         string
		context         map[string]any
		expectedValue   string
		expectedVariant string
		expectedReason  string
		expectedError   error
	}{
		"two strings provided - match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants: colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"ends_with": ["user@faas.com", "faas.com"]
											  },
											  "red", null
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"email": "user@faas.com",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"resolve target property using nested operation - match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants: colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"ends_with": [{"var": "email"}, "faas.com"]
											  },
											  "red", null
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"email": "user@faas.com",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"two strings provided - no match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants: colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"ends_with": ["user@faas.com", "nope"]
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"email": "user@faas.com",
			},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"resolve target property using nested operation - no match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants: colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"ends_with": [{"var": "email"}, "nope"]
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"email": "user@faas.com",
			},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"error during parsing - return default": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants: colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"ends_with": "no-array"
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"email": "user@faas.com",
			},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
	}

	const reqID = "default"
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			log := logger.NewLogger(nil, false)
			s, err := store.NewStore(log, sources)
			if err != nil {
				t.Fatalf("NewStore failed: %v", err)
			}
			je := NewJSON(log, s)
			je.store.Update(source, tt.flags, model.Metadata{})

			value, variant, reason, _, err := resolve[string](ctx, reqID, tt.flagKey, tt.context, je.evaluateVariant)

			if value != tt.expectedValue {
				t.Errorf("expected value '%s', got '%s'", tt.expectedValue, value)
			}

			if variant != tt.expectedVariant {
				t.Errorf("expected variant '%s', got '%s'", tt.expectedVariant, variant)
			}

			if reason != tt.expectedReason {
				t.Errorf("expected reason '%s', got '%s'", tt.expectedReason, reason)
			}

			if err != tt.expectedError {
				t.Errorf("expected err '%v', got '%v'", tt.expectedError, err)
			}
		})
	}
}

func Test_parseStringComparisonEvaluationData(t *testing.T) {
	type args struct {
		values interface{}
	}
	tests := []struct {
		name            string
		args            args
		wantProperty    string
		wantTargetValue string
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name: "return two string values",
			args: args{
				values: []interface{}{"a", "b"},
			},
			wantProperty:    "a",
			wantTargetValue: "b",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err)
				return false
			},
		},
		{
			name: "provided object is not an array",
			args: args{
				values: "not-an-array",
			},
			wantProperty:    "",
			wantTargetValue: "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NotNil(t, err)
				return true
			},
		},
		{
			name: "provided object does not have two elements",
			args: args{
				values: []interface{}{"a"},
			},
			wantProperty:    "",
			wantTargetValue: "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NotNil(t, err)
				return true
			},
		},
		{
			name: "property is not a string",
			args: args{
				values: []interface{}{1, "b"},
			},
			wantProperty:    "",
			wantTargetValue: "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NotNil(t, err)
				return true
			},
		},
		{
			name: "targetValue is not a string",
			args: args{
				values: []interface{}{"a", 1},
			},
			wantProperty:    "",
			wantTargetValue: "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NotNil(t, err)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseStringComparisonEvaluationData(tt.args.values)
			if !tt.wantErr(t, err, fmt.Sprintf("parseStringComparisonEvaluationData(%v)", tt.args.values)) {
				return
			}
			assert.Equalf(t, tt.wantProperty, got, "parseStringComparisonEvaluationData(%v)", tt.args.values)
			assert.Equalf(t, tt.wantTargetValue, got1, "parseStringComparisonEvaluationData(%v)", tt.args.values)
		})
	}
}

func TestJSONEvaluator_regexMatchEvaluation(t *testing.T) {
	const source = "testSource"
	var sources = []string{source}
	ctx := context.Background()

	tests := map[string]struct {
		flags           []model.Flag
		flagKey         string
		context         map[string]any
		expectedValue   string
		expectedVariant string
		expectedReason  string
		expectedError   error
	}{
		"two strings provided - match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants: colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"regex_match": ["user@faas.com", ".*"]
											  },
											  "red", null
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"resolve target property using nested operation - match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants: colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"regex_match": [{"var": "email"}, ".*@.*"]
											  },
											  "red", null
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"email": "user@faas.com",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"two strings provided - no match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants: colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"regex_match": ["user@faas.com", ".*FAAS.*"]
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"email": "user@faas.com",
			},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"three strings provided - match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants: colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"regex_match": ["user@faas.com", ".*FAAS.*", "i"]
											  },
											  "red", null
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"resolve target property using nested operation - no match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants: colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"regex_match": [{"var": "email"}, "nope"]
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"email": "user@faas.com",
			},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"error during parsing - return default": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants: colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"regex_match": "no-array"
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"email": "user@faas.com",
			},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
	}

	const reqID = "default"
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			log := logger.NewLogger(nil, false)
			s, err := store.NewStore(log, sources)
			if err != nil {
				t.Fatalf("NewStore failed: %v", err)
			}
			je := NewJSON(log, s)
			je.store.Update(source, tt.flags, model.Metadata{})

			value, variant, reason, _, err := resolve[string](ctx, reqID, tt.flagKey, tt.context, je.evaluateVariant)

			if value != tt.expectedValue {
				t.Errorf("expected value '%s', got '%s'", tt.expectedValue, value)
			}

			if variant != tt.expectedVariant {
				t.Errorf("expected variant '%s', got '%s'", tt.expectedVariant, variant)
			}

			if reason != tt.expectedReason {
				t.Errorf("expected reason '%s', got '%s'", tt.expectedReason, reason)
			}

			if err != tt.expectedError {
				t.Errorf("expected err '%v', got '%v'", tt.expectedError, err)
			}
		})
	}
}