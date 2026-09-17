---
title: "JetStream Overview"
section: "02-JETSTREAM"
type: note
tags: [nats, jetstream, overview, persistence, exactly-once, at-least-once, enabling]
source: https://docs.nats.io/nats-concepts/jetstream
created: 2026-06-14
---

# JetStream Overview

> JetStream is the persistence and streaming subsystem built into `nats-server`. It upgrades NATS Core's fire-and-forget messaging with durable storage, replay, acknowledgements, deduplication, and flow-controlled delivery — all in the same binary, same port, same subject space.

See also: [[02-JETSTREAM/README|JetStream section index]] · [[02-JETSTREAM/streams|Streams]] · [[02-JETSTREAM/consumers|Consumers]]

---

## What JetStream Is

JetStream was introduced in NATS server **v2.2.0** (2021) as a response to the real-world need for:

1. **Message durability** — don't lose a message if a subscriber is offline or restarts
2. **Replay** — consume messages from the beginning, a given offset, or a point in time
3. **Back-pressure** — producers should not overwhelm slow consumers
4. **Exactly-once delivery** — prevent duplicate processing in critical workflows
5. **Work queues** — distribute work across a pool of competing consumers with reliable tracking

These are the features that drove organisations to Kafka, Pulsar, or RabbitMQ alongside NATS Core. JetStream collapses all of them into `nats-server` itself.

### Why Not Use NATS Core for These Needs?

| Scenario | NATS Core | JetStream |
|----------|-----------|-----------|
| Consumer offline when message sent | ❌ Message lost | ✅ Stored; delivered when consumer reconnects |
| Replay all messages since yesterday | ❌ Impossible | ✅ `DeliverPolicy: by_start_time` |
| Competing consumers, no duplicate processing | ❌ Not guaranteed | ✅ WorkQueue retention + explicit acks |
| Audit trail / event sourcing | ❌ No storage | ✅ Append-only log, queryable by seq |
| Fan-out to 1,000 slow consumers | ❌ Head-of-line blocking | ✅ Each consumer has independent flow control |
| Publisher rate > consumer rate | ❌ Messages dropped | ✅ Stored; consumer pulls at its own rate |

---

## Delivery Semantics

### At-Most-Once (NATS Core)

NATS Core pub/sub delivers a message to all current subscribers and immediately discards it. If no subscriber is present, the message is gone. If a subscriber crashes mid-processing, the message is gone.

### At-Least-Once (JetStream Default)

JetStream stores messages durably. A consumer must **acknowledge** each message. Until acknowledged:
- The server tracks the message as "pending"
- After `AckWait` timeout, the server redelivers it
- Up to `MaxDeliver` attempts

This guarantees at-least-once delivery: the message will be processed *at least* once. With crashes or network partitions, it may be processed more than once (duplicates are possible).

```
Publisher → Stream (stored) → Consumer delivers → App processes → App ACKs
                                      ↑                               |
                               If no ACK within AckWait ←────────────┘
                               server redelivers (up to MaxDeliver)
```

### Exactly-Once (JetStream with Deduplication)

NATS JetStream supports **exactly-once** publish by combining:
1. **`Nats-Msg-Id` header** on the publisher side — the server deduplicates within the `duplicate_window`
2. **Double-ack protocol** on the consumer side — consumer can verify the ack itself was received

This is the `AckKind: AckTerm` / double-ack pattern. For most use cases, idempotent application logic + at-least-once is sufficient and simpler.

See [[02-JETSTREAM/dedup-and-ordering|Deduplication & Ordering]] for the full mechanism.

---

## How JetStream Relates to NATS Core

### The Unified Subject Space

JetStream does **not** use a separate port, protocol, or connection. A JetStream stream is configured to *capture* messages published to a set of subjects. Any NATS Core publisher writing to those subjects has its messages automatically persisted — transparently, with no code change.

```
nats pub ORDERS.new "payload"   ← This is a NATS Core publish
                                   simultaneously:
                                   → delivered to live Core subscribers (at-most-once)
                                   → captured by ORDERS stream (durable)
```

### JetStream API

JetStream has its own API accessed over NATS Core messages to reserved subjects:

```
$JS.API.*          — management API (stream CRUD, consumer CRUD, info)
$JS.ACK.*          — acknowledgement subjects
$JS.FC.*           — flow control subjects  
$KV.*              — Key-Value layer (built on JetStream)
$OBJ.*             — Object Store layer (built on JetStream)
```

Clients interact with JetStream by publishing JSON-encoded requests to these subjects and receiving JSON responses — the exact same publish/subscribe mechanism as everything else.

---

## Enabling JetStream

### Option 1: CLI Flag (Development)

```shell
# Enable with default settings (file storage in ./jetstream)
nats-server -js

# Enable with custom storage directory
nats-server -js -sd /data/nats
```

### Option 2: Configuration File (Production)

```hcl
# js.conf
jetstream {
    store_dir: "/data/nats"
    max_memory_store: 1073741824    # 1 GiB memory store limit
    max_file_store: 10737418240     # 10 GiB file store limit
}
```

```shell
nats-server -c js.conf
```

### Option 3: Per-Account JetStream (Multi-tenant)

```hcl
accounts {
    A {
        jetstream {}          # JetStream enabled for account A
    }
    B {
        users: [{user: b, password: b}]
        # No jetstream block — account B cannot use JetStream
    }
}
```

### Server-Level JetStream Limits

```hcl
jetstream {
    store_dir: nats
    max_file_store: 10G
    limits {
        max_ack_pending: 10000
        duplicate_window: 600s    # max dedup window any stream can set
    }
    # Sync every message to disk (slower but safest)
    sync_interval: always
}
```

### Verify JetStream is Active

```shell
nats server check jetstream
nats account info               # shows JetStream limits + usage
```

---

## JetStream Architecture at a Glance

```
NATS Server (JetStream enabled)
│
├── Meta-Group (RAFT consensus for stream/consumer config)
│   └── Meta-Leader (one server, elected; processes API requests)
│
├── Stream: ORDERS  (R3 — 3 RAFT replicas)
│   ├── Stream-Leader (one server handles writes)
│   ├── Replica 1
│   ├── Replica 2
│   ├── Consumer: DISPATCH  (pull, durable, explicit ack)
│   └── Consumer: AUDIT     (push, durable, ack none)
│
└── Stream: METRICS  (R1 — single copy)
    └── Consumer: REALTIME  (push, ephemeral)
```

Key facts:
- Each stream has an independently elected RAFT leader
- The meta-group leader coordinates all `$JS.API.*` requests
- Streams without replicas (`num_replicas: 1`) have no RAFT overhead — they're just a write-ahead log on one server

---

## Enabling JetStream: Production Checklist

```
☐ Configure store_dir (SSD preferred)
☐ Set max_file_store based on available disk (leave headroom)
☐ Set max_memory_store if using memory streams
☐ Plan replica counts: R1 for non-critical, R3 for HA, R5 for critical
☐ Configure server_tags for stream placement
☐ Enable monitoring (:8222/jsz) for observability
☐ Set sync_interval appropriately (default vs. always)
```

---

## Sources

- https://docs.nats.io/nats-concepts/jetstream
- https://docs.nats.io/using-nats/jetstream/getting_started
- https://docs.nats.io/running-a-nats-service/configuration/resource_management/jetstream_config
- https://github.com/nats-io/nats.docs/blob/master/nats-concepts/jetstream/README.md
- https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/getting_started/README.md
