---
status: draft
author: Simon Schrottner
created: 2026-04-01
updated: 2026-04-01
---

# ADR: Modular flagd Builder

This document proposes a modular build system for flagd inspired by the [OpenTelemetry Collector Builder (ocb)](https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder).
The goal is to allow users to compose custom flagd binaries containing only the sync providers, service endpoints, evaluators, telemetry, and middleware they need — and to enable users to contribute their own implementations without forking flagd.

## Background

flagd's [multi-sync architecture](./multiple-sync-sources.md) was a foundational decision that decoupled flag storage from the evaluation engine. The `ISync` interface enabled the community to add support for file, HTTP, gRPC, Kubernetes CRDs, and cloud blob storage (GCS, Azure Blob, S3). This extensibility has been a strength, but it has come at a cost: every flagd binary includes every provider and every dependency, regardless of what the user actually needs.

Today, the `core` Go module pulls in approximately 700 transitive dependencies through its `go.mod`.
This includes the full AWS SDK v2, Google Cloud SDK, Azure SDK, Kubernetes client-go, the wazero WebAssembly runtime, and gocloud.dev with all three cloud blob drivers registered via side-effect imports.
The `SyncBuilder` in `core/pkg/sync/builder/syncbuilder.go` unconditionally imports all sync provider packages.
Additionally, `blob_sync.go` registers all cloud drivers via blank imports (`_ "gocloud.dev/blob/s3blob"`, etc.).
There are no build tags or conditional compilation — everything is always compiled in.

Similarly, flagd's three service endpoints — the ConnectRPC/gRPC flag evaluation service, the OFREP REST service, and the gRPC flag sync service — are hardcoded in `flagd/pkg/runtime/from_config.go`.
All three are always instantiated and started.
Users cannot selectively disable endpoints, nor can they add custom endpoints (e.g., a WebSocket adapter, an admin API, or a custom protocol) without forking flagd.

Observability adds another dimension to this problem.
The `core/pkg/telemetry` package pulls in `go.opentelemetry.io/contrib/exporters/autoexport`, which transitively brings in all OTLP exporters (gRPC, HTTP), cloud resource detectors, and the Prometheus translator.
While flagd has an `IMetricsRecorder` abstraction with a `NoopMetricsRecorder` fallback, other telemetry usage is hardcoded: the evaluator calls `otel.Tracer()` directly against the global OTel state, and service endpoints use `otelgrpc.NewServerHandler()`, `otelhttp.NewHandler()`, and `otelconnect.NewInterceptor()` without injection.
This means even a minimal flagd build carries the full OTel SDK and exporter stack.

The OpenTelemetry Collector project faced a nearly identical problem and solved it with the [OpenTelemetry Collector Builder](https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder).
It is a CLI tool that reads a YAML manifest specifying which components to include, generates Go source code from templates, and compiles a custom binary with only the selected dependencies.
Each component is a separate Go module exposing a `NewFactory()` function.
The builder generates a `components.go` that imports and registers only the selected factories, a `main.go` entrypoint, and a `go.mod` with only the required dependencies.
Users can reference their own Go modules in the manifest to add custom components.
This approach is proven at scale across hundreds of OTel Collector distributions.

## Requirements

* **Selective compilation**: Users must be able to choose which sync providers, service endpoints, evaluators, telemetry providers, and middleware are compiled into their flagd binary, eliminating unused dependencies entirely (not just disabling them at runtime).
* **User extensibility**: Users must be able to implement custom sync providers, service endpoints, evaluators, telemetry providers, or middleware in their own Go modules and include them in a flagd build without forking the project.
* **Backward compatibility**: The current monolithic flagd binary must remain available as a "full" distribution. Existing users should not be forced to adopt the builder.
* **Standard distributions**: The project must ship pre-configured distributions for common use cases (minimal, cloud, Kubernetes, full) so that most users never need to run the builder themselves.
* **Clean interfaces**: Each component type (sync, evaluator, service, telemetry, middleware) must have a well-defined factory interface that all implementations — official and user-provided — implement.
* **Standard Go tooling**: The build process must use standard Go modules and `go build`. No custom package managers, no dynamic linking, no Go plugins.

## Considered Options

* **Go build tags**: Use `//go:build` tags to conditionally compile sync providers and services. For example, `//go:build with_s3` would gate the S3 provider.
* **OTel-style builder with code generation**: A builder CLI tool that generates Go source files from a YAML manifest and compiles a custom binary. Each component is a separate Go module.
* **Go plugin system**: Use Go's `plugin` package to dynamically load sync providers and services at runtime as `.so` files.

## Proposal

We propose adopting the **OTel-style builder with code generation** approach, extended to cover not just sync providers but also service endpoints, evaluators, telemetry, and middleware as first-class selectable components.

### Why not build tags?

Build tags are fragile and hard to compose. More critically, they do not eliminate dependencies from `go.mod` — the modules are still referenced even if their code is gated by build tags, and Go's module system still downloads and resolves them. Build tags also require maintaining `_tag.go` and `_notag.go` file pairs, which increases complexity and the surface area for bugs.

### Why not Go plugins?

Go's `plugin` package has severe limitations: it only works on Linux (and limited other Unix-like systems), requires both the host and plugin to be compiled with the exact same Go version and build flags, and does not support Windows or macOS. The OTel Collector project explicitly evaluated and rejected this approach for these reasons.

### Component taxonomy

The builder operates on five component types:

| Component Type | Factory Interface | Current Implementations | Example User Extension |
|---------------|------------------|------------------------|----------------------|
| **Sync Provider** | `SyncFactory` | file, kubernetes, http, grpc, gcs, azblob, s3 | Consul, etcd, Vault |
| **Service Endpoint** | `ServiceFactory` | flag-evaluation (ConnectRPC), ofrep (REST), flag-sync (gRPC) | WebSocket adapter, admin API, custom protocol |
| **Evaluator** | `EvaluatorFactory` | jsonlogic, wasm | Custom evaluation engine |
| **Telemetry** | `TelemetryFactory` | otel (tracing + metrics), prometheus, noop | Datadog, custom metrics backend |
| **Middleware** | `MiddlewareFactory` | cors, h2c, metrics | Rate limiting, auth, logging |

### Factory interfaces

Each component type gets a factory interface in the `core` module. The `core` module becomes lightweight (interfaces only, no external dependencies beyond stdlib):

```go
// core/pkg/sync/factory.go
type SyncDependencies struct {
    Config SourceConfig
    Logger *logger.Logger
}

type SyncFactory interface {
    Type() string
    Schemes() []string
    Create(deps SyncDependencies) (ISync, error)
}

// core/pkg/service/factory.go
type ServiceFactory interface {
    Type() string
    Create(deps ServiceDependencies) (Service, error)
}

type Service interface {
    Start(ctx context.Context) error
    Shutdown(ctx context.Context) error
}

type ServiceDependencies struct {
    Evaluator evaluator.IEvaluator
    Store     store.IStore
    Logger    *logger.Logger
    Telemetry telemetry.Provider
    Config    Configuration
    EventBus  EventBus
}

// core/pkg/evaluator/factory.go
type EvaluatorDependencies struct {
    Store     store.IStore
    Logger    *logger.Logger
    Telemetry telemetry.Provider
}

type EvaluatorFactory interface {
    Type() string
    Create(deps EvaluatorDependencies) (IEvaluator, error)
}

// core/pkg/telemetry/factory.go
type Provider interface {
    Tracer(name string) ITracer
    MetricsRecorder() IMetricsRecorder
    GRPCServerOption() grpc.ServerOption      // e.g. otelgrpc stats handler
    HTTPMiddleware(name string) func(http.Handler) http.Handler  // e.g. otelhttp
    ConnectInterceptor() connect.Interceptor  // e.g. otelconnect
    Shutdown(ctx context.Context) error
}

type TelemetryFactory interface {
    Type() string
    Create(cfg TelemetryConfig, logger *logger.Logger) (Provider, error)
}

// core/pkg/service/middleware.go
type MiddlewareFactory interface {
    Type() string
    Create(deps MiddlewareDependencies) (Middleware, error)
}

type MiddlewareDependencies struct {
    Logger    *logger.Logger
    Telemetry telemetry.Provider
    Config    Configuration
}

type Middleware interface {
    Handler(handler http.Handler) http.Handler
}
```

All factory `Create` methods use a dependencies struct rather than positional arguments.
This ensures interfaces remain backward compatible as new dependencies are introduced — adding a field to a struct is non-breaking, while adding a function parameter is.

The `telemetry.Provider` interface is the key abstraction that decouples all components from specific observability implementations.
Today, the evaluator calls `otel.Tracer()` directly and service endpoints use `otelgrpc`/`otelhttp`/`otelconnect` — all hardcoded.
With this interface, components receive telemetry capabilities via injection.
A `noop` implementation returns no-op tracers, no-op metrics recorders, and passthrough middleware — enabling truly zero-overhead builds when observability is not needed.

Each component module exposes a `NewFactory()` function:

```go
// In github.com/open-feature/flagd/sync/file
func NewFactory() sync.SyncFactory { ... }

// In github.com/open-feature/flagd/service/ofrep
func NewFactory() service.ServiceFactory { ... }

// User-provided: github.com/mycompany/flagd-consul-sync
func NewFactory() sync.SyncFactory { ... }
```

### Module structure

The monolithic `core` module is split so that each component is a separate Go module with its own `go.mod`:

```text
core/                          # Interfaces only (lightweight)
sync/file/                     # depends on core + fsnotify
sync/kubernetes/               # depends on core + k8s.io/client-go
sync/http/                     # depends on core + net/http
sync/grpc/                     # depends on core + google.golang.org/grpc
sync/gcs/                      # depends on core + gocloud.dev/blob/gcsblob
sync/azblob/                   # depends on core + azure-sdk-for-go
sync/s3/                       # depends on core + aws-sdk-go-v2
evaluator/jsonlogic/           # depends on core + jsonlogic
evaluator/wasm/                # depends on core + wazero
service/flag-evaluation/       # depends on core + ConnectRPC
service/ofrep/                 # depends on core + net/http
service/flag-sync/             # depends on core + grpc
telemetry/otel/                # depends on core + OTel SDK + autoexport
telemetry/noop/                # depends on core only (zero dependencies)
```

Blob providers are split individually (not sharing a gocloud.dev base module) so users can include a single cloud provider without pulling all three SDKs.

### Builder manifest

The `flagd-builder` CLI reads a YAML manifest that references components as Go module paths:

```yaml
dist:
  module: github.com/mycompany/custom-flagd
  name: flagd
  version: "0.12.0"
  output_path: ./build

syncs:
  - gomod: "github.com/open-feature/flagd/sync/file v0.12.0"
  - gomod: "github.com/open-feature/flagd/sync/http v0.12.0"
  - gomod: "github.com/mycompany/flagd-consul-sync v1.2.0"  # user extension

evaluators:
  - gomod: "github.com/open-feature/flagd/evaluator/jsonlogic v0.12.0"

services:
  - gomod: "github.com/open-feature/flagd/service/ofrep v0.12.0"
  - gomod: "github.com/mycompany/flagd-admin-api v1.0.0"     # user extension

telemetry:
  - gomod: "github.com/open-feature/flagd/telemetry/otel v0.12.0"
  # Or for a minimal build with no observability overhead:
  # - gomod: "github.com/open-feature/flagd/telemetry/noop v0.12.0"

middleware:
  - gomod: "github.com/open-feature/flagd/middleware/cors v0.12.0"
  - gomod: "github.com/open-feature/flagd/middleware/metrics v0.12.0"

replaces: []
```

### Build process

The builder executes three steps (any of which can be skipped):

1. **Generate**: Execute Go templates to produce `components.go`, `main.go`, and `go.mod`
2. **Get modules**: Run `go mod tidy` to resolve dependencies
3. **Compile**: Run `go build` to produce the binary

### Standard distributions

The project ships pre-configured manifests for common use cases:

| Distribution | Syncs | Evaluators | Services | Telemetry | Target Use Case |
|-------------|-------|------------|----------|-----------|----------------|
| `flagd-minimal` | file, http | jsonlogic | ofrep | noop | Smallest binary, CI/testing, REST-only |
| `flagd-cloud` | file, http, gcs, s3, azblob | jsonlogic | flag-evaluation, ofrep | otel | Cloud storage backends |
| `flagd-kubernetes` | file, http, kubernetes, grpc | jsonlogic | flag-evaluation, ofrep, flag-sync | otel | K8s with OFO |
| `flagd-full` | ALL | ALL | ALL | otel | Current behavior (backward compat) |

### Runtime changes

The `Runtime` struct changes from hardcoded service fields to a dynamic registry:

```go
// Current (hardcoded):
type Runtime struct {
    EvaluationService service.IFlagEvaluationService
    OfrepService      ofrep.IOfrepService
    SyncService       flagsync.ISyncService
    ...
}

// Proposed (dynamic with event bus):
type Runtime struct {
    Services  []service.Service
    Syncs     []sync.ISync
    Evaluator evaluator.IEvaluator
    EventBus  service.EventBus
    ...
}
```

A key concern with a dynamic `[]service.Service` slice is that the current `Runtime` calls typed methods like `r.SyncService.Emit(source)` to notify the sync service of state changes.
With a generic service slice, this typed coordination is lost.

The solution is an **event bus** — a lightweight pub/sub mechanism injected into both the `Runtime` and any service that needs cross-component coordination.
Services that care about sync events (like the flag-sync service) subscribe to the event bus during creation via `ServiceDependencies.EventBus`.
The runtime publishes events to the bus after updating the evaluator state, and subscribed services receive them without the runtime needing to know their concrete types:

```go
// core/pkg/service/event_bus.go
type EventBus interface {
    Publish(event Event)
    Subscribe(eventType string) <-chan Event
}

type Event struct {
    Type   string         // e.g. "sync_complete", "configuration_change"
    Source string         // which sync provider emitted this
    Data   map[string]any
}
```

This keeps the runtime simple — it publishes to the bus, and services self-select which events they consume — while avoiding type assertions or service-specific method calls.

### High availability: snapshot-based resilience

A recurring concern with flagd is what happens when the process restarts and the upstream sync source (gRPC server, S3 bucket, HTTP endpoint) is unavailable.
Today, flagd starts with an empty in-memory store (memdb) and waits for sync providers to deliver data.
If the upstream is down at restart time, flagd has zero flags and every evaluation fails.

The modular architecture — specifically the `EventBus` and the multi-sync precedence model — enables a clean solution: **snapshot sync**.

Two components work together:

**1. Snapshot writer** — an `EventBus` subscriber that persists flag state to a local file on every configuration change.
It listens for `configuration_change` events and writes the current flag state to a configurable path (e.g., `/var/lib/flagd/snapshot.json`).
This component has no external dependencies beyond the filesystem.

**2. Snapshot sync provider** — a read-only `ISync` implementation that loads the snapshot file on startup.
It emits the cached state to the `DataSync` channel during `Init()`, then becomes idle.
Because flagd's multi-sync model uses index-based precedence (lowest index = lowest precedence), the snapshot sync is configured first so that live data from the primary source always takes priority:

```yaml
syncs:
    # Lowest precedence — bootstrap from local cache if upstream is unavailable
    - gomod: "github.com/open-feature/flagd/sync/snapshot v0.12.0"
    # Highest precedence — primary source, overwrites snapshot when available
    - gomod: "github.com/open-feature/flagd/sync/grpc v0.12.0"
```

**Startup behavior:**

1. Snapshot sync loads cached state immediately (milliseconds, local disk)
2. flagd begins serving evaluations with cached flags
3. Primary sync provider connects to upstream (may take seconds or fail entirely)
4. If upstream connects: live data overwrites cached state, snapshot writer updates the cache
5. If upstream is down: flagd continues serving from the snapshot until upstream recovers

This approach has several advantages:

* Users who need HA add two lines to their manifest; users who do not, pay nothing
* It composes with any upstream sync provider (HTTP, gRPC, S3, Kubernetes) without modification
* The snapshot file can be pre-seeded in container images or configuration management for deterministic cold starts
* No changes to the store, evaluator, or runtime — it uses existing `ISync` and `EventBus` interfaces

### Consequences

* Good, because flagd binaries can be dramatically smaller (a file+http+ofrep+noop-telemetry build eliminates all cloud SDKs, Kubernetes client, gRPC, ConnectRPC, wazero, and the entire OTel SDK)
* Good, because users can extend flagd with custom sync providers, service endpoints, evaluators, and middleware without forking
* Good, because the builder pattern is proven at scale by the OpenTelemetry Collector community
* Good, because it uses standard Go modules and tooling — no custom package managers or dynamic linking
* Good, because backward compatibility is maintained via the `flagd-full` distribution
* Good, because the modular architecture enables high-availability patterns (snapshot sync) that were previously difficult to implement without coupling to the core
* Bad, because it is a breaking change for anyone importing `core/pkg/sync/builder` or other internal packages directly
* Bad, because it adds build complexity — users who want custom builds need to learn the builder tool
* Bad, because component modules require coordinated releases (or independent versioning, which adds its own complexity)
* Bad, because the module split creates many more Go modules to maintain in the repository

### Timeline

1. ADR review and acceptance
2. Factory interface design and implementation in core (including `telemetry.Provider`)
3. Module split (sync providers, evaluators, service endpoints, telemetry, middleware)
4. Runtime refactoring (dynamic service registry, injected telemetry)
5. Builder CLI tool implementation
6. Standard distribution manifests
7. CI/CD pipeline updates, Dockerfile changes, documentation

### Open questions

* **Module path structure**: Should component modules live at top-level (`sync/file`, `service/ofrep`) or under existing paths (`core/pkg/sync/file` with separate `go.mod`)?
* **flagd-proxy**: Should `flagd-proxy` also be buildable via the builder, or does it remain as-is?
* **Release versioning**: Should component modules version independently (like OTel contrib) or stay in lockstep with flagd releases?
* **Community component registry**: Should there be a curated list or repository of community-contributed components (analogous to `opentelemetry-collector-contrib`)?
* **Default service set**: Can a build have zero service endpoints (library/embedded mode), or should at least one always be required?
* **Configuration compatibility**: How does the builder interact with the existing `-f`/`--sources` CLI flags?
The runtime needs to know which URI schemes are available from the selected sync factories.
See the [CLI compatibility analysis](#cli-compatibility) below for the recommended approach.
* **Telemetry granularity**: A single `telemetry.Provider` is sufficient for the initial implementation.
Splitting tracing and metrics into independently selectable components can be revisited later if there is demand.

## More Information

### CLI compatibility

A key question is how the existing `-f`/`--sources` CLI flags interact with a modular build where only certain sync providers are compiled in.
If a user builds flagd with only file+http sync but passes `-f s3://my-bucket/flags.json`, the binary must handle this gracefully.

Four options were considered:

**Option A: Keep `-f` as-is, fail fast at runtime.**
The generated code registers each `SyncFactory` with its `Schemes()`.
At startup, URI matching checks against registered schemes only.
An unknown scheme produces a clear error: *"sync provider for scheme 's3://' not available in this build. Available: file:, http://, https://"*.
No CLI changes needed.

* Good, because it is fully backward compatible and introduces zero new concepts
* Bad, because errors are only caught at startup, not at build time

**Option B: Builder validates a runtime config file at build time.**
Add an optional `--validate-config=flagd.yaml` flag to `flagd-builder`.
The builder cross-checks the runtime config's URIs against selected sync factories before compiling.

* Good, because errors are caught before deployment
* Bad, because it only works when the config is known at build time (not always the case with environment variables or dynamic sources)

**Option C: Generated CLI with constrained flags.**
The builder generates CLI code that restricts `--uri` to known schemes, and `--help` self-documents what is available.

* Good, because the binary is self-documenting
* Bad, because it adds complex code generation and breaks the generic `-f` interface

**Option D: Hybrid — runtime registry + introspection command.**
Keep `-f` unchanged (Option A behavior), but also generate a `flagd components` subcommand that lists compiled-in syncs, services, evaluators, and telemetry.
Users can inspect any binary to see what it supports.
Fail-fast errors reference this command.

* Good, because it combines backward compatibility with discoverability
* Bad, because it requires slightly more implementation (though trivial — just iterate registered factories)

**Recommendation: Option A + D combined.**
The `-f` flag stays exactly as-is (no breaking change), the runtime fails fast with a clear error listing available schemes, and `flagd components` gives introspection.
Option B can be layered on later as a nice-to-have.
This mirrors the OpenTelemetry Collector's approach — unknown component types produce a clear error at startup, and `otelcol components` lists what is compiled in.

* [OpenTelemetry Collector Builder (ocb)](https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder) — the reference implementation this proposal is based on
* [ADR: Multiple Sync Sources](./multiple-sync-sources.md) — the prior ADR establishing flagd's multi-sync architecture and the `ISync` interface
* [OTel Collector component model](https://opentelemetry.io/docs/collector/custom-collector/) — documentation on building custom collectors
* [gocloud.dev](https://gocloud.dev/) — the cloud abstraction library currently used for blob sync providers
