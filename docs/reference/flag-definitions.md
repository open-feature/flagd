---
description: flagd flag definition
---

# Flag Definitions

## Flags

`flags` is a **required** property.
The flags property is a top level property that contains a collection of individual flags and their corresponding flag configurations.

```json
{
  "flags": {
    ...
  }
}
```

## Flag Definition

`flag key` is a **required** property.
The flag key **must** uniquely identify a flag so that it can be used during flag evaluation.
The flag key **should** convey the intent of the flag.

```json
{
  "flags": {
    "new-welcome-banner": {
      ...
    }
  }
}
```

## Flag properties

A fully configured flag may look like this.

```json
{
  "flags": {
    "new-welcome-banner": {
      "state": "ENABLED",
      "variants": {
        "true": true,
        "false": false
      },
      "defaultVariant": "false",
      "targeting": { "in": ["@example.com", { "var": "email" }] }
    }
  }
}
```

See below for a detailed description of each property.

### State

`state` is a **required** property.
Validate states are "ENABLED" or "DISABLED".
When the state is set to "DISABLED", flagd will behave like the flag doesn't exist.

Example:

```json
"state": "ENABLED"
```

### Variants

`variants` is a **required** property.
It is an object containing the possible variations supported by the flag.
All the values of the object **must** be the same type (e.g. boolean, numbers, string, JSON).
The type used as the variant value will correspond directly affects how the flag is accessed.
For example, to use a flag configured with boolean values the `/schema.v1.Service/ResolveBoolean` path should be used.
If another path such as `/schema.v1.Service/ResolveString` is called, a type mismatch occurred and an error is returned.

Example:

```json
"variants": {
  "red": "c05543",
  "green": "2f5230",
  "blue": "0d507b"
}
```

Example:

```json
"variants": {
  "on": true,
  "off": false
}
```

Example of an invalid configuration:

```json
"variants": {
  "on": true,
  "off": "false"
}
```

### Default Variant

`defaultVariant` is a **required** property.
The value **must** match the name of one of the variants defined above.
The default variant is always used unless a targeting rule explicitly overrides it.

Example:

```json
"variants": {
  "on": true,
  "off": false
},
"defaultVariant": "off"
```

Example:

```json
"variants": {
  "red": "c05543",
  "green": "2f5230",
  "blue": "0d507b"
},
"defaultVariant": "red"
```

Example of an invalid configuration:

```json
"variants": {
  "red": "c05543",
  "green": "2f5230",
  "blue": "0d507b"
},
"defaultVariant": "purple"
```

### Targeting Rules

`targeting` is an **optional** property.
A targeting rule **must** be valid JSON.
Flagd uses a modified version of [JsonLogic](https://jsonlogic.com/), as well as some custom pre-processing, to evaluate these rules.
The output of the targeting rule **must** match the name of one of the variants defined above.
If an invalid or null value is returned by the targeting rule, the `defaultVariant` value is used.
If no targeting rules are defined, the response reason will always be `STATIC`, this allows for the flag values to be cached, this behavior is described [here](specifications/rpc-providers.md#caching).

#### Evaluation Context

Evaluation context is included as part of the evaluation request.
For example, when accessing flagd via HTTP, the POST body may look like this:

```json
{
  "flagKey": "booleanFlagKey",
  "context": {
    "email": "noreply@example.com"
  }
}
```

The evaluation context can be accessed in targeting rules using the `var` operation followed the evaluation context property name.

| Description                                                    | Example                                              |
| -------------------------------------------------------------- | ---------------------------------------------------- |
| Retrieve property from the evaluation context                  | `#!json { "var": "email" }`                          |
| Retrieve property from the evaluation context or use a default | `#!json { "var": ["email", "noreply@example.com"] }` |
| Retrieve a nested property from the evaluation context         | `#!json { "var": "user.email" }`                     |

> For more information, see the `var` section in the [JsonLogic documentation](https://jsonlogic.com/operations.html#var).

#### Conditions

Conditions can be used to control the logical flow and grouping of targeting rules.

| Conditional | Example                                                                                                                                                                                                  |
| ----------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| If          | Logic: `#!json {"if" : [ true, "yes", "no" ]}`<br>Result: `"yes"`<br><br>Logic: `#!json {"if" : [ false, "yes", "no" ]}`<br>Result: `"no"`                                                               |
| If else     | Logic: `#!json {"if" : [ false, "yes", false, "no", "maybe" ]}`<br>Result: `"maybe"`<br><br>Logic: `#!json {"if" : [ false, "yes", false, "no", false, "maybe", "who knows" ]}`<br>Result: `"who knows"` |
| Or          | Logic: `#!json {"or" : [ true, false ]}`<br>Result: `true`<br><br>Logic: `#!json {"or" : [ false, false ]}`<br>Result: `false`                                                                           |
| And         | Logic: `#!json {"and" : [ true, false ]}`<br>Result: `false`<br><br>Logic: `#!json {"and" : [ true, true ]}`<br>Result: `true`                                                                           |

#### Operations

Operations are used to take action on, or compare properties retrieved from the context.
These are provided out-of-the-box by JsonLogic.

| Operator               | Description                                                          | Context type | Example                                                                                                                                                                |
| ---------------------- | -------------------------------------------------------------------- | ------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Equals                 | Attribute equals the specified value, with type coercion.            | any          | Logic: `#!json { "==" : [1, 1] }`<br>Result: `true`<br><br>Logic: `#!json { "==" : [1, "1"] }`<br>Result: `true`                                                       |
| Strict equals          | Attribute equals the specified value, with strict comparison.        | any          | Logic: `#!json { "===" : [1, 1] }`<br>Result: `true`<br><br>Logic: `#!json { "===" : [1, "1"] }`<br>Result: `false`                                                    |
| Not equals             | Attribute doesn't equal the specified value, with type coercion.     | any          | Logic: `#!json { "!=" : [1, 2] }`<br>Result: `true`<br><br>Logic: `#!json { "!=" : [1, "1"] }`<br>Result: `false`                                                      |
| Strict not equal       | Attribute doesn't equal the specified value, with strict comparison. | any          | Logic: `#!json { "!==" : [1, 2] }`<br>Result: `true`<br><br>Logic: `#!json { "!==" : [1, "1"] }`<br>Result: `true`                                                     |
| Exists                 | Attribute is defined                                                 | any          | Logic: `#!json { "!!": [ "mike" ] }`<br>Result: `true`<br><br>Logic: `#!json { "!!": [ "" ] }`<br>Result: `false`                                                      |
| Not exists             | Attribute is not defined                                             | any          | Logic: `#!json {"!": [ "mike" ] }`<br>Result: `false`<br><br>Logic: `#!json {"!": [ "" ] }`<br>Result: `true`                                                          |
| Greater than           | Attribute is greater than the specified value                        | number       | Logic: `#!json { ">" : [2, 1] }`<br>Result: `true`<br><br>Logic: `#!json { ">" : [1, 2] }`<br>Result: `false`                                                          |
| Greater than or equals | Attribute is greater or equal to the specified value                 | number       | Logic: `#!json { ">=" : [2, 1] }`<br>Result: `true`<br><br>Logic: `#!json { ">=" : [1, 1] }`<br>Result: `true`                                                         |
| Less than              | Attribute is less than the specified value                           | number       | Logic: `#!json { "<" : [1, 2] }`<br>Result: `true`<br><br>Logic: `#!json { "<" : [2, 1] }`<br>Result: `false`                                                          |
| Less than or equals    | Attribute is less or equal to the specified value                    | number       | Logic: `#!json { "<=" : [1, 1] }`<br>Result: `true`<br><br>Logic: `#!json { "<=" : [2, 1] }`<br>Result: `false`                                                        |
| Between                | Attribute between the specified values                               | number       | Logic: `#!json { "<" : [1, 5, 10]}`<br>Result: `true`<br><br>Logic: `#!json { "<" : [1, 11, 10] }`<br>Result: `false`                                                  |
| Between inclusive      | Attribute between or equal to the specified values                   | number       | Logic: `#!json {"<=" : [1, 1, 10] }`<br>Result: `true`<br><br>Logic: `#!json {"<=" : [1, 11, 10] }`<br>Result: `false`                                                 |
| Contains               | Contains string                                                      | string       | Logic: `#!json { "in": ["Spring", "Springfield"] }`<br>Result: `true`<br><br>Logic: `#!json { "in":["Illinois", "Springfield"] }`<br>Result: `false`                   |
| Not contains           | Does not contain a string                                            | string       | Logic: `#!json { "!": { "in":["Spring", "Springfield"] } }`<br>Result: `false`<br><br>Logic: `#!json { "!": { "in":["Illinois", "Springfield"] } }`<br>Result: `true`  |
| In                     | Attribute is in an array of strings                                  | string       | Logic: `#!json { "in" : [ "Mike", ["Bob", "Mike"]] }`<br>Result: `true`<br><br>Logic: `#!json { "in":["Todd", ["Bob", "Mike"]] }`<br>Result: `false`                   |
| Not it                 | Attribute is not in an array of strings                              | string       | Logic: `#!json { "!": { "in" : [ "Mike", ["Bob", "Mike"]] } }`<br>Result: `false`<br><br>Logic: `#!json { "!": { "in":["Todd", ["Bob", "Mike"]] } }`<br>Result: `true` |

#### Custom Operations

These are custom operations specific to flagd and flagd providers.
They are purpose built extensions to JsonLogic in order to support common feature flag use cases.

| Function                           | Description                                         | Example                                                                                                                                                                                                                                                                                  |
| ---------------------------------- | --------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `fractional` (_available v0.6.4+_) | Deterministic, pseudorandom fractional distribution | Logic: `#!json { "fractional" : [ { "var": "email" }, [ "red" , 50], [ "green" , 50 ] ] }` <br>Result: Pseudo randomly `red` or `green` based on the evaluation context property `email`.<br><br>Additional documentation can be found [here](./custom-operations/fractional-operation.md).        |
| `starts_with`                      | Attribute starts with the specified value           | Logic: `#!json { "starts_with" : [ "192.168.0.1", "192.168"] }`<br>Result: `true`<br><br>Logic: `#!json { "starts_with" : [ "10.0.0.1", "192.168"] }`<br>Result: `false`<br>Additional documentation can be found [here](./custom-operations/string-comparison-operation.md).                     |
| `ends_with`                        | Attribute ends with the specified value             | Logic: `#!json { "ends_with" : [ "noreply@example.com", "@example.com"] }`<br>Result: `true`<br><br>Logic: `#!json { ends_with" : [ "noreply@example.com", "@test.com"] }`<br>Result: `false`<br>Additional documentation can be found [here](./custom-operations/string-comparison-operation.md).|
| `sem_ver`                          | Attribute matches a semantic versioning condition   | Logic: `#!json {"sem_ver": ["1.1.2", ">=", "1.0.0"]}`<br>Result: `true`<br><br>Additional documentation can be found [here](./custom-operations/semver-operation.md).                                                                                                                   |

#### $flagd properties in the evaluation context

Flagd adds the following properties to the evaluation context that can be used in the targeting rules.

| Property | Description | From version |
|----------|-------------|--------------|
| `$flagd.flagKey` | The identifier for the flag being evaluated | v0.6.4 |
| `$flagd.timestamp`| A unix timestamp (in seconds) of the time of evaluation | v0.6.7 |

## Shared evaluators

`$evaluators` is an **optional** property.
It's a collection of shared targeting configurations used to reduce the number of duplicated targeting rules.

Example:

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
          },
          "binet",
          null
        ]
      }
    },
    "headerColor": {
      "variants": {
        "red": "#FF0000",
        "blue": "#0000FF",
        "green": "#00FF00",
        "yellow": "#FFFF00"
      },
      "defaultVariant": "red",
      "state": "ENABLED",
      "targeting": {
        "if": [
          {
            "$ref": "emailWithFaas"
          },
          {
            "fractional": [
              { "var": "email" },
              ["red", 25],
              ["blue", 25],
              ["green", 25],
              ["yellow", 25]
            ]
          },
          null
        ]
      }
    }
  },
  "$evaluators": {
    "emailWithFaas": {
      "in": [
        "@faas.com",
        {
          "var": ["email"]
        }
      ]
    }
  }
}
```

## Examples

Sample configurations can be found at <https://github.com/open-feature/flagd/tree/main/config/samples>.
