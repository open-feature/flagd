package evaluator

import (
	"context"
	"fmt"
	"testing"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRolloutEvaluate(t *testing.T) {
	log := logger.NewLogger(nil, false)
	rollout := NewRollout(log)

	startTime := int64(1704067200) // Jan 1, 2024
	endTime := int64(1706745600)   // Jan 31, 2024

	tests := []struct {
		name     string
		values   any
		data     any
		expected any
	}{
		{
			name: "shorthand: before start returns nil (defaultVariant)",
			values: []any{
				float64(startTime),
				float64(endTime),
				"new",
			},
			data: map[string]any{
				targetingKeyKey: "user123",
				flagdPropertiesKey: map[string]any{
					"flagKey":   "testFlag",
					"timestamp": startTime - 1000,
				},
			},
			expected: nil,
		},
		{
			name: "shorthand: after end returns toVariant",
			values: []any{
				float64(startTime),
				float64(endTime),
				"new",
			},
			data: map[string]any{
				targetingKeyKey: "user123",
				flagdPropertiesKey: map[string]any{
					"flagKey":   "testFlag",
					"timestamp": endTime + 1000,
				},
			},
			expected: "new",
		},
		{
			name: "longhand: before start returns fromVariant",
			values: []any{
				float64(startTime),
				float64(endTime),
				"old",
				"new",
			},
			data: map[string]any{
				targetingKeyKey: "user123",
				flagdPropertiesKey: map[string]any{
					"flagKey":   "testFlag",
					"timestamp": startTime - 1000,
				},
			},
			expected: "old",
		},
		{
			name: "longhand: after end returns toVariant",
			values: []any{
				float64(startTime),
				float64(endTime),
				"old",
				"new",
			},
			data: map[string]any{
				targetingKeyKey: "user123",
				flagdPropertiesKey: map[string]any{
					"flagKey":   "testFlag",
					"timestamp": endTime + 1000,
				},
			},
			expected: "new",
		},
		{
			name: "string bucketBy (pre-evaluated by JSONLogic)",
			values: []any{
				"custom-bucket",
				float64(startTime),
				float64(endTime),
				"old",
				"new",
			},
			data: map[string]any{
				targetingKeyKey: "user123",
				flagdPropertiesKey: map[string]any{
					"flagKey":   "testFlag",
					"timestamp": endTime + 1000,
				},
			},
			expected: "new",
		},
		{
			name: "with string bucketBy",
			values: []any{
				// JSONLogic pre-evaluates custom operator args, so bucketBy
				// expressions like {"var": "email"} arrive as their resolved string.
				"test@example.com",
				float64(startTime),
				float64(endTime),
				"old",
				"new",
			},
			data: map[string]any{
				targetingKeyKey: "user123",
				"email":         "test@example.com",
				flagdPropertiesKey: map[string]any{
					"flagKey":   "testFlag",
					"timestamp": endTime + 1000,
				},
			},
			expected: "new",
		},
		{
			name: "bool from variant: before start",
			values: []any{
				float64(startTime),
				float64(endTime),
				true,
				false,
			},
			data: map[string]any{
				targetingKeyKey: "user123",
				flagdPropertiesKey: map[string]any{
					"flagKey":   "testFlag",
					"timestamp": startTime - 1000,
				},
			},
			expected: true,
		},
		{
			name: "bool to variant: after end",
			values: []any{
				float64(startTime),
				float64(endTime),
				true,
				false,
			},
			data: map[string]any{
				targetingKeyKey: "user123",
				flagdPropertiesKey: map[string]any{
					"flagKey":   "testFlag",
					"timestamp": endTime + 1000,
				},
			},
			expected: false,
		},
		{
			name: "numeric from variant: before start",
			values: []any{
				float64(startTime),
				float64(endTime),
				float64(0),
				float64(1),
			},
			data: map[string]any{
				targetingKeyKey: "user123",
				flagdPropertiesKey: map[string]any{
					"flagKey":   "testFlag",
					"timestamp": startTime - 1000,
				},
			},
			expected: float64(0),
		},
		{
			name: "numeric to variant: after end",
			values: []any{
				float64(startTime),
				float64(endTime),
				float64(0),
				float64(1),
			},
			data: map[string]any{
				targetingKeyKey: "user123",
				flagdPropertiesKey: map[string]any{
					"flagKey":   "testFlag",
					"timestamp": endTime + 1000,
				},
			},
			expected: float64(1),
		},
		{
			name:   "invalid: not an array",
			values: "not an array",
			data: map[string]any{
				targetingKeyKey: "user123",
			},
			expected: nil,
		},
		{
			name: "invalid: too few arguments",
			values: []any{
				float64(startTime),
				float64(endTime),
			},
			data:     map[string]any{targetingKeyKey: "user123"},
			expected: nil,
		},
		{
			name: "invalid: endTime before startTime",
			values: []any{
				float64(endTime),
				float64(startTime),
				"new",
			},
			data: map[string]any{
				targetingKeyKey: "user123",
				flagdPropertiesKey: map[string]any{
					"flagKey":   "testFlag",
					"timestamp": int64(1705406400),
				},
			},
			expected: nil,
		},
		{
			name: "during rollout: shorthand returns nil (from=defaultVariant) for this user hash",
			values: []any{
				float64(startTime),
				float64(endTime),
				"new",
			},
			data: map[string]any{
				targetingKeyKey: "user123",
				flagdPropertiesKey: map[string]any{
					"flagKey":   "testFlag",
					"timestamp": startTime + (endTime-startTime)/2, // midpoint
				},
			},
			expected: nil, // This user's hash falls in the "from" bucket during rollout
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rollout.Evaluate(tt.values, tt.data)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRolloutDeterminism(t *testing.T) {
	log := logger.NewLogger(nil, false)
	rollout := NewRollout(log)

	startTime := int64(1704067200)
	endTime := int64(1706745600)
	midTime := startTime + (endTime-startTime)/2

	values := []any{
		float64(startTime),
		float64(endTime),
		"old",
		"new",
	}
	data := map[string]any{
		targetingKeyKey: "same-user@example.com",
		flagdPropertiesKey: map[string]any{
			"flagKey":   "testFlag",
			"timestamp": midTime,
		},
	}

	firstResult := rollout.Evaluate(values, data)

	for i := 0; i < 100; i++ {
		result := rollout.Evaluate(values, data)
		assert.Equal(t, firstResult, result, "result should be deterministic")
	}
}

func TestRolloutTimeProgression(t *testing.T) {
	log := logger.NewLogger(nil, false)
	rollout := NewRollout(log)

	startTime := int64(1704067200)
	endTime := int64(1706745600)
	duration := endTime - startTime

	tests := []struct {
		name        string
		currentTime int64
		minNewPct   float64
		maxNewPct   float64
	}{
		{
			name:        "at start (0%)",
			currentTime: startTime,
			minNewPct:   0,
			maxNewPct:   0,
		},
		{
			name:        "25% through",
			currentTime: startTime + duration/4,
			minNewPct:   0.20,
			maxNewPct:   0.30,
		},
		{
			name:        "50% through",
			currentTime: startTime + duration/2,
			minNewPct:   0.45,
			maxNewPct:   0.55,
		},
		{
			name:        "75% through",
			currentTime: startTime + 3*duration/4,
			minNewPct:   0.70,
			maxNewPct:   0.80,
		},
		{
			name:        "at end (100%)",
			currentTime: endTime,
			minNewPct:   1.0,
			maxNewPct:   1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iterations := 5000
			newCount := 0

			for i := 0; i < iterations; i++ {
				// Use default bucketBy (flagKey + targetingKey)
				values := []any{
					float64(startTime),
					float64(endTime),
					"old",
					"new",
				}
				data := map[string]any{
					targetingKeyKey: fmt.Sprintf("user-%d", i),
					flagdPropertiesKey: map[string]any{
						"flagKey":   "testFlag",
						"timestamp": tt.currentTime,
					},
				}

				result := rollout.Evaluate(values, data)
				if result == "new" {
					newCount++
				}
			}

			actualPct := float64(newCount) / float64(iterations)
			if actualPct < tt.minNewPct || actualPct > tt.maxNewPct {
				t.Errorf("expected %.0f%%-%.0f%% 'new', got %.1f%%",
					tt.minNewPct*100, tt.maxNewPct*100, actualPct*100)
			}
		})
	}
}

func TestRolloutIntegration(t *testing.T) {
	const source = "testSource"
	ctx := context.Background()

	flags := []model.Flag{
		{
			Key:            "rolloutFeature",
			State:          "ENABLED",
			DefaultVariant: "old",
			Variants: map[string]any{
				"old": "old-value",
				"new": "new-value",
			},
			Targeting: []byte(`{
				"rollout": [1704067200, 1706745600, "new"]
			}`),
		},
		{
			Key:            "rolloutFeatureLonghand",
			State:          "ENABLED",
			DefaultVariant: "old",
			Variants: map[string]any{
				"old": "old-value",
				"new": "new-value",
			},
			Targeting: []byte(`{
				"rollout": [1704067200, 1706745600, "old", "new"]
			}`),
		},
		{
			Key:            "conditionalRollout",
			State:          "ENABLED",
			DefaultVariant: "old",
			Variants: map[string]any{
				"old": "old-value",
				"new": "new-value",
			},
			Targeting: []byte(`{
				"if": [
					{"ends_with": [{"var": "email"}, "@internal.com"]},
					"new",
					{"rollout": [1704067200, 1706745600, "old", "new"]}
				]
			}`),
		},
		{
			// bucketBy uses {"var": "orgId"} — a JSONLogic operator that gets pre-evaluated to a string
			Key:            "rolloutOperatorInBucketBy",
			State:          "ENABLED",
			DefaultVariant: "old",
			Variants: map[string]any{
				"old": "old-value",
				"new": "new-value",
			},
			Targeting: []byte(`{
				"rollout": [{"var": "orgId"}, 1704067200, 1706745600, "old", "new"]
			}`),
		},
		{
			// from/to use {"cat": [...]} — operators that get pre-evaluated to strings
			Key:            "rolloutOperatorInVariants",
			State:          "ENABLED",
			DefaultVariant: "default",
			Variants: map[string]any{
				"default": "default-value",
				"old":     "old-value",
				"new":     "new-value",
			},
			Targeting: []byte(`{
				"rollout": [1704067200, 1706745600, {"cat": ["ol","d"]}, {"cat": ["ne","w"]}]
			}`),
		},
	}

	log := logger.NewLogger(nil, false)
	s := store.NewFlags()
	s.Update(source, flags, nil)
	je := NewJSON(log, s)

	testCases := []struct {
		name         string
		flagKey      string
		context      map[string]any
		expectedUser string
	}{
		{
			name:         "shorthand rollout after end returns new",
			flagKey:      "rolloutFeature",
			context:      map[string]any{targetingKeyKey: "user123"},
			expectedUser: "new",
		},
		{
			name:         "longhand rollout after end returns new",
			flagKey:      "rolloutFeatureLonghand",
			context:      map[string]any{targetingKeyKey: "user123"},
			expectedUser: "new",
		},
		{
			name:    "conditional rollout for internal user",
			flagKey: "conditionalRollout",
			context: map[string]any{
				targetingKeyKey: "user123",
				"email":         "test@internal.com",
			},
			expectedUser: "new",
		},
		{
			name:         "operator in bucketBy",
			flagKey:      "rolloutOperatorInBucketBy",
			context:      map[string]any{targetingKeyKey: "user1", "orgId": "org1"},
			expectedUser: "new",
		},
		{
			name:         "operator in from/to variants",
			flagKey:      "rolloutOperatorInVariants",
			context:      map[string]any{targetingKeyKey: "user1"},
			expectedUser: "new",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			value, variant, reason, _, err := je.ResolveStringValue(ctx, tt.name, tt.flagKey, tt.context)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedUser, variant)
			assert.Equal(t, tt.expectedUser+"-value", value)
			assert.Equal(t, model.TargetingMatchReason, reason)
		})
	}
}
