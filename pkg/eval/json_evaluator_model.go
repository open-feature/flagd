package eval

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/open-feature/flagd/pkg/logger"
)

type Flags struct {
	Flags map[string]Flag `json:"flags"`
}

type Evaluators struct {
	Evaluators map[string]json.RawMessage `json:"$evaluators"`
}

func (f Flags) Merge(logger *logger.Logger, source string, ff Flags) (Flags, map[string]interface{}) {
	notifications := map[string]interface{}{}
	result := Flags{Flags: make(map[string]Flag)}
	for k, v := range f.Flags {
		if v.Source == source {
			if _, ok := ff.Flags[k]; !ok {
				// flag has been deleted
				notifications[k] = map[string]interface{}{
					"type":   string(NotificationDelete),
					"source": source,
				}
				continue
			}
		}
		result.Flags[k] = v
	}
	for k, v := range ff.Flags {
		v.Source = source
		val, ok := result.Flags[k]
		if !ok {
			notifications[k] = map[string]interface{}{
				"type":   string(NotificationCreate),
				"source": source,
			}
		} else if !reflect.DeepEqual(val, v) {
			if val.Source != source {
				logger.Warn(
					fmt.Sprintf(
						"key value %s is duplicated across multiple sources this can lead to unexpected behavior: %s, %s",
						k,
						val.Source,
						source,
					),
				)
			}
			notifications[k] = map[string]interface{}{
				"type":   string(NotificationUpdate),
				"source": source,
			}
		}
		result.Flags[k] = v
	}
	return result, notifications
}

type Flag struct {
	State          string          `json:"state"`
	DefaultVariant string          `json:"defaultVariant"`
	Variants       map[string]any  `json:"variants"`
	Targeting      json.RawMessage `json:"targeting,omitempty"`
	Source         string          `json:"source"`
}
