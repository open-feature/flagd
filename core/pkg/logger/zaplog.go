package logger

import (
	"fmt"
	"log"
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

type zapLogger struct {
	requestFields *sync.Map
	Logger        *zap.Logger
	fields        []zap.Field
	reqIDLogging  bool
}

// zapLogger explicitly implements Logger
var _ Logger = &zapLogger{}

func makeFields(args ...any) []zap.Field {
	if len(args) == 0 {
		return nil
	}
	// if there are an odd number of
	if len(args)%2 != 0 {
		lastArg := args[len(args)-1]
		args[len(args)-1] = "!BADKEY"
		args = append(args, lastArg)
	}
	out := make([]zap.Field, 0, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		var k string
		switch v := args[i].(type) {
		case string:
			k = v
		default:
			k = "!BADKEY"
		}
		val := args[i+1]
		out = append(out, zap.Any(k, val))
	}
	return out
}

func (l *zapLogger) DebugWithID(reqID string, msg string, args ...any) {
	if !l.reqIDLogging {
		return
	}
	fields := makeFields(args)
	if ce := l.Logger.Check(zap.DebugLevel, msg); ce != nil {
		fields = append(fields, l.getFieldsForLog(reqID)...)
		ce.Write(fields...)
	}
}

func (l *zapLogger) Debug(msg string, args ...any) {
	fields := makeFields(args)
	if ce := l.Logger.Check(zap.DebugLevel, msg); ce != nil {
		fields = append(fields, l.fields...)
		ce.Write(fields...)
	}
}

func (l *zapLogger) InfoWithID(reqID string, msg string, args ...any) {
	if !l.reqIDLogging {
		return
	}
	fields := makeFields(args)
	if ce := l.Logger.Check(zap.InfoLevel, msg); ce != nil {
		fields = append(fields, l.getFieldsForLog(reqID)...)
		ce.Write(fields...)
	}
}

func (l *zapLogger) Info(msg string, args ...any) {
	fields := makeFields(args)
	if ce := l.Logger.Check(zap.InfoLevel, msg); ce != nil {
		fields = append(fields, l.fields...)
		ce.Write(fields...)
	}
}

func (l *zapLogger) WarnWithID(reqID string, msg string, args ...any) {
	if !l.reqIDLogging {
		return
	}
	fields := makeFields(args)
	if ce := l.Logger.Check(zap.WarnLevel, msg); ce != nil {
		fields = append(fields, l.getFieldsForLog(reqID)...)
		ce.Write(fields...)
	}
}

func (l *zapLogger) Warn(msg string, args ...any) {
	fields := makeFields(args)
	if ce := l.Logger.Check(zap.WarnLevel, msg); ce != nil {
		fields = append(fields, l.fields...)
		ce.Write(fields...)
	}
}

func (l *zapLogger) ErrorWithID(reqID string, msg string, args ...any) {
	if !l.reqIDLogging {
		return
	}
	fields := makeFields(args)
	if ce := l.Logger.Check(zap.ErrorLevel, msg); ce != nil {
		fields = append(fields, l.getFieldsForLog(reqID)...)
		ce.Write(fields...)
	}
}

func (l *zapLogger) Error(msg string, args ...any) {
	fields := makeFields(args)
	if ce := l.Logger.Check(zap.ErrorLevel, msg); ce != nil {
		fields = append(fields, l.fields...)
		ce.Write(fields...)
	}
}

func (l *zapLogger) FatalWithID(reqID string, msg string, args ...any) {
	if !l.reqIDLogging {
		return
	}
	fields := makeFields(args)
	if ce := l.Logger.Check(zap.FatalLevel, msg); ce != nil {
		fields = append(fields, l.getFieldsForLog(reqID)...)
		ce.Write(fields...)
	}
}

func (l *zapLogger) Fatal(msg string, args ...any) {
	fields := makeFields(args)
	if ce := l.Logger.Check(zap.FatalLevel, msg); ce != nil {
		fields = append(fields, l.fields...)
		ce.Write(fields...)
	}
}

// WriteFields adds field key and value pairs to the highest level Logger, they will be applied to all
// subsequent log calls using the matching requestID
func (l *zapLogger) WriteFields(reqID string, fields ...any) {
	if !l.reqIDLogging {
		return
	}
	res := append(l.getFields(reqID), makeFields(fields)...)
	l.requestFields.Store(reqID, res)
}

func (l *zapLogger) getFields(reqID string) []zap.Field {
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

func (l *zapLogger) getFieldsForLog(reqID string) []zap.Field {
	fields := l.getFields(reqID)
	fields = append(fields, zap.String(RequestIDFieldName, reqID))
	fields = append(fields, l.fields...)
	return fields
}

// ClearFields clears all stored fields for a given requestID, important for maintaining performance
func (l *zapLogger) ClearFields(reqID string) {
	if !l.reqIDLogging {
		return
	}
	l.requestFields.Delete(reqID)
}

// NewZapLogger creates a *zap.Logger using the base config
func NewZapLogger(Debug bool, logFormat string) Logger {
	level := zapcore.InfoLevel
	if Debug {
		level = zapcore.DebugLevel
	}
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
		log.Fatalf("cant initialize the zap logger: %v", fmt.Errorf("unable to build logger from config: %w", err))
	}
	return &zapLogger{
		Logger:        l.WithOptions(zap.AddCallerSkip(1)),
		fields:        []zap.Field{},
		requestFields: &sync.Map{},
		reqIDLogging:  Debug,
	}
}

// WithFields creates a new logging wrapper with a predefined base set of fields.
// These fields will be added to each request, but the logger will still
// read/write from the highest level logging wrappers field pool
func (l *zapLogger) With(fields ...any) Logger {
	return &zapLogger{
		Logger:        l.Logger,
		requestFields: l.requestFields,
		fields:        makeFields(fields),
		reqIDLogging:  l.reqIDLogging,
	}
}
