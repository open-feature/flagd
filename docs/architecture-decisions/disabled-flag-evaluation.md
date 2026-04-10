---
# Valid statuses: draft | proposed | rejected | accepted | superseded
status: draft
author: Parth Suthar (@suthar26)
created: 2026-04-01
---
# Treat Disabled Flag Evaluation as Successful with Reason DISABLED

This ADR proposes changing flagd's handling of disabled flags from returning an error (`reason=ERROR`, `errorCode=FLAG_DISABLED`) to returning a successful evaluation with `reason=DISABLED` and the flag's `defaultVariant` value. A flag that does not exist in a flag set should remain a `FLAG_NOT_FOUND` error.
This aligns flagd with the [OpenFeature specification's resolution reasons](https://openfeature.dev/specification/types/#resolution-reason), which defines `DISABLED` as a valid resolution reason for successful evaluations.

## Background

flagd currently treats the evaluation of a disabled flag as an error. When a flag exists in the store but has `state: DISABLED`, the evaluator returns `reason=ERROR` with `errorCode=FLAG_DISABLED`. This error propagates through every surface — gRPC, OFREP, and in-process providers — resulting in the caller receiving an error response rather than a resolved value.

This is problematic for several reasons:

1. **Spec misalignment**: The [OpenFeature specification](https://openfeature.dev/specification/types/#resolution-reason) explicitly defines `DISABLED` as a resolution reason with the description: *"The resolved value was the result of the flag being disabled in the management system."* This implies a successful evaluation that communicates the flag's disabled state, not an error.
2. **OFREP masks the disabled state**: The OFREP response handler in `core/pkg/service/ofrep/models.go` rewrites `FLAG_DISABLED` to `FLAG_NOT_FOUND` in the structured error response (while only preserving the "is disabled" distinction in the free-text `errorDetails` string). This means OFREP clients cannot programmatically distinguish between a flag that doesn't exist and one that was intentionally disabled.
3. **Conflation of "missing" and "disabled"**: gRPC v1 maps both `FLAG_NOT_FOUND` and `FLAG_DISABLED` to `connect.CodeNotFound`. These are semantically different situations: a missing flag is a configuration or deployment error, while a disabled flag is an intentional operational decision (incident remediation, environment-specific rollout, not-yet-ready feature).
4. **Loss of observability**: When disabled flags are treated as errors, they pollute error metrics and alerting. Operators who disable a flag for legitimate reasons (ongoing incident remediation, feature not ready for an environment) see false error signals. Conversely, if they suppress these errors, they lose visibility into flag state entirely. A successful evaluation with `reason=DISABLED` would give operators a clean signal without noise.
5. **Flag set use cases**: In multi-flag-set deployments, a flag may exist in a shared definition but be disabled in certain flag sets (e.g., disabled for `staging` but enabled for `production`). Treating this as an error forces the application into error-handling paths when the flag is simply not active — a normal operational state, not an exceptional one.

Related context:

- [OpenFeature Specification - Resolution Reasons](https://openfeature.dev/specification/types/#resolution-reason)
- [flagd Flag Definitions Reference](https://flagd.dev/reference/flag-definitions/)
- [ADR: Support Explicit Code Default Values](./support-code-default.md) — establishes the pattern of returning `defaultVariant` with appropriate reason codes

## Requirements

- Evaluating a disabled flag must return a successful response with `reason=DISABLED` across all surfaces (gRPC v1, gRPC v2, OFREP, in-process)
- The resolved value for a disabled flag must be the flag's configured `defaultVariant`
- If the flag's `defaultVariant` is `null` (per the code-default ADR), the disabled response must defer to code defaults using the same field-omission pattern, but with `reason=DISABLED` instead of `reason=DEFAULT`
- Evaluating a flag key that does not exist in the store or flag set must remain a `FLAG_NOT_FOUND` error
- `reason=DISABLED` must be distinct from all other reasons (`STATIC`, `DEFAULT`, `SPLIT`, `TARGETING_MATCH`, `ERROR`) and must not trigger error-handling paths in providers or SDKs
- Bulk evaluation (`ResolveAll`) must include disabled flags in the response with `reason=DISABLED`, rather than silently omitting them
- Telemetry and metrics must record disabled flag evaluations as successful (non-error) with the `DISABLED` reason
- Existing flag configurations must continue to work without modification (backward compatible at the configuration level)

## Considered Options

- **Option 1: Successful evaluation with `reason=DISABLED` returning `defaultVariant`** — Disabled flags evaluate successfully, returning the `defaultVariant` value and `reason=DISABLED`
- **Option 2: Return a successful evaluation with** `reason=DEFAULT` — Treat disabled flags as if they had no targeting, collapsing the disabled state into the default reason
- **Option 3: Status quo** — Keep the current error behavior and document it as intentional divergence from the OpenFeature spec

## Proposal

We propose **Option 1: Successful evaluation with `reason=DISABLED` returning `defaultVariant`**.

When a flag exists in the store but has `state: DISABLED`, the evaluator should return a successful evaluation with the following properties:

- **value**: The flag's `defaultVariant` value (from the flag configuration). If `defaultVariant` is `null` or absent, the value field is omitted to signal code-default deferral (see [Interaction with code defaults](#interaction-with-code-defaults-defaultvariant-null)).
- **variant**: The flag's `defaultVariant` key. If `defaultVariant` is `null` or absent, the variant field is omitted.
- **reason**: `DISABLED`
- **error**: `nil` (no error)
- **metadata**: The merged flag set + flag metadata (consistent with current behavior for successful evaluations)

This aligns with the OpenFeature specification's definition of `DISABLED` as a resolution reason and leverages the existing but unused `DisabledReason` constant already defined in flagd.

### Interaction with other resolution reasons

The [OpenFeature specification](https://openfeature.dev/specification/types/#resolution-reason) defines several resolution reasons. Here is how `DISABLED` interacts with each:

| Reason            | Current flagd usage                                                          | Interaction with DISABLED                                                                                                                                                                                                                                                                                                          |
| ----------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `STATIC`          | Returned when a flag has no targeting rules and resolves to `defaultVariant` | When disabled, `DISABLED` takes precedence. The flag's `defaultVariant` is still returned as the value, but the reason is `DISABLED`, not `STATIC`. This distinction matters: `STATIC` tells caching providers the value is safe to cache indefinitely, while `DISABLED` signals the value may change when the flag is re-enabled. |
| `DEFAULT`         | Returned when targeting rules exist but evaluate to the `defaultVariant`     | When disabled, `DISABLED` takes precedence. Targeting rules are not evaluated at all for disabled flags, so `DEFAULT` (which implies targeting was attempted) is not appropriate.                                                                                                                                                  |
| `SPLIT`           | Defined in flagd but used for pseudorandom assignment (fractional targeting) | When disabled, `DISABLED` takes precedence. Fractional targeting rules are not evaluated for disabled flags.                                                                                                                                                                                                                       |
| `TARGETING_MATCH` | Returned when targeting rules match and select a specific variant            | Not applicable. Targeting rules are never evaluated for disabled flags.                                                                                                                                                                                                                                                            |
| `ERROR`           | Currently returned for disabled flags (this is what we are changing)         | `DISABLED` replaces `ERROR` for this case. `ERROR` remains the reason for genuine errors (parse errors, type mismatches, etc.).                                                                                                                                                                                                    |

**Key principle**: `DISABLED` is a terminal reason. When a flag is disabled, no targeting evaluation occurs, so reasons that describe targeting outcomes (`STATIC`, `DEFAULT`, `SPLIT`, `TARGETING_MATCH`) never apply. The evaluation short-circuits to `reason=DISABLED` with the `defaultVariant` value.

### Interaction with code defaults (`defaultVariant: null`)

Per the [code-default ADR](./support-code-default.md), when `defaultVariant` is `null`, the server omits the value and variant fields to signal code-default deferral. This pattern applies to disabled flags as well:

- `defaultVariant` is a string → return the variant value with `reason=DISABLED`
- `defaultVariant` is `null` → omit value/variant fields, return `reason=DISABLED`

The only difference from a normal code-default response is the reason field: `DISABLED` instead of `DEFAULT`.

### API changes

**Evaluator core — `evaluateVariant`** (`core/pkg/evaluator/json.go`):

The `evaluateVariant` function changes from returning an error to returning a successful result. Since `flag.DefaultVariant` is already `""` when no default variant is configured, no additional branch is needed:

```go
// Before
if flag.State == Disabled {
    return "", flag.Variants, model.ErrorReason, metadata, errors.New(model.FlagDisabledErrorCode)
}

// After
if flag.State == Disabled {
    return flag.DefaultVariant, flag.Variants, model.DisabledReason, metadata, nil
}
```

**Evaluator core — `resolve[T]`** (`core/pkg/evaluator/json.go`):

The generic `resolve[T]` function currently only short-circuits for `FallbackReason` when the variant is empty. For any other reason it attempts a variant map lookup, which would produce a `TYPE_MISMATCH` error for disabled flags that defer to code defaults (`defaultVariant: null`). `resolve[T]` must treat `DisabledReason` the same way when the variant is empty:

```go
// Before
if reason == model.FallbackReason {
    var zero T
    return zero, variant, model.FallbackReason, metadata, nil
}

// After
if reason == model.FallbackReason || reason == model.DisabledReason {
    if variant == "" {
        var zero T
        return zero, variant, reason, metadata, nil
    }
}
```

**Bulk evaluation — `ResolveAllValues`** (`core/pkg/evaluator/json.go`):

Disabled flags are no longer skipped. They are evaluated and included in the response. Additionally, when `defaultVariant` is `null`, the type switch on `flag.Variants[flag.DefaultVariant]` hits the `default` case. For gRPC v1 requests (where `ProtoVersionKey` is set), the current `default` branch skips unknown types via `continue`. This must be adjusted so disabled flags are always included regardless of proto version:

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

The `SuccessResponseFrom` function currently rewrites `FallbackReason` to `DefaultReason` and omits value/variant fields to signal code-default deferral. When a disabled flag defers to code defaults, the same field-omission pattern must apply, but the reason must remain `DISABLED` (not be rewritten to `DEFAULT`):

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

// After: handle disabled flags deferring to code defaults
if result.Reason == model.FallbackReason {
    return EvaluationSuccess{
        Value:    nil,
        Key:      result.FlagKey,
        Reason:   model.DefaultReason,
        Variant:  "",
        Metadata: result.Metadata,
    }
}
if result.Reason == model.DisabledReason && result.Variant == "" {
    return EvaluationSuccess{
        Value:    nil,
        Key:      result.FlagKey,
        Reason:   model.DisabledReason,
        Variant:  "",
        Metadata: result.Metadata,
    }
}
```

**gRPC response** (single flag evaluation):

```json
{
  "value": false,
  "variant": "off",
  "reason": "DISABLED",
  "metadata": {
    "flagSetId": "my-app",
    "scope": "production"
  }
}
```

**OFREP response** (single flag evaluation):

```json
{
  "key": "my-feature",
  "value": false,
  "variant": "off",
  "reason": "DISABLED",
  "metadata": {
    "flagSetId": "my-app"
  }
}
```

**OFREP bulk response** (disabled flag now included):

```json
{
  "flags": [
    {
      "key": "my-feature",
      "value": false,
      "variant": "off",
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

### Consequences

- Good, because it aligns flagd with the OpenFeature specification's definition of `DISABLED` as a resolution reason
- Good, because it eliminates the OFREP bug where `FLAG_DISABLED` is silently masked as `FLAG_NOT_FOUND`
- Good, because operators get clean observability: disabled flags appear as successful evaluations with a distinct reason, not polluting error metrics
- Good, because it enables flag-set-based workflows where disabling a flag in one environment is a normal operational state
- Good, because the existing `DisabledReason` constant is finally used as designed
- Good, because it provides visibility into disabled flags in bulk evaluation responses, rather than silently omitting them
- Good, because applications can distinguish between "flag doesn't exist" (a real problem) and "flag is disabled" (an intentional state)
- Bad, because it is a breaking change for clients that rely on `FLAG_DISABLED` error responses for control flow, alerting, or metrics
- Bad, because including disabled flags in `ResolveAll` responses increases payload size for flag sets with many disabled flags
- Bad, because it requires coordinated updates across flagd core, all gRPC/OFREP surfaces, providers, and the testbed
- Neutral, because the `FlagDisabledErrorCode` constant and related error-handling code can be removed (code simplification)

### Timeline

1. Update `evaluateVariant` in `core/pkg/evaluator/json.go` to return `reason=DISABLED` with `defaultVariant` instead of an error
2. Update `resolve[T]` in `core/pkg/evaluator/json.go` to handle `DisabledReason` with empty variants (avoids `TYPE_MISMATCH` when disabled flags defer to code defaults)
3. Remove the disabled-flag skip in `ResolveAllValues` and update the `default` branch of the type switch to include disabled flags for gRPC v1 requests
4. Update `SuccessResponseFrom` in `core/pkg/service/ofrep/models.go` to preserve `reason=DISABLED` (with field omission) when a disabled flag defers to code defaults
5. Remove `FlagDisabledErrorCode` handling from `errFormat`, `errFormatV2`, and `EvaluationErrorResponseFrom`
6. Update gRPC v1 service layer (`flag_evaluator_v1.go`) to handle nil values in `ResolveAll` responses for disabled flags deferring to code defaults
7. Update flagd-testbed with test cases for disabled flag evaluation across all surfaces
8. Update OpenFeature providers to recognize `DISABLED` as a non-error, non-cacheable reason
9. Update provider documentation and migration guides

### Open questions

- Should bulk evaluation include an option to exclude disabled flags for clients that prefer the current behavior (smaller payloads)?
- How should existing dashboards and alerts that key on `FLAG_DISABLED` errors be migrated? Should we provide a deprecation period where both behaviors are available?
- Does this change require a new flagd major version, or can it be introduced in a minor version with appropriate documentation given the spec alignment argument?
- Should the `FlagDisabledErrorCode` constant be retained (but unused) for a deprecation period, or removed immediately?
- How should in-process providers handle the transition? They evaluate locally and would need to be updated to return `DISABLED` reason instead of throwing an error.

## More Information

- [OpenFeature Specification - Resolution Reasons](https://openfeature.dev/specification/types/#resolution-reason)
- [OpenFeature Specification - Error Codes](https://openfeature.dev/specification/types/#error-code) — notably, `FLAG_DISABLED` is not in the spec's error code list
- [flagd Flag Definitions Reference](https://flagd.dev/reference/flag-definitions/)
- [ADR: Support Explicit Code Default Values](./support-code-default.md)
- [flagd Testbed](https://github.com/open-feature/flagd-testbed)
