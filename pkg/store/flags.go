package store

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

	"github.com/open-feature/flagd/pkg/logger"

	"github.com/open-feature/flagd/pkg/model"
)

type Flags struct {
	mx    sync.RWMutex
	Flags map[string]model.Flag `json:"flags"`
}

func NewFlags() *Flags {
	return &Flags{Flags: map[string]model.Flag{}}
}

func (f *Flags) Set(key string, flag model.Flag) {
	f.mx.Lock()
	defer f.mx.Unlock()
	f.Flags[key] = flag
}

func (f *Flags) Get(key string) (model.Flag, bool) {
	f.mx.RLock()
	defer f.mx.RUnlock()
	flag, ok := f.Flags[key]

	return flag, ok
}

func (f *Flags) Delete(key string) {
	f.mx.Lock()
	defer f.mx.Unlock()
	delete(f.Flags, key)
}

func (f *Flags) String() (string, error) {
	f.mx.RLock()
	defer f.mx.RUnlock()
	bytes, err := json.Marshal(f)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// GetAll returns a copy of the store's state (copy in order to be concurrency safe)
func (f *Flags) GetAll() map[string]model.Flag {
	f.mx.RLock()
	defer f.mx.RUnlock()
	state := make(map[string]model.Flag, len(f.Flags))

	for key, flag := range f.Flags {
		state[key] = flag
	}

	return state
}

// Add new flags from source.
func (f *Flags) Add(logger *logger.Logger, source string, flags map[string]model.Flag) map[string]interface{} {
	notifications := map[string]interface{}{}

	for k, newFlag := range flags {
		storedFlag, ok := f.Get(k)
		if ok && storedFlag.Source != source {
			logger.Warn(fmt.Sprintf(
				"flag with key %s from source %s already exist, overriding this with flag from source %s",
				k,
				storedFlag.Source,
				source,
			))
		}

		notifications[k] = map[string]interface{}{
			"type":   string(model.NotificationCreate),
			"source": source,
		}

		// Store the new version of the flag
		newFlag.Source = source
		f.Set(k, newFlag)
	}

	return notifications
}

// Update existing flags from source.
func (f *Flags) Update(logger *logger.Logger, source string, flags map[string]model.Flag) map[string]interface{} {
	notifications := map[string]interface{}{}

	for k, flag := range flags {
		storedFlag, ok := f.Get(k)
		if !ok {
			logger.Warn(
				fmt.Sprintf("failed to update the flag, flag with key %s from source %s does not exist.",
					k,
					source))

			continue
		}
		if storedFlag.Source != source {
			logger.Warn(fmt.Sprintf(
				"flag with key %s from source %s already exist, overriding this with flag from source %s",
				k,
				storedFlag.Source,
				source,
			))
		}

		notifications[k] = map[string]interface{}{
			"type":   string(model.NotificationUpdate),
			"source": source,
		}

		flag.Source = source
		f.Set(k, flag)
	}

	return notifications
}

// DeleteFlags matching flags from source.
func (f *Flags) DeleteFlags(logger *logger.Logger, source string, flags map[string]model.Flag) map[string]interface{} {
	notifications := map[string]interface{}{}

	for k := range flags {
		_, ok := f.Get(k)
		if ok {
			notifications[k] = map[string]interface{}{
				"type":   string(model.NotificationDelete),
				"source": source,
			}

			f.Delete(k)
		} else {
			logger.Warn(
				fmt.Sprintf("failed to remove flag, flag with key %s from source %s does not exisit.",
					k,
					source))
		}
	}

	return notifications
}

// Merge provided flags from source with currently stored flags.
func (f *Flags) Merge(logger *logger.Logger, source string, flags map[string]model.Flag) map[string]interface{} {
	notifications := map[string]interface{}{}

	f.mx.Lock()
	for k, v := range f.Flags {
		if v.Source == source {
			if _, ok := flags[k]; !ok {
				// flag has been deleted
				delete(f.Flags, k)
				notifications[k] = map[string]interface{}{
					"type":   string(model.NotificationDelete),
					"source": source,
				}
				continue
			}
		}
	}
	f.mx.Unlock()

	for k, newFlag := range flags {
		newFlag.Source = source

		storedFlag, ok := f.Get(k)
		if !ok {
			notifications[k] = map[string]interface{}{
				"type":   string(model.NotificationCreate),
				"source": source,
			}
		} else if !reflect.DeepEqual(storedFlag, newFlag) {
			if storedFlag.Source != source {
				logger.Warn(
					fmt.Sprintf(
						"key value: %s is duplicated across multiple sources this can lead to unexpected behavior: %s, %s",
						k,
						storedFlag.Source,
						source,
					),
				)
			}
			notifications[k] = map[string]interface{}{
				"type":   string(model.NotificationUpdate),
				"source": source,
			}
		}

		// Store the new version of the flag
		f.Set(k, newFlag)
	}

	return notifications
}
