---
title: NATS — Master Knowledge Base
domain: NATS
type: index
status: active
owner: Legacy
created: 2026-06-14
source_of_truth: docs.nats.io + Context7 (/nats-io/nats.docs) + nats-io & synadia-io GitHub
nats_server_version: v2.14.2 (2026-06-02)
tags: [nats, jetstream, synadia, messaging, infrastructure, knowledge-base]
---

# NATS — The Complete Knowledge Base

> **Purpose.** A canonical, exhaustive, expert-level reference for **NATS** — what it is, every feature it provides, how it works at every level, its full ecosystem, and **Synadia** (the company behind it). This folder is *universal NATS knowledge*. For **OneMind-specific NATS ops** (our live bus, subject law, bridges, workers), see the companion material under "Cross-References" below — this KB does not duplicate it.

> **Provenance.** Every document cites its sources (official `docs.nats.io`, the `nats-io/nats.docs` repo, Context7 indexed docs, and GitHub). Where a claim is version-specific, the NATS server version is noted. Built 2026-06-14 against **nats-server v2.14.2**.

---

## What is NATS (60-second version)

**NATS** is a high-performance, open-source messaging system — a single ~20MB Go binary, zero external dependencies (no JVM, Erlang, or Zookeeper) — built for microservices, cloud-native apps, edge, and IoT. It unifies multiple messaging patterns in one runtime:

- **Pub/Sub** — subject-based, one-to-many, at-most-once (NATS Core)
- **Request/Reply** — built-in RPC with automatic routing
- **Queue Groups** — load-balanced competing consumers
- **JetStream** — persistence layer: streaming, exactly-once, replay, work queues
- **Key-Value Store** — distributed KV with watch + history (built on JetStream)
- **Object Store** — large blob storage with chunking (built on JetStream)
- **Services / Micro** — service framework with discovery + observability

It is a **CNCF** project, originally written by **Derek Collison** (creator of the TIBCO messaging systems and Cloud Foundry's NATS), now stewarded by **Synadia**. The name historically stood for **Neural Autonomic Transport System**.

| Stat | Value |
|------|-------|
| First commit | 2012 (server repo created 2012-10-29) |
| License | Apache-2.0 |
| Language | Go |
| Binary footprint | ~20 MB, single static binary |
| Core throughput | up to ~15M+ msgs/sec (in-memory pub/sub) |
| Latency | sub-100µs typical end-to-end |
| CNCF status | Incubating → graduated track |
| Latest server | v2.14.2 (2026-06-02) |
| Maintainer | Synadia + open-source community |

---

## How to Use This KB

- **Linear learner?** Read sections `00 → 14` in order. Each builds on the prior.
- **Need one thing fast?** Jump via the map below; every section has its own `README`/index note.
- **Building on the bus?** Pair this with the OneMind doctrine (Cross-References) — this KB is the *theory*, the doctrine is *our law*.
- **Going expert?** The companion **skill group** (`nats/` category in `~/.hermes/skills/`) turns this knowledge into repeatable workflows.

---

## Section Map

| # | Section | What's inside |
|---|---------|---------------|
| **00** | [[NATS/00-OVERVIEW/README\|Overview]] | What NATS is, history, philosophy, comparison vs Kafka/RabbitMQ/Pulsar/Redis/MQTT, when to use it |
| **01** | [[NATS/01-CORE-CONCEPTS/README\|Core Concepts]] | Subjects, wildcards, pub/sub, request-reply, queue groups, headers, no-responders, connection lifecycle |
| **02** | [[NATS/02-JETSTREAM/README\|JetStream]] | Streams, consumers (push/pull), retention & discard policies, acks, dedup, replay, mirrors & sources, KV/OS foundation, clustering |
| **03** | [[NATS/03-KV-OBJECT-STORE/README\|KV & Object Store]] | Key-Value buckets, watchers, history, revisions, TTL, atomic ops; Object Store chunking, metadata, large blobs |
| **04** | [[NATS/04-SERVICES-MICRO/README\|Services / Micro]] | Micro framework, service definitions, endpoints, discovery (`$SRV`), stats, versioning |
| **05** | [[NATS/05-SECURITY/README\|Security]] | TLS, token/user-pass/NKEY/JWT auth, authorization, accounts, multi-tenancy, exports/imports, auth callout, decentralized auth |
| **06** | [[NATS/06-TOPOLOGY-SCALING/README\|Topology & Scaling]] | Single server → clusters → superclusters → gateways → leaf nodes; JetStream domains, mirrors, geo-distribution |
| **07** | [[NATS/07-DEPLOYMENT-OPS/README\|Deployment & Ops]] | Install paths, full config reference, monitoring endpoints, Kubernetes/Helm, NACK, tuning, backup/DR, upgrades |
| **08** | [[NATS/08-CLIENTS-SDKS/README\|Clients & SDKs]] | Go, Python, JS/TS, Rust, Java, C#, C, Ruby, Elixir, Deno — connection, patterns, JetStream API per language |
| **09** | [[NATS/09-CLI-TOOLING/README\|CLI & Tooling]] | `nats` CLI, `nats-top`, `nsc`, `nk`, `nats bench`, contexts, full command catalog |
| **10** | [[NATS/10-PROTOCOL-INTERNALS/README\|Protocol & Internals]] | Wire protocol (client + cluster), server architecture, gossip, sublist routing, NUID, NKEYS crypto |
| **11** | [[NATS/11-ADVANCED/README\|Advanced]] | MQTT support, WebSocket, embedded NATS (Go), leaf-node edge patterns, benchmarking, performance tuning |
| **12** | [[NATS/12-ECOSYSTEM/README\|Ecosystem]] | Full `nats-io` GitHub map (100 repos), surveyor, exporters, k8s operators, integrations, NATS by Example |
| **13** | [[NATS/13-SYNADIA/README\|Synadia (the company)]] | History, leadership, Synadia Cloud (NGS), Synadia Platform, Control Plane, Connect, NEX, Orbit, OSS-vs-commercial split, `synadia-io` GitHub map (40 repos) |
| **14** | [[NATS/14-REFERENCE/README\|Reference]] | Glossary, config cheatsheet, CLI cheatsheet, subject-design patterns, monitoring-endpoint catalog, limits table |

---

## The NATS Feature Tree (everything it provides)

```
NATS
├── NATS Core (at-most-once, in-memory)
│   ├── Publish / Subscribe        — subject-addressed, 1→N fan-out
│   ├── Request / Reply            — built-in RPC, _INBOX routing
│   ├── Queue Groups               — load-balanced competing consumers
│   ├── Wildcards                  — token (*) and full (>) subject matching
│   ├── Headers                    — NATS/1.0 message headers
│   └── No-Responders              — instant 503 when no service listening
│
├── JetStream (at-least-once / exactly-once, persistent)
│   ├── Streams                    — append-only logs over subjects
│   │   ├── Storage: file | memory
│   │   ├── Retention: limits | interest | workqueue
│   │   ├── Discard: old | new
│   │   ├── Dedup window           — Nats-Msg-Id sliding window
│   │   ├── Mirrors & Sources      — stream replication / aggregation
│   │   └── Replicas (R1/R3/R5)    — RAFT-replicated in cluster
│   ├── Consumers
│   │   ├── Push                   — server delivers to subject
│   │   ├── Pull                   — client fetches batches (preferred)
│   │   ├── Durable | Ephemeral
│   │   ├── Ack policies           — none | all | explicit
│   │   ├── Replay                 — instant | original
│   │   └── AckWait / MaxDeliver / Backoff
│   ├── Key-Value Store            — buckets, get/put/update, watch, history, TTL, revisions, atomic CAS
│   └── Object Store               — chunked blobs, metadata, watch, large files
│
├── Services / Micro Framework
│   ├── Service definitions        — name, version, description
│   ├── Endpoints + groups         — subject-mapped handlers
│   ├── Discovery ($SRV)           — INFO / PING / STATS
│   └── Observability              — auto request/error/latency stats
│
├── Connectivity Topologies
│   ├── Single server
│   ├── Cluster                    — full-mesh routes, horizontal scale
│   ├── Supercluster               — gateways connect clusters (geo)
│   ├── Leaf Nodes                 — edge extension, store-and-forward
│   └── JetStream Domains          — isolated JS per leaf/edge
│
├── Security
│   ├── TLS / mTLS
│   ├── AuthN: token | user-pass | NKEY | JWT (decentralized)
│   ├── AuthZ: subject-level allow/deny per user
│   ├── Accounts                   — isolated subject spaces (multi-tenant)
│   ├── Exports / Imports          — cross-account streams & services
│   └── Auth Callout               — external auth delegation
│
├── Protocols & Interop
│   ├── Native NATS wire protocol  — text-based, human-readable
│   ├── WebSocket                  — browser-native
│   ├── MQTT                       — QoS 0/1, IoT devices
│   └── Embedded (Go in-process)
│
└── Operations
    ├── Monitoring (:8222)         — /varz /connz /jsz /healthz /accountz ...
    ├── nats CLI / nats-top / nsc / nk
    ├── Prometheus exporter + surveyor
    └── Kubernetes (Helm + NACK)
```

---

## Cross-References (OneMind-specific — NOT duplicated here)

This KB is *universal NATS*. Our **as-built** material lives elsewhere — read it alongside:

- **Binding doctrine:** `~/onemindos-core/docs/NATS-DOCTRINE.md` — subject law, JetStream rules, agent-layer rules, Synadia product map for OneMind
- **Setup record:** `~/onemindos-core/docs/infra/NATS-SETUP.md`
- **Live stack state:** `SO/03 Infrastructure (The Body)/ONEMIND-OS-STACK.md`
- **Ops skill (40+ references):** `~/.hermes/skills/devops/nats-server/` — our launchd service, bridges, subscribers, workers, CLI cheatsheet, live stream topology
- **Expert skill group (this KB's companion):** `~/.hermes/skills/nats/` — workflow playbooks built on this knowledge

---

## Source Index

| Source | URL | Used for |
|--------|-----|----------|
| Official docs | https://docs.nats.io | Primary — all concepts |
| Docs repo | https://github.com/nats-io/nats.docs | Canonical markdown source |
| Context7 (docs) | `/nats-io/nats.docs` (294K tokens, 2,494 snippets) | Indexed concept + code retrieval |
| Context7 (website) | `/websites/nats_io` (289K tokens, 3,194 snippets) | Supplementary |
| Server repo | https://github.com/nats-io/nats-server | Internals, config, releases |
| nats.io | https://nats.io | Marketing, comparisons, adopters |
| NATS by Example | https://natsbyexample.com | Per-language runnable examples |
| GitHub org (nats-io) | https://github.com/nats-io | Ecosystem map (100 repos) |
| GitHub org (synadia-io) | https://github.com/synadia-io | Synadia products (40 repos) |
| Synadia | https://synadia.com | Company + commercial products |

*Last verified: 2026-06-14. NATS server v2.14.2.*
