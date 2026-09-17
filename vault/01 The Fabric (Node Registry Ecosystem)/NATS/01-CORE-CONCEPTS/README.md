---
title: "01 — Core Concepts"
section: "01-CORE-CONCEPTS"
type: index
tags: [nats, core-concepts, subjects, pubsub, request-reply, queue-groups, headers, connection]
source: https://docs.nats.io/nats-concepts/core-nats
created: 2026-06-14
owner: Legacy
status: active
---

# 01 — Core Concepts

> **What this section covers.** The foundational mechanics of NATS Core — how messages are addressed, delivered, and consumed. Everything here is stateless and in-memory (at-most-once). For persistence and guaranteed delivery, see [[NATS/02-JETSTREAM/README|02 — JetStream]].

---

## Section Map

| File | Topic | One-liner |
|------|-------|-----------|
| [[NATS/01-CORE-CONCEPTS/subjects-and-wildcards\|subjects-and-wildcards]] | Subjects & Wildcards | Subject naming rules, token hierarchy, `*` and `>` wildcards, best practices |
| [[NATS/01-CORE-CONCEPTS/publish-subscribe\|publish-subscribe]] | Publish / Subscribe | Core fan-out model, at-most-once semantics, async/sync subscriptions, auto-unsubscribe |
| [[NATS/01-CORE-CONCEPTS/request-reply\|request-reply]] | Request / Reply | RPC pattern, `_INBOX` subjects, no-responders (503), scatter-gather |
| [[NATS/01-CORE-CONCEPTS/queue-groups\|queue-groups]] | Queue Groups | Load-balanced competing consumers, geo-affinity, queue groups + JetStream |
| [[NATS/01-CORE-CONCEPTS/messages-and-headers\|messages-and-headers]] | Messages & Headers | Message anatomy, NATS/1.0 header format, max payload, no-responder headers |
| [[NATS/01-CORE-CONCEPTS/connection-lifecycle\|connection-lifecycle]] | Connection Lifecycle | Connect/reconnect, ping/pong, drain vs close, slow consumers, auth on connect |

---

## Mental Model

```
Publisher                NATS Server              Subscriber(s)
   │                        │                         │
   │  PUB time.us.east       │                         │
   │────────────────────────►│                         │
   │                        │  SUB time.*.*  ──────────┤ sub-A
   │                        │  SUB time.>    ──────────┤ sub-B
   │                        │  QueueSub workers ───────┤ sub-C  ┐  queue group:
   │                        │  QueueSub workers ───────┤ sub-D  ┘  only ONE fires
   │                        │                         │
```

- **Subjects** are the addressing primitives — dot-separated strings.
- **Pub/Sub** is 1→N: every matching subscriber receives every message.
- **Queue Groups** make it 1→1 within the group: exactly one subscriber fires per message.
- **Request/Reply** layers RPC on top: publisher sends with a reply-inbox subject; a responder replies to it.
- **Headers** carry metadata alongside the payload (NATS/1.0 HTTP-like format).
- **Connection lifecycle** governs how clients connect, survive failures, and shut down gracefully.

---

## Key Facts at a Glance

| Property | Value |
|----------|-------|
| Delivery guarantee (Core) | At-most-once (fire-and-forget) |
| Max message payload (default) | 1 MB (configurable up to 64 MB) |
| Wildcard: single token | `*` |
| Wildcard: multi-token tail | `>` |
| Reply inbox prefix | `_INBOX.` |
| Ping interval (default) | 2 minutes |
| Max pings outstanding (default) | 2 |
| Reconnect wait (default) | 2 seconds |
| Max reconnect attempts (default) | 60 |
| Reconnect buffer (default) | 8 MB |
| Subject token separator | `.` (dot) |
| Spaces in subjects | ❌ Not allowed |

---

## Relationships to Other Sections

- [[NATS/02-JETSTREAM/README|JetStream]] builds on Core: subjects are the same, but messages are persisted in Streams and consumed via durable Consumers.
- [[NATS/05-SECURITY/README|Security]] controls which subjects clients can publish/subscribe to via authorization rules.
- [[NATS/06-TOPOLOGY-SCALING/README|Topology]] describes how subject subscriptions propagate across clusters and superclusters.
- [[NATS/10-PROTOCOL-INTERNALS/README|Protocol Internals]] describes the wire format for `PUB`, `SUB`, `HPUB`, `HMSG`, etc.

---

## Sources

- https://docs.nats.io/nats-concepts/core-nats
- https://docs.nats.io/nats-concepts/subjects
- https://docs.nats.io/nats-concepts/core-nats/pubsub
- https://docs.nats.io/nats-concepts/core-nats/reqreply
- https://docs.nats.io/nats-concepts/core-nats/queue.md
- https://docs.nats.io/using-nats/developer/connecting.md
