---
description: sync configuration overview for flagd and flagd providers
---

# Sync configuration

See [syncs](../concepts/syncs.md) for a conceptual overview.

## URI patterns

Any URI passed to flagd via the `--uri` (`-f`) flag must follow one of the 6 following patterns with prefixes to ensure that
it is passed to the correct implementation:

| Implied Sync Provider                 | Prefix                             | Example                               |
| ------------------------------------- | ---------------------------------- | ------------------------------------- |
| `kubernetes`                          | `core.openfeature.dev`             | `core.openfeature.dev/default/my-crd` |
| `file`                                | `file:`                            | `file:etc/flagd/my-flags.json`        |
| `http`                                | `http(s)://`                       | `https://my-flags.com/flags`          |
| `grpc`                                | `grpc(s)://`                       | `grpc://my-flags-server`              |
| &nbsp;[grpc](#custom-grpc-target-uri) | `[ envoy \| dns \| uds\| xds ]://` | `envoy://localhost:9211/test.service` |
| `gcs`                                 | `gs://`                            | `gs://my-bucket/my-flags.json`        |
| `azblob`                              | `azblob://`                        | `azblob://my-container/my-flags.json` |
| `s3`                                  | `s3://`                            | `s3://my-bucket/my-flags.json`        |

### Data Serialization

The `file`, `http`, `gcs`, `azblob` and `s3` sync providers expect the data to be formatted as JSON or YAML.
The file extension is used to determine the serialization format.
If the file extension hasn't been defined, the [media type](https://en.wikipedia.org/wiki/Media_type) will be used instead.

### Custom gRPC Target URI

Apart from default `dns` resolution, Flagd also support different resolution method e.g. `xds`.
Currently, we are supporting all [core resolver](https://grpc.io/docs/guides/custom-name-resolution/) and one custom resolver for `envoy` proxy resolution.
For more details, please refer the [RFC](https://github.com/open-feature/flagd/blob/main/docs/reference/specifications/proposal/rfc-grpc-custom-name-resolver.md) document.

```shell
./bin/flagd start -x --uri envoy://localhost:9211/test.service
```

## Source Configuration

While a URI may be passed to flagd via the `--uri` (`-f`) flag, some implementations may require further configurations.
In these cases the `--sources` flag should be used.

The flagd accepts a string argument, which should be a JSON representation of an array of `SourceConfig` objects.

Alternatively, these configurations can be passed to flagd via config file, specified using the `--config` flag.

| Field       | Type               | Note                                                                                                                                                                                                             |
| ----------- | ------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| uri         | required `string`  | Flag configuration source of the sync                                                                                                                                                                            |
| provider    | required `string`  | Provider type - `file`, `fsnotify`, `fileinfo`, `kubernetes`, `http`, `grpc`, `gcs` or `azblob`                                                                                                                  |
| authHeader  | optional `string`  | Used for http sync; set this to include the complete `Authorization` header value for any authentication scheme (e.g., "Bearer token_here", "Basic base64_credentials", etc.). Cannot be used with `bearerToken` |
| bearerToken | optional `string`  | (Deprecated) Used for http sync; token gets appended to `Authorization` header with [bearer schema](https://www.rfc-editor.org/rfc/rfc6750#section-2.1). Cannot be used with `authHeader`                        |
| interval    | optional `uint32`  | Used for http, gcs and azblob syncs; requests will be made at this interval. Defaults to 5 seconds.                                                                                                              |
| tls         | optional `boolean` | Enable/Disable secure TLS connectivity. Currently used only by gRPC sync. Default (ex: if unset) is false, which will use an insecure connection                                                                 |
| providerID  | optional `string`  | Value binds to grpc connection's providerID field. gRPC server implementations may use this to identify connecting flagd instance                                                                                |
| selector    | optional `string`  | Value binds to grpc connection's selector field. gRPC server implementations may use this to filter flag configurations                                                                                          |
| certPath    | optional `string`  | Used for grpcs sync when TLS certificate is needed. If not provided, system certificates will be used for TLS connection                                                                                         |
| maxMsgSize  | optional `int`     | Used for gRPC sync to set max receive message size (in bytes) e.g. 5242880 for 5MB. If not provided, the default is [4MB](https://pkg.go.dev/google.golang.org#grpc#MaxCallRecvMsgSize)                          |

The `uri` field values **do not** follow the [URI patterns](#uri-patterns). The provider type is instead derived
from the `provider` field. Only exception is the remote provider where `http(s)://` is expected by default. Incorrect
URIs will result in a flagd start-up failure with errors from the respective sync provider implementation.

The `file` provider type uses either an `fsnotify` notification (on systems that
support it), or a timer-based poller that relies on `os.Stat` and `fs.FileInfo`.
The moniker: `file` defaults to using `fsnotify` when flagd detects it is
running in kubernetes and `fileinfo` in all other cases, but you may explicitly
select either polling back-end by setting the provider value to either
`fsnotify` or `fileinfo`.

Given below are example sync providers, startup command and equivalent config file definition:

Sync providers:

- `file` - config/samples/example_flags.json
- `fsnotify` - config/samples/example_flags.json
- `fileinfo` - config/samples/example_flags.json
- [`http`](#http-configuration) - <http://my-flag-source.com/flags.json>
- `https` - <https://my-secure-flag-source.com/flags.json>
- `kubernetes` - default/my-flag-config
- `grpc`(insecure) - grpc-source:8080
- `grpcs`(secure) - my-flag-source:8080
- `grpc`(envoy) - envoy://localhost:9211/test.service
- `gcs` - gs://my-bucket/my-flags.json
- `azblob` - azblob://my-container/my-flags.json

Startup command:

```sh
./bin/flagd start
--sources='[{"uri":"config/samples/example_flags.json","provider":"file"},
            {"uri":"config/samples/example_flags.json","provider":"fsnotify"},
            {"uri":"config/samples/example_flags.json","provider":"fileinfo"},
            {"uri":"http://my-flag-source/flags.json","provider":"http","bearerToken":"bearer-dji34ld2l"},
            {"uri":"https://secure-remote/bearer-auth/flags.json","provider":"http","authHeader":"Bearer bearer-dji34ld2l"},
            {"uri":"https://secure-remote/basic-auth/flags.json","provider":"http","authHeader":"Basic dXNlcjpwYXNz"},
            {"uri":"default/my-flag-config","provider":"kubernetes"},
            {"uri":"grpc-source:8080","provider":"grpc"},
            {"uri":"my-flag-source:8080","provider":"grpc", "maxMsgSize": 5242880},
            {"uri":"envoy://localhost:9211/test.service", "provider":"grpc"},
            {"uri":"my-flag-source:8080","provider":"grpc", "certPath": "/certs/ca.cert", "tls": true, "providerID": "flagd-weatherapp-sidecar", "selector": "source=database,app=weatherapp"},
            {"uri":"gs://my-bucket/my-flag.json","provider":"gcs"},
            {"uri":"azblob://my-container/my-flag.json","provider":"azblob"}]'
```

Configuration file,

```yaml
sources:
  - uri: config/samples/example_flags.json
    provider: file
  - uri: config/samples/example_flags.json
    provider: fsnotify
  - uri: config/samples/example_flags.json
    provider: fileinfo
  - uri: http://my-flag-source/flags.json
    provider: http
    bearerToken: bearer-dji34ld2l
  - uri: default/my-flag-config
    provider: kubernetes
  - uri: my-flag-source:8080
    provider: grpc
  - uri: my-flag-source:8080
    provider: grpc
    maxMsgSize: 5242880
  - uri: envoy://localhost:9211/test.service
    provider: grpc
  - uri: my-flag-source:8080
    provider: grpc
    certPath: /certs/ca.cert
    tls: true
    providerID: flagd-weatherapp-sidecar
    selector: "source=database,app=weatherapp"
  - uri: gs://my-bucket/my-flag.json
    provider: gcs
  - uri: azblob://my-container/my-flags.json
    provider: azblob
```

### HTTP Configuration

The HTTP Configuration also supports OAuth that allows to securely fetch feature flag configurations from an HTTP endpoint
that requires OAuth-based authentication.

#### CLI-based OAuth Configuration

To enable OAuth, you need to update your Flagd configuration setting the `oauth` object which contains parameters to configure

....

#### File-based OAuth Configuration

the `clientID`, `clientSecret`, and the `tokenURL` for the OAuth Server.

```sh
./bin/flagd start
--sources='[{ 
  "uri": "http://localhost:8180/flags", 
  "provider": "http", 
  "interval": 1,
  "timeoutS": 10,
  "oauth": { 
    "clientID": "test", 
    "clientSecret": "test", 
    "tokenURL": "http://localhost:8180/sso/oauth2/token" 
  }}]'
```

Secrets can also be managed from the file system. This can be handy when, for example, deploying Flagd in Kubernetes. In this case, the client id and secret
will be read from the files `client-id` and `client-secret`, respectively. If the `folder` attribute is set, client id and secret on top level will be ignored.
To support rotating the secrets without restarting flagd, the additional parameter `ReloadDelayS` can be used to force
the reload of the secrets from the filesystem every `ReloadDelayS` seconds.

```sh
./bin/flagd start
--sources='[{ 
  "uri": "http://localhost:8180/flags", 
  "provider": "http", 
  "interval": 1,
  "timeoutS": 10,
  "oauth": { 
    "folder": "/etc/secrets", 
    "ReloadDelayS": 60, 
    "tokenURL": "http://localhost:8180/sso/oauth2/token" 
  }}]'
```
