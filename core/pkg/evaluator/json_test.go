//nolint:wrapcheck
package evaluator_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/open-feature/flagd/core/pkg/evaluator"
	"github.com/open-feature/flagd/core/pkg/logger"
	"github.com/open-feature/flagd/core/pkg/model"
	"github.com/open-feature/flagd/core/pkg/store"
	"github.com/open-feature/flagd/core/pkg/sync"
	"github.com/stretchr/testify/assert"
)

const InvalidFlags = `{
  "flags": {
    "invalidFlag": {
      "notState": "ENABLED",
      "notVariants": {
        "on": true,
        "off": false
      },
      "notDefaultVariant": "on"
    }
  }
}`

const ValidFlags = `{
  "flags": {
    "validFlag": {
      "state": "ENABLED",
      "variants": {
        "on": true,
        "off": false
      },
      "defaultVariant": "on"
    }
  }
}`

const NullDefault = `{
  "flags": {
    "validFlag": {
      "state": "ENABLED",
      "variants": {
        "on": true,
        "off": false
      },
      "defaultVariant": null
    }
  }
}`

const UndefinedDefault = `{
  "flags": {
    "validFlag": {
      "state": "ENABLED",
      "variants": {
        "on": true,
        "off": false
      }
    }
  }
}`

const NullDefaultWithTargetting = `{
  "flags": {
    "validFlag": {
      "state": "ENABLED",
      "variants": {
        "on": true,
        "off": false
      },
      "defaultVariant": null,
	  "targeting": {
        "if": [
          {
            "==": [
              {
                "var": [
                  "key"
                ]
              },
              "value"
            ]
          },
          "on"
        ]
      }
    }
  }
}`

const UndefinedDefaultWithTargetting = `{
  "flags": {
    "validFlag": {
      "state": "ENABLED",
      "variants": {
        "on": true,
        "off": false
      },
	  "targeting": {
        "if": [
          {
            "==": [
              {
                "var": [
                  "key"
                ]
              },
              "value"
            ]
          },
          "on"
        ]
      }
    }
  }
}`

const (
	FlagSetID                  = "testSetId"
	Version                    = "v33"
	ValidFlag                  = "validFlag"
	MissingFlag                = "missingFlag"
	StaticBoolFlag             = "staticBoolFlag"
	StaticBoolValue            = true
	StaticStringFlag           = "staticStringFlag"
	StaticStringValue          = "#CC0000"
	StaticFloatFlag            = "staticFloatFlag"
	StaticFloatValue   float64 = 1
	StaticIntFlag              = "staticIntFlag"
	StaticIntValue     int64   = 1
	StaticObjectFlag           = "staticObjectFlag"
	StaticObjectValue          = `{"abc": 123}`
	DynamicBoolFlag            = "targetingBoolFlag"
	DynamicBoolValue           = true
	DynamicStringFlag          = "targetingStringFlag"
	DynamicStringValue         = "my-string"
	DynamicFloatFlag           = "targetingFloatFlag"
	DynamicFloatValue  float64 = 100
	DynamicIntFlag             = "targetingNumberFlag"
	DynamicIntValue    int64   = 100
	DynamicObjectFlag          = "targetingObjectFlag"
	DynamicObjectValue         = `{ "key": true }`
	ColorProp                  = "color"
	ColorValue                 = "yellow"
	DisabledFlag               = "disabledFlag"
	MetadataFlag               = "metadataFlag"
	VersionOverride            = "v66"
)

var Flags = fmt.Sprintf(`{
	"metadata": {
		"flagSetId": "%s",
		"version": "%s"
	},
  "flags": {
    "%s": {
      "state": "ENABLED",
      "variants": {
        "on": %t,
        "off": false
      },
      "defaultVariant": "on"
    },
		"%s": {
      "state": "ENABLED",
      "variants": {
        "red": "%s",
        "blue": "#0000CC"
      },
      "defaultVariant": "red"
    },
		"%s": {
      "state": "ENABLED",
      "variants": {
        "one": %f,
        "two": 2
      },
      "defaultVariant": "one"
    },
	"%s": {
		"state": "ENABLED",
		"variants": {
		  "one": %d,
		  "two": 2
		},
		"defaultVariant": "one"
	  },
		"%s": {
      "state": "ENABLED",
      "variants": {
        "obj1": %s,
        "obj2": {
					"xyz": true
				}
      },
      "defaultVariant": "obj1"
    },
		"%s": {
      "state": "ENABLED",
      "variants": {
        "bool1": %t,
        "bool2": false
      },
      "defaultVariant": "bool2",
			"targeting": {
        "if": [
          {
            "==": [
              {
                "var": [
                  "%s"
                ]
              },
              "%s"
            ]
          },
          "bool1",
          null
        ]
      }
    },
		"%s": {
      "state": "ENABLED",
      "variants": {
        "str1": "%s",
        "str2": "other"
      },
      "defaultVariant": "str2",
			"targeting": {
        "if": [
          {
            "==": [
              {
                "var": [
                  "%s"
                ]
              },
              "%s"
            ]
          },
          "str1",
          null
        ]
      }
    },
		"%s": {
      "state": "ENABLED",
      "variants": {
        "number1": %f,
        "number2": 200
      },
      "defaultVariant": "number2",
			"targeting": {
        "if": [
          {
            "==": [
              {
                "var": [
                  "%s"
                ]
              },
              "%s"
            ]
          },
          "number1",
          null
        ]
      }
    },
	"%s": {
		"state": "ENABLED",
		"variants": {
		  "number1": %d,
		  "number2": 200
		},
		"defaultVariant": "number2",
			  "targeting": {
		  "if": [
			{
			  "==": [
				{
				  "var": [
					"%s"
				  ]
				},
				"%s"
			  ]
			},
			"number1",
			null
		  ]
		}
	  },
		"%s": {
      "state": "ENABLED",
      "variants": {
        "object1": %s,
        "object2": {}
      },
      "defaultVariant": "object2",
			"targeting": {
        "if": [
          {
            "==": [
              {
                "var": [
                  "%s"
                ]
              },
              "%s"
            ]
          },
          "object1",
          null
        ]
      }
    },
	"%s": {
      "state": "DISABLED",
      "variants": {
        "on": true,
        "off": false
      },
      "defaultVariant": "on"
    },
	"%s": {
      "state": "ENABLED",
      "variants": {
        "on": true,
        "off": false
      },
      "defaultVariant": "on",
			"metadata": {
				"version": "%s"
			}
    }
  }
}`,
	FlagSetID,
	Version,
	StaticBoolFlag,
	StaticBoolValue,
	StaticStringFlag,
	StaticStringValue,
	StaticFloatFlag,
	StaticFloatValue,
	StaticIntFlag,
	StaticIntValue,
	StaticObjectFlag,
	StaticObjectValue,
	DynamicBoolFlag,
	DynamicBoolValue,
	ColorProp,
	ColorValue,
	DynamicStringFlag,
	DynamicStringValue,
	ColorProp,
	ColorValue,
	DynamicFloatFlag,
	DynamicFloatValue,
	ColorProp,
	ColorValue,
	DynamicIntFlag,
	DynamicIntValue,
	ColorProp,
	ColorValue,
	DynamicObjectFlag,
	DynamicObjectValue,
	ColorProp,
	ColorValue,
	DisabledFlag,
	MetadataFlag,
	VersionOverride)

func TestGetState_Valid_ContainsFlag(t *testing.T) {
	evaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())
	_, _, err := evaluator.SetState(sync.DataSync{FlagData: ValidFlags})
	if err != nil {
		t.Fatalf("Expected no error")
	}

	// get the state
	state, err := evaluator.GetState()
	if err != nil {
		t.Fatalf("expected no error")
	}

	// validate it contains the flag
	wants := "validFlag"
	if !strings.Contains(state, wants) {
		t.Fatalf("expected: %s to contain: %s", state, wants)
	}
}

func TestSetState_Invalid_Error(t *testing.T) {
	evaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())

	// set state with an invalid flag definition
	_, _, err := evaluator.SetState(sync.DataSync{FlagData: InvalidFlags})
	if err != nil {
		t.Fatalf("unexpected error")
	}
}

func TestSetState_Valid_NoError(t *testing.T) {
	evaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())

	// set state with a valid flag definition
	_, _, err := evaluator.SetState(sync.DataSync{FlagData: ValidFlags})
	if err != nil {
		t.Fatalf("expected no error")
	}
}

func TestResolveAllValues(t *testing.T) {
	evaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())
	_, _, err := evaluator.SetState(sync.DataSync{FlagData: Flags})
	if err != nil {
		t.Fatalf("expected no error")
	}
	tests := []struct {
		context map[string]interface{}
	}{
		{
			context: map[string]interface{}{},
		},
		{
			context: map[string]interface{}{ColorProp: ColorValue},
		},
	}
	const reqID = "default"
	for _, test := range tests {
		vals, _, err := evaluator.ResolveAllValues(context.TODO(), reqID, test.context)
		if err != nil {
			t.Error("error from resolver", err)
		}

		for _, val := range vals {
			// disabled flag must be ignored from bulk evaluation
			if val.FlagKey == DisabledFlag {
				t.Errorf("disabled flag '%s' is present in evaluation results", DisabledFlag)
			}

			switch vT := val.Value.(type) {
			case bool:
				v, _, reason, _, _ := evaluator.ResolveBooleanValue(context.TODO(), reqID, val.FlagKey, test.context)
				assert.Equal(t, v, vT)
				assert.Equal(t, val.Reason, reason)
				assert.Equalf(t, val.Error, nil, "expected no errors, but got %v for flag key %s", val.Error, val.FlagKey)
			case string:
				v, _, reason, _, _ := evaluator.ResolveStringValue(context.TODO(), reqID, val.FlagKey, test.context)
				assert.Equal(t, v, vT)
				assert.Equal(t, val.Reason, reason)
				assert.Equalf(t, val.Error, nil, "expected no errors, but got %v for flag key %s", val.Error, val.FlagKey)
			case float64:
				v, _, reason, _, _ := evaluator.ResolveFloatValue(context.TODO(), reqID, val.FlagKey, test.context)
				assert.Equal(t, v, vT)
				assert.Equal(t, val.Reason, reason)
				assert.Equalf(t, val.Error, nil, "expected no errors, but got %v for flag key %s", val.Error, val.FlagKey)
			case interface{}:
				v, _, reason, _, _ := evaluator.ResolveObjectValue(context.TODO(), reqID, val.FlagKey, test.context)
				assert.Equal(t, v, vT)
				assert.Equal(t, val.Reason, reason)
				assert.Equalf(t, val.Error, nil, "expected no errors, but got %v for flag key %s", val.Error, val.FlagKey)
			}
		}
	}
}

func TestMetadataResolveType(t *testing.T) {
	tests := []struct {
		flagKey  string
		metadata model.Metadata
	}{
		{StaticBoolFlag, model.Metadata{"flagSetId": FlagSetID, "version": Version}},
		{MetadataFlag, model.Metadata{"flagSetId": FlagSetID, "version": VersionOverride}},
	}
	const reqID = "default"
	evaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())
	_, _, err := evaluator.SetState(sync.DataSync{FlagData: Flags})
	if err != nil {
		t.Fatalf("expected no error")
	}

	for _, test := range tests {
		_, _, _, metadata, _ := evaluator.ResolveBooleanValue(context.TODO(), reqID, test.flagKey, nil)
		if !reflect.DeepEqual(test.metadata, metadata) {
			t.Errorf("expected metadata to be %v, but got %v", test.metadata, metadata)
		}
	}
}

func TestMetadataResolveAll(t *testing.T) {
	expectedFlagSetMetadata := model.Metadata{"flagSetId": FlagSetID, "version": Version}

	tests := []struct {
		flagKey  string
		metadata model.Metadata
	}{
		{StaticBoolFlag, model.Metadata{"flagSetId": FlagSetID, "version": Version}},
		{MetadataFlag, model.Metadata{"flagSetId": FlagSetID, "version": VersionOverride}},
	}
	const reqID = "default"
	evaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())
	_, _, err := evaluator.SetState(sync.DataSync{FlagData: Flags})
	if err != nil {
		t.Fatalf("expected no error")
	}

	for _, test := range tests {
		resolutions, flagSetMetadata, _ := evaluator.ResolveAllValues(context.TODO(), reqID, nil)

		for _, resolved := range resolutions {
			if resolved.FlagKey == test.flagKey {
				if !reflect.DeepEqual(test.metadata, resolved.Metadata) {
					t.Errorf("expected flag metadata to be %v, but got %v", test.metadata, resolved.Metadata)
				}
			}
		}
		if !reflect.DeepEqual(expectedFlagSetMetadata, flagSetMetadata) {
			t.Errorf("expected flag set metadata to be %v, but got %v", expectedFlagSetMetadata, flagSetMetadata)
		}
	}
}

func TestResolveBooleanValue(t *testing.T) {
	tests := []struct {
		flagKey   string
		context   map[string]interface{}
		val       bool
		reason    string
		errorCode string
	}{
		{StaticBoolFlag, nil, StaticBoolValue, model.StaticReason, ""},
		{DynamicBoolFlag, map[string]interface{}{ColorProp: ColorValue}, StaticBoolValue, model.TargetingMatchReason, ""},
		{StaticObjectFlag, nil, StaticBoolValue, model.ErrorReason, model.TypeMismatchErrorCode},
		{MissingFlag, nil, StaticBoolValue, model.ErrorReason, model.FlagNotFoundErrorCode},
		{DisabledFlag, nil, StaticBoolValue, model.ErrorReason, model.FlagDisabledErrorCode},
	}
	const reqID = "default"
	evaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())
	_, _, err := evaluator.SetState(sync.DataSync{FlagData: Flags})
	if err != nil {
		t.Fatalf("expected no error")
	}

	for _, test := range tests {
		val, _, reason, _, err := evaluator.ResolveBooleanValue(context.TODO(), reqID, test.flagKey, test.context)
		if test.errorCode == "" {
			if assert.NoError(t, err) {
				assert.Equal(t, test.val, val)
				assert.Equal(t, test.reason, reason)
			}
		} else {
			assert.Equal(t, model.ErrorReason, reason)
			assert.EqualError(t, err, test.errorCode)
		}
	}
}

func BenchmarkResolveBooleanValue(b *testing.B) {
	tests := []struct {
		flagKey   string
		context   map[string]interface{}
		val       bool
		reason    string
		errorCode string
	}{
		{StaticBoolFlag, nil, StaticBoolValue, model.StaticReason, ""},
		{DynamicBoolFlag, map[string]interface{}{ColorProp: ColorValue}, StaticBoolValue, model.TargetingMatchReason, ""},
		{StaticObjectFlag, nil, StaticBoolValue, model.ErrorReason, model.TypeMismatchErrorCode},
		{MissingFlag, nil, StaticBoolValue, model.ErrorReason, model.FlagNotFoundErrorCode},
		{DisabledFlag, nil, StaticBoolValue, model.ErrorReason, model.FlagDisabledErrorCode},
	}

	evaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())
	_, _, err := evaluator.SetState(sync.DataSync{FlagData: Flags})
	if err != nil {
		b.Fatalf("expected no error")
	}
	reqID := "test"
	for _, test := range tests {
		b.Run(fmt.Sprintf("test %s", test.flagKey), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				val, _, reason, _, err := evaluator.ResolveBooleanValue(context.TODO(), reqID, test.flagKey, test.context)

				if test.errorCode == "" {
					if assert.NoError(b, err) {
						assert.Equal(b, test.val, val)
						assert.Equal(b, test.reason, reason)
					}
				} else {
					assert.Equal(b, model.ErrorReason, reason)
					assert.EqualError(b, err, test.errorCode)
				}
			}
		})
	}
}

func TestResolveStringValue(t *testing.T) {
	tests := []struct {
		flagKey   string
		context   map[string]interface{}
		val       string
		reason    string
		errorCode string
	}{
		{StaticStringFlag, nil, StaticStringValue, model.StaticReason, ""},
		{DynamicStringFlag, map[string]interface{}{ColorProp: ColorValue}, DynamicStringValue, model.TargetingMatchReason, ""},
		{StaticObjectFlag, nil, "", model.ErrorReason, model.TypeMismatchErrorCode},
		{MissingFlag, nil, "", model.ErrorReason, model.FlagNotFoundErrorCode},
		{DisabledFlag, nil, "", model.ErrorReason, model.FlagDisabledErrorCode},
	}
	const reqID = "default"
	evaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())
	_, _, err := evaluator.SetState(sync.DataSync{FlagData: Flags})
	if err != nil {
		t.Fatalf("expected no error")
	}

	for _, test := range tests {
		val, _, reason, _, err := evaluator.ResolveStringValue(context.TODO(), reqID, test.flagKey, test.context)

		if test.errorCode == "" {
			if assert.NoError(t, err) {
				assert.Equal(t, test.val, val)
				assert.Equal(t, test.reason, reason)
			}
		} else {
			assert.Equal(t, model.ErrorReason, reason)
			assert.EqualError(t, err, test.errorCode)
		}
	}
}

func BenchmarkResolveStringValue(b *testing.B) {
	tests := []struct {
		flagKey   string
		context   map[string]interface{}
		val       string
		reason    string
		errorCode string
	}{
		{StaticStringFlag, nil, StaticStringValue, model.StaticReason, ""},
		{DynamicStringFlag, map[string]interface{}{ColorProp: ColorValue}, DynamicStringValue, model.TargetingMatchReason, ""},
		{StaticObjectFlag, nil, "", model.ErrorReason, model.TypeMismatchErrorCode},
		{MissingFlag, nil, "", model.ErrorReason, model.FlagNotFoundErrorCode},
		{DisabledFlag, nil, "", model.ErrorReason, model.FlagDisabledErrorCode},
	}

	evaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())
	_, _, err := evaluator.SetState(sync.DataSync{FlagData: Flags})
	if err != nil {
		b.Fatalf("expected no error")
	}
	reqID := "test"
	for _, test := range tests {
		b.Run(fmt.Sprintf("test %s", test.flagKey), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				val, _, reason, _, err := evaluator.ResolveStringValue(context.TODO(), reqID, test.flagKey, test.context)

				if test.errorCode == "" {
					if assert.NoError(b, err) {
						assert.Equal(b, test.val, val)
						assert.Equal(b, test.reason, reason)
					}
				} else {
					assert.Equal(b, model.ErrorReason, reason)
					assert.EqualError(b, err, test.errorCode)
				}
			}
		})
	}
}

func TestResolveFloatValue(t *testing.T) {
	tests := []struct {
		flagKey   string
		context   map[string]interface{}
		val       float64
		reason    string
		errorCode string
	}{
		{StaticFloatFlag, nil, StaticFloatValue, model.StaticReason, ""},
		{DynamicFloatFlag, map[string]interface{}{ColorProp: ColorValue}, DynamicFloatValue, model.TargetingMatchReason, ""},
		{StaticObjectFlag, nil, 13, model.ErrorReason, model.TypeMismatchErrorCode},
		{MissingFlag, nil, 13, model.ErrorReason, model.FlagNotFoundErrorCode},
		{DisabledFlag, nil, 0, model.ErrorReason, model.FlagDisabledErrorCode},
	}
	const reqID = "default"
	evaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())
	_, _, err := evaluator.SetState(sync.DataSync{FlagData: Flags})
	if err != nil {
		t.Fatalf("expected no error")
	}

	for _, test := range tests {
		val, _, reason, _, err := evaluator.ResolveFloatValue(context.TODO(), reqID, test.flagKey, test.context)

		if test.errorCode == "" {
			if assert.NoError(t, err) {
				assert.Equal(t, test.val, val)
				assert.Equal(t, test.reason, reason)
			}
		} else {
			assert.Equal(t, model.ErrorReason, reason)
			assert.EqualError(t, err, test.errorCode)
		}
	}
}

func BenchmarkResolveFloatValue(b *testing.B) {
	tests := []struct {
		flagKey   string
		context   map[string]interface{}
		val       float64
		reason    string
		errorCode string
	}{
		{StaticFloatFlag, nil, StaticFloatValue, model.StaticReason, ""},
		{DynamicFloatFlag, map[string]interface{}{ColorProp: ColorValue}, DynamicFloatValue, model.TargetingMatchReason, ""},
		{StaticObjectFlag, nil, 13, model.ErrorReason, model.TypeMismatchErrorCode},
		{MissingFlag, nil, 13, model.ErrorReason, model.FlagNotFoundErrorCode},
		{DisabledFlag, nil, 0, model.ErrorReason, model.FlagDisabledErrorCode},
	}

	evaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())
	_, _, err := evaluator.SetState(sync.DataSync{FlagData: Flags})
	if err != nil {
		b.Fatalf("expected no error")
	}
	reqID := "test"
	for _, test := range tests {
		b.Run(fmt.Sprintf("test: %s", test.flagKey), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				val, _, reason, _, err := evaluator.ResolveFloatValue(context.TODO(), reqID, test.flagKey, test.context)

				if test.errorCode == "" {
					if assert.NoError(b, err) {
						assert.Equal(b, test.val, val)
						assert.Equal(b, test.reason, reason)
					}
				} else {
					assert.Equal(b, model.ErrorReason, reason)
					assert.EqualError(b, err, test.errorCode)
				}
			}
		})
	}
}

func TestResolveIntValue(t *testing.T) {
	tests := []struct {
		flagKey   string
		context   map[string]interface{}
		val       int64
		reason    string
		errorCode string
	}{
		{StaticIntFlag, nil, StaticIntValue, model.StaticReason, ""},
		{DynamicIntFlag, map[string]interface{}{ColorProp: ColorValue}, DynamicIntValue, model.TargetingMatchReason, ""},
		{StaticObjectFlag, nil, 13, model.ErrorReason, model.TypeMismatchErrorCode},
		{MissingFlag, nil, 13, model.ErrorReason, model.FlagNotFoundErrorCode},
		{DisabledFlag, nil, 0, model.ErrorReason, model.FlagDisabledErrorCode},
	}
	const reqID = "default"
	evaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())
	_, _, err := evaluator.SetState(sync.DataSync{FlagData: Flags})
	if err != nil {
		t.Fatalf("expected no error")
	}

	for _, test := range tests {
		val, _, reason, _, err := evaluator.ResolveIntValue(context.TODO(), reqID, test.flagKey, test.context)

		if test.errorCode == "" {
			if assert.NoError(t, err) {
				assert.Equal(t, test.val, val)
				assert.Equal(t, test.reason, reason)
			}
		} else {
			assert.Equal(t, model.ErrorReason, reason)
			assert.EqualError(t, err, test.errorCode)
		}
	}
}

func BenchmarkResolveIntValue(b *testing.B) {
	tests := []struct {
		flagKey   string
		context   map[string]interface{}
		val       int64
		reason    string
		errorCode string
	}{
		{StaticIntFlag, nil, StaticIntValue, model.StaticReason, ""},
		{DynamicIntFlag, map[string]interface{}{ColorProp: ColorValue}, DynamicIntValue, model.TargetingMatchReason, ""},
		{StaticObjectFlag, nil, 13, model.ErrorReason, model.TypeMismatchErrorCode},
		{MissingFlag, nil, 13, model.ErrorReason, model.FlagNotFoundErrorCode},
		{DisabledFlag, nil, 0, model.ErrorReason, model.FlagDisabledErrorCode},
	}

	evaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())
	_, _, err := evaluator.SetState(sync.DataSync{FlagData: Flags})
	if err != nil {
		b.Fatalf("expected no error")
	}
	reqID := "test"
	for _, test := range tests {
		b.Run(fmt.Sprintf("test %s", test.flagKey), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				val, _, reason, _, err := evaluator.ResolveIntValue(context.TODO(), reqID, test.flagKey, test.context)

				if test.errorCode == "" {
					if assert.NoError(b, err) {
						assert.Equal(b, test.val, val)
						assert.Equal(b, test.reason, reason)
					}
				} else {
					assert.Equal(b, model.ErrorReason, reason)
					assert.EqualError(b, err, test.errorCode)
				}
			}
		})
	}
}

func TestResolveObjectValue(t *testing.T) {
	tests := []struct {
		flagKey   string
		context   map[string]interface{}
		val       string
		reason    string
		errorCode string
	}{
		{StaticObjectFlag, nil, StaticObjectValue, model.StaticReason, ""},
		{DynamicObjectFlag, map[string]interface{}{ColorProp: ColorValue}, DynamicObjectValue, model.TargetingMatchReason, ""},
		{StaticBoolFlag, nil, "{}", model.ErrorReason, model.TypeMismatchErrorCode},
		{MissingFlag, nil, "{}", model.ErrorReason, model.FlagNotFoundErrorCode},
		{DisabledFlag, nil, "{}", model.ErrorReason, model.FlagDisabledErrorCode},
	}
	const reqID = "default"
	evaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())
	_, _, err := evaluator.SetState(sync.DataSync{FlagData: Flags})
	if err != nil {
		t.Fatalf("expected no error")
	}

	for _, test := range tests {
		val, _, reason, _, err := evaluator.ResolveObjectValue(context.TODO(), reqID, test.flagKey, test.context)

		if test.errorCode == "" {
			if assert.NoError(t, err) {
				marshalled, err := json.Marshal(val)
				if assert.NoError(t, err) {
					assert.JSONEq(t, test.val, string(marshalled))
					assert.Equal(t, test.reason, reason)
				}
			}
		} else {
			assert.Equal(t, model.ErrorReason, reason)
			assert.EqualError(t, err, test.errorCode)
		}
	}
}

func BenchmarkResolveObjectValue(b *testing.B) {
	tests := []struct {
		flagKey   string
		context   map[string]interface{}
		val       string
		reason    string
		errorCode string
	}{
		{StaticObjectFlag, nil, StaticObjectValue, model.StaticReason, ""},
		{DynamicObjectFlag, map[string]interface{}{ColorProp: ColorValue}, DynamicObjectValue, model.TargetingMatchReason, ""},
		{StaticBoolFlag, nil, "{}", model.ErrorReason, model.TypeMismatchErrorCode},
		{MissingFlag, nil, "{}", model.ErrorReason, model.FlagNotFoundErrorCode},
		{DisabledFlag, nil, "{}", model.ErrorReason, model.FlagDisabledErrorCode},
	}

	evaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())
	_, _, err := evaluator.SetState(sync.DataSync{FlagData: Flags})
	if err != nil {
		b.Fatalf("expected no error")
	}
	reqID := "test"
	for _, test := range tests {
		b.Run(fmt.Sprintf("test %s", test.flagKey), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				val, _, reason, _, err := evaluator.ResolveObjectValue(context.TODO(), reqID, test.flagKey, test.context)

				if test.errorCode == "" {
					if assert.NoError(b, err) {
						marshalled, err := json.Marshal(val)
						if assert.NoError(b, err) {
							assert.JSONEq(b, test.val, string(marshalled))
							assert.Equal(b, test.reason, reason)
						}
					}
				} else {
					assert.Equal(b, model.ErrorReason, reason)
					assert.EqualError(b, err, test.errorCode)
				}
			}
		})
	}
}

func TestResolveAsAnyValue(t *testing.T) {
	tests := []struct {
		flagKey   string
		context   map[string]interface{}
		val       string
		reason    string
		errorCode string
	}{
		// success
		{StaticBoolFlag, nil, "{}", model.StaticReason, ""},
		{StaticObjectFlag, nil, StaticObjectValue, model.StaticReason, ""},
		{DynamicObjectFlag, map[string]interface{}{ColorProp: ColorValue}, DynamicObjectValue, model.TargetingMatchReason, ""},
		// errors
		{MissingFlag, nil, "{}", model.ErrorReason, model.FlagNotFoundErrorCode},
		{DisabledFlag, nil, "{}", model.ErrorReason, model.FlagDisabledErrorCode},
	}

	evaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())
	_, _, err := evaluator.SetState(sync.DataSync{FlagData: Flags})
	if err != nil {
		t.Fatalf("expected no error")
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("evaluating flag: %s", test.flagKey), func(t *testing.T) {
			anyResult := evaluator.ResolveAsAnyValue(context.TODO(), "", test.flagKey, test.context)

			if test.errorCode == "" {
				assert.NoError(t, anyResult.Error)
			} else {
				assert.Equal(t, model.ErrorReason, anyResult.Reason)
				assert.EqualError(t, anyResult.Error, test.errorCode)
			}
		})
	}
}

func TestResolve_DefaultVariant(t *testing.T) {
	tests := []struct {
		flags     string
		flagKey   string
		context   map[string]interface{}
		reason    string
		errorCode string
	}{
		{NullDefault, ValidFlag, nil, model.ErrorReason, model.FlagNotFoundErrorCode},
		{UndefinedDefault, ValidFlag, nil, model.ErrorReason, model.FlagNotFoundErrorCode},
		{NullDefaultWithTargetting, ValidFlag, nil, model.ErrorReason, model.FlagNotFoundErrorCode},
		{UndefinedDefaultWithTargetting, ValidFlag, nil, model.ErrorReason, model.FlagNotFoundErrorCode},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			evaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())
			_, _, err := evaluator.SetState(sync.DataSync{FlagData: test.flags})

			if err != nil {
				t.Fatalf("expected no error")
			}

			anyResult := evaluator.ResolveAsAnyValue(context.TODO(), "", test.flagKey, test.context)

			assert.Equal(t, model.ErrorReason, anyResult.Reason)
			assert.EqualError(t, anyResult.Error, test.errorCode)
		})
	}
}

func TestSetState_DefaultVariantValidation(t *testing.T) {
	tests := map[string]struct {
		jsonFlags string
		valid     bool
	}{
		"is valid": {
			jsonFlags: `
				{
				  "flags": {
					"foo": {
					  "state": "ENABLED",
					  "variants": {
						"on": true,
						"off": false
					  },
					  "defaultVariant": "on"
					},
					"bar": {
					  "state": "ENABLED",
					  "variants": {
						"black": "#000000",
						"white": "#FFFFFF"
					  },
					  "defaultVariant": "black"
					}
				  }
			    }
			`,
			valid: true,
		},
		"is not valid": {
			jsonFlags: `
				{
				  "flags": {
					"foo": {
					  "state": "ENABLED",
					  "variants": {
						"black": "#000000",
						"white": "#FFFFFF"
					  },
					  "defaultVariant": "yellow"
					}
				  }
			    }
			`,
			valid: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			jsonEvaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())

			_, _, err := jsonEvaluator.SetState(sync.DataSync{FlagData: tt.jsonFlags})

			if tt.valid && err != nil {
				t.Error(err)
			}
		})
	}
}

func TestState_Evaluator(t *testing.T) {
	tests := map[string]struct {
		inputState          string
		expectedOutputState string
		expectedError       bool
		expectedResync      bool
	}{
		"success": {
			inputState: `
				{
  					"flags": {
						"fibAlgo": {
						  "variants": {
							"recursive": "recursive",
							"memo": "memo",
							"loop": "loop",
							"binet": "binet"
						  },
						  "defaultVariant": "recursive",
						  "state": "ENABLED",
						  "targeting": {
							"if": [
							  {
								"$ref": "emailWithFaas"
							  }, "binet", null
							]
						  }
    					}
					},
					"$evaluators": {
						"emailWithFaas": {
							  "in": ["@faas.com", {
								"var": ["email"]
							  }]
						}
  					}
				}
			`,
			expectedOutputState: `
				{
  					"flags": {
						"fibAlgo": {
						  "variants": {
							"recursive": "recursive",
							"memo": "memo",
							"loop": "loop",
							"binet": "binet"
						  },
						  "defaultVariant": "recursive",
						  "state": "ENABLED",
						  "source":"",
						  "selector":"",
						  "targeting": {
							"if": [
							  {
								"in": ["@faas.com", {
								"var": ["email"]
							  }]
							  }, "binet", null
							]
						  }
    					}
					},
					"flagSources":null
				}
			`,
		},
		"no-indentation": {
			inputState: `
				{
				"flags": {
				"fibAlgo": {
				"variants": {
				"recursive": "recursive",
				"memo": "memo",
				"loop": "loop",
				"binet": "binet"
				},
				"defaultVariant": "recursive",
				"state": "ENABLED",
				"targeting": {
				"if": [
				{
				"$ref": "emailWithFaas"
				}, "binet", null
				]
				}
				}
				},
				"$evaluators": {
				"emailWithFaas": {
				"in": ["@faas.com", {
				"var": ["email"]
				}]
				}
				}
				}
			`,
			expectedOutputState: `
				{
  					"flags": {
						"fibAlgo": {
						  "variants": {
							"recursive": "recursive",
							"memo": "memo",
							"loop": "loop",
							"binet": "binet"
						  },
						  "defaultVariant": "recursive",
						  "state": "ENABLED",
						  "source":"",
						  "selector":"",
						  "targeting": {
							"if": [
							  {
								"in": ["@faas.com", {
								"var": ["email"]
							  }]
							  }, "binet", null
							]
						  }
    					}
					},
					"flagSources":null
				}
			`,
		},
		"invalid evaluator json": {
			inputState: `
				{
  					"flags": {
						"fibAlgo": {
						  "variants": {
							"recursive": "recursive",
							"memo": "memo",
							"loop": "loop",
							"binet": "binet"
						  },
						  "defaultVariant": "recursive",
						  "state": "ENABLED",
						  "targeting": {
							"if": [
							  {
								"$ref": "emailWithFaas"
							  }, "binet", null
							]
						  }
    					}
					},
					"$evaluators": {
						"emailWithFaas": "foo"
  					}
				}
			`,
			expectedError: true,
		},
		"invalid targeting": {
			inputState: `
			{
					"flags": {
					"fibAlgo": {
						"variants": {
						"recursive": "recursive",
						"memo": "memo",
						"loop": "loop",
						"binet": "binet"
						},
						"defaultVariant": "recursive",
						"state": "ENABLED",
						"targeting": {
						"if": [
							{
							"in": ["@faas.com", {
							"var": ["email"]
							}]
							}, "binet", "null", "loop"
						]
						}
						},
					"isColorYellow": {
						"state": "ENABLED",
						"variants": {
							"on": true,
							"off": false
						},
						"defaultVariant": "off",
						"source":"",
						"targeting": {
							"if": [
								{
									"==": [
										{
											"varr": ["color"]
										},
										"yellow"
									]
								},
								"on",
								"off",
								"none"
							]
						}
					}
				},
				"flagSources":null
			}
		`,
			expectedError: false,
			expectedOutputState: `
			{
					"flags": {
					"fibAlgo": {
						"variants": {
						"recursive": "recursive",
						"memo": "memo",
						"loop": "loop",
						"binet": "binet"
						},
						"defaultVariant": "recursive",
						"state": "ENABLED",
						"source":"",
						"selector":"",
						"targeting": {
						"if": [
							{
							"in": ["@faas.com", {
							"var": ["email"]
							}]
							}, "binet", "null", "loop"
						]
						}
						},
					"isColorYellow": {
						"state": "ENABLED",
						"variants": {
							"on": true,
							"off": false
						},
						"defaultVariant": "off",
						"source":"",
						"selector":"",
						"targeting": {
							"if": [
								{
									"==": [
										{
											"varr": ["color"]
										},
										"yellow"
									]
								},
								"on",
								"off",
								"none"
							]
						}
					}
				},
				"flagSources":null
			}
		`,
		},
		"empty evaluator": {
			inputState: `
				{
  					"flags": {
						"fibAlgo": {
						  "variants": {
							"recursive": "recursive",
							"memo": "memo",
							"loop": "loop",
							"binet": "binet"
						  },
						  "defaultVariant": "recursive",
						  "state": "ENABLED",
						  "targeting": {
							"if": [
							  {
								"$ref": "emailWithFaas"
							  }, "binet", null
							]
						  }
    					}
					},
					"$evaluators": {
						"emailWithFaas": ""
  					}
				}
			`,
			expectedError: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			jsonEvaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())

			_, resync, err := jsonEvaluator.SetState(sync.DataSync{FlagData: tt.inputState})
			if err != nil {
				if !tt.expectedError {
					t.Error(err)
				}
				if resync != tt.expectedResync {
					t.Errorf("expected resync %t got %t", tt.expectedResync, resync)
				}
				return
			} else if tt.expectedError {
				t.Error("expected error, got nil")
				return
			}

			if resync != tt.expectedResync {
				t.Errorf("expected resync %t got %t", tt.expectedResync, resync)
			}

			got, err := jsonEvaluator.GetState()
			if err != nil {
				t.Error(err)
			}

			var expectedOutputJSON map[string]interface{}
			if err := json.Unmarshal([]byte(tt.expectedOutputState), &expectedOutputJSON); err != nil {
				t.Fatal(err)
			}

			var gotOutputJSON map[string]interface{}
			if err := json.Unmarshal([]byte(got), &gotOutputJSON); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(expectedOutputJSON["flags"], gotOutputJSON["flags"]) {
				t.Errorf("expected state: %v got state: %v", expectedOutputJSON, gotOutputJSON)
			}
		})
	}
}

func TestFlagStateSafeForConcurrentReadWrites(t *testing.T) {
	tests := map[string]struct {
		flagResolution func(evaluator *evaluator.JSON) error
	}{
		"Add_ResolveAllValues": {
			flagResolution: func(evaluator *evaluator.JSON) error {
				_, _, err := evaluator.ResolveAllValues(context.TODO(), "", nil)
				if err != nil {
					return err
				}
				return nil
			},
		},
		"Update_ResolveAllValues": {
			flagResolution: func(evaluator *evaluator.JSON) error {
				_, _, err := evaluator.ResolveAllValues(context.TODO(), "", nil)
				if err != nil {
					return err
				}
				return nil
			},
		},
		"Delete_ResolveAllValues": {
			flagResolution: func(evaluator *evaluator.JSON) error {
				_, _, err := evaluator.ResolveAllValues(context.TODO(), "", nil)
				if err != nil {
					return err
				}
				return nil
			},
		},
		"Add_ResolveBooleanValue": {
			flagResolution: func(evaluator *evaluator.JSON) error {
				_, _, _, _, err := evaluator.ResolveBooleanValue(context.TODO(), "", StaticBoolFlag, nil)
				return err
			},
		},
		"Update_ResolveStringValue": {
			flagResolution: func(evaluator *evaluator.JSON) error {
				_, _, _, _, err := evaluator.ResolveBooleanValue(context.TODO(), "", StaticStringValue, nil)
				return err
			},
		},
		"Delete_ResolveIntValue": {
			flagResolution: func(evaluator *evaluator.JSON) error {
				_, _, _, _, err := evaluator.ResolveIntValue(context.TODO(), "", StaticIntFlag, nil)
				return err
			},
		},
		"Add_ResolveFloatValue": {
			flagResolution: func(evaluator *evaluator.JSON) error {
				_, _, _, _, err := evaluator.ResolveFloatValue(context.TODO(), "", StaticFloatFlag, nil)
				return err
			},
		},
		"Update_ResolveObjectValue": {
			flagResolution: func(evaluator *evaluator.JSON) error {
				_, _, _, _, err := evaluator.ResolveObjectValue(context.TODO(), "", StaticObjectFlag, nil)
				return err
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			jsonEvaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())

			_, _, err := jsonEvaluator.SetState(sync.DataSync{FlagData: Flags})
			if err != nil {
				t.Fatal(err)
			}

			errChan := make(chan error)

			timeoutDur := 25 * time.Millisecond

			go func() {
				defer func() {
					if r := recover(); r != nil {
						errChan <- fmt.Errorf("%v", r)
					}
				}()
				timeout := time.After(timeoutDur)

				for {
					select {
					case <-timeout:
						errChan <- nil
						return
					default:
						_, _, err := jsonEvaluator.SetState(sync.DataSync{FlagData: Flags})
						if err != nil {
							errChan <- err
							return
						}
					}
				}
			}()

			go func() {
				defer func() {
					if r := recover(); r != nil {
						errChan <- fmt.Errorf("%v", r)
					}
				}()
				timeout := time.After(timeoutDur)

				for {
					select {
					case <-timeout:
						errChan <- nil
						return
					default:
						_ = tt.flagResolution(jsonEvaluator)
					}
				}
			}()

			for i := 0; i < 2; i++ {
				err := <-errChan
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestFlagdAmbientProperties(t *testing.T) {
	t.Run("flagKeyIsInTheContext", func(t *testing.T) {
		evaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())

		_, _, err := evaluator.SetState(sync.DataSync{FlagData: `{
			"flags": {
				"welcome-banner": {
					"state": "ENABLED",
					"variants": {
						"true": true,
						"false": false
					},
					"defaultVariant": "false",
					"targeting": {
						"==": [ { "var": "$flagd.flagKey" }, "welcome-banner" ]
					}
				}
			}
		}`})
		if err != nil {
			t.Fatal(err)
		}

		value, variant, reason, _, err := evaluator.ResolveBooleanValue(context.Background(), "default", "welcome-banner", nil)
		if err != nil {
			t.Fatal(err)
		}

		if !value {
			t.Fatal("expected true, got false")
		}

		if variant != "true" {
			t.Fatal("expected true, got false")
		}

		if reason != model.TargetingMatchReason {
			t.Fatalf("expected %s, got %s", model.TargetingMatchReason, reason)
		}
	})

	t.Run("timestampIsInTheContext", func(t *testing.T) {
		evaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())

		_, _, err := evaluator.SetState(sync.DataSync{FlagData: `{
			"flags": {
				"welcome-banner": {
					"state": "ENABLED",
					"variants": {
						"true": true,
						"false": false
					},
					"defaultVariant": "false",
					"targeting": {
						"<": [ 1696904426, { "var": "$flagd.timestamp" } ]
					}
				}
			}
		}`})
		if err != nil {
			t.Fatal(err)
		}

		value, variant, reason, _, err := evaluator.ResolveBooleanValue(context.Background(), "default", "welcome-banner", nil)
		if err != nil {
			t.Fatal(err)
		}

		if !value || variant != "true" || reason != model.TargetingMatchReason {
			t.Fatal("timestamp was not in the context")
		}
	})
}

func TestTargetingVariantBehavior(t *testing.T) {
	t.Run("missing variant error", func(t *testing.T) {
		evaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())

		_, _, err := evaluator.SetState(sync.DataSync{FlagData: `{
			"flags": {
				"missing-variant": {
					"state": "ENABLED",
					"variants": {
						"foo": true,
						"bar": false
					},
					"defaultVariant": "foo",
					"targeting": {
						"if": [ true, "buz", "baz"]
					}
				}
			}
		}`})
		if err != nil {
			t.Fatal(err)
		}

		_, _, _, _, err = evaluator.ResolveBooleanValue(context.Background(), "default", "missing-variant", nil)
		if err == nil {
			t.Fatal("missing variant did not result in error")
		}
	})

	t.Run("null fallback", func(t *testing.T) {
		evaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())

		_, _, err := evaluator.SetState(sync.DataSync{FlagData: `{
			"flags": {
				"null-fallback": {
					"state": "ENABLED",
					"variants": {
						"foo": true,
						"bar": false
					},
					"defaultVariant": "foo",
					"targeting": {
						"if": [ true, null, "baz"]
					}
				}
			}
		}`})
		if err != nil {
			t.Fatal(err)
		}

		value, variant, reason, _, err := evaluator.ResolveBooleanValue(context.Background(), "default", "null-fallback", nil)
		if err != nil {
			t.Fatal(err)
		}

		if !value || variant != "foo" || reason != model.DefaultReason {
			t.Fatal("did not fallback to defaultValue")
		}
	})

	t.Run("match booleans", func(t *testing.T) {
		evaluator := evaluator.NewJSON(logger.NewLogger(nil, false), store.NewFlags())

		//nolint:dupword
		_, _, err := evaluator.SetState(sync.DataSync{FlagData: `{
			"flags": {
				"match-boolean": {
					"state": "ENABLED",
					"variants": {
						"false": 1,
						"true": 2
					},
					"defaultVariant": "false",
					"targeting": {
						"if": [ true, true, false]
					}
				}
			}
		}`})
		if err != nil {
			t.Fatal(err)
		}

		value, variant, reason, _, err := evaluator.ResolveIntValue(context.Background(), "default", "match-boolean", nil)
		if err != nil {
			t.Fatal(err)
		}

		if value != 2 || variant != "true" || reason != model.TargetingMatchReason {
			t.Fatal("did not map to stringified boolean")
		}
	})
}
