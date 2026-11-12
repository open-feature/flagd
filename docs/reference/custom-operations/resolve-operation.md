---
description: flagd resolve custom operation
---

# Resolve Operation

`resolve` ends a targeting rule with a variant, and optionally a custom reason and metadata for the branch that fired.

`"red"` and `{ "resolve": ["red"] }` select the same variant.

A missing or null reason keeps `TARGETING_MATCH`.

The reason is free-form: any string is passed through as the evaluation reason, including flagd's own reasons such as `FALLBACK` or `DISABLED`.

Metadata is a shallow list of `[key, value]` pairs.

Values are strings, numbers, or booleans.

Rule metadata overwrites flag metadata, which overwrites flag-set metadata.

`resolve` is terminal.

It is valid as the whole targeting rule, or as a result of `if` / `?:`.

flagd marks a flag invalid at load time when `resolve` appears anywhere else, including inside `==`, `cat`, an `if` condition, or a `fractional` bucket, or when its literal arguments are malformed.

Evaluating an invalid flag returns `PARSE_ERROR`. Other flags in the same configuration still load.

`fractional` buckets can't carry a reason or metadata yet.

To tag a split, compare the `fractional` result in an `if` condition and use `resolve` in its branches, as in the [example below](#fractional-split).

Arguments may be computed with other operators, such as `var` or `cat`.

When a computed argument is malformed, such as a variant that is not a string or a metadata number that is `NaN` or infinite, `resolve` returns null like the other custom operations, and the flag falls back to its default variant.

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
  "if": [
    {
      "==": [
        { "fractional": [{ "var": "targetingKey" }, ["on", 10], ["off", 90]] },
        "on"
      ]
    },
    { "resolve": ["on", "SPLIT_RANDOM", [["details", "10pct-rollout"]]] },
    { "resolve": ["off", "SPLIT_FALLTHROUGH", [["details", "holdback"]]] }
  ]
}
```

Both buckets of the split now report which one fired.
