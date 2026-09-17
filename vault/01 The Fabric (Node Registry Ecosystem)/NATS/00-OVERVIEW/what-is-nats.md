---
title: "What is NATS?"
section: "00-OVERVIEW"
type: note
tags:
  - nats
  - overview
  - messaging
  - pub-sub
  - jetstream
  - fundamentals
source: "https://docs.nats.io/nats-concepts/what-is-nats"
created: 2026-06-14
---

# What is NATS?

> *"NATS is a connective technology that powers modern distributed systems."*
> — docs.nats.io

NATS is **message-oriented middleware** — an infrastructure layer that allows software applications, services, and devices to exchange data segmented into **messages**. It is not a database, not a log aggregator, and not an API gateway. It is the **connective tissue** of distributed systems: addressing, discovery, and message exchange in one coherent system.

See also: [[history-and-philosophy]], [[nats-vs-alternatives]], [[when-to-use-nats]]

---

## The Problem NATS Solves

Traditional inter-service communication technologies (HTTP, gRPC) share three architectural constraints:

1. **Location coupling** — services address each other by hostname/IP + port
2. **1:1 topology** — each call has exactly one sender and one receiver
3. **Fragmented security** — each protocol requires its own auth/TLS model

As microservices, functions, and IoT devices proliferate — especially at the edge — these assumptions break down:

- A service moving between pods, regions, or edge nodes changes its IP
- Adding fanout or load-balancing requires bolting on additional infrastructure (load balancers, service meshes, API gateways, sidecars)
- Managing N different security models across N protocols is operationally expensive

**NATS replaces all of this with a single connective fabric.** Subjects (not IP addresses) route messages. M:N communication is the default. Security is uniform and baked in.

---

## The Unified Messaging Model

NATS unifies multiple messaging patterns under a single protocol and server process. Where competing tools require separate systems for pub/sub, queuing, and streaming, NATS supports all of these natively.

### Core NATS (at-most-once delivery)

Core NATS operates on a **fire-and-forget publish-subscribe model**:

```
Publisher ──► NATS Server ──► Subscriber(s)
```

- Messages are held **in memory only** — never written directly to disk in Core mode
- If no subscriber is active when a message arrives, it is **dropped**
- This matches the semantics of TCP/IP: best-effort, low-overhead
- Suitable when you can tolerate occasional loss (metrics, telemetry, live events)

#### Pub/Sub (Fan-out)

One publisher sends to a **subject**; all matching subscribers receive the message simultaneously.

```
nats pub "sensors.temp.room-1" "22.5"
# All subscribers of "sensors.temp.room-1" (or "sensors.>") receive it
```

#### Queue Groups (Load Balancing)

Multiple subscribers share the same **queue group name**. The server delivers each message to **exactly one** member of the group, randomly selected. This is NATS's built-in horizontal scaling for work distribution.

```go
// Workers join queue group "order-processors"
nc.QueueSubscribe("orders.new", "order-processors", handler)
```

No external load balancer required. Adding or removing workers is zero-config.

#### Request / Reply (RPC)

NATS has first-class support for the request/reply pattern. The requestor publishes to a subject with a unique **reply-to** subject; the first available responder publishes back to that inbox.

```go
// Requester
msg, err := nc.Request("auth.check", payload, 2*time.Second)

// Responder
nc.Subscribe("auth.check", func(m *nats.Msg) {
    m.Respond([]byte("ok"))
})
```

This is **not** HTTP — it works across any topology (single server, cluster, leaf node, geo-distributed supercluster) with zero configuration changes. The "location transparency" principle: callers don't know or care where the responder lives.

---

### NATS JetStream (at-least-once / exactly-once delivery)

JetStream is NATS's **built-in persistence engine**, introduced in NATS Server v2.2 (2021). It is compiled into `nats-server` and enabled per-deployment — no separate process, no ZooKeeper, no external storage dependency.

JetStream adds:

| Capability | Description |
|---|---|
| **Streams** | Durable, ordered, replayable message logs bound to one or more subjects |
| **Consumers** | Named read positions into a stream (push or pull delivery) |
| **Key/Value Store** | Atomic get/put/delete map with history and TTL, built on streams |
| **Object Store** | Chunked file transfer and storage, built on streams |
| **At-least-once delivery** | Ack-based redelivery until acknowledged |
| **Exactly-once delivery** | Deduplication window + idempotent publish |
| **Temporal decoupling** | Consumers can start replaying from any point in the stream's history |
| **Flow control** | Decoupled publisher/consumer rates without head-of-line blocking |

JetStream uses a **NATS-optimized RAFT algorithm** for distributed consensus, providing **linearizable writes** and **serializable reads** across a cluster.

```
Stream "ORDERS" bound to subject "orders.>"
  ├── Consumer "audit-log"    (push, deliver-all, no ack needed)
  ├── Consumer "billing"      (pull, deliver-new, ack required)
  └── Consumer "analytics"    (pull, start-from 7 days ago, ack required)
```

The same message set can be consumed independently by multiple consumers at different speeds, with different starting positions, and different delivery guarantees — all without modifying the publisher.

---

## Subject-Based Addressing

Everything in NATS flows through **subjects** — UTF-8 strings with `.` as a hierarchy separator. There are no "topics," "queues," or "channels" at the addressing layer; subjects are the universal primitive.

### Wildcards

| Wildcard | Meaning | Example |
|---|---|---|
| `*` | One token | `sensors.*.temp` matches `sensors.room1.temp` |
| `>` | One or more tokens (must be terminal) | `sensors.>` matches all subjects under `sensors.` |

```
sensors.building-a.floor-2.room-14.temp
         ^          ^       ^       ^
         *          *       *       *  ← each * matches one token
sensors.>  ← matches everything under sensors
```

Subjects are resolved **server-side** with O(n) matching in the subscription tree — no client-side filtering needed.

---

## The Single-Binary / Zero-Dependency Philosophy

> *"The NATS server is a single binary with no external dependencies, consuming minimal memory and CPU resources. The server typically uses less than 20 MB of RAM."*
> — nats.io/about

This is a first-class design goal, not a side effect:

- **No ZooKeeper/etcd** — NATS clustering uses a native RAFT implementation
- **No separate persistence process** — JetStream is compiled into `nats-server`
- **No message schema registry** — NATS is payload-agnostic (raw bytes)
- **No sidecar required** — the server runs standalone everywhere
- **Single binary deployment** — `nats-server` is a ~20 MB Go binary

```bash
# That's the entire installation:
nats-server

# Or with Docker:
docker run -p 4222:4222 nats:latest

# Or on a Raspberry Pi (ARM64):
nats-server -js  # JetStream enabled
```

This philosophy makes NATS viable across the full computing spectrum — from a cloud hyperscaler cluster to a Raspberry Pi on a factory floor to a resource-constrained IoT gateway.

---

## Core Statistics Reference

These figures are sourced from official NATS documentation and nats.io marketing claims. Third-party benchmarks (2025, 4 vCPU / 8 GB VPS configurations) are noted separately.

| Metric | Value | Source / Notes |
|---|---|---|
| **Server binary size** | ~20 MB | nats.io/about (official claim) |
| **Idle RAM footprint** | <20 MB | nats.io documentation |
| **Core throughput** | Millions of msgs/sec | nats.io/about (vendor claim; hardware-dependent) |
| **Core latency** | Sub-millisecond | nats.io/about (vendor claim) |
| **JetStream throughput** | 200K–400K msgs/sec (persisted) | Third-party VPS benchmark, onidel.com (2025) |
| **JetStream latency** | Sub-ms (memory) / 1–5 ms (file) | Third-party VPS benchmark, onidel.com (2025) |
| **License** | Apache 2.0 | github.com/nats-io/nats-server |
| **Server language** | Go (99.7% of codebase) | GitHub |
| **Latest stable** | v2.14.2 (2026-06-02) | github.com/nats-io/nats-server/releases |
| **GitHub stars** | ~20K | github.com/nats-io/nats-server (Jun 2026) |
| **Official client languages** | 40+ | nats.io/download |
| **External dependencies** | Zero (server) | No ZooKeeper, no etcd, no Kafka |
| **Default client port** | 4222 | Standard |
| **Monitoring HTTP port** | 8222 | Standard |
| **Cluster route port** | 6222 | Standard |

> ⚠️ **Benchmark caveat:** Throughput figures are highly workload-dependent (message size, persistence mode, replication factor, hardware). The nats.io "millions of msgs/sec" claim refers to in-memory Core NATS on optimized hardware. Always benchmark your specific workload.

---

## Quality of Service (QoS) Levels

| Mode | Provider | Delivery | Persistence | Notes |
|---|---|---|---|---|
| **At-most-once** | Core NATS | Fire-and-forget | In-memory only | Dropped if no active subscriber; lowest overhead |
| **At-least-once** | JetStream | Ack-based, redelivered until acked | Disk (or memory) | Durable streams, consumer offsets |
| **Exactly-once** | JetStream | Idempotent publish + ack | Disk | Uses `Nats-Msg-Id` dedup header + double-ack protocol |

---

## Client Ecosystem

NATS ships official client libraries maintained by the nats-io organization:

| Language | Repo | Notes |
|---|---|---|
| Go | nats.io/nats.go | Reference client; most feature-complete |
| Rust | nats.rs | High performance, async |
| JavaScript (Node + Browser) | nats.deno / nats.ws | WebSocket support for browsers |
| TypeScript (Deno) | nats.deno | |
| Python | nats.py | asyncio-based |
| Java | nats.java | |
| C# (.NET) | nats.net | NATS.Net package |
| C | cnats | Low-level, C99 |
| Ruby | nats-pure.rb | |
| Elixir | nats.ex | |
| CLI | nats-io/natscli | Full admin + bench + pub/sub CLI |

Community clients exist for 30+ additional languages and runtimes.

---

## The NATS URL Scheme

```
nats://[user:password@]host[:port]
tls://[user:password@]host[:port]   # TLS
ws://host[:port]                    # WebSocket
wss://host[:port]                   # WebSocket + TLS
```

Clients need only a URL (and optional credentials) to connect. No service discovery configuration, no DNS SRV records, no sidecar — the NATS protocol handles topology transparently after initial bootstrap.

---

## Protocol Basics

The NATS wire protocol is a **line-based text protocol** (not binary, not HTTP). It is intentionally simple:

```
# Server greeting
INFO {"server_id":"...","version":"2.14.2","go":"go1.26.3","host":"0.0.0.0","port":4222,...}

# Client connect
CONNECT {"verbose":false,"pedantic":false,"tls_required":false,...}

# Publish
PUB orders.new 12
{"id":"001"}

# Subscribe
SUB orders.new 1

# Server delivers
MSG orders.new 1 12
{"id":"001"}

# Keepalive
PING
PONG
```

This simplicity means NATS clients are easy to implement (40+ languages) and the protocol overhead is extremely low.

---

## See Also

- [[history-and-philosophy]] — Origin story, CNCF, Synadia, design pillars
- [[nats-vs-alternatives]] — How NATS compares to Kafka, RabbitMQ, Pulsar, gRPC
- [[when-to-use-nats]] — Use cases and anti-patterns

---

## Sources

- What is NATS (official docs): <https://docs.nats.io/nats-concepts/what-is-nats>
- NATS Overview (official docs): <https://docs.nats.io/nats-concepts/overview>
- Core NATS concepts: <https://docs.nats.io/nats-concepts/core-nats>
- JetStream concepts: <https://docs.nats.io/nats-concepts/jetstream>
- About NATS (nats.io): <https://nats.io/about/>
- NATS server GitHub (stats): <https://github.com/nats-io/nats-server>
- NATS 2025 Benchmark (onidel.com): <https://onidel.com/blog/nats-jetstream-rabbitmq-kafka-2025-benchmarks>
