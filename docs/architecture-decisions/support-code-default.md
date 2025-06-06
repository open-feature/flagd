---
# Valid statuses: draft | proposed | rejected | accepted | superseded
status: draft
author: @beeme1mr
created: 2025-05-06
updated: 2025-05-06
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

2. **Protobuf Considerations**:

   **No proto changes required** - The existing proto3 definitions already support this behavior:

   ```protobuf
   message ResolveBooleanResponse {
     bool value = 1;
     string reason = 2;
     string variant = 3;
     google.protobuf.Struct metadata = 4;
   }
   ```

   In proto3, all fields are optional by default. To implement code default delegation:
   - Do NOT set the `value` field (not even to false/0/"")
   - Do NOT set the `variant` field
   - Set `reason` to "DEFAULT"
   - Optionally set `metadata`

   **Critical implementation notes**:
   - Ensure flagd omits fields rather than setting zero values
   - Different languages detect field presence differently:
     - Go: Check for zero values or use pointers
     - Java: Use `hasValue()` / `hasVariant()` methods
     - Python: Use `HasField()` or check field presence
     - JavaScript: Check for `undefined` (not null)
   - Consider adding proto comments documenting this behavior
   - Test wire format to confirm fields are actually omitted

   **Example - Correct vs Incorrect Implementation**:

   ```go
   // INCORRECT - sets zero value
   response := &ResolveBooleanResponse{
       Value:   false,  // This sends 'false' on the wire
       Reason:  "DEFAULT",
       Variant: "",     // This sends empty string on the wire
   }
   
   // CORRECT - omits fields
   response := &ResolveBooleanResponse{
       Reason: "DEFAULT",
       // Value and Variant fields are not set at all
   }
   ```

3. **Evaluation Behavior**:
   - When flag has `defaultVariant: null` and targeting returns no match
   - Server responds with both value and variant fields omitted, reason set to "DEFAULT"
   - Client detects the missing value field and uses its code-defined default
   - This same pattern works across all evaluation modes

   **In-Process Mode Special Considerations**:
   - The in-process evaluator must return the same "shape" of response
   - Cannot return null/nil as that's different from "no value"
   - May need a special response type or wrapper to indicate delegation
   - Example approach: Return evaluation result with a "useDefault" flag

4. **Remote Evaluation Protocol Responses**:

   **OFREP (OpenFeature Remote Evaluation Protocol)**:
   - Return HTTP 200 with response body that omits both value and variant fields
   - Reason field set to "DEFAULT" to indicate delegation
   - Clear distinction from error cases (which return 4xx/5xx)

   **flagd RPC**:
   - Omit both the type-specific value field and variant field in response messages
   - Use protobuf field presence to signal "no opinion from server"
   - No changes needed to RPC method signatures

5. **Provider Implementation**:
   - Check for presence/absence of value field in responses
   - When value is absent and reason is "DEFAULT", use code-defined default
   - When value is present (even if null/false/empty), use that value
   - Variant field will also be absent in delegation responses, resulting in undefined variant in resolution details
   - Responses with value fields work as before, maintaining backward compatibility

### Design Rationale

**Using "DEFAULT" reason**: We intentionally reuse the existing "DEFAULT" reason code rather than introducing a new one (like "CODE_DEFAULT"). The distinction between a configuration default and code default is clear from the response structure:

- Configuration default: Has both value and variant fields
- Code default: Omits both value and variant fields

**Field Omission Pattern**: Using field presence/absence is a well-established pattern in protocol design:

- Unambiguous: Cannot confuse "null value" with "no opinion"
- Language agnostic: Works across type systems
- Protocol friendly: Natural in both JSON and Protobuf
- Backward compatible: Existing responses always include values
- Spec compliant: OpenFeature allows undefined variants

The omission of both value and variant fields creates a clear, consistent signal that the server is fully delegating the decision to the client's code default.

This approach maintains compatibility with the established OpenFeature terminology while providing clear semantics through response structure.

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

```json
{
  "key": "my-feature",
  "reason": "DEFAULT",
  "metadata": {}
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
    }
  ]
}
```

Note: Both `value` and `variant` fields are intentionally omitted to signal "use code default"

**flagd RPC Response** (ResolveBooleanResponse):

```protobuf
{
  "reason": "DEFAULT",
  "metadata": {}
}
```

Both the type-specific value field and variant field are omitted

**Provider behavior**:

1. Check response for presence of value field
2. If value field is absent and reason is "DEFAULT", use code-defined default
3. If value field is present (even if null/false/empty), use that value
4. The variant field will also be absent when delegating to code defaults
5. This logic is consistent across in-process, RPC, and OFREP modes

This approach clearly differentiates between:

- Server returning an actual value with a variant (both fields present)
- Server delegating to code default (both fields absent)

**Example evaluation flow**:

```javascript
// Client code
const value = client.getBooleanValue('my-feature', false, context);

// Server evaluates flag with defaultVariant: null
// No targeting match, so server returns:
{
  "reason": "DEFAULT",
  "metadata": {}
  // Note: both "value" and "variant" fields are omitted
}

// Client detects missing value field and uses its default (false)
// Resolution details show reason: "DEFAULT" with undefined variant
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
- Neutral, because existing configurations continue to work unchanged

### Remote Evaluation Considerations

This feature works consistently across all flagd evaluation modes through a unified pattern:

1. **In-Process Mode**: Direct evaluation returns responses without value fields when delegating
2. **RPC Mode (gRPC/HTTP)**: Responses omit the type-specific value field to signal delegation
3. **OFREP (OpenFeature Remote Evaluation Protocol)**: HTTP responses omit the value field

The key insight is using field omission to indicate "use code default":

- Present value field (even if null/false/empty) = use this value
- Absent value field = use your code default
- Works identically across all protocols and evaluation modes
- No protocol changes required

### Implementation Plan

1. Update flagd-schemas with new JSON schema supporting null default variants
2. Implement core logic in flagd to handle null defaults and omit value/variant fields
3. Update flagd-testbed with comprehensive test cases for all evaluation modes
4. Update OpenFeature providers to handle responses without value fields
5. Documentation updates, migration guides, and proto comment additions

Note: No protobuf schema changes are required, but implementation must carefully handle field omission

### Testing Considerations

To ensure correct implementation across all components:

1. **Wire Format Tests**: Verify that protobuf messages with omitted fields are correctly serialized without the fields (not with zero values)
2. **Provider Tests**: Each provider must have tests confirming they detect missing fields correctly in their language
3. **Integration Tests**: End-to-end tests across different language combinations (e.g., Go flagd with Java provider)
4. **OFREP Tests**: Verify JSON responses correctly omit fields (not set to null)
5. **Backward Compatibility Tests**: Ensure old providers handle new responses gracefully
6. **Consistency Tests**: Verify identical behavior across in-process, RPC, and OFREP modes

### Open questions

- How should providers handle responses with missing value fields in strongly-typed languages?
- Should we support both `null` and absent `defaultVariant` fields, or choose one approach?
- What migration path should we recommend for users currently using workarounds?
- Should this feature be gated behind a configuration flag during initial rollout?
- How do we ensure consistent behavior across all provider implementations?
- Should providers validate that the reason is "DEFAULT" when value is omitted, or accept any omitted value as delegation?
- How do we handle edge cases where network protocols might strip empty fields?
- When the client uses its code default after receiving a delegation response, what variant should be reported in telemetry/analytics?
- Should we add explicit proto comments documenting the field omission behavior?

## More Information

- [OpenFeature Specification - Flag Evaluation](https://openfeature.dev/specification/types#flag-evaluation)
- [flagd Flag Definitions Reference](https://flagd.dev/reference/flag-definitions/)
- [flagd JSON Schema Repository](https://github.com/open-feature/flagd-schemas)
- [flagd Testbed](https://github.com/open-feature/flagd-testbed)
