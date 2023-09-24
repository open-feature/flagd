# Docker

:octicons-terminal-24: Install from the command line:

```shell
docker pull ghcr.io/open-feature/flagd:latest
```

:octicons-code-square-24: Use as base image in Dockerfile:

```
FROM ghcr.io/open-feature/flagd:latest
```
## Kubernetes

flagd was designed with cloud-native paradigms in mind.
You can run it as a sidecar, or as a central service in your cluster.
If you're interested in a full-featured solution for using flagd in Kubernetes, consider the [OpenFeature operator](https://github.com/open-feature/open-feature-operator).

<!-- for now, we are just linking out to the OFO README, but eventually we should consider fully documenting OFO here. -->

---

# Binary

:fontawesome-brands-linux::fontawesome-brands-windows::fontawesome-brands-apple: Binaries are available in x86/ARM.

[Release](https://github.com/open-feature/flagd/releases)

## systemd

A systemd wrapper is available [here](https://github.com/open-feature/flagd/blob/main/systemd/flagd.service).

---

# In-Process

## :fontawesome-brands-golang: Go in-process provider

[flagd-in-process/pkg](https://pkg.go.dev/github.com/open-feature/go-sdk-contrib/providers/flagd-in-process/pkg)

```shell
go get github.com/open-feature/go-sdk-contrib/providers/flagd-in-process/pkg@v0.1.1
```

## :fontawesome-brands-java: Java in-process provider

### Maven

```xml
<dependency>
  <groupId>dev.openfeature.contrib.providers</groupId>
  <artifactId>flagd</artifactId>
</dependency>
```

### Gradle

```
implementation 'dev.openfeature.contrib.providers:flagd'
```
