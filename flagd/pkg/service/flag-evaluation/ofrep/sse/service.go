package sse

import (
	"context"
	"net/http"
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
	active            *activeChannels
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
		active:            newActiveChannels(),
		heartbeatInterval: heartbeat,
	}
}

// Tracker exposes the change tracker so the OFREP bulk handler can resolve config versions
// for conditional (ETag/304) evaluation.
func (svc *Service) Tracker() *Tracker { return svc.tracker }

// Handler registers the request's channel in the active set (for heartbeats) and delegates to
// the eventsource server, which streams until the client disconnects.
func (svc *Service) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		channel := channelFromRequest(r)
		svc.active.add(channel)
		defer svc.active.remove(channel)
		svc.es.Handler(channel).ServeHTTP(w, r)
	})
}

// Start runs the change tracker and heartbeat loop until ctx is cancelled, then shuts the
// eventsource server down. It blocks and is intended to run in its own goroutine.
func (svc *Service) Start(ctx context.Context) error {
	go svc.tracker.Run(ctx)

	ticker := time.NewTicker(svc.heartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			svc.es.Close()
			if svc.logger != nil {
				svc.logger.Info("shutting down ofrep sse service")
			}
			return nil
		case <-ticker.C:
			channels := svc.active.snapshot()
			if len(channels) > 0 {
				svc.es.PublishComment(channels, "keep-alive")
			}
		}
	}
}
