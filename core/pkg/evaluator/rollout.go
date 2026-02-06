package evaluator

import (
	"errors"
	"fmt"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/twmb/murmur3"
)

const RolloutEvaluationName = "rollout"

// Rollout is a custom JSONLogic operator for time-based progressive feature rollouts.
// It gradually transitions users from one variant/expression to another over a specified time period.
// For complex distributions, compose with the "fractional" operator.
//
// Syntax:
//
// Shorthand (roll from defaultVariant to target):
//
//	{"rollout": [<startTime>, <endTime>, <to>]}
//
// Longhand (explicit from and to):
//
//	{"rollout": [<startTime>, <endTime>, <from>, <to>]}
//
// With custom bucketBy:
//
//	{"rollout": [<bucketBy>, <startTime>, <endTime>, <from>, <to>]}
//
// Parameters:
//   - bucketBy: Optional. JSONLogic expression (e.g., {"var": "email"}).
//     If omitted, uses flagKey + targetingKey (same as fractional).
//   - startTime: Unix timestamp (int32) when rollout begins (0% on to)
//   - endTime: Unix timestamp (int32) when rollout completes (100% on to)
//   - from: The starting variant/expression. Can be a string, bool, int, or JSONLogic.
//   - to: The target variant/expression. Can be a string, bool, int, or JSONLogic.
//
// The rollout linearly increases the percentage of users seeing toVariant from 0% at
// startTime to 100% at endTime. For more complex distributions at either end, use
// the fractional operator as the from/to value.
//
// Examples:
//
// 1. Simple rollout to "new":
//
//	{"rollout": [1704067200, 1706745600, "new"]}
//
// 2. Explicit from/to:
//
//	{"rollout": [1704067200, 1706745600, "old", "new"]}
//
// 3. With custom bucketBy:
//
//	{"rollout": [{"var": "email"}, 1704067200, 1706745600, "old", "new"]}
//
// 4. Rollout to a fractional split (50% A, 50% B):
//
//	{"rollout": [1704067200, 1706745600, "old", {"fractional": [["a", 1], ["b", 1]]}]}
//
// 5. Phased rollout using if + timestamp:
//
//	{
//	  "if": [
//	    {"<": [{"var": "$flagd.timestamp"}, 1736294400]},
//	    {"rollout": [1735689600, 1736294400, "old", "new"]},
//	    "new"
//	  ]
//	}
type Rollout struct {
	Logger *logger.Logger
}

type rolloutConfig struct {
	bucketBy    string
	from        any // string, bool, int, or JSONLogic node
	to          any // string, bool, int, or JSONLogic node
	startTime   int32
	endTime     int32
	currentTime int32
	data        any // original data context for evaluating nested JSONLogic
}

func NewRollout(logger *logger.Logger) *Rollout {
	return &Rollout{Logger: logger}
}

func (r *Rollout) Evaluate(values, data any) any {
	config, err := r.parseRolloutData(values, data)
	if err != nil {
		r.Logger.Warn(fmt.Sprintf("parse rollout evaluation data: %v", err))
		return nil
	}

	return r.determineVariant(config)
}

func (r *Rollout) parseRolloutData(values, data any) (*rolloutConfig, error) {
	valuesArray, ok := values.([]any)
	if !ok {
		return nil, errors.New("rollout evaluation data is not an array")
	}

	if len(valuesArray) < 3 {
		return nil, errors.New("rollout requires at least 3 arguments: startTime, endTime, toVariant")
	}

	dataMap, ok := data.(map[string]any)
	if !ok {
		return nil, errors.New("data isn't of type map[string]any")
	}

	properties, _ := getFlagdProperties(dataMap)
	config := &rolloutConfig{
		currentTime: int32(properties.Timestamp),
		data:        data,
	}

	if config.currentTime == 0 {
		return nil, errors.New("timestamp not available in context for rollout")
	}

	argIndex := 0

	if _, err := extractInt(valuesArray[0]); err != nil {
		// first arg is not a number — it's a bucketBy value.
		// Note: JSONLogic pre-evaluates all custom operator arguments via parseValues(),
		// so expressions like {"var": "email"} or {"cat": [...]} arrive here already
		// resolved to their string, bool, or numeric result.
		switch v := valuesArray[0].(type) {
		case string:
			config.bucketBy = v
			argIndex = 1
		case bool:
			config.bucketBy = fmt.Sprintf("%v", v)
			argIndex = 1
		case float64:
			config.bucketBy = fmt.Sprintf("%v", v)
			argIndex = 1
		case nil:
			argIndex = 1
		default:
			return nil, fmt.Errorf("invalid first argument type %T: expected string, bool, number, or nil", valuesArray[0])
		}
	}

	// if no explicit bucketBy, use flagKey + targetingKey
	if config.bucketBy == "" {
		targetingKey, ok := dataMap[targetingKeyKey].(string)
		if !ok {
			return nil, errors.New("bucketing value not supplied and no targetingKey in context")
		}
		config.bucketBy = fmt.Sprintf("%s%s", properties.FlagKey, targetingKey)
	}

	// parse startTime
	if argIndex >= len(valuesArray) {
		return nil, errors.New("missing startTime argument")
	}
	startTime, err := extractInt(valuesArray[argIndex])
	if err != nil {
		return nil, fmt.Errorf("startTime: %w", err)
	}
	config.startTime = startTime
	argIndex++

	// parse endTime
	if argIndex >= len(valuesArray) {
		return nil, errors.New("missing endTime argument")
	}
	endTime, err := extractInt(valuesArray[argIndex])
	if err != nil {
		return nil, fmt.Errorf("endTime: %w", err)
	}
	config.endTime = endTime
	argIndex++

	if config.endTime <= config.startTime {
		return nil, errors.New("endTime must be after startTime")
	}

	// parse variants
	remaining := len(valuesArray) - argIndex

	if remaining == 1 {
		// shorthand: just toVariant, from = nil (defaultVariant)
		config.from = nil
		config.to = valuesArray[argIndex]
	} else if remaining >= 2 {
		// longhand: fromVariant, toVariant
		config.from = valuesArray[argIndex]
		config.to = valuesArray[argIndex+1]
	} else {
		return nil, errors.New("missing variant argument(s)")
	}

	return config, nil
}

// extractInt extracts an int32 from a JSON number (always float64 in Go).
func extractInt(val any) (int32, error) {
	v, ok := val.(float64)
	if !ok {
		return 0, fmt.Errorf("must be a number, got %T", val)
	}
	return int32(v), nil
}

// calculateProgress returns the rollout progress as elapsed time since start
func (r *Rollout) calculateElapsed(config *rolloutConfig) (elapsed int64, duration int64) {
	duration = int64(config.endTime - config.startTime)
	if config.currentTime <= config.startTime {
		return 0, duration
	}
	if config.currentTime >= config.endTime {
		return duration, duration
	}
	elapsed = int64(config.currentTime - config.startTime)
	return elapsed, duration
}

// determineVariant uses murmur3 hashing to deterministically select either from or to variant based on the rollout progress.
func (r *Rollout) determineVariant(config *rolloutConfig) any {
	elapsed, duration := r.calculateElapsed(config)

	// edge cases:
	if elapsed <= 0 {
		return config.from
	}
	if elapsed >= duration {
		return config.to
	}

	// use integer math for bucket assignment - maps hash to [0, duration) range, then compares against elapsed
	hashValue := murmur3.StringSum32(config.bucketBy)
	bucket := (uint64(hashValue) * uint64(duration)) >> 32

	if bucket < uint64(elapsed) {
		return config.to
	}
	return config.from
}
