package store

import (
	"reflect"
	"testing"

	"github.com/open-feature/flagd/core/pkg/model"
)

func TestSelector_IsEmpty(t *testing.T) {
	tests := []struct {
		name      string
		selector  *Selector
		wantEmpty bool
	}{
		{
			name:      "nil selector",
			selector:  nil,
			wantEmpty: true,
		},
		{
			name:      "nil indexMap",
			selector:  &Selector{indexMap: nil},
			wantEmpty: true,
		},
		{
			name:      "empty indexMap",
			selector:  &Selector{indexMap: map[string]string{}},
			wantEmpty: true,
		},
		{
			name:      "non-empty indexMap",
			selector:  &Selector{indexMap: map[string]string{"source": "abc"}},
			wantEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.selector.IsEmpty()
			if got != tt.wantEmpty {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.wantEmpty)
			}
		})
	}
}

func TestSelector_WithIndex(t *testing.T) {
	oldS := Selector{indexMap: map[string]string{"source": "abc"}}
	newS := oldS.WithIndex("flagSetId", "1234")

	if newS.indexMap["source"] != "abc" {
		t.Errorf("WithIndex did not preserve existing keys")
	}
	if newS.indexMap["flagSetId"] != "1234" {
		t.Errorf("WithIndex did not add new key")
	}
	// Ensure original is unchanged
	if _, ok := oldS.indexMap["flagSetId"]; ok {
		t.Errorf("WithIndex mutated original selector")
	}
}

func TestSelector_ToQuery(t *testing.T) {
	tests := []struct {
		name       string
		selector   Selector
		wantIndex  string
		wantConstr []interface{}
	}{
		{
			name:       "flagSetId and key primary index special case",
			selector:   Selector{indexMap: map[string]string{"flagSetId": "fsid", "key": "myKey"}},
			wantIndex:  "id",
			wantConstr: []interface{}{"fsid", "myKey"},
		},
		{
			name:       "multiple keys sorted",
			selector:   Selector{indexMap: map[string]string{"source": "src", "flagSetId": "fsid"}},
			wantIndex:  "flagSetId+source",
			wantConstr: []interface{}{"fsid", "src"},
		},
		{
			name:       "single key",
			selector:   Selector{indexMap: map[string]string{"source": "src"}},
			wantIndex:  "source",
			wantConstr: []interface{}{"src"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIndex, gotConstr := tt.selector.ToQuery()
			if gotIndex != tt.wantIndex {
				t.Errorf("ToQuery() index = %v, want %v", gotIndex, tt.wantIndex)
			}
			if !reflect.DeepEqual(gotConstr, tt.wantConstr) {
				t.Errorf("ToQuery() constraints = %v, want %v", gotConstr, tt.wantConstr)
			}
		})
	}
}

func TestSelector_ToMetadata(t *testing.T) {
	tests := []struct {
		name     string
		selector *Selector
		want     model.Metadata
	}{
		{
			name:     "nil selector",
			selector: nil,
			want:     model.Metadata{},
		},
		{
			name:     "nil indexMap",
			selector: &Selector{indexMap: nil},
			want:     model.Metadata{},
		},
		{
			name:     "empty indexMap",
			selector: &Selector{indexMap: map[string]string{}},
			want:     model.Metadata{},
		},
		{
			name:     "flagSetId only",
			selector: &Selector{indexMap: map[string]string{"flagSetId": "fsid"}},
			want:     model.Metadata{"flagSetId": "fsid"},
		},
		{
			name:     "source only",
			selector: &Selector{indexMap: map[string]string{"source": "src"}},
			want:     model.Metadata{"source": "src"},
		},
		{
			name:     "flagSetId and source",
			selector: &Selector{indexMap: map[string]string{"flagSetId": "fsid", "source": "src"}},
			want:     model.Metadata{"flagSetId": "fsid", "source": "src"},
		},
		{
			name:     "flagSetId, source, and key (key should be ignored)",
			selector: &Selector{indexMap: map[string]string{"flagSetId": "fsid", "source": "src", "key": "myKey"}},
			want:     model.Metadata{"flagSetId": "fsid", "source": "src"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.selector.ToMetadata()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSelector(t *testing.T) {
	tests := []struct {
		name                  string
		input                 string
		fallbackExpressionKey string
		wantMap               map[string]string
	}{
		{
			name:    "source and flagSetId",
			input:   "source=abc,flagSetId=1234",
			wantMap: map[string]string{"source": "abc", "flagSetId": "1234"},
		},
		{
			name:    "source",
			input:   "source=abc",
			wantMap: map[string]string{"source": "abc"},
		},
		{
			name:    "no equals, treat as source",
			input:   "mysource",
			wantMap: map[string]string{"source": "mysource"},
		},
		{
			name:                  "no equals, treat as flagSetId",
			input:                 "flagSetId",
			fallbackExpressionKey: "flagSetId",
			wantMap:               map[string]string{"flagSetId": "flagSetId"},
		},
		{
			name:    "empty string",
			input:   "",
			wantMap: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSelectorWithFallback(tt.input, tt.fallbackExpressionKey)
			if !reflect.DeepEqual(s.indexMap, tt.wantMap) {
				t.Errorf("NewSelector(%q) indexMap = %v, want %v", tt.input, s.indexMap, tt.wantMap)
			}
		})
	}
}
