---
status: draft
author: Simon Schrottner
created: 2026-04-01
updated: 2026-04-01
---

# ADR: Modular flagd Builder

This document proposes a modular build system for flagd inspired by the [OpenTelemetry Collector Builder (ocb)](https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder). The goal is to allow users to compose custom flagd binaries containing only the sync providers, service endpoints, evaluators, and middleware they need — and to enable users to contribute their own implementations without forking flagd.

## Background

flagd's [multi-sync architecture](./multiple-sync-sources.md) was a foundational decision that decoupled flag storage from the evaluation engine. The `ISync` interface enabled the community to add support for file, HTTP, gRPC, Kubernetes CRDs, and cloud blob storage (GCS, Azure Blob, S3). This extensibility has been a strength, but it has come at a cost: every flagd binary includes every provider and every dependency, regardless of what the user actually needs.

Today, the `core` Go module pulls in approximately 700 transitive dependencies through its `go.mod`.
This includes the full AWS SDK v2, Google Cloud SDK, Azure SDK, Kubernetes client-go, the wazero WebAssembly runtime, and gocloud.dev with all three cloud blob drivers registered via side-effect imports.
The `SyncBuilder` in `core/pkg/sync/builder/syncbuilder.go` unconditionally imports all sync provider packages.
Additionally, `blob_sync.go` registers all cloud drivers via blank imports (`_ "gocloud.dev/blob/s3blob"`, etc.).
There are no build tags or conditional compilation — everything is always compiled in.

Similarly, flagd's three service endpoints — the ConnectRPC/gRPC flag evaluation service, the OFREP REST service, and the gRPC flag sync service — are hardcoded in `flagd/pkg/runtime/from_config.go`. All three are always instantiated and started. Users cannot selectively disable endpoints, nor can they add custom endpoints (e.g., a WebSocket adapter, an admin API, or a custom protocol) without forking flagd.

The OpenTelemetry Collector project faced a nearly identical problem and solved it with the [OpenTelemetry Collector Builder](https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder).
It is a CLI tool that reads a YAML manifest specifying which components to include, generates Go source code from templates, and compiles a custom binary with only the selected dependencies.
Each component is a separate Go module exposing a `NewFactory()` function.
The builder generates a `components.go` that imports and registers only the selected factories, a `main.go` entrypoint, and a `go.mod` with only the required dependencies.
Users can reference their own Go modules in the manifest to add custom components.
This approach is proven at scale across hundreds of OTel Collector distributions.

## Requirements

* **Selective compilation**: Users must be able to choose which sync providers, service endpoints, evaluators, and middleware are compiled into their flagd binary, eliminating unused dependencies entirely (not just disabling them at runtime).
* **User extensibility**: Users must be able to implement custom sync providers, service endpoints, evaluators, or middleware in their own Go modules and include them in a flagd build without forking the project.
* **Backward compatibility**: The current monolithic flagd binary must remain available as a "full" distribution. Existing users should not be forced to adopt the builder.
* **Standard distributions**: The project must ship pre-configured distributions for common use cases (minimal, cloud, Kubernetes, full) so that most users never need to run the builder themselves.
* **Clean interfaces**: Each component type (sync, evaluator, service, middleware) must have a well-defined factory interface that all implementations — official and user-provided — implement.
* **Standard Go tooling**: The build process must use standard Go modules and `go build`. No custom package managers, no dynamic linking, no Go plugins.

## Considered Options

* **Go build tags**: Use `//go:build` tags to conditionally compile sync providers and services. For example, `//go:build with_s3` would gate the S3 provider.
* **OTel-style builder with code generation**: A builder CLI tool that generates Go source files from a YAML manifest and compiles a custom binary. Each component is a separate Go module.
* **Go plugin system**: Use Go's `plugin` package to dynamically load sync providers and services at runtime as `.so` files.

## Proposal

We propose adopting the **OTel-style builder with code generation** approach, extended to cover not just sync providers but also service endpoints, evaluators, and middleware as first-class selectable components.

### Why not build tags?

Build tags are fragile and hard to compose. More critically, they do not eliminate dependencies from `go.mod` — the modules are still referenced even if their code is gated by build tags, and Go's module system still downloads and resolves them. Build tags also require maintaining `_tag.go` and `_notag.go` file pairs, which increases complexity and the surface area for bugs.

### Why not Go plugins?

Go's `plugin` package has severe limitations: it only works on Linux (and limited other Unix-like systems), requires both the host and plugin to be compiled with the exact same Go version and build flags, and does not support Windows or macOS. The OTel Collector project explicitly evaluated and rejected this approach for these reasons.

### Component taxonomy

The builder operates on four component types:

| Component Type | Factory Interface | Current Implementations | Example User Extension |
|---------------|------------------|------------------------|----------------------|
| **Sync Provider** | `SyncFactory` | file, kubernetes, http, grpc, gcs, azblob, s3 | Consul, etcd, Vault |
| **Service Endpoint** | `ServiceFactory` | flag-evaluation (ConnectRPC), ofrep (REST), flag-sync (gRPC) | WebSocket adapter, admin API, custom protocol |
| **Evaluator** | `EvaluatorFactory` | jsonlogic, wasm | Custom evaluation engine |
| **Middleware** | `MiddlewareFactory` | cors, h2c, metrics | Rate limiting, auth, logging |

### Factory interfaces

Each component type gets a factory interface in the `core` module. The `core` module becomes lightweight (interfaces only, no external dependencies beyond stdlib):

```go
// core/pkg/sync/factory.go
type SyncFactory interface {
    Type() string
    Schemes() []string
    Create(cfg SourceConfig, logger *logger.Logger) (ISync, error)
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
    Config    Configuration
}

// core/pkg/evaluator/factory.go
type EvaluatorFactory interface {
    Type() string
    Create(store store.IStore, logger *logger.Logger) (IEvaluator, error)
}
```

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

| Distribution | Syncs | Evaluators | Services | Target Use Case |
|-------------|-------|------------|----------|----------------|
| `flagd-minimal` | file, http | jsonlogic | ofrep | Smallest binary, CI/testing, REST-only |
| `flagd-cloud` | file, http, gcs, s3, azblob | jsonlogic | flag-evaluation, ofrep | Cloud storage backends |
| `flagd-kubernetes` | file, http, kubernetes, grpc | jsonlogic | flag-evaluation, ofrep, flag-sync | K8s with OFO |
| `flagd-full` | ALL | ALL | ALL | Current behavior (backward compat) |

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

// Proposed (dynamic):
type Runtime struct {
    Services  []service.Service
    Syncs     []sync.ISync
    Evaluator evaluator.IEvaluator
    ...
}
```

### Consequences

* Good, because flagd binaries can be dramatically smaller (a file+http+ofrep build eliminates all cloud SDKs, Kubernetes client, gRPC, ConnectRPC, and wazero)
* Good, because users can extend flagd with custom sync providers, service endpoints, evaluators, and middleware without forking
* Good, because the builder pattern is proven at scale by the OpenTelemetry Collector community
* Good, because it uses standard Go modules and tooling — no custom package managers or dynamic linking
* Good, because backward compatibility is maintained via the `flagd-full` distribution
* Bad, because it is a breaking change for anyone importing `core/pkg/sync/builder` or other internal packages directly
* Bad, because it adds build complexity — users who want custom builds need to learn the builder tool
* Bad, because component modules require coordinated releases (or independent versioning, which adds its own complexity)
* Bad, because the module split creates many more Go modules to maintain in the repository

### Timeline

1. ADR review and acceptance
2. Factory interface design and implementation in core
3. Module split (sync providers, evaluators, service endpoints, middleware)
4. Runtime refactoring (dynamic service registry)
5. Builder CLI tool implementation
6. Standard distribution manifests
7. CI/CD pipeline updates, Dockerfile changes, documentation

### Open questions

* **Module path structure**: Should component modules live at top-level (`sync/file`, `service/ofrep`) or under existing paths (`core/pkg/sync/file` with separate `go.mod`)?
* **flagd-proxy**: Should `flagd-proxy` also be buildable via the builder, or does it remain as-is?
* **Release versioning**: Should component modules version independently (like OTel contrib) or stay in lockstep with flagd releases?
* **Community component registry**: Should there be a curated list or repository of community-contributed components (analogous to `opentelemetry-collector-contrib`)?
* **Default service set**: Can a build have zero service endpoints (library/embedded mode), or should at least one always be required?
* **Configuration compatibility**: How does the builder interact with the existing `-f`/`--sources` CLI flags? The runtime needs to know which URI schemes are available from the selected sync factories.

## More Information

* [OpenTelemetry Collector Builder (ocb)](https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder) — the reference implementation this proposal is based on
* [ADR: Multiple Sync Sources](./multiple-sync-sources.md) — the prior ADR establishing flagd's multi-sync architecture and the `ISync` interface
* [OTel Collector component model](https://opentelemetry.io/docs/collector/custom-collector/) — documentation on building custom collectors
* [gocloud.dev](https://gocloud.dev/) — the cloud abstraction library currently used for blob sync providers
