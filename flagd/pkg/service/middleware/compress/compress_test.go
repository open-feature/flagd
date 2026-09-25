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

// bulkBody stands in for a bulk evaluation response.
var bulkBody = `{"flags":[` + strings.Repeat(`{"key":"flag","value":true},`, 200) + `{"key":"last","value":true}]}`

// singleBody stands in for a single-flag evaluation, which is compressed too: there is no size
// threshold.
const singleBody = `{"key":"my-flag","value":true,"reason":"STATIC","variant":"on"}`

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

// serve runs handler behind the middleware and returns the recorded response.
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
	tests := map[string]string{"a bulk response": bulkBody, "a single evaluation": singleBody}

	for name, body := range tests {
		t.Run(name, func(t *testing.T) {
			recorder := serve(t, jsonHandler(http.StatusOK, body, nil), "gzip")

			require.Equal(t, http.StatusOK, recorder.Code)
			require.Equal(t, "gzip", recorder.Header().Get("Content-Encoding"))
			assert.Contains(t, recorder.Header().Values("Vary"), "Accept-Encoding")

			reader, err := gzip.NewReader(recorder.Body)
			require.NoError(t, err)
			decoded, err := io.ReadAll(reader)
			require.NoError(t, err)
			assert.Equal(t, body, string(decoded))
		})
	}
}

// The middleware must leave the ETag alone. A weak validator covers both encodings, so suffixing
// it would only make the two look like different representations to a cache.
func TestPreservesETag(t *testing.T) {
	const etag = `W/"0123456789abcdef"`

	recorder := serve(t, jsonHandler(http.StatusOK, bulkBody, map[string]string{"ETag": etag}), "gzip")

	assert.Equal(t, etag, recorder.Header().Get("ETag"))
}

func TestLeavesResponseUncompressed(t *testing.T) {
	eventStream := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte(strings.Repeat("data: {}\n\n", 500)))
	})

	tests := map[string]struct {
		handler        http.Handler
		acceptEncoding string
	}{
		"client does not accept gzip": {jsonHandler(http.StatusOK, bulkBody, nil), ""},
		"client refuses gzip by name": {jsonHandler(http.StatusOK, bulkBody, nil), "identity"},
		// SSE is routed around this middleware entirely; the content-type filter is the backstop
		"content type is not JSON": {eventStream, "gzip"},
		"no body to encode":        {jsonHandler(http.StatusNotModified, "", nil), "gzip"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			recorder := serve(t, test.handler, test.acceptEncoding)
			assert.Empty(t, recorder.Header().Get("Content-Encoding"))
		})
	}
}
