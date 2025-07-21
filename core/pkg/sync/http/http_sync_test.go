package http

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	syncmock "github.com/open-feature/flagd/core/pkg/sync/http/mock"
	synctesting "github.com/open-feature/flagd/core/pkg/sync/testing"
	"go.uber.org/mock/gomock"
)

func TestSimpleSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCron := synctesting.NewMockCron(ctrl)
	mockCron.EXPECT().AddFunc(gomock.Any(), gomock.Any()).DoAndReturn(func(_ string, _ func()) error {
		return nil
	})
	mockCron.EXPECT().Start().Times(1)

	mockClient := syncmock.NewMockClient(ctrl)
	responseBody := "test response"
	resp := &http.Response{
		Header:     map[string][]string{"Content-Type": {"application/json"}},
		Body:       io.NopCloser(strings.NewReader(responseBody)),
		StatusCode: http.StatusOK,
	}
	mockClient.EXPECT().Do(gomock.Any()).Return(resp, nil)

	httpSync := Sync{
		URI:         "http://localhost/flags",
		Client:      mockClient,
		Cron:        mockCron,
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

	if data.FlagData != responseBody {
		t.Errorf("expected content: %s, but received content: %s", responseBody, data.FlagData)
	}
}

func TestExtensionWithQSSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCron := synctesting.NewMockCron(ctrl)
	mockCron.EXPECT().AddFunc(gomock.Any(), gomock.Any()).DoAndReturn(func(_ string, _ func()) error {
		return nil
	})
	mockCron.EXPECT().Start().Times(1)

	mockClient := syncmock.NewMockClient(ctrl)
	responseBody := "test response"
	resp := &http.Response{
		Header:     map[string][]string{"Content-Type": {"application/json"}},
		Body:       io.NopCloser(strings.NewReader(responseBody)),
		StatusCode: http.StatusOK,
	}
	mockClient.EXPECT().Do(gomock.Any()).Return(resp, nil)

	httpSync := Sync{
		URI:         "http://localhost/flags.json?env=dev",
		Client:      mockClient,
		Cron:        mockCron,
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

	if data.FlagData != responseBody {
		t.Errorf("expected content: %s, but received content: %s", responseBody, data.FlagData)
	}
}

func TestHTTPSync_Fetch(t *testing.T) {
	ctrl := gomock.NewController(t)

	tests := map[string]struct {
		setup          func(t *testing.T, client *syncmock.MockClient)
		uri            string
		bearerToken    string
		authHeader     string
		lastBodySHA    string
		handleResponse func(*testing.T, Sync, string, error)
	}{
		"success": {
			setup: func(_ *testing.T, client *syncmock.MockClient) {
				client.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Header:     map[string][]string{"Content-Type": {"application/json"}},
					Body:       io.NopCloser(strings.NewReader("test response")),
					StatusCode: http.StatusOK,
				}, nil)
			},
			uri: "http://localhost",
			handleResponse: func(t *testing.T, _ Sync, fetched string, err error) {
				if err != nil {
					t.Fatalf("fetch: %v", err)
				}
				expected := "test response"
				if fetched != expected {
					t.Errorf("expected fetched to be: '%s', got: '%s'", expected, fetched)
				}
			},
		},
		"return an error if no uri": {
			setup: func(_ *testing.T, _ *syncmock.MockClient) {},
			handleResponse: func(t *testing.T, _ Sync, _ string, err error) {
				if err == nil {
					t.Error("expected err, got nil")
				}
			},
		},
		"update last body sha": {
			setup: func(_ *testing.T, client *syncmock.MockClient) {
				client.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Header:     map[string][]string{"Content-Type": {"application/json"}},
					Body:       io.NopCloser(strings.NewReader("test response")),
					StatusCode: http.StatusOK,
				}, nil)
			},
			uri:         "http://localhost",
			lastBodySHA: "",
			handleResponse: func(t *testing.T, httpSync Sync, _ string, err error) {
				if err != nil {
					t.Fatalf("fetch: %v", err)
				}

				expectedLastBodySHA := "UjeJHtCU_wb7OHK-tbPoHycw0TqlHzkWJmH4y6cqg50="
				if httpSync.LastBodySHA != expectedLastBodySHA {
					t.Errorf(
						"expected last body sha to be: '%s', got: '%s'", expectedLastBodySHA, httpSync.LastBodySHA,
					)
				}
			},
		},
		"authorization with bearerToken": {
			setup: func(t *testing.T, client *syncmock.MockClient) {
				expectedToken := "bearer-1234"
				client.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
					actualAuthHeader := req.Header.Get("Authorization")
					if actualAuthHeader != "Bearer "+expectedToken {
						t.Fatalf("expected Authorization header to be 'Bearer %s', got %s", expectedToken, actualAuthHeader)
					}
					return &http.Response{
						Header:     map[string][]string{"Content-Type": {"application/json"}},
						Body:       io.NopCloser(strings.NewReader("test response")),
						StatusCode: http.StatusOK,
					}, nil
				})
			},
			uri:         "http://localhost",
			bearerToken: "bearer-1234",
			lastBodySHA: "",
			handleResponse: func(t *testing.T, httpSync Sync, _ string, err error) {
				if err != nil {
					t.Fatalf("fetch: %v", err)
				}

				expectedLastBodySHA := "UjeJHtCU_wb7OHK-tbPoHycw0TqlHzkWJmH4y6cqg50="
				if httpSync.LastBodySHA != expectedLastBodySHA {
					t.Errorf(
						"expected last body sha to be: '%s', got: '%s'", expectedLastBodySHA, httpSync.LastBodySHA,
					)
				}
			},
		},
		"authorization with authHeader": {
			setup: func(t *testing.T, client *syncmock.MockClient) {
				expectedHeader := "Basic dXNlcjpwYXNz"
				client.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
					actualAuthHeader := req.Header.Get("Authorization")
					if actualAuthHeader != expectedHeader {
						t.Fatalf("expected Authorization header to be '%s', got %s", expectedHeader, actualAuthHeader)
					}
					return &http.Response{
						Header:     map[string][]string{"Content-Type": {"application/json"}},
						Body:       io.NopCloser(strings.NewReader("test response")),
						StatusCode: http.StatusOK,
					}, nil
				})
			},
			uri:         "http://localhost",
			authHeader:  "Basic dXNlcjpwYXNz",
			lastBodySHA: "",
			handleResponse: func(t *testing.T, httpSync Sync, _ string, err error) {
				if err != nil {
					t.Fatalf("fetch: %v", err)
				}

				expectedLastBodySHA := "UjeJHtCU_wb7OHK-tbPoHycw0TqlHzkWJmH4y6cqg50="
				if httpSync.LastBodySHA != expectedLastBodySHA {
					t.Errorf(
						"expected last body sha to be: '%s', got: '%s'", expectedLastBodySHA, httpSync.LastBodySHA,
					)
				}
			},
		},
		"unauthorized request": {
			setup: func(_ *testing.T, client *syncmock.MockClient) {
				client.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Header:     map[string][]string{"Content-Type": {"application/json"}},
					Body:       io.NopCloser(strings.NewReader("test response")),
					StatusCode: http.StatusUnauthorized,
				}, nil)
			},
			uri: "http://localhost",
			handleResponse: func(t *testing.T, _ Sync, _ string, err error) {
				if err == nil {
					t.Fatalf("expected unauthorized request to return an error")
				}
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := syncmock.NewMockClient(ctrl)

			tt.setup(t, mockClient)

			httpSync := Sync{
				URI:         tt.uri,
				Client:      mockClient,
				BearerToken: tt.bearerToken,
				AuthHeader:  tt.authHeader,
				LastBodySHA: tt.lastBodySHA,
				Logger:      logger.NewLogger(nil, false),
			}

			fetched, err := httpSync.Fetch(context.Background())
			tt.handleResponse(t, httpSync, fetched, err)
		})
	}
}

func TestSync_Init(t *testing.T) {
	tests := []struct {
		name        string
		bearerToken string
	}{
		{"with bearerToken", "bearer-1234"},
		{"without bearerToken", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpSync := Sync{
				BearerToken: tt.bearerToken,
				Logger:      logger.NewLogger(nil, false),
			}

			if err := httpSync.Init(context.Background()); err != nil {
				t.Errorf("Init() error = %v", err)
			}
		})
	}
}

func TestHTTPSync_Resync(t *testing.T) {
	ctrl := gomock.NewController(t)
	source := "http://localhost"
	emptyeFlagData := "{}"

	tests := map[string]struct {
		setup             func(t *testing.T, client *syncmock.MockClient)
		uri               string
		bearerToken       string
		lastBodySHA       string
		handleResponse    func(*testing.T, Sync, string, error)
		wantErr           bool
		wantNotifications []sync.DataSync
	}{
		"success": {
			setup: func(_ *testing.T, client *syncmock.MockClient) {
				client.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Header:     map[string][]string{"Content-Type": {"application/json"}},
					Body:       io.NopCloser(strings.NewReader(emptyeFlagData)),
					StatusCode: http.StatusOK,
				}, nil)
			},
			uri: source,
			handleResponse: func(t *testing.T, _ Sync, fetched string, err error) {
				if err != nil {
					t.Fatalf("fetch: %v", err)
				}
				expected := "test response"
				if fetched != expected {
					t.Errorf("expected fetched to be: '%s', got: '%s'", expected, fetched)
				}
			},
			wantErr: false,
			wantNotifications: []sync.DataSync{
				{
					FlagData: emptyeFlagData,
					Source:   source,
				},
			},
		},
		"error response": {
			setup: func(_ *testing.T, _ *syncmock.MockClient) {},
			handleResponse: func(t *testing.T, _ Sync, _ string, err error) {
				if err == nil {
					t.Error("expected err, got nil")
				}
			},
			wantErr:           true,
			wantNotifications: []sync.DataSync{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := syncmock.NewMockClient(ctrl)

			d := make(chan sync.DataSync, len(tt.wantNotifications))

			tt.setup(t, mockClient)

			httpSync := Sync{
				URI:         tt.uri,
				Client:      mockClient,
				BearerToken: tt.bearerToken,
				LastBodySHA: tt.lastBodySHA,
				Logger:      logger.NewLogger(nil, false),
			}

			err := httpSync.ReSync(context.Background(), d)
			if tt.wantErr && err == nil {
				t.Errorf("got no error for %s", name)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("got error for %s %s", name, err.Error())
			}
			for _, dataSync := range tt.wantNotifications {
				select {
				case x := <-d:
					if x.FlagData != dataSync.FlagData || x.Source != dataSync.Source {
						t.Errorf("unexpected datasync received %v vs %v", x, dataSync)
					}
				case <-time.After(2 * time.Second):
					t.Error("expected datasync not received", dataSync)
				}
			}
			select {
			case x := <-d:
				t.Error("unexpected datasync received", x)
			case <-time.After(2 * time.Second):
			}
		})
	}
}
