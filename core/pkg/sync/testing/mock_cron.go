package testing

import (
	"context"
	"sync"
)

// MockPoller is a mock of the polling.Poller interface for testing.
// It captures the callback so tests can trigger it manually via Tick().
// All fields are synchronized via a ready channel and mutex so that
// Tick() safely blocks until Start() has been called.
type MockPoller struct {
	mu       sync.Mutex
	callback func()
	ready    chan struct{}
}

// NewMockPoller creates a new MockPoller.
func NewMockPoller() *MockPoller {
	return &MockPoller{
		ready: make(chan struct{}),
	}
}

// Start captures the callback without blocking (unlike the real CronPoller).
func (m *MockPoller) Start(_ context.Context, callback func()) {
	m.mu.Lock()
	m.callback = callback
	m.mu.Unlock()
	close(m.ready)
}

// Tick blocks until Start has been called, then invokes the captured callback.
func (m *MockPoller) Tick() {
	<-m.ready
	m.mu.Lock()
	cb := m.callback
	m.mu.Unlock()
	if cb != nil {
		cb()
	}
}

// Started returns whether Start was called (non-blocking).
func (m *MockPoller) Started() bool {
	select {
	case <-m.ready:
		return true
	default:
		return false
	}
}

// Offset returns 0 (no offset in tests by default).
func (m *MockPoller) Offset() uint32 {
	return 0
}
