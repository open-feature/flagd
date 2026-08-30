package cors

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/open-feature/flagd/flagd/pkg/service/middleware/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockMw := middlewaremock.NewMockIMiddleware(ctrl)

	handlerFunc := http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
		},
	)

	mockMw.EXPECT().Handler(gomock.Any()).Return(handlerFunc)

	ts := httptest.NewServer(handlerFunc)

	defer ts.Close()

	mw := New([]string{"*"})
	require.NotNil(t, mw)

	// wrap the cors middleware around the mock to make sure the wrapped handler is called by the cors middleware
	ts.Config.Handler = mw.Handler(mockMw.Handler(handlerFunc))

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)

	require.Nil(t, err)

	client := http.DefaultClient
	resp, err := client.Do(req)

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

// Without this a browser cannot read the ETag cross-origin, so it never sends If-None-Match and
// the OFREP bulk endpoint's 304 path is unreachable.
func TestMiddleware_ExposesETag(t *testing.T) {
	handler := New([]string{"*"}).Handler(http.HandlerFunc(
		func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("ETag", `"abc123"`)
			writer.WriteHeader(http.StatusOK)
		},
	))

	req := httptest.NewRequest(http.MethodPost, "/ofrep/v1/evaluate/flags", nil)
	req.Header.Set("Origin", "https://app.example.com")

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	exposed := recorder.Header().Get("Access-Control-Expose-Headers")
	require.Contains(t, strings.ToLower(exposed), "etag")
}
