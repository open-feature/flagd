package sse

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/launchdarkly/eventsource"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
)

// allChannel receives a refetch event on any flag change. It is advertised to OFREP clients
// that use no flagSetId selector, mirroring a bulk request with no selector.
const allChannel = ""

// namespacing prefixes for the internal version map, so a flagSetId and a source with the
// same string value never collide.
const (
	allKey    = "all"
	fsPrefix  = "fs:"
	srcPrefix = "src:"
)

func fsKey(id string) string   { return fsPrefix + id }
func srcKey(src string) string { return srcPrefix + src }

type version struct {
	etag         string
	lastModified int64
}

// Tracker owns a global store.Watch subscription. On every flag-configuration change it
// recomputes a per-channel config fingerprint (used as an ETag), publishes an ADR-0008
// refetchEvaluation event to each affected channel, and answers version lookups for the OFREP
// bulk handler so conditional (ETag/304) evaluation stays consistent with the SSE stream.
type Tracker struct {
	logger *logger.Logger
	store  store.IStore
	es     *eventsource.Server

	mu       sync.RWMutex
	versions map[string]version

	eventID atomic.Int64
}

// NewTracker creates a Tracker. Call Run to begin watching the store.
func NewTracker(log *logger.Logger, s store.IStore, es *eventsource.Server) *Tracker {
	return &Tracker{
		logger:   log,
		store:    s,
		es:       es,
		versions: map[string]version{},
	}
}

// Run subscribes to all flag changes and publishes refetch events until ctx is cancelled.
// It blocks, so it is intended to run in its own goroutine.
func (t *Tracker) Run(ctx context.Context) {
	watcher := make(chan store.FlagQueryResult, 1)
	t.store.Watch(ctx, &store.Selector{}, watcher)

	first := true
	for res := range watcher {
		channels := t.update(res.Flags)
		if first {
			// The first emission is the initialization snapshot; only seed fingerprints.
			first = false
			continue
		}
		for _, ch := range channels {
			t.publish(ch)
		}
	}
}

// Version returns the current config ETag and last-modified time (unix seconds) for the
// channel matching the given selector. ok is false when no version is tracked for the
// selector (e.g. an unknown source), letting the caller skip conditional handling.
func (t *Tracker) Version(selector store.Selector) (etag string, lastModified int64, ok bool) {
	var key string
	switch {
	case selector.IsEmpty():
		key = allKey
	case selector.FlagSetId() != "":
		key = fsKey(selector.FlagSetId())
	case selector.Source() != "":
		key = srcKey(selector.Source())
	default:
		return "", 0, false
	}

	t.mu.RLock()
	defer t.mu.RUnlock()
	v, ok := t.versions[key]
	return v.etag, v.lastModified, ok
}

// update recomputes fingerprints for the catch-all, per-flagSetId and per-source groups,
// swaps them into the version map (preserving lastModified when a fingerprint is unchanged)
// and returns the eventsource channels that should be notified.
func (t *Tracker) update(flags []model.Flag) []string {
	now := time.Now().Unix()

	fsGroups := map[string][]model.Flag{}
	srcGroups := map[string][]model.Flag{}
	for _, f := range flags {
		fsGroups[f.FlagSetId] = append(fsGroups[f.FlagSetId], f)
		srcGroups[f.Source] = append(srcGroups[f.Source], f)
	}

	newVersions := make(map[string]version, len(fsGroups)+len(srcGroups)+1)
	newVersions[allKey] = version{etag: fingerprint(flags)}
	for id, g := range fsGroups {
		newVersions[fsKey(id)] = version{etag: fingerprint(g)}
	}
	for src, g := range srcGroups {
		newVersions[srcKey(src)] = version{etag: fingerprint(g)}
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	// carry lastModified forward when the fingerprint is unchanged
	for key, nv := range newVersions {
		if old, exists := t.versions[key]; exists && old.etag == nv.etag {
			nv.lastModified = old.lastModified
		} else {
			nv.lastModified = now
		}
		newVersions[key] = nv
	}

	changed := t.changedChannels(newVersions, fsGroups)
	t.versions = newVersions
	return changed
}

// changedChannels compares the previous version map (t.versions, still held) with the freshly
// computed one and returns the eventsource channels that clients subscribe to and whose config
// changed: the catch-all channel plus any created/updated/removed flagSetId channels. Source
// channels are tracked for ETag lookups but are not (yet) directly subscribable.
func (t *Tracker) changedChannels(newVersions map[string]version, fsGroups map[string][]model.Flag) []string {
	etagChanged := func(key string) bool {
		old, oldOK := t.versions[key]
		nv, newOK := newVersions[key]
		if oldOK != newOK {
			return true
		}
		return old.etag != nv.etag
	}

	notify := map[string]struct{}{}
	if etagChanged(allKey) {
		notify[allChannel] = struct{}{}
	}

	nilID := store.NilFlagSetId()
	for id := range fsGroups {
		if id == nilID {
			continue // internal flagSetId, never subscribable
		}
		if etagChanged(fsKey(id)) {
			notify[id] = struct{}{}
		}
	}
	// flagSetId channels that existed before but have no flags now (whole set deleted)
	for key := range t.versions {
		id, ok := strings.CutPrefix(key, fsPrefix)
		if !ok || id == nilID {
			continue
		}
		if _, present := fsGroups[id]; !present {
			notify[id] = struct{}{}
		}
	}

	channels := make([]string, 0, len(notify))
	for ch := range notify {
		channels = append(channels, ch)
	}
	return channels
}

// publish emits a refetch event to a single eventsource channel using that channel's current
// ETag and lastModified.
func (t *Tracker) publish(channel string) {
	key := allKey
	if channel != allChannel {
		key = fsKey(channel)
	}

	t.mu.RLock()
	v := t.versions[key]
	t.mu.RUnlock()

	id := strconv.FormatInt(t.eventID.Add(1), 10)
	t.es.Publish([]string{channel}, newRefetchEvent(id, v.etag, v.lastModified))
	if t.logger != nil {
		t.logger.Debug(fmt.Sprintf("published refetch event to channel %q (etag=%s)", channel, v.etag))
	}
}

// fingerprint produces a deterministic, restart-stable hash of a group of flag definitions.
// nilFlagSetId is normalised to "" so identical config yields the same fingerprint across
// restarts.
func fingerprint(flags []model.Flag) string {
	sorted := make([]model.Flag, len(flags))
	copy(sorted, flags)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].FlagSetId != sorted[j].FlagSetId {
			return sorted[i].FlagSetId < sorted[j].FlagSetId
		}
		if sorted[i].Source != sorted[j].Source {
			return sorted[i].Source < sorted[j].Source
		}
		return sorted[i].Key < sorted[j].Key
	})

	nilID := store.NilFlagSetId()
	h := sha256.New()
	for _, f := range sorted {
		fsid := f.FlagSetId
		if fsid == nilID {
			fsid = ""
		}
		h.Write([]byte(fsid))
		h.Write([]byte{0})
		h.Write([]byte(f.Key))
		h.Write([]byte{0})
		// model.Flag.MarshalJSON is stable (map keys sorted) and covers the definition fields.
		b, _ := json.Marshal(f)
		h.Write(b)
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))
}
