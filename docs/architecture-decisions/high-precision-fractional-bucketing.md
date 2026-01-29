---
# Valid statuses: draft | proposed | rejected | accepted | superseded
status: draft
author: Michael Beemer
created: 2025-09-10
updated: 2026-01-29
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
- **Option 4 (Favored): Max 32-bit signed integer buckets** - Use `math.MaxInt32` (2,147,483,647) as the bucket count, equal to the maximum allowed weight sum. This naturally sidesteps minimum allocation guarantees and excess bucket handling

## Proposal: Max Int32 Bucket Count (Favored)

> **Amendment (2026-01-29):** This alternative is now favored over the original 100,000-bucket proposal.

### Rationale

After experimentation comparing static vs dynamic bucket sizes, and considering implementation complexity, a simpler approach emerged: use the maximum 32-bit signed integer value (`math.MaxInt32` = 2,147,483,647) as the bucket count.

This value is already established as the maximum allowed weight sum for cross-language compatibility. By making the bucket count equal to this maximum, we gain significant simplifications.

### Advantages

Since the total weight sum cannot exceed `math.MaxInt32`, and the bucket count equals `math.MaxInt32`, any variant with a weight of at least 1 is guaranteed at least 1 bucket. This **naturally sidesteps** the need for:

- **Minimum Allocation Guarantee** (as described above): A weight of 1 out of any valid total will always yield at least 1 bucket—no special handling required
- **Excess Bucket Management**: Without minimum allocation adjustments, bucket totals don't exceed the bucket count

### Simplified Implementation

```go
const bucketCount = math.MaxInt32 // 2,147,483,647

func distributeValue(value string, feDistribution *fractionalEvaluationDistribution) string {
    if feDistribution.totalWeight == 0 {
        return ""
    }

    hashValue := int32(murmur3.StringSum32(value))
    hashRatio := math.Abs(float64(hashValue)) / math.MaxInt32
    bucket := int64(hashRatio * float64(feDistribution.totalWeight))

    var rangeEnd int64 = 0
    for _, variant := range feDistribution.weightedVariants {
        rangeEnd += int64(variant.weight)
        if bucket < rangeEnd {
            return variant.variant
        }
    }

    return ""
}
```

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

### Benefits Over Original Proposal

- **Simpler**: No minimum allocation guarantee logic needed
- **No minimum allocation guarantee needed**: With smaller fixed bucket counts, a configuration like `["variant-a", 1], ["variant-b", 1000000]` could round variant-a to 0 buckets (0% traffic). Special handling was needed to guarantee at least 1 bucket. With MaxInt32 buckets equal to the max weight sum, any weight ≥1 naturally gets at least 1 bucket.
- **No excess bucket handling**: With fixed bucket counts (100, 10,000, 100,000), minimum allocation adjustments could cause the total allocated buckets to exceed the bucket count, requiring complex logic to redistribute the excess. With MaxInt32 buckets, allocations naturally sum to the total weight.
- **Same validation**: Weight sum validation against `math.MaxInt32` remains unchanged  
- **Backwards compatible**: Existing configurations continue to work
- **Effectively infinite precision**: Precision limited only by the total weight sum (up to ~0.00000005%)
- **~25-35% less user reassignment**: Experimental testing showed reduced "thrashing" compared to purely dynamic bucket sizes when configurations change

### Consequences

- Good, because implementation is significantly simpler
- Good, because it eliminates surprising edge-case behaviors (minimum allocation, excess handling)
- Good, because validation logic remains the same
- Good, because it provides effectively unlimited precision for practical use cases
- Good, because experimental testing showed less user reassignment than dynamic alternatives
- Bad, because it represents a behavioral breaking change for existing configurations (just the bucket assignment, same as original proposal)
- Neutral, performance is comparable—division by large 32-bit values is not meaningfully slower
