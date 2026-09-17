---
title: "02 — JetStream"
section: "02-JETSTREAM"
type: note
tags: [nats, jetstream, persistence, streams, consumers, index]
source: https://docs.nats.io/nats-concepts/jetstream
created: 2026-06-14
---

# 02 — JetStream

> **JetStream** is NATS's built-in persistence and streaming layer — a first-class superset of NATS Core that adds durable storage, replay, exactly-once delivery, and flow-controlled consumption, all without external brokers, Zookeeper, or a second binary. It ships inside `nats-server` and is enabled with a single config stanza.

---

## Mental Model: JetStream = Persistence Layer on NATS Core

```
┌──────────────────────────────────────────────────────────────┐
│                         NATS SERVER                          │
├──────────────────────────────────────────────────────────────┤
│  JetStream Layer  (persistence, replay, acks, ordering)      │
│  ┌────────────────────────────────────────────────────────┐  │
│  │  Streams (append-only logs bound to subjects)          │  │
│  │  ┌──────────────────────────────────────────────────┐  │  │
│  │  │  Consumers (cursors + ack state into a stream)   │  │  │
│  │  └──────────────────────────────────────────────────┘  │  │
│  └────────────────────────────────────────────────────────┘  │
├──────────────────────────────────────────────────────────────┤
│  NATS Core  (pub/sub, request/reply, queue groups)           │
│  — at-most-once, in-memory, zero persistence                 │
└──────────────────────────────────────────────────────────────┘
```

**Three primitives:**

| Primitive | What it is | Analogy |
|-----------|------------|---------|
| **Stream** | An append-only log capturing messages from one or more subjects | Kafka topic / log file |
| **Consumer** | A server-side cursor into a stream; tracks position + ack state | Consumer group in Kafka / cursor |
| **Message** | A unit of data with headers, payload, and a monotonic sequence number | Record in a log |

**Key insight:** A stream is a *view* of a subject namespace at a point in time. Multiple consumers can independently replay or tail the same stream with different starting positions, filters, and ack modes. The stream is the authority; consumers are ephemeral positions into it.

### JetStream vs. NATS Core — At a Glance

| Dimension | NATS Core | JetStream |
|-----------|-----------|-----------|
| Delivery guarantee | At-most-once (fire-and-forget) | At-least-once; optionally exactly-once |
| Persistence | None (in-memory only) | File or memory storage |
| Replay | No | Yes — replay from any sequence or time |
| Acknowledgements | None | Explicit, All, or None |
| Backpressure / flow control | None | Yes (push) / natural (pull) |
| Consumer state | None | Durable cursor tracked server-side |
| Ordering | Subject-level FIFO only | Sequence-guaranteed per consumer |
| Cluster HA | Gossip/route topology | RAFT replication per stream |

### Subject Binding: How Core and JetStream Coexist

When a JetStream stream is configured to capture subject `ORDERS.*`, **every NATS Core publication** to `ORDERS.new`, `ORDERS.fulfilled`, etc. is *simultaneously*:
1. Delivered to any live Core subscribers (at-most-once semantics), **and**
2. Persisted into the stream with a monotonic sequence number (at-least-once / durable)

This means there is **no code change** needed on the publisher side to gain persistence — add a stream and existing publishers benefit automatically.

---

## Section Contents

| File | Topic |
|------|-------|
| [[02-JETSTREAM/jetstream-overview\|jetstream-overview]] | What JetStream is, enabling it, Core vs. JetStream semantics |
| [[02-JETSTREAM/streams\|streams]] | Stream configuration deep-dive — all fields, subject binding, editing |
| [[02-JETSTREAM/consumers\|consumers]] | Push vs. Pull, durable/ephemeral, ack policies, deliver policies, flow control |
| [[02-JETSTREAM/retention-and-limits\|retention-and-limits]] | The three retention policies, discard policies, limits interaction, TTL |
| [[02-JETSTREAM/dedup-and-ordering\|dedup-and-ordering]] | `Nats-Msg-Id` deduplication, duplicate window, exactly-once publish |
| [[02-JETSTREAM/mirrors-and-sources\|mirrors-and-sources]] | Stream mirrors vs. sources, aggregation, cross-domain replication |
| [[02-JETSTREAM/jetstream-clustering\|jetstream-clustering]] | R1/R3/R5, RAFT, meta-group, placement, HA |

---

## Reading Order

```
jetstream-overview          ← Start here (concepts + enabling)
      ↓
streams                     ← Storage model, config fields
      ↓
consumers                   ← How data flows out to apps
      ↓
retention-and-limits        ← When messages are removed
      ↓
dedup-and-ordering          ← Exactly-once semantics
      ↓
mirrors-and-sources         ← Replication + aggregation
      ↓
jetstream-clustering        ← HA + RAFT
```

---

## Quick Start

```shell
# Enable JetStream (simplest — in-memory, no limits)
nats-server -js

# With file storage
nats-server -js -sd /data/nats

# Create a stream
nats stream add ORDERS --subjects "ORDERS.*" --storage file --retention limits --replicas 1

# Publish a message
nats pub ORDERS.new "order-id-1234"

# Check stream state
nats stream info ORDERS

# Create a pull consumer and fetch
nats consumer add ORDERS DISPATCH --pull --ack explicit --deliver all
nats consumer next ORDERS DISPATCH
```

---

## Sources

- https://docs.nats.io/nats-concepts/jetstream
- https://docs.nats.io/using-nats/jetstream
- https://github.com/nats-io/nats.docs/tree/master/nats-concepts/jetstream
- https://docs.nats.io/running-a-nats-service/configuration/resource_management/jetstream_config
