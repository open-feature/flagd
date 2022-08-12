package cmd

import (
	"os"

	"github.com/open-feature/flagd/pkg/runtime"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	serviceProvider   string
	syncProvider      string
	evaluator         string
	uri               []string
	servicePort       int32
	socketServicePath string
	bearerToken       string
	serverCertPath    string
	serverKeyPath     string
)

func init() {
	startCmd.Flags().Int32VarP(
		&servicePort, "port", "p", 8013, "Port to listen on")
	startCmd.Flags().StringVarP(
		&socketServicePath, "socketpath", "d", "/tmp/flagd.sock", "flagd socket path")
	startCmd.Flags().StringVarP(
		&serviceProvider, "service-provider", "s", "http", "Set a serve provider e.g. http or grpc")
	startCmd.Flags().StringVarP(
		&syncProvider, "sync-provider", "y", "filepath", "Set a sync provider e.g. filepath or remote")
	startCmd.Flags().StringVarP(
		&evaluator, "evaluator", "e", "json", "Set an evaluator e.g. json")
	startCmd.Flags().StringVarP(
		&serverCertPath, "server-cert-path", "c", "", "Server side tls certificate path")
	startCmd.Flags().StringVarP(
		&serverKeyPath, "server-key-path", "k", "", "Server side tls key path")
	startCmd.Flags().StringSliceVarP(
		&uri, "uri", "f", []string{}, "Set a sync provider uri to read data from this can be a filepath or url. "+
			"Using multiple providers is supported where collisions between "+
			"flags with the same key, the later will be used.")
	startCmd.Flags().StringVarP(
		&bearerToken, "bearer-token", "b", "", "Set a bearer token to use for remote sync")

	_ = startCmd.MarkFlagRequired("uri")
	rootCmd.AddCommand(startCmd)
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
		rt, err := runtime.RuntimeFromConfig(runtime.RuntimeConfig{
			ServiceProvider:   serviceProvider,
			ServicePort:       servicePort,
			ServiceSocketPath: socketServicePath,
			ServiceCertPath:   serverCertPath,
			ServiceKeyPath:    serverKeyPath,

			SyncProvider:    syncProvider,
			SyncUri:         uri,
			SyncBearerToken: bearerToken,

			Evaluator: evaluator,
		})

		if err != nil {
			log.Error(err)
		}

		rt.Start()
	},
}
