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

Evaluation responses are gzip compressed when the client sends an `Accept-Encoding: gzip` header.
Responses below 1KB are left uncompressed, since gzip framing costs more than it saves at that size; in practice this means bulk evaluation responses compress and single-flag evaluations do not.
The SSE stream is never compressed, so events reach subscribers as soon as they are flushed.

```shell
curl -X POST --compressed 'http://localhost:8016/ofrep/v1/evaluate/flags'
```

Compression does not affect the `ETag` on bulk evaluation responses: the tag is a digest of the uncompressed body, so an `If-None-Match` request still gets its `304 Not Modified` whichever encoding the original response used.

## Monitoring

The OFREP endpoint is instrumented with OpenTelemetry HTTP and flag evaluation metrics.
See the [Monitoring reference](./monitoring.md#http-metrics) for the full list of exposed metrics and their attributes.
