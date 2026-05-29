package model

import "fmt"

const (
	FlagNotFoundErrorCode = "FLAG_NOT_FOUND"
	ParseErrorCode        = "PARSE_ERROR"
	TypeMismatchErrorCode = "TYPE_MISMATCH"
	GeneralErrorCode      = "GENERAL"
	FlagDisabledErrorCode = "FLAG_DISABLED"
	InvalidContextCode    = "INVALID_CONTEXT"
)

var ReadableErrorMessage = map[string]string{
	FlagNotFoundErrorCode: "Flag not found",
	ParseErrorCode:        "Error parsing input or configuration",
	TypeMismatchErrorCode: "Type mismatch error",
	GeneralErrorCode:      "General error",
	FlagDisabledErrorCode: "Flag is disabled",
	InvalidContextCode:    "Invalid context provided",
}

func GetErrorMessage(code string) string {
	if msg, exists := ReadableErrorMessage[code]; exists {
		return msg
	}
	return fmt.Sprintf("Unknown error code: %s", code)
}
