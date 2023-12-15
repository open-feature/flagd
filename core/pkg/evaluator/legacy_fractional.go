// This evaluation type is deprecated and will be removed before v1.
// Do not enhance it or use it for reference.

package evaluator

import (
	"errors"
	"fmt"
	"math"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/zeebo/xxh3"
)

const (
	LegacyFractionEvaluationName = "fractionalEvaluation"
	LegacyFractionEvaluationLink = "https://flagd.dev/concepts/#migrating-from-legacy-fractionalevaluation"
)

// Deprecated: LegacyFractional is deprecated. This will be removed prior to v1 release.
type LegacyFractional struct {
	Logger *logger.Logger
}

type legacyFractionalEvaluationDistribution struct {
	variant    string
	percentage int
}

func NewLegacyFractional(logger *logger.Logger) *LegacyFractional {
	return &LegacyFractional{Logger: logger}
}

func (fe *LegacyFractional) LegacyFractionalEvaluation(values, data interface{}) interface{} {
	fe.Logger.Warn(
		fmt.Sprintf("%s is deprecated, please use %s, see: %s",
			LegacyFractionEvaluationName,
			FractionEvaluationName,
			LegacyFractionEvaluationLink))

	valueToDistribute, feDistributions, err := parseLegacyFractionalEvaluationData(values, data)
	if err != nil {
		fe.Logger.Error(fmt.Sprintf("parse fractional evaluation data: %v", err))
		return nil
	}

	return distributeLegacyValue(valueToDistribute, feDistributions)
}

func parseLegacyFractionalEvaluationData(values, data interface{}) (string,
	[]legacyFractionalEvaluationDistribution, error,
) {
	valuesArray, ok := values.([]interface{})
	if !ok {
		return "", nil, errors.New("fractional evaluation data is not an array")
	}
	if len(valuesArray) < 2 {
		return "", nil, errors.New("fractional evaluation data has length under 2")
	}

	bucketBy, ok := valuesArray[0].(string)
	if !ok {
		return "", nil, errors.New("first element of fractional evaluation data isn't of type string")
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return "", nil, errors.New("data isn't of type map[string]interface{}")
	}

	v, ok := dataMap[bucketBy]
	if !ok {
		return "", nil, nil
	}

	valueToDistribute, ok := v.(string)
	if !ok {
		return "", nil, fmt.Errorf("var: %s isn't of type string", bucketBy)
	}

	feDistributions, err := parseLegacyFractionalEvaluationDistributions(valuesArray)
	if err != nil {
		return "", nil, err
	}

	return valueToDistribute, feDistributions, nil
}

func parseLegacyFractionalEvaluationDistributions(values []interface{}) (
	[]legacyFractionalEvaluationDistribution, error,
) {
	sumOfPercentages := 0
	var feDistributions []legacyFractionalEvaluationDistribution
	for i := 1; i < len(values); i++ {
		distributionArray, ok := values[i].([]interface{})
		if !ok {
			return nil, errors.New("distribution elements aren't of type []interface{}")
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

		feDistributions = append(feDistributions, legacyFractionalEvaluationDistribution{
			variant:    variant,
			percentage: int(percentage),
		})
	}

	if sumOfPercentages != 100 {
		return nil, fmt.Errorf("percentages must sum to 100, got: %d", sumOfPercentages)
	}

	return feDistributions, nil
}

func distributeLegacyValue(value string, feDistribution []legacyFractionalEvaluationDistribution) string {
	hashValue := xxh3.HashString(value)

	hashRatio := float64(hashValue) / math.Pow(2, 64) // divide the hash value by the largest possible value, integer 2^64

	bucket := int(hashRatio * 100) // integer in range [0, 99]

	rangeEnd := 0
	for _, dist := range feDistribution {
		rangeEnd += dist.percentage
		if bucket < rangeEnd {
			return dist.variant
		}
	}

	return ""
}
