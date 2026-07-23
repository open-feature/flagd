---
description: flagd string casing custom operations
---

# Lower / Upper Operation

The `lower`/`upper` operation is a custom JsonLogic operation which transforms a string so that targeting rules can match case-insensitively.
It takes a single argument which must resolve to a string value, and returns that string converted to lower or upper case.
Casing is restricted to ASCII characters (`A-Z` and `a-z`); all other characters are returned unchanged.
This keeps the result deterministic and identical across every flagd provider implementation, and is intended for ASCII data such as email addresses and non-IDN domains.

Because the transform composes with other operations, case-insensitive matching is achieved by wrapping the relevant operands rather than introducing ignore-case variants of every operator.

```js
// lower used to match a property case-insensitively
"==": [
  // transform the context property to lower case
  {"lower": [{"var": "email"}]},
  // comparison target (also lower case)
  "user@example.com"
]
```

## Example

Flags defined as such:

```json
{
  "$schema": "https://flagd.dev/schema/v0/flags.json",
  "flags": {
    "headerColor": {
      "variants": {
        "red": "#FF0000",
        "blue": "#0000FF",
        "green": "#00FF00"
      },
      "defaultVariant": "blue",
      "state": "ENABLED",
      "targeting": {
        "if": [
          {
            "==": [{"lower": [{"var": "email"}]}, "user@example.com"]
          },
          "red", "green"
        ]
      }
    }
  }
}
```

will return variant `red`, if the value of the `email` property equals `user@example.com` ignoring case, and the variant `green` otherwise.

Command:

```shell
curl -X POST "localhost:8013/flagd.evaluation.v1.Service/ResolveString" -d '{"flagKey":"headerColor","context":{"email": "User@Example.com"}}' -H "Content-Type: application/json"
```

Result:

```json
{"value":"#FF0000","reason":"TARGETING_MATCH","variant":"red"}
```

## Upper Operation Configuration

The `upper` evaluation behaves identically to `lower`, returning the ASCII upper-cased form of its argument.
It is commonly used to normalize short codes such as country or currency identifiers before comparison.

```js
// upper used to match a property case-insensitively
"upper": [
  // transform the context property to upper case
  {"var": "country"}
]
```

## Example for 'upper' Operation

Flags defined as such:

```json
{
  "$schema": "https://flagd.dev/schema/v0/flags.json",
  "flags": {
    "headerColor": {
      "variants": {
        "red": "#FF0000",
        "blue": "#0000FF",
        "green": "#00FF00"
      },
      "defaultVariant": "blue",
      "state": "ENABLED",
      "targeting": {
        "if": [
          {
            "==": [{"upper": [{"var": "country"}]}, "US"]
          },
          "red", "green"
        ]
      }
    }
  }
}
```

will return variant `red`, if the value of the `country` property equals `US` ignoring case, and the variant `green` otherwise.

Command:

```shell
curl -X POST "localhost:8013/flagd.evaluation.v1.Service/ResolveString" -d '{"flagKey":"headerColor","context":{"country": "us"}}' -H "Content-Type: application/json"
```

Result:

```json
{"value":"#FF0000","reason":"TARGETING_MATCH","variant":"red"}
```
