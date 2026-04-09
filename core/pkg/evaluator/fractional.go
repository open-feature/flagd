package evaluator

import (
	"errors"
	"fmt"
	"math"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/twmb/murmur3"
)

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
	valueToDistribute, feDistributions, err := parseFractionalEvaluationData(values, data, fe.Logger)
	if err != nil {
		fe.Logger.Warn(fmt.Sprintf("parse fractional evaluation data: %v", err))
		return nil
	}

	hashValue := uint32(murmur3.StringSum32(valueToDistribute))
	return distributeValue(hashValue, feDistributions)
}

func parseFractionalEvaluationData(values, data any, logger *logger.Logger) (string, *fractionalEvaluationDistribution, error) {
	valuesArray, ok := values.([]any)
	if !ok {
		return "", nil, errors.New("fractional evaluation data is not an array")
	}
	if len(valuesArray) < 1 {
		return "", nil, errors.New("fractional evaluation data is empty")
	}

	dataMap, ok := data.(map[string]any)
	if !ok {
		return "", nil, errors.New("data isn't of type map[string]any")
	}

	properties, _ := getFlagdProperties(dataMap)
	flagKey := properties.FlagKey

	bucketBy, ok := valuesArray[0].(string)
	if ok {
		valuesArray = valuesArray[1:]
	} else {
		// check for nil here as custom property could be nil/missing
		if valuesArray[0] == nil {
			valuesArray = valuesArray[1:]
		}

		targetingKey, ok := dataMap[targetingKeyKey].(string)
		if !ok {
			return "", nil, fmt.Errorf("flag %q: bucketing value not supplied and no targetingKey in context", flagKey)
		}

		bucketBy = fmt.Sprintf("%s%s", properties.FlagKey, targetingKey)
	}

	feDistributions, err := parseFractionalEvaluationDistributions(valuesArray, data, logger, flagKey)
	if err != nil {
		return "", nil, err
	}

	return bucketBy, feDistributions, nil
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
		return ""
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
	return ""
}
