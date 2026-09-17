---
title: "Bridges & Connectors"
section: "12-ECOSYSTEM"
type: note
tags: [nats, bridges, connectors, kafka, nginx, spring, integration, tracing]
source: https://github.com/nats-io/nats-kafka
created: 2026-06-14
---

# Bridges & Connectors

How NATS interoperates with the *outside* world — Kafka, web servers, JVM frameworks — when you can't be NATS-native end to end. The `nats-io` org keeps a small set of active bridges and a graveyard of archived ones (most superseded by **Synadia Connect**, see [[NATS/13-SYNADIA/README|Section 13]]).

See also: [[NATS/12-ECOSYSTEM/repo-map|Repo Map]] · [[NATS/00-OVERVIEW/nats-vs-alternatives|NATS vs Alternatives]] · [[NATS/13-SYNADIA/synadia-products|Synadia Products]]

---

## Active bridges

### `nats-kafka` — NATS ⇄ Kafka
Repo: `nats-io/nats-kafka` (★161, Go). A bidirectional bridge connecting NATS/JetStream subjects to Kafka topics. The pragmatic answer when an existing **Kafka** estate must interop with NATS — mirror a Kafka topic into a JetStream stream (or vice-versa) without rewriting producers/consumers.

```text
# conceptual connector config
connect: [
  { type: "NATSToKafka", subject: "events.>", topic: "events" },
  { type: "KafkaToNATS", topic: "orders",     subject: "orders.in" }
]
```

Use it for **migration** (drain Kafka into NATS over time) or **coexistence** (a partner system speaks Kafka). The doctrine remains NATS-first; bridges are the boundary, not the architecture.

### `nginx-nats` — NGINX module
Repo: `nats-io/nginx-nats` (★113, C). An NGINX module that acts as a NATS client, letting the web tier publish/consume on the bus from within NGINX. Niche, but powerful for edge HTTP→NATS ingestion without an app server.

### `spring-nats` — Spring Cloud Stream binder
Repo: `nats-io/spring-nats` (★82, Java). Binds NATS as a **Spring Cloud Stream** transport, so JVM/Spring apps use NATS through the familiar `@StreamListener`/binding abstractions. The on-ramp for enterprise Java shops.

---

## Tracing connectors

`not.go` (★104) and `not.java` (★2) — reference implementations for **distributed tracing** across NATS (propagate trace context in headers). Covered in [[NATS/12-ECOSYSTEM/monitoring-exporters|Monitoring & Exporters]].

---

## Archived bridges † (history / caution)

These exist but are **archived** — recognize them, don't adopt them:

| Repo † | Was | Today |
|--------|-----|-------|
| `nats-mq`† | NATS ⇄ IBM MQ Series | unmaintained |
| `nats-jms-bridge`† | NATS ⇄ JMS request/reply | unmaintained |
| `nats-connector-framework`† | Pluggable bridge framework (Java) | superseded |
| `nats-connector-redis`† | Redis pub/sub connector | superseded |
| `nats-replicator`† | Subject/channel replication | use JetStream **sources/mirrors** |
| `nats-account-server`† | HTTP JWT host for accounts | use built-in **NATS resolver** |
| `nats-rest-config-proxy`† | REST config proxy | superseded |

The modern, supported path for "connect NATS to system X" is increasingly **Synadia Connect** (managed connectors / Wombat runtime) rather than these one-off OSS bridges. See [[NATS/13-SYNADIA/synadia-products|Synadia Products]].

---

## Decision guide

| You need to… | Use |
|--------------|-----|
| Interop with existing **Kafka** | `nats-kafka` (active) |
| Replicate subjects/streams **NATS→NATS** | JetStream **sources & mirrors** (not a bridge!) — [[NATS/02-JETSTREAM/mirrors-and-sources\|see 02]] |
| Bridge to many SaaS/DBs (managed) | **Synadia Connect** — [[NATS/13-SYNADIA/synadia-products\|13]] |
| HTTP edge → NATS | `nginx-nats`, or a small service, or Synadia HTTP Gateway |
| JVM/Spring integration | `spring-nats` |

> **Key insight:** for NATS-to-NATS movement, you almost never want a "bridge" — **mirrors and sources** (native JetStream replication) are the right tool, and they're first-class. Bridges are for *non-NATS* systems.

---

## OneMind Relevance
- **NATS-native first**: intra-fleet data movement uses JetStream **sources/mirrors**, never a bridge. Bridges only appear at the boundary with a foreign system.
- If a third-party platform ever needs Kafka interop, `nats-kafka` is the sanctioned tool — but evaluate **Synadia Connect** first for anything resembling a managed connector.
- Avoid the archived connector zoo; their existence in search results is a trap.

## Sources
- nats-kafka — https://github.com/nats-io/nats-kafka
- nginx-nats — https://github.com/nats-io/nginx-nats
- spring-nats — https://github.com/nats-io/spring-nats
- Synadia Connect — https://docs.synadia.com/cloud/walkthrough/connect
