---
description: flagd flag definition
---

# Flag Definitions

## Flags

`flags` is a **required** property.
The flags property is a top-level property that contains a collection of individual flags and their corresponding flag configurations.

```json
{
  "$schema": "https://flagd.dev/schema/v0/flags.json",
  "flags": {
    ...
  }
}
```

## Flag Definition

`flag key` is a **required** property.
The flag key **must** uniquely identify a flag so it can be used during flag evaluation.
The flag key **should** convey the intent of the flag.

```json
{
  "$schema": "https://flagd.dev/schema/v0/flags.json",
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
  "$schema": "https://flagd.dev/schema/v0/flags.json",
  "flags": {
    "new-welcome-banner": {
      "state": "ENABLED",
      "variants": {
        "on": true,
        "off": false
      },
      "defaultVariant": "off",
      "targeting": { 
        "if": [
          { "ends_with": [{ "var": "email" }, "@example.com"] },
          "on",
          "off"
        ]
      },
      "metadata": {
        "version": "17"
      }
    }
  },
  "metadata": {
    "team": "user-experience",
    "flagSetId": "ecommerce"
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
For example, to use a flag configured with boolean values the `/flagd.evaluation.v1.Service/ResolveBoolean` path should be used.
If another path, such as `/flagd.evaluation.v1.Service/ResolveString` is called, a type mismatch occurs and an error is returned.

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

`defaultVariant` is an **optional** property.
If provided, the value **must** match the name of one of the variants defined above.
The default variant is used unless a targeting rule explicitly overrides it.
If `defaultVariant` is omitted or null, flagd providers will revert to the code default for the flag in question if targeting is not defined or falls through.

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

Example of explicitly using the code default:

```json
"variants": {
  "on": true,
  "off": false
},
"defaultVariant": null
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
If no targeting rules are defined, the response reason will always be `STATIC`, this allows for the flag values to be cached, this behavior is described [here](specifications/providers.md#flag-evaluation-caching).

#### Variants Returned From Targeting Rules

The output of the targeting rule **must** match the name of one of the defined variants.
One exception to the above is that rules may return `true` or `false` which will map to the variant indexed by the equivalent string (`"true"`, `"false"`).
If a null value is returned by the targeting rule, the `defaultVariant` is used.
If `defaultVariant` is not defined, flagd providers will revert to the code default.
This can be useful for conditionally "exiting" targeting rules and falling back to the default (in this case the returned reason will be `DEFAULT`).
If an invalid variant is returned (not a string, `true`, or `false`, or a string that is not in the set of variants) the evaluation is considered erroneous.

See [Boolean Variant Shorthand](#boolean-variant-shorthand).

#### Evaluation Context

Evaluation context provides the attributes used by targeting rules to determine flag values.
Context can come from multiple sources, which are merged together before evaluation.

##### Context Sources

flagd supports three sources of evaluation context:

| Source                | Flag                           | Description                                                   |
| --------------------- | ------------------------------ | ------------------------------------------------------------- |
| Request body          | -                              | Context sent with each evaluation request                     |
| Static context        | `-X` / `--context-value`       | Key-value pairs added at startup, included in all evaluations |
| Header-mapped context | `-H` / `--context-from-header` | Maps HTTP/gRPC request headers to context keys                |

##### Request Body Context

Context included as part of the evaluation request.
For example, when accessing flagd via HTTP, the POST body may look like this:

```json
{
  "flagKey": "booleanFlagKey",
  "context": {
    "email": "noreply@example.com"
  }
}
```

This is the most common approach when the calling application has user or session information available.

##### Static Context (`-X` flag)

Static context values are specified at startup and automatically included in every evaluation.
This is useful for server-wide or environment-specific values that don't change per-request.

```shell
flagd start \
  --uri file:./flags.json \
  -X environment=production \
  -X region=us-east-1 \
  -X service=payment-api
```

**Use cases:**

- **Environment identification**: Different flag behavior for production vs staging (`-X environment=production`)
- **Regional configuration**: Apply region-specific rules (`-X region=eu-west-1`)
- **Service identification**: When multiple services share flag definitions (`-X service=checkout`)
- **Infrastructure metadata**: Cloud provider, datacenter, or cluster information

##### Header-Mapped Context (`-H` flag)

Header-mapped context extracts values from HTTP or gRPC request headers and adds them to the evaluation context.
This enables context-sensitive evaluation without modifying request bodies.

```shell
flagd start \
  --uri file:./flags.json \
  -H "X-User-Id=userId" \
  -H "X-User-Tier=tier"
```

With this configuration:

- A request with headers `X-User-Id: abc123` and `X-User-Tier: premium` will have `userId=abc123` and `tier=premium` in its evaluation context.

**Use cases:**

- **Gateway integration**: Extract user information from headers set by an API gateway or auth proxy
- **Multi-tenancy**: Use tenant ID headers for tenant-specific flag behavior

##### Context Merge Priority

When the same key appears in multiple context sources, values are merged with this priority (highest wins):

1. **Header-mapped context** (`-H` flag) - highest priority
2. **Static context** (`-X` flag)
3. **Request body context** - lowest priority

For example, with this configuration:

```shell
flagd start \
  --uri file:./flags.json \
  -X tier=basic \
  -H "X-User-Tier=tier"
```

| Request Body                    | Header                    | Resulting `tier` value           |
| ------------------------------- | ------------------------- | -------------------------------- |
| `{"context": {"tier": "free"}}` | (none)                    | `basic` (static overrides body)  |
| `{"context": {"tier": "free"}}` | `X-User-Tier: premium`    | `premium` (header overrides all) |
| (none)                          | `X-User-Tier: enterprise` | `enterprise`                     |

This priority order allows operators to enforce certain context values at the infrastructure level while still accepting client-provided context for other attributes.

##### Accessing Context in Targeting Rules

The evaluation context can be accessed in targeting rules using the `var` operation followed by the evaluation context property name.

| Description                                                    | Example                                              |
| -------------------------------------------------------------- | ---------------------------------------------------- |
| Retrieve property from the evaluation context                  | `#!json { "var": "email" }`                          |
| Retrieve property from the evaluation context or use a default | `#!json { "var": ["email", "noreply@example.com"] }` |
| Retrieve a nested property from the evaluation context         | `#!json { "var": "user.email" }`                     |

> For more information, see the `var` section in the [JsonLogic documentation](https://jsonlogic.com/operations.html#var).

See the [cheat sheet](./cheat-sheet.md#context-aware-evaluation) for practical examples of context-sensitive evaluation.

#### Conditions

Conditions can be used to control the logical flow and grouping of targeting rules.

| Conditional | Example                                                                                                                                                                                                  |
| ----------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| If          | Logic: `#!json {"if" : [ true, "yes", "no" ]}`<br>Result: `"yes"`<br><br>Logic: `#!json {"if" : [ false, "yes", "no" ]}`<br>Result: `"no"`                                                               |
| If else     | Logic: `#!json {"if" : [ false, "yes", false, "no", "maybe" ]}`<br>Result: `"maybe"`<br><br>Logic: `#!json {"if" : [ false, "yes", false, "no", false, "maybe", "who knows" ]}`<br>Result: `"who knows"` |
| Or          | Logic: `#!json {"or" : [ true, false ]}`<br>Result: `true`<br><br>Logic: `#!json {"or" : [ false, false ]}`<br>Result: `false`                                                                           |
| And         | Logic: `#!json {"and" : [ true, false ]}`<br>Result: `false`<br><br>Logic: `#!json {"and" : [ true, true ]}`<br>Result: `true`                                                                           |

#### Operations

Operations are used to take action on or compare properties retrieved from the context.
These are provided out-of-the-box by JsonLogic.
It's worth noting that JsonLogic operators never throw exceptions or abnormally terminate due to invalid input.
As long as a JsonLogic operator is structurally valid, it will return a falsy/nullish value.

| Operator               | Description                                                          | Context attribute type | Example                                                                                                                                                                |
| ---------------------- | -------------------------------------------------------------------- | ---------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Equals                 | Attribute equals the specified value, with type coercion.            | any                    | Logic: `#!json { "==" : [1, 1] }`<br>Result: `true`<br><br>Logic: `#!json { "==" : [1, "1"] }`<br>Result: `true`                                                       |
| Strict equals          | Attribute equals the specified value, with strict comparison.        | any                    | Logic: `#!json { "===" : [1, 1] }`<br>Result: `true`<br><br>Logic: `#!json { "===" : [1, "1"] }`<br>Result: `false`                                                    |
| Not equals             | Attribute doesn't equal the specified value, with type coercion.     | any                    | Logic: `#!json { "!=" : [1, 2] }`<br>Result: `true`<br><br>Logic: `#!json { "!=" : [1, "1"] }`<br>Result: `false`                                                      |
| Strict not equal       | Attribute doesn't equal the specified value, with strict comparison. | any                    | Logic: `#!json { "!==" : [1, 2] }`<br>Result: `true`<br><br>Logic: `#!json { "!==" : [1, "1"] }`<br>Result: `true`                                                     |
| Exists                 | Attribute is defined                                                 | any                    | Logic: `#!json { "!!": [ "mike" ] }`<br>Result: `true`<br><br>Logic: `#!json { "!!": [ "" ] }`<br>Result: `false`                                                      |
| Not exists             | Attribute is not defined                                             | any                    | Logic: `#!json {"!": [ "mike" ] }`<br>Result: `false`<br><br>Logic: `#!json {"!": [ "" ] }`<br>Result: `true`                                                          |
| Greater than           | Attribute is greater than the specified value                        | number                 | Logic: `#!json { ">" : [2, 1] }`<br>Result: `true`<br><br>Logic: `#!json { ">" : [1, 2] }`<br>Result: `false`                                                          |
| Greater than or equals | Attribute is greater or equal to the specified value                 | number                 | Logic: `#!json { ">=" : [2, 1] }`<br>Result: `true`<br><br>Logic: `#!json { ">=" : [1, 1] }`<br>Result: `true`                                                         |
| Less than              | Attribute is less than the specified value                           | number                 | Logic: `#!json { "<" : [1, 2] }`<br>Result: `true`<br><br>Logic: `#!json { "<" : [2, 1] }`<br>Result: `false`                                                          |
| Less than or equals    | Attribute is less or equal to the specified value                    | number                 | Logic: `#!json { "<=" : [1, 1] }`<br>Result: `true`<br><br>Logic: `#!json { "<=" : [2, 1] }`<br>Result: `false`                                                        |
| Between                | Attribute between the specified values                               | number                 | Logic: `#!json { "<" : [1, 5, 10]}`<br>Result: `true`<br><br>Logic: `#!json { "<" : [1, 11, 10] }`<br>Result: `false`                                                  |
| Between inclusive      | Attribute between or equal to the specified values                   | number                 | Logic: `#!json {"<=" : [1, 1, 10] }`<br>Result: `true`<br><br>Logic: `#!json {"<=" : [1, 11, 10] }`<br>Result: `false`                                                 |
| Contains               | Contains string                                                      | string                 | Logic: `#!json { "in": ["Spring", "Springfield"] }`<br>Result: `true`<br><br>Logic: `#!json { "in":["Illinois", "Springfield"] }`<br>Result: `false`                   |
| Not contains           | Does not contain a string                                            | string                 | Logic: `#!json { "!": { "in":["Spring", "Springfield"] } }`<br>Result: `false`<br><br>Logic: `#!json { "!": { "in":["Illinois", "Springfield"] } }`<br>Result: `true`  |
| In                     | Attribute is in an array of strings                                  | string                 | Logic: `#!json { "in" : [ "Mike", ["Bob", "Mike"]] }`<br>Result: `true`<br><br>Logic: `#!json { "in":["Todd", ["Bob", "Mike"]] }`<br>Result: `false`                   |
| Not in                 | Attribute is not in an array of strings                              | string                 | Logic: `#!json { "!": { "in" : [ "Mike", ["Bob", "Mike"]] } }`<br>Result: `false`<br><br>Logic: `#!json { "!": { "in":["Todd", ["Bob", "Mike"]] } }`<br>Result: `true` |

#### Custom Operations

These are custom operations specific to flagd and flagd providers.
They are purpose-built extensions to JsonLogic in order to support common feature flag use cases.
Consistent with built-in JsonLogic operators, flagd's custom operators return falsy/nullish values with invalid inputs.

| Function                           | Description                                         | Context attribute type                       | Example                                                                                                                                                                                                                                                                                            |
| ---------------------------------- | --------------------------------------------------- | -------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `fractional` (_available v0.6.4+_) | Deterministic, pseudorandom fractional distribution | string (bucketing value)                     | Logic: `#!json { "fractional" : [ { "var": "email" }, [ "red" , 50], [ "green" , 50 ] ] }` <br>Result: Pseudo randomly `red` or `green` based on the evaluation context property `email`.<br><br>Additional documentation can be found [here](./custom-operations/fractional-operation.md).        |
| `starts_with`                      | Attribute starts with the specified value           | string                                       | Logic: `#!json { "starts_with" : [ "192.168.0.1", "192.168"] }`<br>Result: `true`<br><br>Logic: `#!json { "starts_with" : [ "10.0.0.1", "192.168"] }`<br>Result: `false`<br>Additional documentation can be found [here](./custom-operations/string-comparison-operation.md).                      |
| `ends_with`                        | Attribute ends with the specified value             | string                                       | Logic: `#!json { "ends_with" : [ "noreply@example.com", "@example.com"] }`<br>Result: `true`<br><br>Logic: `#!json { ends_with" : [ "noreply@example.com", "@test.com"] }`<br>Result: `false`<br>Additional documentation can be found [here](./custom-operations/string-comparison-operation.md). |
| `sem_ver`                          | Attribute matches a semantic versioning condition   | string (valid [semver](https://semver.org/)) | Logic: `#!json {"sem_ver": ["1.1.2", ">=", "1.0.0"]}`<br>Result: `true`<br><br>Additional documentation can be found [here](./custom-operations/semver-operation.md).                                                                                                                              |

#### Targeting key

flagd and flagd providers map the [targeting key](https://openfeature.dev/specification/glossary#targeting-key) into the `"targetingKey"` property of the context used in rules.
For example, if the targeting key for a particular evaluation was set to `"5c3d8535-f81a-4478-a6d3-afaa4d51199e"`, the following expression would evaluate to `true`:

```json
"==": [
    {
        "var": "targetingKey"
    },
    "5c3d8535-f81a-4478-a6d3-afaa4d51199e"
]
```

#### $flagd properties in the evaluation context

Flagd adds the following properties to the evaluation context that can be used in the targeting rules.

| Property           | Description                                             | From version |
| ------------------ | ------------------------------------------------------- | ------------ |
| `$flagd.flagKey`   | the identifier for the flag being evaluated             | v0.6.4       |
| `$flagd.timestamp` | a Unix timestamp (in seconds) of the time of evaluation | v0.6.7       |

## Shared evaluators

`$evaluators` is an **optional** property.
It's a collection of shared targeting configurations used to reduce the number of duplicated targeting rules.

Example:

```json
{
  "$schema": "https://flagd.dev/schema/v0/flags.json",
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

## Metadata

Metadata can be defined at both the flag set (as a sibling of [flags](#flags)) and within each flag.
Flag metadata conveys arbitrary information about the flag or flag set, such as a version number, or the business unit that is responsible for the flag.
When flagd resolves flags, the returned [flag metadata](https://openfeature.dev/specification/types/#flag-metadata) is a merged representation of the metadata defined in the flag set, and the metadata defined in the flag, with the metadata defined in the flag taking priority.
See the [playground](/playground/?scenario-name=Flag+metadata) for an interactive example.

## Boolean Variant Shorthand

Since rules that return `true` or `false` map to the variant indexed by the equivalent string (`"true"`, `"false"`), you can use shorthand for these cases.

For example, this:

```json
{
  "$schema": "https://flagd.dev/schema/v0/flags.json",
  "flags": {
    "new-welcome-banner": {
      "state": "ENABLED",
      "variants": {
        "true": true,
        "false": false
      },
      "defaultVariant": "false",
      "targeting": { 
        "if": [
          { "ends_with": [{ "var": "email" }, "@example.com"] },
          "true",
          "false"
        ]
      }
    }
  }
}
```

can be shortened to this:

```json
{
  "$schema": "https://flagd.dev/schema/v0/flags.json",
  "flags": {
    "new-welcome-banner": {
      "state": "ENABLED",
      "variants": {
        "true": true,
        "false": false
      },
      "defaultVariant": "false",
      "targeting": { 
        "ends_with": [{ "var": "email" }, "@example.com"]
      }
    }
  }
}
```

## Examples

Sample configurations can be found at <https://github.com/open-feature/flagd/tree/main/config/samples>.
