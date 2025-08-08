---
status: accepted
author: Michael Beemer
created: 2025-06-06
updated: 2025-08-08
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

The implementation leverages field presence in evaluation responses across all protocols.
When a flag configuration has `defaultVariant: null`, the evaluation response omits the value field entirely and uses the "DEFAULT" reason code, which serves as a programmatic signal to the client to use its code-defined default value.

This approach offers several key advantages:

1. **Semantically Correct**: Uses "DEFAULT" reason code which accurately represents the evaluation outcome
2. **Success Responses**: Treats code default usage as successful evaluation, not an error
3. **Clear Semantics**: Omitted value field = "use your code default"
4. **Backward Compatible**: Existing clients and servers continue to work
5. **Universal Pattern**: Works consistently across all evaluation modes
6. **Accurate Telemetry**: Metrics correctly reflect successful evaluations rather than false errors

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
   - Server responds with reason "DEFAULT" and omits value and variant fields
   - Client detects the omitted fields and uses its code-defined default
   - This pattern works consistently across all evaluation modes

3. **Protobuf Schema Changes**:
   - Update response message definitions to use `optional` fields for `value` and `variant`
   - This enables proper field presence detection for code default signaling

   Example protobuf changes:

   ```protobuf
   message ResolveBooleanResponse {
     // The response value, will be unset when deferring to code defaults
     optional bool value = 1;
   
     // The reason for the given return value
     string reason = 2;
   
     // The variant name, will be unset when deferring to code defaults
     optional string variant = 3;
   
     // Metadata for this evaluation
     google.protobuf.Struct metadata = 4;
   }
   ```

4. **Provider Implementation**:
   - Providers must be updated to check field presence rather than just reading field values

### Design Rationale

**Using "DEFAULT" reason with omitted value fields**: We use the "DEFAULT" reason code to accurately represent that a default value is being used, combined with omitting the value and variant fields to signal code default deferral. This approach leverages recent OFREP improvements and requires updating protobuf definitions to use `optional` fields for proper field presence detection.

Advantages of this approach:

- **Accurate Semantics**: "DEFAULT" reason correctly represents the evaluation outcome
- **Proper Telemetry**: Evaluations are recorded as successful rather than errors  
- **Clear Field Presence**: Optional fields provide unambiguous signaling across all protocols
- **Standards Aligned**: Leverages accepted patterns for optional values
- **Backward Compatible**: Existing clients continue to work while new clients can detect code defaults

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

A single flag evaluation returns a `200` status code:

```json
{
  "key": "my-feature",
  "reason": "DEFAULT",
  "metadata": {}
  // Note: No value field - indicates code default usage
}
```

#### Bulk flag evaluation response

```json
{
  "flags": [
    {
      "key": "my-feature", 
      "reason": "DEFAULT",
      "metadata": {}
      // Note: No value field - indicates code default usage
    }
  ]
}
```

**flagd RPC Response** (ResolveBooleanResponse):

```protobuf
{
  "reason": "DEFAULT",
  "metadata": {}
  // Note: value and variant fields omitted to indicate code default
}
```

### Consequences

- Good, because it eliminates the confusion between code and configuration defaults
- Good, because it provides explicit control over default behavior without workarounds
- Good, because it aligns flagd more closely with OpenFeature specification principles
- Good, because it supports gradual flag rollout patterns more naturally
- Good, because it provides the ability to delegate to whatever is defined in code
- Good, because it uses the "DEFAULT" reason code which accurately represents the evaluation outcome
- Good, because it treats code default usage as successful evaluation with proper telemetry
- Good, because telemetry can distinguish between configured defaults (variant present) and code defaults (variant absent)
- Good, because it uses a simple field presence pattern that works across all protocols
- Good, because it maintains backward compatibility for existing flag configurations
- Bad, because it requires protobuf schema changes to use `optional` fields
- Bad, because it requires updates across multiple components (flagd, providers, testbed)
- Bad, because it introduces a new concept that users need to understand
- Bad, because it creates a breaking change for older clients evaluating flags configured with `defaultVariant: null` (they would receive zero values instead of using code defaults)
- Bad, because providers must be updated to handle field presence detection
- Neutral, because existing configurations continue to work unchanged

### Implementation Plan

1. Update flagd protobuf schemas to use `optional` fields for `value` and `variant` in response messages
2. Update flagd-schemas with new JSON schema supporting null default variants
3. Update flagd-testbed with comprehensive test cases for all evaluation modes
4. Implement core logic in flagd to handle null defaults by conditionally omitting fields in responses
5. Update OpenFeature providers to check field presence rather than just reading field values
6. Regenerate protobuf client libraries for all supported languages with new optional field support
7. Release updated clients before configuring any flags with `defaultVariant: null` to avoid zero-value issues with older clients
8. Update provider documentation with field presence detection patterns for each language
9. Add backward compatibility testing to ensure existing clients continue to work
10. Update CI/CD pipelines to validate protobuf schema changes and field presence behavior
11. Documentation updates, migration guides, and playground examples to demonstrate the new configuration options

### Testing Considerations

To ensure correct implementation across all components:

1. **Provider Tests**: Each component (flagd, providers) must have unit tests verifying the handling of `null` as a default variant
2. **Integration Tests**: End-to-end tests across different language combinations (e.g., Go flagd with Java provider)
3. **Schema Tests**: Verify protobuf schemas correctly define `optional` fields and generate appropriate client code
4. **Field Presence Tests**: Verify that providers can correctly detect field presence vs. absence across all languages
5. **OFREP Tests**: Verify JSON responses correctly omit value fields for code default scenarios
6. **RPC Tests**: Verify protobuf responses correctly omit optional fields for code default scenarios
7. **Backward Compatibility Tests**: Ensure old providers handle new responses gracefully
8. **Consistency Tests**: Verify consistent field presence behavior across all evaluation modes

### Open questions

- How should providers handle responses with missing value fields in strongly-typed languages?
    - We'll handle the same way as with optional fields, using language-specific patterns (e.g., pointers in Go, `hasValue()` in Java).
- Should we support both `null` and absent `defaultVariant` fields, or choose one approach?
    - Yes, we'll support both `null` and absent fields to maximize flexibility. An absent `defaultVariant` will be the equivalent of `null`.
- What migration path should we recommend for users currently using workarounds?
    - Update the flag configurations to use `defaultVariant: null` and remove any misconfigured rulesets that force code defaults.
- How should we handle the breaking change for older clients evaluating `defaultVariant: null` flags?
    - Older clients built without optional protobuf fields will receive zero values (false, 0, "") instead of using code defaults. This requires coordinated rollout: (1) Update and deploy all clients with new protobuf definitions, (2) Only then configure flags with `defaultVariant: null`. Alternatively, maintain separate flag configurations during transition period.
- Should this feature be gated behind a configuration flag during initial rollout?
    - We'll avoid public facing documentation until the feature is fully implemented and tested.
- How do we ensure consistent behavior across all provider implementations?
    - Gherkin tests will be added to the flagd testbed to ensure all providers handle the new behavior consistently.
- How do OFREP providers detect and handle responses with omitted value fields?
    - Providers should check for the absence of the `value` field in successful responses and treat it as a signal to use code defaults.
- Should we maintain backward compatibility for providers that don't yet support omitted value fields?
    - Yes, older providers will continue to work but may not benefit from the code default deferral feature until updated.
- When the client uses its code default after receiving a delegation response, what variant should be reported in telemetry/analytics?
    - No variant will be reported since the variant is unknown when using code defaults. The absence of a variant in telemetry indicates that a code default was used.
- Should we add explicit documentation about the field omission behavior?
    - Yes, clear documentation should explain how omitted value fields signal code default deferral for implementers.

## Revision History

| Date       | Author    | Change Summary                                                                                                                                                              |
| ---------- | --------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 2025-06-06 | @beeme1mr | Initial ADR creation with error-based approach                                                                                                                              |
| 2025-08-08 | @beeme1mr | **Major revision**: Changed from error-based (`FLAG_NOT_FOUND`) to success-based approach (`DEFAULT` reason) following OFREP improvements that enable optional value fields |

## More Information

- [OpenFeature Specification - Flag Evaluation](https://openfeature.dev/specification/types#flag-evaluation)
- [flagd Flag Definitions Reference](https://flagd.dev/reference/flag-definitions/)
- [flagd JSON Schema Repository](https://github.com/open-feature/flagd-schemas)
- [flagd Testbed](https://github.com/open-feature/flagd-testbed)
- [OFREP ADR: Optional value field for code default deferral](https://github.com/open-feature/protocol/blob/main/service/adrs/0006-optional-value-for-code-defaults.md)
