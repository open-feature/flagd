package utils

import "testing"

func TestGenerateSha(t *testing.T) {
	t.Run("ignores object key ordering", func(t *testing.T) {
		a := []byte(`{"flags":{"a":{},"b":{}},"metadata":{"x":1}}`)
		b := []byte(`{"metadata":{"x":1},"flags":{"b":{},"a":{}}}`)
		if GenerateSha(a) != GenerateSha(b) {
			t.Errorf("expected key-reordered JSON to hash equally:\n a=%s\n b=%s", GenerateSha(a), GenerateSha(b))
		}
	})

	t.Run("ignores insignificant whitespace", func(t *testing.T) {
		compact := []byte(`{"flags":{"a":{}}}`)
		spaced := []byte("{\n  \"flags\": {\n    \"a\": {}\n  }\n}")
		if GenerateSha(compact) != GenerateSha(spaced) {
			t.Error("expected formatting-only differences to hash equally")
		}
	})

	t.Run("detects genuine content changes", func(t *testing.T) {
		a := []byte(`{"flags":{"a":{}}}`)
		b := []byte(`{"flags":{"b":{}}}`)
		if GenerateSha(a) == GenerateSha(b) {
			t.Error("expected different content to hash differently")
		}
	})

	t.Run("falls back to raw bytes for non-json", func(t *testing.T) {
		if GenerateSha([]byte("not json")) == GenerateSha([]byte("also not json")) {
			t.Error("expected distinct non-json payloads to hash differently")
		}
	})
}
