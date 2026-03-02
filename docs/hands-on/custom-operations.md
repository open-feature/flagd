---
description: flagd custom JsonLogic operations
---

# Custom Operations

flagd extends JsonLogic with **custom operators** tailored for feature flagging.
They follow the exact same shape: `{ "operator": [ params... ] }`.

## `fractional` — Percentage-based rollouts

The `fractional` operator lets you split users into groups by percentage.
It hashes a value (deterministically) and assigns a variant based on defined weights.

```json title="50/50 split" hl_lines="2 3 4 5 6 7 8 9 10 11"
{
  "fractional": [
    { "cat": [
        { "var": "$flagd.flagKey" },
        { "var": "email" }
    ]},
    [ "red", 50 ],
    [ "green", 50 ]
  ]
}
```

Breaking it down:

| Part | Purpose |
|------|---------|
| `{ "cat": [...] }` | Builds the **bucketing value** by concatenating the flag key and email |
| `[ "red", 50 ]` | 50% chance of variant `"red"` |
| `[ "green", 50 ]` | 50% chance of variant `"green"` |

&nbsp;

!!! tip

    The bucketing value determines which bucket a user falls into.
    Same input → same bucket, every time.
    Using `$flagd.flagKey` + a user identifier prevents collisions across flags.

### Gradual rollout example

Weights don't have to be 50/50:

```json title="10% rollout" hl_lines="7 8"
{
  "fractional": [
    { "cat": [
        { "var": "$flagd.flagKey" },
        { "var": "email" }
    ]},
    [ "new-feature", 10 ],
    [ "old-feature", 90 ]
  ]
}
```

Only 10% of users see `"new-feature"`.

## `starts_with` / `ends_with` — String matching

Check if a context value starts or ends with a given string:

```json title="ends_with" hl_lines="2 3 4"
{
  "ends_with": [
    { "var": "email" },
    "@example.com"
  ]
}
```

Result: `true` if the user's email ends with `@example.com`.

&nbsp;

These are commonly used inside `if` to control variant selection:

```json title="Target by email domain" hl_lines="3 4 5 6"
{
  "if": [
    { "ends_with": [
        { "var": "email" },
        "@beta-testers.com"
    ]},
    "beta-variant",
    null
  ]
}
```

## `sem_ver` — Semantic version comparisons

Compare version strings using semver rules:

```json title="Version gate" hl_lines="2 3 4 5 6"
{
  "sem_ver": [
    { "var": "appVersion" },
    ">=",
    "2.0.0"
  ]
}
```

Result: `true` if context `appVersion` is `2.0.0` or higher.

Supported operators: `=`, `!=`, `>`, `<`, `>=`, `<=`, `~` (match minor), `^` (match major)

## Comparing timestamps

flagd injects `$flagd.timestamp` (Unix seconds) into every evaluation context.
You can use standard comparison operators with it to enable time-based targeting:

```json title="Enable after a date" hl_lines="3 4 5 6 7 8"
{
  "if": [
    { ">=": [
        { "var": "$flagd.timestamp" },
        1735689600
    ]},
    "new-year-variant",
    "default-variant"
  ]
}
```

This returns `"new-year-variant"` after Jan 1, 2025 00:00:00 UTC (Unix timestamp `1735689600`).

&nbsp;

!!! tip

    Combine timestamps with `and` to create time windows:

    ```json
    { "and": [
        { ">=": [ { "var": "$flagd.timestamp" }, 1735689600 ] },
        { "<":  [ { "var": "$flagd.timestamp" }, 1735776000 ] }
    ]}
    ```

    This is `true` only during a 24-hour window.

## Combining custom operations

Custom operators compose with standard JsonLogic just like any other operator.
Here's a targeting rule that uses `fractional` only for internal users:

```json title="Combined targeting" hl_lines="3 4 5 6 7 8 9 10 11 12 13 14"
{
  "if": [
    { "ends_with": [
        { "var": "email" },
        "@mycompany.com"
    ]},
    { "fractional": [
        { "cat": [
            { "var": "$flagd.flagKey" },
            { "var": "email" }
        ]},
        [ "new-ui", 50 ],
        [ "old-ui", 50 ]
    ]},
    "old-ui"
  ]
}
```

- Internal users (`@mycompany.com`) get a 50/50 split between `"new-ui"` and `"old-ui"`
- Everyone else gets `"old-ui"`

For full details on each custom operation, see the [Reference](../reference/flag-definitions.md#targeting-rules).
