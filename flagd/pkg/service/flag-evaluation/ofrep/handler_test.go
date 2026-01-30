package ofrep

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/open-feature/flagd/core/pkg/evaluator"
	mock "github.com/open-feature/flagd/core/pkg/evaluator/mock"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/service/ofrep"
	"github.com/open-feature/flagd/core/pkg/store"
	"go.uber.org/mock/gomock"
)

// testFlagStore is a mock FlagStore for handler tests
type testFlagStore struct {
	flags   []model.Flag
	watchCh chan struct{}
}

func newTestFlagStore() *testFlagStore {
	return &testFlagStore{
		flags:   []model.Flag{{Key: "test-flag", State: "ENABLED"}},
		watchCh: make(chan struct{}),
	}
}

func (m *testFlagStore) GetAll(_ context.Context, _ *store.Selector) ([]model.Flag, model.Metadata, error) {
	return m.flags, model.Metadata{}, nil
}

func (m *testFlagStore) WatchSelector(_ *store.Selector) <-chan struct{} {
	return m.watchCh
}

var flagKey = "key"

var successValue = evaluator.AnyValue{
	Value:    true,
	Variant:  "true",
	Reason:   model.StaticReason,
	FlagKey:  flagKey,
	Metadata: nil,
	Error:    nil,
}

var flagNotFoundValue = evaluator.AnyValue{
	Reason:  model.ErrorReason,
	FlagKey: flagKey,
	Error:   errors.New(model.FlagNotFoundErrorCode),
}

var genericErrorValue = evaluator.AnyValue{
	Reason:  model.ErrorReason,
	FlagKey: flagKey,
	Error:   errors.New(model.GeneralErrorCode),
}

func Test_handler_HandleFlagEvaluation(t *testing.T) {
	log := logger.NewLogger(nil, false)

	tests := []struct {
		name string

		method          string
		path            string
		input           *bytes.Reader
		mockAnyResponse *evaluator.AnyValue

		expectedStatus       int
		expectedResponseType interface{}
	}{
		{
			name:                 "success evaluation",
			method:               http.MethodPost,
			path:                 "/ofrep/v1/evaluate/flags/" + flagKey,
			input:                bytes.NewReader([]byte{}),
			mockAnyResponse:      &successValue,
			expectedStatus:       http.StatusOK,
			expectedResponseType: ofrep.EvaluationSuccess{},
		},
		{
			name:                 "valid context and success",
			method:               http.MethodPost,
			path:                 "/ofrep/v1/evaluate/flags/" + flagKey,
			input:                bytes.NewReader([]byte("{\"context\": {}}")),
			mockAnyResponse:      &successValue,
			expectedStatus:       http.StatusOK,
			expectedResponseType: ofrep.EvaluationSuccess{},
		},
		{
			name:                 "flag not found evaluation",
			method:               http.MethodPost,
			path:                 "/ofrep/v1/evaluate/flags/" + flagKey,
			input:                bytes.NewReader([]byte{}),
			mockAnyResponse:      &flagNotFoundValue,
			expectedStatus:       http.StatusNotFound,
			expectedResponseType: ofrep.EvaluationError{},
		},
		{
			name:                 "general error evaluation",
			method:               http.MethodPost,
			path:                 "/ofrep/v1/evaluate/flags/" + flagKey,
			input:                bytes.NewReader([]byte{}),
			mockAnyResponse:      &genericErrorValue,
			expectedStatus:       http.StatusBadRequest,
			expectedResponseType: ofrep.EvaluationError{},
		},
		{
			name:                 "flag key parsing error - whitespace",
			method:               http.MethodPost,
			path:                 "/ofrep/v1/evaluate/flags/ ",
			mockAnyResponse:      &successValue,
			input:                bytes.NewReader([]byte{}),
			expectedStatus:       http.StatusOK,
			expectedResponseType: ofrep.EvaluationSuccess{},
		},
		{
			name:                 "invalid context payload",
			method:               http.MethodPost,
			path:                 "/ofrep/v1/evaluate/flags/" + flagKey,
			input:                bytes.NewReader([]byte("{some invalid context}")),
			expectedStatus:       http.StatusBadRequest,
			expectedResponseType: ofrep.EvaluationError{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			eval := mock.NewMockIEvaluator(gomock.NewController(t))
			if test.mockAnyResponse != nil {
				eval.EXPECT().
					ResolveAsAnyValue(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(*test.mockAnyResponse)
			}

			h := handler{Logger: log, evaluator: eval}

			request, err := http.NewRequest(test.method, test.path, test.input)
			if err != nil {
				t.Fatalf("error setting up request: %v", err)
			}

			recorder := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc(singleEvaluation, h.HandleFlagEvaluation)
			router.ServeHTTP(recorder, request)

			if test.expectedStatus != recorder.Code {
				t.Errorf("expected status code %d, but got %d", test.expectedStatus, recorder.Code)
			}

			output := test.expectedResponseType
			err = json.NewDecoder(recorder.Result().Body).Decode(&output)
			if err != nil {
				t.Errorf("error parsing response to expected response: %v", err)
			}
		})
	}
}

func Test_handler_HandleBulkEvaluation(t *testing.T) {
	log := logger.NewLogger(nil, false)

	tests := []struct {
		name string

		method          string
		input           *bytes.Reader
		mockAnyResponse []evaluator.AnyValue
		mockAnyMetadata model.Metadata
		mockAnyError    error

		expectedStatus int
	}{
		{
			name:            "success evaluation",
			method:          http.MethodPost,
			input:           bytes.NewReader([]byte{}),
			mockAnyResponse: []evaluator.AnyValue{successValue},
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "success & evaluation errors",
			method:          http.MethodPost,
			input:           bytes.NewReader([]byte{}),
			mockAnyResponse: []evaluator.AnyValue{successValue, genericErrorValue, flagNotFoundValue},
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "only evaluation errors",
			method:          http.MethodPost,
			input:           bytes.NewReader([]byte{}),
			mockAnyResponse: []evaluator.AnyValue{genericErrorValue, flagNotFoundValue},
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "handles internal errors and yield 500",
			method:          "http.MethodPost",
			input:           bytes.NewReader([]byte{}),
			mockAnyResponse: []evaluator.AnyValue{},
			mockAnyError:    errors.New("some internal error from evaluator"),
			expectedStatus:  http.StatusInternalServerError,
		},
		{
			name:            "valid context payload",
			method:          http.MethodPost,
			input:           bytes.NewReader([]byte("{\"context\": {}}")),
			mockAnyResponse: []evaluator.AnyValue{},
			expectedStatus:  http.StatusOK,
		},
		{
			name:           "invalid context payload",
			method:         http.MethodPost,
			input:          bytes.NewReader([]byte("{some invalid context}")),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			eval := mock.NewMockIEvaluator(gomock.NewController(t))
			eval.EXPECT().ResolveAllValues(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(test.mockAnyResponse, test.mockAnyMetadata, test.mockAnyError).MinTimes(0)

			h := handler{Logger: log, evaluator: eval}

			request, err := http.NewRequest(test.method, "/ofrep/v1/evaluate/flags", test.input)
			if err != nil {
				t.Fatalf("error setting up request: %v", err)
			}

			recorder := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc(bulkEvaluation, h.HandleBulkEvaluation)
			router.ServeHTTP(recorder, request)

			if test.expectedStatus != recorder.Code {
				t.Errorf("expected status code %d, but got %d", test.expectedStatus, recorder.Code)
			}
		})
	}
}

func TestWriteJSONResponse(t *testing.T) {
	log := logger.NewLogger(nil, false)
	h := handler{Logger: log}

	tests := []struct {
		name           string
		status         int
		payload        interface{}
		expectedStatus int
	}{
		{
			name:           "success",
			status:         http.StatusOK,
			payload:        successValue,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "marshaling error",
			status:         http.StatusOK,
			payload:        func() {}, // make marshaling to fail
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			h.writeJSONToResponse(test.status, test.payload, recorder)

			if test.expectedStatus != recorder.Code {
				t.Errorf("expected status code %d, but got %d", test.expectedStatus, recorder.Code)
			}

			if test.expectedStatus == http.StatusOK && recorder.Header().Get("Content-Type") != "application/json" {
				t.Error("expected http OK to contain header application/json content type, but header is missing")
			}

			if test.expectedStatus != http.StatusOK {
				// rest of the validations are only for status OK
				return
			}

			b, err := io.ReadAll(recorder.Body)
			if err != nil {
				t.Errorf("error deriving body: %v", err)
			}

			var rsp evaluator.AnyValue
			err = json.Unmarshal(b, &rsp)
			if err != nil {
				t.Errorf("error unmarshaling body: %v", err)
			}

			if !reflect.DeepEqual(test.payload, rsp) {
				t.Errorf("incorrect payload in wire")
			}
		})
	}
}

func TestWriteBulkEvaluationResponse_ETag(t *testing.T) {
	log := logger.NewLogger(nil, false)

	// create version tracker with mock store
	mockStore := newTestFlagStore()
	tracker := NewSelectorVersionTracker(log, mockStore, 0)
	defer tracker.Close()

	h := handler{Logger: log, versionTracker: tracker}

	// test response
	response := ofrep.BulkEvaluationResponse{
		Flags: []interface{}{
			ofrep.EvaluationSuccess{
				Key:     "test-flag",
				Value:   true,
				Reason:  model.StaticReason,
				Variant: "on",
			},
		},
		Metadata: model.Metadata{},
	}

	selectorExpression := "flagSetId=test-set"

	tests := []struct {
		name            string
		ifNoneMatch     string
		expectedStatus  int
		expectedHasETag bool
		expectedHasBody bool
	}{
		{
			name:            "no If-None-Match header returns 200 with body and ETag",
			ifNoneMatch:     "",
			expectedStatus:  http.StatusOK,
			expectedHasETag: true,
			expectedHasBody: true,
		},
		{
			name:            "non-matching If-None-Match header returns 200 with body and ETag",
			ifNoneMatch:     "\"some-invalid-etag-lmao\"",
			expectedStatus:  http.StatusOK,
			expectedHasETag: true,
			expectedHasBody: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/ofrep/v1/evaluate/flags", bytes.NewReader([]byte{}))
			if test.ifNoneMatch != "" {
				request.Header.Set("If-None-Match", test.ifNoneMatch)
			}

			recorder := httptest.NewRecorder()
			h.writeBulkEvaluationResponse(recorder, request, selectorExpression, response)

			if test.expectedStatus != recorder.Code {
				t.Errorf("expected status code %d, but got %d", test.expectedStatus, recorder.Code)
			}

			eTagHeader := recorder.Header().Get("ETag")
			if test.expectedHasETag && eTagHeader == "" {
				t.Error("expected ETag header to be present, but it was missing")
			}
			if !test.expectedHasETag && eTagHeader != "" {
				t.Error("expected ETag header to be absent, but it was present")
			}

			body := recorder.Body.String()
			if test.expectedHasBody && body == "" {
				t.Error("expected response body, but got empty body")
			}
			if !test.expectedHasBody && body != "" {
				t.Errorf("expected no response body, but got: %s", body)
			}

			// for 200 responses, verify Content-Type header is present
			if test.expectedStatus == http.StatusOK && recorder.Header().Get("Content-Type") != "application/json" {
				t.Error("expected Content-Type header to be application/json, but it was missing")
			}
		})
	}

	// test matching If-None-Match - the ETag should match since Track is idempotent
	t.Run("matching If-None-Match header still returns 200 (no 304 in writeBulkEvaluationResponse)", func(t *testing.T) {
		// get the current ETag for this selector
		currentETag := tracker.ETag(selectorExpression)

		request := httptest.NewRequest(http.MethodPost, "/ofrep/v1/evaluate/flags", bytes.NewReader([]byte{}))
		request.Header.Set("If-None-Match", currentETag)

		recorder := httptest.NewRecorder()
		h.writeBulkEvaluationResponse(recorder, request, selectorExpression, response)
	})
}

func TestHandleBulkEvaluation_304NotModified(t *testing.T) {
	log := logger.NewLogger(nil, false)

	mockStore := newTestFlagStore()
	tracker := NewSelectorVersionTracker(log, mockStore, 0)
	defer tracker.Close()

	// pre-track the empty selector
	selectorExpression := ""
	cachedETag := tracker.Track(selectorExpression)

	// create handler with version tracker - evaluator should NOT be called for 304
	eval := mock.NewMockIEvaluator(gomock.NewController(t))

	h := handler{Logger: log, evaluator: eval, versionTracker: tracker}

	// request WITH matching If-None-Match header - should return 304
	request := httptest.NewRequest(http.MethodPost, "/ofrep/v1/evaluate/flags", bytes.NewReader([]byte{}))
	request.Header.Set("If-None-Match", cachedETag)
	recorder := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc(bulkEvaluation, h.HandleBulkEvaluation)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotModified {
		t.Errorf("expected status 304, got %d", recorder.Code)
	}

	if recorder.Header().Get("ETag") == "" {
		t.Error("expected ETag header to be present")
	}
}

func TestHandleBulkEvaluation_NewSelector(t *testing.T) {
	log := logger.NewLogger(nil, false)

	// create version tracker with mock store
	mockStore := newTestFlagStore()
	tracker := NewSelectorVersionTracker(log, mockStore, 0)
	defer tracker.Close()

	// create handler with version tracker
	eval := mock.NewMockIEvaluator(gomock.NewController(t))
	eval.EXPECT().ResolveAllValues(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]evaluator.AnyValue{successValue}, model.Metadata{}, nil)

	h := handler{Logger: log, evaluator: eval, versionTracker: tracker}

	request := httptest.NewRequest(http.MethodPost, "/ofrep/v1/evaluate/flags", bytes.NewReader([]byte{}))
	recorder := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc(bulkEvaluation, h.HandleBulkEvaluation)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", recorder.Code)
	}

	if recorder.Header().Get("ETag") == "" {
		t.Error("expected ETag header to be present")
	}

	etag := tracker.ETag("")
	if etag == "" {
		t.Error("expected selector to be tracked")
	}
}

func TestVersionBumpOnStoreUpdate(t *testing.T) {
	log := logger.NewLogger(nil, false)

	// create a real store
	flagStore, err := store.NewStore(log, []string{"test-source"})
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	// add initial flags
	flagStore.Update("test-source", []model.Flag{
		{
			Key:            "test-flag",
			State:          "ENABLED",
			DefaultVariant: "on",
			Variants:       map[string]any{"on": true, "off": false},
		},
	}, model.Metadata{})

	// create version tracker with watch provider (the store)
	tracker := NewSelectorVersionTracker(log, flagStore, 0)
	defer tracker.Close()

	// track the empty selector (matches all flags)
	etagBefore := tracker.Track("")

	if etagBefore == "" {
		t.Fatal("expected non-empty ETag after tracking")
	}

	// update the store with different flag content
	flagStore.Update("test-source", []model.Flag{
		{
			Key:            "test-flag",
			State:          "ENABLED",
			DefaultVariant: "off", // changed!
			Variants:       map[string]any{"on": true, "off": false},
		},
	}, model.Metadata{})

	// give time for watch goroutine to process the update
	time.Sleep(50 * time.Millisecond)

	// ETag should have changed because content changed
	etagAfter := tracker.ETag("")

	if etagAfter == "" {
		t.Fatal("expected non-empty ETag after update")
	}

	if etagBefore == etagAfter {
		t.Errorf("expected ETag to change after store update, but got same value: %s", etagBefore)
	}
}
