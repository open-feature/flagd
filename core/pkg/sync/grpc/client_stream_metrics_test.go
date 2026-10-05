package grpc

import (
	"context"
	"testing"

	"github.com/open-feature/flagd/core/pkg/sync/syncmetrics"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	msdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

const (
	testURI      = "grpc://features-api.example:8015"
	testSelector = "all-flags"
)

func newTestClientStreamMetrics(t *testing.T) (*clientStreamMetrics, *msdk.ManualReader) {
	t.Helper()
	reader := msdk.NewManualReader()
	mp := msdk.NewMeterProvider(msdk.WithReader(reader))
	return newClientStreamMetrics(mp), reader
}

func collectClientStream(t *testing.T, reader *msdk.ManualReader) metricdata.ResourceMetrics {
	t.Helper()
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(context.Background(), &rm))
	return rm
}

func findClientStreamMetric(rm metricdata.ResourceMetrics, name string) *metricdata.Metrics {
	for i := range rm.ScopeMetrics {
		for j := range rm.ScopeMetrics[i].Metrics {
			if rm.ScopeMetrics[i].Metrics[j].Name == name {
				return &rm.ScopeMetrics[i].Metrics[j]
			}
		}
	}
	return nil
}

func TestNewClientStreamMetrics_NilProviderIsSafe(t *testing.T) {
	m := newClientStreamMetrics(nil)
	require.NotNil(t, m)
	// Method calls must be safe even when the recorder was built without a MeterProvider.
	m.recordStreamOpened(context.Background(), testURI, testSelector)
	m.recordStreamClosed(context.Background(), testURI, testSelector)
	m.recordReconnect(context.Background(), testURI, testSelector)
}

func TestClientStreamMetrics_NilReceiverIsSafe(t *testing.T) {
	var m *clientStreamMetrics
	m.recordStreamOpened(context.Background(), testURI, testSelector)
	m.recordStreamClosed(context.Background(), testURI, testSelector)
	m.recordReconnect(context.Background(), testURI, testSelector)
}

func TestClientStreamMetrics_StreamOpenCloseSymmetry(t *testing.T) {
	m, reader := newTestClientStreamMetrics(t)
	ctx := context.Background()

	// 3 opens - 2 closes = 1 stream currently open.
	m.recordStreamOpened(ctx, testURI, testSelector)
	m.recordStreamOpened(ctx, testURI, testSelector)
	m.recordStreamOpened(ctx, testURI, testSelector)
	m.recordStreamClosed(ctx, testURI, testSelector)
	m.recordStreamClosed(ctx, testURI, testSelector)

	rm := collectClientStream(t, reader)
	metric := findClientStreamMetric(rm, metricClientStreamActive)
	require.NotNil(t, metric, "expected metric %s", metricClientStreamActive)

	sum, ok := metric.Data.(metricdata.Sum[int64])
	require.True(t, ok, "expected UpDownCounter -> Sum[int64]")
	require.False(t, sum.IsMonotonic, "stream.active is an UpDownCounter, not monotonic")
	require.Len(t, sum.DataPoints, 1)
	require.EqualValues(t, 1, sum.DataPoints[0].Value)
}

func TestClientStreamMetrics_ReconnectsAreMonotonic(t *testing.T) {
	m, reader := newTestClientStreamMetrics(t)
	ctx := context.Background()

	for i := 0; i < 4; i++ {
		m.recordReconnect(ctx, testURI, testSelector)
	}

	rm := collectClientStream(t, reader)
	metric := findClientStreamMetric(rm, metricClientStreamReconnects)
	require.NotNil(t, metric, "expected metric %s", metricClientStreamReconnects)

	sum, ok := metric.Data.(metricdata.Sum[int64])
	require.True(t, ok, "expected Counter -> Sum[int64]")
	require.True(t, sum.IsMonotonic, "reconnects counter must be monotonic")
	require.Len(t, sum.DataPoints, 1)
	require.EqualValues(t, 4, sum.DataPoints[0].Value)
}

func TestClientStreamMetrics_EmitsDualAttributeSet(t *testing.T) {
	m, reader := newTestClientStreamMetrics(t)
	m.recordStreamOpened(context.Background(), testURI, testSelector)

	rm := collectClientStream(t, reader)
	metric := findClientStreamMetric(rm, metricClientStreamActive)
	require.NotNil(t, metric)
	set := metric.Data.(metricdata.Sum[int64]).DataPoints[0].Attributes

	// Legacy (kept for back-compat) — source hard-wired to "grpc" in this package.
	assertClientStreamAttr(t, set, syncmetrics.AttrLegacyProvider, syncmetrics.SourceGRPC)
	assertClientStreamAttr(t, set, syncmetrics.AttrLegacyURI, testURI)
	assertClientStreamAttr(t, set, syncmetrics.AttrLegacySelector, testSelector)

	// New OTel-aligned.
	assertClientStreamAttr(t, set, syncmetrics.AttrType, syncmetrics.SourceGRPC)
	assertClientStreamAttr(t, set, syncmetrics.AttrURI, testURI)
	assertClientStreamAttr(t, set, syncmetrics.AttrSelector, testSelector)
}

func TestClientStreamMetrics_SelectorOmittedWhenEmpty(t *testing.T) {
	m, reader := newTestClientStreamMetrics(t)
	m.recordReconnect(context.Background(), testURI, "")

	rm := collectClientStream(t, reader)
	metric := findClientStreamMetric(rm, metricClientStreamReconnects)
	require.NotNil(t, metric)
	set := metric.Data.(metricdata.Sum[int64]).DataPoints[0].Attributes

	_, hasLegacy := set.Value(attribute.Key(syncmetrics.AttrLegacySelector))
	require.False(t, hasLegacy, "expected no %s when selector is empty", syncmetrics.AttrLegacySelector)
	_, hasNew := set.Value(attribute.Key(syncmetrics.AttrSelector))
	require.False(t, hasNew, "expected no %s when selector is empty", syncmetrics.AttrSelector)
}

func assertClientStreamAttr(t *testing.T, set attribute.Set, key, expected string) {
	t.Helper()
	v, ok := set.Value(attribute.Key(key))
	require.True(t, ok, "expected attribute %q", key)
	require.Equal(t, expected, v.AsString(), "attribute %q value", key)
}

// Mirrors the syncmetrics-side redaction test for the stream-lifecycle attribute set.
func TestClientStreamMetrics_RedactsSecretsInURI(t *testing.T) {
	m, reader := newTestClientStreamMetrics(t)

	const secretURI = "grpc://user:hunter2@sync.internal:9090?token=SECRET_TOKEN"
	const cleanURI = "grpc://sync.internal:9090"

	m.recordStreamOpened(context.Background(), secretURI, testSelector)

	rm := collectClientStream(t, reader)
	metric := findClientStreamMetric(rm, metricClientStreamActive)
	require.NotNil(t, metric)
	set := metric.Data.(metricdata.Sum[int64]).DataPoints[0].Attributes

	assertClientStreamAttr(t, set, syncmetrics.AttrURI, cleanURI)
	assertClientStreamAttr(t, set, syncmetrics.AttrLegacyURI, cleanURI)

	for _, kv := range set.ToSlice() {
		require.NotContains(t, kv.Value.AsString(), "hunter2", "userinfo password leaked to %q", kv.Key)
		require.NotContains(t, kv.Value.AsString(), "SECRET_TOKEN", "query token leaked to %q", kv.Key)
	}
}
