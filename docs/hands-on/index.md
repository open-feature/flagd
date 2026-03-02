---
description: flagd hands-on guide
---

# flagd Hands-On

This section walks you through the basics of **JsonLogic** as it's used in flagd's targeting rules, step by step.

flagd uses a modified version of [JsonLogic](https://jsonlogic.com/) to evaluate targeting rules.
If you're new to JsonLogic, it can look a bit alien at first — but it follows a simple, consistent pattern.

## What you'll learn

1. [**Why JsonLogic?**](./why-jsonlogic.md) — Why flagd uses JsonLogic for targeting rules
2. [**The Shape of Operations**](./shape-of-operations.md) — How every JsonLogic operation is structured
3. [**Basic Operators**](./basic-operators.md) — Using `if`, comparisons, and logic
4. [**Variables with `var`**](./variables.md) — Pulling values from the evaluation context
5. [**Custom Operations**](./custom-operations.md) — flagd-specific operators like `fractional` and timestamp comparisons

Work through them in order for the best experience.
