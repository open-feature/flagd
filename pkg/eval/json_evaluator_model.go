package eval

import (
	"encoding/json"
	"fmt"
	"reflect"

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
	Flags map[string]Flag `json:"flags"`
}

// Add new flags from source. The implementation is not thread safe
func (f Flags) Add(logger *logger.Logger, source string, ff Flags) map[string]interface{} {
	notifications := map[string]interface{}{}

	for k, newFlag := range ff.Flags {
		if storedFlag, ok := f.Flags[k]; ok && storedFlag.Source != source {
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
		f.Flags[k] = newFlag
	}

	return notifications
}

// Update existing flags from source. The implementation is not thread safe
func (f Flags) Update(logger *logger.Logger, source string, ff Flags) map[string]interface{} {
	notifications := map[string]interface{}{}

	for k, flag := range ff.Flags {
		if storedFlag, ok := f.Flags[k]; !ok {
			logger.Warn(
				fmt.Sprintf("failed to update the flag, flag with key %s from source %s does not exisit.",
					k,
					source))

			continue
		} else if storedFlag.Source != source {
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
		f.Flags[k] = flag
	}

	return notifications
}

// Delete matching flags from source. The implementation is not thread safe
// If ff.Flags is empty, all flags from the given source are deleted
func (f Flags) Delete(logger *logger.Logger, source string, ff Flags) map[string]interface{} {
	notifications := map[string]interface{}{}

	if len(ff.Flags) == 0 {
		for k, flag := range f.Flags {
			if flag.Source == source {
				notifications[k] = map[string]interface{}{
					"type":   string(NotificationDelete),
					"source": source,
				}
				delete(f.Flags, k)
			}
		}
		return notifications
	}

	for k := range ff.Flags {
		if _, ok := f.Flags[k]; ok {
			notifications[k] = map[string]interface{}{
				"type":   string(NotificationDelete),
				"source": source,
			}

			delete(f.Flags, k)
		} else {
			logger.Warn(
				fmt.Sprintf("failed to remove flag, flag with key %s from source %s does not exist.",
					k,
					source))
		}
	}

	return notifications
}

// Merge provided flags from source with currently stored flags. The implementation is not thread safe
func (f Flags) Merge(logger *logger.Logger, source string, ff Flags) map[string]interface{} {
	notifications := map[string]interface{}{}

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

	for k, newFlag := range ff.Flags {
		newFlag.Source = source

		storedFlag, ok := f.Flags[k]
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

		// Store the new version of the flag
		f.Flags[k] = newFlag
	}

	return notifications
}
