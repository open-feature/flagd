package telemetry

import (
	"fmt"
	"github.com/open-feature/flagd/core/pkg/logger"
	"go.uber.org/zap"
)

type OTelErrorsHandler struct {
	logger *logger.Logger
}

func (h OTelErrorsHandler) Handle(err error) {
	h.logger.Debug(fmt.Sprintf("OpenTelemetry Error: %s", err.Error()))
}

func NewOTelErrorsHandler(log *logger.Logger) *OTelErrorsHandler {
	return &OTelErrorsHandler{
		logger: log.WithFields(zap.String("component", "otel")),
	}
}
