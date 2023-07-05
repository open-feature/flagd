package telemetry

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	msdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

const (
	ProviderName = "flagd"

	FeatureFlagReasonKey = attribute.Key("feature_flag.reason")
	ExceptionTypeKey     = attribute.Key("ExceptionTypeKeyName")

	httpRequestDurationMetric = "http.server.duration"
	httpResponseSizeMetric    = "http.server.response.size"
	httpActiveRequestsMetric  = "http.server.active_requests"
	impressionMetric          = "feature_flag." + ProviderName + ".impression"
	reasonMetric              = "feature_flag." + ProviderName + ".evaluation.reason"
)

type MetricsRecorder struct {
	httpRequestDurHistogram   metric.Float64Histogram
	httpResponseSizeHistogram metric.Float64Histogram
	httpRequestsInflight      metric.Int64UpDownCounter
	impressions               metric.Int64Counter
	reasons                   metric.Int64Counter
}

func (r MetricsRecorder) HTTPAttributes(svcName, url, method, code string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.ServiceNameKey.String(svcName),
		semconv.HTTPURLKey.String(url),
		semconv.HTTPMethodKey.String(method),
		semconv.HTTPStatusCodeKey.String(code),
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

func getDurationView(svcName, viewName string, bucket []float64) msdk.View {
	return msdk.NewView(
		msdk.Instrument{
			// we change aggregation only for instruments with this name and scope
			Name: viewName,
			Scope: instrumentation.Scope{
				Name: svcName,
			},
		},
		msdk.Stream{Aggregation: aggregation.ExplicitBucketHistogram{
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
		msdk.WithView(getDurationView(httpRequestDurationMetric, serviceName, prometheus.DefBuckets)),
		// for response size we want 8 exponential bucket starting from 100 Bytes
		msdk.WithView(getDurationView(httpResponseSizeMetric, serviceName, prometheus.ExponentialBuckets(100, 10, 8))),
		// set entity producing telemetry
		msdk.WithResource(resource),
	)

	meter := provider.Meter(serviceName)

	// we can ignore errors from OpenTelemetry since they could occur if we select the wrong aggregator
	hduration, _ := meter.Float64Histogram(
		httpRequestDurationMetric,
		metric.WithDescription("The latency of the HTTP requests"),
		metric.WithUnit("s"),
	)
	hsize, _ := meter.Float64Histogram(
		httpResponseSizeMetric,
		metric.WithDescription("The size of the HTTP responses"),
		metric.WithUnit("By"),
	)
	reqCounter, _ := meter.Int64UpDownCounter(
		httpActiveRequestsMetric,
		metric.WithDescription("The number of inflight requests being handled at the same time"),
		metric.WithUnit("{request}"),
	)
	impressions, _ := meter.Int64Counter(
		impressionMetric,
		metric.WithDescription("The number of evaluations for a given flag"),
		metric.WithUnit("{impression}"),
	)
	reasons, _ := meter.Int64Counter(
		reasonMetric,
		metric.WithDescription("The number of evaluations for a given reason"),
		metric.WithUnit("{reason}"),
	)
	return &MetricsRecorder{
		httpRequestDurHistogram:   hduration,
		httpResponseSizeHistogram: hsize,
		httpRequestsInflight:      reqCounter,
		impressions:               impressions,
		reasons:                   reasons,
	}
}
