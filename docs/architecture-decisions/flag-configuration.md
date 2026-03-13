---
status: accepted
author: Todd Baert
created: 2025-06-05
updated: 2025-06-05
---

# Flag and Targeting Configuration

## Background

Feature flag systems require a flexible, safe, and portable way to express targeting rules that can evaluate contextual data to determine which variant of a feature to serve.

flagd's targeting system was designed with several key requirements:

## Requirements

- **Language agnostic**: Rules must be portable across different programming languages, ideally relying on existing expression language(s)
- **Safe evaluation**: No arbitrary code execution or system access
- **Deterministic**: Same inputs must always produce same outputs
- **Extensible**: Support for the addition of domain-specific operations relevant to feature flags
- **Developer and machine friendly**: Human-readable, easily validated, and easily serialized

## Proposal

### JSON Logic as the Foundation

flagd chose **JSON Logic** as its core evaluation engine, implementing a modified version with custom extensions.
This provides a secure, portable foundation where rules are expressed as JSON objects with operators as keys and parameters as values.

#### Benefits realized

- Rules can be stored in databases, transmitted over networks, shared between frontend/backend, and embedded in Kubernetes custom resources
- No eval() or code injection risks - computations are deterministic and sand-boxed
- Implementations exist in most languages

#### Overview

The system provides two tiers of operators:

##### Primitive JSON Logic Operators (inherited from the JSONLogic)

- Logical: `and`, `or`, `!`, `!!`
- Comparison: `==`, `!=`, `>`, `<`, etc
- Arithmetic: `+`, `-`, `*`, `/`, `%`
- Array operations: `in`, `map`, `filter`, etc
- String operations: `cat`, `substr`, etc
- Control flow: `if`
- Assignment and extraction: `var`

##### Custom flagd Extensions

- `fractional`: Deterministic percentage-based distribution using murmur3 hashing
- `starts_with`/`ends_with`: String prefix/suffix matching for common patterns
- `regex_match`: String regular expression matching
- `sem_ver`: Semantic version comparisons with standard (npm-style) operators
- `$ref`: Reference to shared evaluators for DRY principle

##### Evaluation Context and Automatic Enrichment

flagd automatically injects critical context values:

##### System-provided context

- `$flagd.flagKey`: The flag being evaluated (available v0.6.4+)
- `$flagd.timestamp`: Unix timestamp of evaluation (available v0.6.7+)

This enables sophisticated targeting rules that can reference the flag itself or time-based conditions without requiring client-side context.

##### Reason Code System for Transparency

flagd returns specific reason codes with every evaluation to indicate how the decision was made:

1. **STATIC**: Flag has no targeting rules, and can be safely cached
2. **TARGETING_MATCH**: Targeting rules matched and returned a variant
3. **DEFAULT**: Targeting rules evaluated to null, fell back to default
4. **CACHED**: Value retrieved from provider cache (RPC mode only)
5. **ERROR**: Evaluation failed due to invalid configuration

This transparency enables:

- Appropriate caching strategies (only STATIC flags are cached)
- Improved debugging, telemetry, and monitoring of flag behavior

##### Shared Evaluators for Reusability

The `$evaluators` top-level property enables shared targeting logic:

```json
{
  "$evaluators": {
    "isEmployee": {
      "ends_with": [{"var": "email"}, "@company.com"]
    }
  },
  "flags": {
    "feature-x": {
      "state": "ENABLED",
      "defaultVariant": "enabled",
      "variants": {
        "enabled": true,
        "disabled": false
      },
      "targeting": {
        "if": [{"$ref": "isEmployee"}, "enabled", "disabled"]
      }
    }
  }
}
```

##### Intelligent Caching Strategy

Only flags with reason **STATIC** are cached, as they have deterministic outputs. This ensures:

- Maximum cache efficiency for simple toggles
- Fresh evaluation for complex targeting rules
- Cache invalidation on configuration changes

##### Schema-Driven Configuration

Two schemas validate flag configurations:

- `https://flagd.dev/schema/v0/flags.json`: Overall flag structure
- `https://flagd.dev/schema/v0/targeting.json`: Targeting rule validation

These enable:

- IDE support with autocomplete
- Run-time and build-time validation
- Separate validation of rules and overall configuration if desired

## Considered Options

- **Custom DSL**: Would require parsers in every language
- **JavaScript/Lua evaluation**: Security risks and language lock-in
- **CEL**: limited number of implementations at time of decision, can't be directly parsed/validated when embedded in Kubernetes resources

## Consequences

### Positive

- Good, because implementations exist across languages
- Good, because, no code injection or system access possible
- Good, because combined with JSON schemas, we have rich IDE support
- Good, because JSON is easily serialized and also can be represented/embedded in YAML

### Negative

- Bad, JSONLogic syntax can be cumbersome when rules are complex
- Bad, hard to debug

## Conclusion

flagd's targeting configuration system represents a thoughtful balance between safety, portability, and capability.
By building on JSON Logic and extending it with feature-flag-specific operators, flagd achieves remarkable flexibility while maintaining security and performance.
