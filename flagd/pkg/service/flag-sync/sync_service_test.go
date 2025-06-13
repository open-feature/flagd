package sync

import (
	"context"
	"fmt"
	"github.com/open-feature/flagd/core/pkg/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

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

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Testing Sync Service %s", tc.title), func(t *testing.T) {
			// given
			port := 18016
			flagStore, sources := getSimpleFlagStore()

			ctx, cancelFunc := context.WithCancel(context.Background())
			defer cancelFunc()

			service, doneChan, err := createAndStartSyncService(port, sources, flagStore, tc.certPath, tc.keyPath, tc.socketPath, ctx, 0)

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

			syncAsMap := syncContext.AsMap()
			if syncAsMap["sources"] == nil {
				t.Fatalf("expected sources in sync_context, but got nil")
			}

			sourcesStr := syncAsMap["sources"].(string)
			sourcesArray := strings.Split(sourcesStr, ",")
			sort.Strings(sourcesArray)

			expectedSources := []string{"A", "B", "C"}
			if !reflect.DeepEqual(sourcesArray, expectedSources) {
				t.Fatalf("sources entry in sync_context does not match expected: got %v, want %v", sourcesArray, expectedSources)
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

			// Emit as a resync
			service.Emit(true, "A")

			select {
			case <-dataReceived:
				t.Fatal("expected no data as this is a resync")
			case <-time.After(1 * time.Second):
				break
			}

			// Emit as a resync
			service.Emit(false, "A")

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
			if err != nil {
				t.Fatal(fmt.Printf("metadata error: %v", err))
				return
			}

			asMap := metadataRsp.GetMetadata().AsMap()

			// expect `sources` to be present
			if asMap["sources"] == nil {
				t.Fatal("expected sources entry in the metadata, but got nil")
			}

			if asMap["sources"] != "A,B,C" {
				t.Fatal("incorrect sources entry in metadata")
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
			flagStore, sources := getSimpleFlagStore()
			certPath := "./test-cert/server-cert.pem"
			keyPath := "./test-cert/server-key.pem"
			socketPath := ""

			ctx, cancelFunc := context.WithCancel(context.Background())
			defer cancelFunc()

			_, _, err := createAndStartSyncService(port, sources, flagStore, certPath, keyPath, socketPath, ctx, tc.deadline)
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

func createAndStartSyncService(port int, sources []string, store *store.State, certPath string, keyPath string, socketPath string, ctx context.Context, deadline time.Duration) (*Service, chan interface{}, error) {
	service, err := NewSyncService(SvcConfigurations{
		Logger:         logger.NewLogger(nil, false),
		Port:           uint16(port),
		Sources:        sources,
		Store:          store,
		CertPath:       certPath,
		KeyPath:        keyPath,
		SocketPath:     socketPath,
		StreamDeadline: deadline,
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
		service.Emit(false, source)
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
