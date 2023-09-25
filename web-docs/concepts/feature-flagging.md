# Feature Flagging

Feature flags are a software development technique that allows teams to enable, disable or change the behavior of certain features or code paths in a product or service, without modifying the source code.

## OpenFeature Compliance

[OpenFeature](https://openfeature.dev/) is an open standard that provides a vendor-agnostic, community-driven API for feature flagging.
The flagd project is fully OpenFeature-compliant.
In fact, flagd was initially conceived as a reference implementation for an OpenFeature backend, but has become a powerful tool in its own right.
For this reason, you'll find flagd's concepts and terminology align with that of the OpenFeature project.
Within the context of an OpenFeature-compliant feature flag solution, flagd artifacts and libraries comprise the [flag management system](https://openfeature.dev/specification/glossary#flag-management-system) and [providers](https://openfeature.dev/specification/glossary#provider).
These artifacts and libraries alone won't allow you to evaluate flags in your application - you'll also need the [OpenFeature SDK](https://openfeature.dev/specification/glossary#feature-flag-sdk) for your language as well, which provides the evaluation API for application developers to use.

## Supported Feature Flagging Use-Cases

Below is a non-exhaustive table of common feature flag use-cases, and how flagd supports them:

| Use case                                  | flagd Feature                                                                                                                                                                                                                                                                    |
| ----------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| dynamic configuration                     | Flag definitions from any sync source are monitored for changes, with some syncs supporting near real time updates.                                                                                                                                                              |
| dynamic (context-sensitive) evaluation    | flagd evaluations are context sensitive. Rules can use arbitrary context attributes as inputs for flag evaluation logic.                                                                                                                                                         |
| fractional evaluation / random assignment | flagd's [fractional](../reference/custom-operations/fractional-operation.md) custom operation supports pseudorandom assignment of flag values.                                                                                                                                   |
| progressive roll-outs                     | Progressive roll-outs of new features can be accomplished by leveraging the [fractional](../reference/custom-operations/fractional-operation.md) custom operation as well as automation in your build pipeline, SCM, or infrastructure which updates the distribution over time. |
