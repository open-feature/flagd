package sync

import (
	"net/http"
	"time"

	corsmw "github.com/open-feature/flagd/flagd/pkg/service/middleware/cors"
	metricsmw "github.com/open-feature/flagd/flagd/pkg/service/middleware/metrics"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// newHTTPServer builds the configuration server. Registering "GET" gets 405s and HEAD for free.
func newHTTPServer(cfg SvcConfigurations, handler http.Handler) *http.Server {
	selectorPath := flagsPath + "/{" + selectorPathVar + "}"

	mux := http.NewServeMux()
	mux.Handle("GET "+flagsPath, httpMetrics(cfg, flagsPath).Handler(handler))
	mux.Handle("GET "+selectorPath, httpMetrics(cfg, selectorPath).Handler(handler))

	h := otelhttp.NewHandler(mux, "flagd.sync.http")
	h = corsmw.New(cfg.CORS).Handler(h)

	return &http.Server{
		Handler:           h,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       5 * time.Second,
		MaxHeaderBytes:    int(cfg.MaxRequestHeaderBytes),
	}
}

// httpMetrics pins the handler ID to the route, since defaulting it to the URL path would make the
// client-supplied selector an unbounded metric dimension.
func httpMetrics(cfg SvcConfigurations, handlerID string) metricsmw.Middleware {
	return metricsmw.NewHTTPMetric(metricsmw.Config{
		Service:        cfg.ServiceName,
		MetricRecorder: cfg.MetricsRecorder,
		Logger:         cfg.Logger,
		HandlerID:      handlerID,
	})
}

// serveHTTP applies TLS on the same both-or-neither condition as the gRPC listener.
func (s *Service) serveHTTP() error {
	if s.certPath != "" && s.keyPath != "" {
		return s.httpServer.ServeTLS(s.httpListener, s.certPath, s.keyPath)
	}
	return s.httpServer.Serve(s.httpListener)
}
