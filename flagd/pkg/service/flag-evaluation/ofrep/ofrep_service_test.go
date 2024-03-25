package ofrep

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	mock "github.com/open-feature/flagd/core/pkg/evaluator/mock"
	"github.com/open-feature/flagd/core/pkg/logger"
	"go.uber.org/mock/gomock"
	"golang.org/x/sync/errgroup"
)

func Test_OfrepServiceStartStop(t *testing.T) {
	port := 18282
	eval := mock.NewMockIEvaluator(gomock.NewController(t))
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

	path, err := url.JoinPath(fmt.Sprintf("http://localhost:%d", port), bulkEvaluation)
	if err != nil {
		t.Fatalf("error creating the path: %v", err)
	}
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
		return 0, fmt.Errorf("error forming the request: %v", err)
	}

	rsp, err := client.Do(request)
	return rsp.StatusCode, fmt.Errorf("error from the request: %v", err)
}
