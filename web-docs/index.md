# What is flagd?

_flagd_ is a _feature flag evaluation engine_.
Think of it as a ready-made, open source, OpenFeature-compliant feature flag backend system.

With flagd you can:

* modify flags in real time
* define flags of various types (boolean, string, number, JSON)
* use context-sensitive rules to target specific users or user-traits
* perform pseudorandom assignments for experimentation
* perform progressive roll-outs of new features
* aggregate flag definitions from multiple sources

It doesn't include a UI, management console or a persistence layer.
It's configurable entirely via a POSIX-style CLI.
Thanks to it's minimalism, it's _extremely flexible_; you can leverage flagd as a sidecar alongside your application, an engine running in your application process, or as a central service evaluating thousands of flags per second.

# How do I deploy flagd?

flagd is designed to fit well into a variety of infrastructures, and can run on various architectures.
It run as a separate process or directly in your application.
It's distributed as a binary, container image, and various libraries.
If you're already leveraging containers in your infrastructure, you can extend the docker image with your required configuration.
You can also run flagd as a service on a VM or a "bare-metal" host.
If you'd prefer not to run an additional process at all, you can run the flagd evaluation engine directly in your application.
No matter how you run flagd, you will need to supply it with feature flags.
The flag definitions supplied to flagd (*sources*) are monitored for changes which will be immediately reflected in flagd's evaluations.
Currently supported sources include files, HTTP endpoints, Kubernetes custom resources, and proto-compliant gRPC services.

<!-- TODO: Link to various deployment sections with grid: https://squidfunk.github.io/mkdocs-material/reference/grids -->

# How do I use flagd?

flagd is fully OpenFeature compliant.
To leverage it in your application you must use the OpenFeature SDK and flagd provider for your language.
You can configure the provider to connect to a flagd instance you deployed earlier (evaluating flags over gRPC) or use the in-process evaluation engine to do flag evaluations directly in your application.
Once you've configured the OpenFeature SDK, you can start evaluating the feature flags configured in your flagd definitions.
