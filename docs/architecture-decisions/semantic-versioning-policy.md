---
# Valid statuses: draft | proposed | rejected | accepted | superseded
status: proposed
author: Maks Osowski (@cupofcat)
created: 2026-01-22
---

# flagd Ecosystem Semantic Versioning Policy

# Introduction

This document outlines the versioning policy for the flagd ecosystem, including the flagd binary, and language-specific providers. The goal is to establish a predictable, stable, and transparent contract between the components and their consumers and ensure that users can confidently adopt updates while allowing the ecosystem to evolve.

The flagd ecosystem consists of distinct components that maintain their own versioning tracks but adhere to a unified policy. The scope for this policy is flagd binaries, and providers. Various API surfaces (e.g. Evaluation API, Sync API, OFREP) or schemas (e.g. Flag Definition schema) are out of scope.

# Legend and definitions

**X.Y.Z** refers to the semantic version of a release component (**X** is the major version, **Y** is the minor version, **Z** is the patch version).

## Major versions (X.)

For example: 1.3.4 → 2.0.0

flagd community is not anticipating new major version releases (2.0.0) of any of the components. These are reserved for fundamental paradigm shifts or massive architectural changes. For the foreseeable future the ecosystem aims to maintain stability within the v1.x.y series.

## Minor versions (Y.)

For example: 1.3.4 → 1.4.0

Minor versions of components include new features and may include breaking changes or introduce incompatibilities. Breaking changes to existing functionality will only be released if adhering to the Deprecation Policy.

## Patch versions (Z.)

For example: 1.3.4 → 1.3.5

Patch releases are intended for critical bug fixes to the latest minor version, such as addressing security vulnerabilities, fixes to problems affecting a large number of users, and severe problems with no workaround.

They should not contain miscellaneous feature additions or improvements, and especially no incompatibilities should be introduced between patch versions of the same minor version.

Dependencies, such as JSON Logic, should also not be changed unless absolutely necessary, and also just to fix critical bugs (so, at most patch version changes, not new major nor minor versions).

# Release Versioning

The flagd ecosystem follows a lightweight, **feature-driven release process**. Releases are driven by feature readiness and stability needs rather than a rigid calendar schedule. This ensures that features reach users quickly without imposing unnecessary process overhead on maintainers.

## Branching Strategy

* **Main Branch (`main`):** The single source of truth. The main branch is expected to be in a stable, deployable state at all times.  
* **Tags:** Releases are strictly defined by Git tags (e.g., `v1.2.0`) created directly on the main branch.  
* **Maintenance Branches (Exception Only):** Long-lived release branches (e.g., `release-1.2`) are **not** created by default. They are only utilized if a critical patch is required for an older version while the main branch has already progressed with breaking changes.

## Release Cadence and Process

### Minor Releases (X.Y.0)

* **Trigger:** Released when a significant set of new features, improvements, or non-breaking changes have been merged to `main` and validated.  
* **Cadence:** Ad-hoc. There are no mandatory waiting periods or fixed release windows. However, the set of changes going into each release is **planned ahead** of the release and the criteria for “ready” are set ahead of time so all the languages can stay aligned.  
* **Process:** A tag is pushed to the main branch, triggering the release pipeline to build artifacts and publish notes.

### Patch Releases (X.Y.Z)

* **Trigger:** Released to fix critical bugs, security vulnerabilities, or regressions found in the current minor version.  
* **Strategy (Roll-Forward):** The preferred method is to merge the fix to `main` and immediately cut the next patch release (e.g., `v1.2.1`). This minimizes branch management overhead.  
* **Strategy (Backporting):** Backporting via a temporary branch is an exception, reserved only for critical fixes required for a specific version where rolling forward is not an option (e.g., the main branch contains breaking changes for the next minor version).

### Pre-Releases (Alpha/Beta/RC)

* **Trigger:** **Optional**. Pre-releases are only utilized when a release contains significant architectural changes or high-risk features that require wider community testing before broad adoption.  
* **Format:** `X.Y.Z-rc.N` (e.g., `v1.5.0-rc.1`).  
* **Process:** Tagged directly from the main branch. If issues are found, fixes are merged to main, and a subsequent RC is tagged. There are no mandatory alpha/beta cycles to streamline the process.

## Artifact Integrity

Released artifacts (binaries, container images) must be immutable. Hashes (e.g., SHA-256) of all released artifacts must be published with the release notes. A specific version tag must always resolve to the same artifact hash.

# Upgrades and SLOs

We expect users to stay reasonably up-to-date with the versions of flagd components they use in production, but understand that it may take time to upgrade, especially for production-critical components.

We expect users to be running approximately the latest patch release of a given minor release; we often include critical bug fixes in patch releases, and so encourage users to upgrade as soon as possible.

We expect to “support” 3 minor releases at a time. "Support" means we expect users to be running that version in production, and we strive to port fixes back into the supported versions. For example, when v1.3 comes out, v1.0 will no longer be “supported” but v1.1 would be expected to contain critical bug fixes discovered when v1.3 is the latest version. Basically, that means that the reasonable response to the question "my v1.0 flagd Go provider isn't working," is, "you should probably upgrade it, (and probably should have some time ago)".

Being an OSS project, we do **NOT** offer any SLA on resolving issues.

We have a “best-effort SLO” of:

* addressing CVEs within 14 days of disclosure  
* addressing severe bugs within 31 days of reporting

# Component Skew

## flagd Providers and APIs

Providers communicate with the Sync and Evaluation API endpoints. Over time the API messages might evolve within the same major version to support new functionalities in backward compatible manners. To ensure system stability, we define the following policy.

**API Compatibility**: Providers expect the API endpoints they are configured with to support the version of the provider. There are no guarantees of forward or backward compatibility in providers. It’s the responsibility of the API owners to not break the clients that use their APIs.

## flagd Providers and OpenFeature SDKs

Each flagd provider implementation is bound to the OpenFeature SDK of its respective language. To ensure stability, we define the following policies:

**Compatibility Declaration**: Each provider release **MUST** declare the [OpenFeature spec version](https://github.com/open-feature/spec/releases) it’s compatible with.

**Version Dependency**:

* A **patch** version upgrade of a flagd provider will **NOT** change the spec version  
* A **minor** version upgrade of a flagd provider **may** **increase** the spec version

# Cross-Provider Alignment

While released independently, providers strive for feature parity.

**Minor Version Alignment**: We aim to align minor versions across providers to represent a consistent feature set and behavior (e.g., Python provider 1.1.x and Java provider 1.1.x should offer similar configuration options).

**Patch Divergence**: Patch versions are released independently as needed for language-specific fixes and do not require alignment.

# Deprecation Policy

To evolve the ecosystem without immediate breaking changes, we employ a strict deprecation process for the flagd binary and providers.

**Announcement**: Features/Behavior must be marked and announced as deprecated in a **minor** release.

**Duration**: Deprecated features must be supported for at least **3 minor releases** or **12 months**, whichever is longer.

**Visibility**: Runtime warnings should be emitted when deprecated features are used.

**Removal**: After the deprecation period, the feature may be removed in a subsequent **minor** release (marked as a breaking change)

# Breaking Changes

## For flagd binary (daemon) and providers

* **Configuration**: Changes to the names, types, or default values of environmental variables, CLI flags (`flagd start` options), or provider constructor options (e.g., `FLAGD_CACHE`).

* **Observability**:  
  * Changes to metric names or types that move, remove, or rename existing parts of the schemas (additions, e.g. of labels, are fine).  
  * Changes to the mapping of evaluation details to OpenTelemetry feature-flag event records.  
  * Please note that the following are considered **not** breaking:  
    * Removing or changing existing logging at any level (ERROR, WARNING, INFO, etc)

* **Runtime Behavior**:  
  * Changes to startup behavior (e.g., fail-open vs. fail-close).  
  * Changes to retry-ability, idempotency, or backoff behavior.

* **Licensing**: Changes to the license texts of the artifacts.

## Additionally, for flagd providers only

* Changes to existing public interfaces (signatures, return types, or behavior for the same input/state) that break the existing client of those interfaces (note that, for example, adding a new optional configuration option to an existing interface is **not** breaking).  
* Changing the minimum supported language runtime or compiler version.  
* Changing the supported OpenFeature spec version

## Additionally, for flagd binary (daemon) only

* Changes to which API versions are exposed by default.  
* Changes to the supported URI patterns for flag sources (e.g., `file:`, `kubernetes:`, `s3:`).  
* Changes to the merge strategy or precedence when using multiple flag sources.  
* Implementing documented future default flips (e.g., the plan to default `--disable-sync-metadata` to true is a breaking change).

# Policy Rollout Checklist (Ahead of 1.0.0)

Before graduating core components to 1.0.0 and finalizing this policy:

- [ ] Battle-testing 1.0.0-rc.1 candidates extensively.  
- [ ] Extraction and independent versioning of the `flagd-core` Go library.  
- [ ] Implementation of CI checks to validate SemVer compliance (e.g., detecting accidental breaking changes in public APIs).  
- [ ] Designation of release stewards for each component.  
- [ ] Publication of the initial compatibility matrix.  
- [ ] Establish the release notes process (where to publish, what format) for minor releases