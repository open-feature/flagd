package h2c

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	middlewaremock "github.com/open-feature/flagd/core/pkg/service/middleware/mock"
	"github.com/stretchr/testify/require"
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

	mw := New()
	require.NotNil(t, mw)

	// wrap the h2c middleware around the mock to make sure the wrapped handler is called by the h2c middleware
	ts.Config.Handler = mw.Handler(mockMw.Handler(handlerFunc))

	resp, err := http.Get(ts.URL)

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
