# NATS can replace TimescaleDB — but only half the job

**NATS JetStream, KV Store, and Object Store can fully handle real-time event streaming, short-to-medium-term time-event storage, replay, distributed state management, and agent-to-agent messaging — eliminating the need for TimescaleDB in those domains.** Where NATS falls short is any use case requiring server-side analytical queries, aggregations across time windows, JOINs, or long-term compressed historical analysis. For the specific use cases outlined — IoT sensor events, distributed multi-location awareness, and OpenClaw agent messaging with no persistent database requirement — NATS alone covers **roughly 80% of needs**, with the remaining 20% (historical analytics) requiring either acceptance of client-side processing or a lightweight complement like DuckDB.

## JetStream turns a message broker into a persistent event store

JetStream transforms NATS from a fire-and-forget pub/sub system into a durable event log with capabilities that directly overlap with time-series storage. Every message stored in a JetStream stream receives a **monotonically increasing sequence number** and a **nanosecond-precision timestamp**, creating an immutable, ordered, time-indexed event log.

**Persistence and replication** operate at production grade. Streams persist to disk (file storage) or memory, with optional **S2 (Snappy) compression** achieving ~86% disk savings. Replication across **1–5 nodes** uses a NATS-optimized RAFT consensus algorithm providing immediate consistency. The storage overhead is minimal — a 5-byte payload message consumes just **39 bytes** on disk including sequence, timestamp, subject, and headers. One caveat: the default `sync_interval` of 2 minutes means non-replicated setups risk data loss on OS crashes (not process crashes). Production deployments should use replication or set `sync_interval: always`.

**Time-based retention** is the killer feature for replacing time-series storage. The `MaxAge` parameter accepts nanosecond-precision durations — set it to 24 hours, 30 days, or 1 year and messages automatically expire. This combines with `MaxBytes` and `MaxMsgs` limits, with **whichever threshold is hit first triggering cleanup**. Three retention modes exist: `LimitsPolicy` (default, respects configured limits), `WorkQueuePolicy` (messages deleted after acknowledgment), and `InterestPolicy` (messages deleted once all active consumers acknowledge). Since NATS 2.11, **per-message TTL** via the `Nats-TTL` header allows individual messages to expire independently of the stream's `MaxAge`.

**Time-based replay** is where JetStream directly substitutes for TimescaleDB's time-range queries — at least for sequential replay. Consumers support six delivery policies:

- **DeliverByStartTime** — replay from a specific timestamp (e.g., "give me everything from 1 hour ago")
- **DeliverByStartSequence** — replay from a specific sequence number
- **DeliverAll** — replay the entire stream from the beginning
- **DeliverLast / DeliverLastPerSubject** — get only the most recent message(s)
- **DeliverNew** — receive only future messages

The `ReplayPolicy` controls replay speed: `ReplayInstant` delivers as fast as the consumer can handle; `ReplayOriginal` replays at the original publication rate — useful for staging production traffic patterns. Combined with subject filtering on consumers, you can replay "all temperature sensor events from building-3 starting at 2 PM yesterday" efficiently.

## KV Store delivers a reactive distributed awareness layer

The NATS Key-Value Store is a **client-side abstraction built directly on JetStream streams**. Each KV bucket maps to a stream named `KV_<bucket>`, and each key maps to a subject `$KV.<bucket>.<key>`. This means KV inherits all of JetStream's persistence, replication, and retention capabilities while exposing a familiar key-value API.

For distributed state management and "awareness layer" use cases, KV provides everything needed. The **Watch** capability is the critical feature: calling `Watch("location1.>")` creates a push-based subscription that first delivers the current value for every matching key, signals "end of initial data," then streams real-time updates as they occur. Multiple watchers can observe the same keys independently — each gets its own ephemeral ordered consumer. In production deployments across **11 global regions**, KV operations achieved sub-60ms response times with connection latency under 10ms.

**TTL and history** provide temporal awareness without a database. Bucket-level TTL (mapping to stream `MaxAge`) auto-expires stale entries — ideal for heartbeat-based liveness detection where a device's key expires if not refreshed. Per-key TTL (NATS 2.11+) allows individual keys to expire at different rates. History depth from **1 to 64 revisions per key** lets you track recent state changes with full timestamps, and each entry includes the operation type (PUT/DEL/PURGE), revision number, and creation timestamp. **Compare-and-swap** via the `Update(key, value, revision)` operation provides optimistic concurrency control for safe concurrent writes.

One significant limitation: **TTL expiry does not currently emit watch events**. When a key expires due to TTL, watchers receive no notification. NATS 2.11's `SubjectDeleteMarkerTTL` partially addresses this by placing delete markers when `MaxAge` causes deletions, but this requires explicit configuration. This is a known gap with active community demand for a fix.

The KV multi-location pattern works as follows: a primary cluster holds the authoritative bucket with RAFT replication, while leaf nodes or mirror clusters in other regions maintain read replicas via JetStream Sources and Mirrors. Services at each location watch their local replica for near-instant reads; writes route to the primary cluster. This architecture maps directly to distributed multi-location awareness state requirements.

## Native scheduling arrived in NATS 2.12

**NATS Server 2.12 introduced native scheduled and delayed message delivery** via ADR-51, eliminating a long-standing gap. Streams opt in with `AllowMsgSchedules: true`, and messages carry scheduling instructions through headers:

- **`Nats-Schedule: @at 2026-03-01T12:00:00Z`** — one-time delayed delivery at a specific time
- **`Nats-Schedule: @every 5m`** or cron expressions like `0 30 14 * * *` — recurring delivery (design finalized in ADR-51; single delayed delivery shipped in 2.12, with recurring/cron support potentially arriving in subsequent releases)
- **`Nats-Schedule-Target`** — the subject where the scheduled message will be delivered
- **`Nats-Schedule-Source`** — reads the last message on a given subject and republishes it (powerful for sensor data sampling/downsampling)

For IoT sensor sampling, this is particularly elegant: publish a schedule message that reads the latest value from `sensors.cnc.temperature` every 5 minutes and republishes it to `sensors.sampled.cnc.temperature` — achieving basic downsampling without any external process. Combined with JetStream's **subject hierarchies** (e.g., `sensors.<type>.<location>.<device>`), wildcard subscriptions, and per-subject message limits (`MaxMsgsPerSubject: 1` for last-value cache semantics), NATS provides a complete real-time IoT data pipeline.

Additional scheduling patterns include **NAK-with-delay** (consumer receives a message, determines it's not time yet, and requeues with `NakWithDelay(remaining_time)`) and external schedulers publishing to NATS via libraries like `robfig/cron/v3`.

## What TimescaleDB provides that NATS fundamentally cannot

The gap between NATS and TimescaleDB is not about storage — it's about **querying**. NATS has no server-side query engine. This creates hard limitations in five categories:

**Aggregation is impossible server-side.** "Average temperature per floor per hour for the last 30 days" requires TimescaleDB's `time_bucket()` function and continuous aggregates. In NATS, you'd consume every message and compute client-side — orders of magnitude slower at scale. TimescaleDB's continuous aggregates are incrementally refreshed materialized views that auto-update, support hierarchical rollups (hourly → daily → weekly), and combine pre-materialized data with the latest raw data.

**Complex filtering requires full replay.** "All sensor readings between 2 PM and 4 PM where temperature > 30" is a trivial indexed query in TimescaleDB with automatic chunk pruning. In NATS, you'd use `DeliverByStartTime` to get messages from 2 PM onward, then deserialize and filter every single message client-side until 4 PM. There is **no way to filter by payload content server-side**.

**JOINs don't exist.** Real-world time-series analysis almost always correlates sensor data with device metadata, location data, or business context. TimescaleDB leverages full PostgreSQL JOINs. NATS can source multiple streams together but creates an interleaved log with no correlation capability.

**Compression favors TimescaleDB for long-term storage.** TimescaleDB's columnstore conversion achieves **90–98% compression** with data tiering to S3 for bottomless storage. NATS S2 compression achieves ~86% but only kicks in after filling 8MB blocks, and critically, **stream limits still use uncompressed message sizes** — so compression doesn't help you store more messages within configured limits. One user reported 48 hours of messages consuming 5.5GB in NATS versus 1GB in TimescaleDB.

**Ad-hoc analytical queries require a query language.** TimescaleDB provides full PostgreSQL SQL with 100+ hyperfunctions. An emerging but immature workaround is the **DuckDB `nats_js` extension** that enables SQL querying of JetStream streams — supporting sequence/timestamp range queries, subject filtering, and JSON extraction — but this has roughly 749 downloads/week and 5 GitHub stars as of early 2026.

## Applying this to IoT, awareness state, and agent messaging

For the three target use cases, here is how NATS alone performs:

**IoT/sensor time events** — NATS handles ingestion, real-time alerting, short-to-medium-term replay, and edge-to-cloud data flow excellently. The leaf node architecture enables store-and-forward at disconnected edge sites. Core NATS handles **~7.5 million messages/second** on modest hardware. JetStream with `MaxAge` retention provides time-windowed event storage. The gap appears only when you need analytical queries over historical data — trend analysis, anomaly detection across time windows, or dashboard aggregations.

**Distributed multi-location awareness state** — NATS KV is purpose-built for this. The Watch mechanism provides exactly the reactive awareness pattern needed: each location maintains a local view of global state via mirrored KV buckets on leaf nodes, with sub-60ms update propagation. History (up to 64 revisions) covers recent state transitions. The `isolate_leafnode_interest` feature (NATS 2.12) prevents unnecessary cross-location traffic. For "what was the state yesterday at 3 PM?" queries beyond 64 revisions, you'd need an external store.

**Agent-to-agent messaging (OpenClaw agents)** — NATS is the natural choice. Request/reply for synchronous agent queries, pub/sub for broadcast awareness, queue groups for load-balanced task distribution, and KV watches for shared world state. Projects like **FastAgency** (NatsAdapter + NatsProvider) and **NATS MCP Server** already demonstrate AI agent communication over NATS. JetStream provides durable task queues with acknowledgment guarantees. The only limitation is searching message history by content — "all messages from agent X containing error keyword" requires sequential replay and client-side filtering.

## Conclusion

NATS JetStream, KV Store, and Object Store collectively provide a **single-binary infrastructure layer** that handles persistent event storage, time-based retention and replay, distributed state management, real-time change notifications, scheduled message delivery, and multi-location edge-to-cloud synchronization. For systems that primarily need real-time event flow, current state awareness, and time-windowed replay — and do not require SQL-style analytical queries — **NATS alone eliminates the need for TimescaleDB entirely**.

The decisive question is whether your system needs to *ask questions about* historical data or merely *replay* it. NATS can replay events from any point in time at any speed. It cannot aggregate, filter by payload content, JOIN across data sources, or compute statistical functions server-side. If your IoT/agent system needs only "show me the current state" and "replay what happened since X," NATS suffices. If it needs "what was the average temperature trend across all locations last month," TimescaleDB (or a lightweight alternative like DuckDB querying JetStream) remains necessary. The emerging DuckDB `nats_js` extension hints at a future where this gap narrows, but today it's the clear architectural boundary between what NATS can and cannot replace.