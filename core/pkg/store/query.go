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

// A Selector represents a set of constraints used to query the store.
type Selector struct {
	indexMap map[string]string
}

// NewSelector creates a new Selector from a selector expression string.
// #1708 Until we decide on the Selector syntax, only a single key=value pair is supported
// For example, to select flags from source "./mySource" or flagSetId "1234", use the expressions:
// "source=./mySource" or "flagSetId=1234"
func NewSelector(selectorExpression string) Selector {
	return Selector{
		indexMap: expressionToMap(selectorExpression),
	}
}

func expressionToMap(sExp string) map[string]string {
	selectorMap := make(map[string]string)
	if sExp == "" {
		return selectorMap
	}

	delimiterIdx := strings.Index(sExp, "=")
	if delimiterIdx == -1 {
		// if no '=' is found, treat the whole string as source (backwards compatibility)
		// we may support interpreting this as a flagSetId in the future as an option
		selectorMap[sourceIndex] = sExp
		return selectorMap
	}

	// split the selector by the first equal sign
	key := sExp[:delimiterIdx]
	value := sExp[delimiterIdx+1:]

	// handle empty flagSetId as nilFlagSetId to query all flags without flagSetId
	if key == "flagSetId" && value == "" {
		value = nilFlagSetId
	}
	selectorMap[key] = value

	return selectorMap
}

// WithIndex creates a new Selector from the current Selector and adds the given key-value-pair
func (s Selector) WithIndex(key string, value string) Selector {
	m := maps.Clone(s.indexMap)
	m[key] = value
	return Selector{
		indexMap: m,
	}
}

func (s *Selector) IsEmpty() bool {
	return s == nil || len(s.indexMap) == 0
}

// ToQuery converts the Selector map to an indexId and constraints for querying the Store.
// For a given index, a specific order and number of constraints are required.
// Both the indexId and constraints are generated based on the keys present in the selector's internal map.
func (s Selector) ToQuery() (indexId string, constraints []interface{}) {

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
	c := make([]interface{}, 0, len(qs))
	for _, v := range qs {
		c = append(c, v)
	}
	constraints = c

	return indexId, constraints
}

// String returns a human-readable representation of the Selector for logging purposes.
func (s *Selector) String() string {
	if s == nil || len(s.indexMap) == 0 {
		return "<empty selector>"
	}

	keys := make([]string, 0, len(s.indexMap))
	for key := range s.indexMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(s.indexMap))
	for _, key := range keys {
		parts = append(parts, key+"="+s.indexMap[key])
	}
	return strings.Join(parts, ", ")
}

// ToMetadata converts the selector's internal map to metadata for logging or tracing purposes.
// Only includes known indices to avoid leaking sensitive information, and is usually returned as the "top level" metadata
func (s *Selector) ToMetadata() model.Metadata {
	meta := model.Metadata{}

	if s == nil || s.indexMap == nil {
		return meta
	}

	if s.indexMap[flagSetIdIndex] != "" {
		meta[flagSetIdIndex] = s.indexMap[flagSetIdIndex]
	}
	if s.indexMap[sourceIndex] != "" {
		meta[sourceIndex] = s.indexMap[sourceIndex]
	}
	return meta
}
