### Configuration

`flagd` supports configuration via config file, environment variables and flags. In cases of conflict, flags have the
highest priority, followed by environment variables and finally config file.

Supported flags are as follows (result of running `./flagd start --help`):

```
  -b, --bearer-token string       Set a bearer token to use for remote sync
  -e, --evaluator string          Set an evaluator e.g. json (default "json")
  -h, --help                      help for start
  -p, --port int32                Port to listen on (default 8013)
  -m, --metrics-port int32        Port to set up metrics listener on (default 8014)
  -c, --server-cert-path string   Server side tls certificate path
  -k, --server-key-path string    Server side tls key path
  -a, --sync-provider-args        Sync provider arguments as key values separated by =
  -d, --socket-path string        Set the flagd socket path.
  -f, --uri strings               Set a sync provider uri to read data from this can be a filepath, url or reference to a kubernetes custom resource. Using multiple providers is supported, however if flag keys are duplicated across multiple sources it may lead to unexpected behavior.
  -C, --cors-origin strings       Set a CORS allow origin header, setting "*" will allow all origins (by default CORS headers are not set)
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