package http

import (
	"context"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

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

	mockClient := syncmock.NewMockClient(ctrl)
	mockClient.EXPECT().Do(gomock.Any()).Return(&http.Response{Body: io.NopCloser(strings.NewReader(resp))}, nil)

	httpSync := Sync{
		URI:         "http://localhost",
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

	if data.FlagData != resp {
		t.Errorf("expected content: %s, but received content: %s", resp, data.FlagData)
	}
}

func TestHTTPSync_Fetch(t *testing.T) {
	ctrl := gomock.NewController(t)

	tests := map[string]struct {
		setup          func(t *testing.T, client *syncmock.MockClient)
		uri            string
		bearerToken    string
		lastBodySHA    string
		handleResponse func(*testing.T, Sync, string, error)
		ready          bool
	}{
		"success": {
			setup: func(t *testing.T, client *syncmock.MockClient) {
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
					t.Errorf("expected fetched to be: '%s', got: '%s'", expected, fetched)
				}
			},
			ready: true,
		},
		"return an error if no uri": {
			setup: func(t *testing.T, client *syncmock.MockClient) {},
			handleResponse: func(t *testing.T, _ Sync, fetched string, err error) {
				if err == nil {
					t.Error("expected err, got nil")
				}
			},
			ready: false,
		},
		"update last body sha": {
			setup: func(t *testing.T, client *syncmock.MockClient) {
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

				expectedLastBodySHA := "UjeJHtCU_wb7OHK-tbPoHycw0TqlHzkWJmH4y6cqg50="
				if httpSync.LastBodySHA != expectedLastBodySHA {
					t.Errorf(
						"expected last body sha to be: '%s', got: '%s'", expectedLastBodySHA, httpSync.LastBodySHA,
					)
				}
			},
			ready: true,
		},
		"authorization header": {
			setup: func(t *testing.T, client *syncmock.MockClient) {
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

				expectedLastBodySHA := "UjeJHtCU_wb7OHK-tbPoHycw0TqlHzkWJmH4y6cqg50="
				if httpSync.LastBodySHA != expectedLastBodySHA {
					t.Errorf(
						"expected last body sha to be: '%s', got: '%s'", expectedLastBodySHA, httpSync.LastBodySHA,
					)
				}
			},
			ready: true,
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
				LastBodySHA: tt.lastBodySHA,
				Logger:      logger.NewLogger(nil, false),
			}

			fetched, err := httpSync.Fetch(context.Background())
			if httpSync.IsReady() != tt.ready {
				t.Errorf("expected httpSync.ready to be: '%v', got: '%v'", tt.ready, httpSync.ready)
			}
			tt.handleResponse(t, httpSync, fetched, err)
		})
	}
}

func TestHTTPSync_Resync(t *testing.T) {
	ctrl := gomock.NewController(t)

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
			setup: func(t *testing.T, client *syncmock.MockClient) {
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
					t.Errorf("expected fetched to be: '%s', got: '%s'", expected, fetched)
				}
			},
			wantErr: false,
			wantNotifications: []sync.DataSync{
				{
					Type:     sync.ALL,
					FlagData: "",
					Source:   "",
				},
			},
		},
		"error response": {
			setup: func(t *testing.T, client *syncmock.MockClient) {},
			handleResponse: func(t *testing.T, _ Sync, fetched string, err error) {
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
					if !reflect.DeepEqual(x.String(), dataSync.String()) {
						t.Error("unexpected datasync received", x, dataSync)
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
