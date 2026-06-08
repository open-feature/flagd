---
# Valid statuses: draft | proposed | rejected | accepted | superseded
status: draft
author: Todd Baert (@toddbaert)
created: 2026-06-08
updated: 2026-06-08
---
# Numeric coercion contract for typed flag accessors

flagd's evaluators don't currently agree on what should happen when a numeric flag is fetched via a different numeric accessor than its parsed JSON type, or when a value doesn't fit the requested type.
This ADR proposes a single contract for all flagd implementations to follow.
Tracking issue: [#1978](https://github.com/open-feature/flagd/issues/1978).

## Background

JSON does not distinguish integers from floats; `10` and `10.0` are the same kind of token.
flagd's wire formats (gRPC and OFREP) and its typed evaluator methods (`ResolveInt`, `ResolveFloat`, etc.) do distinguish them, and each language's flagd-core implementation has made different choices about what to do at that boundary.
The result is observable inconsistency.

Examples seen across implementations today:

| Variant      | Accessor | Go              | Java                                                                      | Python        | .NET          |
| ------------ | -------- | --------------- | ------------------------------------------------------------------------- | ------------- | ------------- |
| `3.14`       | Integer  | `3` (truncates) | `3` (truncates)                                                           | TYPE_MISMATCH | TYPE_MISMATCH |
| `9000000000` | Integer  | ok (int64)      | TYPE_MISMATCH (in-process); `410065408` (RPC, int32 overflow on the wire) | ok            | TYPE_MISMATCH |
| `9000000000` | Float    | ok              | TYPE_MISMATCH (in-process)                                                | ok            | ok            |
| `10.0`       | Integer  | `10`            | `10`                                                                      | TYPE_MISMATCH | TYPE_MISMATCH |

These cases are reachable with simple flag definitions and surface as silent data corruption, not just inconsistent error codes.

## Considered options

1. **Lossless coercion only, with a hard cap at the JSON-safe integer range (2^53 - 1)**: an accessor returns a value if and only if the conversion is lossless.
   Variants outside the safe-integer range are rejected up front so no transport ever has to silently lose precision.
2. **Lossless coercion only, no cap**: same accessor rule as option 1, but variants beyond 2^53 are accepted with documented precision loss in JS / OFREP / JSON paths.
3. **Strict per parsed type**: the variant's parsed JSON type (int vs float) is fixed; cross-type fetches always return `TYPE_MISMATCH`.
4. **Permissive coercion**: any numeric variant is returnable through any numeric accessor; truncation and overflow happen silently, matching today's Go and Java behavior.
5. **Status quo**: leave each implementation as it is.

I propose option 1.

Option 2 is close but leaves a silent-precision-loss case in place for JS clients and any JSON-based wire (OFREP), which contradicts the lossless principle the rest of the contract is built on.

Option 3 is the simplest specification but produces surprising behaviors.
JSON treats `10` and `10.0` interchangeably, so rejecting `10.0` from an integer accessor surprises users who weren't thinking about the JSON parser's typing decisions.
flagd's own JsonLogic engine treats numeric values uniformly during targeting, so strict typing only at the accessor boundary is internally inconsistent.

Option 4 preserves the silent-truncation and silent-overflow behaviors that prompted this ADR.
These are the cases we most need to fix.

Option 5 leaves the inconsistencies in place.

## Proposal

A numeric variant is returnable through a numeric accessor when the conversion is lossless.
Otherwise the evaluator returns `TYPE_MISMATCH`.

flagd additionally caps numeric flag values at the IEEE-754 safe-integer range, `[-(2^53 - 1), 2^53 - 1]`.
Variants whose absolute value exceeds this range are considered invalid per the JSON schema (we'll add this limit there).

| Variant kind                                                  | Fetched as Integer | Fetched as Float  |
| ------------------------------------------------------------- | ------------------ | ----------------- |
| int, fits the maximum of the numeric resolver in use          | value              | value (widened)   |
| int, exceeds the maximum of the numeric resolver in use, but within 2^53-1 | `TYPE_MISMATCH`    | value (widened)   |
| int, exceeds 2^53-1                                           | rejected at load (or `PARSE_ERROR` at evaluation if not validated) | rejected at load (or `PARSE_ERROR` at evaluation if not validated) |
| float, whole-valued and within the resolver's int range (e.g. `10.0`) | value (e.g. `10`)  | value     |
| float, fractional or out of the resolver's int range          | `TYPE_MISMATCH`    | value             |

The contract applies identically across all interfaces (gRPC, OFREP, and in-process evaluation).
"Numeric resolver in use" means the accessor the caller invoked: a 32-bit `Integer` accessor caps at `int32` max, a `Long` accessor (e.g. forthcoming Java `getLong`, .NET `Int64`) caps at `int64` max, a `Float` accessor caps at the safe-integer range.
Each language applies the rule against whichever resolver was called, not against a fixed type.

The 2^53-1 cap matches the interoperable integer range that every JSON parser, including JavaScript's, can faithfully represent.
Capping at this range removes the silent precision-loss case entirely; a value either round-trips exactly through every transport flagd supports, or it is rejected.
The alternative (permitting larger values and documenting precision loss past 2^53) preserves silent corruption in the JS and OFREP paths and contradicts the lossless principle this contract is built on.

## Consequences

The benefit is a single rule that all implementations follow, with no silent truncation, no silent overflow, and no silent precision loss.
The motivating bugs (`3.14` quietly becoming `3`, `9000000000` quietly becoming `410065408`) go away, and a value either round-trips exactly through every transport flagd supports or is rejected.

The cost is that this is a breaking change in two ways.
Go and Java users who currently rely on permissive coercion will see `TYPE_MISMATCH` where they previously got truncated values.
Operators with flag definitions containing values outside `[-(2^53 - 1), 2^53 - 1]` will see those flags fail validation; today such values either work (Go, Python) or fail unpredictably (Java, .NET).
Both changes are detectable; neither silently alters returned values.

A side effect is that the rule requires evaluators to inspect the value, not only the parsed type, when servicing a cross-type request.
This is a small cost; every implementation already has the value in hand at the point the type check occurs.

## Testing

Coverage lives in the [flagd testbed](https://github.com/open-feature/flagd-testbed) so every SDK and provider verifies the same contract.
We will add a new "numeric" suite capturing all these requirements.

## Versioning and migration

flagd is pre-1.0, so this ships as a minor-version bump with the breaking change called out in the release notes.
Provider releases follow, in tight coordination as usual.
