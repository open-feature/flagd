package sync

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	syncv1 "buf.build/gen/go/open-feature/flagd/protocolbuffers/go/flagd/sync/v1"
	"connectrpc.com/connect"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/telemetry"
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
				metaResp, err := handler.GetMetadata(context.Background(), connect.NewRequest(&syncv1.GetMetadataRequest{}))
				if !disableSyncMetadata {
					require.NoError(t, err)
					respMetadata := metaResp.Msg.GetMetadata().AsMap()
					assert.Equal(t, tt.wantMetadata, respMetadata)
				} else {
					assert.NotNil(t, err)
				}

				// Test metadata from sync_context
				stream := newMockSyncFlagsStream()

				go func() {
					err := handler.syncFlags(context.Background(), nil, "", "", stream)
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

// mockSyncFlagsStream stands in for connect.ServerStream, which cannot be constructed in a test.
type mockSyncFlagsStream struct {
	mu        sync.Mutex
	lastResp  *syncv1.SyncFlagsResponse
	respReady chan struct{}
}

func newMockSyncFlagsStream() *mockSyncFlagsStream {
	return &mockSyncFlagsStream{respReady: make(chan struct{}, 1)}
}

func (m *mockSyncFlagsStream) Send(resp *syncv1.SyncFlagsResponse) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lastResp = resp
	select {
	case m.respReady <- struct{}{}:
	default:
	}
	return nil
}

func (m *mockSyncFlagsStream) GetLastResponse() *syncv1.SyncFlagsResponse {
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
			flagStore.Update("header-source", headerFlags, nil, false)
			flagStore.Update("body-source", bodyFlags, nil, false)
			require.NoError(t, err)

			handler := syncHandler{
				store:           flagStore,
				log:             logger.NewLogger(nil, false),
				contextValues:   map[string]any{},
				metricsRecorder: &telemetry.NoopMetricsRecorder{},
			}

			var header http.Header
			if tt.hasHeader {
				header = selectorHeader(tt.headerSelector)
			}

			if strings.Contains(tt.name, "SyncFlags") {
				// Test SyncFlags
				stream := newMockSyncFlagsStream()

				go func() {
					err := handler.syncFlags(context.Background(), header, "", tt.bodySelector, stream)
					assert.NoError(t, err)
				}()

				select {
				case <-stream.respReady:
					config := stream.GetLastResponse().GetFlagConfiguration()
					assert.Contains(t, config, tt.expectedFlag)
					assert.Contains(t, config, tt.expectedSource)
					assert.NotContains(t, config, tt.shouldNotContain)
				case <-time.After(time.Second):
					t.Fatal("timeout waiting for response")
				}
			} else {
				// Test FetchAllFlags
				req := connect.NewRequest(&syncv1.FetchAllFlagsRequest{Selector: tt.bodySelector})
				for key, values := range header {
					req.Header()[key] = values
				}

				resp, err := handler.FetchAllFlags(context.Background(), req)
				require.NoError(t, err)

				assert.Contains(t, resp.Msg.GetFlagConfiguration(), tt.expectedFlag)
				assert.Contains(t, resp.Msg.GetFlagConfiguration(), tt.expectedSource)
				assert.NotContains(t, resp.Msg.GetFlagConfiguration(), tt.shouldNotContain)
			}
		})
	}
}

func TestSyncHandler_InvalidSelector(t *testing.T) {
	const invalidSelector = "invalidKey=val"
	const wantMessage = `invalid selector key "invalidKey", valid keys: "flagSetId", "source"`

	flagStore, err := store.NewStore(logger.NewLogger(nil, false), []string{})
	require.NoError(t, err)

	h := syncHandler{
		store:           flagStore,
		log:             logger.NewLogger(nil, false),
		contextValues:   map[string]any{},
		metricsRecorder: &telemetry.NoopMetricsRecorder{},
	}

	header := selectorHeader(invalidSelector)

	t.Run("SyncFlags", func(t *testing.T) {
		err := h.syncFlags(context.Background(), header, "", "", newMockSyncFlagsStream())
		requireConnectError(t, err, connect.CodeInvalidArgument, wantMessage)
	})

	t.Run("FetchAllFlags", func(t *testing.T) {
		req := connect.NewRequest(&syncv1.FetchAllFlagsRequest{})
		req.Header().Set(selectorHeaderKey, invalidSelector)

		_, err := h.FetchAllFlags(context.Background(), req)
		requireConnectError(t, err, connect.CodeInvalidArgument, wantMessage)
	})
}

func requireConnectError(t *testing.T, err error, wantCode connect.Code, wantMessage string) {
	t.Helper()

	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	assert.Equal(t, wantCode, connectErr.Code())
	assert.Equal(t, wantMessage, connectErr.Message())
}
