---
status: accepted
author: @beeme1mr
created: 2025-06-06
updated: 2025-06-20
---

# Support Explicit Code Default Values in flagd Configuration

This ADR proposes adding support for explicitly configuring flagd to use code-defined default values by allowing `null` as a valid default variant. This change addresses the current limitation where users cannot differentiate between "use the code's default" and "use this configured default" without resorting to workarounds like misconfigured rulesets.

## Background

Currently, flagd requires a default variant to be specified in flag configurations. This creates a fundamental mismatch with the OpenFeature specification and common feature flag usage patterns where code-defined defaults serve as the ultimate fallback.

The current behavior leads to confusion and operational challenges:

1. **Two Sources of Truth**: Applications have default values defined in code (as per OpenFeature best practices), while flagd configurations require their own default variants. This dual-default pattern violates the principle of single source of truth.

2. **State Transition Issues**: When transitioning a flag from DISABLED to ENABLED state, the behavior changes unexpectedly:
   - DISABLED state: Flag evaluation falls through to code defaults
   - ENABLED state: Flag evaluation uses the configured default variant

3. **Workarounds**: Users resort to misconfiguring rulesets (e.g., returning invalid variants) to force fallback to code defaults, which generates confusing error states and complicates debugging.

4. **OpenFeature Alignment**: The OpenFeature specification emphasizes that code defaults should be the ultimate fallback, but flagd's current design doesn't provide a clean way to express this intent.

Related discussions and context can be found in the [OpenFeature specification](https://openfeature.dev/specification/types) and [flagd flag definitions reference](https://flagd.dev/reference/flag-definitions/).

## Requirements

- **Explicit Code Default Support**: Users must be able to explicitly configure a flag to use the code-defined default value as its resolution
- **Backward Compatibility**: Existing flag configurations must continue to work without modification
- **Clear Semantics**: The configuration must clearly indicate when code defaults are being used versus configured defaults
- **Appropriate Reason Codes**: Resolution details must include appropriate reason codes when code defaults are used (e.g., `DEFAULT` or a new specific reason)
- **Schema Validation**: JSON schema must support and validate the new configuration options
- **Provider Compatibility**: All OpenFeature providers must handle the new behavior correctly
- **Testbed Coverage**: flagd testbed must include test cases for the new functionality

## Considered Options

- **Option 1: Allow `null` as Default Variant** - Modify the schema to accept `null` as a valid value for defaultVariant, signaling "use code default"
- **Option 2: Make Default Variant Optional** - Remove the requirement for defaultVariant entirely, with absence meaning "use code default"
- **Option 3: Special Variant Value** - Define a reserved variant name (e.g., `"__CODE_DEFAULT__"`) that signals code default usage
- **Option 4: New Configuration Property** - Add a new property like `useCodeDefault: true` alongside or instead of defaultVariant
- **Option 5: Status Quo with Documentation** - Keep current behavior but improve documentation about workarounds

## Proposal

We propose implementing **Option 1: Allow `null` as Default Variant**, potentially combined with **Option 2: Make Default Variant Optional** for maximum flexibility.

The implementation leverages field presence in evaluation responses across all protocols (in-process, RPC, and OFREP). When a flag configuration has `defaultVariant: null`, the evaluation response omits the value field entirely, which serves as a programmatic signal to the client to use its code-defined default value.

This approach offers several key advantages:

1. **No Protocol Changes**: RPC and OFREP protocols remain unchanged
2. **Clear Semantics**: Omitted value field = "use your code default"
3. **Backward Compatible**: Existing clients and servers continue to work
4. **Universal Pattern**: Works consistently across all evaluation modes

The absence of a value field provides an unambiguous signal that distinguishes between "the server evaluated to null/false/empty" (value field present) and "the server delegates to your code default" (value field absent).

### Implementation Details

1. **Schema Changes**:

   ```json
   {
     "defaultVariant": {
       "oneOf": [
         { "type": "string" },
         { "type": "null" }
       ],
       "description": "Default variant to use. Set to null to use code-defined default."
     }
   }
   ```

2. **Evaluation Behavior**:
   - When flag has `defaultVariant: null` and targeting returns no match
   - Server responds with reason set to reason "ERROR" and error code "FLAG_NOT_FOUND"
   - Client detects this reason value field and uses its code-defined default
   - This same pattern works across all evaluation modes

3. **Provider Implementation**:
   - No changes to existing providers

### Design Rationale

**Using "ERROR" reason**: We intentionally reuse the existing "ERROR" reason code rather than introducing a new one (like "CODE_DEFAULT"). This retains the current behavior of an disabled flag and allows for progressive enablement of a flag without unexpected variations in flag evaluation behavior.

Advantages of this approach:

- The "ERROR" reason is already used for cases where the flag is not found or misconfigured, so it aligns with the intent of using code defaults.
- This approach avoids introducing new reason codes that would require additional handling in providers and clients.

### API changes

**Flag Configuration**:

```yaml
flags:
  my-feature:
    state: ENABLED
    defaultVariant: null  # Explicitly use code default
    variants:
      on: true
      off: false
    targeting:
      if:
        - "===":
            - var: user-type
            - "beta"
        - on
```

**OFREP Response** when code default is indicated:

#### Single flag evaluation response

A single flag evaluation returns a `404` status code.

```json
{
  "key": "my-feature",
  "errorCode": "FLAG_NOT_FOUND",
  // Optional error details
  "errorDetails": "Targeting not matched, using code default",
  "metadata": {}
}
```

#### Bulk flag evaluation response

```json
{
  "flags": [
    // Flag is omitted from bulk response
  ]
}
```

**flagd RPC Response** (ResolveBooleanResponse):

```protobuf
{
  "reason": "ERROR",
  "errorCode": "FLAG_NOT_FOUND",
  "metadata": {}
}
```

### Consequences

- Good, because it eliminates the confusion between code and configuration defaults
- Good, because it provides explicit control over default behavior without workarounds
- Good, because it aligns flagd more closely with OpenFeature specification principles
- Good, because it supports gradual flag rollout patterns more naturally
- Good, because it provides the ability to delegate to whatever is defined in code
- Good, because it requires no changes to existing RPC or protocol signatures
- Good, because it uses established patterns (field presence) for clear semantics
- Good, because it maintains full backward compatibility
- Bad, because it requires updates across multiple components (flagd, providers, testbed)
- Bad, because it introduces a new concept that users need to understand
- Neutral, because existing configurations continue to work unchange

### Implementation Plan

1. Update flagd-schemas with new JSON schema supporting null default variants
2. Update flagd-testbed with comprehensive test cases for all evaluation modes
3. Implement core logic in flagd to handle null defaults and omit value/variant fields
4. Update OpenFeature providers with the latest schema and test harness to ensure they handle the new behavior correctly
5. Documentation updates, migration guides, and playground examples to demonstrate the new configuration options

### Testing Considerations

To ensure correct implementation across all components:

1. **Provider Tests**: Each component (flagd, providers) must have unit tests verifying the handling of `null` as a default variant
2. **Integration Tests**: End-to-end tests across different language combinations (e.g., Go flagd with Java provider)
3. **OFREP Tests**: Verify JSON responses correctly omits flags with a `null` default variant
4. **Backward Compatibility Tests**: Ensure old providers handle new responses gracefully
5. **Consistency Tests**: Verify identical behavior across in-process, RPC, and OFREP modes

### Open questions

- How should providers handle responses with missing value fields in strongly-typed languages?
    - We'll handle the same way as with optional fields, using language-specific patterns (e.g., pointers in Go, `hasValue()` in Java).
- Should we support both `null` and absent `defaultVariant` fields, or choose one approach?
    - Yes, we'll support both `null` and absent fields to maximize flexibility. An absent `defaultVariant` will be the equivalent of `null`.
- What migration path should we recommend for users currently using workarounds?
    - Update the flag configurations to use `defaultVariant: null` and remove any misconfigured rulesets that force code defaults.
- Should this feature be gated behind a configuration flag during initial rollout?
    - We'll avoid public facing documentation until the feature is fully implemented and tested.
- How do we ensure consistent behavior across all provider implementations?
    - Gherkin tests will be added to the flagd testbed to ensure all providers handle the new behavior consistently.
- Should providers validate that the reason is "DEFAULT" when value is omitted, or accept any omitted value as delegation?
    - Providers should accept any omitted value as delegation.
- How do we handle edge cases where network protocols might strip empty fields?
    - It would behaving as expected, as the absence of fields is the intended signal.
- When the client uses its code default after receiving a delegation response, what variant should be reported in telemetry/analytics?
    - The variant will be omitted, indicating that the code default was used.
- Should we add explicit proto comments documenting the field omission behavior?
    - Leave this to the implementers, but it would be beneficial to add comments in the proto files to clarify this behavior for future maintainers.

## More Information

- [OpenFeature Specification - Flag Evaluation](https://openfeature.dev/specification/types#flag-evaluation)
- [flagd Flag Definitions Reference](https://flagd.dev/reference/flag-definitions/)
- [flagd JSON Schema Repository](https://github.com/open-feature/flagd-schemas)
- [flagd Testbed](https://github.com/open-feature/flagd-testbed)
