# Docs

This directory contains all flagd documentation, see table of contents below:

## Usage

There are many ways to get started with flagd, the sections below run through some simple deployment options. Once the flagd service is running flag evaluation requests can either be made via one of the language specific flagd providers, or, directly via curl.

- [Getting started](https://github.com/open-feature/flagd/blob/main/docs/usage/getting_started.md)
- [Flagd providers](https://github.com/open-feature/flagd/blob/main/docs/usage/flagd_providers.md)
- [Evaluation examples](https://github.com/open-feature/flagd/blob/main/docs/usage/evaluation_examples.md)

## Flag Configuration

Flagd is configured via CLI arguments on startup, these configuration options can be found in the flagd configuration section. The remaining sections cover the flag configurations themselves, which are JSON representations of the flag variants and targeting rules.

- [Flagd Configuration](https://github.com/open-feature/flagd/blob/main/docs/configuration/configuration.md)
- [Flag configuration](https://github.com/open-feature/flagd/blob/main/docs/configuration/flag_configuration.md)
- [Fractional evaluation](https://github.com/open-feature/flagd/blob/main/docs/configuration/fractional_evaluation.md)
- [Reusable targeting rules](https://github.com/open-feature/flagd/blob/main/docs/configuration/reusable_targeting_rules.md)

## Help

This section documents any behavior of flagd which may seem unexpected, currently covering 2 topics; why the HTTP int response is a string, and why values may be omitted from the evaluation response.

- [HTTP int response](https://github.com/open-feature/flagd/blob/main/docs/help/http_int_response.md)
- [Omitted value from evaluation response](https://github.com/open-feature/flagd/blob/main/docs/help/omitted_value_from_response.md)

## Other Resources
- [High level architecture](https://github.com/open-feature/flagd/blob/main/docs/other_resources/high_level_architecture.md)
- [Caching](https://github.com/open-feature/flagd/blob/main/docs/other_resources/caching.md)
- [Snap](https://github.com/open-feature/flagd/blob/main/docs/other_resources/snap.md)
- [Systemd service](https://github.com/open-feature/flagd/blob/main/docs/other_resources/systemd_service.md)