---
title: "00-OVERVIEW — Section Index"
section: "00-OVERVIEW"
type: note
tags:
  - nats
  - overview
  - index
  - messaging
source: "https://docs.nats.io/nats-concepts/overview"
created: 2026-06-14
---

# 00-OVERVIEW — NATS Knowledge Base Index

> **This section** is the foundation of the NATS knowledge base. Before diving into configuration, deployment, or advanced patterns, start here. These five notes form the conceptual bedrock you need to reason about NATS correctly.

---

## Section Map

```
00-OVERVIEW/
├── README.md                   ← You are here
├── what-is-nats.md             ← Core definition, unified model, key stats
├── history-and-philosophy.md   ← Origin story, design principles, CNCF/Synadia
├── nats-vs-alternatives.md     ← Kafka, RabbitMQ, Pulsar, Redis, MQTT, gRPC
└── when-to-use-nats.md         ← Ideal use cases and anti-patterns
```

---

## Notes in This Section

### 1. [[what-is-nats]]
**What is NATS?**
The authoritative definition: NATS as message-oriented middleware, the unified messaging model (pub/sub, request/reply, queue groups, streaming), the single-binary/zero-dependency philosophy, and a core stats reference table. Start here if you're new to NATS.

**Key topics:** connective technology, subjects-based addressing, Core NATS vs JetStream, QoS levels, client ecosystem, ~20 MB footprint, sub-millisecond latency, 40+ official clients.

---

### 2. [[history-and-philosophy]]
**History and Philosophy**
The full origin story: Derek Collison's path from TIBCO → Google → VMware → Cloud Foundry → Apcera → Synadia. The "Neural Autonomic Transport System" etymology. CNCF stewardship timeline including the 2025 trademark resolution. The four core design pillars that make NATS architecturally distinct.

**Key topics:** Derek Collison, Ruby origins (2011), Go rewrite, JetStream replacing STAN (2021), CNCF incubating project status, Synadia Communications, simplicity-first philosophy.

---

### 3. [[nats-vs-alternatives]]
**NATS vs Alternatives**
Expert-level comparison against Apache Kafka, RabbitMQ, Apache Pulsar, Redis Pub/Sub, MQTT brokers, and gRPC. Includes a comprehensive feature/performance comparison table and prose guidance on when each technology is the better choice.

**Key topics:** throughput benchmarks, latency profiles, persistence models, operational complexity, delivery guarantees, footprint comparison, use-case fit matrix.

---

### 4. [[when-to-use-nats]]
**When to Use NATS**
Decision guidance: the five ideal NATS use cases (microservices mesh, edge/IoT, command-and-control, request/reply RPC, event streaming) with architectural patterns for each. Anti-patterns section covers where NATS is the wrong tool and what to use instead.

**Key topics:** microservices, IoT telemetry, C2, RPC, streaming, edge/cloud continuum, total-order guarantees, SQL-style queries, batch ETL.

---

## Quick-Reference: NATS in 60 Seconds

| Property | Value |
|---|---|
| **Full name** | Neural Autonomic Transport System |
| **License** | Apache 2.0 |
| **Language** | Go (server) |
| **Latest stable** | v2.14.2 (2026-06-02) |
| **Default port** | 4222 (clients), 8222 (HTTP monitoring) |
| **Footprint** | ~20 MB RAM, single binary, zero external deps |
| **Core throughput** | Millions of msgs/sec per server (per nats.io) |
| **Latency** | Sub-millisecond (core); 1–5 ms (JetStream persistence) |
| **QoS modes** | At-most-once (Core); At-least-once / Exactly-once (JetStream) |
| **Client languages** | 40+ official + community clients |
| **Governance** | CNCF incubating project (Linux Foundation) |
| **Commercial support** | Synadia Communications (synadia.com) |

---

## How to Navigate This Knowledge Base

The NATS knowledge base is organized as follows:

```
NATS/
├── 00-OVERVIEW/        ← Start here (this section)
├── 01-CORE-CONCEPTS/   ← Subjects, subscriptions, pub/sub, queue groups
├── 02-JETSTREAM/       ← Persistence, streams, consumers, KV, object store
├── 03-SECURITY/        ← TLS, NKeys, JWT, decentralized auth
├── 04-CLUSTERING/      ← Clusters, superclusters, leaf nodes, gateways
├── 05-CLIENTS/         ← Per-language client guides
├── 06-OPERATIONS/      ← Deployment, monitoring, tuning, Kubernetes
├── 07-PATTERNS/        ← Architectural recipes and patterns
└── _assets/            ← Raw source docs and images
```

---

## Sources

- NATS official documentation: <https://docs.nats.io>
- NATS about page: <https://nats.io/about/>
- GitHub repository: <https://github.com/nats-io/nats-server>
- CNCF project listing: <https://www.cncf.io/projects/nats/>
