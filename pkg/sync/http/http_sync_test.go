package http

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/open-feature/flagd/pkg/sync"

	"github.com/golang/mock/gomock"
	"github.com/open-feature/flagd/pkg/logger"
	syncmock "github.com/open-feature/flagd/pkg/sync/http/mock"
)

func TestSimpleSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	resp := "test response"

	mockCron := syncmock.NewMockCron(ctrl)
	mockCron.EXPECT().AddFunc(gomock.Any(), gomock.Any()).DoAndReturn(func(spec string, cmd func()) error {
		return nil
	})
	mockCron.EXPECT().Start().Times(1)

	mockClient := syncmock.NewMockHTTPClient(ctrl)
	mockClient.EXPECT().Do(gomock.Any()).Return(&http.Response{Body: io.NopCloser(strings.NewReader(resp))}, nil)

	httpSync := Sync{
		URI:         "http://localhost",
		Client:      mockClient,
		Cron:        mockCron,
		BearerToken: "",
		LastBodySHA: "",
		Logger:      logger.NewLogger(nil, false),
	}

	ctx := context.Background()
	dataSyncChan := make(chan sync.DataSync)

	go func() {
		err := httpSync.Sync(ctx, dataSyncChan)
		if err != nil {
			log.Fatalf("Error start sync: %s", err.Error())
			return
		}
	}()

	data := <-dataSyncChan

	if data.FlagData != resp {
		t.Errorf("Expected content %s, but received content %s", resp, data.FlagData)
	}
}

func TestHTTPSync_Fetch(t *testing.T) {
	ctrl := gomock.NewController(t)

	tests := map[string]struct {
		setup          func(t *testing.T, client *syncmock.MockHTTPClient)
		uri            string
		bearerToken    string
		lastBodySHA    string
		handleResponse func(*testing.T, Sync, string, error)
	}{
		"success": {
			setup: func(t *testing.T, client *syncmock.MockHTTPClient) {
				client.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Body: io.NopCloser(strings.NewReader("test response")),
				}, nil)
			},
			uri: "http://localhost",
			handleResponse: func(t *testing.T, _ Sync, fetched string, err error) {
				if err != nil {
					t.Fatalf("fetch: %v", err)
				}
				expected := "test response"
				if fetched != expected {
					t.Errorf("expected fetched to be '%s', got '%s'", expected, fetched)
				}
			},
		},
		"return an error if no uri": {
			setup: func(t *testing.T, client *syncmock.MockHTTPClient) {},
			handleResponse: func(t *testing.T, _ Sync, fetched string, err error) {
				if err == nil {
					t.Error("expected err, got nil")
				}
			},
		},
		"update last body sha": {
			setup: func(t *testing.T, client *syncmock.MockHTTPClient) {
				client.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Body: io.NopCloser(strings.NewReader("test response")),
				}, nil)
			},
			uri:         "http://localhost",
			lastBodySHA: "",
			handleResponse: func(t *testing.T, httpSync Sync, _ string, err error) {
				if err != nil {
					t.Fatalf("fetch: %v", err)
				}

				expectedLastBodySHA := "fUH6MbDL8tR0nCiC4bag0Rf_6is="
				if httpSync.LastBodySHA != expectedLastBodySHA {
					t.Errorf(
						"expected last body sha to be '%s', got '%s'", expectedLastBodySHA, httpSync.LastBodySHA,
					)
				}
			},
		},
		"authorization header": {
			setup: func(t *testing.T, client *syncmock.MockHTTPClient) {
				client.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Body: io.NopCloser(strings.NewReader("test response")),
				}, nil)
			},
			uri:         "http://localhost",
			lastBodySHA: "",
			handleResponse: func(t *testing.T, httpSync Sync, _ string, err error) {
				if err != nil {
					t.Fatalf("fetch: %v", err)
				}

				expectedLastBodySHA := "fUH6MbDL8tR0nCiC4bag0Rf_6is="
				if httpSync.LastBodySHA != expectedLastBodySHA {
					t.Errorf(
						"expected last body sha to be '%s', got '%s'", expectedLastBodySHA, httpSync.LastBodySHA,
					)
				}
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := syncmock.NewMockHTTPClient(ctrl)

			tt.setup(t, mockClient)

			httpSync := Sync{
				URI:         tt.uri,
				Client:      mockClient,
				BearerToken: tt.bearerToken,
				LastBodySHA: tt.lastBodySHA,
				Logger:      logger.NewLogger(nil, false),
			}

			fetched, err := httpSync.Fetch(context.Background())
			tt.handleResponse(t, httpSync, fetched, err)
		})
	}
}
