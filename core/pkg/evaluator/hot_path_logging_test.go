package evaluator

import (
	"context"
	"testing"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// per-eval failures driven by eval context must not log above debug (flagd#2018)
func TestHotPathLoggingIsDebugOnly(t *testing.T) {
	const source = "testSource"
	const flagKey = "flag"

	tests := map[string]struct {
		targeting string
		context   map[string]any
	}{
		"starts_with absent property": {
			targeting: `{"if": [{"starts_with": [{"var": "email"}, "abc"]}, "on", "off"]}`,
			context:   map[string]any{},
		},
		"starts_with wrong type": {
			targeting: `{"if": [{"starts_with": [{"var": "email"}, "abc"]}, "on", "off"]}`,
			context:   map[string]any{"email": 42},
		},
		"ends_with absent property": {
			targeting: `{"if": [{"ends_with": [{"var": "email"}, "abc"]}, "on", "off"]}`,
			context:   map[string]any{},
		},
		"sem_ver absent property": {
			targeting: `{"if": [{"sem_ver": [{"var": "version"}, ">=", "1.0.0"]}, "on", "off"]}`,
			context:   map[string]any{},
		},
		"sem_ver invalid version": {
			targeting: `{"if": [{"sem_ver": [{"var": "version"}, ">=", "1.0.0"]}, "on", "off"]}`,
			context:   map[string]any{"version": "not-a-version"},
		},
		"fractional missing targetingKey": {
			targeting: `{"fractional": [["on", 50], ["off", 50]]}`,
			context:   map[string]any{},
		},
		"fractional empty targetingKey": {
			targeting: `{"fractional": [["on", 50], ["off", 50]]}`,
			context:   map[string]any{"targetingKey": ""},
		},
		"$flagd in context": {
			targeting: `{"if": [{"==": [{"var": "email"}, "abc"]}, "on", "off"]}`,
			context:   map[string]any{"$flagd": map[string]any{"flagKey": "spoofed"}},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			core, logs := observer.New(zapcore.DebugLevel)
			log := logger.NewLogger(zap.New(core), false)
			s, err := store.NewStore(log, []string{source})
			require.NoError(t, err)
			je := NewJSON(log, s)
			je.store.Update(source, []model.Flag{{
				Key:            flagKey,
				State:          "ENABLED",
				DefaultVariant: "off",
				Variants:       map[string]any{"on": "on", "off": "off"},
				Targeting:      []byte(tt.targeting),
			}}, model.Metadata{}, false)
			logs.TakeAll()

			value, _, _, _, err := resolve[string](context.Background(), "default", flagKey, tt.context, je.evaluateVariant)
			require.NoError(t, err)
			assert.Equal(t, "off", value)

			for _, e := range logs.All() {
				assert.LessOrEqual(t, e.Level, zapcore.DebugLevel, "unexpected %s log: %s", e.Level, e.Message)
			}
		})
	}
}
