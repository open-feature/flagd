package store

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"sync"

	"github.com/hashicorp/go-memdb"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
)

type del = struct{}

var deleteMarker *del

type IStore interface {
	GetAll(ctx context.Context) (map[string]model.Flag, model.Metadata, error)
	Get(ctx context.Context, key string) (model.Flag, model.Metadata, bool)
	SelectorForFlag(ctx context.Context, flag model.Flag) string
}

type Store struct {
	mx                sync.RWMutex
	db                *memdb.MemDB
	logger            *logger.Logger
	FlagSources       []string
	SourceDetails     map[string]SourceDetails  `json:"sourceMetadata,omitempty"`
	MetadataPerSource map[string]model.Metadata `json:"metadata,omitempty"`
}

type SourceDetails struct {
	Source   string
	Selector string
}

func (f *Store) hasPriority(stored string, new string) bool {
	if stored == new {
		return true
	}
	for i := len(f.FlagSources) - 1; i >= 0; i-- {
		switch f.FlagSources[i] {
		case stored:
			return false
		case new:
			return true
		}
	}
	return true
}

func NewStore(logger *logger.Logger) (*Store, error) {

	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"flags": {
				Name: "flags",
				Indexes: map[string]*memdb.IndexSchema{

					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Key", Lowercase: false},
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

	return &Store{
		SourceDetails:     map[string]SourceDetails{},
		MetadataPerSource: map[string]model.Metadata{},
		db:                db,
		logger:            logger,
	}, nil
}

// Deprecated: use NewStore instead
func NewFlags() *Store {
	state, err := NewStore(logger.NewLogger(nil, false))
	if err != nil {
		panic(fmt.Sprintf("unable to create flag store: %v", err))
	}
	return state
}

func (f *Store) Get(_ context.Context, key string) (model.Flag, model.Metadata, bool) {
	f.logger.Debug(fmt.Sprintf("getting flag %s", key))
	txn := f.db.Txn(false)

	raw, err := txn.First("flags", "id", key)
	flag, ok := raw.(model.Flag)
	if err != nil || !ok {
		return model.Flag{}, f.getMetadata(), false
	}
	return flag, f.GetMetadataForSource(flag.Source), true
}

func (f *Store) SelectorForFlag(_ context.Context, flag model.Flag) string {
	f.mx.RLock()
	defer f.mx.RUnlock()

	return f.SourceDetails[flag.Source].Selector
}

func (f *Store) String() (string, error) {
	f.logger.Debug("dumping flags to string")
	f.mx.RLock()
	defer f.mx.RUnlock()

	state, _, err := f.GetAll(context.Background())
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
func (f *Store) GetAll(_ context.Context) (map[string]model.Flag, model.Metadata, error) {
	txn := f.db.Txn(false)

	flags := make(map[string]model.Flag)
	it, err := txn.Get("flags", "id")

	if err != nil {
		return flags, model.Metadata{}, err
	}

	for obj := it.Next(); obj != nil; obj = it.Next() {
		flag := obj.(model.Flag)
		flags[flag.Key] = flag
	}

	return flags, f.getMetadata(), nil
}

// Update the flag state with the provided flags.
func (f *Store) Update(
	source string,
	selector string,
	flags map[string]model.Flag,
	metadata model.Metadata,
) (map[string]interface{}, bool) {
	notifications := map[string]interface{}{}
	resyncRequired := false

	txn := f.db.Txn(true)
	defer txn.Abort()
	storedFlags, _, _ := f.GetAll(context.Background())

	f.mx.Lock()
	f.setSourceMetadata(source, metadata)
	for key, v := range storedFlags {
		if v.Source == source && v.Selector == selector {
			if _, ok := flags[key]; !ok {
				// flag has been deleted
				_, err := txn.DeleteAll("flags", "id", key)
				if err != nil {
					f.logger.Error(fmt.Sprintf("error deleting flag: %s, %v", key, err))
					continue
				}

				notifications[key] = map[string]interface{}{
					"type":   string(model.NotificationDelete),
					"source": source,
				}
				resyncRequired = true
				f.logger.Debug(
					fmt.Sprintf(
						"store resync triggered: flag %s has been deleted from source %s",
						key, source,
					),
				)
				continue
			}
		}
	}
	f.mx.Unlock()
	for key, newFlag := range flags {
		newFlag.Source = source
		newFlag.Selector = selector
		newFlag.Key = key
		storedFlag, _, ok := f.Get(context.Background(), key)
		if ok {
			if !f.hasPriority(storedFlag.Source, source) {
				f.logger.Debug(
					fmt.Sprintf(
						"not merging: flag %s from source %s does not have priority over %s",
						key, source, storedFlag.Source,
					),
				)
				continue
			}
			if reflect.DeepEqual(storedFlag, newFlag) {
				continue
			}
		}
		if !ok {
			notifications[key] = map[string]interface{}{
				"type":   string(model.NotificationCreate),
				"source": source,
			}
		} else {
			notifications[key] = map[string]interface{}{
				"type":   string(model.NotificationUpdate),
				"source": source,
			}
		}
		// Store the new version of the flag
		err := txn.Insert("flags", newFlag)
		if err != nil {
			f.logger.Error(fmt.Sprintf("unable to insert flag %s: %v", key, err))
			continue
		}
	}

	txn.Commit()
	return notifications, resyncRequired
}

func (f *Store) GetMetadataForSource(source string) model.Metadata {
	f.mx.RLock()
	defer f.mx.RUnlock()
	perSource, ok := f.MetadataPerSource[source]
	if ok && perSource != nil {
		return maps.Clone(perSource)
	}
	return model.Metadata{}
}

// TODO: this is a temporary solution to merge metadata in the case of error; properly handle it with https://github.com/open-feature/flagd/issues/1675
func (f *Store) getMetadata() model.Metadata {
	f.mx.RLock()
	defer f.mx.RUnlock()
	metadata := model.Metadata{}
	for _, perSource := range f.MetadataPerSource {
		for key, entry := range perSource {
			_, exists := metadata[key]
			if !exists {
				metadata[key] = entry
			} else {
				metadata[key] = deleteMarker
			}
		}
	}

	// keys that exist across multiple sources are deleted
	maps.DeleteFunc(metadata, func(key string, _ interface{}) bool {
		return metadata[key] == deleteMarker
	})

	return metadata
}

func (f *Store) setSourceMetadata(source string, metadata model.Metadata) {
	if f.MetadataPerSource == nil {
		f.MetadataPerSource = map[string]model.Metadata{}
	}

	f.MetadataPerSource[source] = metadata
}
