---
description: flagd as a gRPC sync service
---

# Overview

flagd can expose a gRPC sync service, allowing in-process providers to obtain their flag definitions.
The gRPC sync stream contains flag definitions currently configured at flagd as [sync-configurations](./sync-configuration.md).

```mermaid
---
title: gRPC sync
---
erDiagram
    flagd ||--o{ "sync (file)" : watches
    flagd ||--o{ "sync (http)" : polls
    flagd ||--o{ "sync (grpc)" : "sync.proto (gRPC/stream)"
    flagd ||--o{ "sync (kubernetes)" : watches
    "In-Process provider" ||--|| flagd : "gRPC sync stream (default port 8015)"
```

You may change the default port of the service using startup flag `--sync-port` (or `-g` shothand flag).

By default, the gRPC stream exposes all the flag configurations, with conflicting flag keys merged following flag's standard merge strategy.
You can read more about the merge strategy in our dedicated [concepts guide on syncs](../concepts/syncs.md).

If you specify a `selector` in the gRPC sync request, the gRPC service will attempt match the provided selector value to a source, and stream just the flags identified in that source.
For example, if `selector` is set to `myFlags.json`, service will stream flags observed from `myFlags.json` file.
Note that, to observe flags from `myFlags.json` file, you may use startup option `uri` like `--uri myFlags.json` or `source` option `--sources='[{"uri":"myFlags.json", provider":"file"}]`.
And the request will fail if there is no flag source matching the requested `selector`.

flagd provider implementations expose the ability to define the `selector` value. Please consider below example for Java,

```java
final FlagdProvider flagdProvider =
        new FlagdProvider(FlagdOptions.builder()
                .resolverType(Config.Evaluator.IN_PROCESS)
                .host("localhost")
                .port(8015)
                .selector("myFlags.json")
                .build());
```

See the [cheat sheet](./cheat-sheet.md#grpc-sync-api-syncproto) for `grpcurl` examples using `FetchAllFlags` and `SyncFlags`.
