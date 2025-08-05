package utils

import (
	"reflect"

	"github.com/open-feature/flagd/core/pkg/model"
)

func BuildNotifications(oldFlags, newFlags map[string]model.Flag) map[string]interface{} {
	notifications := map[string]interface{}{}

	// flags removed
	for key := range oldFlags {
		if _, ok := newFlags[key]; !ok {
			notifications[key] = map[string]interface{}{
				"type": string(model.NotificationDelete),
			}
		}
	}

	// flags added or modified
	for key, newFlag := range newFlags {
		oldFlag, exists := oldFlags[key]
		if !exists {
			notifications[key] = map[string]interface{}{
				"type": string(model.NotificationCreate),
			}
		} else if !flagsEqual(oldFlag, newFlag) {
			notifications[key] = map[string]interface{}{
				"type": string(model.NotificationUpdate),
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
		a.Selector == b.Selector &&
		reflect.DeepEqual(a.Metadata, b.Metadata)
}
