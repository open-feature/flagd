---
# Valid statuses: draft | proposed | rejected | accepted | superseded
status: draft
author: Maks Osowski (@cupofcat)
created: 2025-08-21
updated: 2025-09-01
---

# Harden Hashing Consistency And Add Support For Non-string Attributes in Fractional Evaluation

This proposal aims to enhance the `fractional` operator to:
1. Explicitly ensure hashes are consistent across all providers and platforms.
2. Support non-string values as the hashing input (i.e., the randomization unit).

Currently, all inputs are coerced to strings before hashing, which, in some reare cases, can lead to inconsistent bucketing across different provider implementations (e.g. Java provider running on a non UTF-8 platform). By this change, the targeting attributes of various types will be supported and will always be explicitly encoded in a consistent, language and platform independent, manner in every provider.

This change will be implemented in a *mostly* backward-compatible manner, preserving existing bucketing for all string-based inputs where UTF-8 was used (great majority of cases). There might be some rebucketing for users that were using certain providers on platforms with non UTF-8 locale and a provider that was enforcing UTF-8 encodings for strings.

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

1. **Implicitly:** Providing a string-type `targetingKey` in the evaluation context, which is used if the `fractional` block only contains the variant distribution.
2. **Explicitly:** Providing an expression as the first element of the `fractional` array. Today, this expression *must* evaluate to a string (standard recommendation is to use the `"cat"` operator with `$flagd.flagKey` and `"var"`); that string will be used as hashing input, usually via murmur's `StringSum32` method.

The requirement that the input evaluates to a string has two main drawbacks:

* **Inconsistent Hashing:** Different providers (Go, PHP, Java) may encode the same string into bytes differently (e.g., UTF-8 vs UTF-16). Since hashing functions like MurmurHash3 operate on bytes, this leads to different hash results and thus different bucket assignments for the same logical input across platforms.
* **Unnecessary Coercion:** If a user wishes to bucket based on a numeric ID (e.g., `userId: 12345`), they must first explicitly cast it to a string (`"12345"`) within the flag definition using an operator like `"cat"`.

This proposal seeks to resolve these issues by allowing `fractional` to operate directly on the byte representation of non-string inputs and to explicitly encode values to bytes with deterministic encoders.

## Requirements

* Users must be able to use both string and non-string variables (e.g., integers, booleans) as the primary input for `fractional` evaluation.
* Same "value" (e.g. 57.2, "some text", true, etc) should result in the same bucket assignment no matter the language of the provider and platform used.

## Non-requirements
* This change does not need to be backward-compatible. Nevertheless, care should be taken to minimize disrputions - strings should be encoded using UTF-8 (already a default on most platforms and providers).
* Support advnaced features like salting non-string types in JSON directly (that will be a separate ADR)
* Bucketing improvements (that will be a separate ADR)

## Considered Options

1. **Proposed:** *Type-Aware Hashing:* Extend the current behavior to support non-string types as first arguments to `fractional` and in `targetingKey`.
2. *New Operator:* Introduce a new operator, such as `"bytesVar"`, to explicitly signal that the variable's raw bytes should be hashed.
3. *Operator Overloading:* Reuse an existing operator (e.g., `"merge"`) or structure (e.g., providing a list) to imply byte-based hashing.

Option 1 was chosen for its ergonomics and zero-impact on existing schemas. Option 2 adds unnecessary complexity to the flag definition language, and Option 3 creates confusing and non-obvious semantics.

## Proposal

We will modify the evaluation logic for the `fractional` operator.

When inspecting the first element of the `fractional` array:

1. If the first element in `fractional` evaluates to a non-array type then deterministically encode it to a well defined byte array (using UTF-8 for strings) and hash the bytes.
2. Otherwise, if `targetingKey` is a string, concatenate `flagKey` and `targetingKey`, encode to UTF-8 and hash that (current behavior).
3. Otherwise, if `targetingKey` is non-string, create a 2 element array of [`flagKey`, `targetingKey`] and hash that.

```json
// Will use the new logic
"fractional": [
  {
    "var": "my-non-string-var"
  },
  ["a", 50], ...
]

// Will use the new logic but in a mostly backward-copmpatible way
"fractional": [
  {
    "cat": [{"var" : "$flagd.flagKey"}, {"var" : "some-var"}]
  },
  ["a", 50], ...
]

// Will use the targetingKey
"fractional": [
  ["a", 50], ...
]

// Will use targetingKey
"fractional": [
  {
    "merge": [{"var" : "evaluates-to-some-variant-name"}, {"var" : "evaluates-to-some-int"}]
  },
  ["a", 50], ...
]
```

### API changes

There are **no** changes to the flagd JSON schema. The change is purely semantic, affecting the evaluation logic within providers.

### Consequences

* Good, because any variable can be used for hashing
* Good, because it avoids unnecessary casting
* Bad, because some users will experience rebucketing, if they are on a combination of provider and platform that encodes strings differently than UTF-8

### Timeline

Pre flagd 1.0.

## More Information

Today, flagd recommends to salt the variable with flagKey directly in the `fractional` logic, using the `"cat"` operator. This will not be possible for non-string types. Advanced features like that will be considered in a separate ADR.

### Implementation considerations

The details of how to achieve the requirements of this ADR are left as the implementation detail at the discretion of contributors. However, one option worth considering is to use [CBOR](https://cbor.io/). CBOR libraries in each language ensure that the same values get the same byte encoding and murmur3 libraries will ensure that same byte arrays will get the same hash. This works across any type, even for strings.

That way we have langauge-agnostic, stable, and consistent bucketing.

CBOR is specified in an Internet Standard RFC by IETF, which should mean this stays stable for foreseeable future.
