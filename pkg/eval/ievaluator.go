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

/*
IEvaluator implementations store the state of the flags,
do parsing and validation of the flag state and evaluate flags in response to handlers.
*/
type IEvaluator interface {
	GetState() (string, error)
	SetState(source string, state string) ([]StateChangeNotification, error)

	ResolveBooleanValue(
		flagKey string,
		context *structpb.Struct) (value bool, variant string, reason string, err error)
	ResolveStringValue(
		flagKey string,
		context *structpb.Struct) (value string, variant string, reason string, err error)
	ResolveIntValue(flagKey string,
		context *structpb.Struct) (value int64, variant string, reason string, err error)
	ResolveFloatValue(flagKey string,
		context *structpb.Struct) (value float64, variant string, reason string, err error)
	ResolveObjectValue(
		flagKey string,
		context *structpb.Struct) (value map[string]any, variant string, reasons string, err error)
}

func (s *StateChangeNotification) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type":    string(s.Type),
		"source":  s.Source,
		"flagKey": s.FlagKey,
	}
}
