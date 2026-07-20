package store

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sort"
	"sync/atomic"

	"github.com/hashicorp/go-memdb"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"go.uber.org/zap"
)

var noValidatedSources = []string{}

type SelectorContextKey struct{}

type FlagQueryResult struct {
	Flags []model.Flag
}

// flagSetMembership tracks which (flagSetId, source) combinations include a given flag key.
// Used by the membership table to deduplicate flags across flagSetIds during incremental updates.
type flagSetMembership struct {
	FlagSetId string
	Key       string
	Source    string
}

type IStore interface {
	Get(ctx context.Context, key string, selector *Selector) (model.Flag, model.Metadata, error)
	GetAll(ctx context.Context, selector *Selector) ([]model.Flag, model.Metadata, error)
	Watch(ctx context.Context, selector *Selector, watcher chan<- FlagQueryResult)
	Update(source string, flags []model.Flag, metadata model.Metadata, incrementalUpdate bool)
}

var _ IStore = (*Store)(nil)

type Store struct {
	db      *memdb.MemDB
	logger  *logger.Logger
	sources []string
	// deprecated: has no effect and will be removed soon.
	FlagSources []string
	// hasMembership is set to true after the first incremental update.
	// When false, Get/GetAll/Watch skip membership resolution entirely,
	// ensuring zero behavioral change for non-incremental callers.
	hasMembership atomic.Bool
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
			membershipTable: {
				Name: membershipTable,
				Indexes: map[string]*memdb.IndexSchema{
					idIndex: {
						Name:   idIndex,
						Unique: true,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "FlagSetId", Lowercase: false},
								&memdb.StringFieldIndex{Field: "Key", Lowercase: false},
								&memdb.StringFieldIndex{Field: "Source", Lowercase: false},
							},
						},
					},
					flagSetIdIndex: {
						Name:    flagSetIdIndex,
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "FlagSetId", Lowercase: false},
					},
					membershipFlagSetIdKeyIndex: {
						Name:   membershipFlagSetIdKeyIndex,
						Unique: false,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "FlagSetId", Lowercase: false},
								&memdb.StringFieldIndex{Field: "Key", Lowercase: false},
							},
						},
					},
					flagSetIdSourceCompoundIndex: {
						Name:   flagSetIdSourceCompoundIndex,
						Unique: false,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "FlagSetId", Lowercase: false},
								&memdb.StringFieldIndex{Field: "Source", Lowercase: false},
							},
						},
					},
					keySourceCompoundIndex: {
						Name:   keySourceCompoundIndex,
						Unique: false,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Key", Lowercase: false},
								&memdb.StringFieldIndex{Field: "Source", Lowercase: false},
							},
						},
					},
					sourceIndex: {
						Name:    sourceIndex,
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Source", Lowercase: false},
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
		if err != nil {
			return model.Flag{}, queryMeta, fmt.Errorf("flag %s not found: %w", key, err)
		}
		flag, ok := raw.(model.Flag)
		if ok {
			return flag, queryMeta, nil
		}

		// Flag not found directly — try membership resolution for flagSetId queries.
		// With dedup, flags may be stored under a different flagSetId.
		if s.hasMembership.Load() {
			if flagSetId, hasFSI := selector.HasFlagSetId(); hasFSI {
				memberIt, mErr := txn.Get(membershipTable, membershipFlagSetIdKeyIndex, flagSetId, key)
				if mErr == nil {
					var best model.Flag
					found := false
					for mRaw := memberIt.Next(); mRaw != nil; mRaw = memberIt.Next() {
						m := mRaw.(flagSetMembership)
						flagIt, fErr := txn.Get(flagsTable, keySourceCompoundIndex, m.Key, m.Source)
						if fErr != nil {
							continue
						}
						if fRaw := flagIt.Next(); fRaw != nil {
							candidate := fRaw.(model.Flag)
							if !found || candidate.Priority >= best.Priority {
								best = candidate
								found = true
							}
						}
					}
					if found {
						best.FlagSetId = flagSetId
						if best.Metadata == nil {
							best.Metadata = make(model.Metadata)
						} else {
							patched := make(model.Metadata, len(best.Metadata))
							for k, v := range best.Metadata {
								patched[k] = v
							}
							best.Metadata = patched
						}
						best.Metadata["flagSetId"] = flagSetId
						return best, queryMeta, nil
					}
				}
			}
		}

		return model.Flag{}, queryMeta, fmt.Errorf("flag %s is not a valid flag", key)
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

	// For flagSetId selectors, try membership resolution first
	if s.hasMembership.Load() {
		if flagSetId, hasFSI := selector.HasFlagSetId(); hasFSI {
			txn := s.db.Txn(false)
			if flags := s.collectViaMembership(txn, flagSetId, nil); flags != nil {
				return flags, queryMeta, nil
			}
		}
	}

	// Fall back to direct flags table query
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
// When incrementalUpdate is true, deletion is scoped to only the flagSetIds present in
// this payload (from metadata and flag-level overrides), allowing flags from other
// flagSetIds to accumulate across updates. Flags are deduplicated by (key, source) so
// that identical flags shared across flagSetIds are stored only once.
// When false, all flags for the source are replaced (the default full-snapshot behavior).
// EXPERIMENTAL: incrementalUpdate support may change or be removed in a future release.
func (s *Store) Update(
	source string,
	flags []model.Flag,
	metadata model.Metadata,
	incrementalUpdate bool,
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

	if incrementalUpdate {
		s.updateIncremental(txn, source, newFlags, metadata)
	} else {
		s.updateFullSnapshot(txn, source, newFlags)
	}

	txn.Commit()
}

// updateIncremental handles membership-aware incremental updates with deduplication.
// Flags are stored once per (key, source) in the flags table. A lightweight membership table
// tracks which flagSetIds include which keys, enabling per-flagSetId queries without duplication.
func (s *Store) updateIncremental(txn *memdb.Txn, source string, newFlags map[flagIdentifier]model.Flag, metadata model.Metadata) {
	s.hasMembership.Store(true)

	// Step 1: Determine flagSetIds touched by this payload (from metadata + flag-level overrides)
	seenFlagSetIds := make(map[string]struct{})
	if fsi, ok := metadata["flagSetId"].(string); ok && fsi != "" {
		seenFlagSetIds[fsi] = struct{}{}
	}
	for id := range newFlags {
		seenFlagSetIds[id.flagSetId] = struct{}{}
	}

	// Step 2: Collect old membership entries for each touched flagSetId+source
	oldMembership := make(map[flagIdentifier]struct{})
	for fsi := range seenFlagSetIds {
		it, err := txn.Get(membershipTable, flagSetIdSourceCompoundIndex, fsi, source)
		if err != nil {
			s.logger.Error(fmt.Sprintf("unable to query membership for flagSetId %s: %v", fsi, err))
			continue
		}
		for raw := it.Next(); raw != nil; raw = it.Next() {
			m := raw.(flagSetMembership)
			oldMembership[flagIdentifier{flagSetId: m.FlagSetId, key: m.Key}] = struct{}{}
		}
	}

	// Step 3: Build new membership set
	newMembership := make(map[flagIdentifier]struct{}, len(newFlags))
	for id := range newFlags {
		newMembership[id] = struct{}{}
	}

	// Step 4: Delete stale membership entries and orphaned flags
	for oldId := range oldMembership {
		if _, ok := newMembership[oldId]; ok {
			continue // still present
		}
		// Remove stale membership entry
		if _, err := txn.DeleteAll(membershipTable, idIndex, oldId.flagSetId, oldId.key, source); err != nil {
			s.logger.Error(fmt.Sprintf("error deleting membership: flagSetId=%s key=%s: %v", oldId.flagSetId, oldId.key, err))
		}
		// Check if any other flagSetId still references this key+source
		refIt, err := txn.Get(membershipTable, keySourceCompoundIndex, oldId.key, source)
		if err != nil {
			s.logger.Error(fmt.Sprintf("error checking membership refs for key %s: %v", oldId.key, err))
			continue
		}
		if refIt.Next() == nil {
			// No more references — delete the flag from the flags table
			count, err := txn.DeleteAll(flagsTable, keySourceCompoundIndex, oldId.key, source)
			if err != nil {
				s.logger.Error(fmt.Sprintf("error deleting orphaned flag %s: %v", oldId.key, err))
			} else {
				s.logger.Debug(fmt.Sprintf("deleted %d orphaned flag(s) with key '%s' from source '%s'", count, oldId.key, source))
			}
		}
	}

	// Step 5: Insert/update membership entries and deduplicate flags
	for id, newFlag := range newFlags {
		// Upsert membership entry
		if err := txn.Insert(membershipTable, flagSetMembership{
			FlagSetId: id.flagSetId,
			Key:       id.key,
			Source:    source,
		}); err != nil {
			s.logger.Error(fmt.Sprintf("unable to insert membership for flagSetId=%s key=%s: %v", id.flagSetId, id.key, err))
			continue
		}

		// Dedup: check if a flag with the same key+source already exists (from any flagSetId)
		existingIt, err := txn.Get(flagsTable, keySourceCompoundIndex, newFlag.Key, source)
		if err != nil {
			s.logger.Error(fmt.Sprintf("unable to check existing flag %s: %v", newFlag.Key, err))
			continue
		}
		existing := existingIt.Next()
		if existing != nil {
			existingFlag := existing.(model.Flag)
			if existingFlag.Priority > newFlag.Priority {
				// Higher priority source already owns this flag, skip
				s.logger.Debug(fmt.Sprintf("flag '%s' owned by higher priority source, skipping", newFlag.Key))
				continue
			}
			// Flag already exists at same or lower priority — update in place to pick up content changes.
			// Preserve the existing entry's FlagSetId so the upsert overwrites the canonical row.
			newFlag.FlagSetId = existingFlag.FlagSetId
			s.logger.Debug(fmt.Sprintf("updating existing flag '%s' (canonical flagSetId: %s) for flagSetId '%s'", newFlag.Key, existingFlag.FlagSetId, id.flagSetId))
			if err := txn.Insert(flagsTable, newFlag); err != nil {
				s.logger.Error(fmt.Sprintf("unable to update existing flag %s: %v", newFlag.Key, err))
			}
			continue
		}

		// New flag — insert into flags table
		s.logger.Debug(fmt.Sprintf("storing flag: %s (flagSetId: %s)", newFlag.Key, id.flagSetId))
		if err := txn.Insert(flagsTable, newFlag); err != nil {
			s.logger.Error(fmt.Sprintf("unable to insert flag %s: %v", newFlag.Key, err))
		}
	}
}

// updateFullSnapshot replaces all flags for the source (non-incremental mode).
func (s *Store) updateFullSnapshot(txn *memdb.Txn, source string, newFlags map[flagIdentifier]model.Flag) {
	// Clean up any membership entries for this source (in case a previous incremental update left them)
	if _, err := txn.DeleteAll(membershipTable, sourceIndex, source); err != nil {
		s.logger.Error(fmt.Sprintf("error cleaning membership for source %s: %v", source, err))
	}

	// Get all existing flags for this source
	sel := NewSelector(sourceIndex + "=" + source)
	indexId, constraints := sel.ToQuery()
	it, err := txn.Get(flagsTable, indexId, constraints...)
	var oldFlags []model.Flag
	if err != nil {
		s.logger.Error(fmt.Sprintf("unable to query flags for source %s: %v", source, err))
	} else {
		oldFlags = s.collect(it)
	}

	// Delete flags not in the new set
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
				s.logger.Error(fmt.Sprintf("unable to update flags with key %s and flagSetId %s: higher priority exists", oldFlag.Key, oldFlag.FlagSetId))
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
}

// Watch the result-set of a selector for changes, sending updates to the watcher channel.
func (s *Store) Watch(ctx context.Context, selector *Selector, watcher chan<- FlagQueryResult) {
	go func() {
		for {
			ws := memdb.NewWatchSet()
			txn := s.db.Txn(false)

			var flags []model.Flag

			// For flagSetId selectors, always watch the membership table so we
			// detect new membership additions even when starting with zero entries.
			if s.hasMembership.Load() {
				if flagSetId, hasFSI := selector.HasFlagSetId(); hasFSI {
					memberIt, err := txn.Get(membershipTable, flagSetIdIndex, flagSetId)
					if err == nil {
						ws.Add(memberIt.WatchCh())
					}
					if membershipFlags := s.collectViaMembership(txn, flagSetId, &ws); membershipFlags != nil {
						flags = membershipFlags
					}
				}
			}

			// Fall back to direct flags table query if membership didn't produce results
			if flags == nil {
				it, err := s.selectOrAllWithTxn(txn, selector)
				if err != nil {
					s.logger.WithFields(zap.String("selector", selector.ToLogString()), zap.Error(err)).Error("error getting flags")
					close(watcher)
					return
				}
				ws.Add(it.WatchCh())
				flags = s.collect(it)
			}

			watcher <- FlagQueryResult{
				Flags: flags,
			}

			if err := ws.WatchCtx(ctx); err != nil {
				if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
					s.logger.WithFields(zap.String("selector", selector.ToLogString()), zap.Error(err)).Debug("context cancellation while watching flags")
				} else {
					s.logger.WithFields(zap.String("selector", selector.ToLogString()), zap.Error(err)).Error("context error watching flags")
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
	return s.selectOrAllWithTxn(txn, selector)
}

func (s *Store) selectOrAllWithTxn(txn *memdb.Txn, selector *Selector) (it memdb.ResultIterator, err error) {
	if !selector.IsEmpty() {
		indexId, constraints := selector.ToQuery()
		s.logger.Debug(fmt.Sprintf("getting all flags with query: %s, %v", indexId, constraints))
		return txn.Get(flagsTable, indexId, constraints...)
	} else {
		// no selector, get all flags
		return txn.Get(flagsTable, idIndex)
	}
}

// collectViaMembership resolves flags for a flagSetId through the membership table.
// Each membership entry maps to a canonical flag stored in the flags table by key+source.
// If ws is non-nil, scoped watches are added for each (key, source) lookup so that only
// changes to flags relevant to this flagSetId trigger the watch.
func (s *Store) collectViaMembership(txn *memdb.Txn, flagSetId string, ws *memdb.WatchSet) []model.Flag {
	memberIt, err := txn.Get(membershipTable, flagSetIdIndex, flagSetId)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error querying membership for flagSetId %s: %v", flagSetId, err))
		return nil
	}

	flags := make(map[string]model.Flag) // key -> flag (dedup by key, keep highest priority)
	hasMembership := false
	for raw := memberIt.Next(); raw != nil; raw = memberIt.Next() {
		hasMembership = true
		m := raw.(flagSetMembership)
		flagIt, fErr := txn.Get(flagsTable, keySourceCompoundIndex, m.Key, m.Source)
		if fErr != nil {
			continue
		}
		if ws != nil {
			ws.Add(flagIt.WatchCh())
		}
		for fRaw := flagIt.Next(); fRaw != nil; fRaw = flagIt.Next() {
			flag := fRaw.(model.Flag)
			flag.FlagSetId = flagSetId // patch to match the queried flagSetId
			if flag.Metadata == nil {
				flag.Metadata = make(model.Metadata)
			} else {
				patched := make(model.Metadata, len(flag.Metadata))
				for k, v := range flag.Metadata {
					patched[k] = v
				}
				flag.Metadata = patched
			}
			flag.Metadata["flagSetId"] = flagSetId
			if existing, ok := flags[flag.Key]; ok {
				if flag.Priority < existing.Priority {
					continue
				}
			}
			flags[flag.Key] = flag
		}
	}

	if !hasMembership {
		return nil // signal to caller that no membership exists (fall back to direct query)
	}

	result := make([]model.Flag, 0, len(flags))
	for _, f := range flags {
		result = append(result, f)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Key < result[j].Key
	})
	return result
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
