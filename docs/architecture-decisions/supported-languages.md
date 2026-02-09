---
# Valid statuses: draft | proposed | rejected | accepted | superseded
status: proposed
author: Maks Osowski (@cupofcat)
created: 2026-01-27
---

# flagd Ecosystem Language Support Policy

## Introduction

Reference: flagd Semantic Versioning Policy (TODO)

`flagd` providers are implemented in a wide range of programming languages. However, the maturity, feature completeness, and maintenance capacity for these languages vary.

This document outlines the **Tiered Support Model** for flagd providers. Its goal is to set clear expectations for adopters regarding stability, feature parity, and security response times, while providing a clear path for community contributions to graduate to official support.

**Source of Truth**: The Markdown file version of this policy located in `flagd` repository is the single source of truth. The flagd website and other documentation must reflect the content of this file.

## Governance and Ownership

To ensure long-term sustainability, we define specific roles and responsibilities regarding provider maintenance.

### Roles

* **Core Maintainers:** The maintainers of the flagd core project. They are responsible for vetting new language proposals, approving Tier promotions, and enforcing this policy.  
* **Maintainers of the provider (language):** Individuals or organizations specifically responsible for a language provider (e.g., the "Java Provider Owner").  
* **Community Contributors:** Developers who submit PRs but do not hold long-term maintenance responsibilities.

## Support Tiers and SLAs

Language providers are categorized into three tiers. The classification determines the level of support users can expect.

### Tier 1: Core Supported

**Definition:** Tier 1 providers are "production-ready" and are the primary focus of the ecosystem.

**Criteria:**

* **Conformance**  
  * Must pass 100% of the OpenFeature spec conformance test suite.  
  * Must pass 100% of the Gherkin provider e2e scenarios suite.  
* **Ownership:** Minimum of **2 dedicated provider (language) maintainers**.  
* **Documentation:** Complete API references and usage examples in the repo.  
* **Service Level Objectives (SLOs):** Follow the “best-effort SLO” defined in the Semantic Versioning Policy  
* **Releases:**  
  * All the tier 1 languages follow a **release train** – the Core Maintainers of flagd set a new release date and the list of features that are targeted for that release. The Maintainers of each language are responsible for ensuring that all the targeted features are implemented and tested before the release.  
  * If some language(s) will not make it the Core Maintainers can make one of the two choices:  
    * cut the feature from all the languages for that release  
    * postpone the release date for all the languages  
    * include the feature in some subset of languages as experimental, undocumented feature  
  * The main goal is to ensure feature parity per release for all tier 1 languages  
* **Benchmarking:** All new releases have performance benchmark reports attached  
* **Language and library support:**  
  * The latest versions of the providers support the current LTS version of the language, if one exists  
  * The latest version of the providers do not rely on deprecated language features or libraries

### Tier 2: Community Supported

**Definition:** Tier 2 providers are stable and usable but may lack advanced features or strict SLAs. They rely on the community for ongoing maintenance.

**Criteria:**

* **Conformance**  
  * Must pass 100% of the OpenFeature spec conformance test suite.  
  * Must pass 100% of the Gherkin provider e2e scenarios suite.  
* **Ownership:** Minimum of **1 dedicated provider (language) maintainer**.  
* **Documentation:** Might have some gaps  
* **Service Level Objectives (SLOs):** None  
* **Releases:**   
  * The numbering of the releases should match the numbering of the tier 1 releases (meaning, that the feature set is the same as in tier 1\)  
  * However, the tier 2 languages do NOT need to “catch” the release train, they can lag behind the tier 1  
* **Language and library support:** Best effort

### Tier 3: Experimental / Incubation

**Definition:** Tier 3 includes new providers under active development, proof-of-concept implementations, or providers for niche languages with low usage.

**Criteria:**

* Work in progress or lacking a dedicated owner.  
* May have incomplete API surfaces.  
* **Use at your own risk.** No guarantees of API stability or backward compatibility.  
* These providers may be archived if inactive for more than 6 months.

## Features definitions

* Core feature – a feature that must be present and fully supported in all the providers starting at a specific release (same for all languages)  
* Experimental feature – a feature that can be present but undocumented only in the subset of languages and is not tied to any specific release train  
* Language-specific feature – a feature that makes sense only in a specific language or a subset of languages; it does not need to follow the release train

When a contributor wants to contribute a feature to a specific language but it lacks sponsorship or traction to be implemented in all languages it can be released as experimental until it’s implemented across the board. If the contributor believes their feature should become a core feature they need to build traction and get buy-in / sponsorship from Core Maintainers.

## Lifecycle Management

The status of a language provider is not permanent. It reflects the current reality of the code and community.

### Adding New Languages

We welcome new language providers, but require a structured approach to prevent ecosystem fragmentation.

1. **Proposal:** Open an issue in the `flagd` repository proposing the new language.  
2. **Sponsorship:** A Core Maintainer must sponsor the addition.  
3. **Incubation:** The repository is created (or transferred) and starts at **Tier 3**.  
4. **Development:** The provider must reach a baseline of functionality before being advertised in official docs.

### Promotion (Graduation)

A provider may be promoted (e.g., Tier 2 → Tier 1\) upon request.

* **Procedure:** The Provider Owner submits a "Promotion Request" issue.  
* **Vetting:** Core Maintainers audit the codebase, check conformance tests, and verify the adherence to Tier requirements  
* **Approval:** Requires a majority vote from the Core Maintainers.

### Demotion and Deprecation

A provider may be demoted (e.g., Tier 1 → Tier 2\) or deprecated.

* **Demotion Triggers:**  
  * Loss of maintainers (dropping below the required count).  
  * Release lagging behind the latest  
* **Deprecation (Sunset):**  
  * If a Tier 2 or 3 provider has no activity or owner for **6 months**, it will be marked as **Deprecated**.  
  * Deprecated repos will be archived (read-only) after an additional **3 month** grace period if no new owner steps forward.

# Provider Development Best Practices

To ensure a cohesive ecosystem, all providers should adhere to these development guidelines:

* **Idiomatic Design:** APIs should feel natural to the language (e.g., use `Context` in Go, `Async/Await` in JS/C\#), rather than forcing a direct port of the Java or Go logic.  
* **Minimal Dependencies:** Avoid heavy dependencies. The provider should be lightweight.  
* **Generated Code:** Use `buf` or `protoc` for generating gRPC stubs from the official flagd schemas. Do not manually write protocol buffers code.  
* **CI/CD:** All providers must have a GitHub Actions pipeline that runs:  
  * Linters/Formatters.  
  * Unit Tests.  
  * The standard flagd integration/conformance tests.