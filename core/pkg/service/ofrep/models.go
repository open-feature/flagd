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
	Value    interface{} `json:"value"`
	Key      string      `json:"key"`
	Reason   string      `json:"reason"`
	Variant  string      `json:"variant"`
	Metadata interface{} `json:"metadata"`
}

type BulkEvaluationResponse struct {
	Flags []interface{} `json:"flags,omitempty"`
}

type EvaluationError struct {
	Key          string `json:"key"`
	ErrorCode    string `json:"errorCode"`
	ErrorDetails string `json:"errorDetails"`
}

type BulkEvaluationError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorDetails string `json:"errorDetails"`
}

type InternalError struct {
	ErrorDetails string `json:"errorDetails"`
}

func BulkEvaluationResponseFrom(values []evaluator.AnyValue) BulkEvaluationResponse {
	evaluations := make([]interface{}, 0)

	for _, value := range values {
		if value.Error != nil {
			_, evaluationError := EvaluationErrorResponseFrom(value)
			evaluations = append(evaluations, evaluationError)
		} else {
			evaluations = append(evaluations, SuccessResponseFrom(value))
		}
	}

	return BulkEvaluationResponse{
		evaluations,
	}
}

func SuccessResponseFrom(result evaluator.AnyValue) EvaluationSuccess {
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

func BulkEvaluationContextErrorFrom() BulkEvaluationError {
	return BulkEvaluationError{
		ErrorCode:    model.InvalidContextCode,
		ErrorDetails: "Provider context is not valid",
	}
}

func EvaluationErrorResponseFrom(result evaluator.AnyValue) (int, EvaluationError) {
	payload := EvaluationError{
		Key: result.FlagKey,
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
