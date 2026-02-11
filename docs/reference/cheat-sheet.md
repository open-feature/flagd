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

## Sample Flag Definitions

The examples below use these sample flag definition files. Download them to follow along:

- [cheat-sheet-flags.json](../assets/cheat-sheet-flags.json) - General application flags (flagSetId: `app-flags`)
- [cheat-sheet-flags-payments.json](../assets/cheat-sheet-flags-payments.json) - Payment-related flags (flagSetId: `payment-flags`)

The `app-flags` set includes:

| Flag Key | Type | Description |
|----------|------|-------------|
| `simple-boolean` | boolean | Static boolean flag |
| `simple-string` | string | Static string flag |
| `simple-number` | integer | Static numeric flag |
| `simple-object` | object | Static object flag |
| `user-tier-flag` | string | Context-sensitive flag based on `tier` |
| `email-based-feature` | boolean | Context-sensitive flag based on `email` domain |
| `region-config` | object | Context-sensitive flag based on `region` |

The `payment-flags` set includes:

| Flag Key | Type | Description |
|----------|------|-------------|
| `payment-provider` | string | Static payment provider selection |
| `max-transaction-amount` | integer | Context-sensitive based on `account-verified` |
| `enable-crypto-payments` | boolean | Context-sensitive based on `country` |

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

---

## OFREP API (HTTP)

The [OFREP (OpenFeature Remote Evaluation Protocol)](https://openfeature.dev/docs/reference/other-technologies/ofrep/) API is available on port `8016` by default.

### Evaluate a Single Flag

```shell
curl -X POST 'http://localhost:8016/ofrep/v1/evaluate/flags/simple-boolean'
```

Response:

```json
{
  "key": "simple-boolean",
  "reason": "STATIC",
  "variant": "on",
  "value": true,
  "metadata": {}
}
```

### Evaluate Different Flag Types

```shell
# String flag
curl -X POST 'http://localhost:8016/ofrep/v1/evaluate/flags/simple-string'

# Number flag
curl -X POST 'http://localhost:8016/ofrep/v1/evaluate/flags/simple-number'

# Object flag
curl -X POST 'http://localhost:8016/ofrep/v1/evaluate/flags/simple-object'
```

### Evaluate Multiple Flags (Bulk Evaluation)

```shell
curl -X POST 'http://localhost:8016/ofrep/v1/evaluate/flags'
```

Response:

```json
{
  "flags": [
    {"key": "simple-boolean", "reason": "STATIC", "variant": "on", "value": true, "metadata": {}},
    {"key": "simple-string", "reason": "STATIC", "variant": "greeting", "value": "Hello, World!", "metadata": {}},
    {"key": "simple-number", "reason": "STATIC", "variant": "medium", "value": 50, "metadata": {}}
  ]
}
```

---

## Context-Aware Evaluation

### Context from Request Body

Pass evaluation context in the request body to trigger targeting rules:

```shell
# Evaluate with email context (triggers email-based-feature targeting)
curl -X POST 'http://localhost:8016/ofrep/v1/evaluate/flags/email-based-feature' \
  -H 'Content-Type: application/json' \
  -d '{"context": {"email": "user@example.com"}}'
```

Response (email matches `@example.com`):

```json
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

---

## gRPC Evaluation API (evaluation.proto)

The gRPC evaluation service is available on port `8013` by default.
Use [grpcurl](https://github.com/fullstorydev/grpcurl) to interact with it.

!!! note "Proto files required"
    flagd does not support gRPC reflection. You must provide the proto files to grpcurl so it knows how to serialize/deserialize requests and responses.
    Clone the flagd repo or download the protos from [buf.build/open-feature/flagd](https://buf.build/open-feature/flagd).

```shell
# Clone the repo for proto files
git clone git@github.com:open-feature/flagd-schemas.git
PROTO_DIR="flagd-schemas/protobuf/"
```

### Evaluate Flags by Type

```shell
# Boolean flag
grpcurl -plaintext \
  -import-path "$PROTO_DIR" -proto flagd/evaluation/v1/evaluation.proto \
  -d '{"flagKey": "simple-boolean", "context": {}}' \
  localhost:8013 \
  flagd.evaluation.v1.Service/ResolveBoolean

# String flag
grpcurl -plaintext \
  -import-path "$PROTO_DIR" -proto flagd/evaluation/v1/evaluation.proto \
  -d '{"flagKey": "simple-string", "context": {}}' \
  localhost:8013 \
  flagd.evaluation.v1.Service/ResolveString

# With context
grpcurl -plaintext \
  -import-path "$PROTO_DIR" -proto flagd/evaluation/v1/evaluation.proto \
  -d '{"flagKey": "user-tier-flag", "context": {"tier": "enterprise"}}' \
  localhost:8013 \
  flagd.evaluation.v1.Service/ResolveString
```

### Evaluate All Flags

```shell
grpcurl -plaintext \
  -import-path "$PROTO_DIR" -proto flagd/evaluation/v1/evaluation.proto \
  -d '{"context": {}}' \
  localhost:8013 \
  flagd.evaluation.v1.Service/ResolveAll
```

### Using Selector Header

Filter which flags are evaluated using the `Flagd-Selector` header:

```shell
# Evaluate only flags from the app flag set
grpcurl -plaintext \
  -import-path "$PROTO_DIR" -proto flagd/evaluation/v1/evaluation.proto \
  -H 'Flagd-Selector: flagSetId=app-flags' \
  -d '{"flagKey": "simple-boolean", "context": {}}' \
  localhost:8013 \
  flagd.evaluation.v1.Service/ResolveBoolean

# ResolveAll with selector
grpcurl -plaintext \
  -import-path "$PROTO_DIR" -proto flagd/evaluation/v1/evaluation.proto \
  -H 'Flagd-Selector: flagSetId=payment-flags' \
  -d '{"context": {}}' \
  localhost:8013 \
  flagd.evaluation.v1.Service/ResolveAll
```

---

## gRPC Sync API (sync.proto)

The gRPC sync service is available on port `8015` by default.
This is used by in-process providers to fetch and sync flag configurations.

Use [grpcurl](https://github.com/fullstorydev/grpcurl) to interact with it.

!!! note "Proto files required"
    flagd does not support gRPC reflection. You must provide the proto files to grpcurl so it knows how to serialize/deserialize requests and responses.
    Clone the flagd repo or download the protos from [buf.build/open-feature/flagd](https://buf.build/open-feature/flagd).

### FetchAllFlags

Get all flag configurations as a single response:

```shell
grpcurl -plaintext \
  -import-path "$PROTO_DIR" -proto flagd/sync/v1/sync.proto \
  -d '{}' \
  localhost:8015 \
  flagd.sync.v1.FlagSyncService/FetchAllFlags
```

Response contains the complete flag configuration JSON:

```json
{
  "flagConfiguration": "{\"flags\":{\"simple-boolean\":{...}}}"
}
```

### FetchAllFlags with Selector

Filter which flag source's configuration is returned:

```shell
# Fetch only flags from cheat-sheet-flags.json
grpcurl -plaintext \
  -import-path "$PROTO_DIR" -proto flagd/sync/v1/sync.proto \
  -d '{"selector": "flagSetId=app-flags"}' \
  localhost:8015 \
  flagd.sync.v1.FlagSyncService/FetchAllFlags

# Fetch only payment flags
grpcurl -plaintext \
  -import-path "$PROTO_DIR" -proto flagd/sync/v1/sync.proto \
  -d '{"selector": "flagSetId=payment-flags"}' \
  localhost:8015 \
  flagd.sync.v1.FlagSyncService/FetchAllFlags

# With provider ID for identification
grpcurl -plaintext \
  -import-path "$PROTO_DIR" -proto flagd/sync/v1/sync.proto \
  -d '{"providerId": "my-app-sidecar", "selector": "flagSetId=app-flags"}' \
  localhost:8015 \
  flagd.sync.v1.FlagSyncService/FetchAllFlags
```

### SyncFlags (Streaming)

Establish a server-streaming connection that pushes flag configuration updates in real-time:

```shell
grpcurl -plaintext \
  -import-path "$PROTO_DIR" -proto flagd/sync/v1/sync.proto \
  -d '{}' \
  localhost:8015 \
  flagd.sync.v1.FlagSyncService/SyncFlags
```

The stream outputs the initial flag configuration and continues streaming updates when flags change.

### SyncFlags with Selector

```shell
# Stream only changes to app flags
grpcurl -plaintext \
  -import-path "$PROTO_DIR" -proto flagd/sync/v1/sync.proto \
  -d '{"selector": "flagSetId=app-flags"}' \
  localhost:8015 \
  flagd.sync.v1.FlagSyncService/SyncFlags

# Stream with provider ID
grpcurl -plaintext \
  -import-path "$PROTO_DIR" -proto flagd/sync/v1/sync.proto \
  -d '{"providerId": "my-service", "selector": "flagSetId=app-flags"}' \
  localhost:8015 \
  flagd.sync.v1.FlagSyncService/SyncFlags
```

### Using Selector Header with Sync API

The `Flagd-Selector` header can be used as an alternative to the request body selector field:

```shell
grpcurl -plaintext \
  -import-path "$PROTO_DIR" -proto flagd/sync/v1/sync.proto \
  -H 'Flagd-Selector: flagSetId=app-flags' \
  -d '{}' \
  localhost:8015 \
  flagd.sync.v1.FlagSyncService/FetchAllFlags
```

If both header and request body contain a selector, the header takes precedence.

---

## Quick Reference

### Ports

| Port | Protocol | Service | Description |
|------|----------|---------|-------------|
| 8013 | gRPC | Evaluation | Flag evaluation API (evaluation.proto) |
| 8014 | HTTP | Management | Health checks, metrics |
| 8015 | gRPC | Sync | Flag sync for in-process providers (sync.proto) |
| 8016 | HTTP | OFREP | OpenFeature Remote Evaluation Protocol |

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
