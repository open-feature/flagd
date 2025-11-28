# Migrating to Flag Sets

This guide helps you transition from source-based selector patterns to flag-set-based patterns, taking advantage of flagd's enhanced selector capabilities while maintaining backward compatibility.

## Understanding the Change

### Before: Source-Based Selection

In the traditional approach, providers targeted specific sources:

```yaml
# Provider configuration targeting a source file
selector: "config/my-flags.json"
```

This created tight coupling between providers and sources:

- Providers had to know which source contained their flags
- Moving flags between sources required provider reconfiguration
- One source could only serve one logical set of flags

### After: Flag Set-Based Selection

With flag sets, providers target logical groupings of flags:

```yaml
# Provider configuration targeting a flag set
selector: "flagSetId=my-application"
```

This provides flexibility:

- Providers are decoupled from sources
- Sources can contain multiple flag sets
- Flag sets can span multiple sources
- No breaking changes - old selectors still work

## Migration Process

### Step 1: Add Flag Set IDs to Configurations

Add `flagSetId` to your flag configurations at the set level:

**Before:**

```json
{
  "flags": {
    "feature-a": {
      "state": "ENABLED",
      "variants": {"on": true, "off": false},
      "defaultVariant": "on"
    }
  }
}
```

**After:**

```json
{
  "metadata": {
    "flagSetId": "my-application"
  },
  "flags": {
    "feature-a": {
      "state": "ENABLED",
      "variants": {"on": true, "off": false},
      "defaultVariant": "on"
    }
  }
}
```

### Step 2: Update Provider Configurations

Change provider selectors from source-based to flag set-based:

```java
// Before
new FlagdProvider(FlagdOptions.builder()
    .selector("config/app-flags.json").build());

// After
new FlagdProvider(FlagdOptions.builder()
    .selector("flagSetId=my-application").build());
```

### Step 3: Verify Migration

Test that selectors work correctly and check metadata reflection:

```bash
curl -H "Flagd-Selector: flagSetId=my-application" \
  http://localhost:8014/ofrep/v1/evaluate/flags
```

## Flag Set Organization Patterns

**By Application/Service:**

```yaml
flagSetId: "user-service"    # All user-related flags
flagSetId: "payment-service" # All payment-related flags
```

**By Environment:**

```yaml
flagSetId: "development"     # Dev-specific flags
flagSetId: "production"      # Production flags
```

**By Team:**

```yaml
flagSetId: "frontend-team"   # Frontend features
flagSetId: "backend-team"    # Backend features
```

Choose the pattern that best matches your deployment and organizational structure.

## Common Issues

**No flags returned**: Check that `flagSetId` in selector matches flag configuration exactly

**Wrong flags returned**: Look for flag-level `flagSetId` overrides or header/body selector conflicts

**Selector ignored**: Verify selector syntax is correct (`flagSetId=value`, not `flagSetId:value`)

## Best Practices

- **Group logically**: Organize flags by service, environment, or team
- **Name consistently**: Use clear, descriptive flag set names
- **Test first**: Validate migration in non-production environments
- **Use metadata reflection**: Check reflected metadata for debugging

## FAQ

**Q: Do I have to migrate?**  
A: No, source-based selectors continue to work. Migration is optional but recommended.

**Q: Can flag sets span multiple sources?**  
A: Yes, multiple sources can contribute flags to the same flag set.

## Additional Resources

- [Selector Concepts](../concepts/selectors.md) - Understanding selectors and flag sets
- [Selector Syntax Reference](../reference/selector-syntax.md) - Complete syntax documentation  
- [ADR: Decouple Flag Source and Set](../architecture-decisions/decouple-flag-source-and-set.md) - Technical decision rationale
