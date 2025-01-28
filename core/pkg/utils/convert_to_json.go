package utils

import (
	"fmt"
	"mime"
	"regexp"
	"strings"
)

var alphanumericRegex = regexp.MustCompile("[^a-zA-Z0-9]+")

// ConvertToJSON attempts to convert the content of a file to JSON based on the file extension.
// The media type is used as a fallback in case the file extension is not recognized.
func ConvertToJSON(data []byte, fileExtension string, mediaType string) (string, error) {
	var detectedType string
	if fileExtension != "" {
		// file extension only contains alphanumeric characters
		detectedType = alphanumericRegex.ReplaceAllString(fileExtension, "")
	} else {
		parsedMediaType, _, err := mime.ParseMediaType(mediaType)
		if err != nil {
			return "", fmt.Errorf("unable to determine file format: %w", err)
		}
		detectedType = parsedMediaType
	}

	// Normalize the detected type
	detectedType = strings.ToLower(detectedType)

	switch detectedType {
	case "yaml", "yml", "application/yaml", "application/x-yaml":
		str, err := YAMLToJSON(data)
		if err != nil {
			return "", fmt.Errorf("error converting blob from yaml to json: %w", err)
		}
		return str, nil
	case "json", "application/json":
		return string(data), nil
	default:
		return "", fmt.Errorf("unsupported file format: '%s'", detectedType)
	}
}
