package service

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
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
	// ObserveHTTPRequestDuration measures the duration of an HTTP request.
	ObserveHTTPRequestDuration(props HTTPReqProperties, duration time.Duration)
	// ObserveHTTPResponseSize measures the size of an HTTP response in bytes.
	ObserveHTTPResponseSize(props HTTPReqProperties, sizeBytes int64)
	// AddInflightRequests increments and decrements the number of inflight request being
	// processed.
	AddInflightRequests(props HTTPProperties, quantity int)
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

type MetricsRecorder struct {
	httpRequestDurHistogram   *prometheus.HistogramVec
	httpResponseSizeHistogram *prometheus.HistogramVec
	httpRequestsInflight      *prometheus.GaugeVec
}

func (r MetricsRecorder) ObserveHTTPRequestDuration(p HTTPReqProperties, duration time.Duration,
) {
	r.httpRequestDurHistogram.WithLabelValues(p.Service, p.ID, p.Method, p.Code).Observe(duration.Seconds())
}

func (r MetricsRecorder) ObserveHTTPResponseSize(p HTTPReqProperties, sizeBytes int64) {
	r.httpResponseSizeHistogram.WithLabelValues(p.Service, p.ID, p.Method, p.Code).Observe(float64(sizeBytes))
}

func (r MetricsRecorder) AddInflightRequests(p HTTPProperties, quantity int) {
	r.httpRequestsInflight.WithLabelValues(p.Service, p.ID).Add(float64(quantity))
}

type prometheusConfig struct {
	Prefix          string
	DurationBuckets []float64
	SizeBuckets     []float64
	Registry        prometheus.Registerer
	HandlerIDLabel  string
	StatusCodeLabel string
	MethodLabel     string
	ServiceLabel    string
}

type middlewareConfig struct {
	Recorder               Recorder
	Service                string
	GroupedStatus          bool
	DisableMeasureSize     bool
	DisableMeasureInflight bool
}

type Middleware struct {
	cfg middlewareConfig
}

func (c *middlewareConfig) defaults() {
	if c.Recorder == nil {
		panic("recorder is required")
	}
}

func New(cfg middlewareConfig) Middleware {
	cfg.defaults()

	m := Middleware{cfg: cfg}

	return m
}

func (c *prometheusConfig) defaults() {
	if len(c.DurationBuckets) == 0 {
		c.DurationBuckets = prometheus.DefBuckets
	}

	if len(c.SizeBuckets) == 0 {
		c.SizeBuckets = prometheus.ExponentialBuckets(100, 10, 8)
	}

	if c.Registry == nil {
		c.Registry = prometheus.DefaultRegisterer
	}

	if c.HandlerIDLabel == "" {
		c.HandlerIDLabel = "handler"
	}

	if c.StatusCodeLabel == "" {
		c.StatusCodeLabel = "code"
	}

	if c.MethodLabel == "" {
		c.MethodLabel = "method"
	}

	if c.ServiceLabel == "" {
		c.ServiceLabel = "service"
	}
}

func NewRecorder(cfg prometheusConfig) *MetricsRecorder {
	cfg.defaults()

	r := &MetricsRecorder{
		httpRequestDurHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: cfg.Prefix,
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "the latency of the HTTP requests",
			Buckets:   cfg.DurationBuckets,
		}, []string{cfg.ServiceLabel, cfg.HandlerIDLabel, cfg.MethodLabel, cfg.StatusCodeLabel}),

		httpResponseSizeHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: cfg.Prefix,
			Subsystem: "http",
			Name:      "response_size_bytes",
			Help:      "the size of the HTTP responses",
			Buckets:   cfg.SizeBuckets,
		}, []string{cfg.ServiceLabel, cfg.HandlerIDLabel, cfg.MethodLabel, cfg.StatusCodeLabel}),

		httpRequestsInflight: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: cfg.Prefix,
			Subsystem: "http",
			Name:      "requests_inflight",
			Help:      "the number of inflight requests being handled at the same time",
		}, []string{cfg.ServiceLabel, cfg.HandlerIDLabel}),
	}

	cfg.Registry.MustRegister(
		r.httpRequestDurHistogram,
		r.httpResponseSizeHistogram,
		r.httpRequestsInflight,
	)

	return r
}

func (m Middleware) Measure(handlerID string, reporter Reporter, next func()) {
	// If there isn't predefined handler ID we
	// set that ID as the URL path.
	hid := handlerID
	if handlerID == "" {
		hid = reporter.URLPath()
	}

	// Start the timer and when finishing measure the duration.
	start := time.Now()
	defer func() {
		duration := time.Since(start)

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
		m.cfg.Recorder.ObserveHTTPRequestDuration(props, duration)

		// Measure size of response if required.
		if !m.cfg.DisableMeasureSize {
			m.cfg.Recorder.ObserveHTTPResponseSize(props, reporter.BytesWritten())
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

func HandlerProvider(handlerID string, m Middleware) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return Handler(handlerID, m, next)
	}
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
