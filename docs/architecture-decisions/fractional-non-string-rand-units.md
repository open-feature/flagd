---
# Valid statuses: draft | proposed | rejected | accepted | superseded
status: draft
author: Maks Osowski (@cupofcat)
created: 2025-08-21
updated: 2025-08-21
---

# Support Non-String Inputs for Fractional Bucketing

This proposal aims to enhance the `fractional` operator to support non-string values as the hashing input (i.e., the randomization unit). Currently, all inputs are coerced to strings before hashing, which can lead to inconsistent bucketing across different provider implementations. By this change, if the input value is not a string, flagd providers will hash its raw byte representation directly.

This change will be implemented in a backward-compatible manner, preserving existing bucketing for all string-based inputs.

```json
"fractional": [
  {
     // This will now work for non-string types
     "var": "my-non-string-var"
  },
  ["a", 50],
  ["b", 50]
]
```

## Background

The `fractional` operator in flagd determines bucket allocation (e.g., for percentage rollouts) by hashing an input value. Currently, there are two primary methods for providing this input:

1. **Implicitly:** Providing a `targetingKey` in the evaluation context, which is used if the `fractional` block only contains the variant distribution.
2. **Explicitly:** Providing an expression as the first element of the `fractional` array. This expression *must* evaluate to a string (standard recommendation is to use the `"cat"` operator with `$flagd.flagKey` and `"var"`); that string will be used as hashing input, usually via murmur's `StringSum32` method.

The explicit method (2) is the focus of this proposal. The requirement that the input evaluates to a string has two main drawbacks:

* **Inconsistent Hashing:** Different providers (Go, JS, etc.) may encode the same string into bytes differently (e.g., UTF-8 vs UTF-16). Since hashing functions like MurmurHash3 operate on bytes, this leads to different hash results and thus different bucket assignments for the same logical input across platforms.
* **Unnecessary Coercion:** If a user wishes to bucket based on a numeric ID (e.g., `userId: 12345`), they must first explicitly cast it to a string (`"12345"`) within the flag definition using an operator like `"cat"`.

This proposal seeks to resolve these issues by allowing `fractional` to operate directly on the byte representation of non-string inputs.

## Requirements

* Users must be able to use non-string variables (e.g., integers, booleans) as the primary input for `fractional` evaluation.
* The change must be backward-compatible. Existing flag configurations that use string inputs must continue to bucket users identically, with no re-bucketing.

## Considered Options

1. **Proposed: Type-Aware Hashing:** If the first element in `fractional` evaluates to a non-string type, hash its byte representation directly. If it's a string, use the existing string-hashing logic.
2. **New Operator:** Introduce a new operator, such as `"bytesVar"`, to explicitly signal that the variable's raw bytes should be hashed.
3. **Operator Overloading:** Reuse an existing operator (e.g., `"merge"`) or structure (e.g., providing a list) to imply byte-based hashing.

Option 1 was chosen for its ergonomics and zero-impact on existing schemas. Option 2 adds unnecessary complexity to the flag definition language, and Option 3 creates confusing and non-obvious semantics.

## Proposal

We will modify the evaluation logic for the `fractional` operator.

When evaluating the first element of the `fractional` array:

1. **If the resolved value is a string:** The existing logic will be preserved. Providers will continue to use their respective language's MurmurHash3 `StringSumN` (or equivalent) function. This guarantees backward compatibility.
2. **If the resolved value is a non-string type (e.g., integer, float, boolean):** Providers will hash the standard byte-representation of the value using a MurmurHash3 `SumN` (or equivalent) function.

This approach relies on the strong typing within our providers (e.g., the JS provider uses TypeScript, Go is strongly typed), allowing them to reliably distinguish between string and non-string types at evaluation time.

### API changes

There are **no** changes to the flagd JSON schema. The change is purely semantic, affecting the evaluation logic within providers.

### Consequences

* Good, because any variable can be used for hashing
* Good, because it avoids unnecessary casting
* Bad, because strings will still rely on specific of string encoding for the provider language and use `"StringSumN"` methods

### Timeline

Pre flagd 1.0.

### Open questions

* Should we also treat strings in the same way as other types? (this might result in rebucketing after the change is launched)
* Should we ensure all the providers use prescribed string encoding for all strings (e.g. UTF8)?
* Should we also change the behavior when `targetingKey` is used?

## More Information

Today, flagd recommends to salt the variable with flagKey directly in the `fractional` logic, using the `"cat"` operator. This will not be possible for non-string types. As a separate ADR we can consider an improvement where either a new custom operator or some existing operator can be used to salt variables directly in `fractional` logic, without casting to strings.
