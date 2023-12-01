package builder

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	buildermock "github.com/open-feature/flagd/core/pkg/sync/builder/mock"
	"github.com/open-feature/flagd/core/pkg/sync/file"
	"github.com/open-feature/flagd/core/pkg/sync/kubernetes"
	"github.com/stretchr/testify/require"
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
				mockClientBuilder.EXPECT().GetK8sClients().Times(1).Return(nil, nil, nil)

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
				mockClientBuilder.EXPECT().GetK8sClients().Times(1).Return(nil, nil, errors.New("oops"))

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
