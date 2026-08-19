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

func (f fakeVersioner) Version(_ string) (string, int64, bool) {
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

// The 304 decision must use If-None-Match (the client's cached version), not the ADR-0008
// flagConfigEtag change-trigger metadata.
func TestHandleBulkEvaluation_ConditionalUsesIfNoneMatch(t *testing.T) {
	log := logger.NewLogger(nil, false)
	const current = "etag-v2"

	tests := []struct {
		name           string
		flagConfigEtag string // query param (change-trigger metadata)
		ifNoneMatch    string // header (cache validator)
		wantStatus     int
		wantEval       bool
	}{
		{
			name:           "stale cache echoing new trigger etag is served fresh (regression)",
			flagConfigEtag: "etag-v2",
			ifNoneMatch:    `"etag-v1"`,
			wantStatus:     http.StatusOK,
			wantEval:       true,
		},
		{
			name:        "current cache validator yields 304",
			ifNoneMatch: `"etag-v2"`,
			wantStatus:  http.StatusNotModified,
			wantEval:    false,
		},
		{
			name:           "trigger etag alone (no validator) is served fresh",
			flagConfigEtag: "etag-v2",
			wantStatus:     http.StatusOK,
			wantEval:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eval := mock.NewMockIEvaluator(gomock.NewController(t))
			expect := eval.EXPECT().ResolveAllValues(gomock.Any(), gomock.Any(), gomock.Any()).
				Return([]evaluator.AnyValue{}, model.Metadata{}, nil)
			if tt.wantEval {
				expect.Times(1)
			} else {
				expect.Times(0)
			}

			h := handler{
				Logger:     log,
				evaluator:  eval,
				versioner:  fakeVersioner{etag: current, ok: true},
				sseEnabled: true,
			}

			url := "/ofrep/v1/evaluate/flags"
			if tt.flagConfigEtag != "" {
				url += "?flagConfigEtag=" + tt.flagConfigEtag
			}
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte{}))
			require.NoError(t, err)
			req.Header.Set(svc.FLAGD_SELECTOR_HEADER, "flagSetId=fs1")
			if tt.ifNoneMatch != "" {
				req.Header.Set("If-None-Match", tt.ifNoneMatch)
			}

			recorder := serveBulk(h, req)

			assert.Equal(t, tt.wantStatus, recorder.Code)
			assert.Equal(t, `"`+current+`"`, recorder.Header().Get("ETag"))
		})
	}
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
			wantRequestURI: "/ofrep/v1/sse?channels=flagSetId%3Dfs1",
		},
		{
			// previously fell back to the catch-all channel, and so was woken by every change
			name:           "source selector gets its own channel",
			selectorHeader: "source=src1",
			wantOrigin:     "",
			wantRequestURI: "/ofrep/v1/sse?channels=source%3Dsrc1",
		},
		{
			// advertised verbatim, so the SSE endpoint parses it back to the same selector
			name:           "bare selector expression is advertised verbatim",
			selectorHeader: "mySource",
			wantOrigin:     "",
			wantRequestURI: "/ofrep/v1/sse?channels=mySource",
		},
		{
			// "flags with no flagSetId"; the internal nil flagSetId must never be advertised
			name:           "empty flagSetId selector",
			selectorHeader: "flagSetId=",
			wantOrigin:     "",
			wantRequestURI: "/ofrep/v1/sse?channels=flagSetId%3D",
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
			wantRequestURI: "/ofrep/v1/sse?channels=flagSetId%3Dfs1",
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

// TestHandleBulkEvaluation_NoLiveStream_NoETag pins the contract that follows from tracking
// versions per live subscription: with no stream open there is no cache validator to offer, so
// the response is a plain 200 that still advertises eventStreams.
func TestHandleBulkEvaluation_NoLiveStream_NoETag(t *testing.T) {
	log := logger.NewLogger(nil, false)
	eval := mock.NewMockIEvaluator(gomock.NewController(t))
	eval.EXPECT().ResolveAllValues(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]evaluator.AnyValue{}, model.Metadata{}, nil)

	h := handler{
		Logger:                log,
		evaluator:             eval,
		versioner:             fakeVersioner{ok: false},
		sseEnabled:            true,
		sseInactivityDelaySec: 120,
	}

	recorder := serveBulk(h, newBulkRequest(t, "flagSetId=fs1", `"stale-etag"`))

	require.Equal(t, http.StatusOK, recorder.Code, "without a tracked version we must not serve a 304")
	assert.Empty(t, recorder.Header().Get("ETag"))

	var resp ofrep.BulkEvaluationResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	require.Len(t, resp.EventStreams, 1)
	assert.NotContains(t, resp.Metadata, "flagConfigLastModified")
}
