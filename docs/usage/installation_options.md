# Get Started

## Installation

There are many ways to get started with flagd.
Choose the method that best serves your requirements to get started.

### Docker

1. `docker pull ghcr.io/open-feature/flagd:latest`

### Homebrew 🍺

1. `brew install flagd`

### Go binary

1. Install Go 1.20 or above
1. run `go install github.com/open-feature/flagd/flagd@latest`

### Release binary

1. Download pre-built binaries from <https://github.com/open-feature/flagd/releases>

### Systemd service

Documentation for installing flagd as a systemd service can be found [here](../other_resources/systemd_service.md)

### Open Feature Operator

The OpenFeature Operator is a Kubernetes native operator that allows you to expose feature flags to your applications.
It injects a flagD sidecar into your pod and allows you to poll the flagD server for feature flags in a variety of ways.
To get started with the operator, view the project here: <https://github.com/open-feature/open-feature-operator>

## Next Steps

The documentation in the following pages will help you to correctly configure your flagd service, as well as create and evaluate your own custom flags, either using curl or one of the OpenFeature language specific providers.

- [Configuring flagd](../configuration/configuration.md)
- [Creating custom flag definitions](../configuration/flag_configuration.md)
- [Evaluating flag values using a flagd provider](../usage/flagd_providers.md)
- [Evaluating flag values using curl](../usage/evaluation_examples.md)
