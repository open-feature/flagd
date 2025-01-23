package utils

import (
	"strings"
	"testing"
)

func TestConvertToJSON(t *testing.T) {
	tests := map[string]struct {
		data          []byte
		fileExtension string
		mediaType     string
		want          string
		wantErr       bool
		errContains   string
	}{
		"json file type": {
			data:          []byte(`{"flags": {"foo": "bar"}}`),
			fileExtension: "json",
			mediaType:     "application/json",
			want:          `{"flags": {"foo": "bar"}}`,
			wantErr:       false,
		},
		"json file type with encoding": {
			data:          []byte(`{"flags": {"foo": "bar"}}`),
			fileExtension: "json",
			mediaType:     "application/json; charset=utf-8",
			want:          `{"flags": {"foo": "bar"}}`,
			wantErr:       false,
		},
		"yaml file type": {
			data:          []byte("flags:\n  foo: bar"),
			fileExtension: "yaml",
			mediaType:     "application/yaml",
			want:          `{"flags":{"foo":"bar"}}`,
			wantErr:       false,
		},
		"yaml file type with encoding": {
			data:          []byte("flags:\n  foo: bar"),
			fileExtension: "yaml",
			mediaType:     "application/yaml; charset=utf-8",
			want:          `{"flags":{"foo":"bar"}}`,
			wantErr:       false,
		},
		"yml file type": {
			data:          []byte("flags:\n  foo: bar"),
			fileExtension: "yml",
			mediaType:     "application/x-yaml",
			want:          `{"flags":{"foo":"bar"}}`,
			wantErr:       false,
		},
		"invalid yaml": {
			data:          []byte("invalid: [yaml: content"),
			fileExtension: "yaml",
			mediaType:     "application/yaml",
			wantErr:       true,
			errContains:   "error converting blob from yaml to json",
		},
		"unsupported file type": {
			data:          []byte("some content"),
			fileExtension: "txt",
			mediaType:     "text/plain",
			wantErr:       true,
			errContains:   "unsupported file format",
		},
		"empty file type with valid media type": {
			data:          []byte(`{"flags": {"foo": "bar"}}`),
			fileExtension: "",
			mediaType:     "application/json",
			want:          `{"flags": {"foo": "bar"}}`,
			wantErr:       false,
		},
		"invalid media type": {
			data:          []byte("some content"),
			fileExtension: "",
			mediaType:     "invalid/\\type",
			wantErr:       true,
			errContains:   "unable to determine file format",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ConvertToJSON(tt.data, tt.fileExtension, tt.mediaType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil || !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ConvertToJSON() expected error containing %q, got %v", tt.errContains, err)
				}
				return
			}
			if got != tt.want {
				t.Errorf("ConvertToJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
