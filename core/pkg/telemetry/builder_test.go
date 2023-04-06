package telemetry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

func TestBuildMetricsRecorder(t *testing.T) {
	// Simple happy-path test
	recorder, err := BuildMetricsRecorder("service", Config{
		MetricsExporter: "otel",
		CollectorTarget: "localhost:8080",
	})

	require.Nil(t, err, "expected no error, but got: %v", err)
	require.NotNilf(t, recorder, "expected recorder to be non-nil")
}

func TestBuildMetricReader(t *testing.T) {
	gCtx := context.TODO()

	tests := []struct {
		name  string
		cfg   Config
		error bool
	}{
		{
			name:  "Default configurations produce default reader",
			cfg:   Config{},
			error: false,
		},
		{
			name: "Metric exporter overriding require valid overriding parameter",
			cfg: Config{
				MetricsExporter: "unsupported",
			},
			error: true,
		},
		{
			name: "Metric exporter overriding require valid configuration combination",
			cfg: Config{
				MetricsExporter: metricsExporterOtel,
				CollectorTarget: "", // collector target is unset
			},
			error: true,
		},
		{
			name: "Metric exporter overriding with valid configurations",
			cfg: Config{
				MetricsExporter: metricsExporterOtel,
				CollectorTarget: "localhost:8080",
			},
			error: false,
		},
	}

	for _, test := range tests {
		reader, err := buildMetricReader(gCtx, test.cfg)

		if test.error {
			require.NotNil(t, err, "test %s expected non-nil error", test.name)
			continue
		}

		require.Nilf(t, err, "test %s expected no error, but got: %v", test.name, err)
		require.NotNil(t, reader, "test %s expected non-nil reader", test.name)
	}
}

func TestBuildSpanProcessor(t *testing.T) {
	gCtx := context.TODO()

	tests := []struct {
		name  string
		cfg   Config
		error bool
	}{
		{
			name: "Valid configurations yield a valid processor",
			cfg: Config{
				CollectorTarget: "localhost:8080",
			},
			error: false,
		},
		{
			name:  "Invalid configurations result in an error",
			cfg:   Config{},
			error: true,
		},
	}

	for _, test := range tests {
		spanProcessor, err := BuildSpanProcessor(gCtx, test.cfg)

		if test.error {
			require.NotNil(t, err, "test %s expected non-nil error", test.name)
			continue
		}

		require.Nilf(t, err, "test %s expected no error, but got: %v", test.name, err)
		require.NotNil(t, spanProcessor, "test %s expected non-nil reader", test.name)
	}
}

func TestBuildResourceFor(t *testing.T) {
	svc := "testSvc"

	resource, err := buildResourceFor(context.Background(), svc)
	require.Nil(t, err, "expected no error, but got: %v", err)

	attributes := resource.Attributes()
	require.GreaterOrEqual(t, len(attributes), 1, "expect attributes to contain at least service name")
	require.Containsf(t, attributes, attribute.KeyValue{
		Key:   semconv.ServiceNameKey,
		Value: attribute.StringValue(svc),
	}, "")
}
