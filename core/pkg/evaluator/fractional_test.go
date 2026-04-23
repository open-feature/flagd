package evaluator

import (
	"context"
	"testing"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/stretchr/testify/assert"
)

// Test constants
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
	test4Email  = "test4@faas.com"
	fooEmail    = "foo@foo.com"

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
)

// setupEvaluator creates and initializes a JSON evaluator with the given flags
func setupEvaluator(source string, flags []model.Flag) (*JSON, error) {
	log := logger.NewLogger(nil, false)
	s, err := store.NewStore(log, []string{source})
	if err != nil {
		return nil, err
	}
	je := NewJSON(log, s)
	je.store.Update(source, flags, model.Metadata{}, false)
	return je, nil
}

func TestFractionalEvaluation(t *testing.T) {
	const source = "testSource"
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
			expectedVariant: blueVariant,
			expectedValue:   blueHex,
			expectedReason:  model.TargetingMatchReason,
		},
		monicaEmail: {
			flags:   commonFlags,
			flagKey: "headerColor",
			context: map[string]any{
				emailField: monicaEmail,
			},
			expectedVariant: yellowVariant,
			expectedValue:   yellowHex,
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
			expectedVariant: blueVariant,
			expectedValue:   blueHex,
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
			expectedVariant: blueVariant,
			expectedValue:   blueHex,
			expectedReason:  model.TargetingMatchReason,
		},
		"ross@faas.com with custom seed": {
			flags:   commonFlags,
			flagKey: "customSeededHeaderColor",
			context: map[string]any{
				"email": "ross@faas.com",
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
				"email": "ross@faas.com",
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
			expectedVariant: greenVariant,
			expectedValue:   greenHex,
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
			expectedVariant: blueVariant,
			expectedValue:   blueHex,
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
			expectedVariant: blueVariant,
			expectedValue:   blueHex,
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
			expectedVariant: greenVariant,
			expectedValue:   greenHex,
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
			expectedVariant: blueVariant,
			expectedValue:   blueHex,
			expectedReason:  model.TargetingMatchReason,
		},
		"null targetingKey returns default variant": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: redVariant,
				Variants:       colorVariants,
				Targeting: []byte(`{
							"fractional": [
								["blue", 50],
								["green", 50]
							]
						}`),
			}},
			flagKey: "headerColor",
			context: map[string]any{
				"targetingKey": nil,
			},
			expectedVariant: redVariant,
			expectedValue:   redHex,
			expectedReason:  model.DefaultReason,
		},
		"missing targetingKey returns default variant": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: redVariant,
				Variants:       colorVariants,
				Targeting: []byte(`{
							"fractional": [
								["blue", 50],
								["green", 50]
							]
						}`),
			}},
			flagKey:         "headerColor",
			context:         map[string]any{},
			expectedVariant: redVariant,
			expectedValue:   redHex,
			expectedReason:  model.DefaultReason,
		},
		"empty targetingKey returns default variant": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: redVariant,
				Variants:       colorVariants,
				Targeting: []byte(`{
							"fractional": [
								["blue", 50],
								["green", 50]
							]
						}`),
			}},
			flagKey: "headerColor",
			context: map[string]any{
				"targetingKey": "",
			},
			expectedVariant: redVariant,
			expectedValue:   redHex,
			expectedReason:  model.DefaultReason,
		},
		"single-entry always returns the sole variant": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: redVariant,
				Variants:       colorVariants,
				Targeting: []byte(`{
							"fractional": [
								["blue", 1]
							]
						}`),
			}},
			flagKey: "headerColor",
			context: map[string]any{
				"targetingKey": "any-user",
			},
			expectedVariant: blueVariant,
			expectedValue:   blueHex,
			expectedReason:  model.TargetingMatchReason,
		},
		"single-entry with explicit bucket-by always returns the sole variant": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: redVariant,
				Variants:       colorVariants,
				Targeting: []byte(`{
							"fractional": [
								{"var": "email"},
								["green", 100]
							]
						}`),
			}},
			flagKey: "headerColor",
			context: map[string]any{
				"email": "any@user.com",
			},
			expectedVariant: greenVariant,
			expectedValue:   greenHex,
			expectedReason:  model.TargetingMatchReason,
		},
		"single-entry shorthand without weight always returns the sole variant": {
			flags: []model.Flag{{
				Key:            "headerColor",
				State:          "ENABLED",
				DefaultVariant: redVariant,
				Variants:       colorVariants,
				Targeting: []byte(`{
							"fractional": [
								["yellow"]
							]
						}`),
			}},
			flagKey: "headerColor",
			context: map[string]any{
				"targetingKey": "any-user",
			},
			expectedVariant: yellowVariant,
			expectedValue:   yellowHex,
			expectedReason:  model.TargetingMatchReason,
		},
	}
	const reqID = "default"
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			je, err := setupEvaluator(source, tt.flags)
			if err != nil {
				t.Fatalf("setupEvaluator failed: %v", err)
			}

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
	const source = "testSource"
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
			je, err := setupEvaluator(source, tt.flags)
			if err != nil {
				b.Fatalf("setupEvaluator failed: %v", err)
			}

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
		weight  int32
	}
	type args struct {
		totalWeight int32
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

func TestFractionalEvaluationNegativeClamping(t *testing.T) {
	ctx := context.Background()
	flagKey := "clampedWeightFlag"

	evalContext := map[string]any{
		targetingKeyField: "some-targeting-key",
	}

	commonFlags := []model.Flag{
		{
			Key:            flagKey,
			State:          "ENABLED",
			DefaultVariant: blueVariant,
			Variants:       colorVariants,
			Targeting: []byte(`{
				"fractional": [
					[
						"red",
						-1000
					],
					[
						"green",
						1
					]
				]
			}`),
		},
	}

	je, err := setupEvaluator("testSource", commonFlags)
	if err != nil {
		t.Fatalf("setupEvaluator failed: %v", err)
	}

	value, variant, reason, _, err := resolve[string](ctx, "default", flagKey, evalContext, je.evaluateVariant)
	assert.Equal(t, greenVariant, variant)
	assert.Equal(t, greenHex, value)
	assert.Equal(t, model.TargetingMatchReason, reason)
	assert.NoError(t, err)
}

func TestFractionalEvaluationWithNestedJSONLogic(t *testing.T) {
	const source = "testSource"
	ctx := context.Background()

	commonFlags := []model.Flag{
		{
			Key:            "nestedIfVariant",
			State:          "ENABLED",
			DefaultVariant: redVariant,
			Variants:       colorVariants,
			Targeting: []byte(`{
				"fractional": [
					"email",
					[
						{
							"if": [
								{"in": ["us", {"var": "locale"}]},
								"red",
								"blue"
							]
						},
						25
					],
					[
						"green",
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
					"email",
					[
						{
							"fractional": [
								"tier",
								["red", 25],
								["blue", 25],
								["green", 25],
								["yellow", 25]
							]
						},
						25
					],
					[
						"green",
						75
					]
				]
			}`),
		},
		{
			Key:            "dynamicWeights",
			State:          "ENABLED",
			DefaultVariant: redVariant,
			Variants:       colorVariants,
			Targeting: []byte(`{
				"fractional": [
					"email",
					["red", {"var": "redWeight"}],
					["blue", {"var": "blueWeight"}],
					["green", {"var": "greenWeight"}]
				]
			}`),
		},
	}

	tests := map[string]struct {
		flags             []model.Flag
		flagKey           string
		context           map[string]any
		expectedVariant   string
		expectedValue     string
		expectedReason    string
		expectedErrorCode string
		validVariants     []string // for tests where exact variant depends on hash
	}{
		"nested if - us locale in second bucket returns static variant": {
			flags:   commonFlags,
			flagKey: "nestedIfVariant",
			context: map[string]any{
				emailField:  rachelEmail,
				localeField: usLocale,
			},
			expectedVariant: greenVariant,
			expectedValue:   greenHex,
			expectedReason:  model.TargetingMatchReason,
		},
		"nested if - non-us locale in second bucket returns static value": {
			flags:   commonFlags,
			flagKey: "nestedIfVariant",
			context: map[string]any{
				emailField:  rachelEmail,
				localeField: caLocale,
			},
			expectedVariant: greenVariant,
			expectedValue:   greenHex,
			expectedReason:  model.TargetingMatchReason,
		},
		"nested fractional in second bucket returns one of variants": {
			flags:   commonFlags,
			flagKey: "nestedFractional",
			context: map[string]any{
				emailField: rachelEmail,
				tierField:  premiumTier,
			},
			expectedVariant: greenVariant,
			expectedValue:   greenHex,
			expectedReason:  model.TargetingMatchReason,
		},
		"weights computed from context variables": {
			flags:   commonFlags,
			flagKey: "dynamicWeights",
			context: map[string]any{
				emailField:    testAEmail,
				"redWeight":   float64(0),
				"blueWeight":  float64(1),
				"greenWeight": float64(0),
			},
			expectedVariant: blueVariant,
			expectedValue:   blueHex,
			expectedReason:  model.TargetingMatchReason,
		},
		"weights computed: red": {
			flags:   commonFlags,
			flagKey: "dynamicWeights",
			context: map[string]any{
				emailField:    testAEmail,
				"redWeight":   float64(1),
				"blueWeight":  float64(0),
				"greenWeight": float64(0),
			},
			expectedVariant: redVariant,
			expectedValue:   redHex,
			expectedReason:  model.TargetingMatchReason,
		},
		"weights computed: green": {
			flags:   commonFlags,
			flagKey: "dynamicWeights",
			context: map[string]any{
				emailField:    testAEmail,
				"redWeight":   float64(0),
				"blueWeight":  float64(0),
				"greenWeight": float64(1),
			},
			expectedVariant: greenVariant,
			expectedValue:   greenHex,
			expectedReason:  model.TargetingMatchReason,
		},
	}
	const reqID = "default"
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			je, err := setupEvaluator(source, tt.flags)
			if err != nil {
				t.Fatalf("setupEvaluator failed: %v", err)
			}

			value, variant, reason, _, err := resolve[string](ctx, reqID, tt.flagKey, tt.context, je.evaluateVariant)

			if tt.expectedVariant != "" && variant != tt.expectedVariant {
				t.Errorf("expected variant '%s', got '%s'", tt.expectedVariant, variant)
			}

			// check valid variants if specified (for hash-dependent tests)
			if len(tt.validVariants) > 0 {
				valid := false
				for _, v := range tt.validVariants {
					if variant == v {
						valid = true
						break
					}
				}
				if !valid {
					t.Errorf("expected variant to be one of %v, got '%s'", tt.validVariants, variant)
				}
			}

			if tt.expectedValue != "" && value != tt.expectedValue {
				t.Errorf("expected value '%s', got '%s'", tt.expectedValue, value)
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

func TestFractionalVariantBoolNumericAndOperators(t *testing.T) {
	log := logger.NewLogger(nil, false)
	fractional := NewFractional(log)

	tests := []struct {
		name            string
		values          any
		data            any
		expected        any
		expectedOptions []any // for hash-dependent or operators variants
	}{
		{
			name: "bool variant true",
			values: []any{
				"user123",
				[]any{true, float64(50)},
				[]any{false, float64(50)},
			},
			data: map[string]any{
				flagdPropertiesKey: map[string]any{
					"flagKey":   "test",
					"timestamp": int64(0),
				},
			},
			expectedOptions: []any{true, false},
		},
		{
			name: "numeric variant 0",
			values: []any{
				"user789",
				[]any{float64(0), float64(33)},
				[]any{float64(1), float64(33)},
				[]any{float64(2), float64(34)},
			},
			data: map[string]any{
				flagdPropertiesKey: map[string]any{
					"flagKey":   "test",
					"timestamp": int64(0),
				},
			},
			expectedOptions: []any{float64(0), float64(1), float64(2)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fractional.Evaluate(tt.values, tt.data)

			// Check if result is one of the expected options
			found := false
			for _, v := range tt.expectedOptions {
				if result == v {
					found = true
					break
				}
			}
			assert.True(t, found, "expected one of %v, got %v", tt.expectedOptions, result)
		})
	}
}

func TestFractionalEvaluation_ErrorFallbackWhenUsedDirectly(t *testing.T) {
	const source = "testSource"
	ctx := context.Background()

	tests := map[string]struct {
		targeting string
		context   map[string]any
	}{
		"missing bucket key falls back": {
			targeting: `{
				"fractional": [
					{"var": "missing_key"},
					["one", 50],
					["two", 50]
				]
			}`,
			context: map[string]any{},
		},
		"all zero weights fall back": {
			targeting: `{
				"fractional": [
					{"var": "targetingKey"},
					["one", 0],
					["two", 0]
				]
			}`,
			context: map[string]any{"targetingKey": "any-user"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			je, err := setupEvaluator(source, []model.Flag{{
				Key:            "fractional-op-error-fallback",
				State:          "ENABLED",
				DefaultVariant: "fallback",
				Variants: map[string]any{
					"one":      "one",
					"two":      "two",
					"fallback": "fallback",
				},
				Targeting: []byte(tt.targeting),
			}})
			assert.NoError(t, err)

			value, variant, reason, _, err := resolve[string](ctx, "default", "fractional-op-error-fallback", tt.context, je.evaluateVariant)
			assert.NoError(t, err)
			assert.Equal(t, "fallback", value)
			assert.Equal(t, "fallback", variant)
			assert.Equal(t, model.DefaultReason, reason)
		})
	}
}
