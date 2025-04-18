package sync

import (
	"context"
	"fmt"
	"log"
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
		wantErr        bool
	}{
		{title: "with TLS Connection", certPath: "./test-cert/server-cert.pem", keyPath: "./test-cert/server-key.pem", clientCertPath: "./test-cert/ca-cert.pem", socketPath: "", tls: true, wantErr: false},
		{title: "witout TLS Connection", certPath: "", keyPath: "", clientCertPath: "", socketPath: "", tls: false, wantErr: false},
		{title: "with invalid TLS certificate path", certPath: "./lol/not/a/cert", keyPath: "./test-cert/server-key.pem", clientCertPath: "./test-cert/ca-cert.pem", socketPath: "", tls: true, wantErr: true},
		{title: "with unix socket connection", certPath: "", keyPath: "", clientCertPath: "", socketPath: "/tmp/flagd", tls: false, wantErr: false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Testing Sync Service %s", tc.title), func(t *testing.T) {
			// given
			port := 18016
			store, sources := getSimpleFlagStore()

			service, err := NewSyncService(SvcConfigurations{
				Logger:     logger.NewLogger(nil, false),
				Port:       uint16(port),
				Sources:    sources,
				Store:      store,
				CertPath:   tc.certPath,
				KeyPath:    tc.keyPath,
				SocketPath: tc.socketPath,
			})

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error creating the service!")
				}
				return
			} else if err != nil {
				t.Fatal("unexpected error creating the service: %w", err)
				return
			}

			ctx, cancelFunc := context.WithCancel(context.Background())
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

			// when - derive a client for sync service
			var con *grpc.ClientConn
			if tc.tls {
				tlsCredentials, e := loadTLSClientCredentials(tc.clientCertPath)
				if e != nil {
					log.Fatal("cannot load TLS credentials: ", e)
				}
				con, err = grpc.Dial(fmt.Sprintf("0.0.0.0:%d", port), grpc.WithTransportCredentials(tlsCredentials))
			} else {
				if tc.socketPath != "" {
					con, err = grpc.Dial(
						fmt.Sprintf("unix://%s", tc.socketPath),
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
				return
			}

			serviceClient := syncv1grpc.NewFlagSyncServiceClient(con)

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
