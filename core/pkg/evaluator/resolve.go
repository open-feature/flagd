package evaluator

import (
	"errors"
	"fmt"
	"math"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"go.uber.org/zap"
)

const ResolveEvaluationName = "resolve"

// resolveResult is what the operator returns. JsonLogic passes it through if/?: untouched (it only
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

func parseEvaluatedResolve(values any) (*resolveResult, error) {
	args, ok := values.([]any)
	if !ok || len(args) < 1 || len(args) > 3 {
		return nil, errors.New("resolve arguments must be an array of 1 to 3 items")
	}

	variant, ok := args[0].(string)
	if !ok {
		return nil, errors.New("resolve variant must be a string")
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
		metadata[key] = value
	}
	reason = model.TargetingMatchReason
	if parsed.reason != nil {
		reason = *parsed.reason
	}
	return parsed.variant, reason, metadata, nil
}
