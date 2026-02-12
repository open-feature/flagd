package telemetry

import (
	"fmt"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
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
				semconv.HTTPRouteKey.String(""),
				semconv.HTTPRequestMethodKey.String(""),
				semconv.HTTPResponseStatusCodeKey.String(""),
				semconv.URLSchemeKey.String("http"),
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
				semconv.HTTPRouteKey.String("#123"),
				semconv.HTTPRequestMethodKey.String("POST"),
				semconv.HTTPResponseStatusCodeKey.String("300"),
				semconv.URLSchemeKey.String("http"),
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
				semconv.HTTPRouteKey.String(""),
				semconv.HTTPRequestMethodKey.String(""),
				semconv.HTTPResponseStatusCodeKey.String(""),
				semconv.URLSchemeKey.String("http"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := MetricsRecorder{}
			res := rec.HTTPAttributes(tt.req.Service, tt.req.ID, tt.req.Method, tt.req.Code, "http")
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
	type MetricF func(t *testing.T, exp metric.Reader)
	tests := []struct {
		name       string
		metricFunc MetricF
		metricsLen int
	}{
		{
			name: "HTTPRequestDuration",
			metricFunc: func(t *testing.T, exp metric.Reader) {
				rs := resource.NewWithAttributes("testSchema")
				rec := NewOTelRecorder(exp, rs, svcName)
				for i := 0; i < n; i++ {
					rec.HTTPRequestDuration(t.Context(), 10, attrs)
				}
			},
			metricsLen: 1,
		},
		{
			name: "HTTPResponseSize",
			metricFunc: func(t *testing.T, exp metric.Reader) {
				rs := resource.NewWithAttributes("testSchema")
				rec := NewOTelRecorder(exp, rs, svcName)
				for i := 0; i < n; i++ {
					rec.HTTPResponseSize(t.Context(), 100, attrs)
				}
			},
			metricsLen: 1,
		},
		{
			name: "InFlightRequestStart",
			metricFunc: func(t *testing.T, exp metric.Reader) {
				rs := resource.NewWithAttributes("testSchema")
				rec := NewOTelRecorder(exp, rs, svcName)
				ctx := t.Context()
				for i := 0; i < n; i++ {
					rec.InFlightRequestStart(ctx, attrs)
					rec.InFlightRequestEnd(ctx, attrs)
				}
			},
			metricsLen: 1,
		},
		{
			name: "Impressions",
			metricFunc: func(t *testing.T, exp metric.Reader) {
				rs := resource.NewWithAttributes("testSchema")
				rec := NewOTelRecorder(exp, rs, svcName)
				for i := 0; i < n; i++ {
					rec.Impressions(t.Context(), "reason", "variant", "key")
				}
			},
			metricsLen: 1,
		},
		{
			name: "Reasons",
			metricFunc: func(t *testing.T, exp metric.Reader) {
				rs := resource.NewWithAttributes("testSchema")
				rec := NewOTelRecorder(exp, rs, svcName)
				for i := 0; i < n; i++ {
					rec.Reasons(t.Context(), "keyA", "reason", nil)
				}
				for i := 0; i < n; i++ {
					rec.Reasons(t.Context(), "keyB", "error", fmt.Errorf("err not found"))
				}
			},
			metricsLen: 1,
		},
		{
			name: "RecordEvaluations",
			metricFunc: func(t *testing.T, exp metric.Reader) {
				rs := resource.NewWithAttributes("testSchema")
				rec := NewOTelRecorder(exp, rs, svcName)
				for i := 0; i < n; i++ {
					rec.RecordEvaluation(t.Context(), nil, "reason", "variant", "key")
				}
				for i := 0; i < n; i++ {
					rec.RecordEvaluation(t.Context(), fmt.Errorf("general"), "error", "variant", "key")
				}
				for i := 0; i < n; i++ {
					rec.RecordEvaluation(t.Context(), fmt.Errorf("not found"), "error", "variant", "key")
				}
			},
			metricsLen: 2,
		},
		{
			name: "SyncActiveStreams",
			metricFunc: func(t *testing.T, exp metric.Reader) {
				rs := resource.NewWithAttributes("testSchema")
				rec := NewOTelRecorder(exp, rs, svcName)
				ctx := t.Context()
				for i := 0; i < n; i++ {
					rec.SyncStreamStart(ctx, attrs)
					rec.SyncStreamEnd(ctx, attrs)
				}
			},
			metricsLen: 1,
		},
		{
			name: "SyncStreamDuration",
			metricFunc: func(t *testing.T, exp metric.Reader) {
				rs := resource.NewWithAttributes("testSchema")
				rec := NewOTelRecorder(exp, rs, svcName)
				for i := 0; i < n; i++ {
					rec.SyncStreamDuration(t.Context(), 100*time.Millisecond, attrs)
				}
			},
			metricsLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp := metric.NewManualReader()
			tt.metricFunc(t, exp)
			var data metricdata.ResourceMetrics
			err := exp.Collect(t.Context(), &data)
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

// some really simple tests just to make sure all methods are actually implemented and nothing panics
func TestNoopMetricsRecorder_HTTPAttributes(t *testing.T) {
	no := NoopMetricsRecorder{}
	got := no.HTTPAttributes("", "", "", "", "")
	require.Empty(t, got)
}

func TestNoopMetricsRecorder_HTTPRequestDuration(t *testing.T) {
	no := NoopMetricsRecorder{}
	no.HTTPRequestDuration(t.Context(), 0, nil)
}

func TestNoopMetricsRecorder_InFlightRequestStart(t *testing.T) {
	no := NoopMetricsRecorder{}
	no.InFlightRequestStart(t.Context(), nil)
}

func TestNoopMetricsRecorder_InFlightRequestEnd(t *testing.T) {
	no := NoopMetricsRecorder{}
	no.InFlightRequestEnd(t.Context(), nil)
}

func TestNoopMetricsRecorder_RecordEvaluation(t *testing.T) {
	no := NoopMetricsRecorder{}
	no.RecordEvaluation(t.Context(), nil, "", "", "")
}

func TestNoopMetricsRecorder_Impressions(t *testing.T) {
	no := NoopMetricsRecorder{}
	no.Impressions(t.Context(), "", "", "")
}

func TestNoopMetricsRecorderSyncStreamStart(t *testing.T) {
	no := NoopMetricsRecorder{}
	no.SyncStreamStart(t.Context(), nil)
}

func TestNoopMetricsRecorderSyncStreamEnd(t *testing.T) {
	no := NoopMetricsRecorder{}
	no.SyncStreamEnd(t.Context(), nil)
}

func TestNoopMetricsRecorderSyncStreamDuration(t *testing.T) {
	no := NoopMetricsRecorder{}
	no.SyncStreamDuration(t.Context(), 0, nil)
}

// testHistogramBuckets is a helper function that tests histogram bucket configuration
func testHistogramBuckets(t *testing.T, metricName string, expectedBounds []float64, recordMetric func(t *testing.T, rec *MetricsRecorder, attrs []attribute.KeyValue), assertMsg string) {
	t.Helper()
	const testSvcName = "testService"
	exp := metric.NewManualReader()
	rs := resource.NewWithAttributes("testSchema")
	rec := NewOTelRecorder(exp, rs, testSvcName)

	attrs := []attribute.KeyValue{
		semconv.ServiceNameKey.String(testSvcName),
	}
	recordMetric(t, rec, attrs)

	var data metricdata.ResourceMetrics
	err := exp.Collect(t.Context(), &data)
	require.NoError(t, err)

	require.Len(t, data.ScopeMetrics, 1)
	scopeMetrics := data.ScopeMetrics[0]
	require.Equal(t, testSvcName, scopeMetrics.Scope.Name)

	var foundHistogram bool
	for _, m := range scopeMetrics.Metrics {
		if m.Name == metricName {
			histogram, ok := m.Data.(metricdata.Histogram[float64])
			require.True(t, ok, "Expected metric to be a Histogram")

			require.NotEmpty(t, histogram.DataPoints, "Expected at least one data point")
			require.Equal(t, expectedBounds, histogram.DataPoints[0].Bounds, assertMsg)
			foundHistogram = true
			break
		}
	}
	require.Truef(t, foundHistogram, "Expected to find %s histogram", metricName)
}

func TestHTTPRequestDurationBuckets(t *testing.T) {
	testHistogramBuckets(t,
		httpRequestDurationMetric,
		prometheus.DefBuckets,
		func(t *testing.T, rec *MetricsRecorder, attrs []attribute.KeyValue) {
			rec.HTTPRequestDuration(t.Context(), 100*time.Millisecond, attrs)
		},
		"Expected histogram buckets to match prometheus.DefBuckets",
	)
}

func TestHTTPResponseSizeBuckets(t *testing.T) {
	testHistogramBuckets(t,
		httpResponseSizeMetric,
		prometheus.ExponentialBuckets(100, 10, 8),
		func(t *testing.T, rec *MetricsRecorder, attrs []attribute.KeyValue) {
			rec.HTTPResponseSize(t.Context(), 500, attrs)
		},
		"Expected histogram buckets to match exponential buckets (100, 10, 8)",
	)
}

func TestGRPCSyncStreamDurationBuckets(t *testing.T) {
	testHistogramBuckets(t,
		syncStreamDurationMetric,
		[]float64{30, 60, 120, 300, 480, 600, 1200, 1800, 3600, 10800},
		func(t *testing.T, rec *MetricsRecorder, attrs []attribute.KeyValue) {
			rec.SyncStreamDuration(t.Context(), 100*time.Millisecond, attrs)
		},
		"Expected histogram buckets for long-lived sync streams (30s, 1min, 2min, 5min, 8min, 10min, 20min, 30min, 1h, 3h)",
	)
}
