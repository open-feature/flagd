---
description: installing flagd
---

# Installation

## Docker

:octicons-terminal-24: Install from the command line:

```shell
docker pull ghcr.io/open-feature/flagd:latest
```

:octicons-code-square-24: Use as base image in Dockerfile:

```dockerfile
FROM ghcr.io/open-feature/flagd:latest
```

### Kubernetes

flagd was designed with cloud-native paradigms in mind.
You can run it as a sidecar, or as a central service in your cluster.
If you're interested in a full-featured solution for using flagd in Kubernetes, consider the [OpenFeature operator](https://github.com/open-feature/open-feature-operator).

For more information, see [OpenFeature Operator](./reference/openfeature-operator/installation.md).

---

## Binary

:fontawesome-brands-linux::fontawesome-brands-windows::fontawesome-brands-apple: Binaries are available in x86/ARM.

[Release](https://github.com/open-feature/flagd/releases)

### systemd

A systemd wrapper is available [here](https://github.com/open-feature/flagd/blob/main/systemd/flagd.service).

### Homebrew

```shell
brew install flagd
```

### Go binary

```shell
go install github.com/open-feature/flagd/flagd@latest
```

---

## flagd Providers for OpenFeature

Leverage flagd in your application using [OpenFeature](https://openfeature.dev/) and a flagd [provider](https://openfeature.dev/docs/reference/concepts/provider).
Use the table below to see what provides are available.

| Technology                                                    | RPC              | in-process       |
| ------------------------------------------------------------- | ---------------- | ---------------- |
| :fontawesome-brands-golang: [Go](./providers/go.md)           | :material-check: | :material-check: |
| :fontawesome-brands-java: [Java](./providers/java.md)         | :material-check: | :material-check: |
| :fontawesome-brands-node-js: [Node.JS](./providers/nodejs.md) | :material-check: | :material-check: |
| :simple-php: [PHP](./providers/php.md)                        | :material-check: | :material-close: |
| :simple-dotnet: [.NET](./providers/dotnet.md)                 | :material-check: | :material-close: |
| :material-web: [Web](./providers/web.md)                      | :material-check: | :material-close: |
