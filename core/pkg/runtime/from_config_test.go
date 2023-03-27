package runtime

import (
	"reflect"
	"testing"

	"github.com/open-feature/flagd/core/pkg/logger"
)

func TestParseSource(t *testing.T) {
	test := map[string]struct {
		in        string
		expectErr bool
		out       []SourceConfig
	}{
		"simple": {
			in:        "[{\"uri\":\"config/samples/example_flags.json\",\"provider\":\"file\"}]",
			expectErr: false,
			out: []SourceConfig{
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
					{"uri":"default/my-crd","provider":"kubernetes"}
				]`,
			expectErr: false,
			out: []SourceConfig{
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
			},
		},
		"multiple-syncs-with-options": {
			in: `[{"uri":"config/samples/example_flags.json","provider":"file"},
            		{"uri":"http://my-flag-source.json","provider":"http","bearerToken":"bearer-dji34ld2l"},
            		{"uri":"https://secure-remote","provider":"http","bearerToken":"bearer-dji34ld2l"},
					{"uri":"default/my-flag-config","provider":"kubernetes"},
            		{"uri":"grpc-source:8080","provider":"grpc"},
            		{"uri":"my-flag-source:8080","provider":"grpc", "grpcSecure":true, "certPath": "/certs/ca.cert", "providerID": "flagd-weatherapp-sidecar", "selector": "source=database,app=weatherapp"}]
				`,
			expectErr: false,
			out: []SourceConfig{
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
					GrpcSecure: true,
					CertPath:   "/certs/ca.cert",
					ProviderID: "flagd-weatherapp-sidecar",
					Selector:   "source=database,app=weatherapp",
				},
			},
		},
		"empty": {
			in:        `[]`,
			expectErr: false,
			out:       []SourceConfig{},
		},
		"parse-failure": {
			in:        ``,
			expectErr: true,
			out:       []SourceConfig{},
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
		out       []SourceConfig
	}{
		"simple": {
			in: []string{
				"file:my-file.json",
			},
			expectErr: false,
			out: []SourceConfig{
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
			},
			expectErr: false,
			out: []SourceConfig{
				{
					URI:      "my-file.json",
					Provider: "file",
				},
				{
					URI:      "https://test.com",
					Provider: "http",
				},
				{
					URI:        "host:port",
					Provider:   "grpc",
					GrpcSecure: false,
				},
				{
					URI:        "secure-grpc",
					Provider:   "grpc",
					GrpcSecure: true,
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
			out:       []SourceConfig{},
		},
		"parse-failure": {
			in:        []string{"care.openfeature.dev/will/fail"},
			expectErr: true,
			out:       []SourceConfig{},
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

// Note - K8s configuration require K8s client, hence do not use K8s sync provider in this test
func Test_syncProvidersFromConfig(t *testing.T) {
	lg := logger.NewLogger(nil, false)

	type args struct {
		logger  *logger.Logger
		sources []SourceConfig
	}

	tests := []struct {
		name      string
		args      args
		wantSyncs int // simply check the count of ISync providers yield from configurations
		wantErr   bool
	}{
		{
			name: "Empty",
			args: args{
				logger:  lg,
				sources: []SourceConfig{},
			},
			wantSyncs: 0,
			wantErr:   false,
		},
		{
			name: "Error",
			args: args{
				logger: lg,
				sources: []SourceConfig{
					{
						URI:      "fake",
						Provider: "disk",
					},
				},
			},
			wantSyncs: 0,
			wantErr:   true,
		},
		{
			name: "single",
			args: args{
				logger: lg,
				sources: []SourceConfig{
					{
						URI:        "grpc://host:port",
						Provider:   syncProviderGrpc,
						ProviderID: "myapp",
						CertPath:   "/tmp/ca.cert",
						Selector:   "source=database",
					},
				},
			},
			wantSyncs: 1,
			wantErr:   false,
		},
		{
			name: "combined",
			args: args{
				logger: lg,
				sources: []SourceConfig{
					{
						URI:        "grpc://host:port",
						Provider:   syncProviderGrpc,
						ProviderID: "myapp",
						CertPath:   "/tmp/ca.cert",
						Selector:   "source=database",
					},
					{
						URI:         "https://host:port",
						Provider:    syncProviderHTTP,
						BearerToken: "token",
					},
					{
						URI:      "/tmp/flags.json",
						Provider: syncProviderFile,
					},
				},
			},
			wantSyncs: 3,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			syncs, err := syncProvidersFromConfig(tt.args.logger, tt.args.sources)
			if (err != nil) != tt.wantErr {
				t.Errorf("syncProvidersFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantSyncs != len(syncs) {
				t.Errorf("syncProvidersFromConfig() yielded = %v, but expected %v", len(syncs), tt.wantSyncs)
			}
		})
	}
}
