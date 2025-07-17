package telemetry

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
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
// attempt to register a global TracerProvider backed by batch SpanProcessor.Config. CollectorTarget can be used to
// provide the grpc collector target. Providing empty target results in skipping provider & propagator registration.
// This results in tracers having NoopTracerProvider and propagator having No-Op TextMapPropagator performing no action
func BuildTraceProvider(ctx context.Context, logger *logger.Logger, svc string, svcVersion string, cfg Config) error {
	if cfg.CollectorConfig.Target == "" {
		logger.Debug("skipping trace provider setup as collector target is not set." +
			" Traces will use NoopTracerProvider provider and propagator will use no-Op TextMapPropagator")
		return nil
	}

	exporter, err := buildOtlpExporter(ctx, cfg.CollectorConfig)
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

// BuildConnectOptions is a helper to build connect options based on telemetry configurations
func BuildConnectOptions(cfg Config) ([]connect.HandlerOption, error) {
	options := []connect.HandlerOption{}

	// add interceptor if configuration is available for collector
	if cfg.CollectorConfig.Target != "" {
		interceptor, err := otelconnect.NewInterceptor(otelconnect.WithTrustRemote())
		if err != nil {
			return nil, fmt.Errorf("error creating interceptor, %w", err)
		}

		options = append(options, connect.WithInterceptors(interceptor))
	}

	return options, nil
}

func buildTransportCredentials(_ context.Context, cfg CollectorConfig) (credentials.TransportCredentials, error) {
	creds := insecure.NewCredentials()
	if cfg.KeyPath != "" || cfg.CertPath != "" || cfg.CAPath != "" {
		capool := x509.NewCertPool()
		if cfg.CAPath != "" {
			ca, err := os.ReadFile(cfg.CAPath)
			if err != nil {
				return nil, fmt.Errorf("can't read ca file from %s", cfg.CAPath)
			}
			if !capool.AppendCertsFromPEM(ca) {
				return nil, fmt.Errorf("can't add CA '%s' to pool", cfg.CAPath)
			}
		}

		reloader, err := certreloader.NewCertReloader(certreloader.Config{
			KeyPath:        cfg.KeyPath,
			CertPath:       cfg.CertPath,
			ReloadInterval: cfg.ReloadInterval,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create certreloader: %w", err)
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

		creds = credentials.NewTLS(tlsConfig)
	}

	return creds, nil
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
	if cfg.CollectorConfig.Target == "" {
		return nil, fmt.Errorf("metric exporter is set(%s) without providing otel collector target."+
			" collector target is required for this option", cfg.MetricsExporter)
	}

	transportCredentials, err := buildTransportCredentials(ctx, cfg.CollectorConfig)
	if err != nil {
		return nil, fmt.Errorf("metric export would not build transport credentials: %w", err)
	}

	// Non-blocking, insecure grpc connection
	conn, err := grpc.NewClient(cfg.CollectorConfig.Target, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		return nil, fmt.Errorf("error creating client connection: %w", err)
	}

	// Otel metric exporter
	otelExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("error creating otel metric exporter: %w", err)
	}

	return metric.NewPeriodicReader(otelExporter), nil
}

// buildOtlpExporter is a helper to build grpc backed otlp trace exporter
func buildOtlpExporter(ctx context.Context, cfg CollectorConfig) (*otlptrace.Exporter, error) {
	transportCredentials, err := buildTransportCredentials(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("metric export would not build transport credentials: %w", err)
	}

	// Non-blocking, grpc connection
	conn, err := grpc.NewClient(cfg.Target, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		return nil, fmt.Errorf("error creating client connection: %w", err)
	}

	traceClient := otlptracegrpc.NewClient(otlptracegrpc.WithGRPCConn(conn))
	exporter, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		return nil, fmt.Errorf("error starting otel exporter: %w", err)
	}
	return exporter, nil
}

// buildDefaultMetricReader provides the default metric reader
func buildDefaultMetricReader() (metric.Reader, error) {
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
