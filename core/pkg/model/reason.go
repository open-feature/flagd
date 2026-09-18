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
	// StaleReason indicates the flag was resolved from a store whose sync source is
	// currently disconnected, so the value may no longer reflect the source of truth.
	// See https://openfeature.dev/specification/types#resolution-details
	StaleReason = "STALE"
	// only used internally if no default value could be determined
	// will be translated to DefaultReason in the API response
	FallbackReason = "FALLBACK"
)
