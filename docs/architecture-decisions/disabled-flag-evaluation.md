---
# Valid statuses: draft | proposed | rejected | accepted | superseded
status: draft
author: Parth Suthar (@suthar26)
created: 2026-04-01
updated: 2026-04-15
---
# Treat Disabled Flag Evaluation as Successful with Reason DISABLED

This ADR proposes changing flagd's handling of disabled flags from returning an error (`reason=ERROR`, `errorCode=FLAG_DISABLED`) to returning a successful evaluation with `reason=DISABLED` that defers to code defaults. A disabled flag means the flag management system has nothing to say — the application should use its code-defined default value. A flag that does not exist in a flag set should remain a `FLAG_NOT_FOUND` error.
This aligns flagd with the [OpenFeature specification's resolution reasons](https://openfeature.dev/specification/types/#resolution-reason), which defines `DISABLED` as a valid resolution reason for successful evaluations.

## Background

flagd currently treats the evaluation of a disabled flag as an error. When a flag exists in the store but has `state: DISABLED`, the evaluator returns `reason=ERROR` with `errorCode=FLAG_DISABLED`. This error propagates through every surface — gRPC, OFREP, and in-process providers — resulting in the caller receiving an error response rather than a resolved value.

This is problematic for several reasons:

1. **Spec misalignment**: The [OpenFeature specification](https://openfeature.dev/specification/types/#resolution-reason) explicitly defines `DISABLED` as a resolution reason with the description: *"The resolved value was the result of the flag being disabled in the management system."* This implies a successful evaluation that communicates the flag's disabled state, not an error.
2. **OFREP masks the disabled state**: The OFREP response handler in `core/pkg/service/ofrep/models.go` rewrites `FLAG_DISABLED` to `FLAG_NOT_FOUND` in the structured error response (while only preserving the "is disabled" distinction in the free-text `errorDetails` string). This means OFREP clients cannot programmatically distinguish between a flag that doesn't exist and one that was intentionally disabled.
3. **Conflation of "missing" and "disabled"**: gRPC v1 maps both `FLAG_NOT_FOUND` and `FLAG_DISABLED` to `connect.CodeNotFound`. These are semantically different situations: a missing flag is a configuration or deployment error, while a disabled flag is an intentional operational decision (incident remediation, environment-specific rollout, not-yet-ready feature).
4. **Non-standard error code**: `FLAG_DISABLED` is not in the [OpenFeature specification's error code list](https://openfeature.dev/specification/types/#error-code) (`PROVIDER_NOT_READY`, `FLAG_NOT_FOUND`, `PARSE_ERROR`, `TYPE_MISMATCH`, `TARGETING_KEY_MISSING`, `INVALID_CONTEXT`, `PROVIDER_FATAL`, `GENERAL`), nor is it a valid `errorCode` in the [OFREP `evaluationFailure` schema](https://github.com/open-feature/protocol/blob/main/service/openapi.yaml)
(which only allows `PARSE_ERROR`, `TARGETING_KEY_MISSING`, `INVALID_CONTEXT`, `GENERAL`). Conversely, `DISABLED` *is* already a valid `reason` in OFREP's [`evaluationSuccess` schema](https://github.com/open-feature/protocol/blob/main/service/openapi.yaml). flagd's current treatment of disabled flags as errors is a spec violation on both the OpenFeature and OFREP sides.
5. **Loss of observability**: When disabled flags are treated as errors, they pollute error metrics and alerting. Operators who disable a flag for legitimate reasons (ongoing incident remediation, feature not ready for an environment) see false error signals. Conversely, if they suppress these errors, they lose visibility into flag state entirely. A successful evaluation with `reason=DISABLED` would give operators a clean signal without noise.
6. **Flag set use cases**: In multi-flag-set deployments, a flag may exist in a shared definition but be disabled in certain flag sets (e.g., disabled for `staging` but enabled for `production`). Treating this as an error forces the application into error-handling paths when the flag is simply not active — a normal operational state, not an exceptional one.

Related context:

- [OpenFeature Specification - Resolution Reasons](https://openfeature.dev/specification/types/#resolution-reason)
- [flagd Flag Definitions Reference](https://flagd.dev/reference/flag-definitions/)
- [ADR: Support Explicit Code Default Values](./support-code-default.md) — establishes the field-omission pattern for code-default deferral

## Requirements

- Evaluating a disabled flag must return a successful response with `reason=DISABLED` across all surfaces (gRPC v1, gRPC v2, OFREP, in-process)
- The resolved value for a disabled flag must always defer to the application's code default — the server omits value and variant fields using the same field-omission pattern established in the [code-default ADR](./support-code-default.md), but with `reason=DISABLED` instead of `reason=DEFAULT`
- Evaluating a flag key that does not exist in the store or flag set must remain a `FLAG_NOT_FOUND` error
- `reason=DISABLED` must be distinct from all other reasons (`STATIC`, `DEFAULT`, `SPLIT`, `TARGETING_MATCH`, `ERROR`) and must not trigger error-handling paths in providers or SDKs
- Bulk evaluation (`ResolveAll`) must include disabled flags in the response with `reason=DISABLED`, rather than silently omitting them
- Telemetry and metrics must record disabled flag evaluations as successful (non-error) with the `DISABLED` reason
- Existing flag configurations must continue to work without modification (backward compatible at the configuration level)

## Considered Options

- **Option 1: Successful evaluation with `reason=DISABLED` deferring to code defaults** — Disabled flags evaluate successfully, always deferring to the application's code-defined default value with `reason=DISABLED`
- **Option 2: Successful evaluation with `reason=DISABLED` returning `defaultVariant`** — Disabled flags evaluate successfully, returning the flag configuration's `defaultVariant` value
- **Option 3: Return a successful evaluation with** `reason=DEFAULT` — Treat disabled flags as if they had no targeting, collapsing the disabled state into the default reason
- **Option 4: Status quo** — Keep the current error behavior and document it as intentional divergence from the OpenFeature spec

## Proposal

We propose **Option 1: Successful evaluation with `reason=DISABLED` deferring to code defaults**.

When a flag exists in the store but has `state: DISABLED`, the evaluator should return a successful evaluation that always defers to the application's code-defined default value. A disabled flag means the flag management system is explicitly stepping aside — it has nothing to say about what the value should be. This is semantically different from returning a configured `defaultVariant`, which would still delegate the decision to the flag management system and contradict the meaning of `DISABLED`.

The response has the following properties:

- **value**: Omitted from the response. The SDK/provider uses the application's code-defined default value. In Go this surfaces as the type's zero value; on the wire (protobuf/JSON) the field is left unset; SDKs in languages with optional types should expose it as absent (`None`/`null`/`undefined`).
- **variant**: Omitted from the response, using the same omission semantics as `value` above.
- **reason**: `DISABLED`
- **error**: `nil` (no error)
- **metadata**: The merged flag set + flag metadata (consistent with current behavior for successful evaluations)

This aligns with the OpenFeature specification's definition of `DISABLED` as a resolution reason and leverages the existing but unused `DisabledReason` constant already defined in flagd. The field-omission pattern reuses the mechanism established in the [code-default ADR](./support-code-default.md).

### Interaction with other resolution reasons

The [OpenFeature specification](https://openfeature.dev/specification/types/#resolution-reason) defines several resolution reasons. Here is how `DISABLED` interacts with each:

| Reason            | Current flagd usage                                                          | Interaction with DISABLED                                                                                                                                                                                                                                                                   |
| ----------------- | ---------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `STATIC`          | Returned when a flag has no targeting rules and resolves to `defaultVariant` | When disabled, `DISABLED` takes precedence. No variant is returned — the application uses its code default. `STATIC` tells caching providers the value is safe to cache indefinitely, while `DISABLED` signals the flag management system is stepping aside entirely.                        |
| `DEFAULT`         | Returned when targeting rules exist but evaluate to the `defaultVariant`     | When disabled, `DISABLED` takes precedence. Targeting rules are not evaluated at all for disabled flags, so `DEFAULT` (which implies targeting was attempted) is not appropriate.                                                                                                            |
| `SPLIT`           | Defined in flagd but used for pseudorandom assignment (fractional targeting) | When disabled, `DISABLED` takes precedence. Fractional targeting rules are not evaluated for disabled flags.                                                                                                                                                                                 |
| `TARGETING_MATCH` | Returned when targeting rules match and select a specific variant            | Not applicable. Targeting rules are never evaluated for disabled flags.                                                                                                                                                                                                                      |
| `ERROR`           | Currently returned for disabled flags (this is what we are changing)         | `DISABLED` replaces `ERROR` for this case. `ERROR` remains the reason for genuine errors (parse errors, type mismatches, etc.).                                                                                                                                                              |

**Key principle**: `DISABLED` is a terminal reason. When a flag is disabled, no targeting evaluation occurs and the flag management system defers entirely to the application's code default. Reasons that describe targeting outcomes (`STATIC`, `DEFAULT`, `SPLIT`, `TARGETING_MATCH`) never apply. The evaluation short-circuits to `reason=DISABLED` with value and variant omitted.

### API changes

**Evaluator core — `evaluateVariant`** (`core/pkg/evaluator/json.go`):

The `evaluateVariant` function changes from returning an error to returning a successful result with an empty variant, which signals code-default deferral:

```go
// Before
if flag.State == Disabled {
    return "", flag.Variants, model.ErrorReason, metadata, errors.New(model.FlagDisabledErrorCode)
}

// After — the empty string return is Go's zero value and is interpreted
// downstream as "variant omitted"; it is not a literal empty-string variant.
if flag.State == Disabled {
    return "", flag.Variants, model.DisabledReason, metadata, nil
}
```

**Evaluator core — `resolve[T]`** (`core/pkg/evaluator/json.go`):

The generic `resolve[T]` function currently only short-circuits for `FallbackReason` when the variant is empty. For any other reason it attempts a variant map lookup, which would produce a `TYPE_MISMATCH` error when the variant is empty. `resolve[T]` must treat `DisabledReason` the same way:

```go
// Before
if reason == model.FallbackReason {
    var zero T
    return zero, variant, model.FallbackReason, metadata, nil
}

// After
if reason == model.FallbackReason || reason == model.DisabledReason {
    var zero T
    return zero, variant, reason, metadata, nil
}
```

**Bulk evaluation — `ResolveAllValues`** (`core/pkg/evaluator/json.go`):

Disabled flags are no longer skipped. They are evaluated and included in the response. Since disabled flags always return an empty variant, the type switch on the variant value hits the `default` case. For gRPC v1 requests (where `ProtoVersionKey` is set), the current `default` branch skips unknown types via `continue`. This must be adjusted so disabled flags are always included regardless of proto version:

```go
// Before
if flag.State == Disabled {
    continue
}

// After: remove the skip — disabled flags flow through normal evaluation
// and will be returned with reason=DISABLED
```

```go
// Before (default branch of type switch)
default:
    if ctx.Value(ProtoVersionKey) == nil {
        value, variant, reason, metadata, err = resolve[interface{}](...)
    } else {
        continue
    }

// After: disabled flags must not be skipped even for old proto versions
default:
    if ctx.Value(ProtoVersionKey) == nil {
        value, variant, reason, metadata, err = resolve[interface{}](...)
    } else if flag.State == Disabled {
        value, variant, reason, metadata, err = resolve[interface{}](...)
    } else {
        continue
    }
```

**OFREP success mapping** (`core/pkg/service/ofrep/models.go`):

The `SuccessResponseFrom` function currently rewrites `FallbackReason` to `DefaultReason` and omits value/variant fields to signal code-default deferral. Disabled flags use the same field-omission pattern, but the reason must remain `DISABLED` (not be rewritten to `DEFAULT`):

```go
// Before
if result.Reason == model.FallbackReason {
    return EvaluationSuccess{
        Value:    nil,
        Key:      result.FlagKey,
        Reason:   model.DefaultReason,
        Variant:  "",
        Metadata: result.Metadata,
    }
}

// After: handle both fallback and disabled code-default deferral
if result.Reason == model.FallbackReason || result.Reason == model.DisabledReason {
    return EvaluationSuccess{
        Value:    nil,
        Key:      result.FlagKey,
        Reason:   lo.Ternary(result.Reason == model.FallbackReason, model.DefaultReason, model.DisabledReason),
        Variant:  "",
        Metadata: result.Metadata,
    }
}
```

**gRPC response** (single flag evaluation — value and variant omitted, SDK uses code default):

```json
{
  "reason": "DISABLED",
  "metadata": {
    "flagSetId": "my-app",
    "scope": "production"
  }
}
```

**OFREP response** (single flag evaluation, HTTP 200 — previously HTTP 404):

```json
{
  "key": "my-feature",
  "reason": "DISABLED",
  "metadata": {
    "flagSetId": "my-app"
  }
}
```

**OFREP bulk response** (disabled flag now included, value/variant omitted):

```json
{
  "flags": [
    {
      "key": "my-feature",
      "reason": "DISABLED",
      "metadata": {}
    },
    {
      "key": "active-feature",
      "value": true,
      "variant": "on",
      "reason": "STATIC",
      "metadata": {}
    }
  ]
}
```

### In-process providers

The sync.proto payload already includes disabled flags in the `flag_configuration` JSON with `"state": "DISABLED"`. No wire-format changes are needed. Each language SDK's in-process evaluator must be updated to return `reason=DISABLED` with code-default deferral (omitted value/variant) instead of raising an error when encountering a disabled flag.
This is functionally the same change as the flagd core evaluator, replicated in each SDK. This is a coordinated rollout across all SDK in-process providers and should be tracked alongside the flagd core changes.

### Consequences

- Good, because it aligns flagd with the OpenFeature specification's definition of `DISABLED` as a resolution reason
- Good, because it eliminates the OFREP bug where `FLAG_DISABLED` is silently masked as `FLAG_NOT_FOUND`
- Good, because operators get clean observability: disabled flags appear as successful evaluations with a distinct reason, not polluting error metrics
- Good, because it enables flag-set-based workflows where disabling a flag in one environment is a normal operational state
- Good, because the existing `DisabledReason` constant is finally used as designed
- Good, because it provides visibility into disabled flags in bulk evaluation responses, rather than silently omitting them
- Good, because applications can distinguish between "flag doesn't exist" (a real problem) and "flag is disabled" (an intentional state)
- Good, because always deferring to code defaults gives a clear semantic: disabled means the flag management system has nothing to say
- Bad, because it is a breaking change for clients that rely on `FLAG_DISABLED` error responses for control flow, alerting, or metrics
- Bad, because OFREP single-flag evaluation changes from HTTP 404 to HTTP 200, which is a breaking change for HTTP clients that branch on status codes
- Bad, because including disabled flags in `ResolveAll` responses increases payload size for flag sets with many disabled flags
- Bad, because it requires coordinated updates across flagd core, all gRPC/OFREP surfaces, in-process providers in each language SDK, and the testbed
- Neutral, because the `FlagDisabledErrorCode` constant and related error-handling code can be removed (code simplification)

### Versioning and migration

- This is a behavior-breaking change. Because flagd is pre-1.0, it is shipped as a minor-version bump and called out as breaking in the release notes — there is no dual-mode or deprecation period.
- Operators relying on `FLAG_DISABLED` error signals (in dashboards, alerts, log filters, or HTTP 404 branches) must migrate to keying on successful evaluations with `reason=DISABLED`. Migration guidance is communicated through the release notes rather than a runtime compatibility flag.
- The `FlagDisabledErrorCode` constant and its handling in `errFormat`, `errFormatV2`, and `EvaluationErrorResponseFrom` are removed outright in the same release. Retaining them serves no purpose once the evaluator no longer produces the error path, and there is no straightforward way to keep the old behavior reachable without re-introducing the spec-violating code path.

### Implementation steps

The work breaks down into three groups that must land together for a coherent release, but can be developed in parallel.

**flagd core (single release, behavior-breaking minor bump):**

1. Update `evaluateVariant` in `core/pkg/evaluator/json.go` to return `reason=DISABLED` with an omitted variant (code-default deferral) instead of an error.
2. Update `resolve[T]` in `core/pkg/evaluator/json.go` to handle `DisabledReason` with omitted variants (avoids `TYPE_MISMATCH`).
3. Remove the disabled-flag skip in `ResolveAllValues` and update the `default` branch of the type switch to include disabled flags for gRPC v1 requests.
4. Update `SuccessResponseFrom` in `core/pkg/service/ofrep/models.go` to preserve `reason=DISABLED` (with field omission) for disabled flags deferring to code defaults.
5. Update the gRPC v1 service layer (`flag_evaluator_v1.go`) to handle nil values in `ResolveAll` responses for disabled flags deferring to code defaults.
6. Remove `FlagDisabledErrorCode` handling from `errFormat`, `errFormatV2`, and `EvaluationErrorResponseFrom`.

**Ecosystem (rolled out alongside or shortly after the flagd core release):**

7. Update each language SDK's in-process provider evaluator to return `reason=DISABLED` with code-default deferral instead of raising an error when encountering a disabled flag.
8. Update OpenFeature providers (RPC and in-process) to recognize `DISABLED` as a non-error, non-cacheable reason.

**Validation and documentation:**

9. Update flagd-testbed with test cases for disabled flag evaluation across all surfaces (gRPC v1, gRPC v2, OFREP single, OFREP bulk, in-process).
10. Update provider documentation and call out the behavior change prominently in the flagd release notes.

### Open questions

- Should bulk evaluation include an option to exclude disabled flags for clients that prefer the current behavior (smaller payloads)?

## More Information

- [OpenFeature Specification - Resolution Reasons](https://openfeature.dev/specification/types/#resolution-reason)
- [OpenFeature Specification - Error Codes](https://openfeature.dev/specification/types/#error-code) — notably, `FLAG_DISABLED` is not in the spec's error code list
- [flagd Flag Definitions Reference](https://flagd.dev/reference/flag-definitions/)
- [ADR: Support Explicit Code Default Values](./support-code-default.md)
- [flagd Testbed](https://github.com/open-feature/flagd-testbed)
