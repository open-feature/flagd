package eval

import (
	"fmt"
	"strings"
)

// semVerEvaluation checks if the given property matches a semantic versioning condition.
// It returns 'true', if the value of the given property meets the condition, 'false' if not.
// As an example, it can be used in the following way inside an 'if' evaluation:
//
//	{
//	  "if": [
//			{
//				"starts_with": [{"var": "version"}, ">=" "1.0.0"]
//			},
//			"red", null
//			]
//	}
//
// This rule can be applied to the following data object, where the evaluation will resolve to 'true':
//
// { "email": "2.0.0" }
//
// Note that the 'starts_with'  evaluation rule must contain exactly two items, which both resolve to a
// string value
func (je *JSONEvaluator) startsWithEvaluation(values, _ interface{}) interface{} {
	propertyValue, target, err := parseStringComparisonEvaluationData(values)
	if err != nil {
		je.Logger.Error(fmt.Sprintf("parse starts_with evaluation data: %v", err))
		return nil
	}
	return strings.HasPrefix(propertyValue, target)
}
