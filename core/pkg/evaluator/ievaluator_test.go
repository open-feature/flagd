package evaluator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnyValue(t *testing.T) {
	obj := AnyValue{
		Value:    "val",
		Variant:  "variant",
		Reason:   "reason",
		FlagKey:  "key",
		Metadata: map[string]interface{}{},
		Error:    fmt.Errorf("err"),
	}

	require.Equal(t, obj, NewAnyValue("val", "variant", "reason", "key", map[string]interface{}{}, fmt.Errorf("err")))
}
