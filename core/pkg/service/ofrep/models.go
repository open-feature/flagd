package ofrep

import (
	"fmt"

	"github.com/open-feature/flagd/core/pkg/evaluator"
	"github.com/open-feature/flagd/core/pkg/model"
)

type Request struct {
	Context interface{} `json:"context"`
}

type EvaluationSuccess struct {
	Value    interface{}    `json:"value,omitempty"`
	Key      string         `json:"key"`
	Reason   string         `json:"reason"`
	Variant  string         `json:"variant,omitempty"`
	Metadata model.Metadata `json:"metadata"`
}

type BulkEvaluationResponse struct {
	Flags    []interface{}  `json:"flags"`
	Metadata model.Metadata `json:"metadata"`
}

type EvaluationError struct {
	Key          string         `json:"key"`
	ErrorCode    string         `json:"errorCode"`
	ErrorDetails string         `json:"errorDetails"`
	Metadata     model.Metadata `json:"metadata"`
}

type BulkEvaluationError struct {
	ErrorCode    string         `json:"errorCode"`
	ErrorDetails string         `json:"errorDetails"`
	Metadata     model.Metadata `json:"metadata"`
}

type InternalError struct {
	ErrorDetails string `json:"errorDetails"`
}

func BulkEvaluationResponseFrom(resolutions []evaluator.AnyValue, metadata model.Metadata) BulkEvaluationResponse {
	evaluations := make([]interface{}, 0)

	for _, value := range resolutions {
		if value.Error != nil {
			_, evaluationError := EvaluationErrorResponseFrom(value)
			evaluations = append(evaluations, evaluationError)
		} else {
			evaluations = append(evaluations, SuccessResponseFrom(value))
		}
	}

	return BulkEvaluationResponse{
		evaluations,
		metadata,
	}
}

func SuccessResponseFrom(result evaluator.AnyValue) EvaluationSuccess {
	// if reason is fallback, we want to omit the value and variant from the response, and set reason to default
	if result.Reason == model.FallbackReason {
		return EvaluationSuccess{
			Value:    nil, // not marshalled due to omitempty
			Key:      result.FlagKey,
			Reason:   model.DefaultReason,
			Variant:  "", // not marshalled due to omitempty
			Metadata: result.Metadata,
		}
	}
	return EvaluationSuccess{
		Value:    result.Value,
		Key:      result.FlagKey,
		Reason:   result.Reason,
		Variant:  result.Variant,
		Metadata: result.Metadata,
	}
}

func ContextErrorResponseFrom(key string) EvaluationError {
	return EvaluationError{
		Key:          key,
		ErrorCode:    model.InvalidContextCode,
		ErrorDetails: "Provider context is not valid",
	}
}

func BulkEvaluationContextError() BulkEvaluationError {
	return BulkEvaluationError{
		ErrorCode:    model.InvalidContextCode,
		ErrorDetails: "Provider context is not valid",
	}
}

func BulkEvaluationContextErrorFrom(code string, details string) BulkEvaluationError {
	return BulkEvaluationError{
		ErrorCode:    code,
		ErrorDetails: details,
	}
}

func EvaluationErrorResponseFrom(result evaluator.AnyValue) (int, EvaluationError) {
	payload := EvaluationError{
		Key:      result.FlagKey,
		Metadata: result.Metadata,
	}

	status := 400

	switch result.Error.Error() {
	case model.FlagNotFoundErrorCode:
		status = 404
		payload.ErrorCode = model.FlagNotFoundErrorCode
		payload.ErrorDetails = fmt.Sprintf("flag `%s` does not exist", result.FlagKey)
	case model.FlagDisabledErrorCode:
		status = 404
		payload.ErrorCode = model.FlagNotFoundErrorCode
		payload.ErrorDetails = fmt.Sprintf("flag `%s` is disabled", result.FlagKey)
	case model.ParseErrorCode:
		payload.ErrorCode = model.ParseErrorCode
		payload.ErrorDetails = fmt.Sprintf("error parsing the flag `%s`", result.FlagKey)
	case model.GeneralErrorCode:
		fallthrough
	default:
		payload.ErrorCode = model.GeneralErrorCode
		payload.ErrorDetails = "error processing the flag for evaluation"
	}

	return status, payload
}
