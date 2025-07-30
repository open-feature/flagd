package store

import (
	"context"
	"testing"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/stretchr/testify/require"
)

func TestMergeFlags(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		setup       func(t *testing.T) *Store
		new         map[string]model.Flag
		newSource   string
		newSelector string
		want        map[string]model.Flag
		wantNotifs  map[string]interface{}
		wantResync  bool
	}{
		{
			name: "both nil",
			setup: func(t *testing.T) *Store {
				s, err := NewStore(logger.NewLogger(nil, false))
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}
				return s
			},
			new:        nil,
			want:       map[string]model.Flag{},
			wantNotifs: map[string]interface{}{},
		},
		{
			name: "both empty flags",
			setup: func(t *testing.T) *Store {
				s, err := NewStore(logger.NewLogger(nil, false))
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}
				return s
			},
			new:        map[string]model.Flag{},
			want:       map[string]model.Flag{},
			wantNotifs: map[string]interface{}{},
		},
		{
			name: "empty new",
			setup: func(t *testing.T) *Store {
				s, err := NewStore(logger.NewLogger(nil, false))
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}
				return s
			},
			new:        nil,
			want:       map[string]model.Flag{},
			wantNotifs: map[string]interface{}{},
		},
		{
			name: "merging with new source",
			setup: func(t *testing.T) *Store {
				s, err := NewStore(logger.NewLogger(nil, false))
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}
				s.Update("1", "", map[string]model.Flag{
					"waka": {DefaultVariant: "off"},
				}, model.Metadata{})
				return s
			},
			new: map[string]model.Flag{
				"paka": {DefaultVariant: "on"},
			},
			newSource: "2",
			want: map[string]model.Flag{
				"waka": {Key: "waka", DefaultVariant: "off", Source: "1"},
				"paka": {Key: "paka", DefaultVariant: "on", Source: "2"},
			},
			wantNotifs: map[string]interface{}{"paka": map[string]interface{}{"type": "write", "source": "2"}},
		},
		{
			name: "override by new update",
			setup: func(t *testing.T) *Store {
				s, err := NewStore(logger.NewLogger(nil, false))
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}
				s.Update("", "", map[string]model.Flag{
					"waka": {DefaultVariant: "off"},
					"paka": {DefaultVariant: "off"},
				}, model.Metadata{})
				return s
			},
			new: map[string]model.Flag{
				"waka": {DefaultVariant: "on"},
				"paka": {DefaultVariant: "on"},
			},
			want: map[string]model.Flag{
				"waka": {Key: "waka", DefaultVariant: "on", Source: ""},
				"paka": {Key: "paka", DefaultVariant: "on", Source: ""},
			},
			wantNotifs: map[string]interface{}{
				"waka": map[string]interface{}{"type": "update", "source": ""},
				"paka": map[string]interface{}{"type": "update", "source": ""},
			},
		},
		{
			name: "identical update so empty notifications",
			setup: func(t *testing.T) *Store {
				s, err := NewStore(logger.NewLogger(nil, false))
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}
				s.Update("", "", map[string]model.Flag{
					"hello": {DefaultVariant: "off"},
				}, model.Metadata{})
				return s
			},
			new: map[string]model.Flag{
				"hello": {DefaultVariant: "off"},
			},
			want: map[string]model.Flag{
				"hello": {Key: "hello", DefaultVariant: "off", Source: ""},
			},
			wantNotifs: map[string]interface{}{},
		},
		{
			name: "deleted flag & trigger resync for same source",
			setup: func(t *testing.T) *Store {
				s, err := NewStore(logger.NewLogger(nil, false))
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}
				s.Update("A", "", map[string]model.Flag{
					"hello": {DefaultVariant: "off"},
				}, model.Metadata{})
				return s
			},
			new:       map[string]model.Flag{},
			newSource: "A",
			want:      map[string]model.Flag{},
			wantNotifs: map[string]interface{}{
				"hello": map[string]interface{}{"type": "delete", "source": "A"},
			},
			wantResync: true,
		},
		{
			name: "no deleted & no resync for same source but different selector",
			setup: func(t *testing.T) *Store {
				s, err := NewStore(logger.NewLogger(nil, false))
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}
				s.Update("A", "X", map[string]model.Flag{
					"hello": {DefaultVariant: "off"},
				}, model.Metadata{})
				return s
			},
			new:         map[string]model.Flag{},
			newSource:   "A",
			newSelector: "Y",
			want: map[string]model.Flag{
				"hello": {Key: "hello", DefaultVariant: "off", Source: "A", Selector: "X"},
			},
			wantResync: false,
			wantNotifs: map[string]interface{}{},
		},
		{
			name: "no merge due to low priority",
			setup: func(t *testing.T) *Store {
				s, err := NewStore(logger.NewLogger(nil, false))
				if err != nil {
					t.Fatalf("NewStore failed: %v", err)
				}
				s.FlagSources = []string{"B", "A"}
				s.Update("A", "", map[string]model.Flag{
					"hello": {DefaultVariant: "off"},
				}, model.Metadata{})
				return s
			},
			new:       map[string]model.Flag{"hello": {DefaultVariant: "off"}},
			newSource: "B",
			want: map[string]model.Flag{
				"hello": {Key: "hello", DefaultVariant: "off", Source: "A"},
			},
			wantNotifs: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			state := tt.setup(t)
			gotNotifs, resyncRequired := state.Update(tt.newSource, tt.newSelector, tt.new, model.Metadata{})
			gotFlags, _, _ := state.GetAll(context.Background())

			require.Equal(t, tt.want, gotFlags)
			require.Equal(t, tt.wantNotifs, gotNotifs)
			require.Equal(t, tt.wantResync, resyncRequired)
		})
	}
}
