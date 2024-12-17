package model

import "encoding/json"

type Flag struct {
	State          string          `json:"state"`
	DefaultVariant string          `json:"defaultVariant"`
	Variants       map[string]any  `json:"variants"`
	Targeting      json.RawMessage `json:"targeting,omitempty"`
	Source         string          `json:"source"`
	Selector       string          `json:"selector"`
	Metadata       Metadata        `json:"metadata"`
}

type Metadata struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

type Evaluators struct {
	Evaluators map[string]json.RawMessage `json:"$evaluators"`
}
