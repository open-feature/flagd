package builder

import (
	"reflect"
	"testing"

	"github.com/open-feature/flagd/core/pkg/sync"
)

func TestParseSource(t *testing.T) {
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
					Provider: syncProviderFile,
				},
			},
		},
		"multiple-syncs": {
			in: `[
					{"uri":"config/samples/example_flags.json","provider":"file"},
					{"uri":"http://test.com","provider":"http","bearerToken":":)"},
					{"uri":"host:port","provider":"grpc"},
					{"uri":"default/my-crd","provider":"kubernetes"},
					{"uri":"gs://bucket-name/path/to/file","provider":"gcs"},
					{"uri":"azblob://bucket-name/path/to/file","provider":"azblob"},
					{"uri":"s3://bucket-name/path/to/file","provider":"s3"}
				]`,
			expectErr: false,
			out: []sync.SourceConfig{
				{
					URI:      "config/samples/example_flags.json",
					Provider: syncProviderFile,
				},
				{
					URI:         "http://test.com",
					Provider:    syncProviderHTTP,
					BearerToken: ":)",
				},
				{
					URI:      "host:port",
					Provider: syncProviderGrpc,
				},
				{
					URI:      "default/my-crd",
					Provider: syncProviderKubernetes,
				},
				{
					URI:      "gs://bucket-name/path/to/file",
					Provider: syncProviderGcs,
				},
				{
					URI:      "azblob://bucket-name/path/to/file",
					Provider: syncProviderAzblob,
				},
				{
					URI:      "s3://bucket-name/path/to/file",
					Provider: syncProviderS3,
				},
			},
		},
		"multiple-syncs-with-options": {
			in: `[
				{"uri":"config/samples/example_flags.json","provider":"file"},
				{"uri":"http://my-flag-source.json","provider":"http","bearerToken":"bearer-dji34ld2l"},
				{"uri":"https://secure-remote","provider":"http","bearerToken":"bearer-dji34ld2l"},
				{"uri":"https://secure-remote","provider":"http","authHeader":"Bearer bearer-dji34ld2l"},
				{"uri":"https://secure-remote","provider":"http","authHeader":"Basic dXNlcjpwYXNz"},
				{"uri":"http://site.com","provider":"http","interval":77 },
				{"uri":"default/my-flag-config","provider":"kubernetes"},
				{"uri":"grpc-source:8080","provider":"grpc"},
				{"uri":"my-flag-source:8080","provider":"grpc", "tls":true, "certPath": "/certs/ca.cert", "providerID": "flagd-weatherapp-sidecar", "selector": "source=database,app=weatherapp"}
			]`,
			expectErr: false,
			out: []sync.SourceConfig{
				{
					URI:      "config/samples/example_flags.json",
					Provider: syncProviderFile,
				},
				{
					URI:         "http://my-flag-source.json",
					Provider:    syncProviderHTTP,
					BearerToken: "bearer-dji34ld2l",
				},
				{
					URI:         "https://secure-remote",
					Provider:    syncProviderHTTP,
					BearerToken: "bearer-dji34ld2l",
				},
				{
					URI:        "https://secure-remote",
					Provider:   syncProviderHTTP,
					AuthHeader: "Bearer bearer-dji34ld2l",
				},
				{
					URI:        "https://secure-remote",
					Provider:   syncProviderHTTP,
					AuthHeader: "Basic dXNlcjpwYXNz",
				},
				{
					URI:      "http://site.com",
					Provider: syncProviderHTTP,
					Interval: 77,
				},
				{
					URI:      "default/my-flag-config",
					Provider: syncProviderKubernetes,
				},
				{
					URI:      "grpc-source:8080",
					Provider: syncProviderGrpc,
				},
				{
					URI:        "my-flag-source:8080",
					Provider:   syncProviderGrpc,
					TLS:        true,
					CertPath:   "/certs/ca.cert",
					ProviderID: "flagd-weatherapp-sidecar",
					Selector:   "source=database,app=weatherapp",
				},
			},
		},
		"multiple-auth-options": {
			in: `[
				{"uri":"https://secure-remote","provider":"http","authHeader":"Bearer bearer-dji34ld2l","bearerToken":"bearer-dji34ld2l"}
			]`,
			expectErr: true,
			out: []sync.SourceConfig{
				{
					URI:         "https://secure-remote",
					Provider:    syncProviderHTTP,
					AuthHeader:  "Bearer bearer-dji34ld2l",
					BearerToken: "bearer-dji34ld2l",
					TLS:         false,
					Interval:    0,
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
			out, err := ParseSources(tt.in)
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

func TestParseSyncProviderURIs(t *testing.T) {
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
				"grpcs://secure-grpc",
				"core.openfeature.dev/default/my-crd",
				"gs://bucket-name/path/to/file",
				"azblob://bucket-name/path/to/file",
				"s3://bucket-name/path/to/file",
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
					URI:      "host:port",
					Provider: "grpc",
					TLS:      false,
				},
				{
					URI:      "secure-grpc",
					Provider: "grpc",
					TLS:      true,
				},
				{
					URI:      "default/my-crd",
					Provider: "kubernetes",
				},
				{
					URI:      "gs://bucket-name/path/to/file",
					Provider: syncProviderGcs,
				},
				{
					URI:      "azblob://bucket-name/path/to/file",
					Provider: syncProviderAzblob,
				},
				{
					URI:      "s3://bucket-name/path/to/file",
					Provider: syncProviderS3,
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
			out, err := ParseSyncProviderURIs(tt.in)
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

func TestParseOAuth(t *testing.T) {
	test := map[string]struct {
		in        string
		expectErr bool
		out       []sync.SourceConfig
	}{
		"noOauth": {
			in:        "[{\"uri\":\"https://secure-remote\",\"provider\":\"http\",\"authHeader\":\"Bearer bearer-dji34ld2l\"}]",
			expectErr: false,
			out: []sync.SourceConfig{
				{
					URI:        "https://secure-remote",
					Provider:   "http",
					AuthHeader: "Bearer bearer-dji34ld2l",
				},
			},
		},
		"oauth": {
			in: `[{
	"uri": "https://secure-remote",
	"provider":"http", 
	"oauth": { 
		"clientID": "myID",
		"clientSecret": "mySecret",
		"tokenURL": "myTokenUrl" 
	}}]`,
			expectErr: false,
			out: []sync.SourceConfig{
				{
					URI:      "https://secure-remote",
					Provider: "http",
					OAuth: &sync.OAuthCredentialHandler{
						ClientID:     "myID",
						ClientSecret: "mySecret",
						TokenURL:     "myTokenUrl",
					},
				},
			},
		},
	}

	for name, tt := range test {
		t.Run(name, func(t *testing.T) {
			out, err := ParseSources(tt.in)
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
