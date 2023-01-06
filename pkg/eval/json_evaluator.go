package eval

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/diegoholiveira/jsonlogic/v3"
	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/model"
	schema "github.com/open-feature/schemas/json"
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/structpb"
)

type JSONEvaluator struct {
	state  Flags
	Logger *logger.Logger
}

type constraints interface {
	bool | string | map[string]any | float64
}

const (
	Disabled = "DISABLED"
)

func NewJSONEvaluator(logger *logger.Logger) *JSONEvaluator {
	ev := JSONEvaluator{
		Logger: logger.WithFields(
			zap.String("component", "evaluator"),
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

	s, notifications := je.state.Merge(je.Logger, source, newFlags)
	je.state = s

	return notifications, nil
}

func resolve[T constraints](reqID string, key string, context *structpb.Struct,
	variantEval func(string, string, *structpb.Struct) (string, string, error),
	variants map[string]any) (
	value T,
	variant string,
	reason string,
	err error,
) {
	variant, reason, err = variantEval(reqID, key, context)
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

func (je *JSONEvaluator) ResolveAllValues(reqID string, context *structpb.Struct) []AnyValue {
	values := []AnyValue{}
	var value interface{}
	var variant string
	var reason string
	var err error
	for flagKey, flag := range je.state.Flags {
		defaultValue := flag.Variants[flag.DefaultVariant]
		switch defaultValue.(type) {
		case bool:
			value, variant, reason, err = resolve[bool](
				reqID,
				flagKey,
				context,
				je.evaluateVariant,
				je.state.Flags[flagKey].Variants,
			)
		case string:
			value, variant, reason, err = resolve[string](
				reqID,
				flagKey,
				context,
				je.evaluateVariant,
				je.state.Flags[flagKey].Variants,
			)
		case float64:
			value, variant, reason, err = resolve[float64](
				reqID,
				flagKey,
				context,
				je.evaluateVariant,
				je.state.Flags[flagKey].Variants,
			)
		case map[string]any:
			value, variant, reason, err = resolve[map[string]any](
				reqID,
				flagKey,
				context,
				je.evaluateVariant,
				je.state.Flags[flagKey].Variants,
			)
		}
		if err != nil {
			je.Logger.ErrorWithID(reqID, fmt.Sprintf("Bulk evaluation: key %s returned error %s", flagKey, err.Error()))
			continue
		}
		values = append(values, NewAnyValue(value, variant, reason, flagKey))
	}
	return values
}

func (je *JSONEvaluator) ResolveBooleanValue(reqID string, flagKey string, context *structpb.Struct) (
	value bool,
	variant string,
	reason string,
	err error,
) {
	je.Logger.DebugWithID(reqID, fmt.Sprintf("evaluating boolean flag: %s", reqID))
	return resolve[bool](reqID, flagKey, context, je.evaluateVariant, je.state.Flags[flagKey].Variants)
}

func (je *JSONEvaluator) ResolveStringValue(reqID string, flagKey string, context *structpb.Struct) (
	value string,
	variant string,
	reason string,
	err error,
) {
	je.Logger.DebugWithID(reqID, fmt.Sprintf("evaluating string flag: %s", flagKey))
	return resolve[string](reqID, flagKey, context, je.evaluateVariant, je.state.Flags[flagKey].Variants)
}

func (je *JSONEvaluator) ResolveFloatValue(reqID string, flagKey string, context *structpb.Struct) (
	value float64,
	variant string,
	reason string,
	err error,
) {
	je.Logger.DebugWithID(reqID, fmt.Sprintf("evaluating float flag: %s", flagKey))
	value, variant, reason, err = resolve[float64](
		reqID, flagKey, context, je.evaluateVariant, je.state.Flags[flagKey].Variants)
	return
}

func (je *JSONEvaluator) ResolveIntValue(reqID string, flagKey string, context *structpb.Struct) (
	value int64,
	variant string,
	reason string,
	err error,
) {
	je.Logger.DebugWithID(reqID, fmt.Sprintf("evaluating int flag: %s", flagKey))
	var val float64
	val, variant, reason, err = resolve[float64](
		reqID, flagKey, context, je.evaluateVariant, je.state.Flags[flagKey].Variants)
	value = int64(val)
	return
}

func (je *JSONEvaluator) ResolveObjectValue(reqID string, flagKey string, context *structpb.Struct) (
	value map[string]any,
	variant string,
	reason string,
	err error,
) {
	je.Logger.DebugWithID(reqID, fmt.Sprintf("evaluating object flag: %s", flagKey))
	return resolve[map[string]any](reqID, flagKey, context, je.evaluateVariant, je.state.Flags[flagKey].Variants)
}

// runs the rules (if defined) to determine the variant, otherwise falling through to the default
func (je *JSONEvaluator) evaluateVariant(
	reqID string,
	flagKey string,
	context *structpb.Struct,
) (variant string, reason string, err error) {
	flag, ok := je.state.Flags[flagKey]
	if !ok {
		// flag not found
		je.Logger.DebugWithID(reqID, fmt.Sprintf("requested flag could not be found %s", flagKey))
		return "", model.ErrorReason, errors.New(model.FlagNotFoundErrorCode)
	}

	if flag.State == Disabled {
		je.Logger.DebugWithID(reqID, fmt.Sprintf("requested flag is disabled %s", flagKey))
		return "", model.ErrorReason, errors.New(model.FlagDisabledErrorCode)
	}

	// get the targeting logic, if any
	targeting := flag.Targeting

	if targeting != nil && string(targeting) != "{}" {
		targetingBytes, err := targeting.MarshalJSON()
		if err != nil {
			je.Logger.ErrorWithID(reqID, fmt.Sprintf("Error parsing rules for flag %s, %s", flagKey, err))
			return "", model.ErrorReason, err
		}

		b, err := json.Marshal(context)
		if err != nil {
			je.Logger.ErrorWithID(reqID, fmt.Sprintf("error parsing context for flag %s, %s, %v", flagKey, err, context))

			return "", model.ErrorReason, errors.New(model.ErrorReason)
		}
		var result bytes.Buffer
		// evaluate json-logic rules to determine the variant
		err = jsonlogic.Apply(bytes.NewReader(targetingBytes), bytes.NewReader(b), &result)
		if err != nil {
			je.Logger.ErrorWithID(reqID, fmt.Sprintf("Error applying rules %s", err))
			return "", model.ErrorReason, err
		}
		// strip whitespace and quotes from the variant
		variant = strings.ReplaceAll(strings.TrimSpace(result.String()), "\"", "")

		// if this is a valid variant, return it
		if _, ok := je.state.Flags[flagKey].Variants[variant]; ok {
			return variant, model.TargetingMatchReason, nil
		}

		je.Logger.DebugWithID(reqID, fmt.Sprintf("returning default variant for flagKey %s, variant is not valid", flagKey))
		reason = model.DefaultReason
	} else {
		reason = model.StaticReason
	}

	return je.state.Flags[flagKey].DefaultVariant, reason, nil
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
