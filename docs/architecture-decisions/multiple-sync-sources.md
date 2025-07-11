---
status: accepted
author: Dave Josephsen
created: 2025-05-21
updated: 2025-05-21
---

# ADR: Multiple Sync Sources

It is the Intent of this document to articulate our rationale for supporting multiple flag syncronization sources (grpc, http, blob, local file, etc..) as a core design property of flagd. This document also includes a short discussion of how flagd is engineered to enable the community to extend it to support new sources in the future, to "future proof" the runtime against sources that don't yet exist, or those we may have omitted is a requisite byproduct of this architectural decision.

The goal of first-class multi-sync support generally is to broaden flagd's potential to suit the needs of many different types of users or architecture. By decoupling flag persistence from the runtime, flagd can focus on evaluation and sync, while enabling its user-base to choose a persistence layer that best suits their individual requirements.

## Background

The flagd daemon is a feature flag evaluation engine that forms a core part of the OpenFeature ecosystem as a production-grade reference implementation. Unlike OpenFeature SDK components, which are, by design, agnostic to the specifics of flag structure, evaluation, and persistence, flagd must take an opinionated stance about how feature-flags look, feel, and act.
What schema best describes a flag? How should they be evaluated? And in what sort of persistence layer should they be stored?

This latter-most question -- _how should they be stored_ -- is the most opaque design detail of every commercial flagging product from the perspective its end-users.
As a front-line engineer using a commercial flagging product, I may, for example, be exposed to flag schema by the product's SDK, and become familiar with its evaluation intricacies over time as my needs grow to require advanced features, or as I encounter surprising behavior. Rarely, however, is an end-user exposed to the details of a commercial product's storage backend.
The SaaS vendor is expected to engineer a safe, fast, multi-tenant storage back-end, optimized for its flag schema and operational parameters, and insulate the customer from these details via its SDK.
This presents Flagd, an open-source evaluation engine, with an interesting conundrum: what sort of flag storage best suits the needs of its potential user-base (which is everyone)?

## Requirements

* Support the storage technology that's most likely to meet the needs of current Flagd user-base (Don't be weird. Don't be surprising.)
* Make it "easy" to extend the flagd runtime to support "new" storage systems
* Horizontally scalable persistence layer
* Minimize end-user exposure to persistence "quirks" (replication schemes, leader election, back-end scaling, consistency minutia, etc.. )
* Reliable, Fast, Transparent
* Full CRUD, read-optimized.

## Considered Options

* Be super-opinionated and prescribe a built-in raftesque key-value setup, analogous to the designs of k8s and kafka, which prescribe etcd and zookeeper respectively.
* Roll a single "standard interface" for flag sync (published grpc spec or similar) (??)
* Decouple storage from flagd entirely, by exposing a Golang interface type that "providers" can implement to provide support for any data store.

## Proposal
<!--
Unsure whether we want a diagram in this section or not. Happy to add one if we want one.
-->
The solution to the conundrum posited in the background section of this document is to decouple flag storage entirely from the rest of the runtime, including instead support for myriad commonly used data syncronization interfaces.
This allows Flagd to be agnostic to flag storage, while enabling users to use whichever storage back-end best suits their environment.

To extend Flagd to support a new storage back-end, _sync providers_ implement the _ISync_ interface, detailed below:

```go
type ISync interface {
 // Init is used by the sync provider to initialize its data structures and external dependencies.
 Init(ctx context.Context) error

 // Sync is the contract between Runtime and sync implementation.
 // Note that, it is expected to return the first data sync as soon as possible to fill the store.
 Sync(ctx context.Context, dataSync chan<- DataSync) error

 // ReSync is used to fetch the full flag configuration from the sync
 // This method should trigger an ALL sync operation then exit
 ReSync(ctx context.Context, dataSync chan<- DataSync) error

 // IsReady shall return true if the provider is ready to communicate with the Runtime
 IsReady() bool
}
```

syncronization events "fan-in" from all configured sync providers to flagd's in-memory state-store via a channel carrying [`sync.DataSync`](https://github.com/open-feature/flagd/blob/main/core/pkg/store/flags.go#L19) events.
These events detail the source and type of the change, along with the flag data in question and are merged into the currently held state by the [store](https://github.com/open-feature/flagd/blob/main/core/pkg/store/flags.go#L19).

### Consequences

Because syncronization providers may vary wildly with respect to their implementation details, supporting multiple sync providers means supporting custom configuration parameters for each provider.
As a consequence, Flagd's configuration is itself made more complex, and its bootstrap process, whose goal is to create a [`runtime.Runtime`](https://github.com/open-feature/flagd/blob/main/flagd/pkg/runtime/runtime.go#L21) object from user-provided configuration, spends the preponderance of its time and effort interpreting, configuring, and initializing sync providers.
There is, in fact, a custom bootstrap type, called the `syncbuilder` whose job is to bootstrap sync providers and arrange them into a map, for the runtime to use.

Further, Because sync providers may vary wildly with respect to implementation, the end-user's choice of sync sources can change Flagd's operational parameters. For example, end-users who choose the GRPC provider can expect flag-sync operations to be nearly immediate, because GRPC updates can be pushed to flagd as they occur, compared with end-users who chose the HTTP provider, who must wait for a timer to expire in order to notice updates, because HTTP is a polling-based implementation.

Finally, sync Providers also contribute a great deal of girth to flagd's documentation, because again, their setup, syntax, and runtime idiosyncrasies may differ wildly.
