package syncmetrics

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	msdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

const (
	testURI      = "grpc://features-api.example:8015"
	testSelector = "all-flags"
)

// newTestRecorder wires a Recorder against an in-memory ManualReader so tests can inspect
// the emitted data points without any Prometheus / OTLP export path in the loop.
func newTestRecorder(t *testing.T) (*Recorder, *msdk.ManualReader) {
	t.Helper()
	reader := msdk.NewManualReader()
	mp := msdk.NewMeterProvider(msdk.WithReader(reader))
	return NewRecorder(mp), reader
}

func collect(t *testing.T, reader *msdk.ManualReader) metricdata.ResourceMetrics {
	t.Helper()
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(context.Background(), &rm))
	return rm
}

// findMetric locates a metric by fully-qualified OTel name across all scope metrics.
func findMetric(rm metricdata.ResourceMetrics, name string) *metricdata.Metrics {
	for i := range rm.ScopeMetrics {
		for j := range rm.ScopeMetrics[i].Metrics {
			if rm.ScopeMetrics[i].Metrics[j].Name == name {
				return &rm.ScopeMetrics[i].Metrics[j]
			}
		}
	}
	return nil
}

func TestNewRecorder_NilProviderIsSafe(t *testing.T) {
	// A nil MeterProvider must yield a fully-functional (noop) Recorder — required for
	// tests, library embedders, and any bootstrap path where OTel isn't wired yet.
	r := NewRecorder(nil)
	require.NotNil(t, r)
	// These must not panic even though nothing is being observed.
	r.RecordFlagConfigReceived(context.Background(), SourceGRPC, testURI, testSelector)
}

func TestRecorder_NilReceiverIsSafe(t *testing.T) {
	var r *Recorder
	// Callers that hold a *Recorder must be able to invoke record methods unconditionally
	// even if the value was never constructed.
	r.RecordFlagConfigReceived(context.Background(), SourceGRPC, testURI, testSelector)
}

func TestRecordFlagConfigReceived_IncrementsCounterAndUpdatesTimestamp(t *testing.T) {
	r, reader := newTestRecorder(t)
	ctx := context.Background()

	before := time.Now().Unix()
	r.RecordFlagConfigReceived(ctx, SourceGRPC, testURI, testSelector)
	r.RecordFlagConfigReceived(ctx, SourceGRPC, testURI, testSelector)
	r.RecordFlagConfigReceived(ctx, SourceGRPC, testURI, testSelector)
	after := time.Now().Unix()

	rm := collect(t, reader)

	// Counter: three applies => value 3, monotonic.
	counter := findMetric(rm, metricFlagConfigReceived)
	require.NotNil(t, counter, "expected metric %s", metricFlagConfigReceived)
	cSum, ok := counter.Data.(metricdata.Sum[int64])
	require.True(t, ok, "expected Counter -> Sum[int64]")
	require.True(t, cSum.IsMonotonic, "applied counter must be monotonic")
	require.Len(t, cSum.DataPoints, 1)
	require.EqualValues(t, 3, cSum.DataPoints[0].Value)

	// Gauge: last-applied timestamp within [before, after].
	gauge := findMetric(rm, metricFlagConfigLastReceivedTimestamp)
	require.NotNil(t, gauge, "expected metric %s", metricFlagConfigLastReceivedTimestamp)
	gData, ok := gauge.Data.(metricdata.Gauge[int64])
	require.True(t, ok, "expected Int64Gauge -> Gauge[int64]")
	require.Len(t, gData.DataPoints, 1)
	ts := gData.DataPoints[0].Value
	require.GreaterOrEqual(t, ts, before, "timestamp must not predate the first record call")
	require.LessOrEqual(t, ts, after, "timestamp must not postdate the last record call")
}

func TestRecordFlagConfigReceived_EmitsDualAttributeSet(t *testing.T) {
	r, reader := newTestRecorder(t)
	r.RecordFlagConfigReceived(context.Background(), SourceGRPC, testURI, testSelector)

	rm := collect(t, reader)
	counter := findMetric(rm, metricFlagConfigReceived)
	require.NotNil(t, counter)

	dp := counter.Data.(metricdata.Sum[int64]).DataPoints[0]
	set := dp.Attributes

	// Legacy (kept for back-compat).
	assertAttr(t, set, AttrLegacyProvider, SourceGRPC)
	assertAttr(t, set, AttrLegacyURI, testURI)
	assertAttr(t, set, AttrLegacySelector, testSelector)

	// New OTel-aligned.
	assertAttr(t, set, AttrType, SourceGRPC)
	assertAttr(t, set, AttrURI, testURI)
	assertAttr(t, set, AttrSelector, testSelector)
}

func TestRecordFlagConfigReceived_OmitsSelectorAttributesWhenEmpty(t *testing.T) {
	r, reader := newTestRecorder(t)
	r.RecordFlagConfigReceived(context.Background(), SourceHTTP, testURI, "")

	rm := collect(t, reader)
	counter := findMetric(rm, metricFlagConfigReceived)
	require.NotNil(t, counter)
	set := counter.Data.(metricdata.Sum[int64]).DataPoints[0].Attributes

	// Selector attrs must NOT be present.
	_, hasLegacy := set.Value(AttrLegacySelector)
	require.False(t, hasLegacy, "expected no %s attribute when selector is empty", AttrLegacySelector)
	_, hasNew := set.Value(AttrSelector)
	require.False(t, hasNew, "expected no %s attribute when selector is empty", AttrSelector)

	// URI + type/provider still present.
	assertAttr(t, set, AttrType, SourceHTTP)
	assertAttr(t, set, AttrLegacyProvider, SourceHTTP)
}

func TestRecordFlagConfigReceived_MultipleSourcesEmitDistinctSeries(t *testing.T) {
	r, reader := newTestRecorder(t)
	ctx := context.Background()

	// Two sources, three applies each. Expect two distinct series on both instruments.
	for i := 0; i < 3; i++ {
		r.RecordFlagConfigReceived(ctx, SourceGRPC, "grpc://server-a:8015", "")
		r.RecordFlagConfigReceived(ctx, SourceHTTP, "http://server-b/flags", "")
	}

	rm := collect(t, reader)
	counter := findMetric(rm, metricFlagConfigReceived)
	require.NotNil(t, counter)
	cSum := counter.Data.(metricdata.Sum[int64])
	require.Len(t, cSum.DataPoints, 2, "expected one data point per source type")

	// Each series should have value 3.
	for _, dp := range cSum.DataPoints {
		require.EqualValues(t, 3, dp.Value)
	}
}

// assertAttr fails the test if the attribute set doesn't contain the expected key/value.
func assertAttr(t *testing.T, set attribute.Set, key, expected string) {
	t.Helper()
	v, ok := set.Value(attribute.Key(key))
	require.True(t, ok, "expected attribute %q", key)
	require.Equal(t, expected, v.AsString(), "attribute %q value", key)
}

func TestSanitizeURI(t *testing.T) {
	cases := []struct {
		name, in, want string
	}{
		{"empty", "", ""},
		{"clean grpc", "grpc://sync.internal:9090", "grpc://sync.internal:9090"},
		{"clean http", "https://api.example.com/flags", "https://api.example.com/flags"},
		{"clean file path", "file:./demo.flagd.json", "file:./demo.flagd.json"},

		{"strips query token", "https://api.example.com/flags?token=abc123", "https://api.example.com/flags"},
		{"strips signed s3 url",
			"https://bucket.s3.amazonaws.com/flags.json?X-Amz-Signature=deadbeef&X-Amz-Credential=AKIAEXAMPLE",
			"https://bucket.s3.amazonaws.com/flags.json"},
		{"strips userinfo", "https://user:pass@host/flags", "https://host/flags"},
		{"strips userinfo (username only)", "https://user@host/flags", "https://host/flags"},
		{"strips fragment", "grpc://sync.internal:9090#region-eu", "grpc://sync.internal:9090"},
		{"strips everything combined",
			"https://user:pass@api.example.com/flags?token=abc#top",
			"https://api.example.com/flags"},

		// Fail-closed: url.Parse rejects these, so we return "" rather than a
		// partially-sanitized string that could still preserve userinfo.
		{"malformed escape with userinfo fails closed",
			"https://user:hunter2@api.example.com/flags%ZZ",
			""},
		{"control character fails closed",
			"https://user:hunter2@api.example.com/flags\x7f",
			""},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := SanitizeURI(tc.in)
			require.Equal(t, tc.want, got)
			// Extra guard: even on the fallback path, no known secret marker
			// from a malformed input must survive into the output.
			require.NotContains(t, got, "hunter2", "userinfo password leaked from %q", tc.in)
		})
	}
}

// End-to-end: a secret-bearing URI must not survive into either attribute key.
func TestRecordFlagConfigReceived_RedactsSecretsInURI(t *testing.T) {
	r, reader := newTestRecorder(t)
	ctx := context.Background()

	const secretURI = "https://user:hunter2@api.example.com/flags?token=SECRET_TOKEN"
	const cleanURI = "https://api.example.com/flags"

	r.RecordFlagConfigReceived(ctx, SourceHTTP, secretURI, "")

	rm := collect(t, reader)
	counter := findMetric(rm, metricFlagConfigReceived)
	require.NotNil(t, counter)
	dps := counter.Data.(metricdata.Sum[int64]).DataPoints
	require.Len(t, dps, 1)

	set := dps[0].Attributes
	// Both attribute keys must carry the sanitized URI.
	assertAttr(t, set, AttrURI, cleanURI)
	assertAttr(t, set, AttrLegacyURI, cleanURI)

	// And no attribute value in the whole set should contain any of the secret tokens.
	for _, kv := range set.ToSlice() {
		require.NotContains(t, kv.Value.AsString(), "hunter2", "userinfo password leaked to attribute %q", kv.Key)
		require.NotContains(t, kv.Value.AsString(), "SECRET_TOKEN", "query token leaked to attribute %q", kv.Key)
	}
}
