package eval

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnyValue(t *testing.T) {
	obj := AnyValue{
		Value:   "val",
		Variant: "variant",
		Reason:  "reason",
		FlagKey: "key",
		Error:   fmt.Errorf("err"),
	}

	require.Equal(t, obj, NewAnyValue("val", "variant", "reason", "key", fmt.Errorf("err")))
}
