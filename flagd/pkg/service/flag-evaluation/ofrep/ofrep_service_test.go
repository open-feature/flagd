package ofrep

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/open-feature/flagd/core/pkg/evaluator"
	mock "github.com/open-feature/flagd/core/pkg/evaluator/mock"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	"go.uber.org/mock/gomock"
	"golang.org/x/sync/errgroup"
)

func Test_OfrepServiceStartStop(t *testing.T) {
	port := 18282
	eval := mock.NewMockIEvaluator(gomock.NewController(t))

	eval.EXPECT().ResolveAllValues(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]evaluator.AnyValue{}, model.Metadata{}, nil)

	cfg := SvcConfiguration{
		Logger:          logger.NewLogger(nil, false),
		Port:            uint16(port),
		ServiceName:     "test-service",
		MetricsRecorder: &telemetry.NoopMetricsRecorder{},
	}

	service, err := NewOfrepService(eval, []string{"*"}, cfg, nil, nil)
	if err != nil {
		t.Fatalf("error creating the ofrep service: %v", err)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return service.Start(gCtx)
	})

	// allow time for server startup
	<-time.After(2 * time.Second)

	path := fmt.Sprintf("http://localhost:%d/ofrep/v1/evaluate/flags", port)

	// validate response
	response, err := tryResponse(http.MethodPost, path, []byte{})
	if err != nil {
		t.Fatalf("error from server: %v", err)
	}

	if response == 0 {
		t.Fatal("expected non zero status")
	}

	// cancel the context
	cancelFunc()

	err = group.Wait()
	if err != nil {
		t.Errorf("error from service group: %v", err)
	}
}

func tryResponse(method string, uri string, payload []byte) (int, error) {
	return tryResponseWithHeaders(method, uri, payload, nil)
}

func tryResponseWithHeaders(method string, uri string, payload []byte, headers map[string]string) (int, error) {
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	request, err := http.NewRequest(method, uri, bytes.NewReader(payload))
	if err != nil {
		return 0, fmt.Errorf("error forming the request: %w", err)
	}

	for k, v := range headers {
		request.Header.Set(k, v)
	}

	rsp, err := client.Do(request)
	if err != nil {
		return 0, fmt.Errorf("error from the request: %w", err)
	}
	return rsp.StatusCode, nil
}

func Test_OfrepService_RequestBodySizeLimit(t *testing.T) {
	port := 18283
	eval := mock.NewMockIEvaluator(gomock.NewController(t))

	cfg := SvcConfiguration{
		Logger:              logger.NewLogger(nil, false),
		Port:                uint16(port),
		ServiceName:         "test-service",
		MetricsRecorder:     &telemetry.NoopMetricsRecorder{},
		MaxRequestBodyBytes: 10, // allow only 10 bytes
	}

	service, err := NewOfrepService(eval, []string{"*"}, cfg, nil, nil)
	if err != nil {
		t.Fatalf("error creating the ofrep service: %v", err)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return service.Start(gCtx)
	})
	defer func() {
		cancelFunc()
		_ = group.Wait()
	}()

	// wait for server startup
	<-time.After(2 * time.Second)

	path := fmt.Sprintf("http://localhost:%d/ofrep/v1/evaluate/flags/myFlag", port)
	// Valid JSON whose size exceeds the 10-byte limit, so MaxBytesReader triggers mid-parse.
	largeBody := []byte(`{"context":{"k":"` + strings.Repeat("a", 100) + `"}}`)

	status, err := tryResponse(http.MethodPost, path, largeBody)
	if err != nil {
		t.Fatalf("unexpected request error: %v", err)
	}

	if status != http.StatusRequestEntityTooLarge {
		t.Errorf("expected HTTP 413, got %d", status)
	}
}

func Test_OfrepService_RequestHeaderSizeLimit(t *testing.T) {
	port := 18284
	eval := mock.NewMockIEvaluator(gomock.NewController(t))

	cfg := SvcConfiguration{
		Logger:                logger.NewLogger(nil, false),
		Port:                  uint16(port),
		ServiceName:           "test-service",
		MetricsRecorder:       &telemetry.NoopMetricsRecorder{},
		MaxRequestHeaderBytes: 100, // 10000-byte test header value easily exceeds 100 + slop
	}

	service, err := NewOfrepService(eval, []string{"*"}, cfg, nil, nil)
	if err != nil {
		t.Fatalf("error creating the ofrep service: %v", err)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return service.Start(gCtx)
	})
	defer func() {
		cancelFunc()
		_ = group.Wait()
	}()

	// wait for server startup
	<-time.After(2 * time.Second)

	path := fmt.Sprintf("http://localhost:%d/ofrep/v1/evaluate/flags/myFlag", port)
	// The header value must exceed MaxHeaderBytes + Go's ~4096-byte read buffer slop.
	largeHeaderValue := string(bytes.Repeat([]byte("a"), 10000))

	status, err := tryResponseWithHeaders(http.MethodPost, path, []byte{}, map[string]string{
		"X-Large-Header": largeHeaderValue,
	})
	if err != nil {
		t.Fatalf("unexpected request error: %v", err)
	}

	if status != http.StatusRequestHeaderFieldsTooLarge {
		t.Errorf("expected HTTP 431, got %d", status)
	}
}
