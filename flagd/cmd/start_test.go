package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseHeaderString(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected map[string]string
	}{
		"empty string": {
			input:    "",
			expected: map[string]string{},
		},
		"single header": {
			input:    "X-Proxy-Gateway-Host=b-flags-api.service",
			expected: map[string]string{"X-Proxy-Gateway-Host": "b-flags-api.service"},
		},
		"multiple headers": {
			input: "X-Proxy-Gateway-Host=myhost.service,X-Tenant-ID=tenant1",
			expected: map[string]string{
				"X-Proxy-Gateway-Host": "myhost.service",
				"X-Tenant-ID":           "tenant1",
			},
		},
		"value with equals sign": {
			input:    "Authorization=Bearer=token123",
			expected: map[string]string{"Authorization": "Bearer=token123"},
		},
		"whitespace around pairs": {
			input: "X-Custom=value , X-Other=val2",
			expected: map[string]string{
				"X-Custom": "value",
				"X-Other":  "val2",
			},
		},
		"empty value": {
			input:    "X-Empty=",
			expected: map[string]string{"X-Empty": ""},
		},
		"missing equals is skipped": {
			input:    "invalidentry",
			expected: map[string]string{},
		},
		"mix of valid and invalid": {
			input: "X-Valid=value,invalid,X-Also-Valid=ok",
			expected: map[string]string{
				"X-Valid":      "value",
				"X-Also-Valid": "ok",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := parseHeaderString(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}
