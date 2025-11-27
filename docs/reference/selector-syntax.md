# Selector Syntax Reference

This document provides the complete technical specification for flagd selector syntax, including supported operators, precedence rules, and metadata reflection behavior.

## Syntax Overview

Selectors use a simple key-value syntax to filter flags. Currently, selectors support single key-value pairs with plans to expand to more complex queries in the future.

### Basic Syntax

```text
<key>=<value>
```

### Backward Compatibility Syntax

```text
<value>
```

When no `=` is present, the value is treated as a source selector for backward compatibility.

## Supported Keys

### `flagSetId`

Selects flags belonging to a specific flag set.

**Syntax:**

```text
flagSetId=<set-identifier>
```

**Examples:**

```text
flagSetId=project-42
flagSetId=dev-environment
flagSetId=team-payments
```

**Special Case - Empty Flag Set:**

```text
flagSetId=
```

Selects flags that don't belong to any named flag set (equivalent to the "null" flag set).

### `source`

Selects flags from a specific source.

**Syntax:**

```text
source=<source-identifier>
```

**Examples:**

```text
source=config/flags.json
source=http://flag-server/config
source=./local-flags.yaml
```

## Selector Precedence

When selectors are provided in multiple locations, flagd uses the following precedence order (highest to lowest):

1. **gRPC Header**: `Flagd-Selector` header in gRPC metadata
2. **HTTP Header**: `Flagd-Selector` header in HTTP requests
3. **Request Body**: `selector` field in protobuf/JSON request body

### Example: Header Precedence

```bash
# gRPC request with both header and body selector
# Header takes precedence
grpcurl -H "Flagd-Selector: flagSetId=production" \
  -d '{"selector": "flagSetId=development"}' \
  localhost:8013 flagd.sync.v1.FlagSyncService/FetchAllFlags

# Result: Uses "flagSetId=production" from header
```

## Metadata "Reflection"

Flagd reflects selector information back in response metadata, providing transparency about query execution.

### Reflection Behavior

**Input Selector:**

```text
flagSetId=project-42
```

**Reflected in Response Metadata:**

```json
{
  "metadata": {
    "flagSetId": "project-42"
  }
}
```

### Multiple Metadata Sources

Reflected metadata includes:

- **Selector Information**: The parsed selector key-value pairs
- **Set-Level Metadata**: Metadata from the flag configuration itself
- **Source Context**: Additional context from sync operations

### Metadata Inheritance

Flagd uses a hierarchical metadata system where flags inherit metadata from their flag set:

**Set-Level Metadata (Inherited):**

```json
{
  "metadata": {
    "flagSetId": "payment-service",
    "team": "payments",
    "version": "1.2.0"
  },
  "flags": {
    "checkout-flow": {
      "state": "ENABLED"
      // Inherits all set-level metadata
    }
  }
}
```

**Flag-Level Metadata (Override):**

```json
{
  "metadata": {
    "flagSetId": "payment-service",
    "team": "payments"
  },
  "flags": {
    "experimental-feature": {
      "metadata": {
        "flagSetId": "experiments",  // Overrides set-level
        "owner": "research-team"      // Adds flag-specific metadata
        // Still inherits "team": "payments"
      },
      "state": "DISABLED"
    }
  }
}
```

## Examples

### Flag Set Selection

```bash
# Select flags from the "payments" flag set
curl -H "Flagd-Selector: flagSetId=payments" \
  http://localhost:8014/ofrep/v1/evaluate/flags
```

### Source Selection (Legacy)

```bash
# Select flags from a specific source (backward compatibility)
curl -H "Flagd-Selector: source=config/prod-flags.json" \
  http://localhost:8014/ofrep/v1/evaluate/flags
```

### Empty Flag Set Selection

```bash
# Select flags that don't belong to any named flag set
curl -H "Flagd-Selector: flagSetId=" \
  http://localhost:8014/ofrep/v1/evaluate/flags
```

### Provider SDK Usage

#### Go Provider

```go
import "github.com/open-feature/go-sdk-contrib/providers/flagd"

provider := flagd.NewProvider(
    flagd.WithHost("localhost"),
    flagd.WithPort(8013),
    flagd.WithSelector("flagSetId=user-service"),
)
```

#### Java Provider

```java
FlagdProvider provider = new FlagdProvider(
    FlagdOptions.builder()
        .host("localhost")
        .port(8013)
        .selector("flagSetId=payment-service")
        .build()
);
```

#### JavaScript Provider

```javascript
const provider = new FlagdProvider({
  host: 'localhost',
  port: 8013,
  selector: 'flagSetId=frontend-features'
});
```

## Future Enhancements

The selector syntax is designed to be extensible. Future versions may support:

- **Multiple Criteria**: `flagSetId=app1,source=prod`
- **Complex Queries**: `flagSetId=app1 OR flagSetId=app2`  
- **Filter Expressions**: `metadata.environment=production`
- **Kubernetes-Style Selectors**: `app=frontend,tier=web`

> **Note**: The current implementation supports single key-value pairs only. Complex selectors are planned for future releases.

## API Reference

### gRPC Services

**Sync Service:**

- `SyncFlags(SyncFlagsRequest)`: Supports selector in header and request body
- `FetchAllFlags(FetchAllFlagsRequest)`: Supports selector in header and request body

**Evaluation Service:**

- `ResolveBoolean(ResolveBooleanRequest)`: Supports selector in header
- `ResolveString(ResolveStringRequest)`: Supports selector in header
- `ResolveInt(ResolveIntRequest)`: Supports selector in header
- `ResolveFloat(ResolveFloatRequest)`: Supports selector in header
- `ResolveObject(ResolveObjectRequest)`: Supports selector in header
- `ResolveAll(ResolveAllRequest)`: Supports selector in header

### HTTP/OFREP Services

**OFREP Endpoints:**

- `POST /ofrep/v1/evaluate/flags/{key}`: Supports selector in header
- `POST /ofrep/v1/evaluate/flags`: Supports selector in header

All HTTP endpoints support the `Flagd-Selector` header for selector specification.
