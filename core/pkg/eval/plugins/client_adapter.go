package plugins

import (
	"context"
	"fmt"

	"github.com/open-feature/flagd/core/pkg/eval"
	"github.com/open-feature/flagd/core/pkg/sync"
	"github.com/open-feature/go-sdk/pkg/openfeature"
	"google.golang.org/protobuf/types/known/structpb"
)

type OpenFeatureClientEvaluatorAdapter struct {
	client *openfeature.Client
}

func fromClient(client *openfeature.Client) eval.IEvaluator {
	return OpenFeatureClientEvaluatorAdapter{client}
}

func (of OpenFeatureClientEvaluatorAdapter) ResolveBooleanValue(reqID string, flagKey string, ctx *structpb.Struct) (
	value bool,
	variant string,
	reason string,
	err error,
) {
	details, err := of.client.BooleanValueDetails(context.Background(), flagKey, false, openfeature.EvaluationContext{})

	if err != nil {
		return
	}

	if details.Reason == openfeature.DefaultReason {
		err = fmt.Errorf("Specific flag could not be evaluated")
		return
	}

	value = details.Value
	variant = details.Variant
	reason = string(details.Reason)

	return
}

func (of OpenFeatureClientEvaluatorAdapter) ResolveStringValue(reqID string, flagKey string, ctx *structpb.Struct) (
	value string,
	variant string,
	reason string,
	err error,
) {
	details, err := of.client.StringValueDetails(context.Background(), flagKey, "", openfeature.EvaluationContext{})

	if err != nil {
		return
	}

	if details.Reason == openfeature.DefaultReason {
		err = fmt.Errorf("Specific flag could not be evaluated")
		return
	}

	value = details.Value
	variant = details.Variant
	reason = string(details.Reason)

	return
}

func (of OpenFeatureClientEvaluatorAdapter) ResolveFloatValue(reqID string, flagKey string, ctx *structpb.Struct) (
	value float64,
	variant string,
	reason string,
	err error,
) {
	details, err := of.client.FloatValueDetails(context.Background(), flagKey, 0.0, openfeature.EvaluationContext{})

	if err != nil {
		return
	}

	if details.Reason == openfeature.DefaultReason {
		err = fmt.Errorf("Specific flag could not be evaluated")
		return
	}

	value = details.Value
	variant = details.Variant
	reason = string(details.Reason)

	return
}

func (of OpenFeatureClientEvaluatorAdapter) ResolveIntValue(reqID string, flagKey string, ctx *structpb.Struct) (
	value int64,
	variant string,
	reason string,
	err error,
) {
	details, err := of.client.IntValueDetails(context.Background(), flagKey, 0, openfeature.EvaluationContext{})

	if err != nil {
		return
	}

	if details.Reason == openfeature.DefaultReason {
		err = fmt.Errorf("Specific flag could not be evaluated")
		return
	}

	value = details.Value
	variant = details.Variant
	reason = string(details.Reason)

	return
}

func (of OpenFeatureClientEvaluatorAdapter) ResolveObjectValue(reqID string, flagKey string, ctx *structpb.Struct) (
	value map[string]any,
	variant string,
	reason string,
	err error,
) {
	details, err := of.client.ObjectValueDetails(context.Background(), flagKey, nil, openfeature.EvaluationContext{})

	if err != nil {
		return
	}

	if details.Reason == openfeature.DefaultReason {
		err = fmt.Errorf("Specific flag could not be evaluated")
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

func (eval OpenFeatureClientEvaluatorAdapter) GetState() (string, error) {
	return "", fmt.Errorf("Unsupported action for OpenFeature Client adapters")
}

func (of OpenFeatureClientEvaluatorAdapter) SetState(payload sync.DataSync) (map[string]interface{}, bool, error) {
	return map[string]interface{}{}, false, fmt.Errorf("Unsupported action for OpenFeature Client adapters")
}

func (of OpenFeatureClientEvaluatorAdapter) ResolveAllValues(reqID string, ctx *structpb.Struct) []eval.AnyValue {
	// not truly implemented for OpenFeature adapters

	return []eval.AnyValue{}
}
