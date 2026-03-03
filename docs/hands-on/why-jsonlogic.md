---
description: Why flagd uses JsonLogic for targeting rules
---

# Why JsonLogic?

flagd needs a way to express targeting rules — the logic that decides which variant a user gets.
We chose [JsonLogic](https://jsonlogic.com/) for several reasons:

## It's just JSON

JsonLogic rules are plain JSON. This means:

- **Simplicity**: rules can be stored in files, databases, HTTP responses, or Kubernetes custom resources without any special encoding/decoding
- **Validation is straightforward**: standard JSON Schema tooling works out of the box
- **Parsing is universal**: every language and platform has a JSON parser built in, databases have traversal mechanisms, etc
- **Kubernetes-native**: JSON fits perfectly inside CRD specs (YAML or JSON), no escaping or string-embedding required; can be easily diff'd

## Implemented in many languages

JsonLogic has [implementations in 15+ languages](https://jsonlogic.com/), including Go, JavaScript, Python, Java, .NET, PHP, Ruby, and more.
This means the same targeting rules can be evaluated consistently across different services, SDKs, and platforms.

## Stateless and secure

JsonLogic expressions are **pure data**. They:

- **Cannot execute arbitrary code**: there are no function definitions, loops, or system calls
- **Have no side effects**: evaluating a rule never modifies state
- **Cannot access the filesystem or network**: rules can only operate on the data you explicitly pass in

This makes JsonLogic safe to accept from untrusted sources, safe to evaluate in any environment, and easy to reason about.

## Deterministic

Given the same input, a JsonLogic expression always produces the same output.
This makes targeting rules predictable, testable, and cacheable.

## Extensible

JsonLogic supports custom operators, which flagd uses to add feature-flag-specific functionality like [`fractional`](./custom-operations.md) splits and semantic version comparisons — while keeping the same consistent `{ "operator": [params] }` shape.

---

Now that you know *why* we use JsonLogic, let's learn *how* it works: [The Shape of Operations](./shape-of-operations.md)
