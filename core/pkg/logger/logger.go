package logger

import (
	"log/slog"
	"os"
)

/*
This package wraps both zap and slog API for use across services without passing a shared logger.
Fields can be added to a requestID using the WriteFields method, these will be added to any
subsequent XxxWithID log calls. To preserve performance ClearFields must be called when the
requestID's thread is closed as a sync map is used internally.
Child loggers can be created from a parent logger using the WithFields method, this child logger
will append the provided fields to all logs, whilst maintaining a reference to the top level
request fields pool.

Example:

WrappedLogger := logger.New("slog", debug, "format")
WrappedLogger.WriteFields("my-id", "foo", "bar")
WrappedLogger2 := WrappedLogger.With("ping", "pong")

WrappedLogger.DebugWithID("myID", "my log line")
	=> {"level":"debug","requestID":"myID","foo":"bar","msg":"my log line""}
*/

const RequestIDFieldName = "requestID"

type Logger interface {
	DebugWithID(reqID string, msg string, args ...any)
	Debug(msg string, fields ...any)
	InfoWithID(reqID string, msg string, args ...any)
	Info(msg string, fields ...any)
	WarnWithID(reqID string, msg string, args ...any)
	Warn(msg string, fields ...any)
	ErrorWithID(reqID string, msg string, args ...any)
	Error(msg string, fields ...any)
	FatalWithID(reqID string, msg string, args ...any)
	Fatal(msg string, args ...any)
	With(args ...any) Logger
	WriteFields(reqID string, args ...any)
	ClearFields(reqID string)
}

func New(loggerType string, Debug bool, logFormat string) Logger {
	if loggerType == "slog" {
		level := slog.LevelInfo
		if Debug {
			level = slog.LevelDebug
		}
		return NewSlogLogger(level, logFormat, os.Stdout)
	}
	return NewZapLogger(Debug, logFormat)
}
