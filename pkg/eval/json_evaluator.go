package eval

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/open-feature/flagd/pkg/model"
	log "github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
)

//go:embed flagd-definitions.json
var schema string

type JsonEvaluator struct {
	state Flags
}

func (je *JsonEvaluator) GetState () (string, error) {
	bytes, err := json.Marshal(&je.state)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (je *JsonEvaluator) SetState (state string) error {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	flagStringLoader := gojsonschema.NewStringLoader(state)
	result, err := gojsonschema.Validate(schemaLoader, flagStringLoader)
	if err != nil {
		return err
	} else if !result.Valid() {
		err := errors.New("Invalid JSON file.")
		log.Error(err)
		return err
	}
	json.Unmarshal([]byte(state), &je.state)
	return nil
}

func (je *JsonEvaluator) ResolveBooleanValue (flagKey string, defaultValue bool) (value bool, reason string, err error) {
  variant := je.state.Flags[flagKey].DefaultVariant
	val, ok := je.state.Flags[flagKey].Variants[variant].(bool);
	if (!ok) {
		log.Errorf("Error converting %s to bool", flagKey)
		return defaultValue, model.ErrorReason, errors.New(model.TypeMismatchErrorCode)
	}
	return val, model.StaticReason, nil
}

func (je *JsonEvaluator) ResolveStringValue (flagKey string, defaultValue string) (value string, reason string, err error) {
  variant := je.state.Flags[flagKey].DefaultVariant
	val, ok := je.state.Flags[flagKey].Variants[variant].(string);
	if (!ok) {
		log.Errorf("Error converting %s to string", flagKey)
		return defaultValue, model.ErrorReason, errors.New(model.TypeMismatchErrorCode)
	}
	return val, model.StaticReason, nil
}

func (je *JsonEvaluator) ResolveNumberValue (flagKey string, defaultValue float32) (value float32, reason string, err error) {
  variant := je.state.Flags[flagKey].DefaultVariant
	val, ok := je.state.Flags[flagKey].Variants[variant].(float64);
	if (!ok) {
		log.Errorf("Error converting %s to float32", flagKey)
		return defaultValue, model.ErrorReason, errors.New(model.TypeMismatchErrorCode)
	}
	return float32(val), model.StaticReason, nil
}

func (je *JsonEvaluator) ResolveObjectValue (flagKey string, defaultValue map[string]interface{}) (value map[string]interface{}, reason string, err error) {
  variant := je.state.Flags[flagKey].DefaultVariant
	val, ok := je.state.Flags[flagKey].Variants[variant].(map[string]interface{});
	if (!ok) {
		log.Errorf(fmt.Sprintf("Error converting %s to object", flagKey))
		return defaultValue, model.ErrorReason, errors.New(model.TypeMismatchErrorCode)
	}
	return val, model.StaticReason, nil
}

