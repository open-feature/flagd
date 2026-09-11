---
description: flagd cheat sheet - quick reference for common operations
---

# Cheat sheet

This cheat sheet provides quick reference examples for running flagd and evaluating flags using various protocols and configurations.

Recommended tools:

- [docker](https://docs.docker.com/)
- [curl](https://curl.se/)
- [grpcurl](https://github.com/fullstorydev/grpcurl)
- [jq](https://jqlang.org/) (optional, for formatting)

!!! tip

    These commands assume a unix-like shell.

    Output is generally JSON, and can be pretty-printed by piping into `jq` (ie: `curl ... | jq`)

!!! note "Protocols and tabs"

    flagd exposes the same functionality over multiple protocols, so the examples below are tabbed.
    Not every operation supports every protocol. Streaming methods use native gRPC or Connect streaming over HTTP; OFREP exposes REST evaluation.

    - **HTTP (REST)** - plain `curl` against [OFREP](https://openfeature.dev/docs/reference/other-technologies/ofrep/) (evaluation) or `GET /v1/flags` (sync). The simplest path.
    - **HTTP (Connect)** - the gRPC service methods as HTTP/JSON via the [Connect protocol](https://connectrpc.com/docs/protocol/); plain `curl` with `Content-Type: application/json`, no gRPC client required.
    - **gRPC (grpcurl)** - the native gRPC services, using [grpcurl](https://github.com/fullstorydev/grpcurl) (needs the proto files, see [below](#proto-files-for-grpcurl)).

## Sample Flag Definitions

The examples below use these sample flag definition files. Download them to follow along:

- [cheat-sheet-flags.json](../assets/cheat-sheet-flags.json) - General application flags (flagSetId: `app-flags`)
- [cheat-sheet-flags-payments.json](../assets/cheat-sheet-flags-payments.json) - Payment-related flags (flagSetId: `payment-flags`)

The `app-flags` set includes:

| Flag Key              | Type    | Description                                    |
| --------------------- | ------- | ---------------------------------------------- |
| `simple-boolean`      | boolean | Static boolean flag                            |
| `simple-string`       | string  | Static string flag                             |
| `simple-number`       | integer | Static numeric flag                            |
| `simple-object`       | object  | Static object flag                             |
| `user-tier-flag`      | string  | Context-sensitive flag based on `tier`         |
| `email-based-feature` | boolean | Context-sensitive flag based on `email` domain |
| `region-config`       | object  | Context-sensitive flag based on `region`       |

The `payment-flags` set includes:

| Flag Key                 | Type    | Description                                   |
| ------------------------ | ------- | --------------------------------------------- |
| `payment-provider`       | string  | Static payment provider selection             |
| `max-transaction-amount` | integer | Context-sensitive based on `account-verified` |
| `enable-crypto-payments` | boolean | Context-sensitive based on `country`          |

---

## Running flagd

=== "Docker"

    ```shell
    # Single flag source (local file)
    docker run --rm -it \
      -p 8013:8013 \
      -p 8015:8015 \
      -p 8016:8016 \
      -v $(pwd):/flags \
      ghcr.io/open-feature/flagd:latest start \
      --uri file:./flags/cheat-sheet-flags.json
    ```

    ```shell
    # Multiple flag sources
    docker run --rm -it \
      -p 8013:8013 \
      -p 8015:8015 \
      -p 8016:8016 \
      -v $(pwd):/flags \
      ghcr.io/open-feature/flagd:latest start \
      --uri file:./flags/cheat-sheet-flags.json \
      --uri file:./flags/cheat-sheet-flags-payments.json
    ```

    ```shell
    # HTTP source
    docker run --rm -it \
      -p 8013:8013 \
      -p 8015:8015 \
      -p 8016:8016 \
      -v $(pwd):/flags \
      ghcr.io/open-feature/flagd:latest start \
      --uri https://flagd.dev/assets/cheat-sheet-flags.json
    ```

=== "Binary"

    ```shell
    # Single flag source (local file)
    flagd start --uri file:./cheat-sheet-flags.json
    ```

    ```shell
    # Multiple flag sources
    flagd start \
      --uri file:./cheat-sheet-flags.json \
      --uri file:./cheat-sheet-flags-payments.json
    ```

    ```shell
    # HTTP source
    flagd start --uri https://flagd.dev/assets/cheat-sheet-flags.json
    ```

!!! note
    The remaining examples use Docker, but all CLI flags work identically with the binary.

### Default ports

| Port | Protocol    | Service                                                                          |
| ---- | ----------- | -------------------------------------------------------------------------------- |
| 8013 | gRPC + HTTP | Flag evaluation (evaluation.proto, also served as HTTP/JSON via Connect)         |
| 8014 | HTTP        | Management (health checks, metrics)                                              |
| 8015 | gRPC + HTTP | Flag sync (sync.proto, also served as HTTP/JSON via Connect and `GET /v1/flags`) |
| 8016 | HTTP        | OFREP (OpenFeature Remote Evaluation Protocol)                                   |

### Proto files for grpcurl

flagd does not support gRPC reflection, so `grpcurl` needs the proto files to serialize requests and responses.
Clone the schemas repo (or download the protos from [buf.build/open-feature/flagd](https://buf.build/open-feature/flagd)) and point `$PROTO_DIR` at them:

```shell
git clone git@github.com:open-feature/flagd-schemas.git
PROTO_DIR="flagd-schemas/protobuf/"
```

---

## Evaluating flags

### Evaluate a single flag

=== "HTTP (REST)"

    OFREP, on port `8016`:

    ```shell
    curl -X POST 'http://localhost:8016/ofrep/v1/evaluate/flags/simple-boolean'
    ```

    Response:

    ```jsonc
    {
      "key": "simple-boolean",
      "reason": "STATIC",
      "variant": "on",
      "value": true,
      "metadata": {}
    }
    ```

=== "HTTP (Connect)"

    The evaluation service as HTTP/JSON, on port `8013`:

    ```shell
    curl -X POST 'http://localhost:8013/flagd.evaluation.v2.Service/ResolveBoolean' \
      -H 'Content-Type: application/json' \
      -d '{"flagKey": "simple-boolean", "context": {}}'
    ```

    Response:

    ```jsonc
    {
      "value": true,
      "reason": "STATIC",
      "variant": "on",
      "metadata": {}
    }
    ```

=== "gRPC (grpcurl)"

    ```shell
    grpcurl -plaintext \
      -import-path "$PROTO_DIR" -proto flagd/evaluation/v2/evaluation.proto \
      -d '{"flagKey": "simple-boolean", "context": {}}' \
      localhost:8013 \
      flagd.evaluation.v2.Service/ResolveBoolean
    ```

### Evaluate different flag types

=== "HTTP (REST)"

    ```shell
    # String flag
    curl -X POST 'http://localhost:8016/ofrep/v1/evaluate/flags/simple-string'

    # Number flag
    curl -X POST 'http://localhost:8016/ofrep/v1/evaluate/flags/simple-number'

    # Object flag
    curl -X POST 'http://localhost:8016/ofrep/v1/evaluate/flags/simple-object'
    ```

=== "HTTP (Connect)"

    Each type has its own method. Note `ResolveInt` encodes the value as a JSON string (`"50"`), per proto3 JSON:

    ```shell
    # String flag
    curl -X POST 'http://localhost:8013/flagd.evaluation.v2.Service/ResolveString' \
      -H 'Content-Type: application/json' \
      -d '{"flagKey": "simple-string", "context": {}}'

    # Number flag (ResolveInt or ResolveFloat)
    curl -X POST 'http://localhost:8013/flagd.evaluation.v2.Service/ResolveInt' \
      -H 'Content-Type: application/json' \
      -d '{"flagKey": "simple-number", "context": {}}'

    # Object flag
    curl -X POST 'http://localhost:8013/flagd.evaluation.v2.Service/ResolveObject' \
      -H 'Content-Type: application/json' \
      -d '{"flagKey": "simple-object", "context": {}}'
    ```

=== "gRPC (grpcurl)"

    ```shell
    # String flag
    grpcurl -plaintext \
      -import-path "$PROTO_DIR" -proto flagd/evaluation/v2/evaluation.proto \
      -d '{"flagKey": "simple-string", "context": {}}' \
      localhost:8013 \
      flagd.evaluation.v2.Service/ResolveString

    # Number flag (ResolveInt or ResolveFloat)
    grpcurl -plaintext \
      -import-path "$PROTO_DIR" -proto flagd/evaluation/v2/evaluation.proto \
      -d '{"flagKey": "simple-number", "context": {}}' \
      localhost:8013 \
      flagd.evaluation.v2.Service/ResolveInt

    # Object flag
    grpcurl -plaintext \
      -import-path "$PROTO_DIR" -proto flagd/evaluation/v2/evaluation.proto \
      -d '{"flagKey": "simple-object", "context": {}}' \
      localhost:8013 \
      flagd.evaluation.v2.Service/ResolveObject
    ```

### Evaluate all flags

=== "HTTP (REST)"

    OFREP bulk evaluation returns every flag in one response:

    ```shell
    curl -X POST 'http://localhost:8016/ofrep/v1/evaluate/flags'
    ```

    Response:

    ```jsonc
    {
      "flags": [
        {"key": "simple-boolean", "reason": "STATIC", "variant": "on", "value": true, "metadata": {}},
        {"key": "simple-string", "reason": "STATIC", "variant": "greeting", "value": "Hello, World!", "metadata": {}},
        {"key": "simple-number", "reason": "STATIC", "variant": "medium", "value": 50, "metadata": {}}
      ]
    }
    ```

=== "HTTP (Connect)"

    `ResolveAll` lives in the `v1` evaluation service:

    ```shell
    curl -X POST 'http://localhost:8013/flagd.evaluation.v1.Service/ResolveAll' \
      -H 'Content-Type: application/json' \
      -d '{"context": {}}'
    ```

=== "gRPC (grpcurl)"

    ```shell
    grpcurl -plaintext \
      -import-path "$PROTO_DIR" -proto flagd/evaluation/v1/evaluation.proto \
      -d '{"context": {}}' \
      localhost:8013 \
      flagd.evaluation.v1.Service/ResolveAll
    ```

---

## Context-Aware Evaluation

### Context from Request Body

Pass evaluation context in the request body to trigger targeting rules:

=== "HTTP (REST)"

    ```shell
    # Evaluate with email context (triggers email-based-feature targeting)
    curl -X POST 'http://localhost:8016/ofrep/v1/evaluate/flags/email-based-feature' \
      -H 'Content-Type: application/json' \
      -d '{"context": {"email": "user@example.com"}}'
    ```

    Response (email matches `@example.com`):

    ```jsonc
    {
      "key": "email-based-feature",
      "reason": "TARGETING_MATCH",
      "variant": "on",
      "value": true,
      "metadata": {}
    }
    ```

    ```shell
    # Evaluate with tier context
    curl -X POST 'http://localhost:8016/ofrep/v1/evaluate/flags/user-tier-flag' \
      -H 'Content-Type: application/json' \
      -d '{"context": {"tier": "premium"}}'
    ```

    ```shell
    # Bulk evaluation with multiple context values
    curl -X POST 'http://localhost:8016/ofrep/v1/evaluate/flags' \
      -H 'Content-Type: application/json' \
      -d '{"context": {"email": "admin@example.com", "tier": "enterprise", "region": "us"}}'
    ```

=== "HTTP (Connect)"

    ```shell
    # Evaluate with email context (triggers email-based-feature targeting)
    curl -X POST 'http://localhost:8013/flagd.evaluation.v2.Service/ResolveBoolean' \
      -H 'Content-Type: application/json' \
      -d '{"flagKey": "email-based-feature", "context": {"email": "user@example.com"}}'

    # Evaluate with tier context
    curl -X POST 'http://localhost:8013/flagd.evaluation.v2.Service/ResolveString' \
      -H 'Content-Type: application/json' \
      -d '{"flagKey": "user-tier-flag", "context": {"tier": "premium"}}'
    ```

=== "gRPC (grpcurl)"

    ```shell
    # Evaluate with email context (triggers email-based-feature targeting)
    grpcurl -plaintext \
      -import-path "$PROTO_DIR" -proto flagd/evaluation/v2/evaluation.proto \
      -d '{"flagKey": "email-based-feature", "context": {"email": "user@example.com"}}' \
      localhost:8013 \
      flagd.evaluation.v2.Service/ResolveBoolean

    # Evaluate with tier context
    grpcurl -plaintext \
      -import-path "$PROTO_DIR" -proto flagd/evaluation/v2/evaluation.proto \
      -d '{"flagKey": "user-tier-flag", "context": {"tier": "premium"}}' \
      localhost:8013 \
      flagd.evaluation.v2.Service/ResolveString
    ```

### Context from Static Values

Add static context values using the `-X` flag at startup. These are automatically included in all evaluations:

```shell
docker run --rm -it \
  -p 8013:8013 -p 8015:8015 -p 8016:8016 \
  -v $(pwd):/flags \
  ghcr.io/open-feature/flagd:latest start \
  --uri file:./flags/cheat-sheet-flags.json \
  -X region=eu \
  -X environment=production
```

```shell
# region=eu and environment=production is automatically applied without needing to send context
curl -X POST 'http://localhost:8016/ofrep/v1/evaluate/flags/region-config'
```

### Context from HTTP Headers

Map HTTP headers to evaluation context keys using the `-H` flag at startup:

```shell
docker run --rm -it \
  -p 8013:8013 -p 8015:8015 -p 8016:8016 \
  -v $(pwd):/flags \
  ghcr.io/open-feature/flagd:latest start \
  --uri file:./flags/cheat-sheet-flags.json \
  -H "X-User-Tier=tier" \
  -H "X-User-Email=email"
```

Now context is extracted from request headers:

```shell
# tier context comes from X-User-Tier header
curl -X POST 'http://localhost:8016/ofrep/v1/evaluate/flags/user-tier-flag' \
  -H 'X-User-Tier: enterprise'

# email context comes from X-User-Email header
curl -X POST 'http://localhost:8016/ofrep/v1/evaluate/flags/email-based-feature' \
  -H 'X-User-Email: developer@example.com'
```

### Context Priority

When using multiple context sources, values are merged with this priority (highest to lowest):

1. Header-mapped context values (`-H` flag)
2. Static context values (`-X` flag)
3. Request body context

---

## Using the Selector Header

When using multiple flag sources, the `Flagd-Selector` header restricts which flags are evaluated.

Start flagd with multiple sources:

```shell
docker run --rm -it \
  -p 8013:8013 -p 8015:8015 -p 8016:8016 \
  -v $(pwd):/flags \
  ghcr.io/open-feature/flagd:latest start \
  --uri file:./flags/cheat-sheet-flags.json \
  --uri file:./flags/cheat-sheet-flags-payments.json
```

Filter evaluations by flag set (`flagSetId`):

=== "HTTP (REST)"

    ```shell
    # Evaluate only flags from the app flag set
    curl -X POST 'http://localhost:8016/ofrep/v1/evaluate/flags' \
      -H 'Flagd-Selector: flagSetId=app-flags'

    # Evaluate only flags from the payments flag set
    curl -X POST 'http://localhost:8016/ofrep/v1/evaluate/flags' \
      -H 'Flagd-Selector: flagSetId=payment-flags'

    # Single flag evaluation with selector
    curl -X POST 'http://localhost:8016/ofrep/v1/evaluate/flags/payment-provider' \
      -H 'Flagd-Selector: flagSetId=payment-flags'
    ```

=== "HTTP (Connect)"

    ```shell
    # Single flag evaluation with selector
    curl -X POST 'http://localhost:8013/flagd.evaluation.v2.Service/ResolveBoolean' \
      -H 'Content-Type: application/json' \
      -H 'Flagd-Selector: flagSetId=app-flags' \
      -d '{"flagKey": "simple-boolean", "context": {}}'

    # ResolveAll with selector
    curl -X POST 'http://localhost:8013/flagd.evaluation.v1.Service/ResolveAll' \
      -H 'Content-Type: application/json' \
      -H 'Flagd-Selector: flagSetId=payment-flags' \
      -d '{"context": {}}'
    ```

=== "gRPC (grpcurl)"

    ```shell
    # Single flag evaluation with selector
    grpcurl -plaintext \
      -import-path "$PROTO_DIR" -proto flagd/evaluation/v2/evaluation.proto \
      -H 'Flagd-Selector: flagSetId=app-flags' \
      -d '{"flagKey": "simple-boolean", "context": {}}' \
      localhost:8013 \
      flagd.evaluation.v2.Service/ResolveBoolean

    # ResolveAll with selector
    grpcurl -plaintext \
      -import-path "$PROTO_DIR" -proto flagd/evaluation/v1/evaluation.proto \
      -H 'Flagd-Selector: flagSetId=payment-flags' \
      -d '{"context": {}}' \
      localhost:8013 \
      flagd.evaluation.v1.Service/ResolveAll
    ```

---

## Event Stream

The `EventStream` is a server stream that pushes `configuration_change` events when flags change (used by RPC-mode providers for cache invalidation).
It is a streaming method; over HTTP it is available as Connect streaming, but `grpcurl` is the practical CLI:

```shell
grpcurl -plaintext \
  -import-path "$PROTO_DIR" -proto flagd/evaluation/v1/evaluation.proto \
  -d '{}' \
  localhost:8013 \
  flagd.evaluation.v1.Service/EventStream
```

The `Flagd-Selector` header also filters the stream, so it only reports `configuration_change` events for flags in the selected flag set(s):

```shell
grpcurl -plaintext \
  -import-path "$PROTO_DIR" -proto flagd/evaluation/v1/evaluation.proto \
  -H 'Flagd-Selector: flagSetId=payment-flags' \
  -d '{}' \
  localhost:8013 \
  flagd.evaluation.v1.Service/EventStream
```

For the `EventStream`, the selector can only be supplied via the header; its request body has no selector field.

---

## Syncing Flag Configuration

The sync service (port `8015`) is used by in-process providers to fetch and stream flag configurations.
It serves the gRPC `sync.proto` service and, on the same port, an HTTP `GET /v1/flags` endpoint (the unary equivalent of `FetchAllFlags`).
Disable the HTTP endpoint with `--sync-http-enabled=false`.

### Fetch all flags

=== "HTTP (REST)"

    `GET /v1/flags` returns the flag configuration document itself (the same string `FetchAllFlags` returns in its `flag_configuration` field):

    ```shell
    curl http://localhost:8015/v1/flags
    ```

    Response:

    ```jsonc
    {
      "flags": {
        "simple-boolean": {
          "state": "ENABLED",
          "defaultVariant": "on",
          "variants": {"on": true, "off": false}
        }
        // ...
      }
    }
    ```

=== "HTTP (Connect)"

    `FetchAllFlags` as HTTP/JSON returns the document wrapped in an envelope:

    ```shell
    curl -X POST 'http://localhost:8015/flagd.sync.v1.FlagSyncService/FetchAllFlags' \
      -H 'Content-Type: application/json' \
      -d '{}'
    ```

    Response:

    ```jsonc
    {
      "flagConfiguration": "{\"flags\":{\"simple-boolean\":{...}}}"
    }
    ```

=== "gRPC (grpcurl)"

    ```shell
    grpcurl -plaintext \
      -import-path "$PROTO_DIR" -proto flagd/sync/v1/sync.proto \
      -d '{}' \
      localhost:8015 \
      flagd.sync.v1.FlagSyncService/FetchAllFlags
    ```

### Fetch all flags with a selector

Filter which flag source's configuration is returned.

=== "HTTP (REST)"

    The selector is supplied with the `Flagd-Selector` header:

    ```shell
    # Fetch only the payment flags
    curl -H 'Flagd-Selector: flagSetId=payment-flags' http://localhost:8015/v1/flags
    ```

    The endpoint distinguishes two outcomes:

    | Result                                           | Example                       | Status                    |
    | ------------------------------------------------ | ----------------------------- | ------------------------- |
    | Selector is malformed or names an unknown filter | control characters, `bogus=1` | `400`                     |
    | Valid filter that currently matches no flags     | `flagSetId=empty-set`         | `200` with `{"flags":{}}` |

    An empty result is deliberately not an error: a flag set holding no flags is a normal state, and a downstream flagd syncing from this endpoint should not break when it happens.

=== "HTTP (Connect)"

    The selector may be supplied in the request body or with the `Flagd-Selector` header; if both are present, the header wins:

    ```shell
    curl -X POST 'http://localhost:8015/flagd.sync.v1.FlagSyncService/FetchAllFlags' \
      -H 'Content-Type: application/json' \
      -d '{"selector": "flagSetId=payment-flags"}'
    ```

=== "gRPC (grpcurl)"

    ```shell
    # Fetch only the app flags
    grpcurl -plaintext \
      -import-path "$PROTO_DIR" -proto flagd/sync/v1/sync.proto \
      -d '{"selector": "flagSetId=app-flags"}' \
      localhost:8015 \
      flagd.sync.v1.FlagSyncService/FetchAllFlags

    # With provider ID for identification
    grpcurl -plaintext \
      -import-path "$PROTO_DIR" -proto flagd/sync/v1/sync.proto \
      -d '{"providerId": "my-app-sidecar", "selector": "flagSetId=app-flags"}' \
      localhost:8015 \
      flagd.sync.v1.FlagSyncService/FetchAllFlags
    ```

    The `Flagd-Selector` header works as an alternative to the request body `selector` field; if both are supplied, the header takes precedence.

### Streaming sync

`SyncFlags` establishes a server-streaming connection that pushes the initial configuration and then streams updates whenever flags change.
Like `EventStream` it is a streaming method (Connect streaming over HTTP, or native gRPC); `grpcurl` is the practical CLI:

```shell
grpcurl -plaintext \
  -import-path "$PROTO_DIR" -proto flagd/sync/v1/sync.proto \
  -d '{}' \
  localhost:8015 \
  flagd.sync.v1.FlagSyncService/SyncFlags
```

```shell
# Stream the app flags (initial config, then updates as they change)
grpcurl -plaintext \
  -import-path "$PROTO_DIR" -proto flagd/sync/v1/sync.proto \
  -d '{"selector": "flagSetId=app-flags"}' \
  localhost:8015 \
  flagd.sync.v1.FlagSyncService/SyncFlags
```

### HTTP caching

`GET /v1/flags` responses carry both validators, so pollers can revalidate cheaply:

- `ETag` is computed from the response body, so it is exact for the requested selector.
- `Last-Modified` is the last time flagd observed a change to **any** flag configuration, so it is conservative; it may cost a request a `304` it could have had, but never serves a stale one.

`If-None-Match` and `If-Modified-Since` are both honored, with `If-None-Match` taking precedence when both are sent:

```shell
curl -H 'If-None-Match: "<etag>"' -i http://localhost:8015/v1/flags   # 304 Not Modified
```

### Chaining flagd instances

Because `GET /v1/flags` returns an ordinary flag configuration document, another flagd instance can consume it directly as an HTTP sync source:

```shell
flagd start --uri http://localhost:8015/v1/flags
```

To sync only a subset, set the selector header on the source, which needs the `--sources` form:

```shell
flagd start --sources='[{"uri":"http://localhost:8015/v1/flags","provider":"http","headers":{"Flagd-Selector":"flagSetId=payment-flags"}}]'
```

---

## Quick Reference

### Ports

| Port | Protocol    | Service    | Description                                                            |
| ---- | ----------- | ---------- | ---------------------------------------------------------------------- |
| 8013 | gRPC / HTTP | Evaluation | Flag evaluation API (evaluation.proto, also HTTP/JSON via Connect)     |
| 8014 | HTTP        | Management | Health checks, metrics                                                 |
| 8015 | gRPC / HTTP | Sync       | Flag sync (sync.proto, also HTTP/JSON via Connect and `GET /v1/flags`) |
| 8016 | HTTP        | OFREP      | OpenFeature Remote Evaluation Protocol                                 |

### Health Check

```shell
curl http://localhost:8014/readyz
```

---

## See Also

- [Flag Definitions](./flag-definitions.md) - Complete flag definition reference
- [OFREP Service](./flagd-ofrep.md) - OFREP API details
- [gRPC Sync Service](./grpc-sync-service.md) - Sync service details
- [Sync Configuration](./sync-configuration.md) - Configure flag sources
- [CLI Reference](./flagd-cli/flagd_start.md) - Complete CLI options
