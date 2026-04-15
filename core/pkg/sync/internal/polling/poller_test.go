package polling

import (
	"testing"
	"time"
)

// helper to build a time on a fixed date (2025-01-15) in UTC
func at(h, m, s int) time.Time {
	return time.Date(2025, 1, 15, h, m, s, 0, time.UTC)
}

// helper for times that cross the day boundary (day offset from Jan 15)
func atDay(dayOffset, h, m, s int) time.Time {
	return time.Date(2025, 1, 15+dayOffset, h, m, s, 0, time.UTC)
}

func atNano(h, m, s, ns int) time.Time {
	return time.Date(2025, 1, 15, h, m, s, ns, time.UTC)
}

// pick an arbitrary time and verify that that next poll is as expected
// ie: if interval is 30s, offset is 0, and wall-clock time is 1h:06m:17s, the next poll is 1h:06m:30s
func TestOffsetSchedule_Next(t *testing.T) {
	tests := []struct {
		name     string
		interval uint32
		offset   uint32
		now      time.Time
		expected time.Time
	}{
		// zero offset
		{"zero offset, mid-interval", 30, 0, at(14, 3, 17), at(14, 3, 30)},
		{"zero offset, on boundary", 30, 0, at(14, 3, 30), at(14, 4, 0)},
		{"zero offset, on :00 boundary", 30, 0, at(14, 3, 0), at(14, 3, 30)},

		// with offset
		{"offset=7, mid-interval", 30, 7, at(14, 3, 17), at(14, 3, 37)},
		{"offset=7, on boundary", 30, 7, at(14, 3, 37), at(14, 4, 7)},
		{"offset=7, just before boundary", 30, 7, at(14, 3, 5), at(14, 3, 7)},

		// small interval
		{"interval=5, mid-interval", 5, 0, at(14, 3, 17), at(14, 3, 20)},

		// sub-second truncation
		{"sub-second truncated", 30, 0, atNano(14, 3, 17, 500000000), at(14, 3, 30)},

		// large offset relative to interval
		{"offset=29, interval=30", 30, 29, at(14, 3, 58), at(14, 3, 59)},
		{"offset=29, on boundary", 30, 29, at(14, 3, 59), at(14, 4, 29)},

		// large interval (2 minutes)
		{"interval=120, offset=10", 120, 10, at(14, 3, 17), at(14, 4, 10)},

		// boundary crossings
		{"hour boundary", 30, 0, at(14, 59, 45), at(15, 0, 0)},
		{"day boundary", 30, 0, at(23, 59, 45), atDay(1, 0, 0, 0)},

		// zero interval falls back to 1s
		{"zero interval", 0, 0, at(14, 3, 17), at(14, 3, 18)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := OffsetSchedule{Interval: tt.interval, Offset: tt.offset}
			got := s.Next(tt.now)
			if !got.Equal(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}

	// large interval: verify int64 math prevents uint32 overflow
	t.Run("max uint32 interval does not overflow", func(t *testing.T) {
		now := at(0, 0, 1)
		s := OffsetSchedule{Interval: 1<<32 - 1, Offset: 0}
		got := s.Next(now)
		if !got.After(now) {
			t.Errorf("expected time after %v, got %v", now, got)
		}
	})
}

func TestPollOffset(t *testing.T) {
	tests := []struct {
		name     string
		seed     string
		interval uint32
		wantZero bool
	}{
		{"empty seed", "", 30, true},
		{"zero interval", "test", 0, true},
		{"pod-a", "pod-a", 30, false},
		{"pod-b", "pod-b", 30, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offset := pollOffset(tt.seed, tt.interval)
			if tt.wantZero && offset != 0 {
				t.Errorf("expected 0, got %d", offset)
			}
			if !tt.wantZero && offset >= tt.interval {
				t.Errorf("offset %d should be less than interval %d", offset, tt.interval)
			}
		})
	}
}

func TestPollOffset_Deterministic(t *testing.T) {
	a := pollOffset("my-pod", 30)
	b := pollOffset("my-pod", 30)
	if a != b {
		t.Errorf("expected deterministic offset, got %d and %d", a, b)
	}
}

func TestPollOffset_DifferentSeeds(t *testing.T) {
	a := pollOffset("pod-alpha", 60)
	b := pollOffset("pod-beta", 60)
	if a == b {
		t.Logf("warning: different seeds produced the same offset %d (unlikely but possible)", a)
	}
}

func TestCronPoller_Offset(t *testing.T) {
	tests := []struct {
		name     string
		interval uint32
		seed     string
		wantZero bool
	}{
		{"no seed returns 0", 30, "", true},
		{"zero interval returns 0", 0, "my-pod", true},
		{"with seed returns non-zero", 30, "my-pod", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewCronPoller(tt.interval, tt.seed)
			got := p.Offset()
			if tt.wantZero && got != 0 {
				t.Errorf("expected offset 0, got %d", got)
			}
			if !tt.wantZero && got >= tt.interval {
				t.Errorf("offset %d should be less than interval %d", got, tt.interval)
			}
		})
	}
}

func TestCronPoller_Offset_MatchesPollOffset(t *testing.T) {
	// Offset() should return the same value as pollOffset() for the same inputs
	p := NewCronPoller(60, "my-pod")
	expected := pollOffset("my-pod", 60)
	if p.Offset() != expected {
		t.Errorf("expected Offset() = %d, got %d", expected, p.Offset())
	}
}

func TestCronPoller_Offset_Deterministic(t *testing.T) {
	a := NewCronPoller(30, "my-pod").Offset()
	b := NewCronPoller(30, "my-pod").Offset()
	if a != b {
		t.Errorf("expected deterministic offset, got %d and %d", a, b)
	}
}
