# Flag Configuration

A flagd configuration is represented as a JSON object. Feature flag configurations can be found under `flags` and each item within `flags` represents a flag key (the unique identifier for a flag) and its corresponding configuration.

Sample configurations can be found at https://github.com/open-feature/flagd/tree/main/config/samples.

## Flag configuration properties

### State

`state` is **required** property. Validate states are "ENABLED" or "DISABLED". When the state is set to "DISABLED", flagd will behave like the flag doesn't exist.

Example:

```
"state": "ENABLED"
```

### Variants

`variants` is a **required** property. It is an object containing the possible variations supported by the flag. All the values of the object **must** but the same type (e.g. boolean, numbers, string, JSON). The type used as the variant value will correspond directly affects how the flag is accessed. For example, to use a flag configured with boolean values the `/schema.v1.Service/ResolveBoolean` path should be used. If another path such as `/schema.v1.Service/ResolveString` is called, a type mismatch occurred and an error is returned.

Example:

```
"variants": {
  "red": "c05543",
  "green": "2f5230",
  "blue": "0d507b"
}
```

Example:

```
"variants": {
  "on": true,
  "off": false
}
```

Example of an invalid configuration:

```
"variants": {
  "on": true,
  "off": "false"
}
```

### Default Variant

`defaultVariant` is a **required** property. The value **must** match the name of one of the variants defined above. The default variant is always used unless a targeting rule explicitly overrides it.

Example:

```
"variants": {
  "on": true,
  "off": false
},
"defaultVariant": "off"
```

Example:

```
"variants": {
  "red": "c05543",
  "green": "2f5230",
  "blue": "0d507b"
},
"defaultVariant": "red"
```

Example of an invalid configuration:

```
"variants": {
  "red": "c05543",
  "green": "2f5230",
  "blue": "0d507b"
},
"defaultVariant": "purple"
```

### Targeting Rules

`targeting` is an **optional** property. A targeting rule **must** be valid JSON. Flagd uses a modified version of [JSON Logic](https://jsonlogic.com/), as well as some custom pre-processing, to evaluate these rules. The output of the targeting rule **must** match the name of one of the variants defined above. If an invalid or null value is is returned by the targeting rule, the `defaultVariant` value is used. If no targeting rules are defined, the response reason will always be `STATIC`, this allows for the client side caching of these flag values, this behavior is described [here](../other_resources/caching.md).

The [JSON Logic playground](https://jsonlogic.com/play.html) is a great way to experiment with new targeting rules. The following example shows how a rule could be configured to return `binet` when the email (which comes from evaluation context) contains `@faas.com`. If the email wasn't included in the evaluation context or doesn't contain `@faas.com`, null is returned and the `defaultVariant` is used instead.

<details>
  <summary>Click here to see how this targeting rule would look in the JSON Logic playground.</summary>

1. Open the [JSON Logic playground](https://jsonlogic.com/play.html) in your favorite browser
1. Add the follow JSON as the `Rule`:

    ```json
    {
      "if": [
        {
          "in": [
            "@faas.com",
            {
              "var": ["email"]
            }
          ]
        },
        "binet",
        null
      ]
    }
    ```

1. Add the following JSON as the `Data`:

    ```json
    {
      "email": "test@faas.com"
    }
    ```

1. Click `Compute`
1. confirm the output show `"binet"`
1. Optionally, experiment with different rules and data
