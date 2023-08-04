# Docs

This directory contains all flagd documentation, see table of contents below:

## Quick Start / Basic Usage

See the [main page](../README.md) for the quick start guide.

## Installation Options

[See all flagd installation options](usage/installation_options.md)

## Copy and Paste Evaluation Options

[This page](usage/evaluation_examples.md) provides copy and paste evaluation examples.

## Flagd Configuration

[Flagd is configured via CLI arguments on startup](configuration/configuration.md), this page describes all available options.

## Flag Configuration

This document describes the syntax for feature flag JSON configurations: [Flag configuration](configuration/flag_configuration.md).

## Application Integration

Once flagd is running, your next step is to integrate it into you application. [This page](usage/flagd_providers.md) shows all available integration options (called providers) in a variety of languages.

## Targeting Rules

flagd offers a functionality called targeting rules which rely on the incoming context sent by the client during flag evaluation.

[This page](configuration/reusable_targeting_rules.md) describes how to define targeting rules.

## Fractional Evaluation

flagd supports [fractional evaluation](configuration/fractional_evaluation.md) meaning an incoming property in the context can be sub-divided at "evaluation time" into "buckets".

[This page](configuration/fractional_evaluation.md) explains the concept and describes the technical implementation in detail.

## Starts/Ends With Evaluation

flagd supports [starts/ends_with evaluation](configuration/string_comparison_evaluation.md) meaning an incoming property in the context can be used
to determine whether a certain variant should be returned based on if its value starts or ends with a certain string.

[This page](configuration/string_comparison_evaluation.md) explains the concept and describes the technical implementation in detail.

## SemVer Evaluation

flagd supports [sem_ver evaluation](configuration/sem_ver_evaluation.md) meaning an incoming property
representing a semantic version in the context can be used to determine whether a certain variant should be returned
based on if the version meets a certain criteria.

[This page](configuration/sem_ver_evaluation.md) explains the concept and describes the technical implementation in detail.

## Flag Merging

flagd can retrieve flags from multiple sources simultaneously. [This page](configuration/flag_configuration_merging.md) describes the de-duplication and merging rules that occur if multiple identical flags are found from different flag sources.

## Help

This section documents any behavior of flagd which may seem unexpected:

- [HTTP int response](./help/http_int_response.md)

## Other Resources

- [High level architecture](./other_resources/high_level_architecture.md)
- [Creating providers](./other_resources/creating_providers.md)
- [Caching](./other_resources/caching.md)
- [Snap](./other_resources/snap.md)
- [Systemd service](./other_resources/systemd_service.md)

## Still Stuck?

[Speak to the OpenFeature community](https://openfeature.dev/community) and someone will help.
