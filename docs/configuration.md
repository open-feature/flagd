### Configuration

`flagd` supports configuration via config file, environment variables and flags. In cases of conflict, flags have the
highest priority, followed by environment variables and finally config file.

Supported flags are as follows (result of running `./flagd start --help`):

```
  -b, --bearer-token string       Set a bearer token to use for remote sync
  -e, --evaluator string          Set an evaluator e.g. json (default "json")
  -h, --help                      help for start
  -p, --port int32                Port to listen on (default 8013)
  -c, --server-cert-path string   Server side tls certificate path
  -k, --server-key-path string    Server side tls key path
  -s, --service-provider string   Set a service provider e.g. http or grpc (default "http")
  -a, --sync-provider-args stringToString   Sync provider arguments as key values separated by = (default [])
  -d, --socket-path string        Set the flagd socket path. With grpc the service will become available on this address. With http(s) the grpc-gateway proxy will use this address internally
  -y, --sync-provider string      Set a sync provider e.g. filepath or remote (default "filepath")
  -f, --uri strings               Set a sync provider uri to read data from this can be a filepath or url. Using multiple providers is supported where collisions between flags with the same key, the later will be used.
```

Environment variable keys are uppercased, prefixed with `FLAGD_` and all `-` are replaced with `_`. For example,
`sync-provider` in environment variable form is `FLAGD_SYNC_PROVIDER`.

Config file expects the keys to have the exact naming as the flags.


### Customising sync providers

Custom sync providers can be used to provide flag evaluation logic.
The Kubernetes provider allows flagD to connect to a Kubernetes cluster and evaluate flags against a specified FeatureFlagConfiguration resource as defined within the [open-feature-operator](https://github.com/open-feature/open-feature-operator/blob/main/apis/core/v1alpha1/featureflagconfiguration_types.go) spec.

To use an existing FeatureFlagConfiguration custom resource, start flagD with the following command:

```shell
flagd start --sync-provider=kubernetes --sync-provider-args=featureflagconfiguration=my-example --sync-provider-args=featureflagconfigurationnamespace=default
```