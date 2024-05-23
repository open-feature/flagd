package ofrep

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/open-feature/flagd/core/pkg/evaluator"
	mock "github.com/open-feature/flagd/core/pkg/evaluator/mock"
	"github.com/open-feature/flagd/core/pkg/logger"
	"go.uber.org/mock/gomock"
	"golang.org/x/sync/errgroup"
)

func Test_OfrepServiceStartStop(t *testing.T) {
	port := 18282
	eval := mock.NewMockIEvaluator(gomock.NewController(t))

	eval.EXPECT().ResolveAllValues(gomock.Any(), gomock.Any(), gomock.Any()).Return([]evaluator.AnyValue{})
	cfg := SvcConfiguration{
		Logger: logger.NewLogger(nil, false),
		Port:   uint16(port),
	}

	service, err := NewOfrepService(eval, []string{"*"}, cfg)
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
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	request, err := http.NewRequest(method, uri, bytes.NewReader(payload))
	if err != nil {
		return 0, fmt.Errorf("error forming the request: %w", err)
	}

	rsp, err := client.Do(request)
	if err != nil {
		return 0, fmt.Errorf("error from the request: %w", err)
	}
	return rsp.StatusCode, nil
}
