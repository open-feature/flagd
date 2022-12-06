### Omitted Value From Response

When interacting directly with the flagd http(s) API and requesting a flag where the value is the type's default (e.g. `false` is the boolean type default) the value is omitted from the response.
This behaviour is a feature of [proto3](https://developers.google.com/protocol-buffers/docs/proto3#json), the motivation is to avoid sending default values over the wire as clients are able to interpret that no value received means default value (every bit counts right?).

## Examples

Command:

```sh
curl -s -X POST "localhost:30000/schema.v1.Service/ResolveBoolean" -d '{"flagKey":"simple-flag","context":{}}' -H "Content-Type: application/json"
```

Result:

```sh
{"reason":"DEFAULT","variant":"off"}
```