package builder

import (
	"errors"
	"testing"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	"github.com/open-feature/flagd/core/pkg/sync/blob"
	buildermock "github.com/open-feature/flagd/core/pkg/sync/builder/mock"
	"github.com/open-feature/flagd/core/pkg/sync/file"
	"github.com/open-feature/flagd/core/pkg/sync/grpc"
	"github.com/open-feature/flagd/core/pkg/sync/http"
	"github.com/open-feature/flagd/core/pkg/sync/kubernetes"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSyncBuilder_SyncFromURI(t *testing.T) {
	type args struct {
		uri    string
		logger *logger.Logger
	}
	tests := []struct {
		name       string
		args       args
		injectFunc func(builder *SyncBuilder)
		want       sync.ISync
		wantErr    bool
	}{
		{
			name: "kubernetes sync",
			args: args{
				uri:    "core.openfeature.dev/ff-config",
				logger: logger.NewLogger(nil, false),
			},
			injectFunc: func(builder *SyncBuilder) {
				ctrl := gomock.NewController(t)

				mockClientBuilder := buildermock.NewMockIK8sClientBuilder(ctrl)
				mockClientBuilder.EXPECT().GetK8sClient().Times(1).Return(nil, nil)

				builder.k8sClientBuilder = mockClientBuilder
			},
			want:    &kubernetes.Sync{},
			wantErr: false,
		},
		{
			name: "kubernetes sync - error when retrieving config",
			args: args{
				uri:    "core.openfeature.dev/ff-config",
				logger: logger.NewLogger(nil, false),
			},
			injectFunc: func(builder *SyncBuilder) {
				ctrl := gomock.NewController(t)

				mockClientBuilder := buildermock.NewMockIK8sClientBuilder(ctrl)
				mockClientBuilder.EXPECT().GetK8sClient().Times(1).Return(nil, errors.New("oops"))

				builder.k8sClientBuilder = mockClientBuilder
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "file sync",
			args: args{
				uri:    "file:my-file",
				logger: logger.NewLogger(nil, false),
			},
			want:    &file.Sync{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewSyncBuilder()

			if tt.injectFunc != nil {
				tt.injectFunc(sb)
			}

			got, err := sb.SyncFromURI(tt.args.uri, tt.args.logger)
			if tt.wantErr {
				require.NotNil(t, err)
				require.Nil(t, got)
			} else {
				require.Nil(t, err)
				require.IsType(t, tt.want, got)
			}
		})
	}
}

func Test_k8sClusterConfig(t *testing.T) {
	t.Run("Cannot find KUBECONFIG file", func(tt *testing.T) {
		tt.Setenv("KUBECONFIG", "")
		_, err := k8sClusterConfig()
		if err == nil {
			tt.Error("Expected error but got none")
		}
	})
	t.Run("KUBECONFIG file not existing", func(tt *testing.T) {
		tt.Setenv("KUBECONFIG", "value")
		_, err := k8sClusterConfig()
		if err == nil {
			tt.Error("Expected error but got none")
		}
	})
	t.Run("Default REST Config and missing svc account", func(tt *testing.T) {
		tt.Setenv("KUBERNETES_SERVICE_HOST", "127.0.0.1")
		tt.Setenv("KUBERNETES_SERVICE_PORT", "8080")
		_, err := k8sClusterConfig()
		if err == nil {
			tt.Error("Expected error but got none")
		}
	})
}

func Test_SyncsFromFromConfig(t *testing.T) {
	lg := logger.NewLogger(nil, false)

	type args struct {
		logger  *logger.Logger
		sources []sync.SourceConfig
	}

	tests := []struct {
		name       string
		args       args
		injectFunc func(builder *SyncBuilder)
		wantSyncs  []sync.ISync
		wantErr    bool
	}{
		{
			name: "Empty",
			args: args{
				logger:  lg,
				sources: []sync.SourceConfig{},
			},
			wantSyncs: nil,
			wantErr:   false,
		},
		{
			name: "Error",
			args: args{
				logger: lg,
				sources: []sync.SourceConfig{
					{
						URI:      "fake",
						Provider: "disk",
					},
				},
			},
			wantSyncs: nil,
			wantErr:   true,
		},
		{
			name: "single",
			args: args{
				logger: lg,
				sources: []sync.SourceConfig{
					{
						URI:        "grpc://host:port",
						Provider:   syncProviderGrpc,
						ProviderID: "myapp",
						CertPath:   "/tmp/ca.cert",
						Selector:   "source=database",
					},
				},
			},
			wantSyncs: []sync.ISync{
				&grpc.Sync{},
			},
			wantErr: false,
		},
		{
			name: "grpc-with-msg-size",
			args: args{
				logger: lg,
				sources: []sync.SourceConfig{
					{
						URI:        "grpc://host:port",
						Provider:   syncProviderGrpc,
						ProviderID: "myapp",
						CertPath:   "/tmp/ca.cert",
						Selector:   "source=database",
						MaxMsgSize: 10,
					},
				},
			},
			wantSyncs: []sync.ISync{
				&grpc.Sync{},
			},
			wantErr: false,
		},
		{
			name: "combined",
			injectFunc: func(builder *SyncBuilder) {
				t.Setenv("AZURE_STORAGE_ACCOUNT", "myaccount")

				ctrl := gomock.NewController(t)

				mockClientBuilder := buildermock.NewMockIK8sClientBuilder(ctrl)
				mockClientBuilder.EXPECT().GetK8sClient().Times(1).Return(nil, nil)

				builder.k8sClientBuilder = mockClientBuilder
			},
			args: args{
				logger: lg,
				sources: []sync.SourceConfig{
					{
						URI:        "grpc://host:port",
						Provider:   syncProviderGrpc,
						ProviderID: "myapp",
						CertPath:   "/tmp/ca.cert",
						Selector:   "source=database",
					},
					{
						URI:        "https://host:port",
						Provider:   syncProviderHTTP,
						AuthHeader: "scheme credentials/token",
					},
					{
						URI:      "/tmp/flags.json",
						Provider: syncProviderFile,
					},
					{
						URI:      "my-namespace/my-flags",
						Provider: syncProviderKubernetes,
					},
					{
						URI:      "gs://bucket/path/to/file",
						Provider: syncProviderGcs,
					},
					{
						URI:      "azblob://bucket/path/to/file",
						Provider: syncProviderAzblob,
					},
					{
						URI:      "s3://bucket/path/to/file",
						Provider: syncProviderS3,
					},
				},
			},
			wantSyncs: []sync.ISync{
				&grpc.Sync{},
				&http.Sync{},
				&file.Sync{},
				&kubernetes.Sync{},
				&blob.Sync{},
				&blob.Sync{},
				&blob.Sync{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewSyncBuilder()

			if tt.injectFunc != nil {
				tt.injectFunc(sb)
			}
			syncs, err := sb.SyncsFromConfig(tt.args.sources, tt.args.logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("syncProvidersFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Len(t, syncs, len(tt.wantSyncs))

			// check if we got the expected sync types
			for index, wantType := range tt.wantSyncs {
				require.IsType(t, wantType, syncs[index])
			}
		})
	}
}

func Test_GcsConfig(t *testing.T) {
	lg := logger.NewLogger(nil, false)
	defaultInterval := uint32(5)
	tests := []struct {
		name             string
		uri              string
		interval         uint32
		expectedBucket   string
		expectedObject   string
		expectedInterval uint32
	}{
		{
			name:             "simple path",
			uri:              "gs://bucket/path/to/object",
			interval:         10,
			expectedBucket:   "gs://bucket/",
			expectedObject:   "path/to/object",
			expectedInterval: 10,
		},
		{
			name:             "default interval",
			uri:              "gs://bucket/path/to/object",
			expectedBucket:   "gs://bucket/",
			expectedObject:   "path/to/object",
			expectedInterval: defaultInterval,
		},
		{
			name:             "no object set", // Blob syncer will return error when fetching
			uri:              "gs://bucket/",
			expectedBucket:   "gs://bucket/",
			expectedObject:   "",
			expectedInterval: defaultInterval,
		},
		{
			name:             "malformed uri", // Blob syncer will return error when opening bucket
			uri:              "malformed",
			expectedBucket:   "",
			expectedObject:   "malformed",
			expectedInterval: defaultInterval,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gcsSync := NewSyncBuilder().newGcs(sync.SourceConfig{
				URI:      tt.uri,
				Interval: tt.interval,
			}, lg)
			require.Equal(t, tt.expectedBucket, gcsSync.Bucket)
			require.Equal(t, tt.expectedObject, gcsSync.Object)
			require.Equal(t, int(tt.expectedInterval), int(gcsSync.Interval))
		})
	}
}

func Test_AzblobConfig(t *testing.T) {
	lg := logger.NewLogger(nil, false)
	defaultInterval := uint32(5)
	tests := []struct {
		name             string
		uri              string
		interval         uint32
		storageAccount   string
		expectedBucket   string
		expectedObject   string
		expectedInterval uint32
		wantErr          bool
	}{
		{
			name:             "simple path",
			uri:              "azblob://bucket/path/to/object",
			interval:         10,
			storageAccount:   "myaccount",
			expectedBucket:   "azblob://bucket/",
			expectedObject:   "path/to/object",
			expectedInterval: 10,
			wantErr:          false,
		},
		{
			name:             "default interval",
			uri:              "azblob://bucket/path/to/object",
			storageAccount:   "myaccount",
			expectedBucket:   "azblob://bucket/",
			expectedObject:   "path/to/object",
			expectedInterval: defaultInterval,
			wantErr:          false,
		},
		{
			name:             "no object set", // Blob syncer will return error when fetching
			uri:              "azblob://bucket/",
			storageAccount:   "myaccount",
			expectedBucket:   "azblob://bucket/",
			expectedObject:   "",
			expectedInterval: defaultInterval,
			wantErr:          false,
		},
		{
			name:             "malformed uri", // Blob syncer will return error when opening bucket
			uri:              "malformed",
			storageAccount:   "myaccount",
			expectedBucket:   "",
			expectedObject:   "malformed",
			expectedInterval: defaultInterval,
			wantErr:          false,
		},
		{
			name:           "storage account not set", // Sync builder will fail and return error
			uri:            "azblob://bucket/path/to/object",
			storageAccount: "",
			wantErr:        true,
		},
		{
			name:           "storage account contains whitespace", // Sync builder will fail and return error
			uri:            "azblob://bucket/path/to/object",
			storageAccount: "my account",
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("AZURE_STORAGE_ACCOUNT", tt.storageAccount)
			azblobSync, err := NewSyncBuilder().newAzblob(sync.SourceConfig{
				URI:      tt.uri,
				Interval: tt.interval,
			}, lg)

			if (err != nil) != tt.wantErr {
				t.Errorf("newAzblob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (err != nil) && (tt.wantErr == true) {
				return
			}

			require.Equal(t, tt.expectedBucket, azblobSync.Bucket)
			require.Equal(t, tt.expectedObject, azblobSync.Object)
			require.Equal(t, int(tt.expectedInterval), int(azblobSync.Interval))
		})
	}
}

func Test_S3Config(t *testing.T) {
	lg := logger.NewLogger(nil, false)
	defaultInterval := uint32(5)
	tests := []struct {
		name             string
		uri              string
		interval         uint32
		expectedBucket   string
		expectedObject   string
		expectedInterval uint32
	}{
		{
			name:             "simple path",
			uri:              "s3://bucket/path/to/object",
			interval:         10,
			expectedBucket:   "s3://bucket/",
			expectedObject:   "path/to/object",
			expectedInterval: 10,
		},
		{
			name:             "default interval",
			uri:              "s3://bucket/path/to/object",
			expectedBucket:   "s3://bucket/",
			expectedObject:   "path/to/object",
			expectedInterval: defaultInterval,
		},
		{
			name:             "no object set", // Blob syncer will return error when fetching
			uri:              "s3://bucket/",
			expectedBucket:   "s3://bucket/",
			expectedObject:   "",
			expectedInterval: defaultInterval,
		},
		{
			name:             "malformed uri", // Blob syncer will return error when opening bucket
			uri:              "malformed",
			expectedBucket:   "",
			expectedObject:   "malformed",
			expectedInterval: defaultInterval,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s3Sync := NewSyncBuilder().newS3(sync.SourceConfig{
				URI:      tt.uri,
				Interval: tt.interval,
			}, lg)
			require.Equal(t, tt.expectedBucket, s3Sync.Bucket)
			require.Equal(t, tt.expectedObject, s3Sync.Object)
			require.Equal(t, int(tt.expectedInterval), int(s3Sync.Interval))
		})
	}
}
