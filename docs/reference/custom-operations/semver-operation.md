---
description: flagd semver custom operation
---

# Semantic Version Operation

The `sem_ver` evaluation checks if the given property matches a semantic versioning condition.
It returns 'true', if the value of the given property meets the condition, 'false' if not.
Note that the 'sem_ver' evaluation rule must contain exactly three items:

1. Target property: this needs which both resolve to a semantic versioning string
2. Operator: One of the following: `=`, `!=`, `>`, `<`, `>=`, `<=`, `~` (match minor version), `^` (match major version)
3. Target value: this needs which both resolve to a semantic versioning string

The `sem_ver` evaluation returns a boolean, indicating whether the condition has been met.

```js
// sem_ver property name used in a targeting rule
"sem_ver": [
  // Evaluation context property to be evaluated
  {"var": "version"},
  // Operator to use for comparison
  ">=",
  // Target value to compare against
  "1.0.0"
]
```

!!! tip

    Version strings may include a `v` or `V` prefix (e.g. `v1.0.0`), which is stripped before comparison.
    Partial versions such as `1.0` or `1` are also accepted and padded with `.0` to form a complete version.
    Numeric context values (e.g. integer `1`) are coerced to strings before parsing.
    Build metadata (e.g. `1.0.0+build`) is ignored during comparison, per the SemVer specification.

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
            "sem_ver": [{"var": "version"}, ">=", "1.0.0"]
          },
          "red", "green"
        ]
      }
    }
  }
}
```

will return variant `red`, if the value of the `version` is a semantic version that is greater than or equal to `1.0.0`.

Command:

```shell
curl -X POST "localhost:8013/flagd.evaluation.v1.Service/ResolveString" -d '{"flagKey":"headerColor","context":{"version": "1.0.1"}}' -H "Content-Type: application/json"
```

Result:

```json
{"value":"#00FF00","reason":"TARGETING_MATCH","variant":"red"}
```

Command:

```shell
curl -X POST "localhost:8013/flagd.evaluation.v1.Service/ResolveString" -d '{"flagKey":"headerColor","context":{"version": "0.1.0"}}' -H "Content-Type: application/json"
```

Result:

```shell
{"value":"#0000FF","reason":"TARGETING_MATCH","variant":"green"}
```
