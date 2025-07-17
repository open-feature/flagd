//nolint:wrapcheck
package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"testing"
	"time"

	"buf.build/gen/go/open-feature/flagd/grpc/go/flagd/sync/v1/syncv1grpc"
	v1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/sync/v1"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	credendialsmock "github.com/open-feature/flagd/core/pkg/sync/grpc/credentials/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

func Test_InitWithMockCredentialBuilder(t *testing.T) {
	tests := []struct {
		name                       string
		mockCredentials            credentials.TransportCredentials
		mockCredentialBuilderError error
		shouldError                bool
	}{
		{
			name:                       "Initializer - happy path",
			mockCredentialBuilderError: nil,
			mockCredentials:            insecure.NewCredentials(),
			shouldError:                false,
		},
		{
			name:                       "Initializer fails on nil credentials",
			mockCredentialBuilderError: nil,
			mockCredentials:            nil,
			shouldError:                true,
		},
		{
			name:                       "Initializer handles credential builder error",
			mockCredentialBuilderError: errors.New("could not create transport credentials"),
			mockCredentials:            nil,
			shouldError:                true,
		},
	}

	for _, test := range tests {
		mockCtrl := gomock.NewController(t)
		mockCredentialBulder := credendialsmock.NewMockBuilder(mockCtrl)

		mockCredentialBulder.EXPECT().
			Build(gomock.Any(), gomock.Any()).
			Return(test.mockCredentials, test.mockCredentialBuilderError)

		grpcSync := Sync{
			URI:               "grpc-target",
			Logger:            logger.NewLogger(nil, false),
			CredentialBuilder: mockCredentialBulder,
		}

		err := grpcSync.Init(context.Background())

		if test.shouldError {
			require.NotNilf(t, err, "%s: expected an error", test.name)
		} else {
			require.Nilf(t, err, "%s: expected no error, but got non nil error", test.name)
			require.NotNilf(t, grpcSync.client, "%s: expected client to be initialized", test.name)
		}
	}
}

func Test_InitWithSizeOverride(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)

	mockCtrl := gomock.NewController(t)
	mockCredentialBulder := credendialsmock.NewMockBuilder(mockCtrl)

	mockCredentialBulder.EXPECT().
		Build(gomock.Any(), gomock.Any()).
		Return(insecure.NewCredentials(), nil)

	grpcSync := Sync{
		URI:               "grpc-target",
		Logger:            logger.NewLogger(observedLogger, false),
		CredentialBuilder: mockCredentialBulder,
		MaxMsgSize:        10,
	}

	err := grpcSync.Init(context.Background())

	require.Nilf(t, err, "%s: expected no error, but got non nil error", t.Name())
	require.Equal(t, "setting max receive message size 10 bytes default 4MB", observedLogs.All()[0].Message)
}

func Test_ReSyncTests(t *testing.T) {
	const target = "localBufCon"

	tests := []struct {
		name          string
		res           *v1.FetchAllFlagsResponse
		err           error
		shouldError   bool
		notifications []sync.DataSync
	}{
		{
			name: "happy-path",
			res: &v1.FetchAllFlagsResponse{
				FlagConfiguration: "success",
			},
			notifications: []sync.DataSync{
				{
					FlagData: "success",
				},
			},
			shouldError: false,
		},
		{
			name:          "happy-path",
			res:           &v1.FetchAllFlagsResponse{},
			err:           errors.New("internal disaster"),
			notifications: []sync.DataSync{},
			shouldError:   true,
		},
	}

	for _, test := range tests {
		bufCon := bufconn.Listen(5)

		bufServer := bufferedServer{
			listener:              bufCon,
			fetchAllFlagsResponse: test.res,
			fetchAllFlagsError:    test.err,
		}

		// start server
		go serve(&bufServer)

		// initialize client
		dial, err := grpc.Dial(target,
			grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
				return bufCon.Dial()
			}),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Errorf("Error setting up client connection: %s", err.Error())
		}

		c := syncv1grpc.NewFlagSyncServiceClient(dial)

		grpcSync := Sync{
			URI:        target,
			ProviderID: "",
			Logger:     logger.NewLogger(nil, false),
			client:     c,
		}

		syncChan := make(chan sync.DataSync, 1)

		err = grpcSync.ReSync(context.Background(), syncChan)
		if test.shouldError && err == nil {
			t.Errorf("test %s should have returned error but did not", test.name)
		}
		if !test.shouldError && err != nil {
			t.Errorf("test %s should not have returned error, but did: %s", test.name, err.Error())
		}

		for _, expected := range test.notifications {
			out := <-syncChan

			if expected.FlagData != out.FlagData {
				t.Errorf("Returned sync data = %v, wanted %v", out.FlagData, expected.FlagData)
			}
		}

		// channel must be empty
		if len(syncChan) != 0 {
			t.Errorf("Data sync channel must be empty after all test syncs. But received non empty: %d", len(syncChan))
		}
	}
}

func Test_StreamListener(t *testing.T) {
	const target = "localBufCon"

	metadata, err := structpb.NewStruct(map[string]any{"sources": "A,B,C"})
	if err != nil {
		t.Fatalf("Failed to create sync context: %v", err)
	}

	tests := []struct {
		name   string
		input  []serverPayload
		output []sync.DataSync
	}{
		{
			name: "Single send",
			input: []serverPayload{
				{
					flags: "{\"flags\": {}}",
				},
			},
			output: []sync.DataSync{
				{
					FlagData:    "{\"flags\": {}}",
					SyncContext: metadata,
				},
			},
		},
		{
			name: "Multiple send",
			input: []serverPayload{
				{
					flags: "{}",
				},
				{
					flags: "{\"flags\": {}}",
				},
			},
			output: []sync.DataSync{
				{
					FlagData:    "{}",
					SyncContext: metadata,
				},
				{
					FlagData:    "{\"flags\": {}}",
					SyncContext: metadata,
				},
			},
		},
	}

	for _, test := range tests {
		bufCon := bufconn.Listen(5)

		bufServer := bufferedServer{
			listener:      bufCon,
			mockResponses: test.input,
		}

		// start server
		go serve(&bufServer)

		// initialize client
		dial, err := grpc.Dial(target,
			grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
				return bufCon.Dial()
			}),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Errorf("Error setting up client connection: %s", err.Error())
		}

		serviceClient := syncv1grpc.NewFlagSyncServiceClient(dial)

		grpcSync := Sync{
			URI:        target,
			ProviderID: "",
			Logger:     logger.NewLogger(nil, false),

			client: serviceClient,
		}

		syncChan := make(chan sync.DataSync, 1)

		// listen to stream
		go func() {
			err := grpcSync.Sync(context.TODO(), syncChan)
			if err != nil {
				// must ignore EOF as this is returned for stream end
				if err != io.EOF {
					t.Errorf("Error from stream listener:  %s", err.Error())
				}
			}
		}()

		for _, expected := range test.output {
			out := <-syncChan

			if expected.FlagData != out.FlagData {
				t.Errorf("Returned sync data = %v, wanted %v", out.FlagData, expected.FlagData)
			}

			if !proto.Equal(expected.SyncContext, out.SyncContext) {
				t.Errorf("Returned sync context = %v, wanted = %v", out.SyncContext, expected.SyncContext)
			}
		}

		// channel must be empty
		if len(syncChan) != 0 {
			t.Errorf("Data sync channel must be empty after all test syncs. But received non empty: %d", len(syncChan))
		}
	}
}

// Test_ConnectWithRetry is an attempt to validate grpc.connectWithRetry behavior
func Test_ConnectWithRetry(t *testing.T) {
	target := "grpc://local"
	bufListener := bufconn.Listen(1)
	// buffer based server. response ignored purposefully
	bServer := bufferedServer{listener: bufListener}

	// generate a client connection backed with bufconn
	clientConn, err := grpc.Dial(target,
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return bufListener.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("error initiating the connection: %s", err.Error())
	}

	// minimal sync provider
	grpcSync := Sync{
		Logger: logger.NewLogger(nil, false),
		client: syncv1grpc.NewFlagSyncServiceClient(clientConn),
	}

	// test must complete within an acceptable timeframe
	tCtx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	// channel for connection
	clientChan := make(chan syncv1grpc.FlagSyncService_SyncFlagsClient)

	// start connection retry attempts
	go func() {
		client, ok := grpcSync.connectWithRetry(tCtx)
		if !ok {
			clientChan <- nil
		}

		clientChan <- client
	}()

	// Wait for retries in the background
	select {
	case <-time.After(2 * time.Second):
		break
	case <-tCtx.Done():
		// We should not reach this with correct test setup, but in case we do
		cancelFunc()
		t.Errorf("timeout occurred while waiting for conditions to fulfil")
	}

	// start the server - fulfill connection after the wait
	go serve(&bServer)

	// Wait for client connection
	var con syncv1grpc.FlagSyncService_SyncFlagsClient

	select {
	case con = <-clientChan:
		break
	case <-tCtx.Done():
		cancelFunc()
		t.Errorf("timeout occurred while waiting for conditions to fulfil")
	}

	if con == nil {
		t.Errorf("received a nil value when expecting a non-nil return")
	}
}

// Test_SyncRetry validates sync and retry attempts
func Test_SyncRetry(t *testing.T) {
	// Setup
	target := "grpc://local"
	bufListener := bufconn.Listen(1)
	emptyFlagData := "{}"

	// buffer based server. response ignored purposefully
	bServer := bufferedServer{listener: bufListener, mockResponses: []serverPayload{
		{
			flags: "{}",
		},
	}}

	// generate a client connection backed by bufListener
	clientConn, err := grpc.Dial(target,
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return bufListener.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("error initiating the connection: %s", err.Error())
	}

	// minimal sync provider
	grpcSync := Sync{
		Logger: logger.NewLogger(nil, false),
		client: syncv1grpc.NewFlagSyncServiceClient(clientConn),
	}

	// channel for data sync
	syncChan := make(chan sync.DataSync, 1)

	// Testing

	// Initial mock server - start mock server backed by a error group. Allow connection and disconnect with a timeout
	tCtx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()

	group, _ := errgroup.WithContext(tCtx)
	group.Go(func() error {
		serve(&bServer)
		return nil
	})

	// Start Sync for grpc streaming
	go func() {
		err := grpcSync.Sync(context.Background(), syncChan)
		if err != nil {
			t.Errorf("sync start error: %s", err.Error())
		}
	}()

	// Check for timeout (not ideal) or data sync (ideal) and cancel the context
	select {
	case <-tCtx.Done():
		t.Errorf("timeout waiting for conditions to fulfil")
		break
	case data := <-syncChan:
		if data.FlagData != emptyFlagData {
			t.Errorf("sync start error: %s", err.Error())
		}
	}

	// cancel make error group to complete, making background mock server to exit
	cancelFunc()

	// Follow up mock server start - start mock server after initial shutdown
	tCtx, cancelFunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	// Restart the server
	go serve(&bServer)

	// validate connection re-establishment
	select {
	case <-tCtx.Done():
		cancelFunc()
		t.Error("timeout waiting for conditions to fulfil")
	case data := <-syncChan:
		if data.FlagData != emptyFlagData {
			t.Errorf("sync start error: %s", err.Error())
		}
	}
}

// Mock implementations

// serve serves a bufferedServer. This is a blocking call
func serve(bServer *bufferedServer) {
	server := grpc.NewServer()

	syncv1grpc.RegisterFlagSyncServiceServer(server, bServer)

	if err := server.Serve(bServer.listener); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}

type serverPayload struct {
	flags string
}

// bufferedServer - a mock grpc service backed by buffered connection
type bufferedServer struct {
	listener              *bufconn.Listener
	mockResponses         []serverPayload
	fetchAllFlagsResponse *v1.FetchAllFlagsResponse
	fetchAllFlagsError    error
}

func (b *bufferedServer) SyncFlags(_ *v1.SyncFlagsRequest, stream syncv1grpc.FlagSyncService_SyncFlagsServer) error {
	for _, response := range b.mockResponses {
		metadata, _ := structpb.NewStruct(map[string]any{"sources": "A,B,C"})
		err := stream.Send(&v1.SyncFlagsResponse{
			FlagConfiguration: response.flags,
			SyncContext:       metadata,
		})
		if err != nil {
			fmt.Printf("Error with stream: %s", err.Error())
			return err
		}
	}

	return nil
}

func (b *bufferedServer) FetchAllFlags(_ context.Context, _ *v1.FetchAllFlagsRequest) (*v1.FetchAllFlagsResponse, error) {
	return b.fetchAllFlagsResponse, b.fetchAllFlagsError
}

func (b *bufferedServer) GetMetadata(_ context.Context, _ *v1.GetMetadataRequest) (*v1.GetMetadataResponse, error) {
	return &v1.GetMetadataResponse{}, nil
}
