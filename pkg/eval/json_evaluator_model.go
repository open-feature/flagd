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

func (f Flags) Merge(logger *logger.Logger, source string, ff Flags) (Flags, []StateChangeNotification) {
	notifications := []StateChangeNotification{}
	result := Flags{Flags: make(map[string]Flag)}
	for k, v := range f.Flags {
		if v.Source == source {
			if _, ok := ff.Flags[k]; !ok {
				// flag has been deleted
				notifications = append(notifications, StateChangeNotification{
					Type:    NotificationDelete,
					Source:  source,
					FlagKey: k,
				})
				continue
			}
		}
		result.Flags[k] = v
	}
	for k, v := range ff.Flags {
		v.Source = source
		val, ok := result.Flags[k]
		if !ok {
			notifications = append(notifications, StateChangeNotification{
				Type:    NotificationCreate,
				Source:  source,
				FlagKey: k,
			})
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
			notifications = append(notifications, StateChangeNotification{
				Type:    NotificationUpdate,
				Source:  source,
				FlagKey: k,
			})
		}
		result.Flags[k] = v
	}
	return result, notifications
}

type Flag struct {
	State          string          `json:"state"`
	DefaultVariant string          `json:"defaultVariant"`
	Variants       map[string]any  `json:"variants"`
	Targeting      json.RawMessage `json:"targeting"`
	Source         string          `json:"source"`
}
