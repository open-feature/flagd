package logger

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestFieldStorageAndRetrieval(t *testing.T) {
	tests := map[string]struct {
		fields []zap.Field
	}{
		"happyPath": {
			fields: []zap.Field{
				zap.String("this", "that"),
				zap.Strings("this2", []string{"that2", "that2"}),
			},
		},
	}
	for name, test := range tests {
		l := NewLogger(&zap.Logger{}, true)
		l.WriteFields(name, test.fields...)
		returnedFields := l.getFields(name)
		if !reflect.DeepEqual(returnedFields, test.fields) {
			t.Error("returned fields to not match the input", test.fields, returnedFields)
		}
	}
}

func TestLoggerChildOperation(t *testing.T) {
	id := "test"
	// create parent logger
	p := NewLogger(&zap.Logger{}, true)
	// add field 1
	field1 := zap.Int("field", 1)
	p.WriteFields(id, field1)

	// create child logger with field 2
	field2 := zap.Int("field", 2)
	c := p.WithFields(field2)

	if !reflect.DeepEqual(c.getFields(id), []zapcore.Field{field1}) {
		t.Error("1: child logger contains incorrect fields ", c.getFieldsForLog(id))
	}
	if !reflect.DeepEqual(p.getFields(id), []zapcore.Field{field1}) {
		t.Error("1: parent logger contains incorrect fields ", c.getFields(id))
	}

	// add field 3 to the child, should be present in both
	field3 := zap.Int("field", 3)
	c.WriteFields(id, field3)

	if !reflect.DeepEqual(c.getFields(id), []zapcore.Field{field1, field3}) {
		t.Error("1: child logger contains incorrect fields ", c.getFieldsForLog(id))
	}
	if !reflect.DeepEqual(p.getFields(id), []zapcore.Field{field1, field3}) {
		t.Error("1: parent logger contains incorrect fields ", c.getFields(id))
	}

	// ensure child logger appends field 2
	logFields := c.getFieldsForLog(id)
	field2Found := false
	for _, field := range logFields {
		if field == field2 {
			field2Found = true
		}
	}
	if !field2Found {
		t.Error("field 2 is missing from the child logger getFieldsForLog response")
	}

	// ensure parent logger does not
	logFields = p.getFieldsForLog(id)
	field2Found = false
	for _, field := range logFields {
		if field == field2 {
			field2Found = true
		}
	}
	if field2Found {
		t.Error("field 2 is present in the parent logger getFieldsForLog response")
	}
}

func TestNewZapLogger(t *testing.T) {
	tests := []struct {
		name      string
		level     zapcore.Level
		logFormat string
		wantErr   bool
	}{
		{
			name:      "valid json format",
			level:     zapcore.InfoLevel,
			logFormat: "json",
			wantErr:   false,
		},
		{
			name:      "valid console format",
			level:     zapcore.DebugLevel,
			logFormat: "console",
			wantErr:   false,
		},
		{
			name:      "invalid format returns error",
			level:     zapcore.InfoLevel,
			logFormat: "unsupported-format",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewZapLogger(tt.level, tt.logFormat)
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, logger)
				assert.Contains(t, err.Error(), "unable to build logger from config")
			} else {
				require.NoError(t, err)
				assert.NotNil(t, logger)
			}
		})
	}
}

func TestNewLogger(t *testing.T) {
	t.Run("nil logger falls back to noop logger and disables reqIDLogging", func(t *testing.T) {
		l := NewLogger(nil, true)
		require.NotNil(t, l)
		assert.False(t, l.reqIDLogging)
		assert.NotNil(t, l.Logger)
		assert.NotNil(t, l.requestFields)
	})

	t.Run("custom logger with reqIDLogging enabled", func(t *testing.T) {
		l := NewLogger(zap.NewNop(), true)
		require.NotNil(t, l)
		assert.True(t, l.reqIDLogging)
		assert.NotNil(t, l.Logger)
	})

	t.Run("custom logger with reqIDLogging disabled", func(t *testing.T) {
		l := NewLogger(zap.NewNop(), false)
		require.NotNil(t, l)
		assert.False(t, l.reqIDLogging)
	})
}

func TestClearFields(t *testing.T) {
	t.Run("clears fields when reqIDLogging is enabled", func(t *testing.T) {
		l := NewLogger(zap.NewNop(), true)
		reqID := "req-clear-1"

		l.WriteFields(reqID, zap.String("key", "val"))
		assert.Len(t, l.getFields(reqID), 1)

		l.ClearFields(reqID)
		assert.Empty(t, l.getFields(reqID))
	})

	t.Run("noop when reqIDLogging is disabled", func(t *testing.T) {
		l := NewLogger(zap.NewNop(), false)
		reqID := "req-clear-2"

		// ClearFields should safely return early
		l.ClearFields(reqID)
		assert.Empty(t, l.getFields(reqID))
	})
}

func TestWriteFields_Disabled(t *testing.T) {
	l := NewLogger(zap.NewNop(), false)
	reqID := "req-disabled"

	l.WriteFields(reqID, zap.String("key", "val"))
	assert.Empty(t, l.getFields(reqID))
}

func TestGetFields_InvalidTypeFallback(t *testing.T) {
	l := NewLogger(zap.NewNop(), true)
	reqID := "req-corrupt"

	// Store non-[]zap.Field type directly in sync.Map
	l.requestFields.Store(reqID, "invalid-type")
	fields := l.getFields(reqID)
	assert.Empty(t, fields)
}

func TestStandardLogLevels(t *testing.T) {
	core, recorded := observer.New(zapcore.DebugLevel)
	zapLog := zap.New(core, zap.WithFatalHook(zapcore.WriteThenPanic))
	l := NewLogger(zapLog, true)
	child := l.WithFields(zap.String("childKey", "childVal"))

	tests := []struct {
		name     string
		logFunc  func(msg string, fields ...zap.Field)
		expected zapcore.Level
		msg      string
		panics   bool
	}{
		{
			name:     "Debug",
			logFunc:  child.Debug,
			expected: zapcore.DebugLevel,
			msg:      "debug message",
		},
		{
			name:     "Info",
			logFunc:  child.Info,
			expected: zapcore.InfoLevel,
			msg:      "info message",
		},
		{
			name:     "Warn",
			logFunc:  child.Warn,
			expected: zapcore.WarnLevel,
			msg:      "warn message",
		},
		{
			name:     "Error",
			logFunc:  child.Error,
			expected: zapcore.ErrorLevel,
			msg:      "error message",
		},
		{
			name:     "Fatal",
			logFunc:  child.Fatal,
			expected: zapcore.FatalLevel,
			msg:      "fatal message",
			panics:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorded.TakeAll() // drain
			if tt.panics {
				assert.Panics(t, func() {
					tt.logFunc(tt.msg, zap.String("callKey", "callVal"))
				})
			} else {
				tt.logFunc(tt.msg, zap.String("callKey", "callVal"))
			}

			entries := recorded.TakeAll()
			require.Len(t, entries, 1)
			assert.Equal(t, tt.expected, entries[0].Level)
			assert.Equal(t, tt.msg, entries[0].Message)

			fieldMap := entries[0].ContextMap()
			assert.Equal(t, "childVal", fieldMap["childKey"])
			assert.Equal(t, "callVal", fieldMap["callKey"])
		})
	}

	t.Run("level check skips logging when level disabled", func(t *testing.T) {
		warnCore, warnRecorded := observer.New(zapcore.WarnLevel)
		lWarn := NewLogger(zap.New(warnCore), true)

		lWarn.Debug("should not appear")
		lWarn.Info("should not appear")
		assert.Equal(t, 0, warnRecorded.Len())
	})
}

func TestWithIDLogLevels(t *testing.T) {
	reqID := "test-request-id"

	t.Run("logs with requestID and stored fields when reqIDLogging enabled", func(t *testing.T) {
		core, recorded := observer.New(zapcore.DebugLevel)
		zapLog := zap.New(core, zap.WithFatalHook(zapcore.WriteThenPanic))
		l := NewLogger(zapLog, true)
		l.WriteFields(reqID, zap.String("contextKey", "contextVal"))
		child := l.WithFields(zap.String("childKey", "childVal"))

		tests := []struct {
			name     string
			logFunc  func(reqID string, msg string, fields ...zap.Field)
			expected zapcore.Level
			msg      string
			panics   bool
		}{
			{
				name:     "DebugWithID",
				logFunc:  child.DebugWithID,
				expected: zapcore.DebugLevel,
				msg:      "debug with id message",
			},
			{
				name:     "InfoWithID",
				logFunc:  child.InfoWithID,
				expected: zapcore.InfoLevel,
				msg:      "info with id message",
			},
			{
				name:     "WarnWithID",
				logFunc:  child.WarnWithID,
				expected: zapcore.WarnLevel,
				msg:      "warn with id message",
			},
			{
				name:     "ErrorWithID",
				logFunc:  child.ErrorWithID,
				expected: zapcore.ErrorLevel,
				msg:      "error with id message",
			},
			{
				name:     "FatalWithID",
				logFunc:  child.FatalWithID,
				expected: zapcore.FatalLevel,
				msg:      "fatal with id message",
				panics:   true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				recorded.TakeAll()
				if tt.panics {
					assert.Panics(t, func() {
						tt.logFunc(reqID, tt.msg, zap.String("callKey", "callVal"))
					})
				} else {
					tt.logFunc(reqID, tt.msg, zap.String("callKey", "callVal"))
				}

				entries := recorded.TakeAll()
				require.Len(t, entries, 1)
				assert.Equal(t, tt.expected, entries[0].Level)
				assert.Equal(t, tt.msg, entries[0].Message)

				fieldMap := entries[0].ContextMap()
				assert.Equal(t, reqID, fieldMap[RequestIDFieldName])
				assert.Equal(t, "contextVal", fieldMap["contextKey"])
				assert.Equal(t, "childVal", fieldMap["childKey"])
				assert.Equal(t, "callVal", fieldMap["callKey"])
			})
		}
	})

	t.Run("noop when reqIDLogging is disabled", func(t *testing.T) {
		core, recorded := observer.New(zapcore.DebugLevel)
		zapLog := zap.New(core, zap.WithFatalHook(zapcore.WriteThenPanic))
		l := NewLogger(zapLog, false)

		l.DebugWithID(reqID, "debug")
		l.InfoWithID(reqID, "info")
		l.WarnWithID(reqID, "warn")
		l.ErrorWithID(reqID, "error")
		l.FatalWithID(reqID, "fatal")

		assert.Equal(t, 0, recorded.Len())
	})

	t.Run("level check skips logging when level disabled", func(t *testing.T) {
		errorCore, errorRecorded := observer.New(zapcore.ErrorLevel)
		lError := NewLogger(zap.New(errorCore), true)

		lError.DebugWithID(reqID, "debug skipped")
		lError.InfoWithID(reqID, "info skipped")
		lError.WarnWithID(reqID, "warn skipped")

		assert.Equal(t, 0, errorRecorded.Len())
	})
}
