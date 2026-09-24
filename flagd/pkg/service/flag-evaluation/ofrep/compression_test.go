package ofrep

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/open-feature/flagd/core/pkg/evaluator"
	mock "github.com/open-feature/flagd/core/pkg/evaluator/mock"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	compressmw "github.com/open-feature/flagd/flagd/pkg/service/middleware/compress"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/mock/gomock"
)

// sizeRecordingMetricsRecorder captures what the metrics middleware reports as the response size,
// which is the figure operators read as bytes served.
type sizeRecordingMetricsRecorder struct {
	telemetry.NoopMetricsRecorder
	sizes []int64
}

func (r *sizeRecordingMetricsRecorder) HTTPResponseSize(_ context.Context, sizeBytes int64, _ []attribute.KeyValue) {
	r.sizes = append(r.sizes, sizeBytes)
}

// flagCount is enough flags for a bulk response to clear DefaultMinSize.
const flagCount = 40

// compressibleEvaluations builds a bulk evaluation result of n flags.
func compressibleEvaluations(n int) []evaluator.AnyValue {
	values := make([]evaluator.AnyValue, 0, n)
	for i := range n {
		values = append(values, evaluator.AnyValue{
			Value:   true,
			Variant: "on",
			Reason:  model.StaticReason,
			FlagKey: fmt.Sprintf("my-service-feature-flag-%d", i),
		})
	}

	return values
}

// compressingHandler builds the production route chain with compression enabled.
func compressingHandler(t *testing.T, metrics telemetry.IMetricsRecorder, evaluations []evaluator.AnyValue) http.Handler {
	t.Helper()

	eval := mock.NewMockIEvaluator(gomock.NewController(t))
	eval.EXPECT().ResolveAllValues(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(evaluations, model.Metadata{}, nil).AnyTimes()

	compression, err := compressmw.New(compressmw.Config{MinSize: compressmw.DefaultMinSize})
	require.NoError(t, err)

	return NewOfrepHandler(logger.NewLogger(nil, false), eval, nil, nil, metrics, "flagd", SSEConfig{}, compression)
}

func serveCompressedBulk(handler http.Handler, acceptEncoding, ifNoneMatch string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/ofrep/v1/evaluate/flags", nil)
	req.Header.Set("Accept-Encoding", acceptEncoding)
	if ifNoneMatch != "" {
		req.Header.Set("If-None-Match", ifNoneMatch)
	}

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	return recorder
}

func TestBulkEvaluationIsCompressed(t *testing.T) {
	handler := compressingHandler(t, &telemetry.NoopMetricsRecorder{}, compressibleEvaluations(flagCount))

	recorder := serveCompressedBulk(handler, "gzip", "")

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "gzip", recorder.Header().Get("Content-Encoding"))

	reader, err := gzip.NewReader(recorder.Body)
	require.NoError(t, err)
	decoded, err := io.ReadAll(reader)
	require.NoError(t, err)

	var body struct {
		Flags []map[string]any `json:"flags"`
	}
	require.NoError(t, json.Unmarshal(decoded, &body))
	assert.Len(t, body.Flags, flagCount)
}

// The compressed body is a different representation, so it must not reuse the uncompressed
// response's strong validator (RFC 9110 8.8.1), and flagd has to strip the marker back off to
// keep answering conditional requests.
func TestCompressedBulkEvaluationETags(t *testing.T) {
	handler := compressingHandler(t, &telemetry.NoopMetricsRecorder{}, compressibleEvaluations(flagCount))

	gzipETag := serveCompressedBulk(handler, "gzip", "").Header().Get("ETag")
	identityETag := serveCompressedBulk(handler, "identity", "").Header().Get("ETag")

	require.NotEmpty(t, identityETag)
	require.NotEqual(t, identityETag, gzipETag, "the two encodings must not share one strong validator")
	require.Equal(t, identityETag, compressmw.TrimETagSuffix(gzipETag), "the marker should hide a recoverable digest")

	tests := map[string]struct {
		acceptEncoding string
		ifNoneMatch    string
		expected       int
	}{
		"gzip tag, still accepting gzip":     {"gzip", gzipETag, http.StatusNotModified},
		"identity tag, still identity":       {"identity", identityETag, http.StatusNotModified},
		"identity tag, now accepting gzip":   {"gzip", identityETag, http.StatusNotModified},
		"gzip tag, no longer accepting gzip": {"identity", gzipETag, http.StatusOK},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			recorder := serveCompressedBulk(handler, test.acceptEncoding, test.ifNoneMatch)
			assert.Equal(t, test.expected, recorder.Code)
		})
	}
}

// The metrics middleware sits outside compression, so the size it reports is what actually goes on
// the wire rather than the body the evaluator produced.
func TestRecordedResponseSizeIsTheCompressedSize(t *testing.T) {
	metrics := &sizeRecordingMetricsRecorder{}
	handler := compressingHandler(t, metrics, compressibleEvaluations(flagCount))

	recorder := serveCompressedBulk(handler, "gzip", "")

	require.Len(t, metrics.sizes, 1)
	assert.Equal(t, int64(recorder.Body.Len()), metrics.sizes[0])
}

// The SSE stream must stay uncompressed so each event reaches subscribers as it is flushed, which
// means compression has to be wired onto the evaluate routes rather than around the whole mux.
func TestSSEStreamIsNotCompressed(t *testing.T) {
	eval := mock.NewMockIEvaluator(gomock.NewController(t))
	flagStore, err := store.NewStore(logger.NewLogger(nil, false), []string{"src1"})
	require.NoError(t, err)

	service, err := NewOfrepService(eval, flagStore, []string{"*"}, SvcConfiguration{
		Logger:             logger.NewLogger(nil, false),
		Port:               18286,
		ServiceName:        testServiceName,
		MetricsRecorder:    &telemetry.NoopMetricsRecorder{},
		SSEEnabled:         true,
		CompressionEnabled: true,
		CompressionMinSize: 0, // compress everything, so an exemption is the only way to stay plain
	}, nil, nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req := httptest.NewRequest(http.MethodGet, ssePath, nil).WithContext(ctx)
	req.Header.Set("Accept-Encoding", "gzip")

	recorder := httptest.NewRecorder()
	service.server.Handler.ServeHTTP(recorder, req)

	assert.Empty(t, recorder.Header().Get("Content-Encoding"))
}
