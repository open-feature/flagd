package sync

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
)

// Multiplexer abstract subscription handling and storage processing.
// Flag configurations will be lazy loaded using reFill logic upon the calls to publish.
type Multiplexer struct {
	store   *store.Flags
	sources []string

	subs         map[interface{}]subscription            // subscriptions on all sources
	selectorSubs map[string]map[interface{}]subscription // source specific subscriptions

	allFlags      string            // pre-calculated all flags in store as a string
	selectorFlags map[string]string // pre-calculated selector scoped flags in store as strings

	mu sync.RWMutex
}

type subscription struct {
	id      interface{}
	channel chan payload
}

type payload struct {
	flags string
}

// NewMux creates a new sync multiplexer
func NewMux(store *store.Flags, sources []string) (*Multiplexer, error) {
	m := &Multiplexer{
		store:         store,
		sources:       sources,
		subs:          map[interface{}]subscription{},
		selectorSubs:  map[string]map[interface{}]subscription{},
		selectorFlags: map[string]string{},
	}

	return m, m.reFill()
}

// Register a subscription
func (r *Multiplexer) Register(id interface{}, source string, con chan payload) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if source != "" && !slices.Contains(r.sources, source) {
		return fmt.Errorf("no flag watcher setup for source %s", source)
	}

	var initSync string

	if source == "" {
		// subscribe for flags from all source
		r.subs[id] = subscription{
			id:      id,
			channel: con,
		}

		initSync = r.allFlags
	} else {
		// subscribe for specific source
		s, ok := r.selectorSubs[source]
		if ok {
			s[id] = subscription{
				id:      id,
				channel: con,
			}
		} else {
			r.selectorSubs[source] = map[interface{}]subscription{
				id: {
					id:      id,
					channel: con,
				},
			}
		}

		initSync = r.selectorFlags[source]
	}

	// Initial sync
	con <- payload{flags: initSync}
	return nil
}

// Publish sync updates to subscriptions
func (r *Multiplexer) Publish() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// perform a refill prior to publishing
	err := r.reFill()
	if err != nil {
		return err
	}

	// push to all source subs
	for _, sub := range r.subs {
		sub.channel <- payload{r.allFlags}
	}

	// push to selector subs
	for source, flags := range r.selectorFlags {
		for _, s := range r.selectorSubs[source] {
			s.channel <- payload{flags}
		}
	}

	return nil
}

// Unregister a subscription
func (r *Multiplexer) Unregister(id interface{}, selector string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var from map[interface{}]subscription

	if selector == "" {
		from = r.subs
	} else {
		from = r.selectorSubs[selector]
	}

	delete(from, id)
}

// GetAllFlags per specific source
func (r *Multiplexer) GetAllFlags(source string) (string, error) {
	r.mu.RLocker()
	defer r.mu.RUnlock()

	if source == "" {
		return r.allFlags, nil
	}

	if !slices.Contains(r.sources, source) {
		return "", fmt.Errorf("no flag watcher setup for source %s", source)
	}

	return r.selectorFlags[source], nil
}

// SourcesAsMetadata returns all known sources, comma separated to be used as service metadata
func (r *Multiplexer) SourcesAsMetadata() string {
	r.mu.RLocker()
	defer r.mu.RUnlock()

	return strings.Join(r.sources, ",")
}

// reFill local configuration values
func (r *Multiplexer) reFill() error {
	clear(r.selectorFlags)

	all := r.store.GetAll()
	bytes, err := json.Marshal(map[string]interface{}{"flags": all})
	if err != nil {
		return fmt.Errorf("error from marshallin: %w", err)
	}

	r.allFlags = string(bytes)

	collector := map[string]map[string]model.Flag{}

	for key, flag := range all {
		c, ok := collector[flag.Source]
		if ok {
			c[key] = flag
		} else {
			collector[flag.Source] = map[string]model.Flag{
				key: flag,
			}
		}
	}

	for source, flags := range collector {
		bytes, err := json.Marshal(map[string]interface{}{"flags": flags})
		if err != nil {
			return fmt.Errorf("unable to marshal flags: %w", err)
		}

		r.selectorFlags[source] = string(bytes)
	}

	return nil
}
