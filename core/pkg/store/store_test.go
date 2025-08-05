package store

import (
	"context"
	"testing"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateFlags(t *testing.T) {

	const source1 = "source1"
	const source2 = "source2"
	var sources = []string{source1, source2}

	t.Parallel()
	tests := []struct {
		name        string
		setup       func(t *testing.T) *Store
		newFlags    map[string]model.Flag
		source      string
		wantFlags   map[string]model.Flag
		setMetadata model.Metadata
		wantNotifs  map[string]interface{}
		wantResync  bool
	}{
		{
			name: "both nil",
			setup: func(t *testing.T) *Store {
				s, err := NewStore(logger.NewLogger(nil, false), sources)
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}
				return s
			},
			source:     source1,
			newFlags:   nil,
			wantFlags:  map[string]model.Flag{},
			wantNotifs: map[string]interface{}{},
		},
		{
			name: "both empty flags",
			setup: func(t *testing.T) *Store {
				s, err := NewStore(logger.NewLogger(nil, false), sources)
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}
				return s
			},
			source:     source1,
			newFlags:   map[string]model.Flag{},
			wantFlags:  map[string]model.Flag{},
			wantNotifs: map[string]interface{}{},
		},
		{
			name: "empty new",
			setup: func(t *testing.T) *Store {
				s, err := NewStore(logger.NewLogger(nil, false), sources)
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}
				return s
			},
			source:     source1,
			newFlags:   nil,
			wantFlags:  map[string]model.Flag{},
			wantNotifs: map[string]interface{}{},
		},
		{
			name: "update from source 1 (old flag removed)",
			setup: func(t *testing.T) *Store {
				s, err := NewStore(logger.NewLogger(nil, false), sources)
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}
				s.Update(source1, map[string]model.Flag{
					"waka": {DefaultVariant: "off"},
				}, nil)
				return s
			},
			newFlags: map[string]model.Flag{
				"paka": {DefaultVariant: "on"},
			},
			source: source1,
			wantFlags: map[string]model.Flag{
				"paka": {Key: "paka", DefaultVariant: "on", Source: source1, FlagSetId: nilFlagSetId, Priority: 0},
			},
			wantNotifs: map[string]interface{}{
				"paka": map[string]interface{}{"type": "write"},
				"waka": map[string]interface{}{"type": "delete"},
			},
		},
		{
			name: "update from source 1 (new flag added)",
			setup: func(t *testing.T) *Store {
				s, err := NewStore(logger.NewLogger(nil, false), sources)
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}
				s.Update(source1, map[string]model.Flag{
					"waka": {DefaultVariant: "off"},
				}, nil)
				return s
			},
			newFlags: map[string]model.Flag{
				"paka": {DefaultVariant: "on"},
			},
			source: source2,
			wantFlags: map[string]model.Flag{
				"waka": {Key: "waka", DefaultVariant: "off", Source: source1, FlagSetId: nilFlagSetId, Priority: 0},
				"paka": {Key: "paka", DefaultVariant: "on", Source: source2, FlagSetId: nilFlagSetId, Priority: 1},
			},
			wantNotifs: map[string]interface{}{"paka": map[string]interface{}{"type": "write"}},
		},
		{
			name: "flag set inheritance",
			setup: func(t *testing.T) *Store {
				s, err := NewStore(logger.NewLogger(nil, false), sources)
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}
				s.Update(source1, map[string]model.Flag{}, model.Metadata{})
				return s
			},
			setMetadata: model.Metadata{
				"flagSetId": "topLevelSet", // top level set metadata, including flagSetId
			},
			newFlags: map[string]model.Flag{
				"waka": {DefaultVariant: "on"},
				"paka": {DefaultVariant: "on", Metadata: model.Metadata{"flagSetId": "flagLevelSet"}}, // overrides set level flagSetId
			},
			source: source1,
			wantFlags: map[string]model.Flag{
				"waka": {Key: "waka", DefaultVariant: "on", Source: source1, FlagSetId: "topLevelSet", Priority: 0, Metadata: model.Metadata{"flagSetId": "topLevelSet"}},
				"paka": {Key: "paka", DefaultVariant: "on", Source: source1, FlagSetId: "flagLevelSet", Priority: 0, Metadata: model.Metadata{"flagSetId": "flagLevelSet"}},
			},
			wantNotifs: map[string]interface{}{
				"paka": map[string]interface{}{"type": "write"},
				"waka": map[string]interface{}{"type": "write"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			store := tt.setup(t)
			gotNotifs, resyncRequired := store.Update(tt.source, tt.newFlags, tt.setMetadata)
			gotFlags, _, _ := store.GetAll(context.Background(), nil, nil)

			require.Equal(t, tt.wantFlags, gotFlags)
			require.Equal(t, tt.wantNotifs, gotNotifs)
			require.Equal(t, tt.wantResync, resyncRequired)
		})
	}
}

func TestGet(t *testing.T) {

	sourceA := "sourceA"
	sourceB := "sourceB"
	sourceC := "sourceC"
	flagSetIdB := "flagSetIdA"
	flagSetIdC := "flagSetIdC"
	var sources = []string{sourceA, sourceB, sourceC}

	sourceASelector := NewSelector("source=" + sourceA)
	flagSetIdCSelector := NewSelector("flagSetId=" + flagSetIdC)

	t.Parallel()
	tests := []struct {
		name      string
		key       string
		selector  *Selector
		wantFlag  model.Flag
		wantFound bool
	}{
		{
			name:      "nil selector",
			key:       "flagA",
			selector:  nil,
			wantFlag:  model.Flag{Key: "flagA", DefaultVariant: "off", Source: sourceA, FlagSetId: nilFlagSetId, Priority: 0},
			wantFound: true,
		},
		{
			name:      "flagSetId selector",
			key:       "dupe",
			selector:  &flagSetIdCSelector,
			wantFlag:  model.Flag{Key: "dupe", DefaultVariant: "off", Source: sourceC, FlagSetId: flagSetIdC, Priority: 2, Metadata: model.Metadata{"flagSetId": flagSetIdC}},
			wantFound: true,
		},
		{
			name:      "source selector",
			key:       "dupe",
			selector:  &sourceASelector,
			wantFlag:  model.Flag{Key: "dupe", DefaultVariant: "on", Source: sourceA, FlagSetId: nilFlagSetId, Priority: 0},
			wantFound: true,
		},
		{
			name:      "flag not found with source selector",
			key:       "flagB",
			selector:  &sourceASelector,
			wantFlag:  model.Flag{Key: "flagB", DefaultVariant: "off", Source: sourceB, FlagSetId: flagSetIdB, Priority: 1, Metadata: model.Metadata{"flagSetId": flagSetIdB}},
			wantFound: false,
		},
		{
			name:      "flag not found with flagSetId selector",
			key:       "flagB",
			selector:  &flagSetIdCSelector,
			wantFlag:  model.Flag{Key: "flagB", DefaultVariant: "off", Source: sourceB, FlagSetId: flagSetIdB, Priority: 1, Metadata: model.Metadata{"flagSetId": flagSetIdB}},
			wantFound: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sourceAFlags := map[string]model.Flag{
				"flagA": {Key: "flagA", DefaultVariant: "off"},
				"dupe":  {Key: "dupe", DefaultVariant: "on"},
			}
			sourceBFlags := map[string]model.Flag{
				"flagB": {Key: "flagB", DefaultVariant: "off", Metadata: model.Metadata{"flagSetId": flagSetIdB}},
			}
			sourceCFlags := map[string]model.Flag{
				"flagC": {Key: "flagC", DefaultVariant: "off", Metadata: model.Metadata{"flagSetId": flagSetIdC}},
				"dupe":  {Key: "dupe", DefaultVariant: "off", Metadata: model.Metadata{"flagSetId": flagSetIdC}},
			}

			store, err := NewStore(logger.NewLogger(nil, false), sources)
			if err != nil {
				t.Fatalf("NewStore failed: %v", err)
			}

			store.Update(sourceA, sourceAFlags, nil)
			store.Update(sourceB, sourceBFlags, nil)
			store.Update(sourceC, sourceCFlags, nil)
			gotFlag, _, found := store.Get(context.Background(), tt.key, tt.selector)

			require.Equal(t, tt.wantFound, found, "expected found to be %v, got %v", tt.wantFound, found)

			if tt.wantFound {
				require.Equal(t, tt.wantFlag, gotFlag)
			}
		})
	}
}

func TestGetAllNoWatcher(t *testing.T) {

	sourceA := "sourceA"
	sourceB := "sourceB"
	sourceC := "sourceC"
	flagSetIdB := "flagSetIdA"
	flagSetIdC := "flagSetIdC"
	sources := []string{sourceA, sourceB, sourceC}

	sourceASelector := NewSelector("source=" + sourceA)
	flagSetIdCSelector := NewSelector("flagSetId=" + flagSetIdC)

	t.Parallel()
	tests := []struct {
		name      string
		selector  *Selector
		wantFlags map[string]model.Flag
	}{
		{
			name:     "nil selector",
			selector: nil,
			wantFlags: map[string]model.Flag{
				// "dupe" should be overwritten by higher priority flag
				"flagA": {Key: "flagA", DefaultVariant: "off", Source: sourceA, FlagSetId: nilFlagSetId, Priority: 0},
				"flagB": {Key: "flagB", DefaultVariant: "off", Source: sourceB, FlagSetId: flagSetIdB, Priority: 1, Metadata: model.Metadata{"flagSetId": flagSetIdB}},
				"flagC": {Key: "flagC", DefaultVariant: "off", Source: sourceC, FlagSetId: flagSetIdC, Priority: 2, Metadata: model.Metadata{"flagSetId": flagSetIdC}},
				"dupe":  {Key: "dupe", DefaultVariant: "off", Source: sourceC, FlagSetId: flagSetIdC, Priority: 2, Metadata: model.Metadata{"flagSetId": flagSetIdC}},
			},
		},
		{
			name:     "source selector",
			selector: &sourceASelector,
			wantFlags: map[string]model.Flag{
				// we should get the "dupe" from sourceA
				"flagA": {Key: "flagA", DefaultVariant: "off", Source: sourceA, FlagSetId: nilFlagSetId, Priority: 0},
				"dupe":  {Key: "dupe", DefaultVariant: "on", Source: sourceA, FlagSetId: nilFlagSetId, Priority: 0},
			},
		},
		{
			name:     "flagSetId selector",
			selector: &flagSetIdCSelector,
			wantFlags: map[string]model.Flag{
				// we should get the "dupe" from flagSetIdC
				"flagC": {Key: "flagC", DefaultVariant: "off", Source: sourceC, FlagSetId: flagSetIdC, Priority: 2, Metadata: model.Metadata{"flagSetId": flagSetIdC}},
				"dupe":  {Key: "dupe", DefaultVariant: "off", Source: sourceC, FlagSetId: flagSetIdC, Priority: 2, Metadata: model.Metadata{"flagSetId": flagSetIdC}},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sourceAFlags := map[string]model.Flag{
				"flagA": {Key: "flagA", DefaultVariant: "off"},
				"dupe":  {Key: "dupe", DefaultVariant: "on"},
			}
			sourceBFlags := map[string]model.Flag{
				"flagB": {Key: "flagB", DefaultVariant: "off", Metadata: model.Metadata{"flagSetId": flagSetIdB}},
			}
			sourceCFlags := map[string]model.Flag{
				"flagC": {Key: "flagC", DefaultVariant: "off", Metadata: model.Metadata{"flagSetId": flagSetIdC}},
				"dupe":  {Key: "dupe", DefaultVariant: "off", Metadata: model.Metadata{"flagSetId": flagSetIdC}},
			}

			store, err := NewStore(logger.NewLogger(nil, false), sources)
			if err != nil {
				t.Fatalf("NewStore failed: %v", err)
			}

			store.Update(sourceA, sourceAFlags, nil)
			store.Update(sourceB, sourceBFlags, nil)
			store.Update(sourceC, sourceCFlags, nil)
			gotFlags, _, _ := store.GetAll(context.Background(), tt.selector, nil)

			require.Equal(t, len(tt.wantFlags), len(gotFlags))
			require.Equal(t, tt.wantFlags, gotFlags)
		})
	}
}

func TestGetAllWithWatcher(t *testing.T) {

	sourceA := "sourceA"
	sourceB := "sourceB"
	sourceC := "sourceC"
	myFlagSetId := "myFlagSet"
	var sources = []string{sourceA, sourceB, sourceC}

	sourceASelector := NewSelector("source=" + sourceA)
	flagSetIdCSelector := NewSelector("flagSetId=" + myFlagSetId)
	emptySelector := NewSelector("")
	sourceCSelector := NewSelector("source=" + sourceC)

	tests := []struct {
		name        string
		selector    *Selector
		wantUpdates int
	}{
		{
			name:        "flag source selector (1 update)",
			selector:    &sourceASelector,
			wantUpdates: 1,
		},
		{
			name:        "flag set selector (3 updates)",
			selector:    &flagSetIdCSelector,
			wantUpdates: 3,
		},
		{
			name:        "no selector (all updates)",
			selector:    &emptySelector,
			wantUpdates: 4,
		},
		{
			name:        "flag source selector for unchanged source (no updates)",
			selector:    &sourceCSelector,
			wantUpdates: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sourceAFlags := map[string]model.Flag{
				"flagA": {Key: "flagA", DefaultVariant: "off"},
			}
			sourceBFlags := map[string]model.Flag{
				"flagB": {Key: "flagB", DefaultVariant: "off", Metadata: model.Metadata{"flagSetId": myFlagSetId}},
			}
			sourceCFlags := map[string]model.Flag{
				"flagC": {Key: "flagC", DefaultVariant: "off"},
			}

			store, err := NewStore(logger.NewLogger(nil, false), sources)
			if err != nil {
				t.Fatalf("NewStore failed: %v", err)
			}

			// setup initial flags
			store.Update(sourceA, sourceAFlags, model.Metadata{})
			store.Update(sourceB, sourceBFlags, model.Metadata{})
			store.Update(sourceC, sourceCFlags, model.Metadata{})
			watcher := make(chan Payload, 10)
			_, _, err = store.GetAll(context.Background(), tt.selector, watcher)
			if err != nil {
				t.Fatalf("GetAll failed: %v", err)
			}

			// perform updates
			go func() {

				// changing a flag default variant should trigger an update
				store.Update(sourceA, map[string]model.Flag{
					"flagA": {Key: "flagA", DefaultVariant: "on"},
				}, model.Metadata{})

				time.Sleep(10 * time.Millisecond)

				// changing a flag default variant should trigger an update
				store.Update(sourceB, map[string]model.Flag{
					"flagB": {Key: "flagB", DefaultVariant: "on", Metadata: model.Metadata{"flagSetId": myFlagSetId}},
				}, model.Metadata{})

				time.Sleep(10 * time.Millisecond)

				// removing a flag set id should trigger an update (even for flag set id selectors; it should remove the flag from the set)
				store.Update(sourceB, map[string]model.Flag{
					"flagB": {Key: "flagB", DefaultVariant: "on"},
				}, model.Metadata{})

				time.Sleep(10 * time.Millisecond)

				// adding a flag set id should trigger an update
				store.Update(sourceB, map[string]model.Flag{
					"flagB": {Key: "flagB", DefaultVariant: "on", Metadata: model.Metadata{"flagSetId": myFlagSetId}},
				}, model.Metadata{})
			}()

			updates := 0

			for true {
				select {
				case <-watcher:
					updates++
				case <-time.After(1000 * time.Millisecond):
					assert.Equal(t, tt.wantUpdates, updates, "expected %d updates, got %d", tt.wantUpdates, updates)
					return
				}
			}
		})
	}
}

func TestQueryMetadata(t *testing.T) {

	sourceA := "sourceA"
	otherSource := "otherSource"
	nonExistingFlagSetId := "nonExistingFlagSetId"
	var sources = []string{sourceA}
	sourceAFlags := map[string]model.Flag{
		"flagA": {Key: "flagA", DefaultVariant: "off"},
		"flagB": {Key: "flagB", DefaultVariant: "on"},
	}

	store, err := NewStore(logger.NewLogger(nil, false), sources)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}

	// setup initial flags
	store.Update(sourceA, sourceAFlags, model.Metadata{})

	selector := NewSelector("source=" + otherSource + ",flagSetId=" + nonExistingFlagSetId)
	_, metadata, _ := store.GetAll(context.Background(), &selector, nil)
	assert.Equal(t, metadata, model.Metadata{"source": otherSource, "flagSetId": nonExistingFlagSetId}, "metadata did not match expected")

	selector = NewSelector("source=" + otherSource + ",flagSetId=" + nonExistingFlagSetId)
	_, metadata, _ = store.Get(context.Background(), "key", &selector)
	assert.Equal(t, metadata, model.Metadata{"source": otherSource, "flagSetId": nonExistingFlagSetId}, "metadata did not match expected")
}
