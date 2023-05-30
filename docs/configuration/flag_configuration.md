# Flag Configuration

- [Flag Configuration](#flag-configuration)
  - [Flags](#flags)
  - [Flag configuration](#flag-configuration-1)
  - [Flag properties](#flag-properties)
    - [State](#state)
    - [Variants](#variants)
    - [Default Variant](#default-variant)
    - [Targeting Rules](#targeting-rules)
      - [Evaluation Context](#evaluation-context)
      - [Conditions](#conditions)
      - [Operators](#operators)
    - [Functions](#functions)
  - [Shared evaluators](#shared-evaluators)
  - [Examples](#examples)

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

## Flag configuration

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
      "targeting": { "in": ["@example.com", { "var": "email" } ] }
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
Flagd uses a modified version of [JSON Logic](https://jsonlogic.com/), as well as some custom pre-processing, to evaluate these rules.
The output of the targeting rule **must** match the name of one of the variants defined above.
If an invalid or null value is is returned by the targeting rule, the `defaultVariant` value is used.
If no targeting rules are defined, the response reason will always be `STATIC`, this allows for the client side caching of these flag values, this behavior is described [here](../other_resources/caching.md).

#### Evaluation Context

Evaluation context can be accessed in targeting rules using the `var` property followed the evaluation context property name.

| Description                                                        | Example                                    |
| ------------------------------------------------------------------ | ------------------------------------------ |
| Retrieve property from the evaluation context                      | { "var": "email" }                         |
| Retrieve property from the evaluation context or use default value | { "var": ["email", "noreply@example.com] } |
| Retrieve a nested property from the evaluation context             | { "var": "user.email" }                    |

> For more information, see the `var` section in the [JSON Logic documentation](https://jsonlogic.com/operations.html#var).

#### Conditions

| Conditional | Example                                                                                                                                                                                |
| ----------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| If          | Logic: `{"if" : [ true, "yes", "no" ]}`<br>Result: "yes"<br><br>Logic: `{"if" : [ false, "yes", "no" ]}`<br>Result: "no"                                                               |
| If else     | Logic: `{"if" : [ false, "yes", false, "no", "maybe" ]}`<br>Result: "maybe"<br><br>Logic: `{"if" : [ false, "yes", false, "no", false, "maybe", "who knows" ]}`<br>Result: "who knows" |
| Or          | Logic: `{"or" : [ true, false ]}`<br>Result: true<br><br>Logic: `{"or" : [ false, false ]}`<br>Result: false                                                                           |
| And         | Logic: `{"and" : [ true, false ]}`<br>Result: false<br><br>Logic: `{"and" : [ true, true ]}`<br>Result: true                                                                           |


#### Operators

| Operator               | Description                                                          | Context type | Example                                                                                                                                              |
| ---------------------- | -------------------------------------------------------------------- | ------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| Equals                 | Attribute equals the specified value, with type coercion.            | any          | Logic: `{ "==" : [1, 1] }`<br>Result: true<br><br>Logic: `{ "==" : [1, "1"] }`<br>Result: true                                                       |
| Strict equals          | Attribute equals the specified value, with strict comparison.        | any          | Logic: `{ "===" : [1, 1] }`<br>Result: true<br><br>Logic: `{ "===" : [1, "1"] }`<br>Result: false                                                    |
| Not equals             | Attribute doesn't equal the specified value, with type coercion.     | any          | Logic: `{ "!=" : [1, 2] }`<br>Result: true<br><br>Logic: `{ "!=" : [1, "1"] }`<br>Result: false                                                      |
| Strict not equal       | Attribute doesn't equal the specified value, with strict comparison. | any          | Logic: `{ "!==" : [1, 2] }`<br>Result: true<br><br>Logic: `{ "!==" : [1, "1"] }`<br>Result: true                                                     |
| Exists                 | Attribute is defined                                                 | any          | Logic: `{ "!!": [ "mike" ] }`<br>Result: true<br><br>Logic: `{ "!!": [ "" ] }`<br>Result: false                                                      |
| Not exists             | Attribute is not defined                                             | any          | Logic: `{"!": [ "mike" ] }`<br>Result: false<br><br>Logic: `{"!": [ "" ] }`<br>Result: true                                                          |
| Greater than           | Attribute is greater than the specified value                        | number       | Logic: `{ ">" : [2, 1] }`<br>Result: true<br><br>Logic: `{ ">" : [1, 2] }`<br>Result: false                                                          |
| Greater than or equals | Attribute is greater or equal to the specified value                 | number       | Logic: `{ ">=" : [2, 1] }`<br>Result: true<br><br>Logic: `{ ">=" : [1, 1] }`<br>Result: true                                                         |
| Less than              | Attribute is less than the specified value                           | number       | Logic: `{ "<" : [1, 2] }`<br>Result: true<br><br>Logic: `{ "<" : [2, 1] }`<br>Result: false                                                          |
| Less than or equals    | Attribute is less or equal to the specified value                    | number       | Logic: `{ "<=" : [1, 1] }`<br>Result: true<br><br>Logic: `{ "<=" : [2, 1] }`<br>Result: false                                                        |
| Between                | Attribute between the specified values                               | number       | Logic: `{ "<" : [1, 5, 10]}`<br>Result: true<br><br>Logic: `{ "<" : [1, 11, 10] }`<br>Result: false                                                  |
| Between inclusive      | Attribute between or equal to the specified values                   | number       | Logic: `{"<=" : [1, 1, 10] }`<br>Result: true<br><br>Logic: `{"<=" : [1, 11, 10] }`<br>Result: false                                                 |
| Contains               | Contains string                                                      | string       | Logic: `{ "in": ["Spring", "Springfield"] }`<br>Result: true<br><br>Logic: `{ "in":["Illinois", "Springfield"] }`<br>Result: false                   |
| Not contains           | Does not contain a string                                            | string       | Logic: `{ "!": { "in":["Spring", "Springfield"] } }`<br>Result: false<br><br>Logic: `{ "!": { "in":["Illinois", "Springfield"] } }`<br>Result: true  |
| In                     | Attribute is in an array of strings                                  | string       | Logic: `{ "in" : [ "Mike", ["Bob", "Mike"]] }`<br>Result: true<br><br>Logic: `{ "in":["Todd", ["Bob", "Mike"]] }`<br>Result: false                   |
| Not it                 | Attribute is not in an array of strings                              | string       | Logic: `{ "!": { "in" : [ "Mike", ["Bob", "Mike"]] } }`<br>Result: false<br><br>Logic: `{ "!": { "in":["Todd", ["Bob", "Mike"]] } }`<br>Result: true |

### Functions

| Conditional           | Description                               | Example                                                                                                                                                                      |
| --------------------- | ----------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Fractional evaluation | A deterministic percentage-based rollout  | A complete example can be found [here](./reusable_targeting_rules.md)                                                                                                        |
| Starts with           | Attribute starts with the specified value | Logic: `{ "starts_with" : [ "192.168.0.1", "192.168"] }`<br>Result: true<br><br>Logic: `{ "starts_with" : [ "10.0.0.1", "192.168"] }`<br>Result: false                       |
| Ends with             | Attribute ends with the specified value   | Logic: `{ "ends_with" : [ "noreply@example.com", "@example.com"] }`<br>Result: true<br><br>Logic: `{ "ends_with" : [ "noreply@example.com", "@test.com"] }`<br>Result: false |

## Shared evaluators

`$evaluators` is an **optional** property.
It's a collection of shared targeting configurations.
It can be used to reduce the number of duplicated targeting rule configurations.

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
          }, "binet", null
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
            "fractionalEvaluation": [
              "email",
              [
                "red",
                25
              ],
              [
                "blue",
                25
              ],
              [
                "green",
                25
              ],
              [
                "yellow",
                25
              ]
            ]
          }, null
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

## Examples

Sample configurations can be found at <https://github.com/open-feature/flagd/tree/main/config/samples>.
