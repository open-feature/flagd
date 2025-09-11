package evaluator

import (
	"context"
	"errors"
	"testing"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/stretchr/testify/require"
)

func TestSemVerOperator_Compare(t *testing.T) {
	type args struct {
		v1 string
		v2 string
	}
	tests := []struct {
		name    string
		svo     SemVerOperator
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "invalid version",
			svo:  Greater,
			args: args{
				v1: "invalid",
				v2: "v1.0.0",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "preview version vs non preview version",
			svo:  Greater,
			args: args{
				v1: "v1.0.0-preview.1.2",
				v2: "v1.0.0",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "preview version vs preview version",
			svo:  Greater,
			args: args{
				v1: "v1.0.0-preview.1.3",
				v2: "v1.0.0-preview.1.2",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "no prefixed v left greater",
			svo:  Greater,
			args: args{
				v1: "0.0.1",
				v2: "v0.0.2",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "no prefixed v right greater",
			svo:  Greater,
			args: args{
				v1: "v0.0.1",
				v2: "0.0.2",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "no prefixed v right equals",
			svo:  Equals,
			args: args{
				v1: "v0.0.1",
				v2: "0.0.1",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "no prefixed v both",
			svo:  Greater,
			args: args{
				v1: "0.0.1",
				v2: "0.0.2",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "invalid operator",
			svo:  "",
			args: args{
				v1: "v0.0.1",
				v2: "v0.0.2",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "less with large number",
			svo:  Less,
			args: args{
				v1: "v1234.0.1",
				v2: "v1235.0.2",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "less",
			svo:  Less,
			args: args{
				v1: "v0.0.1",
				v2: "v0.0.2",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "no minor version",
			svo:  Less,
			args: args{
				v1: "v1.0",
				v2: "v1.2",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "not less",
			svo:  Less,
			args: args{
				v1: "v0.0.1",
				v2: "v0.0.1",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "less or equal",
			svo:  LessOrEqual,
			args: args{
				v1: "v0.0.1",
				v2: "v0.0.2",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "less or equal 2",
			svo:  LessOrEqual,
			args: args{
				v1: "v0.0.1",
				v2: "v0.0.1",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "equal",
			svo:  Equals,
			args: args{
				v1: "v0.0.1",
				v2: "v0.0.1",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "not equal",
			svo:  Equals,
			args: args{
				v1: "v0.0.2",
				v2: "v0.0.1",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "unequal",
			svo:  NotEqual,
			args: args{
				v1: "v0.0.2",
				v2: "v0.0.1",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "not unequal",
			svo:  NotEqual,
			args: args{
				v1: "v0.0.1",
				v2: "v0.0.1",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "greater or equal 1",
			svo:  GreaterOrEqual,
			args: args{
				v1: "v0.0.2",
				v2: "v0.0.1",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "greater or equal 2",
			svo:  GreaterOrEqual,
			args: args{
				v1: "v0.0.1",
				v2: "v0.0.1",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "not greater or equal",
			svo:  GreaterOrEqual,
			args: args{
				v1: "v0.0.1",
				v2: "v0.0.2",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "greater",
			svo:  Greater,
			args: args{
				v1: "v0.0.2",
				v2: "v0.0.1",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "not greater",
			svo:  Greater,
			args: args{
				v1: "v0.0.1",
				v2: "v0.0.1",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "matching major version",
			svo:  MatchMajor,
			args: args{
				v1: "v1.3.4",
				v2: "v1.5.3",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "not matching major version",
			svo:  MatchMajor,
			args: args{
				v1: "v2.1.1",
				v2: "v1.1.1",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "matching minor version",
			svo:  MatchMinor,
			args: args{
				v1: "v1.3.4",
				v2: "v1.3.1",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "not matching minor version",
			svo:  MatchMinor,
			args: args{
				v1: "v2.2.1",
				v2: "v2.1.1",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var operatorInterface interface{} = string(tt.svo)
			actualVersion, targetVersion, operator, err := parseSemverEvaluationData([]interface{}{tt.args.v1, operatorInterface, tt.args.v2})
			if err != nil {
				require.Truef(t, tt.wantErr, "Error parsing semver evaluation data. actualVersion: %s, targetVersion: %s, operator: %s, err: %s", actualVersion, targetVersion, operator, err)
				return
			}

			got, err := operator.compare(actualVersion, targetVersion)

			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
				require.Equalf(t, tt.want, got, "compare(%v, %v) operator: %s", tt.args.v1, tt.args.v2, operator)
			}
		})
	}
}

func TestJSONEvaluator_semVerEvaluation(t *testing.T) {
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
		"versions and operator provided - match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"sem_ver": ["1.0.0", ">", "0.1.0"]
											  },
											  "red", null
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"version": "1.0.0",
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
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"sem_ver": [{"var": "version"}, ">", "1.0.0"]
											  },
											  "red", null
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"version": "1.0.1",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"versions and operator provided - no match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"sem_ver": ["1.0.0", ">", "1.0.0"]
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"version": "1.0.0",
			},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"versions and major-version operator provided - match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"sem_ver": ["1.2.3", "^", "1.5.6"]
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"version": "1.0.0",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"versions and minor-version operator provided - match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"sem_ver": ["1.2.3", "~", "1.2.6"]
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"version": "1.0.0",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"versions given as double - match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"sem_ver": [1.2, "=", "1.2"]
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"version": "1.0.0",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"versions given as int - match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"sem_ver": [1, "=", "v1.0.0"]
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"version": "1.0.0",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"versions and minor-version without patch version operator provided - match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"sem_ver": [1.2, "=", "1.2"]
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"version": "1.0.0",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"versions with prefixed v operator provided - match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"sem_ver": [{"var": "version"}, "<", "v1.2"]
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"version": "v1.0.0",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"versions and major-version operator provided - no match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"sem_ver": ["2.2.3", "^", "1.2.3"]
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"version": "1.0.0",
			},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"versions and minor-version operator provided - no match": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"sem_ver": ["1.3.3", "~", "1.2.6"]
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"version": "1.0.0",
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
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"sem_ver": [{"var": "version"}, ">", "1.0.0"]
											  },
											  "red", "green"
											]
										  }`),
			},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"version": "0.0.1",
			},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"error during parsing (not an array) - return default": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"sem_ver": "not an array"
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
		"error during parsing (wrong number of items in array) - return default": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"sem_ver": ["not", "enough"]
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
		"error during parsing (invalid property value) - return default": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"sem_ver": ["invalid", ">", "1.0.0"]
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
		"error during parsing (invalid property type) - return default": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"sem_ver": [1.0, ">", "1.0.0"]
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
		"error during parsing (invalid operator) - return default": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"sem_ver": ["1.0.0", "invalid", "1.0.0"]
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
		"error during parsing (invalid operator type) - return default": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"sem_ver": ["1.0.0", 1, "1.0.0"]
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
		"error during parsing (invalid target version) - return default": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: "red",
				Variants:       colorVariants,
				Targeting: []byte(`{
											"if": [
											  {
												"sem_ver": ["1.0.0", ">", "invalid"]
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
