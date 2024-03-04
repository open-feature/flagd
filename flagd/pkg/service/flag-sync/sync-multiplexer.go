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

type syncMultiplexer struct {
	store   *store.Flags
	sources []string

	subs         map[interface{}]subscription            // subscriptions on all sources
	selectorSubs map[string]map[interface{}]subscription // source specific subscriptions

	allFlags      string            // pre-calculated all flags in store as a string
	selectorFlags map[string]string // pre-calculated selector scoped flags

	mu sync.Mutex
}

type subscription struct {
	id      interface{}
	channel chan payload
}

type payload struct {
	flags string
}

func newMux(store *store.Flags, sources []string) *syncMultiplexer {
	return &syncMultiplexer{
		store:         store,
		sources:       sources,
		subs:          map[interface{}]subscription{},
		selectorSubs:  map[string]map[interface{}]subscription{},
		selectorFlags: map[string]string{},
	}
}

func (r *syncMultiplexer) register(id interface{}, source string, con chan payload) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if source != "" && !slices.Contains(r.sources, source) {
		return fmt.Errorf("no flag watcher setup for source %s", source)
	}

	var initSync string
	var err error

	if source == "" {
		// subscribe for flags from all source
		r.subs[id] = subscription{
			id:      id,
			channel: con,
		}

		initSync, err = r.store.String()
		if err != nil {
			return fmt.Errorf("errpr getting all flags: %w", err)
		}
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

func (r *syncMultiplexer) pushUpdates() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	err := r.extract()
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

func (r *syncMultiplexer) unregister(id interface{}, selector string) {
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

func (r *syncMultiplexer) getALlFlags(source string) (string, error) {
	if source != "" && !slices.Contains(r.sources, source) {
		return "", fmt.Errorf("no flag watcher setup for source %s", source)
	}

	if source == "" {
		return r.allFlags, nil
	}

	return r.selectorFlags[source], nil
}

func (r *syncMultiplexer) sourcesAsMetadata() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	return strings.Join(r.store.FlagSources, ",")
}

func (r *syncMultiplexer) extract() error {
	clear(r.selectorFlags)

	all := r.store.GetAll()
	bytes, err := json.Marshal(all)
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
		bytes, err := json.Marshal(flags)
		if err != nil {
			return fmt.Errorf("unable to marshal flags: %w", err)
		}

		r.selectorFlags[source] = string(bytes)
	}

	return nil
}
