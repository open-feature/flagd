# Fractional Evaluation

OpenFeature allows clients to pass contextual information which can then be used during a flag evaluation. For example, a client could pass the email address of the user.

In some scenarios, it is desirable to use that contextual information to segment the user population further and thus return dynamic values.

Look at the [headerColor](https://github.com/open-feature/flagd/blob/main/samples/example_flags.flagd.json#L88-#L133) flag. The `defaultVariant` is `red`, but it contains a [targeting rule](reusable_targeting_rules.md), meaning a fractional evaluation occurs for flag evaluation with a `context` object containing `email` and where that `email` value contains `@faas.com`.

In this case, `25%` of the email addresses will receive `red`, `25%` will receive `blue`, and so on.

Importantly, the evaluations are "sticky" meaning that the same `email` address will always belong to the same "bucket" and thus always receive the same color.

## Fractional Evaluation: Technical Description

The `fractional` operation is a custom JsonLogic operation which deterministically selects a variant based on
the defined distribution of each variant (as a percentage).
This works by hashing ([murmur3](https://github.com/aappleby/smhasher/blob/master/src/MurmurHash3.cpp))
the given data point, converting it into an int in the range [0, 99].
Whichever range this int falls in decides which variant
is selected.
As hashing is deterministic we can be sure to get the same result every time for the same data point.

## Fractional evaluation configuration

The `fractional` operation can be added as part of a targeting definition.
The value is an array and the first element is the name of the property to use from the evaluation context.
This value should typically be something that remains consistent for the duration of a users session (e.g. email or session ID).
The other elements in the array are nested arrays with the first element representing a variant and the second being the percentage that this option is selected.
There is no limit to the number of elements but the configured percentages must add up to 100.

```js
// Factional evaluation property name used in a targeting rule
"fractional": [
  // Evaluation context property used to determine the split
  { "var": "email" },
  // Split definitions contain an array with a variant and percentage
  // Percentages must add up to 100
  [
    // Must match a variant defined in the flag configuration
    "red",
    // The probability this configuration is selected
    50
  ],
  [
    // Must match a variant defined in the flag configuration
    "green",
    // The probability this configuration is selected
    50
  ]
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

### Migrating from legacy fractionalEvaluation

If you are using a legacy fractional evaluation (`fractionalEvaluation`), it's recommended you migrate to `fractional`.
The new `fractional` evaluator supports nested properties and json-logic expressions.
To migrate, simply use a json-logic variable declaration for the bucketing property, instead of a string:

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
