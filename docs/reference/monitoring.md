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

flagd exposes the following metrics:

- `http.server.request.duration` - Measures the duration of inbound HTTP requests
- `http.server.response.body.size` - Measures the size of HTTP response messages
- `http.server.active_requests` - Measures the number of concurrent HTTP requests that are currently in-flight
- `feature_flag.flagd.impression` - Measures the number of evaluations for a given flag
- `feature_flag.flagd.evaluation.reason` - Measures the number of evaluations for a given reason

> Please note that metric names may vary based on the consuming monitoring tool naming requirements.
> For example, the transformation of OTLP metrics to Prometheus is described [here](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/compatibility/prometheus_and_openmetrics.md#otlp-metric-points-to-prometheus).

### HTTP Metric Attributes

flagd uses OpenTelemetry Semantic Conventions v1.34.0 for HTTP metrics. The following attributes are included with HTTP metrics:

- `service.name` - The name of the service
- `http.route` - The matched route (path template)
- `http.request.method` - The HTTP request method (GET, POST, etc.)
- `http.response.status_code` - The HTTP response status code
- `url.scheme` - The URI scheme (http or https)

## Traces

flagd creates the following spans as part of a trace:

- `flagEvaluationService(resolveX)` - SpanKind server
    - `jsonEvaluator(resolveX)` - SpanKind internal
- `jsonEvaluator(setState)` - SpanKind internal

## Export to OTEL collector

flagd can be configured to connect to [OTEL collector](https://opentelemetry.io/docs/collector/). This requires startup
flag `metrics-exporter` to be `otel` and a valid `otel-collector-uri`. For example,

`flagd start --uri file:/flags.json --metrics-exporter otel --otel-collector-uri localhost:4317`

### Configure local collector setup

To configure a local collector setup along with Jaeger and Prometheus, you can use following sample docker-compose
file and configuration files.

Note - content is adopted from
official [OTEL collector example](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/examples/demo)

#### docker-compose.yaml

```yaml
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
        endpoint: 0.0.0.0:4317
exporters:
  prometheus:
    endpoint: "0.0.0.0:8889"
    const_labels:
      label1: value1
  otlp/jaeger:
    endpoint: jaeger-all-in-one:4317
    tls:
      insecure: true
processors:
  batch:
service:
  pipelines:
    traces:
      receivers: [ otlp ]
      processors: [ batch ]
      exporters: [ otlp/jaeger ]
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

## Metadata

[Flag metadata](https://openfeature.dev/specification/types/#flag-metadata) comprises auxiliary data pertaining to feature flags; it's highly valuable in telemetry signals.
Flag metadata might consist of attributes indicating the version of the flag, an identifier for the flag set, ownership information about the flag, or other documentary information.
flagd supports flag metadata in all its [gRPC protocols](../reference/specifications//protos.md), in [OFREP](../reference/flagd-ofrep.md), and in its [flag definitions](./flag-definitions.md#metadata).
These attributes are returned with flag evaluations, and can be added to telemetry signals as outlined in the [OpenFeature specification](https://openfeature.dev/specification/appendix-d).
