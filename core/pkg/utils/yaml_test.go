package utils

import "testing"

func TestYAMLToJSON(t *testing.T) {
	tests := map[string]struct {
		input          []byte
		expected       string
		expectedError  bool
	}{
		"empty": {
			input:         []byte(""),
			expected:      "",
			expectedError: false,
		},
		"simple yaml": {
			input:         []byte("key: value"),
			expected:      `{"key":"value"}`,
			expectedError: false,
		},
		"nested yaml": {
			input:         []byte("parent:\n  child: value"),
			expected:      `{"parent":{"child":"value"}}`,
			expectedError: false,
		},
		"invalid yaml": {
			input:         []byte("invalid: yaml: : :"),
			expectedError: true,
		},
		"array yaml": {
			input:         []byte("items:\n  - item1\n  - item2"),
			expected:      `{"items":["item1","item2"]}`,
			expectedError: false,
		},
		"complex yaml": {
			input:         []byte("bool: true\nnum: 123\nstr: hello\nobj:\n  nested: value\narr:\n  - 1\n  - 2"),
			expected:      `{"arr":[1,2],"bool":true,"num":123,"obj":{"nested":"value"},"str":"hello"}`,
			expectedError: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := YAMLToJSON(tt.input)

			if tt.expectedError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if output != tt.expected {
				t.Errorf("expected output '%v', got '%v'", tt.expected, output)
			}
		})
	}
}