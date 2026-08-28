package sync

import (
	"net/http"
	"time"

	corsmw "github.com/open-feature/flagd/flagd/pkg/service/middleware/cors"
	h2cmw "github.com/open-feature/flagd/flagd/pkg/service/middleware/h2c"
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
	// TLS is terminated at the listener, so http.Server sees a plain conn and would speak HTTP/1.1
	// only. ALPN still has to advertise h2 for gRPC, so h2 clients land here and need h2c to be
	// understood.
	h = h2cmw.New().Handler(h)

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
