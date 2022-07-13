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
		context *structpb.Struct) (value bool, variant string, reason string, err error)
	ResolveStringValue(
		flagKey string,
		context *structpb.Struct) (value string, variant string, reason string, err error)
	ResolveNumberValue(flagKey string,
		context *structpb.Struct) (value float32, variant string, reason string, err error)
	ResolveObjectValue(
		flagKey string,
		context *structpb.Struct) (value map[string]interface{}, variant string, reasons string, err error)
}
