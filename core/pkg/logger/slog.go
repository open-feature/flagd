package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync"
)

type slogger struct {
	requestAttrs *sync.Map
	Logger       *slog.Logger
	attrs        []slog.Attr
	reqIDLogging bool
}

var _ Logger = &slogger{}

func (l *slogger) DebugWithID(reqID string, msg string, args ...any) {
	if !l.reqIDLogging {
		return
	}
	attrs := append(makeAttrs(args...), l.getAttrs(reqID)...)
	l.Logger.LogAttrs(context.TODO(), slog.LevelDebug, msg, attrs...)
}

func (l *slogger) Debug(msg string, args ...any) {
	l.Logger.LogAttrs(context.TODO(), slog.LevelDebug, msg, makeAttrs(args...)...)
}

func (l *slogger) InfoWithID(reqID string, msg string, args ...any) {
	if !l.reqIDLogging {
		return
	}
	attrs := append(makeAttrs(args...), l.getAttrs(reqID)...)
	l.Logger.LogAttrs(context.TODO(), slog.LevelInfo, msg, attrs...)
}

func (l *slogger) Info(msg string, args ...any) {
	l.Logger.LogAttrs(context.TODO(), slog.LevelInfo, msg, makeAttrs(args...)...)
}

func (l *slogger) WarnWithID(reqID string, msg string, args ...any) {
	if !l.reqIDLogging {
		return
	}
	attrs := append(makeAttrs(args...), l.getAttrs(reqID)...)
	l.Logger.LogAttrs(context.TODO(), slog.LevelWarn, msg, attrs...)
}

func (l *slogger) Warn(msg string, args ...any) {
	l.Logger.LogAttrs(context.TODO(), slog.LevelWarn, msg, makeAttrs(args...)...)
}

func (l *slogger) ErrorWithID(reqID string, msg string, args ...any) {
	if !l.reqIDLogging {
		return
	}
	attrs := append(makeAttrs(args...), l.getAttrs(reqID)...)
	l.Logger.LogAttrs(context.TODO(), slog.LevelError, msg, attrs...)
}

func (l *slogger) Error(msg string, args ...any) {
	l.Logger.LogAttrs(context.TODO(), slog.LevelError, msg, makeAttrs(args...)...)
}

func (l *slogger) FatalWithID(reqID string, msg string, args ...any) {
	if !l.reqIDLogging {
		return
	}
	attrs := append(makeAttrs(args...), l.getAttrs(reqID)...)
	l.Logger.LogAttrs(context.TODO(), slog.LevelError, msg, attrs...)
	os.Exit(1)
}

func (l *slogger) Fatal(msg string, args ...any) {
	l.Logger.LogAttrs(context.TODO(), slog.LevelError, msg, makeAttrs(args...)...)
	os.Exit(1)
}

// ============================================

func makeAttrs(args ...any) []slog.Attr {
	// if there are an odd number of
	if len(args)%2 != 0 {
		lastArg := args[len(args)-1]
		args[len(args)-1] = "!BADKEY"
		args = append(args, lastArg)
	}
	out := make([]slog.Attr, 0, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		var k string
		switch v := args[i].(type) {
		// TODO: this does not yet support zap.Field types
		case string:
			k = v
		default:
			k = "!BADKEY"
		}
		val := args[i+1]
		out = append(out, slog.Any(k, val))
	}
	return out
}

func (l *slogger) WriteFields(reqID string, fields ...any) {
	if !l.reqIDLogging {
		return
	}
	res := append(l.getAttrs(reqID), makeAttrs(fields)...)
	l.requestAttrs.Store(reqID, res)
}

func (l *slogger) getAttrs(reqID string) []slog.Attr {
	res := []slog.Attr{}
	f, ok := l.requestAttrs.Load(reqID)
	if ok {
		r, ok := f.([]slog.Attr)
		if ok {
			res = r
		}
	}
	return res
}

//func (l *slogger) getAttrsForLog(reqID string) []slog.Attr {
//	fields := l.getAttrs(reqID)
//	fields = append(fields, slog.Any("requestID", reqID))
//	fields = append(fields, l.attrs...)
//	return fields
//}

// ClearFields clears all stored fields for a given requestID, important for maintaining performance
func (l *slogger) ClearFields(reqID string) {
	if !l.reqIDLogging {
		return
	}
	l.requestAttrs.Delete(reqID)
}

// New SlogLogger creates a *slogger using the base config
func NewSlogLogger(level slog.Level, logFormat string, writeFile io.Writer) *slogger {
	var sl *slog.Logger
	opts := &slog.HandlerOptions{
		Level: level,
	}

	if logFormat == "json" {
		sl = slog.New(slog.NewJSONHandler(writeFile, opts))
	} else {
		sl = slog.New(slog.NewTextHandler(writeFile, opts))
	}

	return &slogger{
		requestAttrs: &sync.Map{},
		Logger:       sl,
		attrs:        []slog.Attr{},
		reqIDLogging: false,
	}
}

// Flagd uses WithFields, we'll need to implement.
func (l *slogger) With(attrs ...any) Logger {
	return &slogger{
		Logger:       l.Logger,
		requestAttrs: l.requestAttrs,
		attrs:        makeAttrs(attrs),
		reqIDLogging: l.reqIDLogging,
	}
}
