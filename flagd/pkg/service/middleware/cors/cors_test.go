package cors

import (
	"net/http"
	"net/http/httptest"
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
