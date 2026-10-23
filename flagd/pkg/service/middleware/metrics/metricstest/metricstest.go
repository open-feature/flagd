// Package metricstest provides helpers for asserting on recorded HTTP metrics in tests.
package metricstest

import (
	"testing"

	"github.com/open-feature/flagd/core/pkg/telemetry"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
)

// NewRecorder returns an in-memory recorder and reader; it sets the global meter provider until cleanup.
func NewRecorder(t *testing.T, svcName string) (telemetry.IMetricsRecorder, metric.Reader) {
	t.Helper()

	prev := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(prev) })

	reader := metric.NewManualReader()
	return telemetry.NewOTelRecorder(reader, resource.NewWithAttributes("testSchema"), svcName), reader
}

// RequestDurationPoints returns the scopes that recorded http.server.request.duration, and its data points.
func RequestDurationPoints(
	t *testing.T, reader metric.Reader,
) (map[string]bool, []metricdata.HistogramDataPoint[float64]) {
	t.Helper()

	var data metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &data))

	scopes := map[string]bool{}
	var points []metricdata.HistogramDataPoint[float64]
	for _, sm := range data.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name == "http.server.request.duration" {
				scopes[sm.Scope.Name] = true
				points = append(points, m.Data.(metricdata.Histogram[float64]).DataPoints...)
			}
		}
	}
	return scopes, points
}
