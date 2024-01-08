package telemetry

import (
	"context"
	"fmt"
	"os"

	"github.com/open-feature/flagd/core/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.uber.org/zap"
)

const (
	otelMetricsExporter = "OTEL_METRICS_EXPORTER"
	otelTracesExporter  = "OTEL_TRACES_EXPORTER"

	otelExporterOtlpProtocol        = "OTEL_EXPORTER_OTLP_PROTOCOL"
	otelExporterOtlpMetricsProtocol = "OTEL_EXPORTER_OTLP_METRICS_PROTOCOL"
	otelExporterOtlpTracesProtocol  = "OTEL_EXPORTER_OTLP_TRACES_PROTOCOL"

	otelExporterNone       = "none"
	otelExporterOtlp       = "otlp"
	otelExporterPrometheus = "prometheus"

	otelExporterOtlpProtocolGrpc         = "grpc"
	otelExporterOtlpProtocolHTTPProtobuf = "http/protobuf"
)

var (
	nilReader       metric.Reader
	nilSpanExporter trace.SpanExporter
)

func RegisterErrorHandling(log *logger.Logger) {
	otel.SetErrorHandler(otelErrorsHandler{
		logger: log,
	})
}

// BuildMetricsRecorder is a helper to build telemetry.MetricsRecorder based on configurations
func BuildMetricsRecorder(
	ctx context.Context, logger *logger.Logger, svcName string, svcVersion string,
) (*MetricsRecorder, error) {
	// Build metric reader based on configurations
	mReader, err := buildMetricReader(ctx, logger)
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
// attempt to register a global TracerProvider backed by batch SpanProcessor.Config.
func BuildTraceProvider(ctx context.Context, logger *logger.Logger, svc string, svcVersion string) error {
	exporter, err := buildOtlpExporter(ctx, logger)
	if err != nil {
		return err
	}

	res, err := buildResourceFor(ctx, svc, svcVersion)
	if err != nil {
		return err
	}

	provider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithSpanProcessor(trace.NewBatchSpanProcessor(exporter)),
		trace.WithResource(res))

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return nil
}

// buildMetricReader builds a metric reader based on provided configurations
func buildMetricReader(ctx context.Context, logger *logger.Logger) (metric.Reader, error) {
	switch k, v := getOtelExporter(logger, otelMetricsExporter); v {
	case otelExporterNone:
		logger.Debug(fmt.Sprintf("skipping setup for metrics due to %s=%s", k, v))
		return nilReader, nil
	case otelExporterOtlp:
		logger.Debug(fmt.Sprintf("setting up metrics based on %s=%s", k, v))
		return buildMetricReaderOtlp(ctx, logger)
	case otelExporterPrometheus:
		logger.Debug(fmt.Sprintf("setting up metrics based on %s=%s", k, v))
		return buildMetricReaderPrometheus()
	default:
		logger.Debug(fmt.Sprintf("skipping unsupported value for %s=%s", k, v))
		return nilReader, nil
	}
}

// buildOtlpExporter is a helper to build grpc backed otlp trace exporter
func buildOtlpExporter(ctx context.Context, logger *logger.Logger) (trace.SpanExporter, error) {
	switch k, v := getOtelExporter(logger, otelTracesExporter); v {
	case otelExporterNone:
		logger.Debug(fmt.Sprintf("skipping setup for traces due to %s=%s", k, v))
		return nilSpanExporter, nil
	case otelExporterOtlp:
		break
	default:
		logger.Debug(fmt.Sprintf("skipping unsupported value for %s=%s", k, v))
		return nilSpanExporter, nil
	}

	var client otlptrace.Client
	var err error

	switch k, v := getOtelExporterProtocol(logger, otelExporterOtlpTracesProtocol, otelExporterOtlpProtocol); v {
	case otelExporterOtlpProtocolGrpc:
		client = otlptracegrpc.NewClient()
	case otelExporterOtlpProtocolHTTPProtobuf:
		client = otlptracehttp.NewClient()
	default:
		logger.Debug(fmt.Sprintf("skipping unsupported value %s for %s", v, k))
		return nilSpanExporter, nil
	}

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("error starting otel v: %w", err)
	}
	return exporter, nil
}

// buildMetricReaderOtlp provides an OTLP metric reader
func buildMetricReaderOtlp(ctx context.Context, logger *logger.Logger) (metric.Reader, error) {
	var exporter metric.Exporter
	var err error

	switch k, v := getOtelExporterProtocol(logger, otelExporterOtlpMetricsProtocol, otelExporterOtlpProtocol); v {
	case otelExporterOtlpProtocolGrpc:
		if exporter, err = otlpmetricgrpc.New(ctx); err == nil {
			return metric.NewPeriodicReader(exporter), nil
		} else {
			return nil, fmt.Errorf("error creating otel metric exporter based on %s=%s: %w", k, v, err)
		}
	case otelExporterOtlpProtocolHTTPProtobuf:
		if exporter, err = otlpmetrichttp.New(ctx); err == nil {
			return metric.NewPeriodicReader(exporter), nil
		} else {
			return nil, fmt.Errorf("error creating otel metric exporter based on %s=%s: %w", k, v, err)
		}
	default:
		logger.Debug(fmt.Sprintf("skipping unsupported value %s for %s", v, k))
		return nilReader, nil
	}
}

// buildMetricReaderPrometheus provides a prometheus metric reader
func buildMetricReaderPrometheus() (metric.Reader, error) {
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

// OTelErrorsHandler is a custom error interceptor for OpenTelemetry
type otelErrorsHandler struct {
	logger *logger.Logger
}

func (h otelErrorsHandler) Handle(err error) {
	msg := fmt.Sprintf("OpenTelemetry Error: %s", err.Error())
	h.logger.WithFields(zap.String("component", "otel")).Debug(msg)
}

func getOtelExporter(logger *logger.Logger, signalEnvVar string) (string, string) {
	if v := os.Getenv(signalEnvVar); v != "" {
		logger.Debug(fmt.Sprintf("resolved %s=%s", signalEnvVar, v))
		return signalEnvVar, v
	}
	return "", otelExporterOtlp
}

func getOtelExporterProtocol(logger *logger.Logger, signalEnvVar string, commonEnvVar string) (string, string) {
	if v := os.Getenv(signalEnvVar); v != "" {
		logger.Debug(fmt.Sprintf("resolved %s=%s", signalEnvVar, v))
		return signalEnvVar, v
	}
	if v := os.Getenv(commonEnvVar); v != "" {
		logger.Debug(fmt.Sprintf("resolved %s=%s", commonEnvVar, v))
		return commonEnvVar, v
	}
	return "", otelExporterOtlpProtocolGrpc
}
