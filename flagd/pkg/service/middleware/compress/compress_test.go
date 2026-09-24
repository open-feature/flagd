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

// bulkBody stands in for a bulk evaluation response: comfortably over DefaultMinSize, so the size
// threshold is never what is under test.
var bulkBody = `{"flags":[` + strings.Repeat(`{"key":"flag","value":true},`, 200) + `{"key":"last","value":true}]}`

// singleBody stands in for a single-flag evaluation, which falls below DefaultMinSize.
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

func serve(t *testing.T, cfg Config, handler http.Handler, acceptEncoding string) *httptest.ResponseRecorder {
	t.Helper()

	mw, err := New(cfg)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/ofrep/v1/evaluate/flags", nil)
	if acceptEncoding != "" {
		req.Header.Set("Accept-Encoding", acceptEncoding)
	}

	recorder := httptest.NewRecorder()
	mw.Handler(handler).ServeHTTP(recorder, req)

	return recorder
}

// The ETag assertion is the load-bearing one: OFREP's bulk tag digests the uncompressed body, so it
// must survive compression untouched for a client's If-None-Match to still match on the next request.
func TestCompressesJSONWhenClientAcceptsGzip(t *testing.T) {
	const etag = `"0123456789abcdef"`

	recorder := serve(t, Config{MinSize: DefaultMinSize},
		jsonHandler(http.StatusOK, bulkBody, map[string]string{"ETag": etag}), "gzip")

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "gzip", recorder.Header().Get("Content-Encoding"))
	assert.Contains(t, recorder.Header().Values("Vary"), "Accept-Encoding")
	assert.Equal(t, etag, recorder.Header().Get("ETag"))
	assert.Less(t, recorder.Body.Len(), len(bulkBody), "compressed body should be smaller than the original")

	reader, err := gzip.NewReader(recorder.Body)
	require.NoError(t, err)
	decoded, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, bulkBody, string(decoded))
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
		// gzip framing would cost more than it saves on a single-flag evaluation
		"body below the minimum size": {jsonHandler(http.StatusOK, singleBody, nil), "gzip"},
		// SSE is routed around this middleware entirely; the content-type filter is the backstop
		"content type is not JSON": {eventStream, "gzip"},
		"no body to encode":        {jsonHandler(http.StatusNotModified, "", nil), "gzip"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			recorder := serve(t, Config{MinSize: DefaultMinSize}, test.handler, test.acceptEncoding)
			assert.Empty(t, recorder.Header().Get("Content-Encoding"))
		})
	}
}

func TestMinSize(t *testing.T) {
	require.Less(t, len(singleBody), DefaultMinSize)

	tests := map[string]struct {
		minSize  int
		compress bool
	}{
		"zero compresses whatever the size": {0, true},
		"a body at exactly the minimum":     {len(singleBody), true},
		"a body one byte under":             {len(singleBody) + 1, false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			recorder := serve(t, Config{MinSize: test.minSize}, jsonHandler(http.StatusOK, singleBody, nil), "gzip")

			if !test.compress {
				assert.Empty(t, recorder.Header().Get("Content-Encoding"))
				return
			}

			require.Equal(t, "gzip", recorder.Header().Get("Content-Encoding"))
			reader, err := gzip.NewReader(recorder.Body)
			require.NoError(t, err)
			decoded, err := io.ReadAll(reader)
			require.NoError(t, err)
			assert.Equal(t, singleBody, string(decoded))
		})
	}

	_, err := New(Config{MinSize: -1})
	require.Error(t, err, "a negative minimum should be rejected rather than silently clamped")
}
