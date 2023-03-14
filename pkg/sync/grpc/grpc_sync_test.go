package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	grpcmock "github.com/open-feature/flagd/pkg/sync/grpc/mock"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net"
	"testing"
	"time"

	"buf.build/gen/go/open-feature/flagd/grpc/go/sync/v1/syncv1grpc"
	v1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/sync/v1"

	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/sync"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

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

		grpcSync := Sync{
			Target:     target,
			ProviderID: "",
			Logger:     logger.NewLogger(nil, false),
			options: []grpc.DialOption{
				grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
					return bufCon.Dial()
				}),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			},
		}

		err := grpcSync.connectClient(context.Background())

		require.Nil(t, err)

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

func TestUrlToGRPCTarget(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "With Prefix",
			url:  "grpc://test.com/endpoint",
			want: "test.com/endpoint",
		},
		{
			name: "Without Prefix",
			url:  "test.com/endpoint",
			want: "test.com/endpoint",
		},
		{
			name: "Empty is empty",
			url:  "",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := URLToGRPCTarget(tt.url); got != tt.want {
				t.Errorf("URLToGRPCTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSync_BasicFlagSyncStates(t *testing.T) {
	grpcSyncImpl := Sync{
		Target:     "grpc://test",
		ProviderID: "",
		Logger:     logger.NewLogger(nil, false),
	}

	mockError := errors.New("could not sync")

	tests := []struct {
		name      string
		stream    syncv1grpc.FlagSyncService_SyncFlagsClient
		setup     func(t *testing.T, client *grpcmock.MockFlagSyncServiceSyncFlagsClient)
		want      sync.Type
		wantError error
		ready     bool
	}{
		{
			name: "State All maps to Sync All",
			setup: func(t *testing.T, client *grpcmock.MockFlagSyncServiceSyncFlagsClient) {
				gomock.InOrder(
					client.EXPECT().Recv().Return(
						&v1.SyncFlagsResponse{
							FlagConfiguration: "{}",
							State:             v1.SyncState_SYNC_STATE_ALL,
						},
						nil,
					),
					client.EXPECT().Recv().Return(
						nil, io.EOF,
					),
				)
			},
			want:  sync.ALL,
			ready: true,
		},
		{
			name: "State Add maps to Sync Add",
			setup: func(t *testing.T, client *grpcmock.MockFlagSyncServiceSyncFlagsClient) {
				gomock.InOrder(
					client.EXPECT().Recv().Return(
						&v1.SyncFlagsResponse{
							FlagConfiguration: "{}",
							State:             v1.SyncState_SYNC_STATE_ADD,
						},
						nil,
					),
					client.EXPECT().Recv().Return(
						nil, io.EOF,
					),
				)
			},
			want:  sync.ADD,
			ready: true,
		},
		{
			name: "State Update maps to Sync Update",
			setup: func(t *testing.T, client *grpcmock.MockFlagSyncServiceSyncFlagsClient) {
				gomock.InOrder(
					client.EXPECT().Recv().Return(
						&v1.SyncFlagsResponse{
							FlagConfiguration: "{}",
							State:             v1.SyncState_SYNC_STATE_UPDATE,
						},
						nil,
					),
					client.EXPECT().Recv().Return(
						nil, io.EOF,
					),
				)
			},
			want:  sync.UPDATE,
			ready: true,
		},
		{
			name: "State Delete maps to Sync Delete",
			setup: func(t *testing.T, client *grpcmock.MockFlagSyncServiceSyncFlagsClient) {
				gomock.InOrder(
					client.EXPECT().Recv().Return(
						&v1.SyncFlagsResponse{
							FlagConfiguration: "{}",
							State:             v1.SyncState_SYNC_STATE_DELETE,
						},
						nil,
					),
					client.EXPECT().Recv().Return(
						nil, io.EOF,
					),
				)
			},
			want:  sync.DELETE,
			ready: true,
		},
		{
			name: "Error during flag sync",
			setup: func(t *testing.T, client *grpcmock.MockFlagSyncServiceSyncFlagsClient) {
				client.EXPECT().Recv().Return(
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

			mockClient := grpcmock.NewMockFlagSyncServiceSyncFlagsClient(ctrl)
			test.setup(t, mockClient)

			waitChan := make(chan struct{})
			go func() {
				grpcSyncImpl.syncClient = mockClient
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()
				err := grpcSyncImpl.Sync(ctx, syncChan)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
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

		grpcSync := Sync{
			Target:     target,
			ProviderID: "",
			Logger:     logger.NewLogger(nil, false),
			options: []grpc.DialOption{
				grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
					return bufCon.Dial()
				}),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			},
		}

		err := grpcSync.connectClient(context.Background())

		require.Nil(t, err)

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

// Mock implementations

// serve serves a bufferedServer
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
