---
status: proposed
author: @aepfli
created: 2026-01-20
updated: 2026-01-20
---

# Unified Evaluator:  Centralizing Flag Evaluation Logic Across the OpenFeature Ecosystem

Adopt a unified evaluation engine to replace fragmented, language-specific evaluation implementations across OpenFeature SDKs, reducing maintenance burden and ensuring behavioral consistency.

## Background

The OpenFeature ecosystem currently consists of multiple SDK implementations across different programming languages (Java, Python, Go, . NET, Rust, PHP, etc.). Each SDK that integrates with flagd-compatible systems must implement:

1. **Evaluation logic**:  JsonLogic parsing, targeting rules, fractional evaluation
2. **Storage/caching**: In-process flag storage and state management
3. **Validation**:  JSON schema validation for flag configurations
4. **Sync protocols**: Communication with flagd daemon or other sync sources

This distributed approach has led to:

- **Fragmented implementations**: Each language maintains its own evaluator with potential behavioral differences
- **Maintenance overhead**: Bug fixes, feature additions, and specification changes must be replicated across 10+ codebases
- **Inconsistent behavior**: Subtle differences in evaluation logic, validation, and storage can cause production issues
- **Duplicated effort**: The same logic is written, tested, and debugged multiple times
- **Divergent features**: New features may be available in some SDKs but not others
- **Quality variance**: Testing rigor and edge case handling varies by implementation

The fundamental question:  **What if we had a single, canonical implementation that all SDKs could use?**

## Problem Statement

We need to decide:

1. **What responsibilities should be centralized?** (Evaluation?  Storage? Validation? All of the above?)
2. **How do we achieve cross-language compatibility?** (WASM? Native bindings? Code generation?)
3. **What performance trade-offs are acceptable?** (Native speed vs. maintenance burden)
4. **How do we ensure this doesn't fragment the ecosystem further?** (Adoption strategy, migration path)

## Requirements

### Functional Requirements

1. **In-process evaluation**: Must run within the application process (not as separate service)
2. **Single source of truth**: One canonical implementation of core functionality
3. **Cross-language compatibility**: Must work with all major SDK languages
4. **Feature parity**: Must support all current flagd evaluation features
5. **Behavioral consistency**: Same inputs must produce identical outputs across all SDKs
6. **Specification compliance**: Must adhere to OpenFeature and flagd specifications
7. **Extensibility**: Must support future feature additions without SDK changes

### Non-Functional Requirements

1. **Performance**: Acceptable overhead compared to native implementations
2. **Memory efficiency**:  Minimal memory footprint
3. **Security**: Proper isolation and sandboxing
4. **Debugging**: Ability to troubleshoot issues across language boundaries
5. **Distribution**: Easy to package and distribute with SDKs
6. **Versioning**: Clear versioning and compatibility guarantees

### Non-Requirements

1. **Complete SDK replacement**: SDKs still handle OpenFeature API, provider patterns, and language-specific concerns
2. **Zero performance overhead**: Some overhead is acceptable for maintenance benefits
3. **Immediate adoption**: Migration can be gradual and experimental

## Scope:  What Should the Unified Evaluator Be Responsible For?

This is a critical architectural decision.  We need to determine the boundaries of the unified evaluator.

### Option A: Evaluation Logic Only (Minimal Scope)

**Responsibilities**:

- JsonLogic evaluation
- Targeting rule processing
- Fractional evaluation algorithms
- Default variant resolution

**Not Included**:

- Flag storage/caching (SDK-specific)
- Schema validation (SDK-specific)
- Sync protocols (SDK-specific)
- Error handling translation (SDK-specific)

**Pros**:

- Smallest surface area
- Easier to integrate
- SDKs retain control over storage and caching strategies

**Cons**:

- Storage implementations still diverge
- Validation logic still duplicated
- Partial solution to the fragmentation problem

### Option B:  Evaluation + Storage (Medium Scope)

**Responsibilities**:

- All from Option A
- In-memory flag storage
- Flag state management
- Query/lookup interface

**Not Included**:

- Schema validation (happens before storage)
- Sync protocols (SDK-specific)
- Network communication

**Pros**:

- Consistent storage behavior across SDKs
- Unified caching strategies
- Better performance (no serialization for stored flags)

**Cons**:

- More complex integration
- Larger memory footprint
- Less SDK flexibility for storage optimization

### Option C: Evaluation + Storage + Validation (Maximal Scope)

**Responsibilities**:

- All from Option B
- JSON schema validation
- Configuration validation
- Flag conflict detection

**Not Included**:

- Sync protocols (SDK-specific)
- Network communication
- SDK-specific error translation

**Pros**:

- Complete consistency across SDKs
- Single source of truth for all core logic
- Validation errors consistent everywhere
- Maximum maintenance burden reduction

**Cons**:

- Largest integration complexity
- Most data crossing language boundaries
- Least SDK flexibility

### Recommendation:  Option C (Evaluation + Storage + Validation)

**Rationale**:  The current PoC implementation already includes storage and validation, demonstrating feasibility. The maintenance benefits of centralizing all core logic outweigh the integration complexity.  SDKs become thin adapters, focusing on OpenFeature API compliance and language-specific concerns.

**Evaluator Responsibilities** (Unified Component):

```
┌─────────────────────────────────────────────┐
│        Unified Evaluator (Core Logic)       │
├─────────────────────────────────────────────┤
│ 1. Flag Storage & State Management          │
│    - In-memory flag store                   │
│    - Flag lifecycle management              │
│    - Query interface                        │
│                                             │
│ 2. JSON Schema Validation                   │
│    - Flag configuration validation          │
│    - Schema compliance checking             │
│    - Error reporting                        │
│                                             │
│ 3. Evaluation Logic                         │
│    - JsonLogic evaluation                   │
│    - Targeting rule processing              │
│    - Fractional evaluation                  │
│    - Context resolution                     │
│    - Default variant handling               │
└─────────────────────────────────────────────┘
```

**SDK Responsibilities** (Language-Specific Adapter):

```
┌─────────────────────────────────────────────┐
│           SDK (Language-Specific)           │
├─────────────────────────────────────────────┤
│ 1. OpenFeature API Implementation           │
│    - Provider interface                     │
│    - Client/evaluation API                  │
│    - Hooks and event system                 │
│                                             │
│ 2. Evaluator Integration                    │
│    - Runtime initialization                 │
│    - Type marshalling/unmarshalling         │
│    - Error translation                      │
│                                             │
│ 3. Sync Protocol Implementation             │
│    - gRPC, HTTP, file watching, etc.        │
│    - Network communication                  │
│    - Connection management                  │
│                                             │
│ 4. Language-Specific Concerns               │
│    - Threading/concurrency                  │
│    - Logging integration                    │
│    - Metrics/telemetry                      │
└─────────────────────────────────────────────┘
```

## Considered Options

### Option 1: Unified Rust-Based Evaluator with Multiple Compilation Targets (Proposed)

**Approach**: Implement evaluator in Rust, compile to multiple targets based on SDK needs:

- **WASM** for maximum portability and security
- **Native bindings** (Python, JavaScript/Node.js) for enhanced performance where available

**Why Rust?**

- Memory safety without garbage collection
- Excellent cross-compilation support (WASM, native bindings)
- Strong type system for correctness
- High performance
- Existing battle-tested implementation in Rust SDK
- Active ecosystem for embedded/WASM scenarios
- Can compile to Python native extensions (PyO3)
- Can compile to JavaScript/Node.js native modules (neon, napi-rs)

**Pros**:

- Single source of truth for all core logic
- Guaranteed behavioral consistency
- **Flexible compilation targets**:  Choose WASM for portability or native for performance
- WASM provides sandboxing and security when needed
- Native bindings offer near-native performance for Python and JavaScript
- Platform-independent WASM works everywhere as fallback
- Reduced maintenance burden (fix once, deploy everywhere)
- Strong isolation between SDK and evaluator (WASM)
- **In-process execution** - no IPC/network overhead
- Best of both worlds: portability (WASM) + performance (native where available)

**Cons**:

- WASM has performance overhead (WASM boundary crossing, serialization)
- Native bindings require platform-specific builds (but only for specific SDKs)
- Additional dependency on WASM runtime or native module
- Debugging across language boundaries more complex
- Learning curve for maintainers unfamiliar with WASM/Rust
- Memory must be copied across boundaries

**Compilation Target Strategy**:

| SDK Language | Primary Target | Fallback | Rationale |
|--------------|---------------|----------|-----------|
| Python | Native (PyO3) | WASM | PyO3 provides excellent performance, native Python objects |
| JavaScript/Node.js | Native (napi-rs) | WASM | Native modules common in Node.js ecosystem |
| Java | WASM | - | JNI complexity not worth performance gain |
| . NET | WASM | - | P/Invoke complexity not worth performance gain |
| Go | WASM | - | cgo complexity and cross-compilation challenges |
| PHP | WASM | - | FFI less mature |
| Rust | Direct library | - | No FFI needed, use Rust library directly |

**Technical Considerations**:

- WASM runtimes available for all major languages
- Mature tooling (wasmtime, wasmer)
- PyO3 provides seamless Rust-Python integration
- napi-rs provides seamless Rust-Node.js integration
- Growing ecosystem and community support
- WASI for future extensibility

### Option 2: Shared Native Library (C/C++ or Rust via FFI)

**Approach**: Implement evaluator as native library, use FFI bindings in each SDK.

**Pros**:

- Near-native performance
- No serialization overhead
- Mature FFI tooling in most languages
- Direct memory access
- **In-process execution**

**Cons**:

- Platform-specific binaries (Windows/Linux/macOS, x86/ARM)
- Complex build and distribution (must compile for each platform)
- Memory safety concerns (especially with C/C++)
- Less isolation (shared memory space)
- FFI bindings are language-specific and error-prone
- Rust FFI still requires careful memory management
- Cross-compilation complexity
- Dependency hell (system libraries, ABI compatibility)

**Verdict**: More complex distribution, platform-specific builds for all SDKs, and less safety than WASM.

### Option 3: Code Generation from Specification

**Approach**: Define evaluation logic in DSL/specification, generate code for each language.

**Pros**:

- Native performance
- No runtime dependencies
- Traditional debugging
- **In-process execution**

**Cons**:

- Generated code may not be idiomatic
- Still requires testing in each language
- Code generation tooling adds complexity
- Behavioral drift still possible (platform-specific number handling, edge cases)
- Does not solve storage/validation duplication
- Generator becomes a complex, critical dependency

**Verdict**: Doesn't fully solve the fragmentation problem; complexity shifted to generator.

### Option 4: Reference Implementation with Conformance Tests

**Approach**:  Maintain language-specific implementations with comprehensive cross-SDK tests.

**Pros**:

- No architectural changes
- Maximum flexibility per SDK
- Native performance
- **In-process execution**

**Cons**:

- Maintenance burden continues to grow
- Behavioral drift still possible
- Does not address the root problem
- Testing alone cannot prevent implementation divergence
- Still requires implementing in 10+ languages

**Verdict**: Status quo - the problem we're trying to solve.

## Proposal

### Architecture Overview

```
┌────────────────────────────────────────────────────────────────┐
│                     Application Layer                          │
│                  (User's Application Code)                     │
└────────────────────────────────┬───────────────────────────────┘
                                 │
                                 ▼
┌────────────────────────────────────────────────────────────────┐
│                  OpenFeature SDK (Language-Specific)           │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  OpenFeature API (Client, Provider Interface, Hooks)     │  │
│  └────────────────────────────┬─────────────────────────────┘  │
│                               │                                │
│  ┌────────────────────────────▼─────────────────────────────┐  │
│  │         Language-Specific Adapter Layer                  │  │
│  │  • Type marshalling (SDK types ↔ JSON/binary/native)    │  │
│  │  • Error translation (evaluator errors → SDK exceptions)│  │
│  │  • Runtime initialization (WASM or native module)       │  │
│  └────────────────────────────┬─────────────────────────────┘  │
└─────────────────────────────────┼────────────────────────────────┘
                                  │
        ════════��═════════════════╧══════════════════════════
          Language Boundary (WASM or Native FFI, In-Process)
        ══════════════════════════╤══════════════════════════
                                  │
┌─────────────────────────────────▼────────────────────────────────┐
│    Runtime Layer (WASM Runtime OR Native Module Loader)          │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │         Unified Evaluator (Rust → WASM/Native)             │  │
│  │                                                            │  │
│  │  ┌──────────────────────────────────────────────────────┐ │  │
│  │  │  Storage Layer                                       │ │  │
│  │  │  • In-memory flag store                             │ │  │
│  │  │  • State management                                 │ │  │
│  │  │  • Query interface                                  │ │  │
│  │  └──────────────────────────────────────────────────────┘ │  │
│  │                                                            │  │
│  │  ┌──────────────────────────────────────────────────────┐ │  │
│  │  │  Validation Layer                                    │ │  │
│  │  │  • JSON schema validation                           │ │  │
│  │  │  • Configuration validation                         │ │  │
│  │  │  • Conflict detection                               │ │  │
│  │  └──────────────────────────────────────────────────────┘ │  │
│  │                                                            │  │
│  │  ┌──────────────────────────────────────────────────────┐ │  │
│  │  │  Evaluation Engine                                   │ │  │
│  │  │  • JsonLogic evaluation                             │ │  │
│  │  │  • Targeting rules                                  │ │  │
│  │  │  • Fractional evaluation                            │ │  │
│  │  │  • Context resolution                               │ │  │
│  │  └──────────────────────────────────────────────────────┘ │  │
│  └────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
           │
           │  ALL IN SAME PROCESS - No network/IPC overhead
           │
```

### Core Principle:  Separation of Concerns

**Unified Evaluator (Rust)** - Core business logic:

- Flag storage, validation, and evaluation
- Specification compliance
- Behavioral correctness

**SDK Layer (Language-Specific)** - Integration and adaptation:

- OpenFeature API compliance
- Language idioms and patterns
- Sync protocol implementation
- Telemetry and observability

### Interface Design

The evaluator exposes a minimal, stable interface:

```rust
// Conceptual interface (not final API)

// Initialize evaluator instance
fn create_evaluator() -> EvaluatorHandle;

// Store/update flags (with validation)
fn upsert_flags(handle: EvaluatorHandle, config: &str) -> Result<()>;

// Evaluate a flag
fn evaluate_flag(
    handle: EvaluatorHandle,
    flag_key: &str,
    context: &str,          // JSON evaluation context
) -> Result<EvaluationResult>;

// Query stored flags
fn list_flags(handle: EvaluatorHandle) -> Result<Vec<String>>;

// Remove flags
fn remove_flags(handle: EvaluatorHandle, keys: &[String]) -> Result<()>;

// Clean up
fn destroy_evaluator(handle: EvaluatorHandle);
```

**Design Principles**:

- Simple, stable API surface
- JSON for complex data in WASM (maximizes compatibility)
- Native objects for Python/JavaScript (maximizes performance)
- Clear error propagation
- Minimal state exposure
- Future-proof (easy to add new functions)

### Implementation Technology:  Rust with Multiple Compilation Targets

**Why Rust?**

1. **Safety**: Memory safety without GC prevents entire classes of bugs
2. **Performance**:  Zero-cost abstractions, minimal runtime overhead
3. **Multi-Target**: Best-in-class tooling for WASM, PyO3, napi-rs
4. **Existing Code**: Rust SDK already has battle-tested evaluation logic
5. **Correctness**: Strong type system catches errors at compile time
6. **Ecosystem**: Rich library ecosystem (serde, json-logic-rs, etc.)

**Compilation Targets**:

1. **WASM (WebAssembly)**:
   - Portability: Runs on any platform with a WASM runtime (no recompilation)
   - Security: Sandboxed execution with capability-based security
   - Performance:  Near-native speed (typically 1. 5-2x native in modern runtimes)
   - Language Support:  WASM runtimes exist for all major languages
   - In-Process: No IPC or network overhead
   - Deterministic: Same bytecode runs identically everywhere

2. **PyO3 (Python Native Extension)**:
   - Performance: Native speed, no serialization overhead
   - Integration: Direct Python object access
   - Distribution: Wheels for major platforms
   - Ecosystem: Well-established in Python community

3. **napi-rs (Node.js Native Module)**:
   - Performance: Native speed, no serialization overhead
   - Integration: Direct JavaScript object access
   - Distribution: Pre-built binaries for major platforms
   - Ecosystem: Standard approach for performance-critical Node.js modules

### Distribution Strategy

**Multi-Target Distribution**:

1. **WASM Binary** (universal fallback):
   - Publish versioned `.wasm` binary to GitHub Releases
   - Typical size: ~500KB (compressed)
   - Works on all platforms

2. **Python Wheels** (PyO3):
   - Publish to PyPI as platform-specific wheels
   - Pre-built for:  Linux (x86_64, ARM64), macOS (x86_64, ARM64), Windows (x86_64)
   - Falls back to WASM if platform not supported

3. **Node.js Native Modules** (napi-rs):
   - Publish to npm with pre-built binaries
   - Pre-built for: Linux (x86_64, ARM64), macOS (x86_64, ARM64), Windows (x86_64)
   - Falls back to WASM if platform not supported

**Versioning**:

- Semantic versioning (MAJOR.MINOR.PATCH)
- SDKs declare compatible evaluator version range
- Breaking changes require SDK updates
- All targets built from same source, same version

## Performance Considerations

### Expected Performance by Target

| Target | Overhead vs Pure Native | Use Case |
|--------|------------------------|----------|
| PyO3 (Python) | ~1.1-1.3x | Performance-critical Python applications |
| napi-rs (Node.js) | ~1.1-1.3x | Performance-critical JavaScript applications |
| WASM (all languages) | ~3-6x | Universal compatibility, security isolation |

### WASM Performance

Based on initial benchmarks:

| Scenario | Expected Overhead | Impact |
|----------|------------------|---------|
| Simple flags (no targeting) | ~3-4µs baseline | Negligible |
| Moderate targeting | 3-5x slower | Acceptable for most use cases |
| Complex targeting (large context) | 4-6x slower | Serialization dominates |

### Native Bindings Performance

Python (PyO3) and JavaScript (napi-rs) benchmarks show:

- Minimal overhead (~1.1-1.3x vs pure native)
- No serialization overhead for simple types
- Direct object access eliminates copies

### Context

- Typical web request: 10-100ms
- Flag evaluation: 4-20µs (WASM) vs 1-4µs (native) vs 1. 2-5µs (PyO3/napi-rs)
- Overhead:  <0.1% of total request time

### Optimization Opportunities

1. **Binary serialization**: Replace JSON with MessagePack/CBOR (WASM targets)
2. **Context caching**: Avoid re-serializing unchanged context
3. **Batch evaluation**: Evaluate multiple flags per call
4. **Short-circuit paths**: Bypass complex logic for trivial cases
5. **Native object handling**: Use PyO3/napi-rs for zero-copy access (Python/JS)

## Migration and Adoption Path

### Experimental Rollout Strategy

The unified evaluator will be introduced as an **experimental, opt-in feature** to minimize risk and gather real-world feedback before making it the default.

#### Phase 1: Experimental Release

**Goal**: Gather real-world feedback with minimal risk

**SDK Integration Pattern**:

```java
// Java example - opt-in via configuration
FlagdOptions options = FlagdOptions.builder()
    .useWasmEvaluator(true)  // Experimental:  defaults to false
    .build();

FlagdProvider provider = new FlagdProvider(options);
```

```python
# Python example - can use native or WASM
provider = FlagdProvider(
    experimental_unified_evaluator=True,  # Defaults to False
    evaluator_backend="native"  # or "wasm", auto-detects best available
)
```

```javascript
// JavaScript/Node.js example
const provider = new FlagdProvider({
  experimentalUnifiedEvaluator:  true,  // Defaults to false
  evaluatorBackend: 'native'  // or 'wasm', auto-detects best available
});
```

**Characteristics**:

- **Opt-in only**: Users must explicitly enable (defaults to existing implementation)
- **Auto-detection**: Python/JS SDKs automatically choose native or WASM based on availability
- **Feature flag**: Easy to disable if issues arise
- **Parallel implementations**: Both old and new code paths maintained
- **Clear documentation**: Marked as "experimental" in all docs
- **Feedback channels**: Clear path for users to report issues
- **Metrics/telemetry**: Track adoption and performance in the wild

#### Phase 2: Community Validation

- Gather feedback on architecture and scope
- Refine interface based on SDK integration experience
- Address performance and complexity concerns
- Establish governance model
- Expand to more SDKs (Go, .NET, PHP, etc.)

#### Phase 3: Stabilization

- Comprehensive test suite (unit, integration, conformance)
- Performance optimization based on real-world data
- Production-grade error handling
- Security audit
- Documentation and examples
- Remove "experimental" label

#### Phase 4: Default Implementation

**Make unified evaluator the default** (opt-out instead of opt-in)

```java
// Unified evaluator now default
FlagdProvider provider = new FlagdProvider();  // Uses unified evaluator

// Opt-out if needed (for migration period)
FlagdOptions options = FlagdOptions.builder()
    .useLegacyEvaluator(true)  // Temporary compatibility
    .build();
```

```python
# Python:  Native binding by default if available, WASM fallback
provider = FlagdProvider()  # Auto-selects best backend

# Explicit control if needed
provider = FlagdProvider(evaluator_backend="wasm")  # Force WASM
```

#### Phase 5: Ecosystem Standard

- All SDKs use unified evaluator by default
- Native implementations deprecated and removed
- Specification changes drive evaluator updates
- Single implementation for testing and validation

### Rollback Strategy

At any phase, if critical issues arise:

1. **Immediate**: Users disable via configuration flag
2. **SDK-level**: SDK maintainers can disable by default in patch release
3. **Evaluator-level**: Roll back to previous evaluator version
4. **Complete**: Pause migration, return to native implementations

## Consequences

### The Good

1. **Unified behavior**: Guaranteed consistency across all SDKs
2. **Reduced maintenance**: Fix once, benefit everywhere (10+ SDKs)
3. **Faster innovation**: New features implemented once, available immediately
4. **Better testing**: Concentrate testing efforts on single implementation
5. **Quality improvement**: Single implementation gets more scrutiny
6. **Specification compliance**: Easier to stay aligned with OpenFeature spec
7. **Onboarding**: New SDK languages easier to bootstrap
8. **Security**:  WASM sandboxing provides isolation when needed
9. **Cross-platform**: Works on any platform (desktop, mobile, embedded, web)
10. **Low-risk adoption**: Experimental rollout allows gradual validation
11. **Performance flexibility**: Native bindings for Python/JS, WASM for others
12. **Best of both worlds**: Portability + performance where it matters most

### The Bad

1. **WASM performance overhead**: 3-6x slower for complex evaluations (when not using native)
2. **Additional dependencies**:  WASM runtime or native module loader required
3. **Debugging complexity**: Cross-language boundary issues harder to debug
4. **Serialization cost**: Data marshalling adds overhead (WASM targets)
5. **Learning curve**: Maintainers need WASM/Rust knowledge
6. **Binary size**: Adds ~500KB (WASM) or platform-specific binaries (native)
7. **Build complexity**: Must build for multiple targets
8. **Migration effort**: SDKs must refactor to adopt evaluator
9. **Dual maintenance**: During experimental phase, maintain both implementations

### The Ugly

1. **Adoption uncertainty**: Not all SDK maintainers may adopt
2. **Fragmentation risk**: Could create "unified" vs "native" SDK split during transition
3. **Performance perception**: Users may overweight benchmark numbers
4. **Debugging tools**:  WASM debugging still less mature than native
5. **Dependency risk**: WASM runtime bugs affect entire ecosystem
6. **Breaking changes**: Interface changes require coordinated SDK updates
7. **Governance challenges**: Who decides evaluator direction?
8. **Platform coverage**: Native bindings not available for all platforms initially

### Risk Mitigation

1. **Performance**:  Provide native bindings for performance-critical languages (Python, JS)
2. **Adoption**: Make adoption optional and gradual with experimental phase
3. **Debugging**:  Invest in logging, error reporting, and tooling
4. **Governance**:  Establish clear decision-making process
5. **Migration**: Provide clear guides and support
6. **Testing**: Comprehensive test suite before promotion to default
7. **Fragmentation**: Clear communication and support during transition
8. **Rollback**: Easy opt-out at every phase
9. **Platform coverage**:  WASM fallback ensures universal compatibility

## Governance and Ownership

### Repository Ownership

- **Repository**: `open-feature-forking/flagd-evaluator` (consider moving to `open-feature/` org)
- **Maintenance**: OpenFeature community maintainers
- **Decision-making**: RFC process for significant changes

### Authorship and Attribution

The evaluator logic is derived from the existing Rust SDK flagd provider implementation. This work:

- Acknowledges original authors and contributors
- Preserves design decisions and architectural choices
- Builds upon battle-tested, validated logic
- Extends for cross-language use

### Contribution Model

- Standard OpenFeature contribution guidelines
- RFC process for interface changes
- Performance regression testing
- Semantic versioning for releases

## References

### Related Work

- **flagd-evaluator repository**: <https://github.com/open-feature-forking/flagd-evaluator>
- **flagd-evaluator v0.1.1 release**: <https://github.com/open-feature-forking/flagd-evaluator/releases/tag/v0.1.1>
- **Python SDK PoC**: <https://github.com/open-feature/python-sdk-contrib/pull/328>
- **Java SDK PoC**: <https://github.com/open-feature/java-sdk-contrib/pull/1672>
- **Rust SDK PoC**: <https://github.com/open-feature/rust-sdk-contrib/pull/94>
- **Original Issue**: <https://github.com/open-feature/flagd/issues/1842>

### Specifications

- **OpenFeature Specification**: <https://openfeature.dev/specification>
- **flagd Flag Configuration**: <https://flagd.dev/reference/flag-definitions/>
- **WASM Specification**: <https://webassembly.org/specs/>

### Technologies

- **PyO3**: <https://pyo3.rs/> - Rust bindings for Python
- **napi-rs**: <https://napi.rs/> - Rust bindings for Node.js
- **wasmtime**: <https://wasmtime.dev/> - Fast and secure WASM runtime
- **wasmer**: <https://wasmer.io/> - Universal WASM runtime

### Prior Art

- **WebAssembly System Interface (WASI)**: Component model for portable modules
- **Envoy WASM filters**: Envoy Proxy uses WASM for extensibility (in-process)
- **Figma plugins**: Desktop app uses WASM for sandboxed plugins (in-process)
- **Fastly Compute@Edge**: Edge computing platform powered by WASM
- **Shopify Scripts**: Ruby VMs replaced with WASM for performance and safety
- **Ruff**: Python linter written in Rust, distributed as PyO3 extension
- **swc**: JavaScript/TypeScript compiler written in Rust, distributed as napi-rs module
