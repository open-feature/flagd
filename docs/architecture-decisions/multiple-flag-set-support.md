---
status: proposed
author: @alexandraoberaigner
created: 2025-05-28
updated: -
---

# Add support for dynamic usage of Flag Sets to `flagd`

The goal of this decision document is to establish flag sets as a first class concept in `flagd`, and support the dynamic addition/update/removal of flag sets at runtime.

## Background

`flagd` is a language-agnostic feature flagging engine that forms a core part of the OpenFeature ecosystem.
Flag configurations can be stored in different locations so called `sources`. These are specified at startup, e.g.:

````shell
flagd start \
  --port 8013 \
  --uri file:etc/flagd/my-flags-1.json \
  --uri https://my-flags-2.com/flags
````

Currently, we see 2 limitations:

1. Once `flagd` is running, it is not possible to dynamically add flag sources.
2. Logical groups of flags (for example the flags belonging to a particular development team or org unit) are tightly coupled to sources, and clients (usually flagd providers) must be aware of the source in order to select it.

Consequently, it's challenging to implement basic "multi-tenancy" in flagd.
A client should only need concern itself with a logical identifier for its set of flags (ie: `marketing-team-flags`) instead of the resource (source) that happens to correspond to that group of flags (ie: `/etc/flags/marketing-team-flags.json`)
Furthermore, if a new set of flags is introduced, flagd must be restarted with a new source (`--uri /etc/flags/holiday-sale-flags.json`)
These limitations are acceptable when flagd is deployed close to its workload (for example as a kubenetes sidecar or daemon-set) however, to enhance its viability when deployed as a centralized feature-flag service with a variety of client workloads, improvements are necessary.

## Requirements

* Should enable the dynamic usage of flag sets.
* Should support configurations without flag sets.
* Should adhere to existing OpenFeature and flagd terminology and concepts

## Considered Options

1. Addition of flag set support in the [flags schema](https://flagd.dev/reference/schema/#targeting) and associated enhancements to `flagd` storage layer
2. Support for dynamically adding/removing flag sources through some kind of runtime configuration API
3. Support for dynamically adding/removing flag sources through some kind of "discovery" protocol or endpoint (ie: point flagd at a resource that would enumerate a mutable collection of secondary resources which represent flag sets)

## Proposal

To support the dynamic usage of flag sets we propose to adapt the flag schema & storage layer in `flagd`.
The changes will decouple flag sets from flag sources by supporting multiple flag sets within single flag sources.
Dynamic updates to flag sources is already a feature of `flagd`.

### New Schema Structure

The proposed changes to the current flagd schema would allow the following json structure for **sources**:

````json
{
  "$schema": "https://flagd.dev/schema/v1/flagsets.json",
  "flagSets": {
    "my-project-1": {
      "metadata": {
        ...
      },
      "flags": {
        "my-flag-1": {
          "metadata": {
            ...
          },
          ...
        },
        ...
      },
      "$evaluators": {
        ...
      }
    },
    "my-project-2": {
      ...
    }
  }
}
````

We propose to introduce a 3rd json schema `flagSets.json`, which references to `flags.json`:

1. flagSets.json (new)
2. flags.json
3. targeting.json

We don't want to support merging of flag sets, due to implementation efforts & potential confusing behaviour of the
merge strategy.
Therefore, we propose for the initial implementation, `flagSetId`s must be unique across different sources or the configuration is considered invalid.
In the future, it might be useful to support and implement multiple "strategies" for merging flagSets from different sources, but that's beyond the scope of this proposal.

### New Data Structure

The storage layer in `flagd` requires refactoring to better support multiple flag sets within one source.

````go
package store

type State struct {
    FlagSets map[string]FlagSet `json:"flagSets"` // key = flagSetId
}

type FlagSet struct {
    Flags    map[string]model.Flag `json:"flags"` // key = flagKey
    Metadata Metadata              `json:"metadata,omitempty"`
}

type Flag struct {
    State          string          `json:"state"`
    DefaultVariant string          `json:"defaultVariant"`
    Variants       map[string]any  `json:"variants"`
    Targeting      json.RawMessage `json:"targeting,omitempty"`
    Metadata       Metadata        `json:"metadata,omitempty"`
}

type Metadata = map[string]interface{}
````

### OpenFeature Provider Implications

Currently, creating a new flagd provider can look like follows:

````java
final FlagdProvider flagdProvider =
        new FlagdProvider(FlagdOptions.builder()
                .resolverType(Config.Evaluator.IN_PROCESS)
                .host("localhost")
                .port(8015)
                .selector("myFlags.json")
                .build());
````

* With the proposed solution the `flagSetId` should be passed to the builder as selector argument instead of the source.
* `null` is now a valid selector value, referencing flags which do not belong to a flag set. The default/fallback `flagSetId` should be `null`.

### Consequences

* Good, because it decouples flag sets from the sources
* Good, because we will refactor the flagd storage layer (which is currently storing duplicate data & difficult to
  understand)
* Good, because we can support backwards compatibility with the v0 schema
* Good, because the "null" flag set is logically treated as any other flag set, reducing overall implementation complexity.
* Bad, because there's additional complexity to support this new config schema as well as the current.

### Other Options

We evaluated _options 2 + 3: support for dynamically adding/removing flag sources_ and decided against this option because it requires much more implementation effort than _option 1_. Required changes include:

* flagd/core/sync: dynamic mode, which allows specifying the sync type that should be added/removed at runtime
* flagd/flagd: startup dynamic sync configuration
* make sure to still support static syncs

## More Information

* Current flagd schema: [flags.json](https://flagd.dev/schema/v0/flags.json)
* flagd storage layer
  implementation: [store/flags.go](https://github.com/open-feature/flagd/blob/main/core/pkg/store/flags.go)
* [flagd GitHub Repository](https://github.com/open-feature/flagd)
* [OpenFeature Project Overview](https://openfeature.dev/)
