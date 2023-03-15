package metrics

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"

	semconv "go.opentelemetry.io/otel/semconv/v1.13.0"
)

var (
	_ http.ResponseWriter = &responseWriterInterceptor{}
	_ http.Hijacker       = &responseWriterInterceptor{}
	_ http.Flusher        = &responseWriterInterceptor{}
)

type HTTPReqProperties struct {
	Service string
	ID      string
	Method  string
	Code    string
}

type Recorder interface {
	// OTelObserveHTTPRequestDuration measures the duration of an HTTP request.
	OTelObserveHTTPRequestDuration(props HTTPReqProperties, duration time.Duration)
	// OTelObserveHTTPResponseSize measures the size of an HTTP response in bytes.
	OTelObserveHTTPResponseSize(props HTTPReqProperties, sizeBytes int64)

	// OTelInFlightRequestStart count the active requests.
	OTelInFlightRequestStart(props HTTPReqProperties)
	// OTelInFlightRequestEnd count the finished requests.
	OTelInFlightRequestEnd(props HTTPReqProperties)
}

type Reporter interface {
	Method() string
	URLPath() string
	StatusCode() int
	BytesWritten() int64
}

type HTTPProperties struct {
	Service string
	ID      string
}

type OTelMetricsRecorder struct {
	httpRequestDurHistogram   instrument.Float64Histogram
	httpResponseSizeHistogram instrument.Float64Histogram
	httpRequestsInflight      instrument.Int64UpDownCounter
}

func (r OTelMetricsRecorder) setAttributes(p HTTPReqProperties) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.ServiceNameKey.String(p.Service),
		semconv.HTTPURLKey.String(p.ID),
		semconv.HTTPMethodKey.String(p.Method),
		semconv.HTTPStatusCodeKey.String(p.Code),
	}
}

func (r OTelMetricsRecorder) OTelObserveHTTPRequestDuration(p HTTPReqProperties, duration time.Duration) {
	r.httpRequestDurHistogram.Record(context.TODO(), duration.Seconds(), r.setAttributes(p)...)
}

func (r OTelMetricsRecorder) OTelObserveHTTPResponseSize(p HTTPReqProperties, sizeBytes int64) {
	r.httpResponseSizeHistogram.Record(context.TODO(), float64(sizeBytes), r.setAttributes(p)...)
}

func (r OTelMetricsRecorder) OTelInFlightRequestStart(p HTTPReqProperties) {
	r.httpRequestsInflight.Add(context.TODO(), 1, r.setAttributes(p)...)
}

func (r OTelMetricsRecorder) OTelInFlightRequestEnd(p HTTPReqProperties) {
	r.httpRequestsInflight.Add(context.TODO(), -1, r.setAttributes(p)...)
}

type MiddlewareConfig struct {
	recorder           Recorder
	MetricReader       metric.Reader
	Logger             *logger.Logger
	Service            string
	GroupedStatus      bool
	DisableMeasureSize bool
}

type Middleware struct {
	cfg MiddlewareConfig
}

func New(cfg MiddlewareConfig) Middleware {
	cfg.defaults()
	m := Middleware{cfg: cfg}
	return m
}

func (cfg *MiddlewareConfig) defaults() {
	if cfg.Logger == nil {
		log.Fatal("missing logger")
	}
	if cfg.MetricReader == nil {
		log.Fatal("missing MetricReader/Exporter")
	}
	cfg.recorder = cfg.newOTelRecorder(cfg.MetricReader)
}

func (cfg *MiddlewareConfig) getDurationView(name string, bucket []float64) metric.View {
	return metric.NewView(
		metric.Instrument{
			// we change aggregation only for instruments with this name and scope
			Name: name,
			Scope: instrumentation.Scope{
				Name: cfg.Service,
			},
		},
		metric.Stream{Aggregation: aggregation.ExplicitBucketHistogram{
			Boundaries: bucket,
		}},
	)
}

func (cfg *MiddlewareConfig) newOTelRecorder(exporter metric.Reader) *OTelMetricsRecorder {
	const requestDurationName = "http_request_duration_seconds"
	const responseSizeName = "http_response_size_bytes"

	// create a metric provider with custom bucket size for histograms
	provider := metric.NewMeterProvider(
		metric.WithReader(exporter),
		metric.WithView(cfg.getDurationView(requestDurationName, prometheus.DefBuckets)),
		metric.WithView(cfg.getDurationView(responseSizeName, prometheus.ExponentialBuckets(100, 10, 8))),
	)
	meter := provider.Meter(cfg.Service)
	// we can ignore errors from OpenTelemetry since they could occur if we select the wrong aggregator
	hduration, _ := meter.Float64Histogram(
		requestDurationName,
		instrument.WithDescription("The latency of the HTTP requests"),
	)
	hsize, _ := meter.Float64Histogram(
		responseSizeName,
		instrument.WithDescription("The size of the HTTP responses"),
		instrument.WithUnit(unit.Bytes),
	)
	reqCounter, _ := meter.Int64UpDownCounter(
		"http_requests_inflight",
		instrument.WithDescription("The number of inflight requests being handled at the same time"),
	)
	return &OTelMetricsRecorder{
		httpRequestDurHistogram:   hduration,
		httpResponseSizeHistogram: hsize,
		httpRequestsInflight:      reqCounter,
	}
}

func (m Middleware) Measure(handlerID string, reporter Reporter, next func()) {
	// If there isn't predefined handler ID we
	// set that ID as the URL path.
	hid := handlerID
	if handlerID == "" {
		hid = reporter.URLPath()
	}

	// If we need to group the status code, it uses the
	// first number of the status code because is the least
	// required identification way.
	var code string
	if m.cfg.GroupedStatus {
		code = fmt.Sprintf("%dxx", reporter.StatusCode()/100)
	} else {
		code = strconv.Itoa(reporter.StatusCode())
	}
	props := HTTPReqProperties{
		Service: m.cfg.Service,
		ID:      hid,
		Method:  reporter.Method(),
		Code:    code,
	}

	m.cfg.recorder.OTelInFlightRequestStart(props)
	defer m.cfg.recorder.OTelInFlightRequestEnd(props)

	// Start the timer and when finishing measure the duration.
	start := time.Now()
	defer func() {
		duration := time.Since(start)

		m.cfg.recorder.OTelObserveHTTPRequestDuration(props, duration)

		// Measure size of response if required.
		if !m.cfg.DisableMeasureSize {
			m.cfg.recorder.OTelObserveHTTPResponseSize(props, reporter.BytesWritten())
		}
	}()

	// Call the wrapped logic.
	next()
}

// Handler returns an measuring standard http.Handler.
func Handler(handlerID string, m Middleware, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wi := &responseWriterInterceptor{
			statusCode:     http.StatusOK,
			ResponseWriter: w,
		}
		reporter := &stdReporter{
			w: wi,
			r: r,
		}

		m.Measure(handlerID, reporter, func() {
			h.ServeHTTP(wi, r)
		})
	})
}

type stdReporter struct {
	w *responseWriterInterceptor
	r *http.Request
}

func (s *stdReporter) Method() string { return s.r.Method }

func (s *stdReporter) URLPath() string { return s.r.URL.Path }

func (s *stdReporter) StatusCode() int { return s.w.statusCode }

func (s *stdReporter) BytesWritten() int64 { return int64(s.w.bytesWritten) }

// responseWriterInterceptor is a simple wrapper to intercept set data on a
// ResponseWriter.
type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (w *responseWriterInterceptor) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterInterceptor) Write(p []byte) (int, error) {
	w.bytesWritten += len(p)
	return w.ResponseWriter.Write(p)
}

func (w *responseWriterInterceptor) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("type assertion failed http.ResponseWriter not a http.Hijacker")
	}
	return h.Hijack()
}

func (w *responseWriterInterceptor) Flush() {
	f, ok := w.ResponseWriter.(http.Flusher)
	if !ok {
		return
	}

	f.Flush()
}
