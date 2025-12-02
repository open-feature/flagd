package store

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type source struct {
	Name  string
	flags []model.Flag
}

func TestUpdateFlags(t *testing.T) {

	const source1 = "source1"
	const source2 = "source2"
	var sources = []string{source1, source2}

	t.Parallel()
	tests := []struct {
		name        string
		setup       func(t *testing.T) IStore
		newFlags    []model.Flag
		source      string
		wantFlags   []model.Flag
		setMetadata model.Metadata
	}{
		{
			name: "both nil",
			setup: func(t *testing.T) IStore {
				s, err := NewStore(logger.NewLogger(nil, false), sources)
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}
				return s
			},
			source:    source1,
			newFlags:  nil,
			wantFlags: []model.Flag{},
		},
		{
			name: "both empty flags",
			setup: func(t *testing.T) IStore {
				s, err := NewStore(logger.NewLogger(nil, false), sources)
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}
				return s
			},
			source:    source1,
			newFlags:  []model.Flag{},
			wantFlags: []model.Flag{},
		},
		{
			name: "empty new",
			setup: func(t *testing.T) IStore {
				s, err := NewStore(logger.NewLogger(nil, false), sources)
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}
				return s
			},
			source:    source1,
			newFlags:  nil,
			wantFlags: []model.Flag{},
		},
		{
			name: "update from source 1 (old flag removed)",
			setup: func(t *testing.T) IStore {
				s, err := NewStore(logger.NewLogger(nil, false), sources)
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}
				s.Update(source1, []model.Flag{
					{Key: "waka", DefaultVariant: "off"},
				}, nil)
				return s
			},
			newFlags: []model.Flag{
				{Key: "paka", DefaultVariant: "on"},
			},
			source: source1,
			wantFlags: []model.Flag{
				{Key: "paka", DefaultVariant: "on", Source: source1, FlagSetId: nilFlagSetId, Priority: 0},
			},
		},
		{
			name: "update from source 1 (new flag added)",
			setup: func(t *testing.T) IStore {
				s, err := NewStore(logger.NewLogger(nil, false), sources)
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}
				s.Update(source1, []model.Flag{
					{Key: "waka", DefaultVariant: "off"},
				}, nil)
				return s
			},
			newFlags: []model.Flag{
				{Key: "paka", DefaultVariant: "on"},
			},
			source: source2,
			wantFlags: []model.Flag{
				{Key: "waka", DefaultVariant: "off", Source: source1, FlagSetId: nilFlagSetId, Priority: 0},
				{Key: "paka", DefaultVariant: "on", Source: source2, FlagSetId: nilFlagSetId, Priority: 1},
			},
		},
		{
			name: "flag set inheritance",
			setup: func(t *testing.T) IStore {
				s, err := NewStore(logger.NewLogger(nil, false), sources)
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}
				s.Update(source1, []model.Flag{}, model.Metadata{})
				return s
			},
			setMetadata: model.Metadata{
				"flagSetId": "topLevelSet", // top level set metadata, including flagSetId
			},
			newFlags: []model.Flag{
				{Key: "waka", DefaultVariant: "on"},
				{Key: "paka", DefaultVariant: "on", Metadata: model.Metadata{"flagSetId": "flagLevelSet"}}, // overrides set level flagSetId
			},
			source: source1,
			wantFlags: []model.Flag{
				{Key: "waka", DefaultVariant: "on", Source: source1, FlagSetId: "topLevelSet", Priority: 0, Metadata: model.Metadata{"flagSetId": "topLevelSet"}},
				{Key: "paka", DefaultVariant: "on", Source: source1, FlagSetId: "flagLevelSet", Priority: 0, Metadata: model.Metadata{"flagSetId": "flagLevelSet"}},
			},
		},
		{
			name: "flag set same for different sets",
			setup: func(t *testing.T) IStore {
				s, err := NewStore(logger.NewLogger(nil, false), sources)
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}
				s.Update(source1, []model.Flag{}, model.Metadata{})
				return s

			},
			setMetadata: model.Metadata{},
			newFlags: []model.Flag{
				{Key: "paka", DefaultVariant: "on"},
				{Key: "paka", DefaultVariant: "on", Metadata: model.Metadata{"flagSetId": "flagLevelSet1"}}, // overrides set level flagSetId
				{Key: "paka", DefaultVariant: "on", Metadata: model.Metadata{"flagSetId": "flagLevelSet2"}}, // overrides set level flagSetId
				{Key: "paka", DefaultVariant: "on", Metadata: model.Metadata{"flagSetId": "flagLevelSet3"}}, // overrides set level flagSetId
			},
			source: source1,
			wantFlags: []model.Flag{
				{Key: "paka", DefaultVariant: "on", Source: source1, FlagSetId: "flagLevelSet3", Priority: 0, Metadata: model.Metadata{"flagSetId": "flagLevelSet3"}},
				{Key: "paka", DefaultVariant: "on", Source: source1, FlagSetId: "flagLevelSet2", Priority: 0, Metadata: model.Metadata{"flagSetId": "flagLevelSet2"}},
				{Key: "paka", DefaultVariant: "on", Source: source1, FlagSetId: "flagLevelSet1", Priority: 0, Metadata: model.Metadata{"flagSetId": "flagLevelSet1"}},
				{Key: "paka", DefaultVariant: "on", Source: source1, FlagSetId: nilFlagSetId, Priority: 0, Metadata: model.Metadata{}},
			},
		},
		{
			name: "flag set same for different sets - toplevelflagset",
			setup: func(t *testing.T) IStore {
				s, err := NewStore(logger.NewLogger(nil, false), sources)
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}
				s.Update(source1, []model.Flag{}, model.Metadata{})
				return s
			},
			setMetadata: model.Metadata{
				"flagSetId": "topLevelSet", // top level set metadata, including flagSetId
			},
			newFlags: []model.Flag{
				{Key: "paka", DefaultVariant: "on"},
				{Key: "paka", DefaultVariant: "on", Metadata: model.Metadata{"flagSetId": "flagLevelSet1"}}, // overrides set level flagSetId
				{Key: "paka", DefaultVariant: "on", Metadata: model.Metadata{"flagSetId": "flagLevelSet2"}}, // overrides set level flagSetId
				{Key: "paka", DefaultVariant: "on", Metadata: model.Metadata{"flagSetId": "flagLevelSet3"}}, // overrides set level flagSetId
			},
			source: source1,
			wantFlags: []model.Flag{
				{Key: "paka", DefaultVariant: "on", Source: source1, FlagSetId: "topLevelSet", Priority: 0, Metadata: model.Metadata{"flagSetId": "topLevelSet"}},
				{Key: "paka", DefaultVariant: "on", Source: source1, FlagSetId: "flagLevelSet3", Priority: 0, Metadata: model.Metadata{"flagSetId": "flagLevelSet3"}},
				{Key: "paka", DefaultVariant: "on", Source: source1, FlagSetId: "flagLevelSet2", Priority: 0, Metadata: model.Metadata{"flagSetId": "flagLevelSet2"}},
				{Key: "paka", DefaultVariant: "on", Source: source1, FlagSetId: "flagLevelSet1", Priority: 0, Metadata: model.Metadata{"flagSetId": "flagLevelSet1"}},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			store := tt.setup(t)
			store.Update(tt.source, tt.newFlags, tt.setMetadata)
			gotFlags, _, _ := store.GetAll(context.Background(), nil)
			sort.Slice(tt.wantFlags, func(i, j int) bool {
				return tt.wantFlags[i].FlagSetId+"|"+tt.wantFlags[i].Key > tt.wantFlags[j].FlagSetId+"|"+tt.wantFlags[j].Key
			})
			sort.Slice(gotFlags, func(i, j int) bool {
				return gotFlags[i].FlagSetId+"|"+gotFlags[i].Key > gotFlags[j].FlagSetId+"|"+gotFlags[j].Key
			})
			require.EqualValues(t, tt.wantFlags, gotFlags)
		})
	}
}

func TestGet(t *testing.T) {

	flagSetIdB := "flagSetIdA"
	flagSetIdC := "flagSetIdC"

	sourceA := source{
		Name: "sourceA",
		flags: []model.Flag{
			{Key: "flagA", DefaultVariant: "off"},
			{Key: "dupe", DefaultVariant: "on"},
		},
	}
	sourceB := source{
		Name: "sourceB",
		flags: []model.Flag{
			{Key: "flagB", DefaultVariant: "off", Metadata: model.Metadata{"flagSetId": flagSetIdB}},
			{Key: "dupeMultiSource", DefaultVariant: "on", Metadata: model.Metadata{"flagSetId": flagSetIdC}},
		},
	}
	sourceC := source{
		Name: "sourceC",
		flags: []model.Flag{
			{Key: "flagC", DefaultVariant: "off", Metadata: model.Metadata{"flagSetId": flagSetIdC}},
			{Key: "dupe", DefaultVariant: "off", Metadata: model.Metadata{"flagSetId": flagSetIdC}},
			{Key: "dupeMultiSource", DefaultVariant: "off", Metadata: model.Metadata{"flagSetId": flagSetIdC}},
		},
	}

	sources := []string{sourceA.Name, sourceB.Name, sourceC.Name}

	sourceASelector := NewSelector("source=" + sourceA.Name)
	flagSetIdCSelector := NewSelector("flagSetId=" + flagSetIdC)

	t.Parallel()
	tests := []struct {
		name     string
		key      string
		selector *Selector
		wantFlag model.Flag
		wantErr  bool
	}{
		{
			name:     "nil selector",
			key:      "flagA",
			selector: nil,
			wantFlag: model.Flag{Key: "flagA", DefaultVariant: "off", Source: sourceA.Name, FlagSetId: nilFlagSetId, Priority: 0},
			wantErr:  false,
		},
		{
			name:     "flagSetId selector",
			key:      "dupe",
			selector: &flagSetIdCSelector,
			wantFlag: model.Flag{Key: "dupe", DefaultVariant: "off", Source: sourceC.Name, FlagSetId: flagSetIdC, Priority: 2, Metadata: model.Metadata{"flagSetId": flagSetIdC}},
			wantErr:  false,
		},
		{
			name:     "flagSetId selector - MultiSource",
			key:      "dupeMultiSource",
			selector: &flagSetIdCSelector,
			wantFlag: model.Flag{Key: "dupeMultiSource", FlagSetId: "flagSetIdC", Priority: 2, State: "", DefaultVariant: "off", Source: sourceC.Name, Metadata: map[string]interface{}{"flagSetId": "flagSetIdC"}},
			wantErr:  false,
		},
		{
			name:     "source selector",
			key:      "dupe",
			selector: &sourceASelector,
			wantFlag: model.Flag{Key: "dupe", DefaultVariant: "on", Source: sourceA.Name, FlagSetId: nilFlagSetId, Priority: 0},
			wantErr:  false,
		},
		{
			name:     "flag not found with source selector",
			key:      "flagB",
			selector: &sourceASelector,
			wantFlag: model.Flag{Key: "flagB", DefaultVariant: "off", Source: sourceB.Name, FlagSetId: flagSetIdB, Priority: 1, Metadata: model.Metadata{"flagSetId": flagSetIdB}},
			wantErr:  true,
		},
		{
			name:     "flag not found with flagSetId selector",
			key:      "flagB",
			selector: &flagSetIdCSelector,
			wantFlag: model.Flag{Key: "flagB", DefaultVariant: "off", Source: sourceB.Name, FlagSetId: flagSetIdB, Priority: 1, Metadata: model.Metadata{"flagSetId": flagSetIdB}},
			wantErr:  true,
		},
	}

	sourceOrder := []struct {
		name  string
		order []source
	}{
		{
			name:  "normal",
			order: []source{sourceA, sourceB, sourceC},
		},
		{
			name:  "inverted",
			order: []source{sourceC, sourceB, sourceA},
		},
		{
			name:  "random1",
			order: []source{sourceB, sourceA, sourceC},
		},
		{
			name:  "random2",
			order: []source{sourceB, sourceC, sourceA},
		},
		{
			name:  "normal, loading sourceA twice",
			order: []source{sourceA, sourceB, sourceC, sourceA},
		},
	}
	for _, tt := range tests {
		for _, s := range sourceOrder {
			tt := tt
			t.Run(tt.name+" - "+s.name, func(t *testing.T) {
				t.Parallel()

				store, err := NewStore(logger.NewLogger(nil, false), sources)
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}

				for _, source := range s.order {
					store.Update(source.Name, source.flags, nil)
				}
				gotFlag, _, err := store.Get(context.Background(), tt.key, tt.selector)

				if !tt.wantErr {
					require.Equal(t, tt.wantFlag, gotFlag)
				} else {
					require.Error(t, err, "expected an error for key %s with selector %v", tt.key, tt.selector)
				}
			})
		}
	}
}

func TestGetAllNoWatcher(t *testing.T) {

	flagSetIdC := "flagSetIdC"

	sourceA := source{
		Name: "sourceA",
		flags: []model.Flag{
			{Key: "flagA", DefaultVariant: "off"},
			{Key: "dupe", DefaultVariant: "on"},
			{Key: "dupe", DefaultVariant: "on", Metadata: model.Metadata{"flagSetId": flagSetIdC}},
			{Key: "dupeSingleSource", DefaultVariant: "on", Metadata: model.Metadata{"flagSetId": flagSetIdC}},
			{Key: "dupeSingleSource", DefaultVariant: "off", Metadata: model.Metadata{"flagSetId": flagSetIdC}},
			{Key: "dupeMultiSource", DefaultVariant: "on", Metadata: model.Metadata{"flagSetId": flagSetIdC}},
			{Key: "dupeMultiSource", DefaultVariant: "off", Metadata: model.Metadata{"flagSetId": flagSetIdC}},
		},
	}
	sourceB := source{
		Name: "sourceB",
		flags: []model.Flag{
			{Key: "flagB", DefaultVariant: "off"},
		},
	}
	sourceC := source{
		Name: "sourceC",
		flags: []model.Flag{
			{Key: "flagC", DefaultVariant: "off"},
			{Key: "dupe", DefaultVariant: "off", Metadata: model.Metadata{"flagSetId": flagSetIdC}},
			{Key: "dupeMultiSource", DefaultVariant: "both", Metadata: model.Metadata{"flagSetId": flagSetIdC}},
		},
	}

	sources := []string{sourceA.Name, sourceB.Name, sourceC.Name}

	sourceASelector := NewSelector("source=" + sourceA.Name)
	flagSetIdCSelector := NewSelector("flagSetId=" + flagSetIdC)
	// #1708 Until we decide on the Selector syntax, only a single key=value pair is supported
	//flagSetIdAndCSelector := NewSelector("flagSetId=" + flagSetIdC + ",source=" + sourceC.Name)

	t.Parallel()
	tests := []struct {
		name      string
		selector  *Selector
		wantFlags []model.Flag
	}{
		{
			name:     "nil selector",
			selector: nil,
			wantFlags: []model.Flag{
				// "dupe" should be overwritten by higher priority flag
				{Key: "flagC", DefaultVariant: "off", Source: sourceC.Name, FlagSetId: nilFlagSetId, Priority: 2},
				{Key: "dupeSingleSource", DefaultVariant: "off", Source: sourceA.Name, FlagSetId: flagSetIdC, Metadata: model.Metadata{"flagSetId": flagSetIdC}, Priority: 0},
				{Key: "dupeMultiSource", DefaultVariant: "both", Source: sourceC.Name, FlagSetId: flagSetIdC, Metadata: model.Metadata{"flagSetId": flagSetIdC}, Priority: 2},
				{Key: "dupe", DefaultVariant: "off", Source: sourceC.Name, FlagSetId: flagSetIdC, Priority: 2, Metadata: model.Metadata{"flagSetId": flagSetIdC}},
				{Key: "flagB", DefaultVariant: "off", Source: sourceB.Name, FlagSetId: nilFlagSetId, Priority: 1},
				{Key: "flagA", DefaultVariant: "off", Source: sourceA.Name, FlagSetId: nilFlagSetId, Priority: 0},
				{Key: "dupe", DefaultVariant: "on", Source: sourceA.Name, FlagSetId: nilFlagSetId, Priority: 0},
			},
		},
		{
			name:     "source selector",
			selector: &sourceASelector,
			wantFlags: []model.Flag{
				// we should get the "dupe" from sourceAName
				{Key: "dupe", FlagSetId: nilFlagSetId, Priority: 0, State: "", DefaultVariant: "on", Source: sourceA.Name},
				{Key: "flagA", FlagSetId: nilFlagSetId, Priority: 0, State: "", DefaultVariant: "off", Source: sourceA.Name},
				{Key: "dupeSingleSource", FlagSetId: "flagSetIdC", Priority: 0, State: "", DefaultVariant: "off", Source: sourceA.Name, Metadata: map[string]interface{}{"flagSetId": "flagSetIdC"}}},
		},
		{
			name:     "flagSetId selector",
			selector: &flagSetIdCSelector,
			wantFlags: []model.Flag{
				// we should get the "dupe" from flagSetIdC
				{Key: "dupeSingleSource", DefaultVariant: "off", Source: sourceA.Name, FlagSetId: flagSetIdC, Metadata: model.Metadata{"flagSetId": flagSetIdC}, Priority: 0},
				{Key: "dupeMultiSource", DefaultVariant: "both", Source: sourceC.Name, FlagSetId: flagSetIdC, Metadata: model.Metadata{"flagSetId": flagSetIdC}, Priority: 2},
				{Key: "dupe", DefaultVariant: "off", Source: sourceC.Name, FlagSetId: flagSetIdC, Priority: 2, Metadata: model.Metadata{"flagSetId": flagSetIdC}},
			},
		},
		// #1708 Until we decide on the Selector syntax, only a single key=value pair is supported
		/*
			{
				name:     "flagSetId and source selector",
				selector: &flagSetIdAndCSelector,
				wantFlags: []model.Flag{
					{Key: "dupeMultiSource", DefaultVariant: "both", Source: sourceC.Name, FlagSetId: flagSetIdC, Metadata: model.Metadata{"flagSetId": flagSetIdC}, Priority: 2},
					{Key: "dupe", DefaultVariant: "off", Source: sourceC.Name, FlagSetId: flagSetIdC, Priority: 2, Metadata: model.Metadata{"flagSetId": flagSetIdC}},
				},
			},
		*/
	}

	sourceOrder := []struct {
		name  string
		order []source
	}{
		{
			name:  "normal",
			order: []source{sourceA, sourceB, sourceC},
		},
		{
			name:  "inverted",
			order: []source{sourceC, sourceB, sourceA},
		},
		{
			name:  "random1",
			order: []source{sourceB, sourceA, sourceC},
		},
		{
			name:  "random2",
			order: []source{sourceB, sourceC, sourceA},
		},
		{
			name:  "normal, loading sourceA twice",
			order: []source{sourceA, sourceB, sourceC, sourceA},
		},
	}

	for _, tt := range tests {
		for _, s := range sourceOrder {
			wantFlags := make([]model.Flag, len(tt.wantFlags))
			copy(wantFlags, tt.wantFlags)
			t.Run(tt.name+" - "+s.name, func(t *testing.T) {
				t.Parallel()

				store, err := NewStore(logger.NewLogger(nil, false), sources)
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}

				for _, source := range s.order {
					store.Update(source.Name, source.flags, nil)
				}
				gotFlags, _, _ := store.GetAll(context.Background(), tt.selector)

				require.Equal(t, len(wantFlags), len(gotFlags))
				sort.Slice(wantFlags, func(i, j int) bool {
					if wantFlags[i].FlagSetId != wantFlags[j].FlagSetId {
						return wantFlags[i].FlagSetId < wantFlags[j].FlagSetId
					}
					return wantFlags[i].Key < wantFlags[j].Key
				})
				sort.Slice(gotFlags, func(i, j int) bool {
					if gotFlags[i].FlagSetId != gotFlags[j].FlagSetId {
						return gotFlags[i].FlagSetId < gotFlags[j].FlagSetId
					}
					return gotFlags[i].Key < gotFlags[j].Key
				})
				require.Equal(t, wantFlags, gotFlags)
			})
		}
	}
}

func TestWatch(t *testing.T) {

	sourceA := "sourceA"
	sourceB := "sourceB"
	sourceC := "sourceC"
	myFlagSetId := "myFlagSet"
	var sources = []string{sourceA, sourceB, sourceC}
	pauseTime := 100 * time.Millisecond // time for updates to settle
	timeout := 1000 * time.Millisecond  // time to make sure we get enough updates, and no extras

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
			name:        "flag source selector (initial, plus 1 update)",
			selector:    &sourceASelector,
			wantUpdates: 2,
		},
		{
			name:        "flag set selector (initial, plus 3 updates)",
			selector:    &flagSetIdCSelector,
			wantUpdates: 4,
		},
		{
			name:        "no selector (all updates)",
			selector:    &emptySelector,
			wantUpdates: 5,
		},
		{
			name:        "flag source selector for unchanged source (initial, plus no updates)",
			selector:    &sourceCSelector,
			wantUpdates: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sourceAFlags := []model.Flag{
				{Key: "flagA", DefaultVariant: "off"},
			}
			sourceBFlags := []model.Flag{{Key: "flagB", DefaultVariant: "off", Metadata: model.Metadata{"flagSetId": myFlagSetId}}}
			sourceCFlags := []model.Flag{
				{Key: "flagC", DefaultVariant: "off"},
			}

			store, err := NewStore(logger.NewLogger(nil, false), sources)
			if err != nil {
				t.Fatalf("NewStore failed: %v", err)
			}

			// setup initial flags
			store.Update(sourceA, sourceAFlags, model.Metadata{})
			store.Update(sourceB, sourceBFlags, model.Metadata{})
			store.Update(sourceC, sourceCFlags, model.Metadata{})
			watcher := make(chan FlagQueryResult, 1)
			time.Sleep(pauseTime)

			ctx, cancel := context.WithCancel(context.Background())
			store.Watch(ctx, tt.selector, watcher)

			// perform updates
			go func() {

				time.Sleep(pauseTime)

				// changing a flag default variant should trigger an update
				store.Update(sourceA, []model.Flag{
					{Key: "flagA", DefaultVariant: "on"},
				}, model.Metadata{})

				time.Sleep(pauseTime)

				// changing a flag default variant should trigger an update
				store.Update(sourceB, []model.Flag{
					{Key: "flagB", DefaultVariant: "on", Metadata: model.Metadata{"flagSetId": myFlagSetId}},
				}, model.Metadata{})

				time.Sleep(pauseTime)

				// removing a flag set id should trigger an update (even for flag set id selectors; it should remove the flag from the set)
				// TODO: challenge this test and behaviour
				store.Update(sourceB, []model.Flag{
					{Key: "flagB", DefaultVariant: "on"},
				}, model.Metadata{})

				time.Sleep(pauseTime)

				// adding a flag set id should trigger an update
				store.Update(sourceB, []model.Flag{
					{Key: "flagB", DefaultVariant: "on", Metadata: model.Metadata{"flagSetId": myFlagSetId}},
				}, model.Metadata{})
			}()

			updates := 0

			for {
				select {
				case <-time.After(timeout):
					assert.Equal(t, tt.wantUpdates, updates, "expected %d updates, got %d", tt.wantUpdates, updates)
					cancel()
					_, open := <-watcher
					assert.False(t, open, "watcher channel should be closed after cancel")
					return
				case q := <-watcher:
					if q.Flags != nil {
						updates++
					}
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
	sourceAFlags := []model.Flag{
		{Key: "flagA", DefaultVariant: "off"},
		{Key: "flagB", DefaultVariant: "on"},
	}

	store, err := NewStore(logger.NewLogger(nil, false), sources)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}

	// setup initial flags
	store.Update(sourceA, sourceAFlags, model.Metadata{})

	// #1708 Until we decide on the Selector syntax, only a single key=value pair is supported
	// 		 these tests should then also cover more complex selectors

	selector := NewSelector("flagSetId=" + nonExistingFlagSetId)
	_, metadata, _ := store.GetAll(context.Background(), &selector)
	assert.Equal(t, metadata, model.Metadata{"flagSetId": nonExistingFlagSetId}, "metadata did not match expected")

	selector = NewSelector("flagSetId=" + nonExistingFlagSetId)
	_, metadata, _ = store.Get(context.Background(), "key", &selector)
	assert.Equal(t, metadata, model.Metadata{"flagSetId": nonExistingFlagSetId}, "metadata did not match expected")

	selector = NewSelector("source=" + otherSource)
	_, metadata, _ = store.Get(context.Background(), "key", &selector)
	assert.Equal(t, metadata, model.Metadata{"source": otherSource}, "metadata did not match expected")
}
