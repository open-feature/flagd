package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/service"
	syncServer "github.com/open-feature/flagd/core/pkg/service/sync"
	sync_store "github.com/open-feature/flagd/core/pkg/sync-store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

// start

const (
	logFormatFlagName   = "log-format"
	metricsPortFlagName = "metrics-port"
	portFlagName        = "port"
)

func init() {
	flags := startCmd.Flags()

	// allows environment variables to use _ instead of -
	flags.Int32P(portFlagName, "p", 8013, "Port to listen on")
	flags.Int32P(metricsPortFlagName, "m", 8014, "Metrics port to listen on")
	flags.StringP(logFormatFlagName, "z", "console", "Set the logging format, e.g. console or json ")

	_ = viper.BindPFlag(logFormatFlagName, flags.Lookup(logFormatFlagName))
	_ = viper.BindPFlag(metricsPortFlagName, flags.Lookup(metricsPortFlagName))
	_ = viper.BindPFlag(portFlagName, flags.Lookup(portFlagName))
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

		ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

		store := sync_store.NewSyncStore(ctx, logger)

		s := syncServer.SyncServer{
			SyncStore: store,
			Configuration: syncServer.SyncServerConfiguration{
				Port:        viper.GetUint16(portFlagName),
				MetricsPort: viper.GetUint16(metricsPortFlagName),
			},
			Logger: logger,
		}

		go store.Cleanup()
		go s.Serve(ctx, service.Configuration{
			ReadinessProbe: func() bool { return true },
		})

		logger.Info(fmt.Sprintf("listening for connections on %d...", s.Configuration.Port))

		<-ctx.Done()
	},
}
