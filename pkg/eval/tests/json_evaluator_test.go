package eval_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/open-feature/flagd/pkg/eval"
	"github.com/open-feature/flagd/pkg/model"
	"github.com/stretchr/testify/assert"
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
	}

	evaluator := eval.JSONEvaluator{}
	err := evaluator.SetState(Flags)
	if err != nil {
		t.Fatalf("Expected no error")
	}

	for _, test := range tests {
		apStruct, err := structpb.NewStruct(test.context)
		if err != nil {
			t.Fatal(err)
		}
		val, _, reason, err := evaluator.ResolveBooleanValue(test.flagKey, apStruct)
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
	}

	evaluator := eval.JSONEvaluator{}
	err := evaluator.SetState(Flags)
	if err != nil {
		t.Fatalf("Expected no error")
	}

	for _, test := range tests {
		apStruct, err := structpb.NewStruct(test.context)
		if err != nil {
			t.Fatal(err)
		}
		val, _, reason, err := evaluator.ResolveStringValue(test.flagKey, apStruct)

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
		flagKey   string
		context   map[string]interface{}
		val       float32
		reason    string
		errorCode string
	}{
		{StaticNumberFlag, nil, StaticNumberValue, model.StaticReason, ""},
		{DynamicNumberFlag, map[string]interface{}{ColorProp: ColorValue}, DynamicNumberValue, model.TargetingMatchReason, ""},
		{StaticObjectFlag, nil, 13, model.ErrorReason, model.TypeMismatchErrorCode},
		{MissingFlag, nil, 13, model.ErrorReason, model.FlagNotFoundErrorCode},
	}

	evaluator := eval.JSONEvaluator{}
	err := evaluator.SetState(Flags)
	if err != nil {
		t.Fatalf("Expected no error")
	}

	for _, test := range tests {
		apStruct, err := structpb.NewStruct(test.context)
		if err != nil {
			t.Fatal(err)
		}
		val, _, reason, err := evaluator.ResolveNumberValue(test.flagKey, apStruct)

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
	}

	evaluator := eval.JSONEvaluator{}
	err := evaluator.SetState(Flags)
	if err != nil {
		t.Fatalf("Expected no error")
	}

	for _, test := range tests {
		apStruct, err := structpb.NewStruct(test.context)
		if err != nil {
			t.Fatal(err)
		}
		val, _, reason, err := evaluator.ResolveObjectValue(test.flagKey, apStruct)

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
