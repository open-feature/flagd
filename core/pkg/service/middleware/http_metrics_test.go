package middleware

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/otel"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/zap/zapcore"
)

func TestMiddlewareExposesMetrics(t *testing.T) {
	const svcName = "mySvc"
	exp := metric.NewManualReader()
	l, _ := logger.NewZapLogger(zapcore.DebugLevel, "")
	m := NewHttpMetric(Config{
		MetricRecorder: otel.NewOTelRecorder(exp, svcName),
		Service:        svcName,
		Logger:         logger.NewLogger(l, true),
	})
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("answer"))
	})
	svr := httptest.NewServer(Handler("id", m, handler))
	defer svr.Close()
	resp, err := http.Get(svr.URL)
	if err != nil {
		t.Errorf("Got %v", err)
	}
	_, _ = io.ReadAll(resp.Body)
	data, err := exp.Collect(context.TODO())
	if err != nil {
		t.Errorf("Got %v", err)
	}
	if len(data.ScopeMetrics) != 1 {
		t.Errorf("A single scope is expected, got %d", len(data.ScopeMetrics))
	}
	scopeMetrics := data.ScopeMetrics[0]
	if !reflect.DeepEqual(scopeMetrics.Scope.Name, svcName) {
		t.Errorf("Scope name %s, want %s", scopeMetrics.Scope.Name, svcName)
	}
	if len(scopeMetrics.Metrics) != 3 {
		t.Errorf("Expected 3 metrics, got %d", len(scopeMetrics.Metrics))
	}
}

func TestMeasure(t *testing.T) {
	exp := metric.NewManualReader()
	l, _ := logger.NewZapLogger(zapcore.DebugLevel, "")

	next := func() {}
	ctx := context.TODO()
	type expected struct {
		urlCalled    bool
		methodCalled bool
		statusCalled bool
		bytesCalled  bool
	}
	tests := []struct {
		name        string
		rep         MockReporter
		id          string
		groupStatus bool
		measureSize bool
		exp         expected
	}{
		{
			name: "empty id",
			id:   "",
			rep: MockReporter{
				URL:    "myURL",
				Meth:   "GET",
				Status: 100,
				Bytes:  0,
			},
			groupStatus: true,
			measureSize: true,
			exp: expected{
				urlCalled:    true,
				methodCalled: true,
				statusCalled: true,
				bytesCalled:  true,
			},
		},
		{
			name: "id provided",
			id:   "mySpecialHandler",
			rep: MockReporter{
				URL:    "myURL",
				Meth:   "GET",
				Status: 100,
				Bytes:  0,
			},
			groupStatus: true,
			measureSize: true,
			exp: expected{
				urlCalled:    false,
				methodCalled: true,
				statusCalled: true,
				bytesCalled:  true,
			},
		},
		{
			name: "id provided - no report of size",
			id:   "mySpecialHandler",
			rep: MockReporter{
				URL:    "myURL",
				Meth:   "GET",
				Status: 100,
				Bytes:  0,
			},
			groupStatus: true,
			measureSize: false,
			exp: expected{
				urlCalled:    false,
				methodCalled: true,
				statusCalled: true,
				bytesCalled:  false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// test the middleware correctly
			rep := tt.rep
			m := NewHttpMetric(Config{
				MetricRecorder:     otel.NewOTelRecorder(exp, tt.name),
				Service:            tt.name,
				Logger:             logger.NewLogger(l, true),
				GroupedStatus:      tt.groupStatus,
				DisableMeasureSize: !tt.measureSize,
			})
			m.Measure(ctx, tt.id, &rep, next)
			if rep.urlCalled != tt.exp.urlCalled {
				t.Errorf("Expected %v for URLPath but got %v", tt.exp.urlCalled, rep.urlCalled)
			}
			if rep.methodCalled != tt.exp.methodCalled {
				t.Errorf("Expected %v for Method but got %v", tt.exp.methodCalled, rep.methodCalled)
			}
			if rep.statusCalled != tt.exp.statusCalled {
				t.Errorf("Expected %v for StatusCode but got %v", tt.exp.statusCalled, rep.statusCalled)
			}
			if rep.bytesCalled != tt.exp.bytesCalled {
				t.Errorf("Expected %v for BytesWritten but got %v", tt.exp.bytesCalled, rep.bytesCalled)
			}
		})
	}
}

type MockReporter struct {
	URL          string
	Meth         string
	Status       int
	Bytes        int64
	urlCalled    bool
	methodCalled bool
	statusCalled bool
	bytesCalled  bool
}

func (m *MockReporter) Method() string {
	m.methodCalled = true
	return m.Meth
}

func (m *MockReporter) URLPath() string {
	m.urlCalled = true
	return m.URL
}

func (m *MockReporter) StatusCode() int {
	m.statusCalled = true
	return m.Status
}

func (m *MockReporter) BytesWritten() int64 {
	m.bytesCalled = true
	return m.Bytes
}

func (m *MockReporter) URLCalled() bool    { return m.urlCalled }
func (m *MockReporter) MethodCalled() bool { return m.methodCalled }
func (m *MockReporter) StatusCalled() bool { return m.statusCalled }
func (m *MockReporter) BytesCalled() bool  { return m.bytesCalled }

func TestNewHttpMetric(t *testing.T) {
	l, _ := logger.NewZapLogger(zapcore.DebugLevel, "")
	log := logger.NewLogger(l, true)
	exp := metric.NewManualReader()
	const svcName = "mySvc"
	const groupedStatus = false
	const disableMeasureSize = false

	mdw := NewHttpMetric(Config{
		MetricRecorder:     otel.NewOTelRecorder(exp, svcName),
		Logger:             log,
		Service:            svcName,
		GroupedStatus:      groupedStatus,
		DisableMeasureSize: disableMeasureSize,
	})
	if mdw.cfg.MetricRecorder == nil {
		t.Errorf("Expected OpenTelemetry to be configured, got nil")
	}
	if mdw.cfg.Logger == nil {
		t.Errorf("Expected logger to be configured, got nil")
	}
	if mdw.cfg.Service != svcName {
		t.Errorf("Expected Service to be configured with %s, got %s", svcName, mdw.cfg.Service)
	}
	if mdw.cfg.GroupedStatus != groupedStatus {
		t.Errorf("Expected GroupedStatus to be configured with %v, got %v", groupedStatus, mdw.cfg.GroupedStatus)
	}
	if mdw.cfg.DisableMeasureSize != disableMeasureSize {
		t.Errorf("Expected DisableMeasureSize to be configured with %v, got %v", disableMeasureSize, mdw.cfg.DisableMeasureSize)
	}
}
