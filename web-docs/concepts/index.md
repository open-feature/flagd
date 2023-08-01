# flagd basics

Your flagd journey will start by defining your feature flags.

flagd will then read those feature flags and make them available to your application.

Your application will interact with flagd via the [OpenFeature SDK](https://example.com) to retrieve flag values via the flagd API.

![flagd architecture](../images/flagd-logical-architecture.jpg)

## Defining feature flags

Flags can be defined in either [JSON](https://github.com/open-feature/flagd/blob/main/samples/example_flags.flagd.json) or [YAML](https://github.com/open-feature/flagd/blob/main/samples/example_flags.flagd.yaml) syntax and the values can be of different types.

Here are two flags, `flagOne` has `boolean` values and `flagTwo` has `string` values.

### flags represented as JSON

```json
{
  "flags": {
    "flagOne": {
      "state": "ENABLED",
      "variants": {
        "on": true,
        "off": false
      },
      "defaultVariant": "on",
      "targeting": {}
    },
    "flagTwo": {
      "state": "ENABLED",
      "variants": {
        "key1": "val1",
        "key2": "val2"
      },
      "defaultVariant": "key1",
      "targeting": {}
    }
  }
}
```

### flags represented as YAML

```yaml
flags:
  flagOne:
    state: ENABLED
    variants:
      'on': true
      'off': false
    defaultVariant: 'on'
    targeting:
  flagTwo:
    state: ENABLED
    variants:
      key1: val1
      key2: val2
    defaultVariant: 'key1'
    targeting:
```

## The structure of a flag

Each flag has:

- A flag key: `flagOne` and `flagTwo` above
- A state: `ENABLED` or `DISABLED`
- One or more possible `variants`. These are the possible values that a flag key can take.
- An optional `targeting` rule (explained below)

## Targeting rules

Imagine you are introducing a new feature. You create a flag with two possible variants: `on` and `off`. You want to safely roll out the feature.
Therefore the flags `defaultValue` is set to `off` for all users.

In other words, the new feature is disabled by default.

Now imagine you want to enable the feature, but only when both of these conditions are true:

- Logged in users
- The user email ends in `@example.com`

Rather than codifying that in your application, flagd targeting rules can be used. The flag definition below models this behaviour.

Your application is responsible for sending the `email` address via OpenFeature's context parameter (see below) and flagd will return the correct flag.

```json
{
    "flags": {
        "isFeatureEnabled": {
            "state": "ENABLED",
            "variants": {
                "on": true,
                "off": false
                },
            "defaultVariant": "off",
            "targeting": {
                "if": [{
                    "in": [
                        "@example.com",
                        {
                            "var": ["email"]
                        }]
                },
                "on", null]
            }
         }
    }
}
```

### Pseudo-code of application passing context

```js
// The second parameter is the default in case flagd is unavailable
featureAvailable = openFeature.getBooleanValue("isFeatureEnabled", false, {}) // false

// isFeatureEnabled for a logged in user with an email example@gmail.com
featureAvailable = openFeature.getBooleanValue("isFeatureEnabled", false, {"email": "example@gmail.com"}) // false

// isFeatureEnabled for a logged in user with an email someone@example.com
featureAvailable = openFeature.getBooleanValue("isFeatureEnabled", false, {"email": "someone@example.com"}) // true
```

## Fractional Evaluation

In some scenarios, it is desirable to use contextual information to segment the user population further and thus return dynamic values.

Look at the `headerColor` flag below. The `defaultVariant` is `red`, but the flag contains a targeting rule, meaning a fractional evaluation occurs when a context is passed and a key of `email` contains the value `@example.com`.

In this case, `25%` of the email addresses will receive `red`, `25%` will receive `blue`, and so on.

```json
{
    "flags": {
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
                "if": [{
                    "emailWithFaas": {
                        "in": ["@faas.com", {
                            "var": ["email"]
                            }]
                    }
                },
                {
                    "fractionalEvaluation": [ "email",
                        [ "red", 25 ], [ "blue", 25 ], [ "green", 25 ], [ "yellow", 25 ]
                    ]
                }, null
                ]
            }
        }
    }
}
```

### Fractional evaluations are sticky

Fractional evaluations are "sticky" and deterministic meaning that the same email address will always belong to the same "bucket" and thus always receive the same color.

This is true even if you run multiple flagd APIs completely independently.

See this page for more information on [flagd fractional evaluation logic](https://github.com/open-feature/flagd/blob/main/docs/configuration/fractional_evaluation.md).

## Other target specifiers

The example above shows the `in` keyword being used, but flagd is also compatible with:

- [starts_with](https://github.com/open-feature/flagd/blob/main/docs/configuration/string_comparison_evaluation.md#startswith-evaluation-configuration)
- [ends_with](https://github.com/open-feature/flagd/blob/main/docs/configuration/string_comparison_evaluation.md#endswith-evaluation-configuration)
- [sem_ver comparisons](https://github.com/open-feature/flagd/blob/main/docs/configuration/sem_ver_evaluation.md)

## flagd OpenTelemetry

flagd is fully compatible with OpenTelemetry:

- flagd exposes metrics at `http://localhost:8014/metrics`
- flagd can export metrics and traces to an OpenTelemetry collector.

See the [flagd OpenTelemetry](opentelemetry.md) page for more information.