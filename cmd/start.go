package cmd

import (
	"log"
	"strings"

	"github.com/open-feature/flagd/pkg/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	portFlagName           = "port"
	metricsPortFlagName    = "metrics-port"
	socketPathFlagName     = "socket-path"
	syncProviderFlagName   = "sync-provider"
	providerArgsFlagName   = "sync-provider-args"
	evaluatorFlagName      = "evaluator"
	serverCertPathFlagName = "server-cert-path"
	serverKeyPathFlagName  = "server-key-path"
	uriFlagName            = "uri"
	bearerTokenFlagName    = "bearer-token"
	corsFlagName           = "cors-origin"
)

func init() {
	flags := startCmd.Flags()

	// allows environment variables to use _ instead of -
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_")) // sync-provider becomes SYNC_PROVIDER
	viper.SetEnvPrefix("FLAGD")                            // port becomes FLAGD_PORT
	flags.Int32P(metricsPortFlagName, "m", 8014, "Port to serve metrics on")
	flags.Int32P(portFlagName, "p", 8013, "Port to listen on")
	flags.StringP(socketPathFlagName, "d", "", "Flagd socket path. "+
		"With grpc the service will become available on this address. "+
		"With http(s) the grpc-gateway proxy will use this address internally.")
	flags.StringP(
		syncProviderFlagName, "y", "filepath", "Set a sync provider e.g. filepath or remote",
	)
	flags.StringP(evaluatorFlagName, "e", "json", "Set an evaluator e.g. json")
	flags.StringP(serverCertPathFlagName, "c", "", "Server side tls certificate path")
	flags.StringP(serverKeyPathFlagName, "k", "", "Server side tls key path")
	flags.StringToStringP(providerArgsFlagName,
		"a", nil, "Sync provider arguments as key values separated by =")
	flags.StringSliceP(
		uriFlagName, "f", []string{}, "Set a sync provider uri to read data from this can be a filepath or url. "+
			"Using multiple providers is supported where collisions between "+
			"flags with the same key, the later will be used.")
	flags.StringP(
		bearerTokenFlagName, "b", "", "Set a bearer token to use for remote sync")
	flags.StringSliceP(corsFlagName, "C", []string{}, "CORS allowed origins, * will allow all origins")

	_ = viper.BindPFlag(portFlagName, flags.Lookup(portFlagName))
	_ = viper.BindPFlag(metricsPortFlagName, flags.Lookup(metricsPortFlagName))
	_ = viper.BindPFlag(socketPathFlagName, flags.Lookup(socketPathFlagName))
	_ = viper.BindPFlag(syncProviderFlagName, flags.Lookup(syncProviderFlagName))
	_ = viper.BindPFlag(providerArgsFlagName, flags.Lookup(providerArgsFlagName))
	_ = viper.BindPFlag(evaluatorFlagName, flags.Lookup(evaluatorFlagName))
	_ = viper.BindPFlag(serverCertPathFlagName, flags.Lookup(serverCertPathFlagName))
	_ = viper.BindPFlag(serverKeyPathFlagName, flags.Lookup(serverKeyPathFlagName))
	_ = viper.BindPFlag(uriFlagName, flags.Lookup(uriFlagName))
	_ = viper.BindPFlag(bearerTokenFlagName, flags.Lookup(bearerTokenFlagName))
	_ = viper.BindPFlag(corsFlagName, flags.Lookup(corsFlagName))
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start flagd",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Configure loggers -------------------------------------------------------
		// log.SetFormatter(&log.JSONFormatter{})
		// log.SetOutput(os.Stdout)
		// if Debug {
		// 	log.SetLevel(log.DebugLevel)
		// 	log.SetReportCaller(true)
		// } else {
		// 	log.SetLevel(log.InfoLevel)
		// }
		var logger *zap.Logger
		var err error
		if Debug {
			logger, err = zap.NewDevelopment()
		} else {
			logger, err = zap.NewProduction()
		}
		if err != nil {
			log.Fatalf("can't initialize zap logger: %v", err)
		}
		defer logger.Sync()

		rtLogger := logger.With(zap.String("component", "start"))
		// Build Runtime -----------------------------------------------------------
		rt, err := runtime.FromConfig(logger, runtime.Config{
			ServicePort:       viper.GetInt32(portFlagName),
			MetricsPort:       viper.GetInt32(metricsPortFlagName),
			ServiceSocketPath: viper.GetString(socketPathFlagName),
			ServiceCertPath:   viper.GetString(serverCertPathFlagName),
			ServiceKeyPath:    viper.GetString(serverKeyPathFlagName),

			SyncProvider:    viper.GetString(syncProviderFlagName),
			ProviderArgs:    viper.GetStringMapString(providerArgsFlagName),
			SyncURI:         viper.GetStringSlice(uriFlagName),
			SyncBearerToken: viper.GetString(bearerTokenFlagName),

			Evaluator: viper.GetString(evaluatorFlagName),

			CORS: viper.GetStringSlice(corsFlagName),
		})
		if err != nil {
			rtLogger.Fatal(err.Error())
		}

		if err := rt.Start(); err != nil {
			rtLogger.Fatal(err.Error())
		}
	},
}
