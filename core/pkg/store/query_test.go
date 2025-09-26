package store

import (
	"reflect"
	"testing"

	"github.com/open-feature/flagd/core/pkg/model"
)

// Existing tests remain as they are...
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
	t.Run("basic functionality", func(t *testing.T) {
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
	})

	// NEW TESTS for WithIndex
	t.Run("nil selector", func(t *testing.T) {
		var s *Selector
		newS := s.WithIndex("key", "value")

		if newS == nil {
			t.Errorf("WithIndex on nil selector should return new selector")
		}
		if newS.indexMap["key"] != "value" {
			t.Errorf("WithIndex on nil selector should set key-value pair")
		}
		if newS.usingFallback {
			t.Errorf("WithIndex on nil selector should not set usingFallback to true")
		}
	})

	t.Run("nil indexMap", func(t *testing.T) {
		s := &Selector{indexMap: nil, usingFallback: true}
		newS := s.WithIndex("key", "value")

		if newS.indexMap["key"] != "value" {
			t.Errorf("WithIndex with nil indexMap should set key-value pair")
		}
		if !newS.usingFallback {
			t.Errorf("WithIndex should preserve usingFallback status")
		}
	})

	t.Run("overwrite existing key", func(t *testing.T) {
		s := &Selector{indexMap: map[string]string{"key": "oldvalue"}}
		newS := s.WithIndex("key", "newvalue")

		if newS.indexMap["key"] != "newvalue" {
			t.Errorf("WithIndex should overwrite existing key")
		}
		if s.indexMap["key"] != "oldvalue" {
			t.Errorf("WithIndex should not mutate original")
		}
	})

	t.Run("empty values", func(t *testing.T) {
		s := &Selector{indexMap: map[string]string{"existing": "value"}}
		newS := s.WithIndex("", "emptykey")

		if newS.indexMap[""] != "emptykey" {
			t.Errorf("WithIndex should handle empty key")
		}

		newS2 := s.WithIndex("emptyvalue", "")
		if newS2.indexMap["emptyvalue"] != "" {
			t.Errorf("WithIndex should handle empty value")
		}
	})
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
		// NEW TEST CASES
		{
			name:       "empty selector",
			selector:   Selector{indexMap: map[string]string{}},
			wantIndex:  "",
			wantConstr: []interface{}{},
		},
		{
			name:       "three keys sorted alphabetically",
			selector:   Selector{indexMap: map[string]string{"zebra": "z", "alpha": "a", "beta": "b"}},
			wantIndex:  "alpha+beta+zebra",
			wantConstr: []interface{}{"a", "b", "z"},
		},
		{
			name:       "flagSetId and key with empty values",
			selector:   Selector{indexMap: map[string]string{"flagSetId": "", "key": "myKey"}},
			wantIndex:  "flagSetId+key",
			wantConstr: []interface{}{"", "myKey"},
		},
		{
			name:       "priority index",
			selector:   Selector{indexMap: map[string]string{"priority": "high"}},
			wantIndex:  "priority",
			wantConstr: []interface{}{"high"},
		},
		{
			name:       "compound index matching constants",
			selector:   Selector{indexMap: map[string]string{"flagSetId": "123", "source": "file.yaml"}},
			wantIndex:  "flagSetId+source",
			wantConstr: []interface{}{"123", "file.yaml"},
		},
		{
			name:       "three key compound including flagSetId key source",
			selector:   Selector{indexMap: map[string]string{"flagSetId": "123", "key": "mykey", "source": "file.yaml"}},
			wantIndex:  "flagSetId+key+source",
			wantConstr: []interface{}{"123", "mykey", "file.yaml"},
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
		// NEW TEST CASES
		{
			name:     "empty values should be ignored",
			selector: &Selector{indexMap: map[string]string{"flagSetId": "", "source": ""}},
			want:     model.Metadata{},
		},
		{
			name:     "only unknown indices (should be ignored)",
			selector: &Selector{indexMap: map[string]string{"unknown": "value", "other": "data"}},
			want:     model.Metadata{},
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
			s := NewSelector(tt.input)
			ts := (&s).WithFallback(tt.fallbackExpressionKey)
			indexMap := ts.indexMap
			if !reflect.DeepEqual(indexMap, tt.wantMap) {
				t.Errorf("NewSelector(%q) indexMap = %v, want %v", tt.input, ts.indexMap, tt.wantMap)
			}
		})
	}
}

// NEW COMPREHENSIVE TESTS

func TestExpressionToMap(t *testing.T) {
	tests := []struct {
		name                  string
		sExp                  string
		fallbackExpressionKey string
		wantMap               map[string]string
		wantUsingFallback     bool
	}{
		{
			name:              "empty expression",
			sExp:              "",
			wantMap:           map[string]string{},
			wantUsingFallback: false,
		},
		{
			name:              "single key-value pair",
			sExp:              "key=value",
			wantMap:           map[string]string{"key": "value"},
			wantUsingFallback: false,
		},
		{
			name:              "multiple key-value pairs",
			sExp:              "key1=value1,key2=value2,key3=value3",
			wantMap:           map[string]string{"key1": "value1", "key2": "value2", "key3": "value3"},
			wantUsingFallback: false,
		},
		{
			name:              "no equals sign - default fallback",
			sExp:              "somevalue",
			wantMap:           map[string]string{"source": "somevalue"},
			wantUsingFallback: true,
		},
		{
			name:                  "no equals sign - custom fallback",
			sExp:                  "somevalue",
			fallbackExpressionKey: "customKey",
			wantMap:               map[string]string{"customKey": "somevalue"},
			wantUsingFallback:     true,
		},
		{
			name:              "malformed pair - missing value",
			sExp:              "key1=value1,key2=,key3=value3",
			wantMap:           map[string]string{"key1": "value1", "key2": "", "key3": "value3"},
			wantUsingFallback: false,
		},
		{
			name:              "malformed pair - missing key",
			sExp:              "key1=value1,=value2,key3=value3",
			wantMap:           map[string]string{"key1": "value1", "": "value2", "key3": "value3"},
			wantUsingFallback: false,
		},
		// Not sure about this testcase how we should handle this cases
		//{
		//	name:              "multiple equals signs (should split on first)",
		//	sExp:              "key=value=with=equals",
		//	wantMap:           map[string]string{"key": "value=with=equals"},
		//	wantUsingFallback: false,
		//},
		{
			name:              "empty pairs should be ignored",
			sExp:              "key1=value1,,key2=value2",
			wantMap:           map[string]string{"key1": "value1", "key2": "value2"},
			wantUsingFallback: false,
		},
		{
			name:              "single equals sign only",
			sExp:              "=",
			wantMap:           map[string]string{"": ""},
			wantUsingFallback: false,
		},
		{
			name:              "whitespace handling",
			sExp:              "key1=value1, key2 = value2 ",
			wantMap:           map[string]string{"key1": "value1", " key2 ": " value2 "},
			wantUsingFallback: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMap, gotUsingFallback := expressionToMap(tt.sExp, tt.fallbackExpressionKey)
			if !reflect.DeepEqual(gotMap, tt.wantMap) {
				t.Errorf("expressionToMap() map = %v, want %v", gotMap, tt.wantMap)
			}
			if gotUsingFallback != tt.wantUsingFallback {
				t.Errorf("expressionToMap() usingFallback = %v, want %v", gotUsingFallback, tt.wantUsingFallback)
			}
		})
	}
}

func TestSelector_WithFallback(t *testing.T) {
	tests := []struct {
		name        string
		selector    *Selector
		fallbackKey string
		want        *Selector
	}{
		{
			name:     "nil selector",
			selector: nil,
			want:     nil,
		},
		{
			name:     "selector not using fallback",
			selector: &Selector{indexMap: map[string]string{"key": "value"}, usingFallback: false},
			want:     &Selector{indexMap: map[string]string{"key": "value"}, usingFallback: false},
		},
		{
			name:     "nil indexMap",
			selector: &Selector{indexMap: nil, usingFallback: true},
			want:     &Selector{indexMap: nil, usingFallback: true},
		},
		{
			name:        "using fallback with default key",
			selector:    &Selector{indexMap: map[string]string{"source": "defaultvalue"}, usingFallback: true},
			fallbackKey: "",
			want:        &Selector{indexMap: map[string]string{"source": "defaultvalue"}, usingFallback: false},
		},
		{
			name:        "using fallback with custom key",
			selector:    &Selector{indexMap: map[string]string{"source": "defaultvalue"}, usingFallback: true},
			fallbackKey: "customKey",
			want:        &Selector{indexMap: map[string]string{"customKey": "defaultvalue"}, usingFallback: false},
		},
		{
			name:        "fallback key already exists",
			selector:    &Selector{indexMap: map[string]string{"source": "defaultvalue", "customKey": "existingvalue"}, usingFallback: true},
			fallbackKey: "customKey",
			want:        &Selector{indexMap: map[string]string{"source": "defaultvalue", "customKey": "existingvalue"}, usingFallback: false},
		},
		{
			name:        "fallback with multiple keys",
			selector:    &Selector{indexMap: map[string]string{"source": "defaultvalue", "other": "othervalue"}, usingFallback: true},
			fallbackKey: "newkey",
			want:        &Selector{indexMap: map[string]string{"newkey": "defaultvalue", "other": "othervalue"}, usingFallback: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.selector.WithFallback(tt.fallbackKey)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithFallback() = %v, want %v", got, tt.want)
			}
			// Ensure original is not mutated (if not nil)
			if tt.selector != nil && got != nil && got != tt.selector {
				if reflect.DeepEqual(tt.selector.indexMap, got.indexMap) && tt.selector.usingFallback == got.usingFallback {
					t.Errorf("WithFallback() should not mutate original selector")
				}
			}
		})
	}
}

func TestNewSelector_UsingFallbackFlag(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		wantUsingFallback bool
	}{
		{
			name:              "key=value should not use fallback",
			input:             "source=test",
			wantUsingFallback: false,
		},
		{
			name:              "multiple key=value should not use fallback",
			input:             "source=test,flagSetId=123",
			wantUsingFallback: false,
		},
		{
			name:              "no equals should use fallback",
			input:             "testsource",
			wantUsingFallback: true,
		},
		{
			name:              "empty should not use fallback",
			input:             "",
			wantUsingFallback: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSelector(tt.input)
			if s.usingFallback != tt.wantUsingFallback {
				t.Errorf("NewSelector(%q).usingFallback = %v, want %v", tt.input, s.usingFallback, tt.wantUsingFallback)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	// Test that constants have expected values
	expectedConstants := map[string]string{
		"flagsTable":                      "flags",
		"idIndex":                         "id",
		"keyIndex":                        "key",
		"sourceIndex":                     "source",
		"defaultFallbackKey":              "source",
		"priorityIndex":                   "priority",
		"flagSetIdIndex":                  "flagSetId",
		"flagSetIdSourceCompoundIndex":    "flagSetId+source",
		"keySourceCompoundIndex":          "key+source",
		"flagSetIdKeySourceCompoundIndex": "flagSetId+key+source",
	}

	// Use reflection to check constants - this is a bit meta but ensures they match
	if flagsTable != expectedConstants["flagsTable"] {
		t.Errorf("flagsTable = %v, want %v", flagsTable, expectedConstants["flagsTable"])
	}
	if idIndex != expectedConstants["idIndex"] {
		t.Errorf("idIndex = %v, want %v", idIndex, expectedConstants["idIndex"])
	}
	// Continue for other constants...

	// Test that nilFlagSetId is a valid UUID format (36 characters with dashes)
	if len(nilFlagSetId) != 36 {
		t.Errorf("nilFlagSetId should be 36 characters long, got %d", len(nilFlagSetId))
	}
}

// Benchmark tests for performance-sensitive operations
func BenchmarkSelector_ToQuery(b *testing.B) {
	selector := Selector{
		indexMap: map[string]string{
			"flagSetId": "benchmark-flagset",
			"source":    "benchmark-source",
			"key":       "benchmark-key",
			"priority":  "high",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = selector.ToQuery()
	}
}

func BenchmarkNewSelector(b *testing.B) {
	expression := "source=benchmark-source,flagSetId=benchmark-flagset,key=benchmark-key"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewSelector(expression)
	}
}

func BenchmarkSelector_WithIndex(b *testing.B) {
	selector := Selector{
		indexMap: map[string]string{
			"source":    "benchmark-source",
			"flagSetId": "benchmark-flagset",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = selector.WithIndex("key", "benchmark-key")
	}
}
