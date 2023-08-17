package eval

import (
	"testing"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
)

func TestFractionalEvaluation(t *testing.T) {
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
												"fractionalEvaluation": [
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
		flags           Flags
		flagKey         string
		context         map[string]any
		expectedValue   string
		expectedVariant string
		expectedReason  string
		expectedError   error
	}{
		"test@faas.com": {
			flags:   flags,
			flagKey: "headerColor",
			context: map[string]any{
				"email": "test@faas.com",
			},
			expectedVariant: "green",
			expectedValue:   "#00FF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"test2@faas.com": {
			flags:   flags,
			flagKey: "headerColor",
			context: map[string]any{
				"email": "test2@faas.com",
			},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"test3@faas.com": {
			flags:   flags,
			flagKey: "headerColor",
			context: map[string]any{
				"email": "test3@faas.com",
			},
			expectedVariant: "yellow",
			expectedValue:   "#FFFF00",
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
												"fractionalEvaluation": [
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
							"fractionalEvaluation": [
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
		"fallback to default variant if invalid variant as result of fractional evaluation": {
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
							"fractionalEvaluation": [
								"email",
								[
								"black",
								100
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
			expectedReason:  model.DefaultReason,
		},
		"fallback to default variant if percentages don't sum to 100": {
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
							"fractionalEvaluation": [
								"email",
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
			expectedReason:  model.DefaultReason,
		},
	}
	const reqID = "default"
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			log := logger.NewLogger(nil, false)
			je := NewJSONEvaluator(
				log,
				store.NewFlags(),
				WithEvaluator(
					"fractionalEvaluation",
					NewFractionalEvaluator(log).FractionalEvaluation,
				),
			)
			je.store.Flags = tt.flags.Flags

			value, variant, reason, _, err := resolve[string](reqID, tt.flagKey, tt.context, je.evaluateVariant)

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

func BenchmarkFractionalEvaluation(b *testing.B) {
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
						"fractionalEvaluation": [
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
		flags           Flags
		flagKey         string
		context         map[string]any
		expectedValue   string
		expectedVariant string
		expectedReason  string
		expectedError   error
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
			je := NewJSONEvaluator(
				log,
				&store.Flags{Flags: tt.flags.Flags},
				WithEvaluator(
					"fractionalEvaluation",
					NewFractionalEvaluator(log).FractionalEvaluation,
				),
			)
			for i := 0; i < b.N; i++ {
				value, variant, reason, _, err := resolve[string](reqID, tt.flagKey, tt.context, je.evaluateVariant)

				if value != tt.expectedValue {
					b.Errorf("expected value '%s', got '%s'", tt.expectedValue, value)
				}

				if variant != tt.expectedVariant {
					b.Errorf("expected variant '%s', got '%s'", tt.expectedVariant, variant)
				}

				if reason != tt.expectedReason {
					b.Errorf("expected reason '%s', got '%s'", tt.expectedReason, reason)
				}

				if err != tt.expectedError {
					b.Errorf("expected err '%v', got '%v'", tt.expectedError, err)
				}
			}
		})
	}
}
