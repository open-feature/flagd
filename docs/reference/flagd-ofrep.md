---
description: OpenFeature Remote Evaluation Protocol with flagd
---

# Overview

![EXPERIMENTAL](https://img.shields.io/badge/status-experimental-red)

flagd supports the [OpenFeature Remote Evaluation Protocol](https://github.com/open-feature/protocol) for flag evaluations.
The service starts on port `8016` by default and this can be changed using startup flag `--ofrep-port` (or `-r` shothand flag).

## Usage

Given flagd is running with flag configuration for `myBoolFlag`, you can evaluate the flag with OFREP API with following curl request,

```shell
curl -X POST 'http://localhost:8016/ofrep/v1/evaluate/flags/myBoolFlag'
```

To evaluate all flags currently configured at flagd, use OFREP bulk evaluation request,

```shell
curl -X POST 'http://localhost:8016/ofrep/v1/evaluate/flags'
```

## Evaluation Caching

The bulk evaluation endpoint caches responses per selector to avoid redundant evaluations. Clients can use the `If-None-Match` header with a previously received `ETag` to check if the cache is still valid. When the ETag matches, flagd returns the cached response without re-evaluating.

**Important**: The ETag only corresponds the flag configuration version, not the evaluation context. Clients must not send a cached ETag when their evaluation context has changed, otherwise they may receive stale results.

See the [cheat sheet](./cheat-sheet.md#ofrep-api-http) for more OFREP examples including context-sensitive evaluation and selectors.
