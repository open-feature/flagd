package eval

import (
	"errors"
	"fmt"
	"math"

	"github.com/zeebo/xxh3"
)

type fractionalEvaluationDistribution struct {
	variant    string
	percentage int
}

func (je *JSONEvaluator) fractionalEvaluation(values, data interface{}) interface{} {
	valueToDistribute, feDistributions, err := parseFractionalEvaluationData(values, data)
	if err != nil {
		je.Logger.Error(fmt.Sprintf("parseFractionalEvaluationData: %v", err))
		return nil
	}

	return distributeValue(valueToDistribute, feDistributions)
}

func parseFractionalEvaluationData(values, data interface{}) (string, []fractionalEvaluationDistribution, error) {
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
		return "", nil, fmt.Errorf("%s isn't a found var in data", bucketBy)
	}

	valueToDistribute, ok := v.(string)
	if !ok {
		return "", nil, fmt.Errorf("var %s isn't of type string", bucketBy)
	}

	feDistributions, err := parseFractionalEvaluationDistributions(valuesArray)
	if err != nil {
		return "", nil, err
	}

	return valueToDistribute, feDistributions, nil
}

func parseFractionalEvaluationDistributions(values []interface{}) ([]fractionalEvaluationDistribution, error) {
	sumOfPercentages := 0
	var feDistributions []fractionalEvaluationDistribution
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

		feDistributions = append(feDistributions, fractionalEvaluationDistribution{
			variant:    variant,
			percentage: int(percentage),
		})
	}

	if sumOfPercentages != 100 {
		return nil, fmt.Errorf("percentages must sum to 100, got %d", sumOfPercentages)
	}

	return feDistributions, nil
}

func distributeValue(value string, feDistribution []fractionalEvaluationDistribution) string {
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
