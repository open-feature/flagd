---
# Valid statuses: draft | proposed | rejected | accepted | superseded
status: accepted
author: Parth Suthar
created: 2026-07-30
updated: 2026-11-26
---

# Per-evaluation metadata from targeting rules

Let a targeting rule terminate with `{ "resolve": ["<key>", "<reason>", [["<metadata-key>", <value>], ...]] }` in place of a plain variant string, so the _branch that fired_ can annotate the evaluation with extra metadata, custom reasons, and other auxiliary information. Existing string returns are unchanged.

## Background

An evaluation today returns `value`, `variant`, `reason`, and `metadata`. The `reason` is a coarse enum (`TARGETING_MATCH`, `DEFAULT`, `STATIC`, …); `metadata` carries only the static blocks defined at flag-set and flag level. Neither answers the debugging question we hit most often: **which branch of the targeting expression fired?**

For a nested `if` / `and` / `or` tree the resolver just returns the winning variant key. Two rules landing on the same variant are indistinguishable in the response — the only way to know _why_ today is to fetch the flag config and re-run the logic by hand. Encoding the branch identity into the variant key (`"clubs-eu-rollout-a"`) is the workaround, and it pollutes the variant space with debug info consumers then have to parse back out.

The plumbing to carry metadata already exists on most paths: `AnyValue.Metadata` (`core/pkg/evaluator/ievaluator.go`) is threaded through the evaluator and surfaced via OFREP (single and bulk) and single-resolution gRPC responses as `flagMetadata`. Two paths don't carry it yet — bulk gRPC `ResolveAll` and `RecordEvaluation` telemetry — but wiring those up is separate follow-up work.

## Proposal

Add a `resolve` operator to `schemas/json/targeting.json`, alongside `fractional`, `starts_with`, `ends_with`, and `sem_ver`. Plain string returns are untouched:

```json
{
  "resolve": {
    "description": "Terminates targeting with a variant, and optionally a custom reason and metadata for the branch that fired.",
    "type": "array",
    "minItems": 1,
    "maxItems": 3,
    "items": [
      {
        "description": "Variant: string, boolean coerced to \"true\"/\"false\", or null to fall back like a plain null result. Literal numbers are rejected. A computed rule is also accepted; it is coerced like a plain targeting result.",
        "oneOf": [
          { "type": "string" },
          { "type": "boolean" },
          { "type": "null" },
          { "$ref": "#/definitions/anyRule" }
        ]
      },
      {
        "description": "Reason: free-form string, null, or a computed rule. A non-null reason replaces the reason the equivalent plain return would report.",
        "oneOf": [
          { "type": ["string", "null"] },
          { "$ref": "#/definitions/anyRule" }
        ]
      },
      {
        "description": "Metadata: [key, value] pairs. The list, each pair, each key, and each value may also be a computed rule. flagSetId and version are reserved for flag-set metadata.",
        "oneOf": [
          { "$ref": "#/definitions/anyRule" },
          {
            "type": "array",
            "items": {
              "oneOf": [
                { "$ref": "#/definitions/anyRule" },
                {
                  "type": "array",
                  "minItems": 2,
                  "maxItems": 2,
                  "items": [
                    {
                      "oneOf": [
                        { "type": "string", "not": { "enum": ["flagSetId", "version"] } },
                        { "$ref": "#/definitions/anyRule" }
                      ]
                    },
                    {
                      "oneOf": [
                        { "type": ["string", "number", "boolean"] },
                        { "$ref": "#/definitions/anyRule" }
                      ]
                    }
                  ]
                }
              ]
            }
          }
        ]
      }
    ]
  }
}
```

### Arguments

- **`variant`** (1st): a literal string; a literal boolean, coerced to `"true"` or `"false"`; `null`; or a computed expression. The only **required** argument. A string is the exact equivalent of the plain-string return: `"red"` and `{ "resolve": ["red"] }` are interchangeable terminals selecting the variant `red`. Literal numbers are rejected. Computed variants are coerced like plain targeting results: booleans become `"true"`/`"false"`, and the behavior of numbers is not defined.
  A `null` variant falls back exactly like a plain `null` targeting result: to `defaultVariant` if the flag has one, otherwise to the caller's code default (no value or variant). Unlike a plain `null`, it can carry a custom reason and metadata.
- **`reason`** (2nd): string, `null`, or a computed expression, optional. Free-form and operator-defined; overrides the top-level evaluation `reason` for this branch (see [Narrowing the coarse `reason`](#narrowing-the-coarse-reason)).
  A non-null `reason` replaces the reason the equivalent plain return would report; a `null` or omitted one keeps it. Standard reasons such as `DISABLED` or `ERROR` are allowed, and are passed through like any other string (see [Standard reasons](#standard-reasons)).
- **`metadata`** (3rd): array of `[key, value]` pairs, optional, shallow. The whole list, each pair, each key, and each value may also be a computed expression. Literal values are `string | number | boolean`, matching flag and flag-set metadata. `flagSetId` and `version` are reserved for flag-set metadata and cannot be overridden by rules; the schema rejects them as literal keys, and implementations skip computed keys that match.
- **`details`**: a well-known metadata key: free-text explanation of why this branch fired.

`resolve` is **terminal**: it is only valid as the whole targeting rule; the result branches of `if` (recursively); and `fractional` bucket variants, in both the shorthand and the bucketing-value form. Every bucket variant is evaluated, and the selected one is returned.

The `resolve` operation is validated the same way as others: via JSON schema at load time only. The schema enforces terminality through separate `terminalIfRule` and `terminalFractionalRule` definitions reachable only from terminal positions, and by restricting the first `if` argument (always a condition) to plain args.
What the schema cannot enforce: `resolve` as a later else-if condition (positions 2, 4, … in a chained `if`), and `resolve` reached through `$ref` shared evaluators. Those cases are undefined behavior.

### Evaluator changes

In `evaluateVariant` (`core/pkg/evaluator/json.go`), the current unconditional quote-strip of the JsonLogic result gains a `resolve` case ahead of it:

- `resolve` result: evaluate exactly as the equivalent plain return of its `variant` would (a string selects that variant, with booleans coerced to `"true"`/`"false"` officially and numbers coerced like plain terminals, behavior not defined; `null` falls back to `defaultVariant`, or the caller's code default if there is none), on every surface (gRPC v1 and v2, OFREP, in-process; single and bulk).
  Then shallow-merge its `metadata` pairs into the returned metadata, and if `reason` is non-null, use it in place of the reason that return would report (`TARGETING_MATCH` for a variant, `DEFAULT` for `null`)
- everything else → unchanged from today: `null` exits targeting to `defaultVariant`; strings, booleans, and numbers are quote-stripped into a variant key (so a boolean `true` still resolves the `"true"` variant, as `schemas/json/targeting.json` documents and the existing tests cover); a key absent from `variants` errors exactly as it does now

`reason` never drives evaluation logic: whether a value is returned depends only on the variant, and whether an evaluation failed depends only on the presence of an `errorCode`. A rule-returned `ERROR` without an `errorCode` is a successful evaluation. Reasons are auxiliary data.

Metadata merge precedence, lowest > highest, so more specific wins: flag-set metadata < flag metadata < **rule-returned metadata**. The reserved keys (`flagSetId`, `version`) are the exception: the schema rejects them as literal rule metadata keys, so they come from the existing flag evaluation flow.

!!! warning "Forward compatibility"

    Existing configs keep working on a new binary, but the reverse does not hold: a config using `resolve`, evaluated by an older flagd binary, fails with an unknown-operator error. flagd has to be upgraded fleet-wide before any config starts using `resolve`. This needs an explicit migration note in the release.

    flagd-schemas is pre-1.0, so `resolve` merges into the `v0` schema without a version bump; flagd core decides which schema version it uses by pinning the embedded snapshot it ships with.

    **Hold the hosted schema back.** `https://flagd.dev/schema/v0/flags.json` is what editors and IDE tooling resolve against. Do not publish `resolve` there until it has landed in flagd and a few providers — otherwise editors will happily autocomplete a construct most of the fleet cannot evaluate.

Why an operator with array arguments rather than a returned object: JsonLogic engines disagree on objects. Some return a multi-key object as data (Go, Python, .NET), others reject it (JS, Rust, and Java outside `if` results), and C++ applies only its first key. A single-key object as an operator call, and arrays of primitives as arguments, behave the same in every engine flagd implementations use.
Metadata is an array of pairs for the same reason: an object argument fails to parse in Java, JS, and Rust.

### Example — a full flag config

```jsonc
{
  "$schema": "https://flagd.dev/schema/v0/flags.json",
  "flags": {
    "enable-mainframe-access": {
      "state": "ENABLED",
      "defaultVariant": "off",
      "variants": {
        "on": true,
        "off": false,
      },
      "targeting": {
        "if": [
          {
            "ends_with": [{ "var": "email" }, "@ingen.com"],
          },
          // new terminus — this node used to be the plain string "on"
          {
            "resolve": [
              "on", // required
              "my_super_cool_custom_reason", // optional
              [], // custom metadata pairs, optionally including "details"
            ],
          },
        ],
      },
    },
  },
}
```

### Example — chained `if`

Each branch tags the evaluation with the rule that fired:

```json
{
  "acceptable-feature-stability": {
    "state": "ENABLED",
    "defaultVariant": "ga",
    "variants": { "alpha": "alpha", "beta": "beta", "ga": "ga" },
    "targeting": {
      "if": [
        { "===": [{ "var": "customerId" }, "customer-A"] },
        {
          "resolve": [
            "alpha",
            "TARGETING_MATCH",
            [["details", "explicit allowlist for enterprise pilot"]]
          ]
        },
        { "in": [{ "var": "customerId" }, ["customer-B1", "customer-B2"]] },
        { "resolve": ["beta", "TARGETING_MATCH"] },
        { "resolve": ["ga", "FALLTHROUGH"] }
      ]
    }
  }
}
```

### Example — `if` with a fractional split

`resolve` is valid as a `fractional` bucket variant, so each bucket can carry its own reason and metadata directly:

```json
{
  "targeting": {
    "if": [
      { "==": [{ "var": "locale" }, "en-US"] },
      {
        "fractional": [
          { "var": "targetingKey" },
          [{ "resolve": ["on", "SPLIT_RANDOM", [["details", "us-10pct-rollout"]]] }, 10],
          [{ "resolve": ["off", "SPLIT_FALLTHROUGH", [["details", "us-holdback"]]] }, 90]
        ]
      },
      { "resolve": ["off", "FALLTHROUGH", [["details", "non-us-off"]]] }
    ]
  }
}
```

The 90% US holdback and the non-US off case both serve `off` but carry different `reason` and `details` values, so they stay distinguishable in telemetry.

### Example: disabled with overrides

A management system disables a flag for a stage, but some users self-target it with overrides. Rather than emitting `state: DISABLED` (which would also disable the overrides), it emits the flag as `ENABLED`, with no `defaultVariant`, and:

```json
{
  "state": "ENABLED",
  "variants": { "on": true, "off": false },
  "targeting": {
    "if": [
      { "in": [{ "var": "targetingKey" }, ["alice"]] },
      "on",
      { "resolve": [null, "DISABLED"] }
    ]
  }
}
```

Override users get `on`; everyone else gets the caller's code default with `reason=DISABLED`, exactly as if the flag were disabled. Adding or removing an override doesn't change what non-targeted users see.

### Narrowing the coarse `reason`

flagd's top-level `reason` today (`TARGETING_MATCH`, `DEFAULT`, `STATIC`, `DISABLED`, `ERROR`) says _what class_ of evaluation happened but not _which specific path_ inside it — every rule-driven result collapses to `TARGETING_MATCH`. A rule-scoped `reason` plus `metadata.details` sub-classifies that. The schema deliberately does not fix an enum: `reason` is an arbitrary string the operator defines, so each fleet picks whatever taxonomy fits it.

Note that the OpenFeature spec ([evaluation details, requirement 6.1](https://openfeature.dev/specification/sections/evaluation-context)) types `reason` as a free-form string — the values in the spec (`TARGETING_MATCH`, `SPLIT`, `DEFAULT`, …) are _recommended_, not exhaustive, and providers are explicitly allowed to emit their own.
That means letting a branch-scoped `reason` override the top-level `reason` is spec-legal without any SDK contract change; SDK hooks and telemetry sinks already treat the field as opaque.

The user-facing reason does not match what the flagd engine actually evaluated for different variants of `SPLIT`. Today `SPLIT` covers both fractional and gradual rollout. Letting the rule override the reason is an easy win here, and the same applies to the several distinct kinds of `OVERRIDE`.

```json
{ "resolve": ["on", "SPLIT_GRADUAL", [["details", "5%→50% ramp, week 3"]]] }
{ "resolve": ["on", "SPLIT_STEPPED", [["details", "stage 2 of 4, cohort=eu-west"]]] }
{ "resolve": ["on", "SPLIT_RANDOM", [["details", "50/50 A/B, salt=exp-4231"]]] }
```

Seeing these reasons instead of a bare `SPLIT` allows for much richer telemetry signals.

#### Standard reasons

A rule may return any reason, including those defined by OpenFeature or flagd (`DISABLED`, `DEFAULT`, `STATIC`, `ERROR`, `CACHED`, etc.). They're passed through as plain strings: flagd attaches no behavior to a rule-returned reason, so the evaluation is otherwise exactly what the `variant` makes it.

It's up to the rule author to use standard reasons sensibly. For example, a rule-returned `ERROR` comes with no `errorCode`, so consumers that expect one may misbehave (see [Consequences](#consequences)).

Other patterns the free-form field unlocks inside a single `TARGETING_MATCH`:

- **Split vs. allowlist under the same top-level reason** — two paths reach the same variant, telemetry can tell them apart:

  ```json
  { "resolve": ["on", "SPLIT", [["details", "10% cohort bucketed by email"]]] }
  { "resolve": ["on", "ALLOWLIST", [["details", "customerId in enterprise-pilot"]]] }
  ```

- **Kill switch / override tagged distinctly** — a branch short-circuiting an incident stays visible in traces:

  ```json
  { "resolve": ["off", "OVERRIDE", [["details", "incident-4231 kill switch"]]] }
  ```

- **Audience-segment attribution** — which specific condition inside a large predicate fired:

  ```json
  {
    "resolve": [
      "beta",
      "AUDIENCE_MATCH",
      [["details", "country=US AND app_version>=2.0"]]
    ]
  }
  ```

### Testing

Beyond Go unit tests in `core/pkg/evaluator`, the flagd gherkin suite gets extended to cover the new terminus end-to-end: testbed flag configs using `resolve`, and scenarios in the shared feature files asserting the resolved `variant`, the overridden `reason`, and the merged `flagMetadata` (including `details`). These include:

- the reason rule: with no reason, `resolve` reports exactly what the equivalent plain return reports (`TARGETING_MATCH` for a variant, `DEFAULT` for `null`); with a reason, that reason replaces it
- a `null` variant, with and without a `defaultVariant`: the same fallback as a plain `null`
- a `null` variant with a custom reason (e.g. `OFF_BY_POLICY`) on a flag without a `defaultVariant`: no value or variant (the caller's code default), and the custom reason, across RPC, OFREP and in-process
- a rule-returned `DISABLED`: with a `null` variant, matching a disabled flag, in both single and bulk evaluation; with a selected variant, that variant is served
- a rule-returned `ERROR`: passed through as the reason, with no error code, and the variant served

Those scenarios run across the gRPC, OFREP, and in-process paths that `test/integration` already exercises, which keeps every provider honest about the new shape rather than testing it only inside flagd.

## Consequences

- **Good** — direct answer to "which branch fired" without overloading the variant key; no wire-format changes. `flagMetadata` is already emitted on OpenFeature evaluation events and picked up by SDK hooks (OpenTelemetry, logging, custom exporters), so anything reading it — OTel spans, A/B dashboards, debug logs — picks up rule-level attribution automatically.
- **Good** — backwards compatible for existing configs: string, boolean, number, and null returns are byte-for-byte unchanged on an upgraded flagd binary.
- **Good** — no OpenFeature spec change required. `reason` is already a first-class spec field, and `details` rides `flagMetadata`, so no SDK gains a new resolution field.
- **Good**: a `null` variant lets a rule fall back like a plain `null` while still saying why, so "disabled except for overrides" can be expressed without disabling the overrides.
- **Bad**: rules can claim any standard reason, including `DISABLED` or `ERROR`, and hooks and telemetry sinks that branch on them will take it at face value. Operators own the truthfulness of their rules' reasons; flagd only guarantees that a rule-returned reason never changes the served variant.
- **Bad** — not forward compatible: a config using `resolve`, evaluated by an older flagd binary, fails with an unknown-operator error. Rollouts must upgrade flagd before configs start using `resolve`. Needs an explicit migration note in the release. No schema version bump is needed — flagd-schemas is pre-1.0 and `resolve` merges into the existing `v0` schema; flagd core pins the snapshot it embeds.
- **Bad** — automatic telemetry attribution requires the separate `RecordEvaluation` work called out in Background; without it, rule ids ride `flagMetadata` for SDK-side hooks but don't reach flagd's own OTel metrics.

## Open questions

None outstanding. Making `variant` optional (`{ "resolve": [] }`) was considered and rejected: an explicit `null` mirrors the plain `null` return, so the fallback is visible in the rule. Schema version bump and fractional bucket support were both resolved during schema review (flagd-schemas#238, flagd#2127): no version bump is needed (pre-1.0), and `resolve` is now valid as a `fractional` bucket variant.
