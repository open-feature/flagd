package notifications

import (
	"testing"

	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestNewFromFlags(t *testing.T) {
	flagA := model.Flag{
		Key:            "flagA",
		State:          "ENABLED",
		DefaultVariant: "on",
		Source:         "source1",
	}
	flagAUpdated := model.Flag{
		Key:            "flagA",
		State:          "DISABLED",
		DefaultVariant: "on",
		Source:         "source1",
	}
	flagB := model.Flag{
		Key:            "flagB",
		State:          "ENABLED",
		DefaultVariant: "off",
		Source:         "source1",
	}

	tests := []struct {
		name     string
		oldFlags map[string]model.Flag
		newFlags map[string]model.Flag
		want     Notifications
	}{
		{
			name:     "flag added",
			oldFlags: map[string]model.Flag{},
			newFlags: map[string]model.Flag{"flagA": flagA},
			want: Notifications{
				"flagA": map[string]interface{}{
					"type": string(model.NotificationCreate),
				},
			},
		},
		{
			name:     "flag deleted",
			oldFlags: map[string]model.Flag{"flagA": flagA},
			newFlags: map[string]model.Flag{},
			want: Notifications{
				"flagA": map[string]interface{}{
					"type": string(model.NotificationDelete),
				},
			},
		},
		{
			name:     "flag changed",
			oldFlags: map[string]model.Flag{"flagA": flagA},
			newFlags: map[string]model.Flag{"flagA": flagAUpdated},
			want: Notifications{
				"flagA": map[string]interface{}{
					"type": string(model.NotificationUpdate),
				},
			},
		},
		{
			name:     "flag unchanged",
			oldFlags: map[string]model.Flag{"flagA": flagA},
			newFlags: map[string]model.Flag{"flagA": flagA},
			want:     Notifications{},
		},
		{
			name: "mixed changes",
			oldFlags: map[string]model.Flag{
				"flagA": flagA,
				"flagB": flagB,
			},
			newFlags: map[string]model.Flag{
				"flagA": flagAUpdated, // updated
				"flagC": flagA,        // added
			},
			want: Notifications{
				"flagA": map[string]interface{}{
					"type": string(model.NotificationUpdate),
				},
				"flagB": map[string]interface{}{
					"type": string(model.NotificationDelete),
				},
				"flagC": map[string]interface{}{
					"type": string(model.NotificationCreate),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFromFlags(tt.oldFlags, tt.newFlags)
			assert.Equal(t, tt.want, got)
		})
	}
}
