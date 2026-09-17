---
title: "When to Use NATS (and When Not To)"
section: "00-OVERVIEW"
type: note
tags: [nats, use-cases, architecture, decision-guide, anti-patterns]
source: https://docs.nats.io/nats-concepts/overview
created: 2026-06-14
---

# When to Use NATS (and When Not To)

See also: [[NATS/00-OVERVIEW/nats-vs-alternatives|NATS vs Alternatives]] · [[NATS/06-TOPOLOGY-SCALING/README|Topology & Scaling]]

---

## Ideal Use Cases

### 1. Microservice Communication
The canonical fit. Services publish/subscribe on subjects and use **request-reply** for RPC without service-discovery glue — the subject *is* the address. Add **queue groups** to load-balance N instances of a service automatically. One service can scale from 1→100 instances with zero caller changes.

### 2. Command & Control / Device Fabric
Fan a command out to many devices with one publish (`cmd.fleet.shutdown` → every subscriber, microseconds). Request-reply gets a confirmation from a single device (`{"ok": true}`). This is the OneMind HA bridge model: ~1–5ms round trips over the bus.

### 3. Edge & IoT
**Leaf nodes** extend the fabric to edge devices with **store-and-forward** during disconnects and JetStream catch-up on reconnect. **Native MQTT** (`:1883`) lets constrained devices speak MQTT while the backend uses native NATS. **Embedded NATS** (Go in-process) turns any binary into a node.

### 4. Event Streaming & Event Sourcing (moderate scale)
**JetStream streams** give append-only logs with replay, work-queue semantics, and exactly-once publish — event sourcing, audit trails, task queues, and durable event distribution without standing up Kafka.

### 5. Request-Reply RPC with Resilience
Location-transparent RPC: caller publishes to a subject with a reply inbox; any healthy instance answers. Get the **no-responders** signal (instant 503) when nothing is listening — no hanging timeouts. Scatter-gather collects multiple replies for fan-out queries.

### 6. Distributed KV & Configuration
**JetStream KV** provides a replicated key-value store with **watch** (live change notifications), **history**, **revisions**, **TTL**, and atomic compare-and-swap — perfect for live config, feature flags, presence/context (OneMind's `zeus-context` bucket), and coordination.

### 7. Large-Object Distribution
**JetStream Object Store** chunks and stores large blobs (firmware images, model files, media) with metadata and watch — distribute artifacts across the fabric without a separate object store for moderate sizes.

### 8. Real-Time Browser Apps
**WebSocket** support (`:8080`) lets browser apps connect *directly* to the bus — live dashboards, collaborative UIs, real-time feeds with no backend polling proxy (OneMind's live bus viewer + Flutter `NatsMonitorService`).

### 9. Multi-Tenant Platforms
**Accounts** give isolated subject spaces with explicit **exports/imports** for controlled cross-tenant sharing — true multi-tenancy in one server.

---

## Anti-Patterns — When NATS Is the Wrong Choice

| Situation | Why NATS struggles | Use instead |
|-----------|--------------------|-------------|
| **Petabyte-scale durable analytics / data lake** | JetStream stores in file/memory (RAFT-replicated), not tiered/offloaded infinite storage | Kafka + tiered storage, Pulsar, or a lakehouse |
| **Rich stream-processing ecosystem needed** | NATS has no built-in SQL/stream-processing engine | Kafka (Streams/ksqlDB/Flink) |
| **Complex broker-side routing (AMQP exchanges, DLQ policies)** | NATS keeps the server thin; routing is subject-pattern, not broker logic | RabbitMQ |
| **Strict per-message broker-enforced ordering across huge partitions** | Ordering is per-stream/per-subject; not Kafka-partition semantics | Kafka |
| **You only need point-to-point typed RPC to a known host** | NATS adds a bus you don't need | gRPC |
| **Giant individual messages (hundreds of MB+ as messages)** | `max_payload` defaults to ~1MB (raise cautiously); huge blobs belong in Object Store, not as messages | JetStream Object Store, or S3-class storage |
| **You need full MQTT 5.0 advanced features at massive device fleets** | NATS MQTT supports QoS 0/1 (not the full 5.0 surface) | EMQX / HiveMQ for those specific needs |

---

## Decision Checklist

Choose **NATS** if most of these are true:
- [ ] Many services/devices that need to talk across a network, possibly edge/geo
- [ ] You value **operational simplicity** (one binary, no Zookeeper/Erlang) highly
- [ ] You need **more than one pattern** (pub/sub + req-reply + streams + KV)
- [ ] **Low latency** (sub-ms–low-ms) matters more than infinite retention
- [ ] You want **self-hostable, Apache-2.0, sovereignty-friendly** infrastructure
- [ ] **Edge / IoT / disconnected** operation is in scope (leaf nodes, MQTT)

Reconsider if:
- [ ] Your core need is **massive durable log analytics** with a processing ecosystem
- [ ] You need **broker-side complex routing** semantics (AMQP)
- [ ] You require **tiered/offloaded infinite storage**

---

## OneMind OS Verdict

NATS is the ratified **central nervous system** of OneMind OS precisely because the use-case profile matches: hundreds of heterogeneous nodes (Macs, Pis, robots, drones, vehicles, phones, SaaS webhooks, AI agents) needing low-latency, location-transparent, multi-pattern communication that survives disconnection and runs on a single self-hosted binary. The anti-patterns (data-lake analytics) are handled *beside* NATS by ClickHouse, not *by* the bus — the bus stays thin and fast. See the binding doctrine at `~/onemindos-core/docs/NATS-DOCTRINE.md`.

## Sources
- NATS Overview — https://docs.nats.io/nats-concepts/overview
- Compare NATS — https://docs.nats.io/nats-concepts/overview/compare-nats
- JetStream use cases — https://docs.nats.io/nats-concepts/jetstream
- Leaf nodes — https://docs.nats.io/running-a-nats-service/configuration/leafnodes
