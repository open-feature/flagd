package evaluator

import (
	"context"
	"fmt"
	"testing"

	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestJSONEvaluator_lowerEvaluation(t *testing.T) {
	const source = "testSource"
	var sources = []string{source}
	ctx := context.Background()

	tests := map[string]stringFlagEvalTestCase{
		"lower composed with == - match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"==": [{"lower": [{"var": "email"}]}, "user@example.com"]
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"email": "User@Example.com",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"lower composed with starts_with - match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"starts_with": [{"lower": [{"var": "email"}]}, "user@faas"]
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"email": "USER@FAAS.com",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"lower of bare string argument - match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"==": [{"lower": {"var": "email"}}, "user@example.com"]
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"email": "USER@EXAMPLE.COM",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"lower is ASCII-only, non-ASCII unchanged - no match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"==": [{"lower": [{"var": "name"}]}, "stra\u00dfe"]
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				// "STRAßE" lowered with ASCII-only rules stays "straße" (ß unchanged),
				// not "strasse"; equality target is "straße" so this matches.
				"name": "STRA\u00dfE",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"non-string input falls back to default": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"==": [{"lower": [{"var": "num"}]}, "1"]
											  },
											  "blue", "red"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"num": 1,
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
	}

	runStringFlagEvalTests(t, ctx, source, sources, tests)
}

func TestJSONEvaluator_upperEvaluation(t *testing.T) {
	const source = "testSource"
	var sources = []string{source}
	ctx := context.Background()

	tests := map[string]stringFlagEvalTestCase{
		"upper composed with == - match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"==": [{"upper": [{"var": "country"}]}, "US"]
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"country": "us",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"upper composed with == - no match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"==": [{"upper": [{"var": "country"}]}, "US"]
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"country": "ca",
			},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"upper is ASCII-only, non-ASCII unchanged - no match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"==": [{"upper": [{"var": "name"}]}, "STRASSE"]
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				// "straße" uppered with ASCII-only rules stays "STRAßE", never "STRASSE".
				"name": "stra\u00dfe",
			},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
	}

	runStringFlagEvalTests(t, ctx, source, sources, tests)
}

func Test_parseStringCasingEvaluationData(t *testing.T) {
	type args struct {
		values interface{}
	}
	tests := []struct {
		name         string
		args         args
		wantProperty string
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name:         "bare string value",
			args:         args{values: "a"},
			wantProperty: "a",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err)
				return true
			},
		},
		{
			name:         "single-element array",
			args:         args{values: []interface{}{"a"}},
			wantProperty: "a",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err)
				return true
			},
		},
		{
			name:         "array with more than one element",
			args:         args{values: []interface{}{"a", "b"}},
			wantProperty: "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NotNil(t, err)
				return true
			},
		},
		{
			name:         "value is not a string",
			args:         args{values: []interface{}{1}},
			wantProperty: "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NotNil(t, err)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseStringCasingEvaluationData(tt.args.values)
			if !tt.wantErr(t, err, fmt.Sprintf("parseStringCasingEvaluationData(%v)", tt.args.values)) {
				return
			}
			assert.Equalf(t, tt.wantProperty, got, "parseStringCasingEvaluationData(%v)", tt.args.values)
		})
	}
}

func TestStringCasingEvaluation_ErrorFallbackWhenUsedDirectly(t *testing.T) {
	const source = "testSource"
	ctx := context.Background()

	tests := map[string]errorFallbackTestCase{
		"lower invalid input falls back": {
			targeting: `{"lower": [{"var": "num"}]}`,
			context:   map[string]any{"num": 123.0},
		},
		"upper invalid input falls back": {
			targeting: `{"upper": [{"var": "num"}]}`,
			context:   map[string]any{"num": 123.0},
		},
	}

	runErrorFallbackTests(t, ctx, source, "string-casing-error-fallback", tests)
}
