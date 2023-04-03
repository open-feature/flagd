package telemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/metric/global"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

type MetricsRecorder struct {
	httpRequestDurHistogram   instrument.Float64Histogram
	httpResponseSizeHistogram instrument.Float64Histogram
	httpRequestsInflight      instrument.Int64UpDownCounter
	impressions               instrument.Int64Counter
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
	r.httpRequestDurHistogram.Record(ctx, duration.Seconds(), attrs...)
}

func (r MetricsRecorder) HTTPResponseSize(ctx context.Context, sizeBytes int64, attrs []attribute.KeyValue) {
	r.httpResponseSizeHistogram.Record(ctx, float64(sizeBytes), attrs...)
}

func (r MetricsRecorder) InFlightRequestStart(ctx context.Context, attrs []attribute.KeyValue) {
	r.httpRequestsInflight.Add(ctx, 1, attrs...)
}

func (r MetricsRecorder) InFlightRequestEnd(ctx context.Context, attrs []attribute.KeyValue) {
	r.httpRequestsInflight.Add(ctx, -1, attrs...)
}

func (r MetricsRecorder) Impressions(ctx context.Context, key, variant string) {
	r.impressions.Add(ctx, 1, []attribute.KeyValue{
		semconv.FeatureFlagKey(key),
		semconv.FeatureFlagVariant(variant),
		semconv.FeatureFlagProviderName("flagd"),
	}...)
}

func getDurationView(svcName, viewName string, bucket []float64) metric.View {
	return metric.NewView(
		metric.Instrument{
			// we change aggregation only for instruments with this name and scope
			Name: viewName,
			Scope: instrumentation.Scope{
				Name: svcName,
			},
		},
		metric.Stream{Aggregation: aggregation.ExplicitBucketHistogram{
			Boundaries: bucket,
		}},
	)
}

func NewOTelRecorder(serviceName string) *MetricsRecorder {
	meter := global.Meter(serviceName)

	// we can ignore errors from OpenTelemetry since they could occur if we select the wrong aggregator
	hduration, _ := meter.Float64Histogram(
		requestDurationName,
		instrument.WithDescription("The latency of the HTTP requests"),
	)
	hsize, _ := meter.Float64Histogram(
		responseSizeName,
		instrument.WithDescription("The size of the HTTP responses"),
		instrument.WithUnit("By"),
	)
	reqCounter, _ := meter.Int64UpDownCounter(
		"http_requests_inflight",
		instrument.WithDescription("The number of inflight requests being handled at the same time"),
	)
	impressions, _ := meter.Int64Counter(
		"impressions",
		instrument.WithDescription("The number of evaluation for a given flag"),
	)
	return &MetricsRecorder{
		httpRequestDurHistogram:   hduration,
		httpResponseSizeHistogram: hsize,
		httpRequestsInflight:      reqCounter,
		impressions:               impressions,
	}
}
