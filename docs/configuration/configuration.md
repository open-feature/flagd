# Configuration

`flagd` supports configuration via config file, environment variables and flags. In cases of conflict, flags have the
highest priority, followed by environment variables and finally config file.

Supported flags are as follows (result of running `./flagd start --help`):

```
  -b, --bearer-token string                 Set a bearer token to use for remote sync
  -C, --cors-origin strings                 CORS allowed origins, * will allow all origins
  -h, --help                                help for start
  -m, --metrics-port int32                  Port to serve metrics on (default 8014)
  -p, --port int32                          Port to listen on (default 8013)
  -c, --server-cert-path string             Server side tls certificate path
  -k, --server-key-path string              Server side tls key path
  -d, --socket-path string                  Flagd socket path. With grpc the service will become available on this address. With http(s) the grpc-gateway proxy will use this address internally.
  -a, --sync-provider-args stringToString   Sync provider arguments as key values separated by = (default [])
  -f, --uri .yaml/.yml/.json                Set a sync provider uri to read data from, this can be a filepath,url or FeatureFlagConfiguration. Using multiple providers is supported however if flag keys are duplicated across multiple sources it may lead to unexpected behavior. Please note that if you are using filepath, flagd only supports files with .yaml/.yml/.json extension.
  -e, --evaluator string                    DEPRECATED: Set an evaluator e.g. json, yaml/yml. Please note that yaml/yml and json evaluations work the same (yaml/yml files are converted to json internally) (default "json")
  -y, --sync-provider string                DEPRECATED: Set a sync provider e.g. filepath or remote
```

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