---
status: draft
author: @toddbaert
created: 2026-02-06
updated: 2026-02-06
---

# Rollout Operator

The rollout operator enables time-based progressive feature rollouts.

## Background

Progressive rollouts are a fundamental feature flag use case: gradually shifting traffic from one variant to another over time.
While stepped progression can be approximated in flagd by manually updating `fractional` weights on a schedule, or building a ruleset with multiple discrete timestamp checks each with different fractional distributions, true linear progression, where the percentage changes continuously over time, requires a time-aware operator.
The proposed rollout operator provides this.

The rollout operator complements (but does not make obsolete) the existing `fractional` operator by featuring a time dimension.
Where `fractional` distributes traffic across variants at a point in time, `rollout` transitions between any two variants over a time window, _including nested JSONLogic like `fractional` splits or conditional rules_.

## Requirements

- **Time-based**: traffic distribution must change automatically based on the current timestamp
- **Deterministic**: same user must get consistent results at a given point in time (no re-bucketing mid-request)
- **Composable**: must support nested JSONLogic (e.g., rollout to a `fractional` distribution)
- **Consistent hashing**: must use the same hashing strategy as `fractional` (MurmurHash3-32)
- **Cross-language portable**: must use only integer arithmetic (no floating-point operations)
- **JSONLogic conventions**: must follow established patterns for custom operators

## Proposal

### Operator Syntax

Three forms are supported, following JSONLogic array conventions:

```jsonc
// shorthand: roll from defaultVariant to "new"
{"rollout": [1704067200, 1706745600, "new"]}

// longhand: explicit from and to - from "old" to "new"
{"rollout": [1704067200, 1706745600, "old", "new"]}

// with custom bucketBy
{"rollout": [{"var": "email"}, 1704067200, 1706745600, "old", "new"]}
```

Parameters:

- `bucketBy` (optional): JSONLogic expression for bucketing value; defaults to `flagKey + targetingKey`, consistent with existing `fractional`
- `startTime`: Unix timestamp (seconds) when rollout begins (0% on `to`). Must be less than `endTime`.
- `endTime`: Unix timestamp (seconds) when rollout completes (100% on `to`). Must be greater than `startTime`.
- `from`: Starting variant or expression (omit for shorthand to use `defaultVariant`)
- `to`: Target variant or expression

**Timestamp validation**: To prevent accidental use of millisecond timestamps (which would schedule rollouts thousands of years in the future), the JSON Schema should enforce reasonable bounds on `startTime` and `endTime` (e.g., `minimum: 0`, `maximum: 3000000000`, approximately year 2065).

### Hashing Consistency

The rollout operator uses the same hashing strategy as `fractional` with one exception:

- MurmurHash3 (32-bit)
- Same default bucketing value: `flagKey + targetingKey`
- Same `bucketBy` expression support
- After the bucketing value is retrieved, before hashing, the UTF-8 byte representation of `rollout` is appended to the bucketing value
    - This injects entropy to ensure users in fractional rules nested within a `rollout` don't bucket identically (we want to ensure users early in a rollout don't always end up in the first fractional bucket).
    - The `rollback` operator also appends `rollout` (not `rollback`) to preserve correlation with rollout timing; this ensures the "first-in-last-out" property is maintained.

### Integer-Only Arithmetic

Per the [High-Precision Fractional Bucketing ADR](high-precision-fractional-bucketing.md), we avoid floating-point operations entirely.

Implementations must validate that `endTime > startTime` (strict inequality) at parse time, rejecting the configuration otherwise.
This also rejects `startTime == endTime` (duration = 0), which would be a degenerate case; an "instant rollout" is better expressed as a direct variant assignment.
Additionally, `elapsed` must be clamped to `[0, duration]` to prevent overflow from negative values or times beyond the window:

```go
duration := endTime - startTime // validated > 0 at parse time
elapsed  := currentTime - startTime

// before startTime: everyone gets "from"; after endTime → everyone gets "to"
if elapsed <= 0 {
    return from
}
if elapsed >= duration {
    return to
}

// Maps hash to [0, duration) range using integer math only
bucket := (uint64(hashValue) * uint64(duration)) >> 32

if bucket < uint64(elapsed) {
    return to
}
return from
```

This is mathematically equivalent to `(hash/2^32) < (elapsed/duration)` but uses only:

- 64-bit multiplication
- 32-bit right shift
- Integer comparison

These operations are portable across all languages without floating-point precision concerns.

### Nested JSONLogic Support

Variants can be JSONLogic expressions, enabling composition:

```jsonc
// Rollout to a fractional split
{
  "rollout": [
    1704067200, 1706745600,
    "old",
    {"fractional": [["a", 50], ["b", 50]]}
  ]
}

// Conditional logic within rollout
{
  "rollout": [
    1704067200, 1706745600,
    "old",
    {"if": [{"==": [{"var": "tier"}, "premium"]}, "premium-new", "basic-new"]}
  ]
}
```

### Rollback Operator

The `rollback` operator enables graceful reversal of a rollout, transitioning users back in **reverse order**: first adopters are last to revert, and users who've never been assigned new functionality never see it.

```jsonc
// Rollback: same parameters as rollout
{"rollback": [1704067200, 1706745600, "old", "new"]}
```

**Implementation**: `rollback` inverts the hash using bitwise NOT (`~`) before bucketing and swaps the return values:

```go
bucket := (uint64(^hashValue) * uint64(duration)) >> 32

if bucket < uint64(elapsed) {
    return from
}
return to
```

This inverts user ordering:

- Hash `0x00000000` (first in rollout) → `0xFFFFFFFF` (last in rollback)
- Hash `0xFFFFFFFF` (last in rollout) → `0x00000000` (first in rollback)

**Example**: Given a rollout to a fractional split:

```jsonc
{"rollout": [1704067200, 1706745600, "old", {"fractional": [["a", 50], ["b", 50]]}]}
```

- **Alice** (hash `0x66666666`, 40% position): transitions to "a" at t=40%
- **Fred** (hash `0xFFFFFFFF`, 100% position): never reaches the fractional

If switched to rollback mid-way:

```jsonc
{"rollback": [1704067200, 1706745600, "old", {"fractional": [["a", 50], ["b", 50]]}]}
```

- **Alice** (inverted: `0x99999999`, 60% position): reverts to "old" at t=60%
- **Fred**: already "reverted" immediately (was never transitioned)

Users revert in the exact reverse order they adopted. Nested operators (like `fractional`) are **not affected** by the hash inversion — only the rollback timing decision uses the inverted hash, preserving stable bucket assignments within the `to` expression.

### Future-proofing

Later, we may want to support additional non-linear rollouts.
This can be done with an additional, optional, configuration parameter before the times params (similar to custom bucketing).

```jsonc
{"rollout": [{"var": "email"}, "linear|exponential", 1704067200, 1706745600, "old", "new"]}
```

```jsonc
{"rollout": [{"var": "email"}, { some-json-logic-lambda }, 1704067200, 1706745600, "old", "new"]}
```

**Implementation of non-linear rollouts is out of the scope of this proposal.**

### Considered Alternative: Enhanced `fractional` with Dynamic Weights

An alternative to a dedicated operator was proposed: use `fractional` with JSONLogic expressions as weights, combined with `$flagd.timestamp`, to achieve time-based progression without any new operator:

```jsonc
{
  "fractional": [
    { "var": "targetingKey" },
    ["on",  { "-": [{ "var": "$flagd.timestamp" }, 1740000000] }],
    ["off", { "-": [1800000000, { "var": "$flagd.timestamp" }] }]
  ]
}
```

As time advances, the weight of `"on"` grows and `"off"` shrinks, producing a progressive rollout using only existing primitives.
This requires allowing the `fractional` weight argument to be a JSONLogic expression (currently it must be a hard-coded integer), in addition to support for non-string/nested variants ([#1877](https://github.com/open-feature/flagd/pull/1877)) and high-precision bucketing.

This approach is elegant and avoids a new operator, but was not adopted for two reasons:

1. **No FILO rollback.** With `fractional`, swapping the weight expressions to reverse a rollout produces FIFO ordering: early adopters revert first, not last. Worse, if weights are swapped mid-rollout, the formula starts at 100% "new" and ramps down, meaning users who _never saw the new variant_ are suddenly exposed to it before being reverted. This is exactly the wrong behavior during an incident. The `rollback` operator avoids this via hash inversion (`~hash`), which is not expressible in JSONLogic.

2. **No automatic hash decorrelation.** The `rollout` operator appends `"rollout"` (or some other salt) to the bucketing value before hashing, ensuring that a user's position in the rollout timeline does not correlate with their bucket in a nested `fractional`.
With the pure-fractional approach, the outer and inner `fractional` share the same hash, so early-rollout users systematically land in the first inner bucket. Users can work around this by manually adding a salt via `cat`, but this is error-prone and non-obvious.

### Consequences

- Good, because this enables functionality present in many other systems
- Good, because time-based rollouts are declarative and require no external automation
- Good, because hashing is consistent with `fractional`
- Good, because integer-only math ensures cross-language portability
- Good, because nested JSONLogic enables complex rollout scenarios
- Good, because timestamp usage, array parameter style, and shorthand are consistent with other operators
- Good, because `rollback` enables graceful reversal without subjecting users to unnecessary thrashing.
- Bad, because it's more definition surface area to understand
- Bad, because additional timed mechanisms may represent changes in behavior ("time-bombs") that can be difficult to trace
- Bad, because consistently testing a time-sensitive operator might be somewhat challenging
