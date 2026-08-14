package service

import (
	"context"
	"time"

	"connectrpc.com/connect"
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
	ReadinessProbe             ReadinessProbe
	Port                       uint16
	ManagementPort             uint16
	ServiceName                string
	CertPath                   string
	KeyPath                    string
	SocketPath                 string
	CORS                       []string
	Options                    []connect.HandlerOption
	ContextValues              map[string]any
	HeaderToContextKeyMappings map[string]string
	StreamDeadline             time.Duration
	MaxRequestHeaderBytes      int64
	MaxRequestBodyBytes        int64

	// KeepAliveMinTime is the minimum interval the gRPC server permits between
	// client keepalive pings. Pings arriving more frequently are rejected with
	// GOAWAY (ENHANCE_YOUR_CALM). Honoured by services exposing a gRPC sync
	// server; ignored by the connect-based flag evaluation service.
	KeepAliveMinTime time.Duration
	// KeepAlivePermitWithoutStream allows clients to send keepalive pings even
	// when there is no active stream. Honoured by the same services as
	// KeepAliveMinTime.
	KeepAlivePermitWithoutStream bool
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
