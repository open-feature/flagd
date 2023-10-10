---
description: sync configuration overview for flagd and flagd providers
---

# Sync configuration

See [syncs](../concepts/syncs.md) for a conceptual overview.

## URI patterns

Any URI passed to flagd via the `--uri` (`-f`) flag must follow one of the 4 following patterns with prefixes to ensure that
it is passed to the correct implementation:

| Implied Sync Provider | Prefix                 | Example                               |
| --------------------- | ---------------------- | ------------------------------------- |
| `kubernetes`          | `core.openfeature.dev` | `core.openfeature.dev/default/my-crd` |
| `file`                | `file:`                | `file:etc/flagd/my-flags.json`        |
| `http`                | `http(s)://`           | `https://my-flags.com/flags`          |
| `grpc`                | `grpc(s)://`           | `grpc://my-flags-server`              |

## Source Configuration

While a URI may be passed to flagd via the `--uri` (`-f`) flag, some implementations may require further configurations.
In these cases the `--sources` flag should be used.

The flagd accepts a string argument, which should be a JSON representation of an array of `SourceConfig` objects.

Alternatively, these configurations can be passed to flagd via config file, specified using the `--config` flag.

| Field       | Type               | Note                                                                                                                                             |
| ----------- | ------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------ |
| uri         | required `string`  | Flag configuration source of the sync                                                                                                            |
| provider    | required `string`  | Provider type - `file`, `kubernetes`, `http` or `grpc`                                                                                           |
| bearerToken | optional `string`  | Used for http sync; token gets appended to `Authorization` header with [bearer schema](https://www.rfc-editor.org/rfc/rfc6750#section-2.1)       |
| interval    | optional `uint32`  | Used for http sync; requests will be made at this interval. Defaults to 5 seconds.                                                               |
| tls         | optional `boolean` | Enable/Disable secure TLS connectivity. Currently used only by gRPC sync. Default(ex:- if unset) is false, which will use an insecure connection |
| providerID  | optional `string`  | Value binds to grpc connection's providerID field. gRPC server implementations may use this to identify connecting flagd instance                |
| selector    | optional `string`  | Value binds to grpc connection's selector field. gRPC server implementations may use this to filter flag configurations                          |
| certPath    | optional `string`  | Used for grpcs sync when TLS certificate is needed. If not provided, system certificates will be used for TLS connection                         |

The `uri` field values **do not** follow the [URI patterns](#uri-patterns). The provider type is instead derived
from the `provider` field. Only exception is the remote provider where `http(s)://` is expected by default. Incorrect
URIs will result in a flagd start-up failure with errors from the respective sync provider implementation.

Given below are example sync providers, startup command and equivalent config file definition:

Sync providers:

- `file` - config/samples/example_flags.json
- `http` - <http://my-flag-source.json/>
- `https` - <https://my-secure-flag-source.json/>
- `kubernetes` - default/my-flag-config
- `grpc`(insecure) - grpc-source:8080
- `grpcs`(secure) - my-flag-source:8080

Startup command:

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
