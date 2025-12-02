---
description: troubleshooting flagd
---

# Troubleshooting flagd

## Debugging Evaluations

If a flag or targeting rule isn't proceeding the way you'd expect, you may want to enable more verbose logging.

flagd and flagd providers typically have debug or verbose logging modes that you can use for this sort of troubleshooting.
You can do this in the standalone version of flagd by starting it with the `--debug` flag (see [CLI](./reference/flagd-cli/flagd.md) for more information) or through `FLAGD_DEBUG` env variable.

_In-process_ providers which embed the flag evaluation engine use a logging consistent with their implementation language and SDK.
See your provider's documentation for details on how to enable verbose logging.

The [detailed evaluation](https://openfeature.dev/docs/reference/concepts/evaluation-api#detailed-evaluation) functions can also be helpful in understanding why an evaluation proceeded a particular way.

---

## HTTP Integer Response Behavior

Why is my `int` response a `string`?
Command:

```sh
curl -X POST "localhost:8013/flagd.evaluation.v1.Service/ResolveInt" -d '{"flagKey":"myIntFlag","context":{}}' -H "Content-Type: application/json"
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
curl -X POST "localhost:8013/flagd.evaluation.v1.Service/ResolveFloat" -d '{"flagKey":"myIntFlag","context":{}}' -H "Content-Type: application/json"
```

Result:

```sh
{"value":1.23,"reason":"DEFAULT","variant":"one"}
```

---

## Extra (Duplicate) Events in File Syncs

When using the file sync, you may notice more events than you'd expect.
This is attributable to nuances in text editors and the OS/filesystem.
Most editors will cause a few filesystem events on a save, for instance, not just one as you might expect.
Additionally, most filesystem operations are not atomic and also will result in multiple events.
Generally speaking, updating a symbolic link will result in only a single event, and may even be atomic on your filesystem/OS.
In fact, this is how Kubernetes handles changes to mounted ConfigMaps (OpenFeature Operator takes advantage of this fact in it's `file` mode).
**It's recommended if you're using the file sync in production, you use a symbolic link for the watched file, and update its contents by changing its target.**

---

## Received unexpected EOS on empty DATA frame from server

This error message indicates that HTTP2 requests are being blocked (gRPC runs over HTTP2).
You may need to explicitly allow HTTP2 or gRPC in your platform if you're using the [sync gRPC service](./reference/specifications/protos.md#syncv1sync_serviceproto).

!!! note

    HTTP2 _is not_ strictly for the flag [evaluation gRPC service](./reference/specifications/protos.md#schemav1schemaproto), which is exposed both as a gRPC service and a RESTful HTTP/1.1 service, thanks to the [connect protocol](https://connectrpc.com/docs/protocol/).

---

## Selector Issues

### No Flags Returned with Selector

**Problem**: Provider returns no flags when using a selector.

**Debugging Steps:**

- Verify `flagSetId` in selector matches flag configuration exactly
- Check selector syntax: `flagSetId=my-app` (not `flagSetId:my-app`)
- Test without selector to confirm flags exist

### Wrong Flags Returned

**Problem**: Selector returns unexpected flags.

**Debugging Steps:**

- Check for flag-level `flagSetId` overrides in individual flags
- Verify header precedence: `Flagd-Selector` header overrides request body
- Use metadata reflection to see what selector was actually applied

### Selector Ignored

**Problem**: Selector appears to be ignored, all flags returned.

**Debugging Steps:**

- Verify selector syntax is correct (`key=value` format)
- Check if provider configuration has a selector that overrides requests
- Ensure selector value is not empty (`flagSetId=` returns all flags without flagSetId)

**Debug with metadata reflection:**

```bash
curl -H "Flagd-Selector: flagSetId=my-app" \
  http://localhost:8014/ofrep/v1/evaluate/flags
# Check response metadata to see parsed selector
```
