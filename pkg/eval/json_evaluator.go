package eval

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/diegoholiveira/jsonlogic/v3"
	"github.com/open-feature/flagd/pkg/model"
	schema "github.com/open-feature/schemas/json"
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/structpb"
)

type JSONEvaluator struct {
	state  Flags
	Logger *zap.Logger
}

type constraints interface {
	bool | string | map[string]any | float64
}

const (
	Disabled = "DISABLED"
)

func NewJSONEvaluator(logger *zap.Logger) *JSONEvaluator {
	ev := JSONEvaluator{
		Logger: logger.With(
			zap.String("component", "service"),
			zap.String("evaluator", "json"),
		),
	}
	jsonlogic.AddOperator("fractionalEvaluation", ev.fractionalEvaluation)
	return &ev
}

func (je *JSONEvaluator) GetState() (string, error) {
	data, err := json.Marshal(&je.state)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (je *JSONEvaluator) SetState(source string, state string) ([]StateChangeNotification, error) {
	schemaLoader := gojsonschema.NewStringLoader(schema.FlagdDefinitions)
	flagStringLoader := gojsonschema.NewStringLoader(state)
	result, err := gojsonschema.Validate(schemaLoader, flagStringLoader)

	if err != nil {
		return nil, err
	} else if !result.Valid() {
		err := errors.New("invalid JSON file")
		return nil, err
	}

	state, err = je.transposeEvaluators(state)
	if err != nil {
		return nil, fmt.Errorf("transpose evaluators: %w", err)
	}

	var newFlags Flags
	err = json.Unmarshal([]byte(state), &newFlags)
	if err != nil {
		return nil, fmt.Errorf("unmarshal new state: %w", err)
	}
	if err := validateDefaultVariants(newFlags); err != nil {
		return nil, err
	}

	s, notifications := je.state.Merge(source, newFlags)
	je.state = s

	return notifications, nil
}

func resolve[T constraints](key string, context *structpb.Struct,
	variantEval func(string, *structpb.Struct) (string, string, error),
	variants map[string]any) (
	value T,
	variant string,
	reason string,
	err error,
) {
	variant, reason, err = variantEval(key, context)
	if err != nil {
		return value, variant, reason, err
	}

	var ok bool
	value, ok = variants[variant].(T)
	if !ok {
		return value, variant, model.ErrorReason, errors.New(model.TypeMismatchErrorCode)
	}

	return value, variant, reason, nil
}

func (je *JSONEvaluator) ResolveBooleanValue(flagKey string, context *structpb.Struct) (
	value bool,
	variant string,
	reason string,
	err error,
) {
	return resolve[bool](flagKey, context, je.evaluateVariant, je.state.Flags[flagKey].Variants)
}

func (je *JSONEvaluator) ResolveStringValue(flagKey string, context *structpb.Struct) (
	value string,
	variant string,
	reason string,
	err error,
) {
	return resolve[string](flagKey, context, je.evaluateVariant, je.state.Flags[flagKey].Variants)
}

func (je *JSONEvaluator) ResolveFloatValue(flagKey string, context *structpb.Struct) (
	value float64,
	variant string,
	reason string,
	err error,
) {
	value, variant, reason, err = resolve[float64](flagKey, context, je.evaluateVariant, je.state.Flags[flagKey].Variants)
	return
}

func (je *JSONEvaluator) ResolveIntValue(flagKey string, context *structpb.Struct) (
	value int64,
	variant string,
	reason string,
	err error,
) {
	var val float64
	val, variant, reason, err = resolve[float64](flagKey, context, je.evaluateVariant, je.state.Flags[flagKey].Variants)
	value = int64(val)
	return
}

func (je *JSONEvaluator) ResolveObjectValue(flagKey string, context *structpb.Struct) (
	value map[string]any,
	variant string,
	reason string,
	err error,
) {
	return resolve[map[string]any](flagKey, context, je.evaluateVariant, je.state.Flags[flagKey].Variants)
}

// runs the rules (if defined) to determine the variant, otherwise falling through to the default
func (je *JSONEvaluator) evaluateVariant(
	flagKey string,
	context *structpb.Struct,
) (variant string, reason string, err error) {
	flag, ok := je.state.Flags[flagKey]
	if !ok {
		// flag not found
		return "", model.ErrorReason, errors.New(model.FlagNotFoundErrorCode)
	}

	if flag.State == Disabled {
		return "", model.ErrorReason, errors.New(model.FlagDisabledErrorCode)
	}

	// get the targeting logic, if any
	targeting := flag.Targeting

	if targeting != nil {
		targetingBytes, err := targeting.MarshalJSON()
		if err != nil {
			je.Logger.Error(fmt.Sprintf("Error parsing rules for flag %s, %s", flagKey, err))
			return "", model.ErrorReason, err
		}

		b, err := json.Marshal(context)
		if err != nil {
			je.Logger.Error(fmt.Sprintf("error parsing context for flag %s, %s, %v", flagKey, err, context))

			return "", model.ErrorReason, errors.New(model.ErrorReason)
		}
		var result bytes.Buffer
		// evaluate json-logic rules to determine the variant
		err = jsonlogic.Apply(bytes.NewReader(targetingBytes), bytes.NewReader(b), &result)
		if err != nil {
			je.Logger.Error(fmt.Sprintf("Error applying rules %s", err))
			return "", model.ErrorReason, err
		}
		// strip whitespace and quotes from the variant
		variant = strings.ReplaceAll(strings.TrimSpace(result.String()), "\"", "")
	}

	// if this is a valid variant, return it
	if _, ok := je.state.Flags[flagKey].Variants[variant]; ok {
		return variant, model.TargetingMatchReason, nil
	}

	// if it's not a valid variant, use the default value
	return je.state.Flags[flagKey].DefaultVariant, model.DefaultReason, nil
}

// validateDefaultVariants returns an error if any of the default variants aren't valid
func validateDefaultVariants(flags Flags) error {
	for name, flag := range flags.Flags {
		if _, ok := flag.Variants[flag.DefaultVariant]; !ok {
			return fmt.Errorf(
				"default variant '%s' isn't a valid variant of flag '%s'", flag.DefaultVariant, name,
			)
		}
	}

	return nil
}

func (je *JSONEvaluator) transposeEvaluators(state string) (string, error) {
	var evaluators Evaluators
	if err := json.Unmarshal([]byte(state), &evaluators); err != nil {
		return "", fmt.Errorf("unmarshal: %w", err)
	}

	for evalName, evalRaw := range evaluators.Evaluators {
		// replace any occurrences of "evaluator": "evalName"
		regex, err := regexp.Compile(fmt.Sprintf(`"\$ref":(\s)*"%s"`, evalName))
		if err != nil {
			return "", fmt.Errorf("compile regex: %w", err)
		}

		marshalledEval, err := evalRaw.MarshalJSON()
		if err != nil {
			return "", fmt.Errorf("marshal evaluator: %w", err)
		}

		evalValue := string(marshalledEval)
		if len(evalValue) < 3 {
			return "", errors.New("evaluator object is empty")
		}
		evalValue = evalValue[1 : len(evalValue)-2] // remove first { and last }

		state = regex.ReplaceAllString(state, evalValue)
	}

	return state, nil
}
