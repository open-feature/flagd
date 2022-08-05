package eval

import (
	"encoding/json"
)

type Flags struct {
	Flags map[string]Flag `json:"flags"`
}

type Evaluators struct {
	Evaluators map[string]json.RawMessage `json:"evaluators"`
}

func (f Flags) Merge(ff Flags) Flags {
	result := Flags{Flags: make(map[string]Flag)}
	for k, v := range f.Flags {
		result.Flags[k] = v
	}
	for k, v := range ff.Flags {
		result.Flags[k] = v
	}
	return result
}

type Flag struct {
	State          string          `json:"state"`
	DefaultVariant string          `json:"defaultVariant"`
	Variants       map[string]any  `json:"variants"`
	Targeting      json.RawMessage `json:"targeting"`
}
