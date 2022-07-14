package eval_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/open-feature/flagd/pkg/eval"
	gen "github.com/open-feature/flagd/pkg/generated"
	"github.com/open-feature/flagd/pkg/model"
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
	StaticNumberFlag           = "staticNumberFlag"
	StaticNumberValue  float32 = 1
	StaticObjectFlag           = "staticObjectFlag"
	StaticObjectValue          = `{"abc": 123}`
	DynamicBoolFlag            = "targetingBoolFlag"
	DynamicBoolValue           = true
	DynamicStringFlag          = "targetingStringFlag"
	DynamicStringValue         = "my-string"
	DynamicNumberFlag          = "targetingNumberFlag"
	DynamicNumberValue float32 = 100
	DynamicObjectFlag          = "targetingObjectFlag"
	DynamicObjectValue         = `{ "key": true }`
	ColorProp                  = "color"
	ColorValue                 = "yellow"
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
    }
  }
}`,
	StaticBoolFlag,
	StaticBoolValue,
	StaticStringFlag,
	StaticStringValue,
	StaticNumberFlag,
	StaticNumberValue,
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
	DynamicNumberFlag,
	DynamicNumberValue,
	ColorProp,
	ColorValue,
	DynamicObjectFlag,
	DynamicObjectValue,
	ColorProp,
	ColorValue)

func TestGetState_Valid_ContainsFlag(t *testing.T) {
	evaluator := eval.JSONEvaluator{}
	err := evaluator.SetState(ValidFlags)
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
	evaluator := eval.JSONEvaluator{}

	// set state with an invalid flag definition
	err := evaluator.SetState(InvalidFlags)
	if err == nil {
		t.Fatalf("Expected error")
	}
}

func TestSetState_Valid_NoError(t *testing.T) {
	evaluator := eval.JSONEvaluator{}

	// set state with a valid flag definition
	err := evaluator.SetState(ValidFlags)
	if err != nil {
		t.Fatalf("Expected no error")
	}
}

func TestResolveBooleanValue(t *testing.T) {
	tests := []struct {
		flagKey      string
		defaultValue bool
		context      gen.Context
		val          bool
		reason       string
		errorCode    string
	}{
		{StaticBoolFlag, false, gen.Context{}, StaticBoolValue, model.StaticReason, ""},
		{DynamicBoolFlag, false, gen.Context{AdditionalProperties: map[string]interface{}{
			ColorProp: ColorValue,
		}}, StaticBoolValue, model.TargetingMatchReason, ""},
		{StaticObjectFlag, false, gen.Context{}, StaticBoolValue, model.ErrorReason, model.TypeMismatchErrorCode},
		{MissingFlag, false, gen.Context{}, StaticBoolValue, model.ErrorReason, model.FlagNotFoundErrorCode},
	}

	evaluator := eval.JSONEvaluator{}
	err := evaluator.SetState(Flags)
	if err != nil {
		t.Fatalf("Expected no error")
	}

	for _, test := range tests {
		val, reason, err := evaluator.ResolveBooleanValue(test.flagKey, test.defaultValue, test.context)

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

func TestResolveStringValue(t *testing.T) {
	tests := []struct {
		flagKey      string
		defaultValue string
		context      gen.Context
		val          string
		reason       string
		errorCode    string
	}{
		{StaticStringFlag, "default", gen.Context{}, StaticStringValue, model.StaticReason, ""},
		{DynamicStringFlag, "default", gen.Context{AdditionalProperties: map[string]interface{}{
			ColorProp: ColorValue,
		}}, DynamicStringValue, model.TargetingMatchReason, ""},
		{StaticObjectFlag, "default", gen.Context{}, "", model.ErrorReason, model.TypeMismatchErrorCode},
		{MissingFlag, "default", gen.Context{}, "", model.ErrorReason, model.FlagNotFoundErrorCode},
	}

	evaluator := eval.JSONEvaluator{}
	err := evaluator.SetState(Flags)
	if err != nil {
		t.Fatalf("Expected no error")
	}

	for _, test := range tests {
		val, reason, err := evaluator.ResolveStringValue(test.flagKey, test.defaultValue, test.context)

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

func TestResolveNumberValue(t *testing.T) {
	tests := []struct {
		flagKey      string
		defaultValue float32
		context      gen.Context
		val          float32
		reason       string
		errorCode    string
	}{
		{StaticNumberFlag, 13, gen.Context{}, StaticNumberValue, model.StaticReason, ""},
		{DynamicNumberFlag, 13, gen.Context{AdditionalProperties: map[string]interface{}{
			ColorProp: ColorValue,
		}}, DynamicNumberValue, model.TargetingMatchReason, ""},
		{StaticObjectFlag, 13, gen.Context{}, 13, model.ErrorReason, model.TypeMismatchErrorCode},
		{MissingFlag, 13, gen.Context{}, 13, model.ErrorReason, model.FlagNotFoundErrorCode},
	}

	evaluator := eval.JSONEvaluator{}
	err := evaluator.SetState(Flags)
	if err != nil {
		t.Fatalf("Expected no error")
	}

	for _, test := range tests {
		val, reason, err := evaluator.ResolveNumberValue(test.flagKey, test.defaultValue, test.context)

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

func TestResolveObjectValue(t *testing.T) {
	tests := []struct {
		flagKey      string
		defaultValue map[string]interface{}
		context      gen.Context
		val          string
		reason       string
		errorCode    string
	}{
		{StaticObjectFlag, map[string]interface{}{}, gen.Context{}, StaticObjectValue, model.StaticReason, ""},
		{DynamicObjectFlag, map[string]interface{}{}, gen.Context{AdditionalProperties: map[string]interface{}{
			ColorProp: ColorValue,
		}}, DynamicObjectValue, model.TargetingMatchReason, ""},
		{StaticBoolFlag, map[string]interface{}{}, gen.Context{}, "{}", model.ErrorReason, model.TypeMismatchErrorCode},
		{MissingFlag, map[string]interface{}{}, gen.Context{}, "{}", model.ErrorReason, model.FlagNotFoundErrorCode},
	}

	evaluator := eval.JSONEvaluator{}
	err := evaluator.SetState(Flags)
	if err != nil {
		t.Fatalf("Expected no error")
	}

	for _, test := range tests {
		val, reason, err := evaluator.ResolveObjectValue(test.flagKey, test.defaultValue, test.context)

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

func TestMergeFlags(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		current eval.Flags
		new     eval.Flags
		want    eval.Flags
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
				"waka": {DefaultVariant: "off"},
			}},
			new: eval.Flags{Flags: map[string]eval.Flag{
				"paka": {DefaultVariant: "on"},
			}},
			want: eval.Flags{Flags: map[string]eval.Flag{
				"waka": {DefaultVariant: "off"},
				"paka": {DefaultVariant: "on"},
			}},
		},
		{
			name: "ignore override",
			current: eval.Flags{Flags: map[string]eval.Flag{
				"waka": {DefaultVariant: "off"},
			}},
			new: eval.Flags{Flags: map[string]eval.Flag{
				"waka": {DefaultVariant: "on"},
				"paka": {DefaultVariant: "on"},
			}},
			want: eval.Flags{Flags: map[string]eval.Flag{
				"waka": {DefaultVariant: "off"},
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
			got := tt.current.Merge(tt.new)
			require.Equal(t, got, tt.want)
		})
	}
}
