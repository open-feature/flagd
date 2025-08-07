package store

import (
	"maps"
	"sort"
	"strings"

	uuid "github.com/google/uuid"
	"github.com/open-feature/flagd/core/pkg/model"
)

// flags table and index constants
const flagsTable = "flags"

const idIndex = "id"
const keyIndex = "key"
const sourceIndex = "source"
const priorityIndex = "priority"
const flagSetIdIndex = "flagSetId"

// compound indices; maintain sub-indexes alphabetically; order matters; these must match what's generated in the SelectorMapToQuery func.
const flagSetIdSourceCompoundIndex = flagSetIdIndex + "+" + sourceIndex
const keySourceCompoundIndex = keyIndex + "+" + sourceIndex
const flagSetIdKeySourceCompoundIndex = flagSetIdIndex + "+" + keyIndex + "+" + sourceIndex

// flagSetId defaults to a UUID generated at startup to make our queries consistent
// any flag without a "flagSetId" is assigned this one; it's never exposed externally
var nilFlagSetId = uuid.New().String()

// A selector represents a set of constraints used to query the store.
type Selector interface {
	SelectorMapToQuery() (indexId string, constraints []interface{})
	SelectorToMetadata() model.Metadata
	IsEmpty() bool
	WithIndex(key string, value string) Selector
}

var _ Selector = selector{}

type selector struct {
	indexMap map[string]string
}

// NewSelector creates a new Selector from a selector expression string.
// For example, to select flags from source "./mySource" and flagSetId "1234", use the expression:
// "source=./mySource,flagSetId=1234"
func NewSelector(selectorExpression string) Selector {
	return selector{
		indexMap: expressionToMap(selectorExpression),
	}
}

func expressionToMap(selector string) map[string]string {
	selectorMap := make(map[string]string)
	if selector == "" {
		return selectorMap
	}

	// Split the selector by commas
	pairs := strings.Split(selector, ",")
	for _, pair := range pairs {
		// Split each pair by the first equal sign
		parts := strings.Split(pair, "=")
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			selectorMap[key] = value
		}
	}
	return selectorMap
}

func (s selector) WithIndex(key string, value string) Selector {
	m := maps.Clone(s.indexMap)
	m[key] = value
	return selector{
		indexMap: m,
	}
}

func (s selector) IsEmpty() bool {
	return len(s.indexMap) == 0
}

// SelectorMapToQuery converts the selector map to an indexId and constraints for querying the store.
// For a given index, a specific order and number of constraints are required.
// Both the indexId and constraints are generated based on the keys present in the selector's internal map.
func (s selector) SelectorMapToQuery() (indexId string, constraints []interface{}) {

	if len(s.indexMap) == 2 && s.indexMap[flagSetIdIndex] != "" && s.indexMap[keyIndex] != "" {
		// special case for flagSetId and key (this is the "id" index)
		return idIndex, []interface{}{s.indexMap[flagSetIdIndex], s.indexMap[keyIndex]}
	}

	qs := []string{}
	keys := make([]string, 0, len(s.indexMap))

	for key := range s.indexMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		indexId += key + "+"
		qs = append(qs, s.indexMap[key])
	}

	indexId = strings.TrimSuffix(indexId, "+")
	// Convert []string to []interface{}
	c := make([]interface{}, 0, len(qs)+1)
	for _, v := range qs {
		c = append(c, v)
	}
	constraints = c

	return indexId, constraints
}

// SelectorToMetadata converts the selector's internal map to metadata for logging or tracing purposes.
// Only includes known indices to avoid leaking sensitive information, and is usually returned as the "top level" metadata
func (s selector) SelectorToMetadata() model.Metadata {
	meta := model.Metadata{}

	if s.indexMap[flagSetIdIndex] != "" {
		meta[flagSetIdIndex] = s.indexMap[flagSetIdIndex]
	}
	if s.indexMap[sourceIndex] != "" {
		meta[sourceIndex] = s.indexMap[sourceIndex]
	}
	return meta
}
