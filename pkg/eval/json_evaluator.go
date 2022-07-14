package eval

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/diegoholiveira/jsonlogic/v3"
	"github.com/open-feature/flagd/pkg/model"
	log "github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
	"google.golang.org/protobuf/types/known/structpb"
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
func (je *JSONEvaluator) ResolveBooleanValue(flagKey string, context *structpb.Struct) (
	value bool,
	variant string,
	reason string,
	err error,
) {
	variant, reason, err = je.evaluateVariant(flagKey, context)
	if err != nil {
		log.Errorf("Error evaluating flag, %s", err.Error())
		return false, variant, reason, err
	}

	val, ok := je.state.Flags[flagKey].Variants[variant].(bool)
	if !ok {
		log.Errorf("Error converting %s to bool", flagKey)
		return false, variant, model.ErrorReason, errors.New(model.TypeMismatchErrorCode)
	}
	return val, variant, reason, nil
}

func (je *JSONEvaluator) ResolveStringValue(flagKey string, context *structpb.Struct) (
	value string,
	variant string,
	reason string,
	err error,
) {
	variant, reason, err = je.evaluateVariant(flagKey, context)
	if err != nil {
		log.Errorf("Error evaluating flag, %s", err.Error())
		return "", variant, reason, err
	}

	val, ok := je.state.Flags[flagKey].Variants[variant].(string)
	if !ok {
		log.Errorf("Error converting %s to string", flagKey)
		return "", variant, model.ErrorReason, errors.New(model.TypeMismatchErrorCode)
	}
	return val, variant, reason, nil
}

func (je *JSONEvaluator) ResolveNumberValue(flagKey string, context *structpb.Struct) (
	value float32,
	variant string,
	reason string,
	err error,
) {
	variant, reason, err = je.evaluateVariant(flagKey, context)
	if err != nil {
		log.Errorf("Error evaluating flag, %s", err.Error())
		return 0, variant, reason, err
	}

	val, ok := je.state.Flags[flagKey].Variants[variant].(float64)
	if !ok {
		log.Errorf("Error converting %s to float32", flagKey)
		return 0, variant, model.ErrorReason, errors.New(model.TypeMismatchErrorCode)
	}
	return float32(val), variant, reason, nil
}

func (je *JSONEvaluator) ResolveObjectValue(flagKey string, context *structpb.Struct) (
	value map[string]interface{},
	variant string,
	reason string,
	err error,
) {
	variant, reason, err = je.evaluateVariant(flagKey, context)
	if err != nil {
		log.Errorf("Error evaluating flag, %s", err.Error())
		return nil, variant, reason, err
	}

	val, ok := je.state.Flags[flagKey].Variants[variant].(map[string]interface{})
	if !ok {
		log.Errorf(fmt.Sprintf("Error converting %s to object", flagKey))
		return nil, variant, model.ErrorReason, errors.New(model.TypeMismatchErrorCode)
	}
	return val, variant, reason, nil
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

	// get the targeting logic, if any
	targeting := flag.Targeting

	if targeting != nil {
		targetingBytes, err := targeting.MarshalJSON()
		if err != nil {
			log.Errorf("Error parsing rules for flag %s, %s", flagKey, err)
			return "", model.ErrorReason, err
		}

		b, err := json.Marshal(context)
		if err != nil {
			log.Errorf("error parsing context for flag %s, %s, %v", flagKey, err, context)

			return "", model.ErrorReason, errors.New(model.ErrorReason)
		}
		var result bytes.Buffer
		// evaluate json-logic rules to determine the variant
		err = jsonlogic.Apply(bytes.NewReader(targetingBytes), bytes.NewReader(b), &result)
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
