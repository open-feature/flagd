package evaluator

import (
	"context"
	"testing"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var colorVariants = map[string]any{
	"red":    "#FF0000",
	"blue":   "#0000FF",
	"green":  "#00FF00",
	"yellow": "#FFFF00",
}

type stringFlagEvalTestCase struct {
	flags           []model.Flag
	flagKey         string
	context         map[string]any
	expectedValue   string
	expectedVariant string
	expectedReason  string
	expectedError   error
}

func runStringFlagEvalTests(t *testing.T, ctx context.Context, source string, sources []string, tests map[string]stringFlagEvalTestCase) {
	t.Helper()
	const reqID = "default"
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			log := logger.NewLogger(nil, false)
			s, err := store.NewStore(log, sources)
			require.NoError(t, err)
			je := NewJSON(log, s)
			je.store.Update(source, tt.flags, model.Metadata{}, false)

			value, variant, reason, _, err := resolve[string](ctx, reqID, tt.flagKey, tt.context, je.evaluateVariant)

			assert.Equal(t, tt.expectedValue, value)
			assert.Equal(t, tt.expectedVariant, variant)
			assert.Equal(t, tt.expectedReason, reason)
			assert.ErrorIs(t, err, tt.expectedError)
		})
	}
}

type errorFallbackTestCase struct {
	targeting string
	context   map[string]any
}

func runErrorFallbackTests(t *testing.T, ctx context.Context, source, flagKey string, tests map[string]errorFallbackTestCase) {
	t.Helper()
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			log := logger.NewLogger(nil, false)
			s, err := store.NewStore(log, []string{source})
			require.NoError(t, err)
			je := NewJSON(log, s)
			je.store.Update(source, []model.Flag{{
				Key:            flagKey,
				State:          "ENABLED",
				DefaultVariant: "fallback",
				Variants: map[string]any{
					"true":     "true",
					"false":    "false",
					"fallback": "fallback",
				},
				Targeting: []byte(tt.targeting),
			}}, model.Metadata{}, false)

			value, variant, reason, _, err := resolve[string](ctx, "default", flagKey, tt.context, je.evaluateVariant)
			assert.NoError(t, err)
			assert.Equal(t, "fallback", value)
			assert.Equal(t, "fallback", variant)
			assert.Equal(t, model.DefaultReason, reason)
		})
	}
}
