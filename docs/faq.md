---
description: flagd faq
---

# Frequently Asked Questions

> Why do I need this? Can't I just use environment variables?

Feature flags are not environment variables.
If you need to update your flag values without restarting your application, target specific users, randomly assign values for experimentation, or perform scheduled roll-outs, you should consider using feature flags.
If the values are always static, an environment variable or static configuration may be sufficient.

For more information on feature-flagging concepts, see [feature-flagging](./concepts/feature-flagging.md).

---

> Why is it called "flagd"?

Please see [naming](./reference/naming.md).

---

> What is flagd's relationship to OpenFeature?

flagd is sub-project of OpenFeature and aims to be fully [OpenFeature-compliant](./concepts/feature-flagging.md#openfeature-compliance).

---

> How do I run flagd?

You can run flagd as a standalone application, accessible over HTTP or gRPC, or you can embed it into your application.
Please see [architecture](./architecture.md) and [installation](./installation.md) for more information.
For quick command examples, see the [cheat sheet](./reference/cheat-sheet.md).

---

> How can I access the SBOM for flagd?

SBOMs for the flagd binary are available as assets on the [GitHub release page](https://github.com/open-feature/flagd/releases).
Container SBOMs can be inspected using the Docker CLI.

An example of inspecting the SBOM for the latest flagd `linux/amd64` container image:

```shell
docker buildx imagetools inspect ghcr.io/open-feature/flagd:latest \
    --format '{{ json (index .SBOM "linux/amd64").SPDX }}'
```

An example of inspecting the SBOM for the latest flagd `linux/arm64` container image:

```shell
docker buildx imagetools inspect ghcr.io/open-feature/flagd:latest \
    --format '{{ json (index .SBOM "linux/arm64").SPDX }}'
```

---

> Why doesn't flagd support {_my desired feature_}?

Because you haven't opened a PR or created an issue!

We're always adding new functionality to flagd, and welcome additions and ideas from new contributors.
Don't hesitate to [open an issue](https://github.com/open-feature/flagd/issues)!
