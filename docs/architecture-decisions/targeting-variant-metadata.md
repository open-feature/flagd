---
# Valid statuses: draft | proposed | rejected | accepted | superseded
status: proposed
author: Parth Suthar
created: 2026-07-30
updated: 2026-09-02
---

# Per-evaluation metadata from targeting rules

Let a targeting rule return `{ "variant": "<key>", "reason": "<string>", "metadata": { ... } }` in place of a plain variant string, so the _branch that fired_ can annotate the evaluation with extra metadata, custom reasons, and other auxiliary information. Existing string returns are unchanged.

## Background

An evaluation today returns `value`, `variant`, `reason`, and `metadata`. The `reason` is a coarse enum (`TARGETING_MATCH`, `DEFAULT`, `STATIC`, …); `metadata` carries only the static blocks defined at flag-set and flag level. Neither answers the debugging question we hit most often: **which branch of the targeting expression fired?**

For a nested `if` / `and` / `or` tree the resolver just returns the winning variant key. Two rules landing on the same variant are indistinguishable in the response — the only way to know _why_ today is to fetch the flag config and re-run the logic by hand. Encoding the branch identity into the variant key (`"clubs-eu-rollout-a"`) is the workaround, and it pollutes the variant space with debug info consumers then have to parse back out.

The plumbing to carry metadata already exists on most paths: `AnyValue.Metadata` (`core/pkg/evaluator/ievaluator.go`) is threaded through the evaluator and surfaced via OFREP (single and bulk) and single-resolution gRPC responses as `flagMetadata`. Two paths don't carry it yet — bulk gRPC `ResolveAll` and `RecordEvaluation` telemetry — but wiring those up is separate follow-up work.

## Proposal

In `definitions.primitive` in `schemas/json/targeting.json`, the `string` entry (the variant-key return) becomes a choice of the plain string or the tagged object — the other `primitive` entries (`null`, `boolean`, `number`, `array`) are untouched:

```json
{
  "description": "When returned from rules, strings are used as keys to retrieve the associated value from the \"variants\" object. Be sure that the returned string is present as a key in the variants! As of <version>, an object of the form { \"variant\": ..., \"reason\": ..., \"metadata\": ... } may be used instead, to attach a custom reason and arbitrary metadata to the branch that fired.",
  "oneOf": [
    { "type": "string" },
    {
      "type": "object",
      "required": ["variant"],
      "additionalProperties": false,
      "properties": {
        "variant": { "type": "string" },
        "reason": { "type": "string" },
        "metadata": {
          "allOf": [
            {
              "$ref": "https://flagd.dev/schema/v0/flags.json#/definitions/metadata"
            },
            {
              "properties": {
                "details": {
                  "type": "string",
                  "description": "Free-text explanation of why this branch fired."
                }
              }
            }
          ]
        }
      }
    }
  ]
}
```

### Field types

- **`variant`** — string. The only **required** property, and the exact equivalent of the plain-string return: `"red"` and `{ "variant": "red" }` are interchangeable terminals selecting the variant `red`.
- **`reason`** — string, optional. Free-form, operator-defined; overrides the top-level evaluation `reason` for this branch. The schema does not constrain it to an enum (see [Narrowing the coarse `reason`](#narrowing-the-coarse-reason)).
- **`metadata`** — object, optional, shallow. Reuses the existing `metadata` definition, so its values stay restricted to `string | number | boolean`, matching flag and flag-set metadata. Nested objects and arrays are not permitted.
- **`metadata.details`** — string, optional. A well-known key, declared in the schema so editors autocomplete and validate it, but carried inside `metadata` rather than promoted to a peer of `variant`.

### Evaluator changes

In `evaluateVariant` (`core/pkg/evaluator/json.go`), the current unconditional quote-strip of the JsonLogic result gains an object case ahead of it:

- object with a `variant` field → use it as the variant key; shallow-merge the rule's `metadata` object into the returned metadata; if `reason` is present, use it in place of `TARGETING_MATCH`
- everything else → unchanged from today: `null` exits targeting to `defaultVariant`; strings, booleans, and numbers are quote-stripped into a variant key (so a boolean `true` still resolves the `"true"` variant, as `schemas/json/targeting.json` documents and the existing tests cover); a key absent from `variants` errors exactly as it does now

Metadata merge precedence, lowest → highest, so more specific wins: flag-set metadata → flag metadata → **rule-returned metadata**.

!!! warning "Breaking change to the config schema"

    Existing configs keep working on a new binary, but the reverse does not hold: a config using the tagged-object return, evaluated by an older flagd binary, fails — the old resolver quote-strips the JsonLogic result and gets a non-existent variant key (or panics on the single-key `variant` map, which it has no operator registered for).

    There's no schema/feature version negotiation today, so the config schema version has to gate it: operators on an older binary must stay on the older schema version, and flagd has to be upgraded fleet-wide before any config starts emitting the new shape. This needs a schema version bump and a migration note in the release.

    Gating on the schema version is sufficient on its own — providers and flagd embed a snapshot of the schema in each version, so an older release rejects the new shape at validation time rather than misbehaving at evaluation time.

    **Hold the hosted schema back.** `https://flagd.dev/schema/v0/flags.json` is what editors and IDE tooling resolve against. Do not publish the new shape there until it has landed in flagd and a few providers — otherwise editors will happily autocomplete a construct most of the fleet cannot evaluate.

One wrinkle: the Go JsonLogic engine (`github.com/diegoholiveira/jsonlogic`) treats a returned map as a data literal only when it has **more than one key** (`apply()` in `jsonlogic.go`); a single-key map is looked up as an operator and panics if unregistered. Register `variant` as a passthrough JsonLogic operator — same way `fractional`, `starts_with`, `ends_with`, and `semver` are already registered — so all shapes reach the typed decode instead of panicking.

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
            "variant": "on", // required
            "reason": "my_super_cool_custom_reason", // optional
            "metadata": {}, // custom metadata, optionally including "details"
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
          "variant": "alpha",
          "reason": "TARGETING_MATCH",
          "metadata": { "details": "explicit allowlist for enterprise pilot" }
        },
        { "in": [{ "var": "customerId" }, ["customer-B1", "customer-B2"]] },
        { "variant": "beta", "reason": "TARGETING_MATCH" },
        { "variant": "ga", "reason": "FALLTHROUGH" }
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
          {
            "variant": "on",
            "reason": "SPLIT_RANDOM",
            "metadata": { "details": "us-10pct-rollout" }
          },
          {
            "variant": "off",
            "reason": "SPLIT_FALLTHROUGH",
            "metadata": { "details": "us-holdback" }
          }
        ]
      },
      {
        "variant": "off",
        "reason": "FALLTHROUGH",
        "metadata": { "details": "non-us-off" }
      }
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
{ "variant": "on", "reason": "SPLIT_GRADUAL", "metadata": { "details": "5%→50% ramp, week 3" } }
{ "variant": "on", "reason": "SPLIT_STEPPED", "metadata": { "details": "stage 2 of 4, cohort=eu-west" } }
{ "variant": "on", "reason": "SPLIT_RANDOM",  "metadata": { "details": "50/50 A/B, salt=exp-4231" } }
```

Seeing these reasons instead of a bare `SPLIT` allows for much richer telemetry signals.

Other patterns the free-form field unlocks inside a single `TARGETING_MATCH`:

- **Split vs. allowlist under the same top-level reason** — two paths reach the same variant, telemetry can tell them apart:

  ```json
  { "variant": "on", "reason": "SPLIT", "metadata": { "details": "10% cohort bucketed by email" } }
  { "variant": "on", "reason": "ALLOWLIST", "metadata": { "details": "customerId in enterprise-pilot" } }
  ```

- **Kill switch / override tagged distinctly** — a branch short-circuiting an incident stays visible in traces:

  ```json
  {
    "variant": "off",
    "reason": "OVERRIDE",
    "metadata": { "details": "incident-4231 kill switch" }
  }
  ```

- **Audience-segment attribution** — which specific condition inside a large predicate fired:

  ```json
  {
    "variant": "beta",
    "reason": "AUDIENCE_MATCH",
    "metadata": { "details": "country=US AND app_version>=2.0" }
  }
  ```

### Testing

Beyond Go unit tests in `core/pkg/evaluator`, the flagd gherkin suite gets extended to cover the new terminus end-to-end: testbed flag configs using the tagged-object return, and scenarios in the shared feature files asserting the resolved `variant`, the overridden `reason`, and the merged `flagMetadata` (including `details`).

Those scenarios run across the gRPC, OFREP, and in-process paths that `test/integration` already exercises, which keeps every provider honest about the new shape rather than testing it only inside flagd.

## Consequences

- **Good** — direct answer to "which branch fired" without overloading the variant key; no wire-format changes. `flagMetadata` is already emitted on OpenFeature evaluation events and picked up by SDK hooks (OpenTelemetry, logging, custom exporters), so anything reading it — OTel spans, A/B dashboards, debug logs — picks up rule-level attribution automatically.
- **Good** — backwards compatible for existing configs: string, boolean, number, and null returns are byte-for-byte unchanged on an upgraded flagd binary.
- **Good** — no OpenFeature spec change required. `reason` is already a first-class spec field, and `details` rides `flagMetadata`, so no SDK gains a new resolution field.
- **Bad** — not forward compatible: a config using the new object-return shape, evaluated by an older flagd binary, fails (the returned object doesn't quote-strip into a valid variant key). No schema/feature version negotiation exists today, so rollouts must upgrade flagd before configs start using the new shape. Needs an explicit migration note.
- **Bad** — automatic telemetry attribution requires the separate `RecordEvaluation` work called out in Background; without it, rule ids ride `flagMetadata` for SDK-side hooks but don't reach flagd's own OTel metrics.

## Open questions

- **Config version.** Proposal: bump the targeting/flags schema to `v0.1` and gate the new terminus on it, leaving `v0` untouched for existing fleets. Sequencing of the hosted schema at `flagd.dev` is covered in the migration warning above.
- **Fractional buckets.** `fractional` returns a variant string directly, so it can't tag the picked bucket with metadata without a separate operator extension (e.g. an optional third element per weight tuple). Proposal: defer — ship the base object-return shape first. Nesting metadata into fractional is meaningfully more complicated and shouldn't block this.
- **JsonLogic operator registration for `variant`.** Registering `variant` as a passthrough operator is a workaround for a dependency quirk, not a documented contract of `github.com/diegoholiveira/jsonlogic`. Pin the behavior with an integration test so a future engine swap or upgrade doesn't silently regress single-key tagged-object returns.
