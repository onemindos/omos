---
title: "Clients & SDKs — Section Index"
section: "08-CLIENTS-SDKS"
type: index
tags: [nats, clients, sdks, go, python, javascript, rust, java]
created: 2026-06-14
---

# 08 — Clients & SDKs

NATS has **first-party** clients in every major language, all speaking the same wire protocol with consistent concepts (connect, pub/sub, request-reply, JetStream, KV, OS, Micro). Pick your language; the mental model is identical.

## Notes
- [[NATS/08-CLIENTS-SDKS/go|Go (nats.go)]] — the reference client; embedded server lives here too
- [[NATS/08-CLIENTS-SDKS/python|Python (nats-py)]] — asyncio; OneMind's primary agent language
- [[NATS/08-CLIENTS-SDKS/javascript|JavaScript/TypeScript (nats.js)]] — Node, Deno, **browser** (WebSocket)
- [[NATS/08-CLIENTS-SDKS/rust|Rust (async-nats)]] — high-performance async
- [[NATS/08-CLIENTS-SDKS/other-languages|Java, C#, C, Ruby, Elixir, …]] — the rest of the official family

## The universal model
> Every client: **connect** (with creds/TLS/reconnect opts) → **publish/subscribe** on subjects → **request/reply** → get a **JetStream context** for streams/consumers/KV/OS → optionally build a **Micro** service. Method names differ slightly; concepts don't.

## Official clients (maturity)
| Language | Repo | Notes |
|----------|------|-------|
| **Go** | `nats-io/nats.go` | Reference impl; embeds the server |
| **Python** | `nats-io/nats.py` | asyncio (`nats-py`) |
| **JS/TS** | `nats-io/nats.js` | Node + Deno + browser (WS) |
| **Rust** | `nats-io/nats.rs` | `async-nats` |
| **Java** | `nats-io/nats.java` | JVM |
| **C#/.NET** | `nats-io/nats.net` | modern .NET |
| **C** | `nats-io/nats.c` | native, embeddable |
| **Ruby** | `nats-io/nats-pure.rb` | pure Ruby |
| **Elixir / others** | community + official | Zig, Rust, etc. |

## See also
- [[NATS/01-CORE-CONCEPTS/connection-lifecycle|Connection Lifecycle]] (applies to all clients)
- [[NATS/10-PROTOCOL-INTERNALS/wire-protocol|Wire Protocol]] (what every client speaks)
- **NATS by Example** — https://natsbyexample.com — runnable snippets in every language
