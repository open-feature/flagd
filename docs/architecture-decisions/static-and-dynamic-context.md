---
# Valid statuses: draft | proposed | rejected | accepted | superseded
status: accepted
author: @leakonvalinka
created: 2026-01-21
updated: 2026-01-21
---

# Static and Dynamic Context Enrichment in flagd

This document retrospectively records the introduction of static context enrichment (`--context-value` or `-X`) and dynamic context enrichment (`--context-from-header` or `-H`) for flag evaluations in the flagd daemon.
It explains the purpose and basic use cases for both options and defines the merge priority if the same context item can be found in more than one provided context source.

## Background

flagd is an open-source feature flagging engine that allows targeted flag evaluations based on optional context data.
A simple example for such context data would be the current user's email: if it ends with `@example.com`, the evaluation returns true, otherwise false.
Initially, context could only be provided by including it in the evaluation request's body like this:

```json
{
  "flagKey": "booleanFlagKey",
  "context": {
    "email": "noreply@example.com"
  }
}
```

## Requirements

* support HTTP header mapping in OFREP requests
* support HTTP header mapping in evaluation service v2 via connect
* clearly defined priority list for merging context values from multiple sources

## Considered Options

* **Static context enrichment**: `--context-value` or `-X` \
On startup, static context data can be provided in the form of concrete key-value pairs, which is then used as context in every evaluation.
This is aimed at attributes that do not change during a flagd instance’s lifetime (such as the server region or cloud provider), reducing effort on the client-side.
For example, `flagd start -X region=europe ...` adds a `region` attribute with value `europe` as a context to every evaluation.

* **Dynamic context enrichment**: `--context-from-header` or `-H`\
On startup, specific request headers can be configured so that, for each incoming request, their values are automatically extracted and included as context attributes for the evaluation.
This is targeted at dynamic values that likely vary per request (like the email of the user making the request).
For example, `flagd start -H X-User-Email:userEmail ...` tells flagd to extract the value of the `X-User-Email` header from each incoming request and include it as the `userEmail` attribute in the evaluation context, if provided.

* **Merge priority definition**: \
In case the same context key can be found in more than one context source, this priority list defines which value takes precedence (from highest to lowest):
  1. Dynamic context from request headers (`-H`)
  2. Static context from startup options (`-X`)
  3. Context provided in the evaluation request body

## Proposal

All options in the section above were accepted and implemented as described.

### Consequences

* Good, because static context values reduce effort on the client-side for attributes that do not change often.
* Good, because this allows a more targeted way of providing context depending on whether the attribute is static or dynamic.
* Bad, because the merging of context values from multiple sources introduces new complexity that needs to be learned and understood by users.

### Timeline

* Static context enrichment: The issue was created in October 2024, the final PR was merged in December 2024.
* Dynamic context enrichment: The issue was created in March 2025, the final PR was merged in June 2025.

## More Information

* [Static context enrichment Issue](https://github.com/open-feature/flagd/issues/1435)
* [Dynamic context enrichment Issue](https://github.com/open-feature/flagd/issues/1583)