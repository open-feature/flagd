---
status: proposed
author: Faisal Rafiuddin
created: 2026-10-05
---

# Per-flag delta sync over FlagSyncService

## Details

|                        |                                      |
|------------------------|--------------------------------------|
| **Feature Name**      | Per-flag delta sync                  |
| **Type**              | enhancement                           |
| **Related components** | FlagSyncService, flagd sync clients |

## Background

Current `FlagSyncService.SyncFlags` always returns a full `flag_configuration`
snapshot per message. At scale, this is wasteful:

- Large projects commonly carry 2k+ flags. A single flag toggle causes the
  entire set to be re-serialized, pushed over the network, re-parsed on every
  connected client, and re-written into every client's store.
- For Booking.com's flag sync backend, the sync service fans out to thousands
  of clients across bare-metal hosts and containers. The current model of
  selector-flags full-config drives high memory and CPU use on the server,
  generates significant outbound network traffic, and creates stream
  backpressure on busy updates, sometimes resulting in OOMKills.
- Clients pay a symmetric cost on every message: re-parse, re-store, re-watch,
  and re-evaluate dependent caches.

A per-flag delta protocol would let the server emit `O(changed_flags)` rather
than `O(total_flags)` on each notification.

## Proposal

Extend `flagd.sync.v1` additively:

```proto
message SyncFlagsRequest {
  string provider_id = 1;
  string selector = 2 [deprecated = true];

  // Client advertises that it can parse the delta fields on SyncFlagsResponse
  // (is_delta + flag_ops) and the revision stamp (revision). Servers MUST NOT
  // emit delta messages to a client that leaves this unset.
  bool supports_deltas = 3;
}

message SyncFlagsResponse {
  string flag_configuration = 1;
  optional google.protobuf.Struct sync_context = 2;

  // When true, this message carries per-flag operations in flag_ops rather
  // than a full snapshot in flag_configuration. Only sent to clients that
  // advertised supports_deltas=true.
  bool is_delta = 3;

  // Per-flag operations. Only meaningful when is_delta is true.
  repeated FlagOp flag_ops = 4;

  // Strictly monotonically increasing revision for every message within a
  // stream. The first message of a stream MAY start at 1; 
  // every subsequent message MUST be exactly previous+1. The
  // client uses the revision to detect lost messages; on a gap, the client
  // reconnects to re-establish a fresh baseline. The server does NOT buffer
  // history or honor a since_revision parameter on the request.
  //
  // revision=0 means the server did not stamp this message (older servers,
  // pre-rollout). Clients MUST tolerate unstamped messages and skip gap
  // detection for the remainder of such a stream.
  uint64 revision = 5;
}

message FlagOp {
  enum Op {
    OP_UNSPECIFIED = 0;

    // Add or replace one flag by (flag_set_id, flag_key). Idempotent.
    OP_UPSERT = 1;

    // Remove one flag by (flag_set_id, flag_key). flag_json is ignored.
    OP_DELETE = 2;
  }

  Op op = 1;

  // Flag key.
  string flag_key = 2;

  // Optional flagSetId scope. Empty means the source-level (nil) flagSetId.
  optional string flag_set_id = 3;

  // Single-flag JSON body matching the flagd flag schema. Required for
  // OP_UPSERT, ignored for OP_DELETE.
  string flag_json = 4;
}
```

## Contract

1. **Capability handshake and single opt-in gate.** All new fields in
   `SyncFlagsResponse` (`is_delta`, `flag_ops`, `revision`) MUST only be set
   when the client advertised `SyncFlagsRequest.supports_deltas=true`. For
   any other client, the server MUST build responses that leave these fields
   unset — i.e., byte-identical to the pre-feature wire format.

2. **Initial snapshot first.** For each successful `SyncFlags` RPC invocation,
   the server MUST send a full snapshot as the first response before sending
   any delta. For clients that advertised `supports_deltas=true`, this
   baseline may start with revision 1.

3. **Reconnect re-establishes baseline.** Every reconnect gets a fresh initial
   snapshot and a fresh revision sequence. Deltas lost during a disconnect are
   superseded by the new baseline — the server does not back-fill.

4. **Op semantics.** `OP_UPSERT` is add-or-replace; `OP_DELETE` is remove.
   Clients MUST skip `OP_UNSPECIFIED` and any unknown op values to preserve
   forward compatibility.

5. **Batching.** A single response MAY carry N ops in `flag_ops`.

6. **Mixed state is forbidden.** Within one message, `is_delta=true` means
   `flag_configuration` is unused; `is_delta=false` means `flag_ops` is
   unused. Servers MUST NOT populate both.

7. **Revision stamping (opt-in clients only).** For clients that advertised
   `supports_deltas=true`, every `SyncFlagsResponse`, snapshot or delta,
   carries a revision. Revisions are per-stream and strictly monotonic:

   `revision_{n+1} = revision_n + 1`

   On a new stream (reconnect), the counter restarts at 1 and the first
   message is always a full snapshot. This is orthogonal to `is_delta`: the
   one-message-one-revision-bump invariant applies to snapshots and deltas
   alike.

8. **Gap detection and recovery.** A client that observes
   `revision_n != previous+1` MUST close the stream and reconnect.

## Config knob

flagd (and other delta-aware clients) exposes a per-source boolean
`supportsDeltas` on `SourceConfig`. When true:

- The client advertises `supports_deltas=true` on `SyncFlagsRequest`.
- The client processes incoming delta responses via a new per-flag apply
  path in the store.
- The client enables per-stream revision gap detection and reconnect.

Default is `false` for backward compatibility. Enabling it is opt-in per
source.

> **Note:** An existing boolean `incrementalUpdates` on `SourceConfig`
> (flagd core) today gates flagSetId-scoped snapshot replacement, which
> still transfers full per-flagSetId payloads. `supportsDeltas` reflects
> the broader semantics now that per-flag deltas and revision gap detection
> are also gated by it. `incrementalUpdates` stays as a deprecated alias
> for one release cycle.

## Compatibility

Additive proto change at new tag numbers (3, 4, 5). Older servers/clients
ignore unknown fields. Mixed-version fleets stay correct in every
combination:

- **Old server + new client (opted in):** The client sets
  `supports_deltas=true`, but the server does not know the field. The server
  emits snapshots only, with no new fields set. The client sees
  `revision=0` on every message, so gap detection is disabled and the delta
  path never fires. Snapshots work as before.

- **Old server + new client (not opted in):** Identical to the row above.

- **New server + old client:** The client does not set the field
  (zero-value `false`). The server reads it as "not delta-capable", gates all
  new fields off, and emits byte-identical wire format to the pre-feature
  server. The old client is behaviorally and byte-wise unaware the feature
  exists.

- **New server + new client (not opted in):** Same as the row above — the
  client is new but has not advertised capability, so the server behaves as
  if it were old. There is zero visible change for the client.

- **New server + new client (opted in):** Initial snapshot at `revision=1`,
  deltas flow thereafter with contiguous revisions. On a gap or disconnect,
  the client reconnects and starts over at `revision=1`.
