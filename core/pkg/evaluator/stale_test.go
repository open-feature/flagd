//nolint:wrapcheck
package evaluator_test

import (
	"context"
	"testing"

	flagdEvaluator "github.com/open-feature/flagd/core/pkg/evaluator"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/sync"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const staleSource = "testSource"

func staleEvaluator(t *testing.T, state *store.SourceState) *flagdEvaluator.JSON {
	t.Helper()

	e := flagdEvaluator.NewJSON(
		logger.NewLogger(nil, false),
		store.NewFlags(),
		flagdEvaluator.WithSourceState(state),
	)
	require.NoError(t, e.SetState(sync.DataSync{FlagData: flagConfig, Source: staleSource}))

	return e
}

// Flags are still served while their source is disconnected -- serving
// last-known-good data beats failing evaluations open -- but the reason tells
// the caller the value may no longer match the source of truth.
func TestStale_ReplacesSuccessfulReasons(t *testing.T) {
	tests := []struct {
		name        string
		flagKey     string
		evalCtx     map[string]interface{}
		freshReason string
		expectedVal bool
	}{
		{
			name:        "static resolution",
			flagKey:     StaticBoolFlag,
			freshReason: model.StaticReason,
			expectedVal: StaticBoolValue,
		},
		{
			name:        "targeting match",
			flagKey:     DynamicBoolFlag,
			evalCtx:     map[string]interface{}{ColorProp: ColorValue},
			freshReason: model.TargetingMatchReason,
			expectedVal: StaticBoolValue,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := store.NewSourceState()
			e := staleEvaluator(t, state)

			// baseline: a connected source keeps its usual reason
			val, _, reason, _, err := e.ResolveBooleanValue(context.TODO(), "default", test.flagKey, test.evalCtx)
			require.NoError(t, err)
			assert.Equal(t, test.freshReason, reason)
			assert.Equal(t, test.expectedVal, val)

			// the source drops; the value is unchanged but now reported stale
			state.SetStale(staleSource, true)

			val, _, reason, _, err = e.ResolveBooleanValue(context.TODO(), "default", test.flagKey, test.evalCtx)
			require.NoError(t, err)
			assert.Equal(t, model.StaleReason, reason, "a disconnected source must resolve as STALE")
			assert.Equal(t, test.expectedVal, val, "the last-known-good value must still be served")

			// the source recovers
			state.SetStale(staleSource, false)

			_, _, reason, _, err = e.ResolveBooleanValue(context.TODO(), "default", test.flagKey, test.evalCtx)
			require.NoError(t, err)
			assert.Equal(t, test.freshReason, reason, "reconnecting must restore the original reason")
		})
	}
}

// STALE reports uncertainty about a value that was resolved. An evaluation that
// failed has no value to be uncertain about, so its error reason must survive.
func TestStale_DoesNotMaskErrors(t *testing.T) {
	tests := []struct {
		name      string
		flagKey   string
		errorCode string
	}{
		{"missing flag", MissingFlag, model.FlagNotFoundErrorCode},
		{"type mismatch", StaticObjectFlag, model.TypeMismatchErrorCode},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := store.NewSourceState()
			state.SetStale(staleSource, true)
			e := staleEvaluator(t, state)

			_, _, reason, _, err := e.ResolveBooleanValue(context.TODO(), "default", test.flagKey, nil)

			assert.EqualError(t, err, test.errorCode)
			assert.Equal(t, model.ErrorReason, reason, "errors must keep ERROR, not be rewritten to STALE")
		})
	}
}

func TestStale_OnlyAffectsTheDisconnectedSource(t *testing.T) {
	state := store.NewSourceState()
	e := staleEvaluator(t, state)

	state.SetStale("some-other-source", true)

	_, _, reason, _, err := e.ResolveBooleanValue(context.TODO(), "default", StaticBoolFlag, nil)
	require.NoError(t, err)
	assert.Equal(t, model.StaticReason, reason, "an unrelated stale source must not affect this flag")
}

// Source tracking is opt-in. Without it flagd must behave exactly as it did
// before stale reporting existed.
func TestStale_NoSourceStateConfigured(t *testing.T) {
	e := flagdEvaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())
	require.NoError(t, e.SetState(sync.DataSync{FlagData: flagConfig, Source: staleSource}))

	_, _, reason, _, err := e.ResolveBooleanValue(context.TODO(), "default", StaticBoolFlag, nil)
	require.NoError(t, err)
	assert.Equal(t, model.StaticReason, reason)
}
