package telemetry

import (
	"errors"
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

	"github.com/open-feature/flagd/core/pkg/model"
)

const svcName = "mySvc"

// Verifies the metric attribute set, including method canonicalization.
func TestHTTPAttributes(t *testing.T) {
	type HTTPReqProperties struct {
		Service string
		ID      string
		Method  string
		Code    string
	}

	want := func(service, route, method, code string) []attribute.KeyValue {
		return []attribute.KeyValue{
			semconv.ServiceNameKey.String(service),
			semconv.HTTPRouteKey.String(route),
			semconv.HTTPRequestMethodKey.String(method),
			semconv.HTTPResponseStatusCodeKey.String(code),
			semconv.URLSchemeKey.String("http"),
		}
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
			want: want("", "", "_OTHER", ""),
		},
		{
			name: "some values",
			req: HTTPReqProperties{
				Service: "myService",
				ID:      "#123",
				Method:  "POST",
				Code:    "300",
			},
			want: want("myService", "#123", "POST", "300"),
		},
		{
			name: "mixed case standard method",
			req: HTTPReqProperties{
				Method: "pOsT",
			},
			want: want("", "", "POST", ""),
		},
		{
			name: "unknown method",
			req: HTTPReqProperties{
				Method: "PROPFIND",
			},
			want: want("", "", "_OTHER", ""),
		},
		{
			name: "special chars",
			req: HTTPReqProperties{
				Service: "!@#$%^&*()_+|}{[];',./<>",
				ID:      "",
				Method:  "",
				Code:    "",
			},
			want: want("!@#$%^&*()_+|}{[];',./<>", "", "_OTHER", ""),
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
	newRec := func(exp metric.Reader) *MetricsRecorder {
		return NewOTelRecorder(exp, resource.NewWithAttributes("testSchema"), svcName)
	}
	tests := []struct {
		name       string
		metricFunc MetricF
		metricsLen int
	}{
		{
			name: "HTTPRequestDuration",
			metricFunc: func(t *testing.T, exp metric.Reader) {
				rec := newRec(exp)
				for i := 0; i < n; i++ {
					rec.HTTPRequestDuration(t.Context(), 10, attrs)
				}
			},
			metricsLen: 1,
		},
		{
			name: "HTTPResponseSize",
			metricFunc: func(t *testing.T, exp metric.Reader) {
				rec := newRec(exp)
				for i := 0; i < n; i++ {
					rec.HTTPResponseSize(t.Context(), 100, attrs)
				}
			},
			metricsLen: 1,
		},
		{
			name: "InFlightRequestStart",
			metricFunc: func(t *testing.T, exp metric.Reader) {
				rec := newRec(exp)
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
				rec := newRec(exp)
				for i := 0; i < n; i++ {
					rec.Impressions(t.Context(), "reason", "variant", "key", nil)
				}
			},
			metricsLen: 1,
		},
		{
			name: "Reasons",
			metricFunc: func(t *testing.T, exp metric.Reader) {
				rec := newRec(exp)
				for i := 0; i < n; i++ {
					rec.Reasons(t.Context(), "keyA", "reason", nil, nil)
				}
				for i := 0; i < n; i++ {
					rec.Reasons(t.Context(), "keyB", "error", nil, fmt.Errorf("err not found"))
				}
			},
			metricsLen: 1,
		},
		{
			name: "RecordEvaluations",
			metricFunc: func(t *testing.T, exp metric.Reader) {
				rec := newRec(exp)
				for i := 0; i < n; i++ {
					rec.RecordEvaluation(t.Context(), nil, "reason", "variant", "key", nil)
				}
				for i := 0; i < n; i++ {
					rec.RecordEvaluation(t.Context(), fmt.Errorf("general"), "error", "variant", "key", nil)
				}
				for i := 0; i < n; i++ {
					rec.RecordEvaluation(t.Context(), fmt.Errorf("not found"), "error", "variant", "key", nil)
				}
			},
			metricsLen: 2,
		},
		{
			name: "SyncServerActiveStreams",
			metricFunc: func(t *testing.T, exp metric.Reader) {
				rec := newRec(exp)
				ctx := t.Context()
				for i := 0; i < n; i++ {
					rec.SyncServerStreamStart(ctx, attrs)
					rec.SyncServerStreamEnd(ctx, attrs)
				}
			},
			metricsLen: 1,
		},
		{
			name: "SyncServerStreamDuration",
			metricFunc: func(t *testing.T, exp metric.Reader) {
				rec := newRec(exp)
				for i := 0; i < n; i++ {
					rec.SyncServerStreamDuration(t.Context(), 100*time.Millisecond, attrs)
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
func TestNoopMetricsRecorder(t *testing.T) {
	no := NoopMetricsRecorder{}

	t.Run("HTTPAttributes", func(t *testing.T) {
		got := no.HTTPAttributes("", "", "", "", "")
		require.Empty(t, got)
	})

	t.Run("HTTPRequestDuration", func(t *testing.T) {
		no.HTTPRequestDuration(t.Context(), 0, nil)
	})

	t.Run("InFlightRequestStart", func(t *testing.T) {
		no.InFlightRequestStart(t.Context(), nil)
	})

	t.Run("InFlightRequestEnd", func(t *testing.T) {
		no.InFlightRequestEnd(t.Context(), nil)
	})

	t.Run("RecordEvaluation", func(t *testing.T) {
		no.RecordEvaluation(t.Context(), nil, "", "", "", nil)
	})

	t.Run("Impressions", func(t *testing.T) {
		no.Impressions(t.Context(), "", "", "", nil)
	})

	t.Run("SyncServerStreamStart", func(t *testing.T) {
		no.SyncServerStreamStart(t.Context(), nil)
	})

	t.Run("SyncServerStreamEnd", func(t *testing.T) {
		no.SyncServerStreamEnd(t.Context(), nil)
	})

	t.Run("SyncServerStreamDuration", func(t *testing.T) {
		no.SyncServerStreamDuration(t.Context(), 0, nil)
	})
}

func TestRecordEvaluationFlagSetID(t *testing.T) {
	exp := metric.NewManualReader()
	rec := NewOTelRecorder(exp, resource.NewWithAttributes("testSchema"), svcName)

	rec.RecordEvaluation(t.Context(), nil, "STATIC", "on", "flag-a", model.Metadata{"flagSetId": "set-a"})
	rec.RecordEvaluation(t.Context(), nil, "STATIC", "on", "flag-b", model.Metadata{"flagSetId": "set-b"})
	rec.RecordEvaluation(
		t.Context(), errors.New(model.FlagNotFoundErrorCode), "ERROR", "", "missing",
		model.Metadata{"flagSetId": "caller-controlled"},
	)
	rec.RecordEvaluation(
		t.Context(), errors.New(model.GeneralErrorCode), "ERROR", "", "found-error",
		model.Metadata{"flagSetId": "configured-set"},
	)
	rec.RecordEvaluation(t.Context(), nil, "STATIC", "on", "without-set", model.Metadata{"flagSetId": 42})

	var data metricdata.ResourceMetrics
	require.NoError(t, exp.Collect(t.Context(), &data))
	require.Len(t, data.ScopeMetrics, 1)

	want := map[string]map[string]bool{
		impressionMetric: {
			"set-a": true,
			"set-b": true,
			"":      true,
		},
		reasonMetric: {
			"set-a":          true,
			"set-b":          true,
			"configured-set": true,
			"":               true,
		},
	}

	for _, m := range data.ScopeMetrics[0].Metrics {
		expected, ok := want[m.Name]
		if !ok {
			continue
		}
		sum, ok := m.Data.(metricdata.Sum[int64])
		require.True(t, ok)
		got := make(map[string]bool, len(sum.DataPoints))
		for _, point := range sum.DataPoints {
			value, present := point.Attributes.Value(semconv.FeatureFlagSetIDKey)
			if present {
				got[value.AsString()] = true
			} else {
				got[""] = true
			}
		}
		require.Equal(t, expected, got, m.Name)
		delete(want, m.Name)
	}
	require.Empty(t, want)
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

func TestGRPCSyncServerStreamDurationBuckets(t *testing.T) {
	testHistogramBuckets(t,
		syncServerStreamDurationMetric,
		[]float64{30, 60, 120, 300, 480, 600, 1200, 1800, 3600, 10800},
		func(t *testing.T, rec *MetricsRecorder, attrs []attribute.KeyValue) {
			rec.SyncServerStreamDuration(t.Context(), 100*time.Millisecond, attrs)
		},
		"Expected histogram buckets for long-lived sync streams (30s, 1min, 2min, 5min, 8min, 10min, 20min, 30min, 1h, 3h)",
	)
}
