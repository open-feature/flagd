# Configuration

<!-- TOC -->
* [Configuration](#configuration)
  * [Sync providers](#sync-providers)
    * [Kubernetes provider](#kubernetes-provider)
    * [Filepath provider](#filepath-provider)
    * [Remote provider](#remote-provider)
    * [GRPC provider](#grpc-provider)
  * [Sync provider configurations](#sync-provider-configurations)
    * [URI patterns](#uri-patterns)
    * [Source Configuration](#source-configuration)
<!-- TOC -->

`flagd` supports configuration via config file, environment variables and start-up flags. In cases of a conflict,
start-up flags have the highest priority, followed by environment variables and config file.

Supported start-up flags are documented (auto-generated) [here](./flagd_start.md).

Environment variable keys are uppercase, prefixed with `FLAGD_` and all `-` are replaced with `_`. For example,
`metrics-port` flag in environment variable form is `FLAGD_METRICS_PORT`.

Config file expects the keys to have the exact naming as startup-flags flags.

## Sync providers

Sync providers are a core part of flagd. They are the sources of feature flag configurations for evaluations. flagd
currently support the following sync providers,

* [Kubernetes provider](#kubernetes-provider)
* [Filepath Configuration](#filepath-provider)
* [Remote Configuration](#remote-provider)
* [GRPC Configuration](#grpc-provider)

### Kubernetes provider

The Kubernetes sync provider allows flagD to connect to a Kubernetes cluster and evaluate flags against a specified
FeatureFlagConfiguration resource as defined within
the [open-feature-operator](https://github.com/open-feature/open-feature-operator/blob/main/apis/core/v1alpha1/featureflagconfiguration_types.go)
spec.

To use an existing FeatureFlagConfiguration custom resource, start flagD with the following command:

```shell
flagd start --uri core.openfeature.dev/default/my_example
```

### Filepath provider

File path sync provider reads and watch the source file for updates(ex:- changes and deletions).

```shell
flagd start --uri file:etc/featureflags.json
```

### Remote provider

Remote sync provider fetch flags from a remote source and periodically poll the source for flag configuration updates.

```shell
flagd start --uri https://my-flag-source.json
```

### GRPC provider

GRPC sync provider stream flag configurations from a grpc sync provider implementation. This stream connection is ruled
by
the [sync service protobuf definition](https://github.com/open-feature/schemas/blob/main/protobuf/sync/v1/sync_service.proto).

```shell
flagd start --uri grpc://grpc-sync-source
```

There are two mechanisms to provide configurations of sync providers,

* [URI patterns](#uri-patterns)
* [Source Configuration](#source-configuration)

## Sync provider configurations

### URI patterns

Any URI passed to flagd via the `--uri` flag must follow one of the 4 following patterns with prefixes to ensure that
it is passed to the correct implementation:

| Sync       | Prefix                 | Example                               |
|------------|------------------------|---------------------------------------|
| Kubernetes | `core.openfeature.dev` | `core.openfeature.dev/default/my-crd` |
| Filepath   | `file:`                | `file:etc/flagd/my-flags.json`        |
| Remote     | `http(s)://`           | `https://my-flags.com/flags`          |
| Grpc       | `grpc(s)://`           | `grpc://my-flags-server`              |

### Source Configuration

While a URI may be passed to flagd via the `--uri` flag, some implementations may require further configurations.
In these cases the `--sources` flag should be used.

The flagd accepts a string argument, which should be a JSON representation of an array of `SourceConfig` objects.

Alternatively, these configurations should be passed to flagd via config file, specified using the `--config` flag.

| Field       | Type               | Note                                                                                                                                         |
|-------------|--------------------|----------------------------------------------------------------------------------------------------------------------------------------------|
| uri         | required `string`  | Flag configuration source of the provider                                                                                                    |
| provider    | required `string`  | Provider type - `file`, `kubernetes`, `http` or `grpc`                                                                                       |
| bearerToken | optional `string`  | Used for http sync and token get appended to `Authorization` header with [bearer schema](https://www.rfc-editor.org/rfc/rfc6750#section-2.1) |
| providerID  | optional `string`  | Value binds to grpc connection's providerID field. GRPC server implementations may use this to identify connecting flagd instance            |
| selector    | optional `string`  | Value binds to grpc connection's selector field. GRPC server implementations may use this to filter flag configurations                      |
| grpcSecure  | optional `boolean` | Used to enable secure TLS connectivity for grpc sync. Default(ex:- if unset) is false, which will use insecure grpc connection               |
| certPath    | optional `string`  | Used for grpcs sync when TLS certificate is needed. If not provided, system certificates will be used for TLS connection                     |

The `uri` field values **do not** follow the [URI patterns](#uri-patterns). The provider type is instead derived
from the `provider` field. Only exception is the remote provider where `http(s)://` is expected by default. Incorrect 
URIs will result in a flagd start-up failure with errors from the respective sync provider implementation.

Example start command using a filepath sync provider and the equivalent config file definition:

```sh
./bin/flagd start 
--sources='[{"uri":"config/samples/example_flags.json","provider":"file"},
            {"uri":"http://my-flag-source.json","provider":"http","bearerToken":"bearer-dji34ld2l"},
            {"uri":"default/my-flag-config","provider":"kubernetes"},
            {"uri":"grpc-source:8080","provider":"grpc"},
            {"uri":"my-flag-source:8080","provider":"grpc", "certPath": "/certs/ca.cert", "grpcSecure": "true", "providerID": "flagd-weatherapp-sidecar", "selector": "source=database,app=weatherapp"}]'
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
  - uri: my-flag-source:8080
    provider: grpc
  - uri: my-flag-source:8080
    provider: grpc
    certPath: /certs/ca.cert
    grpcSecure: true
    providerID: flagd-weatherapp-sidecar
    selector: 'source=database,app=weatherapp'
```
