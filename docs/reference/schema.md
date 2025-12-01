---
description: flagd flag and targeting schemas 
---

# Schema

## Flags

A comprehensive JSON schema is available for flagd configuration at [https://flagd.dev/schema/v0/flags.json](https://flagd.dev/schema/v0/flags.json).
It comprises definitions for flags as well as targeting.
You can use this schema to validate flagd configurations by using any JSON Schema validation library compliant with [JSON Schema draft-07](https://json-schema.org/draft-07/schema#).
_Additionally, most IDEs will automatically validate JSON documents if the document contains a `$schema` key and the schema is available at the specified URL_.

The example below is automatically validated in most IDEs:

```json hl_lines="2"
{
  "$schema": "https://flagd.dev/schema/v0/flags.json",
  "flags": {
    "basic-flag": {
      "state": "ENABLED",
      "variants": {
        "on": true,
        "off": false
      },
      "defaultVariant": "on"
    },
    "fractional-flag": {
      "state": "ENABLED",
      "variants": {
        "clubs": "clubs",
        "diamonds": "diamonds",
        "hearts": "hearts",
        "spades": "spades",
        "wild": "wild"
      },
      "defaultVariant": "wild",
      "targeting": {
        "fractional": [
          { "var": "email" },
          ["clubs", 25],
          ["diamonds", 25],
          ["hearts", 25],
          ["spades", 25]
        ]
      }
    }
  }
}
```

## Flagd Config

Flagd will return its flag configuration adhering to the `flagdConfig` schema [https://flagd.dev/schema/v0/flagd.json](https://flagd.dev/schema/v0/flagd.json).
This schema is defined in and refers to the flags schema above through the path `flags.json#/definitions/flagdConfig`.

```json hl_lines="2"
{
  "$schema": "https://flagd.dev/schema/v0/flagd.json",
  "flags": {
    "$schema": "https://flagd.dev/schema/v0/flags.json#/definitions/flagdConfig",
    "flags": [
      {
        "key": "basic-flag",
        "state": "ENABLED",
        "variants": {
          "on": true,
          "off": false
        },
        "defaultVariant": "on"
      },
      {
        "key": "fractional-flag",
        "state": "ENABLED",
        "defaultVariant": "default",
        "targeting": {
          "fractional": [
            { "var": "email" },
            ["clubs", 25],
            ["diamonds", 25],
            ["hearts", 25],
            ["spades", 25]
          ]
        },
        "variants": {
          "clubs": "clubs",
          "diamonds": "diamonds",
          "hearts": "hearts",
          "spades": "spades",
          "wild": "wild"
        }
      }
    ]
  }
}
```

## Targeting

In addition to the _flags_ schema, there's a schema available specifically for flagd _targeting rules_ at [https://flagd.dev/schema/v0/targeting.json](https://flagd.dev/schema/v0/targeting.json).
This validates only the `targeting` property of a flag.
**Please note that the flags schema also validates the targeting for each flag**, so it's not necessary to specifically use the targeting schema unless you wish to validate a targeting field individually.
