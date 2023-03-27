package sync

import (
	"context"

	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/metric"
)

const (
	serviceName = "openfeature/kube-flagd-proxy"
)

func (s *Server) captureMetrics(ctx context.Context) error {
	exporter, err := prometheus.New()
	if err != nil {
		return err
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter(serviceName)

	syncGuage, err := meter.Int64ObservableGauge(
		"sync_active_streams",
		instrument.WithDescription("number of open sync subscriptions"),
	)
	if err != nil {
		return err
	}

	_, err = meter.RegisterCallback(func(_ context.Context, o api.Observer) error {
		o.ObserveInt64(syncGuage, s.handler.syncStore.GetSyncMetrics())
		return nil
	}, syncGuage)
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}
