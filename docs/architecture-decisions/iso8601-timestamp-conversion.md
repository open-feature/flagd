---
# Valid statuses: draft | proposed | rejected | accepted | superseded
status: proposed
author: Srinivas Chippagiri
created: 2026-10-28
updated: 2026-10-28
---

# ISO 8601 Timestamp Conversion Operator

This ADR proposes a `to_timestamp` custom operation that converts ISO 8601 dates and date-times into Unix timestamps (seconds), so that time-based targeting rules can be written with human-readable times while still using the existing numeric comparison operators and `$flagd.timestamp`.

## Background

Time-based targeting in flagd works on Unix timestamps: `$flagd.timestamp` is the evaluation time in epoch seconds, and rules compare it with numeric literals, for example `{">": [{"var": "$flagd.timestamp"}, 1775674200]}`.
Numbers like `1775674200` are hard to read and review, and easy to get wrong (milliseconds instead of seconds, wrong timezone when converting by hand).
Users have asked to express these times as ISO 8601 strings instead ([#1927](https://github.com/open-feature/flagd/issues/1927)).

Two approaches were discussed on the issue: a conversion operator that turns an ISO 8601 string into a Unix timestamp, and an ISO 8601 "sub-type" with dedicated comparison operators (`is_before`, `is_after`, ...).
There was consensus on the conversion operator, as it reuses every existing numeric operator, and adds a single operator to implement across the flagd evaluation engines in all supported languages.
This ADR records that decision and specifies the behavior precisely enough for consistent cross-language implementations.

## Requirements

- **Readable**: times in targeting rules can be written as ISO 8601 dates and date-times.
- **Composable**: the result works with all existing numeric JSONLogic operators and with `$flagd.timestamp`.
- **Deterministic**: the same rule and input produce the same timestamp in every flagd implementation (flagd, the in-process providers and the language SDK evaluators), regardless of the host the evaluation runs on.
- **Cross-language portable**: the accepted input format can be parsed identically in Go, Java, JavaScript, Python and .NET without relying on parser behavior that differs between runtimes.
- **JSONLogic conventions**: invalid input yields `null` (as other custom operations do) rather than failing the evaluation.

## Considered Options

1. **Conversion operator** (`to_timestamp`): converts an ISO 8601 string to a Unix timestamp; rules use existing operators.
2. **ISO 8601 sub-type** with dedicated operators (`is_before`, `is_after`, ...), similar to `sem_ver`.

Option 2 has high semantic clarity, but requires a family of new operators and schema changes, makes the existing timestamp-based rules second-class, and has no natural limit on the operators it invites (`is_weekday`, `is_leap_year`, ...).
**Option 1 is proposed.**

## Proposal

### Operator Syntax

```jsonc
// literal
{ "to_timestamp": "2026-04-08T18:50:00Z" }

// any JSONLogic expression that yields a string, e.g. a context attribute
{ "to_timestamp": { "var": "subscription_end" } }

// array form, as with other operations
{ "to_timestamp": ["2026-04-08"] }
```

The result is an integer number of seconds since the Unix epoch, the same unit as `$flagd.timestamp`, so rules compose directly:

```jsonc
{
  "if": [
    { ">=": [{ "var": "$flagd.timestamp" }, { "to_timestamp": "2026-04-08T18:50:00Z" }] },
    "on",
    "off"
  ]
}
```

### Accepted Formats

To guarantee identical results across runtimes, `to_timestamp` accepts a defined profile of ISO 8601, aligned with RFC 3339, rather than "whatever the platform parser accepts":

| Form                       | Example                     | Meaning                                                          |
| -------------------------- | --------------------------- | ---------------------------------------------------------------- |
| Date                       | `2026-04-08`                | start of that day (`00:00:00`) in the evaluation timezone        |
| Date-time with offset      | `2026-04-08T18:50:00+02:00` | that instant                                                     |
| Date-time in UTC           | `2026-04-08T18:50:00Z`      | that instant                                                     |
| Date-time without offset   | `2026-04-08T18:50:00`       | that local time in the evaluation timezone                       |
| Optional seconds/fractions | `2026-04-08T18:50`, `...:00.250Z` | seconds default to `0`; fractional seconds are truncated  |

Rules:

- Separators are required (extended format): `YYYY-MM-DD`, `T`, `HH:MM[:SS[.fraction]]`, and an offset of `Z` or `±HH:MM`. `T` and `Z` may be lowercase, as RFC 3339 allows.
- Week dates (`2026-W15-3`), ordinal dates (`2026-098`), basic format without separators, durations, intervals and reduced precision other than the forms above are **not** accepted.
- Fractional seconds are truncated (not rounded), so the result is always an integer.
- A number (integer) argument is returned unchanged, so rules can accept either form from context.
- Anything else (invalid format, out-of-range fields such as month `13`, a leap second `23:59:60`, non-string non-number values) returns `null` and logs a warning, following the convention of the other custom operations. A `null` makes numeric comparisons false, so a malformed time never enables a flag by accident.

### Evaluation Timezone

Inputs with an explicit offset (`Z`, `±HH:MM`) are unambiguous.
Inputs without one (dates and date-times without offset) are interpreted in the **evaluation timezone**:

- The default is **UTC**.
- It can be overridden in configuration (for flagd, a startup flag/environment variable such as `--timezone` / `FLAGD_TIMEZONE`; for in-process providers, the equivalent provider option), using an IANA zone name such as `America/Chicago`.

The issue discussion suggested defaulting to the host OS timezone.
This ADR proposes UTC instead, because flag rules are evaluated in many places (the flagd server and in-process providers in each application) whose hosts may run in different timezones; an OS default would make the same rule switch at different instants depending on where it is evaluated, which breaks the **Deterministic** requirement.
With UTC as the default, the result only changes when an operator opts in to a timezone explicitly. See [Open Questions](#open-questions).

### Daylight Saving Time

For local times (no offset) in a timezone with daylight saving time:

- **Ambiguous times** (repeated during a backward transition): the **earlier** instant is used.
- **Non-existent times** (skipped during a forward transition): the time is shifted **forward** by the length of the gap (e.g. `02:30` on a day that skips `02:00`–`03:00` becomes `03:30`).

These match the defaults of `java.time.ZonedDateTime.of`.
Other runtimes need to apply them explicitly: Go's `time.Date`, for example, documents that it does not guarantee which instant it returns in a transition, so the Go implementation has to detect both cases and resolve them as above.

### Cross-Language Implementation Notes

The accepted profile is deliberately narrow so that each implementation can use a strict parser:

- **Go**: a regular expression for the profile, then `time.Date` in the evaluation `*time.Location` (or the parsed fixed offset). Avoid `time.Parse` layouts with optional parts.
- **Java**: `java.time` (`OffsetDateTime`, `LocalDateTime`, `LocalDate`), after validating the profile with a regular expression, as `DateTimeFormatter.ISO_*` accepts more than the profile.
- **JavaScript**: a regular expression plus `Date.UTC` and an offset computed from the evaluation timezone with `Intl.DateTimeFormat`; do **not** use `Date.parse`, whose handling of date-only and offset-less strings differs from the rules above.
- **Python**: a regular expression plus `datetime` and `zoneinfo`; do not rely on `datetime.fromisoformat`, which accepts different inputs across Python versions.
- **.NET**: a regular expression plus `DateTimeOffset` and `TimeZoneInfo`, rather than culture-sensitive parsing.

### API Changes

- `to_timestamp` is added to the targeting JSON schema in `flagd-schemas` as a custom operation taking a single string, number or JSONLogic expression (or a one-element array of these).
- flagd and the in-process providers get the evaluation timezone option described above.
- Documentation: a new page under `docs/reference/custom-operations/`, and a cross-reference from the `$flagd.timestamp` documentation.

### Test Plan

Shared test cases in the flagd testbed (Gherkin), so every implementation is checked against the same expectations:

- each accepted form, with `Z`, positive and negative offsets, and fractional seconds;
- date-only and offset-less inputs with the default (UTC) and an overridden timezone;
- an ambiguous and a non-existent local time across a DST transition;
- invalid inputs: wrong separators, week and ordinal dates, month `13`, `23:59:60`, empty string, non-string values; all yield `null`;
- a number argument returned unchanged;
- composition with `$flagd.timestamp` and numeric operators.

### Consequences

- Positive: readable, reviewable time-based rules without a new family of operators; existing timestamp rules keep working.
- Positive: one operator to implement and test across all flagd evaluation engines.
- Negative: rules need an explicit conversion, which is more verbose than a dedicated `is_after`.
- Negative: implementations must carry a small, strict parser instead of using the platform's ISO 8601 parser.

## Open Questions

1. Should the default evaluation timezone be UTC (proposed, for determinism across evaluation sites) or the host OS timezone (as initially suggested on the issue)?
2. Is `to_timestamp` the preferred name, or should it be `to_unix` to make the unit explicit?
3. Should durations (e.g. adding `PT2H` to a time, raised on the issue) be a follow-up operator, or stay out of scope?
