package service

import (
	"context"
	"sort"
	"strings"

	"github.com/open-feature/flagd/core/pkg/evaluator"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
)

func ResolveAllWithSelectorMerge(
	ctx context.Context,
	reqID string,
	eval evaluator.IEvaluator,
	evaluationContext map[string]any,
	selectorExpression string,
) ([]evaluator.AnyValue, model.Metadata, error) {
	selectors := splitSelectorExpression(selectorExpression)

	switch len(selectors) {
	case 0:
		return eval.ResolveAllValues(ctx, reqID, evaluationContext)
	case 1:
		selector := store.NewSelector(selectors[0])
		selectorCtx := context.WithValue(ctx, store.SelectorContextKey{}, selector)
		return eval.ResolveAllValues(selectorCtx, reqID, evaluationContext)
	default:
		mergedValues := map[string]evaluator.AnyValue{}
		mergedMetadata := model.Metadata{}

		for _, selectorExpression := range selectors {
			selector := store.NewSelector(selectorExpression)
			selectorCtx := context.WithValue(ctx, store.SelectorContextKey{}, selector)
			values, metadata, err := eval.ResolveAllValues(selectorCtx, reqID, evaluationContext)
			if err != nil {
				return nil, nil, err
			}

			for key, value := range metadata {
				mergedMetadata[key] = value
			}
			for _, value := range values {
				mergedValues[value.FlagKey] = value
			}
		}

		keys := make([]string, 0, len(mergedValues))
		for key := range mergedValues {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		resolutions := make([]evaluator.AnyValue, 0, len(keys))
		for _, key := range keys {
			resolutions = append(resolutions, mergedValues[key])
		}

		return resolutions, mergedMetadata, nil
	}
}

func splitSelectorExpression(selectorExpression string) []string {
	if strings.TrimSpace(selectorExpression) == "" {
		return nil
	}

	parts := strings.Split(selectorExpression, ",")
	selectors := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		selectors = append(selectors, trimmed)
	}
	return selectors
}
