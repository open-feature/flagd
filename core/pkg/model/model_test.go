package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetErrorMessage(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{name: "flag not found", code: FlagNotFoundErrorCode, want: "Flag not found"},
		{name: "parse error", code: ParseErrorCode, want: "Error parsing input or configuration"},
		{name: "type mismatch", code: TypeMismatchErrorCode, want: "Type mismatch error"},
		{name: "general error", code: GeneralErrorCode, want: "General error"},
		{name: "flag disabled", code: FlagDisabledErrorCode, want: "Flag is disabled"},
		{name: "invalid context", code: InvalidContextCode, want: "Invalid context provided"},
		{name: "unknown code", code: "NOT_A_CODE", want: "Unknown error code: NOT_A_CODE"},
		{name: "empty code", code: "", want: "Unknown error code: "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GetErrorMessage(tt.code))
		})
	}
}

func TestFlagMarshalJSON(t *testing.T) {
	flag := Flag{
		Key:            "my-flag",
		FlagSetId:      "set-1",
		Priority:       7,
		State:          "ENABLED",
		DefaultVariant: "on",
		Variants:       map[string]any{"on": true, "off": false},
		Targeting:      json.RawMessage(`{"if":[true,"on","off"]}`),
		Source:         "file.json",
		Metadata:       Metadata{"team": "platform"},
	}

	b, err := json.Marshal(flag)
	require.NoError(t, err)

	var got map[string]any
	require.NoError(t, json.Unmarshal(b, &got))

	// MarshalJSON intentionally omits the key; locking this in so a refactor
	// cannot silently start emitting it.
	assert.NotContains(t, got, "key")
	// FlagSetId and Priority are json:"-" and must never serialize.
	assert.NotContains(t, got, "FlagSetId")
	assert.NotContains(t, got, "Priority")

	assert.Equal(t, "ENABLED", got["state"])
	assert.Equal(t, "on", got["defaultVariant"])
	assert.Equal(t, "file.json", got["source"])
	assert.Contains(t, got, "variants")
	assert.Contains(t, got, "targeting")
	assert.Contains(t, got, "metadata")
}

func TestFlagMarshalJSON_OmitsEmptyOptionalFields(t *testing.T) {
	b, err := json.Marshal(Flag{State: "ENABLED", DefaultVariant: "on", Source: "file.json"})
	require.NoError(t, err)

	var got map[string]any
	require.NoError(t, json.Unmarshal(b, &got))

	assert.NotContains(t, got, "key")
	assert.NotContains(t, got, "targeting")
	assert.NotContains(t, got, "metadata")
}
