package eval

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/diegoholiveira/jsonlogic/v3"
	gen "github.com/open-feature/flagd/pkg/generated"
	"github.com/open-feature/flagd/pkg/model"
	log "github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
)

//go:embed flagd-definitions.json
var schema string

type JSONEvaluator struct {
	state Flags
}

func (je *JSONEvaluator) GetState() (string, error) {
	bytes, err := json.Marshal(&je.state)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (je *JSONEvaluator) SetState(state string) error {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	flagStringLoader := gojsonschema.NewStringLoader(state)
	result, err := gojsonschema.Validate(schemaLoader, flagStringLoader)

	// TODO: we can add validation for all rules by calling jsonlogic.IsValid() on each

	if err != nil {
		return err
	} else if !result.Valid() {
		err := errors.New("invalid JSON file")
		log.Error(err)
		return err
	}

	var newFlags Flags
	err = json.Unmarshal([]byte(state), &newFlags)
	if err != nil {
		return fmt.Errorf("unmarshal new state: %w", err)
	}
	je.state = je.state.Merge(newFlags)

	return nil
}

// TODO: might be able to simplify some of this with generics.
func (je *JSONEvaluator) ResolveBooleanValue(flagKey string, defaultValue bool, context gen.Context) (
	value bool,
	reason string,
	err error,
) {
	variant, reason, err := je.evaluateVariant(flagKey, context)
	if err != nil {
		log.Errorf("Error evaluating flag, %s", err.Error())
		return defaultValue, reason, err
	}

	val, ok := je.state.Flags[flagKey].Variants[variant].(bool)
	if !ok {
		log.Errorf("Error converting %s to bool", flagKey)
		return defaultValue, model.ErrorReason, errors.New(model.TypeMismatchErrorCode)
	}
	return val, reason, nil
}

func (je *JSONEvaluator) ResolveStringValue(flagKey string, defaultValue string, context gen.Context) (
	value string,
	reason string,
	err error,
) {
	variant, reason, err := je.evaluateVariant(flagKey, context)
	if err != nil {
		log.Errorf("Error evaluating flag, %s", err.Error())
		return defaultValue, reason, err
	}

	val, ok := je.state.Flags[flagKey].Variants[variant].(string)
	if !ok {
		log.Errorf("Error converting %s to string", flagKey)
		return defaultValue, model.ErrorReason, errors.New(model.TypeMismatchErrorCode)
	}
	return val, reason, nil
}

func (je *JSONEvaluator) ResolveNumberValue(flagKey string, defaultValue float32, context gen.Context) (
	value float32,
	reason string,
	err error,
) {
	variant, reason, err := je.evaluateVariant(flagKey, context)
	if err != nil {
		log.Errorf("Error evaluating flag, %s", err.Error())
		return defaultValue, reason, err
	}

	val, ok := je.state.Flags[flagKey].Variants[variant].(float64)
	if !ok {
		log.Errorf("Error converting %s to float32", flagKey)
		return defaultValue, model.ErrorReason, errors.New(model.TypeMismatchErrorCode)
	}
	return float32(val), reason, nil
}

func (je *JSONEvaluator) ResolveObjectValue(flagKey string, defaultValue map[string]interface{}, context gen.Context) (
	value map[string]interface{},
	reason string,
	err error,
) {
	variant, reason, err := je.evaluateVariant(flagKey, context)
	if err != nil {
		log.Errorf("Error evaluating flag, %s", err.Error())
		return defaultValue, reason, err
	}

	val, ok := je.state.Flags[flagKey].Variants[variant].(map[string]interface{})
	if !ok {
		log.Errorf(fmt.Sprintf("Error converting %s to object", flagKey))
		return defaultValue, model.ErrorReason, errors.New(model.TypeMismatchErrorCode)
	}
	return val, reason, nil
}

// runs the rules (if defined) to determine the variant, otherwise falling through to the default
func (je *JSONEvaluator) evaluateVariant(
	flagKey string,
	context gen.Context,
) (variant string, reason string, err error) {
	flag, ok := je.state.Flags[flagKey]
	if !ok {
		// flag not found
		return "", model.ErrorReason, errors.New(model.FlagNotFoundErrorCode)
	}

	// get the targeting logic, if any
	targeting := flag.Targeting

	if targeting != nil {
		targetingBytes, err := targeting.MarshalJSON()
		if err != nil {
			log.Errorf("Error parsing rules for flag %s, %s", flagKey, err)
			return "", model.ErrorReason, err
		}
		contextBytes, err := context.MarshalJSON()
		if err != nil {
			log.Errorf("Error parsing context %s", err)
			return "", model.ErrorReason, err
		}

		var result bytes.Buffer

		// evaluate json-logic rules to determine the variant
		err = jsonlogic.Apply(bytes.NewReader(targetingBytes), bytes.NewReader(contextBytes), &result)
		if err != nil {
			log.Errorf("Error applying rules %s", err)
			return "", model.ErrorReason, err
		}
		// strip whitespace and quotes from the variant
		variant = strings.ReplaceAll(strings.TrimSpace(result.String()), "\"", "")
	}

	// if this is a valid variant, return it
	if _, ok := je.state.Flags[flagKey].Variants[variant]; ok {
		return variant, model.TargetingMatchReason, nil
	}

	// if it's not a valid variant, use the default (static) value
	return je.state.Flags[flagKey].DefaultVariant, model.StaticReason, nil
}

func MergeJSON(left, right []byte) ([]byte, error) {
	l := make(map[string]any)
	r := make(map[string]any)

	if len(left) > 0 {
		if err := json.Unmarshal(left, &l); err != nil {
			return nil, fmt.Errorf("invalid left: %w", err)
		}
	}

	if len(right) > 0 {
		if err := json.Unmarshal(right, &r); err != nil {
			return nil, fmt.Errorf("invalid right: %w", err)
		}
	}

	res := make(map[string]any)
	for k, v := range r {
		res[k] = v
	}
	for k, v := range l {
		res[k] = v
	}
	result, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("marshal result: %w", err)
	}

	return result, nil
}
