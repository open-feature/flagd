package runtime_test

import (
	"reflect"
	"testing"

	"github.com/open-feature/flagd/pkg/runtime"
	"github.com/open-feature/flagd/pkg/sync"
)

func TestSyncProviderArgParse(t *testing.T) {
	test := map[string]struct {
		in        string
		expectErr bool
		out       []sync.SourceConfig
	}{
		"simple": {
			in:        "[{\"uri\":\"config/samples/example_flags.json\",\"provider\":\"file\"}]",
			expectErr: false,
			out: []sync.SourceConfig{
				{
					URI:      "config/samples/example_flags.json",
					Provider: "file",
				},
			},
		},
		"multiple-syncs": {
			in: `[
					{"uri":"config/samples/example_flags.json","provider":"file"},
					{"uri":"http://test.com","provider":"http","bearerToken":":)"},
					{"uri":"host:port","provider":"grpc"},
					{"uri":"default/my-crd","provider":"kubernetes"}
				]`,
			expectErr: false,
			out: []sync.SourceConfig{
				{
					URI:      "config/samples/example_flags.json",
					Provider: "file",
				},
				{
					URI:         "http://test.com",
					Provider:    "http",
					BearerToken: ":)",
				},
				{
					URI:      "host:port",
					Provider: "grpc",
				},
				{
					URI:      "default/my-crd",
					Provider: "kubernetes",
				},
			},
		},
		"empty": {
			in:        `[]`,
			expectErr: false,
			out:       []sync.SourceConfig{},
		},
		"parse-failure": {
			in:        ``,
			expectErr: true,
			out:       []sync.SourceConfig{},
		},
	}

	for name, tt := range test {
		t.Run(name, func(t *testing.T) {
			out, err := runtime.SyncProviderArgParse(tt.in)
			if tt.expectErr {
				if err == nil {
					t.Error("expected error, got none")
				}
			} else if err != nil {
				t.Errorf("did not expect error: %s", err.Error())
			}
			if !reflect.DeepEqual(out, tt.out) {
				t.Errorf("unexpected output, expected %v, got %v", tt.out, out)
			}
		})
	}
}

func TestSyncProvidersFromURIs(t *testing.T) {
	test := map[string]struct {
		in        []string
		expectErr bool
		out       []sync.SourceConfig
	}{
		"simple": {
			in: []string{
				"file:my-file.json",
			},
			expectErr: false,
			out: []sync.SourceConfig{
				{
					URI:      "my-file.json",
					Provider: "file",
				},
			},
		},
		"multiple-uris": {
			in: []string{
				"file:my-file.json",
				"https://test.com",
				"grpc://host:port",
				"core.openfeature.dev/default/my-crd",
			},
			expectErr: false,
			out: []sync.SourceConfig{
				{
					URI:      "my-file.json",
					Provider: "file",
				},
				{
					URI:      "https://test.com",
					Provider: "http",
				},
				{
					URI:      "grpc://host:port",
					Provider: "grpc",
				},
				{
					URI:      "default/my-crd",
					Provider: "kubernetes",
				},
			},
		},
		"empty": {
			in:        []string{},
			expectErr: false,
			out:       []sync.SourceConfig{},
		},
		"parse-failure": {
			in:        []string{"care.openfeature.dev/will/fail"},
			expectErr: true,
			out:       []sync.SourceConfig{},
		},
	}

	for name, tt := range test {
		t.Run(name, func(t *testing.T) {
			out, err := runtime.SyncProvidersFromURIs(tt.in)
			if tt.expectErr {
				if err == nil {
					t.Error("expected error, got none")
				}
			} else if err != nil {
				t.Errorf("did not expect error: %s", err.Error())
			}
			if !reflect.DeepEqual(out, tt.out) {
				t.Errorf("unexpected output, expected %v, got %v", tt.out, out)
			}
		})
	}
}
