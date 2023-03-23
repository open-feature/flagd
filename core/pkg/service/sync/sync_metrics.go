package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func init() {
	prometheus.Register(totalRequests)
}

var totalRequests = promauto.NewGauge(
	prometheus.GaugeOpts{
		Name: "sync_subscriptions_gauge",
		Help: "Number of open sync subscriptions.",
	},
)

func (s *Server) captureMetrics(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
				fmt.Println("fetching metrics")
				syncs := s.handler.syncStore.GetSyncMetrics()
				totalRequests.Set(syncs)
			}
		}
	}()
}
