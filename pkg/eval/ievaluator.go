package eval

import (
	"google.golang.org/protobuf/types/known/structpb"
)

type StateChangeNotificationType string

const (
	NotificationDelete StateChangeNotificationType = "delete"
	NotificationCreate StateChangeNotificationType = "write"
	NotificationUpdate StateChangeNotificationType = "update"
)

type StateChangeNotification struct {
	Type    StateChangeNotificationType `json:"type"`
	Source  string                      `json:"source"`
	FlagKey string                      `json:"flagKey"`
}

type AnyValue struct {
	Value   interface{}
	Variant string
	Reason  string
	FlagKey string
}

func NewAnyValue(value interface{}, variant string, reason string, flagKey string) AnyValue {
	return AnyValue{
		Value:   value,
		Variant: variant,
		Reason:  reason,
		FlagKey: flagKey,
	}
}

/*
IEvaluator implementations store the state of the flags,
do parsing and validation of the flag state and evaluate flags in response to handlers.
*/
type IEvaluator interface {
	GetState() (string, error)
	SetState(source string, state string) (map[string]interface{}, error)

	ResolveBooleanValue(
		reqID string,
		flagKey string,
		context *structpb.Struct) (value bool, variant string, reason string, err error)
	ResolveStringValue(
		reqID string,
		flagKey string,
		context *structpb.Struct) (value string, variant string, reason string, err error)
	ResolveIntValue(
		reqID string,
		flagKey string,
		context *structpb.Struct) (value int64, variant string, reason string, err error)
	ResolveFloatValue(
		reqID string,
		flagKey string,
		context *structpb.Struct) (value float64, variant string, reason string, err error)
	ResolveObjectValue(
		reqID string,
		flagKey string,
		context *structpb.Struct) (value map[string]any, variant string, reason string, err error)
	ResolveAllValues(
		reqID string,
		context *structpb.Struct) (values []AnyValue)
}
