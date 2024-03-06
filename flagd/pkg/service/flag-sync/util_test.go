package sync

import (
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
)

// getSimpleFlagStore returns a flag store pre-filled with flags from sources A & B
func getSimpleFlagStore() (*store.Flags, []string) {
	variants := map[string]any{
		"true":  true,
		"false": false,
	}

	flagStore := store.NewFlags()

	flagStore.Set("flagA", model.Flag{
		State:          "ENABLED",
		DefaultVariant: "false",
		Variants:       variants,
		Source:         "A",
	})

	flagStore.Set("flagB", model.Flag{
		State:          "ENABLED",
		DefaultVariant: "true",
		Variants:       variants,
		Source:         "B",
	})

	return flagStore, []string{"A", "B"}
}
