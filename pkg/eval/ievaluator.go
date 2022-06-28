package eval

import gen "github.com/open-feature/flagd/pkg/generated"

/*
IEvaluator implementations store the state of the flags,
do parsing and validation of the flag state and evaluate flags in response to handlers.
*/
type IEvaluator interface {
	GetState() (string, error)
	SetState(state string) error
	ResolveBooleanValue(flagKey string, defaultValue bool, context gen.Context) (value bool, reason string, err error)
	ResolveStringValue(flagKey string, defaultValue string, context gen.Context) (value string, reason string, err error)
	ResolveNumberValue(flagKey string, defaultValue float32, context gen.Context) (value float32, reason string, err error)
	ResolveObjectValue(
		flagKey string,
		defaultValue map[string]interface{},
		context gen.Context) (value map[string]interface{}, reasons string, err error)
}
