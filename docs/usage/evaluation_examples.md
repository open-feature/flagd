### Evaluation Examples Using curl

## Setup

1. Download sample flag configuration:

    ```shell
    curl https://raw.githubusercontent.com/open-feature/flagd/main/config/samples/example_flags.json -o example_flags.json
    ```

1. Run one of the following commands, depending on how [flagd was installed](https://github.com/open-feature/flagd/blob/main/docs/usage/getting_started.md):
   - binary:

    ```shell
    flagd start -f example_flags.json
    ```

   - Docker:

    ```shell
    docker run -p 8013:8013 -v $(pwd)/:/etc/flagd/ -it --pull=always ghcr.io/open-feature/flagd:latest start --uri file:./etc/flagd/example_flags.json
    ```

1. Changes made in `example_flags.json` will immediately take affect. Go ahead, give a shot!

Flagd is now ready to perform flag evaluations over either HTTP or gRPC. In this example, we'll utilize HTTP via cURL.

### Resolve a boolean value

Command:

```sh
curl -X POST "localhost:8013/schema.v1.Service/ResolveBoolean" -d '{"flagKey":"myBoolFlag","context":{}}' -H "Content-Type: application/json"
```

Result:

```sh
{"value":true,"reason":"DEFAULT","variant":"on"}
```

<br />

### Resolve a string value

Command:

```sh
curl -X POST "localhost:8013/schema.v1.Service/ResolveString" -d '{"flagKey":"myStringFlag","context":{}}' -H "Content-Type: application/json"
```

Result:

```sh
{"value":"val1","reason":"DEFAULT","variant":"key1"}
```

<br />

### Resolve a integer value

Command:

```sh
curl -X POST "localhost:8013/schema.v1.Service/ResolveInt" -d '{"flagKey":"myIntFlag","context":{}}' -H "Content-Type: application/json"
```

Result:

```sh
{"value":"1","reason":"DEFAULT","variant":"one"}
```
[Why is this int response a string](https://github.com/open-feature/flagd/blob/main/docs/help/http_int_response.md)
<br />
<br />

### Resolve a float value

Command:

```sh
curl -X POST "localhost:8013/schema.v1.Service/ResolveFloat" -d '{"flagKey":"myFloatFlag","context":{}}' -H "Content-Type: application/json"
```

Result:

```sh
{"value":1.23,"reason":"DEFAULT","variant":"one"}
```

<br />

### Resolve an object value

Command:

```sh
curl -X POST "localhost:8013/schema.v1.Service/ResolveObject" -d '{"flagKey":"myObjectFlag","context":{}}' -H "Content-Type: application/json"
```

Result:

```sh
{"value":{"key":"val"},"reason":"DEFAULT","variant":"object1"}
```

<br />

### Resolve a boolean value with evaluation context

Command:

```sh
curl -X POST "localhost:8013/schema.v1.Service/ResolveBoolean" -d '{"flagKey":"isColorYellow","context":{"color":"yellow"}}' -H "Content-Type: application/json"
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
curl -X POST "localhost:8013/schema.v1.Service/ResolveString" -d '{"flagKey":"myBoolFlag","context":{}}' -H "Content-Type: application/json"
```

Result:

```sh
{"code":"invalid_argument","message":"TYPE_MISMATCH"}
```

<br />

### Return flag not found error

The flag not found error is returned when flag key in the request doesn't match any configured flags.

Command:

```sh
curl -X POST "localhost:8013/schema.v1.Service/ResolveBoolean" -d '{"flagKey":"aMissingFlag","context":{}}' -H "Content-Type: application/json"
```

Result:

```sh
{"code":"not_found","message":"FLAG_NOT_FOUND"}
```