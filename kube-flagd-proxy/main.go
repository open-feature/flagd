package main

import (
	"log"

	"github.com/open-feature/flagd/core/pkg/logger"
	"go.uber.org/zap/zapcore"
)

func main() {
	l, err := logger.NewZapLogger(zapcore.DebugLevel, "console")
	if err != nil {
		log.Fatalf("initialize zap logger: %v", err)
	}

	l.Info("kube-flagd-proxy")
}
