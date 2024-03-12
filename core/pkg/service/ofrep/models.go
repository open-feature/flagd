package ofrep

import (
	"fmt"

	"github.com/open-feature/flagd/core/pkg/evaluator"
	"github.com/open-feature/flagd/core/pkg/model"
)

const (
	InvalidContextReason = "INVALID_CONTEXT"
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
	flags := make([]interface{}, 0)

	for _, value := range values {
		if value.Error != nil {
			_, evaluationError := ErrorResponseAndStatus(value)
			flags = append(flags, evaluationError)
		} else {
			flags = append(flags, SuccessResponseFrom(value))
		}
	}

	return BulkEvaluationResponse{
		flags,
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

func ContextErrorResponseAndStatus(key string) EvaluationError {
	return EvaluationError{
		Key:          key,
		ErrorCode:    InvalidContextReason,
		ErrorDetails: "provided context is invalid",
	}
}

func ContextErrorResponseForBulkAndStatus() BulkEvaluationError {
	return BulkEvaluationError{
		ErrorCode:    InvalidContextReason,
		ErrorDetails: "provided context is invalid",
	}
}

func ErrorResponseAndStatus(result evaluator.AnyValue) (int, EvaluationError) {
	payload := EvaluationError{
		Key: result.FlagKey,
	}

	status := 400

	switch result.Error.Error() {
	case model.FlagNotFoundErrorCode:
		status = 404
		payload.ErrorCode = model.FlagNotFoundErrorCode
		payload.ErrorDetails = fmt.Sprintf("flag `%s` does not exisit", result.FlagKey)
	case model.ParseErrorCode:
		payload.ErrorCode = model.ParseErrorCode
		payload.ErrorDetails = "error parsing the flag"
	case model.GeneralErrorCode:
	default:
		payload.ErrorCode = model.GeneralErrorCode
		payload.ErrorDetails = "error processing the flag for evaluation"
	}

	return status, payload
}
