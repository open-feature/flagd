package main

import (
	"context"
	"fmt"
	"time"

	pb "buf.build/gen/go/open-feature/flagd/grpc/go/sync/v1/syncv1grpc"
	schemav1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/sync/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//nolint:staticcheck
func doRequests(grpcClient pb.FlagSyncServiceClient, waitSecondsBetweenRequests int) error {
	ctx := context.TODO()
	stream, err := grpcClient.SyncFlags(ctx, &schemav1.SyncFlagsRequest{
		ProviderId: "zd",
		Selector:   "file:/etc/flagd/config.json",
	})
	if err != nil {
		return fmt.Errorf("%s", "error SyncFlags(): "+err.Error())
	}

	for {
		// We do not care about the message received, only the error and then we try to re-connect.
		// If the re-connection fails; the server is down and ZD test should fail
		_, err = stream.Recv()
		if err != nil {
			fmt.Println("error Recv(): " + err.Error())
			stream, err = grpcClient.SyncFlags(ctx, &schemav1.SyncFlagsRequest{
				ProviderId: "zd",
				Selector:   "file:/etc/flagd/config.json",
			})
			if err != nil {
				return fmt.Errorf("%s", "error SyncFlags(): "+err.Error())
			}
		}
		<-time.After(time.Duration(waitSecondsBetweenRequests) * time.Second)
	}
}

func establishGrpcConnection(url string) (*grpc.ClientConn, pb.FlagSyncServiceClient) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err.Error())
	}
	client := pb.NewFlagSyncServiceClient(conn)
	return conn, client
}
