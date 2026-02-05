package store

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sort"

	"github.com/hashicorp/go-memdb"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
)

var noValidatedSources = []string{}

type SelectorContextKey struct{}

type FlagQueryResult struct {
	Flags []model.Flag
}

type IStore interface {
	Get(ctx context.Context, key string, selector *Selector) (model.Flag, model.Metadata, error)
	GetAll(ctx context.Context, selector *Selector) ([]model.Flag, model.Metadata, error)
	Watch(ctx context.Context, selector *Selector, watcher chan<- FlagQueryResult)
	Update(source string, flags []model.Flag, metadata model.Metadata)
}

var _ IStore = (*Store)(nil)

type Store struct {
	db      *memdb.MemDB
	logger  *logger.Logger
	sources []string
	// deprecated: has no effect and will be removed soon.
	FlagSources []string
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

func (s *Store) Get(_ context.Context, key string, selector *Selector) (model.Flag, model.Metadata, error) {
	s.logger.Debug(fmt.Sprintf("getting flag %s", key))
	txn := s.db.Txn(false)
	queryMeta := selector.ToMetadata()

	// if present, use the selector to query the flags
	if !selector.IsEmpty() {
		selector := selector.WithIndex("key", key)
		indexId, constraints := selector.ToQuery()
		s.logger.Debug(fmt.Sprintf("getting flag with query: %s, %v", indexId, constraints))
		raw, err := txn.First(flagsTable, indexId, constraints...)
		flag, ok := raw.(model.Flag)
		if err != nil {
			return model.Flag{}, queryMeta, fmt.Errorf("flag %s not found: %w", key, err)
		}
		if !ok {
			return model.Flag{}, queryMeta, fmt.Errorf("flag %s is not a valid flag", key)
		}
		return flag, queryMeta, nil
	}

	// otherwise, get all flags with the given key, and keep the last one with the highest priority
	s.logger.Debug(fmt.Sprintf("getting highest priority flag with key: %s", key))
	it, err := txn.Get(flagsTable, keyIndex, key)
	if err != nil {
		return model.Flag{}, queryMeta, fmt.Errorf("flag %s not found: %w", key, err)
	}
	flag := model.Flag{}
	found := false
	for raw := it.Next(); raw != nil; raw = it.Next() {
		nextFlag, ok := raw.(model.Flag)
		if !ok {
			continue
		}
		found = true
		if nextFlag.Priority >= flag.Priority {
			flag = nextFlag
		} else {
			s.logger.Debug(fmt.Sprintf("discarding flag %s from lower priority source %s in favor of flag from source %s", nextFlag.Key, s.sources[nextFlag.Priority], s.sources[flag.Priority]))
		}
	}

	if !found {
		return flag, queryMeta, fmt.Errorf("flag %s not found", key)
	}
	return flag, queryMeta, nil
}

// GetAll returns a copy of the store's state (copy in order to be concurrency safe)
func (s *Store) GetAll(ctx context.Context, selector *Selector) ([]model.Flag, model.Metadata, error) {
	var flags []model.Flag
	queryMeta := selector.ToMetadata()
	it, err := s.selectOrAll(selector)

	if err != nil {
		s.logger.Error(fmt.Sprintf("flag query error: %v", err))
		return flags, queryMeta, err
	}
	flags = s.collect(it)
	return flags, queryMeta, nil
}

type flagIdentifier struct {
	flagSetId string
	key       string
}

// Update the flag state with the provided flags.
func (s *Store) Update(
	source string,
	flags []model.Flag,
	metadata model.Metadata,
) {
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
	newFlags := make(map[flagIdentifier]model.Flag)
	for _, newFlag := range flags {
		s.logger.Debug(fmt.Sprintf("got metadata %v", metadata))

		newFlag.Source = source
		newFlag.Priority = priority
		newFlag.Metadata = patchMetadata(metadata, newFlag.Metadata)

		// flagSetId defaults to a UUID generated at startup to make our queries isomorphic
		flagSetId := nilFlagSetId
		// flagSetId is inherited from the set, but can be overridden by the flag
		setFlagSetId, ok := newFlag.Metadata["flagSetId"].(string)
		if ok {
			flagSetId = setFlagSetId
		}
		newFlag.FlagSetId = flagSetId
		newFlags[flagIdentifier{flagSetId: newFlag.FlagSetId, key: newFlag.Key}] = newFlag
	}

	txn := s.db.Txn(true)
	defer txn.Abort()

	// get all flags for the source we are updating
	selector := NewSelector(sourceIndex + "=" + source)
	oldFlags, _, _ := s.GetAll(context.Background(), &selector)

	for _, oldFlag := range oldFlags {
		if _, ok := newFlags[flagIdentifier{flagSetId: oldFlag.FlagSetId, key: oldFlag.Key}]; !ok {
			// flag has been deleted
			s.logger.Debug(fmt.Sprintf("flag '%s' and flagSetId '%s' has been deleted from source '%s'", oldFlag.Key, oldFlag.FlagSetId, source))

			count, err := txn.DeleteAll(flagsTable, flagSetIdKeySourceCompoundIndex, oldFlag.FlagSetId, oldFlag.Key, source)
			s.logger.Debug(fmt.Sprintf(
				"deleted %d flags with key '%s' and flagSetId '%s' from source '%s'",
				count,
				oldFlag.Key,
				oldFlag.FlagSetId,
				source,
			))

			if err != nil {
				s.logger.Error(fmt.Sprintf("error deleting flag: %s, %v", oldFlag.Key, err))
			}
			continue
		}
	}

	for _, newFlag := range newFlags {
		raw, err := txn.First(flagsTable, idIndex, newFlag.FlagSetId, newFlag.Key)
		if err != nil {
			s.logger.Error(fmt.Sprintf("unable to get flag %s from source %s: %v", newFlag.Key, newFlag.FlagSetId, err))
			continue
		}
		oldFlag, ok := raw.(model.Flag)
		// If we already have a flag with the same key and source, we need to check if it has the same flagSetId
		if ok {
			if oldFlag.Priority > newFlag.Priority {
				// if the old flag has a higher prio, we should not try to write it
				s.logger.Error(fmt.Sprintf("unable to delete flags with key %s and flagSetId %s: %v", oldFlag.Key, oldFlag.FlagSetId, err))
				continue
			}
		}

		// Store the new version of the flag
		s.logger.Debug(fmt.Sprintf("storing flag: %v", newFlag))

		err = txn.Insert(flagsTable, newFlag)
		if err != nil {
			s.logger.Error(fmt.Sprintf("unable to insert flag %s: %v", newFlag.Key, err))
			continue
		}
	}

	txn.Commit()
}

// Watch the result-set of a selector for changes, sending updates to the watcher channel.
func (s *Store) Watch(ctx context.Context, selector *Selector, watcher chan<- FlagQueryResult) {
	go func() {
		for {
			ws := memdb.NewWatchSet()
			it, err := s.selectOrAll(selector)
			if err != nil {
				s.logger.Error(fmt.Sprintf("error getting flags for selector %s: %v", selector.ToLogString(), err))
				close(watcher)
				return
			}
			ws.Add(it.WatchCh())

			flags := s.collect(it)

			watcher <- FlagQueryResult{
				Flags: flags,
			}

			if err = ws.WatchCtx(ctx); err != nil {
				if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
					s.logger.Debug(fmt.Sprintf("while watching flags for selector %s: %v", selector.ToLogString(), err))
				} else {
					s.logger.Error(fmt.Sprintf("context error watching flags for selector %s: %v", selector.ToLogString(), err))
				}
				close(watcher)
				return
			}
		}
	}()
}

// returns an iterator for the given selector, or all flags if the selector is nil or empty
func (s *Store) selectOrAll(selector *Selector) (it memdb.ResultIterator, err error) {
	txn := s.db.Txn(false)
	if !selector.IsEmpty() {
		indexId, constraints := selector.ToQuery()
		s.logger.Debug(fmt.Sprintf("getting all flags with query: %s, %v", indexId, constraints))
		return txn.Get(flagsTable, indexId, constraints...)
	} else {
		// no selector, get all flags
		return txn.Get(flagsTable, idIndex)
	}
}

// collects flags from an iterator, ensuring that only the highest priority flag is kept when there are duplicates
func (s *Store) collect(it memdb.ResultIterator) []model.Flag {
	flags := make(map[flagIdentifier]model.Flag)
	for raw := it.Next(); raw != nil; raw = it.Next() {
		flag := raw.(model.Flag)

		// checking for multiple flags with the same key, as they can be defined multiple times in different sources
		if existing, ok := flags[flagIdentifier{flagSetId: flag.FlagSetId, key: flag.Key}]; ok {
			if flag.Priority < existing.Priority {
				s.logger.Debug(fmt.Sprintf("discarding duplicate flag with key '%s' and flagSetId '%s' from lower priority source '%s' in favor of flag from source '%s'", flag.Key, flag.FlagSetId, s.sources[flag.Priority], s.sources[existing.Priority]))
				continue // we already have a higher priority flag
			}
			s.logger.Debug(fmt.Sprintf("overwriting duplicate flag with key '%s' and flagSetId '%s' from lower priority source '%s' in favor of flag from source '%s'", flag.Key, flag.FlagSetId, s.sources[existing.Priority], s.sources[flag.Priority]))
		}

		flags[flagIdentifier{flagSetId: flag.FlagSetId, key: flag.Key}] = flag
	}

	flattenedFlags := make([]model.Flag, 0, len(flags))
	for _, value := range flags {
		flattenedFlags = append(flattenedFlags, value)
	}
	// we should order to keep the same order all the time in our response
	sort.Slice(flattenedFlags, func(i, j int) bool {
		if flattenedFlags[i].FlagSetId != flattenedFlags[j].FlagSetId {
			return flattenedFlags[i].FlagSetId < flattenedFlags[j].FlagSetId
		}
		return flattenedFlags[i].Key < flattenedFlags[j].Key
	})
	return flattenedFlags
}

func patchMetadata(original, patch model.Metadata) model.Metadata {
	patched := make(model.Metadata)
	if original == nil && patch == nil {
		return nil
	}
	for key, value := range original {
		patched[key] = value
	}
	for key, value := range patch { // patch values overwrite m1 values on key conflict
		patched[key] = value
	}
	return patched
}
