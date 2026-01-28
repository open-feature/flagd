package ofrep

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"github.com/open-feature/flagd/core/pkg/evaluator"
	mock "github.com/open-feature/flagd/core/pkg/evaluator/mock"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/service/ofrep"
	"go.uber.org/mock/gomock"
)

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
				t.Errorf("error unmarshelling body: %v", err)
			}

			if !reflect.DeepEqual(test.payload, rsp) {
				t.Errorf("incorrect payload in wire")
			}
		})
	}
}
func TestFlagdContextInvalidContextType(t *testing.T) {
	log := logger.NewLogger(nil, false)

	result := flagdContext(
		log,
		"test-request-id",
		ofrep.Request{Context: "not a map"}, // invalid: string instead of map
		map[string]any{"staticKey": "staticValue"},
		http.Header{},
		map[string]string{},
	)

	if val, exists := result["staticKey"]; !exists || val != "staticValue" {
		t.Errorf("expected static context to be included even with invalid request context")
	}
}

func TestFlagdContextDelegatesContextMerging(t *testing.T) {
	log := logger.NewLogger(nil, false)

	h := http.Header{}
	h.Set("X-User-Tier", "premium")

	result := flagdContext(
		log,
		"test-request-id",
		ofrep.Request{Context: map[string]any{"requestKey": "requestValue"}},
		map[string]any{"staticKey": "staticValue"},
		h,
		map[string]string{"X-User-Tier": "userTier"},
	)

	expected := map[string]any{
		"requestKey": "requestValue",
		"staticKey":  "staticValue",
		"userTier":   "premium",
	}

	for k, v := range expected {
		if result[k] != v {
			t.Errorf("expected key '%s' to have value '%s', but got '%v'", k, v, result[k])
		}
	}
}
