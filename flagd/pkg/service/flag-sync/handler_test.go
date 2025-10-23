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
	flagdService "github.com/open-feature/flagd/flagd/pkg/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
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

// TestSyncHandler_SelectorFromHeader tests that the selector is correctly extracted from the header
func TestSyncHandler_SelectorFromHeader(t *testing.T) {
	flagStore, err := store.NewStore(logger.NewLogger(nil, false), []string{})
	require.NoError(t, err)

	// Create a logger with observer to capture log messages
	observedZapCore, observedLogs := observer.New(zapcore.WarnLevel)
	observedLogger := zap.New(observedZapCore)
	log := logger.NewLogger(observedLogger, false)

	handler := syncHandler{
		store:         flagStore,
		log:           log,
		contextValues: map[string]any{},
	}

	// Create context with metadata containing the selector header
	md := metadata.New(map[string]string{
		flagdService.FLAGD_SELECTOR_HEADER: "source:my-source",
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	// Test with SyncFlags
	stream := &mockSyncFlagsServer{
		ctx:       ctx,
		mu:        sync.Mutex{},
		respReady: make(chan struct{}, 1),
	}

	go func() {
		// Use empty request body selector to verify header is used
		err := handler.SyncFlags(&syncv1.SyncFlagsRequest{Selector: ""}, stream)
		assert.NoError(t, err)
	}()

	select {
	case <-stream.respReady:
		// Verify no deprecation warning was logged
		logs := observedLogs.All()
		for _, log := range logs {
			assert.NotContains(t, log.Message, "deprecated", "Should not log deprecation warning when using header")
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for response")
	}
}

// TestSyncHandler_SelectorFromRequestBody tests backward compatibility with request body selector
func TestSyncHandler_SelectorFromRequestBody(t *testing.T) {
	flagStore, err := store.NewStore(logger.NewLogger(nil, false), []string{})
	require.NoError(t, err)

	// Create a logger with observer to capture log messages
	observedZapCore, observedLogs := observer.New(zapcore.WarnLevel)
	observedLogger := zap.New(observedZapCore)
	log := logger.NewLogger(observedLogger, false)

	handler := syncHandler{
		store:         flagStore,
		log:           log,
		contextValues: map[string]any{},
	}

	// Create context without metadata (no header)
	ctx := context.Background()

	// Test with SyncFlags
	stream := &mockSyncFlagsServer{
		ctx:       ctx,
		mu:        sync.Mutex{},
		respReady: make(chan struct{}, 1),
	}

	go func() {
		// Use request body selector
		err := handler.SyncFlags(&syncv1.SyncFlagsRequest{Selector: "source:legacy-source"}, stream)
		assert.NoError(t, err)
	}()

	select {
	case <-stream.respReady:
		// Verify deprecation warning was logged
		logs := observedLogs.All()
		require.Greater(t, len(logs), 0, "Expected at least one log entry")
		found := false
		for _, log := range logs {
			if log.Level == zapcore.WarnLevel {
				assert.Contains(t, log.Message, "deprecated", "Should log deprecation warning when using request body selector")
				assert.Contains(t, log.Message, "Flagd-Selector", "Deprecation message should mention the header name")
				found = true
				break
			}
		}
		assert.True(t, found, "Expected to find deprecation warning in logs")
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for response")
	}
}

// TestSyncHandler_SelectorHeaderTakesPrecedence tests that header takes precedence over request body
func TestSyncHandler_SelectorHeaderTakesPrecedence(t *testing.T) {
	flagStore, err := store.NewStore(logger.NewLogger(nil, false), []string{})
	require.NoError(t, err)

	// Create a logger with observer to capture log messages
	observedZapCore, observedLogs := observer.New(zapcore.WarnLevel)
	observedLogger := zap.New(observedZapCore)
	log := logger.NewLogger(observedLogger, false)

	handler := syncHandler{
		store:         flagStore,
		log:           log,
		contextValues: map[string]any{},
	}

	// Create context with metadata containing the selector header
	md := metadata.New(map[string]string{
		flagdService.FLAGD_SELECTOR_HEADER: "source:header-source",
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	// Test with SyncFlags
	stream := &mockSyncFlagsServer{
		ctx:       ctx,
		mu:        sync.Mutex{},
		respReady: make(chan struct{}, 1),
	}

	go func() {
		// Provide both header and request body selector
		err := handler.SyncFlags(&syncv1.SyncFlagsRequest{Selector: "source:body-source"}, stream)
		assert.NoError(t, err)
	}()

	select {
	case <-stream.respReady:
		// Verify no deprecation warning was logged (header was used)
		logs := observedLogs.All()
		for _, log := range logs {
			assert.NotContains(t, log.Message, "deprecated", "Should not log deprecation warning when header is present")
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for response")
	}
}

// TestSyncHandler_FetchAllFlags_SelectorFromHeader tests FetchAllFlags with header selector
func TestSyncHandler_FetchAllFlags_SelectorFromHeader(t *testing.T) {
	flagStore, err := store.NewStore(logger.NewLogger(nil, false), []string{})
	require.NoError(t, err)

	// Create a logger with observer to capture log messages
	observedZapCore, observedLogs := observer.New(zapcore.WarnLevel)
	observedLogger := zap.New(observedZapCore)
	log := logger.NewLogger(observedLogger, false)

	handler := syncHandler{
		store:         flagStore,
		log:           log,
		contextValues: map[string]any{},
	}

	// Create context with metadata containing the selector header
	md := metadata.New(map[string]string{
		flagdService.FLAGD_SELECTOR_HEADER: "source:my-source",
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	// Call FetchAllFlags with empty request body selector
	_, err = handler.FetchAllFlags(ctx, &syncv1.FetchAllFlagsRequest{Selector: ""})
	require.NoError(t, err)

	// Verify no deprecation warning was logged
	logs := observedLogs.All()
	for _, log := range logs {
		assert.NotContains(t, log.Message, "deprecated", "Should not log deprecation warning when using header")
	}
}

// TestSyncHandler_FetchAllFlags_SelectorFromRequestBody tests FetchAllFlags with request body selector
func TestSyncHandler_FetchAllFlags_SelectorFromRequestBody(t *testing.T) {
	flagStore, err := store.NewStore(logger.NewLogger(nil, false), []string{})
	require.NoError(t, err)

	// Create a logger with observer to capture log messages
	observedZapCore, observedLogs := observer.New(zapcore.WarnLevel)
	observedLogger := zap.New(observedZapCore)
	log := logger.NewLogger(observedLogger, false)

	handler := syncHandler{
		store:         flagStore,
		log:           log,
		contextValues: map[string]any{},
	}

	// Create context without metadata (no header)
	ctx := context.Background()

	// Call FetchAllFlags with request body selector
	_, err = handler.FetchAllFlags(ctx, &syncv1.FetchAllFlagsRequest{Selector: "source:legacy-source"})
	require.NoError(t, err)

	// Verify deprecation warning was logged
	logs := observedLogs.All()
	require.Greater(t, len(logs), 0, "Expected at least one log entry")
	found := false
	for _, log := range logs {
		if log.Level == zapcore.WarnLevel {
			assert.Contains(t, log.Message, "deprecated", "Should log deprecation warning when using request body selector")
			assert.Contains(t, log.Message, "Flagd-Selector", "Deprecation message should mention the header name")
			found = true
			break
		}
	}
	assert.True(t, found, "Expected to find deprecation warning in logs")
}
