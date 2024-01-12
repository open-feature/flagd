---
title: OpenFeature Providers
description: Overview of the available flagd providers compatible with OpenFeature.
---

flagd was built from the ground up to be [Openfeature-compliant](../concepts/feature-flagging.md#openfeature-compliance).
To use it in your application, you must use the [OpenFeature SDK](https://openfeature.dev/docs/reference/technologies/) for your language, along with the associated OpenFeature _provider_.
For more information about Openfeature providers, see the [OpenFeature documentation](https://openfeature.dev/docs/reference/concepts/provider).

## Providers

Providers for flagd come in two flavors: those that are built to communicate with a flagd instance (over HTTP or gRPC) and those that embed flagd's evaluation engine directly (note that some providers are capable of operating in either mode). For more information on how to deploy and use flagd, see [architecture](../architecture.md) and [installation](../installation.md).

The following table lists all the available flagd providers.

| Technology                                                    | RPC              | in-process       |
| ------------------------------------------------------------- | ---------------- | ---------------- |
| :fontawesome-brands-golang: [Go](./go.md)           | :material-check: | :material-check: |
| :fontawesome-brands-java: [Java](./java.md)         | :material-check: | :material-check: |
| :fontawesome-brands-node-js: [Node.JS](./nodejs.md) | :material-check: | :material-check: |
| :simple-php: [PHP](./php.md)                        | :material-check: | :material-close: |
| :simple-dotnet: [.NET](./dotnet.md)                 | :material-check: | :material-close: |
| :material-web: [Web](./web.md)                      | :material-check: | :material-close: |

For information on implementing a flagd provider, see the specifications for [RPC](../reference/specifications/rpc-providers.md) and [in-process](../reference/specifications/in-process-providers.md) providers.