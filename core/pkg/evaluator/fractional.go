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
	variant    string
	percentage int
}

func NewFractional(logger *logger.Logger) *Fractional {
	return &Fractional{Logger: logger}
}

func (fe *Fractional) Evaluate(values, data any) any {
	valueToDistribute, feDistributions, err := parseFractionalEvaluationData(values, data)
	if err != nil {
		fe.Logger.Error(fmt.Sprintf("parse fractional evaluation data: %v", err))
		return nil
	}

	return distributeValue(valueToDistribute, feDistributions)
}

func parseFractionalEvaluationData(values, data any) (string, []fractionalEvaluationDistribution, error) {
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
		bucketBy, ok = dataMap[targetingKeyKey].(string)
		if !ok {
			return "", nil, errors.New("bucketing value not supplied and no targetingKey in context")
		}
	}

	feDistributions, err := parseFractionalEvaluationDistributions(valuesArray)
	if err != nil {
		return "", nil, err
	}

	return fmt.Sprintf("%s%s", properties.FlagKey, bucketBy), feDistributions, nil
}

func parseFractionalEvaluationDistributions(values []any) ([]fractionalEvaluationDistribution, error) {
	sumOfPercentages := 0
	var feDistributions []fractionalEvaluationDistribution
	for i := 0; i < len(values); i++ {
		distributionArray, ok := values[i].([]any)
		if !ok {
			return nil, errors.New("distribution elements aren't of type []any")
		}

		if len(distributionArray) != 2 {
			return nil, errors.New("distribution element isn't length 2")
		}

		variant, ok := distributionArray[0].(string)
		if !ok {
			return nil, errors.New("first element of distribution element isn't string")
		}

		percentage, ok := distributionArray[1].(float64)
		if !ok {
			return nil, errors.New("second element of distribution element isn't float")
		}

		sumOfPercentages += int(percentage)

		feDistributions = append(feDistributions, fractionalEvaluationDistribution{
			variant:    variant,
			percentage: int(percentage),
		})
	}

	if sumOfPercentages != 100 {
		return nil, fmt.Errorf("percentages must sum to 100, got: %d", sumOfPercentages)
	}

	return feDistributions, nil
}

// distributeValue calculate hash for given hash key and find the bucket distributions belongs to
func distributeValue(value string, feDistribution []fractionalEvaluationDistribution) string {
	hashValue := int32(murmur3.StringSum32(value))
	hashRatio := math.Abs(float64(hashValue)) / math.MaxInt32
	bucket := int(hashRatio * 100) // in range [0, 100]

	rangeEnd := 0
	for _, dist := range feDistribution {
		rangeEnd += dist.percentage
		if bucket < rangeEnd {
			return dist.variant
		}
	}

	return ""
}
