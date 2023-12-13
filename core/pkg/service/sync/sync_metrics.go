package sync

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

const (
	serviceName = "flagd-proxy"
)

func (s *Server) captureMetrics() error {
	exporter, err := prometheus.New()
	if err != nil {
		return fmt.Errorf("unable to create prometheus exporter: %w", err)
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter(serviceName)

	syncGuage, err := meter.Int64ObservableGauge(
		"sync_active_streams",
		api.WithDescription("number of open sync subscriptions"),
	)
	if err != nil {
		return fmt.Errorf("unable to create active subscription metric gauge: %w", err)
	}

	_, err = meter.RegisterCallback(func(_ context.Context, o api.Observer) error {
		o.ObserveInt64(syncGuage, s.oldHandler.syncStore.GetActiveSubscriptionsInt64())
		return nil
	}, syncGuage)
	if err != nil {
		return fmt.Errorf("unable to register active subscription metric callback: %w", err)
	}

	return nil
}
