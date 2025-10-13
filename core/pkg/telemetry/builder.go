package telemetry

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"github.com/open-feature/flagd/core/pkg/certreloader"
	"github.com/open-feature/flagd/core/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	metricsExporterOtel = "otel"
	metricsExporterSDK  = "otel-sdk"
)

// CollectorConfig holds the configuration for connecting to an OpenTelemetry collector
type CollectorConfig struct {
	Target         string        // The collector endpoint (e.g., "localhost:4317")
	CertPath       string        // Path to the TLS certificate file
	KeyPath        string        // Path to the TLS key file
	CAPath         string        // Path to the CA certificate file
	Headers        string        // Additional headers in OTEL format (key1=value1,key2=value2)
	Protocol       string        // Protocol to use (e.g., "grpc", "http")
	ReloadInterval time.Duration // Interval for reloading certificates
	Timeout        time.Duration // Timeout for exporter operations
}

// Config of the telemetry runtime. These are expected to be mapped to start-up arguments
type Config struct {
	MetricsExporter string          // Type of metrics exporter ("otel" or empty for default)
	CollectorConfig CollectorConfig // Configuration for the collector
}

// RegisterErrorHandling sets up a global error handler for OpenTelemetry errors
func RegisterErrorHandling(log *logger.Logger) {
	otel.SetErrorHandler(otelErrorsHandler{
		logger: log,
	})
}

// ============================================================================
// Public API - Builders
// ============================================================================

// BuildMetricsProvider is a helper to build telemetry.MetricsRecorder based on configurations
func BuildMetricsProvider(
	ctx context.Context, svcName string, svcVersion string, config Config,
) (IMetricsRecorder, error) {
	// Build metric exporter based on configurations
	mReader, err := buildMetricExporter(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup metric exporter: %w", err)
	}

	// Build telemetry resource identifier
	rsc, err := buildResource(ctx, svcName, svcVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to setup resource identifier: %w", err)
	}

	return NewOTelRecorder(mReader, rsc, svcName), nil
}

// BuildTraceProvider builds and registers the trace provider and propagator for the caller runtime.
// This method attempts to register a global TracerProvider backed by batch SpanProcessor.
// If CollectorConfig.Target is empty, provider & propagator registration is skipped, resulting in
// tracers having NoopTracerProvider and propagator having No-Op TextMapPropagator.
func BuildTraceProvider(ctx context.Context, logger *logger.Logger, svc string, svcVersion string, cfg Config) error {
	if cfg.CollectorConfig.Target == "" {
		logger.Debug("skipping trace provider setup as collector target is not set." +
			" Traces will use NoopTracerProvider provider and propagator will use no-Op TextMapPropagator")
		return nil
	}

	exporter, err := buildGrpcTraceExporter(ctx, cfg.CollectorConfig)
	if err != nil {
		return fmt.Errorf("failed to build trace exporter: %w", err)
	}

	res, err := buildResource(ctx, svc, svcVersion)
	if err != nil {
		return fmt.Errorf("failed to build resource: %w", err)
	}

	provider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithSpanProcessor(trace.NewBatchSpanProcessor(exporter)),
		trace.WithResource(res))

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	return nil
}

// BuildConnectOptions is a helper to build connect options based on telemetry configurations
func BuildConnectOptions(cfg Config) ([]connect.HandlerOption, error) {
	options := []connect.HandlerOption{}

	// add interceptor if configuration is available for collector
	if cfg.CollectorConfig.Target != "" {
		interceptor, err := otelconnect.NewInterceptor(otelconnect.WithTrustRemote())
		if err != nil {
			return nil, fmt.Errorf("error creating interceptor: %w", err)
		}

		options = append(options, connect.WithInterceptors(interceptor))
	}

	return options, nil
}

// ============================================================================
// Internal Helpers - Transport & Credentials
// ============================================================================

// buildTransportCredentials creates gRPC transport credentials based on the collector configuration.
// Returns insecure credentials if no TLS configuration is provided.
func buildTransportCredentials(_ context.Context, cfg CollectorConfig) (credentials.TransportCredentials, error) {
	// Use insecure credentials by default
	if cfg.KeyPath == "" && cfg.CertPath == "" && cfg.CAPath == "" {
		return insecure.NewCredentials(), nil
	}

	// Build TLS configuration
	capool, err := buildCAPool(cfg.CAPath)
	if err != nil {
		return nil, err
	}

	reloader, err := buildCertReloader(cfg)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		RootCAs:    capool,
		MinVersion: tls.VersionTLS12,
		GetCertificate: func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
			certs, err := reloader.GetCertificate()
			if err != nil {
				return nil, fmt.Errorf("failed to reload certs: %w", err)
			}
			return certs, nil
		},
	}

	return credentials.NewTLS(tlsConfig), nil
}

// buildCAPool creates a certificate pool from the provided CA file path.
// Returns an empty pool if no CA path is provided.
func buildCAPool(caPath string) (*x509.CertPool, error) {
	capool := x509.NewCertPool()
	if caPath == "" {
		return capool, nil
	}

	ca, err := os.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("can't read ca file from %s: %w", caPath, err)
	}

	if !capool.AppendCertsFromPEM(ca) {
		return nil, fmt.Errorf("can't add CA '%s' to pool", caPath)
	}

	return capool, nil
}

// buildCertReloader creates a certificate reloader for automatic certificate rotation
func buildCertReloader(cfg CollectorConfig) (*certreloader.CertReloader, error) {
	reloader, err := certreloader.NewCertReloader(certreloader.Config{
		KeyPath:        cfg.KeyPath,
		CertPath:       cfg.CertPath,
		ReloadInterval: cfg.ReloadInterval,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create certreloader: %w", err)
	}
	return reloader, nil
}

// buildGrpcConnection creates a gRPC client connection with the appropriate credentials
func buildGrpcConnection(ctx context.Context, target string, cfg CollectorConfig) (*grpc.ClientConn, error) {
	transportCredentials, err := buildTransportCredentials(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to build transport credentials: %w", err)
	}

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		return nil, fmt.Errorf("error creating client connection: %w", err)
	}

	return conn, nil
}

// ============================================================================
// Internal Helpers - Headers & Options
// ============================================================================

// parseOTelHeaders parses the OTEL_EXPORTER_OTLP_HEADERS format (key1=value1,key2=value2)
// into a map[string]string
func parseOTelHeaders(headersStr string) map[string]string {
	headers := make(map[string]string)
	if headersStr == "" {
		return headers
	}

	// Split by comma to get individual key=value pairs
	pairs := strings.Split(headersStr, ",")
	for _, pair := range pairs {
		// Split by = to separate key and value
		kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(kv) == 2 {
			headers[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return headers
}

// ============================================================================
// Internal Helpers - Metric Exporter
// ============================================================================

// buildMetricExporter builds a metric exporter based on provided configurations.
// Returns a Prometheus exporter by default, or an OTLP exporter if configured.
func buildMetricExporter(ctx context.Context, cfg Config) (metric.Reader, error) {
	// Use default (Prometheus) if no exporter is specified
	if cfg.MetricsExporter == "" {
		return buildDefaultMetricExporter()
	}

	// Validate exporter type
	if !isValidMetricsExporter(cfg.MetricsExporter) {
		return nil, fmt.Errorf("provided metrics exporter %s is not supported. currently only support %s or %s",
			cfg.MetricsExporter, metricsExporterOtel, metricsExporterSDK)
	}

	// Validate collector configuration
	if cfg.CollectorConfig.Target == "" {
		return nil, fmt.Errorf("metric exporter is set(%s) without providing otel collector target. "+
			"collector target is required for this option", cfg.MetricsExporter)
	}

	return buildOTLPMetricExporter(ctx, cfg.CollectorConfig)
}

// isValidMetricsExporter checks if the provided exporter type is supported
func isValidMetricsExporter(exporter string) bool {
	return exporter == metricsExporterOtel || exporter == metricsExporterSDK
}

// buildOTLPMetricExporter creates an OTLP metric exporter with gRPC
func buildOTLPMetricExporter(ctx context.Context, cfg CollectorConfig) (metric.Reader, error) {
	conn, err := buildGrpcConnection(ctx, cfg.Target, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	exporterOpts := buildMetricExporterOptions(cfg, conn)

	otelExporter, err := otlpmetricgrpc.New(ctx, exporterOpts...)
	if err != nil {
		return nil, fmt.Errorf("error creating otel metric exporter: %w", err)
	}

	return metric.NewPeriodicReader(otelExporter), nil
}

// buildMetricExporterOptions builds the options for OTLP metric exporter
func buildMetricExporterOptions(cfg CollectorConfig, conn *grpc.ClientConn) []otlpmetricgrpc.Option {
	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithGRPCConn(conn),
	}

	if cfg.Headers != "" {
		headers := parseOTelHeaders(cfg.Headers)
		opts = append(opts, otlpmetricgrpc.WithHeaders(headers))
	}

	if cfg.Timeout > 0 {
		opts = append(opts, otlpmetricgrpc.WithTimeout(cfg.Timeout))
	}

	return opts
}

// buildDefaultMetricExporter provides the default metric exporter (Prometheus)
func buildDefaultMetricExporter() (metric.Reader, error) {
	p, err := prometheus.New()
	if err != nil {
		return nil, fmt.Errorf("unable to create default metric exporter: %w", err)
	}
	return p, nil
}

// ============================================================================
// Internal Helpers - Trace Exporter
// ============================================================================

// buildGrpcTraceExporter builds a gRPC-backed OTLP trace exporter
func buildGrpcTraceExporter(ctx context.Context, cfg CollectorConfig) (*otlptrace.Exporter, error) {
	conn, err := buildGrpcConnection(ctx, cfg.Target, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	traceOpts := buildTraceExporterOptions(cfg, conn)

	traceClient := otlptracegrpc.NewClient(traceOpts...)
	exporter, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		return nil, fmt.Errorf("error creating otel trace exporter: %w", err)
	}
	return exporter, nil
}

// buildTraceExporterOptions builds the options for OTLP trace exporter
func buildTraceExporterOptions(cfg CollectorConfig, conn *grpc.ClientConn) []otlptracegrpc.Option {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithGRPCConn(conn),
	}

	if cfg.Headers != "" {
		headers := parseOTelHeaders(cfg.Headers)
		opts = append(opts, otlptracegrpc.WithHeaders(headers))
	}

	if cfg.Timeout > 0 {
		opts = append(opts, otlptracegrpc.WithTimeout(cfg.Timeout))
	}

	return opts
}

// ============================================================================
// Internal Helpers - Resource
// ============================================================================

// buildResource builds a resource identifier with set of resources and service key as attributes
func buildResource(ctx context.Context, serviceName string, serviceVersion string) (*resource.Resource, error) {
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

// ============================================================================
// Error Handler
// ============================================================================

// otelErrorsHandler is a custom error interceptor for OpenTelemetry
type otelErrorsHandler struct {
	logger *logger.Logger
}

func (h otelErrorsHandler) Handle(err error) {
	msg := fmt.Sprintf("OpenTelemetry Error: %s", err.Error())
	h.logger.WithFields(zap.String("component", "otel")).Debug(msg)
}
