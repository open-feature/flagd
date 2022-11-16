package logger

import (
	"sync"

	"go.uber.org/zap"
)

type Logger struct {
	requestFields *sync.Map
	Logger        *zap.Logger
	fields        []zap.Field
}

func (l *Logger) DebugWithID(reqID string, msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, l.getFieldsForLog(reqID)...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	f := append(l.fields, fields...)
	l.Logger.Debug(msg, f...)
}

func (l *Logger) InfoWithID(reqID string, msg string, fields ...zap.Field) {
	l.Logger.Info(msg, l.getFieldsForLog(reqID)...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	f := append(l.fields, fields...)
	l.Logger.Info(msg, f...)
}

func (l *Logger) WarnWithID(reqID string, msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, l.getFieldsForLog(reqID)...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	f := append(l.fields, fields...)
	l.Logger.Warn(msg, f...)
}

func (l *Logger) ErrorWithID(reqID string, msg string, fields ...zap.Field) {
	l.Logger.Error(msg, l.getFieldsForLog(reqID)...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	f := append(l.fields, fields...)
	l.Logger.Error(msg, f...)
}

func (l *Logger) FatalWithID(reqID string, msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, l.getFieldsForLog(reqID)...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	f := append(l.fields, fields...)
	l.Logger.Debug(msg, f...)
}

func (l *Logger) WriteFields(reqID string, fields ...zap.Field) {
	res := append(l.getFields(reqID), fields...)
	l.requestFields.Store(reqID, res)
}

func (l *Logger) getFields(reqID string) []zap.Field {
	res := []zap.Field{}
	f, ok := l.requestFields.Load(reqID)
	if ok {
		r, ok := f.([]zap.Field)
		if ok {
			res = r
		}
	}
	return res
}

func (l *Logger) getFieldsForLog(reqID string) []zap.Field {
	fields := l.getFields(reqID)
	fields = append(fields, zap.String("requestID", reqID))
	fields = append(fields, l.fields...)
	return fields
}

func (l *Logger) ClearFields(reqID string) {
	l.requestFields.Delete(reqID)
}

func NewLogger(logger *zap.Logger) *Logger {
	if logger == nil {
		logger = zap.New(nil)
	}
	return &Logger{
		Logger:        logger.WithOptions(zap.AddCallerSkip(1)),
		requestFields: &sync.Map{},
	}
}

func (l *Logger) WithFields(fields ...zap.Field) *Logger {
	return &Logger{
		Logger:        l.Logger,
		requestFields: l.requestFields,
		fields:        fields,
	}
}
