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
		{
			name:      "non-empty indexMap, empty value",
			selector:  &Selector{indexMap: map[string]string{"flagSetId": ""}},
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

func TestSelector_WithSourceAndFlagSetId(t *testing.T) {
	s := Selector{}.WithSource("abc")
	if s.indexMap[sourceIndex] != "abc" {
		t.Errorf("WithSource did not set source")
	}

	s2 := s.WithFlagSetId("1234")
	if s2.indexMap[sourceIndex] != "abc" {
		t.Errorf("WithFlagSetId did not preserve source")
	}
	if s2.indexMap[flagSetIdIndex] != "1234" {
		t.Errorf("WithFlagSetId did not set flagSetId")
	}

	// Ensure original is unchanged
	if _, ok := s.indexMap[flagSetIdIndex]; ok {
		t.Errorf("WithFlagSetId mutated original selector")
	}
}

func TestSelector_ToQuery(t *testing.T) {
	tests := []struct {
		name       string
		selector   Selector
		wantIndex  string
		wantConstr []interface{}
	}{
		// #1708 Until we decide on the Selector syntax, only a single key=value pair is supported
		/*
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
		*/
		{
			name:       "single key",
			selector:   Selector{indexMap: map[string]string{"source": "src"}},
			wantIndex:  "source",
			wantConstr: []interface{}{"src"},
		},
		{
			name:       "flagSetId null",
			selector:   Selector{indexMap: map[string]string{"flagSetId": ""}},
			wantIndex:  "flagSetId",
			wantConstr: []interface{}{""},
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
		// #1708 Until we decide on the Selector syntax, only a single key=value pair is supported
		/*
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
		*/
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
		name    string
		input   string
		wantMap map[string]string
		wantErr bool
	}{
		// #1708 Until we decide on the Selector syntax, only a single key=value pair is supported
		/*
			{
				name:    "source and flagSetId",
				input:   "source=abc,flagSetId=1234",
				wantMap: map[string]string{"source": "abc", "flagSetId": "1234"},
			},
		*/
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
			name:    "null flagSetId",
			input:   "flagSetId=",
			wantMap: map[string]string{"flagSetId": nilFlagSetId},
		},
		{
			name:    "empty string",
			input:   "",
			wantMap: map[string]string{},
		},
		{
			name:    "invalid key",
			input:   "flagSetIds=abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewSelector(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSelector(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(s.indexMap, tt.wantMap) {
				t.Errorf("NewSelector(%q) indexMap = %v, want %v", tt.input, s.indexMap, tt.wantMap)
			}
		})
	}
}
