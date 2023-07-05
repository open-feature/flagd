package cmd

import (
	"log"
	"testing"

	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Test_getPort(t *testing.T) {
	logger := getLogger()
	type args struct {
		flagName     string
		value        uint16
		defaultValue uint16
	}
	tests := []struct {
		name        string
		args        args
		want        uint16
		envVarValue string
		envVar      string
	}{
		{
			name: "get provided value",
			args: args{
				flagName:     portFlagName,
				value:        1234,
				defaultValue: defaultServicePort,
			},
			want: 1234,
		},
		{
			name: "get provided value",
			args: args{
				flagName:     portFlagName,
				value:        0,
				defaultValue: defaultServicePort,
			},
			want: defaultServicePort,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("FLAGD_PORT", tt.envVarValue)
			if got := getPortValueOrDefault(tt.args.flagName, tt.args.value, tt.args.defaultValue, logger); got != tt.want {
				t.Errorf("getPort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getLogger() *logger.Logger {
	l, err := logger.NewZapLogger(zapcore.DebugLevel, viper.GetString(logFormatFlagName))
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	logger := logger.NewLogger(l, Debug)
	rtLogger := logger.WithFields(zap.String("component", "start"))
	return rtLogger
}
