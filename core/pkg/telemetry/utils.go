package telemetry

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

// utils contain common utilities to help with telemetry

const provider = "flagd"

// SemConvFeatureFlagAttributes is helper to derive semantic convention adhering feature flag attributes
// refer - https://opentelemetry.io/docs/reference/specification/trace/semantic_conventions/feature-flags/
func SemConvFeatureFlagAttributes(ffKey string, ffVariant string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.FeatureFlagKey(ffKey),
		semconv.FeatureFlagResultVariant(ffVariant),
		semconv.FeatureFlagProviderName(provider),
	}
}
