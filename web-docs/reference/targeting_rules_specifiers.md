# Reusable Targeting Rules and Specifiers

## Reusable targeting rules

At the same level as the `flags` key one can define an `$evaluators` object.
Each object defined under `$evaluators` is
a reusable targeting rule.
In any targeting rule one can reference a defined reusable targeting rule, foo, like so:
`"$ref": "foo"`

## Example (using `in`)

Flags/evaluators defined as such:

```json
{
    "flags": {
        "fibAlgo": {
          "variants": {
            "recursive": "recursive",
            "memo": "memo",
            "loop": "loop",
            "binet": "binet"
          },
          "defaultVariant": "recursive",
          "state": "ENABLED",
          "targeting": {
            "if": [
              {
                "$ref": "emailWithFaas"
              }, "binet", null
            ]
          }
        }
    },
    "$evaluators": {
        "emailWithFaas": {
              "in": ["@faas.com", {
                "var": ["email"]
              }]
        }
    }
}
```

becomes (once the `$evaluators` have been substituted):

```json
{
    "flags": {
        "fibAlgo": {
          "variants": {
            "recursive": "recursive",
            "memo": "memo",
            "loop": "loop",
            "binet": "binet"
          },
          "defaultVariant": "recursive",
          "state": "ENABLED",
          "targeting": {
            "if": [
              {
                "in": ["@faas.com", {
                "var": ["email"]
              }]
              }, "binet", null
            ]
          }
        }
    }
}
```

## StartsWith/EndsWith Evaluation

OpenFeature allows clients to pass contextual information which can then be used during a flag evaluation. For example, a client could pass the email address of the user.

In some scenarios, it is desirable to use that contextual information to segment the user population further and thus return dynamic values.

### StartsWith/EndsWith Evaluation: Technical Description

The `starts_with`/`ends_with` operation is a custom JsonLogic operation which selects a variant based on
whether the specified property starts/ends with a certain value.

### StartsWith Evaluation Configuration

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

### Example for 'starts_with' Evaluation

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

### EndsWith Evaluation Configuration

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

### Example for 'ends_with' Evaluation

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

## SemVer Evaluation

OpenFeature allows clients to pass contextual information which can then be used during a flag evaluation. For example, a client could pass the email address of the user.

In some scenarios, it is desirable to use that contextual information to segment the user population further and thus return dynamic values.

### SemVer Evaluation: Technical Description

The `sem_ver` evaluation checks if the given property matches a semantic versioning condition.
It returns 'true', if the value of the given property meets the condition, 'false' if not.

### SemVer Evaluation Configuration

The `sem_ver` evaluation can be added as part of a targeting definition.
Note that the 'sem_ver' evaluation rule must contain exactly three items:

1. Target property: this needs which both resolve to a semantic versioning string
1. Operator: One of the following: `=`, `!=`, `>`, `<`, `>=`, `<=`, `~` (match minor version), `^` (match major version)
1. Target value: this needs which both resolve to a semantic versioning string

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

### Example for 'sem_ver' Evaluation

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
