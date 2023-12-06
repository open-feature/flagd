package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/service"
	syncServer "github.com/open-feature/flagd/core/pkg/service/sync"
	"github.com/open-feature/flagd/core/pkg/subscriptions"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

// start

const (
	logFormatFlagName      = "log-format"
	metricsPortFlagName    = "metrics-port" // deprecated
	managementPortFlagName = "management-port"
	portFlagName           = "port"
	defaultManagementPort  = 8016
)

func init() {
	flags := startCmd.Flags()

	// allows environment variables to use _ instead of -
	flags.Int32P(portFlagName, "p", 8015, "Port to listen on")
	flags.Int32(metricsPortFlagName, defaultManagementPort, "DEPRECATED: Superseded by --management-port.")
	flags.Int32P(managementPortFlagName, "m", defaultManagementPort, "Management port")
	flags.StringP(logFormatFlagName, "z", "console", "Set the logging format, e.g. console or json")

	_ = viper.BindPFlag(logFormatFlagName, flags.Lookup(logFormatFlagName))
	_ = viper.BindPFlag(metricsPortFlagName, flags.Lookup(metricsPortFlagName))
	_ = viper.BindPFlag(managementPortFlagName, flags.Lookup(managementPortFlagName))
	_ = viper.BindPFlag(portFlagName, flags.Lookup(portFlagName))
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start flagd-proxy",
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

		if viper.GetUint16(metricsPortFlagName) != defaultManagementPort {
			logger.Warn("DEPRECATED: The --metrics-port flag has been deprecated and is superseded by --management-port.")
		}

		ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

		syncStore := subscriptions.NewManager(ctx, logger)
		s := syncServer.NewServer(logger, syncStore)

		// If --management-port is set use that value. If not and
		// --metrics-port is set use that value. Otherwise use the default
		// value.
		managementPort := uint16(defaultManagementPort)
		if viper.GetUint16(managementPortFlagName) != defaultManagementPort {
			managementPort = viper.GetUint16(managementPortFlagName)
		} else if viper.GetUint16(metricsPortFlagName) != defaultManagementPort {
			managementPort = viper.GetUint16(metricsPortFlagName)
		}

		cfg := service.Configuration{
			ReadinessProbe: func() bool { return true },
			Port:           viper.GetUint16(portFlagName),
			ManagementPort: managementPort,
		}

		errChan := make(chan error, 1)
		go func() {
			if err := s.Serve(ctx, cfg); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errChan <- err
			}
		}()

		logger.Info(fmt.Sprintf("listening for connections on %d", cfg.Port))

		defer func() {
			logger.Info("Shutting down server...")
			s.Shutdown()
			logger.Info("Server successfully shutdown.")
		}()

		select {
		case <-ctx.Done():
			return
		case err := <-errChan:
			logger.Fatal(err.Error())
		}
	},
}
