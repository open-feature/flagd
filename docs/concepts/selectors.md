# Selectors

Selectors are query expressions that allow you to filter flag configurations from flagd's sync service. They enable providers to request only specific subsets of flags instead of receiving all flags, making flagd more efficient and flexible for complex deployments.

## Overview

In flagd, **selectors** provide a way to query flags based on different criteria. This is particularly powerful because flagd decouples **flag sources** from **flag sets**, allowing for more granular control over which flags are synchronized and evaluated.

### Key Concepts

- **Flag Source**: Where flag configuration data comes from (file, HTTP endpoint, gRPC service, etc.)
- **Flag Set**: A logical grouping of flags identified by a `flagSetId`
- **Selector**: A query expression that filters flags by source, flag set, or other criteria
- **Flag Set Metadata**: The selector information is "reflected" back in response metadata for transparency

See the [cheat sheet](../reference/cheat-sheet.md#using-the-selector-header) for practical examples of using selectors.

!!! tip

    The `flagSetId` + `key` combination represents the unique identifier for a flag.
    Be sure not to create duplicates, or unexpected behavior may result.
    See [Array-Based Flag Definitions](#array-based-flag-definitions) for how this enables flags with the same key to coexist in different flag sets.

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

### Array-Based Flag Definitions

Flags can be defined as an array instead of an object, with each flag specifying its `key` explicitly:

```json
{
  "flags": [
    {
      "key": "checkout-flow",
      "state": "ENABLED",
      "variants": {"on": true, "off": false},
      "defaultVariant": "on",
      "metadata": { "flagSetId": "payment-service" }
    },
    {
      "key": "checkout-flow",
      "state": "DISABLED",
      "variants": {"on": true, "off": false},
      "defaultVariant": "off",
      "metadata": { "flagSetId": "user-service" }
    }
  ]
}
```

This format is useful for systems that generate large flag configurations programmatically. It also allows flags with the same key to coexist when they belong to different flag sets, since the `flagSetId` + `key` combination represents the unique identifier for a flag.

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
    "bill-buddy-integration": {
      "state": "DISABLED",
      "variants": { "on": true, "off": false },
      "defaultVariant": "off"
    }
  }
}
```

In this example, both `new-checkout-flow` and `bill-buddy-integration` flags belong to the `payment-service` flag set.

### Flag-Level Configuration

Alternatively, the `flagSetId` can be defined on flag level:

```json
{
  "metadata": {
    "version": "v1.2.0"
  },
  "flags": {
    "new-checkout-flow": {
      "state": "ENABLED",
      "variants": {
        "on": true,
        "off": false
      },
      "defaultVariant": "on",
      "metadata": {
        "flagSetId": "webshop",
        "version": "v1.2.0"
      }
    },
    "bill-buddy-integration": {
      "state": "DISABLED",
      "variants": { "on": true, "off": false },
      "defaultVariant": "off",
      "metadata": {
        "flagSetId": "payment-service",
        "version": "v1.2.0"
      },
    }
  }
}
```

In this example the two flags `new-checkout-flow` and `bill-buddy-integration` flags belong to different flag sets.

### Metadata Integration

Selectors work closely with flagd's metadata system. For advanced patterns like flag-level overrides of `flagSetId` or complex metadata inheritance, see the [Metadata concepts](metadata.md) section.

## Metadata Reflection

When you make a request with a selector, flagd "reflects" the selector information back in the response metadata for transparency and debugging. For complete details on metadata selector reflection, inheritance, and configuration patterns, see the [Metadata concepts](metadata.md) section.

## Use Cases

### Multi-Tenant Applications

```yaml
# Tenant A's flags
selector: "flagSetId=tenant-a"

# Tenant B's flags
selector: "flagSetId=tenant-b"
```

### Component Separation

```yaml
# Web service
selector: "flagSetId=payment-service"
# Web application
selector: "flagSetId=webshop"
```

### Environment Separation

```yaml
# Development environment
selector: "flagSetId=dev-features"

# Production environment
selector: "flagSetId=prod-features"
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
4. **Document Your Schema**: Clearly document your flag set naming conventions for your team
5. **Do Not Duplicate Flags Across Sources**: Make sure that flags with the same key and flagSetId do not exist in multiple sources (relative priority of flags in such configurations is not defined).

## Migration Considerations

The selector enhancement maintains full backward compatibility. See the [migration guide](../guides/migrating-to-flag-sets.md) for detailed guidance on transitioning from source-based to flag-set-based selection patterns.
