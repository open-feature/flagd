package plugins

import (
	"context"
	"fmt"

	"github.com/open-feature/flagd/core/pkg/eval"
	"github.com/open-feature/flagd/core/pkg/sync"
	"github.com/open-feature/go-sdk/pkg/openfeature"
	"google.golang.org/protobuf/types/known/structpb"
)

type OpenFeatureProviderEvaluatorAdapter struct {
	provider openfeature.FeatureProvider
}

func fromProvider(provider openfeature.FeatureProvider) eval.IEvaluator {
	// not truly implemented for OpenFeature adapters

	return OpenFeatureProviderEvaluatorAdapter{provider}
}

func (of OpenFeatureProviderEvaluatorAdapter) ResolveBooleanValue(reqID string, flagKey string, ctx *structpb.Struct) (
	value bool,
	variant string,
	reason string,
	err error,
) {
	details := of.provider.BooleanEvaluation(context.Background(), flagKey, false, openfeature.FlattenedContext{})

	err = details.Error()
	if err != nil {
		return
	}

	value = details.Value
	variant = details.Variant
	reason = string(details.Reason)

	return
}

func (of OpenFeatureProviderEvaluatorAdapter) ResolveStringValue(reqID string, flagKey string, ctx *structpb.Struct) (
	value string,
	variant string,
	reason string,
	err error,
) {
	details := of.provider.StringEvaluation(context.Background(), flagKey, "", openfeature.FlattenedContext{})

	err = details.Error()
	if err != nil {
		return
	}

	value = details.Value
	variant = details.Variant
	reason = string(details.Reason)

	return
}

func (of OpenFeatureProviderEvaluatorAdapter) ResolveFloatValue(reqID string, flagKey string, ctx *structpb.Struct) (
	value float64,
	variant string,
	reason string,
	err error,
) {
	details := of.provider.FloatEvaluation(context.Background(), flagKey, 0.0, openfeature.FlattenedContext{})

	err = details.Error()
	if err != nil {
		return
	}

	value = details.Value
	variant = details.Variant
	reason = string(details.Reason)

	return
}

func (of OpenFeatureProviderEvaluatorAdapter) ResolveIntValue(reqID string, flagKey string, ctx *structpb.Struct) (
	value int64,
	variant string,
	reason string,
	err error,
) {
	details := of.provider.IntEvaluation(context.Background(), flagKey, 0, openfeature.FlattenedContext{})

	err = details.Error()
	if err != nil {
		return
	}

	value = details.Value
	variant = details.Variant
	reason = string(details.Reason)

	return
}

func (of OpenFeatureProviderEvaluatorAdapter) ResolveObjectValue(reqID string, flagKey string, ctx *structpb.Struct) (
	value map[string]any,
	variant string,
	reason string,
	err error,
) {
	details := of.provider.ObjectEvaluation(context.Background(), flagKey, nil, openfeature.FlattenedContext{})

	err = details.Error()
	if err != nil {
		return
	}

	value, ok := details.Value.(map[string]any)

	if !ok {
		err = fmt.Errorf("Invalid map type returned for object evaluation")
		return
	}
	variant = details.Variant
	reason = string(details.Reason)

	return
}

//// Unsupported actions for OpenFeature adapters

func (eval OpenFeatureProviderEvaluatorAdapter) GetState() (string, error) {
	return "", fmt.Errorf("Unsupported action for OpenFeature Provider adapters")
}

func (of OpenFeatureProviderEvaluatorAdapter) SetState(payload sync.DataSync) (map[string]interface{}, bool, error) {
	return map[string]interface{}{}, false, fmt.Errorf("Unsupported action for OpenFeature Provider adapters")
}

func (of OpenFeatureProviderEvaluatorAdapter) ResolveAllValues(reqID string, ctx *structpb.Struct) []eval.AnyValue {
	return []eval.AnyValue{}
}
