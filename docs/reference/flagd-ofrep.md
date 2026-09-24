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

See the [cheat sheet](./cheat-sheet.md#evaluating-flags) for more OFREP examples including context-sensitive evaluation and selectors.

## Compression

Evaluation responses are gzip compressed when the client sends an `Accept-Encoding: gzip` header, and served uncompressed otherwise.

```shell
curl -X POST --compressed 'http://localhost:8016/ofrep/v1/evaluate/flags'
```

A single-flag evaluation is around 150 bytes and barely shrinks, while a 10-flag bulk response compresses roughly 3.5x and a 100-flag response roughly 7x. Gzip costs a near-fixed few microseconds per response, so the small bodies do not repay it, but the saving on the bulk responses that matter is large enough that there is nothing to configure here.

The SSE stream is never compressed, so events reach subscribers as soon as they are flushed.

Bulk evaluation responses carry a weak `ETag`, which covers the compressed and uncompressed forms of the same body alike. A client revalidating with `If-None-Match` gets its `304 Not Modified` whether or not its `Accept-Encoding` has changed since.

## Monitoring

The OFREP endpoint is instrumented with OpenTelemetry HTTP and flag evaluation metrics.
See the [Monitoring reference](./monitoring.md#http-metrics) for the full list of exposed metrics and their attributes.
