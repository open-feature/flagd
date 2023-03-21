package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/otel"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/metric"
	"os"
	"testing"
	"time"

	schemaGrpcV1 "buf.build/gen/go/open-feature/flagd/grpc/go/schema/v1/schemav1grpc"
	schemaV1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/schema/v1"
	"github.com/golang/mock/gomock"
	mock "github.com/open-feature/flagd/core/pkg/eval/mock"
	"github.com/open-feature/flagd/core/pkg/model"
	iservice "github.com/open-feature/flagd/core/pkg/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestConnectService_UnixConnection(t *testing.T) {
	type evalFields struct {
		result  bool
		variant string
		reason  string
		err     error
	}

	tests := []struct {
		name       string
		socketPath string
		evalFields evalFields
		req        *schemaV1.ResolveBooleanRequest
		want       *schemaV1.ResolveBooleanResponse
		wantErr    error
	}{
		{
			name:       "happy path",
			socketPath: "/tmp/flagd.sock",
			evalFields: evalFields{
				result:  true,
				variant: "on",
				reason:  model.DefaultReason,
				err:     nil,
			},
			req: &schemaV1.ResolveBooleanRequest{
				FlagKey: "myBoolFlag",
				Context: &structpb.Struct{},
			},
			want: &schemaV1.ResolveBooleanResponse{
				Value:   true,
				Reason:  model.DefaultReason,
				Variant: "on",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// try to ensure the socket file doesn't exist
			_ = os.Remove(tt.socketPath)
			ctrl := gomock.NewController(t)
			eval := mock.NewMockIEvaluator(ctrl)
			eval.EXPECT().ResolveBooleanValue(gomock.Any(), tt.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.evalFields.err,
			).AnyTimes()
			// configure OTel Metrics
			exp := metric.NewManualReader()
			metricRecorder := otel.NewOTelRecorder(exp, tt.name)
			svc := ConnectService{
				ConnectServiceConfiguration: &ConnectServiceConfiguration{
					ServerSocketPath: tt.socketPath,
				},
				Logger:  logger.NewLogger(nil, false),
				Metrics: metricRecorder,
			}
			serveConf := iservice.Configuration{
				ReadinessProbe: func() bool {
					return true
				},
			}
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			go func() {
				err := svc.Serve(ctx, eval, serveConf)
				fmt.Println(err)
			}()
			conn, err := grpc.Dial(
				fmt.Sprintf("unix://%s", tt.socketPath),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithBlock(),
				grpc.WithTimeout(2*time.Second),
			)
			if err != nil {
				t.Errorf("grpc - fail to dial: %v", err)
				return
			}
			client := schemaGrpcV1.NewServiceClient(
				conn,
			)

			res, err := client.ResolveBoolean(ctx, tt.req)
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("ConnectService.ResolveBoolean() error = %v, wantErr %v", err, tt.wantErr)
			}
			require.Equal(t, tt.want.Reason, res.Reason)
			require.Equal(t, tt.want.Value, res.Value)
			require.Equal(t, tt.want.Variant, res.Variant)
		})
	}
}
