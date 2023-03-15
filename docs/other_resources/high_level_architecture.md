# High level architecture

## Component overview

Flagd consists of four main components - Service, Evaluator engine, Runtime and Sync.

The service component exposes the evaluator engine for client libraries.
It further exposes an interface
for flag configuration change notifications.

The sync component has implementations to update flag configurations from various sources.
The current implementation
contain sync providers for files, K8s resources and HTTP endpoints.

The evaluation engine's role is twofold, it acts as an intermediary between configuration changes and the state store by interpreting change events and forwarding the necessary changes to the state store.
It also performs the feature flag evaluations based on evaluation requests coming from feature flag libraries.

The Runtime stays in between these components and coordinates operations.

<img src="../images/of-flagd-0.png" width="560">

## Sync component

The Sync component contains implementations of the ISync interface.
The interface contract simply allows updating
flag configurations watched by the respective implementation.
For example, the file sync provider watches for a change
(ex: - add, modify, remove) of a specific file in the file system.

The update provided by sync implementation is pushed to the evaluator engine, which interprets the event and forwards it to the state store.
Change notifications generated in the
process gets pushed to event subscribers.

<img src="../images/of-flagd-1.png" width="560">

## Readiness & Liveness probes

Flagd exposes HTTP liveness and readiness probes.
These probes can be used for K8s deployments.
With default
start-up configurations, these probes are exposed at the following URLs,

- Liveness: <http://localhost:8014/healthz>
- Readiness: <http://localhost:8014/readyz>

### Definition of Liveness

The liveness probe becomes active and HTTP 200 status is served as soon as Flagd service is up and running.

### Definition of Readiness

The readiness probe becomes active similar to the liveness probe as soon as Flagd service is up and running.
However,
the probe emits HTTP 412 until all sync providers are ready.
This status changes to HTTP 200 when all sync providers at
least have one successful data sync.
The status does not change from there on.
