package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// jsonBody is comfortably over minSize so the middleware's size threshold is not what is under test.
var jsonBody = `{"flags":[` + strings.Repeat(`{"key":"flag","value":true},`, 200) + `{"key":"last","value":true}]}`

// jsonHandler writes body as JSON, plus any extra headers, at the given status.
func jsonHandler(status int, body string, headers map[string]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		for k, v := range headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	})
}

func serve(t *testing.T, handler http.Handler, acceptEncoding string) *httptest.ResponseRecorder {
	t.Helper()

	mw, err := New()
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/ofrep/v1/evaluate/flags", nil)
	if acceptEncoding != "" {
		req.Header.Set("Accept-Encoding", acceptEncoding)
	}

	recorder := httptest.NewRecorder()
	mw.Handler(handler).ServeHTTP(recorder, req)

	return recorder
}

func TestCompressesJSONWhenClientAcceptsGzip(t *testing.T) {
	recorder := serve(t, jsonHandler(http.StatusOK, jsonBody, nil), "gzip")

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "gzip", recorder.Header().Get("Content-Encoding"))
	assert.Contains(t, recorder.Header().Values("Vary"), "Accept-Encoding")
	assert.Less(t, recorder.Body.Len(), len(jsonBody), "compressed body should be smaller than the original")

	reader, err := gzip.NewReader(recorder.Body)
	require.NoError(t, err)
	decoded, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, jsonBody, string(decoded))
}

func TestLeavesBodyAloneWithoutAcceptEncoding(t *testing.T) {
	recorder := serve(t, jsonHandler(http.StatusOK, jsonBody, nil), "")

	require.Equal(t, http.StatusOK, recorder.Code)
	assert.Empty(t, recorder.Header().Get("Content-Encoding"))
	assert.Equal(t, jsonBody, recorder.Body.String())
}

// A single-flag evaluation is far below minSize, where gzip framing would make the response bigger.
func TestSkipsResponsesBelowMinSize(t *testing.T) {
	small := `{"key":"my-flag","value":true,"reason":"STATIC","variant":"on"}`
	require.Less(t, len(small), minSize)

	recorder := serve(t, jsonHandler(http.StatusOK, small, nil), "gzip")

	assert.Empty(t, recorder.Header().Get("Content-Encoding"))
	assert.Equal(t, small, recorder.Body.String())
}

// OFREP's bulk ETag digests the uncompressed body, so it must survive compression untouched for a
// client's If-None-Match to still match on the next request.
func TestPreservesETag(t *testing.T) {
	const etag = `"0123456789abcdef"`

	recorder := serve(t, jsonHandler(http.StatusOK, jsonBody, map[string]string{"ETag": etag}), "gzip")

	require.Equal(t, "gzip", recorder.Header().Get("Content-Encoding"))
	assert.Equal(t, etag, recorder.Header().Get("ETag"))
}

// A 304 carries no body, so there is nothing to encode and the client must not be told otherwise.
func TestLeavesNotModifiedAlone(t *testing.T) {
	recorder := serve(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotModified)
	}), "gzip")

	require.Equal(t, http.StatusNotModified, recorder.Code)
	assert.Empty(t, recorder.Header().Get("Content-Encoding"))
	assert.Empty(t, recorder.Body.Bytes())
}

// text/event-stream is not in the compressed content types; SSE is additionally routed around this
// middleware entirely, but the filter is the backstop if that ever changes.
func TestSkipsNonJSONContentTypes(t *testing.T) {
	recorder := serve(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte(strings.Repeat("data: {}\n\n", 500)))
	}), "gzip")

	assert.Empty(t, recorder.Header().Get("Content-Encoding"))
}
