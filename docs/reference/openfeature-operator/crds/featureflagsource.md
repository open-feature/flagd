# FeatureFlagSources

`FeatureFlagSource` support multiple flag sources. Sources are configured as a list and given below are supported sources and their configurations,

## kubernetes aka `FeatureFlag`

This is `FeatureFlag` custom resource backed flagd feature flag definition.
Read more on the custom resource at the dedicated documentation of [FeatureFlag](https://github.com/open-feature/open-feature-operator/blob/main/docs/feature_flag_configuration.md)

To refer this custom resource in `FlagSource`, provider type `kubernetes` is used as below example,

```yaml
sources:                        
  - source: flags/sample-flags  # FeatureFlag - namespace/custom_resource_name
    provider: kubernetes        # kubernetes flag source backed by FeatureFlag custom resource
```

## flagd-proxy

`flagd-proxy` is an alternative to direct resource access on `FeatureFlag` custom resources.
This source type is useful when there is a need for restricting workload permissions and/or to reduce k8s API load.

Read more about proxy approach to access kubernetes resources: [flagd-proxy](https://github.com/open-feature/open-feature-operator/blob/main/docs/flagd_proxy.md)

## filepath

Injected sidecar can use volume mounted files as flag sources.
For this, provider type `filepath` is used as below example,

```yaml
sources:                        
  - source: /etc/flagd/config.json 
    provider: filepath          
```

## http

Feature flags can be sources from a http endpoint using provider type `http`,

```yaml
sources:
  - source: http://my-flag-source.json
    provider: http
    httpSyncBearerToken: token                  # optional bearer token for the http connection
```

## grpc

Given below is an example configuration with provider type `grpc` and supported options,

```yaml
sources:                        
  - source: my-flag-source:8080
    provider: grpc
    certPath: /certs/ca.cert                    # certificate for tls connectivity
    tls: true                                   # enforce tls connectivity
    providerID: flagd-weatherapp-sidecar        # identifier for this connection 
    selector: 'source=database,app=weatherapp'  # flag filtering options
```

## Sidecar configurations

`FeatureFlagSource` further allows to provide configurations to the injected flagd sidecar.
Table given below is non-exhaustive list of overriding options,

| Configuration | Explanation                   | Default                    |
|---------------|-------------------------------|----------------------------|
| port          | Flag evaluation endpoint port | 8013                       |
| metricsPort   | Metrics port                  | 8014                       |
| evaluator     | Evaluator to use              | json                       |
| image         | flagD image                   | ghcr.io/open-feature/flagd |
| tag           | flagD image tag               | Latest tag                 |
| probesEnabled | Enable/Disable health probes  | true                       |

## Merging of configurations

The annotation value is a comma separated list of values following one of two patterns: `{NAME}` or `{NAMESPACE}/{NAME}`.
If no namespace is provided, it is assumed that the CR is within the same namespace as the deployed pod, for example:

```yaml
    metadata:
        namespace: test-ns
        annotations:
            openfeature.dev/enabled: "true"
            openfeature.dev/featureflagsource: "config-A, test-ns-2/config-B"
```

In this example, 2 CRs are being used to configure the injected container (by default the operator uses the `flagd:main` image), `config-A` (which is assumed to be in the namespace `test-ns`) and `config-B` from the `test-ns-2` namespace, with `config-B` taking precedence in the configuration merge.

The `FeatureFlagSource` version `v1alpha3` CRD defines a CR with the following example structure, the documentation for this CRD can be found [here](https://github.com/open-feature/open-feature-operator/blob/main/docs/crds.md#featureflagsource):

```yaml
apiVersion: core.openfeature.dev/v1alpha3
kind: FeatureFlagSource
metadata:
    name: flag-source-sample
spec:
    metricsPort: 8080
    Port: 80
    evaluator: json
    image: my-custom-sidecar-image
    defaultSyncProvider: filepath
    tag: main
    sources:
    - source: namespace/name
      provider: kubernetes
    - source: namespace/name2
    - source: not-a-real-host.com
      provider: http
    envVars:
    - name: MY_ENV_VAR
      value: my-env-value
    probesEnabled: true
    debugLogging: false
```

The relevant `FeatureFlagSources` are passed to the operator by setting the `openfeature.dev/featureflagsource` annotation, and is responsible for providing the full configuration of the injected sidecar.

## Configuration Merging

When multiple `FeatureFlagSources` are provided, the configurations are merged. The last `CR` takes precedence over the first, with any configuration from the deprecated `FlagDSpec` field of the `FeatureFlag` CRD taking the lowest priority.

```mermaid
flowchart LR
    FeatureFlagSource-values  -->|highest priority| environment-variables -->|lowest priority| defaults
```

An example of this behavior:

```yaml
    metadata:
        annotations:
            openfeature.dev/enabled: "true"
            openfeature.dev/featureflagsource:"config-A, config-B"
```

Config-A:

```yaml
apiVersion: core.openfeature.dev/v1alpha2
kind: FeatureFlagSource
metadata:
    name: config-A
spec:
    metricsPort: 8080
    tag: latest
```

Config-B:

```yaml
apiVersion: core.openfeature.dev/v1alpha2
kind: FeatureFlagSource
metadata:
    name: config-B
spec:
    port: 8000
    tag: main
```

Results in the following configuration:

```yaml
spec:
    metricsPort: 8080
    port: 8000
    tag: main
```
