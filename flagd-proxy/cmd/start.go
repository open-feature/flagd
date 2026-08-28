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
	"time"

	"github.com/open-feature/flagd/core/pkg/fips"
	"github.com/open-feature/flagd/core/pkg/logger"
	iService "github.com/open-feature/flagd/core/pkg/service"
	"github.com/open-feature/flagd/flagd-proxy/pkg/service"
	"github.com/open-feature/flagd/flagd-proxy/pkg/service/subscriptions"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

// start

const (
	logFormatFlagName                    = "log-format"
	managementPortFlagName               = "management-port"
	portFlagName                         = "port"
	keepAliveMinTimeFlagName             = "keep-alive-min-time"
	keepAlivePermitWithoutStreamFlagName = "keep-alive-permit-without-stream"
)

func init() {
	flags := startCmd.Flags()

	// allows environment variables to use _ instead of -
	flags.Int32P(portFlagName, "p", 8015, "Port to listen on")
	flags.Int32P(managementPortFlagName, "m", 8016, "Management port")
	flags.StringP(logFormatFlagName, "z", "console", "Set the logging format, e.g. console or json")
	flags.Duration(keepAliveMinTimeFlagName, 30*time.Second, "Minimum interval the flag sync gRPC server permits between client keepalive pings. Pings arriving more frequently than this are rejected with GOAWAY (ENHANCE_YOUR_CALM). Defaults to 30s.")
	flags.Bool(keepAlivePermitWithoutStreamFlagName, true, "Permit clients of the flag sync gRPC server to send keepalive pings even when there is no active stream. Defaults to true.")

	_ = viper.BindPFlag(logFormatFlagName, flags.Lookup(logFormatFlagName))
	_ = viper.BindPFlag(managementPortFlagName, flags.Lookup(managementPortFlagName))
	_ = viper.BindPFlag(portFlagName, flags.Lookup(portFlagName))
	_ = viper.BindPFlag(keepAliveMinTimeFlagName, flags.Lookup(keepAliveMinTimeFlagName))
	_ = viper.BindPFlag(keepAlivePermitWithoutStreamFlagName, flags.Lookup(keepAlivePermitWithoutStreamFlagName))
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

		logger.Info(fmt.Sprintf("build variant: %s, FIPS 140-3 mode: %s", fips.Variant, fips.Current()))
		if err := fips.Check(); err != nil {
			logger.Fatal(fmt.Sprintf("refusing to start: %s", err))
		}

		ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

		syncStore := subscriptions.NewManager(ctx, logger)
		s := service.NewServer(ctx, logger, syncStore)

		cfg := iService.Configuration{
			ReadinessProbe: func() bool { return true },
			Port:           viper.GetUint16(portFlagName),
			ManagementPort: viper.GetUint16(managementPortFlagName),

			KeepAliveMinTime:             viper.GetDuration(keepAliveMinTimeFlagName),
			KeepAlivePermitWithoutStream: viper.GetBool(keepAlivePermitWithoutStreamFlagName),
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
