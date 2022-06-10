package eval

import (
	_ "embed"
	"encoding/json"
	"errors"

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

func (je *JsonEvaluator) ResolveBooleanValue (flagKey string, defaultValue bool) (bool, error) {
	var variant = je.state.BooleanFlags[flagKey].DefaultVariant
	return je.state.BooleanFlags[flagKey].Variants[variant], nil
}

func (je *JsonEvaluator) ResolveStringValue (flagKey string, defaultValue string) (string, error) {
	var variant = je.state.StringFlags[flagKey].DefaultVariant
	return je.state.StringFlags[flagKey].Variants[variant], nil
}

func (je *JsonEvaluator) ResolveNumberValue (flagKey string, defaultValue float32) (float32, error) {
	var variant = je.state.NumericFlags[flagKey].DefaultVariant
	return je.state.NumericFlags[flagKey].Variants[variant], nil
}

func (je *JsonEvaluator) ResolveObjectValue (flagKey string, defaultValue map[string]interface{}) (map[string]interface{}, error) {
	var variant = je.state.ObjectFlags[flagKey].DefaultVariant
	return je.state.ObjectFlags[flagKey].Variants[variant], nil
}

