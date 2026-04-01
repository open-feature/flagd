package evaluator

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/store"
	"go.uber.org/zap"
)

const (
	StartsWithEvaluationName = "starts_with"
	EndsWithEvaluationName   = "ends_with"
	RegexMatchEvaluationName = "regex_match"
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

type RegexMatchEvaluator struct {
	Logger *logger.Logger
	// RegexCache caches compiled regex patterns for reuse
	RegexCache *map[string]*regexp.Regexp
	// PrevRegexCache holds the previous cache to allow for cache retention across config reloads
	PrevRegexCache *map[string]*regexp.Regexp
}

func NewRegexMatchEvaluator(log *logger.Logger, s store.IStore) *RegexMatchEvaluator {
	self := &RegexMatchEvaluator{
		Logger: 		log,
		RegexCache:     &map[string]*regexp.Regexp{},
		PrevRegexCache: &map[string]*regexp.Regexp{},
	}

	watcher := make(chan store.FlagQueryResult, 1)
	go func() {
		for range watcher {
			// On config change, rotate the regex caches
			// If the current cache is empty, do nothing, to keep the previous cache intact in case it is still helpful
			if len(*self.RegexCache) == 0 {
				continue
			}
			self.PrevRegexCache = self.RegexCache
			self.RegexCache = &map[string]*regexp.Regexp{}
		}
	}()
	selector := store.NewSelector("")
	s.Watch(context.Background(), &selector, watcher)

	return self
}

// RegexMatchEvaluation checks if the given property matches a certain regex pattern.
// It returns 'true', if the value of the given property matches the pattern, 'false' if not.
// As an example, it can be used in the following way inside an 'if' evaluation:
//
//	{
//	  "if": [
//			{
//				"regex_match": [{"var": "email"}, ".*@faas\\.com"]
//			},
//			"red", null
//			]
//	}
//
// This rule can be applied to the following data object, where the evaluation will resolve to 'true':
//
// { "email": "user@faas.com" }
//
// Note that the 'regex_match' evaluation rule must contain two or three items, all of which resolve to a
// string value.
// The first item is the property to check, the second item is the regex pattern to match against,
// and an optional third item can contain regex flags (e.g. "i" for case-insensitive matching).
func (rme *RegexMatchEvaluator) RegexMatchEvaluation(values, _ interface{}) interface{} {
	propertyValue, pattern, flags, err := parseRegexMatchEvaluationData(values)
	if err != nil {
		rme.Logger.Error("error parsing regex_match evaluation data: %v", zap.Error(err))
		return false
	}

	re, err := rme.getRegex(pattern, flags)
	if err != nil {
		rme.Logger.Error("error compiling regex pattern: %v", zap.Error(err))
		return false
	}

	return re.MatchString(propertyValue)
}

var validFlagsStringRe *regexp.Regexp = regexp.MustCompile("[imsU]+")

func parseRegexMatchEvaluationData(values interface{}) (string, string, string, error) {
	parsed, ok := values.([]interface{})
	if !ok {
		return "", "", "", errors.New("regex_match evaluation is not an array")
	}

	if len(parsed) != 2 && len(parsed) != 3 {
		return "", "", "", errors.New("regex_match evaluation must contain a value, a regex pattern, and (optionally) regex flags")
	}

	property, ok := parsed[0].(string)
	if !ok {
		return "", "", "", errors.New("regex_match evaluation: property did not resolve to a string value")
	}

	pattern, ok := parsed[1].(string)
	if !ok {
		return "", "", "", errors.New("regex_match evaluation: pattern did not resolve to a string value")
	}

	flags := ""
	if (len(parsed) == 3) {
		flags, ok = parsed[2].(string)
		if !ok {
			return "", "", "", errors.New("regex_match evaluation: flags did not resolve to a string value")
		}
		if !validFlagsStringRe.MatchString(flags) {
			return "", "", "", errors.New("regex_match evaluation: flags value is invalid")
		}
	}

	return property, pattern, flags, nil
}

func (rme *RegexMatchEvaluator) getRegex(pattern string, flags string) (*regexp.Regexp, error) {
	finalPattern := pattern
	if flags != "" {
		finalPattern = fmt.Sprintf("(?%s)%s", flags, pattern)
	}

	if cached, exists := (*rme.RegexCache)[finalPattern]; exists {
		return cached, nil
	}

	// Check previous cache to allow for cache retention across config reloads
	if cached, exists := (*rme.PrevRegexCache)[finalPattern]; exists {
		(*rme.RegexCache)[finalPattern] = cached
		delete(*rme.PrevRegexCache, finalPattern)
		return cached, nil
	}

	regexp, err := regexp.Compile(finalPattern)
	if err != nil {
		return nil, err
	}
	return regexp, nil
}