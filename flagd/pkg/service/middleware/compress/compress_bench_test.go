package compress

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// bulkPayload builds an OFREP bulk evaluation response with n flags, shaped like what flagd
// actually writes: repeated keys, a small set of repeated reasons and variants, and per-flag
// metadata. The repetition is the reason these payloads compress so well.
func bulkPayload(n int) []byte {
	reasons := []string{"STATIC", "TARGETING_MATCH", "DEFAULT", "SPLIT"}
	variants := []string{"on", "off", "control", "treatment"}

	flags := make([]map[string]any, 0, n)
	for i := 0; i < n; i++ {
		flags = append(flags, map[string]any{
			"key":     fmt.Sprintf("my-service-feature-flag-%d", i),
			"value":   i%2 == 0,
			"reason":  reasons[i%len(reasons)],
			"variant": variants[i%len(variants)],
			"metadata": map[string]any{
				"flagSetId": "my-flag-set",
				"version":   fmt.Sprintf("v1.%d.0", i%20),
			},
		})
	}

	body, err := json.Marshal(map[string]any{"flags": flags})
	if err != nil {
		panic(err)
	}

	return body
}

// benchSizes span a single-flag evaluation up to a large bulk response.
var benchSizes = []int{1, 2, 5, 10, 25, 50, 100, 250, 500, 1000}

// BenchmarkCompress measures the end-to-end cost of serving a bulk payload through the middleware
// with the client accepting gzip, against the same payload served with no middleware at all. The
// reported compressed-size ratio comes from a one-shot run before the timed loop.
func BenchmarkCompress(b *testing.B) {
	for _, flags := range benchSizes {
		body := bulkPayload(flags)

		b.Run(fmt.Sprintf("flags=%d/bytes=%d/gzip", flags, len(body)), func(b *testing.B) {
			mw, err := New(Config{MinSize: 0}) // force compression at every size
			require.NoError(b, err)
			handler := mw.Handler(jsonHandler(http.StatusOK, string(body), nil))

			benchServe(b, handler, "gzip")
		})

		b.Run(fmt.Sprintf("flags=%d/bytes=%d/identity", flags, len(body)), func(b *testing.B) {
			handler := jsonHandler(http.StatusOK, string(body), nil)

			benchServe(b, handler, "")
		})
	}
}

// servedSize returns the body size a single request through handler produces.
func servedSize(b *testing.B, handler http.Handler, acceptEncoding string) int {
	b.Helper()

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, benchRequest(acceptEncoding))

	return recorder.Body.Len()
}

func benchRequest(acceptEncoding string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/ofrep/v1/evaluate/flags", nil)
	if acceptEncoding != "" {
		req.Header.Set("Accept-Encoding", acceptEncoding)
	}

	return req
}

func benchServe(b *testing.B, handler http.Handler, acceptEncoding string) {
	b.Helper()

	req := benchRequest(acceptEncoding)
	out := servedSize(b, handler, acceptEncoding)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)
	}
	// reported after the loop: ResetTimer discards user metrics set before it
	b.ReportMetric(float64(out), "out-bytes")
}
