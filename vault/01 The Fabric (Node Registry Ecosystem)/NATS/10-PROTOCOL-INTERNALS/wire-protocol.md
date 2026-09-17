---
title: "Wire Protocol (Client ↔ Server)"
section: "10-PROTOCOL-INTERNALS"
type: note
tags: [nats, protocol, wire-protocol, pub, sub, msg, hpub, hmsg, headers]
source: https://docs.nats.io/reference/nats-protocol/nats-protocol
created: 2026-06-14
---

# Wire Protocol (Client ↔ Server)

The NATS **client protocol** is a small, text-based, line-oriented protocol over TCP. Control lines terminate with `\r\n` (CRLF). It is human-readable — you can drive a server by hand with `telnet`/`nats-box`. Every client SDK in [[NATS/08-CLIENTS-SDKS/README|Section 08]] is, at bottom, an implementation of this protocol.

See also: [[NATS/10-PROTOCOL-INTERNALS/cluster-protocol|Cluster Protocol]] · [[NATS/01-CORE-CONCEPTS/messages-and-headers|Messages & Headers]] · [[NATS/01-CORE-CONCEPTS/connection-lifecycle|Connection Lifecycle]]

---

## The verbs (complete set)

| Verb | Sent by | Purpose |
|------|---------|---------|
| `INFO` | Server | Sent on connect (and on topology change) — server identity + capabilities |
| `CONNECT` | Client | Client options: auth, verbose, TLS, protocol level, headers support |
| `PUB` | Client | Publish a message to a subject (optional reply-to) |
| `HPUB` | Client | Publish **with headers** (NATS/1.0) |
| `SUB` | Client | Subscribe to a subject (optional queue group), assign `sid` |
| `UNSUB` | Client | Unsubscribe `sid` (optionally after N more messages) |
| `MSG` | Server | Deliver a message to a subscription |
| `HMSG` | Server | Deliver a message **with headers** |
| `PING` | Both | Keepalive / RTT probe |
| `PONG` | Both | Keepalive response |
| `+OK` | Server | Acknowledgement (only when `verbose:true`) |
| `-ERR` | Server | Protocol/auth error (some fatal — connection closed) |

That's the entire protocol. Everything NATS does — pub/sub, request-reply, queue groups, JetStream — is built from these primitives. (JetStream adds no new verbs; it rides on `PUB`/`SUB`/`MSG` against `$JS.API.*` subjects with JSON payloads.)

---

## The handshake

A client begins as a plain TCP socket. The server speaks first:

```
# 1. Server → client, immediately on accept
INFO {"server_id":"NDx…","server_name":"omos-mbp","version":"2.14.2","go":"go1.24",
      "host":"0.0.0.0","port":4222,"headers":true,"max_payload":1048576,
      "proto":1,"jetstream":true,"auth_required":true,"tls_required":false}\r\n

# 2. Client → server
CONNECT {"verbose":false,"pedantic":false,"tls_required":false,"name":"legacy-bridge",
         "lang":"go","version":"1.49.0","protocol":1,"headers":true,
         "user":"legacy","pass":"…"}\r\n

# 3. Client → server (optional immediate flush probe)
PING\r\n
# 3. Server → client
PONG\r\n
```

Key `INFO` fields: `max_payload` (hard cap on message size, default 1 MB), `headers` (whether `HPUB`/`HMSG` are allowed), `proto` (1 enables dynamic server discovery via subsequent `INFO`s), `auth_required`, `tls_required`, `jetstream`. The `CONNECT` is the client's side of the contract — it must respect `tls_required`/`auth_required` from `INFO`.

---

## Publish / Subscribe / Deliver

```
# PUB <subject> [reply-to] <#bytes>\r\n<payload>\r\n
PUB orders.new 5\r\nhello\r\n
PUB time.req _INBOX.aZ9k 0\r\n\r\n          # request: empty body, reply-to inbox

# SUB <subject> [queue-group] <sid>\r\n
SUB orders.* 1\r\n
SUB orders.new workers 7\r\n               # queue subscription

# MSG <subject> <sid> [reply-to] <#bytes>\r\n<payload>\r\n   (server → client)
MSG orders.new 1 5\r\nhello\r\n

# UNSUB <sid> [max_msgs]\r\n
UNSUB 1\r\n
UNSUB 7 10\r\n                              # auto-unsub after 10 more msgs
```

The byte count is the payload length **excluding** the trailing `\r\n`. This length-prefixing is what makes the protocol binary-safe despite being text-framed: payloads can contain any bytes, including CRLF.

---

## Headers — `HPUB` / `HMSG`

When both server and client set `headers:true`, messages carry a header block (the `NATS/1.0` version line + RFC-822-style key/values). `PUB`→`HPUB`, `MSG`→`HMSG`, with an extra header-length field:

```
# HPUB <subject> [reply-to] <#header-bytes> <#total-bytes>\r\n<headers>\r\n\r\n<payload>\r\n
HPUB orders.new 34 39\r\nNATS/1.0\r\nNats-Msg-Id: 42\r\n\r\nhello\r\n
```

Headers power: JetStream dedup (`Nats-Msg-Id`), message rollup, republish metadata, no-responder status (`503`), and JetStream per-message metadata. See [[NATS/01-CORE-CONCEPTS/messages-and-headers|Messages & Headers]].

---

## Keepalive & flow

- **`PING`/`PONG`** — either side may PING; the peer must PONG. The server PINGs on `ping_interval` (default 2 min); after `max_pings_out` (default 2) missed PONGs it declares a **stale connection** and closes. Clients also use a `PING`+`PONG` round-trip to implement `Flush()` (confirm the server has processed everything sent so far).
- **`+OK`** — only emitted when the client set `verbose:true` in `CONNECT`. Production clients use `verbose:false` (no per-op ack) for throughput.
- **`pedantic`** — when `true`, the server strictly validates subjects on every `PUB`/`SUB`. Off by default.

---

## Errors — `-ERR`

Two classes (full list in [[NATS/14-REFERENCE/limits-and-defaults|Reference: Limits & Defaults]]):

**Fatal (server closes connection):** `Unknown Protocol Operation`, `Authorization Violation`, `Authorization Timeout`, `Invalid Client Protocol`, `Maximum Control Line Exceeded`, `Secure Connection - TLS Required`, `Stale Connection`, `Maximum Connections Exceeded`.

**Non-fatal (connection stays open):** `Maximum Payload Violation` (the offending publish is rejected), `Permissions Violation for Subscription/Publish` (authz denied a specific subject), `Slow Consumer` (server-side pending buffer for the connection hit its limit — the *connection* may then be closed if it can't drain).

---

## Drive it by hand

```bash
# Raw protocol with nats-box (or telnet host 4222)
nats-box
> telnet localhost 4222
INFO {"server_id":...}        # server greets you
CONNECT {"verbose":true}      # you must send valid JSON
+OK
SUB foo 1
+OK
PUB foo 5
hello
+OK
MSG foo 1 5
hello
PING
PONG
```

This is the single best way to *understand* NATS — and to debug a misbehaving client.

---

## OneMind Relevance
- Every reply in the fleet uses a `_INBOX.<NUID>` reply-to in a `PUB`; understanding the request-reply framing here is what makes our service workers debuggable.
- `max_payload` (1 MB default) is why large blobs go through the **Object Store** (chunked), never a raw `PUB`. See [[NATS/03-KV-OBJECT-STORE/object-store|Object Store]].
- When a bridge logs `-ERR 'Authorization Violation'`, this doc is the decoder ring.

## Sources
- NATS protocol reference — https://docs.nats.io/reference/nats-protocol/nats-protocol
- Protocol demo (telnet walkthrough) — https://docs.nats.io/reference/nats-protocol/nats-protocol-demo
- Headers (ADR-4) — https://github.com/nats-io/nats-architecture-and-design
