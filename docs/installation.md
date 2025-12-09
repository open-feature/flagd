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

## Kubernetes

flagd was designed with cloud-native paradigms in mind.
You can run it as a sidecar, or as a central service in your cluster.
If you're interested in a full-featured solution for using flagd in Kubernetes, consider the [OpenFeature operator](https://github.com/open-feature/open-feature-operator).

For more information, see [OpenFeature Operator](./reference/openfeature-operator/overview.md).

You can also choose to run a Kubernetes service in front of a deploymnent with multiple Flagd pods connecting to the same datasource. However, if doing so, be aware that synchronization is not instant. The service may return different values after a change until all pods have synchronized with the data source. To be fair, this is typically a small amount of time.

---

## Binary

:fontawesome-brands-linux::fontawesome-brands-windows::fontawesome-brands-apple: Binaries are available in x86/ARM.

[Releases](https://github.com/open-feature/flagd/releases)

!!! note

    Installing flagd with `go install github.com/open-feature/flagd/flagd@latest` is not recommended, because the module depends on potentially unpublished, local workspace modules.
    Please use one of the official, versioned binary releases above.

### systemd

A systemd wrapper is available [here](https://github.com/open-feature/flagd/blob/main/systemd/flagd.service).

### Homebrew

```shell
brew install flagd
```

## Summary

Once flagd is installed, you can start using it within your application.
Check out the [OpenFeature providers page](./providers/index.md) to learn more.
