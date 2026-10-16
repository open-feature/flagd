---
# Valid statuses: draft | proposed | rejected | accepted | superseded
status: accepted
author: Parth Suthar
created: 2026-07-30
updated: 2026-10-16
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
    "prefixItems": [
      { "type": "string" },
      { "type": ["string", "null"] },
      {
        "type": "array",
        "items": {
          "type": "array",
          "prefixItems": [
            { "type": "string" },
            { "type": ["string", "number", "boolean"] }
          ],
          "minItems": 2,
          "maxItems": 2
        }
      }
    ]
  }
}
```

### Arguments

- **`variant`** (1st): string. The only **required** argument, and the exact equivalent of the plain-string return: `"red"` and `{ "resolve": ["red"] }` are interchangeable terminals selecting the variant `red`.
- **`reason`** (2nd): string or `null`, optional. Free-form, operator-defined; overrides the top-level evaluation `reason` for this branch. The schema does not constrain it to an enum (see [Narrowing the coarse `reason`](#narrowing-the-coarse-reason)).
- **`metadata`** (3rd): array of `[key, value]` pairs, optional, shallow. Values are restricted to `string | number | boolean`, matching flag and flag-set metadata.
- **`details`**: a well-known metadata key: free-text explanation of why this branch fired.

`resolve` is **terminal**: it is only valid as the whole targeting rule, or as a result branch of `if` / `?:` (recursively). Implementations reject it anywhere else when the flag is loaded, so its result always passes through to the resolver unchanged.

### Evaluator changes

In `evaluateVariant` (`core/pkg/evaluator/json.go`), the current unconditional quote-strip of the JsonLogic result gains a `resolve` case ahead of it:

- `resolve` result: use its `variant` as the variant key; shallow-merge its `metadata` pairs into the returned metadata; if `reason` is present, use it in place of `TARGETING_MATCH` (a `null` or omitted `reason` keeps `TARGETING_MATCH`)
- everything else → unchanged from today: `null` exits targeting to `defaultVariant`; strings, booleans, and numbers are quote-stripped into a variant key (so a boolean `true` still resolves the `"true"` variant, as `schemas/json/targeting.json` documents and the existing tests cover); a key absent from `variants` errors exactly as it does now

Metadata merge precedence, lowest → highest, so more specific wins: flag-set metadata → flag metadata → **rule-returned metadata**.

!!! warning "Breaking change to the config schema"

    Existing configs keep working on a new binary, but the reverse does not hold: a config using `resolve`, evaluated by an older flagd binary, fails with an unknown-operator error.

    There's no schema/feature version negotiation today, so the config schema version has to gate it: operators on an older binary must stay on the older schema version, and flagd has to be upgraded fleet-wide before any config starts using `resolve`. This needs a schema version bump and a migration note in the release.

    Gating on the schema version is sufficient on its own — providers and flagd embed a snapshot of the schema in each version, so an older release rejects `resolve` at validation time rather than misbehaving at evaluation time.

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

### Example — nested `if` with a fractional split

`fractional` returns a plain string (`"on"` or `"off"`), and both are truthy under JsonLogic — so it can't be used directly as an `if` condition; wrapping it in `==` against a target bucket label makes it a real boolean:

```json
{
  "targeting": {
    "if": [
      { "==": [{ "var": "locale" }, "en-US"] },
      {
        "if": [
          {
            "==": [
              {
                "fractional": [
                  { "var": "targetingKey" },
                  ["on", 10],
                  ["off", 90]
                ]
              },
              "on"
            ]
          },
          { "resolve": ["on", "SPLIT_RANDOM", [["details", "us-10pct-rollout"]]] },
          { "resolve": ["off", "SPLIT_FALLTHROUGH", [["details", "us-holdback"]]] }
        ]
      },
      { "resolve": ["off", "FALLTHROUGH", [["details", "non-us-off"]]] }
    ]
  }
}
```

The 90% US holdback and the non-US off case both serve `off` but carry different `reason` and `details` values, so they stay distinguishable in telemetry.

### Narrowing the coarse `reason`

flagd's top-level `reason` today (`TARGETING_MATCH`, `DEFAULT`, `STATIC`, `DISABLED`, `ERROR`, `FALLBACK`) says _what class_ of evaluation happened but not _which specific path_ inside it — every rule-driven result collapses to `TARGETING_MATCH`. A rule-scoped `reason` plus `metadata.details` sub-classifies that. The schema deliberately does not fix an enum: `reason` is an arbitrary string the operator defines, so each fleet picks whatever taxonomy fits it.

Note that the OpenFeature spec ([evaluation details, requirement 6.1](https://openfeature.dev/specification/sections/evaluation-context)) types `reason` as a free-form string — the values in the spec (`TARGETING_MATCH`, `SPLIT`, `DEFAULT`, …) are _recommended_, not exhaustive, and providers are explicitly allowed to emit their own.
That means letting a branch-scoped `reason` override the top-level `reason` is spec-legal without any SDK contract change; SDK hooks and telemetry sinks already treat the field as opaque.

The user-facing reason does not match what the flagd engine actually evaluated for different variants of `SPLIT`. Today `SPLIT` covers both fractional and gradual rollout. Letting the rule override the reason is an easy win here, and the same applies to the several distinct kinds of `OVERRIDE`.

```json
{ "resolve": ["on", "SPLIT_GRADUAL", [["details", "5%→50% ramp, week 3"]]] }
{ "resolve": ["on", "SPLIT_STEPPED", [["details", "stage 2 of 4, cohort=eu-west"]]] }
{ "resolve": ["on", "SPLIT_RANDOM", [["details", "50/50 A/B, salt=exp-4231"]]] }
```

Seeing these reasons instead of a bare `SPLIT` allows for much richer telemetry signals.

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

Beyond Go unit tests in `core/pkg/evaluator`, the flagd gherkin suite gets extended to cover the new terminus end-to-end: testbed flag configs using `resolve`, configs using it outside a terminal position, and scenarios in the shared feature files asserting the resolved `variant`, the overridden `reason`, and the merged `flagMetadata` (including `details`).

Those scenarios run across the gRPC, OFREP, and in-process paths that `test/integration` already exercises, which keeps every provider honest about the new shape rather than testing it only inside flagd.

## Consequences

- **Good** — direct answer to "which branch fired" without overloading the variant key; no wire-format changes. `flagMetadata` is already emitted on OpenFeature evaluation events and picked up by SDK hooks (OpenTelemetry, logging, custom exporters), so anything reading it — OTel spans, A/B dashboards, debug logs — picks up rule-level attribution automatically.
- **Good** — backwards compatible for existing configs: string, boolean, number, and null returns are byte-for-byte unchanged on an upgraded flagd binary.
- **Good** — no OpenFeature spec change required. `reason` is already a first-class spec field, and `details` rides `flagMetadata`, so no SDK gains a new resolution field.
- **Bad** — not forward compatible: a config using `resolve`, evaluated by an older flagd binary, fails with an unknown-operator error. No schema/feature version negotiation exists today, so rollouts must upgrade flagd before configs start using `resolve`. Needs an explicit migration note.
- **Bad** — automatic telemetry attribution requires the separate `RecordEvaluation` work called out in Background; without it, rule ids ride `flagMetadata` for SDK-side hooks but don't reach flagd's own OTel metrics.

## Open questions

- **Config version.** Proposal: bump the targeting/flags schema to `v0.1` and gate the new terminus on it, leaving `v0` untouched for existing fleets. Sequencing of the hosted schema at `flagd.dev` is covered in the migration warning above.
- **Fractional buckets.** `fractional` returns a variant string directly, so it can't tag the picked bucket with metadata without a separate operator extension (e.g. an optional third element per weight tuple). Proposal: defer — ship the base `resolve` operator first. Nesting metadata into fractional is meaningfully more complicated and shouldn't block this.
