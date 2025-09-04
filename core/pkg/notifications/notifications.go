package notifications

import (
	"reflect"

	"github.com/open-feature/flagd/core/pkg/model"
)

const typeField = "type"

// Use to represent change notifications for  mode PROVIDER_CONFIGURATION_CHANGE events.
type Notifications map[string]any

// Generate notifications (deltas) from old and new flag sets for use in RPC mode PROVIDER_CONFIGURATION_CHANGE events.
func NewFromFlags(oldFlags, newFlags map[string]model.Flag) Notifications {
	notifications := map[string]interface{}{}

	// NOTE: we do not care about the events here for flagd itself, those are only needed for the provider, except for the openTelemetry information!!!!
	// flags removed
	for key := range oldFlags {
		if _, ok := newFlags[key]; !ok {
			notifications[key] = map[string]interface{}{
				typeField: string(model.NotificationDelete),
			}
		}
	}

	// flags added or modified
	for key, newFlag := range newFlags {
		oldFlag, exists := oldFlags[key]
		if !exists {
			notifications[key] = map[string]interface{}{
				typeField: string(model.NotificationCreate),
			}
		} else if !flagsEqual(oldFlag, newFlag) {
			notifications[key] = map[string]interface{}{
				typeField: string(model.NotificationUpdate),
			}
		}
	}

	return notifications
}

func flagsEqual(a, b model.Flag) bool {
	return a.State == b.State &&
		a.DefaultVariant == b.DefaultVariant &&
		reflect.DeepEqual(a.Variants, b.Variants) &&
		reflect.DeepEqual(a.Targeting, b.Targeting) &&
		a.Source == b.Source &&
		reflect.DeepEqual(a.Metadata, b.Metadata)
}
