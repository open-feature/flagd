package plugins

import (
	"fmt"
	"plugin"

	"github.com/open-feature/flagd/core/pkg/eval"
	"github.com/open-feature/go-sdk/pkg/openfeature"
)

// FromPath opens an evaluator plugin from a target path
// The Evaluator must be defined as the public export "Evaluator"
// If the "Evaluator" export is defined, it must be an IEvaluator,
// OpenFeature Client, or OpenFeature Provider.
func FromPath(path string) (*eval.IEvaluator, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("Invalid path for plugin")
	}

	p, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	plugin, err := p.Lookup("Evaluator")
	if err != nil {
		return nil, err
	}

	typeEvaluator, ok := plugin.(*eval.IEvaluator)
	if ok {
		return typeEvaluator, nil
	}

	typeProvider, ok := plugin.(*openfeature.FeatureProvider)
	if ok {
		eval := fromProvider(*typeProvider)

		return &eval, nil
	}

	typeClient, ok := plugin.(*openfeature.Client)
	if ok {
		eval := fromClient(typeClient)

		return &eval, nil
	}

	return nil, fmt.Errorf("Failed to load plugin as evaluator")
}
