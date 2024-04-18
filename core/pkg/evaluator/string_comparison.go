package evaluator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/open-feature/flagd/core/pkg/logger"
)

const (
	StartsWithEvaluationName = "starts_with"
	EndsWithEvaluationName   = "ends_with"
)

type StringComparisonEvaluator struct {
	Logger *logger.Logger
}

func NewStringComparisonEvaluator(log *logger.Logger) *StringComparisonEvaluator {
	return &StringComparisonEvaluator{Logger: log}
}

// StartsWithEvaluation checks if the given property starts with a certain prefix.
// It returns 'true', if the value of the given property starts with the prefix, 'false' if not.
// As an example, it can be used in the following way inside an 'if' evaluation:
//
//	{
//	  "if": [
//			{
//				"starts_with": [{"var": "email"}, "user@faas"]
//			},
//			"red", null
//			]
//	}
//
// This rule can be applied to the following data object, where the evaluation will resolve to 'true':
//
// { "email": "user@faas.com" }
//
// Note that the 'starts_with' evaluation rule must contain exactly two items, which both resolve to a
// string value
func (sce *StringComparisonEvaluator) StartsWithEvaluation(values, _ interface{}) interface{} {
	propertyValue, target, err := parseStringComparisonEvaluationData(values)
	if err != nil {
		sce.Logger.Error(fmt.Sprintf("parse starts_with evaluation data: %v", err))
		return nil
	}
	return strings.HasPrefix(propertyValue, target)
}

// EndsWithEvaluation checks if the given property ends with a certain prefix.
// It returns 'true', if the value of the given property starts with the prefix, 'false' if not.
// As an example, it can be used in the following way inside an 'if' evaluation:
//
//	{
//	  "if": [
//			{
//				"ends_with": [{"var": "email"}, "faas.com"]
//			},
//			"red", null
//			]
//	}
//
// This rule can be applied to the following data object, where the evaluation will resolve to 'true':
//
// { "email": "user@faas.com" }
//
// Note that the 'ends_with'  evaluation rule must contain exactly two items, which both resolve to a
// string value
func (sce *StringComparisonEvaluator) EndsWithEvaluation(values, _ interface{}) interface{} {
	propertyValue, target, err := parseStringComparisonEvaluationData(values)
	if err != nil {
		sce.Logger.Error(fmt.Sprintf("parse ends_with evaluation data: %v", err))
		return false
	}
	return strings.HasSuffix(propertyValue, target)
}

// parseStringComparisonEvaluationData tries to parse the input for the starts_with/ends_with evaluation.
// this evaluator requires an array containing exactly two strings.
// Note that, when used with jsonLogic, those two items can also have been objects in the original 'values' object,
// which have been resolved to string values by jsonLogic before this function is called.
// As an example, the following values object:
//
//	{
//	  "if": [
//			{
//				"starts_with": [{"var": "email"}, "user@faas"]
//			},
//			"red", null
//			]
//	}
//
// with the following data object:
//
// { "email": "user@faas.com" }
//
// will have been resolved to
//
// ["user@faas.com", "user@faas"]
//
// at the time this function is reached.
func parseStringComparisonEvaluationData(values interface{}) (string, string, error) {
	parsed, ok := values.([]interface{})
	if !ok {
		return "", "", errors.New("[start/end]s_with evaluation is not an array")
	}

	if len(parsed) != 2 {
		return "", "", errors.New("[start/end]s_with evaluation must contain a value and a comparison target")
	}

	property, ok := parsed[0].(string)
	if !ok {
		return "", "", errors.New("[start/end]s_with evaluation: property did not resolve to a string value")
	}

	targetValue, ok := parsed[1].(string)
	if !ok {
		return "", "", errors.New("[start/end]s_with evaluation: target value did not resolve to a string value")
	}

	return property, targetValue, nil
}
