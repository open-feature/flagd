package sse

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/launchdarkly/eventsource"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/store"
)

// defaultHeartbeatInterval must stay well below the advertised inactivity delay so idle
// connections are not closed by clients or proxies.
const defaultHeartbeatInterval = 30 * time.Second

type Config struct {
	Logger            *logger.Logger
	HeartbeatInterval time.Duration
}

// Service serves the OFREP SSE stream. It owns an eventsource server, the change Tracker that
// drives it, and a heartbeat loop that keeps idle connections alive.
type Service struct {
	logger            *logger.Logger
	es                *eventsource.Server
	tracker           *Tracker
	heartbeatInterval time.Duration
}

// New builds an SSE Service backed by the shared flag store.
func New(s store.IStore, cfg Config) *Service {
	es := eventsource.NewServer()
	// CORS is handled by the OFREP server's middleware; avoid emitting duplicate headers.
	es.AllowCORS = false
	es.ReplayAll = false

	heartbeat := cfg.HeartbeatInterval
	if heartbeat <= 0 {
		heartbeat = defaultHeartbeatInterval
	}

	return &Service{
		logger:            cfg.Logger,
		es:                es,
		tracker:           NewTracker(cfg.Logger, s, es),
		heartbeatInterval: heartbeat,
	}
}

// Tracker exposes the change tracker so the OFREP bulk handler can resolve config versions
// for conditional (ETag/304) evaluation.
func (svc *Service) Tracker() *Tracker { return svc.tracker }

// Register mounts the stream on mux as {prefix}/{channel}, where the channel is a selector
// expression. The bare prefix is registered too, so subscribing to every flag does not depend on
// a trailing-slash redirect.
func (svc *Service) Register(mux *http.ServeMux, prefix string) {
	h := svc.Handler()
	mux.Handle(prefix, h)
	mux.Handle(prefix+"/{"+channelPathVar+"}", h)
}

// ChannelPath returns the stream path for a channel under prefix. The selector expression is a
// single path segment, so it is escaped: source selectors routinely contain "/".
func ChannelPath(prefix, channel string) string {
	if channel == "" {
		return prefix
	}
	return prefix + "/" + url.PathEscape(channel)
}

// Start runs the heartbeat loop until ctx is cancelled, then shuts down. It blocks, so it is
// intended to run in its own goroutine.
func (svc *Service) Start(ctx context.Context) error {
	ticker := time.NewTicker(svc.heartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Order matters: the tracker's watch goroutines are the only publishers, and
			// eventsource.Server.Publish blocks forever once the server is closed.
			svc.tracker.Close()
			svc.es.Close()
			if svc.logger != nil {
				svc.logger.Info("shutting down ofrep sse service")
			}
			return nil
		case <-ticker.C:
			if channels := svc.tracker.Channels(); len(channels) > 0 {
				svc.es.PublishComment(channels, "keep-alive")
			}
		}
	}
}
