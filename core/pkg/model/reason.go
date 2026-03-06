package model

type EvaluationReason string

const (
	TargetingMatchReason = "TARGETING_MATCH"
	SplitReason          = "SPLIT"
	DisabledReason       = "DISABLED"
	DefaultReason        = "DEFAULT"
	UnknownReason        = "UNKNOWN"
	ErrorReason          = "ERROR"
	StaticReason         = "STATIC"
	// only used internally if no default value could be determined
	// will be translated to DefaultReason in the API response
	FallbackReason = "FALLBACK"
)
