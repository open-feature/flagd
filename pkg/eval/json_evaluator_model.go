package eval

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

	"github.com/open-feature/flagd/pkg/logger"
)

type Flag struct {
	State          string          `json:"state"`
	DefaultVariant string          `json:"defaultVariant"`
	Variants       map[string]any  `json:"variants"`
	Targeting      json.RawMessage `json:"targeting,omitempty"`
	Source         string          `json:"source"`
}

type Evaluators struct {
	Evaluators map[string]json.RawMessage `json:"$evaluators"`
}

type Flags struct {
	mx    *sync.RWMutex
	Flags map[string]Flag `json:"flags"`
}

// Add new flags from source. The implementation is not thread safe
func (f Flags) Add(logger *logger.Logger, source string, ff Flags) map[string]interface{} {
	notifications := map[string]interface{}{}

	for k, newFlag := range ff.Flags {
		f.mx.RLock()
		storedFlag, ok := f.Flags[k]
		f.mx.RUnlock()
		if ok && storedFlag.Source != source {
			logger.Warn(fmt.Sprintf(
				"flag with key %s from source %s already exist, overriding this with flag from source %s",
				k,
				storedFlag.Source,
				source,
			))
		}

		notifications[k] = map[string]interface{}{
			"type":   string(NotificationCreate),
			"source": source,
		}

		// Store the new version of the flag
		newFlag.Source = source
		f.mx.Lock()
		f.Flags[k] = newFlag
		f.mx.Unlock()
	}

	return notifications
}

// Update existing flags from source. The implementation is not thread safe
func (f Flags) Update(logger *logger.Logger, source string, ff Flags) map[string]interface{} {
	notifications := map[string]interface{}{}

	for k, flag := range ff.Flags {
		f.mx.RLock()
		storedFlag, ok := f.Flags[k]
		f.mx.RUnlock()
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
			"type":   string(NotificationUpdate),
			"source": source,
		}

		flag.Source = source
		f.mx.Lock()
		f.Flags[k] = flag
		f.mx.Unlock()
	}

	return notifications
}

// Delete matching flags from source. The implementation is not thread safe
func (f Flags) Delete(logger *logger.Logger, source string, ff Flags) map[string]interface{} {
	notifications := map[string]interface{}{}

	for k := range ff.Flags {
		f.mx.RLock()
		_, ok := f.Flags[k]
		f.mx.RUnlock()
		if ok {
			notifications[k] = map[string]interface{}{
				"type":   string(NotificationDelete),
				"source": source,
			}

			f.mx.Lock()
			delete(f.Flags, k)
			f.mx.Unlock()
		} else {
			logger.Warn(
				fmt.Sprintf("failed to remove flag, flag with key %s from source %s does not exisit.",
					k,
					source))
		}
	}

	return notifications
}

// Merge provided flags from source with currently stored flags. The implementation is not thread safe
func (f Flags) Merge(logger *logger.Logger, source string, ff Flags) map[string]interface{} {
	notifications := map[string]interface{}{}

	f.mx.Lock()
	for k, v := range f.Flags {
		if v.Source == source {
			if _, ok := ff.Flags[k]; !ok {
				// flag has been deleted
				delete(f.Flags, k)
				notifications[k] = map[string]interface{}{
					"type":   string(NotificationDelete),
					"source": source,
				}
				continue
			}
		}
	}
	f.mx.Unlock()

	for k, newFlag := range ff.Flags {
		newFlag.Source = source

		f.mx.RLock()
		storedFlag, ok := f.Flags[k]
		f.mx.RUnlock()
		if !ok {
			notifications[k] = map[string]interface{}{
				"type":   string(NotificationCreate),
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
				"type":   string(NotificationUpdate),
				"source": source,
			}
		}

		f.mx.Lock()
		// Store the new version of the flag
		f.Flags[k] = newFlag
		f.mx.Unlock()
	}

	return notifications
}
