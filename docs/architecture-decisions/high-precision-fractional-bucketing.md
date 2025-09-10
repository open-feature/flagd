---
# Valid statuses: draft | proposed | rejected | accepted | superseded
status: draft
author: Michael Beemer
created: 2025-09-10
updated: 2025-09-10
---

# High-Precision Fractional Bucketing for Sub-Percent Traffic Allocation

This ADR proposes enhancing the fractional operation to support high-precision traffic allocation down to 0.001% granularity by increasing the internal bucket count from 100 to 100,000 while maintaining the existing weight-based API.

## Background

The current fractional operation in flagd uses a 100-bucket system that maps hash values to percentages in the range [0, 100].
This approach works well for most use cases but has significant limitations in high-throughput environments where precise sub-percent traffic allocation is required.

Currently, the smallest allocation possible is 1%, which is insufficient for:

- Gradual rollouts in ultra-high-traffic systems where 1% could represent millions of users
- A/B testing scenarios requiring precise control over small experimental groups
- Canary deployments where operators need to start with very small traffic percentages (e.g., 0.1% or 0.01%)

The current implementation in `fractional.go` calculates bucket assignment using:

```go
bucket := hashRatio * 100 // in range [0, 100]
```

This limits granularity to 1% increments, making it impossible to achieve the precision required for sophisticated traffic management strategies.

## Requirements

- Support traffic allocation precision down to 0.001% (3 decimal places)
- Maintain backwards compatibility with existing weight-based API
- Preserve deterministic bucketing behavior (same hash input always produces same bucket)
- Ensure consistent bucket assignment across different programming languages
- Support weight values up to a reasonable maximum that works across multiple languages
- Maintain current performance characteristics
- Prevent users from being moved between buckets when only distribution percentages change
- Guarantee that any variant with weight > 0 receives some traffic allocation
- Handle edge cases gracefully without silent failures
- Validate weight configurations and provide clear error messages for invalid inputs

## Considered Options

- **Option 1: 10,000 buckets (0.01% precision)** - 1 in every 10,000 users, better but still not sufficient for many high-throughput use cases
- **Option 2: 100,000 buckets (0.001% precision)** - 1 in every 100,000 users, meets most high-precision needs
- **Option 3: 1,000,000 buckets (0.0001% precision)** - 1 in every 1,000,000 users, likely overkill and could impact performance

## Proposal

Implement a 100,000-bucket system that provides 0.001% precision while maintaining the existing integer weight-based API.

### API changes

No API changes are required. The existing fractional operation syntax remains unchanged:

```json
"fractional": [
  { "cat": [{ "var": "$flagd.flagKey" }, { "var": "email" }] },
  ["red", 50],
  ["blue", 30], 
  ["green", 20]
]
```

### Implementation Changes

1. **Bucket Count**: Change from 100 to 100,000 buckets by modifying bucket calculation from `hashRatio * 100` to `hashRatio * 100000`
2. **Minimum Allocation Guarantee**: Any variant with weight > 0 receives at least 1 bucket (0.001%)
3. **Excess Bucket Handling**: Remove excess buckets from the largest variant to maintain exactly 100,000 total buckets
4. **Weight Sum Validation**: Reject configurations where total weight exceeds maximum safe integer value
5. **Maximum Weight Sum**: Use language-specific maximum 32-bit signed integer constants for cross-platform compatibility

### Minimum Allocation Guarantee

To prevent silent configuration failures, any variant with a positive weight will receive at least 0.001% allocation (1 bucket), even if the calculated percentage would round to zero. This ensures predictable behavior where positive weights always result in some traffic allocation.

**Example**: Configuration `["variant-a", 1], ["variant-b", 1000000]`

- Without guarantee: variant-a gets 0% (never selected)
- With guarantee: variant-a gets 0.001%, variant-b gets 99.999%

### Excess Bucket Management

When minimum allocations cause the total to exceed 100,000 buckets, excess buckets are removed from the variant with the largest allocation.
This approach:

- Maintains the minimum guarantee for small variants
- Has minimal impact on large variants (small relative reduction)
- Preserves deterministic behavior
- Prevents bucket count overflow

### Weight Sum Validation

When the total weight sum exceeds the maximum safe integer value, the fractional evaluation will return a validation error with a clear message.
This prevents integer overflow issues and provides immediate feedback to users about invalid configurations.

```go
import "math"

func validateWeightSum(variants []fractionalEvaluationVariant) error {
    var totalWeight int64 = 0
    for _, variant := range variants {
        totalWeight += int64(variant.weight)
        if totalWeight > math.MaxInt32 {
            return fmt.Errorf("total weight sum %d exceeds maximum allowed value %d", 
                totalWeight, math.MaxInt32)
        }
    }
    return nil
}
```

Implementations should prefer built-in language constants (e.g., `math.MaxInt32` in Go, `Integer.MAX_VALUE` in Java, `int.MaxValue` in C#) rather than hardcoded values to ensure maintainability and clarity.

### Edge Case Handling

The implementation addresses several edge cases:

1. **All weights are 0**: Returns empty string (maintains current behavior)
2. **Negative weights**: Treated as 0 (maintains current validation behavior)
3. **Single variant**: Receives all 100,000 buckets regardless of weight value
4. **Empty variants**: Returns error (maintains current validation behavior)
5. **Weight sum overflow**: Returns validation error with clear message
6. **Multiple variants with minimum allocation**: Excess distributed fairly among largest variants

### Maximum Weight Considerations

To ensure cross-language compatibility, we establish a maximum total weight sum equal to the maximum 32-bit signed integer value (2,147,483,647). This limit:

- Works reliably across all target languages (Go, Java, .NET, JavaScript, Python)
- Provides more than sufficient range for any practical use case
- Prevents integer overflow issues in 32-bit signed integer systems
- Allows for extremely fine-grained control (individual weights can be 1 out of 2+ billion)
- Uses language-native constants for better maintainability

### Code Changes

The following shows how the core logic in `fractional.go` would be modified.

```go
const bucketCount = 100000

// bucketAllocation represents the number of buckets allocated to a variant
type bucketAllocation struct {
    variant string
    buckets int
}

func (fe *Fractional) Evaluate(values, data any) any {
    valueToDistribute, feDistributions, err := parseFractionalEvaluationData(values, data)
    if err != nil {
        fe.Logger.Warn(fmt.Sprintf("parse fractional evaluation data: %v", err))
        return nil
    }

    if err := validateWeightSum(feDistributions.weightedVariants); err != nil {
        fe.Logger.Warn(fmt.Sprintf("weight validation failed: %v", err))
        return nil
    }

    return distributeValue(valueToDistribute, feDistributions)
}

func validateWeightSum(variants []fractionalEvaluationVariant) error {
    var totalWeight int64 = 0
    for _, variant := range variants {
        totalWeight += int64(variant.weight)
        if totalWeight > math.MaxInt32 {
            return fmt.Errorf("total weight sum %d exceeds maximum allowed value %d", 
                totalWeight, math.MaxInt32)
        }
    }
    return nil
}

func calculateBucketAllocations(variants []fractionalEvaluationVariant, totalWeight int) []bucketAllocation {
    allocations := make([]bucketAllocation, len(variants))
    totalAllocated := 0
    
    // Calculate initial allocations
    for i, variant := range variants {
        if variant.weight == 0 {
            allocations[i] = bucketAllocation{variant: variant.variant, buckets: 0}
        } else {
            // Calculate proportional allocation
            proportional := int((int64(variant.weight) * bucketCount) / int64(totalWeight))
            // Ensure minimum allocation of 1 bucket for any positive weight
            buckets := max(1, proportional)
            allocations[i] = bucketAllocation{variant: variant.variant, buckets: buckets}
        }
        totalAllocated += allocations[i].buckets
    }
    
    // Handle excess buckets by removing from largest allocation
    excess := totalAllocated - bucketCount
    if excess > 0 {
        // Sort indices by bucket count (descending) to find largest allocation
        indices := make([]int, len(allocations))
        for i := range indices {
            indices[i] = i
        }
        sort.Slice(indices, func(i, j int) bool {
            if allocations[indices[i]].buckets == allocations[indices[j]].buckets {
                return allocations[indices[i]].variant < allocations[indices[j]].variant // Tie-break by variant name
            }
            return allocations[indices[i]].buckets > allocations[indices[j]].buckets
        })
        
        // Remove excess from largest allocation, respecting minimum guarantee
        for _, idx := range indices {
            if excess <= 0 {
                break
            }
            
            // Don't reduce below 1 bucket if original weight > 0
            minAllowed := 0
            if variants[idx].weight > 0 {
                minAllowed = 1
            }
            
            canRemove := allocations[idx].buckets - minAllowed
            toRemove := min(excess, canRemove)
            allocations[idx].buckets -= toRemove
            excess -= toRemove
        }
    }
    
    return allocations
}
```

**5. Replace the distribution logic:**

```go
func distributeValue(value string, feDistribution *fractionalEvaluationDistribution) string {
    if feDistribution.totalWeight == 0 {
        return ""
    }
    
    allocations := calculateBucketAllocations(feDistribution.weightedVariants, feDistribution.totalWeight)
    
    hashValue := int32(murmur3.StringSum32(value))
    hashRatio := math.Abs(float64(hashValue)) / math.MaxInt32
    bucket := int(hashRatio * bucketCount) // in range [0, bucketCount)

    currentBucket := 0
    for _, allocation := range allocations {
        currentBucket += allocation.buckets
        if bucket < currentBucket {
            return allocation.variant
        }
    }

    return ""
}
```

### Consequences

- Good, because it enables precise traffic control for high-throughput environments
- Good, because it matches industry-standard precision offered by leading vendors
- Good, because it maintains API backwards compatibility
- Good, because integer weights remain simple to understand and configure
- Good, because it prevents silent configuration failures through minimum allocation guarantee
- Good, because excess handling is predictable and fair
- Good, because weight validation provides clear error messages for invalid configurations
- Bad, because it represents a behavioral breaking change for existing configurations
- Bad, because it slightly increases memory usage for bucket calculations
- Bad, because actual percentages may differ slightly from configured weights due to minimum allocations

### Implementation Plan

1. Update flagd-testbed with comprehensive test cases for high-precision fractional bucketing across all evaluation modes
2. Implement core logic in flagd to support 100,000-bucket system with minimum allocation guarantee and excess handling
3. Update flagd providers to ensure consistent behavior and testing across language implementations
4. Documentation updates, migration guides, and example configurations to demonstrate the new precision capabilities
