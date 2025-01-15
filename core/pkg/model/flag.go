package model

import "encoding/json"

type Flag struct {
	State          string                 `json:"state"`
	DefaultVariant string                 `json:"defaultVariant"`
	Variants       map[string]any         `json:"variants"`
	Targeting      json.RawMessage        `json:"targeting,omitempty"`
	Source         string                 `json:"source"`
	Selector       string                 `json:"selector"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

type Evaluators struct {
	Evaluators map[string]json.RawMessage `json:"$evaluators"`
}
