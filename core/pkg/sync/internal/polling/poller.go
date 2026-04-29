package polling

import (
	"context"
	"fmt"
	"hash/fnv"
	"time"

	"github.com/robfig/cron/v3"
)

// MaxInterval is the largest supported polling interval in seconds (1 day).
// OffsetSchedule uses seconds-since-midnight math, so intervals beyond a day
// would wrap around and produce incorrect fire times.
const MaxInterval uint32 = 86400

// OffsetSchedule is a cron.Schedule that fires every `interval` seconds,
// aligned to wall-clock time but shifted by `offset` seconds.
// For example, with interval=30 and offset=7, it fires at :07 and :37 of
// every minute. With offset=0, it is equivalent to the cron expression
// "*/interval * * * *" (using the 6-field format where the first field is seconds).
//
// This allows multiple instances with different offsets (derived from a seed)
// to avoid polling at the same instant (thundering herd), while remaining
// deterministic across restarts.
type OffsetSchedule struct {
	Interval uint32
	Offset   uint32
}

// Next returns the next activation time after t.
func (s OffsetSchedule) Next(t time.Time) time.Time {
	if s.Interval == 0 {
		return t.Add(time.Second)
	}

	// seconds since midnight in the local timezone
	hour, min, sec := t.Clock()
	now := int64(hour*3600 + min*60 + sec)
	interval := int64(s.Interval)
	offset := int64(s.Offset)

	// seconds since the last fire
	sinceLastFire := (now - offset%interval + interval) % interval
	lastFire := now - sinceLastFire

	// the next fire time
	nextFire := lastFire + interval

	delta := nextFire - now
	// truncate to the start of the current second, then add delta seconds
	return t.Truncate(time.Second).Add(time.Duration(delta) * time.Second)
}

// pollOffset computes a deterministic offset from a seed string.
// Returns 0 if seed is empty.
// fnv32a is a fast, well-distributed non-cryptographic hash;
// we only need even distribution across [0, interval), not collision resistance.
func pollOffset(seed string, interval uint32) uint32 {
	if seed == "" || interval == 0 {
		return 0
	}
	h := fnv.New32a()
	h.Write([]byte(seed))
	return h.Sum32() % interval
}

// Poller schedules a recurring callback. Start blocks until ctx is cancelled.
type Poller interface {
	Start(ctx context.Context, callback func())
	// Offset returns the schedule offset in seconds (0 when no seed is configured).
	Offset() uint32
}

// CronPoller is a Poller backed by a wall-clock-aligned cron schedule with a
// deterministic offset derived from a seed.
type CronPoller struct {
	cr       *cron.Cron
	interval uint32
	offset   uint32
}

// NewCronPoller creates a CronPoller. If intervalSeed is empty, offset defaults to
// 0 (equivalent to the legacy wall-clock-aligned behavior).
// Returns an error if interval exceeds MaxInterval.
func NewCronPoller(interval uint32, intervalSeed string) (*CronPoller, error) {
	if interval > MaxInterval {
		return nil, fmt.Errorf("polling interval %ds exceeds maximum %ds", interval, MaxInterval)
	}
	return &CronPoller{
		cr:       cron.New(),
		interval: interval,
		offset:   pollOffset(intervalSeed, interval),
	}, nil
}

// Offset returns the computed schedule offset in seconds.
func (p *CronPoller) Offset() uint32 {
	return p.offset
}

// Start schedules the callback and blocks until ctx is cancelled.
func (p *CronPoller) Start(ctx context.Context, callback func()) {
	schedule := OffsetSchedule{
		Interval: p.interval,
		Offset:   p.offset,
	}
	p.cr.Schedule(schedule, cron.FuncJob(callback))
	p.cr.Start()

	<-ctx.Done()
	<-p.cr.Stop().Done()
}
