package eval

import (
	"testing"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
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
			got, err := tt.svo.compare(tt.args.v1, tt.args.v2)

			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
				require.Equalf(t, tt.want, got, "compare(%v, %v)", tt.args.v1, tt.args.v2)
			}
		})
	}
}

func TestJSONEvaluator_semVerEvaluation(t *testing.T) {
	tests := map[string]struct {
		flags           Flags
		flagKey         string
		context         *structpb.Struct
		expectedValue   string
		expectedVariant string
		expectedReason  string
		expectedError   error
	}{
		"versions and operator provided - match": {
			flags: Flags{
				Flags: map[string]model.Flag{
					"headerColor": {
						State:          "ENABLED",
						DefaultVariant: "red",
						Variants: map[string]any{
							"red":    "#FF0000",
							"blue":   "#0000FF",
							"green":  "#00FF00",
							"yellow": "#FFFF00",
						},
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
			},
			flagKey: "headerColor",
			context: &structpb.Struct{Fields: map[string]*structpb.Value{
				"version": {Kind: &structpb.Value_StringValue{
					StringValue: "1.0.0",
				}},
			}},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"resolve target property using nested operation - match": {
			flags: Flags{
				Flags: map[string]model.Flag{
					"headerColor": {
						State:          "ENABLED",
						DefaultVariant: "red",
						Variants: map[string]any{
							"red":    "#FF0000",
							"blue":   "#0000FF",
							"green":  "#00FF00",
							"yellow": "#FFFF00",
						},
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
			},
			flagKey: "headerColor",
			context: &structpb.Struct{Fields: map[string]*structpb.Value{
				"version": {Kind: &structpb.Value_StringValue{
					StringValue: "1.0.1",
				}},
			}},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"versions and operator provided - no match": {
			flags: Flags{
				Flags: map[string]model.Flag{
					"headerColor": {
						State:          "ENABLED",
						DefaultVariant: "red",
						Variants: map[string]any{
							"red":    "#FF0000",
							"blue":   "#0000FF",
							"green":  "#00FF00",
							"yellow": "#FFFF00",
						},
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
			},
			flagKey: "headerColor",
			context: &structpb.Struct{Fields: map[string]*structpb.Value{
				"version": {Kind: &structpb.Value_StringValue{
					StringValue: "1.0.0",
				}},
			}},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"versions and major-version operator provided - match": {
			flags: Flags{
				Flags: map[string]model.Flag{
					"headerColor": {
						State:          "ENABLED",
						DefaultVariant: "red",
						Variants: map[string]any{
							"red":    "#FF0000",
							"blue":   "#0000FF",
							"green":  "#00FF00",
							"yellow": "#FFFF00",
						},
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
			},
			flagKey: "headerColor",
			context: &structpb.Struct{Fields: map[string]*structpb.Value{
				"version": {Kind: &structpb.Value_StringValue{
					StringValue: "1.0.0",
				}},
			}},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"versions and minor-version operator provided - match": {
			flags: Flags{
				Flags: map[string]model.Flag{
					"headerColor": {
						State:          "ENABLED",
						DefaultVariant: "red",
						Variants: map[string]any{
							"red":    "#FF0000",
							"blue":   "#0000FF",
							"green":  "#00FF00",
							"yellow": "#FFFF00",
						},
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
			},
			flagKey: "headerColor",
			context: &structpb.Struct{Fields: map[string]*structpb.Value{
				"version": {Kind: &structpb.Value_StringValue{
					StringValue: "1.0.0",
				}},
			}},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"versions and major-version operator provided - no match": {
			flags: Flags{
				Flags: map[string]model.Flag{
					"headerColor": {
						State:          "ENABLED",
						DefaultVariant: "red",
						Variants: map[string]any{
							"red":    "#FF0000",
							"blue":   "#0000FF",
							"green":  "#00FF00",
							"yellow": "#FFFF00",
						},
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
			},
			flagKey: "headerColor",
			context: &structpb.Struct{Fields: map[string]*structpb.Value{
				"version": {Kind: &structpb.Value_StringValue{
					StringValue: "1.0.0",
				}},
			}},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"versions and minor-version operator provided - no match": {
			flags: Flags{
				Flags: map[string]model.Flag{
					"headerColor": {
						State:          "ENABLED",
						DefaultVariant: "red",
						Variants: map[string]any{
							"red":    "#FF0000",
							"blue":   "#0000FF",
							"green":  "#00FF00",
							"yellow": "#FFFF00",
						},
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
			},
			flagKey: "headerColor",
			context: &structpb.Struct{Fields: map[string]*structpb.Value{
				"version": {Kind: &structpb.Value_StringValue{
					StringValue: "1.0.0",
				}},
			}},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"resolve target property using nested operation - no match": {
			flags: Flags{
				Flags: map[string]model.Flag{
					"headerColor": {
						State:          "ENABLED",
						DefaultVariant: "red",
						Variants: map[string]any{
							"red":    "#FF0000",
							"blue":   "#0000FF",
							"green":  "#00FF00",
							"yellow": "#FFFF00",
						},
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
			},
			flagKey: "headerColor",
			context: &structpb.Struct{Fields: map[string]*structpb.Value{
				"version": {Kind: &structpb.Value_StringValue{
					StringValue: "0.0.1",
				}},
			}},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"error during parsing (not an array) - return default": {
			flags: Flags{
				Flags: map[string]model.Flag{
					"headerColor": {
						State:          "ENABLED",
						DefaultVariant: "red",
						Variants: map[string]any{
							"red":    "#FF0000",
							"blue":   "#0000FF",
							"green":  "#00FF00",
							"yellow": "#FFFF00",
						},
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
			},
			flagKey: "headerColor",
			context: &structpb.Struct{Fields: map[string]*structpb.Value{
				"email": {Kind: &structpb.Value_StringValue{
					StringValue: "user@faas.com",
				}},
			}},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"error during parsing (wrong number of items in array) - return default": {
			flags: Flags{
				Flags: map[string]model.Flag{
					"headerColor": {
						State:          "ENABLED",
						DefaultVariant: "red",
						Variants: map[string]any{
							"red":    "#FF0000",
							"blue":   "#0000FF",
							"green":  "#00FF00",
							"yellow": "#FFFF00",
						},
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
			},
			flagKey: "headerColor",
			context: &structpb.Struct{Fields: map[string]*structpb.Value{
				"email": {Kind: &structpb.Value_StringValue{
					StringValue: "user@faas.com",
				}},
			}},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"error during parsing (invalid property value) - return default": {
			flags: Flags{
				Flags: map[string]model.Flag{
					"headerColor": {
						State:          "ENABLED",
						DefaultVariant: "red",
						Variants: map[string]any{
							"red":    "#FF0000",
							"blue":   "#0000FF",
							"green":  "#00FF00",
							"yellow": "#FFFF00",
						},
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
			},
			flagKey: "headerColor",
			context: &structpb.Struct{Fields: map[string]*structpb.Value{
				"email": {Kind: &structpb.Value_StringValue{
					StringValue: "user@faas.com",
				}},
			}},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"error during parsing (invalid property type) - return default": {
			flags: Flags{
				Flags: map[string]model.Flag{
					"headerColor": {
						State:          "ENABLED",
						DefaultVariant: "red",
						Variants: map[string]any{
							"red":    "#FF0000",
							"blue":   "#0000FF",
							"green":  "#00FF00",
							"yellow": "#FFFF00",
						},
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
			},
			flagKey: "headerColor",
			context: &structpb.Struct{Fields: map[string]*structpb.Value{
				"email": {Kind: &structpb.Value_StringValue{
					StringValue: "user@faas.com",
				}},
			}},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"error during parsing (invalid operator) - return default": {
			flags: Flags{
				Flags: map[string]model.Flag{
					"headerColor": {
						State:          "ENABLED",
						DefaultVariant: "red",
						Variants: map[string]any{
							"red":    "#FF0000",
							"blue":   "#0000FF",
							"green":  "#00FF00",
							"yellow": "#FFFF00",
						},
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
			},
			flagKey: "headerColor",
			context: &structpb.Struct{Fields: map[string]*structpb.Value{
				"email": {Kind: &structpb.Value_StringValue{
					StringValue: "user@faas.com",
				}},
			}},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"error during parsing (invalid operator type) - return default": {
			flags: Flags{
				Flags: map[string]model.Flag{
					"headerColor": {
						State:          "ENABLED",
						DefaultVariant: "red",
						Variants: map[string]any{
							"red":    "#FF0000",
							"blue":   "#0000FF",
							"green":  "#00FF00",
							"yellow": "#FFFF00",
						},
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
			},
			flagKey: "headerColor",
			context: &structpb.Struct{Fields: map[string]*structpb.Value{
				"email": {Kind: &structpb.Value_StringValue{
					StringValue: "user@faas.com",
				}},
			}},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"error during parsing (invalid target version) - return default": {
			flags: Flags{
				Flags: map[string]model.Flag{
					"headerColor": {
						State:          "ENABLED",
						DefaultVariant: "red",
						Variants: map[string]any{
							"red":    "#FF0000",
							"blue":   "#0000FF",
							"green":  "#00FF00",
							"yellow": "#FFFF00",
						},
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
			},
			flagKey: "headerColor",
			context: &structpb.Struct{Fields: map[string]*structpb.Value{
				"email": {Kind: &structpb.Value_StringValue{
					StringValue: "user@faas.com",
				}},
			}},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
	}

	const reqID = "default"
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			je := NewJSONEvaluator(logger.NewLogger(nil, false), store.NewFlags())
			je.store.Flags = tt.flags.Flags

			value, variant, reason, err := resolve[string](
				reqID, tt.flagKey, tt.context, je.evaluateVariant, je.store.Flags[tt.flagKey].Variants,
			)

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
