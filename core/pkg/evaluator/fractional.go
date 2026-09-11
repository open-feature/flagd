package evaluator

import (
    "errors"
    "fmt"
    "math"

    "github.com/fxamacker/cbor/v2"
    "github.com/open-feature/flagd/core/pkg/logger"
    "github.com/twmb/murmur3"
)

var cborEncMode, _ = cbor.CoreDetEncOptions().EncMode()

const maxWeightSum = math.MaxInt32 // 2,147,483,647

const FractionEvaluationName = "fractional"

type Fractional struct {
    Logger *logger.Logger
}

type fractionalEvaluationDistribution struct {
    totalWeight      int32
    weightedVariants []fractionalEvaluationVariant
    data             any
    logger           *logger.Logger
}

type fractionalEvaluationVariant struct {
    variant any // string, bool, number or nil
    weight  int32
}

func (v fractionalEvaluationVariant) getPercentage(totalWeight int32) float64 {
    if totalWeight == 0 {
        return 0
    }

    return 100 * float64(v.weight) / float64(totalWeight)
}

func NewFractional(logger *logger.Logger) *Fractional {
    return &Fractional{Logger: logger}
}

func (fe *Fractional) Evaluate(values, data any) any {
    bytesToDistribute, feDistributions, err := parseFractionalEvaluationData(values, data, fe.Logger)
    if err != nil {
        fe.Logger.Warn(fmt.Sprintf("parse fractional evaluation data: %v", err))
        return nil
    }

    if feDistributions == nil {
        return nil
    }

    hashValue := murmur3.Sum32(bytesToDistribute)
    return distributeValue(hashValue, feDistributions)
}

func normalizeValue(val any) any {
    switch v := val.(type) {
    case float64:
        if math.IsNaN(v) || math.IsInf(v, 0) {
            return v
        }
        if v == math.Trunc(v) {
            if v >= 0 && v <= float64(math.MaxUint64) {
                return uint64(v)
            }
            if v < 0 && v >= float64(math.MinInt64) {
                return int64(v)
            }
        }
        return v
    case float32:
        return normalizeValue(float64(v))
    case int:
        if v >= 0 {
            return uint64(v)
        }
        return int64(v)
    case int64:
        if v >= 0 {
            return uint64(v)
        }
        return v
    case uint:
        return uint64(v)
    case uint32:
        return uint64(v)
    case uint64:
        return v
    case map[string]any:
        res := make(map[string]any, len(v))
        for k, item := range v {
            res[k] = normalizeValue(item)
        }
        return res
    case []any:
        res := make([]any, len(v))
        for i, item := range v {
            res[i] = normalizeValue(item)
        }
        return res
    default:
        return v
    }
}

func encodeDeterministicCBOR(val any) ([]byte, error) {
    normalized := normalizeValue(val)
    return cborEncMode.Marshal(normalized)
}

func parseFractionalEvaluationData(values, data any, logger *logger.Logger) ([]byte, *fractionalEvaluationDistribution, error) {
    valuesArray, ok := values.([]any)
    if !ok {
        return nil, nil, errors.New("fractional evaluation data is not an array")
    }
    if len(valuesArray) < 1 {
        return nil, nil, errors.New("fractional evaluation data must contain at least one distribution")
    }

    dataMap, ok := data.(map[string]any)
    if !ok {
        return nil, nil, errors.New("data isn't of type map[string]any")
    }

    properties, _ := getFlagdProperties(dataMap)
    flagKey := properties.FlagKey

    // If first element evaluates to null/nil, report an error and return nil.
    if valuesArray[0] == nil {
        return nil, nil, fmt.Errorf("flag %q: first element of fractional evaluation data is null", flagKey)
    }

    // If first element is a non-array type, use it as explicit hashing input.
    if _, isArray := valuesArray[0].([]any); !isArray {
        hashingInput := valuesArray[0]
        valuesArray = valuesArray[1:]

        bytesToHash, err := encodeDeterministicCBOR(hashingInput)
        if err != nil {
            return nil, nil, fmt.Errorf("flag %q: failed to encode hashing input: %w", flagKey, err)
        }

        feDistributions, err := parseFractionalEvaluationDistributions(valuesArray, data, logger, flagKey)
        if err != nil {
            return nil, nil, err
        }

        return bytesToHash, feDistributions, nil
    }

    // First element is an array ([]any), meaning no explicit hashing input was provided.
    // We fall back to implicit targetingKey rules.
    rawTargetingKey, exists := dataMap[targetingKeyKey]
    if !exists || rawTargetingKey == nil {
        return nil, nil, fmt.Errorf("flag %q: bucketing value not supplied and no targetingKey in context", flagKey)
    }

    targetingKey, isString := rawTargetingKey.(string)
    if !isString {
        return nil, nil, fmt.Errorf("flag %q: targetingKey is not a string", flagKey)
    }

    if targetingKey == "" {
        return nil, nil, fmt.Errorf("flag %q: targetingKey is empty", flagKey)
    }

    // Build 2-element array [flagKey, targetingKey] and encode to CBOR.
    implicitInput := []any{flagKey, targetingKey}
    bytesToHash, err := encodeDeterministicCBOR(implicitInput)
    if err != nil {
        return nil, nil, fmt.Errorf("flag %q: failed to encode implicit targetingKey: %w", flagKey, err)
    }

    feDistributions, err := parseFractionalEvaluationDistributions(valuesArray, data, logger, flagKey)
    if err != nil {
        return nil, nil, err
    }

    return bytesToHash, feDistributions, nil
}

func parseFractionalEvaluationDistributions(values []any, data any, logger *logger.Logger, flagKey string) (*fractionalEvaluationDistribution, error) {
    feDistributions := &fractionalEvaluationDistribution{
        totalWeight:      0,
        weightedVariants: make([]fractionalEvaluationVariant, len(values)),
        data:             data,
        logger:           logger,
    }

    // parse all weights first to validate the sum
    var totalWeightInt64 int64 = 0

    for i := 0; i < len(values); i++ {
        distributionArray, ok := values[i].([]any)
        if !ok {
            return nil, fmt.Errorf("flag %q: distribution elements aren't of type []any. "+
                "please check your rule in flag definition", flagKey)
        }

        if len(distributionArray) == 0 {
            return nil, fmt.Errorf("flag %q: distribution element needs at least one element", flagKey)
        }

        // JSONLogic pre-evaluates all arguments before they reach fractional.
        // Pre-evaluated operators become primitive values (strings, numbers, etc.), never map[string]any nodes.
        var variant any
        switch v := distributionArray[0].(type) {
        case string:
            variant = v
        case bool:
            variant = v
        case float64:
            variant = v
        case nil:
            variant = nil
        default:
            return nil, fmt.Errorf("flag %q: first element of distribution element must be a string, bool, number, or nil", flagKey)
        }

        weight := int64(1)
        if len(distributionArray) >= 2 {
            // parse as float64 first since that's what JSON gives us
            distributionWeight, ok := distributionArray[1].(float64)
            if !ok && distributionArray[1] != nil {
                return nil, fmt.Errorf("flag %q: weight must be a number", flagKey)
            }
            if ok {
                weight = int64(distributionWeight)
            }
        }

        // validate weight is a whole number
        if len(distributionArray) >= 2 {
            distributionWeight, ok := distributionArray[1].(float64)
            if ok && distributionWeight != float64(int64(distributionWeight)) {
                return nil, fmt.Errorf("flag %q: weights must be integers", flagKey)
            }
        }

        // validate individual weight doesn't exceed int32
        if weight > math.MaxInt32 {
            return nil, fmt.Errorf("flag %q: weight %d exceeds maximum allowed value %d", flagKey, weight, math.MaxInt32)
        }

        // clamp negative weights to 0
        if weight < 0 {
            // negative weights can be the result of rollout calculations, so we log and clamp to 0 rather than returning an error
            logger.Debug(fmt.Sprintf("flag %q: negative weight %d clamped to 0", flagKey, weight))
            weight = 0
        }

        totalWeightInt64 += weight
        feDistributions.weightedVariants[i] = fractionalEvaluationVariant{
            variant: variant,
            weight:  int32(weight),
        }
    }

    // validate total weight doesn't exceed MaxInt32
    if totalWeightInt64 > int64(maxWeightSum) {
        return nil, fmt.Errorf("flag %q: sum of all weights (%d) exceeds maximum allowed value (%d)", flagKey, totalWeightInt64, maxWeightSum)
    }

    feDistributions.totalWeight = int32(totalWeightInt64)
    return feDistributions, nil
}

// distributeValue accepts a pre-computed 32-bit hash value and distributes it to a variant using high-precision integer arithmetic.
// It maps a 32-bit hash to the range [0, totalWeight) and finds the variant bucket that contains that value.
func distributeValue(hashValue uint32, feDistribution *fractionalEvaluationDistribution) any {
    if feDistribution.totalWeight == 0 {
        return nil
    }

    bucket := (uint64(hashValue) * uint64(feDistribution.totalWeight)) >> 32

    var rangeEnd uint64 = 0
    for _, variant := range feDistribution.weightedVariants {
        rangeEnd += uint64(variant.weight)
        if bucket < rangeEnd {
            return variant.variant
        }
    }

    // unreachable given validation
    return nil
}
