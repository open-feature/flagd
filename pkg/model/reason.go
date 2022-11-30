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
)
