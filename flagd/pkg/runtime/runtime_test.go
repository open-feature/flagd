package runtime

import (
	"context"
	"testing"

	"github.com/open-feature/flagd/core/pkg/evaluator"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/sync"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validFlags = `{
  "flags": {
    "validFlag": {
      "state": "ENABLED",
      "variants": {
        "on": true,
        "off": false
      },
      "defaultVariant": "on"
    }
  }
}`

const invalidFlags = `{
  "flags": {
    "invalidFlag": {
      "notState": "ENABLED",
      "notVariants": {
        "on": true,
        "off": false
      },
      "notDefaultVariant": "on"
    }
  }
}`

// fakeSyncService is a minimal stand-in for flagsync.ISyncService so that
// updateAndEmit can be exercised without spinning up the gRPC sync server.
type fakeSyncService struct {
	emitted []string
}

func (f *fakeSyncService) Start(_ context.Context) error { return nil }

func (f *fakeSyncService) Emit(source string) {
	f.emitted = append(f.emitted, source)
}

func newRuntimeWithEvaluator(t *testing.T, strict bool) (*Runtime, *fakeSyncService) {
	t.Helper()
	log := logger.NewLogger(nil, false)
	flagStore := store.NewFlags()
	var opts []evaluator.JSONEvaluatorOption
	if strict {
		opts = append(opts, evaluator.WithStrictValidation())
	}
	eval := evaluator.NewJSON(log, flagStore, opts...)
	fakeSync := &fakeSyncService{}
	return &Runtime{
		Logger:           log,
		Evaluator:        eval,
		SyncService:      fakeSync,
		StrictValidation: strict,
	}, fakeSync
}

func TestUpdateAndEmit_StrictValidation_FailsOnInvalidInitialConfig(t *testing.T) {
	r, fakeSync := newRuntimeWithEvaluator(t, true)

	err := r.updateAndEmit(sync.DataSync{FlagData: invalidFlags, Source: "source-A"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "strict validation")
	assert.Contains(t, err.Error(), "source-A")
	assert.Empty(t, fakeSync.emitted, "no emit should happen when initial config is invalid")
}

func TestUpdateAndEmit_StrictValidation_TolaratesRuntimeUpdateFailure(t *testing.T) {
	r, fakeSync := newRuntimeWithEvaluator(t, true)

	// first, a valid bootstrap payload
	require.NoError(t, r.updateAndEmit(sync.DataSync{FlagData: validFlags, Source: "source-A"}))
	require.Equal(t, []string{"source-A"}, fakeSync.emitted)

	// later, an invalid runtime update from the same (already bootstrapped) source
	// must NOT cause an error to propagate; the prior valid state is preserved.
	err := r.updateAndEmit(sync.DataSync{FlagData: invalidFlags, Source: "source-A"})
	require.NoError(t, err)
	require.Equal(t, []string{"source-A"}, fakeSync.emitted, "emit must not fire for the failed update")
}

func TestUpdateAndEmit_StrictValidation_PerSourceBootstrap(t *testing.T) {
	r, _ := newRuntimeWithEvaluator(t, true)

	// source-A bootstraps successfully
	require.NoError(t, r.updateAndEmit(sync.DataSync{FlagData: validFlags, Source: "source-A"}))

	// source-B is brand new; an invalid first payload from it must still fail-fast
	// even though source-A has already bootstrapped.
	err := r.updateAndEmit(sync.DataSync{FlagData: invalidFlags, Source: "source-B"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "source-B")
}

func TestUpdateAndEmit_NoStrictValidation_TolaratesInvalidInitial(t *testing.T) {
	r, fakeSync := newRuntimeWithEvaluator(t, false)

	// without strict-validation, the legacy behavior is preserved: invalid initial
	// payloads are logged and the runtime continues. The store is empty, so Emit
	// still fires (mirroring pre-existing behavior of accepting with a warning).
	err := r.updateAndEmit(sync.DataSync{FlagData: invalidFlags, Source: "source-A"})
	require.NoError(t, err)
	// in non-strict mode, the underlying evaluator returns nil, so we treat the
	// payload as bootstrapped and emit.
	assert.Equal(t, []string{"source-A"}, fakeSync.emitted)
}

// fakeSync is a minimal sync.ISync used to exercise isReady without spinning
// up real sync providers.
type fakeSync struct {
	ready bool
}

func (f *fakeSync) Init(_ context.Context) error                       { return nil }
func (f *fakeSync) Sync(_ context.Context, _ chan<- sync.DataSync) error { return nil }
func (f *fakeSync) ReSync(_ context.Context, _ chan<- sync.DataSync) error {
	return nil
}
func (f *fakeSync) IsReady() bool { return f.ready }

func TestIsReady_StrictValidation_AllSourcesValidated_True(t *testing.T) {
	r, _ := newRuntimeWithEvaluator(t, true)
	r.Syncs = []sync.ISync{&fakeSync{ready: true}, &fakeSync{ready: true}}
	r.ExpectedSources = []string{"source-A", "source-B"}

	require.NoError(t, r.updateAndEmit(sync.DataSync{FlagData: validFlags, Source: "source-A"}))
	require.NoError(t, r.updateAndEmit(sync.DataSync{FlagData: validFlags, Source: "source-B"}))

	assert.True(t, r.isReady())
}

func TestIsReady_StrictValidation_OneSourceUnvalidated_False(t *testing.T) {
	r, _ := newRuntimeWithEvaluator(t, true)
	r.Syncs = []sync.ISync{&fakeSync{ready: true}, &fakeSync{ready: true}}
	r.ExpectedSources = []string{"source-A", "source-B"}

	// only source-A delivers a valid payload; source-B never has.
	require.NoError(t, r.updateAndEmit(sync.DataSync{FlagData: validFlags, Source: "source-A"}))

	assert.False(t, r.isReady(), "must not be ready until every configured source has validated")
}

func TestIsReady_NoStrictValidation_IgnoresValidation(t *testing.T) {
	r, _ := newRuntimeWithEvaluator(t, false)
	r.Syncs = []sync.ISync{&fakeSync{ready: true}, &fakeSync{ready: true}}
	r.ExpectedSources = []string{"source-A", "source-B"}

	// no source has validated, but transports are ready and strict mode is off.
	assert.True(t, r.isReady(), "non-strict mode must not gate readiness on validation")
}

func TestIsReady_StrictValidation_TransportNotReady_False(t *testing.T) {
	r, _ := newRuntimeWithEvaluator(t, true)
	r.Syncs = []sync.ISync{&fakeSync{ready: true}, &fakeSync{ready: false}}
	r.ExpectedSources = []string{"source-A", "source-B"}

	// even with all sources validated, transport-level not-ready wins.
	require.NoError(t, r.updateAndEmit(sync.DataSync{FlagData: validFlags, Source: "source-A"}))
	require.NoError(t, r.updateAndEmit(sync.DataSync{FlagData: validFlags, Source: "source-B"}))

	assert.False(t, r.isReady())
}
