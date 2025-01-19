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

	r, err := json.Marshal(ms)
	if err != nil {
		return "", fmt.Errorf("error marshaling json: %w", err)
	}

	return string(r), err
}
