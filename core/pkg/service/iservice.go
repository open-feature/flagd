package service

import (
	"context"

	"github.com/bufbuild/connect-go"
)

type NotificationType string

const (
	ConfigurationChange NotificationType = "configuration_change"
	Shutdown            NotificationType = "provider_shutdown"
	ProviderReady       NotificationType = "provider_ready"
	KeepAlive           NotificationType = "keep_alive"
)

type Notification struct {
	Type NotificationType       `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type ReadinessProbe func() bool

type Configuration struct {
	ReadinessProbe ReadinessProbe
	Port           uint16
	MetricsPort    uint16
	ServiceName    string
	CertPath       string
	KeyPath        string
	SocketPath     string
	CORS           []string
	Options        []connect.HandlerOption
}

/*
IFlagEvaluationService implementations define handlers for a particular transport,
which call the IEvaluator implementation.
*/
type IFlagEvaluationService interface {
	Serve(ctx context.Context, svcConf Configuration) error
	Notify(n Notification)
	Shutdown()
}

/*
IFlagEvaluationService implementations define handlers for a particular transport,
which call the IEvaluator implementation.
*/
type IKubeSyncService interface {
	Serve(ctx context.Context, svcConf Configuration) error
}
