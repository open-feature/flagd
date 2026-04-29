# Metadata

Metadata in flagd provides contextual information about flags and flag sets.
It enables rich observability, logical separation, and debugging capabilities.

## Overview

Flagd supports metadata at two levels:

- **Flag Set-Level Metadata**: Applied to entire flag configurations
- **Flag-Level Metadata**: Applied to individual flags  

## Metadata Inheritance

Flagd uses a hierarchical metadata system where flags inherit metadata from their containing flag set, with the ability to override specific values at the flag level.

### Flag Set-Level Metadata

The most common pattern is defining metadata at the configuration level, where all flags inherit it:

```json
{
  "metadata": {
    "flagSetId": "payment-service",
    "team": "payments",
    "version": "v1.2.0",
    "environment": "production"
  },
  "flags": {
    "checkout-flow": {
      "state": "ENABLED",
      "variants": {"on": true, "off": false},
      "defaultVariant": "on"
      // Inherits all set-level metadata
    },
    "payment-gateway": {
      "state": "DISABLED", 
      "variants": {"on": true, "off": false},
      "defaultVariant": "off"
      // Also inherits all set-level metadata
    }
  }
}
```

### Flag-Level Overrides

Individual flags can override inherited metadata or add flag-specific metadata:

```json
{
  "metadata": {
    "flagSetId": "payment-service",
    "team": "payments",
    "version": "v1.2.0"
  },
  "flags": {
    "standard-feature": {
      "state": "ENABLED",
      "variants": {"on": true, "off": false},
      "defaultVariant": "on"
      // Inherits: flagSetId="payment-service", team="payments", version="v1.2.0"
    },
    "experimental-feature": {
      "state": "DISABLED",
      "variants": {"on": true, "off": false},
      "defaultVariant": "off",
      "metadata": {
        // Still inherits: flagSetId="payment-service", version="v1.2.0"
        "team": "marketing",        // Override: different flag set
        "owner": "Tom",             // Addition: flag-specific metadata
        "experimental": true        // Addition: flag-specific metadata
      }
    }
  }
}
```

### Inheritance Behavior

1. **Default Inheritance**: Flags inherit all set-level metadata
2. **Selective Override**: Flag-level metadata overrides specific inherited values
3. **Additive Enhancement**: Flag-level metadata can add new keys not present at set level
4. **Preserved Inheritance**: Non-overridden set-level metadata remains inherited

## Metadata Reflection

Metadata reflection provides transparency by echoing selector and configuration information back in API responses. This enables debugging, auditing, and verification of flag targeting.

### Selector Reflection

When making requests with selectors, flagd "reflects" the parsed selector information in the "top-level" `metadata` field:

**Request:**

```bash
curl -H "Flagd-Selector: flagSetId=payment-service" \
  http://localhost:8014/ofrep/v1/evaluate/flags
```

**Response includes reflected metadata:**

```json
{
  "flags": {
    "checkout-flow": {
      "key": "checkout-flow",
      "value": true,
      "variant": "on",
      "metadata": {
        "flagSetId": "payment-service",
        "team": "payments"
      }
    }
  },
  "metadata": {
    "flagSetId": "payment-service"  // Reflected from selector
  }
}
```

### Configuration Reflection

Flag evaluation responses include the complete merged metadata for each flag:

```json
{
  "key": "experimental-feature",
  "value": false,
  "variant": "off", 
  "metadata": {
    "flagSetId": "experiments",       // Overridden at flag level
    "owner": "research-team",         // Added at flag level
    "experimental": true,             // Added at flag level
    "team": "payments",               // Inherited from set level
    "version": "v1.2.0"               // Inherited from set level
  }
}
```

## Common Metadata Fields

### Standard Fields

Some metadata fields are defined in the flag-definition schema for common use-cases:

- **`flagSetId`**: Logical grouping identifier for selectors
- **`version`**: Configuration or flag version

### Custom Fields

You can define any custom metadata fields relevant to your use case:

```json
{
  "metadata": {
    "flagSetId": "user-service",
    "version": "v34",
    "costCenter": "engineering",
    "compliance": "pci-dss",
    "lastReviewed": "2024-01-15",
    "approver": "team-lead"
  }
}
```

## Retrieving Metadata in the OpenFeature SDK

Flag metadata is available in evaluation details returned by flag evaluations.

### Go

```go
details, err := client.BooleanValueDetails(ctx, "new-checkout-flow", false, evalCtx)

// Access metadata from evaluation details
metadata := details.FlagMetadata
flagSetId := metadata["flagSetId"]
team := metadata["team"]
```

### Java

```java
FlagEvaluationDetails<Boolean> details = client.getBooleanDetails(
    "new-checkout-flow", false, new ImmutableContext());

// Access metadata from evaluation details
ImmutableMetadata metadata = details.getFlagMetadata();
String flagSetId = metadata.getString("flagSetId");
String team = metadata.getString("team");
```

### JavaScript

```javascript
const details = await client.getBooleanDetails('new-checkout-flow', false, {});

// Access metadata from evaluation details
const metadata = details.flagMetadata;
const flagSetId = metadata.flagSetId;
const team = metadata.team;
```

## Use Cases

**Debugging**: Metadata reflection shows which selectors were used and how inheritance resolved, making it easier to troubleshoot flag targeting issues.

**Governance**: Track team ownership, compliance requirements, and approval workflows through custom metadata fields.

**Environment Management**: Use metadata for version tracking, environment identification, and change management across deployments.

**Multi-Tenancy**: Isolate tenants through flag sets and maintain tenant-specific configurations and governance.

**Observability**: Metadata attributes can be used in telemetry spans and metrics, providing operational visibility into flag usage patterns and configuration context.  
