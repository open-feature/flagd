package sse

import (
	"net/http"
	"sync"
)

// channelParam carries the ADR-0008 channel token (the flagSetId); empty selects the catch-all.
const channelParam = "channels"

func channelFromRequest(r *http.Request) string {
	return r.URL.Query().Get(channelParam)
}

// activeChannels is a reference-counted set of channels with at least one live subscriber, so
// the heartbeat loop knows which channels to ping (the eventsource server exposes no registry).
type activeChannels struct {
	mu     sync.Mutex
	counts map[string]int
}

func newActiveChannels() *activeChannels {
	return &activeChannels{counts: map[string]int{}}
}

func (a *activeChannels) add(channel string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.counts[channel]++
}

func (a *activeChannels) remove(channel string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.counts[channel] <= 1 {
		delete(a.counts, channel)
		return
	}
	a.counts[channel]--
}

func (a *activeChannels) snapshot() []string {
	a.mu.Lock()
	defer a.mu.Unlock()
	channels := make([]string, 0, len(a.counts))
	for ch := range a.counts {
		channels = append(channels, ch)
	}
	return channels
}
