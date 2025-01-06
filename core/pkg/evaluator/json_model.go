package evaluator

import (
	"encoding/json"

	"github.com/open-feature/flagd/core/pkg/model"
)

type Evaluators struct {
	Evaluators map[string]json.RawMessage `json:"$evaluators"`
}

type Metadata struct {
	FlagSetID      string `json:"flagSetId"`
	FlagSetVersion string `json:"flagSetVersion"`
}

type ConfigWithMetadata struct {
	Flags    map[string]model.Flag `json:"flags"`
	MetaData Metadata              `json:"metadata"`
}

type Flags struct {
	Flags map[string]model.Flag `json:"flags"`
}
