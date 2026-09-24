// Package syncmetrics defines source-agnostic client-side OpenTelemetry metrics for flagd's
// sync channel — the direction where flagd pulls flag configurations from an upstream source
// (flag-server, flagd-proxy, HTTP endpoint, file etc.).
//
// Every sync provider records into the same instruments; the sync source is identified via
// attributes. Provider-specific instruments (e.g. gRPC stream lifecycle) live in the
// respective provider's own package alongside their implementations.
//
// The stream-lifecycle instruments (open/close, reconnects) are intentionally excluded from
// this package because they don't apply to every sync source — a polling HTTP client or a
// filesystem watcher has no "stream" to open. Those live under
// core/pkg/sync/grpc/ (and, later, kubernetes/ for the watch API).
package syncmetrics

import (
	"context"
	"net/url"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

// Source identifies which sync provider produced an observation. Values are lowercase and
// match the config-time provider names used by SyncBuilder so operators see the same term
// on their dashboards as in their flagd config.
const (
	SourceGRPC       = "grpc"
	SourceHTTP       = "http"
	SourceFile       = "file"
	SourceKubernetes = "kubernetes"
	SourceBlob       = "blob"
)

// Instrumentation scope for the OTel meter. Follows the "Go import path of the instrumentation"
// convention so `otel_scope_name` uniquely identifies where these metrics come from.
const meterName = "github.com/open-feature/flagd/core/pkg/sync/syncmetrics"

// Metric names — the canonical (dotted) OTel names. The Prometheus exporter converts "."
// to "_" in label names and drops select unit suffixes depending on collector config.
const (
	metricFlagConfigReceived              = "feature_flag.flagd.sync.client.flag_config.received"
	metricFlagConfigLastReceivedTimestamp = "feature_flag.flagd.sync.client.flag_config.last_received_timestamp"
)

// Attribute keys. Two sets are emitted on every data point so existing dashboards keep
// working during the migration:
//
//   - flagd.sync.*                    (pre-existing convention, kept for back-compat)
//   - feature_flag.flagd.sync.*       (new OTel-aligned taxonomy)
//
// Drop-old is a coordinated breaking release later, not part of this rollout.
const (
	AttrLegacyProvider = "flagd.sync.provider"
	AttrLegacyURI      = "flagd.sync.uri"
	AttrLegacySelector = "flagd.sync.selector"

	AttrType     = "feature_flag.flagd.sync.type"
	AttrURI      = "feature_flag.flagd.sync.uri"
	AttrSelector = "feature_flag.flagd.sync.selector"
)

// Recorder owns the universal client-side sync instruments. One instance is shared across
// all sync providers in a flagd binary; providers pass their own source type / uri / selector
// on each record call.
//
// Recorder methods are nil-receiver safe — callers that don't wire a Recorder (tests,
// library embedders that don't configure OTel) simply get no-ops without needing to check.
type Recorder struct {
	flagConfigReceived    metric.Int64Counter
	lastReceivedTimestamp metric.Int64Gauge
}

// NewRecorder builds a Recorder against the given MeterProvider. A nil MeterProvider falls
// back to the global one set by telemetry.NewOTelRecorder (via otel.SetMeterProvider), and
// if no global provider has been registered either it degrades to a noop MeterProvider —
// so the returned Recorder is always safe to use.
func NewRecorder(mp metric.MeterProvider) *Recorder {
	if mp == nil {
		mp = otel.GetMeterProvider()
	}
	if mp == nil {
		mp = noop.NewMeterProvider()
	}
	m := mp.Meter(meterName)

	// The gauge carries a real unit ("s"); the counter deliberately omits WithUnit
	// so the OTel Prometheus exporter doesn't mangle an annotation unit into the name.
	received, _ := m.Int64Counter(
		metricFlagConfigReceived,
		metric.WithDescription(
			"Total number of flag-configuration payloads received from a sync source and "+
				"handed off downstream. Does not imply the payload was successfully applied "+
				"to the evaluator's store — see the runtime-level apply metric for that."),
	)
	lastReceivedTs, _ := m.Int64Gauge(
		metricFlagConfigLastReceivedTimestamp,
		metric.WithDescription(
			"Unix timestamp (seconds) of the most recent flag-configuration payload received "+
				"from a sync source. Query staleness as time() - <value>."),
		metric.WithUnit("s"),
	)

	return &Recorder{
		flagConfigReceived:    received,
		lastReceivedTimestamp: lastReceivedTs,
	}
}

// RecordFlagConfigReceived bumps the received counter and updates the last-received
// timestamp for the given sync source. Providers call this after a DataSync has been
// handed off downstream — the record does NOT imply the evaluator successfully applied
// the payload; use a runtime-level apply metric for that.
//
// Attributes: dual set (legacy `flagd.sync.*` and new `feature_flag.flagd.sync.*`);
// selector is omitted when empty.
func (r *Recorder) RecordFlagConfigReceived(ctx context.Context, sourceType, uri, selector string) {
	if r == nil {
		return
	}
	attrs := attributes(sourceType, uri, selector)
	if r.flagConfigReceived != nil {
		r.flagConfigReceived.Add(ctx, 1, attrs)
	}
	if r.lastReceivedTimestamp != nil {
		r.lastReceivedTimestamp.Record(ctx, time.Now().Unix(), attrs)
	}
}

// SanitizeURI strips userinfo, query, and fragment from a sync-source URI so
// token- or credential-bearing URIs don't land as Prometheus attribute values.
// Fails closed: any input url.Parse rejects returns "" rather than a partially
// sanitized string that could still preserve userinfo.
func SanitizeURI(raw string) string {
	if raw == "" {
		return ""
	}
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	u.User = nil
	u.RawQuery = ""
	u.Fragment = ""
	return u.String()
}

// attributes builds the dual attribute set for a sync-source observation. Kept as a
// free function (not a method) so it's cheap to test in isolation without wiring a full
// Recorder.
func attributes(sourceType, uri, selector string) metric.MeasurementOption {
	safeURI := SanitizeURI(uri)
	kvs := make([]attribute.KeyValue, 0, 6)
	kvs = append(kvs,
		attribute.String(AttrLegacyProvider, sourceType),
		attribute.String(AttrLegacyURI, safeURI),
		attribute.String(AttrType, sourceType),
		attribute.String(AttrURI, safeURI),
	)
	if selector != "" {
		kvs = append(kvs,
			attribute.String(AttrLegacySelector, selector),
			attribute.String(AttrSelector, selector),
		)
	}
	return metric.WithAttributes(kvs...)
}
