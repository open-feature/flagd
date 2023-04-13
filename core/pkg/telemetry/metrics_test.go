package telemetry

import (
	"context"
	"fmt"
	"testing"

	"go.opentelemetry.io/otel/sdk/resource"

	"go.opentelemetry.io/otel/sdk/metric/metricdata"

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
	rs := resource.NewWithAttributes("testSchema")
	rec := NewOTelRecorder(exp, rs, svcName)
	require.NotNil(t, rec, "Expected object to be created")
	require.NotNil(t, rec.httpRequestDurHistogram, "Expected httpRequestDurHistogram to be created")
	require.NotNil(t, rec.httpResponseSizeHistogram, "Expected httpResponseSizeHistogram to be created")
	require.NotNil(t, rec.httpRequestsInflight, "Expected httpRequestsInflight to be created")
}

func TestMetrics(t *testing.T) {
	attrs := []attribute.KeyValue{
		semconv.ServiceNameKey.String(svcName),
	}
	const n = 5
	type MetricF func(exp metric.Reader)
	tests := []struct {
		name       string
		metricFunc MetricF
		metricsLen int
	}{
		{
			name: "HTTPRequestDuration",
			metricFunc: func(exp metric.Reader) {
				rs := resource.NewWithAttributes("testSchema")
				rec := NewOTelRecorder(exp, rs, svcName)
				for i := 0; i < n; i++ {
					rec.HTTPRequestDuration(context.TODO(), 10, attrs)
				}
			},
			metricsLen: 1,
		},
		{
			name: "HTTPResponseSize",
			metricFunc: func(exp metric.Reader) {
				rs := resource.NewWithAttributes("testSchema")
				rec := NewOTelRecorder(exp, rs, svcName)
				for i := 0; i < n; i++ {
					rec.HTTPResponseSize(context.TODO(), 100, attrs)
				}
			},
			metricsLen: 1,
		},
		{
			name: "InFlightRequestStart",
			metricFunc: func(exp metric.Reader) {
				rs := resource.NewWithAttributes("testSchema")
				rec := NewOTelRecorder(exp, rs, svcName)
				ctx := context.TODO()
				for i := 0; i < n; i++ {
					rec.InFlightRequestStart(ctx, attrs)
					rec.InFlightRequestEnd(ctx, attrs)
				}
			},
			metricsLen: 1,
		},
		{
			name: "Impressions",
			metricFunc: func(exp metric.Reader) {
				rs := resource.NewWithAttributes("testSchema")
				rec := NewOTelRecorder(exp, rs, svcName)
				for i := 0; i < n; i++ {
					rec.Impressions(context.TODO(), "reason", "variant", "key")
				}
			},
			metricsLen: 1,
		},
		{
			name: "Reasons",
			metricFunc: func(exp metric.Reader) {
				rs := resource.NewWithAttributes("testSchema")
				rec := NewOTelRecorder(exp, rs, svcName)
				for i := 0; i < n; i++ {
					rec.Reasons(context.TODO(), "reason", nil)
				}
				for i := 0; i < n; i++ {
					rec.Reasons(context.TODO(), "error", fmt.Errorf("err not found"))
				}
			},
			metricsLen: 1,
		},
		{
			name: "RecordEvaluations",
			metricFunc: func(exp metric.Reader) {
				rs := resource.NewWithAttributes("testSchema")
				rec := NewOTelRecorder(exp, rs, svcName)
				for i := 0; i < n; i++ {
					rec.RecordEvaluation(context.TODO(), nil, "reason", "variant", "key")
				}
				for i := 0; i < n; i++ {
					rec.RecordEvaluation(context.TODO(), fmt.Errorf("general"), "error", "variant", "key")
				}
				for i := 0; i < n; i++ {
					rec.RecordEvaluation(context.TODO(), fmt.Errorf("not found"), "error", "variant", "key")
				}
			},
			metricsLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp := metric.NewManualReader()
			tt.metricFunc(exp)
			var data metricdata.ResourceMetrics
			err := exp.Collect(context.TODO(), &data)
			if err != nil {
				t.Errorf("Got %v", err)
			}
			if len(data.ScopeMetrics) != 1 {
				t.Errorf("A single scope is expected, got %d", len(data.ScopeMetrics))
			}
			scopeMetrics := data.ScopeMetrics[0]
			require.Equal(t, svcName, scopeMetrics.Scope.Name)
			require.Equal(t, tt.metricsLen, len(scopeMetrics.Metrics))

			r := data.Resource
			require.NotEmptyf(t, r.SchemaURL(), "Expected non-empty schema for metric resource")
		})
	}
}
