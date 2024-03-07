package sync

import (
	"context"
	"fmt"
	"testing"
	"time"

	"buf.build/gen/go/open-feature/flagd/grpc/go/flagd/sync/v1/syncv1grpc"
	v1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/sync/v1"
	"github.com/open-feature/flagd/core/pkg/logger"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestSyncServiceEndToEnd(t *testing.T) {
	// given
	port := 18016
	store, sources := getSimpleFlagStore()

	service, err := NewSyncService(SvcConfigurations{
		Logger:  logger.NewLogger(nil, false),
		Port:    uint16(port),
		Sources: sources,
		Store:   store,
	})
	if err != nil {
		t.Fatal("error creating the service: %w", err)
		return
	}

	group, ctx := errgroup.WithContext(context.Background())
	group.Go(func() error {
		err := service.Start()
		if err != nil {
			return err
		}
		return nil
	})

	// when - derive a client for sync service
	con, err := grpc.DialContext(ctx, fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(fmt.Printf("error creating grpc dial ctx: %v", err))
		return
	}

	serviceClient := syncv1grpc.NewFlagSyncServiceClient(con)

	// then

	// sync flags
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

	if asMap["sources"] != "A,B" {
		t.Fatal("incorrect sources entry in metadata")
	}

	//

	// validate shutdown
	go func() {
		service.Shutdown()
	}()

	select {
	case <-ctx.Done():
		return
	case <-time.After(2 * time.Second):
		t.Fatal("service did not exist within sufficient timeframe")
	}
}
