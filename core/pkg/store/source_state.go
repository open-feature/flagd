package store

import "sync"

// SourceState tracks, per sync source, whether that source is currently
// disconnected. Flags already held in the store remain servable while a source
// is down -- flagd deliberately keeps serving last-known-good data rather than
// failing evaluations -- but consumers deserve to know the data may be out of
// date. Evaluations resolved from a disconnected source are reported with
// model.StaleReason.
//
// It is written by the runtime (driven by sync.DataSync payloads) and read on
// the evaluation hot path, so reads are cheap and lock-free-ish via RWMutex.
// The zero value is not usable; construct with NewSourceState.
type SourceState struct {
	mu    sync.RWMutex
	stale map[string]bool
}

func NewSourceState() *SourceState {
	return &SourceState{stale: map[string]bool{}}
}

// SetStale records whether the given source is currently disconnected.
func (s *SourceState) SetStale(source string, stale bool) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if stale {
		s.stale[source] = true
		return
	}
	delete(s.stale, source)
}

// IsStale reports whether the given source is currently disconnected. A nil
// receiver reports false so that callers which never wire up source tracking
// (tests, embedders) behave exactly as before.
func (s *SourceState) IsStale(source string) bool {
	if s == nil {
		return false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.stale[source]
}

// StaleSources returns the sources currently marked disconnected. Intended for
// diagnostics and tests; order is not guaranteed.
func (s *SourceState) StaleSources() []string {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	sources := make([]string, 0, len(s.stale))
	for source := range s.stale {
		sources = append(sources, source)
	}
	return sources
}
