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
Please see [architecture](./architecture.md) and [deployment](./deployment.md) for more information.

---

> Why doesn't flagd support {_my desired feature_}?

Because you haven't opened a PR or created an issue!

We're always adding new functionality to flagd, and welcome additions and ideas from new contributors.
Don't hesitate to [open an issue](https://github.com/open-feature/flagd/issues)!