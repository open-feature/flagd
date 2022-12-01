package eval_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/open-feature/flagd/pkg/eval"
	"github.com/open-feature/flagd/pkg/logger"
	"github.com/open-feature/flagd/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
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

const (
	MissingFlag                = "missingFlag"
	StaticBoolFlag             = "staticBoolFlag"
	StaticBoolValue            = true
	StaticStringFlag           = "staticStringFlag"
	StaticStringValue          = "#CC0000"
	StaticFloatFlag            = "staticFoatFlag"
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
)

var Flags = fmt.Sprintf(`{
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
    }
  }
}`,
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
	DisabledFlag)

func TestGetState_Valid_ContainsFlag(t *testing.T) {
	evaluator := eval.JSONEvaluator{Logger: logger.NewLogger(nil)}
	_, err := evaluator.SetState("", ValidFlags)
	if err != nil {
		t.Fatalf("Expected no error")
	}

	// get the state
	state, err := evaluator.GetState()
	if err != nil {
		t.Fatalf("Expected no error")
	}

	// validate it contains the flag
	wants := "validFlag"
	if !strings.Contains(state, wants) {
		t.Fatalf("Expected %s to contain %s", state, wants)
	}
}

func TestSetState_Invalid_Error(t *testing.T) {
	evaluator := eval.JSONEvaluator{Logger: logger.NewLogger(nil)}

	// set state with an invalid flag definition
	_, err := evaluator.SetState("", InvalidFlags)
	if err == nil {
		t.Fatalf("Expected error")
	}
}

func TestSetState_Valid_NoError(t *testing.T) {
	evaluator := eval.JSONEvaluator{Logger: logger.NewLogger(nil)}

	// set state with a valid flag definition
	_, err := evaluator.SetState("", ValidFlags)
	if err != nil {
		t.Fatalf("Expected no error")
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
	evaluator := eval.JSONEvaluator{Logger: logger.NewLogger(nil)}
	_, err := evaluator.SetState("", Flags)
	if err != nil {
		t.Fatalf("Expected no error")
	}

	for _, test := range tests {
		apStruct, err := structpb.NewStruct(test.context)
		if err != nil {
			t.Fatal(err)
		}
		val, _, reason, err := evaluator.ResolveBooleanValue(reqID, test.flagKey, apStruct)
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
		{StaticBoolFlag, nil, StaticBoolValue, model.DefaultReason, ""},
		{DynamicBoolFlag, map[string]interface{}{ColorProp: ColorValue}, StaticBoolValue, model.TargetingMatchReason, ""},
		{StaticObjectFlag, nil, StaticBoolValue, model.ErrorReason, model.TypeMismatchErrorCode},
		{MissingFlag, nil, StaticBoolValue, model.ErrorReason, model.FlagNotFoundErrorCode},
		{DisabledFlag, nil, StaticBoolValue, model.ErrorReason, model.FlagDisabledErrorCode},
	}

	evaluator := eval.JSONEvaluator{Logger: logger.NewLogger(nil)}
	_, err := evaluator.SetState("", Flags)
	if err != nil {
		b.Fatalf("Expected no error")
	}
	reqID := "test"
	for _, test := range tests {
		apStruct, err := structpb.NewStruct(test.context)
		if err != nil {
			b.Fatal(err)
		}
		b.Run(fmt.Sprintf("test %s", test.flagKey), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				val, _, reason, err := evaluator.ResolveBooleanValue(reqID, test.flagKey, apStruct)

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
	evaluator := eval.JSONEvaluator{Logger: logger.NewLogger(nil)}
	_, err := evaluator.SetState("", Flags)
	if err != nil {
		t.Fatalf("Expected no error")
	}

	for _, test := range tests {
		apStruct, err := structpb.NewStruct(test.context)
		if err != nil {
			t.Fatal(err)
		}
		val, _, reason, err := evaluator.ResolveStringValue(reqID, test.flagKey, apStruct)

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
		{StaticStringFlag, nil, StaticStringValue, model.DefaultReason, ""},
		{DynamicStringFlag, map[string]interface{}{ColorProp: ColorValue}, DynamicStringValue, model.TargetingMatchReason, ""},
		{StaticObjectFlag, nil, "", model.ErrorReason, model.TypeMismatchErrorCode},
		{MissingFlag, nil, "", model.ErrorReason, model.FlagNotFoundErrorCode},
		{DisabledFlag, nil, "", model.ErrorReason, model.FlagDisabledErrorCode},
	}

	evaluator := eval.JSONEvaluator{Logger: logger.NewLogger(nil)}
	_, err := evaluator.SetState("", Flags)
	if err != nil {
		b.Fatalf("Expected no error")
	}
	reqID := "test"
	for _, test := range tests {
		apStruct, err := structpb.NewStruct(test.context)
		if err != nil {
			b.Fatal(err)
		}
		b.Run(fmt.Sprintf("test %s", test.flagKey), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				val, _, reason, err := evaluator.ResolveStringValue(reqID, test.flagKey, apStruct)

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
	evaluator := eval.JSONEvaluator{Logger: logger.NewLogger(nil)}
	_, err := evaluator.SetState("", Flags)
	if err != nil {
		t.Fatalf("Expected no error")
	}

	for _, test := range tests {
		apStruct, err := structpb.NewStruct(test.context)
		if err != nil {
			t.Fatal(err)
		}
		val, _, reason, err := evaluator.ResolveFloatValue(reqID, test.flagKey, apStruct)

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
		{StaticFloatFlag, nil, StaticFloatValue, model.DefaultReason, ""},
		{DynamicFloatFlag, map[string]interface{}{ColorProp: ColorValue}, DynamicFloatValue, model.TargetingMatchReason, ""},
		{StaticObjectFlag, nil, 13, model.ErrorReason, model.TypeMismatchErrorCode},
		{MissingFlag, nil, 13, model.ErrorReason, model.FlagNotFoundErrorCode},
		{DisabledFlag, nil, 0, model.ErrorReason, model.FlagDisabledErrorCode},
	}

	evaluator := eval.JSONEvaluator{Logger: logger.NewLogger(nil)}
	_, err := evaluator.SetState("", Flags)
	if err != nil {
		b.Fatalf("Expected no error")
	}
	reqID := "test"
	for _, test := range tests {
		apStruct, err := structpb.NewStruct(test.context)
		if err != nil {
			b.Fatal(err)
		}
		b.Run(fmt.Sprintf("test %s", test.flagKey), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				val, _, reason, err := evaluator.ResolveFloatValue(reqID, test.flagKey, apStruct)

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
	evaluator := eval.JSONEvaluator{Logger: logger.NewLogger(nil)}
	_, err := evaluator.SetState("", Flags)
	if err != nil {
		t.Fatalf("Expected no error")
	}

	for _, test := range tests {
		apStruct, err := structpb.NewStruct(test.context)
		if err != nil {
			t.Fatal(err)
		}
		val, _, reason, err := evaluator.ResolveIntValue(reqID, test.flagKey, apStruct)

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
		{StaticIntFlag, nil, StaticIntValue, model.DefaultReason, ""},
		{DynamicIntFlag, map[string]interface{}{ColorProp: ColorValue}, DynamicIntValue, model.TargetingMatchReason, ""},
		{StaticObjectFlag, nil, 13, model.ErrorReason, model.TypeMismatchErrorCode},
		{MissingFlag, nil, 13, model.ErrorReason, model.FlagNotFoundErrorCode},
		{DisabledFlag, nil, 0, model.ErrorReason, model.FlagDisabledErrorCode},
	}

	evaluator := eval.JSONEvaluator{Logger: logger.NewLogger(nil)}
	_, err := evaluator.SetState("", Flags)
	if err != nil {
		b.Fatalf("Expected no error")
	}
	reqID := "test"
	for _, test := range tests {
		apStruct, err := structpb.NewStruct(test.context)
		if err != nil {
			b.Fatal(err)
		}
		b.Run(fmt.Sprintf("test %s", test.flagKey), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				val, _, reason, err := evaluator.ResolveIntValue(reqID, test.flagKey, apStruct)

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
	evaluator := eval.JSONEvaluator{Logger: logger.NewLogger(nil)}
	_, err := evaluator.SetState("", Flags)
	if err != nil {
		t.Fatalf("Expected no error")
	}

	for _, test := range tests {
		apStruct, err := structpb.NewStruct(test.context)
		if err != nil {
			t.Fatal(err)
		}
		val, _, reason, err := evaluator.ResolveObjectValue(reqID, test.flagKey, apStruct)

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
		{StaticObjectFlag, nil, StaticObjectValue, model.DefaultReason, ""},
		{DynamicObjectFlag, map[string]interface{}{ColorProp: ColorValue}, DynamicObjectValue, model.TargetingMatchReason, ""},
		{StaticBoolFlag, nil, "{}", model.ErrorReason, model.TypeMismatchErrorCode},
		{MissingFlag, nil, "{}", model.ErrorReason, model.FlagNotFoundErrorCode},
		{DisabledFlag, nil, "{}", model.ErrorReason, model.FlagDisabledErrorCode},
	}

	evaluator := eval.JSONEvaluator{Logger: logger.NewLogger(nil)}
	_, err := evaluator.SetState("", Flags)
	if err != nil {
		b.Fatalf("Expected no error")
	}
	reqID := "test"
	for _, test := range tests {
		apStruct, err := structpb.NewStruct(test.context)
		if err != nil {
			b.Fatal(err)
		}
		b.Run(fmt.Sprintf("test %s", test.flagKey), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				val, _, reason, err := evaluator.ResolveObjectValue(reqID, test.flagKey, apStruct)

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

func TestMergeFlags(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		current   eval.Flags
		new       eval.Flags
		newSource string
		want      eval.Flags
	}{
		{
			name:    "both nil",
			current: eval.Flags{Flags: nil},
			new:     eval.Flags{Flags: nil},
			want:    eval.Flags{Flags: map[string]eval.Flag{}},
		},
		{
			name:    "both empty flags",
			current: eval.Flags{Flags: map[string]eval.Flag{}},
			new:     eval.Flags{Flags: map[string]eval.Flag{}},
			want:    eval.Flags{Flags: map[string]eval.Flag{}},
		},
		{
			name:    "empty current",
			current: eval.Flags{Flags: nil},
			new:     eval.Flags{Flags: map[string]eval.Flag{}},
			want:    eval.Flags{Flags: map[string]eval.Flag{}},
		},
		{
			name:    "empty new",
			current: eval.Flags{Flags: map[string]eval.Flag{}},
			new:     eval.Flags{Flags: nil},
			want:    eval.Flags{Flags: map[string]eval.Flag{}},
		},
		{
			name: "extra fields on each",
			current: eval.Flags{Flags: map[string]eval.Flag{
				"waka": {
					DefaultVariant: "off",
					Source:         "1",
				},
			}},
			new: eval.Flags{Flags: map[string]eval.Flag{
				"paka": {
					DefaultVariant: "on",
				},
			}},
			newSource: "2",
			want: eval.Flags{Flags: map[string]eval.Flag{
				"waka": {
					DefaultVariant: "off",
					Source:         "1",
				},
				"paka": {
					DefaultVariant: "on",
					Source:         "2",
				},
			}},
		},
		{
			name: "override",
			current: eval.Flags{Flags: map[string]eval.Flag{
				"waka": {DefaultVariant: "off"},
			}},
			new: eval.Flags{Flags: map[string]eval.Flag{
				"waka": {DefaultVariant: "on"},
				"paka": {DefaultVariant: "on"},
			}},
			want: eval.Flags{Flags: map[string]eval.Flag{
				"waka": {DefaultVariant: "on"},
				"paka": {DefaultVariant: "on"},
			}},
		},
		{
			name: "identical",
			current: eval.Flags{Flags: map[string]eval.Flag{
				"hello": {DefaultVariant: "off"},
			}},
			new: eval.Flags{Flags: map[string]eval.Flag{
				"hello": {DefaultVariant: "off"},
			}},
			want: eval.Flags{Flags: map[string]eval.Flag{
				"hello": {DefaultVariant: "off"},
			}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, _ := tt.current.Merge(logger.NewLogger(nil), tt.newSource, tt.new)
			require.Equal(t, tt.want, got)
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
			jsonEvaluator := eval.JSONEvaluator{}

			_, err := jsonEvaluator.SetState("", tt.jsonFlags)

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
					}
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
			jsonEvaluator := eval.JSONEvaluator{}

			_, err := jsonEvaluator.SetState("", tt.inputState)
			if err != nil {
				if !tt.expectedError {
					t.Error(err)
				}
				return
			} else if tt.expectedError {
				t.Error("expected error, got nil")
				return
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

			if !reflect.DeepEqual(expectedOutputJSON, gotOutputJSON) {
				t.Errorf("expected state: %v got state: %v", expectedOutputJSON, gotOutputJSON)
			}
		})
	}
}
