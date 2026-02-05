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
- **Option 4 (Favored): Max 32-bit signed integer buckets** - Use `math.MaxInt32` (2,147,483,647) as the maximum allowed weight sum. This naturally sidesteps minimum allocation guarantees and excess bucket handling

## Proposal: Max Int32 Weight Sum (Favored)

> **Amendment (2026-01-29):** This alternative is now favored over the original 100,000-bucket proposal.

### Rationale

After experimentation comparing static vs dynamic bucket sizes, and considering implementation complexity, a simpler approach emerged: use the maximum 32-bit signed integer value (`math.MaxInt32` = 2,147,483,647) as the maximum allowed weight sum.

This value ensures cross-language compatibility. The bucket calculation requires multiplying a 32-bit hash by the total weight, producing a 64-bit intermediate product. The max product (`MaxUint32 × MaxInt32` = 9.22 × 10¹⁸) fits within Java's signed `long` with ~6 billion headroom. Java is the limiting factor — using `MaxUint32` for the weight sum would overflow `long`. JavaScript's `Number` type cannot safely represent the max product, so `BigInt` is required. See [Cross-Language Implementation Notes](#cross-language-implementation-notes) for details.

### Constraints

- The sum of all variant weights must not exceed `math.MaxInt32` (2,147,483,647)
- Weights must be defined as integers

### Advantages

Since the total weight sum cannot exceed `math.MaxInt32`, any variant with a weight of at least 1 is guaranteed at least 1 bucket. This **naturally sidesteps** the need for:

- **Minimum Allocation Guarantee** (as described above): A weight of 1 out of any valid total will always yield at least 1 bucket—no special handling required
- **Excess Bucket Management**: Without minimum allocation adjustments, bucket totals don't exceed the bucket count

### Simplified Implementation

This implementation is designed to be compatible with the "Harden Hashing" ADR, accepting a pre-computed hash value rather than performing string hashing internally. This decouples fractional bucketing from the hashing strategy.

```go
const maxWeightSum = math.MaxInt32 // 2,147,483,647

// distributeValue accepts the hash calculated by the "Harden Hashing" ADR logic.
// It relies purely on integer math, avoiding floating-point precision issues.
// Note: hashValue is uint32 (full 32-bit hash range), while weights are int32
// (max sum of MaxInt32 for cross-language compatibility).
func distributeValue(hashValue uint32, feDistribution *fractionalEvaluationDistribution) string {
    // 0. Validation: Handle empty distribution
    if feDistribution.totalWeight == 0 {
        return ""
    }

    // 1. Use the hash provided 32-bit hash

    // 2. Projection: Map 32-bit hash to [0, totalWeight)
    //    We cast to uint64 to ensure the multiplication does not overflow.
    //    Shifting right by 32 bits is mathematically equivalent to dividing by 2^32.
    //    This logic is safe across major languages because it relies on fundamental
    //    binary operations.
    bucket := (uint64(hashValue) * uint64(feDistribution.totalWeight)) >> 32

    // 3. Selection: Find which variant range the bucket falls into
    var rangeEnd uint64 = 0
    for _, variant := range feDistribution.weightedVariants {
        rangeEnd += uint64(variant.weight) // this would be a Java long, or JS BigInt - needs to handle max product: 9.223372030 × 10^18 (9,223,372,030,412,324,865) 
        if bucket < rangeEnd {
            return variant.variant
        }
    }

    // Unreachable given strict validation of weights (integers, sum <= MaxInt32)
    return ""
}
```

> **Note:** This implementation uses pure integer arithmetic to avoid floating-point precision issues entirely. The expression `(uint64(hashValue) * uint64(totalWeight)) >> 32` is mathematically equivalent to `(hashValue / 2^32) * totalWeight`, but performed in integer space. The Go code uses `uint64`, but each language uses its own 64-bit type (e.g., Java uses `long`). The `MaxInt32` weight constraint ensures the intermediate product fits within Java's more limited signed `long` range, while Go's `uint64` handles it with additional headroom. The right-shift by 32 bits provides exact division by 2^32. This approach is portable across all major languages since it relies only on fundamental binary operations.

### Cross-Language Implementation Notes

MurmurHash3-32 always produces a 32-bit value, but languages differ in how they represent it. The algorithm requires:

1. Treating the hash as an **unsigned** 32-bit integer
2. Performing the multiplication in a 64-bit integer type
3. The `MaxInt32` weight constraint ensures the product fits in Java's signed `long` (the most restrictive common 64-bit type)

| Language | Hash Type | Conversion to Unsigned | 64-bit Multiply | Right-Shift |
|----------|-----------|------------------------|-----------------|-------------|
| **Go** | `uint32` | None needed | `uint64(hash)` | `>> 32` |
| **Java** | `int` (signed) | `hash & 0xFFFFFFFFL` | Use `long` | `>>> 32` (unsigned) |
| **JavaScript** | `Number` | `BigInt(hash)` | Use `BigInt` | `>> 32n` |
| **Python** | `int` | None needed (arbitrary precision) | Native | `>> 32` |
| **C/C++** | `uint32_t` | None needed | `(uint64_t)` | `>> 32` |
| **C#/.NET** | `uint` | None needed | `(ulong)` | `>> 32` |

**Java example:**

```java
int hash = murmur3_32(value);  // signed int
long hashUnsigned = hash & 0xFFFFFFFFL;  // treat as unsigned
long bucket = (hashUnsigned * totalWeight) >>> 32;  // unsigned right-shift
```

**JavaScript example:**

```javascript
const hash = murmur3_32(value);  // Number
const bucket = (BigInt(hash) * BigInt(totalWeight)) >> 32n;
```

### API changes

No API changes are required. The existing fractional operation syntax remains unchanged:

```yaml
# Constraint: The sum of all variant weights must not exceed math.MaxInt32 (2,147,483,647).
# Constraint: Weights must be defined as Integers (can be enforced by JSON schema).
"fractional": [
  { "cat": [{ "var": "$flagd.flagKey" }, { "var": "email" }] },
  ["red", 50],
  ["blue", 30], 
  ["green", 20]
]
```

### Benefits Over Original Proposal

- **Simpler**: No minimum allocation guarantee logic needed
- **No minimum allocation guarantee needed**: With smaller fixed bucket counts, a configuration like `["variant-a", 1], ["variant-b", 1000000]` could round variant-a to 0 buckets (0% traffic). Special handling was needed to guarantee at least 1 bucket. With the integer math approach, any weight ≥1 naturally gets proportional traffic.
- **No excess bucket handling**: With fixed bucket counts (100, 10,000, 100,000), minimum allocation adjustments could cause the total allocated buckets to exceed the bucket count, requiring complex logic to redistribute the excess. With integer math, allocations naturally sum to the total weight.
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
