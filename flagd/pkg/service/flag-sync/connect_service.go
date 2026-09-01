package sync

import (
	"net/http"
	"time"

	"buf.build/gen/go/open-feature/flagd/connectrpc/go/flagd/sync/v1/syncv1connect"
	"connectrpc.com/connect"
	flagdService "github.com/open-feature/flagd/flagd/pkg/service"
	corsmw "github.com/open-feature/flagd/flagd/pkg/service/middleware/cors"
	h2cmw "github.com/open-feature/flagd/flagd/pkg/service/middleware/h2c"
	metricsmw "github.com/open-feature/flagd/flagd/pkg/service/middleware/metrics"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/protobuf/encoding/protojson"
)

// idleTimeout only bounds connections holding no streams, so it can be generous.
const idleTimeout = 5 * time.Minute

// newConnectServer serves gRPC, gRPC-Web, Connect and the flag configuration route on one port.
func newConnectServer(
	cfg SvcConfigurations, handler syncv1connect.FlagSyncServiceHandler, flagsHandler http.Handler,
) *http.Server {
	opts := append([]connect.HandlerOption{}, cfg.Options...)
	opts = append(opts,
		// json parsing configuration - we emit "unpopulated" fields (falsy fields are not dropped)
		flagdService.WithJSON(
			protojson.MarshalOptions{EmitUnpopulated: true},
			protojson.UnmarshalOptions{DiscardUnknown: true},
		),
	)
	// connect has no default receive limit; zero means unlimited, as the flag documents.
	if cfg.MaxRequestBodyBytes > 0 {
		opts = append(opts, connect.WithReadMaxBytes(int(cfg.MaxRequestBodyBytes)))
	}

	mux := http.NewServeMux()

	// Traced by the interceptor in cfg.Options; an otelhttp wrapper here would double every span.
	connectPath, connectHandler := syncv1connect.NewFlagSyncServiceHandler(handler, opts...)
	mux.Handle(connectPath, httpMetrics(cfg, connectPath).Handler(connectHandler))

	if flagsHandler != nil {
		// Registering "GET" gets 405s and HEAD for free.
		mux.Handle("GET "+flagsPath, instrumentFlagsRoute(cfg, flagsPath, flagsHandler))
	}

	var h http.Handler = mux
	h = corsmw.New(cfg.CORS).Handler(h)
	if cfg.CertPath == "" || cfg.KeyPath == "" {
		// no ALPN to negotiate over plaintext, so gRPC clients need h2c
		h = h2cmw.New().Handler(h)
	}

	return &http.Server{
		Handler:           h,
		ReadHeaderTimeout: 3 * time.Second,
		// slowloris protection; only closes the request body, so SyncFlags streams are unaffected
		ReadTimeout: 5 * time.Second,
		// Must be set: http2.ConfigureServer otherwise promotes ReadTimeout to the h2 idle timeout,
		// GOAWAYing streamless connections after 5s. See issue #1998.
		IdleTimeout:    idleTimeout,
		MaxHeaderBytes: int(cfg.MaxRequestHeaderBytes),
	}
}

// instrumentFlagsRoute wraps a flag configuration route, which the connect interceptor does not cover.
func instrumentFlagsRoute(cfg SvcConfigurations, route string, handler http.Handler) http.Handler {
	return otelhttp.NewHandler(httpMetrics(cfg, route).Handler(handler), "flagd.sync.http")
}

// httpMetrics pins the handler ID to the route; the URL path default would make a path-segment
// selector an unbounded metric dimension.
func httpMetrics(cfg SvcConfigurations, handlerID string) metricsmw.Middleware {
	return metricsmw.NewHTTPMetric(metricsmw.Config{
		Service:        cfg.ServiceName,
		MetricRecorder: cfg.MetricsRecorder,
		Logger:         cfg.Logger,
		HandlerID:      handlerID,
	})
}
