package telemetry

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestBuildMetricsRecorder(t *testing.T) {
	// Simple happy-path test
	recorder, err := BuildMetricsRecorder(context.Background(), "service", "0.0.1", Config{
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
			name:  "Empty configurations does not result in error",
			cfg:   Config{},
			error: false,
		},
	}

	for _, test := range tests {
		err := BuildTraceProvider(gCtx, logger.NewLogger(nil, false), "svc", "0.0.1", test.cfg)

		if test.error {
			require.NotNil(t, err, "test %s expected non-nil error", test.name)
			continue
		}

		require.Nilf(t, err, "test %s expected no error, but got: %v", test.name, err)
	}
}

func TestBuildConnectOptions(t *testing.T) {
	tests := []struct {
		name        string
		cfg         Config
		optionCount int
	}{
		{
			name:        "No options for empty/default configurations",
			cfg:         Config{},
			optionCount: 0,
		},
		{
			name: "Connect option is set when telemetry target is set",
			cfg: Config{
				CollectorTarget: "localhost:8080",
			},
			optionCount: 1,
		},
	}

	for _, test := range tests {
		options := BuildConnectOptions(test.cfg)

		require.Len(t, options, test.optionCount, "option count mismatch for test %s", test.name)
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

	// we should have some logs that were intercepted
	require.True(t, observedLogs.FilterField(zap.String("component", "otel")).Len() > 0)
}

// errorExp is an exporter that always fails
type errorExp struct{}

func (e *errorExp) Temporality(k metric.InstrumentKind) metricdata.Temporality {
	return metric.DefaultTemporalitySelector(k)
}

func (e *errorExp) Aggregation(_ metric.InstrumentKind) aggregation.Aggregation {
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
