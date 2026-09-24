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
// actually writes: a key, a value, a reason, a variant and per-flag metadata.
func bulkPayload(n int) []byte {
	reasons := []string{"STATIC", "TARGETING_MATCH", "DEFAULT", "SPLIT"}
	variants := []string{"on", "off", "control", "treatment"}

	flags := make([]map[string]any, 0, n)
	for i := range n {
		flags = append(flags, map[string]any{
			"key":      fmt.Sprintf("my-service-feature-flag-%d", i),
			"value":    i%2 == 0,
			"reason":   reasons[i%len(reasons)],
			"variant":  variants[i%len(variants)],
			"metadata": map[string]any{"flagSetId": "my-flag-set", "version": fmt.Sprintf("v1.%d.0", i%20)},
		})
	}

	body, err := json.Marshal(map[string]any{"flags": flags})
	if err != nil {
		panic(err)
	}

	return body
}

// BenchmarkCompress is what DefaultMinSize is argued from: it serves payloads from a single-flag
// evaluation up to a large bulk response, with and without the middleware, and reports the
// resulting body size alongside the time. Compression costs a near-fixed few microseconds up to a
// few KB, which is why the smallest responses are not worth compressing.
func BenchmarkCompress(b *testing.B) {
	for _, flags := range []int{1, 5, 10, 100, 1000} {
		body := string(bulkPayload(flags))

		for _, acceptEncoding := range []string{"gzip", "identity"} {
			b.Run(fmt.Sprintf("flags=%d/bytes=%d/%s", flags, len(body), acceptEncoding), func(b *testing.B) {
				handler := jsonHandler(http.StatusOK, body, nil)
				if acceptEncoding == "gzip" {
					mw, err := New(Config{MinSize: 0}) // compress at every size, including below the default
					require.NoError(b, err)
					handler = mw.Handler(handler)
				}

				req := httptest.NewRequest(http.MethodPost, "/ofrep/v1/evaluate/flags", nil)
				req.Header.Set("Accept-Encoding", acceptEncoding)

				sized := httptest.NewRecorder()
				handler.ServeHTTP(sized, req)
				out := sized.Body.Len()

				b.ReportAllocs()
				b.ResetTimer()
				for range b.N {
					handler.ServeHTTP(httptest.NewRecorder(), req)
				}
				// reported after the loop: ResetTimer discards user metrics set before it
				b.ReportMetric(float64(out), "out-bytes")
			})
		}
	}
}
