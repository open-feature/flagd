package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/runtime"
	"github.com/open-feature/flagd/pkg/sync"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	bearerTokenFlagName    = "bearer-token"
	corsFlagName           = "cors-origin"
	evaluatorFlagName      = "evaluator"
	logFormatFlagName      = "log-format"
	metricsPortFlagName    = "metrics-port"
	portFlagName           = "port"
	providerArgsFlagName   = "sync-provider-args"
	serverCertPathFlagName = "server-cert-path"
	serverKeyPathFlagName  = "server-key-path"
	socketPathFlagName     = "socket-path"
	syncProvidersFlagName  = "sync-providers"
	syncProviderFlagName   = "sync-provider"
	uriFlagName            = "uri"
)

func init() {
	flags := startCmd.Flags()

	// allows environment variables to use _ instead of -
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_")) // sync-provider-args becomes SYNC_PROVIDER_ARGS
	viper.SetEnvPrefix("FLAGD")                            // port becomes FLAGD_PORT
	flags.Int32P(metricsPortFlagName, "m", 8014, "Port to serve metrics on")
	flags.Int32P(portFlagName, "p", 8013, "Port to listen on")
	flags.StringP(socketPathFlagName, "d", "", "Flagd socket path. "+
		"With grpc the service will become available on this address. "+
		"With http(s) the grpc-gateway proxy will use this address internally.")
	flags.StringP(evaluatorFlagName, "e", "json", "DEPRECATED: Set an evaluator e.g. json, yaml/yml."+
		"Please note that yaml/yml and json evaluations work the same (yaml/yml files are converted to json internally)")
	flags.StringP(serverCertPathFlagName, "c", "", "Server side tls certificate path")
	flags.StringP(serverKeyPathFlagName, "k", "", "Server side tls key path")
	flags.StringToStringP(providerArgsFlagName,
		"a", nil, "Sync provider arguments as key values separated by =")
	flags.StringSliceP(
		uriFlagName, "f", []string{}, "Set a sync provider uri to read data from, this can be a filepath,"+
			"url (http and grpc) or FeatureFlagConfiguration. Using multiple providers is supported however if"+
			" flag keys are duplicated across multiple sources it may lead to unexpected behavior. "+
			"Please note that if you are using filepath, flagd only supports files with `.yaml/.yml/.json` extension.",
	)
	flags.StringP(
		bearerTokenFlagName, "b", "", "Set a bearer token to use for remote sync")
	flags.StringSliceP(corsFlagName, "C", []string{}, "CORS allowed origins, * will allow all origins")
	flags.StringP(
		syncProviderFlagName, "y", "", "DEPRECATED: Set a sync provider e.g. filepath or remote",
	)
	flags.StringP(
		syncProvidersFlagName, "y", "", "JSON representation of an array of SyncProviderConfig objects. This object contains "+
			"2 required fields, uri (string) and provider (string). Documentation for this object can be found here: ",
	)
	flags.StringP(logFormatFlagName, "z", "console", "Set the logging format, e.g. console or json ")

	_ = viper.BindPFlag(bearerTokenFlagName, flags.Lookup(bearerTokenFlagName))
	_ = viper.BindPFlag(corsFlagName, flags.Lookup(corsFlagName))
	_ = viper.BindPFlag(evaluatorFlagName, flags.Lookup(evaluatorFlagName))
	_ = viper.BindPFlag(logFormatFlagName, flags.Lookup(logFormatFlagName))
	_ = viper.BindPFlag(metricsPortFlagName, flags.Lookup(metricsPortFlagName))
	_ = viper.BindPFlag(portFlagName, flags.Lookup(portFlagName))
	_ = viper.BindPFlag(providerArgsFlagName, flags.Lookup(providerArgsFlagName))
	_ = viper.BindPFlag(serverCertPathFlagName, flags.Lookup(serverCertPathFlagName))
	_ = viper.BindPFlag(serverKeyPathFlagName, flags.Lookup(serverKeyPathFlagName))
	_ = viper.BindPFlag(socketPathFlagName, flags.Lookup(socketPathFlagName))
	_ = viper.BindPFlag(syncProviderFlagName, flags.Lookup(syncProviderFlagName))
	_ = viper.BindPFlag(syncProvidersFlagName, flags.Lookup(syncProviderFlagName))
	_ = viper.BindPFlag(uriFlagName, flags.Lookup(uriFlagName))
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

		if viper.GetString(evaluatorFlagName) != "json" {
			rtLogger.Warn("DEPRECATED: The --evaluator flag has been deprecated. " +
				"Docs: https://github.com/open-feature/flagd/blob/main/docs/configuration/configuration.md")
		}

		syncProviders, err := runtime.SyncProvidersFromURIs(viper.GetStringSlice(uriFlagName))
		if err != nil {
			log.Fatal(err)
		}
		syncProviders2 := []sync.SyncProviderConfig{}
		if cfgFile == "" {
			syncProviders2, err = runtime.SyncProviderArgPass(viper.GetString(syncProvidersFlagName))
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err = viper.UnmarshalKey(syncProvidersFlagName, &syncProviders2)
			if err != nil {
				log.Fatal(err)
			}
		}
		syncProviders = append(syncProviders, syncProviders2...)

		// Build Runtime -----------------------------------------------------------
		rt, err := runtime.FromConfig(logger, runtime.Config{
			CORS:              viper.GetStringSlice(corsFlagName),
			MetricsPort:       viper.GetInt32(metricsPortFlagName),
			ServiceCertPath:   viper.GetString(serverCertPathFlagName),
			ServiceKeyPath:    viper.GetString(serverKeyPathFlagName),
			ServicePort:       viper.GetInt32(portFlagName),
			ServiceSocketPath: viper.GetString(socketPathFlagName),
			SyncProviders:     syncProviders,
		})
		if err != nil {
			rtLogger.Fatal(err.Error())
		}

		if err := rt.Start(); err != nil {
			rtLogger.Fatal(err.Error())
		}
	},
}
