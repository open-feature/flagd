package ofrep

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func serveWithETag(t *testing.T, status int, body, ifNoneMatch string) *httptest.ResponseRecorder {
	t.Helper()

	handler := conditionalETag(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))

	req := httptest.NewRequest(http.MethodPost, "/ofrep/v1/evaluate/flags", nil)
	if ifNoneMatch != "" {
		req.Header.Set("If-None-Match", ifNoneMatch)
	}

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	return recorder
}

func TestConditionalETag_HashesWrittenBytes(t *testing.T) {
	const body = `{"flags":[{"key":"a","value":true}],"metadata":{}}`

	first := serveWithETag(t, http.StatusOK, body, "")
	require.Equal(t, http.StatusOK, first.Code)
	assert.Equal(t, body, first.Body.String(), "the body must pass through untouched")

	etag := first.Header().Get("ETag")
	require.NotEmpty(t, etag)
	assert.Equal(t, quoteETag(bodyETag([]byte(body))), etag, "the tag must be the hash of the written bytes")

	t.Run("identical bytes yield the same tag", func(t *testing.T) {
		assert.Equal(t, etag, serveWithETag(t, http.StatusOK, body, "").Header().Get("ETag"))
	})

	t.Run("any byte difference yields a new tag", func(t *testing.T) {
		changed := serveWithETag(t, http.StatusOK, `{"flags":[{"key":"a","value":false}],"metadata":{}}`, "")
		assert.NotEqual(t, etag, changed.Header().Get("ETag"))
	})
}

func TestConditionalETag_NotModified(t *testing.T) {
	const body = `{"flags":[]}`
	etag := serveWithETag(t, http.StatusOK, body, "").Header().Get("ETag")
	require.NotEmpty(t, etag)

	recorder := serveWithETag(t, http.StatusOK, body, etag)

	assert.Equal(t, http.StatusNotModified, recorder.Code)
	assert.Equal(t, etag, recorder.Header().Get("ETag"), "a 304 must still name the representation")
	assert.Empty(t, recorder.Body.String(), "a 304 must not carry a body")
	assert.Empty(t, recorder.Header().Get("Content-Type"), "a 304 must not describe a body it does not have")
}

func TestConditionalETag_UnquotedValidatorMatches(t *testing.T) {
	const body = `{"flags":[]}`
	etag := serveWithETag(t, http.StatusOK, body, "").Header().Get("ETag")

	recorder := serveWithETag(t, http.StatusOK, body, normalizeETag(etag))
	assert.Equal(t, http.StatusNotModified, recorder.Code)
}

func TestConditionalETag_StaleValidatorIsServedFresh(t *testing.T) {
	recorder := serveWithETag(t, http.StatusOK, `{"flags":[]}`, `"stale"`)

	require.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, `{"flags":[]}`, recorder.Body.String())
	assert.NotEqual(t, `"stale"`, recorder.Header().Get("ETag"))
}

// A stale validator must not 304 an error body.
func TestConditionalETag_NonOKResponsesAreNotTagged(t *testing.T) {
	for _, status := range []int{
		http.StatusBadRequest,
		http.StatusInternalServerError,
		http.StatusRequestEntityTooLarge,
	} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			body := `{"errorDetails":"nope"}`
			recorder := serveWithETag(t, status, body, quoteETag(bodyETag([]byte(body))))

			assert.Equal(t, status, recorder.Code)
			assert.Empty(t, recorder.Header().Get("ETag"))
			assert.Equal(t, body, recorder.Body.String())
		})
	}
}

func TestConditionalETag_ImplicitOK(t *testing.T) {
	handler := conditionalETag(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"flags":[]}`))
	}))

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/ofrep/v1/evaluate/flags", nil))

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.NotEmpty(t, recorder.Header().Get("ETag"))
	assert.Equal(t, `{"flags":[]}`, recorder.Body.String())
}

// An empty 200 is still a representation, so it gets a validator.
func TestConditionalETag_EmptyBody(t *testing.T) {
	recorder := serveWithETag(t, http.StatusOK, "", "")

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, quoteETag(bodyETag(nil)), recorder.Header().Get("ETag"))
}
