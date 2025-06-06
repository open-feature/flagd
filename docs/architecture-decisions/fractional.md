---
status: draft
author: @toddbaert
created: 2025-06-06
updated: 2025-06-06
---

# Fractional Operator

The fractional operator enables deterministic, fractional feature flag distribution.

## Background

Nearly all feature flag systems require pseudorandom assignment support to facilitate key use cases, including experimentation and fractional progressive rollouts.
Since flagd seeks to implement a full feature flag evaluation engine, such a feature is required.

## Requirements

- **Deterministic**: must be consistent given the same input (so users aren't re-assigned with each page view, for example)
- **Performant**: must be quick; we want "predictable randomness", but with a relatively low performance cost
- **Ease of use**: must be easy to use and understand for basic use-cases
- **Customization**: must support customization, such as specifying a particular context attribute to "bucket" on
- **Stability**: adding new variants should result in new assignments for as small a section of the audience as possible
- **Strong avalanche effect**: slight input changes should result in relatively high chance of differential bucket assignment

## Considered Options

- We considered various "more common" hash algos, such as `sha1` and `md5`, but they were frequently slower than `Murmur3`, and didn't offer better performance for our purposes
- Initially we required weights to sum to 100, but we've since revoked that requirement

## Proposal

### MurmurHash3 + numeric weights + optional targeting-key-based bucketing value

#### The fractional operator mechanism

The fractional operator facilitates **deterministic A/B testing and gradual rollouts** through a custom JSONLogic extension introduced in flagd version 0.6.4+.
This operator splits feature flag variants into "buckets", based the `targetingKey` (or another optionally specified key), ensuring users consistently receive the same variant across sessions through sticky evaluation.

The core algorithm involves four steps: extracting a bucketing property from the evaluation context, hashing this value using MurmurHash3, mapping the hash to a [0, 100] range, and selecting variants based on cumulative weight thresholds.
This approach guarantees that identical inputs always produce identical outputs (excepting the case of rules involving the `$flag.timestamp`), which is crucial for maintaining a consistent user experience.

#### MurmurHash3: The chosen algorithm

flagd specifically employs **MurmurHash3 (32-bit variant)** for its fractional operator, prioritizing performance and distribution quality over cryptographic security.
This non-cryptographic hash function provides excellent performance and good avalanche properties (small input changes produce dramatically different outputs) while maintaining deterministic behavior essential for sticky evaluations.
Its wide language implementation ensures identical results across different flagd providers, no matter the language in question.

#### Bucketing value

The bucking value is an optional first value to the operator (it may be a JSONLogic expression, other than an array).
This allows enables targeting based on arbitrary attributes (individual users, companies/tenants, etc).
If not specified, the bucketing value is a JSONLogic expression concatenating the `$flagd.flagKey` and the extracted [targeting key](https://openfeature.dev/specification/glossary/#targeting-key) (`targetingKey`) from the context (the inclusion of the flag key prevents users from landing in the same "bucket index" for all flags with the same number of buckets).
If the bucking value does not resolve to a string, or the `targeting key` is undefined, the evaluation is considered erroneous.

```json
// Default bucketing value
{
  "cat": [
    {"var": "$flagd.flagKey"},
    {"var": "targetingKey"}
  ]
}
```

#### Bucketing strategy implementation

After retrieving the bucketing value, and hashing it to a [0, 100] range, the algorithm iterates through variants, accumulating their relative weights until finding the bucket containing the hash value.

```go
// Simplified implementation structure
hashValue := murmur3Hash(bucketingValue) % 100
currentWeight := 0
for _, distribution := range variants {
    currentWeight += (distribution.weight * 100) / sumOfWeights
    if hashValue < currentWeight {
        return distribution.variant
    }
}
```

This approach supports flexible weight ratios; weights of [25, 50, 25] translate to 25%, 50%, and 25% distribution respectively as do [1, 2, 1].
It's worth noting that the maximum bucket resolution is 1/100, meaning that the maximum ratio between variant distributions is 1:99 (ie: a weight distribution of [1, 100000] behaves the same as [1, 100]).

#### Format flexibility: Shorthand vs longhand

flagd provides two syntactic options for defining fractional distributions, balancing simplicity with precision. **Shorthand format** enables equal distribution by specifying variants as single-element arrays (in this case, an equal weight of 1 is automatically assumed):

```json
{
  "fractional": [
    ["red"],
    ["blue"],
    ["green"]
  ]
}
```

**Longhand format** allows precise weight control through two-element arrays:

Note that in this example, we've also specified a custom bucketing value.

```json
{
  "fractional": [
    { "var": "email" },
    ["red", 50],
    ["blue", 20],
    ["green", 30]
  ]
}
```
### Consequences

* Good, because Murmur3 is fast, has good avalanche properties, and we don't need "cryptographic" randomness
* Good, because we have flexibility but also simple shorthand
* Good, because our bucketing algorithm is relatively stable when new variants are added
* Bad, because we only support string bucketing values
* Bad, because we don't have bucket resolution finer than 1:99
* Bad because we don't support JSONLogic expressions within bucket definitions
