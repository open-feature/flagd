package evaluator

import (
	"context"
	"testing"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/stretchr/testify/assert"
)

const (
	emailField        = "email"
	localeField       = "locale"
	tierField         = "tier"
	targetingKeyField = "targetingKey"

	rachelEmail = "rachel@faas.com"
	monicaEmail = "monica@faas.com"
	joeyEmail   = "joey@faas.com"
	rossEmail   = "ross@faas.com"
	testAEmail  = "test_a@faas.com"
	testBEmail  = "test_b@faas.com"
	testCEmail  = "test_c@faas.com"
	testDEmail  = "test_d@faas.com"

	usLocale    = "us"
	caLocale    = "ca"
	premiumTier = "premium"

	redVariant    = "red"
	blueVariant   = "blue"
	greenVariant  = "green"
	yellowVariant = "yellow"

	redHex    = "#FF0000"
	blueHex   = "#0000FF"
	greenHex  = "#00FF00"
	yellowHex = "#FFFF00"

	testSource   = "testSource"
	defaultReqID = "default"
)

// setupTestEvaluator creates a new JSON evaluator with the given flags
func setupTestEvaluator(t *testing.T, flags []model.Flag) *JSON {
	log := logger.NewLogger(nil, false)
	s, err := store.NewStore(log, []string{testSource})
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}
	je := NewJSON(log, s)
	je.store.Update(testSource, flags, model.Metadata{})
	return je
}

func TestFractionalEvaluation(t *testing.T) {
	ctx := context.Background()

	commonFlags := []model.Flag{
		{
			Key:            "headerColor",
			State:          "ENABLED",
			DefaultVariant: redVariant,
			Variants:       colorVariants,
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
		{
			Key:            "customSeededHeaderColor",
			State:          "ENABLED",
			DefaultVariant: redVariant,
			Variants:       colorVariants,
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
	}

	tests := map[string]struct {
		flags             []model.Flag
		flagKey           string
		context           map[string]any
		expectedValue     string
		expectedVariant   string
		expectedReason    string
		expectedErrorCode string
	}{
		rachelEmail: {
			flags:   commonFlags,
			flagKey: "headerColor",
			context: map[string]any{
				emailField: rachelEmail,
			},
			expectedVariant: yellowVariant,
			expectedValue:   yellowHex,
			expectedReason:  model.TargetingMatchReason,
		},
		monicaEmail: {
			flags:   commonFlags,
			flagKey: "headerColor",
			context: map[string]any{
				emailField: monicaEmail,
			},
			expectedVariant: blueVariant,
			expectedValue:   blueHex,
			expectedReason:  model.TargetingMatchReason,
		},
		joeyEmail: {
			flags:   commonFlags,
			flagKey: "headerColor",
			context: map[string]any{
				emailField: joeyEmail,
			},
			expectedVariant: redVariant,
			expectedValue:   redHex,
			expectedReason:  model.TargetingMatchReason,
		},
		rossEmail: {
			flags:   commonFlags,
			flagKey: "headerColor",
			context: map[string]any{
				emailField: rossEmail,
			},
			expectedVariant: greenVariant,
			expectedValue:   greenHex,
			expectedReason:  model.TargetingMatchReason,
		},
		"rachel@faas.com with custom seed": {
			flags:   commonFlags,
			flagKey: "customSeededHeaderColor",
			context: map[string]any{
				"email": "rachel@faas.com",
			},
			expectedVariant: greenVariant,
			expectedValue:   greenHex,
			expectedReason:  model.TargetingMatchReason,
		},
		"monica@faas.com with custom seed": {
			flags:   commonFlags,
			flagKey: "customSeededHeaderColor",
			context: map[string]any{
				"email": "monica@faas.com",
			},
			expectedVariant: redVariant,
			expectedValue:   redHex,
			expectedReason:  model.TargetingMatchReason,
		},
		"joey@faas.com with custom seed": {
			flags:   commonFlags,
			flagKey: "customSeededHeaderColor",
			context: map[string]any{
				"email": "joey@faas.com",
			},
			expectedVariant: greenVariant,
			expectedValue:   greenHex,
			expectedReason:  model.TargetingMatchReason,
		},
		"ross@faas.com with custom seed": {
			flags:   commonFlags,
			flagKey: "customSeededHeaderColor",
			context: map[string]any{
				"email": rossEmail,
			},
			expectedVariant: greenVariant,
			expectedValue:   greenHex,
			expectedReason:  model.TargetingMatchReason,
		},
		"ross@faas.com with different flag key": {
			flags: []model.Flag{{
				Key:            "footerColor",
				State:          "ENABLED",
				DefaultVariant: redVariant,
				Variants:       colorVariants,
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
			flagKey: "footerColor",
			context: map[string]any{
				"email": rossEmail,
			},
			expectedVariant: redVariant,
			expectedValue:   redHex,
			expectedReason:  model.TargetingMatchReason,
		},
		"non even split": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: redVariant,
				Variants:       colorVariants,
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
			flagKey: "headerColor",
			context: map[string]any{
				"email": "test4@faas.com",
			},
			expectedVariant: redVariant,
			expectedValue:   redHex,
			expectedReason:  model.TargetingMatchReason,
		},
		"fallback to default variant if no email provided": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: redVariant,
				Variants:       colorVariants,
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
			flagKey:         "headerColor",
			context:         map[string]any{},
			expectedVariant: redVariant,
			expectedValue:   redHex,
			expectedReason:  model.DefaultReason,
		},
		"get variant for non-percentage weight values": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: redVariant,
				Variants:       colorVariants,
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
			flagKey: "headerColor",
			context: map[string]any{
				"email": "foo@foo.com",
			},
			expectedVariant: redVariant,
			expectedValue:   redHex,
			expectedReason:  model.TargetingMatchReason,
		},
		"get variant for non-specified weight values": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: redVariant,
				Variants:       colorVariants,
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
			flagKey: "headerColor",
			context: map[string]any{
				"email": "foo@foo.com",
			},
			expectedVariant: redVariant,
			expectedValue:   redHex,
			expectedReason:  model.TargetingMatchReason,
		},
		"default to targetingKey if no bucket key provided": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: redVariant,
				Variants:       colorVariants,
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
			flagKey: "headerColor",
			context: map[string]any{
				"targetingKey": "foo@foo.com",
			},
			expectedVariant: blueVariant,
			expectedValue:   blueHex,
			expectedReason:  model.TargetingMatchReason,
		},
		"missing email - parser should ignore nil/missing custom variables and continue": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: redVariant,
				Variants:       colorVariants,
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
			flagKey: "headerColor",
			context: map[string]any{
				"targetingKey": "foo@foo.com",
			},
			expectedVariant: redVariant,
			expectedValue:   redHex,
			expectedReason:  model.TargetingMatchReason,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			je := setupTestEvaluator(t, tt.flags)
			value, variant, reason, _, err := resolve[string](ctx, defaultReqID, tt.flagKey, tt.context, je.evaluateVariant)

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

func TestFractionalEvaluationWithNestedJSONLogic(t *testing.T) {
	ctx := context.Background()

	commonFlags := []model.Flag{
		{
			Key:            "nestedIfVariant",
			State:          "ENABLED",
			DefaultVariant: redVariant,
			Variants:       colorVariants,
			Targeting: []byte(`{
				"fractional": [
					{"var": "email"},
					[
						{
							"if": [
								{"in": ["us", {"var": "locale"}]},
								"blue",
								"red"
							]
						},
						25
					],
					[
						"red",
						75
					]
				]
			}`),
		},
		{
			Key:            "nestedFractional",
			State:          "ENABLED",
			DefaultVariant: redVariant,
			Variants:       colorVariants,
			Targeting: []byte(`{
				"fractional": [
					{"var": "email"},
					[
						{
							"fractional": [
								{"var": "tier"},
								["blue", 1]
							]
						},
						25
					],
					[
						"red",
						75
					]
				]
			}`),
		},
	}

	tests := map[string]struct {
		flags           []model.Flag
		flagKey         string
		context         map[string]any
		expectedVariant string
		expectedValue   string
		expectedReason  string
	}{
		"nested if - us locale in first bucket returns blue from nested if": {
			flags:   commonFlags,
			flagKey: "nestedIfVariant",
			context: map[string]any{
				emailField:  rossEmail,
				localeField: usLocale,
			},
			expectedVariant: blueVariant,
			expectedValue:   blueHex,
			expectedReason:  model.TargetingMatchReason,
		},
		"nested if - non-us locale in first bucket returns red from nested if": {
			flags:   commonFlags,
			flagKey: "nestedIfVariant",
			context: map[string]any{
				emailField:  rossEmail,
				localeField: caLocale,
			},
			expectedVariant: redVariant,
			expectedValue:   redHex,
			expectedReason:  model.TargetingMatchReason,
		},
		"nested fractional in first bucket returns blue from nested fractional": {
			flags:   commonFlags,
			flagKey: "nestedFractional",
			context: map[string]any{
				emailField: rossEmail,
				tierField:  premiumTier,
			},
			expectedVariant: blueVariant,
			expectedValue:   blueHex,
			expectedReason:  model.TargetingMatchReason,
		},
	}
	je := setupTestEvaluator(t, commonFlags)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			value, variant, reason, _, err := resolve[string](ctx, defaultReqID, tt.flagKey, tt.context, je.evaluateVariant)

			if tt.expectedVariant != "" && variant != tt.expectedVariant {
				t.Errorf("expected variant '%s', got '%s'", tt.expectedVariant, variant)
			}

			if tt.expectedValue != "" && value != tt.expectedValue {
				t.Errorf("expected value '%s', got '%s'", tt.expectedValue, value)
			}

			if reason != tt.expectedReason {
				t.Errorf("expected reason '%s', got '%s'", tt.expectedReason, reason)
			}

			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestFractionalVariantBoolNumericAndOperators(t *testing.T) {
	ctx := context.Background()

	commonFlags := []model.Flag{
		{
			Key:            "boolVariantTest",
			State:          "ENABLED",
			DefaultVariant: "error",
			Variants: map[string]any{
				"pass": "pass",
				"fail": "fail",
			},
			// this only return if true is returned from fractional
			Targeting: []byte(`{
				"if": [
					{
						"fractional": [
							{"var": "targetingKey"},
							[true, 1]
						]
					},
					"pass",
					"fail"
				]
			}`),
		},
		{
			Key:            "numericVariantTest",
			State:          "ENABLED",
			DefaultVariant: "error",
			Variants: map[string]any{
				"pass": "pass",
				"fail": "fail",
			},
			// this only passes if 1 is returned from fractional
			Targeting: []byte(`{
				"if": [
					{"===": [
						{
							"fractional": [
								{"var": "targetingKey"},
								[1, 1]
							]
						},
						1
					]},
					"pass",
					"fail"
				]
			}`),
		},
	}

	tests := map[string]struct {
		flags           []model.Flag
		flagKey         string
		context         map[string]any
		expectedVariant string
		expectedReason  string
	}{
		"bool variant returns true": {
			flags:   commonFlags,
			flagKey: "boolVariantTest",
			context: map[string]any{
				targetingKeyField: "test_user",
			},
			expectedVariant: "pass",
			expectedReason:  model.TargetingMatchReason,
		},
		"numeric variant returns one": {
			flags:   commonFlags,
			flagKey: "numericVariantTest",
			context: map[string]any{
				targetingKeyField: "another_user",
			},
			expectedVariant: "pass",
			expectedReason:  model.TargetingMatchReason,
		},
	}

	je := setupTestEvaluator(t, commonFlags)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_, variant, reason, _, err := resolve[string](ctx, defaultReqID, tt.flagKey, tt.context, je.evaluateVariant)

			if variant != tt.expectedVariant {
				t.Errorf("expected variant '%s', got '%s'", tt.expectedVariant, variant)
			}

			if reason != tt.expectedReason {
				t.Errorf("expected reason '%s', got '%s'", tt.expectedReason, reason)
			}

			if err != nil {
				t.Error(err)
			}
		})
	}
}

func BenchmarkFractionalEvaluation(b *testing.B) {
	ctx := context.Background()

	flags := []model.Flag{{
		Key:            "headerColor",
		State:          "ENABLED",
		DefaultVariant: redVariant,
		Variants:       colorVariants,
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
	}

	tests := map[string]struct {
		flags             []model.Flag
		flagKey           string
		context           map[string]any
		expectedValue     string
		expectedVariant   string
		expectedReason    string
		expectedErrorCode string
	}{
		testAEmail: {
			flags:   flags,
			flagKey: "headerColor",
			context: map[string]any{
				emailField: testAEmail,
			},
			expectedVariant: blueVariant,
			expectedValue:   blueHex,
			expectedReason:  model.TargetingMatchReason,
		},
		testBEmail: {
			flags:   flags,
			flagKey: "headerColor",
			context: map[string]any{
				emailField: testBEmail,
			},
			expectedVariant: redVariant,
			expectedValue:   redHex,
			expectedReason:  model.TargetingMatchReason,
		},
		testCEmail: {
			flags:   flags,
			flagKey: "headerColor",
			context: map[string]any{
				emailField: testCEmail,
			},
			expectedVariant: greenVariant,
			expectedValue:   greenHex,
			expectedReason:  model.TargetingMatchReason,
		},
		testDEmail: {
			flags:   flags,
			flagKey: "headerColor",
			context: map[string]any{
				emailField: testDEmail,
			},
			expectedVariant: blueVariant,
			expectedValue:   blueHex,
			expectedReason:  model.TargetingMatchReason,
		},
	}
	reqID := "test"
	for name, tt := range tests {
		b.Run(name, func(b *testing.B) {
			log := logger.NewLogger(nil, false)
			s, err := store.NewStore(log, []string{testSource})
			if err != nil {
				b.Fatalf("NewStore failed: %v", err)
			}
			je := NewJSON(log, s)
			je.store.Update(testSource, tt.flags, model.Metadata{})

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
