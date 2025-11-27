# Selectors

Selectors are query expressions that allow you to filter flag configurations from flagd's sync service. They enable providers to request only specific subsets of flags instead of receiving all flags, making flagd more efficient and flexible for complex deployments.

## Overview

In flagd, **selectors** provide a way to query flags based on different criteria. This is particularly powerful because flagd decouples **flag sources** from **flag sets**, allowing for more granular control over which flags are synchronized and evaluated.

### Key Concepts

- **Flag Source**: Where flag configuration data comes from (file, HTTP endpoint, gRPC service, etc.)
- **Flag Set**: A logical grouping of flags identified by a `flagSetId`
- **Selector**: A query expression that filters flags by source, flag set, or other criteria
- **Flag Set Metadata**: The selector information is "reflected" back in response metadata for transparency

## Source vs Flag Set Decoupling

### Before: Tight Coupling

Historically, each source provided exactly one flag set, and providers had to target specific sources:

```yaml
# Old approach - targeting a specific source
selector: "my-flag-source.json"
```

### After: Flexible Flag Sets

Now, sources and flag sets are decoupled. A single source can contain multiple flag sets, and flag sets can span multiple sources:

```yaml
# New approach - targeting a logical flag set
selector: "flagSetId=project-42"
```

## Flag Set Configuration

Flag sets are typically configured at the top level of a flag configuration, with all flags in that configuration inheriting the same `flagSetId`. This is the recommended approach for most use cases.

### Set-Level Configuration

The most common pattern is to set the `flagSetId` at the configuration level, where all flags inherit it:

```json
{
  "metadata": {
    "flagSetId": "payment-service",
    "version": "v1.2.0"
  },
  "flags": {
    "new-checkout-flow": {
      "state": "ENABLED",
      "variants": {
        "on": true,
        "off": false
      },
      "defaultVariant": "on"
    },
    "stripe-integration": {
      "state": "DISABLED",
      "variants": { "on": true, "off": false },
      "defaultVariant": "off"
    }
  }
}
```

In this example, both `new-checkout-flow` and `stripe-integration` flags belong to the `payment-service` flag set.

### Metadata Inheritance and Override

Flagd uses a hierarchical metadata system:

1. **Set-Level Metadata**: Defined in the top-level `metadata` section, inherited by all flags
2. **Flag-Level Metadata**: Defined in individual flag `metadata`, overrides set-level values for that flag
3. **Merged Result**: Flag evaluations return merged metadata with flag-level taking precedence

### Flag-Level Overrides (Advanced)

For advanced use cases, individual flags can override the set-level `flagSetId`:

```json
{
  "metadata": {
    "flagSetId": "payment-service",
    "team": "payments"
  },
  "flags": {
    "standard-feature": {
      "state": "ENABLED",
      "variants": { "on": true, "off": false },
      "defaultVariant": "on"
      // Inherits flagSetId: "payment-service"
    },
    "experimental-feature": {
      "metadata": {
        "flagSetId": "experiments",  // Override: belongs to different set
        "owner": "research-team"
      },
      "state": "DISABLED",
      "variants": { "on": true, "off": false },
      "defaultVariant": "off"
    }
  }
}
```

In this example:
- `standard-feature` inherits `flagSetId: "payment-service"` from set level
- `experimental-feature` overrides to `flagSetId: "experiments"`
- Both flags inherit `team: "payments"` (unless overridden at flag level)

## Flag Set Metadata "Reflection"

When you make a request with a selector, flagd "reflects" the selector information back in the response metadata. This provides transparency about what was actually queried and helps with debugging.

### Example

**Request with selector:**
```
Selector: "flagSetId=project-42"
```

**Response includes reflected metadata:**
```json
{
  "flags": { /* ... */ },
  "metadata": {
    "flagSetId": "project-42"
  }
}
```

This helps you:
- Verify that your selector was parsed correctly
- Debug complex selector queries
- Understand exactly what flags were returned
- Audit flag access patterns

## Use Cases

### Multi-Tenant Applications

```yaml
# Tenant A's flags
selector: "flagSetId=tenant-a"

# Tenant B's flags  
selector: "flagSetId=tenant-b"
```

### Environment Separation

```yaml
# Development environment
selector: "flagSetId=dev-features"

# Production environment
selector: "flagSetId=prod-features"
```

### Feature Team Isolation

```yaml
# Payment team's flags
selector: "flagSetId=payments"

# User interface team's flags
selector: "flagSetId=ui-components"
```

### Legacy Source-Based Selection

```yaml
# Still supported for backward compatibility
selector: "source=legacy-config.json"
```

## Best Practices

1. **Use Flag Sets for Logical Grouping**: Prefer `flagSetId` over `source` for new deployments
2. **Plan Your Flag Set Strategy**: Design flag sets around logical boundaries (teams, features, environments)
3. **Leverage Metadata**: Use metadata for debugging and auditing
5. **Document Your Schema**: Clearly document your flag set naming conventions for your team

## Migration Considerations

The selector enhancement maintains full backward compatibility. See the [migration guide](../guides/migrating-to-flag-sets.md) for detailed guidance on transitioning from source-based to flag-set-based selection patterns.
