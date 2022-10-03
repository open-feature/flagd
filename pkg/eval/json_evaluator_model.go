package eval

import (
	"encoding/json"
	"fmt"
	"reflect"
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
		existing, ok := result.Flags[k]
		if !ok {
			// flag has been set as new
			fmt.Printf("flag value set %s with source %s\n", k, source)
			result.Flags[k] = v
		} else {
			if !reflect.DeepEqual(existing, v) {
				// flag has been updated
				fmt.Printf("flag value updated %s with source %s\n", k, source)
				result.Flags[k] = v
			}
		}
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
