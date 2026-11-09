package evaluator

import (
	"context"
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

func TestResolveRejectedAtLoad(t *testing.T) {
	tests := map[string]string{
		"inside equality":       `{"==": [{"resolve": ["red"]}, "red"]}`,
		"inside cat":            `{"cat": [{"resolve": ["red"]}, "x"]}`,
		"if condition":          `{"if": [{"resolve": ["red"]}, "red", "blue"]}`,
		"consumed by outer if":  `{"if": [{"if": [true, {"resolve": ["red"]}, "blue"]}, "red", "blue"]}`,
		"inside fractional":     `{"fractional": ["key", [{"resolve": ["red"]}, 1], ["blue", 1]]}`,
		"nested in resolve arg": `{"resolve": [{"resolve": ["red"]}]}`,
		"result array":          `{"if": [true, [{"resolve": ["red"]}], "blue"]}`,
		"no arguments":          `{"resolve": []}`,
		"too many arguments":    `{"resolve": ["red", "ALLOWLIST", [], "extra"]}`,
		"non-string variant":    `{"resolve": [true]}`,
		"non-string reason":     `{"resolve": ["red", 1]}`,
		"metadata object":       `{"resolve": ["red", null, {"details": "x"}]}`,
		"short pair":            `{"resolve": ["red", null, [["details"]]]}`,
		"non-primitive value":   `{"resolve": ["red", null, [["details", {"nested": true}]]]}`,
		"non-string key":        `{"resolve": ["red", null, [[1, "x"]]]}`,
		"inlined bad evaluator": "",
	}

	for name, targeting := range tests {
		t.Run(name, func(t *testing.T) {
			flags := resolveFlag(targeting)
			if name == "inlined bad evaluator" {
				flags = `{
					"flags": {
						"color": {
							"state": "ENABLED",
							"defaultVariant": "yellow",
							"variants": {"red": "#FF0000", "yellow": "#FFFF00"},
							"targeting": {"==": [{"$ref": "bad"}, "red"]}
						}
					},
					"$evaluators": {"bad": {"resolve": ["red"]}}
				}`
			}
			log := logger.NewLogger(nil, false)
			je := NewJSON(log, store.NewFlags())
			err := je.SetState(sync.DataSync{FlagData: flags, Source: "testSource"})
			require.Error(t, err)
			assert.Contains(t, err.Error(), "color")
		})
	}
}

func TestParseEvaluatedResolve(t *testing.T) {
	reason := "ALLOWLIST"
	parsed, err := parseEvaluatedResolve([]any{"red", reason, []any{[]any{"details", "x"}, []any{"n", 2}}})
	require.NoError(t, err)
	assert.Equal(t, "red", parsed.Variant)
	require.NotNil(t, parsed.Reason)
	assert.Equal(t, reason, *parsed.Reason)
	assert.Equal(t, map[string]any{"details": "x", "n": float64(2)}, parsed.Metadata)

	_, err = parseEvaluatedResolve([]any{"red", nil, "nope"})
	require.Error(t, err)
}

func resolveFlag(targeting string) string {
	if targeting == "" {
		targeting = "null"
	}
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
