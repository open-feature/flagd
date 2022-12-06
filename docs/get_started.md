### Get Started

There are many ways to get started with flagd. Choose the method that best serves your requirements to get started.

### Docker

1. `docker pull ghcr.io/open-feature/flagd:latest`

### Go binary

1. Install Go 1.18 or above
1. run `go install github.com/open-feature/flagd@latest`

### Release binary

1. Download pre-built binaries from https://github.com/open-feature/flagd/releases

### Open Feature Operator
The OpenFeature Operator is a Kubernetes native operator that allows you to expose feature flags to your applications. It injects a flagD sidecar into your pod and allows you to poll the flagD server for feature flags in a variety of ways.
To get started with the operator, view the project here: https://github.com/open-feature/open-feature-operator