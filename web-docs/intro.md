# What is flagd?

_flagd_ is a _feature flag evaluation engine_.
Think of it as a ready-made, open source, OpenFeature-compliant feature flag backend system.
It allows you to dynamically evaluate feature flags.

With flagd you can:

* modify flags in real time
* define flags of various types (boolean, string, number, JSON)
* use context-sensitive rules to target specific users or user-traits
* perform pseudorandom assignments for experimentation
* perform progressive roll-outs of new features
* aggregate flag definitions from multiple sources 

It doesn't include a UI, management console or a persistence layer.
It's configurable entirely via a POSIX-style CLI.
Thanks to it's minimalism, it's _extremely flexible_; you can run flagd as a sidecar alongside your application, or as a central service evaluating thousands of flags per second.
