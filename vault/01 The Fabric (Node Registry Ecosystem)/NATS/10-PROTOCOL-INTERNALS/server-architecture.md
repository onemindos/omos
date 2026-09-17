---
title: "Server Architecture"
section: "10-PROTOCOL-INTERNALS"
type: note
tags: [nats, internals, sublist, routing, architecture, accounts, event-loop]
source: https://github.com/nats-io/nats-server
created: 2026-06-14
---

# Server Architecture

What's inside the single ~20 MB `nats-server` Go binary: the **sublist** (the subject-routing data structure that makes NATS fast), the **connection types** it multiplexes, **account** isolation, and the concurrency model.

See also: [[NATS/10-PROTOCOL-INTERNALS/wire-protocol|Wire Protocol]] · [[NATS/10-PROTOCOL-INTERNALS/cluster-protocol|Cluster Protocol]] · [[NATS/01-CORE-CONCEPTS/subjects-and-wildcards|Subjects & Wildcards]]

---

## The sublist — subject routing as a trie

The core hot-path structure is the **sublist**: a tokenized prefix tree (trie) mapping **subjects → interested subscriptions**. Subjects are split on `.` into tokens, and the tree has special branches for the two wildcards:

```
publish "orders.us.new"  →  tokens [orders, us, new]
                            walk the trie:
orders ─┬─ us ──┬─ new ──▶ {sub A, queue-group workers}
        │       └─ *  ───▶ {sub B}        (* = single-token wildcard)
        └─ >  ──────────▶ {sub C}         (> = full multi-token wildcard)
```

- **`*`** matches exactly one token at that position.
- **`>`** matches one-or-more remaining tokens (only valid as the last token).
- Lookups are **O(tokens)**, not O(subscriptions) — adding subscribers doesn't slow matching.
- The sublist has a **cache** of recently-resolved subjects → result sets, so repeated publishes to the same subject skip the walk entirely.

This is the engineering reason NATS sustains millions of msgs/sec and why **subject design matters**: deep, well-structured hierarchies are cheap; a subscription on bare `>` (everything) forces every publish to match it. See our subject law in the OneMind doctrine.

---

## One binary, many connection types

A running server multiplexes several connection kinds, each with its own port and protocol variant:

| Connection | Port (default) | Protocol |
|------------|---------------|----------|
| **Client** | 4222 | Client wire protocol ([[NATS/10-PROTOCOL-INTERNALS/wire-protocol\|wire-protocol]]) |
| **Route** (cluster peer) | 6222 | Route protocol (`RS+`/`RMSG`) |
| **Gateway** (supercluster) | 7222 | Gateway protocol (optimized interest) |
| **Leaf node** | 7422 | Leafnode protocol (single link) |
| **Monitoring** | 8222 | HTTP/JSON (`/varz`, `/connz`, `/jsz`…) |
| **WebSocket** | 443/8080 | WS upgrade → client protocol |
| **MQTT** | 1883 | MQTT 3.1.1 (mapped onto JetStream) |

All in one process. The server is an **event-driven, goroutine-per-connection** Go program: each connection has read/write goroutines; publishing into the sublist hands messages to interested subscriptions' write paths. There is no external broker state — NATS Core holds **no messages at rest** (interest-only). Persistence is the JetStream subsystem.

---

## Accounts — the isolation boundary

An **account** is an isolated subject space. Each account has its **own sublist**: subjects in account A are invisible to account B unless explicitly **exported/imported**. This is how a single server (or cluster) is safely multi-tenant. The `$SYS` account is special — it carries server events, `$SYS.>` monitoring, and is where tools like surveyor listen. See [[NATS/05-SECURITY/accounts-multitenancy|Accounts]].

```
server
├── account ACME      → sublist_ACME   (orders.>, …)
├── account WIDGETS   → sublist_WIDGETS (parts.>, …)
└── account $SYS      → system events, $SYS.>
```

Routing operations (`RS+`/`RMSG`) and JetStream assets are **always account-scoped** — the account is a first-class part of every internal address.

---

## JetStream subsystem

When `jetstream{}` is enabled, the server runs an additional subsystem that:
- owns a **storage directory** (file or memory) for streams,
- exposes the **`$JS.API.>`** management surface (create stream/consumer, etc. — JSON request/reply, no new protocol verbs),
- runs **RAFT groups** (`$NRG.>` / `$JS` system subjects) for R3/R5 replication and consumer state,
- enforces stream limits, retention, dedup windows, and ack tracking.

JetStream is layered *on top of* Core — it's subjects and messages all the way down, just with a durable, replicated state machine behind the management subjects. See [[NATS/02-JETSTREAM/jetstream-overview|JetStream Overview]].

---

## Memory & flow control

- **Slow consumers** — each connection has a pending write buffer (`write_deadline`, max pending). A client that can't keep up becomes a "slow consumer"; rather than let it back-pressure the whole server, NATS drops it (Core is at-most-once). JetStream consumers get durable, flow-controlled delivery instead.
- **`max_payload`** — single message size cap (default 1 MB), advertised in `INFO`.
- **`max_connections`**, **`max_subscriptions`**, **`max_control_line`** — per-server guard rails (see [[NATS/14-REFERENCE/limits-and-defaults|Reference: Limits]]).

---

## OneMind Relevance
- The sublist + cache model is the technical justification for our **structured subject law** — hierarchical, specific subjects keep the fabric fast at fleet scale.
- One binary doing client+route+leaf+monitoring+MQTT+WS is why a single `nats-server` (or embedded `omd`) can be a barn's entire messaging backbone with zero extra infra.
- Account-per-sublist isolation is the mechanism behind our FABRIC/INTERFACE/etc. account separation.

## Sources
- nats-server source — https://github.com/nats-io/nats-server
- Sublist (history) — https://github.com/nats-io/sublist
- Accounts & multi-tenancy — https://docs.nats.io/running-a-nats-service/configuration/securing_nats/accounts
- Monitoring internals — https://docs.nats.io/running-a-nats-service/nats_admin/monitoring
