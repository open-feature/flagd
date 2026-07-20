package utils

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

// converts YAML byte array to JSON string
func YAMLToJSON(rawFile []byte) (string, error) {
	if len(rawFile) == 0 {
		return "", nil
	}

	var ms map[string]interface{}
	if err := yaml.Unmarshal(rawFile, &ms); err != nil {
		return "", fmt.Errorf("error unmarshaling yaml: %w", err)
	}

	r, err := json.Marshal(convertMapKeys(ms))
	if err != nil {
		return "", fmt.Errorf("error marshaling json: %w", err)
	}

	return string(r), err
}

// convertMapKeys recursively converts any map[interface{}]interface{} into
// map[string]interface{} so it can be marshaled by encoding/json. gopkg.in/yaml.v3
// decodes a nested mapping that has a non-string key (e.g. an unquoted numeric key
// in an object variant) as map[interface{}]interface{}, which json.Marshal cannot
// handle; stringifying the keys mirrors how YAML-to-JSON conversion is expected to
// behave and preserves YAML 1.2 value semantics.
func convertMapKeys(in interface{}) interface{} {
	switch v := in.(type) {
	case map[string]interface{}:
		for key, val := range v {
			v[key] = convertMapKeys(val)
		}
		return v
	case map[interface{}]interface{}:
		out := make(map[string]interface{}, len(v))
		for key, val := range v {
			out[fmt.Sprintf("%v", key)] = convertMapKeys(val)
		}
		return out
	case []interface{}:
		for i, val := range v {
			v[i] = convertMapKeys(val)
		}
		return v
	default:
		return in
	}
}
