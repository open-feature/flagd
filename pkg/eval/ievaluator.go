package eval

import (
	"google.golang.org/protobuf/types/known/structpb"
)

/*
IEvaluator implementations store the state of the flags,
do parsing and validation of the flag state and evaluate flags in response to handlers.
*/
type IEvaluator interface {
	GetState() (string, error)
	SetState(state string) error
	ResolveBooleanValue(
		flagKey string,
		context *structpb.Struct) (value bool, reason string, err error)
	ResolveStringValue(
		flagKey string,
		context *structpb.Struct) (value string, reason string, err error)
	ResolveNumberValue(flagKey string,
		context *structpb.Struct) (value float32, reason string, err error)
	ResolveObjectValue(
		flagKey string,
		context *structpb.Struct) (value map[string]interface{}, reasons string, err error)
}
