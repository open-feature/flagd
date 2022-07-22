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
	data, err := json.Marshal(&je.state)
	if err != nil {
		return "", err
	}
	return string(data), nil
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

func resolve[T any](key string, context *structpb.Struct,
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

func (je *JSONEvaluator) ResolveNumberValue(flagKey string, context *structpb.Struct) (
	value float32,
	variant string,
	reason string,
	err error,
) {
	var val float64
	val, variant, reason, err = resolve[float64](flagKey, context, je.evaluateVariant, je.state.Flags[flagKey].Variants)
	value = float32(val)
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
