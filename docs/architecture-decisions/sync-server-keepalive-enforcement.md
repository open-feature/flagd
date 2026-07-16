---
# Valid statuses: draft | proposed | rejected | accepted | superseded
status: proposed
author: Jason Benedicic (@jabenedicic)
created: 2026-07-14
updated: 2026-07-14
---
# gRPC keepalive enforcement on the flag-sync server

The flag-sync gRPC server currently sets no keepalive enforcement policy, so it
inherits grpc-go's defaults, which are tuned for generic request/response RPC
rather than a long-lived streaming sync. Clients that ping to keep the sync
stream alive can be disconnected with `GOAWAY ENHANCE_YOUR_CALM`. This ADR
proposes setting a streaming-appropriate enforcement policy on the server, with
sane defaults and operator-facing overrides.
Tracking issue: [#1998](https://github.com/open-feature/flagd/issues/1998).

## Background

flagd's in-process mode has providers consume a long-lived `SyncFlags` gRPC
stream from the flag-sync server. gRPC clients keep such streams healthy by
sending periodic keepalive pings, especially during idle periods when no flag
changes are flowing.

gRPC servers enforce a *keepalive enforcement policy* that decides how often a
client is allowed to ping. grpc-go's default is `MinTime: 5m` with
`PermitWithoutStream: false`: a client may ping at most once every five minutes,
and only while it has an active stream. A client that exceeds this — pings more
often, or pings during a gap with no active RPC — accrues "ping strikes", and
after a small number the server sends `GOAWAY` with `ENHANCE_YOUR_CALM` /
`too_many_pings` and closes the connection. Those defaults suit short unary RPC
traffic; they are a poor fit for a service whose primary job is a long-lived,
often-idle stream.

The flag-sync server is constructed with a bare `grpc.NewServer(...)` and no
enforcement policy, so it inherits those defaults. Provider SDKs choose their
own keepalive cadence, and several ping far more frequently than every five
minutes to detect dead connections quickly. The Go provider is the clearest
example — it pings every 30s with `PermitWithoutStream: true` and exposes no
override — so it is disconnected roughly every 60–90s and must reconnect. This
is not specific to one provider: any client, in any language, whose keepalive
interval is shorter than five minutes is exposed to the same teardown, and there
is no server-side setting today to widen the tolerance.

## Requirements

* The flag-sync server must tolerate the keepalive cadences that provider SDKs
  use in practice, without emitting `GOAWAY ENHANCE_YOUR_CALM`.
* The fix must apply to every client language, not one provider.
* The policy must be operator-configurable (CLI flag and environment variable),
  so deployments can tune it without a code change.
* The default must not regress any client that works against flagd today.

## Considered Options

* Set a fixed, streaming-appropriate enforcement policy on the server, exposed
  via CLI flags / environment variables (chosen).
* Set a fixed enforcement policy on the server with no configuration surface.
* Change each language provider's client keepalive settings instead of the
  server.
* Leave the server on grpc-go's defaults (status quo).

## Proposal

Set a `KeepaliveEnforcementPolicy` on the flag-sync gRPC server, configurable
through two new options with defaults chosen to match how providers actually
keep the stream alive:

* `--keep-alive-min-time` (env `FLAGD_KEEP_ALIVE_MIN_TIME`, duration, default
  `30s`) — the minimum interval the server permits between client keepalive
  pings.
* `--keep-alive-permit-without-stream` (env
  `FLAGD_KEEP_ALIVE_PERMIT_WITHOUT_STREAM`, bool, default `true`) — whether the
  server permits pings while the client holds no active stream.

The `30s` default matches the tightest keepalive cadence a shipped provider uses
today; because a ping sent every 30s arrives at the server `30s + RTT` apart, it
clears a 30s minimum with margin. `PermitWithoutStream: true` is required
because providers ping during idle gaps when, transiently, no stream is active.
Both defaults are strictly more permissive than grpc-go's `5m` /
`PermitWithoutStream: false`, so any client that works against flagd today
continues to work; clients that were being disconnected stop being disconnected.

Choosing the server as the point of control is deliberate. Keepalive enforcement
is a property of the server — it is the party that decides how often a client
may ping — so the server is the correct and complete place to fix a
too-aggressive default. The alternative, changing keepalive in each provider,
would require coordinated releases across every provider repository, would not
help third-party or future clients, and would leave the server rejecting
reasonable keepalives out of the box. One server-side change covers every client
in every language.

### API changes

Two new CLI flags on `flagd start`, with matching `FLAGD_`-prefixed environment
variables, documented in the auto-generated CLI reference. No changes to the
gRPC or flag-configuration APIs. No wire-format or schema changes.

### Consequences

* Good, because a single server-side change resolves `ENHANCE_YOUR_CALM`
  disconnects for every provider language at once, at the layer that owns the
  policy.
* Good, because the defaults are more permissive than grpc-go's, so no
  currently-working client regresses.
* Good, because deployments — including operator-managed flagd — can tune the
  policy through environment variables or arguments on the flagd container,
  giving estate-wide, cross-language keepalive tuning from one place, without
  editing any gRPC client in any workload.
* Bad, because it adds two configuration options to the flagd surface that
  operators may need to understand.
* Bad, because a misconfigured very-low `--keep-alive-min-time` could let a
  misbehaving client ping more aggressively than intended; the default is
  conservative and this is opt-in.

### Open questions

* Should the OpenFeature Operator expose these as first-class `Flagd` CRD fields,
  or is passing them through container env/args sufficient? (Follow-up; out of
  scope for this change.)

## More Information

* grpc-go keepalive documentation:
  <https://github.com/grpc/grpc-go/blob/master/Documentation/keepalive.md>
* Keepalive enforcement and `ENHANCE_YOUR_CALM`:
  <https://github.com/grpc/grpc/blob/master/doc/keepalive.md>
