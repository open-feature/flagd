package logger

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
This package wraps the zap logging API for use across services without passing a shared logger.
Fields can be added to a requestID using the WriteFields method, these will be added to any
subsequent XxxWithID log calls. To preserve performance ClearFields must be called when the
requestID's thread is closed as a sync map is used internally.
Child loggers can be created from a parent logger using the WithFields method, this child logger
will append the provided fields to all logs, whilst maintaining a reference to the top level
request fields pool.

Example:

WrappedLogger := NewLogger(myLogger)
WrappedLogger.WriteFields("my-id", zap.String("foo", "bar"))
WrappedLogger2 := WrappedLogger.WithFields(zap.String("ping", "pong"))

WrappedLogger.DebugWithID("myID", "my log line")
	=> {"level":"debug","requestID":"myID","foo":"bar","msg":"my log line""}

WrappedLogger2.DebugWithID("myID", "my log line")
	=> {"level":"debug","requestID":"myID","foo":"bar","ping":"pong","msg":"my log line""}

WrappedLogger2.WriteFields("myID", zap.String("food", "bars"))

WrappedLogger.DebugWithID("myID", "my log line")
	=> {"level":"debug","requestID":"myID","foo":"bar","food":"bars","msg":"my log line""}
*/

const RequestIDFieldName = "requestID"

type Logger struct {
	requestFields *sync.Map
	Logger        *zap.Logger
	fields        []zap.Field
	reqIDLogging  bool
}

func (l *Logger) DebugWithID(reqID string, msg string, fields ...zap.Field) {
	if !l.reqIDLogging {
		return
	}
	if ce := l.Logger.Check(zap.DebugLevel, msg); ce != nil {
		fields = append(fields, l.getFieldsForLog(reqID)...)
		ce.Write(fields...)
	}
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	if ce := l.Logger.Check(zap.DebugLevel, msg); ce != nil {
		fields = append(fields, l.fields...)
		ce.Write(fields...)
	}
}

func (l *Logger) InfoWithID(reqID string, msg string, fields ...zap.Field) {
	if !l.reqIDLogging {
		return
	}
	if ce := l.Logger.Check(zap.InfoLevel, msg); ce != nil {
		fields = append(fields, l.getFieldsForLog(reqID)...)
		ce.Write(fields...)
	}
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	if ce := l.Logger.Check(zap.InfoLevel, msg); ce != nil {
		fields = append(fields, l.fields...)
		ce.Write(fields...)
	}
}

func (l *Logger) WarnWithID(reqID string, msg string, fields ...zap.Field) {
	if !l.reqIDLogging {
		return
	}
	if ce := l.Logger.Check(zap.WarnLevel, msg); ce != nil {
		fields = append(fields, l.getFieldsForLog(reqID)...)
		ce.Write(fields...)
	}
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	if ce := l.Logger.Check(zap.WarnLevel, msg); ce != nil {
		fields = append(fields, l.fields...)
		ce.Write(fields...)
	}
}

func (l *Logger) ErrorWithID(reqID string, msg string, fields ...zap.Field) {
	if !l.reqIDLogging {
		return
	}
	if ce := l.Logger.Check(zap.ErrorLevel, msg); ce != nil {
		fields = append(fields, l.getFieldsForLog(reqID)...)
		ce.Write(fields...)
	}
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	if ce := l.Logger.Check(zap.ErrorLevel, msg); ce != nil {
		fields = append(fields, l.fields...)
		ce.Write(fields...)
	}
}

func (l *Logger) FatalWithID(reqID string, msg string, fields ...zap.Field) {
	if !l.reqIDLogging {
		return
	}
	if ce := l.Logger.Check(zap.FatalLevel, msg); ce != nil {
		fields = append(fields, l.getFieldsForLog(reqID)...)
		ce.Write(fields...)
	}
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	if ce := l.Logger.Check(zap.FatalLevel, msg); ce != nil {
		fields = append(fields, l.fields...)
		ce.Write(fields...)
	}
}

// WriteFields adds field key and value pairs to the highest level Logger, they will be applied to all
// subsequent log calls using the matching requestID
func (l *Logger) WriteFields(reqID string, fields ...zap.Field) {
	if !l.reqIDLogging {
		return
	}
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
	fields = append(fields, zap.String(RequestIDFieldName, reqID))
	fields = append(fields, l.fields...)
	return fields
}

// ClearFields clears all stored fields for a given requestID, important for maintaining performance
func (l *Logger) ClearFields(reqID string) {
	if !l.reqIDLogging {
		return
	}
	l.requestFields.Delete(reqID)
}

// NewZapLogger creates a *zap.Logger using the base config
func NewZapLogger(level zapcore.Level, logFormat string) (*zap.Logger, error) {
	cfg := zap.Config{
		Encoding:         logFormat,
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		DisableCaller: false,
	}
	l, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("unable to build logger from config: %w", err)
	}
	return l, nil
}

// NewLogger returns the logging wrapper for a given *zap.logger.
// Noop logger bypasses the setting of fields, improving performance.
// If *zap.Logger is nil a noop logger is set
// and the reqIDLogging argument is overwritten to false
func NewLogger(logger *zap.Logger, reqIDLogging bool) *Logger {
	if logger == nil {
		reqIDLogging = false
		logger = zap.New(nil)
	}
	return &Logger{
		Logger:        logger.WithOptions(zap.AddCallerSkip(1)),
		requestFields: &sync.Map{},
		reqIDLogging:  reqIDLogging,
	}
}

// WithFields creates a new logging wrapper with a predefined base set of fields.
// These fields will be added to each request, but the logger will still
// read/write from the highest level logging wrappers field pool
func (l *Logger) WithFields(fields ...zap.Field) *Logger {
	return &Logger{
		Logger:        l.Logger,
		requestFields: l.requestFields,
		fields:        fields,
		reqIDLogging:  l.reqIDLogging,
	}
}
