---
title: "Connection Lifecycle"
section: "01-CORE-CONCEPTS"
type: note
tags: [nats, connection, reconnect, ping-pong, drain, slow-consumer]
source: https://docs.nats.io/using-nats/developer/connecting
created: 2026-06-14
---

# Connection Lifecycle

See also: [[NATS/01-CORE-CONCEPTS/publish-subscribe|Publish/Subscribe]] · [[NATS/05-SECURITY/README|Security]] · [[NATS/10-PROTOCOL-INTERNALS/README|Protocol & Internals]]

---

## The Handshake

When a client connects, the server immediately sends an **`INFO`** protocol message describing itself (server id, version, max payload, whether TLS is required, cluster URLs for discovery). The client replies with **`CONNECT`** carrying its options (verbose, pedantic, TLS, auth credentials, name, protocol level). From then on the connection is a long-lived TCP (or TLS/WebSocket) stream over which `PUB`/`SUB`/`MSG`/`PING`/`PONG`/`HPUB`/`HMSG` flow.

```
Client ──TCP connect──▶ Server
Server ──INFO {...}────▶ Client      (server capabilities + cluster URLs)
Client ──CONNECT {...}─▶ Server      (client options + credentials)
Client ──PING─────────▶ Server
Server ──PONG─────────▶ Client       (connection live)
        ◀──── normal PUB/SUB/MSG traffic ────▶
```

> Because the server advertises **cluster URLs** in `INFO`, clients learn the full topology automatically and can fail over to other servers without being pre-configured with every address.

---

## Connecting (client examples)

```go
// Go
nc, _ := nats.Connect("nats://localhost:4222",
    nats.Name("my-service"),
    nats.MaxReconnects(-1),            // retry forever
    nats.ReconnectWait(2*time.Second),
    nats.UserCredentials("zeus.creds"))
defer nc.Drain()
```

```python
# Python (asyncio)
nc = await nats.connect(
    servers=["nats://localhost:4222"],
    name="my-service",
    max_reconnect_attempts=-1,
    reconnect_time_wait=2,
)
```

Common connect options across clients: server list, name, credentials (token / user-pass / NKEY / JWT creds file), TLS config, max reconnects, reconnect wait + jitter, ping interval, max pings out, no-randomize, drain timeout.

---

## Auto-Reconnect

NATS clients reconnect **transparently**. On connection loss the client:
1. Walks its server list (plus servers learned from `INFO`), optionally randomized.
2. Retries with `ReconnectWait` + jitter, up to `MaxReconnects` (`-1` = forever).
3. On success, **re-establishes all subscriptions automatically** and flushes buffered outbound messages.

During the disconnected window, published messages may be buffered (bounded by the client's reconnect buffer) — beyond that they error. Application code is notified via callbacks/events: `DisconnectedCB`, `ReconnectedCB`, `ClosedCB`, `ErrorCB`.

> **Core caveat:** auto-reconnect resubscribes, but NATS **Core** is at-most-once — messages published by *others* while you were disconnected are gone. If you need gap-free delivery across disconnects, use **JetStream** (a durable consumer resumes from its last ack'd sequence).

---

## PING/PONG Keepalive

The server and client exchange **PING/PONG** at a configurable interval (`ping_interval`, default ~2m) to detect half-open connections. If a client misses `max_pings_out` consecutive PONGs, the connection is considered stale and torn down (triggering reconnect). This is how NATS detects dead peers fast without relying solely on TCP timeouts.

---

## Drain vs Close

| Operation | Behavior |
|-----------|----------|
| **Close** | Immediately tears down the connection. In-flight messages and unprocessed subscription messages may be **lost**. |
| **Drain** | **Graceful.** Unsubscribes, then *processes all already-received messages* and flushes pending outbound before closing. The correct way to shut a service down without dropping work. |

```go
nc.Drain()   // graceful: finish in-flight, then close
// vs
nc.Close()   // abrupt
```

Always **drain** for clean service shutdown — especially queue-group workers, so in-flight requests complete.

---

## Slow Consumers

If a subscriber can't keep up, the server buffers pending messages per-subscription up to a **pending limit** (count and bytes). Exceeding it makes the client a **slow consumer** — the server drops messages for that subscription and raises a `SlowConsumer` error. Mitigations: increase pending limits, process faster, use **queue groups** to spread load, or move to **JetStream** pull consumers (flow-controlled, pull-based — the consumer asks for what it can handle, so it can't be overwhelmed the same way).

---

## Authentication on Connect

Credentials are presented in the `CONNECT` message (or via TLS client cert). Supported schemes (see [[NATS/05-SECURITY/README|Security]]):
- **Token** — single shared bearer token
- **User/Password**
- **NKEY** — ed25519 challenge/response (server sends a nonce, client signs it)
- **JWT** — decentralized; user JWT + NKEY seed in a `.creds` file
- **TLS client certificate**

A rejected auth closes the connection with an `-ERR 'Authorization Violation'`.

---

## Connection Events (observability)

Clients expose lifecycle callbacks you should wire to logs/metrics:

| Event | Meaning |
|-------|---------|
| `Connected` | Initial connection established |
| `Disconnected` | Lost connection (reconnect loop begins) |
| `Reconnected` | Re-established (subs restored) |
| `Closed` | Permanently closed (no more reconnects) |
| `Error` | Async error (e.g. slow consumer, permission violation) |
| `Discovered servers` | Cluster topology changed (new server learned) |

## Sources
- Connecting — https://docs.nats.io/using-nats/developer/connecting
- Reconnect — https://docs.nats.io/using-nats/developer/connecting/reconnect
- Drain — https://docs.nats.io/using-nats/developer/receiving/drain
- Slow consumers — https://docs.nats.io/using-nats/developer/events/slow
- Client protocol — https://docs.nats.io/reference/reference-protocols/nats-protocol
