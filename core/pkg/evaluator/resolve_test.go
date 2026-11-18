package evaluator

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/sync"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveEvaluation(t *testing.T) {
	const source = "testSource"
	ctx := context.Background()
	nestedBranches := resolveFlag(`{
		"if": [
			{"==": [{"var": "id"}, "a"]},
			{"resolve": ["red", "ALLOWLIST", [["details", "customer-a"]]]},
			{"==": [{"var": "id"}, "b"]},
			{"resolve": ["red", "SPLIT", [["details", "customer-b"]]]},
			{"resolve": ["blue", "FALLTHROUGH", [["details", "else"]]]}
		]
	}`)

	tests := map[string]struct {
		flags    string
		metadata model.Metadata
		flagKey  string
		context  map[string]any
		value    string
		variant  string
		reason   string
		meta     map[string]any
		err      error
	}{
		"whole rule selects variant and keeps targeting match": {
			flags:   resolveFlag(`{"resolve": ["red"]}`),
			flagKey: "color",
			value:   "#FF0000",
			variant: "red",
			reason:  model.TargetingMatchReason,
		},
		"custom reason": {
			flags:   resolveFlag(`{"resolve": ["blue", "ALLOWLIST"]}`),
			flagKey: "color",
			value:   "#0000FF",
			variant: "blue",
			reason:  "ALLOWLIST",
		},
		// the reason is free-form; flagd's own reason names must not change the value or variant
		"FALLBACK reason keeps value and variant": {
			flags:   resolveFlag(`{"resolve": ["blue", "FALLBACK"]}`),
			flagKey: "color",
			value:   "#0000FF",
			variant: "blue",
			reason:  model.FallbackReason,
		},
		"DISABLED reason keeps value and variant": {
			flags:   resolveFlag(`{"resolve": ["blue", "DISABLED"]}`),
			flagKey: "color",
			value:   "#0000FF",
			variant: "blue",
			reason:  model.DisabledReason,
		},
		"ERROR reason keeps value and variant": {
			flags:   resolveFlag(`{"resolve": ["blue", "ERROR"]}`),
			flagKey: "color",
			value:   "#0000FF",
			variant: "blue",
			reason:  model.ErrorReason,
		},
		"null reason keeps targeting match and merges metadata": {
			flags:   resolveFlag(`{"resolve": ["green", null, [["details", "pilot"]]]}`),
			flagKey: "color",
			value:   "#00FF00",
			variant: "green",
			reason:  model.TargetingMatchReason,
			meta:    map[string]any{"details": "pilot"},
		},
		"rule metadata wins over flag and flag set": {
			flags: `{
				"metadata": {"details": "set", "owner": "set"},
				"flags": {
					"color": {
						"state": "ENABLED",
						"defaultVariant": "yellow",
						"variants": {"red": "#FF0000", "blue": "#0000FF", "green": "#00FF00", "yellow": "#FFFF00"},
						"metadata": {"details": "flag", "version": "flag"},
						"targeting": {"resolve": ["red", "ALLOWLIST", [["details", "rule"], ["onlyRule", true]]]}
					}
				}
			}`,
			flagKey: "color",
			value:   "#FF0000",
			variant: "red",
			reason:  "ALLOWLIST",
			meta: map[string]any{
				"details":  "rule",
				"owner":    "set",
				"version":  "flag",
				"onlyRule": true,
			},
		},
		// the first two branches serve the same variant, so only the reason and details tell them apart
		"nested if first branch": {
			flags:   nestedBranches,
			flagKey: "color",
			context: map[string]any{"id": "a"},
			value:   "#FF0000",
			variant: "red",
			reason:  "ALLOWLIST",
			meta:    map[string]any{"details": "customer-a"},
		},
		"nested if second branch": {
			flags:   nestedBranches,
			flagKey: "color",
			context: map[string]any{"id": "b"},
			value:   "#FF0000",
			variant: "red",
			reason:  "SPLIT",
			meta:    map[string]any{"details": "customer-b"},
		},
		"nested if else branch": {
			flags:   nestedBranches,
			flagKey: "color",
			context: map[string]any{"id": "z"},
			value:   "#0000FF",
			variant: "blue",
			reason:  "FALLTHROUGH",
			meta:    map[string]any{"details": "else"},
		},
		"single if clause is terminal": {
			flags:   resolveFlag(`{"if": [{"resolve": ["green", "ALLOWLIST"]}]}`),
			flagKey: "color",
			value:   "#00FF00",
			variant: "green",
			reason:  "ALLOWLIST",
		},
		// a fixed bucketing value and zero weights pick the bucket deterministically
		"fractional bucket variant is terminal": {
			flags: resolveFlag(`{"fractional": ["key",
				[{"resolve": ["red", "SPLIT", [["details", "treatment"]]]}, 1],
				[{"resolve": ["blue", "SPLIT", [["details", "holdback"]]]}, 0]
			]}`),
			flagKey: "color",
			value:   "#FF0000",
			variant: "red",
			reason:  "SPLIT",
			meta:    map[string]any{"details": "treatment"},
		},
		"fractional mixes plain and resolve buckets": {
			flags:   resolveFlag(`{"fractional": ["key", [{"resolve": ["red", "SPLIT"]}, 0], ["blue", 1]]}`),
			flagKey: "color",
			value:   "#0000FF",
			variant: "blue",
			reason:  model.TargetingMatchReason,
		},
		"fractional in an if result": {
			flags: resolveFlag(`{"if": [true,
				{"fractional": ["key", [{"resolve": ["green", "SPLIT", [["details", "treatment"]]]}, 1]]},
				"blue"
			]}`),
			flagKey: "color",
			value:   "#00FF00",
			variant: "green",
			reason:  "SPLIT",
			meta:    map[string]any{"details": "treatment"},
		},
		"computed variant reason and details": {
			flags: resolveFlag(`{
				"resolve": [
					{"if": [true, "red", "blue"]},
					{"cat": ["ALLOW", "LIST"]},
					[["details", {"cat": ["us-", {"var": "bucket"}]}]]
				]
			}`),
			flagKey: "color",
			context: map[string]any{"bucket": "10pct"},
			value:   "#FF0000",
			variant: "red",
			reason:  "ALLOWLIST",
			meta:    map[string]any{"details": "us-10pct"},
		},
		"boolean variant is coerced like a plain result": {
			flags:   coercedKeysFlag(`{"resolve": [true, "BETA"]}`),
			flagKey: "coerced",
			value:   "#00FF00",
			variant: "true",
			reason:  "BETA",
		},
		"computed number variant is coerced like a plain result": {
			flags:   coercedKeysFlag(`{"resolve": [{"var": "pick"}]}`),
			flagKey: "coerced",
			context: map[string]any{"pick": 1},
			value:   "#0000FF",
			variant: "1",
			reason:  model.TargetingMatchReason,
		},
		"computed metadata list": {
			flags:   resolveFlag(`{"resolve": ["red", null, {"var": "meta"}]}`),
			flagKey: "color",
			context: map[string]any{"meta": []any{[]any{"details", "from-context"}}},
			value:   "#FF0000",
			variant: "red",
			reason:  model.TargetingMatchReason,
			meta:    map[string]any{"details": "from-context"},
		},
		"computed metadata pair": {
			flags:   resolveFlag(`{"resolve": ["red", null, [{"var": "pair"}]]}`),
			flagKey: "color",
			context: map[string]any{"pair": []any{"details", "from-context"}},
			value:   "#FF0000",
			variant: "red",
			reason:  model.TargetingMatchReason,
			meta:    map[string]any{"details": "from-context"},
		},
		"computed metadata key": {
			flags:   resolveFlag(`{"resolve": ["red", null, [[{"cat": ["az-", {"var": "zone"}]}, "served"]]]}`),
			flagKey: "color",
			context: map[string]any{"zone": "a"},
			value:   "#FF0000",
			variant: "red",
			reason:  model.TargetingMatchReason,
			meta:    map[string]any{"az-a": "served"},
		},
		"numeric and boolean metadata": {
			flags:   resolveFlag(`{"resolve": ["red", "SPLIT", [["percent", 10], ["enabled", false]]]}`),
			flagKey: "color",
			value:   "#FF0000",
			variant: "red",
			reason:  "SPLIT",
			meta:    map[string]any{"percent": float64(10), "enabled": false},
		},
		"missing variant matches a plain string miss": {
			flags:   resolveFlag(`{"resolve": ["missing"]}`),
			flagKey: "color",
			reason:  model.ErrorReason,
			err:     errors.New(model.GeneralErrorCode),
		},
		"shared evaluator inlines a terminal resolve": {
			flags: `{
				"flags": {
					"color": {
						"state": "ENABLED",
						"defaultVariant": "yellow",
						"variants": {"red": "#FF0000", "blue": "#0000FF", "green": "#00FF00", "yellow": "#FFFF00"},
						"targeting": {"$ref": "picked"}
					}
				},
				"$evaluators": {"picked": {"resolve": ["blue", "ALLOWLIST"]}}
			}`,
			flagKey: "color",
			value:   "#0000FF",
			variant: "blue",
			reason:  "ALLOWLIST",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			log := logger.NewLogger(nil, false)
			s, err := store.NewStore(log, []string{source})
			require.NoError(t, err)
			je := NewJSON(log, s)
			require.NoError(t, je.SetState(sync.DataSync{FlagData: tt.flags, Source: source}))

			value, variant, reason, metadata, err := je.ResolveStringValue(ctx, "default", tt.flagKey, tt.context)
			assert.Equal(t, tt.value, value)
			assert.Equal(t, tt.variant, variant)
			assert.Equal(t, tt.reason, reason)
			assert.Equal(t, tt.err, err)
			for key, want := range tt.meta {
				assert.Equal(t, want, metadata[key])
			}
		})
	}
}

func TestResolveRuntimeErrors(t *testing.T) {
	tests := map[string]struct {
		targeting string
		context   map[string]any
		err       error
	}{
		// an empty variant fails like any missing variant, so a FALLBACK reason cannot make it look like flagd's own
		"empty variant": {
			targeting: `{"resolve": ["", "FALLBACK"]}`,
			err:       errors.New(model.GeneralErrorCode),
		},
		"computed empty variant": {
			targeting: `{"resolve": [{"var": "pick"}, "FALLBACK"]}`,
			context:   map[string]any{"pick": ""},
			err:       errors.New(model.GeneralErrorCode),
		},
		// context data shaped like an operator result must not be treated as one
		"context cannot forge a resolve result": {
			targeting: `{"var": "pick"}`,
			context: map[string]any{"pick": map[string]any{
				"variant": "red", "reason": "INJECTED", "metadata": map[string]any{"flagSetId": "spoofed"},
			}},
			err: errors.New(model.GeneralErrorCode),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			log := logger.NewLogger(nil, false)
			s, err := store.NewStore(log, []string{"testSource"})
			require.NoError(t, err)
			je := NewJSON(log, s)
			require.NoError(t, je.SetState(sync.DataSync{FlagData: resolveFlag(tt.targeting), Source: "testSource"}))

			_, variant, reason, metadata, err := je.ResolveStringValue(context.Background(), "default", "color", tt.context)
			assert.Empty(t, variant)
			assert.Equal(t, model.ErrorReason, reason)
			assert.Equal(t, tt.err, err)
			assert.NotContains(t, metadata, "flagSetId")
			_, encodeErr := json.Marshal(metadata)
			assert.NoError(t, encodeErr, "metadata must stay JSON-encodable")
		})
	}
}

// a single-key context map like {"resolve": [...]} returned from a branch must never reach the real resolve
// operator. jsonlogic before v3.10.2 re-applies such maps as operator calls (the flag then falls back to its
// default variant); later versions return them as values (an invalid variant). Either way the context's variant,
// reason and metadata must not come through.
func TestResolveContextCannotInvokeResolve(t *testing.T) {
	injected := map[string]any{"pick": map[string]any{
		"resolve": []any{"red", "INJECTED", []any{[]any{"flagSetId", "spoofed"}}},
	}}

	tests := map[string]string{
		"nested if":      `{"if": [true, {"var": "pick"}, "blue"]}`,
		"nested if else": `{"if": [false, "blue", {"var": "pick"}]}`,
		"nested ?:":      `{"?:": [true, {"var": "pick"}, "blue"]}`,
		"inside resolve": `{"resolve": [{"if": [true, {"var": "pick"}, "blue"]}]}`,
	}

	for name, targeting := range tests {
		t.Run(name, func(t *testing.T) {
			log := logger.NewLogger(nil, false)
			s, err := store.NewStore(log, []string{"testSource"})
			require.NoError(t, err)
			je := NewJSON(log, s)
			require.NoError(t, je.SetState(sync.DataSync{FlagData: resolveFlag(targeting), Source: "testSource"}))

			_, variant, reason, metadata, _ := je.ResolveStringValue(context.Background(), "default", "color", injected)
			assert.NotEqual(t, "red", variant)
			assert.NotEqual(t, "INJECTED", reason)
			assert.NotContains(t, metadata, "flagSetId")
		})
	}
}

// malformed computed arguments return null like the other custom operators, so the default variant is selected
func TestResolveMalformedComputedArguments(t *testing.T) {
	tests := map[string]struct {
		targeting string
		context   map[string]any
	}{
		"null variant": {targeting: `{"resolve": [{"var": "pick"}]}`},
		"object variant": {
			targeting: `{"resolve": [{"var": "pick"}]}`,
			context:   map[string]any{"pick": map[string]any{"a": 1}},
		},
		"non-string reason": {targeting: `{"resolve": ["red", {"var": "why"}]}`, context: map[string]any{"why": 1}},
		"non-primitive metadata": {
			targeting: `{"resolve": ["red", null, [["details", {"var": "obj"}]]]}`,
			context:   map[string]any{"obj": map[string]any{"a": 1}},
		},
		"inside an if result": {targeting: `{"if": [true, {"resolve": [{"var": "pick"}]}, "blue"]}`},
		"inside a fractional bucket": {
			targeting: `{"fractional": ["key", [{"resolve": [{"var": "pick"}]}, 1], ["blue", 0]]}`,
		},
		"+Inf metadata": {targeting: nonFiniteMetadataRule, context: map[string]any{"hits": 1, "total": 0}},
		"NaN metadata":  {targeting: nonFiniteMetadataRule, context: map[string]any{"hits": 0, "total": 0}},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			log := logger.NewLogger(nil, false)
			s, err := store.NewStore(log, []string{"testSource"})
			require.NoError(t, err)
			je := NewJSON(log, s)
			require.NoError(t, je.SetState(sync.DataSync{FlagData: resolveFlag(tt.targeting), Source: "testSource"}))

			value, variant, reason, metadata, err := je.ResolveStringValue(context.Background(), "default", "color", tt.context)
			require.NoError(t, err)
			assert.Equal(t, "#FFFF00", value)
			assert.Equal(t, "yellow", variant)
			assert.Equal(t, model.DefaultReason, reason)
			assert.NotContains(t, metadata, "share")
			assert.NotContains(t, metadata, "details")
		})
	}
}

// a computed ratio that is +Inf for total 0 and NaN for 0/0
const nonFiniteMetadataRule = `{"resolve": ["red", "SPLIT", [["share", {"/": [{"var": "hits"}, {"var": "total"}]}]]]}`

// one flag with non-finite computed metadata must not break encoding of a bulk evaluation
func TestResolveNonFiniteMetadataInBulk(t *testing.T) {
	log := logger.NewLogger(nil, false)
	s, err := store.NewStore(log, []string{"testSource"})
	require.NoError(t, err)
	je := NewJSON(log, s)
	flags := `{"flags": {
		"ratio": {"state": "ENABLED", "defaultVariant": "red", "variants": {"red": "#FF0000"},
			"targeting": ` + nonFiniteMetadataRule + `},
		"sibling": {"state": "ENABLED", "defaultVariant": "on", "variants": {"on": true, "off": false}}
	}}`
	require.NoError(t, je.SetState(sync.DataSync{FlagData: flags, Source: "testSource"}))

	values, _, err := je.ResolveAllValues(context.Background(), "default", map[string]any{"hits": 1, "total": 0})
	require.NoError(t, err)
	require.Len(t, values, 2)

	_, err = json.Marshal(values)
	require.NoError(t, err)
	for _, value := range values {
		if value.FlagKey == "sibling" {
			assert.Equal(t, true, value.Value)
			assert.NoError(t, value.Error)
		}
		if value.FlagKey == "ratio" {
			assert.NoError(t, value.Error)
			assert.Equal(t, model.DefaultReason, value.Reason)
			assert.NotContains(t, value.Metadata, "share")
		}
	}
}

func TestParseEvaluatedResolve(t *testing.T) {
	reason := "ALLOWLIST"
	parsed, err := parseEvaluatedResolve([]any{"red", reason, []any{[]any{"details", "x"}, []any{"n", 2}}})
	require.NoError(t, err)
	assert.Equal(t, "red", parsed.variant)
	require.NotNil(t, parsed.reason)
	assert.Equal(t, reason, *parsed.reason)
	assert.Equal(t, map[string]any{"details": "x", "n": float64(2)}, parsed.metadata)

	_, err = parseEvaluatedResolve([]any{"red", nil, "nope"})
	require.Error(t, err)
}

func resolveFlag(targeting string) string {
	return `{
		"flags": {
			"color": {
				"state": "ENABLED",
				"defaultVariant": "yellow",
				"variants": {
					"red": "#FF0000",
					"blue": "#0000FF",
					"green": "#00FF00",
					"yellow": "#FFFF00"
				},
				"targeting": ` + targeting + `
			}
		}
	}`
}

// coercedKeysFlag has variants named "true", "false" and "1", which only a coerced variant can select
func coercedKeysFlag(targeting string) string {
	return `{
		"flags": {
			"coerced": {
				"state": "ENABLED",
				"defaultVariant": "false",
				"variants": {"true": "#00FF00", "false": "#FF0000", "1": "#0000FF"},
				"targeting": ` + targeting + `
			}
		}
	}`
}
