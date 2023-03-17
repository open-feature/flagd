package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"

	"buf.build/gen/go/open-feature/flagd/grpc/go/sync/v1/syncv1grpc"
	v1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/sync/v1"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	grpcmock "github.com/open-feature/flagd/core/pkg/sync/grpc/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const sampleCert = `-----BEGIN CERTIFICATE-----
MIIEnDCCAoQCCQCHcl3hGXwRQzANBgkqhkiG9w0BAQsFADAQMQ4wDAYDVQQDDAVm
bGFnZDAeFw0yMzAyMTAxODM1NDVaFw0zMzAyMDcxODM1NDVaMBAxDjAMBgNVBAMM
BWZsYWdkMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAwDLEAUti/kG9
MhJLtO7oAy7diHxWKDFmsIHrE+z2IzTxjXxVHQLv1HiYB/UN75y7qlb3MwvzSc+C
BoLuoiM0PDiMio9/o9X5j0U+v3H1JpUU5LardkvsprFqJWmHF+D7aRdM0LBLn2X6
HQOhSnPyH9Qjl2l2tyPiPTZ6g0i2+rXZsNUoTs4fm6ThhZ0LeXR8KDmCTun3ze1d
hXA7ydxwILH2OVc+Wnzl30+BRvOiLQbc9nYnwSREFeIy8sFbhrTHqSNn3eY79ssZ
T6f4tN3jEV1d7NqoFk9KFLJKJhMt7smMB9NLwVWi581Zj1krYirNlP6mtmPrn3kJ
lsgT15kFftShMVcYFSHqOSLiy4SspHGK8KJaFoEVx0wp/weRwrWXi6vWg7tuHATH
fw7gW/9CyV+ylc0pJ002wtPAgzJYUaOrna0R2r3yQsSzRcDnqsm4FLkPHLoyjrwQ
vshKcEqjhGml1M+lTDEo3RO5ZoQ3ZN2AZKPDrK2zGG4wFJjHRu9FtutOEZkYYOzA
emTQWW8US3q8WVQqGl/EwQqzXk9Lco7uhLdXmqVOvAi6z01gehQJPnjhH7iqAPVp
1tlOBHit1F3sTAQIO/2zff3LCKiD2d27KINh4aFEyDbDmglPA8VPO3BMQVSjFlxj
K1s2G1IDBixXK76VmBP+ZpvxOaQtYIUCAwEAATANBgkqhkiG9w0BAQsFAAOCAgEA
K9+wnl5gpkfNBa+OSxlhOn3CKhcaW/SWZ4aLw2yK1NZNnNjpwUcLQScUDBKDoJJR
5roc3PIImX7hdnobZWqFhD23laaAlu5XLk9P7n51uMEiNjQQc2WaaBZDTRJfki1C
MvPskXqptgPsVyuPJc0DxfaCz7pDYjq/CtJ+osaj404P5mlO1QJ8W91QSx+aq2x4
uUTUWuyr/8flIcxiX0o8VTb2LcUvWZBMGa3CdeLnPHrOjovfjJFy0Ysk3SGEACLL
9mpbNbv23v9UXVfyFffHpyzvyUJIOsNXG0O1AYf5t9bukqHolGR/RQUN4yGd3M62
mFR5bOST36DjNSzTrx1eyCLv22+h9VVlWFPrebFnq1W5SSi8PtsGSMjhvX7dB1kS
t0yJtlj2HwBAvI1zVKG76q6neSU51UXFQUbO0OA0sxjicEOlNfXnShM/kY2lobpX
hrCysWpqoSS0S3UBvmuRiraLWkP1KueC0XHoAi8yuwMAdM6Y+h2OJpnO0PdpUmrp
lAqdxbyICnB1Nsm5QGGm6Pxd8lEbQ9ZSwFjgqApjT2zVhuaaUC7jdlEP1H5snt9n
8FQR06lrzGyW04ud9pd6MXJup1oghAlvnzXioAH2Az0IXcHvqUGZQattFv27OXqj
QZ6ayNO119SNscvC6Qe9GLlbBEHDQWKPiftnS2Mh6Do=
-----END CERTIFICATE-----`

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
					Type:     sync.ALL,
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
			if expected.Type != out.Type {
				t.Errorf("Returned sync type = %v, wanted %v", out.Type, expected.Type)
			}

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

func TestSourceToGRPCTarget(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
		ok   bool
	}{
		{
			name: "With Prefix",
			url:  "grpc://test.com/endpoint",
			want: "test.com/endpoint",
			ok:   true,
		},
		{
			name: "With secure Prefix",
			url:  "grpcs://test.com/endpoint",
			want: "test.com/endpoint",
			ok:   true,
		},
		{
			name: "Empty is error",
			url:  "",
			want: "",
			ok:   false,
		},
		{
			name: "Invalid is error",
			url:  "https://test.com/endpoint",
			want: "",
			ok:   false,
		},
		{
			name: "Prefix is not enough I",
			url:  Prefix,
			want: "",
			ok:   false,
		},
		{
			name: "Prefix is not enough II",
			url:  PrefixSecure,
			want: "",
			ok:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := sourceToGRPCTarget(tt.url)

			if tt.ok != ok {
				t.Errorf("URLToGRPCTarget() returned = %v, want %v", ok, tt.ok)
			}

			if got != tt.want {
				t.Errorf("URLToGRPCTarget() returned = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSync_BasicFlagSyncStates(t *testing.T) {
	grpcSyncImpl := Sync{
		URI:        "grpc://test",
		ProviderID: "",
		Logger:     logger.NewLogger(nil, false),
	}

	mockError := errors.New("could not sync")

	tests := []struct {
		name      string
		stream    syncv1grpc.FlagSyncService_SyncFlagsClient
		setup     func(t *testing.T, client *grpcmock.MockFlagSyncServiceClient, clientResponse *grpcmock.MockFlagSyncServiceClientResponse)
		want      sync.Type
		wantError error
		ready     bool
	}{
		{
			name: "State All maps to Sync All",
			setup: func(t *testing.T, client *grpcmock.MockFlagSyncServiceClient, clientResponse *grpcmock.MockFlagSyncServiceClientResponse) {
				client.EXPECT().SyncFlags(gomock.Any(), gomock.Any(), gomock.Any()).Return(clientResponse, nil)
				gomock.InOrder(
					clientResponse.EXPECT().Recv().Return(
						&v1.SyncFlagsResponse{
							FlagConfiguration: "{}",
							State:             v1.SyncState_SYNC_STATE_ALL,
						},
						nil,
					),
					clientResponse.EXPECT().Recv().Return(
						nil, io.EOF,
					),
				)
			},
			want:  sync.ALL,
			ready: true,
		},
		{
			name: "State Add maps to Sync Add",
			setup: func(t *testing.T, client *grpcmock.MockFlagSyncServiceClient, clientResponse *grpcmock.MockFlagSyncServiceClientResponse) {
				client.EXPECT().SyncFlags(gomock.Any(), gomock.Any(), gomock.Any()).Return(clientResponse, nil)
				gomock.InOrder(
					clientResponse.EXPECT().Recv().Return(
						&v1.SyncFlagsResponse{
							FlagConfiguration: "{}",
							State:             v1.SyncState_SYNC_STATE_ADD,
						},
						nil,
					),
					clientResponse.EXPECT().Recv().Return(
						nil, io.EOF,
					),
				)
			},
			want:  sync.ADD,
			ready: true,
		},
		{
			name: "State Update maps to Sync Update",
			setup: func(t *testing.T, client *grpcmock.MockFlagSyncServiceClient, clientResponse *grpcmock.MockFlagSyncServiceClientResponse) {
				client.EXPECT().SyncFlags(gomock.Any(), gomock.Any(), gomock.Any()).Return(clientResponse, nil)
				gomock.InOrder(
					clientResponse.EXPECT().Recv().Return(
						&v1.SyncFlagsResponse{
							FlagConfiguration: "{}",
							State:             v1.SyncState_SYNC_STATE_UPDATE,
						},
						nil,
					),
					clientResponse.EXPECT().Recv().Return(
						nil, io.EOF,
					),
				)
			},
			want:  sync.UPDATE,
			ready: true,
		},
		{
			name: "State Delete maps to Sync Delete",
			setup: func(t *testing.T, client *grpcmock.MockFlagSyncServiceClient, clientResponse *grpcmock.MockFlagSyncServiceClientResponse) {
				client.EXPECT().SyncFlags(gomock.Any(), gomock.Any(), gomock.Any()).Return(clientResponse, nil)
				gomock.InOrder(
					clientResponse.EXPECT().Recv().Return(
						&v1.SyncFlagsResponse{
							FlagConfiguration: "{}",
							State:             v1.SyncState_SYNC_STATE_DELETE,
						},
						nil,
					),
					clientResponse.EXPECT().Recv().Return(
						nil, io.EOF,
					),
				)
			},
			want:  sync.DELETE,
			ready: true,
		},
		{
			name: "Error during flag sync",
			setup: func(t *testing.T, client *grpcmock.MockFlagSyncServiceClient, clientResponse *grpcmock.MockFlagSyncServiceClientResponse) {
				client.EXPECT().SyncFlags(gomock.Any(), gomock.Any(), gomock.Any()).Return(clientResponse, nil)
				clientResponse.EXPECT().Recv().Return(
					nil,
					mockError,
				)
			},
			ready: true,
			want:  -1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			syncChan := make(chan sync.DataSync, 1)

			mockClient := grpcmock.NewMockFlagSyncServiceClient(ctrl)
			mockClientResponse := grpcmock.NewMockFlagSyncServiceClientResponse(ctrl)
			test.setup(t, mockClient, mockClientResponse)

			waitChan := make(chan struct{})
			go func() {
				grpcSyncImpl.client = mockClient
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()
				err := grpcSyncImpl.Sync(ctx, syncChan)
				if err != nil {
					t.Errorf("Error handling flag sync: %v", err)
				}
				close(waitChan)
			}()
			<-waitChan

			if test.want < 0 {
				require.Empty(t, syncChan)
				return
			}
			data := <-syncChan

			if grpcSyncImpl.IsReady() != test.ready {
				t.Errorf("expected grpcSyncImpl.ready to be: '%v', got: '%v'", test.ready, grpcSyncImpl.ready)
			}

			if data.Type != test.want {
				t.Errorf("Returned data sync state = %v, wanted %v", data.Type, test.want)
			}
		})
	}
}

func Test_StreamListener(t *testing.T) {
	const target = "localBufCon"

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
					state: v1.SyncState_SYNC_STATE_ALL,
				},
			},
			output: []sync.DataSync{
				{
					FlagData: "{\"flags\": {}}",
					Type:     sync.ALL,
				},
			},
		},
		{
			name: "Multiple send",
			input: []serverPayload{
				{
					flags: "{}",
					state: v1.SyncState_SYNC_STATE_ALL,
				},
				{
					flags: "{\"flags\": {}}",
					state: v1.SyncState_SYNC_STATE_DELETE,
				},
			},
			output: []sync.DataSync{
				{
					FlagData: "{}",
					Type:     sync.ALL,
				},
				{
					FlagData: "{\"flags\": {}}",
					Type:     sync.DELETE,
				},
			},
		},
		{
			name: "Pings are ignored & not written to channel",
			input: []serverPayload{
				{
					flags: "",
					state: v1.SyncState_SYNC_STATE_PING,
				},
				{
					flags: "",
					state: v1.SyncState_SYNC_STATE_PING,
				},
				{
					flags: "{\"flags\": {}}",
					state: v1.SyncState_SYNC_STATE_DELETE,
				},
			},
			output: []sync.DataSync{
				{
					FlagData: "{\"flags\": {}}",
					Type:     sync.DELETE,
				},
			},
		},
		{
			name: "Unknown states are & not written to channel",
			input: []serverPayload{
				{
					flags: "",
					state: 42,
				},
				{
					flags: "",
					state: -1,
				},
				{
					flags: "{\"flags\": {}}",
					state: v1.SyncState_SYNC_STATE_ALL,
				},
			},
			output: []sync.DataSync{
				{
					FlagData: "{\"flags\": {}}",
					Type:     sync.ALL,
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

			if expected.Type != out.Type {
				t.Errorf("Returned sync type = %v, wanted %v", out.Type, expected.Type)
			}

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

func Test_BuildTCredentials(t *testing.T) {
	// "insecure" is a hardcoded term at insecure.NewCredentials
	const insecure = "insecure"
	// "tls" is a hardcoded term at tlsCreds.Info
	const tls = "tls"
	// local test file with valid certificate
	const validCertFile = "valid.cert"
	// local test file with invalid certificate
	const invalidCertFile = "invalid.cert"

	// init cert files for tests & cleanup with a deffer
	err := os.WriteFile(validCertFile, []byte(sampleCert), 0o600)
	if err != nil {
		t.Errorf("error creating valid certificate file: %s", err)
	}

	err = os.WriteFile(invalidCertFile, []byte("--certificate--"), 0o600)
	if err != nil {
		t.Errorf("error creating invalid certificate file: %s", err)
	}

	defer func() {
		errV := os.Remove(validCertFile)
		errI := os.Remove(invalidCertFile)
		if errV != nil || errI != nil {
			t.Errorf("error removing cerificate files: %v, %v", errV, errI)
		}
	}()

	tests := []struct {
		name           string
		source         string
		certPath       string
		expectSecProto string
		error          bool
	}{
		{
			name:           "Insecure source results in insecure connection",
			source:         Prefix + "some.domain",
			certPath:       "",
			expectSecProto: insecure,
		},
		{
			name:           "Secure source results in secure connection",
			source:         PrefixSecure + "some.domain",
			certPath:       validCertFile,
			expectSecProto: tls,
		},
		{
			name:           "Secure source with no certificate results in a secure connection",
			source:         PrefixSecure + "some.domain",
			expectSecProto: tls,
		},
		{
			name:     "Invalid cert path results in an error",
			source:   PrefixSecure + "some.domain",
			certPath: "invalid/path",
			error:    true,
		},
		{
			name:     "Invalid certificate results in an error",
			source:   PrefixSecure + "some.domain",
			certPath: invalidCertFile,
			error:    true,
		},
		{
			name:   "Invalid prefix results in an error",
			source: "http://some.domain",
			error:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tCred, err := buildTransportCredentials(test.source, test.certPath)

			if test.error {
				if err == nil {
					t.Errorf("test expected non error execution. But resulted in an error: %s", err.Error())
				}

				// Test expected an error. Nothing to validate further
				return
			}

			// check for errors to be certain
			if err != nil {
				t.Errorf("unexpected error: %s", err.Error())
			}

			protoc := tCred.Info().SecurityProtocol
			if protoc != test.expectSecProto {
				t.Errorf("buildTransportCredentials() returned protocol= %v, want %v", protoc, test.expectSecProto)
			}
		})
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

	expectType := sync.ALL

	// buffer based server. response ignored purposefully
	bServer := bufferedServer{listener: bufListener, mockResponses: []serverPayload{
		{
			flags: "{}",
			state: v1.SyncState_SYNC_STATE_ALL,
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
		if data.Type != expectType {
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
	case rsp := <-syncChan:
		if rsp.Type != expectType {
			t.Errorf("expected response: %s, but got: %s", expectType, rsp.Type)
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
	state v1.SyncState
}

// bufferedServer - a mock grpc service backed by buffered connection
type bufferedServer struct {
	listener              *bufconn.Listener
	mockResponses         []serverPayload
	fetchAllFlagsResponse *v1.FetchAllFlagsResponse
	fetchAllFlagsError    error
}

func (b *bufferedServer) SyncFlags(req *v1.SyncFlagsRequest, stream syncv1grpc.FlagSyncService_SyncFlagsServer) error {
	for _, response := range b.mockResponses {
		err := stream.Send(&v1.SyncFlagsResponse{
			FlagConfiguration: response.flags,
			State:             response.state,
		})
		if err != nil {
			fmt.Printf("Error with stream: %s", err.Error())
			return err
		}
	}

	return nil
}

func (b *bufferedServer) FetchAllFlags(ctx context.Context, req *v1.FetchAllFlagsRequest) (*v1.FetchAllFlagsResponse, error) {
	return b.fetchAllFlagsResponse, b.fetchAllFlagsError
}
