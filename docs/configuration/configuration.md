# Configuration

`flagd` supports configuration via config file, environment variables and flags. In cases of conflict, flags have the
highest priority, followed by environment variables and finally config file.

Supported flags are documented (auto-generated) [here](./flagd_start.md).

Environment variable keys are uppercased, prefixed with `FLAGD_` and all `-` are replaced with `_`. For example,
`sync-provider-args` in environment variable form is `FLAGD_SYNC_PROVIDER_ARGS`.

Config file expects the keys to have the exact naming as the flags.

### URI patterns

Any URI passed to flagd via the `--uri` flag must follow one of the 3 following patterns to ensure that it is passed to the correct implementation: 

| Sync      | Pattern | Example |
| ----------- | ----------- | ----------- |
| Kubernetes      | `core.openfeature.dev/namespace/name`       | `core.openfeature.dev/default/my-crd`       |
| Filepath   | `file:path/to/my/flag`        | `file:etc/flagd/my-flags.json`       |
| Remote   | `http(s)://flag-source-url`        | `https://my-flags.com/flags`       |



### Customising sync providers

Custom sync providers can be used to provide flag evaluation logic.

#### Kubernetes provider 

The Kubernetes provider allows flagD to connect to a Kubernetes cluster and evaluate flags against a specified FeatureFlagConfiguration resource as defined within the [open-feature-operator](https://github.com/open-feature/open-feature-operator/blob/main/apis/core/v1alpha1/featureflagconfiguration_types.go) spec.

To use an existing FeatureFlagConfiguration custom resource, start flagD with the following command:

```shell
flagd start --uri core.openfeature.dev/default/my_example
```