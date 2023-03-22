package otel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.13.0"
)

const svcName = "mySvc"

func TestHTTPAttributes(t *testing.T) {
	type HTTPReqProperties struct {
		Service string
		ID      string
		Method  string
		Code    string
	}

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
			rec := MetricsRecorder{}
			res := rec.HTTPAttributes(tt.req.Service, tt.req.ID, tt.req.Method, tt.req.Code)
			require.Equal(t, tt.want, res)
		})
	}
}

func TestNewOTelRecorder(t *testing.T) {
	exp := metric.NewManualReader()
	rec := NewOTelRecorder(exp, svcName)
	require.NotNil(t, rec, "Expected object to be created")
	require.NotNil(t, rec.httpRequestDurHistogram, "Expected httpRequestDurHistogram to be created")
	require.NotNil(t, rec.httpResponseSizeHistogram, "Expected httpResponseSizeHistogram to be created")
	require.NotNil(t, rec.httpRequestsInflight, "Expected httpRequestsInflight to be created")
}

func TestMetrics(t *testing.T) {
	exp := metric.NewManualReader()
	rec := NewOTelRecorder(exp, svcName)
	ctx := context.TODO()
	attrs := []attribute.KeyValue{
		semconv.ServiceNameKey.String(svcName),
	}
	const n = 5
	type MetricF func()
	tests := []struct {
		name       string
		metricFunc MetricF
	}{
		{
			name: "HTTPRequestDuration",
			metricFunc: func() {
				for i := 0; i < n; i++ {
					rec.HTTPRequestDuration(ctx, 10, attrs)
				}
			},
		},
		{
			name: "HTTPResponseSize",
			metricFunc: func() {
				for i := 0; i < n; i++ {
					rec.HTTPResponseSize(ctx, 100, attrs)
				}
			},
		},
		{
			name: "InFlightRequestStart",
			metricFunc: func() {
				for i := 0; i < n; i++ {
					rec.InFlightRequestStart(ctx, attrs)
					rec.InFlightRequestEnd(ctx, attrs)
				}
			},
		},
	}
	i := 0
	for _, tt := range tests {
		i++
		tt.metricFunc()
		data, err := exp.Collect(context.TODO())
		if err != nil {
			t.Errorf("Got %v", err)
		}
		if len(data.ScopeMetrics) != 1 {
			t.Errorf("A single scope is expected, got %d", len(data.ScopeMetrics))
		}
		scopeMetrics := data.ScopeMetrics[0]
		require.Equal(t, svcName, scopeMetrics.Scope.Name)
		require.Equal(t, i, len(scopeMetrics.Metrics))
	}
}
