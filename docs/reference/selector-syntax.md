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

1. **gRPC Metadata**: `Flagd-Selector` header in gRPC metadata
2. **HTTP Header**: `Flagd-Selector` header in HTTP requests
3. **URL Path**: selector path segment, on endpoints that accept one
4. **Query String**: `selector` query parameter, on endpoints that accept one
5. **Request Body**: `selector` field in protobuf/JSON request body

### Example: Header Precedence

```bash
# gRPC request with both header and body selector
# Header takes precedence
grpcurl -H "Flagd-Selector: flagSetId=production" \
  -d '{"selector": "flagSetId=development"}' \
  localhost:8013 flagd.sync.v1.FlagSyncService/FetchAllFlags

# Result: Uses "flagSetId=production" from header
```

## Metadata Reflection

Flagd reflects selector information back in response metadata, providing transparency about query execution. For complete details on metadata selector reflection, inheritance patterns, and configuration examples, see the [Metadata concepts](../concepts/metadata.md) section.

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

### No Flag Set Selection

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

**Flag Configuration Endpoint:**

- `GET /v1/flags`: Supports selector in header or `selector` query parameter
- `GET /v1/flags/{selector}`: Supports selector as a URL-escaped path segment

All HTTP endpoints support the `Flagd-Selector` header for selector specification.

### Path Segment Encoding

The selector occupies a single path segment, so it must be URL-escaped.
Source selectors routinely contain `/`, which would otherwise split the segment:

```bash
# source=./flags.json
curl http://localhost:8015/v1/flags/source%3D.%2Fflags.json
```

No file extension is stripped, so a selector whose value ends in `.json` is preserved as written.

### Error Responses

An unresolvable selector is reported differently depending on why it failed:

| Result                                       | Example                           | Status                     |
|----------------------------------------------|-----------------------------------|----------------------------|
| Selector string is not well-formed           | control characters, invalid UTF-8 | `400`                      |
| Well-formed, but names an unknown filter     | `bogus=1`                         | `404`                      |
| Valid filter that currently matches no flags | `flagSetId=empty-set`             | `200` with an empty result |

Note that the flag configuration endpoint returns `404` for an unknown filter, while the OFREP endpoints return `400` and the sync service RPCs return `invalid_argument` for the same input. Error responses never echo the submitted expression back.
