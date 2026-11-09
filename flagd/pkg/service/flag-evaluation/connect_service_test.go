package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	evaluationV1 "buf.build/gen/go/open-feature/flagd/connectrpc/go/flagd/evaluation/v1/evaluationv1connect"
	evaluationV2 "buf.build/gen/go/open-feature/flagd/connectrpc/go/flagd/evaluation/v2/evaluationv2connect"
	schemaConnectV1 "buf.build/gen/go/open-feature/flagd/connectrpc/go/schema/v1/schemav1connect"
	schemaGrpcV1 "buf.build/gen/go/open-feature/flagd/grpc/go/schema/v1/schemav1grpc"
	schemaV1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/schema/v1"
	mock "github.com/open-feature/flagd/core/pkg/evaluator/mock"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	iservice "github.com/open-feature/flagd/core/pkg/service"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	middlewaremock "github.com/open-feature/flagd/flagd/pkg/service/middleware/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
)

const resolveAllURLFmt = "http://localhost:%d/flagd.evaluation.v1.Service/ResolveAll"

// Verifies bounded (route, method) labels on HTTP server metrics.
func TestConnectServiceHTTPMetricLabels(t *testing.T) {
	reader := metric.NewManualReader()
	recorder := telemetry.NewOTelRecorder(reader, resource.NewWithAttributes("testSchema"), "testSvc")
	svc := NewConnectService(logger.NewLogger(nil, false), nil, nil, recorder)
	svc.setupHandler(iservice.Configuration{ServiceName: "testSvc"})

	requests := []struct {
		method string
		path   string
		status int
	}{
		{method: http.MethodGet, path: schemaConnectV1.ServiceResolveAllProcedure},
		{method: http.MethodGet, path: evaluationV1.ServiceResolveBooleanProcedure},
		{method: "PROPFIND", path: evaluationV2.ServiceResolveStringProcedure},
		{method: http.MethodGet, path: "/unregistered/first", status: http.StatusNotFound},
		{method: http.MethodGet, path: evaluationV1.ServiceResolveBooleanProcedure + "/extra", status: http.StatusNotFound},
		{method: "PROPFIND", path: "/unregistered/second", status: http.StatusNotFound},
		{method: "CUSTOM", path: "/unregistered/third", status: http.StatusNotFound},
	}
	for _, req := range requests {
		response := httptest.NewRecorder()
		svc.server.Handler.ServeHTTP(response, httptest.NewRequest(req.method, req.path, nil))
		if req.status != 0 {
			require.Equal(t, req.status, response.Code, "%s %s", req.method, req.path)
		}
	}

	var data metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &data))
	want := map[[2]string]bool{
		{schemaConnectV1.ServiceResolveAllProcedure, http.MethodGet}:  true,
		{evaluationV1.ServiceResolveBooleanProcedure, http.MethodGet}: true,
		{evaluationV2.ServiceResolveStringProcedure, "_OTHER"}:        true,
		{"/", http.MethodGet}: true,
		{"/", "_OTHER"}:       true,
	}
	seen := collectHTTPMetricLabelPairs(t, data, want)
	require.Len(t, seen, 3)
}

// httpServerMetrics are the OTel HTTP server metric names emitted by the recorder.
var httpServerMetrics = map[string]bool{
	"http.server.request.duration":   true,
	"http.server.response.body.size": true,
	"http.server.active_requests":    true,
}

// collectHTTPMetricLabelPairs asserts each HTTP metric has exactly the want (route, method) pairs; returns names seen.
func collectHTTPMetricLabelPairs(t *testing.T, data metricdata.ResourceMetrics, want map[[2]string]bool) map[string]bool {
	t.Helper()

	seen := make(map[string]bool)
	for _, scope := range data.ScopeMetrics {
		for _, m := range scope.Metrics {
			if !httpServerMetrics[m.Name] {
				continue
			}
			require.Equal(t, want, metricRouteMethodPairs(t, m), m.Name)
			seen[m.Name] = true
		}
	}
	return seen
}

// metricRouteMethodPairs returns the set of (route, method) label pairs on the metric's data points.
func metricRouteMethodPairs(t *testing.T, m metricdata.Metrics) map[[2]string]bool {
	t.Helper()

	got := make(map[[2]string]bool)
	collect := func(attrs attribute.Set) {
		route, _ := attrs.Value(semconv.HTTPRouteKey)
		method, _ := attrs.Value(semconv.HTTPRequestMethodKey)
		got[[2]string{route.AsString(), method.AsString()}] = true
	}
	switch points := m.Data.(type) {
	case metricdata.Histogram[float64]:
		for _, point := range points.DataPoints {
			collect(point.Attributes)
		}
	case metricdata.Sum[int64]:
		for _, point := range points.DataPoints {
			collect(point.Attributes)
		}
	default:
		t.Fatalf("unexpected data type for %s: %T", m.Name, m.Data)
	}
	return got
}

func TestConnectServiceUnixConnection(t *testing.T) {
	type evalFields struct {
		result   bool
		variant  string
		reason   string
		metadata map[string]interface{}
		err      error
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
				tt.evalFields.metadata,
				tt.evalFields.err,
			).AnyTimes()
			// configure OTel Metrics
			exp := metric.NewManualReader()
			rs := resource.NewWithAttributes("testSchema")
			metricRecorder := telemetry.NewOTelRecorder(exp, rs, tt.name)
			svc := NewConnectService(logger.NewLogger(nil, false), eval, nil, metricRecorder)
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

	svc := NewConnectService(logger.NewLogger(nil, false), nil, nil, metricRecorder)

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
		resp, err := http.Get(fmt.Sprintf(resolveAllURLFmt, port))
		if err == nil && resp != nil {
			resp.Body.Close()
		}
		// with the default http handler we should get a method not allowed (405) when attempting a GET request
		return err == nil && resp.StatusCode == http.StatusMethodNotAllowed
	}, 3*time.Second, 100*time.Millisecond)

	svc.AddMiddleware(mwMock)

	// with the injected middleware, the GET method should work
	resp, err := http.Get(fmt.Sprintf(resolveAllURLFmt, port))

	require.Nil(t, err)
	defer resp.Body.Close()
	// verify that the status we return in the mocked middleware
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestConnectServiceNotify(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	eval := mock.NewMockIEvaluator(ctrl)
	sources := []string{"source1", "source2"}
	log := logger.NewLogger(nil, false)
	s, err := store.NewStore(log, sources)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}

	exp := metric.NewManualReader()
	rs := resource.NewWithAttributes("testSchema")
	metricRecorder := telemetry.NewOTelRecorder(exp, rs, "my-exporter")

	service := NewConnectService(logger.NewLogger(nil, false), eval, s, metricRecorder)

	sChan := make(chan iservice.Notification, 1)
	eventing := service.eventingConfiguration
	eventing.Subscribe(context.Background(), "key", nil, sChan)

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

func TestConnectServiceWatcher(t *testing.T) {
	sources := []string{"source1", "source2"}
	log := logger.NewLogger(nil, false)
	s, err := store.NewStore(log, sources)

	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}

	sChan := make(chan iservice.Notification, 1)
	eventing := eventingConfiguration{
		store:  s,
		logger: log,
		mu:     &sync.RWMutex{},
		subs:   make(map[any]*subscription),
	}

	// subscribe and wait for for the sub to be active
	eventing.Subscribe(context.Background(), "anything", nil, sChan)
	time.Sleep(100 * time.Millisecond)

	// make a change
	s.Update(sources[0], []model.Flag{
		{
			Key:            "flag1",
			DefaultVariant: "off",
		},
	}, model.Metadata{}, false)

	// notification type
	ofType := iservice.ConfigurationChange

	timeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	select {
	case n := <-sChan:
		require.Equal(t, ofType, n.Type, "expected notification type: %s, but received %s", ofType, n.Type)
		flags := n.Data["flags"].(map[string]interface{})
		flag1, ok := flags["flag1"].(map[string]interface{})
		require.True(t, ok, "flag1 notification should be a map[string]interface{}")
		require.Equal(t, flag1["type"], string(model.NotificationCreate), "expected notification type: %s, but received %s", model.NotificationCreate, flag1["type"])
	case <-timeout.Done():
		t.Error("timeout while waiting for notifications")
	}
}

func TestConnectServiceShutdown(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	eval := mock.NewMockIEvaluator(ctrl)
	sources := []string{"source1", "source2"}
	log := logger.NewLogger(nil, false)
	s, err := store.NewStore(log, sources)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}

	exp := metric.NewManualReader()
	rs := resource.NewWithAttributes("testSchema")
	metricRecorder := telemetry.NewOTelRecorder(exp, rs, "my-exporter")

	service := NewConnectService(logger.NewLogger(nil, false), eval, s, metricRecorder)

	sChan := make(chan iservice.Notification, 1)
	eventing := service.eventingConfiguration
	eventing.Subscribe(context.Background(), "key", nil, sChan)

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

// startConnectService creates a ConnectService with a mock evaluator and metric recorder,
// starts it in a background goroutine with the given configuration, and waits until it is ready.
// It returns the port the service is listening on.
func startConnectService(t *testing.T, port uint16, conf iservice.Configuration) {
	t.Helper()

	ctrl := gomock.NewController(t)
	eval := mock.NewMockIEvaluator(ctrl)

	exp := metric.NewManualReader()
	rs := resource.NewWithAttributes("testSchema")
	metricRecorder := telemetry.NewOTelRecorder(exp, rs, "limit-test")

	svc := NewConnectService(logger.NewLogger(nil, false), eval, nil, metricRecorder)

	conf.ReadinessProbe = func() bool { return true }
	conf.Port = port

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	go func() {
		_ = svc.Serve(ctx, conf)
	}()

	require.Eventually(t, func() bool {
		resp, err := http.Get(fmt.Sprintf(resolveAllURLFmt, port))
		if err == nil && resp != nil {
			resp.Body.Close()
		}
		return err == nil && resp != nil
	}, 3*time.Second, 100*time.Millisecond)
}

func TestConnectServiceRequestBodySizeLimit(t *testing.T) {
	const port = 18291

	startConnectService(t, port, iservice.Configuration{
		MaxRequestBodyBytes: 10, // allow only 10 bytes
	})

	// Valid JSON that exceeds the 10-byte body limit, so MaxBytesReader fires mid-parse.
	largeBody := []byte(`{"flagKey":"` + strings.Repeat("a", 100) + `"}`)
	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("http://localhost:%d/flagd.evaluation.v1.Service/ResolveBoolean", port),
		bytes.NewReader(largeBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	// connect-go maps MaxBytesError (resource exhausted) to HTTP 429.
	require.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
}

func TestConnectServiceRequestHeaderSizeLimit(t *testing.T) {
	const port = 18292

	startConnectService(t, port, iservice.Configuration{
		MaxRequestHeaderBytes: 100, // 10000-byte test header value easily exceeds 100 + slop
	})

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("http://localhost:%d/flagd.evaluation.v1.Service/ResolveBoolean", port),
		bytes.NewReader([]byte("{}")))
	require.NoError(t, err)
	// Use valid ASCII to avoid client-side rejection; value exceeds MaxHeaderBytes + slop.
	req.Header.Set("X-Large-Header", strings.Repeat("a", 10000))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusRequestHeaderFieldsTooLarge, resp.StatusCode)
}
