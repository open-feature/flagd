package model

import "encoding/json"

const Key = "Key"
const FlagSetId = "FlagSetId"
const Source = "Source"
const Priority = "Priority"

type Flag struct {
	Key            string          `json:"key"`
	FlagSetId      string          `json:"-"` // not serialized, used only for indexing
	Priority       int             `json:"-"` // not serialized, used only for indexing
	State          string          `json:"state"`
	DefaultVariant string          `json:"defaultVariant"`
	Variants       map[string]any  `json:"variants"`
	Targeting      json.RawMessage `json:"targeting,omitempty"`
	Source         string          `json:"source"`
	Metadata       Metadata        `json:"metadata,omitempty"`
}

type Evaluators struct {
	Evaluators map[string]json.RawMessage `json:"$evaluators"`
}

type Metadata = map[string]interface{}

func (f Flag) MarshalJSON() ([]byte, error) {
	type flagAlias Flag
	return json.Marshal(struct {
		*flagAlias
		Key interface{} `json:"key,omitempty"`
	}{
		flagAlias: (*flagAlias)(&f),
		Key:       nil,
	})
}
