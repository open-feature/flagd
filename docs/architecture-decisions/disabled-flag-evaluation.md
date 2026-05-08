---
# Valid statuses: draft | proposed | rejected | accepted | superseded
status: draft
author: Parth Suthar (@suthar26)
created: 2026-04-01
updated: 2026-04-28
---
# Treat Disabled Flag Evaluation as Successful with Reason DISABLED

Today, evaluating a disabled flag in flagd produces an error (`reason=ERROR`, `errorCode=FLAG_DISABLED`). We propose returning a successful evaluation with `reason=DISABLED` and no value, so the calling SDK falls back to the application's code default. A flag that does not exist still produces `FLAG_NOT_FOUND`. This matches how the [OpenFeature specification](https://openfeature.dev/specification/types/#resolution-reason) defines `DISABLED`: a successful evaluation, not a failure.

## Background

flagd's current behavior treats `state: DISABLED` as an error and surfaces that error through gRPC, OFREP, and in-process providers. Several issues follow from this.

The OpenFeature specification lists `DISABLED` as a resolution reason and describes it as *"the resolved value was the result of the flag being disabled in the management system."* Errors are described separately. Treating disabled as an error therefore conflicts with the spec.

`FLAG_DISABLED` is also not a valid error code anywhere it is used. It is missing from the OpenFeature error code list and from the OFREP `evaluationFailure` schema, which only allows `PARSE_ERROR`, `TARGETING_KEY_MISSING`, `INVALID_CONTEXT`, and `GENERAL`. The OFREP success schema, on the other hand, does allow `reason=DISABLED`. The current behavior violates both specs at once.

The error path also conflates two different situations.
A missing flag is usually a deployment or configuration mistake that an operator wants to know about.
A disabled flag is an intentional operational state, often used during incident remediation, environment-scoped rollouts, or features that are not yet ready.
Today both surface as `connect.CodeNotFound` on gRPC v1, and OFREP rewrites `FLAG_DISABLED` into `FLAG_NOT_FOUND` in its structured error response, leaving the disabled distinction visible only in a free-text field.
Clients cannot reliably tell the two apart.

These collapsed error paths hurt observability. Operators who disable a flag deliberately see false error signals in dashboards and alerts; if they suppress those alerts, they lose visibility into flag state altogether. The same problem appears in flag-set-based deployments, where a flag may legitimately be disabled in one set and active in another, and treating that as an exception forces normal operations through error-handling code.

Related reading: [OpenFeature resolution reasons](https://openfeature.dev/specification/types/#resolution-reason), the [flagd flag definitions reference](https://flagd.dev/reference/flag-definitions/), and the prior [ADR on explicit code defaults](./support-code-default.md), which establishes the field-omission pattern reused below.

## Requirements

A disabled flag should evaluate successfully with `reason=DISABLED` on every surface: gRPC v1, gRPC v2, OFREP, and in-process.
The resolved value should follow the same field-omission pattern as the code-default ADR, so the SDK uses the application's code default; only the `reason` differs.
Unknown flag keys must continue to return `FLAG_NOT_FOUND`.
The `DISABLED` reason must not feed into provider or SDK error paths, and bulk evaluation must include disabled flags in the response rather than skipping them.
Telemetry should record these as successful evaluations.
No change to existing flag configuration files is required.

## Considered options

1. Successful evaluation with `reason=DISABLED`, value omitted so the SDK falls back to code defaults.
2. Successful evaluation with `reason=DISABLED`, returning the configured `defaultVariant` value.
3. Successful evaluation with `reason=DEFAULT`, treating disabled as a special case of "no targeting matched".
4. Keep the current error behavior and document the spec divergence.

We propose option 1. Option 2 still lets the management system pick a value, which contradicts the OpenFeature description of `DISABLED` and prevents the SDK from using its real fallback path. Option 3 hides the disabled state from clients and metrics, removing the very signal that motivated the change. Option 4 leaves the OFREP and OpenFeature spec violations in place and keeps the missing-vs-disabled confusion described above.

## Proposal

When a flag exists with `state: DISABLED`, the evaluator returns a successful result with no value and no variant, `reason=DISABLED`, and the usual flag and flag-set metadata.
The omission of `value` and `variant` is the same mechanism used in the code-default ADR; the SDK treats omission as a signal to use the application default.
Targeting rules are not evaluated, so reasons that describe targeting outcomes (`STATIC`, `DEFAULT`, `SPLIT`, `TARGETING_MATCH`) never apply to a disabled flag.
`ERROR` continues to mean a real failure such as a parse error or type mismatch.

The behavior change is uniform across surfaces.
The single-flag and bulk evaluation paths both include disabled flags with `reason=DISABLED` instead of erroring or skipping them.
OFREP returns a success payload rather than an error response shaped like `FLAG_NOT_FOUND`.
On the wire, gRPC and OFREP omit the value and variant fields.
In-process providers already receive `"state": "DISABLED"` in the sync payload, so the change there is in the per-language evaluator: it must treat that state the same way as the core flagd evaluator.
The provider and core changes need to ship together so that integrators see consistent behavior.

A typical OFREP single-flag response looks like this. The status moves from HTTP 404 (the current `FLAG_NOT_FOUND` rewrite) to HTTP 200, since the evaluation now succeeds.

```json
{
  "key": "my-feature",
  "reason": "DISABLED",
  "metadata": { "flagSetId": "my-app" }
}
```

File-level changes are out of scope for this ADR and will be tracked in the implementation PRs.

## Consequences

The main benefits are spec alignment with both OpenFeature and OFREP, a clear distinction between missing and disabled flags, less noisy error metrics, and visibility into disabled flags in bulk responses. Operators get a clean signal that a flag is intentionally off, and applications can apply their normal default-value logic without going through an error branch.

The main cost is that this is a breaking change. Clients that switch on `FLAG_DISABLED` in error handling, alerting, or HTTP 404 responses from OFREP single-flag evaluation will need to change. Bulk responses also grow when a flag set contains many disabled flags. The rollout has to be coordinated across the flagd core, language SDKs and providers, and the testbed.

As a side effect, the existing `FlagDisabledErrorCode` plumbing in the error formatters can be removed once the evaluator no longer produces it.

## Testing

Coverage for this change should live in the [flagd testbed](https://github.com/open-feature/flagd-testbed) so every SDK and provider can verify behavior against the same scenarios. We need cases for single-flag and bulk evaluation on gRPC v1, gRPC v2, OFREP, and in-process, including the case where a flag is disabled in one flag set and enabled in another.

## Versioning and migration

flagd is pre-1.0, so this ships as a minor-version bump with the breaking change called out in the release notes rather than as a long-running compatibility mode. Operators and client authors should:

- Replace `FLAG_DISABLED` error handling with checks for a successful evaluation whose reason is `DISABLED`.
- Update OFREP and HTTP clients that branched on a 404 status for disabled single-flag evaluation.
- Audit dashboards, alerts, and log parsers keyed on disabled-flag errors.

The obsolete error-code paths are removed in the same release. Keeping them around does not preserve any reachable behavior once the evaluator stops producing the error.

## Open questions

- Should bulk evaluation expose an option to omit disabled flags, for clients that prefer smaller payloads over visibility?

## More information

- [OpenFeature resolution reasons](https://openfeature.dev/specification/types/#resolution-reason)
- [OpenFeature error codes](https://openfeature.dev/specification/types/#error-code)
- [flagd flag definitions](https://flagd.dev/reference/flag-definitions/)
- [ADR: Support explicit code default values](./support-code-default.md)
- [flagd testbed](https://github.com/open-feature/flagd-testbed)
