package evaluator

import (
	"context"

	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/sync"
)

type AnyValue struct {
	Value    interface{}
	Variant  string
	Reason   string
	FlagKey  string
	Metadata model.Metadata
	Error    error
}

func NewAnyValue(
	value interface{}, variant string, reason string, flagKey string, metadata model.Metadata,
	err error,
) AnyValue {
	return AnyValue{
		Value:    value,
		Variant:  variant,
		Reason:   reason,
		FlagKey:  flagKey,
		Metadata: metadata,
		Error:    err,
	}
}

/*
IEvaluator is an extension of IResolver, allowing storage updates and retrievals
*/
type IEvaluator interface {
	GetState() (string, error)
	SetState(payload sync.DataSync) (model.Metadata, bool, error)
	IResolver
}

// IResolver focuses on resolving of the known flags
type IResolver interface {
	ResolveBooleanValue(
		ctx context.Context,
		reqID string,
		flagKey string,
		context map[string]any) (value bool, variant string, reason string, metadata model.Metadata, err error)
	ResolveStringValue(
		ctx context.Context,
		reqID string,
		flagKey string,
		context map[string]any) (
		value string, variant string, reason string, metadata model.Metadata, err error)
	ResolveIntValue(
		ctx context.Context,
		reqID string,
		flagKey string,
		context map[string]any) (
		value int64, variant string, reason string, metadata model.Metadata, err error)
	ResolveFloatValue(
		ctx context.Context,
		reqID string,
		flagKey string,
		context map[string]any) (
		value float64, variant string, reason string, metadata model.Metadata, err error)
	ResolveObjectValue(
		ctx context.Context,
		reqID string,
		flagKey string,
		context map[string]any) (
		value map[string]any, variant string, reason string, metadata model.Metadata, err error)
	ResolveAsAnyValue(
		ctx context.Context,
		reqID string,
		flagKey string,
		context map[string]any) AnyValue
	ResolveAllValues(
		ctx context.Context,
		reqID string,
		context map[string]any) (resolutions []AnyValue, metadata model.Metadata, err error)
}
