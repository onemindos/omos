---
title: "JetStream Streams"
section: "02-JETSTREAM"
type: note
tags: [nats, jetstream, streams, storage, retention, configuration, append-only-log]
source: https://docs.nats.io/nats-concepts/jetstream/streams
created: 2026-06-14
---

# JetStream Streams

> A **stream** is an append-only log that captures and retains messages published to one or more NATS subjects. Every message appended gets a monotonically increasing **sequence number**. Consumers are cursors into that log — the stream itself just stores; it does not push.

See also: [[02-JETSTREAM/jetstream-overview|JetStream Overview]] · [[02-JETSTREAM/consumers|Consumers]] · [[02-JETSTREAM/retention-and-limits|Retention & Limits]]

---

## Streams as Append-Only Logs

```
Subject:  ORDERS.new   ORDERS.fulfilled   ORDERS.new   ORDERS.cancelled
              │               │               │               │
              ▼               ▼               ▼               ▼
Stream ORDERS: [seq=1]──[seq=2]──[seq=3]──[seq=4]──► (always growing)
               msg1     msg2     msg3     msg4

Consumer DISPATCH: reads from seq=1, acks each message
Consumer AUDIT:    reads from seq=1, no acks (replay)
Consumer NEW_ONLY: filter_subject="ORDERS.new", reads seq=1, seq=3
```

Key properties:
- **Immutable writes** — messages are never modified in place; they are appended
- **Monotonic sequences** — every server restart preserves the sequence counter
- **Independent consumers** — each consumer maintains its own read position; they don't interfere with each other
- **Retention is stream-level** — when messages are evicted (by age, count, bytes, or interest), they are removed from the front of the log

---

## Creating a Stream

### Interactive CLI

```shell
nats stream add my_stream
? Subjects foo
? Storage file
? Replication 1
? Retention Policy Limits
? Discard Policy Old
? Stream Messages Limit -1
? Per Subject Messages Limit -1
? Total Stream Size -1
? Message TTL -1
? Max Message Size -1
? Duplicate tracking time window 2m0s
? Allow message Roll-ups No
? Allow message deletion Yes
? Purge subjects or the entire stream Yes
Stream my_stream was created
```

### CLI with Flags

```shell
nats str add ORDERS \
  --subjects "ORDERS.*" \
  --ack \
  --max-msgs=-1 \
  --max-bytes=-1 \
  --max-age=1y \
  --storage file \
  --retention limits \
  --max-msg-size=-1 \
  --discard old \
  --dupe-window="2m" \
  --replicas 1
```

### From JSON Config File

```shell
nats str add ORDERS --config orders.json
```

---

## Full Stream Configuration Reference

These are the **real field names** as used in the NATS API JSON and the `nats` CLI. All fields correspond to `StreamConfig` in the server API (nats-server v2.x).

| Field | JSON key | CLI flag | Type | Default | Description |
|-------|----------|----------|------|---------|-------------|
| **Name** | `name` | `--name` | `string` | *(required)* | Stream name. Alphanumeric + `-` + `_`. Must be unique per account. |
| **Subjects** | `subjects` | `--subjects` | `[]string` | *(required unless mirror)* | Subject list this stream captures. Supports wildcards (`ORDERS.*`, `EVENTS.>`). |
| **Storage** | `storage` | `--storage` | `"file"` \| `"memory"` | `"file"` | File = persisted to disk; Memory = RAM only (faster, lost on restart unless replicated). |
| **Retention** | `retention` | `--retention` | `"limits"` \| `"interest"` \| `"workqueue"` | `"limits"` | How messages are kept. See [[02-JETSTREAM/retention-and-limits\|Retention & Limits]]. |
| **Discard** | `discard` | `--discard` | `"old"` \| `"new"` | `"old"` | When a limit is hit: `old` drops the oldest message; `new` rejects the incoming message. |
| **MaxMsgs** | `max_msgs` | `--max-msgs` | `int64` | `-1` (unlimited) | Max total messages in the stream. When exceeded, `discard` policy applies. |
| **MaxMsgsPerSubject** | `max_msgs_per_subject` | `--max-msgs-per-subject` | `int64` | `-1` (unlimited) | Per-subject message limit (limits retention policy only). |
| **MaxBytes** | `max_bytes` | `--max-bytes` | `int64` | `-1` (unlimited) | Max total bytes stored. When exceeded, `discard` policy applies. |
| **MaxAge** | `max_age` | `--max-age` | `duration` (nanoseconds in API) | `0` (unlimited) | Messages older than this are removed. e.g. `86400000000000` = 24h. CLI accepts `1y`, `24h`, `2m0s`. |
| **MaxMsgSize** | `max_msg_size` | `--max-msg-size` | `int32` | `-1` (unlimited) | Max size of a single message in bytes. Messages larger are rejected. |
| **MaxConsumers** | `max_consumers` | `--max-consumers` | `int` | `-1` (unlimited) | Max number of consumers that can be created. |
| **Replicas** | `num_replicas` | `--replicas` | `int` | `1` | Number of RAFT replicas (1, 3, or 5). Requires a cluster with ≥ N servers. |
| **DuplicateWindow** | `duplicate_window` | `--dupe-window` | `duration` | `2m` (120s) | Sliding window for `Nats-Msg-Id` deduplication. Max bounded by server `duplicate_window` limit. |
| **NoAck** | `no_ack` | `--no-ack` | `bool` | `false` | If `true`, the stream does not send a `PubAck` on publish. Useful for high-throughput fire-and-forget ingestion. |
| **AllowRollupHdrs** | `allow_rollup_hdrs` | `--allow-rollup` | `bool` | `false` | Allow the `Nats-Rollup` header for subject-level or stream-level rollup (compact). |
| **DenyDelete** | `deny_delete` | `--deny-delete` | `bool` | `false` | Prevents deletion of individual messages from the stream. |
| **DenyPurge** | `deny_purge` | `--deny-purge` | `bool` | `false` | Prevents `stream purge` operations. |
| **Sealed** | `sealed` | `--sealed` | `bool` | `false` | Seals the stream — no new messages accepted, no config changes, no consumer creation. |
| **AllowDirect** | `allow_direct` | `--allow-direct` | `bool` | `false` | Allows direct `GetMsg` API calls that bypass the consumer machinery (lower latency reads). |
| **MirrorDirect** | `mirror_direct` | — | `bool` | `false` | For mirror streams: allow direct `GetMsg` calls to be served by any mirror replica. |
| **Compression** | `compression` | `--compression` | `"none"` \| `"s2"` | `"none"` | Storage compression algorithm. `s2` = Snappy/S2 (fast, ~30% space saving). |
| **Placement** | `placement` | `--tag` | `object` | — | Cluster placement: `cluster` name and/or `tags` list to target specific servers. |
| **SubjectTransform** | `subject_transform` | — | `object` | — | Transform subjects as messages enter the stream (`src` pattern → `dest` pattern). |
| **Republish** | `republish` | — | `object` | — | After storing a message, republish it to another subject. Non-blocking. |
| **Mirror** | `mirror` | `--mirror` | `object` | — | Mirror config — see [[02-JETSTREAM/mirrors-and-sources\|Mirrors & Sources]]. |
| **Sources** | `sources` | `--source` | `[]object` | — | Sources config — see [[02-JETSTREAM/mirrors-and-sources\|Mirrors & Sources]]. |

### Example: Full JSON Stream Config

```json
{
  "name": "ORDERS",
  "subjects": ["ORDERS.*"],
  "storage": "file",
  "retention": "limits",
  "discard": "old",
  "max_msgs": -1,
  "max_msgs_per_subject": -1,
  "max_bytes": -1,
  "max_age": 86400000000000,
  "max_msg_size": -1,
  "max_consumers": -1,
  "num_replicas": 3,
  "duplicate_window": 120000000000,
  "no_ack": false,
  "allow_rollup_hdrs": false,
  "deny_delete": false,
  "deny_purge": false,
  "sealed": false,
  "allow_direct": true,
  "compression": "s2"
}
```

---

## Storage: File vs. Memory

| Aspect | `file` | `memory` |
|--------|--------|----------|
| Persistence | ✅ Survives restarts | ❌ Lost on restart (unless replicated) |
| Speed | Good (NVMe recommended) | ⚡ Fastest |
| Capacity | Limited by disk | Limited by RAM |
| Use case | Production / durable streams | Caching, ephemeral pipelines, test |
| Clustering | Replicated via RAFT to peers' **file** stores | Replicated in-memory across peers |

> **Tip:** Memory streams with `num_replicas: 3` ARE HA — loss of one server doesn't lose data because two peers still hold the in-memory copy. But if all replicas restart simultaneously, data is lost.

---

## Subject Binding

A stream captures subjects matching its `subjects` list. Rules:

1. **Wildcards work** — `ORDERS.*` captures `ORDERS.new`, `ORDERS.fulfilled`, etc.
2. **One subject → one stream** — the same subject cannot be captured by two different streams in the same account
3. **Publishing to a captured subject is transparent** — the publisher uses regular `nats.Publish()`, no JetStream-specific API needed
4. **Subject mapping / transforms** — a `subject_transform` can rewrite the stored subject (e.g., strip a routing prefix before storage)

```
NATS Core publish:  nats pub ORDERS.new "payload"
                           ↓
JetStream captures: ORDERS stream stores seq=42, subject="ORDERS.new"
                           ↓  
Core subscribers:   all live ORDERS.new subscribers get the message too
```

### Subject Overlap & Filtering

```shell
# Capture all ORDERS subjects
nats stream add ORDERS --subjects "ORDERS.*"

# A consumer can further filter within the stream
nats consumer add ORDERS NEW_ONLY --filter-subject "ORDERS.new" ...
```

---

## Editing Streams

Most stream fields can be updated without interruption. Subjects, retention, limits, replicas can all be changed live.

```shell
# Edit interactively
nats stream edit ORDERS

# Edit specific fields via flags
nats stream edit ORDERS --subjects "ORDERS.*" --max-age 2y

# Copy a stream (duplicate config with modifications)
nats str cp ORDERS ARCHIVE --subjects "ORDERS_ARCHIVE.*" --max-age 2y

# Inspect config as JSON
nats str info ORDERS -j | jq .config

# Seal a stream (no more writes)
nats stream seal ORDERS

# Purge all messages
nats stream purge ORDERS

# Purge a specific subject within the stream  
nats stream purge ORDERS --subject ORDERS.cancelled

# Delete a specific message by sequence
nats stream rmsg ORDERS 42
```

### Config Update via JSON

```shell
# Export current config, edit, re-apply
nats str info ORDERS -j | jq .config > orders.json
# Edit orders.json
nats str edit ORDERS --config orders.json
```

---

## Stream Lifecycle

```
nats stream add     → stream created, subjects bound
nats stream info    → inspect state (msgs, bytes, consumers, cluster)
nats stream edit    → update config (live, non-disruptive for most fields)
nats stream purge   → delete all messages (stream persists, consumers reset)
nats stream seal    → make read-only forever
nats stream rm      → delete stream entirely (all messages + consumers gone)
```

---

## Checking Stream State

```shell
nats stream info ORDERS
```

```
Information for Stream ORDERS created 2026-06-14T10:00:00Z

             Subjects: ORDERS.*
             Replicas: 3
              Storage: File
            Retention: Limits
              Max Msgs: unlimited
         Max Msg Size: unlimited
            Max Bytes: unlimited
              Max Age: 24h0m0s
  Maximum Consumers: unlimited
  Duplicate Window: 2m0s

State:
             Messages: 1,234
                Bytes: 2.5 MiB
        First Seq: 1 @ 2026-06-14T10:00:01Z
         Last Seq: 1234 @ 2026-06-14T12:30:00Z
       Num Subjects: 3
       Num Deleted: 0

Cluster Information:
                 Name: nats-cluster
               Leader: server-1
              Replicas:
                server-2 (current, seen 0.3s ago)
                server-3 (current, seen 0.1s ago)
```

---

## Subject Transform & Republish

### Subject Transform (on ingest)

Rewrites the stored subject as messages enter the stream. Useful to strip routing prefixes or normalize subjects:

```json
{
  "name": "orders",
  "subjects": ["orders.local.*"],
  "subject_transform": {
    "src": "orders.local.*",
    "dest": "orders.{{wildcard(1)}}"
  }
}
```

Message published to `orders.local.new` is stored under `orders.new`.

### Republish (post-store)

After a message is stored, republish a copy to another subject for live consumers who don't want to use a consumer:

```json
{
  "republish": {
    "src": "orders.*",
    "dest": "orders.trace.{{wildcard(1)}}"
  }
}
```

---

## Sources

- https://docs.nats.io/nats-concepts/jetstream/streams
- https://docs.nats.io/running-a-nats-service/nats_admin/jetstream_admin/streams
- https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jetstream_admin/streams.md
- https://docs.nats.io/nats-concepts/subject_mapping
- https://docs.nats.io/using-nats/developing-with-nats/js/streams
