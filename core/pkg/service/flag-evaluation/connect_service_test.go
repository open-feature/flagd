package service

import (
	"context"
	"errors"
	"fmt"
	middlewaremock "github.com/open-feature/flagd/core/pkg/service/middleware/mock"
	"net/http"
	"os"
	"testing"
	"time"

	schemaGrpcV1 "buf.build/gen/go/open-feature/flagd/grpc/go/schema/v1/schemav1grpc"
	schemaV1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/schema/v1"
	"github.com/golang/mock/gomock"
	mock "github.com/open-feature/flagd/core/pkg/eval/mock"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/otel"
	iservice "github.com/open-feature/flagd/core/pkg/service"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/metric"
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

func TestAddMiddleware(t *testing.T) {
	const port = 12345
	ctrl := gomock.NewController(t)

	mwMock := middlewaremock.NewMockIMiddleware(ctrl)

	mwMock.EXPECT().Handler(gomock.Any()).Return(
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusTeapot)
		}))

	exp := metric.NewManualReader()
	metricRecorder := otel.NewOTelRecorder(exp, "my-exporter")

	svc := ConnectService{
		ConnectServiceConfiguration: &ConnectServiceConfiguration{},
		Logger:                      logger.NewLogger(nil, false),
		Metrics:                     metricRecorder,
	}

	serveConf := iservice.Configuration{
		ReadinessProbe: func() bool {
			return true
		},
		Port: port,
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		err := svc.Serve(ctx, nil, serveConf)
		fmt.Println(err)
	}()

	require.Eventually(t, func() bool {
		return svc.server.Handler != nil
	}, 3*time.Second, 100*time.Millisecond)

	svc.AddMiddleware(mwMock)

	// call an endpoint provided by the server
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/schema.v1.Service/ResolveAll", port))

	require.Nil(t, err)
	// verify that the status we return in the mocked middleware
	require.Equal(t, http.StatusTeapot, resp.StatusCode)
}
