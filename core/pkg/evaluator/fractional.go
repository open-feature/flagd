package evaluator

import (
	"errors"
	"fmt"
	"math"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/twmb/murmur3"
)

const FractionEvaluationName = "fractional"

type Fractional struct {
	Logger *logger.Logger
}

type fractionalEvaluationDistribution struct {
	totalWeight      int
	weightedVariants []fractionalEvaluationVariant
}

type fractionalEvaluationVariant struct {
	variant string
	weight  int
}

func (v fractionalEvaluationVariant) getPercentage(totalWeight int) float64 {
	if totalWeight == 0 {
		return 0
	}

	return 100 * float64(v.weight) / float64(totalWeight)
}

func NewFractional(logger *logger.Logger) *Fractional {
	return &Fractional{Logger: logger}
}

func (fe *Fractional) Evaluate(values, data any) any {
	valueToDistribute, feDistributions, err := parseFractionalEvaluationData(values, data)
	if err != nil {
		fe.Logger.Warn(fmt.Sprintf("parse fractional evaluation data: %v", err))
		return nil
	}

	return distributeValue(valueToDistribute, feDistributions)
}

func parseFractionalEvaluationData(values, data any) (string, *fractionalEvaluationDistribution, error) {
	valuesArray, ok := values.([]any)
	if !ok {
		return "", nil, errors.New("fractional evaluation data is not an array")
	}
	if len(valuesArray) < 2 {
		return "", nil, errors.New("fractional evaluation data has length under 2")
	}

	dataMap, ok := data.(map[string]any)
	if !ok {
		return "", nil, errors.New("data isn't of type map[string]any")
	}

	// Ignore the error as we can't really do anything if the properties are
	// somehow missing.
	properties, _ := getFlagdProperties(dataMap)

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
			return "", nil, errors.New("bucketing value not supplied and no targetingKey in context")
		}

		bucketBy = fmt.Sprintf("%s%s", properties.FlagKey, targetingKey)
	}

	feDistributions, err := parseFractionalEvaluationDistributions(valuesArray)
	if err != nil {
		return "", nil, err
	}

	return bucketBy, feDistributions, nil
}

func parseFractionalEvaluationDistributions(values []any) (*fractionalEvaluationDistribution, error) {
	feDistributions := &fractionalEvaluationDistribution{
		totalWeight:      0,
		weightedVariants: make([]fractionalEvaluationVariant, len(values)),
	}
	for i := 0; i < len(values); i++ {
		distributionArray, ok := values[i].([]any)
		if !ok {
			return nil, errors.New("distribution elements aren't of type []any. " +
				"please check your rule in flag definition")
		}

		if len(distributionArray) == 0 {
			return nil, errors.New("distribution element needs at least one element")
		}

		variant, ok := distributionArray[0].(string)
		if !ok {
			return nil, errors.New("first element of distribution element isn't string")
		}

		weight := 1.0
		if len(distributionArray) >= 2 {
			distributionWeight, ok := distributionArray[1].(float64)
			if ok {
				// default the weight to 1 if not specified explicitly
				weight = distributionWeight
			}
		}

		feDistributions.totalWeight += int(weight)
		feDistributions.weightedVariants[i] = fractionalEvaluationVariant{
			variant: variant,
			weight:  int(weight),
		}
	}

	return feDistributions, nil
}

// distributeValue calculate hash for given hash key and find the bucket distributions belongs to
func distributeValue(value string, feDistribution *fractionalEvaluationDistribution) string {
	hashValue := int32(murmur3.StringSum32(value))
	hashRatio := math.Abs(float64(hashValue)) / math.MaxInt32
	bucket := hashRatio * 100 // in range [0, 100]

	rangeEnd := float64(0)
	for _, weightedVariant := range feDistribution.weightedVariants {
		rangeEnd += weightedVariant.getPercentage(feDistribution.totalWeight)
		if bucket < rangeEnd {
			return weightedVariant.variant
		}
	}

	return ""
}
