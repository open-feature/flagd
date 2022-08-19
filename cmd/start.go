package cmd

import (
	"os"
	"strings"

	"github.com/open-feature/flagd/pkg/runtime"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	portFlagName            = "port"
	serviceProviderFlagName = "service-provider"
	socketPathFlagName      = "socket-path"
	syncProviderFlagName    = "sync-provider"
	evaluatorFlagName       = "evaluator"
	serverCertPathFlagName  = "server-cert-path"
	serverKeyPathFlagName   = "server-key-path"
	uriFlagName             = "uri"
	bearerTokenFlagName     = "bearer-token"
)

func init() {
	flags := startCmd.Flags()

	// allows environment variables to use _ instead of -
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_")) // sync-provider becomes SYNC_PROVIDER
	viper.SetEnvPrefix("FLAGD")                            // port becomes FLAGD_PORT

	flags.Int32P(portFlagName, "p", 8013, "Port to listen on")
	flags.StringP(socketPathFlagName, "d", "", "flagd socket path, only available when using the gRPC service provider")
	flags.StringP(serviceProviderFlagName, "s", "http", "Set a service provider e.g. http or grpc")
	flags.StringP(
		syncProviderFlagName, "y", "filepath", "Set a sync provider e.g. filepath or remote",
	)
	flags.StringP(evaluatorFlagName, "e", "json", "Set an evaluator e.g. json")
	flags.StringP(serverCertPathFlagName, "c", "", "Server side tls certificate path")
	flags.StringP(serverKeyPathFlagName, "k", "", "Server side tls key path")
	flags.StringSliceP(
		uriFlagName, "f", []string{}, "Set a sync provider uri to read data from this can be a filepath or url. "+
			"Using multiple providers is supported where collisions between "+
			"flags with the same key, the later will be used.")
	flags.StringP(
		bearerTokenFlagName, "b", "", "Set a bearer token to use for remote sync")

	_ = viper.BindPFlag(portFlagName, flags.Lookup(portFlagName))
	_ = viper.BindPFlag(socketPathFlagName, flags.Lookup(socketPathFlagName))
	_ = viper.BindPFlag(serviceProviderFlagName, flags.Lookup(serviceProviderFlagName))
	_ = viper.BindPFlag(syncProviderFlagName, flags.Lookup(syncProviderFlagName))
	_ = viper.BindPFlag(evaluatorFlagName, flags.Lookup(evaluatorFlagName))
	_ = viper.BindPFlag(serverCertPathFlagName, flags.Lookup(serverCertPathFlagName))
	_ = viper.BindPFlag(serverKeyPathFlagName, flags.Lookup(serverKeyPathFlagName))
	_ = viper.BindPFlag(uriFlagName, flags.Lookup(uriFlagName))
	_ = viper.BindPFlag(bearerTokenFlagName, flags.Lookup(bearerTokenFlagName))

	_ = startCmd.MarkFlagRequired("uri")
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start flagd",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Configure loggers -------------------------------------------------------
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		log.SetLevel(log.DebugLevel)

		// Build Runtime -----------------------------------------------------------
		rt, err := runtime.FromConfig(runtime.Config{
			ServiceProvider:   viper.GetString(serviceProviderFlagName),
			ServicePort:       viper.GetInt32(portFlagName),
			ServiceSocketPath: viper.GetString(socketPathFlagName),
			ServiceCertPath:   viper.GetString(serverCertPathFlagName),
			ServiceKeyPath:    viper.GetString(serverKeyPathFlagName),

			SyncProvider:    viper.GetString(syncProviderFlagName),
			SyncURI:         viper.GetStringSlice(uriFlagName),
			SyncBearerToken: viper.GetString(bearerTokenFlagName),

			Evaluator: viper.GetString(evaluatorFlagName),
		})
		if err != nil {
			log.Error(err)
		}

		rt.Start()
	},
}
