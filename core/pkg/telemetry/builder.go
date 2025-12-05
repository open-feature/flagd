package telemetry

import (
	"context"
	"fmt"
	"os"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"github.com/open-feature/flagd/core/pkg/logger"
	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.uber.org/zap"
)

const (
	metricsExporterOtel = "otel"
)

type CollectorConfig struct {
	Target         string
	CertPath       string
	KeyPath        string
	ReloadInterval time.Duration
	CAPath         string
}

// Config of the telemetry runtime. These are expected to be mapped to start-up arguments
type Config struct {
	MetricsExporter string
	CollectorConfig CollectorConfig
}

func RegisterErrorHandling(log *logger.Logger) {
	otel.SetErrorHandler(otelErrorsHandler{
		logger: log,
	})
}

// BuildMetricsRecorder is a helper to build telemetry.MetricsRecorder based on configurations
func BuildMetricsRecorder(
	ctx context.Context, svcName string, svcVersion string, config Config,
) (IMetricsRecorder, error) {
	// Build metric reader based on configurations
	mReader, err := buildMetricReader(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup metric reader: %w", err)
	}

	// Build telemetry resource identifier
	rsc, err := buildResourceFor(ctx, svcName, svcVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to setup resource identifier: %w", err)
	}

	return NewOTelRecorder(mReader, rsc, svcName), nil
}

// BuildTraceProvider build and register the trace provider and propagator for the caller runtime. This method
// uses autoexport to automatically handle OTEL environment variables for trace exporters.
// Providing empty collector target results in using environment variables or falling back to noop.
func BuildTraceProvider(ctx context.Context, logger *logger.Logger, svc string, svcVersion string, cfg Config) error {
	// For backwards compatibility: set environment variables from flagd configuration
	// before calling autoexport if they are provided via flags
	if cfg.CollectorConfig.Target != "" {
		setEnvIfNotSet("OTEL_TRACES_EXPORTER", "otlp")
		setEnvIfNotSet("OTEL_EXPORTER_OTLP_ENDPOINT", cfg.CollectorConfig.Target)
		setEnvIfNotSet("OTEL_EXPORTER_OTLP_PROTOCOL", "grpc")
	}

	// Use autoexport to handle OTEL environment variables
	exporter, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		return fmt.Errorf("failed to create span exporter: %w", err)
	}

	// Skip if noop exporter (when no configuration is provided)
	if autoexport.IsNoneSpanExporter(exporter) {
		logger.Debug("skipping trace provider setup as no exporter is configured." +
			" Traces will use NoopTracerProvider provider and propagator will use no-Op TextMapPropagator")
		return nil
	}

	res, err := buildResourceFor(ctx, svc, svcVersion)
	if err != nil {
		return fmt.Errorf("failed to build resource: %w", err)
	}

	provider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithSpanProcessor(trace.NewBatchSpanProcessor(exporter)),
		trace.WithResource(res))

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return nil
}

// BuildConnectOptions is a helper to build connect options based on telemetry configurations
func BuildConnectOptions(_ Config) ([]connect.HandlerOption, error) {
	options := []connect.HandlerOption{}

	// Always add interceptor - autoexport will handle whether traces are enabled
	interceptor, err := otelconnect.NewInterceptor(otelconnect.WithTrustRemote())
	if err != nil {
		return nil, fmt.Errorf("error creating interceptor, %w", err)
	}

	options = append(options, connect.WithInterceptors(interceptor))

	return options, nil
}

// buildMetricReader builds a metric reader based on provided configurations
// Uses autoexport to automatically handle OTEL environment variables
func buildMetricReader(ctx context.Context, cfg Config) (metric.Reader, error) {
	// For backwards compatibility: set environment variables from flagd configuration
	// before calling autoexport if they are provided via flags
	if cfg.MetricsExporter == metricsExporterOtel && cfg.CollectorConfig.Target != "" {
		// Set OTEL environment variables from configuration if not already set
		setEnvIfNotSet("OTEL_METRICS_EXPORTER", "otlp")
		setEnvIfNotSet("OTEL_EXPORTER_OTLP_ENDPOINT", cfg.CollectorConfig.Target)
		setEnvIfNotSet("OTEL_EXPORTER_OTLP_PROTOCOL", "grpc")
	}

	// Use autoexport with Prometheus as fallback for backwards compatibility
	return autoexport.NewMetricReader(
		ctx,
		autoexport.WithFallbackMetricReader(buildDefaultMetricReader),
	)
}

// buildDefaultMetricReader provides the default metric reader
func buildDefaultMetricReader(_ context.Context) (metric.Reader, error) {
	p, err := prometheus.New()
	if err != nil {
		return nil, fmt.Errorf("unable to create default metric reader: %w", err)
	}
	return p, nil
}

// buildResourceFor builds a resource identifier with set of resources and service key as attributes
func buildResourceFor(ctx context.Context, serviceName string, serviceVersion string) (*resource.Resource, error) {
	r, err := resource.New(
		ctx,
		resource.WithOS(),
		resource.WithHost(),
		resource.WithProcessRuntimeVersion(),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion)),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create resource identifier: %w", err)
	}
	return r, nil
}

// setEnvIfNotSet sets an environment variable only if it's not already set
func setEnvIfNotSet(key, value string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, value)
	}
}

// OTelErrorsHandler is a custom error interceptor for OpenTelemetry
type otelErrorsHandler struct {
	logger *logger.Logger
}

func (h otelErrorsHandler) Handle(err error) {
	msg := fmt.Sprintf("OpenTelemetry Error: %s", err.Error())
	h.logger.WithFields(zap.String("component", "otel")).Debug(msg)
}
