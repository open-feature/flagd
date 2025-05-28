---
status: accepted
author: @toddbaert
created: 2025-05-16
updated: --
---

# Adoption of Cucumber/Gherkin for `flagd` Testing Suite

This decision document outlines the rationale behind adopting the Cucumber/Gherkin testing framework for the `flagd` project’s testing suite. The goal is to establish a clear, maintainable, and language-agnostic approach for writing integration and behavior-driven tests.

By leveraging Gherkin’s natural language syntax and Cucumber’s mature ecosystem, we aim to improve test clarity and accessibility across teams, enabling both developers and non-developers to contribute to test case development and validation.

## Background

`flagd` is an open-source feature flagging engine that forms a core part of the OpenFeature ecosystem. As such, it includes many clients (providers) written in multiple languages and it needs robust, readable, and accessible testing frameworks that allow for scalable behavior-driven testing.

Previously, test cases for `flagd` providers were written in language-specific test frameworks, which created fragmentation and limited contributions from engineers who weren’t familiar with the language in question. Furthermore, the ability to validate consistent feature flag behavior across multiple SDKs and environments became increasingly important as adoption grew, and in-process evaluation was implemented.

To address this, the engineering team investigated frameworks that would enable:

- Behavior-driven development (BDD) to validate consistent flag evaluation behavior, configuration, and provider life-cycle (connection, etc).
- High cross-language support to integrate with multiple SDKs and tools.
- Ease of use for writing, understanding, enhancing and maintaining tests.

After evaluating our options and experimenting with prototypes, we adopted Cucumber with Gherkin syntax for our testing strategy.

## Requirements

- Must be supported across a wide variety of programming languages.
- Must offer mature tooling and documentation.
- Must enable the writing of easily understandable, high-level test cases.
- Must be open source.
- Should support automated integration in CI pipelines.
- Should support parameterized and reusable test definitions.

## Considered Options

- Adoption of Cucumber/Gherkin testing framework

## Proposal

We adopted the Cucumber testing framework, using Gherkin syntax to define feature specifications and test behaviors. Gherkin offers a structured and readable DSL (domain-specific language) that enables concise expression of feature behaviors in plain English, making test scenarios accessible to both technical and non-technical contributors.

We use Cucumber’s tooling in combination with language bindings (e.g., Go, JavaScript, Python) to execute these scenarios across different environments and SDKs. Step definitions are implemented using the idiomatic tools of each language, while test scenarios remain shared and version-controlled.

### API changes

N/A – this decision does not introduce API-level changes but applies to test infrastructure and development workflow.

### Consequences

#### Pros

- Test scenarios are readable and accessible to a broad range of contributors.
- Cucumber and Gherkin are supported in most major programming languages.
- Tests are partially decoupled from the underlying implementation language.
- Parameterized and reuseable test definitions mean new validations and assertions can often be added in providers without writing any code.

#### Cons

- Adding a new framework introduces some complexity and a learning curve.
- In some cases/runtimes, debugging failed tests in Gherkin can be more difficult than traditional unit tests.

### Timeline

N/A - this is a retrospective document, timeline was not recorded.

### Open questions

- Should we enforce Gherkin for all providers?

## More Information

- [flagd Testbed Repository](https://github.com/open-feature/flagd-testbed)
- [Cucumber Documentation](https://cucumber.io/docs/)
- [Gherkin Syntax Guide](https://cucumber.io/docs/gherkin/)
- [flagd GitHub Repository](https://github.com/open-feature/flagd)
- [OpenFeature Project Overview](https://openfeature.dev/)
