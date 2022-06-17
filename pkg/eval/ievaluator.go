package eval

/*
IEvaluator implementations store the state of the flags, do parsing and validation of the flag state and evaluate flags in response to handlers.
*/
type IEvaluator interface {
	// TODO: add context param when rule evaluator is implemented
	GetState() (string, error)
	SetState(state string) error
	ResolveBooleanValue(flagKey string, defaultValue bool) (bool, string, error)
	ResolveStringValue(flagKey string, defaultValue string) (string, string, error)
	ResolveNumberValue(flagKey string, defaultValue float32) (float32, string, error)
	ResolveObjectValue(flagKey string, defaultValue map[string]interface{}) (map[string]interface{}, string, error)
}