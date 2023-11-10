# Feature Flag Configuration

The `FeatureFlag` version `v1beta1` CRD defines a CR with the following example structure:

```yaml
apiVersion: core.openfeature.dev/v1alpha2
kind: FeatureFlag
metadata:
  name: featureflag-sample
spec:
  featureFlagSpec:
    flags:
      foo:
        state: "ENABLED"
        variants:
          bar: "BAR"
          baz: "BAZ"
        defaultVariant: "bar"
```

In the example above, we have defined a `String` type feature flag named `foo` and it is in the `ENABLED` state.
It has variants of `bar` and `baz`, referring to respected values of `BAR` and `BAZ`.
The default variant is set to`bar`.

## featureFlagSpec

The `featureFlagSpec` is an object representing the flag configurations themselves.
The documentation for this object can be found [here](https://github.com/open-feature/flagd/blob/main/docs/configuration/flag_configuration.md).
