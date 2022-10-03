package eval

import (
	"encoding/json"
)

type Flags struct {
	Flags map[string]Flag `json:"flags"`
}

type Evaluators struct {
	Evaluators map[string]json.RawMessage `json:"$evaluators"`
}

func (f Flags) Merge(source string, ff Flags) Flags {
	result := Flags{Flags: make(map[string]Flag)}
	for k, v := range f.Flags {
		if v.Source == source {
			if _, ok := ff.Flags[k]; !ok {
				// flag has been deleted
				continue
			}
		}
		result.Flags[k] = v
	}
	for k, v := range ff.Flags {
		v.Source = source
		result.Flags[k] = v
	}
	return result
}

type Flag struct {
	State          string          `json:"state"`
	DefaultVariant string          `json:"defaultVariant"`
	Variants       map[string]any  `json:"variants"`
	Targeting      json.RawMessage `json:"targeting"`
	Source         string          `json:"source"`
}
