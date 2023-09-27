---
name: Quick Start
---

# Quick Start

Learn the basics of flagd from the comfort of your terminal.

## What you'll need

- Docker
- cURL

## Let's get started

### Download the flag definition

```shell
wget https://raw.githubusercontent.com/open-feature/flagd/main/docs/assets/demo.flagd.json
```

The flag definition file includes two feature flags.
The first one has the flag key `show-welcome-banner` and is a boolean type.
These types of feature flags are commonly used to gate access to a new feature using a conditional in code.
The second flag has the key `background-color` and is a multi-variant string.
These are commonly used for A/B/(n) testing and experimentation.

### Start flagd

```shell
docker run \
  --rm -it \
  --name flagd \
  -p 8013:8013 \
  -v $(pwd):/etc/flagd \
  ghcr.io/open-feature/flagd:latest start \
  --uri file:./etc/flagd/demo.flagd.json
```

??? "Tips for Windows users"
    In Windows, use WSL system for both the file location and Docker runtime.
    Mixed file systems does not work and this is a [limitation of Docker](https://github.com/docker/for-win/issues/8479).

### Evaluating a feature flag

Test it out by running the following cURL command in a separate terminal:

```shell
curl -X POST "http://localhost:8013/schema.v1.Service/ResolveBoolean" \
  -d '{"flagKey":"show-welcome-banner","context":{}}' -H "Content-Type: application/json"
```

You should see the following result:

```json
{
  "value": false,
  "reason": "STATIC",
  "variant": "off",
  "metadata": {}
}
```

### Enable the welcome banner

Open the `demo.flagd.json` file in a text editor and change the `defaultVariant` value from `off` to `on`.

Save and rerun the following cURL command:

```shell
curl -X POST "http://localhost:8013/schema.v1.Service/ResolveBoolean" \
  -d '{"flagKey":"show-welcome-banner","context":{}}' -H "Content-Type: application/json"
```

You should see the updated results:

```json
{
  "value": true,
  "reason": "STATIC",
  "variant": "on",
  "metadata": {}
}
```

!!! note ""

    Notice that flagd picked up the new flag definition without requiring a restart.

### Multi-variant feature flags

In some situations, a boolean value may not be enough.
That's where a multi-variant feature flag comes in handy.
In this section, we'll talk about a multi-variant feature flag can be used to control the background color of an application.

Save and rerun the following cURL command:

```shell
curl -X POST "http://localhost:8013/schema.v1.Service/ResolveString" \
  -d '{"flagKey":"background-color","context":{}}' -H "Content-Type: application/json"
```

You should see the updated results:

```json
{
  "value": "#FF0000",
  "reason": "STATIC",
  "variant": "red",
  "metadata": {}
}
```

### Add a targeting rule

Imagine that we're testing out a new color scheme internally.
Employees should see the green background color while customers should continue seeing red.
This can be accomplished in flagd using targeting rules.

Open the `demo.flagd.json` file in a text editor and extend the `background-color` to include a targeting rule.

``` json hl_lines="19-32"
{
  "flags": {
    "show-welcome-banner": {
      "state": "ENABLED",
      "variants": {
        "on": true,
        "off": false
      },
      "defaultVariant": "off"
    },
    "background-color": {
      "state": "ENABLED",
      "variants": {
        "red": "#FF0000",
        "blue": "#0000FF",
        "green": "#00FF00",
        "yellow": "#FFFF00"
      },
      "defaultVariant": "red",
      "targeting": {
        "if": [
          {
            "===": [
              {
                "var": "company"
              },
              "initech"
            ]
          },
          "green"
        ]
      }
    }
  }
}
```

The evaluation context contains arbitrary attributes that targeting rules can operate on, and can be included in each feature flag evaluation.
This rule will return the `green` variant if the `company` included in the _evaluation context_ matches `initech`.
If there isn't a match, the `defaultVariant` is returned.

#### Test as a customer

Let's confirm that customers are still seeing the `red` variant.

```shell
curl -X POST "http://localhost:8013/schema.v1.Service/ResolveString" \
  -d '{"flagKey":"background-color","context":{"company": "stark industries"}}' -H "Content-Type: application/json"
```

You should see the updated results:

```json
{
  "value": "#FF0000",
  "reason": "DEFAULT",
  "variant": "red",
  "metadata": {}
}
```

#### Test as an employee

Let's confirm that employees of Initech are seeing the updated variant.

Run the following cURL command in the terminal.

```shell
curl -X POST "http://localhost:8013/schema.v1.Service/ResolveString" \
  -d '{"flagKey":"background-color","context":{"company": "initech"}}' -H "Content-Type: application/json"
```

You should see the updated results:

```json
{
  "value": "#00FF00",
  "reason": "TARGETING_MATCH",
  "variant": "green",
  "metadata": {}
}
```

Notice that the `green` variant is returned and the reason is `TARGETING_MATCH`.

## Summary

In this guide, we configured flagd to use a local flag configuration.
We then performed flag evaluation using cURL to see how updating the flag definition affects the output.
We also explored how evaluation context can be used within a targeting rule to personalize the output.
This is just scratching the surface of flagd's capabilities.
Check out the [concepts section](./concepts//feature-flagging.md) to learn about the use cases enabled by flagd.
