package evaluator

import (
	"encoding/json"

	"github.com/open-feature/flagd/core/pkg/model"
)

type Evaluators struct {
	Evaluators map[string]json.RawMessage `json:"$evaluators"`
}

type Flags struct {
	Flags map[string]model.Flag `json:"flags"`
}
