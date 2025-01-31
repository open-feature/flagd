package store

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"sync"

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

type State struct {
	mx                sync.RWMutex
	Flags             map[string]model.Flag `json:"flags"`
	FlagSources       []string
	SourceDetails     map[string]SourceDetails  `json:"sourceMetadata,omitempty"`
	MetadataPerSource map[string]model.Metadata `json:"metadata,omitempty"`
}

type SourceDetails struct {
	Source   string
	Selector string
}

func (f *State) hasPriority(stored string, new string) bool {
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

func NewFlags() *State {
	return &State{
		Flags:             map[string]model.Flag{},
		SourceDetails:     map[string]SourceDetails{},
		MetadataPerSource: map[string]model.Metadata{},
	}
}

func (f *State) Set(key string, flag model.Flag) {
	f.mx.Lock()
	defer f.mx.Unlock()
	f.Flags[key] = flag
}

func (f *State) Get(_ context.Context, key string) (model.Flag, model.Metadata, bool) {
	f.mx.RLock()
	defer f.mx.RUnlock()
	metadata := f.getMetadata()
	flag, ok := f.Flags[key]
	if ok {
		metadata = f.getMetadataForSource(flag.Source)
	}

	return flag, metadata, ok
}

func (f *State) SelectorForFlag(_ context.Context, flag model.Flag) string {
	f.mx.RLock()
	defer f.mx.RUnlock()

	return f.SourceDetails[flag.Source].Selector
}

func (f *State) Delete(key string) {
	f.mx.Lock()
	defer f.mx.Unlock()
	delete(f.Flags, key)
}

func (f *State) String() (string, error) {
	f.mx.RLock()
	defer f.mx.RUnlock()
	bytes, err := json.Marshal(f)
	if err != nil {
		return "", fmt.Errorf("unable to marshal flags: %w", err)
	}

	return string(bytes), nil
}

// GetAll returns a copy of the store's state (copy in order to be concurrency safe)
func (f *State) GetAll(_ context.Context) (map[string]model.Flag, model.Metadata, error) {
	f.mx.RLock()
	defer f.mx.RUnlock()
	flags := make(map[string]model.Flag, len(f.Flags))

	for key, flag := range f.Flags {
		flags[key] = flag
	}

	return flags, f.getMetadata(), nil
}

// Add new flags from source.
func (f *State) Add(logger *logger.Logger, source string, selector string, flags map[string]model.Flag,
) map[string]interface{} {
	notifications := map[string]interface{}{}

	for k, newFlag := range flags {
		storedFlag, _, ok := f.Get(context.Background(), k)
		if ok && !f.hasPriority(storedFlag.Source, source) {
			logger.Debug(
				fmt.Sprintf(
					"not overwriting: flag %s from source %s does not have priority over %s",
					k,
					source,
					storedFlag.Source,
				),
			)
			continue
		}

		notifications[k] = map[string]interface{}{
			"type":   string(model.NotificationCreate),
			"source": source,
		}

		// Store the new version of the flag
		newFlag.Source = source
		newFlag.Selector = selector
		f.Set(k, newFlag)
	}

	return notifications
}

// Update existing flags from source.
func (f *State) Update(logger *logger.Logger, source string, selector string, flags map[string]model.Flag,
) map[string]interface{} {
	notifications := map[string]interface{}{}

	for k, flag := range flags {
		storedFlag, _, ok := f.Get(context.Background(), k)
		if !ok {
			logger.Warn(
				fmt.Sprintf("failed to update the flag, flag with key %s from source %s does not exist.",
					k,
					source))

			continue
		}
		if !f.hasPriority(storedFlag.Source, source) {
			logger.Debug(
				fmt.Sprintf(
					"not updating: flag %s from source %s does not have priority over %s",
					k,
					source,
					storedFlag.Source,
				),
			)
			continue
		}

		notifications[k] = map[string]interface{}{
			"type":   string(model.NotificationUpdate),
			"source": source,
		}

		flag.Source = source
		flag.Selector = selector
		f.Set(k, flag)
	}

	return notifications
}

// DeleteFlags matching flags from source.
func (f *State) DeleteFlags(logger *logger.Logger, source string, flags map[string]model.Flag) map[string]interface{} {
	logger.Debug(
		fmt.Sprintf(
			"store resync triggered: delete event from source %s",
			source,
		),
	)
	ctx := context.Background()

	_, ok := f.MetadataPerSource[source]
	if ok {
		delete(f.MetadataPerSource, source)
	}

	notifications := map[string]interface{}{}
	if len(flags) == 0 {
		allFlags, _, err := f.GetAll(ctx)
		if err != nil {
			logger.Error(fmt.Sprintf("error while retrieving flags from the store: %v", err))
			return notifications
		}
		for key, flag := range allFlags {
			if flag.Source != source {
				continue
			}
			notifications[key] = map[string]interface{}{
				"type":   string(model.NotificationDelete),
				"source": source,
			}
			f.Delete(key)
		}
	}

	for k := range flags {
		flag, _, ok := f.Get(ctx, k)
		if ok {
			if !f.hasPriority(flag.Source, source) {
				logger.Debug(
					fmt.Sprintf(
						"not deleting: flag %s from source %s cannot be deleted by %s",
						k,
						flag.Source,
						source,
					),
				)
				continue
			}
			notifications[k] = map[string]interface{}{
				"type":   string(model.NotificationDelete),
				"source": source,
			}
			f.Delete(k)
		} else {
			logger.Warn(
				fmt.Sprintf("failed to remove flag, flag with key %s from source %s does not exist.",
					k,
					source))
		}
	}
	return notifications
}

// Merge provided flags from source with currently stored flags.
// nolint: funlen
func (f *State) Merge(
	logger *logger.Logger,
	source string,
	selector string,
	flags map[string]model.Flag,
	metadata model.Metadata,
) (map[string]interface{}, bool) {
	notifications := map[string]interface{}{}
	resyncRequired := false
	f.mx.Lock()
	f.setSourceMetadata(source, metadata)

	for k, v := range f.Flags {
		if v.Source == source && v.Selector == selector {
			if _, ok := flags[k]; !ok {
				// flag has been deleted
				delete(f.Flags, k)
				notifications[k] = map[string]interface{}{
					"type":   string(model.NotificationDelete),
					"source": source,
				}
				resyncRequired = true
				logger.Debug(
					fmt.Sprintf(
						"store resync triggered: flag %s has been deleted from source %s",
						k, source,
					),
				)
				continue
			}
		}
	}
	f.mx.Unlock()
	for k, newFlag := range flags {
		newFlag.Source = source
		newFlag.Selector = selector
		storedFlag, _, ok := f.Get(context.Background(), k)
		if ok {
			if !f.hasPriority(storedFlag.Source, source) {
				logger.Debug(
					fmt.Sprintf(
						"not merging: flag %s from source %s does not have priority over %s",
						k, source, storedFlag.Source,
					),
				)
				continue
			}
			if reflect.DeepEqual(storedFlag, newFlag) {
				continue
			}
		}
		if !ok {
			notifications[k] = map[string]interface{}{
				"type":   string(model.NotificationCreate),
				"source": source,
			}
		} else {
			notifications[k] = map[string]interface{}{
				"type":   string(model.NotificationUpdate),
				"source": source,
			}
		}
		// Store the new version of the flag
		f.Set(k, newFlag)
	}
	return notifications, resyncRequired
}

func (f *State) getMetadataForSource(source string) model.Metadata {
	perSource, ok := f.MetadataPerSource[source]
	if ok && perSource != nil {
		return maps.Clone(perSource)
	}
	return model.Metadata{}
}

func (f *State) getMetadata() model.Metadata {
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

func (f *State) setSourceMetadata(source string, metadata model.Metadata) {
	if f.MetadataPerSource == nil {
		f.MetadataPerSource = map[string]model.Metadata{}
	}

	f.MetadataPerSource[source] = metadata
}
