package eval

import (
	"testing"

	"github.com/open-feature/flagd/core/pkg/store"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestJSONEvaluator_startsWithEvaluation(t *testing.T) {
	tests := map[string]struct {
		flags           Flags
		flagKey         string
		context         *structpb.Struct
		expectedValue   string
		expectedVariant string
		expectedReason  string
		expectedError   error
	}{
		"two strings provided - match": {
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
												"starts_with": ["user@faas.com", "user@faas"]
											  },
											  "red", null
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
												"starts_with": [{"var": "email"}, "user@faas"]
											  },
											  "red", null
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
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"two strings provided - no match": {
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
												"starts_with": ["user@faas.com", "nope"]
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
												"starts_with": [{"var": "email"}, "nope"]
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

func TestJSONEvaluator_endsWithEvaluation(t *testing.T) {
	tests := map[string]struct {
		flags           Flags
		flagKey         string
		context         *structpb.Struct
		expectedValue   string
		expectedVariant string
		expectedReason  string
		expectedError   error
	}{
		"two strings provided - match": {
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
												"ends_with": ["user@faas.com", "faas.com"]
											  },
											  "red", null
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
												"ends_with": [{"var": "email"}, "faas.com"]
											  },
											  "red", null
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
			expectedVariant: "red",
			expectedValue:   "#FF0000",
			expectedReason:  model.TargetingMatchReason,
		},
		"two strings provided - no match": {
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
												"ends_with": ["user@faas.com", "nope"]
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
												"ends_with": [{"var": "email"}, "nope"]
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
