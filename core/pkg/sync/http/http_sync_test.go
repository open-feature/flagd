package http

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	syncmock "github.com/open-feature/flagd/core/pkg/sync/http/mock"
	synctesting "github.com/open-feature/flagd/core/pkg/sync/testing"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func buildHeaders(m map[string][]string) http.Header {
	h := http.Header{}
	for k, v := range m {
		for _, val := range v {
			h.Add(k, val)
		}
	}
	return h
}

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
		Header:     buildHeaders(map[string][]string{"Content-Type": {"application/json"}}),
		Body:       io.NopCloser(strings.NewReader(responseBody)),
		StatusCode: http.StatusOK,
	}
	mockClient.EXPECT().Do(gomock.Any()).Return(resp, nil)

	httpSync := Sync{
		uri:         "http://localhost/flags",
		client:      mockClient,
		cron:        mockCron,
		lastBodySHA: "",
		logger:      logger.NewLogger(nil, false),
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
		Header:     buildHeaders(map[string][]string{"Content-Type": {"application/json"}}),
		Body:       io.NopCloser(strings.NewReader(responseBody)),
		StatusCode: http.StatusOK,
	}
	mockClient.EXPECT().Do(gomock.Any()).Return(resp, nil)

	httpSync := Sync{
		uri:         "http://localhost/flags.json?env=dev",
		client:      mockClient,
		cron:        mockCron,
		lastBodySHA: "",
		logger:      logger.NewLogger(nil, false),
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
		authHeader     string
		eTagHeader     string
		lastBodySHA    string
		handleResponse func(*testing.T, Sync, string, error)
	}{
		"success": {
			setup: func(_ *testing.T, client *syncmock.MockClient) {
				client.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Header:     buildHeaders(map[string][]string{"Content-Type": {"application/json"}}),
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
					Header:     buildHeaders(map[string][]string{"Content-Type": {"application/json"}}),
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
				if httpSync.lastBodySHA != expectedLastBodySHA {
					t.Errorf(
						"expected last body sha to be: '%s', got: '%s'", expectedLastBodySHA, httpSync.lastBodySHA,
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
						Header:     buildHeaders(map[string][]string{"Content-Type": {"application/json"}}),
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
				if httpSync.lastBodySHA != expectedLastBodySHA {
					t.Errorf(
						"expected last body sha to be: '%s', got: '%s'", expectedLastBodySHA, httpSync.lastBodySHA,
					)
				}
			},
		},
		"unauthorized request": {
			setup: func(_ *testing.T, client *syncmock.MockClient) {
				client.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Header:     buildHeaders(map[string][]string{"Content-Type": {"application/json"}}),
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
		"not modified response etag matched": {
			setup: func(t *testing.T, client *syncmock.MockClient) {
				expectedIfNoneMatch := `"1af17a664e3fa8e419b8ba05c2a173169df76162a5a286e0c405b460d478f7ef"`
				client.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
					actualIfNoneMatch := req.Header.Get("If-None-Match")
					if actualIfNoneMatch != expectedIfNoneMatch {
						t.Fatalf("expected If-None-Match header to be '%s', got %s", expectedIfNoneMatch, actualIfNoneMatch)
					}
					return &http.Response{
						Header:     buildHeaders(map[string][]string{"ETag": {expectedIfNoneMatch}}),
						Body:       io.NopCloser(strings.NewReader("")),
						StatusCode: http.StatusNotModified,
					}, nil
				})
			},
			uri:        "http://localhost",
			eTagHeader: `"1af17a664e3fa8e419b8ba05c2a173169df76162a5a286e0c405b460d478f7ef"`,
			handleResponse: func(t *testing.T, httpSync Sync, _ string, err error) {
				if err != nil {
					t.Fatalf("fetch: %v", err)
				}

				expectedLastBodySHA := ""
				expectedETag := `"1af17a664e3fa8e419b8ba05c2a173169df76162a5a286e0c405b460d478f7ef"`
				if httpSync.lastBodySHA != expectedLastBodySHA {
					t.Errorf(
						"expected last body sha to be: '%s', got: '%s'", expectedLastBodySHA, httpSync.lastBodySHA,
					)
				}
				if httpSync.eTag != expectedETag {
					t.Errorf(
						"expected last etag to be: '%s', got: '%s'", expectedETag, httpSync.eTag,
					)
				}
			},
		},
		"modified response etag mismatched": {
			setup: func(t *testing.T, client *syncmock.MockClient) {
				expectedIfNoneMatch := `"1af17a664e3fa8e419b8ba05c2a173169df76162a5a286e0c405b460d478f7ef"`
				client.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
					actualIfNoneMatch := req.Header.Get("If-None-Match")
					if actualIfNoneMatch != expectedIfNoneMatch {
						t.Fatalf("expected If-None-Match header to be '%s', got %s", expectedIfNoneMatch, actualIfNoneMatch)
					}

					newContent := "\"Hey there!\""
					newETag := `"c2e01ce63d90109c4c7f4f6dcea97ed1bb2b51e3647f36caf5acbe27413a24bb"`

					return &http.Response{
						Header: buildHeaders(map[string][]string{
							"Content-Type": {"application/json"},
							"ETag":         {newETag},
						}),
						Body:       io.NopCloser(strings.NewReader(newContent)),
						StatusCode: http.StatusOK,
					}, nil
				})
			},
			uri:        "http://localhost",
			eTagHeader: `"1af17a664e3fa8e419b8ba05c2a173169df76162a5a286e0c405b460d478f7ef"`,
			handleResponse: func(t *testing.T, httpSync Sync, _ string, err error) {
				if err != nil {
					t.Fatalf("fetch: %v", err)
				}

				expectedLastBodySHA := "wuAc5j2QEJxMf09tzql-0bsrUeNkfzbK9ay-J0E6JLs="
				expectedETag := `"c2e01ce63d90109c4c7f4f6dcea97ed1bb2b51e3647f36caf5acbe27413a24bb"`
				if httpSync.lastBodySHA != expectedLastBodySHA {
					t.Errorf(
						"expected last body sha to be: '%s', got: '%s'", expectedLastBodySHA, httpSync.lastBodySHA,
					)
				}
				if httpSync.eTag != expectedETag {
					t.Errorf(
						"expected last etag to be: '%s', got: '%s'", expectedETag, httpSync.eTag,
					)
				}
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := syncmock.NewMockClient(ctrl)

			tt.setup(t, mockClient)

			httpSync := Sync{
				uri:         tt.uri,
				client:      mockClient,
				authHeader:  tt.authHeader,
				lastBodySHA: tt.lastBodySHA,
				logger:      logger.NewLogger(nil, false),
				eTag:        tt.eTagHeader,
			}

			fetched, err := httpSync.Fetch(context.Background())
			tt.handleResponse(t, httpSync, fetched, err)
		})
	}
}

func TestHTTPSync_Resync(t *testing.T) {
	ctrl := gomock.NewController(t)
	source := "http://localhost"
	emptyFlagData := "{}"

	tests := map[string]struct {
		setup             func(t *testing.T, client *syncmock.MockClient)
		uri               string
		lastBodySHA       string
		handleResponse    func(*testing.T, Sync, string, error)
		wantErr           bool
		wantNotifications []sync.DataSync
	}{
		"success": {
			setup: func(_ *testing.T, client *syncmock.MockClient) {
				client.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Header:     buildHeaders(map[string][]string{"Content-Type": {"application/json"}}),
					Body:       io.NopCloser(strings.NewReader(emptyFlagData)),
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
					FlagData: emptyFlagData,
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
				uri:         tt.uri,
				client:      mockClient,
				lastBodySHA: tt.lastBodySHA,
				logger:      logger.NewLogger(nil, false),
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

func TestHTTPSync_getClient(t *testing.T) {
	oauth := &sync.OAuthCredentialHandler{
		ClientID:     "myClientID",
		ClientSecret: "myClientSecret",
		TokenURL:     "http://localhost",
	}
	oauthDelay := &sync.OAuthCredentialHandler{
		ClientID:     "myClientID",
		ClientSecret: "myClientSecret",
		TokenURL:     "http://localhost",
		ReloadDelayS: 10000,
	}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	oauthClientCredential := &clientcredentials.Config{
		ClientID:     oauth.ClientID,
		ClientSecret: oauth.ClientSecret,
		TokenURL:     oauth.TokenURL,
		AuthStyle:    oauth2.AuthStyleAutoDetect,
	}
	tests := map[string]struct {
		config sync.SourceConfig
		client *http.Client
	}{
		"no http client no oauth": {
			config: sync.SourceConfig{},
		},
		"no http client yes oauth": {
			config: sync.SourceConfig{
				OAuth: oauth,
			},
		},
		"no http client yes oauth reload": {
			config: sync.SourceConfig{
				OAuth: oauthDelay,
			},
		},
		"yes http client no oauth": {
			config: sync.SourceConfig{},
			client: client,
		},
		"yes http client yes oauth": {
			config: sync.SourceConfig{
				OAuth: oauth,
			},
			client: client,
		},
		"yes http client yes oauth reload": {
			config: sync.SourceConfig{
				OAuth: oauthDelay,
			},
			client: oauthClientCredential.Client(context.Background()),
		},
	}

	l := logger.NewLogger(nil, false)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			httpSync := NewHTTP(tt.config, l)
			if tt.client != nil {
				// we have a cached HTTP client already
				httpSync.client = tt.client
				httpClient, ok := httpSync.getClient().(*http.Client)
				require.True(t, ok, "expected http client")
				if tt.config.OAuth != nil {
					// we use oauth so client should be different
					require.IsType(t, &oauth2.Transport{}, httpClient.Transport)
				} else {
					// we don't use oauth so client should be the same
					require.Equal(t, tt.client, httpClient)
				}

			}
			require.NotNil(t, httpSync.getClient())
		})
	}
}

type oauthHttpMock struct {
	count      int
	lastParams url.Values
	lastHeader http.Header
}

func TestHTTPSync_OAuth(t *testing.T) {
	// given
	const clientID = "clientID"
	const clientSecret = "clientSecret"
	const oauthPath = "/oauth"
	const bearerToken = "mySecretBearerToken"

	tests := map[string]struct {
		oauthResponse          int
		expectedHttpCallCount  int
		expectedOauthCallCount int
		expectedBeaerToken     string
	}{
		"success": {
			oauthResponse:          http.StatusOK,
			expectedHttpCallCount:  2,
			expectedOauthCallCount: 1,
			expectedBeaerToken:     fmt.Sprintf("Bearer %s", bearerToken),
		},
		"oauth error": {
			oauthResponse:          http.StatusInternalServerError,
			expectedHttpCallCount:  0,
			expectedOauthCallCount: 2,
			expectedBeaerToken:     "",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			oauthMock := &oauthHttpMock{}
			httpMock := &oauthHttpMock{}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.HasSuffix(r.URL.Path, oauthPath) {
					body, _ := io.ReadAll(r.Body)
					params, _ := url.ParseQuery(string(body))
					oauthMock.lastParams = params
					oauthMock.lastHeader = r.Header
					oauthMock.count++
					w.WriteHeader(tt.oauthResponse)
					w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
					_, err := w.Write([]byte(fmt.Sprintf("access_token=%s&scope=mockscope&token_type=bearer", bearerToken)))
					if err != nil {
						t.Fatalf("cannot write response: %v", err)
					}
					return
				}
				httpMock.lastHeader = r.Header
				httpMock.count++
				w.WriteHeader(http.StatusInternalServerError)
			}))
			defer ts.Close()
			l := logger.NewLogger(nil, false)
			s := NewHTTP(sync.SourceConfig{
				URI:         ts.URL,
				AuthHeader: "Bearer it_should_be_replaced_by_oauth",
				OAuth: &sync.OAuthCredentialHandler{
					ClientID:     clientID,
					ClientSecret: clientSecret,
					TokenURL:     ts.URL + oauthPath,
					ReloadDelayS: 10000,
				},
			}, l)
			d := make(chan sync.DataSync, 1)
			// when we call resync multiple times
			err := s.ReSync(context.Background(), d)
			err2 := s.ReSync(context.Background(), d)
			// then

			// we only get errors when fetching the HTTP data
			require.ErrorContains(t, err, "500 Internal Server Error")
			require.ErrorContains(t, err2, "500 Internal Server Error")

			// the OAuth endpoint is the right amount of times
			require.Equal(t, tt.expectedOauthCallCount, oauthMock.count)

			// the Beaerer token is replaced by the OAuth values
			require.Equal(t, tt.expectedBeaerToken, httpMock.lastHeader.Get("Authorization"))
			require.Equal(t, tt.expectedHttpCallCount, httpMock.count)
		})
	}
}

func TestHTTPSync_OAuthFolderSecrets(t *testing.T) {
	// given
	const (
		clientID               = "clientID"
		clientSecret           = "clientSecret"
		oauthPath              = "/oauth"
		bearerToken            = "mySecretBearerToken"
		folderName             = "flagd_oauth_test"
		secretClientID         = "mySecretClientID"
		secretClientSecret     = "mySecretClientSecret"
		secretClientID_new     = "newClientID"
		secretClientSecret_new = "newClientSecret"
	)
	// given an oauth server
	oauthMock := &oauthHttpMock{}
	flagsPath := "/flags.json"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, oauthPath) {
			// mock token response
			body, _ := io.ReadAll(r.Body)
			params, _ := url.ParseQuery(string(body))
			oauthMock.lastParams = params
			oauthMock.lastHeader = r.Header
			oauthMock.count++
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
			_, err := w.Write([]byte(fmt.Sprintf("access_token=%s&scope=mockscope&token_type=bearer", bearerToken)))
			if err != nil {
				t.Fatalf("cannot write response: %v", err)
			}
			return
		} else if strings.HasSuffix(r.URL.Path, flagsPath) {
			// mock flags response
			io.ReadAll(r.Body)

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(fmt.Sprintf(`{"flagKey": {"default": true}}`)))
			if err != nil {
				t.Fatalf("cannot write response: %v", err)
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()
	// given credentials stored on a folder
	dir, err := os.MkdirTemp("", folderName)
	defer func(path string) {
		e := os.RemoveAll(path)
		if e != nil {
			fmt.Printf("Cannot delete %s: %v", path, e)
		}
	}(dir)
	require.NoError(t, err)

	err = os.WriteFile(dir+"/client-id", []byte(secretClientID), 0644)
	require.NoError(t, err)
	err = os.WriteFile(dir+"/client-secret", []byte(secretClientSecret), 0644)
	require.NoError(t, err)

	l := logger.NewLogger(nil, false)
	s := NewHTTP(sync.SourceConfig{
		URI:         ts.URL + flagsPath,
		AuthHeader: "Bearer it_should_be_replaced_by_oauth",
		OAuth: &sync.OAuthCredentialHandler{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Folder:       dir,
			TokenURL:     ts.URL + oauthPath,
			ReloadDelayS: 0, // we force loading the secret at each req
		},
	}, l)
	d := make(chan sync.DataSync, 2)
	// when we fire the HTTP call
	err = s.ReSync(context.Background(), d)
	if err != nil {
		t.Fatalf("resync failed: %v", err)
	}

	// then the right secrets are used
	require.Equal(t, secretClientID, oauthMock.lastParams.Get("client_id"))
	require.Equal(t, secretClientSecret, oauthMock.lastParams.Get("client_secret"))

	// when we change the secrets
	err = os.WriteFile(dir+"/client-id", []byte(secretClientID_new), 0644)
	require.NoError(t, err)
	err = os.WriteFile(dir+"/client-secret", []byte(secretClientSecret_new), 0644)
	require.NoError(t, err)

	// then the new HTTP call will use the new value
	err = s.ReSync(context.Background(), d)
	if err != nil {
		t.Fatalf("resync failed: %v", err)
	}
	require.Equal(t, secretClientID_new, oauthMock.lastParams.Get("client_id"))
	require.Equal(t, secretClientSecret_new, oauthMock.lastParams.Get("client_secret"))
}
