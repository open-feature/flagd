package runtime

import (
	"context"
	"fmt"
	"time"

	"github.com/open-feature/flagd/core/pkg/evaluator"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/service"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/sync"
	syncbuilder "github.com/open-feature/flagd/core/pkg/sync/builder"
	"github.com/open-feature/flagd/core/pkg/telemetry"
	flageval "github.com/open-feature/flagd/flagd/pkg/service/flag-evaluation"
	"github.com/open-feature/flagd/flagd/pkg/service/flag-evaluation/ofrep"
	flagsync "github.com/open-feature/flagd/flagd/pkg/service/flag-sync"
	"go.uber.org/zap"
)

// from_config is a collection of structures and parsers responsible for deriving flagd runtime

const svcName = "flagd"

// Config is the configuration structure derived from startup arguments.
type Config struct {
	MetricExporter        string
	ManagementPort        uint16
	OfrepServicePort      uint16
	OtelCollectorURI      string
	OtelCertPath          string
	OtelKeyPath           string
	OtelCAPath            string
	OtelReloadInterval    time.Duration
	ServiceCertPath       string
	ServiceKeyPath        string
	ServicePort           uint16
	ServiceSocketPath     string
	SyncServicePort       uint16
	SyncServiceSocketPath string
	StreamDeadline        time.Duration
	DisableSyncMetadata   bool

	SyncProviders []sync.SourceConfig
	CORS          []string

	SelectorFallbackKey string

	ContextValues              map[string]any
	HeaderToContextKeyMappings map[string]string
}

// FromConfig builds a runtime from startup configurations
// nolint: funlen
func FromConfig(logger *logger.Logger, version string, config Config) (*Runtime, error) {
	telCfg := telemetry.Config{
		MetricsExporter: config.MetricExporter,
		CollectorConfig: telemetry.CollectorConfig{
			Target:         config.OtelCollectorURI,
			CertPath:       config.OtelCertPath,
			KeyPath:        config.OtelKeyPath,
			CAPath:         config.OtelCAPath,
			ReloadInterval: config.OtelReloadInterval,
		},
	}

	// register error handling for OpenTelemetry
	telemetry.RegisterErrorHandling(logger)

	// register trace provider for the runtime
	err := telemetry.BuildTraceProvider(context.Background(), logger, svcName, version, telCfg)
	if err != nil {
		// log the error but continue
		logger.Error(fmt.Sprintf("error building trace provider: %v", err))
	}

	// build metrics recorder with startup configurations
	recorder, err := telemetry.BuildMetricsRecorder(context.Background(), svcName, version, telCfg)
	if err != nil {
		// log the error but continue
		logger.Error(fmt.Sprintf("error building metrics recorder: %v", err))
	}

	sources := []string{}

	for _, provider := range config.SyncProviders {
		sources = append(sources, provider.URI)
	}

	// build flag store, collect flag sources & fill sources details
	store, err := store.NewStore(logger, store.StoreConfig{Sources: sources, SelectorFallbackKey: config.SelectorFallbackKey})
	if err != nil {
		return nil, fmt.Errorf("error creating flag store: %w", err)
	}

	// derive evaluator
	jsonEvaluator := evaluator.NewJSON(logger, store)

	// derive services

	// connect service
	connectService := flageval.NewConnectService(
		logger.WithFields(zap.String("component", "service")),
		jsonEvaluator,
		store,
		recorder)

	// ofrep service
	ofrepService, err := ofrep.NewOfrepService(jsonEvaluator, config.CORS, ofrep.SvcConfiguration{
		Logger:          logger.WithFields(zap.String("component", "OFREPService")),
		Port:            config.OfrepServicePort,
		ServiceName:     svcName,
		MetricsRecorder: recorder,
	},
		config.ContextValues,
		config.HeaderToContextKeyMappings,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating OFREP service: %w", err)
	}

	// flag sync service
	flagSyncService, err := flagsync.NewSyncService(flagsync.SvcConfigurations{
		Logger:              logger.WithFields(zap.String("component", "FlagSyncService")),
		Port:                config.SyncServicePort,
		Sources:             sources,
		Store:               store,
		ContextValues:       config.ContextValues,
		KeyPath:             config.ServiceKeyPath,
		CertPath:            config.ServiceCertPath,
		SocketPath:          config.SyncServiceSocketPath,
		StreamDeadline:      config.StreamDeadline,
		DisableSyncMetadata: config.DisableSyncMetadata,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating sync service: %w", err)
	}

	// build sync providers
	syncLogger := logger.WithFields(zap.String("component", "sync"))
	iSyncs, err := syncProvidersFromConfig(syncLogger, config.SyncProviders)
	if err != nil {
		return nil, err
	}

	options, err := telemetry.BuildConnectOptions(telCfg)
	if err != nil {
		// log the error but continue
		logger.Error(fmt.Sprintf("failed to build connect options, %v", err))
	}

	return &Runtime{
		Logger:            logger.WithFields(zap.String("component", "runtime")),
		Evaluator:         jsonEvaluator,
		SyncService:       flagSyncService,
		OfrepService:      ofrepService,
		EvaluationService: connectService,
		ServiceConfig: service.Configuration{
			Port:                       config.ServicePort,
			ManagementPort:             config.ManagementPort,
			ServiceName:                svcName,
			KeyPath:                    config.ServiceKeyPath,
			CertPath:                   config.ServiceCertPath,
			SocketPath:                 config.ServiceSocketPath,
			CORS:                       config.CORS,
			Options:                    options,
			ContextValues:              config.ContextValues,
			HeaderToContextKeyMappings: config.HeaderToContextKeyMappings,
			StreamDeadline:             config.StreamDeadline,
			SelectorFallbackKey:        config.SelectorFallbackKey,
		},
		Syncs: iSyncs,
	}, nil
}

// syncProvidersFromConfig is a helper to build ISync implementations from SourceConfig
func syncProvidersFromConfig(logger *logger.Logger, sources []sync.SourceConfig) ([]sync.ISync, error) {
	builder := syncbuilder.NewSyncBuilder()
	syncs, err := builder.SyncsFromConfig(sources, logger)
	if err != nil {
		return nil, fmt.Errorf("could not create sync sources from config: %w", err)
	}

	return syncs, nil
}
