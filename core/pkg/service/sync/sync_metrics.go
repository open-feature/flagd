package sync

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func init() {
	prometheus.Register(totalRequests)
}

var totalRequests = promauto.NewGauge(
	prometheus.GaugeOpts{
		Namespace: "sync",
		Name:      "active_streams",
		Help:      "Number of open sync subscriptions.",
	},
)

func (s *Server) captureMetrics(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
				syncs := s.handler.syncStore.GetSyncMetrics()
				totalRequests.Set(syncs)
			}
		}
	}()
}
