package eval

import (
	"context"

	"github.com/open-feature/flagd/core/pkg/sync"
	"google.golang.org/protobuf/types/known/structpb"
)

type AnyValue struct {
	Value   interface{}
	Variant string
	Reason  string
	FlagKey string
	Error   error
}

func NewAnyValue(value interface{}, variant string, reason string, flagKey string, err error) AnyValue {
	return AnyValue{
		Value:   value,
		Variant: variant,
		Reason:  reason,
		FlagKey: flagKey,
		Error:   err,
	}
}

/*
IEvaluator implementations store the state of the flags,
do parsing and validation of the flag state and evaluate flags in response to handlers.
*/
type IEvaluator interface {
	GetState() (string, error)
	SetState(payload sync.DataSync) (map[string]interface{}, bool, error)

	ResolveBooleanValue(
		ctx context.Context,
		reqID string,
		flagKey string,
		context *structpb.Struct) (value bool, variant string, reason string, err error)
	ResolveStringValue(
		ctx context.Context,
		reqID string,
		flagKey string,
		context *structpb.Struct) (value string, variant string, reason string, err error)
	ResolveIntValue(
		ctx context.Context,
		reqID string,
		flagKey string,
		context *structpb.Struct) (value int64, variant string, reason string, err error)
	ResolveFloatValue(
		ctx context.Context,
		reqID string,
		flagKey string,
		context *structpb.Struct) (value float64, variant string, reason string, err error)
	ResolveObjectValue(
		ctx context.Context,
		reqID string,
		flagKey string,
		context *structpb.Struct) (value map[string]any, variant string, reason string, err error)
	ResolveAllValues(
		ctx context.Context,
		reqID string,
		context *structpb.Struct) (values []AnyValue)
}
