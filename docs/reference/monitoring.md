---
description: monitoring and telemetry flagd and flagd providers
---

# Monitoring

## Readiness & Liveness probes

### HTTP

Flagd exposes HTTP liveness and readiness probes.
These probes can be used for K8s deployments.
With default start-up configurations, these probes are exposed on the management port (default: 8014) at the following URLs,

- Liveness: <http://localhost:8014/healthz>
- Readiness: <http://localhost:8014/readyz>

### gRPC

Flagd exposes a [standard gRPC liveness check](https://github.com/grpc/grpc/blob/master/doc/health-checking.md) on the management port (default: 8014).

### Definition of Liveness

The liveness probe becomes active and HTTP 200 status is served as soon as Flagd service is up and running.

### Definition of Readiness

The readiness probe becomes active similar to the liveness probe as soon as Flagd service is up and running.
However,
the probe emits HTTP 412 until all sync providers are ready.
This status changes to HTTP 200 when all sync providers at
least have one successful data sync.
The status does not change from there on.

## OpenTelemetry

flagd provides telemetry data out of the box. This telemetry data is compatible with OpenTelemetry.

By default, the Prometheus exporter is used for metrics which can be accessed via the `/metrics` endpoint. For example,
with default startup flags, metrics are exposed at `http://localhost:8014/metrics`.

Given below is the current implementation overview of flagd telemetry internals,

![flagd telemetry](../images/flagd-telemetry.png)

## Metrics

flagd expose following metrics,

- `http.server.duration`
- `http.server.response.size`
- `http.server.active_requests`
- `feature_flag.flagd.impression`
- `feature_flag.flagd.evaluation.reason`

## Traces

flagd expose following traces,

- `flagEvaluationService(resolveX)` - SpanKind server
  - `jsonEvaluator(resolveX)` - SpanKind internal
- `jsonEvaluator(setState)` - SpanKind internal

## Export to OTEL collector

flagd can be configured to export telemetry to the [OpenTelemetry (OTel) Collector](https://opentelemetry.io/docs/collector/) using standard
[OTel SDK environment variables](https://opentelemetry.io/docs/concepts/sdk-configuration/otlp-exporter-configuration). For example,

```shell
export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
export OTEL_EXPORTER_OTLP_PROTOCOL=grpc
export OTEL_METRICS_EXPORTER=otlp

flagd start --uri file:/flags.json
```

### Configure local collector setup

To configure a local collector setup along with Jaeger and Prometheus, you can use following sample docker-compose
file and configuration files.

Note - content is adopted from
official [OTEL collector example](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/examples/demo)

#### docker-compose.yaml

```yaml
version: "3"
services:
  # Jaeger
  jaeger-all-in-one:
    image: jaegertracing/all-in-one:latest
    restart: always
    ports:
      - "16686:16686"
      - "14268"
      - "14250"
  # Collector
  otel-collector:
    image: otel/opentelemetry-collector:latest
    restart: always
    command: [ "--config=/etc/otel-collector-config.yaml" ]
    volumes:
      - ./otel-collector-config.yaml:/etc/otel-collector-config.yaml
    ports:
      - "1888:1888"   # pprof extension
      - "8888:8888"   # Prometheus metrics exposed by the collector
      - "8889:8889"   # Prometheus exporter metrics
      - "13133:13133" # health_check extension
      - "4317:4317"   # OTLP gRPC receiver
      - "55679:55679" # zpages extension
    depends_on:
      - jaeger-all-in-one
  prometheus:
    container_name: prometheus
    image: prom/prometheus:latest
    restart: always
    volumes:
      - ./prometheus.yaml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"
```

#### otel-collector-config.yaml

```yaml
receivers:
  otlp:
    protocols:
      grpc:
exporters:
  prometheus:
    endpoint: "0.0.0.0:8889"
    const_labels:
      label1: value1
  jaeger:
    endpoint: jaeger-all-in-one:14250
    tls:
      insecure: true
processors:
  batch:
service:
  pipelines:
    traces:
      receivers: [ otlp ]
      processors: [ batch ]
      exporters: [ jaeger ]
    metrics:
      receivers: [ otlp ]
      processors: [ batch ]
      exporters: [ prometheus ]
```

#### prometheus.yml

```yaml
scrape_configs:
  - job_name: 'otel-collector'
    scrape_interval: 10s
    static_configs:
      - targets: [ 'otel-collector:8889' ]
      - targets: [ 'otel-collector:8888' ]
```

Once, configuration files are ready, use `docker-compose up` to start the local setup. With successful startup, you can
access metrics through [Prometheus](http://localhost:9090/graph) & traces through [Jaeger](http://localhost:16686/).
