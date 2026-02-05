package client

import (
	"fmt"

	syncv1 "buf.build/gen/go/open-feature-forking/flagd/grpc/go/sync/v1/syncv1grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Host string
	Port uint16
}

func NewClient(config Config) (syncv1.FlagSyncServiceClient, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf(
			"%s:%d",
			config.Host,
			config.Port,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create client connection: %w", err)
	}
	return syncv1.NewFlagSyncServiceClient(conn), nil
}
