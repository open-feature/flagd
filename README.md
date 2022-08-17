# flagD

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
### Resolve a boolean value

Command:
```sh
curl -X POST "localhost:8013/flags/myBoolFlag/resolve/boolean"
```
Result:
```sh
{"value":true,"reason":"STATIC","variant":"on"}
```
<br />

### Resolve a string value

Command:
```sh
curl -X POST "localhost:8013/flags/myStringFlag/resolve/string"
```
Result:
```sh
{"value":"val1","reason":"STATIC","variant":"key1"}
```
<br />

### Resolve a integer value

Command:
```sh
curl -X POST "localhost:8013/flags/myIntFlag/resolve/int"
```
Result:
```sh
{"value":"1","reason":"STATIC","variant":"one"}
```
[Why is this `int` response a `string`?](./docs/http_int_response.md)
<br />
<br />

### Resolve a float value

Command:
```sh
curl -X POST "localhost:8013/flags/myFloatFlag/resolve/float"
```
Result:
```sh
{"value":1.23,"reason":"STATIC","variant":"one"}
```
<br />

### Resolve an object value

Command:
```sh
curl -X POST "localhost:8013/flags/myObjectFlag/resolve/object"
```
Result:
```sh
{"value":{"key":"val"},"reason":"STATIC","variant":"object1"}
```
<br />

### Resolve a boolean value with evaluation context

Command:
```sh
curl -X POST "localhost:8013/flags/isColorYellow/resolve/boolean" -d '{"color": "yellow"}'
```
Result:
```sh
{"value":true,"reason":"TARGETING_MATCH","variant":"on"}
```
<br />

### Return value type mismatch error

A type mismatch error is returned when the resolved value of a flag does not match the type requested. In the example below, the resolved value of `myBoolFlag` is a `boolean` but the request expects a `string` to be returned.

Command:
```sh
curl -X POST "localhost:8013/flags/myBoolFlag/resolve/string"
```
Result:
```sh
{"error_code":"TYPE_MISMATCH","reason":"ERROR"}
```
<br />

### Return flag not found error

The flag not found error is returned when flag key in the request doesn't match any configured flags.

Command:
```sh
curl -X POST "localhost:8013/flags/aMissingFlag/resolve/string"
```
Result:
```sh
{"error_code":"FLAG_NOT_FOUND","reason":"ERROR"}
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

### [Configuration](./docs/configuration.md)

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

### [Reusable targeting rules](./docs/reusable_targeting_rules.md)

### [Fractional Evaluation](./docs/fractional_evaluation.md)

### The people who make flagD great 💜

<a href="https://github.com/open-feature/flagd/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=open-feature/flagd" />
</a>
