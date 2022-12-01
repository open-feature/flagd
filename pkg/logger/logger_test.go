package logger

import (
	"reflect"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
