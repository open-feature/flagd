package evaluator

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"go.uber.org/zap"
)

const ResolveEvaluationName = "resolve"

// resolveOperatorName is the operator the rule's own resolve calls are bound to before evaluation.
// JsonLogic re-applies a map that some operators (if, ?:, reduce...) get back from a branch, so a context value like
// {"resolve": [...]} is dispatched as an operator call. The invalid UTF-8 byte keeps that from reaching this
// operator: json.Unmarshal replaces it with U+FFFD, so no decoded context key can equal this name.
const resolveOperatorName = ResolveEvaluationName + "\xff"

// reservedMetadataKeys belong to flag set metadata. The schema rejects them as literal rule metadata keys; this
// also skips computed ones.
var reservedMetadataKeys = map[string]struct{}{"flagSetId": {}, "version": {}}

// resolveResult is what the operator returns. JsonLogic passes it through if untouched (it only
// re-applies maps), and since no JSON input decodes to it, context data cannot forge one.
// A nil reason keeps TARGETING_MATCH.
type resolveResult struct {
	variant  string
	reason   *string
	metadata map[string]any
}

// ResolveEvaluator implements the terminal resolve operator.
// It selects a variant and optionally overrides the evaluation reason and metadata.
type ResolveEvaluator struct {
	Logger *logger.Logger
}

func NewResolveEvaluator(log *logger.Logger) *ResolveEvaluator {
	return &ResolveEvaluator{Logger: log}
}

// Evaluate returns a *resolveResult. Arguments are already evaluated by JsonLogic.
// Malformed evaluated arguments return null, like the other custom operators, so the flag falls back to its
// default variant.
func (re *ResolveEvaluator) Evaluate(values, _ any) any {
	parsed, err := parseEvaluatedResolve(values)
	if err != nil {
		re.Logger.Debug("resolve evaluation failed", zap.Error(err))
		return nil
	}
	return parsed
}

// rejectFromContext handles "resolve" calls that are left after bindResolveOperator, which can only come from
// context data. It returns null, so the flag falls back to its default variant.
func (re *ResolveEvaluator) rejectFromContext(_, _ any) any {
	re.Logger.Debug("ignoring resolve operator supplied by evaluation context")
	return nil
}

// bindResolveOperator renames the rule's resolve calls to resolveOperatorName, in place. Only single-key maps are
// operator calls; JsonLogic returns larger maps as literal data without evaluating them, so their contents are left
// unchanged. Call it only on a freshly decoded rule, never on context data.
func bindResolveOperator(rule any) any {
	switch node := rule.(type) {
	case map[string]any:
		if len(node) != 1 {
			return rule
		}
		if args, ok := node[ResolveEvaluationName]; ok {
			delete(node, ResolveEvaluationName)
			node[resolveOperatorName] = args
		}
		for _, value := range node {
			bindResolveOperator(value)
		}
	case []any:
		for _, item := range node {
			bindResolveOperator(item)
		}
	}
	return rule
}

func parseEvaluatedResolve(values any) (*resolveResult, error) {
	args, ok := values.([]any)
	if !ok || len(args) < 1 || len(args) > 3 {
		return nil, errors.New("resolve arguments must be an array of 1 to 3 items")
	}

	variant, err := resolveVariantKey(args[0])
	if err != nil {
		return nil, err
	}
	parsed := &resolveResult{variant: variant}

	if len(args) >= 2 && args[1] != nil {
		reason, ok := args[1].(string)
		if !ok {
			return nil, errors.New("resolve reason must be a string or null")
		}
		parsed.reason = &reason
	}

	if len(args) == 3 && args[2] != nil {
		metadata, err := parseEvaluatedMetadata(args[2])
		if err != nil {
			return nil, err
		}
		parsed.metadata = metadata
	}
	return parsed, nil
}

// resolveVariantKey coerces the variant like a plain targeting result: booleans become "true"/"false", and numbers
// take their JSON form (unofficially, as for plain results)
func resolveVariantKey(value any) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case bool, float64, float32, int, int32, int64, uint, uint32, uint64:
		encoded, err := json.Marshal(v)
		if err != nil {
			return "", fmt.Errorf("resolve variant %v has no variant key: %w", v, err)
		}
		return string(encoded), nil
	default:
		return "", errors.New("resolve variant must be a string, boolean, or number")
	}
}

func parseEvaluatedMetadata(value any) (map[string]any, error) {
	pairs, ok := value.([]any)
	if !ok {
		return nil, errors.New("resolve metadata must be an array of [key, value] pairs")
	}
	metadata := make(map[string]any, len(pairs))
	for _, pair := range pairs {
		key, raw, err := metadataPair(pair)
		if err != nil {
			return nil, err
		}
		primitive, ok := metadataPrimitive(raw)
		if !ok {
			return nil, errors.New("resolve metadata value must be a string, finite number, or boolean")
		}
		metadata[key] = primitive
	}
	if len(metadata) == 0 {
		return nil, nil
	}
	return metadata, nil
}

func metadataPair(pair any) (string, any, error) {
	items, ok := pair.([]any)
	if !ok || len(items) != 2 {
		return "", nil, errors.New("resolve metadata pair must be a [key, value] array")
	}
	key, ok := items[0].(string)
	if !ok {
		return "", nil, errors.New("resolve metadata key must be a string")
	}
	return key, items[1], nil
}

// metadataPrimitive rejects NaN and ±Inf (e.g. a computed 1/0): JSON can't encode them, so returning them would
// fail OFREP and Connect JSON responses, taking a whole bulk response with them
func metadataPrimitive(value any) (any, bool) {
	switch n := value.(type) {
	case string, bool:
		return n, true
	case float64:
		return n, isFinite(n)
	case float32:
		return float64(n), isFinite(float64(n))
	case int:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint32:
		return float64(n), true
	case uint64:
		return float64(n), true
	default:
		return nil, false
	}
}

func isFinite(f float64) bool {
	return !math.IsNaN(f) && !math.IsInf(f, 0)
}

func (je *Resolver) applyResolveResult(
	reqID, flagKey string, flag model.Flag, metadata map[string]any, parsed *resolveResult,
) (variant, reason string, meta map[string]any, err error) {
	if metadata == nil {
		metadata = map[string]any{}
	}
	// variant keys are non-empty, and an empty variant is how flagd marks its own fallback and disabled
	// results, so an empty one is always missing
	if _, exists := flag.Variants[parsed.variant]; !exists || parsed.variant == "" {
		je.Logger.ErrorWithID(reqID, fmt.Sprintf(
			"invalid or missing variant: %s for flagKey: %s, variant is not valid", parsed.variant, flagKey))
		return "", model.ErrorReason, metadata, errors.New(model.GeneralErrorCode)
	}
	for key, value := range parsed.metadata {
		if _, reserved := reservedMetadataKeys[key]; reserved {
			je.Logger.DebugWithID(reqID, fmt.Sprintf("ignoring reserved resolve metadata key: %s", key))
			continue
		}
		metadata[key] = value
	}
	reason = model.TargetingMatchReason
	if parsed.reason != nil {
		reason = *parsed.reason
	}
	return parsed.variant, reason, metadata, nil
}
