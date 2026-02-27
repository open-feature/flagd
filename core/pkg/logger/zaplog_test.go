package logger

import (
	"strings"
	"sync"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// setupTestZapLogger creates a zapLogger backed by an observer core for capturing output.
func setupTestZapLogger(level zapcore.Level) (*zapLogger, *observer.ObservedLogs) {
	core, logs := observer.New(level)
	l := zap.New(core)
	return &zapLogger{
		Logger:        l,
		fields:        []zap.Field{},
		requestFields: &sync.Map{},
		reqIDLogging:  true,
	}, logs
}

func TestZapLoggerLoggingMethods(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		reqID          string
		msg            string
		args           []any
		expectOutput   bool
		expectInOutput string
	}{
		{
			name:           "Debug method",
			method:         "Debug",
			msg:            "test log",
			args:           []any{},
			expectOutput:   true,
			expectInOutput: "test log",
		},
		{
			name:           "Debug with args",
			method:         "Debug",
			msg:            "test log with args",
			args:           []any{"key", "value"},
			expectOutput:   true,
			expectInOutput: "test log with args",
		},
		{
			name:           "DebugWithID method",
			method:         "DebugWithID",
			reqID:          "test-request-id",
			msg:            "test log",
			args:           []any{},
			expectOutput:   true,
			expectInOutput: "test log",
		},
		{
			name:           "Info method",
			method:         "Info",
			msg:            "test log",
			args:           []any{},
			expectOutput:   true,
			expectInOutput: "test log",
		},
		{
			name:           "Info with args",
			method:         "Info",
			msg:            "test log with args",
			args:           []any{"key", "value"},
			expectOutput:   true,
			expectInOutput: "test log with args",
		},
		{
			name:           "InfoWithID method",
			method:         "InfoWithID",
			reqID:          "test-request-id",
			msg:            "test log",
			args:           []any{},
			expectOutput:   true,
			expectInOutput: "test log",
		},
		{
			name:           "Warn method",
			method:         "Warn",
			msg:            "test log",
			args:           []any{},
			expectOutput:   true,
			expectInOutput: "test log",
		},
		{
			name:           "Warn with args",
			method:         "Warn",
			msg:            "test log with args",
			args:           []any{"key", "value"},
			expectOutput:   true,
			expectInOutput: "test log with args",
		},
		{
			name:           "WarnWithID method",
			method:         "WarnWithID",
			reqID:          "test-request-id",
			msg:            "test log",
			args:           []any{},
			expectOutput:   true,
			expectInOutput: "test log",
		},
		{
			name:           "Error method",
			method:         "Error",
			msg:            "test log",
			args:           []any{},
			expectOutput:   true,
			expectInOutput: "test log",
		},
		{
			name:           "Error with args",
			method:         "Error",
			msg:            "test log with args",
			args:           []any{"key", "value"},
			expectOutput:   true,
			expectInOutput: "test log with args",
		},
		{
			name:           "ErrorWithID method",
			method:         "ErrorWithID",
			reqID:          "test-request-id",
			msg:            "test log",
			args:           []any{},
			expectOutput:   true,
			expectInOutput: "test log",
		},
		{
			name:           "DebugWithID disabled reqIDLogging",
			method:         "DebugWithID",
			reqID:          "test-request-id",
			msg:            "test log",
			args:           []any{},
			expectOutput:   false,
			expectInOutput: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, logs := setupTestZapLogger(zapcore.DebugLevel)

			if tt.name == "DebugWithID disabled reqIDLogging" {
				logger.reqIDLogging = false
			}

			switch tt.method {
			case "Debug":
				logger.Debug(tt.msg, tt.args...)
			case "DebugWithID":
				logger.DebugWithID(tt.reqID, tt.msg, tt.args...)
			case "Info":
				logger.Info(tt.msg, tt.args...)
			case "InfoWithID":
				logger.InfoWithID(tt.reqID, tt.msg, tt.args...)
			case "Warn":
				logger.Warn(tt.msg, tt.args...)
			case "WarnWithID":
				logger.WarnWithID(tt.reqID, tt.msg, tt.args...)
			case "Error":
				logger.Error(tt.msg, tt.args...)
			case "ErrorWithID":
				logger.ErrorWithID(tt.reqID, tt.msg, tt.args...)
			}

			if tt.expectOutput {
				if logs.Len() == 0 {
					t.Errorf("Expected output but got none")
				}
				if tt.expectInOutput != "" {
					entry := logs.All()[0]
					if !strings.Contains(entry.Message, tt.expectInOutput) {
						t.Errorf("Expected output to contain %q, but got: %s", tt.expectInOutput, entry.Message)
					}
				}
			} else {
				if logs.Len() != 0 {
					t.Errorf("Expected no output but got %d entries", logs.Len())
				}
			}
		})
	}
}

func TestZapLoggerUtilityMethods(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "WriteFields and getFields",
			testFunc: func(t *testing.T) {
				logger, _ := setupTestZapLogger(zapcore.DebugLevel)
				reqID := "test-request-id"

				logger.WriteFields(reqID, "key1", "value1", "key2", "value2")

				fields := logger.getFields(reqID)
				if len(fields) == 0 {
					t.Error("Expected fields to be stored, but got none")
				}
			},
		},
		{
			name: "WriteFields disabled reqIDLogging",
			testFunc: func(t *testing.T) {
				logger, _ := setupTestZapLogger(zapcore.DebugLevel)
				logger.reqIDLogging = false

				reqID := "test-request-id"
				logger.WriteFields(reqID, "key1", "value1")

				fields := logger.getFields(reqID)
				if len(fields) != 0 {
					t.Error("Expected no fields when reqIDLogging is disabled")
				}
			},
		},
		{
			name: "ClearFields",
			testFunc: func(t *testing.T) {
				logger, _ := setupTestZapLogger(zapcore.DebugLevel)
				reqID := "test-request-id"

				logger.WriteFields(reqID, "key1", "value1")

				fields := logger.getFields(reqID)
				if len(fields) == 0 {
					t.Error("Expected fields to be stored before clearing")
				}

				logger.ClearFields(reqID)

				fields = logger.getFields(reqID)
				if len(fields) != 0 {
					t.Error("Expected fields to be cleared")
				}
			},
		},
		{
			name: "ClearFields disabled reqIDLogging",
			testFunc: func(t *testing.T) {
				logger, _ := setupTestZapLogger(zapcore.DebugLevel)
				logger.reqIDLogging = false

				// Should not panic
				logger.ClearFields("test-request-id")
			},
		},
		{
			name: "With method",
			testFunc: func(t *testing.T) {
				logger, _ := setupTestZapLogger(zapcore.DebugLevel)

				child := logger.With("childKey", "childValue")

				if child == logger {
					t.Error("Expected With() to return a new logger instance")
				}

				childZap, ok := child.(*zapLogger)
				if !ok {
					t.Error("Expected With() to return *zapLogger")
				}

				if childZap.requestFields != logger.requestFields {
					t.Error("Expected child logger to share requestFields with parent")
				}

				if len(childZap.fields) == 0 {
					t.Error("Expected child logger to have additional fields")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func TestMakeFields(t *testing.T) {
	tests := []struct {
		name     string
		args     []any
		expected int
	}{
		{
			name:     "even number of args",
			args:     []any{"key1", "value1", "key2", "value2"},
			expected: 2,
		},
		{
			name:     "odd number of args",
			args:     []any{"key1", "value1", "key2"},
			expected: 2, // Should add !BADKEY for the orphaned value
		},
		{
			name:     "empty args",
			args:     []any{},
			expected: 0,
		},
		{
			name:     "single arg",
			args:     []any{"lone_value"},
			expected: 1, // Should become !BADKEY, lone_value
		},
		{
			name:     "non-string key",
			args:     []any{123, "value1"},
			expected: 1, // Should use !BADKEY for non-string key
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := makeFields(tt.args...)

			if len(fields) != tt.expected {
				t.Errorf("Expected %d fields, got %d", tt.expected, len(fields))
			}
		})
	}
}

func TestMakeFieldsKeyValues(t *testing.T) {
	tests := []struct {
		name         string
		args         []any
		expectedKeys []string
		expectedVals []any
	}{
		{
			name:         "string keys are captured correctly",
			args:         []any{"component", "evaluator", "service", "flagd"},
			expectedKeys: []string{"component", "service"},
			expectedVals: []any{"evaluator", "flagd"},
		},
		{
			name:         "non-string key gets BADKEY",
			args:         []any{42, "value"},
			expectedKeys: []string{"!BADKEY"},
			expectedVals: []any{"value"},
		},
		{
			name:         "odd args get BADKEY key for orphan",
			args:         []any{"key1", "value1", "orphan"},
			expectedKeys: []string{"key1", "!BADKEY"},
			expectedVals: []any{"value1", "orphan"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := makeFields(tt.args...)

			if len(fields) != len(tt.expectedKeys) {
				t.Fatalf("Expected %d fields, got %d", len(tt.expectedKeys), len(fields))
			}

			for i, f := range fields {
				if f.Key != tt.expectedKeys[i] {
					t.Errorf("Field %d: expected key %q, got %q", i, tt.expectedKeys[i], f.Key)
				}
			}
		})
	}
}
