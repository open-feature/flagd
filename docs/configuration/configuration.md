# Configuration

`flagd` supports configuration via config file, environment variables and flags. In cases of conflict, flags have the
highest priority, followed by environment variables and finally config file.

Supported flags are documented (auto-generated) [here](./flagd_start.md).

Environment variable keys are uppercased, prefixed with `FLAGD_` and all `-` are replaced with `_`. For example,
`sync-provider-args` in environment variable form is `FLAGD_SYNC_PROVIDER_ARGS`.

Config file expects the keys to have the exact naming as the flags.

### <a name="uri-patterns"></a> URI patterns

Any URI passed to flagd via the `--uri` flag must follow one of the 4 following patterns to ensure that it is passed to the correct implementation: 

| Sync       | Pattern                               | Example                               |
|------------|---------------------------------------|---------------------------------------|
| Kubernetes | `core.openfeature.dev/namespace/name` | `core.openfeature.dev/default/my-crd` |
| Filepath   | `file:path/to/my/flag`                | `file:etc/flagd/my-flags.json`        |
| Remote     | `http(s)://flag-source-url`           | `https://my-flags.com/flags`          |
| Grpc       | `grpc(s)://flag-source-url`           | `grpc://my-flags-server`              | 


### Customising sync providers

Custom sync providers can be used to provide flag evaluation logic.

#### Kubernetes provider 

The Kubernetes provider allows flagD to connect to a Kubernetes cluster and evaluate flags against a specified FeatureFlagConfiguration resource as defined within the [open-feature-operator](https://github.com/open-feature/open-feature-operator/blob/main/apis/core/v1alpha1/featureflagconfiguration_types.go) spec.

To use an existing FeatureFlagConfiguration custom resource, start flagD with the following command:

```shell
flagd start --uri core.openfeature.dev/default/my_example
```

### Source Configuration

While a URI may be passed to flagd via the `--uri` flag, some implementations may require further configurations. In these cases the `--sources` flag should be used.
The flag takes a string argument, which should be a JSON representation of an array of `SourceConfig` objects. Alternatively, these configurations should be passed to
flagd via config file, specified using the `--config` flag.

| Field       | Type                                                       | Note                                               | 
|-------------|------------------------------------------------------------|----------------------------------------------------|
| uri         | required `string`                                          |                                                    |
| provider    | required `string` (`file`, `kubernetes`, `http` or `grpc`) |                                                    |
| bearerToken | optional `string`                                          | Used for http sync                                 |
| certPath    | optional `string`                                          | Used for grpcs sync when TLS certificate is needed |

The `uri` field values do not need to follow the [URI patterns](#uri-patterns), the provider type is instead derived from the provider field. If the prefix is supplied, it will be removed on startup without error.

Example start command using a filepath sync provider and the equivalent config file definition:
```sh
./bin/flagd start --sources='[{"uri":"config/samples/example_flags.json","provider":"file"},{"uri":"http://my-flag-source.json","provider":"http","bearerToken":"bearer-dji34ld2l"}]{"uri":"default/my-flag-config","provider":"kubernetes"},{"uri":"grpc://my-flag-source:8080","provider":"grpc"}'
```

```yaml
sources:
- uri: config/samples/example_flags.json
  provider: file
- uri: http://my-flag-source.json
  provider: http
  bearerToken: bearer-dji34ld2l
- uri: default/my-flag-config
  provider: kubernetes
- uri: http://my-flag-source.json
  provider: kubernetes
- uri: grpc://my-flag-source:8080
  provider: grpc
- uri: grpcs://my-flag-source:8080
  provider: grpc
  certPath: /certs/ca.cert
```
