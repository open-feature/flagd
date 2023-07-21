package service

import (
	"context"
	"errors"
	"fmt"
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
	iservice "github.com/open-feature/flagd/core/pkg/service"
	middlewaremock "github.com/open-feature/flagd/core/pkg/service/middleware/mock"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestConnectService_UnixConnection(t *testing.T) {
	type evalFields struct {
		result  bool
		variant string
		reason  string
		metadat map[string]interface{}
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
			eval.EXPECT().ResolveBooleanValue(gomock.Any(), gomock.Any(), tt.req.FlagKey, gomock.Any()).Return(
				tt.evalFields.result,
				tt.evalFields.variant,
				tt.evalFields.reason,
				tt.evalFields.metadat,
				tt.evalFields.err,
			).AnyTimes()
			// configure OTel Metrics
			exp := metric.NewManualReader()
			rs := resource.NewWithAttributes("testSchema")
			metricRecorder := telemetry.NewOTelRecorder(exp, rs, tt.name)
			svc := NewConnectService(logger.NewLogger(nil, false), eval, metricRecorder)
			serveConf := iservice.Configuration{
				ReadinessProbe: func() bool {
					return true
				},
				SocketPath: tt.socketPath,
			}
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			go func() {
				err := svc.Serve(ctx, serveConf)
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
			writer.WriteHeader(http.StatusOK)
		}))

	exp := metric.NewManualReader()
	rs := resource.NewWithAttributes("testSchema")
	metricRecorder := telemetry.NewOTelRecorder(exp, rs, "my-exporter")

	svc := NewConnectService(logger.NewLogger(nil, false), nil, metricRecorder)

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
		err := svc.Serve(ctx, serveConf)
		fmt.Println(err)
	}()

	require.Eventually(t, func() bool {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/schema.v1.Service/ResolveAll", port))
		// with the default http handler we should get a method not allowed (405) when attempting a GET request
		return err == nil && resp.StatusCode == http.StatusMethodNotAllowed
	}, 3*time.Second, 100*time.Millisecond)

	svc.AddMiddleware(mwMock)

	// with the injected middleware, the GET method should work
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/schema.v1.Service/ResolveAll", port))

	require.Nil(t, err)
	// verify that the status we return in the mocked middleware
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestConnectServiceNotify(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	eval := mock.NewMockIEvaluator(ctrl)

	exp := metric.NewManualReader()
	rs := resource.NewWithAttributes("testSchema")
	metricRecorder := telemetry.NewOTelRecorder(exp, rs, "my-exporter")

	service := NewConnectService(logger.NewLogger(nil, false), eval, metricRecorder)

	sChan := make(chan iservice.Notification, 1)
	eventing := service.eventingConfiguration
	eventing.subs["key"] = sChan

	// notification type
	ofType := iservice.ConfigurationChange

	// emit notification in routine
	go func() {
		service.Notify(iservice.Notification{
			Type: ofType,
			Data: map[string]interface{}{},
		})
	}()

	// wait for notification
	timeout, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()

	select {
	case n := <-sChan:
		require.Equal(t, ofType, n.Type, "expected notification type: %s, but received %s", ofType, n.Type)
	case <-timeout.Done():
		t.Error("timeout while waiting for notifications")
	}
}

func TestConnectServiceShutdown(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	eval := mock.NewMockIEvaluator(ctrl)

	exp := metric.NewManualReader()
	rs := resource.NewWithAttributes("testSchema")
	metricRecorder := telemetry.NewOTelRecorder(exp, rs, "my-exporter")

	service := NewConnectService(logger.NewLogger(nil, false), eval, metricRecorder)

	sChan := make(chan iservice.Notification, 1)
	eventing := service.eventingConfiguration
	eventing.subs["key"] = sChan

	// notification type
	ofType := iservice.Shutdown

	// emit notification in routine
	go func() {
		service.Notify(iservice.Notification{
			Type: ofType,
			Data: map[string]interface{}{},
		})
	}()

	// wait for notification
	timeout, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()

	require.False(t, service.readinessEnabled)

	select {
	case n := <-sChan:
		require.Equal(t, ofType, n.Type, "expected notification type: %s, but received %s", ofType, n.Type)
	case <-timeout.Done():
		t.Error("timeout while waiting for notifications")
	}
}
