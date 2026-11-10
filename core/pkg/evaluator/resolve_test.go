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
		"nested if branches stay distinguishable": {
			flags: resolveFlag(`{
				"if": [
					{"==": [{"var": "id"}, "a"]},
					{"resolve": ["red", "ALLOWLIST", [["details", "customer-a"]]]},
					{"==": [{"var": "id"}, "b"]},
					{"resolve": ["red", "SPLIT", [["details", "customer-b"]]]},
					{"resolve": ["blue", "FALLTHROUGH", [["details", "else"]]]}
				]
			}`),
			flagKey: "color",
			context: map[string]any{"id": "b"},
			value:   "#FF0000",
			variant: "red",
			reason:  "SPLIT",
			meta:    map[string]any{"details": "customer-b"},
		},
		"ternary result is terminal": {
			flags: resolveFlag(`{
				"?:": [
					{"==": [{"var": "id"}, "a"]},
					{"resolve": ["green", "ALLOWLIST"]},
					{"resolve": ["yellow", "FALLTHROUGH"]}
				]
			}`),
			flagKey: "color",
			context: map[string]any{"id": "z"},
			value:   "#FFFF00",
			variant: "yellow",
			reason:  "FALLTHROUGH",
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
		"computed non-string variant is a parse error": {
			flags:   resolveFlag(`{"resolve": [{"var": "pick"}]}`),
			flagKey: "color",
			context: map[string]any{"pick": 1},
			reason:  model.ErrorReason,
			err:     errors.New(model.ParseErrorCode),
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

func TestResolveInvalidFlagAtLoad(t *testing.T) {
	tests := map[string]struct {
		targeting  string
		evaluators string
	}{
		"inside equality":       {targeting: `{"==": [{"resolve": ["red"]}, "red"]}`},
		"inside cat":            {targeting: `{"cat": [{"resolve": ["red"]}, "x"]}`},
		"if condition":          {targeting: `{"if": [{"resolve": ["red"]}, "red", "blue"]}`},
		"consumed by outer if":  {targeting: `{"if": [{"if": [true, {"resolve": ["red"]}, "blue"]}, "red", "blue"]}`},
		"inside fractional":     {targeting: `{"fractional": ["key", [{"resolve": ["red"]}, 1], ["blue", 1]]}`},
		"nested in resolve arg": {targeting: `{"resolve": [{"resolve": ["red"]}]}`},
		"result array":          {targeting: `{"if": [true, [{"resolve": ["red"]}], "blue"]}`},
		"no arguments":          {targeting: `{"resolve": []}`},
		"too many arguments":    {targeting: `{"resolve": ["red", "ALLOWLIST", [], "extra"]}`},
		"non-string variant":    {targeting: `{"resolve": [true]}`},
		"non-string reason":     {targeting: `{"resolve": ["red", 1]}`},
		"metadata object":       {targeting: `{"resolve": ["red", null, {"details": "x"}]}`},
		"short pair":            {targeting: `{"resolve": ["red", null, [["details"]]]}`},
		"non-primitive value":   {targeting: `{"resolve": ["red", null, [["details", {"nested": true}]]]}`},
		"non-string key":        {targeting: `{"resolve": ["red", null, [[1, "x"]]]}`},
		"inlined bad evaluator": {
			targeting:  `{"==": [{"$ref": "bad"}, "red"]}`,
			evaluators: `{"bad": {"resolve": ["red"]}}`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			log := logger.NewLogger(nil, false)
			s, err := store.NewStore(log, []string{"testSource"})
			require.NoError(t, err)
			je := NewJSON(log, s)
			flags := resolveFlagWithSibling(tt.targeting, tt.evaluators)
			require.NoError(t, je.SetState(sync.DataSync{FlagData: flags, Source: "testSource"}))

			_, _, reason, _, err := je.ResolveStringValue(context.Background(), "default", "color", nil)
			assert.Equal(t, model.ErrorReason, reason)
			assert.Equal(t, errors.New(model.ParseErrorCode), err)

			// the rest of the payload still loads
			value, _, reason, _, err := je.ResolveBooleanValue(context.Background(), "default", "sibling", nil)
			require.NoError(t, err)
			assert.True(t, value)
			assert.Equal(t, model.StaticReason, reason)
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
		"computed +Inf metadata": {
			targeting: nonFiniteMetadataRule,
			context:   map[string]any{"hits": 1, "total": 0},
			err:       errors.New(model.ParseErrorCode),
		},
		"computed NaN metadata": {
			targeting: nonFiniteMetadataRule,
			context:   map[string]any{"hits": 0, "total": 0},
			err:       errors.New(model.ParseErrorCode),
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
			assert.Equal(t, errors.New(model.ParseErrorCode), value.Error)
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

func resolveFlagWithSibling(targeting, evaluators string) string {
	if evaluators == "" {
		evaluators = "{}"
	}
	return `{
		"flags": {
			"color": {
				"state": "ENABLED",
				"defaultVariant": "yellow",
				"variants": {"red": "#FF0000", "blue": "#0000FF", "yellow": "#FFFF00"},
				"targeting": ` + targeting + `
			},
			"sibling": {
				"state": "ENABLED",
				"defaultVariant": "on",
				"variants": {"on": true, "off": false}
			}
		},
		"$evaluators": ` + evaluators + `
	}`
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
