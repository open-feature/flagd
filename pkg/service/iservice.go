package service

import (
	"context"

	"github.com/open-feature/flagd/pkg/eval"
)

type NotificationType string

const (
	ConfigurationChange NotificationType = "configuration_change"
	ProviderReady       NotificationType = "provider_ready"
)

type IServiceConfiguration interface{}

/*
IService implementations define handlers for a particular transport, which call the IEvaluator implementation.
*/
type IService interface {
	Serve(ctx context.Context, eval eval.IEvaluator) error
	Notify(n NotificationType)
}
