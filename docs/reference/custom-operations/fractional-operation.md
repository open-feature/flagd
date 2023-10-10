---
description: flagd fractional custom operation
---

# Fractional Operation

OpenFeature allows clients to pass contextual information which can then be used during a flag evaluation. For example, a client could pass the email address of the user.

```js
// Factional evaluation property name used in a targeting rule
"fractional": [
  // Evaluation context property used to determine the split
  { "var": "email" },
  // Split definitions contain an array with a variant and percentage
  // Percentages must add up to 100
  [
    // Must match a variant defined in the flag definition
    "red",
    // The probability this variant is selected
    50
  ],
  [
    // Must match a variant defined in the flag definition
    "green",
    // The probability this variant is selected
    50
  ]
]
```

See the [headerColor](https://github.com/open-feature/flagd/blob/main/samples/example_flags.flagd.json#L88-#L133) flag.
The `defaultVariant` is `red`, but it contains a [targeting rule](../flag-definitions.md#targeting-rules), meaning a fractional evaluation occurs for flag evaluation with a `context` object containing `email` and where that `email` value contains `@faas.com`.

In this case, `25%` of the evaluations will receive `red`, `25%` will receive `blue`, and so on.

Assignment is deterministic (sticky) based on the expression supplied as the first parameter (`{ "var": "email" }`, in this case).
The value retrieved by this expression is referred to as the "bucketing value".
The bucketing value expression can be omitted, in which case a concatenation of the `targetingKey` and the `flagKey` will be used.

The `fractional` operation is a custom JsonLogic operation which deterministically selects a variant based on
the defined distribution of each variant (as a percentage).
This works by hashing ([murmur3](https://github.com/aappleby/smhasher/blob/master/src/MurmurHash3.cpp))
the given data point, converting it into an int in the range [0, 99].
Whichever range this int falls in decides which variant
is selected.
As hashing is deterministic we can be sure to get the same result every time for the same data point.

The `fractional` operation can be added as part of a targeting definition.
The value is an array and the first element is the name of the property to use from the evaluation context.
This value should typically be something that remains consistent for the duration of a users session (e.g. email or session ID).
The other elements in the array are nested arrays with the first element representing a variant and the second being the percentage that this option is selected.
There is no limit to the number of elements but the configured percentages must add up to 100.

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
      "defaultVariant": "red",
      "state": "ENABLED",
      "targeting": {
        "fractional": [
          { "var": "email" },
          [
            "red",
            50
          ],
          [
            "blue",
            20
          ],
          [
            "green",
            30
          ]
        ]
      }
    }
  }
}
```

will return variant `red` 50% of the time, `blue` 20% of the time & `green` 30% of the time.

Command:

```shell
curl -X POST "localhost:8013/schema.v1.Service/ResolveString" -d '{"flagKey":"headerColor","context":{"email": "foo@bar.com"}}' -H "Content-Type: application/json"
```

Result:

```shell
{"value":"#0000FF","reason":"TARGETING_MATCH","variant":"blue"}
```

Command:

```shell
curl -X POST "localhost:8013/schema.v1.Service/ResolveString" -d '{"flagKey":"headerColor","context":{"email": "foo@test.com"}}' -H "Content-Type: application/json"
```

Result:

```json
{"value":"#00FF00","reason":"TARGETING_MATCH","variant":"green"}
```

Notice that rerunning either curl command will always return the same variant and value.
The only way to get a different value is to change the email or update the `fractional` configuration.

### Migrating from legacy "fractionalEvaluation"

If you are using a legacy fractional evaluation (`fractionalEvaluation`), it's recommended you migrate to `fractional`.
The new `fractional` evaluator supports nested properties and JsonLogic expressions.
To migrate, simply use a JsonLogic variable declaration for the bucketing property, instead of a string:

old:

```json
"fractionalEvaluation": [
    "email",
    [ "red", 25 ], [ "blue", 25 ], [ "green", 25 ], [ "yellow", 25 ]
]
```

new:

```json
"fractional": [
    { "var": "email" },
    [ "red", 25 ], [ "blue", 25 ], [ "green", 25 ], [ "yellow", 25 ]
]
```
