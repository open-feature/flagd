---
# Valid statuses: draft | proposed | rejected | accepted | superseded
status: draft
author: Parth Suthar
created: 2026-07-30
updated: 2026-07-30
---
# Per-evaluation metadata from targeting rules

Let a targeting rule return `{ "variant": "<key>", "metadata": { ... } }` in place of a plain variant string, so the *branch that fired* can annotate the evaluation with a rule id or reason.
Existing string returns are unchanged.

## Background

An evaluation today returns `value`, `variant`, `reason`, and `metadata`.
`reason` is a coarse enum (`TARGETING_MATCH`, `DEFAULT`, `STATIC`, …); `metadata` carries only the static blocks defined at flag-set and flag level.
Neither answers the debugging question we hit most often: **which branch of the targeting expression fired?**

For a nested `if` / `and` / `or` tree the resolver just returns the winning variant key.
Two rules landing on the same variant are indistinguishable in the response — the only way to know *why* is to fetch the flag config and re-run the logic against the same context by hand.
Encoding the branch identity into the variant key (`"clubs-eu-rollout-a"`) is the workaround, and it pollutes the variant space with debug info that consumers then have to parse back out.

The plumbing to carry metadata already exists end to end: `AnyValue.Metadata` (`core/pkg/evaluator/ievaluator.go`) is `map[string]interface{}`, threaded through every resolver path and out through gRPC and OFREP as evaluation metadata that OpenFeature SDKs surface as `flagMetadata`.
The only thing missing is a way for a *rule branch* to contribute to it.

## Proposal

Extend `definitions.primitive` in `schemas/json/targeting.json` with a tagged-object return shape:

```json
{
  "type": "object",
  "required": ["variant"],
  "additionalProperties": false,
  "properties": {
    "variant":  { "type": "string" },
    "metadata": { "$ref": "https://flagd.dev/schema/v0/flags.json#/definitions/metadata" }
  }
}
```

Reusing the existing `metadata` definition keeps values restricted to `string | number | boolean`, matching flag and flag-set metadata.

In `evaluateVariant` (`core/pkg/evaluator/json.go`), replace the current string-strip of the JsonLogic result with a typed decode:

* string → variant key, as today
* object with `variant` field → use that as the variant key; shallow-merge `metadata` into the metadata already being returned
* anything else → `PARSE_ERROR`, same as an unrecognized return today

Merge precedence, lowest → highest, so more specific wins:
flag-set metadata → flag metadata → **rule-returned metadata**.

Example — two rules serve the same variant but leave a breadcrumb identifying which one fired:

```json
"targeting": {
  "if": [
    { "==": [{ "var": "region" }, "EU"] },
    { "variant": "on", "metadata": { "rule": "eu-rollout"     } },
    { "variant": "on", "metadata": { "rule": "global-default" } }
  ]
}
```

The EU caller sees `reason=TARGETING_MATCH`, `variant=on`, `metadata.rule=eu-rollout`; the non-EU caller sees the same variant but `metadata.rule=global-default`.
No config re-parsing, no re-derivation in application code.

### Telemetry win

`flagMetadata` is already emitted on OpenFeature evaluation events and picked up by hooks (OpenTelemetry, logging, custom exporters).
Once rule metadata rides that channel, every downstream tool gets rule-level attribution *for free*:

* OTel spans / evaluation events carry a `rule` (or whatever key the config author picked) attribute — filter and group by rule id in traces and logs without touching application code.
* A/B and experimentation dashboards can bucket by rule id emitted by flagd, instead of re-computing which rule fired from raw context.
* Debug logs on the caller side get an immediate "why" alongside the "what" — variant plus rule id in the same log line, no config lookup.

This turns targeting attribution from a config-diving exercise into a queryable dimension anywhere `flagMetadata` already flows.

## Consequences

* **Good** — direct answer to "which branch fired" without overloading the variant key; no wire-format changes; every SDK already surfaces `flagMetadata`.
* **Good** — fully backwards compatible; string returns are unchanged byte-for-byte.
* **Bad** — in-process implementations (Java, JS, Kotlin, Python, .NET) each need the same object-return handling to stay conformant. Covered the usual way — a new suite in [flagd-testbed](https://github.com/open-feature/flagd-testbed).

## Open questions

* **Fractional buckets.** `fractional` returns a variant string directly, so it can't tag the picked bucket with metadata without a separate operator extension (e.g. an optional third element per weight tuple). Proposal: defer — ship the base object-return shape first.
* **JsonLogic literal-object semantics.** We rely on the current Go engine (`github.com/diegoholiveira/jsonlogic`) treating an object with no operator keys as a data literal at return position. Pin with an integration test so a future engine swap doesn't regress it.
* **Metadata size.** No metadata field is capped today; consistency says leave uncapped, but a soft warning at config load is cheap if reviewers want it.
