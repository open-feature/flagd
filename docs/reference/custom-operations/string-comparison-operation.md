---
description: flagd string custom operations
---

# Starts-With / Ends-With Operation

OpenFeature allows clients to pass contextual information which can then be used during a flag evaluation. For example, a client could pass the email address of the user.

In some scenarios, it is desirable to use that contextual information to segment the user population further and thus return dynamic values.

The `starts_with`/`ends_with` operation is a custom JsonLogic operation which selects a variant based on
whether the specified property starts/ends with a certain value.
The value is an array consisting of exactly two items, which both need to resolve to a string value.
The first entry of the array represents the property to be considered, while the second entry represents
the target value, i.e. the prefix that needs to be present in the value of the referenced property.
This value should typically be something that remains consistent for the duration of a users session (e.g. email or session ID).
The `starts_with` evaluation returns a boolean, indicating whether the condition has been met.

```js
// starts_with property name used in a targeting rule
"starts_with": [
  // Evaluation context property the be evaluated
  {"var": "email"},
  // prefix that has to be present in the value of the referenced property  
  "user@faas"
]
```

## Example for 'starts_with' Operation

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
            "starts_with": [{"var": "email"}, "user@faas"]
          },
          "red", "green"
        ]
      }
    }
  }
}
```

will return variant `red`, if the value of the `email` property starts with `user@faas`, and the variant `green` otherwise.

Command:

```shell
curl -X POST "localhost:8013/flagd.evaluation.v1.Service/ResolveString" -d '{"flagKey":"headerColor","context":{"email": "user@faas.com"}}' -H "Content-Type: application/json"
```

Result:

```json
{"value":"#00FF00","reason":"TARGETING_MATCH","variant":"red"}
```

Command:

```shell
curl -X POST "localhost:8013/flagd.evaluation.v1.Service/ResolveString" -d '{"flagKey":"headerColor","context":{"email": "foo@bar.com"}}' -H "Content-Type: application/json"
```

Result:

```shell
{"value":"#0000FF","reason":"TARGETING_MATCH","variant":"green"}
```

## EndsWith Operation Configuration

The `ends_with` evaluation can be added as part of a targeting definition.
The value is an array consisting of exactly two items, which both need to resolve to a string value.
The first entry of the array represents the property to be considered, while the second entry represents
the target value, i.e. the suffix that needs to be present in the value of the referenced property.
This value should typically be something that remains consistent for the duration of a users session (e.g. email or session ID).
The `ends_with` evaluation returns a boolean, indicating whether the condition has been met.

```js
// starts_with property name used in a targeting rule
"ends_with": [
  // Evaluation context property the be evaluated
  {"var": "email"},
  // suffix that has to be present in the value of the referenced property  
  "faas.com"
]
```

## Example for 'ends_with' Operation

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
            "ends_with": [{"var": "email"}, "faas.com"]
          },
          "red", "green"
        ]
      }
    }
  }
}
```

will return variant `red`, if the value of the `email` property ends with `faas.com`, and the variant `green` otherwise.

Command:

```shell
curl -X POST "localhost:8013/flagd.evaluation.v1.Service/ResolveString" -d '{"flagKey":"headerColor","context":{"email": "user@faas.com"}}' -H "Content-Type: application/json"
```

Result:

```json
{"value":"#00FF00","reason":"TARGETING_MATCH","variant":"red"}
```

Command:

```shell
curl -X POST "localhost:8013/flagd.evaluation.v1.Service/ResolveString" -d '{"flagKey":"headerColor","context":{"email": "foo@bar.com"}}' -H "Content-Type: application/json"
```

Result:

```shell
{"value":"#0000FF","reason":"TARGETING_MATCH","variant":"green"}
```
