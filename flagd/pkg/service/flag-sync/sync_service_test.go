package sync

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"buf.build/gen/go/open-feature/flagd/grpc/go/flagd/sync/v1/syncv1grpc"
	v1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/sync/v1"
	"github.com/open-feature/flagd/core/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestSyncServiceEndToEnd(t *testing.T) {
	testCases := []struct {
		title          string
		certPath       string
		keyPath        string
		clientCertPath string
		socketPath     string
		tls            bool
		wantStartErr   bool
	}{
		{title: "with TLS Connection", certPath: "./test-cert/server-cert.pem", keyPath: "./test-cert/server-key.pem", clientCertPath: "./test-cert/ca-cert.pem", socketPath: "", tls: true, wantStartErr: false},
		{title: "without TLS Connection", certPath: "", keyPath: "", clientCertPath: "", socketPath: "", tls: false, wantStartErr: false},
		{title: "with invalid TLS certificate path", certPath: "./lol/not/a/cert", keyPath: "./test-cert/server-key.pem", clientCertPath: "./test-cert/ca-cert.pem", socketPath: "", tls: true, wantStartErr: true},
		{title: "with unix socket connection", certPath: "", keyPath: "", clientCertPath: "", socketPath: "/tmp/flagd", tls: false, wantStartErr: false},
	}

	for _, disableSyncMetadata := range []bool{true, false} {
		for _, tc := range testCases {
			t.Run(fmt.Sprintf("Testing Sync Service %s", tc.title), func(t *testing.T) {
				// given
				port := 18016
				flagStore, sources := getSimpleFlagStore(t)

				ctx, cancelFunc := context.WithCancel(context.Background())
				defer cancelFunc()

				_, doneChan, err := createAndStartSyncService(
					port,
					sources,
					flagStore,
					tc.certPath,
					tc.keyPath,
					tc.socketPath,
					ctx,
					0,
					disableSyncMetadata,
				)

				if tc.wantStartErr {
					if err == nil {
						t.Fatal("expected error creating the service!")
					}
					return
				} else if err != nil {
					t.Fatal("unexpected error creating the service: %w", err)
					return
				}

				// when - derive a client for sync service
				serviceClient := getSyncClient(t, tc.clientCertPath, tc.socketPath, tc.tls, port, ctx)

				// then

				// sync flags request
				flags, err := serviceClient.SyncFlags(ctx, &v1.SyncFlagsRequest{})
				if err != nil {
					t.Fatal(fmt.Printf("error from sync request: %v", err))
					return
				}

				syncRsp, err := flags.Recv()
				if err != nil {
					t.Fatal(fmt.Printf("stream error: %v", err))
					return
				}

				if len(syncRsp.GetFlagConfiguration()) == 0 {
					t.Error("expected non empty sync response, but got empty")
				}

				// checks sync context actually set
				syncContext := syncRsp.GetSyncContext()
				if syncContext == nil {
					t.Fatal("expected sync_context in SyncFlagsResponse, but got nil")
				}

				// validate emits
				dataReceived := make(chan interface{})
				go func() {
					_, err := flags.Recv()
					if err != nil {
						return
					}

					dataReceived <- nil
				}()

				// make a change
				flagStore.Update(testSource1, testSource1Flags, model.Metadata{
					"keyDuped": "value",
					"keyA":     "valueA",
				}, false)

				select {
				case <-dataReceived:
					break
				case <-time.After(1 * time.Second):
					t.Fatal("expected data but timeout waiting for sync")
				}

				// fetch all flags
				allRsp, err := serviceClient.FetchAllFlags(ctx, &v1.FetchAllFlagsRequest{})
				if err != nil {
					t.Fatal(fmt.Printf("fetch all error: %v", err))
					return
				}

				if allRsp.GetFlagConfiguration() != syncRsp.GetFlagConfiguration() {
					t.Errorf("expected both sync and fetch all responses to be same, but got %s from sync & %s from fetch all",
						syncRsp.GetFlagConfiguration(), allRsp.GetFlagConfiguration())
				}

				// metadata request
				metadataRsp, err := serviceClient.GetMetadata(ctx, &v1.GetMetadataRequest{})

				if disableSyncMetadata {
					if err == nil {
						t.Fatal(fmt.Printf("getMetadata disabled, error should not be nil"))
						return
					}
				} else {
					asMap := metadataRsp.GetMetadata().AsMap()
					assert.NotNil(t, asMap, "expected metadata to be non-nil")
				}

				// validate shutdown from context cancellation
				go func() {
					cancelFunc()
				}()

				select {
				case <-doneChan:
					// exit successful
					return
				case <-time.After(2 * time.Second):
					t.Fatal("service did not exist within sufficient timeframe")
				}
			})
		}
	}
}

func TestSyncServiceDeadlineEndToEnd(t *testing.T) {
	testCases := []struct {
		title    string
		deadline time.Duration
	}{
		{title: "without deadline", deadline: 0},
		{title: "with deadline", deadline: 2 * time.Second},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Testing Sync Service %s", tc.title), func(t *testing.T) {

			// given
			port := 18016
			flagStore, sources := getSimpleFlagStore(t)
			certPath := "./test-cert/server-cert.pem"
			keyPath := "./test-cert/server-key.pem"
			socketPath := ""

			ctx, cancelFunc := context.WithCancel(context.Background())
			defer cancelFunc()

			_, _, err := createAndStartSyncService(port, sources, flagStore, certPath, keyPath, socketPath, ctx, tc.deadline, false)
			if err != nil {
				t.Fatal("error creating sync service")
			}

			// when - derive a client for sync service
			serviceClient := getSyncClient(t, "./test-cert/ca-cert.pem", "", true, port, nil)

			// then

			// sync flags request
			flags, err := serviceClient.SyncFlags(ctx, &v1.SyncFlagsRequest{})
			if err != nil {
				t.Fatal(fmt.Printf("error from sync request: %v", err))
				return
			}

			dataChan := make(chan any)
			errorChan := make(chan error)

			go func() {
				for {
					data, err := flags.Recv()
					dataChan <- data

					if err != nil {
						errorChan <- err
						return
					}
				}
			}()

			for {
				select {
				case <-dataChan:
					// received data, continuing..
					break
				case err := <-errorChan:
					st, _ := status.FromError(err)
					if st.Code() == codes.DeadlineExceeded {
						if tc.deadline == 0 {
							t.Fatal("ran into deadline exceeded error even though no deadline was configured.")
						}
						// expected error due to deadline
						return
					}
					t.Fatal("unexpected error: ", err)
				case <-time.After(tc.deadline + 1*time.Second):
					if tc.deadline == 0 {
						return
					}
					t.Fatal("not expected as the deadline should result in other cases.")
				}
			}
		})
	}
}

func TestKeepAliveEnforcementPolicy(t *testing.T) {
	tests := []struct {
		name                string
		minTime             time.Duration
		permitWithoutStream bool
	}{
		{name: "flag defaults", minTime: 30 * time.Second, permitWithoutStream: true},
		{name: "custom min time", minTime: 10 * time.Second, permitWithoutStream: true},
		{name: "permit without stream disabled", minTime: 30 * time.Second, permitWithoutStream: false},
		{name: "zero min time passed through unchanged", minTime: 0, permitWithoutStream: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := keepAliveEnforcementPolicy(SvcConfigurations{
				KeepAliveMinTime:             tt.minTime,
				KeepAlivePermitWithoutStream: tt.permitWithoutStream,
			})

			assert.Equal(t, tt.minTime, policy.MinTime)
			assert.Equal(t, tt.permitWithoutStream, policy.PermitWithoutStream)
		})
	}
}

// TestSyncServiceKeepAliveEnforcement proves the enforcement policy is applied
// to the sync gRPC server. The Go gRPC client clamps its keepalive interval to a
// 10s minimum, so we drive the HTTP/2 connection directly with a raw framer to
// flood keepalive pings and observe how the server's configured policy reacts:
// a permissive policy tolerates the pings, while a strict one tears the
// connection down with GOAWAY ENHANCE_YOUR_CALM.
func TestSyncServiceKeepAliveEnforcement(t *testing.T) {
	tests := []struct {
		name                string
		minTime             time.Duration
		permitWithoutStream bool
		wantEnhanceYourCalm bool
	}{
		{
			name:                "permissive policy tolerates frequent pings",
			minTime:             time.Millisecond,
			permitWithoutStream: true,
			wantEnhanceYourCalm: false,
		},
		{
			name:                "strict min time rejects frequent pings with GOAWAY",
			minTime:             time.Hour,
			permitWithoutStream: true,
			wantEnhanceYourCalm: true,
		},
		{
			name:                "pings without an active stream are rejected when not permitted",
			minTime:             time.Millisecond,
			permitWithoutStream: false,
			wantEnhanceYourCalm: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			port := 18017
			flagStore, sources := getSimpleFlagStore(t)

			ctx, cancelFunc := context.WithCancel(context.Background())

			service, err := NewSyncService(SvcConfigurations{
				Logger:                       logger.NewLogger(nil, false),
				Port:                         uint16(port),
				Sources:                      sources,
				Store:                        flagStore,
				KeepAliveMinTime:             tt.minTime,
				KeepAlivePermitWithoutStream: tt.permitWithoutStream,
			})
			require.NoError(t, err)

			serverDone := make(chan struct{})
			go func() {
				// error ignored, the assertions below fail if start does not succeed
				_ = service.Start(ctx)
				close(serverDone)
			}()
			// stop the server and wait for the listener to be released before the
			// next subtest reuses the port
			t.Cleanup(func() {
				cancelFunc()
				<-serverDone
			})
			for _, source := range sources {
				service.Emit(source)
			}

			gotEnhanceYourCalm := floodKeepalivePings(t, fmt.Sprintf("localhost:%d", port))
			assert.Equal(t, tt.wantEnhanceYourCalm, gotEnhanceYourCalm)
		})
	}
}

// floodKeepalivePings opens a raw HTTP/2 (h2c) connection to the sync server and
// sends keepalive PING frames in rapid succession without opening a stream. It
// returns true if the server responds with GOAWAY ENHANCE_YOUR_CALM within a
// short window, and false if the connection is left healthy.
func floodKeepalivePings(t *testing.T, addr string) bool {
	t.Helper()

	conn := dialWithRetry(t, addr)
	defer conn.Close()

	_, err := io.WriteString(conn, http2.ClientPreface)
	require.NoError(t, err)

	framer := http2.NewFramer(conn, conn)

	var writeMu sync.Mutex
	write := func(fn func() error) {
		writeMu.Lock()
		defer writeMu.Unlock()
		_ = fn()
	}

	// client connection preface must be followed by a SETTINGS frame
	write(func() error { return framer.WriteSettings() })

	enhanceYourCalm := watchForGoAway(framer, write)

	// flood pings faster than any strict MinTime; grpc-go GOAWAYs after more
	// than two "ping strikes", so a handful of rapid pings is enough
	var pingData [8]byte
	for i := 0; i < 8; i++ {
		write(func() error { return framer.WritePing(false, pingData) })
		select {
		case got := <-enhanceYourCalm:
			return got
		case <-time.After(10 * time.Millisecond):
		}
	}

	select {
	case got := <-enhanceYourCalm:
		return got
	case <-time.After(2 * time.Second):
		return false
	}
}

// dialWithRetry dials addr, retrying briefly because the listener is created in
// NewSyncService but Serve is delayed until the startup tracker completes.
func dialWithRetry(t *testing.T, addr string) net.Conn {
	t.Helper()

	var conn net.Conn
	var err error
	for i := 0; i < 50; i++ {
		conn, err = net.Dial("tcp", addr)
		if err == nil {
			return conn
		}
		time.Sleep(20 * time.Millisecond)
	}
	require.NoError(t, err)
	return conn
}

// watchForGoAway reads frames from the framer in the background, acking server
// SETTINGS, and reports on the returned channel whether the first GOAWAY carries
// the ENHANCE_YOUR_CALM code.
func watchForGoAway(framer *http2.Framer, write func(func() error)) <-chan bool {
	enhanceYourCalm := make(chan bool, 1)
	go func() {
		for {
			frame, err := framer.ReadFrame()
			if err != nil {
				return
			}
			switch f := frame.(type) {
			case *http2.SettingsFrame:
				if !f.IsAck() {
					write(func() error { return framer.WriteSettingsAck() })
				}
			case *http2.GoAwayFrame:
				enhanceYourCalm <- f.ErrCode == http2.ErrCodeEnhanceYourCalm
				return
			}
		}
	}()
	return enhanceYourCalm
}

func createAndStartSyncService(
	port int,
	sources []string,
	store store.IStore,
	certPath string,
	keyPath string,
	socketPath string,
	ctx context.Context,
	deadline time.Duration,
	disableSyncMetadata bool,
) (*Service, chan interface{}, error) {
	service, err := NewSyncService(SvcConfigurations{
		Logger:              logger.NewLogger(nil, false),
		Port:                uint16(port),
		Sources:             sources,
		Store:               store,
		CertPath:            certPath,
		KeyPath:             keyPath,
		SocketPath:          socketPath,
		StreamDeadline:      deadline,
		DisableSyncMetadata: disableSyncMetadata,
	})
	if err != nil {
		return nil, nil, err
	}

	doneChan := make(chan interface{})
	go func() {
		// error ignored, tests will fail if start is not successful
		_ = service.Start(ctx)
		close(doneChan)
	}()
	// trigger manual emits matching sources, so that service can start
	for _, source := range sources {
		service.Emit(source)
	}
	return service, doneChan, err
}

func getSyncClient(t *testing.T, clientCertPath string, socketPath string, tls bool, port int, ctx context.Context) syncv1grpc.FlagSyncServiceClient {
	var con *grpc.ClientConn
	var err error
	if tls {
		tlsCredentials, e := loadTLSClientCredentials(clientCertPath)
		if e != nil {
			log.Fatal("cannot load TLS credentials: ", e)
		}
		con, err = grpc.Dial(fmt.Sprintf("0.0.0.0:%d", port), grpc.WithTransportCredentials(tlsCredentials))
	} else {
		if socketPath != "" {
			con, err = grpc.Dial(
				fmt.Sprintf("unix://%s", socketPath),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithBlock(),
				grpc.WithTimeout(2*time.Second),
			)
		} else {
			con, err = grpc.DialContext(ctx, fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		}
	}
	if err != nil {
		t.Fatal(fmt.Printf("error creating grpc dial ctx: %v", err))
	}

	serviceClient := syncv1grpc.NewFlagSyncServiceClient(con)
	return serviceClient
}
