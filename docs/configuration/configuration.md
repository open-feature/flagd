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

The config file expects the keys to have the exact naming as startup-flags flags.

## Sync providers

Sync providers are a core part of flagd; they are the abstraction that enables different sources for feature flag configurations.
flagd currently support the following sync providers:

* [Kubernetes provider](#kubernetes-provider)
* [Filepath Configuration](#filepath-provider)
* [Remote Configuration](#remote-provider)
* [GRPC Configuration](#grpc-provider)

### Kubernetes provider

The Kubernetes sync provider allows flagd to connect to a Kubernetes cluster and evaluate flags against a specified
FeatureFlagConfiguration resource as defined within
the [open-feature-operator](https://github.com/open-feature/open-feature-operator/blob/main/apis/core/v1alpha1/featureflagconfiguration_types.go)
spec.

To use an existing FeatureFlagConfiguration custom resource, start flagD with the following command:

```shell
flagd start --uri core.openfeature.dev/default/my_example
```

In this example, `default/my_example` expected to be a valid FeatureFlagConfiguration resource, where `default` is the
namespace and `my_example` being the resource name.

### Filepath provider

The file path sync provider reads and watch the source file for updates(ex:- changes and deletions).

```shell
flagd start --uri file:etc/featureflags.json
```

In this example, `etc/featureflags.json` is a valid feature flag configuration file accessible by the flagd runtime.

### Remote provider

The HTTP sync provider fetch flags from a remote source and periodically poll the source for flag configuration updates.

```shell
flagd start --uri https://my-flag-source.json
```

In this example, `https://my-flag-source.json` is a remote endpoint responding valid feature flag configurations when
invoked with **HTTP GET** request.

### GRPC provider

The GRPC sync provider stream flag configurations from a grpc sync provider implementation. This stream connection is ruled
by
the [sync service protobuf definition](https://github.com/open-feature/schemas/blob/main/protobuf/sync/v1/sync_service.proto).

```shell
flagd start --uri grpc://grpc-sync-source
```

In this example, `grpc-sync-source` is a grpc target implementing flagd protobuf definition.

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

Alternatively, these configurations can be passed to flagd via config file, specified using the `--config` flag.

| Field       | Type               | Note                                                                                                                                             |
|-------------|--------------------|--------------------------------------------------------------------------------------------------------------------------------------------------|
| uri         | required `string`  | Flag configuration source of the provider                                                                                                        |
| provider    | required `string`  | Provider type - `file`, `kubernetes`, `http` or `grpc`                                                                                           |
| bearerToken | optional `string`  | Used for http sync and token get appended to `Authorization` header with [bearer schema](https://www.rfc-editor.org/rfc/rfc6750#section-2.1)     |
| tls         | optional `boolean` | Enable/Disable secure TLS connectivity. Currently used only by GRPC sync. Default(ex:- if unset) is false, which will use an insecure connection |
| providerID  | optional `string`  | Value binds to grpc connection's providerID field. GRPC server implementations may use this to identify connecting flagd instance                |
| selector    | optional `string`  | Value binds to grpc connection's selector field. GRPC server implementations may use this to filter flag configurations                          |
| certPath    | optional `string`  | Used for grpcs sync when TLS certificate is needed. If not provided, system certificates will be used for TLS connection                         |

The `uri` field values **do not** follow the [URI patterns](#uri-patterns). The provider type is instead derived
from the `provider` field. Only exception is the remote provider where `http(s)://` is expected by default. Incorrect
URIs will result in a flagd start-up failure with errors from the respective sync provider implementation.

Given below are example sync providers, startup command and equivalent config file definition:

Sync providers,

* `file` - config/samples/example_flags.json
* `http` - <http://my-flag-source.json/>
* `kubernetes` - default/my-flag-config
* `grpc`(insecure) - grpc-source:8080
* `grpc`(secure) - my-flag-source:8080

Startup command,

```sh
./bin/flagd start 
--sources='[{"uri":"config/samples/example_flags.json","provider":"file"},
            {"uri":"http://my-flag-source.json","provider":"http","bearerToken":"bearer-dji34ld2l"},
            {"uri":"default/my-flag-config","provider":"kubernetes"},
            {"uri":"grpc-source:8080","provider":"grpc"},
            {"uri":"my-flag-source:8080","provider":"grpc", "certPath": "/certs/ca.cert", "tls": true, "providerID": "flagd-weatherapp-sidecar", "selector": "source=database,app=weatherapp"}]'
```

Configuration file,

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
    tls: true
    providerID: flagd-weatherapp-sidecar
    selector: 'source=database,app=weatherapp'
```
