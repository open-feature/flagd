package sync

import (
	"context"
	"sync"
	"testing"
	"time"

	"buf.build/gen/go/open-feature/flagd/grpc/go/flagd/sync/v1/syncv1grpc"
	syncv1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/sync/v1"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			require.NoError(t, err)

			handler := syncHandler{
				mux:           mp,
				contextValues: tt.contextValues,
				log:           logger.NewLogger(nil, false),
			}

			// Test getting metadata from `GetMetadata` (deprecated)
			// remove when `GetMetadata` is full removed and deprecated
			metaResp, err := handler.GetMetadata(context.Background(), &syncv1.GetMetadataRequest{})
			require.NoError(t, err)
			respMetadata := metaResp.GetMetadata().AsMap()
			assert.Equal(t, tt.wantMetadata, respMetadata)

			// Test metadata from sync_context
			stream := &mockSyncFlagsServer{
				ctx:       context.Background(),
				mu:        sync.Mutex{},
				respReady: make(chan struct{}, 1),
			}

			go func() {
				err := handler.SyncFlags(&syncv1.SyncFlagsRequest{}, stream)
				assert.NoError(t, err)
			}()

			select {
			case <-stream.respReady:
				syncResp := stream.GetLastResponse()
				assert.NotNil(t, syncResp)

				syncMetadata := syncResp.GetSyncContext().AsMap()
				assert.Equal(t, tt.wantMetadata, syncMetadata)

				// Check the two metadatas are equal
				// remove when `GetMetadata` is full removed and deprecated
				assert.Equal(t, respMetadata, syncMetadata)
			case <-time.After(time.Second):
				t.Fatal("timeout waiting for response")
			}

		})
	}
}

// Mock server for testing
type mockSyncFlagsServer struct {
	syncv1grpc.FlagSyncService_SyncFlagsServer
	ctx       context.Context
	mu        sync.Mutex
	lastResp  *syncv1.SyncFlagsResponse
	respReady chan struct{}
}

func (m *mockSyncFlagsServer) Context() context.Context {
	return m.ctx
}

func (m *mockSyncFlagsServer) Send(resp *syncv1.SyncFlagsResponse) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lastResp = resp
	select {
	case m.respReady <- struct{}{}:
	default:
	}
	return nil
}

func (m *mockSyncFlagsServer) GetLastResponse() *syncv1.SyncFlagsResponse {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.lastResp
}
