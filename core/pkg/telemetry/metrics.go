package telemetry

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	msdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

const (
	ProviderName = "flagd"

	FeatureFlagReasonKey = attribute.Key("feature_flag.reason")
	ExceptionTypeKey     = attribute.Key("ExceptionTypeKeyName")

	httpRequestDurationMetric = "http.server.request.duration"
	httpResponseSizeMetric    = "http.server.response.body.size"
	httpActiveRequestsMetric  = "http.server.active_requests"
	impressionMetric          = "feature_flag." + ProviderName + ".impression"
	reasonMetric              = "feature_flag." + ProviderName + ".result.reason"

	grpcSyncActiveStreamsMetric  = "grpc.sync.active_streams"
	grpcSyncStreamDurationMetric = "grpc.sync.stream.duration"
	grpcSyncEventsSentMetric     = "grpc.sync.events.sent"
	grpcFetchAllFlagsMetric      = "grpc.sync.fetchall.requests"
)

type IMetricsRecorder interface {
	HTTPAttributes(svcName, url, method, code, scheme string) []attribute.KeyValue
	HTTPRequestDuration(ctx context.Context, duration time.Duration, attrs []attribute.KeyValue)
	HTTPResponseSize(ctx context.Context, sizeBytes int64, attrs []attribute.KeyValue)
	InFlightRequestStart(ctx context.Context, attrs []attribute.KeyValue)
	InFlightRequestEnd(ctx context.Context, attrs []attribute.KeyValue)
	RecordEvaluation(ctx context.Context, err error, reason, variant, key string)
	Impressions(ctx context.Context, reason, variant, key string)
	// gRPC Sync metrics
	SyncStreamStart(ctx context.Context, attrs []attribute.KeyValue)
	SyncStreamEnd(ctx context.Context, attrs []attribute.KeyValue)
	SyncStreamDuration(ctx context.Context, duration time.Duration, attrs []attribute.KeyValue)
	SyncEventSent(ctx context.Context, attrs []attribute.KeyValue)
	FetchAllFlagsRequest(ctx context.Context, attrs []attribute.KeyValue)
}

type NoopMetricsRecorder struct{}

func (NoopMetricsRecorder) HTTPAttributes(_, _, _, _, _ string) []attribute.KeyValue {
	return []attribute.KeyValue{}
}

func (NoopMetricsRecorder) HTTPRequestDuration(_ context.Context, _ time.Duration, _ []attribute.KeyValue) {
}

func (NoopMetricsRecorder) HTTPResponseSize(_ context.Context, _ int64, _ []attribute.KeyValue) {
}

func (NoopMetricsRecorder) InFlightRequestStart(_ context.Context, _ []attribute.KeyValue) {
}

func (NoopMetricsRecorder) InFlightRequestEnd(_ context.Context, _ []attribute.KeyValue) {
}

func (NoopMetricsRecorder) RecordEvaluation(_ context.Context, _ error, _, _, _ string) {
}

func (NoopMetricsRecorder) Impressions(_ context.Context, _, _, _ string) {
}

func (NoopMetricsRecorder) SyncStreamStart(_ context.Context, _ []attribute.KeyValue) {
}

func (NoopMetricsRecorder) SyncStreamEnd(_ context.Context, _ []attribute.KeyValue) {
}

func (NoopMetricsRecorder) SyncStreamDuration(_ context.Context, _ time.Duration, _ []attribute.KeyValue) {
}

func (NoopMetricsRecorder) SyncEventSent(_ context.Context, _ []attribute.KeyValue) {
}

func (NoopMetricsRecorder) FetchAllFlagsRequest(_ context.Context, _ []attribute.KeyValue) {
}

type MetricsRecorder struct {
	httpRequestDurHistogram   metric.Float64Histogram
	httpResponseSizeHistogram metric.Float64Histogram
	httpRequestsInflight      metric.Int64UpDownCounter
	impressions               metric.Int64Counter
	reasons                   metric.Int64Counter
	// gRPC Sync metrics
	syncActiveStreams     metric.Int64UpDownCounter
	syncStreamDuration    metric.Float64Histogram
	syncEventsSent        metric.Int64Counter
	fetchAllFlagsRequests metric.Int64Counter
}

func (r MetricsRecorder) HTTPAttributes(svcName, url, method, code, scheme string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.ServiceNameKey.String(svcName),
		semconv.HTTPRouteKey.String(url),
		semconv.HTTPRequestMethodKey.String(method),
		semconv.HTTPResponseStatusCodeKey.String(code),
		semconv.URLSchemeKey.String(scheme),
	}
}

func (r MetricsRecorder) HTTPRequestDuration(ctx context.Context, duration time.Duration, attrs []attribute.KeyValue) {
	r.httpRequestDurHistogram.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

func (r MetricsRecorder) HTTPResponseSize(ctx context.Context, sizeBytes int64, attrs []attribute.KeyValue) {
	r.httpResponseSizeHistogram.Record(ctx, float64(sizeBytes), metric.WithAttributes(attrs...))
}

func (r MetricsRecorder) InFlightRequestStart(ctx context.Context, attrs []attribute.KeyValue) {
	r.httpRequestsInflight.Add(ctx, 1, metric.WithAttributes(attrs...))
}

func (r MetricsRecorder) InFlightRequestEnd(ctx context.Context, attrs []attribute.KeyValue) {
	r.httpRequestsInflight.Add(ctx, -1, metric.WithAttributes(attrs...))
}

func (r MetricsRecorder) RecordEvaluation(ctx context.Context, err error, reason, variant, key string) {
	if err == nil {
		r.Impressions(ctx, reason, variant, key)
	}
	r.Reasons(ctx, key, reason, err)
}

func (r MetricsRecorder) Impressions(ctx context.Context, reason, variant, key string) {
	r.impressions.Add(ctx,
		1,
		metric.WithAttributes(append(SemConvFeatureFlagAttributes(key, variant), FeatureFlagReason(reason))...))
}

func (r MetricsRecorder) Reasons(ctx context.Context, key string, reason string, err error) {
	attrs := []attribute.KeyValue{
		semconv.FeatureFlagProviderName(ProviderName),
		FeatureFlagReason(reason),
	}
	if err == nil {
		// record flag key only if evaluation is successful
		attrs = append(attrs, semconv.FeatureFlagKey(key))
	} else {
		attrs = append(attrs, ExceptionType(err.Error()))
	}

	r.reasons.Add(ctx, 1, metric.WithAttributes(attrs...))
}

func (r MetricsRecorder) SyncStreamStart(ctx context.Context, attrs []attribute.KeyValue) {
	r.syncActiveStreams.Add(ctx, 1, metric.WithAttributes(attrs...))
}

func (r MetricsRecorder) SyncStreamEnd(ctx context.Context, attrs []attribute.KeyValue) {
	r.syncActiveStreams.Add(ctx, -1, metric.WithAttributes(attrs...))
}

func (r MetricsRecorder) SyncStreamDuration(ctx context.Context, duration time.Duration, attrs []attribute.KeyValue) {
	r.syncStreamDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

func (r MetricsRecorder) SyncEventSent(ctx context.Context, attrs []attribute.KeyValue) {
	r.syncEventsSent.Add(ctx, 1, metric.WithAttributes(attrs...))
}

func (r MetricsRecorder) FetchAllFlagsRequest(ctx context.Context, attrs []attribute.KeyValue) {
	r.fetchAllFlagsRequests.Add(ctx, 1, metric.WithAttributes(attrs...))
}

func getDurationView(svcName, viewName string, bucket []float64) msdk.View {
	return msdk.NewView(
		msdk.Instrument{
			// we change aggregation only for instruments with this name and scope
			Name: viewName,
			Scope: instrumentation.Scope{
				Name: svcName,
			},
		},
		msdk.Stream{Aggregation: msdk.AggregationExplicitBucketHistogram{
			Boundaries: bucket,
		}},
	)
}

func FeatureFlagReason(val string) attribute.KeyValue {
	return FeatureFlagReasonKey.String(val)
}

func ExceptionType(val string) attribute.KeyValue {
	return ExceptionTypeKey.String(val)
}

// NewOTelRecorder creates a MetricsRecorder based on the provided metric.Reader. Note that, metric.NewMeterProvider is
// created here but not registered globally as this is the only place we derive a metric.Meter. Consider global provider
// registration if we need more meters
func NewOTelRecorder(exporter msdk.Reader, resource *resource.Resource, serviceName string) *MetricsRecorder {
	// create a metric provider with custom bucket size for histograms
	provider := msdk.NewMeterProvider(
		msdk.WithReader(exporter),
		// for the request duration metric we use the default bucket size which are tailored for response time in seconds
		msdk.WithView(getDurationView(serviceName, httpRequestDurationMetric, prometheus.DefBuckets)),
		// for response size we want 8 exponential bucket starting from 100 Bytes
		msdk.WithView(getDurationView(serviceName, httpResponseSizeMetric, prometheus.ExponentialBuckets(100, 10, 8))),
		// for gRPC sync stream duration we use the default bucket size which are tailored for response time in seconds
		msdk.WithView(getDurationView(serviceName, grpcSyncStreamDurationMetric, prometheus.DefBuckets)),
		// set entity producing telemetry
		msdk.WithResource(resource),
	)

	meter := provider.Meter(serviceName)

	// we can ignore errors from OpenTelemetry since they could occur if we select the wrong aggregator
	hduration, _ := meter.Float64Histogram(
		httpRequestDurationMetric,
		metric.WithDescription("Measures the duration of inbound HTTP requests."),
		metric.WithUnit("s"),
	)
	hsize, _ := meter.Float64Histogram(
		httpResponseSizeMetric,
		metric.WithDescription("Measures the size of HTTP request messages (compressed)."),
		metric.WithUnit("By"),
	)
	reqCounter, _ := meter.Int64UpDownCounter(
		httpActiveRequestsMetric,
		metric.WithDescription("Measures the number of concurrent HTTP requests that are currently in-flight."),
		metric.WithUnit("{request}"),
	)
	impressions, _ := meter.Int64Counter(
		impressionMetric,
		metric.WithDescription("Measures the number of evaluations for a given flag."),
		metric.WithUnit("{impression}"),
	)
	reasons, _ := meter.Int64Counter(
		reasonMetric,
		metric.WithDescription("Measures the number of evaluations for a given reason."),
		metric.WithUnit("{reason}"),
	)

	// gRPC Sync metrics
	syncActiveStreams, _ := meter.Int64UpDownCounter(
		grpcSyncActiveStreamsMetric,
		metric.WithDescription("Measures the number of currently active gRPC sync streaming connections."),
		metric.WithUnit("{stream}"),
	)
	syncStreamDuration, _ := meter.Float64Histogram(
		grpcSyncStreamDurationMetric,
		metric.WithDescription("Measures the duration of gRPC sync streaming connections."),
		metric.WithUnit("s"),
	)
	syncEventsSent, _ := meter.Int64Counter(
		grpcSyncEventsSentMetric,
		metric.WithDescription("Measures the number of flag configuration updates sent to sync clients."),
		metric.WithUnit("{event}"),
	)
	fetchAllFlagsRequests, _ := meter.Int64Counter(
		grpcFetchAllFlagsMetric,
		metric.WithDescription("Measures the number of unary FetchAllFlags requests."),
		metric.WithUnit("{request}"),
	)

	return &MetricsRecorder{
		httpRequestDurHistogram:   hduration,
		httpResponseSizeHistogram: hsize,
		httpRequestsInflight:      reqCounter,
		impressions:               impressions,
		reasons:                   reasons,
		syncActiveStreams:         syncActiveStreams,
		syncStreamDuration:        syncStreamDuration,
		syncEventsSent:            syncEventsSent,
		fetchAllFlagsRequests:     fetchAllFlagsRequests,
	}
}
