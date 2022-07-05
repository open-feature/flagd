package eval

import (
	"encoding/json"
)

type Flags struct {
	Flags map[string]Flag `json:"flags"`
}

// we could make this more type-safe with generics when we upgrade to 1.18.
type Flag struct {
	State          string                 `json:"state"`
	DefaultVariant string                 `json:"defaultVariant"`
	Variants       map[string]interface{} `json:"variants"`
	Targeting      json.RawMessage        `json:"targeting"`
}
