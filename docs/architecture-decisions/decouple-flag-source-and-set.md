---
status: accpted
author: @tangenti
created: 2025-06-16
updated: 2025-06-16
---

# Decouple flag sync sources and flag sets

The goal is to support dynamic flag sets for flagd providers and decouple sources and flag sets.

## Background

Flagd daemon syncs flag configurations from multiple sources. A single source provides a single config, which has an optional flag set ID that may or may not change in the following syncs of the same source.

The in-process provider uses `selector` to specify the desired source.  In order to get a desired flag set, a provider has to stick to a source that provides that flag set. In this case, the flagd daemon cannot remove a source without breaking the dependant flagd providers.

Assumptions of the current model

- `flagSetId`s must be unique across different sources or the configuration is considered invalid.
- In-process providers request at most one flag set.

## Requirements

- Flagd daemon can remove a source without breaking in-process providers that depend on the flag set the source provides.
- In-process providers can select based on flag sets.
- No breaking changes for the current usage of `selector`

## Proposal

### API change

#### Flag Configuration Schema

Add an optional field `flagsetID` under `flag` or `flag.metadata`. The flag set ID cannot be specified if a flag set ID is specified for the config.

### Flagd Sync Selector

Selector will be extended for generic flags selection, starting with checking the equivalence of `source` and `flagsetID` of flags.

Example

```yaml
# Flags from the source `override`
selector: override

# Flags from the source `override`
selector: source=override

# Flags from the flag set `project-42`
selector: flagsetID=project-42
```

The semantic can later be extended with a more complex design, such as AIP-160 filter or Kubernetes selections. This is out of the scope of this ADR.

### Flagd Daemon Storage

1. Flagd will have separate stores for `flags` and `sources`

2. `selector` will be removed from the store

3. `flagSetID` will be added as part of `model.Flag` or under `model.Flag.Metadata` for better consistency with the API.

### Flags Sync

Sync server would count the extended syntax of `selector` and filter the list of flags on-the-fly answering the requests from the providers.

The existing conflict resolving based on sources remains the same. Resyncs on removing flags remains unchanged as well.

## Consequences

### The good

- One source can have multiple flag sets.
- `selector` works on a more grandular level.
- No breaking change
- Sync servers and clients now hold the same understanding of the `selector` semantic.
