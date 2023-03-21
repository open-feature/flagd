package otel

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

type MetricsRecorder struct {
	httpRequestDurHistogram   instrument.Float64Histogram
	httpResponseSizeHistogram instrument.Float64Histogram
	httpRequestsInflight      instrument.Int64UpDownCounter
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

func NewOTelRecorder(exporter metric.Reader, serviceName string) *MetricsRecorder {
	const requestDurationName = "http_request_duration_seconds"
	const responseSizeName = "http_response_size_bytes"

	// create a metric provider with custom bucket size for histograms
	provider := metric.NewMeterProvider(
		metric.WithReader(exporter),
		// for the request duration metric we use the default bucket size which are tailored for response time in seconds
		metric.WithView(getDurationView(requestDurationName, serviceName, prometheus.DefBuckets)),
		// for response size we want 8 exponential bucket starting from 100 Bytes
		metric.WithView(getDurationView(responseSizeName, serviceName, prometheus.ExponentialBuckets(100, 10, 8))),
	)
	meter := provider.Meter(serviceName)
	// we can ignore errors from OpenTelemetry since they could occur if we select the wrong aggregator
	hduration, _ := meter.Float64Histogram(
		requestDurationName,
		instrument.WithDescription("The latency of the HTTP requests"),
	)
	hsize, _ := meter.Float64Histogram(
		responseSizeName,
		instrument.WithDescription("The size of the HTTP responses"),
		instrument.WithUnit(unit.Bytes),
	)
	reqCounter, _ := meter.Int64UpDownCounter(
		"http_requests_inflight",
		instrument.WithDescription("The number of inflight requests being handled at the same time"),
	)
	return &MetricsRecorder{
		httpRequestDurHistogram:   hduration,
		httpResponseSizeHistogram: hsize,
		httpRequestsInflight:      reqCounter,
	}
}
