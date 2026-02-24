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
	data             any
	logger           *logger.Logger
}

type fractionalEvaluationVariant struct {
	variant any // string, bool, number, JSONLogic node (map[string]any from nested sub-array operators), or nil
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
	valueToDistribute, feDistributions, err := parseFractionalEvaluationData(values, data, fe.Logger)
	if err != nil {
		fe.Logger.Warn(fmt.Sprintf("parse fractional evaluation data: %v", err))
		return nil
	}

	return distributeValue(valueToDistribute, feDistributions)
}

func parseFractionalEvaluationData(values, data any, logger *logger.Logger) (string, *fractionalEvaluationDistribution, error) {
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

	feDistributions, err := parseFractionalEvaluationDistributions(valuesArray, data, logger)
	if err != nil {
		return "", nil, err
	}

	return bucketBy, feDistributions, nil
}

func parseFractionalEvaluationDistributions(values []any, data any, logger *logger.Logger) (*fractionalEvaluationDistribution, error) {
	feDistributions := &fractionalEvaluationDistribution{
		totalWeight:      0,
		weightedVariants: make([]fractionalEvaluationVariant, len(values)),
		data:             data,
		logger:           logger,
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

		// variants can be any JSONLogic node.
		// github.com/diegoholiveira/jsonlogic/v3@v3.8.5+ pre-evaluates JSONLogic nodes, so we don't need to directly evaluate them ourselves.
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
			return nil, errors.New("first element of distribution element must be a string, bool, number, or nil")
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
func distributeValue(value string, feDistribution *fractionalEvaluationDistribution) any {
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
