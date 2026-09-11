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

// errPropertyAbsent signals that the property operand resolved to nil, meaning
// the evaluation context simply did not carry the attribute. Targeting rules
// routinely reference optional attributes, so this is an ordinary condition
// rather than a misconfiguration, and is reported separately from a genuine
// type error so the two can be logged differently.
var errPropertyAbsent = errors.New(
	"[start/end]s_with evaluation: property is absent from the evaluation context")

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
	return sce.evaluate(StartsWithEvaluationName, values, strings.HasPrefix)
}

// evaluate applies cmp to the parsed operands, returning nil -- which jsonLogic
// treats as falsy -- when they cannot be used.
//
// The logging severity here is deliberate. This runs on every evaluation of a
// rule, so its volume tracks request rate rather than the number of bad rules,
// and error level additionally attaches a stack trace to each occurrence. A rule
// referencing an attribute that a client did not send would therefore emit tens
// of log lines per evaluation, for a condition that is entirely normal. An
// absent property is logged at debug; a genuinely malformed rule is logged at
// warn, which still surfaces it but without the stack trace.
func (sce *StringComparisonEvaluator) evaluate(
	name string, values interface{}, cmp func(string, string) bool,
) interface{} {
	propertyValue, target, err := parseStringComparisonEvaluationData(values)
	if err != nil {
		if errors.Is(err, errPropertyAbsent) {
			sce.Logger.Debug(fmt.Sprintf("%s: %v", name, err))
		} else {
			sce.Logger.Warn(fmt.Sprintf("parse %s evaluation data: %v", name, err))
		}
		return nil
	}
	return cmp(propertyValue, target)
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
	return sce.evaluate(EndsWithEvaluationName, values, strings.HasSuffix)
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

	// jsonLogic resolves a `var` referencing a missing attribute to nil, so this
	// distinguishes "the context did not carry it" from "it was the wrong type".
	if parsed[0] == nil {
		return "", "", errPropertyAbsent
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
