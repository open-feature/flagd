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

```shell
curl -X POST --compressed 'http://localhost:8016/ofrep/v1/evaluate/flags'
```

By default, responses smaller than 1024 bytes are left uncompressed. Gzip costs a near-fixed few microseconds per response regardless of size, so on small bodies that CPU buys very little: a single-flag evaluation is around 150 bytes and barely compresses at all, while a 10-flag bulk response compresses roughly 3.5x and a 100-flag response roughly 7x. Use `--ofrep-compression-min-size` to move the threshold, or set it to `0` to compress every response whatever its size.

Compression can be turned off entirely with `--ofrep-compression=false`, which is worth doing when a proxy in front of flagd already compresses.

The SSE stream is never compressed, so events reach subscribers as soon as they are flushed.

A compressed bulk evaluation response carries a distinct `ETag`, suffixed with `-gzip`, because the compressed and uncompressed bodies are different representations and must not share one validator. flagd strips the suffix when matching `If-None-Match`, so clients keep getting their `304 Not Modified` and do not need to do anything special.

## Monitoring

The OFREP endpoint is instrumented with OpenTelemetry HTTP and flag evaluation metrics.
See the [Monitoring reference](./monitoring.md#http-metrics) for the full list of exposed metrics and their attributes.
