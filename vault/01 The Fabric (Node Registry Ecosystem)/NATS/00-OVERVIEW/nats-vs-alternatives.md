---
title: "NATS vs Alternatives — Kafka, RabbitMQ, Pulsar, Redis, MQTT, gRPC"
section: "00-OVERVIEW"
type: note
tags: [nats, comparison, kafka, rabbitmq, pulsar, redis, mqtt, grpc]
source: https://docs.nats.io/nats-concepts/overview/compare-nats
created: 2026-06-14
---

# NATS vs Alternatives

How NATS compares to the systems it is most often weighed against, and when to pick which. The honest summary: **NATS optimizes for operational simplicity, multi-pattern unification, and edge/geo reach.** Kafka optimizes for huge-volume durable log analytics; RabbitMQ for complex routing/broker semantics; Pulsar for tiered-storage multi-tenant streaming. NATS deliberately is *not* trying to be a data lake.

See also: [[NATS/00-OVERVIEW/when-to-use-nats|When to Use NATS]] · [[NATS/02-JETSTREAM/README|JetStream]]

---

## At-a-Glance Comparison

| Dimension | **NATS** (+ JetStream) | **Kafka** | **RabbitMQ** | **Pulsar** | **Redis Pub/Sub** | **MQTT brokers** |
|-----------|------------------------|-----------|--------------|------------|-------------------|------------------|
| **Core model** | Subject pub/sub + req-reply + streams + KV/OS | Partitioned commit log | Exchanges/queues (AMQP) | Topics + segments | Channel pub/sub | Topic pub/sub |
| **Persistence** | Optional (JetStream) | Always (log) | Optional (durable queues) | Always (BookKeeper) | None (in-mem) | Optional (retained/QoS) |
| **Delivery** | at-most-once (Core); at-least / exactly-once (JS) | at-least / exactly-once | at-most / at-least-once | at-least / effectively-once | at-most-once | QoS 0/1/2 |
| **Latency** | sub-100µs typical | low-ms | sub-ms–ms | low-ms | sub-ms | ms |
| **Throughput** | very high (~15M+ msg/s core) | very high (durable) | high | very high | very high | moderate |
| **Footprint** | **~20MB single Go binary** | JVM + **Zookeeper/KRaft** | Erlang VM | JVM + **BookKeeper + Zookeeper** | small | small–moderate |
| **External deps** | **none** | Zookeeper or KRaft | Erlang | ZK + BookKeeper | none | varies |
| **Ops complexity** | **low** | high | medium | **high** | low | low |
| **Geo / edge** | **native** (leaf nodes, superclusters, gateways) | MirrorMaker (bolt-on) | Federation/shovel | geo-replication | none | bridging |
| **Multi-tenancy** | **native accounts** | ACLs/quotas | vhosts | namespaces/tenants | DBs (weak) | ACLs |
| **Built-in KV / Object** | **yes** (JetStream KV + OS) | no (Kafka Streams/ksqlDB ext.) | no | no (state via Functions) | yes (it's Redis) | no |
| **Built-in req/reply** | **yes** | no | via RPC pattern | no | no | no |
| **Browser/IoT** | WebSocket + **native MQTT** | no | MQTT plugin | no | no | **native MQTT** |
| **License** | Apache-2.0 | Apache-2.0 | MPL-2.0 | Apache-2.0 | BSD/RSAL | varies |

> Numbers are order-of-magnitude and workload-dependent. Treat the *shape* of each system as the signal, not exact figures.

---

## vs Apache Kafka

**Kafka** is a distributed, partitioned, replicated **commit log** built for high-volume durable event streaming and stream processing (analytics, CDC, log aggregation, replayable event sourcing at massive scale).

- **Pick Kafka when:** you need very long retention of huge data volumes, a rich stream-processing ecosystem (Kafka Streams, ksqlDB, Connect, Flink), partition-ordered replay at terabyte scale, and you already have the ops muscle for JVM + Zookeeper/KRaft.
- **Pick NATS when:** you want messaging *and* request-reply *and* KV *and* edge reach from one ~20MB binary with no Zookeeper; you value operational simplicity and sub-ms latency over data-lake-scale retention.
- **Key contrast:** Kafka makes you run a cluster + coordination layer before "hello world." NATS runs a server with zero config. JetStream gives you Kafka-like streams/replay/work-queues, but JetStream is a *messaging persistence layer*, not a petabyte analytics store.

## vs RabbitMQ

**RabbitMQ** is a mature **AMQP broker** with sophisticated routing (direct, topic, fanout, headers exchanges), per-message TTLs, dead-letter queues, and broker-side delivery logic.

- **Pick RabbitMQ when:** you need complex broker-side routing topologies, AMQP/STOMP protocol compatibility, mature per-queue policies, and classic enterprise task-queue semantics.
- **Pick NATS when:** you want far simpler ops (no Erlang VM), higher raw throughput, native request-reply and service discovery, built-in KV/Object, and edge/geo topology without federation plugins.
- **Key contrast:** RabbitMQ puts routing intelligence in the *broker* (exchanges/bindings). NATS puts addressing in the *subject* and keeps the server thin — routing is subject-pattern matching, which is faster and simpler but less "smart" broker-side.

## vs Apache Pulsar

**Pulsar** is a multi-tenant streaming platform with a **two-layer architecture** (brokers + Apache BookKeeper for storage) and tiered/offloaded storage.

- **Pick Pulsar when:** you need multi-tenant streaming with tiered storage offload (e.g. to S3), geo-replication at scale, and per-tenant namespaces with strong isolation, and can operate ZK + BookKeeper + brokers.
- **Pick NATS when:** you want a *single* component, not a three-tier storage stack; lower operational burden; and the same multi-pattern unification NATS gives.
- **Key contrast:** Pulsar's storage/serving split is powerful but heavy. NATS + JetStream keeps storage *in* the server (file/memory, RAFT-replicated), trading infinite tiered retention for radical simplicity.

## vs Redis Pub/Sub (and Streams)

**Redis** is an in-memory data structure store; its **Pub/Sub** is fire-and-forget with no persistence, and **Redis Streams** add a lightweight log.

- **Pick Redis when:** you're already using Redis for caching/data structures and want incidental pub/sub or a simple stream in the same process.
- **Pick NATS when:** messaging is a first-class concern — you need delivery guarantees, queue groups, request-reply, durable streams, multi-node clustering for *messaging*, and edge reach. Redis Pub/Sub drops messages to disconnected subscribers; NATS Core has the same at-most-once caveat but JetStream fixes it properly.
- **Note:** NATS *also* ships a **KV store** (JetStream KV), so for the "I just need a distributed KV with watch" case, NATS covers it without adding Redis.

## vs MQTT brokers (Mosquitto, EMQX, HiveMQ)

**MQTT** is the dominant **IoT device protocol** (QoS 0/1/2, retained messages, last-will).

- **Pick a dedicated MQTT broker when:** you need full MQTT 5.0 features (QoS 2, shared subscriptions nuances, broker-specific IoT tooling) at large device fleets.
- **Pick NATS when:** you want IoT *and* your backend fabric in one system — **NATS speaks MQTT natively** (QoS 0/1) on `:1883`, so devices connect via MQTT while services use native NATS/JetStream on the same bus. This collapses "MQTT broker + backend bus + bridge" into one server. (This is exactly why OneMind dropped EMQX — NATS handles MQTT at its scale.)

## vs gRPC

**gRPC** is point-to-point RPC over HTTP/2 with protobuf — not a message bus.

- **Pick gRPC when:** you want strongly-typed, schema-first, point-to-point service calls with streaming and you know the callee's address.
- **Pick NATS when:** you want *location-transparent* request-reply (caller doesn't know/care where the service is or how many instances), pub/sub fan-out, load-balanced queue groups, and one fabric for many patterns. NATS request-reply gives you RPC ergonomics *plus* discovery, scaling, and resilience for free.

---

## The One-Paragraph Decision Rule

> If your problem is **"connect many moving services/devices across a network (including edge/geo) with low latency and minimal ops, using several messaging patterns"** → **NATS**. If your problem is **"durably store and reprocess enormous event volumes with a rich stream-processing ecosystem"** → **Kafka/Pulsar** (or NATS JetStream if the volume is moderate and you value simplicity). If your problem is **"complex broker-side routing / AMQP compatibility"** → **RabbitMQ**. NATS's superpower is *unifying* messaging + req-reply + KV + IoT + edge in one tiny self-hostable binary.

## Sources
- Compare NATS — https://docs.nats.io/nats-concepts/overview/compare-nats
- NATS Overview — https://docs.nats.io/nats-concepts/overview
- nats.io (adopters, comparisons) — https://nats.io
- JetStream — https://docs.nats.io/nats-concepts/jetstream
