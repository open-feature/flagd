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
// provider implementations; see asciiCase for details.
func (sce *StringCasingEvaluator) LowerEvaluation(values, _ interface{}) interface{} {
	return sce.evaluateCasing(LowerEvaluationName, values, toASCIILower)
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
// provider implementations; see asciiCase for details.
func (sce *StringCasingEvaluator) UpperEvaluation(values, _ interface{}) interface{} {
	return sce.evaluateCasing(UpperEvaluationName, values, toASCIIUpper)
}

// evaluateCasing parses the single string argument shared by the lower/upper operators and applies
// the given per-byte case mapping. It returns nil (so jsonLogic falls back to the default variant)
// when the argument cannot be resolved to a string.
func (sce *StringCasingEvaluator) evaluateCasing(op string, values interface{}, mapByte func(byte) byte) interface{} {
	value, err := parseStringCasingEvaluationData(values)
	if err != nil {
		sce.Logger.Error(fmt.Sprintf("parse %s evaluation data: %v", op, err))
		return nil
	}
	return asciiCase(value, mapByte)
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

// asciiCase returns s with every byte remapped by mapByte. The mappers used here (toASCIILower and
// toASCIIUpper) only alter ASCII letters and leave every other byte unchanged, so the result is
// identical across SDKs that use different (locale-sensitive or full) Unicode case mappings;
// strings.ToLower and strings.ToUpper are intentionally avoided for this reason.
func asciiCase(s string, mapByte func(byte) byte) string {
	b := []byte(s)
	for i, c := range b {
		b[i] = mapByte(c)
	}
	return string(b)
}

// toASCIILower folds an ASCII upper-case letter (A-Z) to lower case and leaves every other byte
// unchanged.
func toASCIILower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + ('a' - 'A')
	}
	return c
}

// toASCIIUpper folds an ASCII lower-case letter (a-z) to upper case and leaves every other byte
// unchanged.
func toASCIIUpper(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c - ('a' - 'A')
	}
	return c
}
