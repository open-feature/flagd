package sync

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"buf.build/gen/go/open-feature/flagd/grpc/go/flagd/sync/v1/syncv1grpc"
	syncv1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/sync/v1"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	flagdService "github.com/open-feature/flagd/flagd/pkg/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestSyncHandler_SyncFlags(t *testing.T) {
	tests := []struct {
		name          string
		sources       []string
		contextValues map[string]any
		wantMetadata  map[string]any
	}{
		{
			name: "with sources and context",
			contextValues: map[string]any{
				"env":    "prod",
				"region": "us-west",
			},
			wantMetadata: map[string]any{
				"env":    "prod",
				"region": "us-west",
			},
		},
		{
			name: "with empty sources",
			contextValues: map[string]any{
				"env": "dev",
			},
			wantMetadata: map[string]any{
				"env": "dev",
			},
		},
		{
			name:          "with empty context",
			contextValues: map[string]any{},
			wantMetadata:  map[string]any{},
		},
	}

	for _, disableSyncMetadata := range []bool{true, false} {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Shared handler for testing both GetMetadata & SyncFlags methods
				flagStore, err := store.NewStore(logger.NewLogger(nil, false), tt.sources)
				require.NoError(t, err)

				handler := syncHandler{
					store:               flagStore,
					contextValues:       tt.contextValues,
					log:                 logger.NewLogger(nil, false),
					disableSyncMetadata: disableSyncMetadata,
					metricsRecorder:     &telemetry.NoopMetricsRecorder{},
				}

				// Test getting metadata from `GetMetadata` (deprecated)
				// remove when `GetMetadata` is full removed and deprecated
				metaResp, err := handler.GetMetadata(context.Background(), &syncv1.GetMetadataRequest{})
				if !disableSyncMetadata {
					require.NoError(t, err)
					respMetadata := metaResp.GetMetadata().AsMap()
					assert.Equal(t, tt.wantMetadata, respMetadata)
				} else {
					assert.NotNil(t, err)
				}

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
				case <-time.After(time.Second):
					t.Fatal("timeout waiting for response")
				}
			})
		}
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

// Test that selector from header takes precedence over selector from request body for FetchAllFlags and SyncFlags methods.
func TestSyncHandler_SelectorLocationPrecedence(t *testing.T) {
	headerFlags := []model.Flag{
		{
			Key:            "header-flag",
			State:          "ENABLED",
			DefaultVariant: "true",
			Variants:       testVariants,
		},
	}

	bodyFlags := []model.Flag{
		{
			Key:            "body-flag",
			State:          "DISABLED",
			DefaultVariant: "false",
			Variants:       testVariants,
		},
	}

	tests := []struct {
		name             string
		hasHeader        bool
		headerSelector   string
		bodySelector     string
		expectedFlag     string
		expectedSource   string
		shouldNotContain string
	}{
		{
			name:             "SyncFlags with request body selector only",
			hasHeader:        false,
			bodySelector:     "source=body-source",
			expectedFlag:     "body-flag",
			expectedSource:   "body-source",
			shouldNotContain: "header-flag",
		},
		{
			name:             "SyncFlags header takes precedence over request body",
			hasHeader:        true,
			headerSelector:   "source=header-source",
			bodySelector:     "source=body-source",
			expectedFlag:     "header-flag",
			expectedSource:   "header-source",
			shouldNotContain: "body-flag",
		},
		{
			name:             "FetchAllFlags with request body selector only",
			hasHeader:        false,
			bodySelector:     "source=body-source",
			expectedFlag:     "body-flag",
			expectedSource:   "body-source",
			shouldNotContain: "header-flag",
		},
		{
			name:             "FetchAllFlags header takes precedence over request body",
			hasHeader:        true,
			headerSelector:   "source=header-source",
			bodySelector:     "source=body-source",
			expectedFlag:     "header-flag",
			expectedSource:   "header-source",
			shouldNotContain: "body-flag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flagStore, err := store.NewStore(logger.NewLogger(nil, false), []string{})
			flagStore.Update("header-source", headerFlags, nil)
			flagStore.Update("body-source", bodyFlags, nil)
			require.NoError(t, err)

			handler := syncHandler{
				store:           flagStore,
				log:             logger.NewLogger(nil, false),
				contextValues:   map[string]any{},
				metricsRecorder: &telemetry.NoopMetricsRecorder{},
			}

			// Create context with or without header metadata
			var ctx context.Context
			if tt.hasHeader {
				md := metadata.New(map[string]string{
					flagdService.FLAGD_SELECTOR_HEADER: tt.headerSelector,
				})
				ctx = metadata.NewIncomingContext(context.Background(), md)
			} else {
				ctx = context.Background()
			}

			if strings.Contains(tt.name, "SyncFlags") {
				// Test SyncFlags
				stream := &mockSyncFlagsServer{
					ctx:       ctx,
					mu:        sync.Mutex{},
					respReady: make(chan struct{}, 1),
				}

				go func() {
					err := handler.SyncFlags(&syncv1.SyncFlagsRequest{Selector: tt.bodySelector}, stream)
					assert.NoError(t, err)
				}()

				select {
				case <-stream.respReady:
					assert.Contains(t, stream.lastResp.FlagConfiguration, tt.expectedFlag)
					assert.Contains(t, stream.lastResp.FlagConfiguration, tt.expectedSource)
					assert.NotContains(t, stream.lastResp.FlagConfiguration, tt.shouldNotContain)
				case <-time.After(time.Second):
					t.Fatal("timeout waiting for response")
				}
			} else {
				// Test FetchAllFlags
				resp, err := handler.FetchAllFlags(ctx, &syncv1.FetchAllFlagsRequest{Selector: tt.bodySelector})
				require.NoError(t, err)

				assert.Contains(t, resp.FlagConfiguration, tt.expectedFlag)
				assert.Contains(t, resp.FlagConfiguration, tt.expectedSource)
				assert.NotContains(t, resp.FlagConfiguration, tt.shouldNotContain)
			}
		})
	}
}
