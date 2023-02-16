package cmd

import (
	"fmt"
	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"log"
)

const (
	address  = "address"
	secure   = "secure"
	certPath = "cert-path"
	keyPath  = "key-path"
	sources  = "source"
)

// NewServerCmd is the command to start flagd in server mode
func NewServerCmd() *cobra.Command {
	flagdCmd := &cobra.Command{
		Use:   "server",
		Short: "Start flagd as a server",
		Run:   runServer,
	}

	setupServer(flagdCmd.Flags())
	return flagdCmd
}

// setupServer setup flags of the command
func setupServer(flags *pflag.FlagSet) {
	flags.StringP(address, "p", "localhost:9090", "Path this server binds to")
	flags.BoolP(secure, "s", false, "Start secure server")
	flags.StringP(certPath, "c", "", "TLS certificate path")
	flags.StringP(keyPath, "k", "", "TLS key path of the certificate")
	flags.StringP(sources, "f", "", "CRD with feature flag configurations")

	_ = viper.BindPFlag(address, flags.Lookup(address))
	_ = viper.BindPFlag(secure, flags.Lookup(secure))
	_ = viper.BindPFlag(certPath, flags.Lookup(certPath))
	_ = viper.BindPFlag(keyPath, flags.Lookup(keyPath))
	_ = viper.BindPFlag(sources, flags.Lookup(sources))
}

func runServer(cmd *cobra.Command, args []string) {
	zapLogger, err := logger.NewZapLogger(zapcore.DebugLevel, "console")
	if err != nil {
		log.Fatalf("error setting up the logger: %s", err)
	}

	logWrapper := logger.NewLogger(zapLogger, true)

	err = viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		logWrapper.Fatal(fmt.Sprintf("error parsing flags: %s", err.Error()))
	}

	serverConfig := runtime.ServerConfig{
		Address:     viper.GetString(address),
		Secure:      viper.GetBool(secure),
		CertPath:    viper.GetString(certPath),
		KeyPath:     viper.GetString(keyPath),
		SyncSources: viper.GetStringSlice(sources),
	}

	serverRuntime, err := runtime.NewServerRuntime(serverConfig, logWrapper)
	if err != nil {
		logWrapper.Fatal(fmt.Sprintf("error creating the server runtime: %s", err.Error()))
	}

	err = serverRuntime.Start()
	if err != nil {
		logWrapper.Fatal(fmt.Sprintf("error from server runtime: %s", err.Error()))
	}
}
