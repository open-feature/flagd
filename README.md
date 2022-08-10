# flagD

<img src="images/icon.png" width="350px"/>

![build](https://img.shields.io/github/workflow/status/open-feature/flagd/ci)
![goversion](https://img.shields.io/github/go-mod/go-version/open-feature/flagd/main)
![version](https://img.shields.io/badge/version-pre--alpha-green)
![status](https://img.shields.io/badge/status-not--for--production-red)

Flagd is a simple command line tool for fetching and presenting feature flags to services. It is designed to conform to OpenFeature schema for flag definitions.

<img src="images/of-flagd-0.png" width="560">

## Example usage

1. Generate the prerequisites `make generate`
2. Build the flagd binary: `make build`
3. Start the process: `./flagd start -f config/samples/example_flags.json --service-provider http --sync-provider filepath`

This now provides an accessible http or [https](#https) endpoint for the flags:

```
$ curl -X POST "localhost:8013/flags/myBoolFlag/resolve/boolean"
// {"value":true,"reason":"STATIC","variant":"on"}

$ curl -X POST "localhost:8013/flags/myStringFlag/resolve/string"
// {"value":"val1","reason":"STATIC","variant":"key1"}

$ curl -X POST "localhost:8013/flags/myIntFlag/resolve/int"
// {"value":"1","reason":"STATIC","variant":"one"}

$ curl -X POST "localhost:8013/flags/myFloatFlag/resolve/float"
// {"value":1.23,"reason":"STATIC","variant":"one"}

$ curl -X POST "localhost:8013/flags/myObjectFlag/resolve/object"
// {"value":{"key":"val"},"reason":"STATIC","variant":"object1"}

$ curl -X POST "localhost:8013/flags/isColorYellow/resolve/boolean" -d '{"color": "yellow"}'
// {"value":true,"reason":"TARGETING_MATCH","variant":"on"}

$ curl -X POST "localhost:8013/flags/myBoolFlag/resolve/string"
// {"error_code":"TYPE_MISMATCH","reason":"ERROR"}

$ curl -X POST "localhost:8013/flags/aMissingFlag/resolve/string"
// {"error_code":"FLAG_NOT_FOUND","reason":"ERROR"}
```

### https

When it is desired to use TLS for increased security, flagD can be started with the following cert and key information.

`flagd start --server-cert-path ./example-cert.pem --server-key-path ./example-key.pem`

This enables you to use an upgraded connection for the previous example requests, such as the following:

```
$ curl -X POST "https://localhost:8013/flags/myBoolFlag/resolve/boolean"
// {"value":true,"reason":"STATIC","variant":"on"}
```

## Multiple source example
Multiple providers can be supplied as the following:
```
./flagd start -f config/samples/example_flags.json -f config/samples/example_flags_secondary.json --service-provider http --sync-provider filepath
```
In case of collision between flags definition the priority goes to the later (e.g. example_flags < example_flags_secondary).


## Installation

### Systemd

To install as a systemd service run `sudo make install` this will place the binary by default in `/usr/local/bin`

There will also be a default provider and sync enabled ( http / filepath ) both of which can be modified in the flagd.service.

Validation can be run with `systemctl status flagd`
And result similar to below will be seen

```
● flagd.service - "A generic feature flag daemon"
     Loaded: loaded (/etc/systemd/system/flagd.service; disabled; vendor preset: enabled)
     Active: active (running) since Mon 2022-05-30 12:19:55 BST; 5min ago
   Main PID: 64610 (flagd)
      Tasks: 7 (limit: 4572)
     Memory: 1.4M
     CGroup: /system.slice/flagd.service
             └─64610 /usr/local/bin/flagd start -f=/etc/flagd/flags.json

May 30 12:19:55 foo systemd[1]: Started "A generic feature flag daemon".
```

### Running in a container

1. `IMG=flagd-local make docker-build`
2. `docker run -p 8013:8013 -it flagd-local start --uri ./examples/example_flags.json`

## Targeting Rules

The `flag` object has a field named `"targeting"`, this can be populated with a [JsonLogic](https://jsonlogic.com/) rule. Any data
in the body of a flag evaluation call is processed by the JsonLogic rule to determine the result of flag evaluation.
If this result is `null` or an invalid (undefined) variant then the default variant is returned.

JsonLogic provides a [playground](https://jsonlogic.com/play.html) for evaluating your rules against data.
 
<u>Example</u>

A flag is defined as such:
```json
{
  "flags": {
    "isColorYellowFlag": {
      "state": "ENABLED",
      "variants": {
        "on": true,
        "off": false
      },
      "defaultVariant": "off",
      "targeting": {
        "if": [
          {
            "==": [
              {
                "var": [
                  "color"
                ]
              },
              "yellow"
            ]
          },
          "on",
          "off"
        ]
      }
    }
  }
}
```

The rule provided returns `"on"` if `var color == "yellow"` and `"off"` otherwise:

```shell
$ curl -X POST "localhost:8013/flags/isColorYellow/resolve/boolean" -d '{"color": "yellow"}'
```
returns
```json
{"value":true,"reason":"TARGETING_MATCH","variant":"on"}
```

whereas 
```shell
$ curl -X POST "localhost:8013/flags/isColorYellow/resolve/boolean" -d '{"color": "white"}'
```
returns
```json
{"value":true,"reason":"TARGETING_MATCH","variant":"off"}
```

### Reusable targeting rules

At the same level as the `flags` key one can define an `$evaluators` object. Each object defined under `$evaluators` is 
a reusable targeting rule. In any targeting rule one can reference a defined reusable targeting rule, foo, like so:
`"$ref": "foo"`

<u>Example</u>

Flags/evaluators defined as such:

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

becomes (once the `$evaluators` have been substituted):

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
                "in": ["@faas.com", {
                "var": ["email"]
              }]
              }, "binet", null
            ]
          }
        }
    }
}
```

### The people who make flagD great 💜

<a href="https://github.com/open-feature/flagd/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=open-feature/flagd" />
</a>
