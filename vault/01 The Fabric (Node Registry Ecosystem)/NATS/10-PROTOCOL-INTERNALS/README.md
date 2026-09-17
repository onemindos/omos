---
title: "Protocol & Internals — Section Index"
section: "10-PROTOCOL-INTERNALS"
type: index
tags: [nats, protocol, internals, wire-protocol, gossip, sublist, nuid, nkeys]
created: 2026-06-14
---

# 10 — Protocol & Internals

How NATS actually works under the hood: the **text-based wire protocol** every client and server speaks, the **cluster protocol** that meshes servers, the **server architecture** (sublist routing, gossip discovery, accounts), and the two crypto/identity primitives — **NUID** and **NKEYS**.

> This is the "open the hood" section. You don't need it to *use* NATS — but you need it to debug a stuck connection with `telnet`, reason about routing performance, or understand why NATS scales the way it does.

See also: [[NATS/01-CORE-CONCEPTS/README|Core Concepts]] · [[NATS/06-TOPOLOGY-SCALING/README|Topology & Scaling]] · [[NATS/05-SECURITY/jwt-nkeys-decentralized|JWT & NKEYS]]

## Notes
- [[NATS/10-PROTOCOL-INTERNALS/wire-protocol|Wire Protocol]] — the client↔server text protocol: `CONNECT`, `INFO`, `PUB`/`SUB`/`MSG`, `HPUB`/`HMSG`, `PING`/`PONG`, `+OK`/`-ERR`
- [[NATS/10-PROTOCOL-INTERNALS/cluster-protocol|Cluster Protocol]] — how servers route between each other (`RS+`/`RS-`/`RMSG`), full-mesh gossip, interest propagation
- [[NATS/10-PROTOCOL-INTERNALS/server-architecture|Server Architecture]] — the sublist (subject routing trie), connection types, account isolation, the event loop
- [[NATS/10-PROTOCOL-INTERNALS/nuid-and-nkeys|NUID & NKEYS]] — unique-ID generation (inboxes, CIDs) and Ed25519 identity crypto

## The protocol in one breath

NATS speaks a **human-readable, line-oriented text protocol** over TCP (or TLS/WebSocket). You can literally `telnet` a server and talk to it. There are only ~12 verbs. Control lines end in `\r\n`; message bodies are length-prefixed. This deliberate simplicity is *why* NATS has clients in 40+ languages — the protocol is trivial to implement.

```
CLIENT                          SERVER
  │  ──── TCP connect ────▶        │
  │  ◀──── INFO {...} ─────         │   server announces itself
  │  ──── CONNECT {...} ──▶         │   client declares options/auth
  │  ──── SUB foo 1 ──────▶         │   subscribe, sid=1
  │  ──── PUB foo 5\r\nhello ─▶     │   publish 5-byte payload
  │  ◀──── MSG foo 1 5\r\nhello ─   │   delivered to subscription 1
  │  ◀──── PING ──────────          │   keepalive
  │  ──── PONG ──────────▶          │
```

## Why the internals matter for OneMind
- **Debuggability** — when a bridge or worker won't connect, knowing the `INFO`/`CONNECT`/`-ERR` handshake lets us diagnose with `nats-box` or raw `openssl s_client`.
- **Performance reasoning** — the sublist trie explains why wildcard-heavy subject schemes are cheap, and why `>` at the root is expensive. Informs our subject law.
- **Identity** — every `_INBOX.<NUID>` reply subject and every agent NKEY identity in the fleet comes from the primitives documented here.

## Sources
- NATS protocol reference — https://docs.nats.io/reference/nats-protocol/nats-protocol
- Protocol demo — https://docs.nats.io/reference/nats-protocol/nats-protocol-demo
- nats-server source — https://github.com/nats-io/nats-server
