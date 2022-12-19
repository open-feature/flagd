# Reusable targeting rules

At the same level as the `flags` key one can define an `$evaluators` object. Each object defined under `$evaluators` is
a reusable targeting rule. In any targeting rule one can reference a defined reusable targeting rule, foo, like so:
`"$ref": "foo"`

<u>Example</u>

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
