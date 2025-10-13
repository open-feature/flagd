package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/sync"
	syncbuilder "github.com/open-feature/flagd/core/pkg/sync/builder"
	"github.com/open-feature/flagd/flagd/pkg/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	corsFlagName               = "cors-origin"
	logFormatFlagName          = "log-format"
	managementPortFlagName     = "management-port"
	metricsExporter            = "metrics-exporter"
	ofrepPortFlagName          = "ofrep-port"
	otelCollectorURI           = "otel-collector-uri"
	otelCertPathFlagName       = "otel-cert-path"
	otelKeyPathFlagName        = "otel-key-path"
	otelCAPathFlagName         = "otel-ca-path"
	otelReloadIntervalFlagName = "otel-reload-interval"
	portFlagName               = "port"
	serverCertPathFlagName     = "server-cert-path"
	serverKeyPathFlagName      = "server-key-path"
	socketPathFlagName         = "socket-path"
	sourcesFlagName            = "sources"
	syncPortFlagName           = "sync-port"
	syncSocketPathFlagName     = "sync-socket-path"
	uriFlagName                = "uri"
	disableSyncMetadata        = "disable-sync-metadata"
	contextValueFlagName       = "context-value"
	headerToContextKeyFlagName = "context-from-header"
	streamDeadlineFlagName     = "stream-deadline"
)

func init() {
	flags := startCmd.Flags()
	// allows environment variables to use _ instead of -
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_")) // sync-provider-args becomes SYNC_PROVIDER_ARGS
	viper.SetEnvPrefix("FLAGD")                            // port becomes FLAGD_PORT

	flags.Int32P(managementPortFlagName, "m", 8014, "Port for management operations")
	flags.Int32P(portFlagName, "p", 8013, "Port to listen on")
	flags.Int32P(syncPortFlagName, "g", 8015, "gRPC Sync port")
	flags.Int32P(ofrepPortFlagName, "r", 8016, "ofrep service port")

	flags.StringP(socketPathFlagName, "d", "", "Flagd unix socket path. "+
		"With grpc the evaluations service will become available on this address. "+
		"With http(s) the grpc-gateway proxy will use this address internally.")
	flags.StringP(syncSocketPathFlagName, "e", "", "Flagd sync service socket path. "+
		"With grpc the sync service will be available on this address.")
	flags.StringP(serverCertPathFlagName, "c", "", "Server side tls certificate path")
	flags.StringP(serverKeyPathFlagName, "k", "", "Server side tls key path")
	flags.StringSliceP(
		uriFlagName, "f", []string{}, "Set a sync provider uri to read data from, this can be a filepath,"+
			" URL (HTTP and gRPC), FeatureFlag custom resource, or GCS or Azure Blob. "+
			"When flag keys are duplicated across multiple providers the "+
			"merge priority follows the index of the flag arguments, as such flags from the uri at index 0 take the "+
			"lowest precedence, with duplicated keys being overwritten by those from the uri at index 1. "+
			"Please note that if you are using filepath, flagd only supports files with `.yaml/.yml/.json` extension.",
	)
	flags.StringSliceP(corsFlagName, "C", []string{}, "CORS allowed origins, * will allow all origins")
	flags.StringP(
		sourcesFlagName, "s", "", "JSON representation of an array of SourceConfig objects. This object contains "+
			"2 required fields, uri (string) and provider (string). Documentation for this object: "+
			"https://flagd.dev/reference/sync-configuration/#source-configuration",
	)
	flags.StringP(logFormatFlagName, "z", "console", "Set the logging format, e.g. console or json")
	flags.StringP(metricsExporter, "t", "", "Set the metrics exporter. Default(if unset) is Prometheus."+
		" Can be override to otel - OpenTelemetry metric exporter. Overriding to otel require otelCollectorURI to"+
		" be present")
	flags.StringP(otelCollectorURI, "o", "", "Set the grpc URI of the OpenTelemetry collector "+
		"for flagd runtime. If unset, the collector setup will be ignored and traces will not be exported.")
	flags.StringP(otelCertPathFlagName, "D", "", "tls certificate path to use with OpenTelemetry collector")
	flags.StringP(otelKeyPathFlagName, "K", "", "tls key path to use with OpenTelemetry collector")
	flags.StringP(otelCAPathFlagName, "A", "", "tls certificate authority path to use with OpenTelemetry collector")
	flags.DurationP(otelReloadIntervalFlagName, "I", time.Hour, "how long between reloading the otel tls certificate "+
		"from disk")
	flags.StringToStringP(contextValueFlagName, "X", map[string]string{}, "add arbitrary key value pairs "+
		"to the flag evaluation context")
	flags.StringToStringP(headerToContextKeyFlagName, "H", map[string]string{}, "add key-value pairs to map "+
		"header values to context values, where key is Header name, value is context key")
	flags.Duration(streamDeadlineFlagName, 0, "Set a server-side deadline for flagd sync and event streams (default 0, means no deadline).")
	flags.Bool(disableSyncMetadata, false, "Disables the getMetadata endpoint of the sync service. Defaults to false, but will default to true in later versions.")

	bindFlags(flags)
}

func bindFlags(flags *pflag.FlagSet) {
	_ = viper.BindPFlag(corsFlagName, flags.Lookup(corsFlagName))
	_ = viper.BindPFlag(logFormatFlagName, flags.Lookup(logFormatFlagName))
	_ = viper.BindPFlag(metricsExporter, flags.Lookup(metricsExporter))
	_ = viper.BindPFlag(managementPortFlagName, flags.Lookup(managementPortFlagName))
	_ = viper.BindPFlag(otelCollectorURI, flags.Lookup(otelCollectorURI))
	_ = viper.BindPFlag(otelCertPathFlagName, flags.Lookup(otelCertPathFlagName))
	_ = viper.BindPFlag(otelKeyPathFlagName, flags.Lookup(otelKeyPathFlagName))
	_ = viper.BindPFlag(otelCAPathFlagName, flags.Lookup(otelCAPathFlagName))
	_ = viper.BindPFlag(portFlagName, flags.Lookup(portFlagName))
	_ = viper.BindPFlag(serverCertPathFlagName, flags.Lookup(serverCertPathFlagName))
	_ = viper.BindPFlag(serverKeyPathFlagName, flags.Lookup(serverKeyPathFlagName))
	_ = viper.BindPFlag(socketPathFlagName, flags.Lookup(socketPathFlagName))
	_ = viper.BindPFlag(sourcesFlagName, flags.Lookup(sourcesFlagName))
	_ = viper.BindPFlag(uriFlagName, flags.Lookup(uriFlagName))
	_ = viper.BindPFlag(syncPortFlagName, flags.Lookup(syncPortFlagName))
	_ = viper.BindPFlag(syncSocketPathFlagName, flags.Lookup(syncSocketPathFlagName))
	_ = viper.BindPFlag(ofrepPortFlagName, flags.Lookup(ofrepPortFlagName))
	_ = viper.BindPFlag(contextValueFlagName, flags.Lookup(contextValueFlagName))
	_ = viper.BindPFlag(headerToContextKeyFlagName, flags.Lookup(headerToContextKeyFlagName))
	_ = viper.BindPFlag(streamDeadlineFlagName, flags.Lookup(streamDeadlineFlagName))
	_ = viper.BindPFlag(disableSyncMetadata, flags.Lookup(disableSyncMetadata))
}

func overrideMetricsExporter() string {
	var metricsExporter = viper.GetString(metricsExporter)
	if metricsExporter != "" {
		return metricsExporter
	}

	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" {
		metricsExporter = "otel-sdk"
	}

	return metricsExporter
}

func overrideOtelUri() string {
	var collectorUri = viper.GetString(otelCollectorURI)
	if collectorUri != "" {
		return collectorUri
	}

	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" {
		collectorUri = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	}

	return collectorUri
}

func getOtelHeaders() string {
	return os.Getenv("OTEL_EXPORTER_OTLP_HEADERS")
}

func getOtelProtocol() string {
	return os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL")
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start flagd",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Configure loggers -------------------------------------------------------
		var level zapcore.Level
		var err error
		if Debug {
			level = zapcore.DebugLevel
		} else {
			level = zapcore.InfoLevel
		}
		l, err := logger.NewZapLogger(level, viper.GetString(logFormatFlagName))
		if err != nil {
			log.Fatalf("can't initialize zap logger: %v", err)
		}
		logger := logger.NewLogger(l, Debug)
		rtLogger := logger.WithFields(zap.String("component", "start"))

		rtLogger.Info(fmt.Sprintf("flagd version: %s (%s), built at: %s", Version, Commit, Date))

		syncProviders, err := syncbuilder.ParseSyncProviderURIs(viper.GetStringSlice(uriFlagName))
		if err != nil {
			log.Fatal(err)
		}

		syncProvidersFromConfig := []sync.SourceConfig{}
		if cfgFile == "" && viper.GetString(sourcesFlagName) != "" {
			syncProvidersFromConfig, err = syncbuilder.ParseSources(viper.GetString(sourcesFlagName))
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err = viper.UnmarshalKey(sourcesFlagName, &syncProvidersFromConfig)
			if err != nil {
				log.Fatal(err)
			}
		}
		syncProviders = append(syncProviders, syncProvidersFromConfig...)

		contextValuesToMap := make(map[string]any)
		for k, v := range viper.GetStringMapString(contextValueFlagName) {
			contextValuesToMap[k] = v
		}

		headerToContextKeyMappings := make(map[string]string)
		for k, v := range viper.GetStringMapString(headerToContextKeyFlagName) {
			headerToContextKeyMappings[k] = v
		}

		var metricsExporter = overrideMetricsExporter()
		var collectorUri = overrideOtelUri()
		var otelHeaders = getOtelHeaders()
		var otelProtocol = getOtelProtocol()

		// Build Runtime -----------------------------------------------------------
		rt, err := runtime.FromConfig(logger, Version, runtime.Config{
			CORS:                       viper.GetStringSlice(corsFlagName),
			MetricExporter:             metricsExporter,
			ManagementPort:             viper.GetUint16(managementPortFlagName),
			OfrepServicePort:           viper.GetUint16(ofrepPortFlagName),
			OtelCollectorURI:           collectorUri,
			OtelCertPath:               viper.GetString(otelCertPathFlagName),
			OtelKeyPath:                viper.GetString(otelKeyPathFlagName),
			OtelReloadInterval:         viper.GetDuration(otelReloadIntervalFlagName),
			OtelHeaders:                otelHeaders,
			OtelProtocol:               otelProtocol,
			OtelCAPath:                 viper.GetString(otelCAPathFlagName),
			ServiceCertPath:            viper.GetString(serverCertPathFlagName),
			ServiceKeyPath:             viper.GetString(serverKeyPathFlagName),
			ServicePort:                viper.GetUint16(portFlagName),
			ServiceSocketPath:          viper.GetString(socketPathFlagName),
			SyncServicePort:            viper.GetUint16(syncPortFlagName),
			SyncServiceSocketPath:      viper.GetString(syncSocketPathFlagName),
			StreamDeadline:             viper.GetDuration(streamDeadlineFlagName),
			DisableSyncMetadata:        viper.GetBool(disableSyncMetadata),
			SyncProviders:              syncProviders,
			ContextValues:              contextValuesToMap,
			HeaderToContextKeyMappings: headerToContextKeyMappings,
		})
		if err != nil {
			rtLogger.Fatal(err.Error())
		}

		if err := rt.Start(); err != nil {
			rtLogger.Fatal(err.Error())
		}
	},
}
