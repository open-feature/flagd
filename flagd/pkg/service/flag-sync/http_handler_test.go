package sync

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	syncv1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/sync/v1"
	"connectrpc.com/connect"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestHTTPHandler seeds a last-modified stamp so conditional requests are exercised.
func newTestHTTPHandler(t *testing.T) (httpHandler, store.IStore) {
	t.Helper()

	flagStore, _ := getSimpleFlagStore(t)
	mt := &modTime{}
	mt.set(time.Now())

	return httpHandler{
		store:   flagStore,
		log:     logger.NewLogger(nil, false),
		modTime: mt,
	}, flagStore
}

// serve routes through a mux so PathValue and the method restriction behave as they do in production.
func serve(h httpHandler, req *http.Request) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	mux.Handle("GET "+flagsPath, h)
	mux.Handle("GET "+flagsPath+"/{"+selectorPathVar+"}", h)

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

// The core guarantee: the HTTP body is byte-for-byte what gRPC FetchAllFlags returns.
func TestHTTPHandler_MatchesFetchAllFlags(t *testing.T) {
	h, flagStore := newTestHTTPHandler(t)

	grpcHandler := syncHandler{
		store:           flagStore,
		log:             logger.NewLogger(nil, false),
		metricsRecorder: &telemetry.NoopMetricsRecorder{},
	}

	for _, selector := range []string{"", "source=" + testSource1, "source=" + testSource2} {
		t.Run("selector="+selector, func(t *testing.T) {
			grpcResp, err := grpcHandler.FetchAllFlags(context.Background(),
				connect.NewRequest(&syncv1.FetchAllFlagsRequest{Selector: selector}))
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, flagsPath+"?selector="+url.QueryEscape(selector), nil)
			rec := serve(h, req)

			require.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, grpcResp.Msg.GetFlagConfiguration(), rec.Body.String())
		})
	}
}

func TestHTTPHandler_SelectorSources(t *testing.T) {
	h, _ := newTestHTTPHandler(t)

	tests := []struct {
		name    string
		request func() *http.Request
	}{
		{
			name: "header",
			request: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, flagsPath, nil)
				req.Header.Set(selectorHeaderKey, "source="+testSource1)
				return req
			},
		},
		{
			name: "path segment",
			request: func() *http.Request {
				return httptest.NewRequest(http.MethodGet,
					flagsPath+"/"+url.PathEscape("source="+testSource1), nil)
			},
		},
		{
			name: "query param",
			request: func() *http.Request {
				return httptest.NewRequest(http.MethodGet,
					flagsPath+"?selector="+url.QueryEscape("source="+testSource1), nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serve(h, tt.request())

			require.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "flagA")
			assert.NotContains(t, rec.Body.String(), "flagB")
		})
	}
}

// Pins header > path > query, the order the gRPC and OFREP services already document.
func TestHTTPHandler_SelectorPrecedence(t *testing.T) {
	h, _ := newTestHTTPHandler(t)

	t.Run("header beats path and query", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet,
			flagsPath+"/"+url.PathEscape("source="+testSource2)+"?selector="+url.QueryEscape("source="+testSource2), nil)
		req.Header.Set(selectorHeaderKey, "source="+testSource1)

		rec := serve(h, req)

		require.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "flagA")
		assert.NotContains(t, rec.Body.String(), "flagB")
	})

	t.Run("path beats query", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet,
			flagsPath+"/"+url.PathEscape("source="+testSource1)+"?selector="+url.QueryEscape("source="+testSource2), nil)

		rec := serve(h, req)

		require.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "flagA")
		assert.NotContains(t, rec.Body.String(), "flagB")
	})
}

// Malformed is the client's mistake, an unknown filter is a miss, and a valid filter matching
// nothing is an ordinary empty result.
func TestHTTPHandler_SelectorStatuses(t *testing.T) {
	h, _ := newTestHTTPHandler(t)

	tests := []struct {
		name       string
		selector   string
		wantStatus int
		wantBody   string
	}{
		{name: "control character", selector: "bad\x01char", wantStatus: http.StatusBadRequest},
		{name: "invalid utf-8", selector: "bad\xffutf8", wantStatus: http.StatusBadRequest},
		{name: "unknown filter key", selector: "bogus=1", wantStatus: http.StatusNotFound},
		{
			name:       "valid filter matching nothing",
			selector:   "flagSetId=does-not-exist",
			wantStatus: http.StatusOK,
			wantBody:   `{"flags":{}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, flagsPath, nil)
			req.Header.Set(selectorHeaderKey, tt.selector)

			rec := serve(h, req)

			require.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != "" {
				assert.Equal(t, tt.wantBody, rec.Body.String())
			}
			if tt.wantStatus >= 400 {
				// unescaped client input must not be reflected back
				assert.NotContains(t, rec.Body.String(), tt.selector)
			}
		})
	}
}

// Source selectors routinely contain "/", which the path-segment form has to survive.
func TestHTTPHandler_SelectorContainingSlash(t *testing.T) {
	sources := []string{"./flags.json"}
	flagStore, err := store.NewStore(logger.NewLogger(nil, false), sources)
	require.NoError(t, err)
	flagStore.Update(sources[0], testSource1Flags, model.Metadata{}, false)

	mt := &modTime{}
	mt.set(time.Now())
	h := httpHandler{store: flagStore, log: logger.NewLogger(nil, false), modTime: mt}

	req := httptest.NewRequest(http.MethodGet,
		flagsPath+"/"+url.PathEscape("source=./flags.json"), nil)

	rec := serve(h, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "flagA")
}

func TestHTTPHandler_MethodNotAllowed(t *testing.T) {
	h, _ := newTestHTTPHandler(t)

	rec := serve(h, httptest.NewRequest(http.MethodPost, flagsPath, nil))

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestHTTPHandler_ConditionalRequests(t *testing.T) {
	h, flagStore := newTestHTTPHandler(t)

	first := serve(h, httptest.NewRequest(http.MethodGet, flagsPath, nil))
	require.Equal(t, http.StatusOK, first.Code)

	etag := first.Header().Get("ETag")
	lastModified := first.Header().Get("Last-Modified")
	require.NotEmpty(t, etag)
	require.NotEmpty(t, lastModified)

	t.Run("matching ETag yields a bodyless 304", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, flagsPath, nil)
		req.Header.Set("If-None-Match", etag)

		rec := serve(h, req)

		assert.Equal(t, http.StatusNotModified, rec.Code)
		assert.Empty(t, rec.Body.String())
	})

	t.Run("stale ETag yields 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, flagsPath, nil)
		req.Header.Set("If-None-Match", `"stale"`)

		rec := serve(h, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("matching If-Modified-Since yields 304", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, flagsPath, nil)
		req.Header.Set("If-Modified-Since", lastModified)

		rec := serve(h, req)

		assert.Equal(t, http.StatusNotModified, rec.Code)
	})

	t.Run("older If-Modified-Since yields 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, flagsPath, nil)
		req.Header.Set("If-Modified-Since", h.modTime.get().Add(-time.Hour).UTC().Format(http.TimeFormat))

		rec := serve(h, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	// RFC 9110 gives If-None-Match outright precedence; pinned so a hand-rolled replacement for
	// ServeContent cannot silently regress it.
	t.Run("ETag wins over a stale date", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, flagsPath, nil)
		req.Header.Set("If-None-Match", etag)
		req.Header.Set("If-Modified-Since", time.Unix(0, 0).UTC().Format(http.TimeFormat))

		rec := serve(h, req)

		assert.Equal(t, http.StatusNotModified, rec.Code)
	})

	t.Run("stale ETag wins over a matching date", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, flagsPath, nil)
		req.Header.Set("If-None-Match", `"stale"`)
		req.Header.Set("If-Modified-Since", lastModified)

		rec := serve(h, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("ETag changes when the configuration changes", func(t *testing.T) {
		flagStore.Update(testSource1, []model.Flag{{
			Key:            "flagA",
			State:          "ENABLED",
			DefaultVariant: "true",
			Variants:       testVariants,
		}}, model.Metadata{}, false)

		rec := serve(h, httptest.NewRequest(http.MethodGet, flagsPath, nil))

		require.Equal(t, http.StatusOK, rec.Code)
		assert.NotEqual(t, etag, rec.Header().Get("ETag"))
	})
}

// We must not advertise a modification time we never observed.
func TestHTTPHandler_NoLastModifiedBeforeFirstSnapshot(t *testing.T) {
	flagStore, _ := getSimpleFlagStore(t)
	h := httpHandler{store: flagStore, log: logger.NewLogger(nil, false), modTime: &modTime{}}

	rec := serve(h, httptest.NewRequest(http.MethodGet, flagsPath, nil))

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, rec.Header().Get("Last-Modified"))
	assert.NotEmpty(t, rec.Header().Get("ETag"))
}

func TestWellFormedSelector(t *testing.T) {
	tests := []struct {
		expression string
		want       bool
	}{
		{expression: "", want: true},
		{expression: "source=./flags.json", want: true},
		{expression: "flagSetId=payments", want: true},
		{expression: "source=http://host/config", want: true},
		{expression: "bad\x01char", want: false},
		{expression: "bad\nnewline", want: false},
		{expression: "bad\xffutf8", want: false},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, wellFormedSelector(tt.expression), "expression %q", tt.expression)
	}
}

// The bare-source shorthand and its explicit form must be indistinguishable on the wire, ETag
// included, so either spelling revalidates against the same cache entry.
func TestHTTPHandler_SelectorEquivalence(t *testing.T) {
	h, _ := newTestHTTPHandler(t)

	get := func(expression string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, flagsPath, nil)
		req.Header.Set(selectorHeaderKey, expression)
		return serve(h, req)
	}

	bare := get(testSource1)
	explicit := get("source=" + testSource1)

	require.Equal(t, http.StatusOK, bare.Code)
	require.Equal(t, http.StatusOK, explicit.Code)
	assert.Equal(t, explicit.Body.String(), bare.Body.String())
	assert.Equal(t, explicit.Header().Get("ETag"), bare.Header().Get("ETag"))

	// an ETag from one form must satisfy a conditional request made with the other
	req := httptest.NewRequest(http.MethodGet, flagsPath, nil)
	req.Header.Set(selectorHeaderKey, "source="+testSource1)
	req.Header.Set("If-None-Match", bare.Header().Get("ETag"))
	assert.Equal(t, http.StatusNotModified, serve(h, req).Code)
}
