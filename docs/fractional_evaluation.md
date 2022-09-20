### Fractional Evaluation

The `fractionalEvaluation` operation is a custom JsonLogic operation which deterministically selects a variant based on
the defined distribution of each variant (as a percentage). This works by hashing ([xxHash](https://cyan4973.github.io/xxHash/))
the given data point, converting it into an int in the range [0, 99]. Whichever range this int falls in decides which variant
is selected. As hashing is deterministic we can be sure to get the same result every time for the same data point.

<u>Example</u>

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
              "fractionalEvaluation": [
                "email",
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

```shell
$ curl -X POST "localhost:8013/schema.v1.Service/ResolveString" -d ''{"flagKey":"myFloatFlag","context":{"email": "foo@bar.com"}}'' -H "Content-Type: application/json"
{"value":"#0000FF","reason":"TARGETING_MATCH","variant":"blue"}%

$ curl -X POST "localhost:8013/schema.v1.Service/ResolveString" -d ''{"flagKey":"myFloatFlag","context":{"email": "foo@test.com"}}'' -H "Content-Type: application/json"
{"value":"#00FF00","reason":"TARGETING_MATCH","variant":"green"}%
```
