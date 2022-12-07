# Docs

This directory contains all flagd documentation, see table of contents below:

## Usage

There are many ways to get started with flagd, the sections below run through some simple deployment options. Once the flagd service is running flag evaluation requests can either be made via one of the language specific flagd providers, or, directly via curl.

- [Getting started](./usage/getting_started.md)
- [Flagd providers](./usage/flagd_providers.md)
- [Evaluation examples](./usage/evaluation_examples.md)

## Flag Configuration

Flagd is configured via CLI arguments on startup, these configuration options can be found in the flagd configuration section. The remaining sections cover the flag configurations themselves, which are JSON representations of the flag variants and targeting rules.

- [Flagd Configuration](./configuration/configuration.md)
- [Flag configuration](./configuration/flag_configuration.md)
- [Fractional evaluation](./configuration/fractional_evaluation.md)
- [Reusable targeting rules](./configuration/reusable_targeting_rules.md)

## Help

This section documents any behavior of flagd which may seem unexpected, currently covering 2 topics; why the HTTP int response is a string, and why values may be omitted from the evaluation response.

- [HTTP int response](./help/http_int_response.md)
- [Omitted value from evaluation response](./help/omitted_value_from_response.md)

## Other Resources
- [High level architecture](./other_resources/high_level_architecture.md)
- [Caching](./other_resources/caching.md)
- [Snap](./other_resources/snap.md)
- [Systemd service](./other_resources/systemd_service.md)


