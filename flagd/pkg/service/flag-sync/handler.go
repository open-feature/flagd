package flag_sync

import (
	syncv1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/sync/v1"
	"connectrpc.com/connect"
	"context"
)

// Request handler
type syncHandler struct {
}

func (s syncHandler) SyncFlags(ctx context.Context, c *connect.Request[syncv1.SyncFlagsRequest], c2 *connect.ServerStream[syncv1.SyncFlagsResponse]) error {
	//TODO implement me
	panic("implement me")
}

func (s syncHandler) FetchAllFlags(ctx context.Context, c *connect.Request[syncv1.FetchAllFlagsRequest]) (*connect.Response[syncv1.FetchAllFlagsResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (s syncHandler) GetMetadata(ctx context.Context, c *connect.Request[syncv1.GetMetadataRequest]) (*connect.Response[syncv1.GetMetadataResponse], error) {
	//TODO implement me
	panic("implement me")
}
