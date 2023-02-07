package grpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"testing"

	"buf.build/gen/go/open-feature/flagd/grpc/go/sync/v1/syncv1grpc"
	v1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/sync/v1"

	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/sync"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

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

	tests := []struct {
		name   string
		stream syncv1grpc.FlagSyncService_SyncFlagsClient
		want   sync.Type
	}{
		{
			name: "State All maps to Sync All",
			stream: &SimpleRecvMock{
				mockResponse: v1.SyncFlagsResponse{
					FlagConfiguration: "{}",
					State:             v1.SyncState_SYNC_STATE_ALL,
				},
			},
			want: sync.ALL,
		},
		{
			name: "State Add maps to Sync Add",
			stream: &SimpleRecvMock{
				mockResponse: v1.SyncFlagsResponse{
					FlagConfiguration: "{}",
					State:             v1.SyncState_SYNC_STATE_ADD,
				},
			},
			want: sync.ADD,
		},
		{
			name: "State Update maps to Sync Update",
			stream: &SimpleRecvMock{
				mockResponse: v1.SyncFlagsResponse{
					FlagConfiguration: "{}",
					State:             v1.SyncState_SYNC_STATE_UPDATE,
				},
			},
			want: sync.UPDATE,
		},
		{
			name: "State Delete maps to Sync Delete",
			stream: &SimpleRecvMock{
				mockResponse: v1.SyncFlagsResponse{
					FlagConfiguration: "{}",
					State:             v1.SyncState_SYNC_STATE_DELETE,
				},
			},
			want: sync.DELETE,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			syncChan := make(chan sync.DataSync)

			go func() {
				err := grpcSyncImpl.handleFlagSync(test.stream, syncChan)
				if err != nil {
					t.Errorf("Error handling flag sync: %s", err.Error())
				}
			}()
			data := <-syncChan

			if data.Type != test.want {
				t.Errorf("Reuturned data sync state = %v, wanted %v", data.Type, test.want)
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
		}

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
		syncClient, err := serviceClient.SyncFlags(context.Background(), &v1.SyncFlagsRequest{ProviderId: grpcSync.ProviderID})
		if err != nil {
			t.Errorf("Error opening client stream: %s", err.Error())
		}

		syncChan := make(chan sync.DataSync, 1)

		// listen to stream
		go func() {
			err := grpcSync.streamListener(context.Background(), syncClient, syncChan)
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
				t.Errorf("Reuturned sync type = %v, wanted %v", out.Type, expected.Type)
			}

			if expected.FlagData != out.FlagData {
				t.Errorf("Reuturned sync data = %v, wanted %v", out.FlagData, expected.FlagData)
			}
		}

		// channel must be empty
		if len(syncChan) != 0 {
			t.Errorf("Data sync channel must be empty after all test syncs. But received non empty: %d", len(syncChan))
		}
	}
}

// Mock implementations

type SimpleRecvMock struct {
	grpc.ClientStream
	mockResponse v1.SyncFlagsResponse
}

func (s *SimpleRecvMock) Recv() (*v1.SyncFlagsResponse, error) {
	return &s.mockResponse, nil
}

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
	listener      *bufconn.Listener
	mockResponses []serverPayload
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
