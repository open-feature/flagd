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

// BuildMetricsRecorder is a helper to build telemetry.MetricsRecorder based on configurations
func BuildMetricsRecorder(svcName string, config Config) (*MetricsRecorder, error) {
	// Build metric reader based on configurations
	mReader, err := buildMetricReader(context.TODO(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup metric reader: %w", err)
	}

	// Build telemetry resource identifier
	rsc, err := buildResourceFor(context.Background(), svcName)
	if err != nil {
		return nil, fmt.Errorf("failed to setup resource identifier: %w", err)
	}

	return NewOTelRecorder(mReader, rsc, svcName), nil
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

// buildMetricReader builds a metric reader based on provided configurations
func buildMetricReader(ctx context.Context, cfg Config) (metric.Reader, error) {
	if cfg.MetricsExporter == "" {
		return buildDefaultMetricReader()
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

// buildDefaultMetricReader provides the default metric reader
func buildDefaultMetricReader() (metric.Reader, error) {
	return prometheus.New()
}

// buildResourceFor builds a resource identifier with set of resources and service key as attributes
func buildResourceFor(ctx context.Context, serviceName string) (*resource.Resource, error) {
	return resource.New(
		ctx,
		resource.WithOS(),
		resource.WithHost(),
		resource.WithProcessRuntimeVersion(),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(semconv.ServiceNameKey.String(serviceName)),
	)
}
