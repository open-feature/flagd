# StartsWith/EndsWith Evaluation

OpenFeature allows clients to pass contextual information which can then be used during a flag evaluation. For example, a client could pass the email address of the user.

In some scenarios, it is desirable to use that contextual information to segment the user population further and thus return dynamic values.

Look at the [headerColor](https://github.com/open-feature/flagd/blob/main/samples/example_flags.flagd.json#L88-#L133) flag. The `defaultVariant` is `red`, but it contains a [targeting rule](reusable_targeting_rules.md), meaning a fractional evaluation occurs for flag evaluation with a `context` object containing `email` and where that `email` value contains `@faas.com`.

In this case, `25%` of the email addresses will receive `red`, `25%` will receive `blue`, and so on.

Importantly, the evaluations are "sticky" meaning that the same `email` address will always belong to the same "bucket" and thus always receive the same color.

## StartsWith/EndsWith Evaluation: Technical Description

The `starts_with`/`ends_with` operation is a custom JsonLogic operation which selects a variant based on
whether the specified property starts/ends with a certain value.

## StartsWith Evaluation Configuration

The `starts_with` evaluation can be added as part of a targeting definition.
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

## Example

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
curl -X POST "localhost:8013/schema.v1.Service/ResolveString" -d '{"flagKey":"headerColor","context":{"email": "user@faas.com"}}' -H "Content-Type: application/json"
```

Result:

```json
{"value":"#00FF00","reason":"TARGETING_MATCH","variant":"red"}
```

Command:

```shell
curl -X POST "localhost:8013/schema.v1.Service/ResolveString" -d '{"flagKey":"headerColor","context":{"email": "foo@bar.com"}}' -H "Content-Type: application/json"
```

Result:

```shell
{"value":"#0000FF","reason":"TARGETING_MATCH","variant":"green"}
```

## EndsWith Evaluation Configuration

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

## Example

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
curl -X POST "localhost:8013/schema.v1.Service/ResolveString" -d '{"flagKey":"headerColor","context":{"email": "user@faas.com"}}' -H "Content-Type: application/json"
```

Result:

```json
{"value":"#00FF00","reason":"TARGETING_MATCH","variant":"red"}
```

Command:

```shell
curl -X POST "localhost:8013/schema.v1.Service/ResolveString" -d '{"flagKey":"headerColor","context":{"email": "foo@bar.com"}}' -H "Content-Type: application/json"
```

Result:

```shell
{"value":"#0000FF","reason":"TARGETING_MATCH","variant":"green"}
```
