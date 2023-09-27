# Troubleshooting flagd

## Debugging Evaluations

If a flag or targeting rule isn't proceeding the way you'd expect, you may want to enable more verbose logging.

flagd and flagd providers typically have debug or verbose logging modes that you can use for this sort of troubleshooting.
You can do this in the standalone version of flagd by starting it with the `--debug` flag (see [CLI](./reference/flagd-cli/flagd/flagd.md) for more information).

_In-process_ providers which embed the flag evaluation engine use a logging consistent with their implementation language and SDK.
See your provider's documentation for details on how to enable verbose logging.

The [detailed evaluation](https://openfeature.dev/docs/reference/concepts/evaluation-api#detailed-evaluation) functions can also be helpful in understanding why an evaluation proceeded a particular way.

---

## HTTP Integer Response Behavior

Why is my `int` response a `string`?
Command:

```sh
curl -X POST "localhost:8013/schema.v1.Service/ResolveInt" -d '{"flagKey":"myIntFlag","context":{}}' -H "Content-Type: application/json"
```

Result:

```sh
{"value":"1","reason":"DEFAULT","variant":"one"}
```

When interacting directly with the flagd HTTP api and requesting an `int` the response type will be a `string`.
This behaviour is introduced by [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway), which uses [proto3 json mapping](https://developers.google.com/protocol-buffers/docs/proto3#json) to build the response object.
If a number value is required, and none of the provided SDK's can be used, then it is recommended to use the `float64` endpoint instead:  

Command:

```sh
curl -X POST "localhost:8013/schema.v1.Service/ResolveFloat" -d '{"flagKey":"myIntFlag","context":{}}' -H "Content-Type: application/json"
```

Result:

```sh
{"value":1.23,"reason":"DEFAULT","variant":"one"}
```
