package eval

import (
	"testing"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
)

func TestLegacyFractionalEvaluation(t *testing.T) {
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
			flagKey:           "headerColor",
			context:           map[string]any{},
			expectedVariant:   "",
			expectedValue:     "",
			expectedReason:    model.ErrorReason,
			expectedErrorCode: model.ParseErrorCode,
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
			expectedVariant:   "",
			expectedValue:     "",
			expectedReason:    model.ErrorReason,
			expectedErrorCode: model.ParseErrorCode,
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
					NewLegacyFractionalEvaluator(log).LegacyFractionalEvaluation,
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

			if err != nil {
				errorCode := err.Error()
				if errorCode != tt.expectedErrorCode {
					t.Errorf("expected err '%v', got '%v'", tt.expectedErrorCode, err)
				}
			}
		})
	}
}
