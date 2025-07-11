package telemetry

import (
	"testing"

	"github.com/stretchr/testify/require"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

func TestSemConvFeatureFlagAttributes(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		variant string
	}{
		{
			name:    "simple flag",
			key:     "flagA",
			variant: "bool",
		},
		{
			name: "empty variant flag",
			key:  "flagB",
		},
		{
			name: "empty key & variant does not panic",
		},
	}

	for _, test := range tests {
		attributes := SemConvFeatureFlagAttributes(test.key, test.variant)

		for _, attribute := range attributes {
			switch attribute.Key {
			case semconv.FeatureFlagKeyKey:
				require.Equal(t, test.key, attribute.Value.AsString(),
					"expected flag key: %s, but received: %s", test.key, attribute.Value.AsString())
			case semconv.FeatureFlagResultVariantKey:
				require.Equal(t, test.variant, attribute.Value.AsString(),
					"expected flag variant: %s, but received %s", test.variant, attribute.Value.AsString())
			case semconv.FeatureFlagProviderNameKey:
				require.Equal(t, provider, attribute.Value.AsString(),
					"expected flag provider: %s, but received %s", provider, attribute.Value.AsString())
			default:
				t.Errorf("attributes contains unexpected attribute. with key: %v, with type: %v",
					attribute.Key, attribute.Value.Type())
			}
		}
	}
}
