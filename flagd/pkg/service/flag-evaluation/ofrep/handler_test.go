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
	h := handler{Logger: log}

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

	tests := []struct {
		name            string
		eTagGenerator   func() (string, error) // Optional function to generate the ETag
		expectedStatus  int
		expectedHasETag bool
		expectedHasBody bool
	}{
		{
			name:            "no If-None-Match header returns 200 with body and ETag",
			eTagGenerator:   nil, // Will use no ETag in request
			expectedStatus:  http.StatusOK,
			expectedHasETag: true,
			expectedHasBody: true,
		},
		{
			name: "matching If-None-Match header returns 304 Not Modified",
			eTagGenerator: func() (string, error) {
				return calculateETag(response)
			},
			expectedStatus:  http.StatusNotModified,
			expectedHasETag: true,
			expectedHasBody: false,
		},
		{
			name: "non-matching If-None-Match header returns 200 with body and ETag",
			eTagGenerator: func() (string, error) {
				return "\"some-invalid-etag-lmao\"", nil
			},
			expectedStatus:  http.StatusOK,
			expectedHasETag: true,
			expectedHasBody: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Calculate the ETag for this specific test case
			var ifNoneMatchHeader string
			if test.eTagGenerator != nil {
				eTag, err := test.eTagGenerator()
				if err != nil {
					t.Fatalf("error generating ETag: %v", err)
				}
				ifNoneMatchHeader = eTag
			}

			request := httptest.NewRequest(http.MethodPost, "/ofrep/v1/evaluate/flags", bytes.NewReader([]byte{}))
			if ifNoneMatchHeader != "" {
				request.Header.Set("If-None-Match", ifNoneMatchHeader)
			}

			recorder := httptest.NewRecorder()
			h.writeBulkEvaluationResponse(recorder, request, response)

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

			// Verify ETag behavior: if matching ETag in request, should get 304
			if ifNoneMatchHeader != "" && eTagHeader != "" && ifNoneMatchHeader == eTagHeader {
				if test.expectedStatus != http.StatusNotModified {
					t.Errorf("expected 304 status when ETag matches, but got %d", test.expectedStatus)
				}
			}

			// For 200 responses, verify Content-Type header is present
			if test.expectedStatus == http.StatusOK && recorder.Header().Get("Content-Type") != "application/json" {
				t.Error("expected Content-Type header to be application/json, but it was missing")
			}
		})
	}
}

func TestCalculateETag(t *testing.T) {
	tests := []struct {
		name           string
		response       ofrep.BulkEvaluationResponse
		expectedErrNil bool
	}{
		{
			name: "valid response produces ETag",
			response: ofrep.BulkEvaluationResponse{
				Flags: []interface{}{
					ofrep.EvaluationSuccess{
						Key:     "test-flag",
						Value:   true,
						Reason:  model.StaticReason,
						Variant: "on",
					},
				},
				Metadata: model.Metadata{},
			},
			expectedErrNil: true,
		},
		{
			name: "empty response produces ETag",
			response: ofrep.BulkEvaluationResponse{
				Flags:    []interface{}{},
				Metadata: model.Metadata{},
			},
			expectedErrNil: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			eTag, err := calculateETag(test.response)

			if test.expectedErrNil && err != nil {
				t.Errorf("expected no error, but got: %v", err)
			}
			if !test.expectedErrNil && err == nil {
				t.Error("expected error, but got none")
			}

			if test.expectedErrNil {
				// Verify ETag format: quoted hex string
				if len(eTag) < 2 || eTag[0] != '"' || eTag[len(eTag)-1] != '"' {
					t.Errorf("expected ETag to be quoted, but got: %s", eTag)
				}

				// Calculate again to ensure deterministic output
				eTag2, _ := calculateETag(test.response)
				if eTag != eTag2 {
					t.Errorf("expected deterministic ETag, but got different values: %s vs %s", eTag, eTag2)
				}
			}
		})
	}
}
