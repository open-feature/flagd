package evaluator

import (
	"context"
	"testing"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/stretchr/testify/assert"
)

func TestFractionalEvaluation(t *testing.T) {
	ctx := context.Background()

	commonFlags := Flags{
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
												"in": ["@faas.com", {
														"var": ["email"]
													  }]
											  },
											  {
												"fractional": [
												  {"cat": [{"var": "$flagd.flagKey"}, {"var": "email"}]},
												  [
												  "red",
												  25
												  ],
												  [
												  "blue",
												  25
												  ],
												  [
												  "green",
												  25
												  ],
												  [
												  "yellow",
												  25
												  ]
												]
											  }, null
											]
										  }`),
			},
			"customSeededHeaderColor": {
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
						"in": ["@faas.com", {
								"var": ["email"]
								}]
						},
						{
						"fractional": [
							{"cat": ["my-seed", {"var": "email"}]},
							["red",25],
							["blue",25],
							["green",25],
							["yellow",25]
						]
						}, null
					]				  
				}`),
			},
		},
	}

	tests := map[string]struct {
		flags             Flags
		flagKey           string
		context           map[string]any
		expectedValue     string
		expectedVariant   string
		expectedReason    string
		expectedErrorCode string
	}{
		"rachel@faas.com": {
			flags:   commonFlags,
			flagKey: "headerColor",
			context: map[string]any{
				"email": "rachel@faas.com",
			},
			expectedVariant: "yellow",
			expectedValue:   "#FFFF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"monica@faas.com": {
			flags:   commonFlags,
			flagKey: "headerColor",
			context: map[string]any{
				"email": "monica@faas.com",
			},
			expectedVariant: "blue",
			expectedValue:   "#0000FF",
			expectedReason:  model.TargetingMatchReason,
		},
		"joey@faas.com": {
			flags:   commonFlags,
			flagKey: "headerColor",
			context: map[string]any{
				"email": "joey@faas.com",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"ross@faas.com": {
			flags:   commonFlags,
			flagKey: "headerColor",
			context: map[string]any{
				"email": "ross@faas.com",
			},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"rachel@faas.com with custom seed": {
			flags:   commonFlags,
			flagKey: "customSeededHeaderColor",
			context: map[string]any{
				"email": "rachel@faas.com",
			},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"monica@faas.com with custom seed": {
			flags:   commonFlags,
			flagKey: "customSeededHeaderColor",
			context: map[string]any{
				"email": "monica@faas.com",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"joey@faas.com with custom seed": {
			flags:   commonFlags,
			flagKey: "customSeededHeaderColor",
			context: map[string]any{
				"email": "joey@faas.com",
			},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"ross@faas.com with custom seed": {
			flags:   commonFlags,
			flagKey: "customSeededHeaderColor",
			context: map[string]any{
				"email": "ross@faas.com",
			},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"ross@faas.com with different flag key": {
			flags: Flags{
				Flags: map[string]model.Flag{
					"footerColor": {
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
									"in": ["@faas.com", {
										"var": ["email"]
									}]
								},
								{
									"fractional": [
										{"var": "email"},
										[
										"red",
										25
									  	],
									  	[
										"blue",
										25
										],
										[
										"green",
										25
									  	],
									  	[
										"yellow",
										25
									  	]
									]
								}, null
							]
						}`),
					},
				},
			},
			flagKey: "footerColor",
			context: map[string]any{
				"email": "ross@faas.com",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"non even split": {
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
												"in": ["@faas.com", {
														"var": ["email"]
													  }]
											  },
											  {
												"fractional": [
												  "email",
												  [
													"red",
													50
												  ],
												  [
													"blue",
													25
												  ],
												  [
													"green",
													25
												  ]
												]
											  }, null
											]
										  }`),
					},
				},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"email": "test4@faas.com",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"fallback to default variant if no email provided": {
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
							"fractional": [
								{"var": "email"},
								[
								"red",
								25
								],
								[
								"blue",
								25
								],
								[
								"green",
								25
								],
								[
								"yellow",
								25
								]
							]
							}`),
					},
				},
			},
			flagKey:         "headerColor",
			context:         map[string]any{},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.DefaultReason,
		},
		"get variant for non-percentage weight values": {
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
							"fractional": [
								{"var": "email"},
								[
								"red",
								25
								],
								[
								"blue",
								25
								]
							]
							}`),
					},
				},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"email": "foo@foo.com",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"get variant for non-specified weight values": {
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
							"fractional": [
								{"var": "email"},
								[
								"red"
								],
								[
								"blue"
								]
							]
							}`),
					},
				},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"email": "foo@foo.com",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"default to targetingKey if no bucket key provided": {
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
							"fractional": [
								[
								"blue",
								50
								],
								[
								"green",
								50
								]
							]
							}`),
					},
				},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"targetingKey": "foo@foo.com",
			},
			expectedVariant: "blue",
			expectedValue:   "#0000FF",
			expectedReason:  model.TargetingMatchReason,
		},
		"missing email - parser should ignore nil/missing custom variables and continue": {
			flags: Flags{
				Flags: map[string]model.Flag{
					"headerColor": {
						State:          "ENABLED",
						DefaultVariant: "red",
						Variants: map[string]any{
							"red":  "#FF0000",
							"blue": "#0000FF",
						},
						Targeting: []byte(
							`{
								"fractional": [
									{"var": "email"},
									["red",50],
									["blue",50]
								]
							}`),
					},
				},
			},
			flagKey: "headerColor",
			context: map[string]any{
				"targetingKey": "foo@foo.com",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
	}
	const reqID = "default"
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			log := logger.NewLogger(nil, false)
			je := NewJSON(log, store.NewFlags())
			je.store.Flags = tt.flags.Flags

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

			if err != nil {
				errorCode := err.Error()
				if errorCode != tt.expectedErrorCode {
					t.Errorf("expected err '%v', got '%v'", tt.expectedErrorCode, err)
				}
			}
		})
	}
}

func BenchmarkFractionalEvaluation(b *testing.B) {
	ctx := context.Background()

	flags := Flags{
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
						"in": ["@faas.com", {
								"var": ["email"]
								}]
						},
						{
						"fractional": [
							"email",
							[
							"red",
							25
							],
							[
							"blue",
							25
							],
							[
							"green",
							25
							],
							[
							"yellow",
							25
							]
						]
						}, null
					]
					}`),
			},
		},
	}

	tests := map[string]struct {
		flags             Flags
		flagKey           string
		context           map[string]any
		expectedValue     string
		expectedVariant   string
		expectedReason    string
		expectedErrorCode string
	}{
		"test@faas.com": {
			flags:   flags,
			flagKey: "headerColor",
			context: map[string]any{
				"email": "test@faas.com",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"test2@faas.com": {
			flags:   flags,
			flagKey: "headerColor",
			context: map[string]any{
				"email": "test2@faas.com",
			},
			expectedVariant: "yellow",
			expectedValue:   "#FFFF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"test3@faas.com": {
			flags:   flags,
			flagKey: "headerColor",
			context: map[string]any{
				"email": "test3@faas.com",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"test4@faas.com": {
			flags:   flags,
			flagKey: "headerColor",
			context: map[string]any{
				"email": "test4@faas.com",
			},
			expectedVariant: "blue",
			expectedValue:   "#0000FF",
			expectedReason:  model.TargetingMatchReason,
		},
	}
	reqID := "test"
	for name, tt := range tests {
		b.Run(name, func(b *testing.B) {
			log := logger.NewLogger(nil, false)
			je := NewJSON(log, &store.State{Flags: tt.flags.Flags})
			for i := 0; i < b.N; i++ {
				value, variant, reason, _, err := resolve[string](
					ctx, reqID, tt.flagKey, tt.context, je.evaluateVariant)

				if value != tt.expectedValue {
					b.Errorf("expected value '%s', got '%s'", tt.expectedValue, value)
				}

				if variant != tt.expectedVariant {
					b.Errorf("expected variant '%s', got '%s'", tt.expectedVariant, variant)
				}

				if reason != tt.expectedReason {
					b.Errorf("expected reason '%s', got '%s'", tt.expectedReason, reason)
				}

				if err != nil {
					errorCode := err.Error()
					if errorCode != tt.expectedErrorCode {
						b.Errorf("expected err '%v', got '%v'", tt.expectedErrorCode, err)
					}
				}
			}
		})
	}
}

func Test_fractionalEvaluationVariant_getPercentage(t *testing.T) {
	type fields struct {
		variant string
		weight  int
	}
	type args struct {
		totalWeight int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name: "get percentage",
			fields: fields{
				weight: 10,
			},
			args: args{
				totalWeight: 20,
			},
			want: 50,
		},
		{
			name: "total weight 0",
			fields: fields{
				weight: 10,
			},
			args: args{
				totalWeight: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := fractionalEvaluationVariant{
				variant: tt.fields.variant,
				weight:  tt.fields.weight,
			}
			assert.Equalf(t, tt.want, v.getPercentage(tt.args.totalWeight), "getPercentage(%v)", tt.args.totalWeight)
		})
	}
}
