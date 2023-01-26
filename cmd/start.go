package cmd

import (
	"log"
	"strings"

	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	portFlagName           = "port"
	metricsPortFlagName    = "metrics-port"
	socketPathFlagName     = "socket-path"
	providerArgsFlagName   = "sync-provider-args"
	evaluatorFlagName      = "evaluator"
	serverCertPathFlagName = "server-cert-path"
	serverKeyPathFlagName  = "server-key-path"
	uriFlagName            = "uri"
	bearerTokenFlagName    = "bearer-token"
	corsFlagName           = "cors-origin"
	syncProviderFlagName   = "sync-provider"
	prettyLogFlagName      = "log-format"
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
			"url or FeatureFlagConfiguration. Using multiple providers is supported however if"+
			" flag keys are duplicated across multiple sources it may lead to unexpected behavior. "+
			"Please note that if you are using filepath, flagd only supports files with `.yaml/.yml/.json` extension.",
	)
	flags.StringP(
		bearerTokenFlagName, "b", "", "Set a bearer token to use for remote sync")
	flags.StringSliceP(corsFlagName, "C", []string{}, "CORS allowed origins, * will allow all origins")
	flags.StringP(
		syncProviderFlagName, "y", "", "DEPRECATED: Set a sync provider e.g. filepath or remote",
	)
	flags.StringP(prettyLogFlagName, "z", "console", "Set the logging format, e.g. console or json ")

	_ = viper.BindPFlag(portFlagName, flags.Lookup(portFlagName))
	_ = viper.BindPFlag(metricsPortFlagName, flags.Lookup(metricsPortFlagName))
	_ = viper.BindPFlag(socketPathFlagName, flags.Lookup(socketPathFlagName))
	_ = viper.BindPFlag(providerArgsFlagName, flags.Lookup(providerArgsFlagName))
	_ = viper.BindPFlag(evaluatorFlagName, flags.Lookup(evaluatorFlagName))
	_ = viper.BindPFlag(serverCertPathFlagName, flags.Lookup(serverCertPathFlagName))
	_ = viper.BindPFlag(serverKeyPathFlagName, flags.Lookup(serverKeyPathFlagName))
	_ = viper.BindPFlag(uriFlagName, flags.Lookup(uriFlagName))
	_ = viper.BindPFlag(bearerTokenFlagName, flags.Lookup(bearerTokenFlagName))
	_ = viper.BindPFlag(corsFlagName, flags.Lookup(corsFlagName))
	_ = viper.BindPFlag(syncProviderFlagName, flags.Lookup(syncProviderFlagName))
	_ = viper.BindPFlag(prettyLogFlagName, flags.Lookup(prettyLogFlagName))
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
		l, err := logger.NewZapLogger(level, viper.GetString(prettyLogFlagName))
		if err != nil {
			log.Fatalf("can't initialize zap logger: %v", err)
		}
		logger := logger.NewLogger(l, Debug)
		rtLogger := logger.WithFields(zap.String("component", "start"))

		if viper.GetString(syncProviderFlagName) != "" {
			rtLogger.Warn("DEPRECATED: the --sync-provider flag has been deprecated " +
				"Docs: https://github.com/open-feature/flagd/blob/main/docs/configuration.md")
		}

		if viper.GetString(evaluatorFlagName) != "" {
			rtLogger.Warn("DEPRECATED: the --evaluator flag has been deprecated " +
				"Docs: https://github.com/open-feature/flagd/blob/main/docs/configuration.md")
		}
		// Build Runtime -----------------------------------------------------------
		rt, err := runtime.FromConfig(logger, runtime.Config{
			ServicePort:       viper.GetInt32(portFlagName),
			MetricsPort:       viper.GetInt32(metricsPortFlagName),
			ServiceSocketPath: viper.GetString(socketPathFlagName),
			ServiceCertPath:   viper.GetString(serverCertPathFlagName),
			ServiceKeyPath:    viper.GetString(serverKeyPathFlagName),
			ProviderArgs:      viper.GetStringMapString(providerArgsFlagName),
			SyncURI:           viper.GetStringSlice(uriFlagName),
			SyncBearerToken:   viper.GetString(bearerTokenFlagName),
			CORS:              viper.GetStringSlice(corsFlagName),
		})
		if err != nil {
			rtLogger.Fatal(err.Error())
		}

		if err := rt.Start(); err != nil {
			rtLogger.Fatal(err.Error())
		}
	},
}
