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

## HTTP Caching

The bulk evaluation endpoint supports HTTP caching via ETags (MD5-based). Clients can use the `If-None-Match` header with a previously received `ETag` to validate cached responses. When the response hasn't changed, the server returns `304 Not Modified` without a body, reducing bandwidth.

```shell
curl -X POST 'http://localhost:8016/ofrep/v1/evaluate/flags' \
  -H 'If-None-Match: "a1b2c3d4e5f6..."'
```

See the [cheat sheet](./cheat-sheet.md#ofrep-api-http) for more OFREP examples including context-sensitive evaluation and selectors.
