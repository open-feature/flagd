package evaluator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"go.uber.org/zap"
)

const ResolveEvaluationName = "resolve"

// reservedMetadataKeys belong to flag set metadata; skipped even when computed
var reservedMetadataKeys = map[string]struct{}{"flagSetId": {}, "version": {}}

// resolveResult is returned by the resolve operator
type resolveResult struct {
	variant  string
	null     bool // a null variant falls back like a plain null result
	reason   *string
	metadata map[string]any
}

// ResolveEvaluator implements the terminal resolve operator.
type ResolveEvaluator struct {
	Logger *logger.Logger
}

func NewResolveEvaluator(log *logger.Logger) *ResolveEvaluator {
	return &ResolveEvaluator{Logger: log}
}

// Evaluate returns a *resolveResult, or null for malformed arguments, like other custom operators.
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

	parsed := &resolveResult{null: args[0] == nil}
	if !parsed.null {
		variant, err := resolveVariantKey(args[0])
		if err != nil {
			return nil, err
		}
		parsed.variant = variant
	}

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

// resolveVariantKey coerces the variant like a plain targeting result
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

// NaN and Inf are rejected since JSON can't encode them
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
	ctx context.Context, reqID, flagKey string, flag model.Flag, metadata map[string]any, parsed *resolveResult,
) (variant, reason string, meta map[string]any, err error) {
	if metadata == nil {
		metadata = map[string]any{}
	}
	if parsed.null {
		variant, reason, err = nullResult(ctx, flag)
		if err != nil {
			return variant, reason, metadata, err
		}
	} else {
		// an empty variant means no variant, so it's never a valid key
		if _, exists := flag.Variants[parsed.variant]; !exists || parsed.variant == "" {
			je.Logger.ErrorWithID(reqID, fmt.Sprintf(
				"invalid or missing variant: %s for flagKey: %s, variant is not valid", parsed.variant, flagKey))
			return "", model.ErrorReason, metadata, errors.New(model.GeneralErrorCode)
		}
		variant, reason = parsed.variant, model.TargetingMatchReason
	}
	for key, value := range parsed.metadata {
		if _, reserved := reservedMetadataKeys[key]; reserved {
			je.Logger.DebugWithID(reqID, fmt.Sprintf("ignoring reserved resolve metadata key: %s", key))
			continue
		}
		metadata[key] = value
	}
	if parsed.reason != nil {
		reason = *parsed.reason
	}
	return variant, reason, metadata, nil
}
