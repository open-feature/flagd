package sync

import (
	"context"
	"testing"
	"time"

	"buf.build/gen/go/open-feature/flagd/grpc/go/flagd/sync/v1/syncv1grpc"
	syncv1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/sync/v1"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/stretchr/testify/assert"
)

func TestSyncHandler_SyncFlags(t *testing.T) {
	tests := []struct {
		name          string
		sources       []string
		contextValues map[string]any
		wantMetadata  map[string]any
	}{
		{
			name:    "with sources and context",
			sources: []string{"A, B, C"},
			contextValues: map[string]any{
				"env":    "prod",
				"region": "us-west",
			},
			wantMetadata: map[string]any{
				"sources": "A, B, C",
				"env":     "prod",
				"region":  "us-west",
			},
		},
		{
			name:    "with empty sources",
			sources: []string{},
			contextValues: map[string]any{
				"env": "dev",
			},
			wantMetadata: map[string]any{
				"env": "dev",
			},
		},
		{
			name:          "with empty context",
			sources:       []string{"A,B,C"},
			contextValues: map[string]any{},
			wantMetadata: map[string]any{
				"sources": "A,B,C",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Shared handler for testing both GetMetadata & SyncFlags methods
			flagStore := store.NewFlags()
			mp, err := NewMux(flagStore, tt.sources)
			if err != nil {
				t.Fatalf("Failed to create multiplexer: %v", err)
			}

			handler := syncHandler{
				mux:           mp,
				contextValues: tt.contextValues,
				log:           logger.NewLogger(nil, false),
			}

			// Test getting metadata from `GetMetadata` (deprecated)
			// remove when `GetMetadata` is full removed and deprecated
			resp, err := handler.GetMetadata(context.Background(), &syncv1.GetMetadataRequest{})
			assert.NoError(t, err)
			respMetadata := resp.GetMetadata().AsMap()
			assert.Equal(t, tt.wantMetadata, respMetadata)

			// Test metadata from sync_context
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			stream := &mockSyncFlagsServer{ctx: ctx}

			go func() {
				err = handler.SyncFlags(&syncv1.SyncFlagsRequest{}, stream)
				assert.NoError(t, err)
			}()

			// A pause so the handler has time to fully register
			time.Sleep(50 * time.Millisecond)

			syncResp := stream.lastResp
			assert.NotNil(t, syncResp)

			syncMetadata := syncResp.GetSyncContext().AsMap()
			assert.Equal(t, tt.wantMetadata, syncMetadata)

			// Check the two metadatas are equal
			// remove when `GetMetadata` is full removed and deprecated
			assert.Equal(t, respMetadata, syncMetadata)
		})
	}
}

// Mock server for testing
type mockSyncFlagsServer struct {
	syncv1grpc.FlagSyncService_SyncFlagsServer
	ctx      context.Context
	lastResp *syncv1.SyncFlagsResponse
}

func (m *mockSyncFlagsServer) Context() context.Context {
	return m.ctx
}

func (m *mockSyncFlagsServer) Send(resp *syncv1.SyncFlagsResponse) error {
	m.lastResp = resp
	return nil
}
