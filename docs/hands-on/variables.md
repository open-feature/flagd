---
description: Using var to access evaluation context in JsonLogic
---

# Variables with `var`

The `var` operator extracts values from the **evaluation context** — the data your application sends along with each flag evaluation request.

## Basic usage

```json title="Extract a variable" hl_lines="1"
{ "var": "email" }
```

This pulls the value of `email` from the evaluation context.

If the context is:

```json
{ "email": "alice@example.com", "age": 30 }
```

then `{ "var": "email" }` evaluates to `"alice@example.com"`.

## Using `var` in conditions

This is where things get practical. Combine `var` with operators to build real targeting rules:

```json title="Check a user's email domain" hl_lines="4"
{
  "if": [
    { "==": [
      { "var": "email" },
      "admin@example.com"
    ]},
    "admin-variant",
    "default-variant"
  ]
}
```

&nbsp;

Breaking down the highlighted line:

| Expression | Evaluates to |
|-----------|--------------|
| `{ "var": "email" }` | `"admin@example.com"` (from context) |
| `{ "==": [ ..., "admin@example.com" ] }` | `true` |
| `{ "if": [ true, "admin-variant", ... ] }` | `"admin-variant"` |

## Nested properties

You can access nested context values using dot notation:

```json title="Nested property access" hl_lines="1"
{ "var": "user.plan" }
```

With context `{ "user": { "plan": "premium" } }`, this returns `"premium"`.

## Default values

Provide a fallback if the property doesn't exist:

```json title="Default value" hl_lines="1"
{ "var": [ "email", "unknown@example.com" ] }
```

If `email` is not in the context, this returns `"unknown@example.com"`.

!!! note

    When `var` has a default value, the parameters become an **array** instead of a plain string.

## A complete flag definition using `var`

Here's how this looks in a full flagd flag definition:

```json title="Full flag with var" hl_lines="12 13"
{
  "$schema": "https://flagd.dev/schema/v0/flags.json",
  "flags": {
    "welcome-banner": {
      "state": "ENABLED",
      "variants": {
        "new-banner": true,
        "old-banner": false
      },
      "defaultVariant": "old-banner",
      "targeting": {
        "if": [
          { "ends_with": [{ "var": "email" }, "@example.com"] },
          "new-banner",
          null
        ]
      }
    }
  }
}
```

In this example:

- `{ "var": "email" }` extracts the user's email from context
- `ends_with` checks if it ends with `@example.com`
- If yes → `"new-banner"`, otherwise → `null` (falls back to `defaultVariant`)

## Special flagd variables

flagd automatically injects some useful properties into the evaluation context:

| Variable | Description |
|----------|-------------|
| `$flagd.flagKey` | The key of the flag being evaluated |
| `$flagd.timestamp` | Unix timestamp (seconds) of the evaluation time |

Access them just like any other variable:

```json title="flagd built-in variables" hl_lines="1"
{ "var": "$flagd.timestamp" }
```

We'll use `$flagd.timestamp` in the next section.

Next: [Custom Operations](./custom-operations.md)
