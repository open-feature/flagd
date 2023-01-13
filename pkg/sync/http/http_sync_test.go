package http

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/sync"

	"github.com/golang/mock/gomock"
	syncmock "github.com/open-feature/flagd/pkg/sync/http/mock"
)

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

func TestHTTPSync_Notify(t *testing.T) {
	ctrl := gomock.NewController(t)

	tests := map[string]struct {
		setup             func(t *testing.T, cron *syncmock.MockCron, client *syncmock.MockHTTPClient)
		uri               string
		lastBodySHA       string
		expectedEventType sync.DefaultEventType
	}{
		"create event": {
			setup: func(t *testing.T, cron *syncmock.MockCron, client *syncmock.MockHTTPClient) {
				var cronFunc func()
				cron.EXPECT().AddFunc(gomock.Any(), gomock.Any()).DoAndReturn(func(_ string, f func()) error {
					cronFunc = f
					return nil
				})
				cron.EXPECT().Start().DoAndReturn(func() { cronFunc() })
				cron.EXPECT().Stop().AnyTimes()
				client.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Body: io.NopCloser(strings.NewReader("test response")),
				}, nil)
			},
			uri:               "http://localhost",
			expectedEventType: sync.DefaultEventTypeCreate,
		},
		"modify event": {
			setup: func(t *testing.T, cron *syncmock.MockCron, client *syncmock.MockHTTPClient) {
				var cronFunc func()
				cron.EXPECT().AddFunc(gomock.Any(), gomock.Any()).DoAndReturn(func(_ string, f func()) error {
					cronFunc = f
					return nil
				})
				cron.EXPECT().Start().DoAndReturn(func() { cronFunc() })
				cron.EXPECT().Stop().AnyTimes()
				client.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Body: io.NopCloser(strings.NewReader("foo")),
				}, nil)
			},
			uri:               "http://localhost",
			expectedEventType: sync.DefaultEventTypeModify,
			lastBodySHA:       "fUH6MbDL8tR0nCiC4bag0Rf_6is=",
		},
		"delete event": {
			setup: func(t *testing.T, cron *syncmock.MockCron, client *syncmock.MockHTTPClient) {
				var cronFunc func()
				cron.EXPECT().AddFunc(gomock.Any(), gomock.Any()).DoAndReturn(func(_ string, f func()) error {
					cronFunc = f
					return nil
				})
				cron.EXPECT().Start().DoAndReturn(func() { cronFunc() })
				cron.EXPECT().Stop().AnyTimes()
				client.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Body: io.NopCloser(strings.NewReader("")),
				}, nil)
			},
			uri:               "http://localhost",
			expectedEventType: sync.DefaultEventTypeDelete,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// prevent deadlock with a timeout if expected event doesn't arrive
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			mockCron := syncmock.NewMockCron(ctrl)
			mockClient := syncmock.NewMockHTTPClient(ctrl)

			inotifyChan := make(chan sync.INotify)
			tt.setup(t, mockCron, mockClient)

			httpSync := Sync{
				URI:         tt.uri,
				Client:      mockClient,
				Cron:        mockCron,
				LastBodySHA: tt.lastBodySHA,
				Logger:      logger.NewLogger(nil, false),
			}

			go func() {
				httpSync.Notify(ctx, inotifyChan)
			}()

			w := <-inotifyChan // first emitted event by Notify is to signal readiness
			if w.GetEvent().EventType != sync.DefaultEventTypeReady {
				t.Errorf("expected event type to be %d, got %d", sync.DefaultEventTypeReady, w.GetEvent().EventType)
			}

			for {
				select {
				case event, ok := <-inotifyChan:
					if !ok {
						t.Fatal("inotify chan closed")
					}
					if event.GetEvent().EventType != tt.expectedEventType {
						t.Errorf(
							"expected event of type %d, got %d", tt.expectedEventType, event.GetEvent().EventType,
						)
					}
					return
				case <-ctx.Done():
					t.Error("context timed out")
					return
				}
			}
		})
	}
}
