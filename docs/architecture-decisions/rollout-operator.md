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
    - The `rollback` operator also appends `rollout` (not `rollback`) so that bucket assignments match the original rollout; this is essential for the pivot-time gate to correctly identify which users were transitioned.

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

The `rollback` operator enables graceful reversal of a rollout, transitioning users back in **FILO order**: first adopters are last to revert, and users who never transitioned never see the new variant.

This requires a **pivot time** (the moment the rollback was initiated) which encodes how far the original rollout had progressed. The pivot time gives the operator enough "memory" to gate out never-transitioned users and reverse the rest in order, without storing any state. The rollback uses the same time window as the original rollout; the rollback completes at `endTime`.

```jsonc
// Rollback: same start/end as rollout, plus pivotTime
{"rollback": [1704067200, 1704068200, 1704067700, "new", "old"]}
```

Parameters:

- `bucketBy` (optional): Same as `rollout`.
- `startTime`: `startTime` from the original rollout.
- `endTime`: Controls when the rollback completes. Using the original rollout's `endTime` compresses the rollback into the remaining window; setting `endTime` to `pivotTime + (pivotTime - startTime)` rolls back at the same rate as the rollout progressed, etc.
- `pivotTime`: Unix timestamp when the rollback was initiated. Must be between `startTime` and `endTime`.
- `from`: The variant users are currently on (the rollout's `to`).
- `to`: The variant users revert to (the rollout's `from`).

**Implementation**:

```go
duration         := endTime - startTime
elapsedAtPivot   := pivotTime - startTime
rollbackDuration := endTime - pivotTime
bucket := (uint64(hashValue) * uint64(duration)) >> 32

// Gate: user never transitioned during rollout → always gets "to"
if bucket >= uint64(elapsedAtPivot) {
    return to
}

// Rollback progress
rollbackElapsed := currentTime - pivotTime
if rollbackElapsed <= 0 { return from }
if rollbackElapsed >= rollbackDuration { return to }

// Shrinking threshold: highest-bucket users (last adopted) revert first
remaining := rollbackDuration - rollbackElapsed
if bucket * uint64(rollbackDuration) < uint64(elapsedAtPivot) * uint64(remaining) {
    return from // still on rolled-out variant
}
return to // reverted
```

All operations are integer-only, consistent with the `rollout` operator.

**Example**: Rollout `[0, 1000, "old", "new"]` pivoted at t=500 (rollback completes at t=1000):

- **Alice** (adoption time t=200): adopted "new" at t=200. Reverts to "old" at t=800. First in, last out.
- **Bob** (adoption time t=400): adopted "new" at t=400. Reverts to "old" at t=600.
- **Carol** (adoption time t=600): would have adopted at t=600, but pivot was t=500. **Never sees "new".**
- **Fred** (adoption time t=700): would have adopted at t=700, but pivot was t=500. **Never sees "new".**

Without the pivot-time gate, Carol and Fred would temporarily be _exposed_ to "new" during rollback before being reverted, exactly the wrong behavior during an incident. The gate prevents this: any user whose bucket exceeds the elapsed time at pivot is immediately returned to "old" without ever seeing "new".

Nested operators (like `fractional`) are **not affected** — the rollback uses the same hash, so fractional bucket assignments remain stable.

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

### Alternative Proposal: Enhanced `fractional` with Dynamic Weights

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
This requires allowing the `fractional` weight argument to be a JSONLogic expression (currently it must be a hard-coded integer), as well as clamping negative weights to 0, in addition to support for non-string/nested variants ([#1877](https://github.com/open-feature/flagd/pull/1877)) and high-precision bucketing.

This approach is elegant and avoids a new operator. It achieves both forward rollout and FILO rollback using only existing JSONLogic primitives. The two approaches differ in the following ways:

1. **FILO rollback.** With `fractional`, naively swapping the weight expressions to reverse a rollout produces FIFO ordering: early adopters revert first, not last. FILO rollback _is_ achievable by reflecting time around the pivot point. Given a rollout over `[Ts, Te]` pivoted at `Tp`, define `R = 2Tp - Ts` and use:

    ```jsonc
    {
      "fractional": [
        { "var": "targetingKey" },
        ["new", { "-": [R, { "var": "$flagd.timestamp" }] }],
        ["old", { "-": [{ "+": [{ "var": "$flagd.timestamp" }, Te] }, 2Tp] }]
      ]
    }
    ```

    Where `R`, `Te`, `Ts`, and `Tp` are precomputed constants. Note that `R + Ts = (2Tp - Ts) + Ts = 2Tp`, so the `"old"` weight simplifies to `(t + Te) - 2Tp`. The `"new"` weight shrinks from `Tp - Ts` to 0, and the total weight is always `Te - Ts` (the rollout duration), naturally gating out never-transitioned users and reverting the rest in FILO order.
    The rollback completes at `t = R = 2Tp - Ts`, not at the original `Te`; it mirrors the rollout at the same rate, so the rollback takes as long to complete as the rollout had progressed. For example, if the rollout ran for 300s before pivoting, the rollback also takes 300s.

2. **Hash decorrelation.** The `rollout` operator automatically appends `"rollout"` (or some other salt) to the bucketing value before hashing, ensuring that a user's position in the rollout timeline does not correlate with their bucket in a nested `fractional`.
With the pure-fractional approach, the outer and inner `fractional` share the same hash, so early-rollout users systematically land in the first inner bucket. Users can work around this by manually adding a salt via `cat`, but that is non-obvious.

3. **Operator surface area.** The `fractional` approach requires no new operators; only that `fractional` accept JSONLogic expressions as weight arguments (currently hard-coded integers). The dedicated operators are more readable but add new definition surface area that must be implemented across all language SDKs.

#### Direct Comparison: Rollout with 50% Rollback

To make the tradeoffs clear, here is both approaches implementing the same scenario: a linear rollout from `"off"` to `"on"` over `[1740000000, 1800000000]`, followed by a "first in, last out" rollback initiated at the 50% mark (`pivotTime = 1770000000`).

**Using `rollout` / `rollback` operators:**

Rollout:

```jsonc
{"rollout": [1740000000, 1800000000, "off", "on"]}
```

Rollback (initiated at 50%):

```jsonc
{"rollback": [1740000000, 1800000000, 1770000000, "on", "off"]}
```

**Using `fractional` with dynamic weights:**

Rollout:

```jsonc
{
  "fractional": [
    ["on",  { "-": [{ "var": "$flagd.timestamp" }, 1740000000] }],
    ["off", { "-": [1800000000, { "var": "$flagd.timestamp" }] }]
  ]
}
```

Rollback (FILO, initiated at 50%; `R = 2 × 1770000000 − 1740000000 = 1800000000`, `2Tp = 2 × 1770000000 = 3540000000`):

```jsonc
{
  "fractional": [
    ["on",  { "-": [1800000000, { "var": "$flagd.timestamp" }] }],
    ["off", { "-": [{ "+": [{ "var": "$flagd.timestamp" }, 1800000000] }, 3540000000] }]
  ]
}
```

Note: the `"off"` weight simplifies to `t − 1740000000` (i.e., `t − Ts`), which mirrors the original rollout's `"on"` weight. The total weight is always `1800000000 − 1740000000 = 60000000` (the rollout duration), ensuring bucket assignments are preserved at the pivot instant.

**Summary:**

|                          | `rollout` / `rollback`                                                   | `fractional` with dynamic weights                                                             |
| ------------------------ | ------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------- |
| **Rollout definition**   | Single operator with explicit, non-nested parameters                     | Nested arithmetic expressions over `$flagd.timestamp`                                         |
| **Rollback definition**  | Single operator, same as `rollout`, but adds `pivotTime`                 | New set of weight expressions, requires precomputing `R = 2Tp − Ts`                           |
| **Readability**          | Intent is self-describing: time window, pivot, and direction are visible | User must reconstruct rollout semantics from arithmetic weight expressions                    |
| **Hash decorrelation**   | Automatic (appends `"rollout"` salt before hashing)                      | Manual, requires adding a salt via `cat` to avoid correlated bucketing in nested `fractional` |
| **New operator surface** | Yes, `rollout` and `rollback` must be implemented in all SDKs            | No, only requires `fractional` to accept JSONLogic weight expressions                         |
| **FILO correctness**     | Built-in via `pivotTime` parameter                                       | Achievable but non-obvious; naive weight-swap produces FIFO                                   |

### Consequences of Adding Rollout

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
