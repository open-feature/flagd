---
description: flagd resolve custom operation
---

# Resolve Operation

`resolve` ends a targeting rule with a variant, and optionally a custom reason and metadata for the branch that fired.

`"red"` and `{ "resolve": ["red"] }` select the same variant.

The variant is coerced like a plain targeting result: booleans become `"true"` or `"false"`, and the behavior of numbers is not defined.

A missing or null reason keeps `TARGETING_MATCH`.

The reason is free-form, except `ERROR`, `DISABLED`, `DEFAULT`, `FALLBACK`, `STATIC`, `UNKNOWN`, `CACHED` and `STALE`, which are reserved (see the [ADR](../../architecture-decisions/targeting-variant-metadata.md)).

The [targeting schema](../schema.md) rejects them as literals, and at runtime they're ignored, so the evaluation keeps `TARGETING_MATCH`.

Metadata is a shallow list of `[key, value]` pairs.

Values are strings, numbers, or booleans.

Rule metadata overwrites flag metadata, which overwrites flag set metadata.

`flagSetId` and `version` are reserved for flag set metadata and can't be overridden.

`resolve` is terminal.

It is valid as the whole targeting rule, as a result of `if`, or as the variant of a `fractional` bucket, as in the [example below](#fractional-split).

The [targeting schema](../schema.md) rejects `resolve` anywhere else, such as inside `==`, `cat`, or the bucketing value or weight of `fractional`, and rejects malformed literal arguments.

Like any other configuration that doesn't conform to the schema, flagd logs a warning at load, and the evaluation result is undefined.

Every argument may be computed with other operators, such as `var` or `cat`: the variant, the reason, and the metadata list, its pairs, their keys, and their values.

When a computed argument is malformed, such as an object variant or a metadata number that is `NaN` or infinite, `resolve` returns null like the other custom operations, and the flag falls back to its default variant.

Older flagd binaries do not know this operator.

Upgrade flagd before a configuration starts using `resolve`.

```json
{
  "resolve": [
    "on",
    "SPLIT_RANDOM",
    [["details", "us-10pct-rollout"]]
  ]
}
```

## Example

```json
{
  "flags": {
    "acceptable-feature-stability": {
      "state": "ENABLED",
      "defaultVariant": "ga",
      "variants": {
        "alpha": "alpha",
        "beta": "beta",
        "ga": "ga"
      },
      "targeting": {
        "if": [
          { "===": [{ "var": "customerId" }, "customer-A"] },
          {
            "resolve": [
              "alpha",
              "TARGETING_MATCH",
              [["details", "explicit allowlist for enterprise pilot"]]
            ]
          },
          { "in": [{ "var": "customerId" }, ["customer-B1", "customer-B2"]] },
          { "resolve": ["beta", "TARGETING_MATCH"] },
          { "resolve": ["ga", "FALLTHROUGH"] }
        ]
      }
    }
  }
}
```

The `alpha` and `beta` branches both say why they matched.

The final branch selects `ga` with reason `FALLTHROUGH`.

## Fractional split

```json
{
  "fractional": [
    [{ "resolve": ["on", "SPLIT", [["details", "treatment"]]] }, 10],
    [{ "resolve": ["off", "SPLIT", [["details", "holdback"]]] }, 90]
  ]
}
```

Each bucket reports which one fired.

Wrapping a bucket's variant in `resolve` does not change which bucket a user lands in.
