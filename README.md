<h1 align="center">
  <img src="images/flagD.png" width="350px;" >
</h1>

<h2 align="center">A feature flag daemon with a Unix philosophy.</h4>

<p align="center">
  <a href="https://github.com/open-feature/flagd/actions">
    <img src="https://github.com/open-feature/flagd/actions/workflows/build.yaml/badge.svg" alt="Github Actions">
  </a>
  <a href="https://goreportcard.com/report/github.com/open-feature/flagd">
    <img src="https://goreportcard.com/badge/github.com/open-feature/flagd">
  </a>
  <img src="https://img.shields.io/github/go-mod/go-version/open-feature/flagd">
  <a href="https://github.com/open-feature/flagd/releases">
    <img src="https://img.shields.io/github/release/open-feature/flagd/all.svg">
  </a>
</p>

## Features

- OpenFeature compliant with providers available in many languages
- Multiple flag configuration sources including `files`, `http`, and `Kubernetes`
- Accessible over gRPC and HTTP
- Supports subscriptions to real-time flag change events
- Flexible targeting rules based on [JSON Logic](https://jsonlogic.com/)
- Lightweight daemon, with an emphasis on performance
- Native support for metrics using Prometheus

## Get started

Flagd is a simple command line tool for fetching and evaluating feature flags for services. It is designed to conform with the OpenFeature specification. To get started, follow the installation instructions on the [wiki](https://github.com/open-feature/flagd/wiki).

### The people who make flagD great 💜

<a href="https://github.com/open-feature/flagd/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=open-feature/flagd" />
</a>

## License

Apache License 2.0
