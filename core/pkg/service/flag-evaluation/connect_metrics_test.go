package service

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/open-feature/flagd/core/pkg/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.13.0"
	"go.uber.org/zap/zapcore"
)

func TestSetAttributes(t *testing.T) {
	tests := []struct {
		name string
		req  HTTPReqProperties
		want []attribute.KeyValue
	}{
		{
			name: "empty attributes",
			req: HTTPReqProperties{
				Service: "",
				ID:      "",
				Method:  "",
				Code:    "",
			},
			want: []attribute.KeyValue{
				semconv.ServiceNameKey.String(""),
				semconv.HTTPURLKey.String(""),
				semconv.HTTPMethodKey.String(""),
				semconv.HTTPStatusCodeKey.String(""),
			},
		},
		{
			name: "some values",
			req: HTTPReqProperties{
				Service: "myService",
				ID:      "#123",
				Method:  "POST",
				Code:    "300",
			},
			want: []attribute.KeyValue{
				semconv.ServiceNameKey.String("myService"),
				semconv.HTTPURLKey.String("#123"),
				semconv.HTTPMethodKey.String("POST"),
				semconv.HTTPStatusCodeKey.String("300"),
			},
		},
		{
			name: "special chars",
			req: HTTPReqProperties{
				Service: "!@#$%^&*()_+|}{[];',./<>",
				ID:      "",
				Method:  "",
				Code:    "",
			},
			want: []attribute.KeyValue{
				semconv.ServiceNameKey.String("!@#$%^&*()_+|}{[];',./<>"),
				semconv.HTTPURLKey.String(""),
				semconv.HTTPMethodKey.String(""),
				semconv.HTTPStatusCodeKey.String(""),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := OTelMetricsRecorder{}
			res := rec.setAttributes(tt.req)
			if len(res) != 4 {
				t.Errorf("OTelMetricsRecorder.setAttributes() must provide 4 attributes")
			}
			for i := 0; i < 4; i++ {
				if !reflect.DeepEqual(res[i], tt.want[i]) {
					t.Errorf("attribute %d = %v, want %v", i, res[i], tt.want[i])
				}
			}
		})
	}
}

func TestMiddleware(t *testing.T) {
	const svcName = "mySvc"
	exp := metric.NewManualReader()
	l, _ := logger.NewZapLogger(zapcore.DebugLevel, "")
	m := New(middlewareConfig{
		MetricReader: exp,
		Service:      svcName,
		Logger:       logger.NewLogger(l, true),
	})
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("answer"))
	})
	svr := httptest.NewServer(Handler("id", m, handler))
	defer svr.Close()
	resp, err := http.Get(svr.URL)
	if err != nil {
		t.Errorf("Got %v", err)
	}
	_, _ = io.ReadAll(resp.Body)
	data, err := exp.Collect(context.TODO())
	if err != nil {
		t.Errorf("Got %v", err)
	}
	if len(data.ScopeMetrics) != 1 {
		t.Errorf("A single scope is expected, got %d", len(data.ScopeMetrics))
	}
	scopeMetrics := data.ScopeMetrics[0]
	if !reflect.DeepEqual(scopeMetrics.Scope.Name, svcName) {
		t.Errorf("Scope name %s, want %s", scopeMetrics.Scope.Name, svcName)
	}
	if len(scopeMetrics.Metrics) != 3 {
		t.Errorf("Expected 3 metrics, got %d", len(scopeMetrics.Metrics))
	}
}

func TestNew_AutowireOTel(t *testing.T) {
	l, _ := logger.NewZapLogger(zapcore.DebugLevel, "")
	log := logger.NewLogger(l, true)
	exp := metric.NewManualReader()
	mdw := New(middlewareConfig{
		MetricReader:       exp,
		Logger:             log,
		Service:            "mySvc",
		GroupedStatus:      false,
		DisableMeasureSize: false,
	})
	if mdw.cfg.recorder == nil {
		t.Errorf("Expected OpenTelemetry to be configured, got nil")
	}
}
