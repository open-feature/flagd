---
description: flagd semver custom operation
---

# Semantic Version Operation

OpenFeature allows clients to pass contextual information which can then be used during a flag evaluation. For example, a client could pass the email address of the user.

In some scenarios, it is desirable to use that contextual information to segment the user population further and thus return dynamic values.

The `sem_ver` evaluation checks if the given property matches a semantic versioning condition.
It returns 'true', if the value of the given property meets the condition, 'false' if not.
Note that the 'sem_ver' evaluation rule must contain exactly three items:

1. Target property: this needs which both resolve to a semantic versioning string
2. Operator: One of the following: `=`, `!=`, `>`, `<`, `>=`, `<=`, `~` (match minor version), `^` (match major version)
3. Target value: this needs which both resolve to a semantic versioning string

The `sem_ver` evaluation returns a boolean, indicating whether the condition has been met.

```js
{
    "if": [
        {
            "sem_ver": [{"var": "version"}, ">=", "1.0.0"]
        },
        "red", null
    ]
}
```

## Example for 'sem_ver' Evaluation

Flags defined as such:

```json
{
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
curl -X POST "localhost:8013/schema.v1.Service/ResolveString" -d '{"flagKey":"headerColor","context":{"version": "1.0.1"}}' -H "Content-Type: application/json"
```

Result:

```json
{"value":"#00FF00","reason":"TARGETING_MATCH","variant":"red"}
```

Command:

```shell
curl -X POST "localhost:8013/schema.v1.Service/ResolveString" -d '{"flagKey":"headerColor","context":{"version": "0.1.0"}}' -H "Content-Type: application/json"
```

Result:

```shell
{"value":"#0000FF","reason":"TARGETING_MATCH","variant":"green"}
```
