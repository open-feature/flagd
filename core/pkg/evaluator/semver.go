package evaluator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/open-feature/flagd/core/pkg/logger"
	"golang.org/x/mod/semver"
)

const SemVerEvaluationName = "sem_ver"

type SemVerOperator string

const (
	Equals         SemVerOperator = "="
	NotEqual       SemVerOperator = "!="
	Less           SemVerOperator = "<"
	LessOrEqual    SemVerOperator = "<="
	GreaterOrEqual SemVerOperator = ">="
	Greater        SemVerOperator = ">"
	MatchMajor     SemVerOperator = "^"
	MatchMinor     SemVerOperator = "~"
)

func (svo SemVerOperator) compare(v1, v2 string) (bool, error) {
	cmpRes := semver.Compare(v1, v2)
	switch svo {
	case Less:
		return cmpRes == -1, nil
	case Equals:
		return cmpRes == 0, nil
	case NotEqual:
		return cmpRes != 0, nil
	case LessOrEqual:
		return cmpRes == -1 || cmpRes == 0, nil
	case GreaterOrEqual:
		return cmpRes == +1 || cmpRes == 0, nil
	case Greater:
		return cmpRes == +1, nil
	case MatchMinor:
		v1MajorMinor := semver.MajorMinor(v1)
		v2MajorMinor := semver.MajorMinor(v2)
		return semver.Compare(v1MajorMinor, v2MajorMinor) == 0, nil
	case MatchMajor:
		v1Major := semver.Major(v1)
		v2Major := semver.Major(v2)
		return semver.Compare(v1Major, v2Major) == 0, nil
	default:
		return false, errors.New("invalid operator")
	}
}

type SemVerComparison struct {
	Logger *logger.Logger
}

func NewSemVerComparison(log *logger.Logger) *SemVerComparison {
	return &SemVerComparison{Logger: log}
}

// SemVerEvaluation checks if the given property matches a semantic versioning condition.
// It returns 'true', if the value of the given property meets the condition, 'false' if not.
// As an example, it can be used in the following way inside an 'if' evaluation:
//
//	{
//	  "if": [
//			{
//				"sem_ver": [{"var": "version"}, ">=", "1.0.0"]
//			},
//			"red", null
//			]
//	}
//
// This rule can be applied to the following data object, where the evaluation will resolve to 'true':
//
// { "version": "2.0.0" }
//
// Note that the 'sem_ver' evaluation rule must contain exactly three items:
// 1. Target property: this needs which both resolve to a semantic versioning string
// 2. Operator: One of the following: '=', '!=', '>', '<', '>=', '<=', '~', '^'
// 3. Target value: this needs which both resolve to a semantic versioning string
func (je *SemVerComparison) SemVerEvaluation(values, _ interface{}) interface{} {
	actualVersion, targetVersion, operator, err := parseSemverEvaluationData(values)
	if err != nil {
		je.Logger.Error(fmt.Sprintf("parse sem_ver evaluation data: %v", err))
		return false
	}
	res, err := operator.compare(actualVersion, targetVersion)
	if err != nil {
		je.Logger.Error(fmt.Sprintf("sem_ver evaluation: %v", err))
		return false
	}
	return res
}

func parseSemverEvaluationData(values interface{}) (string, string, SemVerOperator, error) {
	parsed, ok := values.([]interface{})
	if !ok {
		return "", "", "", errors.New("sem_ver evaluation is not an array")
	}

	if len(parsed) != 3 {
		return "", "", "", errors.New("sem_ver evaluation must contain a value, an operator and a comparison target")
	}

	actualVersion, err := parseSemanticVersion(parsed[0])
	if err != nil {
		return "", "", "", fmt.Errorf("sem_ver evaluation: could not parse target property value: %w", err)
	}

	operator, err := parseOperator(parsed[1])
	if err != nil {
		return "", "", "", fmt.Errorf("sem_ver evaluation: could not parse operator: %w", err)
	}

	targetVersion, err := parseSemanticVersion(parsed[2])
	if err != nil {
		return "", "", "", fmt.Errorf("sem_ver evaluation: could not parse target value: %w", err)
	}
	return actualVersion, targetVersion, operator, nil
}

func parseSemanticVersion(v interface{}) (string, error) {
	version, ok := v.(string)
	if !ok {
		return "", errors.New("sem_ver evaluation: property did not resolve to a string value")
	}
	// version strings are only valid in the semver package if they start with a 'v'
	// if it's not present in the given value, we prepend it
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	if !semver.IsValid(version) {
		return "", errors.New("not a valid semantic version string")
	}

	return version, nil
}

func parseOperator(o interface{}) (SemVerOperator, error) {
	operatorString, ok := o.(string)
	if !ok {
		return "", errors.New("could not parse operator")
	}

	return SemVerOperator(operatorString), nil
}
