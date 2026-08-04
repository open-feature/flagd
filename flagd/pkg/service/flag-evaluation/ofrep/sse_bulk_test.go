package ofrep

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/open-feature/flagd/core/pkg/evaluator"
	mock "github.com/open-feature/flagd/core/pkg/evaluator/mock"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/service/ofrep"
	"github.com/open-feature/flagd/core/pkg/store"
	svc "github.com/open-feature/flagd/flagd/pkg/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type fakeVersioner struct {
	etag         string
	lastModified int64
	ok           bool
}

func (f fakeVersioner) Version(_ store.Selector) (string, int64, bool) {
	return f.etag, f.lastModified, f.ok
}

func newBulkRequest(t *testing.T, selectorHeader, clientEtag string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, "/ofrep/v1/evaluate/flags", bytes.NewReader([]byte{}))
	require.NoError(t, err)
	req.Host = "flagd.example:8016"
	if selectorHeader != "" {
		req.Header.Set(svc.FLAGD_SELECTOR_HEADER, selectorHeader)
	}
	if clientEtag != "" {
		req.Header.Set("If-None-Match", clientEtag)
	}
	return req
}

func serveBulk(h handler, req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc(bulkEvaluation, h.HandleBulkEvaluation)
	router.ServeHTTP(recorder, req)
	return recorder
}

func TestHandleBulkEvaluation_NotModified(t *testing.T) {
	log := logger.NewLogger(nil, false)
	eval := mock.NewMockIEvaluator(gomock.NewController(t))
	// evaluation must be skipped on a 304
	eval.EXPECT().ResolveAllValues(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	h := handler{
		Logger:                log,
		evaluator:             eval,
		versioner:             fakeVersioner{etag: "abc123", ok: true},
		sseEnabled:            true,
		sseInactivityDelaySec: 120,
	}

	recorder := serveBulk(h, newBulkRequest(t, "flagSetId=fs1", `"abc123"`))

	assert.Equal(t, http.StatusNotModified, recorder.Code)
	assert.Equal(t, `"abc123"`, recorder.Header().Get("ETag"))
}

func TestHandleBulkEvaluation_AdvertisesEventStreams(t *testing.T) {
	log := logger.NewLogger(nil, false)

	tests := []struct {
		name           string
		publicURL      string
		selectorHeader string
		wantOrigin     string // expected endpoint.origin ("" => omitted)
		wantRequestURI string
	}{
		{
			name:           "flagSetId selector, origin omitted",
			selectorHeader: "flagSetId=fs1",
			wantOrigin:     "",
			wantRequestURI: "/ofrep/v1/sse?channels=fs1",
		},
		{
			name:           "no selector (catch-all), origin omitted",
			selectorHeader: "",
			wantOrigin:     "",
			wantRequestURI: "/ofrep/v1/sse",
		},
		{
			name:           "public url sets origin",
			publicURL:      "https://flags.example.com/",
			selectorHeader: "flagSetId=fs1",
			wantOrigin:     "https://flags.example.com",
			wantRequestURI: "/ofrep/v1/sse?channels=fs1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eval := mock.NewMockIEvaluator(gomock.NewController(t))
			eval.EXPECT().ResolveAllValues(gomock.Any(), gomock.Any(), gomock.Any()).
				Return([]evaluator.AnyValue{}, model.Metadata{}, nil)

			h := handler{
				Logger:                log,
				evaluator:             eval,
				versioner:             fakeVersioner{etag: "etag-1", ok: true},
				sseEnabled:            true,
				sseInactivityDelaySec: 120,
				ssePublicURL:          tt.publicURL,
			}

			recorder := serveBulk(h, newBulkRequest(t, tt.selectorHeader, ""))

			require.Equal(t, http.StatusOK, recorder.Code)
			assert.Equal(t, `"etag-1"`, recorder.Header().Get("ETag"))

			var resp ofrep.BulkEvaluationResponse
			require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
			require.Len(t, resp.EventStreams, 1)

			es := resp.EventStreams[0]
			assert.Equal(t, "sse", es.Type)
			assert.Equal(t, 120, es.InactivityDelaySec)
			require.NotNil(t, es.Endpoint)
			assert.Equal(t, tt.wantOrigin, es.Endpoint.Origin)
			assert.Equal(t, tt.wantRequestURI, es.Endpoint.RequestUri)
		})
	}
}

func TestHandleBulkEvaluation_SSEDisabled_NoEventStreams(t *testing.T) {
	log := logger.NewLogger(nil, false)
	eval := mock.NewMockIEvaluator(gomock.NewController(t))
	eval.EXPECT().ResolveAllValues(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]evaluator.AnyValue{}, model.Metadata{}, nil)

	// sseEnabled false, no versioner: behaves like the legacy handler
	h := handler{Logger: log, evaluator: eval}

	recorder := serveBulk(h, newBulkRequest(t, "", ""))

	require.Equal(t, http.StatusOK, recorder.Code)
	assert.Empty(t, recorder.Header().Get("ETag"))

	var resp ofrep.BulkEvaluationResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	assert.Empty(t, resp.EventStreams)
}
