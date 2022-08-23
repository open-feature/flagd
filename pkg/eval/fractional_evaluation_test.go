package eval

import (
	"testing"

	"github.com/open-feature/flagd/pkg/model"

	"google.golang.org/protobuf/types/known/structpb"
)

func TestFractionalEvaluation(t *testing.T) {
	flags := Flags{
		Flags: map[string]Flag{
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
		context         *structpb.Struct
		expectedValue   string
		expectedVariant string
		expectedReason  string
		expectedError   error
	}{
		"test@faas.com": {
			flags:   flags,
			flagKey: "headerColor",
			context: &structpb.Struct{Fields: map[string]*structpb.Value{
				"email": {Kind: &structpb.Value_StringValue{
					StringValue: "test@faas.com",
				}},
			}},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"test2@faas.com": {
			flags:   flags,
			flagKey: "headerColor",
			context: &structpb.Struct{Fields: map[string]*structpb.Value{
				"email": {Kind: &structpb.Value_StringValue{
					StringValue: "test2@faas.com",
				}},
			}},
			expectedVariant: "yellow",
			expectedValue:   "#FFFF00",
			expectedReason:  model.TargetingMatchReason,
		},
		"test3@faas.com": {
			flags:   flags,
			flagKey: "headerColor",
			context: &structpb.Struct{Fields: map[string]*structpb.Value{
				"email": {Kind: &structpb.Value_StringValue{
					StringValue: "test3@faas.com",
				}},
			}},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"test4@faas.com": {
			flags:   flags,
			flagKey: "headerColor",
			context: &structpb.Struct{Fields: map[string]*structpb.Value{
				"email": {Kind: &structpb.Value_StringValue{
					StringValue: "test4@faas.com",
				}},
			}},
			expectedVariant: "blue",
			expectedValue:   "#0000FF",
			expectedReason:  model.TargetingMatchReason,
		},
		"non even split": {
			flags: Flags{
				Flags: map[string]Flag{
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
			context: &structpb.Struct{Fields: map[string]*structpb.Value{
				"email": {Kind: &structpb.Value_StringValue{
					StringValue: "test4@faas.com",
				}},
			}},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"fallback to default variant if no email provided": {
			flags: Flags{
				Flags: map[string]Flag{
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
			context:         &structpb.Struct{},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.StaticReason,
		},
		"fallback to default variant if invalid variant as result of fractional evaluation": {
			flags: Flags{
				Flags: map[string]Flag{
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
			context: &structpb.Struct{Fields: map[string]*structpb.Value{
				"email": {Kind: &structpb.Value_StringValue{
					StringValue: "foo@foo.com",
				}},
			}},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.StaticReason,
		},
		"fallback to default variant if percentages don't sum to 100": {
			flags: Flags{
				Flags: map[string]Flag{
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
			context: &structpb.Struct{Fields: map[string]*structpb.Value{
				"email": {Kind: &structpb.Value_StringValue{
					StringValue: "foo@foo.com",
				}},
			}},
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.StaticReason,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			je := JSONEvaluator{state: tt.flags}

			value, variant, reason, err := resolve[string](
				tt.flagKey, tt.context, je.evaluateVariant, je.state.Flags[tt.flagKey].Variants,
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
