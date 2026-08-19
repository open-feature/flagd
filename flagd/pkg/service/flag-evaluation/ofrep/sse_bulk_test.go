package ofrep

import (
	"bytes"
	"context"
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

// serveBulk mirrors the production route, so the conditionalETag middleware is exercised.
func serveBulk(h handler, req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.Handle(bulkEvaluation, conditionalETag(http.HandlerFunc(h.HandleBulkEvaluation)))
	router.ServeHTTP(recorder, req)
	return recorder
}

// bulkRequestWithContext posts an evaluation context, so tests can vary it between requests.
func bulkRequestWithContext(t *testing.T, evalContext map[string]any, clientEtag string) *http.Request {
	t.Helper()
	body, err := json.Marshal(ofrep.Request{Context: evalContext})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/ofrep/v1/evaluate/flags", bytes.NewReader(body))
	require.NoError(t, err)
	req.Host = "flagd.example:8016"
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(svc.FLAGD_SELECTOR_HEADER, "flagSetId=fs1")
	if clientEtag != "" {
		req.Header.Set("If-None-Match", clientEtag)
	}
	return req
}

func boolFlag(key string, value bool) evaluator.AnyValue {
	variant := "off"
	if value {
		variant = "on"
	}
	return evaluator.AnyValue{Value: value, Variant: variant, Reason: model.StaticReason, FlagKey: key}
}

// The validator is derived from the response, so it is stable while the response is.
func TestHandleBulkEvaluation_ETagCoversResponseBody(t *testing.T) {
	log := logger.NewLogger(nil, false)

	serve := func(values []evaluator.AnyValue, clientEtag string) *httptest.ResponseRecorder {
		eval := mock.NewMockIEvaluator(gomock.NewController(t))
		eval.EXPECT().ResolveAllValues(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(values, model.Metadata{}, nil)

		h := handler{Logger: log, evaluator: eval}
		return serveBulk(h, newBulkRequest(t, "flagSetId=fs1", clientEtag))
	}

	first := serve([]evaluator.AnyValue{boolFlag("a", true)}, "")
	require.Equal(t, http.StatusOK, first.Code)
	etag := first.Header().Get("ETag")
	require.NotEmpty(t, etag)

	t.Run("identical response yields the same validator", func(t *testing.T) {
		again := serve([]evaluator.AnyValue{boolFlag("a", true)}, "")
		assert.Equal(t, etag, again.Header().Get("ETag"))
	})

	t.Run("matching validator yields an empty 304", func(t *testing.T) {
		conditional := serve([]evaluator.AnyValue{boolFlag("a", true)}, etag)
		assert.Equal(t, http.StatusNotModified, conditional.Code)
		assert.Equal(t, etag, conditional.Header().Get("ETag"))
		assert.Empty(t, conditional.Body.String(), "a 304 must not carry a body")
	})

	t.Run("a changed evaluated value invalidates the validator", func(t *testing.T) {
		changed := serve([]evaluator.AnyValue{boolFlag("a", false)}, etag)
		require.Equal(t, http.StatusOK, changed.Code, "a different resolved value must not be served as 304")
		assert.NotEqual(t, etag, changed.Header().Get("ETag"))
	})
}

// Regression test for the config-scoped validator: the reference ofrep-web provider only drops
// its cached ETag when targetingKey changes, so any other context change sent a stale
// If-None-Match and was answered 304 with evaluations for the previous context.
func TestHandleBulkEvaluation_ContextChangeInvalidatesETag(t *testing.T) {
	log := logger.NewLogger(nil, false)

	// same flag configuration throughout; only the context, and so the resolved value, differs
	serve := func(evalContext map[string]any, clientEtag string) *httptest.ResponseRecorder {
		eval := mock.NewMockIEvaluator(gomock.NewController(t))
		eval.EXPECT().ResolveAllValues(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, _ string, ec map[string]any) ([]evaluator.AnyValue, model.Metadata, error) {
				return []evaluator.AnyValue{boolFlag("discount-banner", ec["plan"] == "pro")}, model.Metadata{}, nil
			})

		h := handler{
			Logger:     log,
			evaluator:  eval,
			versioner:  fakeVersioner{etag: "unchanged-config", lastModified: 1700000000, ok: true},
			sseEnabled: true,
		}
		return serveBulk(h, bulkRequestWithContext(t, evalContext, clientEtag))
	}

	free := serve(map[string]any{"targetingKey": "user-1", "plan": "free"}, "")
	require.Equal(t, http.StatusOK, free.Code)
	freeEtag := free.Header().Get("ETag")
	require.NotEmpty(t, freeEtag)

	// targetingKey is unchanged, so the reference provider keeps sending the cached validator
	pro := serve(map[string]any{"targetingKey": "user-1", "plan": "pro"}, freeEtag)

	require.Equal(t, http.StatusOK, pro.Code,
		"a context change must be served fresh even when the flag configuration is untouched")
	assert.NotEqual(t, freeEtag, pro.Header().Get("ETag"))

	var resp ofrep.BulkEvaluationResponse
	require.NoError(t, json.Unmarshal(pro.Body.Bytes(), &resp))
	require.Len(t, resp.Flags, 1)
	flag, ok := resp.Flags[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, true, flag["value"], "the response must carry values for the new context")
}

// The 304 decision must use If-None-Match (the client's cached validator), not the ADR-0008
// flagConfigEtag change-trigger metadata.
func TestHandleBulkEvaluation_ConditionalUsesIfNoneMatch(t *testing.T) {
	log := logger.NewLogger(nil, false)
	values := []evaluator.AnyValue{boolFlag("a", true)}

	serve := func(t *testing.T, query, ifNoneMatch string) *httptest.ResponseRecorder {
		t.Helper()
		eval := mock.NewMockIEvaluator(gomock.NewController(t))
		eval.EXPECT().ResolveAllValues(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(values, model.Metadata{}, nil)

		req, err := http.NewRequest(http.MethodPost, "/ofrep/v1/evaluate/flags"+query, bytes.NewReader([]byte{}))
		require.NoError(t, err)
		req.Header.Set(svc.FLAGD_SELECTOR_HEADER, "flagSetId=fs1")
		if ifNoneMatch != "" {
			req.Header.Set("If-None-Match", ifNoneMatch)
		}

		h := handler{Logger: log, evaluator: eval, versioner: fakeVersioner{etag: "cfg", ok: true}, sseEnabled: true}
		return serveBulk(h, req)
	}

	current := serve(t, "", "").Header().Get("ETag")
	require.NotEmpty(t, current)

	t.Run("current validator yields 304", func(t *testing.T) {
		assert.Equal(t, http.StatusNotModified, serve(t, "", current).Code)
	})

	t.Run("stale validator is served fresh", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, serve(t, "", `"stale"`).Code)
	})

	t.Run("trigger etag alone is not a validator", func(t *testing.T) {
		// flagConfigEtag is change-trigger metadata; on its own it must never produce a 304
		assert.Equal(t, http.StatusOK, serve(t, "?flagConfigEtag=cfg", "").Code)
	})

	t.Run("trigger etag does not override a stale validator", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, serve(t, "?flagConfigEtag=cfg", `"stale"`).Code)
	})
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
			assert.NotEmpty(t, recorder.Header().Get("ETag"))

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

	h := handler{Logger: log, evaluator: eval}

	recorder := serveBulk(h, newBulkRequest(t, "", ""))

	require.Equal(t, http.StatusOK, recorder.Code)
	// the validator describes the response, so it is offered regardless of SSE
	assert.NotEmpty(t, recorder.Header().Get("ETag"))

	var resp ofrep.BulkEvaluationResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	assert.Empty(t, resp.EventStreams)
}

// The validator comes from the response and is always available; flagConfigLastModified
// describes the configuration and needs a live subscription.
func TestHandleBulkEvaluation_NoLiveStream_StillHasETag(t *testing.T) {
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

	require.Equal(t, http.StatusOK, recorder.Code, "a stale validator must be served fresh")
	assert.NotEmpty(t, recorder.Header().Get("ETag"))

	var resp ofrep.BulkEvaluationResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	require.Len(t, resp.EventStreams, 1)
	assert.NotContains(t, resp.Metadata, "flagConfigLastModified")
}
