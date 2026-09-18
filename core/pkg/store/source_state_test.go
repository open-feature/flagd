package store

import "testing"

func TestSourceState_DefaultsToNotStale(t *testing.T) {
	s := NewSourceState()

	if s.IsStale("grpc://example:8015") {
		t.Fatal("a source with no recorded state must not be reported stale")
	}
}

func TestSourceState_SetAndClear(t *testing.T) {
	const source = "grpc://example:8015"
	s := NewSourceState()

	s.SetStale(source, true)
	if !s.IsStale(source) {
		t.Fatal("expected source to be stale after SetStale(true)")
	}

	s.SetStale(source, false)
	if s.IsStale(source) {
		t.Fatal("expected source to be fresh after SetStale(false)")
	}
	if got := len(s.StaleSources()); got != 0 {
		t.Fatalf("expected no stale sources retained, got %d", got)
	}
}

func TestSourceState_IsolatesSources(t *testing.T) {
	s := NewSourceState()
	s.SetStale("a", true)

	if !s.IsStale("a") {
		t.Fatal("expected source a to be stale")
	}
	if s.IsStale("b") {
		t.Fatal("marking source a stale must not affect source b")
	}
}

// A nil SourceState is the zero-configuration case: embedders and tests that
// never wire up source tracking must see exactly the pre-existing behaviour
// rather than a panic.
func TestSourceState_NilReceiverIsSafe(t *testing.T) {
	var s *SourceState

	s.SetStale("a", true)
	if s.IsStale("a") {
		t.Fatal("a nil SourceState must never report a source stale")
	}
	if s.StaleSources() != nil {
		t.Fatal("a nil SourceState must return no stale sources")
	}
}

func TestSourceState_ConcurrentAccess(t *testing.T) {
	s := NewSourceState()
	done := make(chan struct{})

	go func() {
		for i := 0; i < 1000; i++ {
			s.SetStale("a", i%2 == 0)
		}
		close(done)
	}()
	for i := 0; i < 1000; i++ {
		_ = s.IsStale("a")
	}
	<-done
}
