package telemetry

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestBuildMetricsRecorder(t *testing.T) {
	// Simple happy-path test
	recorder, err := BuildMetricsRecorder(context.Background(), logger.NewLogger(nil, false), "service", "0.0.1")

	require.Nil(t, err, "expected no error, but got: %v", err)
	require.NotNilf(t, recorder, "expected recorder to be non-nil")
}

func TestBuildMetricReader(t *testing.T) {
	gCtx := context.TODO()

	container, cleanup := startOtelContainer(gCtx)
	defer cleanup(gCtx, container)

	tests := []struct {
		name  string
		env   map[string]string
		error bool
	}{
		{
			name:  "Default configurations produce default reader",
			env:   map[string]string{},
			error: false,
		},
		{
			name: "Metric exporter overriding with valid configurations",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "http://localhost:4317",
				"OTEL_EXPORTER_OTLP_PROTOCOL": "grpc",
				"OTEL_METRICS_EXPORTER":       "otlp",
			},
			error: false,
		},
	}

	for _, test := range tests {
		for k, v := range test.env {
			_ = os.Setenv(k, v)
		}

		reader, err := buildMetricReader(gCtx, logger.NewLogger(nil, false))

		if test.error {
			require.NotNil(t, err, "test %s expected non-nil error", test.name)
		} else {
			require.Nilf(t, err, "test %s expected no error, but got: %v", test.name, err)
			require.NotNil(t, reader, "test %s expected non-nil reader", test.name)
		}

		for k := range test.env {
			_ = os.Unsetenv(k)
		}
	}
}

func TestBuildSpanProcessor(t *testing.T) {
	gCtx := context.TODO()

	container, cleanup := startOtelContainer(gCtx)
	defer cleanup(gCtx, container)

	tests := []struct {
		name  string
		env   map[string]string
		error bool
	}{
		{
			name:  "Valid configurations yield a valid processor",
			env:   map[string]string{},
			error: false,
		},
		{
			name: "Valid configurations yield a valid processor",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "http://localhost:4317",
				"OTEL_EXPORTER_OTLP_PROTOCOL": "grpc",
			},
			error: false,
		},
		{
			name: "Valid configurations yield a valid processor",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "http://localhost:4318",
				"OTEL_EXPORTER_OTLP_PROTOCOL": "http/protobuf",
			},
			error: false,
		},
		{
			name: "Valid configurations yield a valid processor",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT": "http://localhost:4317",
				"OTEL_EXPORTER_OTLP_TRACES_PROTOCOL": "grpc",
			},
			error: false,
		},
		{
			name: "Valid configurations yield a valid processor",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT": "http://localhost:4318",
				"OTEL_EXPORTER_OTLP_TRACES_PROTOCOL": "http/protobuf",
			},
			error: false,
		},
		{
			name: "Valid configurations yield a valid processor",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT":        "http://localhost:4317",
				"OTEL_EXPORTER_OTLP_TRACES_PROTOCOL": "grpc",
			},
			error: false,
		},
		{
			name: "Valid configurations yield a valid processor",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT":        "http://localhost:4318",
				"OTEL_EXPORTER_OTLP_TRACES_PROTOCOL": "http/protobuf",
			},
			error: false,
		},
		{
			name: "Valid configurations yield a valid processor",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_PROTOCOL":        "grpc",
				"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT": "http://localhost:4317",
			},
			error: false,
		},
		{
			name: "Valid configurations yield a valid processor",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_PROTOCOL":        "http/protobuf",
				"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT": "http://localhost:4318",
			},
			error: false,
		},
	}

	for _, test := range tests {
		for k, v := range test.env {
			_ = os.Setenv(k, v)
		}

		err := BuildTraceProvider(gCtx, logger.NewLogger(nil, false), "svc", "0.0.1")

		if test.error {
			require.NotNil(t, err, "test %s expected non-nil error", test.name)
		} else {
			require.Nilf(t, err, "test %s expected no error, but got: %v", test.name, err)
		}

		for k := range test.env {
			_ = os.Unsetenv(k)
		}
	}
}

func TestBuildResourceFor(t *testing.T) {
	svc := "testSvc"
	svcVersion := "0.0.1"

	resource, err := buildResourceFor(context.Background(), svc, svcVersion)
	require.Nil(t, err, "expected no error, but got: %v", err)

	attributes := resource.Attributes()
	require.GreaterOrEqual(t, len(attributes), 2, "expect attributes to contain service name, version")
	require.Containsf(t, attributes, attribute.KeyValue{
		Key:   semconv.ServiceNameKey,
		Value: attribute.StringValue(svc),
	}, "expected resource to contain service name")
	require.Containsf(t, attributes, attribute.KeyValue{
		Key:   semconv.ServiceVersionKey,
		Value: attribute.StringValue(svcVersion),
	}, "expected resource to contain service version")
}

func TestErrorIntercepted(t *testing.T) {
	// register the OTel error handling
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	log := logger.NewLogger(observedLogger, true)
	RegisterErrorHandling(log)

	// configure a metric reader with an exporter that only returns error
	reader := metric.NewPeriodicReader(&errorExp{}, metric.WithInterval(1*time.Millisecond))
	rs := resource.NewWithAttributes("testSchema")
	NewOTelRecorder(reader, rs, "testSvc")
	var data metricdata.ResourceMetrics
	err := reader.Collect(context.TODO(), &data)
	require.Nil(t, err)

	// we should have some logs that were intercepted, test every 10ms for 1s
	require.Eventually(t, func() bool { return observedLogs.FilterField(zap.String("component", "otel")).Len() > 0 }, 1000*time.Millisecond, 10*time.Millisecond)
}

// errorExp is an exporter that always fails
type errorExp struct{}

func (e *errorExp) Temporality(k metric.InstrumentKind) metricdata.Temporality {
	return metric.DefaultTemporalitySelector(k)
}

func (e *errorExp) Aggregation(_ metric.InstrumentKind) metric.Aggregation {
	return nil
}

func (e *errorExp) Export(_ context.Context, _ *metricdata.ResourceMetrics) error {
	return fmt.Errorf("I am an error")
}

func (e *errorExp) ForceFlush(_ context.Context) error {
	return fmt.Errorf("I am an error")
}

func (e *errorExp) Shutdown(_ context.Context) error {
	return fmt.Errorf("I am an error")
}

func startOtelContainer(ctx context.Context) (testcontainers.Container, func(context.Context, testcontainers.Container)) {
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "ghcr.io/open-telemetry/opentelemetry-collector-releases/opentelemetry-collector",
			ExposedPorts: []string{"4317/tcp", "4318/tcp"},
		},
		Started: true,
	})
	if err != nil {
		panic(err)
	}

	return container, func(ctx context.Context, container testcontainers.Container) {
		if err := container.Terminate(ctx); err != nil {
			panic(err)
		}
	}
}
