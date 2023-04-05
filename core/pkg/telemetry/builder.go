package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	metricsExporterOtel = "otel"
	exportInterval      = 2 * time.Second
)

type Config struct {
	MetricsExporter string
	CollectorTarget string
}

// BuildMetricReader builds a metric reader based on provided configurations
func BuildMetricReader(ctx context.Context, cfg Config) (metric.Reader, error) {
	if cfg.MetricsExporter == "" {
		return BuildDefaultMetricReader()
	}

	// Handle metric reader override
	if cfg.MetricsExporter != metricsExporterOtel {
		return nil, fmt.Errorf("provided metrics operator %s is not supported. currently only support %s",
			cfg.MetricsExporter, metricsExporterOtel)
	}

	// Otel override require target configuration
	if cfg.CollectorTarget == "" {
		return nil, fmt.Errorf("metric exporter is set(%s) without providing otel collector target."+
			" collector target is required for this option", cfg.MetricsExporter)
	}

	// Non-blocking, insecure grpc connection
	conn, err := grpc.DialContext(ctx, cfg.CollectorTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// Otel metric exporter
	otelExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	return metric.NewPeriodicReader(otelExporter, metric.WithInterval(exportInterval)), nil
}

// BuildSpanProcessor builds a span processor based on provided configurations
// todo - consider fallback mechanism(ex:- to console) if collector is unset
func BuildSpanProcessor(ctx context.Context, cfg Config) (trace.SpanProcessor, error) {
	if cfg.CollectorTarget == "" {
		return nil, fmt.Errorf("otel collector target is required for span processor")
	}

	// Non-blocking, insecure grpc connection
	conn, err := grpc.DialContext(ctx, cfg.CollectorTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	traceClient := otlptracegrpc.NewClient(otlptracegrpc.WithGRPCConn(conn))
	traceExp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		return nil, err
	}

	return trace.NewBatchSpanProcessor(traceExp), nil
}

// BuildResourceFor builds a resource identifier with set of resources and service key as attributes
func BuildResourceFor(ctx context.Context, serviceName string) (*resource.Resource, error) {
	return resource.New(
		ctx,
		resource.WithOS(),
		resource.WithHost(),
		resource.WithProcessRuntimeVersion(),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(semconv.ServiceNameKey.String(serviceName)),
	)
}

// BuildDefaultMetricReader provides the default metric reader
func BuildDefaultMetricReader() (metric.Reader, error) {
	return prometheus.New()
}
