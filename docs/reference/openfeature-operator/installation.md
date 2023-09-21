# Installing the OpenFeature Operator

Use the [OpenFeature Operator](https://github.com/open-feature/open-feature-operator) to install and run flagd on a Kubernetes cluster.

The operator includes flagd (no need to install flagd seperately).

## Installation

1. Install [cert-manager](https://cert-manager.io/docs/installation/) if you don't already have it on the cluster
1. Install the OpenFeature Operator (see below)
1. Define CRDs which describe the feature flag source (flagd) and the feature flags themselves
1. Annotate your deployment to enable flagd feature flagging for that pod

## Install OpenFeature operator

```bash
helm repo add openfeature https://open-feature.github.io/open-feature-operator/
helm repo update
helm upgrade --install openfeature openfeature/open-feature-operator
```

## Define feature flags

Create a namespace to house your flags:

```bash
kubectl create namespace flags
```

Next define your feature flag(s) using the [FeatureFlagConfiguration](./crds/featureflagconfiguration.md) CRD.

This example specifies one flag called `foo` which has two variants `bar` and `baz`. The `defaultVariant` is `bar`.

If this doesn't make sense, review the [concepts](../concepts/index.md) section.

```bash
kubectl apply -n flags -f - <<EOF
apiVersion: core.openfeature.dev/v1alpha2
kind: FeatureFlagConfiguration
metadata:
  name: sample-flags
spec:
  featureFlagSpec:
    flags:
      foo:
        state: "ENABLED"
        variants:
          "bar": "BAR"
          "baz": "BAZ"
        defaultVariant: "bar"
        targeting: {}
EOF
```

Next, tell the OpenFeature operator where to find flags.

Do so by creating a [FlagSourceConfiguration](./crds//flagsourceconfiguration.md) CRD.

This example specifies that the CRD called `sample-flags` (created above) can be found in the `flags` namespace and that the provider is `kubernetes`.

The `port` parameter defines the port on which the flagd API will be made available via the sidecar (more on this below).

```bash
kubectl apply -n flags -f - <<EOF
apiVersion: core.openfeature.dev/v1alpha3
kind: FlagSourceConfiguration
metadata:
  name: flag-source-configuration
spec:
  sources:
  - source: flags/sample-flags
    provider: kubernetes
  port: 8080
EOF
```

## Enable your deployment for feature flags

The operator looks for `Deployment` objects annotated with particular annotations.

- `openfeature.dev/enabled: "true"` enables this deployment for flagd
- `openfeature.dev/flagsourceconfiguration: "flags/flag-source-configuration"` makes the given feature flag sources available to this deployment

When these two annotation are added, the OpenFeature operator will inject a sidecar into your workload.

flagd will then be available via `http://localhost` the port specified in the `FlagSourceConfiguration` (eg. `8080`)

Your Deployment YAML might look like this:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: busybox-curl
spec:
  replicas: 1
  selector:
    matchLabels:
      app: my-busybox-curl-app
  template:
    metadata:
      labels:
        app: my-busybox-curl-app
      annotations:
        # here are the annotations for OpenFeature Operator
        openfeature.dev/enabled: "true"
        openfeature.dev/flagsourceconfiguration: "flags/flag-source-configuration"
    spec:
      containers:
        - name: busybox
          image: yauritux/busybox-curl:latest
          ports:
            - containerPort: 80
          args:
            - sleep
            - "30000"
```

## Pseudo-code of an application interacting with flagd sidecar

```bash
// From within the pod
curl --location 'http://localhost:8080/schema.v1.Service/ResolveString' --header 'Content-Type: application/json' --data '{ "flagKey":"foo"}'
```

In a real application, rather than `curl`, you would probably use the OpenFeature SDK with the `flagd` provider. // TODO link to a good example here.

## What does the operator do?

The operator will look for the annotations above and, when found, inject a sidecar into the relevant pods.

flagd reads the feature flag CRD(s) and makes an API endpoint available so your application can interact with flagd.
