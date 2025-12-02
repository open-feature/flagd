---
# Valid statuses: draft | proposed | rejected | accepted | superseded
status: draft
author: Maks Osowski (@cupofcat)
created: 2025-08-21
updated: 2025-09-02
---

# Harden Hashing Consistency And Add Support For Non-string Attributes in Fractional Evaluation

This proposal aims to enhance the `fractional` operator to:

1. Explicitly ensure hashes are consistent across all providers and platforms.
2. Support non-string values as the hashing input (i.e., the randomization unit).

Currently, all inputs are coerced to strings before hashing, which, in some rare cases, can lead to inconsistent bucketing across different provider implementations (e.g. Java provider running on a non UTF-8 platform). With this change, the targeting attributes of various types will be supported and will always be explicitly encoded in a consistent, language- and platform-independent manner in every provider

This change will be backward-compatible in terms of flags schema but will be a breaking behavioral change for 100% of the users due to rebucketing.

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

### 1. Users must be able to use both string and non-string variables (e.g., integers, booleans) as the primary input for `fractional` evaluation

### 2. Same "value" (e.g. 57.2, "some text", true, etc) should result in the same bucket assignment no matter the language of the provider and platform used

Please note:

* some languages (e.g. Python) don't necessarily have standard types by default (e.g. int32 vs int64).
* [OpenFeature spec 312](https://openfeature.dev/specification/sections/evaluation-context/#requirement-312) dictates that evaluation context needs to support `boolean` | `string` | `number` | `structure` | `datetime` types.
* JSON supports 6 fundamental types: `boolean` | `string` | `number` | `object` | `array` | `null`

As such, the encodings for the following types as first argument (either as literals or results of evaluation) will be standardized:

1. boolean
2. string
3. integer (any integer number, Python style)
4. float (any floating point number, Python style)
5. object (structure / map)
6. datetime
7. null

**array / sequence** will be explicitly not supported as the first argument in fractional so it's possible to distinguish between hashing input and variant bucket. Nevertheless, it can be a part of object type and its encoding needs to be standardized as well.

## Non-requirements

* This change does not need to be backward-compatible.
* Support advanced features like salting non-string types in JSON directly (that will be a separate ADR).
* Bucketing improvements (that will be a separate ADR).

## Considered Options

1. **Proposed:** *Type-Aware Hashing:* Extend the current behavior to support non-string types as first arguments to `fractional`.
2. *New Operator:* Introduce a new operator, such as `"bytesVar"`, to explicitly signal that the variable's raw bytes should be hashed.
3. *Operator Overloading:* Reuse an existing operator (e.g., `"merge"`) or structure (e.g., providing a list) to imply byte-based hashing.

Option 1 was chosen for its ergonomics and zero-impact on existing schemas. Option 2 adds unnecessary complexity to the flag definition language, and Option 3 creates confusing and non-obvious semantics.

## Proposal

We will modify the evaluation logic for the `fractional` operator.

When inspecting the first element of the `fractional` array:

1. If the first element in `fractional` evaluates to a non-array type then deterministically encode it to a well defined byte array and hash the bytes.
2. Otherwise, if `targetingKey` is a string, build a 2-elements array of `flagKey` and `targetingKey`, deterministically encode that and hash (**NOTE:** This is different than string concatenation used today).
3. Otherwise, if `targetingKey` is non-string, report an error and return nil (as this breaks the [OpenFeature spec](https://openfeature.dev/specification/glossary/#targeting-key)).
4. Otherwise, if `targetingKey` is missing, report an error and return nil

```json
// Will use the new logic
"fractional": [
  {
    "var": "my-non-string-var"
  },
  ["a", 50], ...
]

// Will use new logic
"fractional": [
  {
    "cat": [{"var" : "$flagd.flagKey"}, {"var" : "some-var"}]
  },
  ["a", 50], ...
]

// Will use targetingKey
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

### Deterministic and consistent byte encodings

To meet requirement (2) [RFC 8949 Concise Binary Object Representation (CBOR)](https://www.rfc-editor.org/rfc/rfc8949.html) will be used to decide on byte encodings.

* `boolean` is major type 7
* `null` is major type 7
* `string` is major type 3
* `integer`:
    * `unsigned integer` is major type 0
    * `negative integer` is major type 1
* `float` is major type 7
* `map` (object, structure, dict) is major type 5
* `array` (list, sequence) is major type 4
* `datetime` is converted to POSIX epoch time (including fractional seconds for sub-second precision) and CBOR Tag 1 is used

**ATTENTION: When encoding strings, CBOR appends the size of the encoding in first bytes. As such, even though the actual encoding of the string is still UTF-8, the resulting byte array will differ from raw UTF-8 encoding. As such, after this change, all hashes will change, which will result in rebucketing.**

Additionally, it is required to use [4.2.1. Core Deterministic Encoding Requirements](https://www.rfc-editor.org/rfc/rfc8949.html#section-4.2.1) (which includes Preferred Serizalization), to ensure:

1. **Map Key Ordering**: Implementations must strictly adhere to the requirement that keys in maps (objects/structures) must be sorted using bytewise lexicographic order of their deterministic encodings.
2. **Preferred Serialization (Numbers)**: CBOR mandates using the shortest possible encoding. Providers must ensure consistency, especially between integer and float representations, and across different precisions. For example, if a value fits within a 32-bit float, it must be used instead of a 64-bit float, regardless of the native type in the provider's language.

### API changes

There are **no** changes to the flagd JSON schema. The change is purely semantic, affecting the evaluation logic within providers.

### Consequences

* Good, because any variable can be used for hashing.
* Good, because it avoids unnecessary casting.
* Bad, because all of the users will experience rebucketing.

### Timeline

Prior to flagd 1.0 launch.

## More Information

Today, flagd recommends salting the variable with flagKey directly in the `fractional` logic, using the `"cat"` operator. This will not be possible for non-string types. Advanced features like that will be considered in a separate ADR.

Salting of the string types will continue to be possible using the `"cat"` operator as it is built directly into JSON Logic.

### Testing considerations

As part of implementation of this ADR, the current Gherkin suite will need to be updated to ensure more in-depth testing of consistency (e.g. by looking at the distribution of buckets for many samples), as well as support for many new types.
