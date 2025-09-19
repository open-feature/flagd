package logger

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"
)

// setupTestLogger creates a slog logger that writes to the provided buffer
func setupTestLogger(w io.Writer) *slog.Logger {
	handler := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return slog.New(handler)
}

func TestSloggerLoggingMethods(t *testing.T) {
	// Create a named buffer for capturing log output
	testBuffer := &bytes.Buffer{}
	// Create slogger instance
	logger := NewSlogLogger(slog.LevelDebug, "text", testBuffer)
	logger.reqIDLogging = true // Enable request ID logging for WithID tests

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
			// Clear buffer before each test
			testBuffer.Reset()

			// Special case for testing disabled reqIDLogging
			if tt.name == "DebugWithID disabled reqIDLogging" {
				logger.reqIDLogging = false
				defer func() { logger.reqIDLogging = true }()
			}

			// Call the appropriate method
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
			case "Fatal":
				logger.Fatal(tt.msg, tt.args...)
			case "FatalWithID":
				logger.FatalWithID(tt.reqID, tt.msg, tt.args...)
			}

			output := testBuffer.String()

			if tt.expectOutput {
				if output == "" {
					t.Errorf("Expected output but got none")
				}
				if tt.expectInOutput != "" && !strings.Contains(output, tt.expectInOutput) {
					t.Errorf("Expected output to contain %q, but got: %s", tt.expectInOutput, output)
				}
			} else {
				if output != "" {
					t.Errorf("Expected no output but got: %s", output)
				}
			}
		})
	}
}

func TestSloggerUtilityMethods(t *testing.T) {
	// Create a named buffer for capturing log output
	testBuffer := &bytes.Buffer{}

	logger := NewSlogLogger(slog.LevelDebug, "text", testBuffer)
	logger.reqIDLogging = true

	tests := []struct {
		name     string
		testFunc func(t *testing.T, l *slogger)
	}{
		{
			name: "WriteFields and getAttrs",
			testFunc: func(t *testing.T, l *slogger) {
				reqID := "test-request-id"

				// Test WriteFields
				l.WriteFields(reqID, "key1", "value1", "key2", "value2")

				// Test getAttrs
				attrs := l.getAttrs(reqID)
				if len(attrs) == 0 {
					t.Error("Expected attributes to be stored, but got none")
				}
			},
		},
		{
			name: "WriteFields disabled reqIDLogging",
			testFunc: func(t *testing.T, l *slogger) {
				// Create a fresh logger with reqIDLogging disabled
				freshLogger := NewSlogLogger(slog.LevelDebug, "text", os.Stdout)
				freshLogger.reqIDLogging = false

				reqID := "test-request-id"
				freshLogger.WriteFields(reqID, "key1", "value1")

				attrs := freshLogger.getAttrs(reqID)
				if len(attrs) != 0 {
					t.Error("Expected no attributes when reqIDLogging is disabled")
				}
			},
		},
		{
			name: "ClearFields",
			testFunc: func(t *testing.T, l *slogger) {
				reqID := "test-request-id"

				// Add some fields
				l.WriteFields(reqID, "key1", "value1")

				// Verify fields exist
				attrs := l.getAttrs(reqID)
				if len(attrs) == 0 {
					t.Error("Expected attributes to be stored before clearing")
				}

				// Clear fields
				l.ClearFields(reqID)

				// Verify fields are cleared
				attrs = l.getAttrs(reqID)
				if len(attrs) != 0 {
					t.Error("Expected attributes to be cleared")
				}
			},
		},
		{
			name: "ClearFields disabled reqIDLogging",
			testFunc: func(t *testing.T, l *slogger) {
				// Create a fresh logger with reqIDLogging disabled
				freshLogger := NewSlogLogger(slog.LevelDebug, "text", testBuffer)
				freshLogger.reqIDLogging = false

				reqID := "test-request-id"

				// This should not panic even when reqIDLogging is disabled
				freshLogger.ClearFields(reqID)
			},
		},
		{
			name: "With method",
			testFunc: func(t *testing.T, l *slogger) {
				// Create child logger with additional attributes
				child := l.With("childKey", "childValue")

				// Verify child logger is different instance
				if child == l {
					t.Error("Expected With() to return a new logger instance")
				}

				// Verify child logger has the same requestAttrs reference
				childSlogger, ok := child.(*slogger)
				if !ok {
					t.Error("Expected With() to return *slogger")
				}

				if childSlogger.requestAttrs != l.requestAttrs {
					t.Error("Expected child logger to share requestAttrs with parent")
				}

				// Verify child logger has additional attributes
				if len(childSlogger.attrs) == 0 {
					t.Error("Expected child logger to have additional attributes")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, logger)
		})
	}
}

func TestMakeAttrs(t *testing.T) {
	tests := []struct {
		name     string
		args     []any
		expected int // expected number of attributes
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
			attrs := makeAttrs(tt.args...)

			if len(attrs) != tt.expected {
				t.Errorf("Expected %d attributes, got %d", tt.expected, len(attrs))
			}
		})
	}
}
