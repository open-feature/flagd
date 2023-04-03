package telemetry

import (
	"context"
	"fmt"
	"time"

	prometheusclient "github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	requestDurationName = "http_request_duration_seconds"
	responseSizeName    = "http_response_size_bytes"
	metricsExporterOtel = "otel"
	exportInterval      = 2 * time.Second
)

// SetupGlobal setup global telemetry configurations based on provided configurations. This needs to be run before
// obtaining any global.Meter or otel.Tracer based on global providers. Note that, metrics exporter works independent
// of the collector target. However, overriding metrics exporter require collector target as we only support otel
// collector
func SetupGlobal(ctx context.Context, serviceName string, metricsExporter string, collectorTarget string) error {
	// Check for collector target
	if collectorTarget == "" {
		if metricsExporter != "" {
			return fmt.Errorf("metric exporter is set(%s) without providing a collector target."+
				" collector target is required for this option", metricsExporter)
		}

		// Setup default metric provider & return
		return setupDefaultMetricProvider(serviceName)
	}

	conn, err := grpc.DialContext(ctx, collectorTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	res, err := buildResourceFor(ctx, serviceName)
	if err != nil {
		return err
	}

	// Setup otel traces backed by otel collector
	err = setupOtelCollectorTraces(conn, res)
	if err != nil {
		return err
	}

	switch metricsExporter {
	case "":
		return setupDefaultMetricProvider(serviceName)
	case metricsExporterOtel:
		return setupOtelCollectorMetrics(ctx, conn, res)
	default:
		return fmt.Errorf("provided metrics operator %s is not supported. currently only support %s",
			metricsExporter, metricsExporterOtel)
	}
}

// SetupMetricProviderWithCustomReader is a Unit Test helper with DI support for a custom metric.Reader
func SetupMetricProviderWithCustomReader(reader metric.Reader) {
	provider := metric.NewMeterProvider(metric.WithReader(reader))
	global.SetMeterProvider(provider)
}

// setupDefaultMetricProvider sets the global metric exporter with default metric exporter(prometheus)
func setupDefaultMetricProvider(serviceName string) error {
	// Default is prometheus metric exporter
	exporter, err := prometheus.New()
	if err != nil {
		return err
	}

	// create a metric provider with custom bucket size for histograms
	provider := metric.NewMeterProvider(
		metric.WithReader(exporter),
		// for the request duration metric we use the default bucket size which are tailored for response time in seconds
		metric.WithView(getDurationView(requestDurationName, serviceName, prometheusclient.DefBuckets)),
		// for response size we want 8 exponential bucket starting from 100 Bytes
		metric.WithView(getDurationView(responseSizeName, serviceName, prometheusclient.ExponentialBuckets(100, 10, 8))),
	)

	global.SetMeterProvider(provider)
	return nil
}

// setupOtelCollectorTraces registers the global traces exporter backed by otel collector
func setupOtelCollectorTraces(conn *grpc.ClientConn, res *resource.Resource) error {
	traceClient := otlptracegrpc.NewClient(otlptracegrpc.WithGRPCConn(conn))
	traceExp, err := otlptrace.New(context.Background(), traceClient)
	if err != nil {
		return err
	}

	bsp := trace.NewBatchSpanProcessor(traceExp)
	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)
	return nil
}

// setupOtelCollectorMetrics registers the global metrics exporter backed by otel collector
func setupOtelCollectorMetrics(ctx context.Context, conn *grpc.ClientConn, res *resource.Resource) error {
	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return err
	}

	provider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(
			metric.NewPeriodicReader(
				metricExporter,
				metric.WithInterval(exportInterval),
			),
		),
	)

	global.SetMeterProvider(provider)
	return nil
}

// buildResourceFor provide a resource.Resource with default resources and service key as attributes
func buildResourceFor(ctx context.Context, serviceName string) (*resource.Resource, error) {
	return resource.New(
		ctx,
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithHost(),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(semconv.ServiceNameKey.String(serviceName)),
	)
}
