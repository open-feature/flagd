---
description: flagd provider specification
---

# flagd Providers

!!! note

    This document serves as both a specification and general documentation for flagd providers.
    For language-specific details, see the `README.md` for the provider in question.

flagd providers are as essential as the flagd daemon itself, acting as the "bridge" between the OpenFeature SDK and flagd.
In fact, flagd providers may be the most crucial part of the flagd framework, as they can be used without an active flagd instance.
This document outlines their behavior and configuration.

## Naming

Consistent with our [naming conventions](../naming.md), the flagd provider name (no matter what language or runtime is in use) is `flagd`.

## Overview of Resolver Types

There are three resolver types for flagd providers:

- RPC Resolver: Evaluates flags remotely by connecting to a flagd instance using the gRPC evaluation protocol.
- In-Process Resolver: Downloads the flag set rules from a flagd instance and evaluates them locally, offering low-latency performance.
- File Resolver: Similar to in-process in terms of performance and latency; operates offline by reading flag definitions from a file, ideal for environments without network connectivity.

Each resolver type has unique characteristics and use cases. Below, we delve into their specifics.

## flagd Provider Lifecycle

flagd providers are built to adhere to the [provider lifecycle](https://openfeature.dev/specification/sections/flag-evaluation/#17-provider-lifecycle-management) defined in the OpenFeature specification.
Understanding the flagd provider lifecycle is helpful in configuring and optimizing your flagd deployment, and critical to implementing a flagd provider.

The lifecycle is summarized below:

- on initialization, attempt to connect the appropriate stream according to the resolver type ([sync](https://buf.build/open-feature/flagd/docs/main:flagd.sync.v1#flagd.sync.v1.FlagSyncService.SyncFlags) stream for in-process vs [event](https://buf.build/open-feature/flagd/docs/main:flagd.evaluation.v1#flagd.evaluation.v1.Service.EventStream) stream for RPC) and in the case of in-process, fetch the [sync-metadata](https://buf.build/open-feature/flagd/docs/main:flagd.sync.v1#flagd.sync.v1.FlagSyncService.GetMetadata)
    - if stream connection succeeds within the time specified by `deadline`, return from initialization (SDK will emit `PROVIDER_READY`) and for in-process providers, store the `flag set` rules
    - if stream connection fails or exceeds the time specified by `deadline`, abort initialization (SDK will emit `PROVIDER_ERROR`), and attempt to [reconnect](#stream-reconnection)
- while connected:
    - flags are resolved according to resolver mode; either by calling evaluation RPCs, or by evaluating the stored `flag set` rules
    - for RPC providers, flags resolved with `reason=STATIC` are [cached](#flag-evaluation-caching)
    - if flags change the associated stream (event or sync) indicates flags have changed, flush cache, or update `flag set` rules respectively and emit `PROVIDER_CONFIGURATION_CHANGED`
- if stream disconnects:
    - [reconnect](#stream-reconnection) with exponential backoff offered by GRPC.
        - if disconnected time <= `retryGracePeriod`
            - emit `PROVIDER_STALE`
            - RPC mode resolves `STALE` from cache where possible
            - in-process mode resolves `STALE` from stored `flag set` rules
        - if disconnected time > `retryGracePeriod`
            - emit `PROVIDER_ERROR`
            - RPC mode evaluation cache is purged
            - in-process mode resolves `STALE` from stored `flag set` rules
- on stream reconnection:
    - emit `PROVIDER_READY` and `PROVIDER_CONFIGURATION_CHANGED`
    - in-process providers store the latest `flag set` rules
- emit `PROVIDER_CONFIGURATION_CHANGED` event and update `flag set` rules when a `configuration_change` message is received on the streaming connection
- on shutdown, close the streaming connection in the`shutdown` function

```mermaid
stateDiagram-v2
    [*] --> NOT_READY
    NOT_READY --> READY: initialize
    NOT_READY --> ERROR: initialize
    READY --> ERROR: disconnected, disconnected period == 0
    READY --> STALE: disconnected, disconnect period < retry grace period
    STALE --> ERROR: disconnect period >= retry grace period
    ERROR --> READY: reconnected
    ERROR --> [*]: shutdown

    note right of STALE
        stream disconnected, attempting to reconnect,
        resolve from cache*
        resolve from flag set rules**
        STALE emitted
    end note

    note right of READY
        stream connected,
        evaluation cache active*,
        flag set rules stored**,
        metadata fetched**
        READY emitted
        CHANGE emitted with stream messages
    end note

    note right of ERROR
        stream disconnected, attempting to reconnect,
        evaluation cache purged*,
        ERROR emitted
    end note

%% * RPC providers only
%% ** In-Process providers only
```

```pseudo
*   RPC providers only
**  In-Process providers only
```

### Stream Reconnection

When either stream (sync or event) disconnects, whether due to the associated deadline being exceeded, network error or any other cause, the provider attempts to re-establish the stream immediately, and then retries with an exponential back-off.
We always rely on the [integrated functionality of GRPC for reconnection](https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md) and utilize [Wait-for-Ready](https://grpc.io/docs/guides/wait-for-ready/) to re-establish the stream.
We are configuring the underlying reconnection mechanism whenever we can, based on our configuration. (not all GRPC implementations support this)

| language/property | min connect timeout               | max backoff              | initial backoff          | jitter | multiplier |
| ----------------- | --------------------------------- | ------------------------ | ------------------------ | ------ | ---------- |
| GRPC property     | grpc.initial_reconnect_backoff_ms | max_reconnect_backoff_ms | min_reconnect_backoff_ms | 0.2    | 1.6        |
| Flagd property    | deadlineMs                        | retryBackoffMaxMs        | retryBackoffMs           | 0.2    | 1.6        |
| ---               | ---                               | ---                      | ---                      | ---    | ---        |
| default [^1]      | ✅                                 | ✅                        | ✅                        | 0.2    | 1.6        |
| js                | ✅                                 | ✅                        | ❌                        | 0.2    | 1.6        |
| java              | ❌                                 | ❌                        | ❌                        | 0.2    | 1.6        |

[^1] : C++, Python, Ruby, Objective-C, PHP, C#, js(deprecated)

When disconnected, if the time since disconnection is less than `retryGracePeriod`, the provider emits `STALE` when it disconnects.
While the provider is in state `STALE` the provider resolves values from its cache or stored flag set rules, depending on its resolver mode.
When the time since the last disconnect first exceeds `retryGracePeriod`, the provider emits `ERROR`.
The provider attempts to reconnect indefinitely, with a maximum interval of `retryBackoffMaxMs`.

## RPC Resolver

RPC providers use the [evaluation protocol](./protos.md#flagdevaluationv1evaluationproto) to connect to flagd, initiate the [event stream](./protos.md#eventstreamresponse), listen for changes in the flag definitions, and evaluate flags remotely by calling flagd.
RPC providers are relatively simple to implement since they essentially call a remote flagd instance with relevant parameters, and then flagd responds with the resolved flag value.
Of course, this means there's latency associated with RPC providers, though this is mitigated somewhat by [caching](#flag-evaluation-caching).

### Flag Evaluation Caching

In RPC mode, `flagd` uses a caching mechanism which greatly reduces latency for static flags (flags without targeting rules).
Evaluations for flags with targeting rules are never cached.

!!! note

    Evaluation caching is only relevant to when the RPC resolver is used; the in-process resolver stores a complete set of rules for a `flag set`, which means evaluation can be done locally, with low latency.

#### Cacheable flags

`flagd` sets the `reason` of a flag evaluation as `STATIC` when no targeting rules are configured for the flag.
A client can safely store the result of a static evaluation in its cache indefinitely (until the configuration of the flag changes, see [cache invalidation](#cache-invalidation)).

Put simply in pseudocode:

```pseudo
if reason == "STATIC" {
    isFlagCacheable = true
}
```

#### Cache invalidation

`flagd` emits events to the server-to-client stream, among these is the `configuration_change` event.
The structure of this event is as such:

```json
{
    "type": "delete", // ENUM:["delete","write","update"]
    "source": "/flag-configuration.json", // the source of the flag configuration
    "flagKey": "foo"
}
```

A client should invalidate the cache of any flag found in a `configuration_change` event to prevent stale data.
If the connection drops all cache values must be cleared (any number of events may have been missed).

### Client Side Providers

Client side flagd providers (used in mobile and front-end web applications) have unique security and performance considerations.
These flagd providers only support the RPC resolver mode (so that `flag set` rules, which might contain sensitive information, are never sent to the client).
Instead, these do bulk evaluations of all flags in the `flag set`, and cache the results until they are invalidated.
Bulk evaluations take place when:

- the provider is initialized
- the context is changed
- a change in the definition notifies the provider it should re-evaluate the flags

This pattern is consistent with OpenFeature's [static context paradigm](https://openfeature.dev/specification/glossary#static-context-paradigm).

!!! note

    To support easy integration with mobile and browser use cases, flagd's [evaluation protocol](./protos.md#flagdevaluationv1evaluationproto) is accessible over both gRPC and HTTP

!!! note

    flagd supports the OFREP protocol, meaning client-side OFREP providers can also be used for client-side use-cases.

<!-- markdownlint-disable MD024 -->
### Provider Metadata
<!-- markdownlint-enable MD024 -->

The provider metadata includes properties returned from the [provider_ready event payload](./protos.md#eventstreamresponse) data.

## In-Process Resolver

In-process providers use the [sync schema](./protos.md#syncflagsresponse) to connect to flagd, initiate the [sync stream](./protos.md#eventstreamresponse), and download the `flag set` rules to evaluate them locally.
In-process providers are relatively complex (compared to RPC providers) to implement since they essentially must implement more of flagd's logic to evaluate flags locally.
Local evaluation has the impact of much lower latency and almost no serialization compared to RPC providers.

### JsonLogic Evaluation

An in-process flagd providers provide the feature set offered by [JsonLogic](https://jsonlogic.com) to evaluate flag resolution requests for a given context.

### Custom JsonLogic Evaluators

In addition to the built-in evaluators provided by JsonLogic, the following custom targeting rules are implemented by the provider:

- [Fractional operation](../../reference/custom-operations/fractional-operation.md)
- [Semantic version evaluation](../../reference/custom-operations/semver-operation.md)
- [StartsWith/EndsWith evaluation](../../reference/custom-operations/string-comparison-operation.md)

### Targeting Key

Similar to the flagd daemon, in-process providers map the [targeting-key](https://openfeature.dev/specification/glossary#targeting-key) into a top level property of the context used in rules, with the key `"targetingKey"`.

### `$flagd` Properties in the Evaluation Context

Similar to the flagd daemon, in-process flagd providers add the following properties to the JsonLogic evaluation context so that users can use them in their targeting rules.
Conflicting properties in the context will be overwritten by the values below.

| Property           | Description                                             |
| ------------------ | ------------------------------------------------------- |
| `$flagd.flagKey`   | the identifier for the flag being evaluated             |
| `$flagd.timestamp` | a unix timestamp (in seconds) of the time of evaluation |

### Changed Flags

When a new flag definition is parsed, the stored flags are compared with the newly parsed flags.
Flags which have been removed, added, or mutated (considering, at a minimum, their `default variant`, `targeting rules`, and `metadata`) have their keys added to the `flags changed` field of the associated `PROVIDER_CONFIGURATION_CHANGED` event.`

### Sync-Metadata Properties in the Evaluation Context

In-process flagd providers also inject any properties returned by the [sync-metadata RPC response](./protos.md#getmetadataresponse) into the context.
This allows for static properties defined in flagd to be added to in-process evaluations.
If only a subset of the sync-metadata response is desired to be injected into the evaluation context, you can define a mapping function with the `contextEnricher` option.

<!-- markdownlint-disable MD024 -->
### Provider Metadata
<!-- markdownlint-enable MD024 -->

The provider metadata includes the top-level metadata properties in the [flag definition](../flag-definitions.md).

## File Resolver (Offline Mode)

The in-process resolver mode can also use a file based [flag definition](../flag-definitions.md).
This does not connect to a flagd instance or gRPC sync implementation, and instead polls a flag definition from a file.
If the file has been modified since the last poll (based on the file metadata) and [flags have changed](#changed-flags), a `PROVIDER_CONFIGURATION_CHANGED` event with the appropriate `changed flags` field is emitted.

The Evaluation uses [JsonLogic](#jsonlogic-evaluation) and [custom JsonLogic evaluators](#custom-jsonlogic-evaluators) like the [InProcess Resolver](#in-process-resolver).

!!! note

    This mode does not support [context enrichment via sync-metadata](#sync-metadata-properties-in-the-evaluation-context).

## Configuration

### Configuration options

Most options can be defined in the constructor, or as environment variables, with constructor options having the highest
precedence.

Below are the supported configuration parameters (note that not all apply to both resolver modes):

| Option name           | Environment variable name      | Explanation                                                                          | Type & Values                | Default                       | Compatible resolver     |
| --------------------- | ------------------------------ | ------------------------------------------------------------------------------------ | ---------------------------- | ----------------------------- | ----------------------- |
| resolver              | FLAGD_RESOLVER                 | mode of operation                                                                    | String - `rpc`, `in-process` | rpc                           | rpc & in-process        |
| host                  | FLAGD_HOST                     | remote host                                                                          | String                       | localhost                     | rpc & in-process        |
| port                  | FLAGD_PORT                     | remote port                                                                          | int                          | 8013 (rpc), 8015 (in-process) | rpc & in-process        |
| targetUri             | FLAGD_TARGET_URI               | alternative to host/port, supporting custom name resolution                          | string                       | null                          | rpc & in-process        |
| tls                   | FLAGD_TLS                      | connection encryption                                                                | boolean                      | false                         | rpc & in-process        |
| socketPath            | FLAGD_SOCKET_PATH              | alternative to host port, unix socket                                                | String                       | null                          | rpc & in-process        |
| certPath              | FLAGD_SERVER_CERT_PATH         | tls cert path                                                                        | String                       | null                          | rpc & in-process        |
| deadlineMs            | FLAGD_DEADLINE_MS              | deadline for unary calls, and timeout for initialization                             | int                          | 500                           | rpc & in-process & file |
| streamDeadlineMs      | FLAGD_STREAM_DEADLINE_MS       | deadline for streaming calls, useful as an application-layer keepalive               | int                          | 600000                        | rpc & in-process        |
| retryBackoffMs        | FLAGD_RETRY_BACKOFF_MS         | initial backoff for stream retry                                                     | int                          | 1000                          | rpc & in-process        |
| retryBackoffMaxMs     | FLAGD_RETRY_BACKOFF_MAX_MS     | maximum backoff for stream retry                                                     | int                          | 120000                        | rpc & in-process        |
| retryGracePeriod      | FLAGD_RETRY_GRACE_PERIOD       | period in seconds before provider moves from STALE to ERROR state                    | int                          | 5                             | rpc & in-process & file |
| keepAliveTime         | FLAGD_KEEP_ALIVE_TIME_MS       | http 2 keepalive                                                                     | long                         | 0                             | rpc & in-process        |
| cache                 | FLAGD_CACHE                    | enable cache of static flags                                                         | String - `lru`, `disabled`   | lru                           | rpc                     |
| maxCacheSize          | FLAGD_MAX_CACHE_SIZE           | max size of static flag cache                                                        | int                          | 1000                          | rpc                     |
| selector              | FLAGD_SOURCE_SELECTOR          | Selector expression to filter flags (e.g., `flagSetId=my-app`, `source=config.json`) | string                       | null                          | in-process              |
| providerId            | FLAGD_PROVIDER_ID              | A unique identifier for flagd(grpc client) initiating the request.                   | string                       | null                          | in-process              |
| offlineFlagSourcePath | FLAGD_OFFLINE_FLAG_SOURCE_PATH | offline, file-based flag definitions, overrides host/port/targetUri                  | string                       | null                          | file                    |
| offlinePollIntervalMs | FLAGD_OFFLINE_POLL_MS          | poll interval for reading offlineFlagSourcePath                                      | int                          | 5000                          | file                    |
| contextEnricher       | -                              | sync-metadata to evaluation context mapping function                                 | function                     | identity function             | in-process              |

### Custom Name Resolution

Some implementations support [gRPC custom name resolution](https://grpc.io/docs/guides/custom-name-resolution/), and abstractions to introduce additional resolvers.
Specifically, a custom resolver for `envoy` has been implemented in some providers, which overrides the authority header with the authority specified in the envoy URL scheme.
Below is an example of a custom target string which will use envoy sidecar proxy for name resolution:

```text
envoy://localhost:9211/flagd-sync.service
```

The custom name resolver provider in this case will use the endpoint name i.e. `flagd-sync.service` as [authority](https://github.com/grpc/grpc-java/blob/master/examples/src/main/java/io/grpc/examples/nameresolve/ExampleNameResolver.java#L55-L61)
and connect to `localhost:9211`.

### Selector Configuration

Providers support selector configuration to filter which flags are synchronized or evaluated. This enables more granular control in multi-tenant or multi-environment deployments.

#### Selector Syntax

Providers accept selector expressions using the following syntax:

- **Flag Set Selection**: `flagSetId=<identifier>` - Target flags belonging to a specific flag set
- **Source Selection**: `source=<identifier>` - Target flags from a specific source (legacy)
- **Backward Compatibility**: `<identifier>` - Treated as source selection

#### Selector Precedence

When selectors are provided in multiple locations, the following precedence applies:

1. **Request Header**: `Flagd-Selector` header (RPC and OFREP requests)
2. **Provider Configuration**: `selector` option in provider constructor

#### Usage Examples

**Flag Set-Based Selection (Recommended):**

```javascript
const provider = new FlagdProvider({
  host: 'localhost',
  port: 8013,
  selector: 'flagSetId=user-service'
});
```

**Source-Based Selection (Legacy):**

```javascript
const provider = new FlagdProvider({
  host: 'localhost', 
  port: 8013,
  selector: 'source=config/app-flags.json'
});
```

**Header-Based Selection:**

```bash
# gRPC request with selector header
grpcurl -H "Flagd-Selector: flagSetId=payment-service" \
  localhost:8013 flagd.evaluation.v1.Service/ResolveBoolean

# OFREP request with selector header
curl -H "Flagd-Selector: flagSetId=frontend-features" \
  http://localhost:8014/ofrep/v1/evaluate/flags/my-flag
```

### Metadata and Metadata Reflection

#### Flag Metadata

When a flag is resolved, the returned [metadata](./flag-definitions.md#metadata) is a merged representation of the metadata defined on the flag set, and on the flag, with the flag metadata taking priority.
Flag metadata is returned on a "best effort" basis when flags are resolved: disabled, missing or erroneous flags return the metadata of the associated flag set whenever possible.
This is particularly important for debugging purposes and error metrics.

#### Selector Metadata Reflection

Flagd "reflects" selector information back in response metadata, providing transparency about query execution. This helps with debugging selector expressions and understanding which flags were actually queried.

**Example - gRPC Response:**

```protobuf
// Request with selector header: "Flagd-Selector: flagSetId=payment-service"
message ResolveBooleanResponse {
  bool value = 1;
  string reason = 2;
  string variant = 3;
  google.protobuf.Struct metadata = 4; // Contains reflected selector info
}
```

**Example - OFREP Response:**

```json
{
  "value": true,
  "reason": "TARGETING_MATCH",
  "variant": "on",
  "metadata": {
    "flagSetId": "payment-service",
    "team": "payments",
    "version": "1.2.0"
  }
}
```

#### Debugging with Metadata Reflection

Use reflected metadata to:

- **Verify Selector Parsing**: Confirm your selector was interpreted correctly
- **Debug Empty Results**: Check if selectors are filtering flags as expected
- **Audit Access Patterns**: Log selector metadata for compliance and monitoring
- **Troubleshoot Configuration**: Identify selector precedence issues
