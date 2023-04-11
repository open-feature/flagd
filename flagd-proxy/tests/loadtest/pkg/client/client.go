package client

import (
	"fmt"

	syncv1 "buf.build/gen/go/open-feature/flagd/grpc/go/sync/v1/syncv1grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientConfig struct {
	Host string
	Port uint16
}

func NewClient(config ClientConfig) (syncv1.FlagSyncServiceClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", config.Host, config.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return syncv1.NewFlagSyncServiceClient(conn), nil
}
