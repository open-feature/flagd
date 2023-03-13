package sync

import (
	"context"
	"fmt"
	"net/http"

	rpc "buf.build/gen/go/open-feature/flagd/bufbuild/connect-go/sync/v1/syncv1connect"
	syncv1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/sync/v1"
	"github.com/bufbuild/connect-go"
	"github.com/open-feature/flagd/core/pkg/sync"
	syncStore "github.com/open-feature/flagd/core/pkg/sync-store"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type SyncServer struct {
	SyncStore syncStore.SyncStore
	Port      uint16
}

func (s *SyncServer) Serve(ctx context.Context) {
	mux := http.NewServeMux()
	path, handler := rpc.NewFlagSyncServiceHandler(s)
	mux.Handle(path, handler)
	http.ListenAndServe(
		fmt.Sprintf(":%d", s.Port),
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}

func (s *SyncServer) FetchAllFlags(context.Context, *connect.Request[syncv1.FetchAllFlagsRequest]) (*connect.Response[syncv1.FetchAllFlagsResponse], error) {
	return nil, nil
}

func (s *SyncServer) SyncFlags(ctx context.Context, req *connect.Request[syncv1.SyncFlagsRequest], stream *connect.ServerStream[syncv1.SyncFlagsResponse]) error {
	errChan := make(chan error)
	dataSync := make(chan sync.DataSync)
	s.SyncStore.RegisterSubscription(ctx, req.Msg.GetProviderId(), req, dataSync, errChan)
	for {
		select {
		case e := <-errChan:
			fmt.Println("\n\nerror\n\n", e)
			return e
		case d := <-dataSync:
			fmt.Println("received data sync type ", d.String())
			if err := stream.Send(&syncv1.SyncFlagsResponse{
				FlagConfiguration: d.FlagData,
				State:             syncv1.SyncState(d.Type + 1),
			}); err != nil {
				fmt.Println("under")
				return err
			}
			fmt.Println("after")
		case <-ctx.Done():
			fmt.Println("connection closed")
			return nil
		}
	}
}
