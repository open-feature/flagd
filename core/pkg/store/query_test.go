package store

import (
	"reflect"
	"testing"
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
	s := Selector{indexMap: map[string]string{"flagSetId": "fsid", "source": "src", "key": "myKey"}}
	meta := s.ToMetadata()
	if meta["flagSetId"] != "fsid" {
		t.Errorf("ToMetadata missing flagSetId")
	}
	if meta["source"] != "src" {
		t.Errorf("ToMetadata missing source")
	}
	if _, ok := meta["key"]; ok {
		t.Errorf("ToMetadata should not include key")
	}
}

func TestNewSelector(t *testing.T) {
	s := NewSelector("source=abc,flagSetId=1234")
	if s.indexMap["source"] != "abc" {
		t.Errorf("NewSelector did not parse source")
	}
	if s.indexMap["flagSetId"] != "1234" {
		t.Errorf("NewSelector did not parse flagSetId")
	}
}
