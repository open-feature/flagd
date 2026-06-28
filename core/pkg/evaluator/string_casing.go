package evaluator

import (
	"errors"
	"fmt"

	"github.com/open-feature/flagd/core/pkg/logger"
)

const (
	LowerEvaluationName = "lower"
	UpperEvaluationName = "upper"
)

type StringCasingEvaluator struct {
	Logger *logger.Logger
}

func NewStringCasingEvaluator(log *logger.Logger) *StringCasingEvaluator {
	return &StringCasingEvaluator{Logger: log}
}

// LowerEvaluation transforms the given string to lower case so that targeting rules can match
// case-insensitively when composed with other operators (e.g. ==, in, starts_with).
// It returns the ASCII lower-cased string, leaving non-ASCII bytes untouched. As an example,
// it can be used in the following way inside an 'if' evaluation:
//
//	{
//	  "if": [
//			{
//				"==": [{"lower": [{"var": "email"}]}, "user@example.com"]
//			},
//			"red", null
//			]
//	}
//
// This rule can be applied to the following data object, where the evaluation will resolve to 'true':
//
// { "email": "User@Example.com" }
//
// Casing is restricted to ASCII (A-Z) to remain deterministic and consistent across all flagd
// provider implementations; see lowerString for details.
func (sce *StringCasingEvaluator) LowerEvaluation(values, _ interface{}) interface{} {
	value, err := parseStringCasingEvaluationData(values)
	if err != nil {
		sce.Logger.Error(fmt.Sprintf("parse lower evaluation data: %v", err))
		return nil
	}
	return lowerString(value)
}

// UpperEvaluation transforms the given string to upper case so that targeting rules can match
// case-insensitively when composed with other operators (e.g. ==, in, starts_with).
// It returns the ASCII upper-cased string, leaving non-ASCII bytes untouched. As an example,
// it can be used in the following way inside an 'if' evaluation:
//
//	{
//	  "if": [
//			{
//				"==": [{"upper": [{"var": "country"}]}, "US"]
//			},
//			"red", null
//			]
//	}
//
// This rule can be applied to the following data object, where the evaluation will resolve to 'true':
//
// { "country": "us" }
//
// Casing is restricted to ASCII (a-z) to remain deterministic and consistent across all flagd
// provider implementations; see upperString for details.
func (sce *StringCasingEvaluator) UpperEvaluation(values, _ interface{}) interface{} {
	value, err := parseStringCasingEvaluationData(values)
	if err != nil {
		sce.Logger.Error(fmt.Sprintf("parse upper evaluation data: %v", err))
		return nil
	}
	return upperString(value)
}

// parseStringCasingEvaluationData tries to parse the input for the lower/upper evaluation.
// lower/upper takes a single argument which must resolve to a string. Because jsonLogic
// resolves nested operations before this function is called, the input arrives either as a bare
// string (e.g. {"lower": {"var": "email"}}) or as a single-element array (e.g.
// {"lower": [{"var": "email"}]}); both forms are accepted.
func parseStringCasingEvaluationData(values interface{}) (string, error) {
	if property, ok := values.(string); ok {
		return property, nil
	}

	parsed, ok := values.([]interface{})
	if !ok {
		return "", errors.New("[lower/upper] evaluation must be a string or an array")
	}

	if len(parsed) != 1 {
		return "", errors.New("[lower/upper] evaluation must contain a single value")
	}

	property, ok := parsed[0].(string)
	if !ok {
		return "", errors.New("[lower/upper] evaluation: value did not resolve to a string value")
	}

	return property, nil
}

// lowerString returns s with ASCII upper-case letters (A-Z) folded to lower case. Non-ASCII bytes
// are left unchanged so the result is identical across SDKs that use different (locale-sensitive or
// full) Unicode case mappings; strings.ToLower is intentionally avoided for this reason.
func lowerString(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + ('a' - 'A')
		}
	}
	return string(b)
}

// upperString returns s with ASCII lower-case letters (a-z) folded to upper case. Non-ASCII bytes
// are left unchanged so the result is identical across SDKs that use different (locale-sensitive or
// full) Unicode case mappings; strings.ToUpper is intentionally avoided for this reason.
func upperString(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'a' && c <= 'z' {
			b[i] = c - ('a' - 'A')
		}
	}
	return string(b)
}
