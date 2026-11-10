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

// resolveResult is what the operator returns. JsonLogic passes it through if/?: untouched (it only
// re-applies maps), and since no JSON input decodes to it, context data cannot forge one.
// A nil reason keeps TARGETING_MATCH. err is set when evaluated arguments are malformed.
type resolveResult struct {
	variant  string
	reason   *string
	metadata map[string]any
	err      error
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
// Invalid evaluated arguments produce an error result rather than null, because a null
// targeting result would silently fall through to the default variant.
func (re *ResolveEvaluator) Evaluate(values, _ any) any {
	parsed, err := parseEvaluatedResolve(values)
	if err != nil {
		re.Logger.Debug("resolve evaluation failed", zap.Error(err))
		return &resolveResult{err: err}
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

// metadataPrimitive rejects NaN and ±Inf (e.g. a computed 1/0): JSON can't encode them, so OFREP and Connect
// JSON responses would fail, taking a whole bulk response with them
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
	if parsed.err != nil {
		je.Logger.ErrorWithID(reqID, fmt.Sprintf("resolve failed for flag %s: %s", flagKey, parsed.err))
		return "", model.ErrorReason, metadata, errors.New(model.ParseErrorCode)
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

// markInvalidResolveFlags marks flags that use resolve outside a terminal position or with malformed literal
// arguments. Only the offending flag is invalid; the rest of the payload still loads.
// Terminal means the operator is the whole targeting rule, or a result branch of if / ?: whose own
// result is terminal. Its value must reach the resolver unchanged.
func (je *JSON) markInvalidResolveFlags(definition *Definition) {
	for i := range definition.Flags {
		flag := &definition.Flags[i]
		if err := validateResolveTargeting(*flag); err != nil {
			flag.Invalid = err.Error()
			je.Logger.Warn(fmt.Sprintf("%s; evaluating the flag returns %s", err, model.ParseErrorCode))
		}
	}
}

func validateResolveTargeting(flag model.Flag) error {
	if len(flag.Targeting) == 0 || string(flag.Targeting) == "null" {
		return nil
	}
	var node any
	if err := json.Unmarshal(flag.Targeting, &node); err != nil {
		return fmt.Errorf("flag %q: invalid targeting: %w", flag.Key, err)
	}
	return walkResolve(node, true, flag.Key)
}

func walkResolve(node any, terminal bool, flagKey string) error {
	rule, ok := node.(map[string]any)
	if !ok {
		// Arrays do not pass a nested resolve through as the targeting result.
		items, isArray := node.([]any)
		if !isArray {
			return nil
		}
		for _, item := range items {
			if err := walkResolve(item, false, flagKey); err != nil {
				return err
			}
		}
		return nil
	}
	if len(rule) != 1 {
		for _, value := range rule {
			if err := walkResolve(value, false, flagKey); err != nil {
				return err
			}
		}
		return nil
	}

	operator, args := onlyRuleEntry(rule)
	if operator == ResolveEvaluationName {
		if !terminal {
			return fmt.Errorf("flag %q: resolve is only valid as the targeting result or an if/?: result", flagKey)
		}
		return validateResolveArguments(args, flagKey)
	}
	if operator == "if" || operator == "?:" {
		return walkConditional(args, terminal, flagKey)
	}
	return walkResolve(args, false, flagKey)
}

func walkConditional(args any, parentTerminal bool, flagKey string) error {
	clauses, ok := args.([]any)
	if !ok {
		return walkResolve(args, false, flagKey)
	}
	length := len(clauses)
	for i, clause := range clauses {
		terminal := false
		// if/?: results are the then-branches (odd indexes) and, when the length is odd, the final else.
		// A single clause is only a condition: a falsy value becomes null, so it is not terminal.
		if length > 1 && (i%2 == 1 || (length%2 == 1 && i == length-1)) {
			terminal = parentTerminal
		}
		if err := walkResolve(clause, terminal, flagKey); err != nil {
			return err
		}
	}
	return nil
}

func validateResolveArguments(args any, flagKey string) error {
	items, ok := args.([]any)
	if !ok || len(items) < 1 || len(items) > 3 {
		return fmt.Errorf("flag %q: resolve arguments must be an array of 1 to 3 items", flagKey)
	}
	if err := validateResolveVariant(items[0], flagKey); err != nil {
		return err
	}
	if len(items) >= 2 {
		if err := validateResolveReason(items[1], flagKey); err != nil {
			return err
		}
	}
	if len(items) == 3 {
		return validateResolveMetadata(items[2], flagKey)
	}
	return nil
}

func validateResolveVariant(node any, flagKey string) error {
	if isJSONLogicRule(node) {
		return walkResolve(node, false, flagKey)
	}
	if _, ok := node.(string); ok {
		return nil
	}
	return fmt.Errorf("flag %q: resolve variant must be a string", flagKey)
}

func validateResolveReason(node any, flagKey string) error {
	if node == nil {
		return nil
	}
	if isJSONLogicRule(node) {
		return walkResolve(node, false, flagKey)
	}
	if _, ok := node.(string); ok {
		return nil
	}
	return fmt.Errorf("flag %q: resolve reason must be a string or null", flagKey)
}

func validateResolveMetadata(node any, flagKey string) error {
	if isJSONLogicRule(node) {
		return walkResolve(node, false, flagKey)
	}
	pairs, ok := node.([]any)
	if !ok {
		return fmt.Errorf("flag %q: resolve metadata must be an array of [key, value] pairs", flagKey)
	}
	for _, pair := range pairs {
		if isJSONLogicRule(pair) {
			if err := walkResolve(pair, false, flagKey); err != nil {
				return err
			}
			continue
		}
		if err := validateMetadataPair(pair, flagKey); err != nil {
			return err
		}
	}
	return nil
}

func validateMetadataPair(pair any, flagKey string) error {
	items, ok := pair.([]any)
	if !ok || len(items) != 2 {
		return fmt.Errorf("flag %q: resolve metadata pair must be a [key, value] array", flagKey)
	}
	if err := validateMetadataKey(items[0], flagKey); err != nil {
		return err
	}
	return validateMetadataValue(items[1], flagKey)
}

func validateMetadataKey(node any, flagKey string) error {
	if isJSONLogicRule(node) {
		return walkResolve(node, false, flagKey)
	}
	if _, ok := node.(string); ok {
		return nil
	}
	return fmt.Errorf("flag %q: resolve metadata key must be a string", flagKey)
}

func validateMetadataValue(node any, flagKey string) error {
	if isJSONLogicRule(node) {
		return walkResolve(node, false, flagKey)
	}
	switch node.(type) {
	case string, float64, bool:
		return nil
	default:
		return fmt.Errorf("flag %q: resolve metadata value must be a string, number, or boolean", flagKey)
	}
}

func isJSONLogicRule(node any) bool {
	rule, ok := node.(map[string]any)
	if !ok || len(rule) != 1 {
		return false
	}
	operator, _ := onlyRuleEntry(rule)
	return knownTargetingOperator(operator)
}

// knownTargetingOperators are the operators JsonLogic will call.
// Any other single-key object is data, not a computed argument.
var knownTargetingOperators = map[string]struct{}{
	ResolveEvaluationName: {}, FractionEvaluationName: {},
	StartsWithEvaluationName: {}, EndsWithEvaluationName: {}, SemVerEvaluationName: {},
	"and": {}, "or": {}, "filter": {}, "map": {}, "reduce": {}, "all": {}, "none": {}, "some": {},
	"in": {}, "missing": {}, "missing_some": {}, "var": {}, "set": {}, "cat": {}, "substr": {}, "merge": {},
	"if": {}, "?:": {}, "max": {}, "min": {}, "+": {}, "-": {}, "*": {}, "/": {}, "%": {}, "abs": {},
	"!": {}, "!!": {}, "===": {}, "!==": {}, "==": {}, "!=": {}, "<": {}, "<=": {}, ">": {}, ">=": {},
	"contains_all": {}, "contains_any": {}, "contains_none": {},
}

func knownTargetingOperator(operator string) bool {
	_, ok := knownTargetingOperators[operator]
	return ok
}

func onlyRuleEntry(rule map[string]any) (string, any) {
	for operator, args := range rule {
		return operator, args
	}
	return "", nil
}
