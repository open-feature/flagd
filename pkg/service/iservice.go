package service

import (
	"context"

	"github.com/open-feature/flagd/pkg/eval"
)

type NotificationType string

const (
	ConfigurationChange NotificationType = "configuration_change"
	ProviderReady       NotificationType = "provider_ready"
	KeepAlive           NotificationType = "keep_alive"
)

type Notification struct {
	Type NotificationType       `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type ReadinessProbe func() bool

type ServiceConfiguration struct {
	ReadinessProbe ReadinessProbe
}

/*
IService implementations define handlers for a particular transport, which call the IEvaluator implementation.
*/
type IService interface {
	Serve(ctx context.Context, eval eval.IEvaluator, svcConf ServiceConfiguration) error
	Notify(n Notification)
}
