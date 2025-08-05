package store

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"sync"

	"github.com/hashicorp/go-memdb"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/utils"
)

var noValidatedSources = []string{}

type SelectorContextKey struct{}

// do we really need this?
type Payload struct {
	Flags map[string]model.Flag
}

type IStore interface {
	GetAll(ctx context.Context, selector *Selector, watcher chan Payload) (map[string]model.Flag, model.Metadata, error)
	Get(ctx context.Context, key string, selector *Selector) (model.Flag, model.Metadata, bool)
}

type Store struct {
	mx      sync.RWMutex
	db      *memdb.MemDB
	logger  *logger.Logger
	sources []string
	// deprecated: has no effect and will be removed soon.
	FlagSources []string
}

type SourceDetails struct {
	Source   string
	Selector string
}

// NewStore creates a new in-memory store with the given sources.
// The order of sources in the slice determines their priority, when queries result in duplicate flags (queries without source or flagSetId), the higher priority source "wins".
func NewStore(logger *logger.Logger, sources []string) (*Store, error) {

	// a unique index must exist for each set of constraints - for example, to look up by key and source, we need a compound index on key+source, etc
	// we maybe want to generate these dynamically in the future to support more robust querying, but for now we will hardcode the ones we need
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			flagsTable: {
				Name: flagsTable,
				Indexes: map[string]*memdb.IndexSchema{
					// primary index; must be unique and named "id"
					idIndex: {
						Name:   idIndex,
						Unique: true,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: model.FlagSetId, Lowercase: false},
								&memdb.StringFieldIndex{Field: model.Key, Lowercase: false},
							},
						},
					},
					// for looking up by source
					sourceIndex: {
						Name:    sourceIndex,
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: model.Source, Lowercase: false},
					},
					// for looking up by priority, used to maintain highest priority flag when there are duplicates and no selector is provided
					priorityIndex: {
						Name:    priorityIndex,
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: model.Priority},
					},
					// for looking up by flagSetId
					flagSetIdIndex: {
						Name:    flagSetIdIndex,
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: model.FlagSetId, Lowercase: false},
					},
					keyIndex: {
						Name:    keyIndex,
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: model.Key, Lowercase: false},
					},
					flagSetIdSourceCompoundIndex: {
						Name:   flagSetIdSourceCompoundIndex,
						Unique: false,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: model.FlagSetId, Lowercase: false},
								&memdb.StringFieldIndex{Field: model.Source, Lowercase: false},
							},
						},
					},
					keySourceCompoundIndex: {
						Name:   keySourceCompoundIndex,
						Unique: false, // duplicate from a single source ARE allowed (they just must have different flag sets)
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: model.Key, Lowercase: false},
								&memdb.StringFieldIndex{Field: model.Source, Lowercase: false},
							},
						},
					},
					// used to query all flags from a specific source so we know which flags to delete if a flag is missing from a source
					flagSetIdKeySourceCompoundIndex: {
						Name:   flagSetIdKeySourceCompoundIndex,
						Unique: true,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: model.FlagSetId, Lowercase: false},
								&memdb.StringFieldIndex{Field: model.Key, Lowercase: false},
								&memdb.StringFieldIndex{Field: model.Source, Lowercase: false},
							},
						},
					},
				},
			},
		},
	}

	// Create a new data base
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize flag database: %w", err)
	}

	// clone the sources to avoid modifying the original slice
	s := slices.Clone(sources)

	return &Store{
		sources: s,
		db:      db,
		logger:  logger,
	}, nil
}

// Deprecated: use NewStore instead - will be removed very soon.
func NewFlags() *Store {
	state, err := NewStore(logger.NewLogger(nil, false), noValidatedSources)

	if err != nil {
		panic(fmt.Sprintf("unable to create flag store: %v", err))
	}
	return state
}

func (s *Store) Get(_ context.Context, key string, selector *Selector) (model.Flag, model.Metadata, bool) {
	s.logger.Debug(fmt.Sprintf("getting flag %s", key))
	txn := s.db.Txn(false)
	queryMeta := model.Metadata{}

	if selector != nil {
		queryMeta = selector.SelectorToMetadata()
		selector := selector.withIndex("key", key)
		indexId, constraints := selector.SelectorMapToQuery()
		s.logger.Debug(fmt.Sprintf("getting flag with query: %s, %v", indexId, constraints))
		raw, err := txn.First(flagsTable, indexId, constraints...)
		flag, ok := raw.(model.Flag)
		if (err != nil) || !ok {
			return model.Flag{}, queryMeta, false
		}
		return flag, queryMeta, true

	} else {
		// get all flags with the given key, and keep the one with the highest priority
		s.logger.Debug(fmt.Sprintf("getting highest priority flag with key: %s", key))
		it, err := txn.Get(flagsTable, keyIndex, key)
		if err != nil {
			return model.Flag{}, queryMeta, false
		}
		flag := model.Flag{}
		var found bool
		for raw := it.Next(); raw != nil; raw = it.Next() {
			found = true
			s.logger.Debug(fmt.Sprintf("got range scan: %v", raw))
			nextFlag, ok := raw.(model.Flag)
			if !ok {
				continue
			}
			if nextFlag.Priority >= flag.Priority {
				flag = nextFlag
			} else {
				s.logger.Debug(fmt.Sprintf("discarding flag %s from lower priority source %s in favor of flag from source %s", nextFlag.Key, s.sources[nextFlag.Priority], s.sources[flag.Priority]))
			}

		}
		return flag, queryMeta, found
	}
}

func (f *Store) String() (string, error) {
	f.logger.Debug("dumping flags to string")
	f.mx.RLock()
	defer f.mx.RUnlock()

	state, _, err := f.GetAll(context.Background(), nil, nil)
	if err != nil {
		return "", fmt.Errorf("unable to get all flags: %w", err)
	}

	bytes, err := json.Marshal(state)
	if err != nil {
		return "", fmt.Errorf("unable to marshal flags: %w", err)
	}

	return string(bytes), nil
}

// GetAll returns a copy of the store's state (copy in order to be concurrency safe)
func (s *Store) GetAll(ctx context.Context, selector *Selector, watcher chan Payload) (map[string]model.Flag, model.Metadata, error) {
	txn := s.db.Txn(false)
	flags := make(map[string]model.Flag)
	queryMeta := model.Metadata{}

	var it memdb.ResultIterator
	var err error
	if selector != nil && !selector.isEmpty() {
		queryMeta = selector.SelectorToMetadata()
		indexId, constraints := selector.SelectorMapToQuery()
		s.logger.Debug(fmt.Sprintf("getting all flags with query: %s, %v", indexId, constraints))
		it, err = txn.Get("flags", indexId, constraints...)
	} else {
		// no selector, get all flags
		it, err = txn.Get(flagsTable, idIndex)
	}

	if err != nil {
		s.logger.Error(fmt.Sprintf("query error: %v", err))
		return flags, queryMeta, err
	}

	for raw := it.Next(); raw != nil; raw = it.Next() {
		flag := raw.(model.Flag)
		if existing, ok := flags[flag.Key]; ok {
			if flag.Priority < existing.Priority {
				s.logger.Debug(fmt.Sprintf("discarding duplicate flag %s from lower priority source %s in favor of flag from source %s", flag.Key, s.sources[flag.Priority], s.sources[existing.Priority]))
				continue // we already have a higher priority flag
			} else {
				s.logger.Debug(fmt.Sprintf("overwriting duplicate flag %s from lower priority source %s in favor of flag from source %s", flag.Key, s.sources[existing.Priority], s.sources[flag.Priority]))
			}
		}
		flags[flag.Key] = flag
	}

	if watcher != nil {

		// a "one-time" watcher that will be notified of changes to the query
		changes := it.WatchCh()

		go func() {
			select {
			case <-changes:
				s.logger.Debug("flags store has changed, notifying watchers")

				// recursively get all flags again in bytes new goroutine to keep the watcher responsive
				// as long as we do this in bytes new goroutine, we don't risk stack overflow
				bytes, _, err := s.GetAll(ctx, selector, watcher)
				if err != nil {
					s.logger.Error(fmt.Sprintf("error getting flags in watcher: %v", err))
					break
				}
				watcher <- Payload{
					Flags: bytes,
				}

			case <-ctx.Done():
				close(watcher)
			}
		}()
	}

	return flags, queryMeta, nil
}

// Update the flag state with the provided flags.
func (s *Store) Update(
	source string,
	flags map[string]model.Flag,
	metadata model.Metadata,
) (map[string]interface{}, bool) {
	resyncRequired := false

	if source == "" {
		panic("source cannot be empty")
	}

	priority := slices.Index(s.sources, source)
	if priority == -1 {
		// this is a hack to allow old constructors that didn't pass sources, remove when we remove "NewFlags" constructor
		if !slices.Equal(s.sources, noValidatedSources) {
			panic(fmt.Sprintf("source %s is not registered in the store", source))
		}
		// same as above - remove when we remove "NewFlags" constructor
		priority = 0
	}

	txn := s.db.Txn(true)
	defer txn.Abort()

	// get all flags for the source we are updating
	selector := NewSelector(sourceIndex + "=" + source)
	oldFlags, _, _ := s.GetAll(context.Background(), &selector, nil)

	s.mx.Lock()
	for key := range oldFlags {
		if _, ok := flags[key]; !ok {
			// flag has been deleted
			s.logger.Debug(fmt.Sprintf("flag %s has been deleted from source %s", key, source))

			count, err := txn.DeleteAll(flagsTable, keySourceCompoundIndex, key, source)
			s.logger.Debug(fmt.Sprintf("deleted %d flags with key %s from source %s", count, key, source))

			if err != nil {
				s.logger.Error(fmt.Sprintf("error deleting flag: %s, %v", key, err))
				continue
			}

			s.logger.Debug(
				fmt.Sprintf(
					"store resync triggered: flag %s has been deleted from source %s",
					key, source,
				),
			)
			continue
		}
	}
	s.mx.Unlock()
	for key, newFlag := range flags {
		s.logger.Debug(fmt.Sprintf("got metadata %v", metadata))

		newFlag.Key = key
		newFlag.Source = source
		newFlag.Priority = priority
		newFlag.Metadata = mergeMetadata(metadata, newFlag.Metadata)

		// flagSetId defaults to a UUID generated at startup to make our quires isomorphic
		flagSetId := nilFlagSetId
		// flagSetId is inherited from the set, but can be overridden by the flag
		setFlagSetId, ok := newFlag.Metadata["flagSetId"].(string)
		if ok {
			flagSetId = setFlagSetId
		}
		newFlag.FlagSetId = flagSetId

		raw, err := txn.First(flagsTable, keySourceCompoundIndex, key, source)
		if err != nil {
			s.logger.Error(fmt.Sprintf("unable to get flag %s from source %s: %v", key, source, err))
			continue
		}
		oldFlag, ok := raw.(model.Flag)
		// if we already have a flag with the same key and source, we need to check if it has the same flagSetId
		if ok {
			if oldFlag.FlagSetId != newFlag.FlagSetId {
				_, err = txn.DeleteAll(flagsTable, idIndex, oldFlag.FlagSetId, key)
				if err != nil {
					s.logger.Error(fmt.Sprintf("unable to delete flags with key %s and flagSetId %s: %v", key, oldFlag.FlagSetId, err))
					continue
				}
			}
		}
		// Store the new version of the flag
		s.logger.Debug(fmt.Sprintf("storing flag: %v", newFlag))
		err = txn.Insert(flagsTable, newFlag)
		if err != nil {
			s.logger.Error(fmt.Sprintf("unable to insert flag %s: %v", key, err))
			continue
		}
	}

	txn.Commit()
	return utils.BuildNotifications(oldFlags, flags), resyncRequired
}

func mergeMetadata(m1, m2 model.Metadata) model.Metadata {
	merged := make(model.Metadata)
	if m1 == nil && m2 == nil {
		return nil
	}
	for key, value := range m1 {
		merged[key] = value
	}
	for key, value := range m2 { // m2 values overwrite m1 values on key conflict
		merged[key] = value
	}
	return merged
}
